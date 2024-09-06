<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue';
import { Send, XIcon, UserIcon, SendIcon } from 'lucide-vue-next';

const props = defineProps<{
  playerName: string;
  messages: { sender: string; text: string }[];
}>();

const emit = defineEmits<{
  (e: 'sendMessage', message: string): void;
  (e: 'toggle'): void;
}>();


const messages = ref([
  { sender: "John Doe", text: "Hey there! How are you?" },
  { sender: "You", text: "I'm doing great, thanks for asking!" },
  { sender: "system", text: "Someone joined" },
  { sender: "John Doe", text: "That's wonderful to hear. Any plans for the weekend?" },
  { sender: "You", text: "Yes, I'm thinking of going hiking. How about you?" },
])

const newMessage = ref('')
const messageList = ref<HTMLDivElement | null>(null);

const sendMessage = () => {
  if (newMessage.value.trim()) {
    messages.value.push({
      sender: "You",
      text: newMessage.value.trim()
    })
    newMessage.value = ''
    scrollToBottom()
  }
}

// const chatContainer = ref<HTMLDivElement | null>(null);


// const sendMessage = () => {
//   if (newMessage.value.trim()) {
//     emit('sendMessage', newMessage.value);
//     newMessage.value = '';
//   }
//   scrollToBottom()
// };

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

const toggleChat = () => {
  // Implement chat toggle functionality here
  console.log("Toggle chat")
}

onMounted(() => {
  scrollToBottom()
})
</script>

<template>
  <div class="fixed right-0 top-0 bottom-0 w-80 bg-gray-900 text-white flex flex-col shadow-xl">

    <div class="p-4 border-b border-gray-800 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Chat</h2>
      <button @click="$emit('toggle')" class="text-gray-400 hover:text-white transition-colors">
        <XIcon :size="20" />
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 space-y-4" ref="messageList">
      <div v-for="(message, index) in messages" :key="index" class="flex flex-col space-y-1">
        <div class="flex items-end" :class="{ 'justify-end': message.sender === 'You' }">
          <div v-if="message.sender !== 'You'"
            class="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center flex-shrink-0 mr-2">
            <UserIcon :size="12" />
          </div>
          <div class="px-4 py-2 rounded-2xl max-w-[80%]" :class="{
            'bg-blue-600 rounded-br-none': message.sender === 'You',
            'bg-gray-700 rounded-bl-none': message.sender !== 'You'
          }">
            <p class="text-sm">{{ message.text }}</p>
          </div>
        </div>
        <span class="text-xs text-gray-500" :class="{ 'self-end': message.sender === 'You' }">
          {{ message.sender }}
        </span>
      </div>
    </div>

    <!-- Message Input -->
    <div class="p-4 border-t border-gray-800">
      <form @submit.prevent="sendMessage" class="flex items-center space-x-2">
        <input v-model="newMessage" type="text" placeholder="Type a message..."
          class="flex-1 bg-gray-800 text-white rounded-full px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-700" />
        <button @click="sendMessage"
          class="bg-blue-600 hover:bg-blue-700 text-white rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors">
          <SendIcon :size="16" />
        </button>
      </form>
    </div>
  </div>
</template>
