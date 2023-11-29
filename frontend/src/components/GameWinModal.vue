<script setup lang="ts">
defineProps<{ gameStatus: { gameWon: boolean; gameOver: boolean, gameDraw: boolean } }>()
defineEmits<{ (event: 'play-again'): void }>()

</script>
<template>
  <template v-if="gameStatus.gameOver">
    <div v-if="gameStatus.gameDraw" aria-modal="true"
      class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50" role="dialog">
      <div class="rounded-lg bg-gray-800 w-full max-w-md mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <h2 class="text-2xl font-bold text-center dark:text-white mb-4">No one wins!</h2>
        <p class="text-center text-gray-700 dark:text-gray-200 mb-6">
          You both tied this game
        </p>
        <div class="flex justify-center">
          <button @click="$emit('play-again')"
            class="py-2 text-white bg-slate-700 rounded-md text-center hover:bg-slate-900/70 inline-flex items-center justify-center text-sm transition-colors px-4"
            type="button">
            Play Again
          </button>
        </div>
      </div>
    </div>
    <div v-else-if="gameStatus.gameWon" aria-modal="true"
      class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50" role="dialog">
      <div class="bg-white rounded-lg dark:bg-gray-800 w-full max-w-md mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <h2 class="text-2xl font-bold text-center dark:text-white mb-4">Congratulations!</h2>
        <p class="text-center text-gray-700 dark:text-gray-200 mb-6">
          You have won the game.
        </p>
        <div class="flex justify-center">
          <button @click="$emit('play-again')"
            class="py-2 text-white bg-green-700 rounded-md text-center hover:bg-slate-900/70 inline-flex items-center justify-center text-sm transition-colors px-4"
            type="button">
            Play Again
          </button>
        </div>
      </div>
    </div>
    <div v-else aria-modal="true" class="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50"
      role="dialog">
      <div class="rounded-lg bg-gray-800 w-full max-w-md mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <h2 class="text-2xl font-bold text-center dark:text-white mb-4">Game Over</h2>
        <p class="text-center text-gray-700 dark:text-gray-200 mb-6">
          Unfortunately, you have lost the game.
        </p>
        <div class="flex justify-center">
          <button @click="$emit('play-again')"
            class="py-2 text-white bg-red-500 hover:bg-slate-900/70 rounded-md text-center inline-flex items-center justify-center text-sm transition-colors h-10 px-4"
            type="button">
            Try Again
          </button>
        </div>
      </div>
    </div>
  </template>
</template>
