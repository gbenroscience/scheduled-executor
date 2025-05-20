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

const SHUT_DOWN = 1

func NewTimedExecutor(initialDelay time.Duration, delay time.Duration) ScheduledExecutor {
	return ScheduledExecutor{
		delay:  delay,
		ticker: *time.NewTicker(initialDelay),
		quit:   make(chan int),
	}
}

// Start .. process() is the function to run periodically , runAsync detects if the function should block the executor when running or not. It blocks when false
func (se ScheduledExecutor) Start(task func(), runAsync bool) {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			fmt.Println("Scheduler stopping...")
			se.close()
			fmt.Println("Scheduler stopped.")
		}()
		firstExec := true
		for {
			fmt.Println("IN the loop")
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
				fmt.Println("case evaluated, other cases will be ignored for now - 1")
			case a := <-se.quit:
				if a == SHUT_DOWN {
					fmt.Printf("returning here - 2, a= %d\n", a)
					return
				}
				fmt.Println("keep idling sweet golang - 2")

			case <-sigs:
				fmt.Println("AWW AWW AWW - 3")
				fmt.Println("breaking out of select here - 3")
				return
			}
		}
		fmt.Println("OUT of the loop - 4")

	}()
	fmt.Println("OUT of goroutine - 5")

}

func (se *ScheduledExecutor) Close() error {
	go func() {
		fmt.Println("Closing scheduler...")
		se.quit <- SHUT_DOWN
	}()
	return nil
}
func (se *ScheduledExecutor) close() {
	se.ticker.Stop()
}
