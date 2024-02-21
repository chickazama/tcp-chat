# TCP Chat

A command-line group chat application, using TCP sockets, written in Go. Terminal GUI implemented with
[awesome-gocui](https://github.com/awesome-gocui/gocui).

## Server

The server program listens for and accepts socket connections using the TCP version 4 protocol. Once a connection
is established, go-routines are spawned which allow fully-duplexed communication between the client and the server.

Initially, the server sends the client the complete chat history (stored in a text file). Upon receipt of messages from
a connected client, it prepends a timestamp, and broadcasts the message to all connected clients.

The maximum number of simultaneously connected clients is set to 10. When a client closes a connection, the server removes
it from its map of connected clients, and decrements a counter variable tracking the number of connections.

All messages, as well as notifications of client entry and exit, are written to the chat history file, and are logged
server-side to stdout.

## Client

The client program dials the server to establish a connection. Once connected, the user is prompted to enter his/her name.
The program then sets up two terminal views - one which updates the chat history upon receipt of a message, and one which
accepts input from the user to send messages to the server.
