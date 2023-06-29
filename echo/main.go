package main

import (
	"fmt"
)

func main() {
	s := NewServer(":8090")
	if err := s.Start(); err != nil {
		fmt.Printf("Error running Websocket Server: %s\n", err.Error())
	}
}
