package main

import (
    "fmt"
    "encoding/json"
    "net"
    "os"
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
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {
    // Listen for incoming connections.
    listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    
    if err != nil {
	fmt.Println("Error listening:", err.Error())
	os.Exit(1)
    }
    
    // Close the listener when the application closes.
    defer listener.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    
    for {
	// Listen for an incoming connection.
	conn, err := listener.Accept()
	if err != nil {
	    fmt.Println("Error accepting: ", err.Error())
	    //os.Exit(1)
	    continue
	}
	// Handle connections in a new goroutine.
	go handleRequest(listener)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    defer conn.Close()
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)
    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
	fmt.Println("Error reading:", err.Error())
    }
    // Send a response back to person contacting us.
    conn.Write([]byte("Message received."))
    // Close the connection when you're done with it.
}
