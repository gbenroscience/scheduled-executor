package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ScheduledExecutor struct {
	delay  time.Duration
	ticker time.Ticker
	quit   chan int
}

func NewTimedExecutor(initialDelay time.Duration, delay time.Duration) ScheduledExecutor {
	return ScheduledExecutor{
		delay:  delay,
		ticker: *time.NewTicker(initialDelay),
		quit:   make(chan int),
	}
}


//Start .. process() is the function to run periodically , runAsync detects if the function should block the executor when running or not. It blocks when false
func (se ScheduledExecutor) Start(task func(), runAsync bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			se.close()
			fmt.Println("Scheduler stopped!!")
		}()
		firstExec := true
		for {
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

				break
			case <-se.quit:
				return
			case <-sigs:
				_ = se.Close()
				break
			}
		}

	}()

}

func (se *ScheduledExecutor) Close() error {
	go func() {
		se.quit <- 1
	}()
	return nil
}
func (se *ScheduledExecutor) close() {
	close(se.quit)
	se.ticker.Stop()
}
