<script setup lang="ts">
import { Copy } from 'lucide-vue-next'
import { onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
    randomGallery: string | undefined
    roomCode: string
}>()

defineEmits<{
    (e: 'cancel'): void
    (e: 'share'): void
    (e: 'add-bot'): void
}>()

const dots = ref('...')
let intervalId: number

onMounted(() => {
    intervalId = setInterval(() => {
        dots.value = dots.value.length >= 3 ? '' : dots.value + '.'
    }, 500)
})

onUnmounted(() => {
    clearInterval(intervalId)
})

const copyRoomCode = () => {
    if (navigator.clipboard) {
        navigator.clipboard.writeText(props.roomCode)
    }
}
</script>

<template>
    <div class="flex flex-col items-center justify-center py-12 px-6 w-full h-full max-w-400 mx-auto text-slate-100">
        <!-- Main Card Content -->
        <main class="w-full max-w-xl">
            <div
                class="bg-card rounded-[3rem] p-5 md:p-14 shadow-2xl border border-border text-center relative overflow-hidden w-full">
                <div class="absolute top-0 right-0 w-32 h-32 bg-secondary/5 rounded-bl-full pointer-events-none"></div>
                <div class="absolute bottom-0 left-0 w-24 h-24 bg-primary/5 rounded-tr-full pointer-events-none"></div>

                <div class="relative mb-12">
                    <div
                        class="relative w-40 h-40 mx-auto bg-muted rounded-3xl border border-border flex items-center justify-center overflow-hidden">
                        <!-- <div className="absolute inset-0 bg-primary rounded-full opacity-20 animate-ping"></div>
                        <div
                            className="relative z-10 w-24 h-24 bg-white dark:bg-surface-dark rounded-full shadow-[0_0_30px_rgba(43,238,121,0.3)] border-4 border-primary/20 flex flex-col items-center justify-center overflow-hidden">
                        </div> -->
                        <img :src="randomGallery" alt="Loading"
                            class="w-full h-full object-cover opacity-80 animate-float" />
                    </div>
                    <div
                        class="absolute -inset-4 border-2 border-primary/20 rounded-[2.5rem] animate-pulse-slow pointer-events-none">
                    </div>
                </div>

                <h2 class="text-2xl md:text-3xl font-bold text-foreground mb-6">
                    Waiting for an opponent{{ dots }}
                </h2>
                <!-- <p class="text-slate-500 dark:text-slate-400 mb-10">
                    Waiting for someone to join your game.
                </p> -->

                <div class="space-y-4 mb-10">
                    <p class="text-sm font-semibold text-foreground/80">
                        Share the room code to play with a friend:
                    </p>
                    <div
                        class="bg-muted/30 rounded-2xl p-2 pl-8 flex items-center justify-between border border-border max-w-sm mx-auto shadow-inner group transition-colors hover:border-primary/50">
                        <span class="truncate text-3xl font-bold text-foreground tracking-widest">{{
                            roomCode
                            }}</span>
                        <button
                            class="bg-card text-muted-foreground px-6 py-3 rounded-xl shadow-sm border border-border hover:text-primary hover:border-primary active:scale-95 transition-colors font-bold text-sm flex items-center gap-2"
                            @click="copyRoomCode">
                            <Copy class="w-4 h-4" />
                            Copy
                        </button>
                    </div>
                </div>

                <button
                    class="w-full max-w-xs py-4 px-6 rounded-2xl bg-primary text-foreground shadow-xl shadow-primary/20 font-bold hover:brightness-105 transition-all text-sm mb-4"
                    @click="$emit('add-bot')">
                    Play with Bot
                </button>
                <button
                    class="w-full max-w-xs py-3.5 rounded-xl bg-muted text-foreground font-bold hover:bg-muted/80 hover:text-foreground transition-all text-sm mb-4"
                    @click="$emit('cancel')">
                    Cancel
                </button>
            </div>
        </main>
    </div>
</template>

<style scoped>
.animate-float {
    animation: float 3s ease-in-out infinite;
}

.animate-pulse-slow {
    animation: pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes float {

    0%,
    100% {
        transform: translateY(0);
    }

    50% {
        transform: translateY(-10px);
    }
}
</style>
