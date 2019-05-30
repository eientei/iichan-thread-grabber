package grabber

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/eientei/iichan-thread-grabber/model"
	"github.com/nfnt/resize"
	"golang.org/x/net/html"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type ThreadGrabOrder struct {
	Grabber     *ThreadGrabber
	Thread      *url.URL
	Complete    int32
	OutputDir   string
	PublicBase  string
	OutputChunk chan *StatusMessage
}

type ImageThumb struct {
	Full  string
	Thumb string
	Error error
}

func matchTagAttr(tokenizer *html.Tokenizer, expectedName, expectedKey string) (string, bool) {
	name, hasAttrs := tokenizer.TagName()
	if string(name) == expectedName && hasAttrs {
		for {
			ak, av, amore := tokenizer.TagAttr()
			if string(ak) == expectedKey {
				return string(av), true
			}
			if !amore {
				break
			}
		}
	}
	return "", false
}

func ParseThread(prefix, boardpath, id string, r io.Reader) ([]*ImageThumb, error) {
	tokenizer := html.NewTokenizer(r)
	depth := 0

	var images []*ImageThumb

	var full string

loop:
	for {
		tok := tokenizer.Next()
		switch tok {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				break loop
			}
			return nil, tokenizer.Err()
		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			switch depth {
			case 0:
				if value, ok := matchTagAttr(tokenizer, "div", "id"); ok {
					if strings.HasPrefix(value, "thread-"+id) {
						switch tok {
						case html.StartTagToken:
							depth++
						}
					}
				}
			case 1:
				if value, ok := matchTagAttr(tokenizer, "a", "href"); ok {
					if strings.HasPrefix(value, "/"+boardpath+"/src/") {
						switch tok {
						case html.StartTagToken:
							full = value
							depth++
						}
					}
				}
			case 2:
				if value, ok := matchTagAttr(tokenizer, "img", "src"); ok {
					if strings.HasPrefix(value, "/"+boardpath+"/thumb/") {
						images = append(images, &ImageThumb{
							Full:  prefix + full,
							Thumb: prefix + value,
						})
						depth--
					}
				}
			}
		}
	}

	return images, nil
}

type ImageDownloadOrder struct {
	ImageThumb       *ImageThumb
	DownloadingOrder *DownloadingOrder
}

func FileStat(path string) os.FileInfo {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil
	}
	return stat
}

func ThumbPath(outdir, image string) string {
	parts := strings.Split(path.Base(image), ".")
	if len(parts) < 2 {
		return ""
	}
	return filepath.Join(outdir, parts[0]+"s."+parts[1])
}

func (order *ThreadGrabOrder) Execute(ctx context.Context) (err error) {
	prefix := order.Thread.Scheme + "://" + order.Thread.Host
	tpath := order.Thread.Path
	if tpath[0] == '/' {
		tpath = tpath[1:]
	}
	tfparts := strings.Split(tpath, "/")
	if len(tfparts) < 2 {
		return ErrThreadUrl
	}
	boardname := tfparts[0]
	tfname := path.Base(order.Thread.Path)
	split := strings.Split(tfname, ".")
	if len(split) != 2 {
		return ErrThreadUrl
	}
	threadid := split[0]

	outdir := filepath.Join(order.OutputDir, boardname, threadid)

	outthreadfile := filepath.Join(outdir, path.Base(order.Thread.Path))

	stat := FileStat(outthreadfile)

	update := false
	if stat == nil || time.Now().Sub(stat.ModTime()) > ThreadCache {
		update = true
		resultchan := make(chan *DownloadingResult)

		order.Grabber.DownloadingQueue.Submit(ctx, &DownloadingOrder{
			URL:             order.Thread,
			OutputDir:       outdir,
			ProcessingOrder: order,
			Result:          resultchan,
		})

		res := <-resultchan

		if res.Error != nil {
			return res.Error
		}

		outthreadfile = res.File
	}

	file, err := os.OpenFile(outthreadfile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()
	images, err := ParseThread(prefix, boardname, threadid, file)
	if err != nil {
		return err
	}

	var fetches []*ImageDownloadOrder

	for _, img := range images {
		outimgfile := filepath.Join(outdir, path.Base(img.Full))
		stat := FileStat(outimgfile)

		if stat == nil || time.Now().Sub(stat.ModTime()) > ImageCache {
			u, err := url.Parse(img.Full)
			if err != nil {
				return err
			}

			thumbimgfile := ThumbPath(outdir, img.Full)
			if FileStat(thumbimgfile) != nil {
				_ = os.Remove(thumbimgfile)
			}

			fetches = append(fetches, &ImageDownloadOrder{
				ImageThumb: img,
				DownloadingOrder: &DownloadingOrder{
					URL:             u,
					OutputDir:       outdir,
					ProcessingOrder: order,
					Result:          make(chan *DownloadingResult),
				},
			})
		}
	}

	if len(fetches) > 0 {
		fetchcomplete := make(chan int)
		total := len(fetches)
		for _, f := range fetches {
			go func(ido *ImageDownloadOrder) {
				order.Grabber.DownloadingQueue.Submit(ctx, ido.DownloadingOrder)
				res := <-ido.DownloadingOrder.Result
				ido.ImageThumb.Error = res.Error
				fetchcomplete <- 1
			}(f)
		}
		for range fetchcomplete {
			total -= 1
			if atomic.CompareAndSwapInt32(&order.Complete, 0, 0) {
				order.OutputChunk <- &StatusMessage{
					Status: "download_progress",
					Data: &DownloadStatus{
						TotalDownload:   len(fetches),
						CurrentDownload: len(fetches) - total,
					},
				}
			}
			if total == 0 {
				close(fetchcomplete)
				break
			}
		}
	}

	_, err = model.GetThread(order.Thread.Path)
	if err != nil || update || len(fetches) > 0 {
		err = model.CreateThread(order.Thread.Path)
		if err != nil {
			return err
		}
	}

	now := time.Now()
	for _, img := range images {
		outimgfile := filepath.Join(outdir, path.Base(img.Full))
		thumbimgfile := ThumbPath(outdir, img.Full)
		thumbstat := FileStat(thumbimgfile)
		if thumbstat == nil || now.Sub(thumbstat.ModTime()) > ImageCache {
			infile, err := os.OpenFile(outimgfile, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
			imgd, fname, err := image.Decode(infile)
			_ = infile.Close()
			if err != nil {
				return err
			}
			imgd = resize.Thumbnail(200, 200, imgd, resize.MitchellNetravali)

			thumbfile, err := os.OpenFile(thumbimgfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}

			switch fname {
			case "jpeg":
				err = jpeg.Encode(thumbfile, imgd, nil)
			case "png":
				err = png.Encode(thumbfile, imgd)
			case "gif":
				err = gif.Encode(thumbfile, imgd, nil)
			}
			_ = thumbfile.Close()
			if err != nil {
				return err
			}
		}
		err = os.Chtimes(outimgfile, now, now)
		if err != nil {
			return err
		}
		err = os.Chtimes(thumbimgfile, now, now)
		if err != nil {
			return err
		}

		infile, err := os.OpenFile(outimgfile, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		buf := make([]byte, 4096)
		hash := md5.New()
		for {
			n, err := infile.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				break
			}
			n, err = hash.Write(buf[:n])
			if err != nil {
				return err
			}
		}
		_ = infile.Close()
		md5hash := hash.Sum(nil)

		err = model.CreateImage(order.Thread.Path, strings.ToLower(hex.EncodeToString(md5hash)), strings.TrimPrefix(outimgfile, order.OutputDir), strings.TrimPrefix(thumbimgfile, order.OutputDir))
		if err != nil {
			return err
		}
	}

	if atomic.CompareAndSwapInt32(&order.Complete, 0, 1) {
		order.OutputChunk <- &StatusMessage{
			Status: "complete",
			Data: &ReadyStatus{
				Base: order.PublicBase + order.Thread.Path,
			},
		}
		close(order.OutputChunk)
	}

	return nil
}

func (order *ThreadGrabOrder) Cancel(err error) {
	if atomic.CompareAndSwapInt32(&order.Complete, 0, 1) {
		order.OutputChunk <- &StatusMessage{
			Status: "error",
			Data: &ErrorStatus{
				Error: err.Error(),
			},
		}
		close(order.OutputChunk)
	}
}

func (order *ThreadGrabOrder) Notify(position int) {
	if atomic.CompareAndSwapInt32(&order.Complete, 0, 0) {
		order.OutputChunk <- &StatusMessage{
			Status: "queue_progress",
			Data: &QueueStatus{
				Position: position,
			},
		}
	}
}
