package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eientei/iichan-thread-grabber/grabber"
	"net/url"
	"testing"
	"time"
)

func TestSubmitThread(t *testing.T) {
	tu, _ := url.Parse("https://iichan.hk/b/res/4897246.html")
	SubmitThread(context.Background(), tu, &ProcessingSubscriber{
		Submit: func(status *grabber.StatusMessage) {
			bs, _ := json.Marshal(status)
			fmt.Println(string(bs))
		},
	})
	time.Sleep(100 * time.Second)
}
