package main

import (
    "fmt"
)

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
    for i, u := range room.Users {
	if u.Name == user.Name {
	    room.Users = append(room.Users[:i], room.Users[i+1:]...) // The ... is essential
	    return
	}
    }
}

func main() {
    user1 := User{Name: "aaron", Role: "host"}
    user2 := User{Name: "aaron2", Role: "user"}
    user3 := User{Name: "aaron3", Role: "user"}
    
    room := new(Room)
    fmt.Println(room.getUsers())
    
    room.addUser(user1)
    fmt.Println(room.getUsers())
    
    room.addUser(user2)
    fmt.Println(room.getUsers())
    
    room.addUser(user3)
    fmt.Println(room.getUsers())
    
    room.removeUser(user2)
    fmt.Println(room.getUsers())
}

