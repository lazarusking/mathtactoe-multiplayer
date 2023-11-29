<script setup lang="ts">
import type { WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import router from '@/router'
import { useWebSocket } from '@vueuse/core'
import { ref, watchEffect } from 'vue'
import HelpModal from './modal/HelpModal.vue'

const roomId = ref('')
const showHelp = ref(false)
const { data, send } = useWebSocket(websocket.url, {
  onMessage(ws, event) {
    handleMessage(event)
  }
})

watchEffect(() => {
  console.log(data)
})
function createRoom() {
  send(JSON.stringify({ action: 'join-room', message: null }))
}

function handleMessage(event: MessageEvent) {
  const data: WSMessage = JSON.parse(event.data)
  switch (data.action) {
    case 'join-room':
      console.log(data)
      router.push(data.message)
      break
    default:
      break
  }
}
function closeModal() {
  showHelp.value = false
}

</script>

<template>
  <div @click="showHelp = !showHelp" class="flex mx-auto container max-w-screen-md  py-1 sm:py-2 sm:px-6 lg:px-8">
    <div class="ml-auto flex flex-wrap items-center justify-between">
      <button type="button"
        class="text-white items-center text-lg p-3 font-semibold rounded hover:bg-blue-500 hover:text-white">
        <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z">
          </path>
        </svg>
        <span class="sr-only">How to play?</span>
      </button>
    </div>
  </div>
  <HelpModal :show-help="showHelp" @close-modal="closeModal" />
  <div class="flex justify-center items-center space-x-2">
    <h1 class="text-2xl font-bold text-center">TicTacToe Math</h1>
    <!-- <img src="@/assets/three-in-a-row.png" class="w-32" /> -->
    <!-- <img src="@/assets/tic-tac-toe.png" class="w-32" /> -->
  </div>
  <div aria-modal="true" class="z-50 flex items-center justify-center flex-1 text-gray-500 bg-opacity-50" role="dialog">
    <div class="w-auto max-w-md px-4 py-8 mx-auto bg-gray-800 rounded-lg sm:w-80 md:w-full sm:px-6 lg:px-8">
      <h2 class="mb-4 text-2xl font-bold text-center text-white">Enter Game</h2>
      <form>
        <div class="mb-3">
          <label class="block text-sm font-medium text-gray-200" htmlFor="room-id">
            Room ID
          </label>
          <input v-model="roomId"
            class="block w-full px-3 py-2 mt-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-gray-500 focus:border-gray-500 sm:text-sm"
            id="room-id" placeholder="Enter Room ID" type="text" />
        </div>
        <div class="items-center justify-between space-y-2 font-medium md:flex md:space-y-0 md:space-x-2">
          <button @click="$router.push(roomId)" class="w-full py-2 text-center text-white bg-blue-700 rounded-md md:w-1/2"
            type="button">
            Join Game
          </button>
          <button @click="createRoom" class="w-full py-2 text-center text-white bg-green-700 rounded-md md:w-1/2"
            type="button">
            Start New Game
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

