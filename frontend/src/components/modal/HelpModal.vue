<script setup lang="ts">
import BaseModal from './BaseModal.vue';
import { CheckIcon } from 'lucide-vue-next';

defineProps<{ showHelp: boolean }>();
defineEmits<{ (event: 'close-modal'): void }>();
</script>

<template>
  <BaseModal @close-modal="$emit('close-modal')" :show-help="showHelp">
    <template #title>
      <h2 class="text-2xl font-bold text-gray-100">How to Play TicTacToe Math</h2>
    </template>
    <template #content>
      <div class="space-y-6">
        <!-- <p class="text-base text-gray-300 leading-relaxed">
          TicTacToe Math combines the classic game of Tic-Tac-Toe with a mathematical twist. The goal is to form a line of three grids in <strong class="text-blue-400">any direction</strong> (horizontal, vertical, or diagonal) where the sum of the numbers equals <strong class="text-green-400">15</strong>.
        </p> -->
        
        <div class="bg-gray-800 p-4 rounded-lg">
          <h3 class="text-lg font-semibold text-gray-100 mb-2">Game Rules:</h3>
          <ul class="list-disc list-inside text-gray-300 space-y-2">
            <li>Players take turns selecting a number from 1 to 9.</li>
            <li>Each number can only be used once.</li>
            <li>The first player to create a line summing to 15 wins.</li>
            <li>If all numbers are used and no line sums to 15, the game is a draw.</li>
          </ul>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-gray-100 mb-2">Example Winning Combination:</h3>
          <div class="grid w-full max-w-md grid-cols-3 grid-rows-3 gap-4 shadow-md">
            <template v-for="y in [7, 4, 1, 3, 2, 9, 5, 6, 8]" :key="y">
              <button type="button"
                :class="[
                  'grid items-center justify-center w-auto h-auto p-2 text-2xl md:text-4xl font-bold text-white transition-colors rounded-lg shadow-md place-content-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:pointer-events-none',
                  [7, 3, 5].includes(y) 
                    ? 'bg-green-600 hover:bg-green-700' 
                    : 'bg-gray-700 hover:bg-gray-600'
                ]"
              >
                <div class="flex items-center justify-center">
                  {{ y }}
                  <CheckIcon v-if="[7, 3, 5].includes(y)" class="h-5 w-5 ml-1 text-white" />
                </div>
              </button>
            </template>
          </div>
          <p class="text-sm text-gray-400 mt-2">
            In this example, 7 + 3 + 5 = 15, forming a winning diagonal line.
          </p>
        </div>

        <!-- <div class="bg-blue-900 bg-opacity-20 p-4 rounded-lg">
          <h3 class="text-lg font-semibold text-blue-300 mb-2">Strategy Tip:</h3>
          <p class="text-gray-300">
            Remember the magic square: 2-7-6, 9-5-1, 4-3-8. These are all the possible combinations that sum to 15.
          </p>
        </div> -->
      </div>
    </template>
  </BaseModal>
</template>

<style scoped>
/* Add any component-specific styles here */
</style>