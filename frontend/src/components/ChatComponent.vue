<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-900 p-4">
    <div class="mb-4 text-2xl font-bold text-white">{{ status }}</div>
    <div
      class="w-full max-w-md aspect-square p-1 bg-gradient-to-br from-blue-500 via-violet-500 to-pink-500 rounded-lg shadow-lg">
      <div class="grid grid-cols-3 w-full h-full p-0.5 rounded-lg overflow-hidden">
        <button v-for="(value, index) in board" :key="index" @click="handleClick(index)"
          class="w-full h-full bg-gray-900 flex items-center justify-center text-4xl font-bold focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-opacity-50 transition-all duration-200 border-purple-500 border-4 border-gradient iagonal-split">
          <span :class="{ 'text-blue-400': value === 'X', 'text-purple-400': value === 'O' }">
            {{ value }}
          </span>
        </button>
      </div>
    </div>
    <div class="mt-6 flex justify-center space-x-4">
      <button v-for="num in [1, 3, 5, 7, 9]" :key="num"
        class="button-gradient w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-full flex items-center justify-center text-xl font-bold text-white shadow-lg transition-all duration-200 hover:from-blue-600 hover:to-purple-600 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-gray-900">
        {{ num }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

type CellValue = 'X' | 'O' | null

const board = ref<CellValue[]>(['X', 'O', 'X', 'O', 'X', 'O', 'O', 'X', 'O'])
const xIsNext = ref(true)

const handleClick = (index: number) => {
  if (board.value[index] || calculateWinner(board.value)) return

  board.value[index] = xIsNext.value ? 'X' : 'O'
  xIsNext.value = !xIsNext.value
}

const calculateWinner = (squares: CellValue[]): CellValue => {
  const lines = [
    [0, 1, 2],
    [3, 4, 5],
    [6, 7, 8],
    [0, 3, 6],
    [1, 4, 7],
    [2, 5, 8],
    [0, 4, 8],
    [2, 4, 6],
  ]
  for (let i = 0; i < lines.length; i++) {
    const [a, b, c] = lines[i]
    if (squares[a] && squares[a] === squares[b] && squares[a] === squares[c]) {
      return squares[a]
    }
  }
  return null
}

const status = computed(() => {
  const winner = calculateWinner(board.value)
  if (winner) {
    return `Winner: ${winner}`
  } else if (board.value.every(Boolean)) {
    return "It's a draw!"
  } else {
    return `Next player: ${xIsNext.value ? 'X' : 'O'}`
  }
})
</script>

<style scoped>
.border-gradient {
  border-image: linear-gradient(to bottom right, #3b82f6, #a855f7, #ec4899) 1;
  /* border-image: linear-gradient(to bottom right, theme( gradientColorStops.pink.500)) 1; */
}

.button-gradient {
  background-image: linear-gradient(60deg, #29323c 0%, #485563 100%);
}

/* .border-gradient {
  border-image: linear-gradient(to bottom right, #3b82f6 50%, #a855f7 50%, #ec4899 50%) 1;
  /* border-image: linear-gradient(to bottom right, theme( gradientColorStops.pink.500)) 1; */
/* } */
.diagonal-split {
  background: linear-gradient(to bottom right, #3b82f6 50%, #a855f7 50%, #ec4899 50%) 1;
}

.cell:nth-child(-n+4) {
  border-image: linear-gradient(to bottom right, #3b82f6 10%, #a855f7 40%, #ec4899 50%) 1;
}

/* 5th to 9th cells */
.cell:nth-child(n+5) {
  background-image: linear-gradient(to bottom right, #3b82f6 50%, #a855f7 40%, #ec4899 10%) 1;
}
</style>