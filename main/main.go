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

	sc.Start(func() {
		totalCount++
		fmt.Printf("%d.%4stime is %d\n", totalCount, " ", timeStampMillis())
		if totalCount > MAX_CYCLES {
			sc.Close()
		}
	}, true)

}
