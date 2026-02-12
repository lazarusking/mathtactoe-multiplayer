export interface WSMessage<T = unknown> {
  action: WSAction
  data: T
  target?: Target
  sender?: Sender
}

export interface Target {
  id: string
  name?: string
}

export interface Sender {
  id?: string
  isBot?: boolean
  name: string
}
export interface Piece {
  id: number
  number: string
  owner: string
}


export interface GameStatus {
  gameWon: boolean;
  gameOver: boolean;
  gameDraw: boolean;
  winner: string | null;
}
interface PlayerInfo {
  pieces: Piece[];
  client: Sender;
}
export interface GameState {
  currentPlayer: string; //initially 0|1 now uuid
  players: { [key: string]: PlayerInfo };
  options: Piece[][];
  gameStatus: GameStatus;
  isFirstMove: boolean;
}

export type Role = "player" | "spectator";

export interface GameStatePayload {
  game: GameState;
  role: Role;
  self: string;
  playerCount: number;
  totalCount: number;
}

export type WSAction =
  | 'join-room'
  | 'create-room'
  | 'start-game'
  | 'leave-room'
  | 'update-game'
  | 'send-game'
  | 'send-message'
  | 'game-status'

  | 'state-sync'
  | 'request-play'
  | 'leave-game'
  | 'add-bot'
