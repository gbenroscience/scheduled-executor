package utils

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ScheduledExecutor struct {
	delay    time.Duration
	ticker   time.Ticker
	sigs     chan os.Signal
	shutdown bool
}

func NewTimedExecutor(initialDelay time.Duration, delay time.Duration) ScheduledExecutor {
	return ScheduledExecutor{
		delay:  delay,
		ticker: *time.NewTicker(initialDelay),
	}
}

// Start .. process() is the function to run periodically , runAsync detects if the function should block the executor when running or not. It blocks when false
func (se *ScheduledExecutor) Start(task func(), runAsync bool) {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	se.shutdown = false
	se.sigs = make(chan os.Signal, 1)
	signal.Notify(se.sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-se.sigs // Block until a signal is received
		cancel()
	}()

	firstExec := true

	defer func() {
		se.close()
		close(se.sigs)
	}()
	for {
		if se.shutdown {
			return
		}
		select {
		case <-se.ticker.C:

			if firstExec {
				se.ticker.Stop()
				se.ticker = *time.NewTicker(se.delay)
				firstExec = false
			}

			if runAsync {
				go task()
			} else {
				task()
			}
		case <-ctx.Done():
			return
		default:
			if se.shutdown {
				return
			}
		}
	}

}

func (se *ScheduledExecutor) Close() error {
	se.shutdown = true
	return nil
}

func (se *ScheduledExecutor) close() {
	se.ticker.Stop()
}
