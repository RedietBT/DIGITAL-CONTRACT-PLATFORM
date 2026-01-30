package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("🚀 Service is running inside Docker...")
		time.Sleep(10 * time.Second)
	}
}