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
    Host string `json:"host"`
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

var rooms map(string)chan UserInfo
var openRoom chan (chan UserInfo)
var closeRoom chan (chan UserInfo)

func main() {
    // Listen for incoming connections.
    listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    queue := make(chan UserInfo, 10) // Buffered channel with capacity of 10
    ins := make(chan Instruction, 10)
    //rooms = make(map[string]Room)
    
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
	go handleTasks(queue) // Potentially need to increase the number of workers
	go handleRooms(openRoom, closeRoom, ins)
	go handleInstructions(conn, ins)
    }
}

// Handles incoming requests and parse response from json to UserInfo struct
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


func handleRooms() {
    for {
	select {
	case room := <- openRoom:
	    go manageRoom(room)
	case room := <- closeRoom:
	    room <- UserInfo{Type:"closeRoom"}
	    delete(rooms, room)
	}
    }
}

func manageRoom(room <-chan UserInfo) {
    defer close(room)
    
    var graph = new(Graph) // TODO: implement Graph
    var tree = new(Graph)
    
    for {
	userInfo := <- room
	
	switch userInfo.Type {
	case "host": 
	    username := userInfo.Host
	    graph.AddNode(username)
	    graph.SetHead(username)
	    ins <- Instruction{Type: "host", Host: username} 
	    
	    if userInfo.Latency { // may be unnecessary
		for _, p := range userInfo.Latency {
		    peername := p.Peer
		    weight := p.Latency
		    graph.AddEdge(peername, username, weight)
		}
	    }
	    
	case "newUser": 
	    username := userInfo.User
	    graph.AddNode(username)
	    for _, p := range userInfo.Latency {
		peername := p.Peer
		weight := p.Latency
		graph.AddEdge(peername, username, weight)
	    }
	    
	case "disconnectedUser": 
	    username := userInfo.User
	    graph.deleteNode(username)
	    if len(graph.Nodes) == 0 {
		return
	    }
	
	case "closeRoom":
	    return
	}
	
	newTree := graph.GetOptimalTree(2) // parameter is the constraint. 2 means a binary tree
	    
	_, _, addedEdges, removedEdges := Graph.Diff(tree, newTree)  // addedNodes, removedEdges, addedEdges, removedEdges := Graph.Diff(tree, newTree) 
	
	for _, edge := range removedEdges {
	    ins <- Instruction{Type:"deletePeerConnection", Parent: edge.Parent, Child: edge.Child}
	}
	
	for _, edge := range addedEdges { // assuming addedEdges are sorted in good orders 
	    ins <- Instruction{Type:"newPeerConnection", Parent: edge.Parent, Child: edge.Child}
	}
	
	tree = newTree
    }
}

func handleTasks(queue <-chan UserInfo) {
    for {
	userInfo := <- queue
	
	switch userInfo.Type {
	case "newUser": newUserHandler(userInfo, ins); break;
	case "host": newHostHandler(userInfo, ins); break;
	case "disconnectedUser": disconnectHandler(userInfo, ins); break;
	}
	fmt.Printf("New task received -> Type: %s  User: %s  Room: %s\n", userInfo.Type, userInfo.User, userInfo.Room)
    }
}

func newUserHandler(userInfo UserInfo, ins chan<- Instruction) {
    roomId := userInfo.Room
    if room, exist := rooms[roomId]; exist {
	room <- userInfo
    } else {
	fmt.Println("ERR: newUserHandler - room doesn't exist")
    }
	/* Send out instructions */
	/* TODO: may need to separate out this part */
	/*
	host := room.getHost()
	if host.Role == "host" { 
	    ins <- Instruction{Type:"newPeerConnection", Parent: host.Name, Child: user.Name}
	} else {
	    fmt.Println("ERR: Host doesn't exist")
	}
    }
    */
}

func newHostHandler(userInfo UserInfo) {
    roomId := userInfo.Room
    if _, exist := rooms[roomId]; !exist {
	room := make(chan UserInfo)
	rooms[roomId] = room
	openRoom <- room
	room <- userInfo
    } else {
	fmt.Println("ERR: newHostHandler - room already exists")
    }
	/*
	user := User{Name: userInfo.User, Role: "host"}
	users := make([]User, 0)
	users = append(users, user)
	room := Room{ID: roomId, Users: users}
	rooms[roomId] = room;
	fmt.Println(room.getUsers())
	ins <- Instruction{Type:"host", Host: user.Name}
	*/
}

func disconnectHandler(userInfo UserInfo, ins chan<- Instruction) {
    roomId := userInfo.Room
    if room, exist := rooms[roomId]; exist {
	room <- userInfo
	//room.removeUser(user)
	
	/* Send out instruction */
	/*
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
	*/
    }
    else {
	fmt.Println("ERR: disconnectHandler - disconnecting from a room non-existing")
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