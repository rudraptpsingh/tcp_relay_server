package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type RelayServer struct {
	host  string
	port  string
	ln    net.Listener
	rooms Rooms
}

type Room struct {
	conn1          net.Conn
	conn2          net.Conn
	timeOfCreation time.Time
}

type Rooms struct {
	roomMap map[string]Room
	sync.Mutex
}

func NewRelayServer(host string, port string) *RelayServer {
	return &RelayServer{
		host: host,
		port: port,
		rooms: Rooms{
			roomMap: make(map[string]Room),
		},
	}
}

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

func (r *RelayServer) CreateRelayServer() {
	ln, err := net.Listen("tcp", r.port)
	if err != nil {
		fmt.Errorf("Failed to create relay serve. Error: %s", err)
		return
	}

	r.ln = ln
}

func (r *RelayServer) AcceptConnections() {
	defer r.ln.Close()
	for {
		conn, err := r.ln.Accept()
		if err != nil {
			fmt.Errorf("Failed to accept connection. Error: %s", err)
			continue
		}

		go r.CreateRoom(conn)
	}
}

func (r *RelayServer) CreateRoom(conn net.Conn) {
	code := make([]byte, 2048)
	// Read the password from the connection
	n, err := conn.Read(code)
	if err != nil {
		fmt.Errorf("Failed to read room code from connection. Error: %f", err)
		return
	}

	roomCode := strings.TrimSpace(string(code[:n]))

	// Check if a room with the same password exists. If so, add the connection to the room.
	r.rooms.Mutex.Lock()
	if room, ok := r.rooms.roomMap[roomCode]; ok {
		room.conn2 = conn
		// notify both the connections that a room have been created.
		room.conn1.Write([]byte("connected"))
		room.conn2.Write([]byte("connected"))
		r.rooms.Mutex.Unlock()
		// now create a pipe between the two connections.
		SetupConnectionPipe(room.conn1, room.conn2)
		return
	}

	// Create a new room as no room exists with the given password.
	r.rooms.roomMap[roomCode] = Room{conn1: conn, timeOfCreation: time.Now()}
	r.rooms.Mutex.Unlock()
}

func main() {
	relay_server := NewRelayServer("", ":1234")
	relay_server.CreateRelayServer()
	relay_server.AcceptConnections()
}
