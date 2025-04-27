package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// 打开/创建日志文件
	logFile, err := os.OpenFile("/var/log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// 将标准输出重定向到日志文件
	os.Stdout = logFile
	i := 0
	for {
		fmt.Printf("%d--> Hello World\n", i+1)
		time.Sleep(1 * time.Second)
		i++
	}
	select {}
}
