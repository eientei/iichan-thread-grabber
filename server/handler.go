package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eientei/iichan-thread-grabber/grabber"
	"github.com/eientei/iichan-thread-grabber/model"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var ErrInternal = errors.New("internal error")

type SubmitThreadRequest struct {
	ThreadUrl string `json:"thread_url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, PublicPrefix)
	if strings.Index(path, "..") > -1 {
		log.Println("dots")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch {
	case path == "/favicon.ico":
		w.WriteHeader(http.StatusBadRequest)
	case path == "/":
		switch r.Method {
		case http.MethodPost:
			handleSubmitThreadRequest(w, r)
		case http.MethodGet:
			handleListThreads(w, r)
		default:
			log.Println("invalid method /")
			w.WriteHeader(http.StatusBadRequest)
		}
	case path == "/submit":
		path = strings.TrimPrefix(path, "/submit")
		switch r.Method {
		case http.MethodPost:
			handleSubmitThread(w, r)
		default:
			log.Println("invalid method /submit")
			w.WriteHeader(http.StatusBadRequest)
		}
	case strings.HasPrefix(path, "/data"):
		path = strings.TrimPrefix(path, "/data")
		switch r.Method {
		case http.MethodGet:
			handleStatic(w, r, path)
		default:
			log.Println("invalid method /data")
			w.WriteHeader(http.StatusBadRequest)
		}
	default:
		switch r.Method {
		case http.MethodPost:
			handleUpdateThread(w, r, path)
		case http.MethodGet:
			handleGetThread(w, r, path)
		default:
			log.Println("invalid method default")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func handleStatic(w http.ResponseWriter, r *http.Request, path string) {
	path = OutputDir + path
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	parts := strings.Split(path, ".")
	w.Header().Set("Content-Type", mime.TypeByExtension("."+parts[len(parts)-1]))
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file)
}

func handleSubmitThread(w http.ResponseWriter, r *http.Request) {
}

type Submit struct {
	Url string `json:"url"`
}

type Image struct {
	Md5       string `json:"md5"`
	FilePath  string `json:"file_path"`
	ThumbPath string `json:"thumb_path"`
	PostId    int    `json:"post_id"`
	Tags      string `json:"tags"`
	Rating    string `json:"rating"`
	ParentMd5 string `json:"parent_md5"`
	Index     int    `json:"index"`
}

type Thread struct {
	Url         string     `json:"url"`
	LastChecked string     `json:"last_checked"`
	Groups      [][]*Image `json:"groups"`
}

func handleUpdateThread(w http.ResponseWriter, r *http.Request, path string) {
	m, err := model.GetThread(path)
	if err != nil {
		log.Println("error getting thread", path)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	t := &Thread{}
	err = json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	images := make(map[string]*model.ImageModel)
	for _, img := range m.Images {
		images[img.Md5] = img
	}

	fmt.Println(path)
	for gid, group := range t.Groups {
		for oid, img := range group {
			if mimg, ok := images[img.Md5]; ok {
				mimg.Group = gid
				if gid == 0 {
					mimg.Order = mimg.Index
				} else {
					mimg.Order = oid
				}
				mimg.ParentMd5 = img.ParentMd5
				mimg.Tags = img.Tags
				switch img.Rating {
				case "q", "e", "s":
					mimg.Rating = img.Rating
				default:
					mimg.Rating = "q"
				}
				err := mimg.Update(path)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte(err.Error()))
					return
				}
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

func handleGetThread(w http.ResponseWriter, r *http.Request, path string) {
	m, err := model.GetThread(path)
	if err != nil {
		log.Println("error getting thread", path)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	t := &Thread{
		Url:         "https://iichan.hk" + m.Url,
		LastChecked: m.Created.Format(time.RFC1123),
	}

	if len(m.Images) > 0 {
		gid := m.Images[0].Group
		var group []*Image
		for _, img := range m.Images {
			if img.Group != gid {
				t.Groups = append(t.Groups, group)
				group = nil
				gid = img.Group
			}
			group = append(group, &Image{
				Md5:       img.Md5,
				FilePath:  PublicBase + "/data" + strings.TrimPrefix(img.Filepath, OutputDir),
				ThumbPath: PublicBase + "/data" + strings.TrimPrefix(img.Thumbpath, OutputDir),
				PostId:    img.Postid,
				Tags:      img.Tags,
				Rating:    img.Rating,
				ParentMd5: img.ParentMd5,
				Index:     img.Index,
			})
		}
		t.Groups = append(t.Groups, group)
	}

	w.WriteHeader(http.StatusOK)
	bs, _ := json.Marshal(&grabber.StatusMessage{Status: "result", Data: t})
	_, _ = w.Write(bs)
}

func handleListThreads(w http.ResponseWriter, r *http.Request) {
	threads, err := model.ListThreads()
	for i, thread := range threads {
		threads[i] = PublicPrefix + thread
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	bs, _ := json.Marshal(threads)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bs)
}

func flushedSubmitThreadRequest(w http.ResponseWriter, r *http.Request, rawthreadurl string) (http.Flusher, error) {
	threadurl, err := url.Parse(rawthreadurl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, ErrInternal
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")

	done := make(chan struct{})
	var proc *ProcessingThread
	var sub *ProcessingSubscriber
	closing := int32(0)
	sub = &ProcessingSubscriber{
		Submit: func(status *grabber.StatusMessage) {
			defer func() {
				if recover() != nil {
					if atomic.CompareAndSwapInt32(&closing, 0, 1) {
						close(done)
					}
				}
			}()
			bs, err := json.Marshal(status)
			if err != nil {
				proc.Unsubscribe(sub)
				bs, _ = json.Marshal(&grabber.StatusMessage{
					Status: "error",
					Data: &grabber.ErrorStatus{
						Error: err.Error(),
					},
				})
			}
			bs = append(bs, '\n')
			_, werr := w.Write(bs)
			flusher.Flush()
			if err != nil || werr != nil {
				if atomic.CompareAndSwapInt32(&closing, 0, 1) {
					close(done)
				}
				return
			}
			switch status.Status {
			case "error", "complete":
				if atomic.CompareAndSwapInt32(&closing, 0, 1) {
					close(done)
				}
			}
		},
	}
	proc = SubmitThread(context.Background(), threadurl, sub)
	<-done
	return flusher, nil
}

func handleSubmitThreadRequest(w http.ResponseWriter, r *http.Request) {
	var rawthreadurl string

	if r.Header.Get("Content-Type") == "multipart/form-data" {
		err := r.ParseMultipartForm(int64(MultipartLimit))
		if err != nil || len(r.MultipartForm.Value["thread_url"]) == 0 || len(r.MultipartForm.Value["thread_url"][0]) == 0 {
			log.Println("decode error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rawthreadurl = r.MultipartForm.Value["thread_url"][0]
	} else {
		err := r.ParseForm()
		if err != nil || len(r.FormValue("thread_url")) == 0 {
			log.Println("decode error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rawthreadurl = r.FormValue("thread_url")
	}

	u, err := url.Parse(rawthreadurl)
	if err != nil || u.Scheme == "" || u.Host == "" || u.Path == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = flushedSubmitThreadRequest(w, r, "https://iichan.hk"+u.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}
