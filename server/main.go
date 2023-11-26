package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:   hub,
		rooms: make(map[*Room]bool),
		conn:  conn,
		send:  make(chan []byte, 256),
		// Add client-specific properties if needed
		ID:   uuid.New(),
		Name: "",
	}
	// Register the client by sending it to the register channel.
	client.hub.register <- client // Send the client to be registered
	// client.joinRoom("")
	go client.listen()
	go client.write()

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	hub := NewHub()                //decl of hub
	go hub.monitorClientsChannel() // Start a goroutine for receiving channel messages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})
	http.ListenAndServe(envPortOr("8080"), nil)
}
func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}

type Player struct {
	ID     uint8 `json:"id"`
	Number uint8 `json:"number"`
}
type GameStatus struct {
	GameWon  bool `json:"gameWon"`
	GameOver bool `json:"gameOver"`
	GameDraw bool `json:"gameDraw"`
}
type GameState struct {
	CurrentPlayer uint8                  `json:"currentPlayer"`
	Players       map[uuid.UUID][]Player `json:"players"`
	Options       [][]Player             `json:"options"`
	Clients       []*Client              `json:"clients"`
	GameStatus    GameStatus             `json:"gameStatus"`
}

func (game *GameState) switchPlayer() {
	x := (game.CurrentPlayer + 1) % 2
	game.CurrentPlayer = x
}

//	func (game *GameState) reset() {
//		game = &GameState{
//			CurrentPlayer: 0,
//			Players:       map[uuid.UUID][]Player{},
//			Options:       [][]Player{},
//			Clients:       []*Client{},
//			GameStatus: GameStatus{
//				GameWon:  false,
//				GameOver: false,
//				GameDraw: false,
//			},
//		}
//	}
func NewGame() *GameState {
	return &GameState{
		CurrentPlayer: uint8(rand.Intn(2)),
		Players:       map[uuid.UUID][]Player{},
		Options:       [][]Player{playerOne, playerTwo},
		Clients:       []*Client{},
		GameStatus: GameStatus{
			GameWon:  false,
			GameOver: false,
			GameDraw: false,
		},
	}
}

var playerOne = []Player{
	{ID: 1, Number: 1},
	{ID: 2, Number: 3},
	{ID: 3, Number: 5},
	{ID: 4, Number: 7},
	{ID: 5, Number: 9},
}
var playerTwo = []Player{
	{ID: 1, Number: 2},
	{ID: 2, Number: 4},
	{ID: 3, Number: 6},
	{ID: 4, Number: 8},
}

type Detail struct {
	Number string `json:"number"`
	ID     int    `json:"id"`
}

func calculateWinner(squares []Detail) (bool, bool) {
	lines := [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8},
		{2, 4, 6},
	}
	draw := true
	for i := 0; i < len(lines); i++ {
		a, b, c := lines[i][0], lines[i][1], lines[i][2]
		// log.Printf("%v:%v:%v", a, b, c)
		for _, v := range []int{a, b, c} {
			if squares[v].Number == "-" {
				draw = false
				break
			}
		}

		sum := parseInt(squares[a].Number) + parseInt(squares[b].Number) + parseInt(squares[c].Number)
		// fmt.Println("sum:", sum)
		if sum == 15 {
			fmt.Println(squares[a].Number, squares[b].Number, squares[c].Number)
			return true, false
		}
	}
	return false, draw
}
func parseInt(s string) int {
	if s == "-" {
		return math.MinInt
	}
	num, _ := strconv.Atoi(s)
	return num
}
