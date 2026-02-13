<script setup lang="ts">
import { SendIcon, XIcon } from 'lucide-vue-next';
import { nextTick, onMounted, ref } from 'vue';

const props = defineProps<{
  playerID: string
  messages: { sender: string; text: string; id: string }[]
  users: number
}>()

const emit = defineEmits<{
  (e: 'sendMessage', message: string): void
  (e: 'toggle'): void
}>()

const newMessage = ref('')
const messageList = ref<HTMLDivElement | null>(null)

const sendMessage = () => {
  if (newMessage.value.trim()) {
    // messages.value.push({
    //   sender: "You",
    //   text: newMessage.value.trim()
    // })
    emit('sendMessage', newMessage.value)

    newMessage.value = ''
  }
  scrollToBottom()
}

// const chatContainer = ref<HTMLDivElement | null>(null);

// const scrollToBottom = () => {
//   if (chatContainer.value) {
//     chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
//   }
// };

const scrollToBottom = async () => {
  await nextTick()
  if (messageList.value) {
    messageList.value.scrollTop = messageList.value.scrollHeight
  }
}

const shouldShowSender = (index: number) => {
  if (index === 0) return true
  const prevMessage = props.messages[index - 1]
  const currentMessage = props.messages[index]
  return (
    prevMessage &&
    currentMessage &&
    (prevMessage.id !== currentMessage.id || prevMessage.sender === 'system')
  )
}

// const shouldShowAvatar = (index: number) => {
//   return shouldShowSender(index) && props.messages[index].sender !== props.playerName
// }

const isLastMessageFromSameSender = (index: number) => {
  const currentMessage = props.messages[index]
  const nextMessage = props.messages[index + 1]
  return (
    index < props.messages.length - 1 &&
    currentMessage &&
    nextMessage &&
    nextMessage.sender === currentMessage.sender
  )
}

onMounted(() => {
  scrollToBottom()
})
</script>

<template>
  <aside class="bg-gray-900 border border-border text-foreground flex flex-col h-full">
    <!-- Chat Header -->
    <header class="p-4 border-b border-border/50 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Chat</h2>
      <div class="flex items-center text-sm text-muted-foreground">
        <div class="w-2 h-2 rounded-full bg-primary mr-2"></div>
        <span>{{ users }} online</span>
      </div>
      <button class="text-muted-foreground hover:text-foreground transition-colors" @click="$emit('toggle')">
        <XIcon :size="20" />
      </button>
    </header>

    <!-- Message List -->
    <ul ref="messageList" class="flex-1 overflow-y-auto p-4 space-y-2 scrollbar-thin">
      <li v-for="(message, index) in messages" :key="index" class="flex flex-col space-y-1">
        <div v-if="message.sender === 'system'" class="w-full flex justify-center">
          <span class="text-xs text-muted-foreground bg-muted px-2 py-1 rounded-md">{{
            message.text
            }}</span>
        </div>
        <div v-else class="flex flex-col space-y-1">
          <div :class="['flex items-end', message.id === playerID ? 'justify-end' : '']">
            <!-- <div v-if="message.sender !== playerName && shouldShowAvatar(index)"
              class="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center flex-shrink-0 mr-2 mb-auto">
              <UserIcon :size="12" />
            </div> -->
            <div class="flex flex-col" :class="{ 'items-end': message.id === playerID }">
              <span v-if="shouldShowSender(index)" class="text-xs mb-1"
                :class="message.id === playerID ? 'text-blue-300' : 'text-gray-400'">
                {{ message.sender }}
              </span>
              <div class="px-4 py-2 rounded-2xl" :class="{
                'bg-blue-600 rounded-br-none': message.id === playerID,
                'bg-gray-700': message.id !== playerID,
                'rounded-br-none': message.id === playerID && !isLastMessageFromSameSender(index),
                'rounded-bl-none': message.id !== playerID && !isLastMessageFromSameSender(index),
                'mt-1': !shouldShowSender(index)
              }">
                <p class="text-sm">{{ message.text }}</p>
              </div>
            </div>
          </div>
        </div>
      </li>
    </ul>

    <!-- Message Input -->
    <footer class="p-4 border-t border-border/50">
      <form class="flex items-center space-x-2" @submit.prevent="sendMessage">
        <input v-model="newMessage" type="text" placeholder="Type a message..."
          class="flex-1 bg-gray-800 text-white rounded-full px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-700"
          @keyup.enter="sendMessage" />
        <button type="submit"
          class="bg-blue-600 hover:bg-blue-700 text-foreground rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors">
          <SendIcon :size="16" />
        </button>
      </form>
    </footer>
  </aside>
</template>

<style>
/* Webkit-based browsers (Chrome, Safari, newer versions of Edge) */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: var(--card);
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: var(--muted);
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: var(--muted-foreground);
}
</style>
