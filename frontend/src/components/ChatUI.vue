<script setup lang="ts">
import { SendIcon, XIcon } from 'lucide-vue-next';
import { nextTick, onMounted, ref } from 'vue';

const props = defineProps<{
  playerName: string;
  messages: { sender: string; text: string }[];
  users: number
}>();

const emit = defineEmits<{
  (e: 'sendMessage', message: string): void;
  (e: 'toggle'): void;
}>();



const newMessage = ref('')
const messageList = ref<HTMLDivElement | null>(null);

const sendMessage = () => {
  if (newMessage.value.trim()) {
    // messages.value.push({
    //   sender: "You",
    //   text: newMessage.value.trim()
    // })
    emit('sendMessage', newMessage.value);

    newMessage.value = ''
  }
  scrollToBottom()
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

const shouldShowSender = (index: number) => {
  if (index === 0) return true
  const prevMessage = props.messages[index - 1]
  return prevMessage.sender !== props.messages[index].sender || prevMessage.sender === 'system'
}

// const shouldShowAvatar = (index: number) => {
//   return shouldShowSender(index) && props.messages[index].sender !== props.playerName
// }

const isLastMessageFromSameSender = (index: number) => {
  return index < props.messages.length - 1 && props.messages[index + 1].sender === props.messages[index].sender
}


onMounted(() => {
  scrollToBottom()
})
</script>

<template>
  <div class="bg-gray-900 text-white flex flex-col shadow-xl">
    <!-- Chat Header -->
    <div class="p-4 border-b border-gray-800 flex items-center justify-between">
      <h2 class="text-lg font-semibold">Chat</h2>
      <div class="flex items-center text-sm text-gray-400">
        <div class="w-2 h-2 rounded-full bg-green-500 mr-2"></div>
        <span>{{ users }} online</span>
      </div>
      <button @click="$emit('toggle')" class="text-gray-400 hover:text-white transition-colors">
        <XIcon :size="20" />
      </button>
    </div>

    <!-- Message List -->
    <div class="flex-1 overflow-y-auto p-4 space-y-2 scrollbar-thin" ref="messageList">
      <div v-for="(message, index) in messages" :key="index" class="flex flex-col space-y-1">
        <div v-if="message.sender === 'system'" class="w-full flex justify-center">
          <p
            class="flex items-center text-sm text-gray-400 before:flex-1 before:border-t before:border-gray-700 before:mr-2 after:flex-1 after:border-t after:border-gray-700 after:ml-2">
            {{ message.text }}
          </p>
        </div>
        <div v-else class="flex flex-col space-y-1">
          <div :class="['flex items-end', message.sender === playerName ? 'justify-end' : '']">
            <!-- <div v-if="message.sender !== playerName && shouldShowAvatar(index)"
              class="w-6 h-6 rounded-full bg-gray-700 flex items-center justify-center flex-shrink-0 mr-2 mb-auto">
              <UserIcon :size="12" />
            </div> -->
            <div class="flex flex-col" :class="{ 'items-end': message.sender === playerName }">
              <span v-if="shouldShowSender(index)" class="text-xs mb-1"
                :class="message.sender === playerName ? 'text-blue-300' : 'text-gray-400'">
                {{ message.sender }}
              </span>
              <div class="px-4 py-2 rounded-2xl" :class="{
                'bg-blue-600 rounded-br-none': message.sender === playerName,
                'bg-gray-700': message.sender !== playerName,
                'rounded-br-none': message.sender === playerName && !isLastMessageFromSameSender(index),
                'rounded-bl-none': message.sender !== playerName && !isLastMessageFromSameSender(index),
                'mt-1': !shouldShowSender(index)
              }">
                <p class="text-sm">{{ message.text }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Message Input -->
    <div class="p-4 border-t border-gray-800">
      <form @submit.prevent="sendMessage" class="flex items-center space-x-2">
        <input v-model="newMessage" @keyup.enter="sendMessage" type="text" placeholder="Type a message..."
          class="flex-1 bg-gray-800 text-white rounded-full px-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-gray-700" />
        <button type="submit"
          class="bg-blue-600 hover:bg-blue-700 text-white rounded-full p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 transition-colors">
          <SendIcon :size="16" />
        </button>
      </form>
    </div>
  </div>
</template>

<style>
/* Webkit-based browsers (Chrome, Safari, newer versions of Edge) */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: #1f2937;
  /* gray-800 */
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: #4b5563;
  /* gray-600 */
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: #6b7280;
  /* gray-500 */
}

/* Firefox */
.scrollbar-thin {
  scrollbar-width: thin;
  scrollbar-color: #4b5563 #1f2937;
  /* thumb and track color */
}
</style>