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

func (se *ScheduledExecutor) Start(task func(), runAsync bool) {
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

	<-se.ctx.Done()
}

func (se *ScheduledExecutor) Close() {
	if se.cancel != nil {
		se.cancel()
	}
	if se.ticker != nil {
		se.ticker.Stop()
	}
}
