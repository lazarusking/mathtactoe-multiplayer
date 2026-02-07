package main

type Hub struct {
	clients    map[*Client]bool
	rooms      map[*Room]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	roomClosed chan string // Signal channel for room deletion
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[*Room]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		roomClosed: make(chan string),
	}
}

func (hub *Hub) monitorClientsChannel() {
	for {
		logger.Printf("The hub rooms %+v ,%v", hub.rooms, len(hub.rooms))
		for v := range hub.rooms {
			logger.Printf("hub room %v", v.clients)
			for c := range v.clients {
				logger.Printf("client %s %s", c.Name, c.ID)

			}
		}
		select {
		case message := <-hub.broadcast:
			for client := range hub.clients {
				client.send <- message.encode()
			}
		case client := <-hub.register:
			hub.clients[client] = true

		case client := <-hub.unregister:
			delete(hub.clients, client)
		case roomID := <-hub.roomClosed:
			hub.deleteRoom(roomID)
		}
	}
}
func (hub *Hub) createRoom(id string) *Room {
	room := NewRoom(id, hub.roomClosed)
	go room.RunRoom()
	hub.rooms[room] = true
	return room
}
func (h *Hub) deleteRoom(roomID string) {
	room := h.findRoomByID(roomID)
	if room != nil {
		delete(h.rooms, room)
		logger.Printf("Room %s has been deleted from hub", room.ID)
	}
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
