package grabber

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestParseThread(t *testing.T) {
	grabber := NewThreadGrabber(NewExecutableQueue(1, 2048, time.Second*30, time.Second, time.Second*2), NewExecutableQueue(4, 16, time.Hour, time.Second, 0))
	url, err := url.Parse("https://iichan.hk/b/res/4899700.html")
	if err != nil {
		t.Fatal(err)
	}

	chunks := make(chan *StatusMessage)
	grabber.ProcessingQueue.Submit(context.Background(), &ThreadGrabOrder{
		Grabber:     grabber,
		Thread:      url,
		OutputDir:   "/tmp/iich",
		PublicBase:  "/iich",
		OutputChunk: chunks,
	})

	for raw := range chunks {
		bs, err := json.Marshal(raw)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(bs))
	}
}
