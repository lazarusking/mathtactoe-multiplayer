<script setup lang="ts">
import { MessageSquare } from 'lucide-vue-next';
import { onMounted, onUnmounted, ref } from 'vue';
defineEmits<{
  (e: 'toggle'): void;
}>();
const buttonPosition = ref({ x: 0, y: 0 })
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })
const savePosition = () => {
  localStorage.setItem('buttonPosition', JSON.stringify(buttonPosition.value));
};
const startDrag = (event: TouchEvent | MouseEvent) => {
  isDragging.value = true;
  let point;
  if (typeof TouchEvent !== 'undefined' && event instanceof TouchEvent) {
    point = event.touches[0];
  } else {
    point = event as MouseEvent;
  }

  dragOffset.value = {
    x: point.clientX - buttonPosition.value.x,
    y: point.clientY - buttonPosition.value.y
  };
};

const onDrag = (event: TouchEvent | MouseEvent) => {
  if (isDragging.value) {
    let point;
    if (typeof TouchEvent !== 'undefined' && event instanceof TouchEvent) {
      point = event.touches[0];
    } else {
      point = event as MouseEvent;
    }

    buttonPosition.value = {
      x: point.clientX - dragOffset.value.x,
      y: point.clientY - dragOffset.value.y
    };
  }
};

const endDrag = () => {
  isDragging.value = false;
  savePosition();
};

onMounted(() => {
  const savedPosition = localStorage.getItem('buttonPosition');
  if (savedPosition) {
    buttonPosition.value = JSON.parse(savedPosition);
  }

  // Touch event listeners
  document.addEventListener('touchmove', onDrag, { passive: false });
  document.addEventListener('touchend', endDrag);

  // Mouse event listeners
  document.addEventListener('mousemove', onDrag);
  document.addEventListener('mouseup', endDrag);
});

onUnmounted(() => {
  // Remove touch event listeners
  document.removeEventListener('touchmove', onDrag);
  document.removeEventListener('touchend', endDrag);

  // Remove mouse event listeners
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', endDrag);
});


const el = ref<HTMLElement | null>(null)

// `style` will be a helper computed for `left: ?px; top: ?px;`
// const { x, y, style, position } = useDraggable(el, {
//   initialValue: { x: 100, y: 100 },
// })
// console.log(style.value, position.value)
</script>
<!-- class="bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900"> -->

<template>
  <button ref="el" @click="$emit('toggle')" :style="{
    // position: 'absolute',
    left: `${buttonPosition.x}px`,
    top: `${buttonPosition.y}px`,
    touchAction: 'none'
  }" style="touch-action: none;" @touchstart.capture="startDrag" @mousedown.capture="startDrag"
    class="fixed md:hidden block bg-blue-600 text-white rounded-full p-3 shadow-lg hover:bg-blue-700 outline-none">
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