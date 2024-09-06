<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
import { Send } from 'lucide-vue-next';

const props = defineProps<{
  playerName: string;
  roomId: string;
  messages: { sender: string; text: string }[];
}>();

const emit = defineEmits<{
  (e: 'sendMessage', message: string): void;
}>();

const newMessage = ref('');
const chatContainer = ref<HTMLDivElement | null>(null);


const sendMessage = () => {
  if (newMessage.value.trim()) {
    emit('sendMessage', newMessage.value);
    newMessage.value = '';
  }
};

const scrollToBottom = () => {
  if (chatContainer.value) {
    chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
  }
};

const systemMessageMap = ref<Map<string, boolean>>(new Map());

const filteredMessages = computed(() => {
  const messages = props.messages.filter(message => {
    if (message.sender === 'system') {
      if (systemMessageMap.value.has(message.text)) {
        return false; // Skip duplicate system message
      } else {
        systemMessageMap.value.set(message.text, true); // Store new system message
      }
    }
    return true; // Include all non-system messages and unique system messages
  });
  return messages;
});

watch(() => props.messages, scrollToBottom, { deep: true });

</script>

<template>
  <div class="flex flex-col h-full">
    <div ref="chatContainer" class="flex-grow overflow-y-auto p-3 space-y-2">
      <div v-for="(message, index) in messages" :key="index" :class="['p-2 rounded-lg',
        message.sender === 'system' ? '' : (message.sender === playerName ? 'bg-blue-600 ml-auto' : 'bg-gray-700'),
        // message.sender === playerName ? 'bg-blue-600 ml-auto' : 'bg-gray-700',
        'max-w-[80%]']">
        <!-- <p class="rounded-lg">
          :class="['p-2 rounded-lg',
        message.sender === playerName ? 'bg-blue-600 ml-auto' : 'bg-gray-700',
        'max-w-[80%]']"
        </p> -->
        <p v-if="message.sender === 'system'"
          class="flex items-center text-sm text-gray-800 before:flex-1 before:border-t before:border-gray-200 before:me-6 after:flex-1 after:border-t after:border-gray-200 after:ms-6 dark:text-white dark:before:border-neutral-600 dark:after:border-neutral-600">
          {{ message.text }}</p>
        <p v-else class="rounded-lg text-sm font-semibold"
          :class="message.sender === playerName ? 'text-blue-200' : 'text-gray-300'">
          {{ message.sender }}
        <p class="text-white text-base">{{ message.text }}</p>
        </p>
      </div>
    </div>
    <div class="p-3 bg-gray-700">
      <div class="flex">
        <input v-model="newMessage" @keyup.enter="sendMessage" type="text" placeholder="Type a message..."
          class="flex-grow bg-gray-600 text-white rounded-l-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500" />
        <button @click="sendMessage"
          class="bg-blue-600 text-white rounded-r-lg px-3 py-2 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500">
          <Send class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>