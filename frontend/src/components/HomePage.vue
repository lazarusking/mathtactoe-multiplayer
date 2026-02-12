<script setup lang="ts">
import type { WSMessage } from '@/interface'
import { websocket } from '@/lib/socket'
import router from '@/router'
import { useWebSocket } from '@vueuse/core'
import { Info, Rocket } from 'lucide-vue-next'
import { defineAsyncComponent, onMounted, ref, watch } from 'vue'

const HelpModal = defineAsyncComponent(() => import('./modal/HelpModal.vue'))

const roomId = ref('')
const playerName = ref('')
const showHelp = ref(false)
const activeTab = ref<'join' | 'create'>('create')

onMounted(() => {
  const storedUsername = localStorage.getItem('username')
  if (storedUsername) {
    playerName.value = storedUsername
  }
})
const { send } = useWebSocket(websocket.url, {
  onMessage(ws, event) {
    handleMessage(event)
  },
  autoReconnect: true,
  immediate: true
})

watch(playerName, () => {
  localStorage.setItem('username', playerName.value)
})

function createRoom() {
  // const data = { action: 'join-room', message: null, sender: { name: playerName.value } }
  const data = { action: 'create-room', data: null, sender: { name: playerName.value } }
  console.log(data, 'createRoom data')

  send(JSON.stringify(data))
  //clicking this sends a msg to ws, ws then joins room and sends room id created
  //client then leaves the room because the handleMessage navigates to a new page
  //new page GameView calls another join-room,creates another new client for the same client
}

function handleMessage(event: MessageEvent) {
  console.log(event.data)

  const data: WSMessage<string> = JSON.parse(event.data)

  switch (data.action) {
    case 'join-room':
      console.log(data)
      console.log('Did this run')
      router.push({ name: 'room', params: { room: data.data } })
      // router.push(data.message)
      break
    case 'create-room':
      console.log(data)
      console.log('Did this run')
      router.push({ name: 'room', params: { room: data.data } })
      // router.push(data.message)
      break
    // case 'send-message': {
    //   // WSState.data = WSMessage
    //   console.log(data)
    //   router.push({ name: "room", params: { room: data.data } })
    //   // console.log(WSMessage);
    //   break
    // }
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
    router.push({ name: 'room', params: { username: playerName.value, room: roomId.value } })
  }
}


</script>

<template>
  <div
    class="min-h-screen flex flex-col items-center justify-center p-6 bg-background text-foreground transition-colors duration-300 relative overflow-hidden font-sans">
    <!-- Background Decor -->
    <div
      class="fixed top-[-10%] left-[-10%] w-[40%] h-[40%] bg-primary/10 blur-[120px] rounded-full -z-10 pointer-events-none">
    </div>
    <div
      class="fixed bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-accent/10 blur-[120px] rounded-full -z-10 pointer-events-none">
    </div>

    <HelpModal :show-help="showHelp" @close-modal="closeModal" />

    <header class="mb-12 text-center">
      <div class="flex items-center justify-center gap-2 mb-4">
        <h1 class="text-5xl md:text-6xl font-extrabold font-display tracking-tight">
          <span class="text-secondary">TicTac</span>
          <span class="text-primary">Math</span>
        </h1>
        <button class="text-muted-foreground hover:text-primary transition-colors" @click="showHelp = true">
          <Info class="w-8 h-8" />
        </button>
      </div>
      <div
        class="inline-flex items-center px-4 py-1.5 rounded-full bg-muted text-muted-foreground text-sm font-semibold tracking-wide uppercase border border-border">
        Multiplayer
      </div>
    </header>

    <main class="w-full max-w-md">
      <div class="bg-card border border-border shadow-2xl rounded-4xl overflow-hidden transition-all duration-300">
        <!-- Tabs -->
        <div class="flex border-b border-border">
          <button :class="[
            'flex-1 py-5 text-sm font-bold transition-colors relative',
            activeTab === 'join'
              ? 'text-primary bg-muted/30'
              : 'text-muted-foreground hover:bg-muted/10'
          ]" @click="activeTab = 'join'">
            Join Game
            <span v-if="activeTab === 'join'"
              class="absolute bottom-0 left-0 w-full h-1 bg-primary rounded-t-full"></span>
          </button>
          <button :class="[
            'flex-1 py-5 text-sm font-bold transition-colors relative',
            activeTab === 'create'
              ? 'text-primary bg-muted/30'
              : 'text-muted-foreground hover:bg-muted/10'
          ]" @click="activeTab = 'create'">
            Create Game
            <span v-if="activeTab === 'create'"
              class="absolute bottom-0 left-0 w-full h-1 bg-primary rounded-t-full"></span>
          </button>
        </div>

        <!-- Tab Content -->
        <div class="p-8">
          <div class="space-y-8">
            <div v-if="activeTab === 'join'" class="space-y-2">
              <label class="block text-sm font-bold text-muted-foreground px-1" for="room-id">
                Room ID
              </label>
              <input id="room-id" v-model="roomId" type="text"
                class="w-full px-5 py-4 bg-muted/20 border-2 border-border rounded-2xl text-foreground focus:border-primary focus:ring-0 transition-all outline-none font-medium placeholder:text-muted-foreground/50"
                placeholder="Enter Room ID" />
            </div>

            <div class="space-y-2">
              <label class="block text-sm font-bold text-muted-foreground px-1" for="player-name">
                Your Name
              </label>
              <input id="player-name" v-model="playerName" type="text"
                class="w-full px-5 py-4 bg-muted/20 border-2 border-border rounded-2xl text-foreground focus:border-primary focus:ring-0 transition-all outline-none font-medium placeholder:text-muted-foreground/50"
                placeholder="Enter your name" />
            </div>

            <button
              class="w-full bg-primary hover:bg-primary/90 text-foreground font-bold py-5 px-6 rounded-2xl text-lg shadow-lg shadow-primary/20 transition-all transform active:scale-[0.98] flex items-center justify-center gap-2"
              @click="activeTab === 'create' ? createRoom() : joinGame()">
              <span>{{ activeTab === 'create' ? 'Start New Match' : 'Find Match' }}</span>
              <Rocket class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
