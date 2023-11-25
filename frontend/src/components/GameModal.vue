<script setup lang="ts">
import type { WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import router from '@/router'
import { useWebSocket } from '@vueuse/core'
import { ref, watchEffect } from 'vue'

const roomId = ref('')
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
</script>

<template>
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
            type="submit">
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

