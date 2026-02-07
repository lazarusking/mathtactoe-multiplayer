package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/google/uuid"
)

type Piece struct {
	ID     uint8 `json:"id"`
	Number uint8 `json:"number"`
}
type Detail struct {
	Number string `json:"number"` // for some frontend reasons this is a string
	ID     uint8  `json:"id"`
	Owner  string `json:"owner"`
}
type GameStatus struct {
	GameWon  bool       `json:"gameWon"`
	GameOver bool       `json:"gameOver"`
	GameDraw bool       `json:"gameDraw"`
	Winner   *uuid.UUID `json:"winner,omitempty"`
}
type PlayerInfo struct {
	Pieces []Piece `json:"pieces"`
	Client Client  `json:"client"`
}
type GameState struct {
	CurrentPlayer uuid.UUID                `json:"currentPlayer"`
	PlayerPieces  map[uuid.UUID]PlayerInfo `json:"players"`
	Options       [][]Piece                `json:"options"`
	// Clients       []*Client             `json:"clients"`
	GameStatus GameStatus `json:"gameStatus"`

	IsFirstMove bool `json:"isFirstMove"`
}

type GameStatePayload struct {
	Game        GameState `json:"game"`
	Role        Role      `json:"role"`
	Self        uuid.UUID `json:"self"`
	PlayerCount int       `json:"playerCount"`
	TotalCount  int       `json:"totalCount"`
}

//	func (game *GameState) switchPlayer() {
//		x := (game.CurrentPlayer + 1) % 2
//		game.CurrentPlayer = x
//	}
func shuffledOptions() [][]Piece {
	opts := [][]Piece{
		append([]Piece{}, playerOne...),
		append([]Piece{}, playerTwo...),
	}

	rand.Shuffle(len(opts), func(i, j int) {
		opts[i], opts[j] = opts[j], opts[i]
	})

	return opts
}

func NewGame() *GameState {
	return &GameState{
		// CurrentPlayer: uuid.New(),
		PlayerPieces: map[uuid.UUID]PlayerInfo{}, //exactly 2 players
		Options:      [][]Piece{playerOne, playerTwo},
		GameStatus: GameStatus{
			GameWon:  false,
			GameOver: false,
			GameDraw: false,
		},
		IsFirstMove: true,
	}
}

var playerOne = []Piece{
	{ID: 1, Number: 1},
	{ID: 2, Number: 3},
	{ID: 3, Number: 5},
	{ID: 4, Number: 7},
	{ID: 5, Number: 9},
}
var playerTwo = []Piece{
	{ID: 1, Number: 0},
	{ID: 2, Number: 2},
	{ID: 3, Number: 4},
	{ID: 4, Number: 6},
	{ID: 5, Number: 8},
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
	for i := range lines {
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
