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

type User struct {
    Name string `json:"name"`
    Role string `json:"role"` //enum: "host", "user"
}

type Room struct {
    ID string `json:"roomID"`
    Users []User `json:"users"`
}

func (room *Room) addUser(user User) {
    room.Users = append(room.Users, user)
}

func (room *Room) getUsers() []User {
    return room.Users
}

func (room *Room) removeUser(user User) {
    room.Users = room.Users
    for i, u range room.Users {
	if u.Name == user.Name {
	    room.Users = append(room.Users[:i], room.Users[i+1:]...) // The ... is essential
	    return
	}
    }
}

const (
    CONN_HOST = "localhost"
    CONN_PORT = "8888"
    CONN_TYPE = "tcp"
)

var rooms map(string)Room

func main() {
    // Listen for incoming connections.
    listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    queue := make(chan UserInfo, 10) // Buffered channel with capacity of 10
    rooms = make(map(string)Room)
    
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
	go handleRequest(conn, queue)
	go handleTasks(conn, queue) // Potentially need to increase the number of workers
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn, queue chan<- UserInfo) {
    defer conn.Close()
    
    input := bufio.NewScanner(conn)
    var userInfo UserInfo
    
    for input.Scan() {
	text := input.Text()
	byte_text := []byte(text)
	err := json.Unmarshal(byte_text, &userInfo)
	if err != nil {
	    continue
	}
	queue <- userInfo // send userInfo to task queue
    }
}

func handleTasks(conn net.Conn, queue <-chan UserInfo) {
    for {
	var userInfo UserInfo
	userInfo = <- queue
	
	switch userInfo.Type {
	case "newUser": fmt.Println("newUser") 
	case "host": fmt.Println("host")
	case "disconnectedUser": fmt.Println("disconnectedUser")
	}
	
	fmt.Fprintf(conn, "Type: %s  User: %s  Room: %s  Host: %s", userInfo.Type, userInfo.User, userInfo.Room, userInfo.Host)
    }
}

func newUserHandler(userInfo UserInfo) {
    roomId := userInfo.Room
    if room, exist := rooms[roomId]; exist {
	user := User{Name: UserInfo.User, Role: "user"}
	room.addUser(user)
    }
}

func newHostHandler(userInfo UserInfo) {
    roomId := userInfo.Room
    if room, exist := rooms[roomId]; !exist {
	user := User{Name: UserInfo.User, Role: "host"}
	users := make([]User)
	users = append(users, user)
	room := Room{ID: roomID, Users: users]
	rooms[roomId] = room;
    }
}

func disconnectHandler(userInfo UserInfo) {
    roomId := UserInfo.Room
    if room, exist := rooms[roomId]; exist {
	room.removeUser(userInfo.User)
	if len(room.getUsers())==0 {
	    delete(rooms, roomId)
	}
    }
}