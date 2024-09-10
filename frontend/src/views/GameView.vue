<script setup lang="ts">
import Chat from '@/components/ChatUI.vue'
import FloatingChatButton from '@/components/FloatingChatButton.vue'
import GameWinModal from '@/components/GameWinModal.vue'
import type { Detail, WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import { useWebSocket } from '@vueuse/core'
import { computed, onUnmounted, reactive, ref, watchEffect } from 'vue'
import { useRoute } from 'vue-router'

const loaders = Object.values(import.meta.glob('@assets/loaders/*.{png,jpg,jpeg,PNG,JPEG,gif}', { eager: true, as: 'url' }))

const randomGallery = loaders[Math.floor(Math.random() * loaders.length)]

const route = useRoute()
const username = ref('')
const chatMessages = ref<{ sender: string, text: string }[]>([])
const showChat = ref(false)

function toggleChat() {
    showChat.value = !showChat.value;
}
const storedUsername = localStorage.getItem('username') || ''
if (storedUsername) {
    username.value = storedUsername
}

function shareRoom() {
    const data: ShareData = { url: route.fullPath, text: "Join me let's play!", title: "TicTacToe Math" }
    let shareSuccess = false

    try {
        if (navigator.canShare && navigator.canShare(data) && navigator.share) {
            navigator.share(data)
            shareSuccess = true
        }
    } catch (error) {
        shareSuccess = false
    }

    if (!shareSuccess) {
        if (navigator.clipboard) {
            navigator.clipboard.writeText(`${data.url}\n${data.text}`)
        } else {
            console.log("Error copying to clipboard");
        }
    }
}

function playAgain() {
    // router.go(0)
    console.log("replaying game");

    const data = {
        action: 'start-game',
        message: "starting new game",
        sender: { ...WSState.data.sender, name: username.value },
        target: { ...WSState.data.target, id: route.params.room as string }

    };
    send(JSON.stringify(data));
}

const { data, send } = useWebSocket(websocket.url, {
    onMessage(ws, event) {
        handleMessage(event)
    }, immediate: true,
    autoReconnect: true
})

const options = () => {
    const arr = []
    for (let index = 1; index <= 9; index++) {
        arr.push({ id: index, number: '-' })
    }
    return arr
}
const tictacGrid = ref(options())


const gameState = reactive({
    isSelecting: false,
    selectedGrid: { id: 0, number: '-' },
    players: {} as { [key: string]: Detail[] },
    currentPlayer: 0,
    activeTurn: false,
    selectGrid(grid: Detail) {
        if (grid.number === '-') {
            this.isSelecting = true
            this.selectedGrid = grid
            console.log(grid)
            return
        }

        if (grid.number === '-' && this.isSelecting) {
            this.selectedGrid = grid
        }
    },
    placeNumber(button: Detail) {
        if (WSState.data.action === 'start-game') {
            console.log(button.number, gameState.selectedGrid);
            if (gameState.selectedGrid.id === 5 && button.number.toString() === '5') {
                gameState.setToast("Can't start with 5 at this position😏")
                return
            }
        }
        if (this.isSelecting && isCurrentPlayer.value) {
            tictacGrid.value = tictacGrid.value.map((grid) => {
                if (grid.id == this.selectedGrid.id && grid.number === '-') {
                    return { ...grid, number: button.number.toString() }
                }
                return grid
            })
            this.players[WSState.clientID] = this.players[WSState.clientID].filter(
                (item) => button.id !== item.id
            )
            console.log(WSState.data)
            this.isSelecting = false
            this.selectedGrid = { id: 0, number: '-' }
            const message = JSON.stringify({
                ...WSState.data,
                action: 'send-game',
                message: JSON.stringify(tictacGrid.value),
                // target: (WSState.data as WSMessage).target
                sender: { ...WSState.data.sender, name: username.value },
                target: { ...WSState.data.target, id: route.params.room as string }

            } as WSMessage)
            send(message)
        }
    },
    gameStatus: {
        gameWon: false,
        gameOver: false,
        gameDraw: false,
    },
    setGameStatus(gameStatus: { gameWon: boolean; gameOver: boolean; gameDraw: boolean }) {
        gameState.gameStatus = gameStatus
    },
    toastMsg: "",
    setToast(msg: string) {
        this.toastMsg = msg
        setTimeout(() => {
            this.toastMsg = ''
        }, 1000);
    }
})
const players = computed(() => (Object.keys(gameState.players)))


const WSState = reactive({
    IsConnected: false,
    data: data as unknown as WSMessage,
    tictacGrid: tictacGrid,
    clientID: '',
    clients: [] as unknown as { id: string; name: string }[],
    roomId: route.params.room
})
// const keyr = ref(0)
// const flag = ref(true)

// function randPlay() {
//     // let a = ref([1, 2, 3, 4, 5, 6, 7, 8, 9])
//     // let b = gameState.players[WSState.clientID]
//     // let flag = true
//     if (WSState.clients.length >= 2 && !calculateWinner(tictacGrid.value)) {
//         while (flag.value) {
//             const randomGrid = Math.floor(Math.random() * tictacGrid.value.length)
//             keyr.value = (keyr.value + 1) % 2
//             console.log(keyr.value)
//             const player = gameState.players[WSState.clients[keyr.value].id]
//             console.log(WSState.clients[keyr.value].id)

//             const randomBtn = Math.floor(Math.random() * player.length)
//             console.log(tictacGrid.value[randomGrid])
//             if (tictacGrid.value[randomGrid].number === '-' && player[randomBtn]) {
//                 console.log(randomBtn)

//                 selectGrid(tictacGrid.value[randomGrid])
//                 placeNumber(player[randomBtn])
//             }
//             // else {
//             //     continue
//             // }
//             if (calculateWinner(tictacGrid.value)) {
//                 console.log('Winner')
//                 flag.value = false
//                 gameState.gameStatus.gameWon = true
//                 gameState.gameStatus.gameOver = true
//                 break
//             }
//         }
//     }
// }


function sendChatMessage(message: string) {
    const chatMessage = {
        action: 'send-message',
        message: message,
        sender: { id: WSState.clientID, name: username.value },
        target: { id: route.params.room as string }
    };
    send(JSON.stringify(chatMessage));
}
watchEffect(() => {
    if (route.params.room) {
        //here doesn't add the sender name
        const data = { action: 'join-room', message: route.params.room, sender: { name: username.value } }
        console.log(WSState.data, "current state after joining room");
        // send({data})
        send(JSON.stringify({ ...data, action: 'join-room', message: route.params.room }))
    }
})

onUnmounted(() => {
    console.log("leaving room");
    // WSState.clients.splice(gameState.currentPlayer, 1)
    console.log(WSState.clients);

    send(JSON.stringify({ action: "leave-room", message: route.params.room }))
})



function handleMessage(event: MessageEvent<string>) {
    const data = event.data.split(/\r?\n/)
    // console.log(data);
    for (let i = 0; i < data.length; i++) {
        const WSMessage: WSMessage = JSON.parse(data[i])
        console.log(WSMessage);

        let message; // Move the declaration here
        switch (WSMessage.action) {
            case 'start-game': {
                message = JSON.parse(WSMessage.message) // Assign the value here
                WSState.data = WSMessage
                console.log(WSState.data.sender, "current sender state");

                tictacGrid.value = options()
                console.log(message)
                WSState.clients = message.clients
                gameState.players = message.players
                gameState.currentPlayer = message.currentPlayer
                gameState.setGameStatus(message.gameStatus)
                console.log(gameState.players)
                break
            }
            case 'update-game': {
                WSState.data = WSMessage
                message = JSON.parse(WSMessage.message) // Assign the value here
                // gameState.currentPlayer = message.currentPlayer
                gameState.currentPlayer = message
                console.log(message)
                console.log(gameState.players)

                break
            }
            case 'send-message': {
                // WSState.data = WSMessage
                console.log(WSMessage);
                chatMessages.value.push({ sender: WSMessage.sender?.name ?? "system", text: WSMessage.message })
                break
            }
            case 'send-game': {
                // WSState.data = WSMessage
                message = JSON.parse(WSMessage.message) // Assign the value here
                console.log(WSMessage);
                tictacGrid.value = message
                break
            }
            case 'join-room':
                WSState.data = WSMessage
                WSState.clientID = WSMessage.sender.id
                //   if (route.params.room != data.message) {
                //     console.log(route)
                //     console.log(data.message)
                //   }
                break
            case 'game-status': {
                WSState.data = WSMessage
                const message: { gameOver: boolean, gameWon: boolean, gameDraw: boolean } = JSON.parse(WSMessage.message)
                if (message.gameDraw) {
                    gameState.gameStatus.gameDraw = true
                    gameState.gameStatus.gameOver = message.gameOver
                    break
                }
                if (WSState.clientID === WSMessage.sender.id) {
                    gameState.gameStatus.gameOver = message.gameOver
                    gameState.gameStatus.gameWon = message.gameWon
                    // console.log("Juxtapose");

                } else {
                    gameState.gameStatus.gameWon = false
                    gameState.gameStatus.gameOver = true
                    // console.log("Reverse side");

                }
                console.log(WSMessage.sender);
                break
            }
            default:
                console.log(JSON.parse(WSMessage.message))
                message = JSON.parse(WSMessage.message) // Assign the value here
                WSState.clients = message.clients
                gameState.players = message.players
                break
        }

    }
}


const isCurrentPlayer = computed(() => {
    const currentClient = WSState.clientID
    //   console.log(
    //     currentClient,
    //     WSState.clients[gameState.currentPlayer].id,
    //     currentClient === WSState.clients[gameState.currentPlayer].id
    //   )
    if (players.value.length >= 2 && currentClient === players.value[gameState.currentPlayer]) {
        return true
    }
    return false
})


const isPlayingClasses = computed(() => {
    return {
        'blur-none opacity-75': gameState.isSelecting,
        'pointer-events-none': !isCurrentPlayer.value
    }
})

const isSpectator = computed(() => {
    return !players.value.includes(WSState.clientID);
});
const playerName = computed(() => {
    console.log(WSState.data.sender);
    const opponent = WSState.clients.find(client => client.id !== WSState.clientID)
    // If the current user is not a player, consider them a spectator
    if (isSpectator.value) {
        return `You are spectating ${opponent?.name ?? "a game"}`;
    }
    // if (opponent) {
    //     playerNames.opponentPlayer = opponent.name
    // }
    const currentUsername = username.value ? ` (${username.value})` : '';
    const opponentName = opponent?.name ? ` (${opponent.name})` : '';

    return isCurrentPlayer.value ? `Your Turn${currentUsername}` : `Opponent's Turn${opponentName}`;
})

</script>
<template>
    <!-- <div v-if="WSState.clients.length < 2" aria-modal="true"
        class="fixed inset-0 z-10 flex items-center justify-center text-gray-500 bg-black bg-opacity-50" role="dialog">
        <div class="max-w-md px-4 py-8 mx-auto space-y-2 text-center text-white rounded-lg sm:px-6 lg:px-8">
            <p class="text-3xl font-bold"> Waiting for an opponent... </p>
            <button name="share" @click="shareRoom"
                class="flex px-4 py-2 mx-auto text-white bg-purple-700 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-1 ">
                <svg class="w-5 h-6 mr-2" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z">
                    </path>
                </svg>
                <span>Share</span>
            </button>
        </div>

    </div> -->

    <Transition>
        <div v-if="gameState.toastMsg"
            class="transition ease-in-out fixed inset-0 z-10 flex items-center justify-center text-gray-500 bg-opacity-50">
            <div
                class="flex px-4 py-2 mx-auto text-white bg-rose-700 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-1  max-w-md space-y-2 text-center sm:px-6 lg:px-8">
                {{ gameState.toastMsg }}
            </div>
        </div>
    </Transition>
    <!-- <ChatComponent /> -->
    <!-- <main :class="{ 'blur': players.length < 2 }" -->
    <main class="relative container flex flex-col items-center justify-center h-screen px-4 py-12 mx-auto"
        :class="WSState.clients.length < 2 && 'inset-0 z-10 flex items-center justify-center text-gray-500 bg-opacity-50'">
        <h2 v-if="players.length >= 2" class="px-4 py-2 mb-4 text-2xl font-bold text-center text-white"
            :class="[isSpectator ? 'bg-gray-500' : (isCurrentPlayer ? 'bg-green-700' : 'bg-red-700')]">
            {{
                playerName
            }}</h2>
        <!-- <button @click="randPlay">Play</button> -->

        <h2 v-if="WSState.clients.length < 2" class="relative space-y-2 text-center mb-8 animate-pulse">
            <p class="text-2xl font-bold mb-2">Waiting for an opponent...</p>
            <button name="share" @click="shareRoom"
                class="flex px-4 py-2 mx-auto text-white bg-purple-700 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-1 ">
                <svg class="w-5 h-6 mr-2" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"
                    xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round"
                        d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z">
                    </path>
                </svg>
                <span>Share</span>
            </button>
            <!-- <div
                class="inline-block w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin">
            </div> -->

            <img :src="randomGallery" alt="Loading" className="mx-auto max-w-xs rounded-lg shadow-lg" />

        </h2>
        <template v-else>
            <!-- <div class=""> -->
            <section aria-label="tictac grid buttons"
                class="w-full max-w-md aspect-square p-1 bg-gradient-to-br from-blue-500 via-violet-500 to-pink-500 rounded-lg shadow-lg"
                :class="isPlayingClasses">
                <div class="grid grid-cols-3 w-full h-full p-0.5 rounded-lg overflow-hidden">
                    <button aria-label="tictac button" v-for="y in tictacGrid" :key="y.id" :id="y.id.toString()"
                        @click="gameState.selectGrid(y)" type="button" class="w-full h-full bg-gray-900 flex items-center justify-center text-4xl font-bold
                        focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-opacity-50 transition-all
                        duration-200 border-purple-500 border-4 border-gradient"
                        :class="{ 'bg-transparent ': gameState.selectedGrid.id === y.id }">
                        {{ y.number }}
                        <!-- <span :class="{ 'text-blue-400': y.number === 'X', 'text-purple-400': y.number === 'O' }">
                        </span> -->
                    </button>
                </div>
            </section>
            <section aria-label="game-buttons" class="mt-6 flex justify-center space-x-4">
                <button @click="gameState.placeNumber(i)" :value="i.id" v-for="i in gameState.players[WSState.clientID]"
                    :key="i.id" type="button"
                    class="button-gradient p-10 w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-xl font-bold text-white shadow-lg transition-all duration-200 hover:from-blue-600 hover:to-purple-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900">
                    {{ i.number }}
                </button>
            </section>
            <!-- <section aria-label="tictac grid buttons"
                class="grid w-full max-w-md grid-cols-3 grid-rows-3 gap-4 shadow-md h-3/5" :class="isPlayingClasses">
                <button aria-label="tictac button" @click="gameState.selectGrid(y)" type="button"
                    :class="{ 'ring ring-offset-2 ring-offset-slate-800 ring-blue-700': gameState.selectedGrid.id === y.id }"
                    class="grid items-center justify-center w-auto h-auto p-8 text-4xl font-black text-white transition-colors bg-gray-600 rounded-lg shadow-md place-content-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none hover:bg-slate-700"
                    v-for="y in tictacGrid" :key="y.id" :id="y.id.toString()">
                    {{ y.number }}
                </button>
            </section> -->
            <!-- <section aria-label="game-buttons"
                class="w-full max-w-md h-1/4 grid grid-flow-col grid-cols-[repeat(auto-fit,_minmax(0,_1fr))] gap-2 font-black text-3xl shadow-md rounded-sm p-2 m-1">
                <button @click="gameState.placeNumber(i)" :value="i.id" v-for="i in gameState.players[WSState.clientID]"
                    :key="i.id" type="button"
                    class="inline-flex items-center justify-center w-full h-auto max-w-md p-2 mt-8 text-2xl font-bold text-white bg-gray-700 rounded-md md:h-10 hover:bg-slate-800 md:p-10">
                    {{ i.number }}
                </button>
            </section> -->
            <!-- </div> -->
        </template>

        <GameWinModal :game-status="gameState.gameStatus" @play-again="playAgain" />
        <!-- <div v-if="showChat" class="lg:w-1/3 h-[400px] lg:h-auto">
            <div class="bg-gray-800 rounded-lg shadow-lg overflow-hidden h-full flex flex-col">
                <div class="flex items-center justify-between p-3 bg-gray-700">
                    <div class="flex items-center">
                        <MessageSquare class="w-5 h-5 text-blue-400 mr-2" />
                        <h2 class="text-lg font-bold text-gray-100">Chat</h2>
                    </div>
                    <button @click="toggleChat" class="lg:hidden text-gray-400 hover:text-white focus:outline-none">
                        <X class="w-5 h-5" />
                    </button>
                </div>
                <ChatComponent :player-name="username" :room-id="route.params.room as string" :messages="chatMessages"
                    @send-message="sendChatMessage" />
            </div>
        </div> -->

        <Chat class="fixed right-0 top-0 bottom-0 hidden md:flex w-80" v-if="showChat" :player-name="username"
            :messages="chatMessages" @send-message="sendChatMessage" @toggle="showChat = false"
            :users="WSState.clients.length" />

        <!-- mobile size -->
        <Chat class="inset-0 md:hidden block absolute bg-gray-900 bg-opacity-95 z-10" v-if="showChat"
            :player-name="username" :messages="chatMessages" @send-message="sendChatMessage" @toggle="showChat = false"
            :users="WSState.clients.length" />
        <!-- <div v-if="showChat" class="lg:hidden block absolute top-0 right-0 w-1/3 h-[400px] lg:h-auto">
            <div class="bg-gray-800 rounded-lg shadow-lg overflow-hidden h-full flex flex-col">
                <div class="flex items-center justify-between p-3 bg-gray-700">
                    <div class="flex items-center">
                        <MessageSquare class="w-5 h-5 text-blue-400 mr-2" />
                        <h2 class="text-lg font-bold text-gray-100">Chat</h2>
                    </div>
                    <button @click="toggleChat" class="lg:hidden text-gray-400 hover:text-white focus:outline-none">
                        <X class="w-5 h-5" />
                    </button>
                </div>
                <ChatComponent :player-name="username" :room-id="route.params.room as string" :messages="chatMessages"
                    @send-message="sendChatMessage" />
            </div>
        </div> -->
        <FloatingChatButton v-if="!showChat" @toggle="toggleChat" />

    </main>
    <!-- <Transition name="slide-up">
        <div v-if="showChat" class="absolute inset-0 bg-gray-900 bg-opacity-95 flex flex-col z-10">
            <div class="bg-gray-800 rounded-t-lg shadow-lg overflow-hidden flex-grow flex flex-col">
            </div>
            <div class="bg-gray-800 p-4 flex justify-between items-center">
                <div class="flex items-center">
                    <MessageSquare class="w-5 h-5 text-blue-400 mr-2" />
                    <h2 class="text-lg font-bold text-gray-100">Chat</h2>
                </div>
                <button @click="toggleChat" class="text-gray-400 hover:text-white focus:outline-none">
                    <X class="w-5 h-5" />
                </button>
            </div>
            <ChatComponent :player-name="username" :room-id="route.params.room as string" :messages="chatMessages"
                @send-message="sendChatMessage" />
        </div>
    </Transition> -->
</template>

<style scoped>
.border-gradient {
    border-image: linear-gradient(to bottom right, #3b82f6, #a855f7, #ec4899) 1;
    /* border-image: linear-gradient(to bottom right, theme( gradientColorStops.pink.500)) 1; */
}

.button-gradient {
    background-image: linear-gradient(60deg, #29323c 0%, #485563 100%);
}

.v-enter-from,
.v-leave-to {
    opacity: .2;
}

.slide-up-enter-active,
.slide-up-leave-active {
    transition: all 0.3s ease-out;
}

.slide-up-enter-from,
.slide-up-leave-to {
    transform: translateY(100%);
    opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>