package main

import (
	"fmt"
	"github.com/gbenroscience/scheduled-executor/utils"
	"time"
)


func timeStampMillis() int {
	return int(time.Now().UnixNano() / 1000000)
}

func main()  {


	totalCount := 0

		utils.NewTimedExecutor(2 * time.Second , 2 * time.Second).Start(func() {
			totalCount++
			fmt.Printf("%d.%4stime is %d\n" ,totalCount , " ", timeStampMillis())
		} , true)


		time.Sleep(time.Minute)



}
