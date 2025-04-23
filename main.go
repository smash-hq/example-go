package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	for i := 0; i < 10; i++ {
		fmt.Printf("%d--> Hello World\n", i+1)
		time.Sleep(1 * time.Second)
	}
	os.Exit(0)
}
