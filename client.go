package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Create new client connction
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Errorf("Failed to dial to the address. Error: %s", err)
	}

	defer conn.Close()
	fmt.Println("Chatting...")

	done := make(chan bool, 1)
	// Read message from user and sent to the other user
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Errorf("Failed to read from stdin. Error: %s", err)
				done <- true
				return
			}

			if strings.TrimSpace(msg) == "exit" {
				done <- true
				return
			}

			conn.Write([]byte(msg))
		}
	}()

	// Read and print message from another user
	go func() {
		msg := make([]byte, 2048)
		for {
			n, err := conn.Read(msg)
			if err != nil {
				fmt.Errorf("Failed to read from stdin. Error: %s", err)
				done <- true
				return
			}

			fmt.Print("<- " + string(msg[:n]))
		}
	}()

	<-done
	fmt.Println("Exiting...")
}
