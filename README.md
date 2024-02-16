# TCP Relay Server
This project implements a TCP relay server. This server listens to two connections and pipes these connections so that the connections can start interacting with each other.
The TCP relay server creates a pipe between the two connections and just relays the messages to and fro.

To start a TCP relay server, run
```
go run relay.go
```

To start the client connections, run the following command in two terminals
```
go run client.go
```

To close a chat session, type exit as message. This closes the current client connection.