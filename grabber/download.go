package grabber

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

type DownloadingResult struct {
	Error error
	File  string
}

type DownloadingOrder struct {
	URL             *url.URL
	OutputDir       string
	ProcessingOrder *ThreadGrabOrder
	Result          chan *DownloadingResult
	Complete        int32
}

func (order *DownloadingOrder) Execute(ctx context.Context) (err error) {
	request, err := http.NewRequest(http.MethodGet, order.URL.String(), nil)
	if err != nil {
		return err
	}
	boardpath := "/"
	ps := strings.Split(order.URL.Path, "/")
	if len(ps) > 0 {
		boardpath = "/" + ps[0]
	}
	request.Header.Set("User-Agent", UserAgent)
	request.Header.Set("Referer", order.URL.Scheme+"://"+order.URL.Host+boardpath)
	request.AddCookie(&http.Cookie{
		Name:    "wakabastyle",
		Value:   "Futaba",
		Domain:  "iichan.hk",
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 24 * 300),
	})

	request = request.WithContext(ctx)

	log.Println("downloading", request.URL.String())
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
		}
	}()
	if resp.StatusCode/100 != 2 {
		return ErrStatus
	}

	err = os.MkdirAll(order.OutputDir, 0755)
	if err != nil {
		return err
	}
	name := path.Base(order.URL.Path)
	fileloc := filepath.Join(order.OutputDir, name)
	file, err := os.OpenFile(fileloc, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()
	n, err := io.Copy(file, resp.Body)
	if err != nil {
		return
	}
	if resp.ContentLength > 0 && n != resp.ContentLength {
		if atomic.CompareAndSwapInt32(&order.Complete, 0, 1) {
			order.Result <- &DownloadingResult{
				Error: ErrStatus,
			}
			close(order.Result)
		}
		return os.Remove(fileloc)
	}
	if atomic.CompareAndSwapInt32(&order.Complete, 0, 1) {
		order.Result <- &DownloadingResult{
			File: fileloc,
		}
		close(order.Result)
	}
	return
}

func (order *DownloadingOrder) Cancel(err error) {
	if atomic.CompareAndSwapInt32(&order.Complete, 0, 1) {
		order.Result <- &DownloadingResult{
			Error: err,
		}
		close(order.Result)
	}
}

func (order *DownloadingOrder) Notify(position int) {

}
