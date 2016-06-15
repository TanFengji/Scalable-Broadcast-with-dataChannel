package main

import (
    "fmt"
    "encoding/json"
    "net"
	"bufio"
)

type PeerInfo struct {
    Peer string `json:"peer"`
    Latency int `json:"latency"`
}

type UserInfo struct {
    Type string `json:"type" `
    User string `json:"user" `
    Room string `json:"room" `
    Host string `json:"host" `
    Latency []PeerInfo `json:"latency"`
}

const (
    CONN_HOST = "localhost"
    CONN_PORT = "8889"
    CONN_TYPE = "tcp"
)

var peer PeerInfo

func main() {
    // Listen for incoming connections.
    listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    
    if err != nil {
		fmt.Println("Error listening:", err.Error())
    }
    
    // Close the listener when the application closes.
    defer listener.Close()
    
    for {
	// Listen for an incoming connection.
	conn, err := listener.Accept()

	if err != nil {
	    fmt.Println("Error accepting: ", err.Error())
	    continue
	}

	// Handle connections in a new goroutine.
	go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    defer conn.Close()

	input := bufio.NewScanner(conn)

	for input.Scan() {
		text := input.Text()
		byte_text := []byte(text)
		json.Unmarshal(byte_text, &peer)
		fmt.Fprintf(conn, "Peer: %s \t Latency: %d \n ", peer.Peer, peer.Latency)
	}
}
