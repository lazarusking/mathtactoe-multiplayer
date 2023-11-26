package main

import (
	"encoding/json"
	"fmt"
	"log"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	clients       map[*Client]bool
	broadcast     chan *Message
	gamebroadcast chan *Message
	register      chan *Client
	unregister    chan *Client
	game          GameState
}

func NewRoom(id string) *Room {
	// if id == uuid.Nil {
	// 	id = uuid.New()
	// }
	return &Room{
		ID:            id,
		Name:          "",
		clients:       make(map[*Client]bool),
		broadcast:     make(chan *Message),
		gamebroadcast: make(chan *Message),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		game:          *NewGame(),
	}
}
func (room *Room) GetRoomSize() int {
	return len(room.clients)
}
func (room *Room) RunRoom() {
	for {
		select {
		case message := <-room.gamebroadcast:
			// game, _ := json.Marshal(room.game)
			log.Println("Game Broadcast")
			log.Println(len(room.game.Players), room.game.CurrentPlayer, room.GetRoomSize())

			for client := range room.clients {
				client.send <- message.encode()
			}
			log.Println("client id", message.Sender.ID)
			log.Println(room.getStatus())
			if !room.getStatus().GameOver {
				room.checkWinState(message)
			}

			room.game.switchPlayer()
			room.updateState()

		case message := <-room.broadcast:
			for client := range room.clients {
				client.send <- message.encode()
			}

		case client := <-room.register:
			log.Println(room.GetRoomSize(), len(room.game.Clients), len(room.game.Players))
			// if len(room.game.Clients) >= 2 || room.GetRoomSize() >= 2 {
			if len(room.game.Players) >= 2 || room.GetRoomSize() >= 2 {
				room.notifyFullRoom(client)
				break
			}
			// room.notifyClientJoined(client)
			room.clients[client] = true
			for client := range room.clients {
				if _, exists := room.clients[client]; exists {
					a := room.game.Options[room.game.CurrentPlayer]
					log.Println("current Player:", room.game.CurrentPlayer)
					room.game.Players[client.ID] = a
					room.game.Clients = append(room.game.Clients, client)
					room.game.switchPlayer()
				}
			}
			// log.Printf("Room Size:%d\n %v", room.GetRoomSize(), room.ID)
			// for k := range room.clients {
			// 	log.Println("client:", k.isInRoom(room), k.rooms, k.ID)
			// }
			if room.GetRoomSize() >= 2 || len(room.game.Players) >= 2 {
				room.startGame(client)
			}

		case client := <-room.unregister:
			// log.Printf("69:%v -%d", client.ID, room.GetRoomSize())

			delete(client.rooms, room)
			delete(room.clients, client)
			delete(room.game.Players, client.ID)
			// log.Printf("102: leaving: %v -%d", client.ID, room.game.Players[client.ID])
			room.setStatus(false, false, false)
			// room.game.reset()
			room.game = *NewGame()
			// room.updateState()
			room.startGame(client)
		}
	}
}
func (room *Room) startGame(client *Client) {
	game, _ := json.Marshal(room.game)
	log.Println("Game Start")
	log.Println(len(room.clients), len(room.game.Players), room.game.CurrentPlayer)
	message := &Message{
		Action:  StartGameAction,
		Message: string(game),
		Target:  room,
		Sender:  client}
	for client := range room.clients {
		client.send <- message.encode()
	}
}
func (room *Room) updateState() {
	if room.getStatus().GameOver {
		return
	}
	// var grids = struct {
	// 	Tictac []Detail `json:"tictac"`
	// 	Button Detail   `json:"button"`
	// }{
	// 	Tictac: []Detail{},
	// 	Button: Detail{},
	// }
	// log.Println(message.Message)
	// if err := json.Unmarshal([]byte(message.Message), &grids); err != nil {
	// 	log.Printf("Error on unmarshal message %s", err)
	// 	return
	// }
	// var newButtons []Player
	// for _, p := range room.game.Players[message.Sender.ID] {
	// 	if p.ID != uint8(grids.Button.ID) {
	// 		newButtons = append(newButtons, p)
	// 	}
	// }
	// room.game.Players[message.Sender.ID] = newButtons

	game, _ := json.Marshal(room.game)
	log.Println("Game Updated")
	log.Println(len(room.game.Players), room.game.CurrentPlayer, len(room.game.Clients))
	updatedStateMessage := &Message{
		Action:  UpdateStateAction,
		Message: string(game),
		Target:  room}
	for client := range room.clients {
		updatedStateMessage.Sender = client
		client.send <- updatedStateMessage.encode()
	}
}
func (room *Room) checkWinState(message *Message) {
	log.Println("Checking win")
	var grids []Detail
	data := message.Message
	// log.Println(data)
	if err := json.Unmarshal([]byte(data), &grids); err != nil {
		log.Printf("Error on unmarshal message %s", err)
		return
	}
	// log.Printf("%T:%v", data, data)
	log.Println(len(room.game.Players), room.game.CurrentPlayer, len(room.game.Clients))
	win, draw := calculateWinner(grids)
	if win {
		room.setStatus(true, true, false)
		game, _ := json.Marshal(room.getStatus())
		fmt.Println("Yesss won")
		message := &Message{
			Action:  GameStatusAction,
			Message: string(game),
			Target:  room,
			Sender:  message.Sender,
		}
		for client := range room.clients {
			client.send <- message.encode()
		}
	}
	if draw {
		room.setStatus(false, true, true)
		game, _ := json.Marshal(room.getStatus())
		fmt.Println("It's a draw!!")
		message := &Message{
			Action:  GameStatusAction,
			Message: string(game),
			Target:  room,
			Sender:  message.Sender,
		}
		for client := range room.clients {
			client.send <- message.encode()
		}
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.GetID()),
	}
	log.Println(client.ID, "Joined room", room.ID)
	for client := range room.clients {
		client.send <- message.encode()
	}
}
func (room *Room) notifyFullRoom(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf("Room filled to %d", room.GetRoomSize()),
	}
	log.Println(client.ID, "Full room", room.ID, room.GetRoomSize())
	// room.broadcastToClientsInRoom(message.encode())
	for client := range room.clients {
		client.send <- message.encode()
	}
}
func (room *Room) GetId() string {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}
func (room *Room) getStatus() GameStatus {
	return room.game.GameStatus
}
func (room *Room) setStatus(gameWon, gameOver, gameDraw bool) {
	room.game.GameStatus = GameStatus{GameWon: gameWon, GameOver: gameOver, GameDraw: gameDraw}
}
