package grabber

import (
	"context"
	"log"
	"testing"
	"time"
)

type TestTask struct {
	complete chan struct{}
	n        int
}

func (order *TestTask) Context() context.Context {
	return context.Background()
}

func (order *TestTask) Execute(ctx context.Context) error {
	close(order.complete)
	return nil
}

func (order *TestTask) Cancel(err error) {
	log.Println("cancel", order.n, err)
}

func (order *TestTask) Notify(position int) {
	if order.n == 99 {
		log.Println("notify", order.n, position)
	}
}

func TestTaskQueue(t *testing.T) {
	queue := NewExecutableQueue(20, 100, time.Second*2, time.Second, time.Second)
	var task *TestTask
	for i := 0; i < 100; i++ {
		task = &TestTask{
			n: i,
		}
		if i == 99 {
			task.complete = make(chan struct{})
		}
		queue.Submit(context.Background(), task)
	}
	<-task.complete
}
