import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ViewSecret from '../views/ViewSecret.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/s/:key',
      name: 'view',
      component: ViewSecret,
      props: true
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/'
    }
  ]
})

// 添加全局导航守卫用于调试
router.beforeEach((to, from) => {
  console.log('[Router Debug] Navigation to:', to.fullPath)
  console.log('[Router Debug] Route params:', to.params)
  console.log('[Router Debug] Matched routes:', to.matched.map(r => r.path))
  return true
})

export default router