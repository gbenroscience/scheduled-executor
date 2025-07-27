package main

import (
	"fmt"
	"time"

	"github.com/gbenroscience/scheduled-executor/utils"
)

func timeStampMillis() int {
	return int(time.Now().UnixNano() / 1000000)
}

func main() {

	totalCount := 0
	const MAX_CYCLES = 10

	sc := utils.NewTimedExecutor(2*time.Second, 500*time.Millisecond)
	sc1 := utils.NewTimedExecutor(2*time.Second, 500*time.Millisecond)

	sc.Start(func() {
		totalCount++
		fmt.Printf("sc:---%d.%4stime is %d\n", totalCount, " ", timeStampMillis())
		if totalCount > MAX_CYCLES {
			sc.Close()
		}
	}, true)

	totalCount1 := 0
	sc1.Start(func() {
		totalCount1++
		fmt.Printf("sc1:---%d.%4stime is %d\n", totalCount1, " ", timeStampMillis())
		if totalCount1 > MAX_CYCLES {
			sc1.Close()
		}
	}, true)

}
