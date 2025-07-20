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
        <p class="text-gray-600 mt-2">新しい瞑想の旅を始めましょう</p>
      </div>

      <!-- 登録フォーム -->
      <div class="card">
        <form @submit.prevent="handleRegister" class="space-y-6">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
              お名前
            </label>
            <input
              id="name"
              v-model="name"
              type="text"
              required
              class="input-field touch-manipulation"
              placeholder="お名前を入力"
              :disabled="loading"
            />
          </div>

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
              minlength="8"
              class="input-field touch-manipulation"
              placeholder="8文字以上のパスワード"
              :disabled="loading"
            />
            <p class="text-xs text-gray-500 mt-2">
              8文字以上で入力してください
            </p>
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              パスワード（確認）
            </label>
            <input
              id="confirmPassword"
              v-model="confirmPassword"
              type="password"
              required
              class="input-field touch-manipulation"
              placeholder="パスワードを再入力"
              :disabled="loading"
            />
          </div>

          <!-- エラーメッセージ -->
          <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4">
            <p class="text-red-600 text-sm">{{ error }}</p>
          </div>

          <!-- 利用規約同意 -->
          <div class="flex items-start">
            <input
              id="terms"
              v-model="acceptTerms"
              type="checkbox"
              required
              class="mt-1 h-4 w-4 text-meditation-main focus:ring-meditation-main border-gray-300 rounded"
              :disabled="loading"
            />
            <label for="terms" class="ml-3 text-sm text-gray-600">
              <span class="font-medium text-meditation-main cursor-pointer">利用規約</span>
              および
              <span class="font-medium text-meditation-main cursor-pointer">プライバシーポリシー</span>
              に同意します
            </label>
          </div>

          <!-- 登録ボタン -->
          <button
            type="submit"
            :disabled="loading || !isFormValid"
            class="w-full btn-meditation touch-manipulation disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登録中...
            </span>
            <span v-else>アカウント作成</span>
          </button>

          <!-- ログインリンク -->
          <div class="text-center">
            <p class="text-sm text-gray-600">
              すでにアカウントをお持ちの方は
              <router-link
                to="/login"
                class="font-medium text-meditation-main hover:text-meditation-dark transition-colors"
              >
                ログイン
              </router-link>
            </p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const acceptTerms = ref(false)
const loading = ref(false)
const error = ref('')

const isFormValid = computed(() => {
  return name.value.length > 0 &&
         email.value.length > 0 &&
         password.value.length >= 8 &&
         password.value === confirmPassword.value &&
         acceptTerms.value
})

const handleRegister = async () => {
  if (!isFormValid.value) {
    error.value = 'すべての項目を正しく入力してください'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'パスワードが一致しません'
    return
  }

  loading.value = true
  error.value = ''

  const success = await authStore.register(email.value, password.value, name.value)
  
  if (success) {
    router.push('/')
  } else {
    error.value = authStore.error || '登録に失敗しました'
  }
  
  loading.value = false
}
</script>