package grabber

import (
	"context"
	"errors"
	"time"
)

var ErrShutdown = errors.New("queue shutdown")
var ErrOverload = errors.New("queue overloaded")

type ExecutableTask interface {
	Execute(ctx context.Context) error
	Cancel(err error)
	Notify(position int)
}

type ExecutableQueue interface {
	Submit(ctx context.Context, task ExecutableTask)
	Stop()
}

type executableTaskDescriptor struct {
	context   context.Context
	notify    time.Time
	notifymin time.Duration
	task      ExecutableTask
}

func (desc *executableTaskDescriptor) Notify(position int) {
	now := time.Now()
	if now.Sub(desc.notify) < desc.notifymin && position > 0 {
		return
	}
	desc.notify = now
	desc.task.Notify(position)
}

type executableQueue struct {
	buffer    int
	context   context.Context
	cancel    context.CancelFunc
	rate      time.Duration
	limit     time.Duration
	notifymin time.Duration
	queue     chan *executableTaskDescriptor
	workers   []*executableQueueWorker
	awaiting  []int
	complete  chan int
}

type executableQueueWorker struct {
	process chan *executableTaskDescriptor
}

func NewExecutableQueue(workers, buffer int, limit, notifymin, rate time.Duration) ExecutableQueue {
	ctx, cancel := context.WithCancel(context.Background())
	ex := &executableQueue{
		buffer:    buffer,
		context:   ctx,
		cancel:    cancel,
		rate:      rate,
		limit:     limit,
		notifymin: notifymin,
		queue:     make(chan *executableTaskDescriptor, workers*buffer),
		complete:  make(chan int),
	}
	for i := 0; i < workers; i++ {
		process := ex.worker(i)
		ex.workers = append(ex.workers, &executableQueueWorker{
			process: process,
		})
		ex.awaiting = append(ex.awaiting, i)
	}

	go func() {
		buffer := make([]*executableTaskDescriptor, workers*buffer)
		var front, back, length int

		drain := func(worker int) {
			ex.awaiting = append(ex.awaiting, worker)
			if length == 0 {
				return
			}
			for i := 0; i < length; i++ {
				buffer[(back+i)%len(buffer)].Notify(i)
			}
			task := buffer[back]
			back = (back + 1) % len(buffer)
			length -= 1

			worker, ex.awaiting = ex.awaiting[0], ex.awaiting[1:]
			ex.workers[worker].process <- task
		}

		for {
			if length < len(buffer) {
				select {
				case task, ok := <-ex.queue:
					if !ok {
						for _, w := range ex.workers {
							close(w.process)
						}
						for _, task := range buffer {
							go task.task.Cancel(ErrShutdown)
						}
						return
					}
					if len(ex.awaiting) > 0 {
						go task.Notify(0)
						var worker int
						worker, ex.awaiting = ex.awaiting[0], ex.awaiting[1:]
						ex.workers[worker].process <- task
					} else {
						buffer[front] = task
						front = (front + 1) % len(buffer)
						length += 1
						go task.Notify(length)
					}
				case worker := <-ex.complete:
					drain(worker)
				}
			} else {
				worker := <-ex.complete
				drain(worker)
			}
		}
	}()

	return ex
}

func (ex *executableQueue) worker(worker int) (process chan *executableTaskDescriptor) {
	process = make(chan *executableTaskDescriptor)
	go func() {
		for desc := range process {
			start := time.Now()
			ctx, cancel := context.WithTimeout(desc.context, ex.limit)
			go func() {
				defer func() {
					_ = recover()
					cancel()
				}()
				err := desc.task.Execute(ctx)
				if err != nil {
					desc.task.Cancel(err)
				}
			}()
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				go desc.task.Cancel(ctx.Err())
			}
			done := time.Now()
			remain := done.Sub(start)
			if remain < ex.rate {
				time.Sleep(ex.rate - remain)
			}
			ex.complete <- worker
		}
	}()
	return
}

func (ex *executableQueue) Submit(ctx context.Context, task ExecutableTask) {
	select {
	case ex.queue <- &executableTaskDescriptor{
		context:   ctx,
		task:      task,
		notifymin: ex.notifymin,
	}:
	default:
		task.Cancel(ErrOverload)
	}
}

func (ex *executableQueue) Stop() {
	close(ex.queue)
}
