export interface WSMessage {
  action: Action
  message: string
  target: Target
  sender: Sender
}

export interface Target {
  id: string
  name: string
}

export interface Sender {
  UserID: string
  id: string
  name: string
}
export interface Detail {
  id: number
  number: string
}
type Action =
  | 'update-game'
  | 'game-status'
  | 'start-game'
  | 'send-message'
  | 'join-room'
  | 'leave-room'
  | 'send-game'
