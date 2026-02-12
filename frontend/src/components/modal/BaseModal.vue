<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
defineProps<{ showHelp: boolean, baseClass?: string }>()
defineSlots<{ title: string, content: HTMLElement }>()
const emit = defineEmits<{ (event: 'close-modal'): void }>()

const dialog = ref(null);
function handleOutsideClick(event: Event) {
    event.stopPropagation();
    const node = dialog.value;
    if (event.target === node) {
        emit('close-modal');
    }
}

const handleCancel = (event: Event) => {
    event.preventDefault();
    // console.log(event);
    if ((event as KeyboardEvent).key === "Escape") {
        emit('close-modal')
    }
};

onMounted(() => {
    const node = dialog.value;
    if (node) {
        window.addEventListener("keydown", handleCancel);
    }
})


onUnmounted(() => {
    window?.removeEventListener("keydown", handleCancel);
})
</script>

<template>
    <Teleport to="body">
        <!-- <div  class="fixed inset-0 flex items-center justify-center z-50 bg-opacity-50"> -->
        <div v-if="showHelp" @click="handleOutsideClick" ref="dialog"
            class="fixed p-4 inset-0 flex items-center justify-center z-50 overflow-y-auto animate-opacity transition ease-in-out duration-500 bg-background/80 backdrop-blur-sm">
            <div
                :class="['w-auto px-4 py-8 mx-auto bg-card border border-border rounded-4xl sm:px-6 lg:px-8 shadow-2xl', baseClass]">
                <header class="flex font-bold text-center text-lg text-foreground mb-6">
                    <h2 class="flex-1">
                        <slot name="title">
                            <!-- Insert title here -->
                        </slot>
                    </h2>
                    <div class="ml-auto">
                        <button @click="emit('close-modal')" aria-label="Close Button"
                            class="text-muted-foreground hover:text-foreground">
                            <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="1.5"
                                viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                                <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M9.75 9.75l4.5 4.5m0-4.5l-4.5 4.5M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                            </svg>
                        </button>
                    </div>
                </header>
                <slot name="content">
                    <!-- Insert content here -->
                </slot>
            </div>

        </div>
    </Teleport>
</template>
<style scoped>
.modal {
    position: fixed;
    z-index: 999;
    top: 20%;
    left: 50%;
    width: 300px;
    margin-left: -150px;
}
</style>