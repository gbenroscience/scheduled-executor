package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/gbenroscience/scheduled-executor/utils"
)

func timeStampMillis() int {
	return int(time.Now().UnixNano() / 1000000)
}

func main() {

	totalCount := 0

	var wg *sync.WaitGroup = &sync.WaitGroup{}

	utils.NewTimedExecutor(2*time.Second, 2*time.Second).Start(func() {
		totalCount++
		fmt.Printf("%d.%4stime is %d\n", totalCount, " ", timeStampMillis())
		wg.Done()
	}, true)

	wg.Add(10)

	wg.Wait()

}
