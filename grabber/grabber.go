package grabber

import (
	"errors"
)

var ErrStatus = errors.New("unexpected response status")
var ErrThreadUrl = errors.New("incorrect thread url")

type ThreadGrabber struct {
	DownloadingQueue ExecutableQueue
	ProcessingQueue  ExecutableQueue
}

func NewThreadGrabber(downloading, processing ExecutableQueue) *ThreadGrabber {
	grabber := &ThreadGrabber{
		DownloadingQueue: downloading,
		ProcessingQueue:  processing,
	}
	return grabber
}
