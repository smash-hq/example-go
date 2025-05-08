package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request on /hello")
	fmt.Fprintln(w, "hello")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	port := ":8848"
	fmt.Printf("Server listening on port %s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
