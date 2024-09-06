<script setup lang="ts">
import { useDraggable } from '@vueuse/core';
import { MessageSquare } from 'lucide-vue-next';
import { onMounted, onUnmounted, ref } from 'vue';
import { UseDraggable } from "@vueuse/components";
const emit = defineEmits<{
  (e: 'toggle'): void;
}>();
const buttonPosition = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })
const savePosition = () => {
  localStorage.setItem('buttonPosition', JSON.stringify(buttonPosition.value));
};
const startDrag = (event: TouchEvent | MouseEvent) => {
  isDragging.value = true
  const touch = event instanceof TouchEvent ? event.touches[0] : event;
  dragOffset.value = {
    x: touch.clientX - buttonPosition.value.x,
    y: touch.clientY - buttonPosition.value.y
  }
}

const onDrag = (event: TouchEvent | MouseEvent) => {
  if (isDragging.value) {
    const touch = event instanceof TouchEvent ? event.touches[0] : event;
    buttonPosition.value = {
      x: touch.clientX - dragOffset.value.x,
      y: touch.clientY - dragOffset.value.y
    }
  }
}

const endDrag = () => {
  isDragging.value = false
  savePosition();
}

onMounted(() => {
  const savedPosition = localStorage.getItem('buttonPosition');
  if (savedPosition) {
    buttonPosition.value = JSON.parse(savedPosition); // Apply saved position
  }
  document.addEventListener('touchmove', onDrag, { passive: false })
  document.addEventListener('touchend', endDrag)

  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', endDrag);
  document.addEventListener('mousedown', endDrag);
  document.addEventListener('mouseenter', endDrag);
  document.addEventListener('mouseover', endDrag);
})

onUnmounted(() => {
  document.removeEventListener('touchmove', onDrag)
  document.removeEventListener('touchend', endDrag)

  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', endDrag);

  document.removeEventListener('mousedown', endDrag);
  document.removeEventListener('mouseenter', endDrag);
  document.removeEventListener('mouseover', endDrag);
})


const el = ref<HTMLElement | null>(null)

// `style` will be a helper computed for `left: ?px; top: ?px;`
const { x, y, style } = useDraggable(el, {
  initialValue: { x: 100, y: 100 },
})
</script>
<!-- class="bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900"> -->

<template>
  <button ref="el" @click="$emit('toggle')" :style="style" style="touch-action: none;" @touchstart.capture="startDrag"
    class="fixed bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 outline-none">
    <MessageSquare class="w-6 h-6" />
  </button>
  <button @click="$emit('toggle')"
    class="hidden md:block absolute bottom-10 right-10 bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 outline-none">
    <MessageSquare class="w-6 h-6" />
  </button>
  <!-- <UseDraggable ref="el" v-slot="{ x, y, style }" storage-key="buttonPosition" storage-type="local">
    <button @click="$emit('toggle')" :style="style" style="touch-action: none;"
      class="absolute bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 outline-none">
      <MessageSquare class="w-6 h-6" />
    </button>
  </UseDraggable> -->
</template>