package main

import (
	"fmt"
	"net"
)

func GetChannelForConnection(conn net.Conn) (ch chan []byte) {
	ch = make(chan []byte, 1)
	go func() {
		msg := make([]byte, 2048)
		for {
			n, err := conn.Read(msg)
			if err != nil {
				fmt.Errorf("Failed to read from connection. Error: %s", err)
				return
			}

			ch <- msg[:n]
		}
	}()

	return
}

func SetupConnectionPipe(conn1 net.Conn, conn2 net.Conn) {

	fmt.Println("Setting up pipe b/w connections")
	chan1 := GetChannelForConnection(conn1)
	chan2 := GetChannelForConnection(conn2)
	for {
		select {
		case msg := <-chan1:
			if msg == nil {
				// close the connection
				return
			}

			_, err := conn2.Write(msg)
			if err != nil {
				fmt.Errorf("Failed to write to connection. Error: %s", err)
			}

		case msg := <-chan2:
			if msg == nil {
				// close the connection
				return
			}

			_, err := conn1.Write(msg)
			if err != nil {
				fmt.Errorf("Failed to write to connection. Error: %s", err)
			}
		}
	}
}

func CreateRelayServer() {
	relay_server, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Errorf("Failed to create relay serve. Error: %s", err)
	}

	conns := make([]net.Conn, 2)

	defer relay_server.Close()
	for i := 0; i < 2; i++ {
		conns[i], err = relay_server.Accept()
		if err != nil {
			fmt.Errorf("Failed to accept connection. Error: %s", err)
			continue
		}
	}

	// pipe both the client connections so that messages can be relayed
	SetupConnectionPipe(conns[0], conns[1])
}

func main() {
	CreateRelayServer()
}
