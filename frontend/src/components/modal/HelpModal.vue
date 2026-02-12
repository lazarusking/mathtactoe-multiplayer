<script setup lang="ts">
import BaseModal from './BaseModal.vue';
import { CheckIcon } from 'lucide-vue-next';

defineProps<{ showHelp: boolean }>();
defineEmits<{ (event: 'close-modal'): void }>();
</script>

<template>
  <BaseModal @close-modal="$emit('close-modal')" :show-help="showHelp">
    <template #title>
      <h2 class="text-2xl font-bold text-foreground">How to Play</h2>
    </template>
    <template #content>
      <section class="space-y-6">
        <div class="bg-muted/50 p-4 rounded-2xl border border-border">
          <h3 class="text-lg font-semibold text-gray-100 mb-2">Game Rules:</h3>
          <ul class="list-disc list-inside text-gray-300 space-y-2">
            <li>Players take turns selecting a number from 1 to 9.</li>
            <li>Each number can only be used once.</li>
            <li>The first player to create a line summing to 15 wins.</li>
            <li>If all numbers are used and no line sums to 15, the game is a draw.</li>
          </ul>
        </div>

        <div>
          <h3 class="text-lg font-semibold text-foreground mb-2">Winning Example (15):</h3>
          <div class="grid w-full justify-self-center max-w-md grid-cols-3 grid-rows-3 gap-3 shadow-md">
            <template v-for="y in [7, 4, 1, 3, 2, 9, 5, 6, 8]" :key="y">
              <button type="button" :class="[
                'grid items-center justify-center w-auto h-auto p-2 text-xl md:text-3xl font-bold text-white transition-colors rounded-xl shadow-md focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-foreground focus-visible:ring-offset-2',
                [7, 3, 5].includes(y) ?
                  'bg-primary/80 hover:bg-primary/60'
                  : 'bg-muted hover:bg-muted/80 '
              ]">
                <div class="flex items-center justify-center">
                  {{ y }}
                  <CheckIcon v-if="[7, 3, 5].includes(y)" class="h-4 w-4 ml-1 text-foreground" />
                </div>
              </button>
            </template>
          </div>
          <p class="text-sm text-muted-foreground mt-2">
            7 + 3 + 5 = 15, forming a winning diagonal.
          </p>
        </div>

      </section>
    </template>
  </BaseModal>
</template>

<style scoped>
/* Add any component-specific styles here */
</style>