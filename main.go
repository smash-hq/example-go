package main

import (
	"fmt"
	"net/http"
	"time"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request on /hello")
	fmt.Fprintln(w, "hello")
}

func main() {
	for i := 0; i < 100; i++ {
		fmt.Println("Hello World")
		time.Sleep(500 * time.Millisecond)
	}
}
