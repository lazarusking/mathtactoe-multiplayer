import HomePageVue from '@/components/HomePage.vue'
import GameView from '@/views/GameView.vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'game',
      component: HomePageVue
    },
    {
      path: '/:room',
      name: 'room',
      component: GameView
    }
  ]
})

export default router
