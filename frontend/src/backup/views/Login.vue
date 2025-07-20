<template>
  <div class="min-h-screen flex items-center justify-center ios-safe-area">
    <div class="w-full max-w-md space-y-8 px-6">
      <!-- ロゴ・タイトル -->
      <div class="text-center">
        <div class="w-20 h-20 bg-gradient-to-br from-meditation-main to-meditation-dark rounded-full mx-auto mb-6 flex items-center justify-center">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 font-meditation">MindfulMe</h1>
        <p class="text-gray-600 mt-2">心の平穏へようこそ</p>
      </div>

      <!-- ログインフォーム -->
      <div class="card">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              メールアドレス
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              required
              class="input-field touch-manipulation"
              placeholder="your@example.com"
              :disabled="loading"
            />
          </div>

          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              パスワード
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              required
              class="input-field touch-manipulation"
              placeholder="パスワードを入力"
              :disabled="loading"
            />
          </div>

          <!-- エラーメッセージ -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4">
            <p class="text-red-600 text-sm">{{ error }}</p>
          </div>

          <!-- ログインボタン -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full btn-meditation touch-manipulation disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              ログイン中...
            </span>
            <span v-else>ログイン</span>
          </button>

          <!-- 登録リンク -->
          <div class="text-center">
            <p class="text-sm text-gray-600">
              アカウントをお持ちでない方は
              <router-link
                to="/register"
                class="font-medium text-meditation-main hover:text-meditation-dark transition-colors"
              >
                新規登録
              </router-link>
            </p>
          </div>
        </form>
      </div>

      <!-- デモアカウント -->
      <div class="card bg-gray-50">
        <div class="text-center">
          <h3 class="text-sm font-medium text-gray-700 mb-3">デモアカウント</h3>
          <button
            @click="loginDemo"
            :disabled="loading"
            class="btn-secondary touch-manipulation disabled:opacity-50"
          >
            デモでお試し
          </button>
        </div>
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

const handleLogin = async () => {
  if (!email.value || !password.value) {
    error.value = 'メールアドレスとパスワードを入力してください'
    return
  }

  loading.value = true
  error.value = ''

  const success = await authStore.login(email.value, password.value)
  
  if (success) {
    router.push('/')
  } else {
    error.value = authStore.error || 'ログインに失敗しました'
  }
  
  loading.value = false
}

const loginDemo = async () => {
  email.value = 'demo@mindfulme.app'
  password.value = 'demo123456'
  await handleLogin()
}
</script>