package main

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client represents a connected WebSocket client.
type Client struct {
	hub   *Hub
	rooms map[*Room]bool
	conn  *websocket.Conn
	send  chan []byte // Buffered channel of outbound messages.
	// Add any client-specific properties here
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal message %s", err)
		return
	}
	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		log.Printf("30%v", message.Target)
		roomId := message.Target.GetId()
		if room := client.hub.findRoomByID(roomId); room != nil {
			room.broadcast <- &message
		}
	case SendGameAction:
		log.Printf("30: %v", message.Target.ID)
		roomId := message.Target.GetId()
		if room := client.hub.findRoomByID(roomId); room != nil {
			room.gamebroadcast <- &message
		}

	case JoinRoomAction:
		roomId := message.Message
		client.joinRoom(roomId)
	case LeaveRoomAction:
		roomId := message.Message
		room := client.hub.findRoomByID(roomId)
		if room == nil {
			return
		}
		log.Println("leaving room:", roomId)
		room.unregister <- client
	default:
		break
	}
}

func (client *Client) GetID() string {
	return client.ID.String()
}
func (client *Client) joinRoom(roomId string) *Room {
	// var room *Room
	room := client.hub.findRoomByID(roomId)
	// if id, err := uuid.Parse(roomId); id == uuid.Nil && err != nil {
	// 	log.Println("tis nil", id, err)
	// 	room = client.hub.createRoom(uuid.NewString())
	// }
	if roomId == "" {
		room = client.hub.createRoom(uuid.NewString())
	}
	if room == nil {
		// log.Println(room)
		room = client.hub.createRoom(roomId)
	}
	// if _, ok := client.rooms[room]; !ok {
	// 	client.rooms[room] = true
	// 	room.register <- client
	// 	client.notifyRoomJoined(room)
	// }
	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.register <- client

		client.notifyRoomJoined(room)
	}
	return room
}
func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}
func (client *Client) notifyRoomJoined(room *Room) {
	message := &Message{
		Action:  JoinRoomAction,
		Message: room.ID,
		Target:  room,
		Sender:  client}
	client.send <- message.encode()

}
func (client *Client) disconnect() {
	client.hub.unregister <- client //unregister client from hub and rooms
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	client.conn.Close()
}

// listen handles incoming messages from a client.
func (c *Client) listen() {
	defer func() {
		// c.hub.unregister <- c // Send the client to be unregistered
		// c.conn.Close()
		c.disconnect()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageContent, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		messageContent = bytes.TrimSpace(bytes.Replace(messageContent, newline, space, -1))
		c.handleNewMessage(messageContent)
	}
}

// send messages to the client.
func (c *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
