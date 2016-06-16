package main

import (
    "fmt"
    "encoding/json"
    "net"
    "bufio"
    "sync"
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

type Instruction struct {
    Type string `json:"type"` //enum: "newPeerConnection" "deletePeerConnection"
    Parent string `json:"parent"`
    Child string `json:"child"` 
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
    for i, u := range room.Users {
	if u.Name == user.Name {
	    room.Users = append(room.Users[:i], room.Users[i+1:]...) // The ... is essential
	    return
	}
    }
}

func (room *Room) getHost() User {
    users := room.getUsers()
    for _, u := range users {
	if u.Role == "host" {
	    return u
	}
    }
    return User{}
}

const (
    CONN_HOST = "localhost"
    CONN_PORT = "8888"
    CONN_TYPE = "tcp"
)

var rooms map[string]Room

var mu sync.Mutex


func main() {
    // Listen for incoming connections.
    listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    queue := make(chan UserInfo, 10) // Buffered channel with capacity of 10
    ins := make(chan Instruction, 10)
    rooms = make(map[string]Room)
    
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
	go handleRequests(conn, queue)
	go handleTasks(queue, ins) // Potentially need to increase the number of workers
	go handleInstructions(conn, ins)
    }
}

// Handles incoming requests.
func handleRequests(conn net.Conn, queue chan<- UserInfo) {
    fmt.Println("handleRequests is working")
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

func handleTasks(queue <-chan UserInfo, ins chan<- Instruction) {
    fmt.Println("handleTasks is working")
    for {
	var userInfo UserInfo
	userInfo = <- queue
	
	switch userInfo.Type {
	case "newUser": newUserHandler(userInfo, ins) 
	case "host": newHostHandler(userInfo, ins)
	case "disconnectedUser": disconnectHandler(userInfo, ins)
	}
	fmt.Printf("New task received -> Type: %s  User: %s  Room: %s\n", userInfo.Type, userInfo.User, userInfo.Room)
    }
}

func newUserHandler(userInfo UserInfo, ins chan<- Instruction) {
    fmt.Println("newUserHandler triggered")
    roomId := userInfo.Room
    //mu.Lock()
    if room, exist := rooms[roomId]; exist {
	fmt.Println(room.getHost())
	user := User{Name: userInfo.User, Role: "user"}
	room.addUser(user)
	
	/* Send out instructions */
	/* TODO: may need to separate out this part */
	host := room.getHost()
	if host.Role == "host" { 
	    ins <- Instruction{Type:"newPeerConnection", Parent: host.Name, Child: user.Name}
	} else {
	    fmt.Println("ERR: Host doesn't exist")
	}
	fmt.Println(room.getUsers())
    }
    //mu.Unlock()
}

func newHostHandler(userInfo UserInfo, ins chan<- Instruction) {
    fmt.Println("newHostHandler triggered")
    roomId := userInfo.Room
    //mu.Lock()
    if _, exist := rooms[roomId]; !exist {
	user := User{Name: userInfo.User, Role: "host"}
	users := make([]User, 0)
	users = append(users, user)
	room := Room{ID: roomId, Users: users}
	rooms[roomId] = room;
	fmt.Println(room.getUsers())
    }
    //mu.Unlock()
}

func disconnectHandler(userInfo UserInfo, ins chan<- Instruction) {
    fmt.Println("disconnectHandler triggered")
    roomId := userInfo.Room
    if room, exist := rooms[roomId]; exist {
	user := User{Name:userInfo.User, Role:"user"}
	room.removeUser(user)
	
	/* Send out instruction */
	host := room.getHost()
	
	if host.Role == "host" {
	    ins <- Instruction{Type:"deletePeerConnection", Parent: host.Name, Child: user.Name}
	} else {
	    fmt.Println("ERR: Host doesn't exist")
	}
	
	if len(room.getUsers())==0 {
	    delete(rooms, roomId)
	}
	fmt.Println(room.getUsers())
    }
}

func handleInstructions(conn net.Conn, ins <-chan Instruction) {
    fmt.Println("handleInstructions is working")
    for {
	instruction := <- ins
	str, err := json.Marshal(instruction)
	if err != nil {
	    fmt.Println("Error listening:", err.Error())
	    continue
	}
	fmt.Fprintf(conn, "%s", string(str))
	fmt.Println("Instruction Sent")
    }
}