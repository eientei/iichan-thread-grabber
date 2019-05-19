package server

import (
	"encoding/json"
	"github.com/eientei/iichan-thread-grabber/grabber"
	"net/http"
	"net/url"
)

type SubmitThreadRequest struct {
	ThreadUrl string `json:"thread_url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	req := &SubmitThreadRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || len(req.ThreadUrl) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	threadurl, err := url.Parse(req.ThreadUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")

	done := make(chan struct{})
	var proc *ProcessingThread
	var sub *ProcessingSubscriber
	sub = &ProcessingSubscriber{
		Submit: func(status *grabber.StatusMessage) {
			defer func() {
				if recover() != nil {
					close(done)
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
				close(done)
				return
			}
			switch status.Status {
			case "error", "complete":
				close(done)
			}
		},
	}
	proc = SubmitThread(r.Context(), threadurl, sub)
	<-done
}
