import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../layouts/framework/FrameworkLayout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          redirect: '/dashboard'
        },
        {
          path: '/dashboard',
          name: 'dashboard',
          component: () => import('@/views/dashboard/DashboardView.vue')
        },
        {
          path: '/setting',
          name: 'setting',
          component: HomeView
        },
        {
          path: '/schedule',
          name: 'schedule',
          component: () => import('@/views/schedule/ScheduleView.vue')
        },
        {
          path: '/notification',
          name: 'notification',
          component: HomeView
        },
        {
          path: '/:pathMatch(.*)',
          name: 'not-found',
          component: () => import('../views/NotFoundView.vue')
        }
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue')
    },
  ]
})

export default router
