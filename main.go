package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	// 启动一个 goroutine 循环打印 hello world
	go func() {
		i := 1
		for {
			fmt.Printf("%d-->world hello\n", i)
			time.Sleep(10 * time.Millisecond)
			i++
		}
	}()

	// 设置 HTTP 处理器
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "hello world")
		if err != nil {
			return
		}
		fmt.Println("request handler, hello...")
	})

	// 启动 HTTP 服务器监听 3000 端口
	log.Println("HTTP server listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
