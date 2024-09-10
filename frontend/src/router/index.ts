import HomePageVue from '@/components/HomePage.vue'
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
      component: import('@/views/GameView.vue')
    }
  ]
})

export default router
