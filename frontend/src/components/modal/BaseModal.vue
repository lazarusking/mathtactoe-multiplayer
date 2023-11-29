<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
defineProps<{ showHelp: boolean }>()
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
        <div v-show="showHelp" @click="handleOutsideClick" ref="dialog"
            class="fixed inset-0 flex items-center justify-center z-50 overflow-y-auto animate-opacity transition ease-in-out duration-500 bg-gray-500 bg-opacity-75 focus:bg-slate-900">
            <div class="w-auto max-w-md px-4 py-8 mx-auto bg-gray-800 rounded-lg sm:w-80 md:w-full sm:px-6 lg:px-8">
                <header className="flex font-bold text-center text-lg">
                    <h2 className="flex-1">
                        <slot name="title">
                            <!-- Insert title here -->
                        </slot>
                    </h2>
                    <div className="ml-auto">
                        <button @click="emit('close-modal')" aria-label="Close Button">
                            <svg class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 24 24"
                                xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
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