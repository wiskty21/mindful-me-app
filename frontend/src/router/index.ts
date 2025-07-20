import { createRouter, createWebHistory } from 'vue-router'
import Login from '@/views/Login.vue'
import Register from '@/views/Register.vue'
import Timer from '@/views/Timer.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Login',
      component: Login,
      meta: { 
        index: 0,
        title: 'ログイン'
      }
    },
    {
      path: '/register',
      name: 'Register',
      component: Register,
      meta: { 
        index: 1,
        title: '新規登録'
      }
    },
    {
      path: '/timer',
      name: 'Timer',
      component: Timer,
      meta: { 
        index: 2,
        title: '瞑想タイマー'
      }
    }
  ]
})

// iOS風ページ遷移方向の計算
router.beforeEach((to, from) => {
  const toIndex = to.meta?.index ?? 0
  const fromIndex = from.meta?.index ?? 0
  
  // 遷移方向を動的に設定（forward: 右から左、backward: 左から右）
  to.meta.transitionDirection = toIndex > fromIndex ? 'forward' : 'backward'
})

export { router }