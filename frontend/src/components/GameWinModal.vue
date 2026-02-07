<script setup lang="ts">

defineProps<{ won: Boolean, draw: Boolean, winner: String | null }>()
defineEmits<{ (event: 'play-again'): void }>()

</script>
<template>
  <div aria-modal="true" class="fixed inset-0 flex items-center justify-center z-50 bg-black/60 backdrop-blur-sm"
    role="dialog">
    <div
      class="bg-white dark:bg-surface-dark rounded-[2.5rem] w-full max-w-sm mx-4 p-10 shadow-2xl border border-slate-100 dark:border-slate-700 text-center scale-up">
      <template v-if="draw">
        <div class="text-6xl mb-6">🤝</div>
        <h2 class="text-3xl font-bold mb-2">It's a Draw!</h2>
        <p class="text-slate-500 mb-8">No one reached 15 this time.</p>
      </template>
      <template v-else-if="won">
        <div class="text-6xl mb-6">🏆</div>
        <h2 class="text-3xl font-bold mb-2 text-primary">Victory!</h2>
        <p class="text-gray-700 dark:text-gray-200 mb-8">You calculated your way to the top!</p>
      </template>
      <template v-else>
        <div class="text-6xl mb-6">😞</div>
        <h2 class="text-3xl font-bold mb-2 text-accent">Defeat</h2>
        <p class="text-slate-500 mb-8">{{ winner }} won the match.</p>
      </template>
      <button @click="$emit('play-again')"
        class="w-full py-4 bg-primary hover:bg-emerald-600 text-white font-bold rounded-2xl transition-all shadow-lg shadow-emerald-500/20 active:scale-95">
        New Match
      </button>
    </div>
  </div>
</template>
