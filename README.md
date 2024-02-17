# TCP Relay Server
This project implements a TCP relay server. This server listens to connections and creates a new room. If a connection provides room info of a room that already exists, the two connections are piped and messages are relayed to and fro.

To start a TCP relay server, run
```
go run relay.go
```

To start the client connections, run the following command in two terminals
```
go run client.go
```

Enter the room code. This room code is shared between two users who want to exchange messages.

To close a chat session, type exit as message. This closes the current client connection.