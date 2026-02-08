package main

import (
	"bytes"
	"encoding/json"
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
	IsBot bool        `json:"isBot"`
	ID    uuid.UUID   `json:"id"`
	Name  string      `json:"name"`
}

func NewBotClient(name string) *Client {
	return &Client{
		ID:    uuid.New(),
		Name:  name,
		IsBot: true,
		send:  make(chan []byte, 256),
	}
}
func (client *Client) GetID() string {
	return client.ID.String()
}

// handleNewMessage processes incoming messages and performs the appropriate actions based on the message's action.
func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		logger.Printf("Error on unmarshal message %v :%v -- %s %v", err, message.Action, message.Data, message.Sender)
		// logger.Panic(err)
		return
	}

	logger.Printf("Handling Message: %v ---> Client:%s", message.Action, client.ID)
	logger.Printf("Message Data: %+v", string(message.Data))
	logger.Printf("Message Sender: %+v", message.Sender)
	if message.Sender != nil {
		client.Name = message.Sender.Name
	}

	message.Sender = client
	//not unmarshalling gave me a headache because of json.rawMessage and quote issues if I used string()
	switch message.Action {
	case SendMessageAction:
		client.sendMessage(message)
	case SendGameAction:
		client.sendGameAction(message)
	case CreateRoomAction:
		client.createAndNotifyRoom()
	case JoinRoomAction:
		var roomId string
		if err := json.Unmarshal(message.Data, &roomId); err != nil {
			logger.Printf("Invalid join-room payload: %v", err)
			return
		}
		client.joinRoom(roomId)
	case LeaveRoomAction:
		var roomId string
		if err := json.Unmarshal(message.Data, &roomId); err != nil {
			logger.Printf("Invalid leave-room payload: %v", err)
			return
		}
		client.leaveRoom(roomId)
	case StartGameAction:
		client.startGame(message)
	case AddBotAction:
		client.addBot(message)
	default:
		logger.Printf("Unknown action: %s", message.Action)
	}
}

func (client *Client) addBot(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.events <- &message
		logger.Println("Found room to add bot", room.ID)
	}
}

func (client *Client) sendMessage(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.events <- &message
		logger.Println("Found room to send message", room.ID, message.Action)
	}
}

func (client *Client) sendGameAction(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		room.events <- &message
	}
}

func (client *Client) createAndNotifyRoom() {
	room := client.createRoom()
	logger.Printf("Notifying client %s of created room %s", client.ID, room.ID)
	roomId, _ := json.Marshal(room.ID)
	message := &Message{
		Action: CreateRoomAction,
		Data:   json.RawMessage(roomId),
		Sender: client,
	}
	client.send <- message.encode()
}

func (client *Client) leaveRoom(roomId string) {
	room := client.hub.findRoomByID(roomId)
	if room == nil {
		return
	}
	logger.Println("Leaving room:", roomId)
	room.unregister <- client
}

func (client *Client) startGame(message Message) {
	roomId := message.Target.GetId()
	if room := client.hub.findRoomByID(roomId); room != nil {
		logger.Printf("Client %s is starting game in room %s", client.ID, roomId)
		room.startGame()
	}
}

func (client *Client) createRoom() *Room {
	roomId := uuid.NewString()            // Generate a new room ID
	room := client.hub.createRoom(roomId) // Create the room
	logger.Printf("Room created: %s by client %s", room.ID, client.ID)
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
	logger.Printf("Client: %s,%s joined room: %s", client.ID, client.Name, room.ID)
	roomId, _ := json.Marshal(room.ID)

	message := &Message{
		Action: JoinRoomAction,
		Data:   json.RawMessage(roomId),
		Target: room,
		Sender: client,
	}
	client.send <- message.encode()
}

func (client *Client) disconnect() {
	client.hub.unregister <- client // Unregister client from hub and rooms
	logger.Println("Disconnect run")
	for room := range client.rooms {
		select {
		case room.unregister <- client:
		default:
			// Room unregister channel is likely closed, or no receiver is ready.
			// This client will be cleaned up by the room's main loop if the room is still active.
			logger.Printf("Attempted to unregister client %s from room %s, but channel was closed or blocked.\n", client.ID, room.ID)
		}
	}
	// for room := range client.rooms {
	// 	room.unregister <- client
	// logger.Printf("%v", room.GetRoomSize())
	// if len(room.clients) == 0 {
	// 	client.hub.deleteRoom(room.ID)
	// 	// room.done <- true
	// 	close(room.done)
	// 	logger.Printf("Room %s has been deleted due to no active clients", room.ID)
	// }
	// logger.Printf("disconnect room logs %v", room.clients)
	// for c := range room.clients {
	// 	logger.Printf("room %s %s", c.Name, c.ID)
	// }

	// }
	// close(client.send)
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
				logger.Printf("error: %v", err)
			}
			break
		}
		messageContent = bytes.TrimSpace(bytes.ReplaceAll(messageContent, newline, space))
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
				logger.Printf("NextWriter error for %s: %v", client.ID, err)
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(client.send)
			for range n {
				w.Write(newline)
				w.Write(<-client.send)
			}
			if err := w.Close(); err != nil {
				logger.Printf("Writer close error for %s: %v", client.ID, err)
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Printf("Ping error for %s: %v", client.ID, err)
				return
			}
		}
	}
}
