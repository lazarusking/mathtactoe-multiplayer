<script setup lang="ts">
import BaseModal from './BaseModal.vue';
defineProps<{ showHelp: boolean }>()
defineEmits<{ (event: 'close-modal'): void }>()

</script>
<template>
    <BaseModal @close-modal="$emit('close-modal')" :show-help="showHelp">
        <template #title>How to play?</template>
        <template #content>
            <p class="text-sm sm:text-base font-normal leading-normal">
                The goal of the game is to form a line of three grids in <strong>any</strong> direction (horizontal,
                vertical, or diagonal)
                where the sum of the numbers is equal to <strong>15</strong>.
            </p>
            <section class="grid w-full max-w-md grid-cols-3 grid-rows-3 gap-4 shadow-md h-3/5">
                <template v-for="y in [7, 4, 1, 3, 2, 9, 5, 6, 8]" :key="y">
                    <button type="button"
                        :class="{ 'ring-1 border-2 border-green-700 bg-slate-950': [7, 3, 5].includes(y) }"
                        class="grid items-center justify-center w-auto h-auto p-2 text-4xl font-black text-white transition-colors bg-gray-600 rounded-lg shadow-md place-content-center focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none hover:bg-slate-700">
                        <div v-if="[7, 3, 5].includes(y)">{{ y }}✅</div>
                        <div v-else>{{ y }}</div>
                    </button>
                </template>
            </section>
        </template>
    </BaseModal>
</template>