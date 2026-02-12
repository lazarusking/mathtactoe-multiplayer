<script setup lang="ts">
import { MessageSquare } from 'lucide-vue-next';
import { onMounted, onUnmounted, ref, useTemplateRef } from 'vue';

const emit = defineEmits<{
  (e: 'toggle'): void;
}>();

const windowWidth = ref(window.innerWidth)
const windowHeight = ref(window.innerHeight)
const buttonPosition = ref({ x: windowWidth.value - 20, y: windowHeight.value - 20 })
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

const savePosition = () => {
  localStorage.setItem('buttonPosition', JSON.stringify(buttonPosition.value));
};

const startDrag = (event: TouchEvent | MouseEvent) => {
  isDragging.value = true;
  const point =
    typeof TouchEvent !== 'undefined' && event instanceof TouchEvent
      ? event.touches[0]
      : (event as MouseEvent);

  if (!point) return;

  dragOffset.value = {
    x: point.clientX - buttonPosition.value.x,
    y: point.clientY - buttonPosition.value.y
  };
};

const onDrag = (event: TouchEvent | MouseEvent) => {
  if (isDragging.value) {
    const point =
      typeof TouchEvent !== 'undefined' && event instanceof TouchEvent
        ? event.touches[0]
        : (event as MouseEvent);

    if (!point) return;

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

const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === '/') {
    // Don't toggle if the user is typing in an input or textarea
    const target = event.target as HTMLElement;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
      return;
    }
    event.preventDefault();
    emit('toggle');
  } else if (event.key === 'Escape') {
    emit('toggle');
  }
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

  // Keyboard shortcut listener
  window.addEventListener('keydown', handleKeyDown);
});

onUnmounted(() => {
  // Remove touch event listeners
  document.removeEventListener('touchmove', onDrag);
  document.removeEventListener('touchend', endDrag);

  // Remove mouse event listeners
  document.removeEventListener('mousemove', onDrag);
  document.removeEventListener('mouseup', endDrag);

  // Remove keyboard shortcut listener
  window.removeEventListener('keydown', handleKeyDown);
});

const el = useTemplateRef('el')

// `style` will be a helper computed for `left: ?px; top: ?px;`
// const { x, y, style, position } = useDraggable(el, {
//   initialValue: { x: 100, y: 100 },
// })
// console.log(style.value, position.value)

</script>

<template>
  <button aria-label="floating chat button" ref="el" @click="$emit('toggle')" :style="buttonPosition.x !== 0 && buttonPosition.y !== 0
    ? { left: `${buttonPosition.x - 100}px`, top: `${buttonPosition.y - 100}px`, position: 'fixed', touchAction: 'none' }
    : { right: `${buttonPosition.x}px`, top: `${buttonPosition.y}px`, position: 'fixed' }"
    @touchstart.capture="startDrag" @mousedown.capture="startDrag"
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
