<script setup lang="ts">

defineProps<{ won: boolean, draw: boolean, winner: string | null }>()
defineEmits<{ (event: 'play-again'): void }>()

</script>
<template>
  <div aria-modal="true"
    class="fixed inset-0 font-display flex items-center justify-center z-50 bg-background/80 backdrop-blur-sm"
    role="dialog">
    <div
      class="bg-card rounded-[2.5rem] w-full max-w-sm mx-4 p-10 shadow-2xl border border-border text-center scale-up">
      <template v-if="draw">
        <div class="text-6xl mb-6">🤝</div>
        <h2 class="text-3xl font-bold mb-2">It's a Draw!</h2>
        <p class="text-muted-foreground mb-8">No one reached 15 this time.</p>
      </template>
      <template v-else-if="won">
        <div class="text-6xl mb-6">🏆</div>
        <h2 class="text-3xl font-bold mb-2 text-primary">Victory!</h2>
        <p class="text-muted-foreground mb-8">You calculated your way to the top!</p>
      </template>
      <template v-else>
        <div class="text-6xl mb-6">😞</div>
        <h2 class="text-3xl font-bold mb-2 text-accent">Defeat</h2>
        <p class="text-muted-foreground mb-8">{{ winner }} won the match.</p>
      </template>
      <button @click="$emit('play-again')"
        class="w-full py-4 bg-primary/80 hover:bg-primary/90 text-foreground font-bold rounded-2xl transition-all shadow-lg shadow-primary/20 active:scale-95">
        New Match
      </button>
    </div>
  </div>
</template>
