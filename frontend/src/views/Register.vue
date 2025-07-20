<template>
  <div class="screen center safe-area-all" 
       style="background: linear-gradient(135deg, var(--color-primary-500) 0%, var(--color-secondary-600) 100%);">
    <div class="container space-y-8">
      <!-- ロゴ・タイトル -->
      <div style="text-align: center;">
        <div class="center" 
             style="width: 80px; height: 80px; margin: 0 auto var(--space-6); background: rgba(255, 255, 255, 0.15); backdrop-filter: blur(10px); border-radius: var(--radius-full);">
          <svg style="width: 40px; height: 40px; color: var(--color-neutral-0);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        <h1 style="font-size: var(--text-4xl); font-weight: var(--font-bold); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
          新規登録
        </h1>
        <p style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-lg); font-weight: var(--font-medium);">
          瞑想の旅を始めましょう
        </p>
      </div>

      <!-- 登録フォーム -->
      <div class="card space-y-6">
        <div>
          <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
            お名前
          </label>
          <input
            v-model="name"
            type="text"
            class="input-field"
            placeholder="山田 太郎"
            :disabled="loading"
          />
        </div>

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
            placeholder="8文字以上で入力"
            :disabled="loading"
          />
        </div>

        <div>
          <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-0); margin-bottom: var(--space-2);">
            パスワード確認
          </label>
          <input
            v-model="confirmPassword"
            type="password"
            class="input-field"
            placeholder="パスワードを再入力"
            :disabled="loading"
          />
        </div>

        <!-- エラーメッセージ -->
        <div v-if="error" class="error-message" style="background: rgba(239, 68, 68, 0.1); border-color: rgba(239, 68, 68, 0.3); color: var(--color-neutral-0);">
          {{ error }}
        </div>

        <!-- 登録ボタン -->
        <button
          @click="handleRegister"
          :disabled="loading"
          class="btn-primary"
          style="width: 100%;"
        >
          <span v-if="loading" class="center">
            <div class="loading-spinner" style="margin-right: var(--space-2);"></div>
            登録中...
          </span>
          <span v-else>アカウント作成</span>
        </button>

        <!-- ログインリンク -->
        <div style="text-align: center;">
          <p style="color: rgba(255, 255, 255, 0.8); font-size: var(--text-sm);">
            すでにアカウントをお持ちの方は
            <router-link
              to="/"
              style="color: var(--color-neutral-0); text-decoration: underline; font-weight: var(--font-medium); transition: opacity var(--transition-fast);"
              class="hover:opacity-80"
            >
              ログイン
            </router-link>
          </p>
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

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

const handleRegister = async () => {
  if (!name.value || !email.value || !password.value || !confirmPassword.value) {
    error.value = 'すべての項目を入力してください'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'パスワードが一致しません'
    return
  }

  if (password.value.length < 8) {
    error.value = 'パスワードは8文字以上で入力してください'
    return
  }

  loading.value = true
  error.value = ''

  const success = await authStore.register(email.value, password.value, name.value)
  
  if (success) {
    router.push('/timer')
  } else {
    error.value = authStore.error || '登録に失敗しました'
  }
  
  loading.value = false
}
</script>