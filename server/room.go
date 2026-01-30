package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room"

type RoomState int

const (
	Waiting RoomState = iota
	InProgress
	Finished
)

type Role string

const (
	RolePlayer    Role = "player"
	RoleSpectator Role = "spectator"
)

func (room *Room) getRole(client *Client) Role {
	if _, ok := room.players[client.ID]; ok {
		return RolePlayer
	}
	return RoleSpectator
}
func (room *Room) buildGameStatePayload(client *Client) GameStatePayload {
	return GameStatePayload{
		Game:        room.game,
		Role:        room.getRole(client),
		Self:        client.ID,
		PlayerCount: len(room.players),
		TotalCount:  len(room.clients),
	}
}

type BotEvent struct {
	Reason string
}
type Room struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	clients        map[*Client]bool
	spectators     map[uuid.UUID]*Client
	players        map[uuid.UUID]*Client //exactly 2 players
	lastBoardState []Detail
	events         chan *Message
	register       chan *Client
	unregister     chan *Client
	game           GameState
	state          RoomState

	enableBot bool
	botStop   map[uuid.UUID]chan struct{}
	botEvents chan BotEvent
	done      chan bool // Stop signal channel
}

func (room *Room) getBot() *Client {
	for _, c := range room.players {
		if c != nil && c.IsBot {
			return c
		}
	}
	return nil
}
func (room *Room) stopBotRun(botID uuid.UUID) {
	if ch, ok := room.botStop[botID]; ok {
		// close the channel to signal the goroutine to exit
		close(ch)
		delete(room.botStop, botID)
		logger.Printf("Stopped bot runner for %s in room %s", botID, room.ID)
	}
	// close(room.botEvents)
	// room.botEvents = make(chan BotEvent, 8) // reset for next game or it panics because it got closed

}

func (room *Room) runBot(bot *Client) {
	if bot == nil {
		return
	}
	if _, ok := room.botStop[bot.ID]; ok {
		// already running
		return
	}
	stop := make(chan struct{})
	room.botStop[bot.ID] = stop
	go func() {
		logger.Println("Bot listener started", bot.ID)
		defer logger.Println("Bot stopped", bot.ID)

		playedPieces := append([]Piece(nil), room.game.PlayerPieces[bot.ID].Pieces...) // copy
		//here im trying to track what the frontends bot might be playing
		//but it messes up on new game start since this part can't reinitialize itself
		//frontend owns the buttons so yh this is a workaround
		for {
			select {
			case <-room.done:
				logger.Println("Bot stopped (room closed)", room.ID)
				close(room.botEvents)
				return
			case <-stop:
				logger.Printf("Bot runner stopped for %s in room %s", bot.ID, room.ID)
				return
			case event := <-room.botEvents:
				switch event.Reason {
				case "start-game":
					// reset played pieces on new game
					playedPieces = append([]Piece(nil), room.game.PlayerPieces[bot.ID].Pieces...)
					logger.Printf("Bot %s reset pieces for new game: %+v", bot.ID, playedPieces)
				default:
					if room.state != InProgress {
						continue
					}
					logger.Println(event)
					if room.game.CurrentPlayer != bot.ID {
						continue
					}

					// short human-like delay
					time.Sleep(600 * time.Millisecond)

					move, ok := randomBotMove(
						room.lastBoardState,
						playedPieces)
					//todo: fix stale pieces from previous match
					logger.Printf("playedPieces: %+v, move: %+v", playedPieces,
						move)

					logger.Println("Bot move generated:", move, "OK:", ok)

					if !ok {
						continue
					}

					newBoard := make([]Detail, len(room.lastBoardState))
					copy(newBoard, room.lastBoardState)
					newBoard[move.location-1] = Detail{Number: strconv.Itoa(int(move.move)), ID: uint8(move.location)}
					var updatedPieces []Piece
					for _, p := range playedPieces {
						if strconv.Itoa(int(p.Number)) != strconv.Itoa(int(move.move)) {
							// updatedPieces = append(updatedPieces, p)
							updatedPieces = append(updatedPieces, Piece{ID: p.ID, Number: p.Number})
						}
					}
					playedPieces = updatedPieces
					payload, _ := json.Marshal(newBoard)

					room.events <- &Message{
						Action: SendGameAction,
						Data:   payload,
						Sender: bot,
						Target: room,
					}
				}
			}

		}
	}()
}

func randomBotMove(
	board []Detail,
	pieces []Piece,
) (struct {
	move     uint
	location uint
}, bool) {

	var empty []uint8
	logger.Println(board)
	for i, cell := range board {
		if cell.Number == "-" || cell.Number == "" {
			empty = append(empty, uint8(i+1))
		}
	}
	if len(empty) == 0 || len(pieces) == 0 {
		// if len(empty) == 0 {
		return struct {
			move     uint
			location uint
		}{}, false
	}

	cell := empty[rand.IntN(len(empty))]

	logger.Printf("empty: %v", empty)
	logger.Printf("Cell:%v Pieces:%v", cell, pieces)
	piecesIndex := rand.IntN(len(pieces))
	piece := pieces[piecesIndex]
	// return Detail{
	// 	ID:     cell,
	// 	Number: strconv.Itoa(int(piece.Number)),
	// }, true
	return struct {
		move     uint
		location uint
	}{
		move:     uint(piece.Number),
		location: uint(cell),
	}, true
}

func NewRoom(id string) *Room {
	// if id == uuid.Nil {
	// 	id = uuid.New()
	// }
	board := make([]Detail, 9)
	for i := range board {
		board[i] = Detail{Number: "-", ID: uint8(i)}
	}

	return &Room{
		ID:             id,
		Name:           "",
		clients:        make(map[*Client]bool),
		players:        make(map[uuid.UUID]*Client),
		lastBoardState: board,
		events:         make(chan *Message),
		spectators:     make(map[uuid.UUID]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		game:           *NewGame(),
		done:           make(chan bool),

		enableBot: os.Getenv("ENABLE_BOT") != "true",
		botStop:   make(map[uuid.UUID]chan struct{}),
		botEvents: make(chan BotEvent),
		state:     Waiting,
	}
}
func (room *Room) GetRoomSize() int {
	return len(room.clients)
}
func (room *Room) RunRoom() {
	// bot := NewBotClient("BOT")
	//enablebot flag that only is true on dev
	// enableBot := os.Getenv("ENABLE_BOT") != "true"
	logger.Printf("Room %s started. Enable bot: %v", room.ID, room.enableBot)
	for {
		select {
		case <-room.done: // Stop the room gracefully
			logger.Printf("Stopping room %s", room.ID)
			close(room.register)
			close(room.unregister)
			// close(room.gamebroadcast)
			return // Exit the loop

		case message := <-room.events:
			room.handleMessage(message)
			// if message.Action == SendGameAction {
			// 	var board []Detail
			// 	if err := json.Unmarshal(message.Data, &board); err != nil {
			// 		logger.Printf("Error on saving board state %s", err)
			// 		return
			// 	}
			// 	room.lastBoardState = board
			// }

			// logger.Println("Game Broadcast")

			// logger.Println("client id", message.Sender.ID)
			// logger.Println(room.getStatus())
			// room.checkWinState(message)
			// if !room.getStatus().GameOver {

			// 	room.switchTurn()
			// }
			// remember to uncomment
			// room.game.switchPlayer()
			// if room.getStatus().GameOver {
			// 	return
			// }
			// room.broadcastTurnUpdate()

		// case message := <-room.broadcast:
		// 	for client := range room.clients {
		// 		client.send <- message.encode()

		// 	}
		case client := <-room.register:
			logger.Println("Client Register")

			room.clients[client] = true
			if len(room.players) < 2 {
				room.players[client.ID] = client
			} else {
				room.spectators[client.ID] = client
			}
			room.notifyClientJoined(client)
			// room.broadcastGameStateToClients()

			logger.Println("Player count:", len(room.players))

			// }
			//#bot auto-add bot if only one human
			if room.enableBot && len(room.players) == 1 && room.getBot() == nil {
				logger.Println("Scheduling bot registration for room", room.ID)
				go func() {
					bot := NewBotClient("BOT")
					// use same registration flow so notifyClientJoined and broadcast happen
					room.register <- bot
				}()
			}
			for user := range room.clients {
				fmt.Println(user.ID, user.Name)
			}
			// if len(room.players) >= 2 || room.GetRoomSize() >= 2 {
			// 	room.notifyFullRoom(client)
			// }
			if len(room.players) == 2 && room.state == Waiting {
				room.startGame()
			}
			room.broadcastGameStateToClients()

			//log the player pieces
			for _, p := range room.game.PlayerPieces {
				logger.Printf("Player %s has pieces: %v", p.Client.Name, p.Client.ID)
			}
			// 	// todo: have to allow multiple clients in for spectating
			logger.Printf("%v room size.....player count %v ,%v", room.GetRoomSize(), len(room.game.PlayerPieces), room.players)

		case client := <-room.unregister:
			// logger.Printf("69:%v -%d", client.ID, room.GetRoomSize())
			logger.Println("Client Unregister")
			delete(room.clients, client)
			logger.Println(room.clients[client], "does it exist?", "is it a bot?", client.IsBot)
			delete(room.spectators, client.ID)
			delete(room.players, client.ID)
			if _, isPlayer := room.players[client.ID]; isPlayer {
				delete(room.players, client.ID)
				delete(room.game.PlayerPieces, client.ID)
				// Reset game if players < 2

			}

			//#bot auto-remove bot if only both is the player
			if room.enableBot && len(room.players) == 1 && room.getBot() != nil {
				logger.Println("Scheduling bot registration for removal", room.ID)
				// this makes sense than whatever I was doing earlier
				go func() {
					bot := room.getBot()
					if bot != nil {
						room.unregister <- bot
						room.stopBotRun(bot.ID)
					}
				}()
			}
			//
			//this removes bot on any client disconnect no?
			//also seems the bot really is the last man standing since
			//client gets removed by the disconnect(), but bot gets injected inside the room
			// so doesn't follow the conventional ws client cycle tsw
			// for c := range room.clients {
			// 	if c.IsBot {
			// 		logger.Println("Removing bot from room", room.ID)
			// 		room.stopBotRun(c.ID)
			// 		delete(room.players, c.ID)
			// 		delete(room.clients, c)
			// 		delete(room.game.PlayerPieces, c.ID)
			// 		close(c.send)
			// 	}
			// }

			logger.Print("Connected Clients")
			for c := range room.clients {
				logger.Printf("%s---%s", c.ID, c.Name)
			}
			logger.Print("Connected players")
			for _, p := range room.players {
				logger.Printf("Player isBot:%v --> ID:%s --> %s", p.IsBot, p.ID, p.Name)
			}
			if len(room.players) < 2 {
				room.game = *NewGame()
				// room.assignPlayerPieces()
				room.state = Waiting
			}
			room.broadcastGameStateToClients()
		}
	}
}

// Give a player their number pieces
func (room *Room) assignPlayerPieces() {
	// if _, exists := room.game.PlayerPieces[client.ID]; !exists {
	// 	a := room.game.Options[len(room.game.PlayerPieces)]
	// 	room.game.PlayerPieces[client.ID] = PlayerInfo{
	// 		Pieces: a,
	// 		Client: *client,
	// 	}
	// 	logger.Printf("Assigned pieces to player %s: %v", client.ID, a)
	// }
	ids := make([]uuid.UUID, 0, len(room.players))
	for id := range room.players {
		ids = append(ids, id)
	}
	options := shuffledOptions()

	for i, id := range ids {
		client := room.players[id]
		room.game.PlayerPieces[id] = PlayerInfo{
			Pieces: options[i],
			Client: *client,
		}
	}
	room.game.CurrentPlayer = ids[0]

}

func (room *Room) startGame() {
	if room.state == InProgress {
		return
	}
	room.state = InProgress
	// room.setStatus(false, false, false)
	room.game = *NewGame()
	// game, _ := json.Marshal(room.game)
	logger.Println("Game Start")
	// logger.Println(len(room.clients), len(room.game.PlayerPieces), room.game.CurrentPlayer)
	message := &Message{
		Action: StartGameAction,
		Data:   nil, //game,
		Target: room,
	}
	for client := range room.clients {
		client.send <- message.encode()
	}

	//simplify emitting broadcasttoclients updates the state for the game to start
	room.game.PlayerPieces = make(map[uuid.UUID]PlayerInfo)

	// Collect players deterministically
	// ids := make([]uuid.UUID, 0, len(room.players))
	// for id := range room.players {
	// 	ids = append(ids, id)
	// }

	// // Randomize both turn + pieces
	// rand.Shuffle(len(ids), func(i, j int) {
	// 	ids[i], ids[j] = ids[j], ids[i]
	// })

	// options := shuffledOptions()

	// for i, id := range ids {
	// 	client := room.players[id]
	// 	room.game.PlayerPieces[id] = PlayerInfo{
	// 		Pieces: options[i],
	// 		Client: *client,
	// 	}
	// }

	// // ✅ Valid turn (always a real player)
	// room.game.CurrentPlayer = ids[0]
	room.assignPlayerPieces()
	logger.Printf("Current player is %s", room.game.CurrentPlayer)
	//#bot run bot
	if bot := room.getBot(); bot != nil {
		room.runBot(bot)

		room.botEvents <- BotEvent{Reason: "start-game"}

	}
	if room.game.CurrentPlayer == uuid.Nil {
		logger.Panic("Current player is nil after startGame")
	}
	room.broadcastGameStateToClients()
}
func (room *Room) broadcastTurnUpdate() {
	logger.Println("Game Updated")
	currentPlayer, _ := json.Marshal(room.game.CurrentPlayer)
	updatedStateMessage := &Message{
		Action: UpdateStateAction,
		Data:   json.RawMessage(currentPlayer),
		Target: room}

	room.broadcast(updatedStateMessage)
}
func (room *Room) checkWinState(message *Message) {
	logger.Println("Checking win")
	var grids []Detail
	data := message.Data
	// logger.Println(data)
	if err := json.Unmarshal(data, &grids); err != nil {
		logger.Printf("Error on unmarshal message %s", err)
		return
	}
	// logger.Printf("%T:%v", data, data)
	// logger.Println(len(room.game.PlayerPieces), room.game.CurrentPlayer)
	win, draw := calculateWinner(grids)
	var winner *uuid.UUID
	if win {
		winner = &message.Sender.ID
	}
	if win || draw {
		room.setStatus(win, true, draw, winner)
		game, _ := json.Marshal(room.getStatus())

		message := &Message{
			Action: GameStatusAction,
			Data:   json.RawMessage(game),
			Target: room,
			Sender: message.Sender,
		}
		logger.Println("Game Won by:", message.Sender.Name)

		for client := range room.clients {
			client.send <- message.encode()
		}
	}

}

func (room *Room) handleMessage(message *Message) {
	switch message.Action {

	case SendGameAction:
		var board []Detail
		if err := json.Unmarshal(message.Data, &board); err != nil {
			logger.Printf("Error on saving board state %s", err)
			return
		}
		room.lastBoardState = board
		// room.broadcast(message)

		room.checkWinState(message)
		if !room.getStatus().GameOver {
			room.switchTurn()

			room.broadcastTurnUpdate()
		}

		room.broadcast(message)

	case StartGameAction:
		room.startGame()

	default:
		logger.Printf("Unknown action: %s", message.Action)
		room.broadcast(message)

	}
}
func (room *Room) broadcastGameStateToClients() {
	// Marshal the current game state
	for client := range room.clients {

		// gameState, err := json.Marshal(room.game)
		payload := room.buildGameStatePayload(client)
		data, err := json.Marshal(payload)

		if err != nil {
			logger.Printf("Error marshalling game state: %s", err)
			return
			//continue
		}

		// Create a message to send to all clients
		message := &Message{
			Action: GameStateSyncAction, // You can create a specific action if necessary
			Data:   json.RawMessage(data),
			Target: room,
			Sender: client,
		}

		// Send the game state to all clients
		client.send <- message.encode()
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	if client == nil {
		logger.Println("Error: client is nil in notifyClientJoined")
		return
	}
	clientName := client.Name
	if clientName == "" {
		clientName = client.GetID()
	}
	// payload := room.buildGameStatePayload(client)
	// data, err := json.Marshal(payload)
	// if err != nil {
	// 	logger.Printf("Error marshalling game state: %s", err)
	// 	return
	// 	//continue
	// }
	a, _ := json.Marshal(fmt.Sprintf(welcomeMessage, clientName))
	message := &Message{
		Action: SendMessageAction,
		Target: room,
		Data:   json.RawMessage(a),
	}
	// this is so i can keep the new joined client up to sync with the game
	//without sending the message to every other client
	// clientStateUpdate := &Message{
	// 	Action: GameStateSyncAction,
	// 	Target: room,
	// 	Data:   json.RawMessage(data),
	// 	Sender: client,
	// }
	// client.send <- clientStateUpdate.encode()
	for client := range room.clients {
		client.send <- message.encode()
	}
	// room.broadcast(message)
}
func (room *Room) notifyFullRoom(client *Client) {
	data, _ := json.Marshal(fmt.Sprintf("Room filled to %d", room.GetRoomSize()))
	message := &Message{
		Action: SendMessageAction,
		Target: room,
		Data:   json.RawMessage(data),
	}
	logger.Println(client.ID, "Full room", room.ID, room.GetRoomSize())
	// room.broadcast(message)
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
func (room *Room) setStatus(gameWon, gameOver, gameDraw bool, winner *uuid.UUID) {
	room.game.GameStatus = GameStatus{
		gameWon,
		gameOver,
		gameDraw, winner}
	if gameOver {
		room.state = Finished
		// room.stopBotRun(room.getBot().ID)

	}
}

func (room *Room) switchTurn() {
	for id := range room.game.PlayerPieces {
		if id != room.game.CurrentPlayer {
			room.game.CurrentPlayer = id
			// return
			break // break so it flows to botevent
		}
	}
	if bot := room.getBot(); bot != nil {
		room.botEvents <- BotEvent{Reason: "turn-switch"}
	}

}

func (room *Room) broadcast(msg *Message) {
	encoded := msg.encode()
	for client := range room.clients {
		client.send <- encoded
	}
}
