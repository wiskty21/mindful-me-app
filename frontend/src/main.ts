import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { router } from './router'
import App from './App.vue'

// デザインシステムのインポート
import './styles/design-tokens.css'
import './styles/components.css'

// iOS風ジェスチャー・ハプティックのインポート
import { createSwipeDirective } from './utils/swipe.js'
import { createHapticDirective } from './utils/haptics.js'

console.log('main.ts starting...')

const app = createApp(App)
const pinia = createPinia()

// ディレクティブ登録
app.directive('swipe', createSwipeDirective())
app.directive('haptic', createHapticDirective())

app.use(pinia)
app.use(router)
console.log('Vue app created with Pinia, router, and iOS directives')

app.mount('#app')
console.log('Vue app mounted')