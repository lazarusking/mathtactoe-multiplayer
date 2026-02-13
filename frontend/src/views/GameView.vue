<script setup lang="ts">
import Chat from '@/components/ChatUI.vue'
import FloatingChatButton from '@/components/FloatingChatButton.vue'
import GameWinModal from '@/components/GameWinModal.vue'
import WaitingScreen from '@/components/WaitingScreen.vue'
import type { GameState, GameStatePayload, GameStatus, Piece, Role, WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import { useWebSocket } from '@vueuse/core'
import {
  ArrowLeft,
  Ban,
  BookOpen,
  Calculator,
  CheckCircle,
  MousePointer2,
  Plus,
  RefreshCw,
  Share2,
  Trophy
} from 'lucide-vue-next'
import { computed, onBeforeMount, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

/**
 * CONSTANTS & REFS
 */
const route = useRoute()
const router = useRouter()
const loaders: string[] = Object.values(
  import.meta.glob('@assets/loaders/*.{png,jpg,jpeg,PNG,JPEG,gif,webp}', {
    eager: true,
    query: '?url',
    import: 'default'
  })
)
const randomGallery = loaders[Math.floor(Math.random() * loaders.length)]

const username = ref(localStorage.getItem('username') || '')
const chatMessages = ref<{ sender: string; text: string; id: string }[]>([])
const showChat = ref(false)

const options = () => Array.from({ length: 9 }, (_, i) => ({ id: i + 1, number: '-', owner: '' }))
const tictacGrid = ref<Piece[]>(options())

/**
 * REACTIVE STATE
 */
const gameState = reactive({
  isSelecting: false,
  selectedGrid: { id: 0, number: '-', owner: '' } as Piece,
  players: {} as GameState['players'],
  currentPlayer: '',
  activeTurn: false,
  gameStatus: {
    gameWon: false,
    gameOver: false,
    gameDraw: false,
    winner: null
  } as GameStatus,
  toastMsg: '',

  selectGrid(grid: Piece) {
    if (grid.number === '-') {
      this.isSelecting = true
      this.selectedGrid = grid
    }
  },

  placeNumber(button: Piece) {
    const currentPlayerId = WSState.clientID
    const currentPlayerInfo = this.players[currentPlayerId]
    if (!currentPlayerInfo) return

    if (WSState.isFirstMove) {
      if (this.selectedGrid.id === 5 && button.number.toString() === '5') {
        this.setToast("Can't start with 5 at this position😏")
        return
      }
      WSState.isFirstMove = false
    }

    if (this.isSelecting && isCurrentPlayer.value) {
      tictacGrid.value = tictacGrid.value.map((grid) => {
        if (grid.id == this.selectedGrid.id && grid.number === '-') {
          return { ...grid, number: button.number.toString(), owner: WSState.clientID }
        }
        return grid
      })

      this.isSelecting = false
      const message = JSON.stringify({
        action: 'send-game',
        data: { location: this.selectedGrid.id, number: button.number, playerID: WSState.clientID },
        sender: { id: WSState.clientID, name: username.value },
        target: { id: route.params.room as string }
      } as WSMessage)

      try {
        send(message)
      } catch (error) {
        console.error('Error sending game state', error)
        return
      }

      currentPlayerInfo.pieces = currentPlayerInfo.pieces.filter((item) => button.id !== item.id)
      this.selectedGrid = { id: 0, number: '-', owner: '' }
    }
  },

  setGameStatus(status: GameStatus) {
    this.gameStatus = status
  },

  setToast(msg: string) {
    this.toastMsg = msg
    setTimeout(() => (this.toastMsg = ''), 1000)
  }
})

/**
 * WEBSOCKET SETUP
 */
const { send, ws } = useWebSocket(websocket.url, {
  onMessage: (_, event) => handleMessage(event),
  immediate: true,
  autoReconnect: true
})

const WSState = reactive({
  role: 'spectator' as Role,
  clientID: '',
  playerCount: 0,
  totalCount: 0,
  isFirstMove: true
})

/**
 * COMPUTED PROPERTIES
 */
const playerIds = computed(() => Object.keys(gameState.players))
const currentPlayerId = computed(() => WSState.clientID)
const isCurrentPlayer = computed(
  () => playerIds.value.length >= 2 && WSState.clientID === gameState.currentPlayer
)
const isSpectator = computed(() => WSState.role === 'spectator')

// Logic to identify the two competitors for display
const p1Info = computed(() => {
  if (isSpectator.value) {
    const id = playerIds.value[0] || ''
    const name = gameState.players[id]?.client.name
    return { id, name: (Boolean(name) && name) || 'Player 1' }
  }
  return { id: WSState.clientID, name: 'You' }
})

const p2Info = computed(() => {
  if (isSpectator.value) {
    const id = playerIds.value[1] || ''
    const name = gameState.players[id]?.client.name
    return { id, name: (Boolean(name) && name) || 'Player 2' }
  }
  const oppId = playerIds.value.find((id) => id !== WSState.clientID) || ''
  return { id: oppId, name: gameState.players[oppId]?.client.name || 'Opponent' }
})

const playerNameLabel = computed(() => {
  if (isSpectator.value) {
    const name = gameState.players[gameState.currentPlayer]?.client.name || '...'
    return `Spectating: ${name}`
  }
  const oppName = p2Info.value.name
  return isCurrentPlayer.value ? `Your Turn` : `Waiting for ${oppName}`
})

const currentTurnLabel = computed(() => {
  if (gameState.gameStatus.gameOver) return 'Game Over'
  return 'Status'
})

const didIWin = computed(
  () => gameState.gameStatus.gameWon && gameState.gameStatus.winner === WSState.clientID
)
const winnerName = computed(() => {
  const winnerId = gameState.gameStatus.winner
  return winnerId ? (gameState.players[winnerId]?.client.name ?? 'Unknown') : null
})

/**
 * HANDLER FUNCTIONS
 */
function toggleChat() {
  showChat.value = !showChat.value
}

function shareRoom() {
  const data: ShareData = {
    url: window.location.href,
    text: "Join me let's play!",
    title: 'TicTacToe Math'
  }
  if (navigator.clipboard) {
    navigator.clipboard.writeText(`${data.url}\n${data.text}`)
    gameState.setToast('Link copied to clipboard!')
  }
  try {
    if (navigator.share && navigator.canShare?.(data)) {
      navigator.share(data)
      return
    }
  } catch (error) {
    console.error('Share failed', error)
  }
}

function playAgain() {
  send(
    JSON.stringify({
      action: 'start-game',
      sender: { id: WSState.clientID, name: username.value },
      target: { id: route.params.room as string }
    })
  )
}

function sendChatMessage(message: string) {
  send(
    JSON.stringify({
      action: 'send-message',
      data: message,
      sender: { id: WSState.clientID, name: username.value },
      target: { id: route.params.room as string }
    })
  )
}

function handleAddBot() {
  send(
    JSON.stringify({
      action: 'add-bot',
      target: { id: route.params.room as string }
    })
  )
}

function handleMessage(event: MessageEvent<string>) {
  const msgLines = event.data.split(/\r?\n/).filter((line) => line.trim() !== '')
  for (const line of msgLines) {
    try {
      const msg: WSMessage = JSON.parse(line)
      console.log(msg)

      switch (msg.action) {
        case 'start-game':
          tictacGrid.value = options()
          WSState.isFirstMove = true
          break
        case 'state-sync': {
          const payload = msg.data as GameStatePayload
          WSState.clientID = payload.self
          gameState.players = payload.game.players
          gameState.currentPlayer = payload.game.currentPlayer
          gameState.setGameStatus(payload.game.gameStatus)
          WSState.role = payload.role
          WSState.playerCount = payload.playerCount
          WSState.totalCount = payload.totalCount
          WSState.isFirstMove = payload.game.isFirstMove
          break
        }
        case 'update-game':
          gameState.currentPlayer = msg.data as string
          break
        case 'send-message':
          chatMessages.value.push({
            sender: msg.sender?.name ?? 'system',
            text: msg.data as string,
            id: msg.sender?.id ?? 'system'
          })
          break
        case 'send-game': {
          const payload = msg.data as { location: number; number: string; playerID: string }
          tictacGrid.value = tictacGrid.value.map((g) =>
            g.id === payload.location
              ? { ...g, number: payload.number.toString(), owner: payload.playerID }
              : g
          )
          WSState.isFirstMove = false
          break
        }
        case 'join-room':
          WSState.clientID = msg.sender?.id || ''
          break
        case 'game-status':
          gameState.setGameStatus(msg.data as GameStatus)
          break
      }
    } catch (e) {
      console.error('Failed to parse WS message', line, e)
    }
  }
}

const handleExit = () => router.push('/')

/**
 * LIFECYCLE & WATCHERS
 */
onBeforeMount(() => {
  if (route.params.room) {
    send(
      JSON.stringify({
        action: 'join-room',
        data: route.params.room,
        sender: { name: username.value }
      })
    )
  }
})

watch(
  [() => route.params.room, () => ws.value?.readyState],
  ([room, state], [oldRoom]) => {
    if (state !== WebSocket.OPEN || !room) return
    if (oldRoom && oldRoom !== room) {
      send(
        JSON.stringify({
          action: 'leave-room',
          data: oldRoom,
          sender: { id: WSState.clientID, name: username.value }
        })
      )
    }
    send(
      JSON.stringify({
        action: 'join-room',
        data: room,
        sender: { name: username.value }
      })
    )
  },
  { immediate: true, deep: true }
)
</script>

<template>
  <div class="h-dynamic-screen bg-background text-foreground font-sans">
    <!-- View 1: Lobby (Waiting for Opponent) -->
    <template v-if="playerIds.length < 2">
      <WaitingScreen :random-gallery="randomGallery" :room-code="route.params.room as string" @cancel="handleExit"
        @share="shareRoom" @add-bot="handleAddBot" />
    </template>

    <!-- View 2: Match In Progress -->
    <template v-else>
      <!-- Mobile Header Navigation -->
      <header
        class="sticky top-0 bg-transparent px-6 py-5 flex items-center justify-between lg:hidden backdrop-blur-md border-border">
        <h1 class="text-2xl font-display font-bold text-foreground tracking-wide">
          TicTac<span class="text-primary">Math</span>
        </h1>
        <div class="flex items-center gap-2">
          <button class="p-2 rounded-xl bg-muted text-muted-foreground active:scale-90 transition-transform"
            @click="shareRoom">
            <Share2 class="w-5 h-5" />
          </button>
        </div>
      </header>

      <div
        class="min-h-screen py-6 lg:py-10 px-4 md:px-8 max-w-400 mx-auto grid grid-cols-1 lg:grid-cols-12 gap-8 relative">
        <!-- Left Sidebar: Hidden on Mobile -->
        <div class="hidden lg:block lg:col-span-3 space-y-6">
          <div class="flex items-center space-x-3 mb-6">
            <div
              class="w-10 h-10 rounded-xl bg-linear-to-br from-secondary to-green-600 flex items-center justify-center shadow-lg">
              <Calculator class="w-5 h-5 text-white" />
            </div>
            <h1 class="text-2xl font-bold text-foreground tracking-wide">
              TicTac<span class="text-secondary">Math</span>
            </h1>
          </div>

          <div class="bg-card rounded-3xl p-6 shadow-xl border-border">
            <h2 class="text-xl font-bold text-foreground mb-6 flex items-center">
              <BookOpen class="w-5 h-5 mr-2 text-secondary" />
              How to Play
            </h2>
            <div class="space-y-6">
              <div class="flex items-start gap-3 group">
                <div
                  class="p-3 rounded-2xl bg-yellow-900/30 text-yellow-400 transition-transform group-hover:scale-110 shadow-sm">
                  <Trophy class="w-4 h-4" />
                </div>
                <div>
                  <h3 class="text-sm font-bold text-foreground mb-1">Objective</h3>
                  <p class="text-xs text-muted-foreground leading-relaxed">
                    Get 3 numbers in a line that sum exactly to 15.
                  </p>
                </div>
              </div>
              <div class="flex items-start gap-3 group">
                <div
                  class="p-3 rounded-2xl bg-blue-900/30 text-blue-400 transition-transform group-hover:scale-110 shadow-sm">
                  <MousePointer2 class="w-4 h-4" />
                </div>
                <div>
                  <h3 class="font-black text-sm tracking-tight text-foreground">Strategy</h3>
                  <p class="text-xs text-muted-foreground mt-1 leading-relaxed">
                    Select a grid slot first, then pick a number from your deck.
                  </p>
                </div>
              </div>
              <div class="flex items-start gap-3 group">
                <div
                  class="p-3 rounded-2xl bg-purple-100 text-purple-600 dark:bg-purple-900/30 dark:text-purple-400 transition-transform group-hover:scale-110 shadow-sm">
                  <Ban class="w-4 h-4" />
                </div>
                <div>
                  <h3 class="font-black text-foreground text-sm tracking-tight">Restriction</h3>
                  <p class="text-xs text-muted-foreground mt-1 leading-relaxed">
                    No starting with 5 in the center.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Main Board Section -->
        <div class="col-span-1 lg:col-span-6 flex flex-col items-center">
          <!-- Scoreboard -->
          <div
            class="w-full max-w-sm bg-muted/20 rounded-4xl p-6 mb-10 flex items-center justify-between border-inherit backdrop-blur-sm">
            <div class="flex flex-col items-center gap-2">
              <span class="text-[10px] font-bold text-slate-500 uppercase tracking-widest">You</span>
              <div :class="[
                'w-14 h-14 rounded-2xl flex items-center justify-center border-2 transition-all',
                gameState.currentPlayer === p1Info.id
                  ? 'border-primary bg-primary/10 shadow-[0_0_15px_rgba(45,238,121,0.2)]'
                  : 'border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800 opacity-50'
              ]">
                <span class="text-2xl font-bold text-primary">{{
                  p1Info.name === 'You' ? username.charAt(0) || '?' : p1Info.name.charAt(0)
                }}</span>
              </div>
            </div>

            <!-- Center: Turn indicator -->
            <div class="flex flex-col items-center flex-1">
              <span class="text-xs font-bold text-muted-foreground mb-2 uppercase tracking-tighter">{{ currentTurnLabel
              }}</span>
              <div :class="[
                'px-6 py-2 rounded-full text-xs font-bold shadow-inner truncate max-w-45',
                playerNameLabel === 'Your Turn' ? 'bg-primary/20' : 'bg-accent/10'
              ]">
                {{ playerNameLabel }}
              </div>
            </div>

            <div class="flex flex-col items-center gap-2">
              <span class="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">Bot</span>
              <div :class="[
                'w-14 h-14 rounded-2xl flex items-center justify-center border-2 transition-all',
                gameState.currentPlayer === p2Info.id
                  ? 'border-accent bg-accent/10 shadow-[0_0_15px_rgba(251,113,133,0.2)]'
                  : 'border-border bg-card opacity-50'
              ]">
                <span class="text-2xl font-bold text-accent">{{ p2Info.name.charAt(0) }}</span>
              </div>
            </div>
          </div>
          <!-- Spectator match results  -->
          <Transition name="fade">
            <div v-if="gameState.gameStatus.gameOver && isSpectator"
              class="mb-6 w-full max-w-sm animate-in zoom-in duration-500">
              <div class="bg-indigo-500/10 border border-indigo-500/30 rounded-3xl py-3 px-6 text-center">
                <span class="text-[10px] font-black uppercase tracking-[0.2em] text-indigo-300 block mb-1">Match
                  Results</span>
                <p class="text-lg font-bold text-foreground">
                  <template v-if="gameState.gameStatus.gameDraw">🤝 Match Drawn</template>
                  <template v-else>🏆 {{ winnerName }} Won!</template>
                </p>
              </div>
            </div>
          </Transition>
          <!-- 3x3 Skeletal Grid -->
          <div class="relative p-6 bg-muted/15 rounded-[4rem] shadow-2xl border-2 border-border backdrop-blur-md">
            <div class="grid grid-cols-3 gap-4 w-72 h-72 sm:w-80 sm:h-80 md:w-96 md:h-96">
              <button v-for="y in tictacGrid" :key="y.id" :disabled="!isCurrentPlayer || gameState.gameStatus.gameOver"
                :class="[
                  'group relative rounded-2xl flex items-center justify-center transition-all duration-300 outline-none border-2 shadow-sm',
                  y.number === '-'
                    ? gameState.selectedGrid.id === y.id
                      ? 'bg-primary/10 border-primary border-solid shadow-[inset_0_0_10px_rgba(45,238,121,0.2)] scale-[1.02]'
                      : 'bg-transparent border-dashed border-border hover:border-muted-foreground/80'
                    : 'bg-card border-border shadow-xl'
                ]" @click="gameState.selectGrid(y)">
                <template v-if="y.number !== '-'">
                  <div class="flex flex-col items-center scale-up">
                    <span :class="[
                      'text-3xl lg:text-6xl font-black',
                      y.owner === p1Info.id ? 'text-primary' : 'text-accent'
                    ]">
                      {{ y.number }}
                    </span>
                    <div :class="[
                      'mt-2 w-8 h-1 rounded-full opacity-60',
                      y.owner === p1Info.id ? 'bg-primary' : 'bg-accent'
                    ]"></div>
                  </div>
                </template>
                <template v-else>
                  <div class="flex items-center justify-center">
                    <CheckCircle v-if="gameState.selectedGrid.id === y.id"
                      class="w-6 h-6 lg:w-8 lg:h-8 text-primary/60" />
                    <Plus v-else
                      class="w-6 h-6 lg:w-8 lg:h-8 text-muted-foreground/20 group-hover:text-muted-foreground/60 transition-opacity" />
                  </div>
                </template>
              </button>
            </div>
          </div>

          <!-- Number Selection Deck -->
          <div class="mt-10 w-full max-w-lg"
            :class="{ 'opacity-50 pointer-events-none grayscale-[0.5]': !gameState.isSelecting }">
            <p class="text-center text-xs font-black text-muted-foreground mb-5 uppercase tracking-[0.3em]">
              {{
                isSpectator
                  ? 'Spectator Mode'
                  : gameState.isSelecting
                    ? 'Pick Your Number'
                    : 'Select a\
              Slot First'
              }}
            </p>
            <div class="flex flex-wrap justify-center gap-4">
              <button v-for="i in gameState.players[currentPlayerId]?.pieces || []" :key="i.id"
                class="w-14 h-14 rounded-2xl text-2xl font-black transition-all shadow-lg transform hover:scale-110 active:scale-90 bg-card border-2 border-border text-foreground hover:border-primary hover:text-primary"
                @click="gameState.placeNumber(i)">
                {{ i.number }}
              </button>
            </div>
          </div>
        </div>

        <!-- Right Sidebar: Hidden on Mobile -->
        <div class="hidden lg:block lg:col-span-3 space-y-6">
          <div class="bg-card rounded-3xl p-6 shadow-xl border border-border">
            <h2 class="text-xl font-black text-foreground mb-6">Controls</h2>
            <div class="space-y-4">
              <button
                class="w-full py-4 px-6 rounded-2xl bg-purple-700 text-foreground border-border font-bold focus:outline-none focus:ring-2 focus:ring-purple-500 hover:ring hover:ring-purple-600 transition-all flex items-center justify-between group"
                @click="shareRoom">
                <span>Share Room</span>
                <Share2 class="w-5 h-5 text-muted-foreground" />
              </button>
              <button
                class="w-full py-4 px-6 rounded-2xl bg-primary/80 text-white shadow-xl shadow-primary/20 font-bold hover:brightness-105 transition-all flex items-center justify-between group"
                @click="playAgain">
                <span>New Game</span>
                <RefreshCw
                  class="w-5 h-5 text-foreground/80 group-hover:rotate-180 transition-transform duration-500" />
              </button>
            </div>
          </div>

          <div :class="[{ hidden: !showChat }, 'flex flex-col h-100 overflow-hidden relative']">
            <Chat v-if="showChat" class="h-full rounded-2xl" :playerID="WSState.clientID" :messages="chatMessages"
              :users="WSState.totalCount" @send-message="sendChatMessage" @toggle="toggleChat" />
          </div>

          <button
            class="w-full py-4 text-muted-foreground hover:text-foreground transition-colors text-xs font-black uppercase tracking-[0.2em] flex items-center justify-center gap-2"
            @click="handleExit">
            <ArrowLeft class="w-4 h-4" />
            Leave game
          </button>
        </div>

        <GameWinModal v-if="gameState.gameStatus.gameOver && !isSpectator" :won="didIWin"
          :draw="gameState.gameStatus.gameDraw" :winner="winnerName" @play-again="playAgain" />
      </div>
    </template>

    <Chat v-if="showChat && playerIds.length < 2" class="fixed right-0 top-0 bottom-0 hidden md:flex w-80"
      :playerID="WSState.clientID" :messages="chatMessages" :users="WSState.totalCount" @send-message="sendChatMessage"
      @toggle="toggleChat" />

    <!-- mobile size -->
    <Chat v-if="showChat" class="h-dynamic-screen inset-0 md:hidden block absolute bg-background/90 z-50"
      :playerID="WSState.clientID" :messages="chatMessages" :users="WSState.totalCount" @send-message="sendChatMessage"
      @toggle="toggleChat" />

    <FloatingChatButton v-if="!showChat" @toggle="toggleChat" />

    <!-- Toast / Notification Layer -->
    <Transition>
      <div v-if="gameState.toastMsg"
        class="fixed top-8 left-1/2 -translate-x-1/2 z-100 px-5 py-2 bg-primary/80 text-foreground font-bold rounded-2xl shadow-2xl animate-in fade-in slide-in-from-top-4">
        {{ gameState.toastMsg }}
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.scale-up {
  animation: scale-up 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
}

@keyframes scale-up {
  from {
    transform: scale(0.8);
    opacity: 0;
  }

  to {
    transform: scale(1);
    opacity: 1;
  }
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
