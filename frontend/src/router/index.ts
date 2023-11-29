import HomepageVue from '@/components/Homepage.vue'
import GameViewVue from '@/views/GameView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'game',
      component: HomepageVue
    },
    {
      path: '/:room',
      name: 'room',
      component: GameViewVue
    }
  ]
})

export default router
