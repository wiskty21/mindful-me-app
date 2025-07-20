<template>
  <div class="screen center safe-area-all" 
       style="background: linear-gradient(135deg, var(--color-primary-500) 0%, var(--color-secondary-600) 100%);">
    <div class="container space-y-8">
      <!-- ロゴ・タイトル -->
      <div class="text-center">
        <div class="w-20 h-20 mx-auto mb-6 center" 
             style="background: rgba(255, 255, 255, 0.15); backdrop-filter: blur(10px); border-radius: var(--radius-full);">
          <svg class="w-10 h-10" style="color: var(--color-neutral-0);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        </div>
        <h1 style="font-size: var(--text-4xl); font-weight: var(--font-bold); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
          MindfulMe
        </h1>
        <p style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-lg); font-weight: var(--font-medium);">
          心の平穏へようこそ
        </p>
      </div>

      <!-- ログインフォーム -->
      <div class="card space-y-6">
        <div>
          <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
            メールアドレス
          </label>
          <input
            v-model="email"
            type="email"
            class="input-field"
            placeholder="your@example.com"
            :disabled="loading"
          />
        </div>

        <div>
          <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
            パスワード
          </label>
          <input
            v-model="password"
            type="password"
            class="input-field"
            placeholder="パスワードを入力"
            :disabled="loading"
          />
        </div>

        <!-- エラーメッセージ -->
        <div v-if="error" class="error-message" style="background: rgba(239, 68, 68, 0.1); border-color: rgba(239, 68, 68, 0.3); color: var(--color-neutral-0);">
          {{ error }}
        </div>

        <!-- デバッグ情報 -->
        <div v-if="debugMode" style="background: rgba(0, 0, 0, 0.3); color: white; padding: var(--space-3); border-radius: var(--radius-lg); font-size: var(--text-xs); margin-bottom: var(--space-4);">
          <div>クリック: {{ clickCount }}</div>
          <div>状態: {{ loading ? 'ローディング' : 'アイドル' }}</div>
          <div>エラー: {{ error || 'なし' }}</div>
        </div>

        <!-- ログインボタン -->
        <button
          @click="handleLogin"
          :disabled="loading"
          class="btn-primary"
          style="width: 100%;"
        >
          <span v-if="loading" class="center">
            <div class="loading-spinner" style="margin-right: var(--space-2);"></div>
            ログイン中...
          </span>
          <span v-else>ログイン</span>
        </button>

        <!-- 登録リンク -->
        <div style="text-align: center;">
          <p style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-sm);">
            アカウントをお持ちでない方は
            <router-link
              to="/register"
              style="color: var(--color-neutral-0); text-decoration: underline; font-weight: var(--font-medium); transition: opacity var(--transition-fast);"
              class="hover:opacity-80"
            >
              新規登録
            </router-link>
          </p>
          
          <!-- デモ・統計リンク -->
          <div style="display: flex; gap: var(--space-4); justify-content: center; margin-top: var(--space-4);">
            <button @click="demoLogin" class="btn-ghost" style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-sm);">
              デモでログイン
            </button>
            <router-link to="/profile" class="btn-ghost" style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-sm);">
              プロフィール
            </router-link>
          </div>
        </div>
      </div>

      <!-- デモアカウント -->
      <div style="text-align: center;">
        <button
          @click="loginDemo"
          :disabled="loading"
          class="btn-secondary"
        >
          デモでお試し
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

// デバッグ用
const debugMode = ref(true) // 実機テスト時はtrueにする
const clickCount = ref(0)

const demoLogin = async () => {
  loading.value = true
  error.value = ''
  
  // 2秒待機してタイマー画面に遷移
  setTimeout(() => {
    loading.value = false
    router.push('/timer')
  }, 2000)
}

const handleLogin = async () => {
  clickCount.value++
  
  // 実機テスト用：バックエンドなしでも動作
  if (debugMode.value) {
    await demoLogin()
    return
  }

  if (!email.value || !password.value) {
    error.value = 'メールアドレスとパスワードを入力してください'
    return
  }

  loading.value = true
  error.value = ''

  const success = await authStore.login(email.value, password.value)
  
  if (success) {
    router.push('/timer')
  } else {
    error.value = authStore.error || 'ログインに失敗しました'
  }
  
  loading.value = false
}

const loginDemo = async () => {
  clickCount.value++
  
  // 実機テスト用：すぐにタイマー画面に遷移
  if (debugMode.value) {
    loading.value = true
    setTimeout(() => {
      loading.value = false
      router.push('/timer')
    }, 1000)
    return
  }
  
  email.value = 'demo@mindfulme.app'
  password.value = 'demo123456'
  await handleLogin()
}
</script>