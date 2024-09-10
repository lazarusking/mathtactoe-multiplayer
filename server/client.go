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
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 512                 // Maximum message size allowed from peer.
)

// Client represents a connected WebSocket client.
type Client struct {
	hub   *Hub
	rooms map[*Room]bool
	conn  *websocket.Conn
	send  chan []byte // Buffered channel of outbound messages.
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
}

func (client *Client) GetID() string {
	return client.ID.String()
}

// handleNewMessage processes incoming messages and performs the appropriate actions based on the message's action.
func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal message %s", err)
		return
	}

	log.Printf("Handling Message: %v........Client:%s", message.Action, client.ID)
	if message.Sender != nil {
		client.Name = message.Sender.Name
		log.Printf("Sender is %s", client.Name)
	} else {
		log.Printf("Sender is nil in the received message")
	}

	message.Sender = client
	switch message.Action {
	case SendMessageAction:
		client.sendMessage(message)
	case SendGameAction:
		client.sendGameAction(message)
	case CreateRoomAction:
		client.createAndNotifyRoom()
	case JoinRoomAction:
		client.joinRoom(message.Message)
	case LeaveRoomAction:
		client.leaveRoom(message.Message)
	case StartGameAction:
		client.startGame(message)
	default:
		break
	}
}

func (client *Client) sendMessage(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.broadcast <- &message
	}
}

func (client *Client) sendGameAction(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.gamebroadcast <- &message
	}
}

func (client *Client) createAndNotifyRoom() {
	room := client.createRoom()
	message := &Message{
		Action:  CreateRoomAction,
		Message: room.ID,
		Sender:  client,
	}
	client.send <- message.encode()
}

func (client *Client) leaveRoom(roomId string) {
	room := client.hub.findRoomByID(roomId)
	if room == nil {
		return
	}
	log.Println("Leaving room:", roomId)
	room.unregister <- client
}

func (client *Client) startGame(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.startGame(client)
	}
}

func (client *Client) createRoom() *Room {
	roomId := uuid.NewString()            // Generate a new room ID
	room := client.hub.createRoom(roomId) // Create the room
	log.Printf("Room created: %s by client %s", room.ID, client.ID)
	return room
}

func (client *Client) joinRoom(roomId string) *Room {
	room := client.hub.findRoomByID(roomId)
	if roomId == "" {
		room = client.hub.createRoom(uuid.NewString())
	}
	if room == nil {
		room = client.hub.createRoom(roomId)
	}
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
		client.notifyRoomJoined(room)
	}
	return room
}

func (client *Client) isInRoom(room *Room) bool {
	_, ok := client.rooms[room]
	return ok
}

func (client *Client) notifyRoomJoined(room *Room) {
	log.Printf("Client: %s,%s joined room: %s", client.ID, client.Name, room.ID)
	message := &Message{
		Action:  JoinRoomAction,
		Message: room.ID,
		Target:  room,
		Sender:  client,
	}
	client.send <- message.encode()
}

func (client *Client) disconnect() {
	client.hub.unregister <- client // Unregister client from hub and rooms
	for room := range client.rooms {
		room.unregister <- client
		if len(room.clients) == 0 {
			client.hub.deleteRoom(room.ID)
			room.done <- true
			log.Printf("Room %s has been deleted due to no active clients", room.ID)
		}
	}
	close(client.send)
	client.conn.Close()
}

// listen handles incoming messages from a client.
func (client *Client) listen() {
	defer client.disconnect()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, messageContent, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		messageContent = bytes.TrimSpace(bytes.Replace(messageContent, newline, space, -1))
		client.handleNewMessage(messageContent)
	}
}

// write sends messages to the client.
func (client *Client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		client.conn.Close()
		ticker.Stop()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
