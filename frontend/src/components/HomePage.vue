<script setup lang="ts">
import type { WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import router from '@/router'
import { useWebSocket } from '@vueuse/core'
import { InfoIcon } from 'lucide-vue-next'
import { defineAsyncComponent, onMounted, ref, watch, watchEffect } from 'vue'
const HelpModal = defineAsyncComponent(() => import('./modal/HelpModal.vue'))

const roomId = ref('')
const playerName = ref('')
const showHelp = ref(false)
const activeTab = ref('create')

onMounted(() => {
  const storedUsername = localStorage.getItem('username')
  if (storedUsername) {
    playerName.value = storedUsername
  }
})
const { data, send } = useWebSocket(websocket.url, {
  onMessage(ws, event) {
    handleMessage(event)
  }, autoReconnect: true, immediate: true
})

watch(playerName, () => {
  localStorage.setItem('username', playerName.value)
})
watchEffect(() => {
  console.log(data)
})

function createRoom() {
  // const data = { action: 'join-room', message: null, sender: { name: playerName.value } }
  const data = { action: 'create-room', message: null, sender: { name: playerName.value } }
  console.log(data, "createRoom data");

  send(JSON.stringify(data))
  //clicking this sends a msg to ws, ws then joins room and sends room id created
  //client then leaves the room because the handleMessage navigates to a new page
  //new page GameView calls another join-room,creates another new client for the same client
}

function handleMessage(event: MessageEvent) {
  const data: WSMessage = JSON.parse(event.data)
  switch (data.action) {
    case 'join-room':
      console.log(data)
      console.log("Did this run")
      router.push({ name: "room", params: { room: data.message } })
      // router.push(data.message)
      break
    case 'create-room':
      console.log(data)
      console.log("Did this run")
      router.push({ name: "room", params: { room: data.message } })
      // router.push(data.message)
      break
    case 'send-message': {
      // WSState.data = WSMessage
      console.log(data)
      router.push({ name: "room", params: { room: data.message } })
      // console.log(WSMessage);
      break
    }
    default:
      break
  }
}

function closeModal() {
  showHelp.value = false
}

function joinGame() {
  if (roomId.value) {
    // const data = { action: 'join-room', message: null, sender: { name: playerName.value } }
    router.push({ name: "room", params: { username: 'eduardo', room: roomId.value } })
  }
}

const showLeaderboard = () => {
  console.log('Showing leaderboard')
}

const showOnlinePlayers = () => {
  console.log('Showing online players')
}

const showSettings = () => {
  console.log('Showing settings')
}

</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 to-gray-800 text-gray-100 flex flex-col">
    <!-- <div @click="showHelp = !showHelp" class="flex mx-auto container max-w-screen-md py-1 sm:py-2 sm:px-6 lg:px-8">
      <div class="ml-auto flex flex-wrap items-center justify-between">
        <button type="button"
          class="text-white items-center text-lg p-3 font-semibold rounded hover:bg-blue-500 hover:text-white">
          <InfoIcon class="h-6 w-6" />
          <span class="sr-only">How to play?</span>
        </button>
      </div>
    </div> -->

    <HelpModal :show-help="showHelp" @close-modal="closeModal" />

    <div class="flex-grow flex items-center justify-center p-4">
      <div class="w-full max-w-md">
        <div class="flex justify-center items-center space-x-2 mb-6">
          <h1
            class="text-4xl font-bold text-center bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-green-400">
            TicTacToe Math
          </h1>
          <button @click="showHelp = true"
            class="text-gray-400 hover:text-gray-200 focus:outline-none transition duration-200"
            aria-label="Information">
            <InfoIcon class="h-6 w-6" />
            <span class="sr-only">How to play?</span>
          </button>
        </div>
        <div class="text-center mb-6">
          <span class="bg-gray-700 text-xs font-semibold px-2.5 py-0.5 rounded-full">Multiplayer</span>
        </div>

        <div class="bg-gray-800 border border-gray-700 rounded-lg shadow-lg p-6">
          <div class="mb-6">
            <div class="flex border-b border-gray-700">
              <button @click="activeTab = 'join'"
                :class="['flex-1 py-2 px-4 focus:outline-none transition-colors duration-200',
                  activeTab === 'join' ? 'border-b-2 border-blue-500 text-blue-400' : 'text-gray-400 hover:text-gray-200']">
                Join Game
              </button>
              <button @click="activeTab = 'create'"
                :class="['flex-1 py-2 px-4 focus:outline-none transition-colors duration-200',
                  activeTab === 'create' ? 'border-b-2 border-green-500 text-green-400' : 'text-gray-400 hover:text-gray-200']">
                Create Game
              </button>
            </div>
          </div>

          <div v-if="activeTab === 'join'">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-200 mb-2" for="room-id">
                Room ID
              </label>
              <input v-model="roomId" required
                class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-200"
                id="room-id" placeholder="Enter Room ID" type="text" />
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-200 mb-2" for="player-name">
                Your Name
              </label>
              <input v-model="playerName" required
                class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 transition duration-200"
                id="player-name" placeholder="Enter your name" type="text" />
            </div>
            <button @click="joinGame"
              class="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline transition duration-200 ease-in-out">
              Join Game
            </button>
          </div>

          <form v-else @submit.prevent="createRoom">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-200 mb-2" for="create-player-name">
                Your Name
              </label>
              <input v-model="playerName" required
                class="block w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500 transition duration-200"
                id="create-player-name" placeholder="Enter your name" type="text" />
            </div>
            <button type="submit"
              class="w-full bg-green-600 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline transition duration-200 ease-in-out">
              Create New Game
            </button>
          </form>
        </div>

        <!-- <div class="mt-8 flex justify-center space-x-4">
          <button @click="showLeaderboard"
            class="text-gray-400 hover:text-gray-200 focus:outline-none transition duration-200"
            aria-label="Leaderboard">
            <Trophy class="h-6 w-6" />
          </button>
          <button @click="showOnlinePlayers"
            class="text-gray-400 hover:text-gray-200 focus:outline-none transition duration-200"
            aria-label="Players Online">
            <Users class="h-6 w-6" />
          </button>
          <button @click="showSettings"
            class="text-gray-400 hover:text-gray-200 focus:outline-none transition duration-200" aria-label="Settings">
            <Settings class="h-6 w-6" />
          </button>
          <button @click="showHelp = true"
            class="text-gray-400 hover:text-gray-200 focus:outline-none transition duration-200"
            aria-label="Information">
            <InfoIcon class="h-6 w-6" />
          </button>
        </div> -->
      </div>
    </div>
  </div>
</template>
