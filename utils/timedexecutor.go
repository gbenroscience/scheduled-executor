package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ScheduledExecutor struct {
	initialDelay time.Duration
	delay        time.Duration
	ticker       *time.Ticker
	sigs         chan os.Signal
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewTimedExecutor(initialDelay time.Duration, delay time.Duration) *ScheduledExecutor {
	return &ScheduledExecutor{
		initialDelay: initialDelay,
		delay:        delay,
	}
}

// Start begins the execution of the task at the specified intervals.
// If runAsync is true, the task will be executed in a separate goroutine.
// If runAsync is false, the task will block the goroutine until it completes.
// If startAsync is false, the executor will start asynchronously. So the caller can continue without waiting for the task to start.
// If startAsync is false, the caller(of Start) will wait until the task starts executing.
// The task will be stopped when the context is canceled.
func (se *ScheduledExecutor) Start(task func(), runAsync bool, startAsync bool) {
	se.ctx, se.cancel = context.WithCancel(context.Background())
	se.sigs = make(chan os.Signal, 1)
	signal.Notify(se.sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-se.sigs:
			se.cancel()
		case <-se.ctx.Done():
			return
		}
	}()

	go func() {
		time.Sleep(se.initialDelay)
		if runAsync {
			go task()
		} else {
			task()
		}
		se.ticker = time.NewTicker(se.delay)
		defer se.ticker.Stop()
		for {
			select {
			case <-se.ticker.C:
				if runAsync {
					go task()
				} else {
					task()
				}
			case <-se.ctx.Done():
				return
			}
		}
	}()

	if !startAsync {
		<-se.ctx.Done()
	}
}

func (se *ScheduledExecutor) Close() {
	if se.cancel != nil {
		se.cancel()
	}
	if se.ticker != nil {
		se.ticker.Stop()
	}
}
