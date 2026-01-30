package main

import (
	"log"
	"net/http"
	"os"

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

var (
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

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
	return "localhost:" + port
}
