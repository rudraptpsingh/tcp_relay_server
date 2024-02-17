package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	// Create new client connction
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Errorf("Failed to dial to the address. Error: %s", err)
	}

	fmt.Println("Enter room code")
	roomInput := bufio.NewReader(os.Stdin)
	roomCode, _ := roomInput.ReadString('\n')
	// Send room code to the relay server
	conn.Write([]byte(roomCode))
	confirmation := make([]byte, 2048)
	var wg sync.WaitGroup
	fmt.Println("Waiting for users in room...")
	n, err := conn.Read(confirmation)
	if err != nil {
		fmt.Errorf("Failed to read from connection. Error: %s", err)
		return
	}

	if strings.TrimSpace(string(confirmation[:n])) == "connected" {
		fmt.Println("Room connected")
		wg.Add(1)
		go StartChat(conn, &wg)
	}

	wg.Wait()
}

func StartChat(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	fmt.Println("Chatting...")

	done := make(chan struct{})
	// Read message from user and sent to the other user
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Errorf("Failed to read from stdin. Error: %s", err)
				done <- struct{}{}
				return
			}

			if strings.TrimSpace(msg) == "exit" {
				done <- struct{}{}
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
				done <- struct{}{}
				return
			}

			fmt.Print("<- " + string(msg[:n]))
		}
	}()

	<-done
	fmt.Println("Exiting...")
}
