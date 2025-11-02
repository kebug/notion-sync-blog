import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/Home.vue'),
    },
    {
      path: '/projects',
      name: 'Projects',
      component: () => import('../views/Projects.vue'),
    },
  ],
})

export default router

