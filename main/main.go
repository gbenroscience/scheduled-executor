package main

import (
	"com.github.gbenroscience/scheduled-executor/utils"
	"fmt"
	"time"
)


func timeStampMillis() int {
	return int(time.Now().UnixNano() / 1000000)
}

func main()  {


	totalCount := 0


		utils.NewTimedExecutor(5 * time.Second , 2 * time.Second).Start(func() {
			totalCount++
			fmt.Printf("%d.%4stime is %d\n" ,totalCount , " ", timeStampMillis())
		} , true)


		time.Sleep(time.Minute)



}
