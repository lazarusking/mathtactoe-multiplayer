package main

import (
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	rooms      map[*Room]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[*Room]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) monitorClientsChannel() {
	for {
		select {
		case message := <-hub.broadcast:
			for client := range hub.clients {
				log.Printf("96%v", len(hub.rooms))
				client.send <- message.encode()
			}
		case client := <-hub.register:
			hub.clients[client] = true
			// user := User{ID: client.UserID, User: "", Status: "setup"}
			// client.conn.WriteJSON(user) //write user struct to client

		case client := <-hub.unregister:
			if ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				// message := &Message{
				// 	Action: UserLeftAction,
				// 	Sender: client,
				// }
				// for client := range hub.clients {
				// 	client.send <- message.encode()
				// }
				// close(client.send)
			}
		}
	}
}
func (hub *Hub) createRoom(id string) *Room {
	room := NewRoom(id)
	go room.RunRoom()
	hub.rooms[room] = true
	return room
}
func (hub *Hub) findRoomByID(id string) *Room {
	var foundRoom *Room
	for room := range hub.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}
	return foundRoom
}
