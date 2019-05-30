package server

import (
	"context"
	"github.com/eientei/iichan-thread-grabber/grabber"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var Grabber = grabber.NewThreadGrabber(
	grabber.NewExecutableQueue(DownloaderWorkers, DownloaderBuffer, DownloaderTimeout, DownloaderNotifyRate, DownloaderExecuteRate),
	grabber.NewExecutableQueue(GrabberWorkers, GrabberBuffer, GrabberTimeout, GrabberNotifyRate, GrabberExecuteRate),
)

var Processing = make(map[string]*ProcessingThread)

var ProcessingMutex = &sync.Mutex{}

func init() {
	if CleanupOutputDir > 0 {
		go func() {
			err := os.MkdirAll(OutputDir, 0755)
			if err != nil {
				panic(err)
			}
			for {
				log.Println("cleanup", OutputDir)
				var dirs []string
				now := time.Now()
				err := filepath.Walk(OutputDir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						dirs = append(dirs, path)
						return nil
					}
					if strings.HasSuffix(path, "html") && now.Sub(info.ModTime()) > grabber.ThreadCache*time.Duration(CleanupFactor) {
						return os.Remove(path)
					} else if now.Sub(info.ModTime()) > grabber.ImageCache*time.Duration(CleanupFactor) {
						return os.Remove(path)
					}
					return nil
				})
				if err != nil {
					log.Println(err)
					continue
				}
				for i := len(dirs) - 1; i >= 0; i-- {
					dir, err := os.Open(dirs[i])
					if err != nil {
						log.Println(err)
						continue
					}
					_, err = dir.Readdirnames(1)
					if err == io.EOF && dirs[i] != OutputDir {
						_ = os.Remove(dirs[i])
					}
				}

				time.Sleep(CleanupOutputDir)
			}
		}()
	}
}

type ProcessingSubscriber struct {
	Submit func(status *grabber.StatusMessage)
}

type ProcessingThread struct {
	Order       *grabber.ThreadGrabOrder
	Subscribers []*ProcessingSubscriber
	Mutex       *sync.Mutex
	LastChunk   *grabber.StatusMessage
}

func SubmitThread(ctx context.Context, threadurl *url.URL, subscriber *ProcessingSubscriber) *ProcessingThread {
	ProcessingMutex.Lock()
	defer ProcessingMutex.Unlock()
	if proc, ok := Processing[threadurl.Path]; ok {
		proc.Subscribe(subscriber)
		return proc
	}

	statuschan := make(chan *grabber.StatusMessage)
	proc := &ProcessingThread{
		Order: &grabber.ThreadGrabOrder{
			Grabber:     Grabber,
			Thread:      threadurl,
			OutputDir:   OutputDir,
			PublicBase:  PublicBase,
			OutputChunk: statuschan,
		},
		Mutex: &sync.Mutex{},
	}
	Processing[threadurl.Path] = proc
	proc.Subscribe(subscriber)
	go proc.Deliver()
	Grabber.ProcessingQueue.Submit(ctx, proc.Order)
	return proc
}

func (proc *ProcessingThread) Subscribe(subscriber *ProcessingSubscriber) {
	proc.Mutex.Lock()
	defer proc.Mutex.Unlock()
	if proc.LastChunk != nil {
		subscriber.Submit(proc.LastChunk)
	}
	proc.Subscribers = append(proc.Subscribers, subscriber)
}

func (proc *ProcessingThread) Unsubscribe(subscriber *ProcessingSubscriber) {
	proc.Mutex.Lock()
	defer proc.Mutex.Unlock()
	for i, sub := range proc.Subscribers {
		if sub == subscriber {
			proc.Subscribers = append(proc.Subscribers[:i], proc.Subscribers[i+1:]...)
			break
		}
	}
}

func (proc *ProcessingThread) Deliver() {
	for status := range proc.Order.OutputChunk {
		proc.Mutex.Lock()
		proc.LastChunk = status
		for _, sub := range proc.Subscribers {
			go sub.Submit(status)
		}
		proc.Mutex.Unlock()
	}
	ProcessingMutex.Lock()
	defer ProcessingMutex.Unlock()
	delete(Processing, proc.Order.Thread.Path)
	proc.LastChunk = nil
}
