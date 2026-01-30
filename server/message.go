// message.go
package main

import (
	"encoding/json"
)

type ActionType string // New custom type

const (
	SendMessageAction     ActionType = "send-message"
	SendGameAction        ActionType = "send-game"
	StartGameAction       ActionType = "start-game"
	UpdateStateAction     ActionType = "update-game"
	GameStatusAction      ActionType = "game-status"
	CreateRoomAction      ActionType = "create-room"
	JoinRoomAction        ActionType = "join-room"
	LeaveRoomAction       ActionType = "leave-room"
	UserJoinedAction      ActionType = "user-join"
	UserLeftAction        ActionType = "user-left"
	JoinRoomPrivateAction ActionType = "join-room-private"
	RoomJoinedAction      ActionType = "room-joined"
	GameStateSyncAction   ActionType = "state-sync"
)

type Message struct {
	Action ActionType      `json:"action"`
	Data   json.RawMessage `json:"data,omitempty"`
	Target *Room           `json:"target,omitempty"`
	Sender *Client         `json:"sender,omitempty"`
}

func (message *Message) encode() []byte {
	// log.Println(121212)
	json, err := json.Marshal(message)
	if err != nil {
		logger.Printf("Error marshaling message: %s %q-%s", message.Action, message.Data, err)
		logger.Panic(err)
	}

	return json
}
