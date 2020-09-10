package main

import (
	"fmt"
	"time"
	"com.github.gbenroscience/scheduled-executor/utils"

)



func main()  {


	totalCount := 0


		utils.NewTimedExecutor(5 * time.Second , 2 * time.Second).Start(func() {
			totalCount++
			fmt.Printf("%d.%4stime is %d\n" ,totalCount , " ", utils.CurrentTimeStamp())
		} , true)


		time.Sleep(time.Minute)



}
