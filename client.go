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
	// Create new client connection
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to dial to the address. Error: %s\n", err)
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("Chatting...")

	var wg sync.WaitGroup
	wg.Add(2)

	// Read message from user and send to the other user
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("-> ")
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read from stdin. Error: %s\n", err)
				return
			}

			msg = strings.TrimSpace(msg)
			if msg == "exit" {
				return
			}

			_, err = conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to send message. Error: %s\n", err)
				return
			}
		}
	}()

	// Read and print message from another user
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			msg := scanner.Text()
			fmt.Printf("<- %s\n", msg)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read from server. Error: %s\n", err)
		}
	}()

	wg.Wait()
	fmt.Println("Exiting...")
}
