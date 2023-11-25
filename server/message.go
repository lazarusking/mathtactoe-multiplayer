// message.go
package main

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const SendGameAction = "send-game"
const StartGameAction = "start-game"
const UpdateStateAction = "update-game"
const GameStatusAction = "game-status"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Target  *Room   `json:"target"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	// log.Println(121212)
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
