<script setup lang="ts">
import { Copy } from 'lucide-vue-next';
import { onMounted, onUnmounted, ref } from 'vue';

const props = defineProps<{
    randomGallery: string,
    roomCode: string
}>()

const emit = defineEmits<{
    (e: 'cancel'): void,
    (e: 'share'): void,
    (e: 'add-bot'): void
}>()

const dots = ref('...')
let intervalId: any

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
    <div
        class="flex flex-col items-center justify-center py-12 px-6 w-full h-full max-w-[1600px] mx-auto text-slate-100">

        <!-- Main Card Content -->
        <main class="w-full max-w-xl">
            <div
                class="bg-white dark:bg-surface-dark rounded-[3rem] p-5 md:p-14 shadow-2xl border border-slate-100 dark:border-slate-700 text-center relative overflow-hidden w-full">
                <div class="absolute top-0 right-0 w-32 h-32 bg-secondary/5 rounded-bl-full pointer-events-none"></div>
                <div class="absolute bottom-0 left-0 w-24 h-24 bg-primary/5 rounded-tr-full pointer-events-none"></div>

                <div class="relative mb-12">
                    <div
                        class="relative w-40 h-40 mx-auto bg-slate-100 dark:bg-slate-900/50 rounded-3xl border border-slate-200 dark:border-slate-800 flex items-center justify-center overflow-hidden">
                        <img :src="randomGallery" alt="Loading"
                            class="w-full h-full object-cover opacity-80 animate-float" />
                    </div>
                    <div
                        class="absolute -inset-4 border-2 border-primary/20 rounded-[2.5rem] animate-pulse-slow pointer-events-none">
                    </div>
                </div>

                <h2 class="text-2xl md:text-3xl font-bold text-slate-800 dark:text-white mb-6">
                    Waiting for an opponent{{ dots }}
                </h2>
                <!-- <p class="text-slate-500 dark:text-slate-400 mb-10">
                    Waiting for someone to join your game.
                </p> -->

                <div class="space-y-4 mb-10">
                    <p class="text-sm font-semibold text-slate-600 dark:text-slate-300">
                        Share the room code to play with a friend:
                    </p>
                    <div
                        class="bg-slate-50 dark:bg-slate-800/50 rounded-2xl p-2 pl-8 flex items-center justify-between border border-slate-200 dark:border-slate-600 max-w-sm mx-auto shadow-inner group transition-colors hover:border-primary/50">
                        <span class="truncate text-3xl font-bold text-slate-800 dark:text-white tracking-widest">{{
                            roomCode }}</span>
                        <button @click="copyRoomCode"
                            class="bg-white dark:bg-surface-dark text-slate-600 dark:text-slate-200 px-6 py-3 rounded-xl shadow-sm border border-slate-100 dark:border-slate-600 hover:text-primary hover:border-primary active:scale-95 transition-colors font-bold text-sm flex items-center gap-2">
                            <Copy class="w-4 h-4" />
                            Copy
                        </button>
                    </div>
                </div>

                <button @click="$emit('add-bot')"
                    class="w-full max-w-xs py-4 px-6 rounded-2xl bg-primary text-white shadow-xl shadow-primary/20 font-bold hover:brightness-105 transition-all text-sm mb-4">
                    Play with Bot
                </button>
                <button @click="$emit('cancel')"
                    class="w-full max-w-xs py-3.5 rounded-xl bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400 font-bold hover:bg-slate-200 dark:hover:bg-slate-600 hover:text-slate-700 dark:hover:text-slate-200 transition-all text-sm mb-4">
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
