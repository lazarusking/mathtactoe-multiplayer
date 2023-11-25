import GameModalVue from '@/components/GameModal.vue'
import { createRouter, createWebHistory } from 'vue-router'
import GameViewVue from '@/views/GameView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'game',
      component: GameModalVue
    },
    {
      path: '/:room',
      name: 'room',
      component: GameViewVue
    }
  ]
})

export default router
