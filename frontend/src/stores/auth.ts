import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI, type User } from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref('')

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const accessToken = computed(() => localStorage.getItem('access_token'))

  // Actions
  const login = async (email: string, password: string): Promise<boolean> => {
    loading.value = true
    error.value = ''

    try {
      const response = await authAPI.login({ email, password })
      const { user: userData, access_token, refresh_token } = response.data

      // ユーザー情報を保存
      user.value = userData
      
      // トークンをローカルストレージに保存
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      console.log('Login successful:', userData)
      return true
    } catch (err: any) {
      console.error('Login error:', err)
      error.value = err.response?.data?.error || 'ログインに失敗しました'
      return false
    } finally {
      loading.value = false
    }
  }

  const register = async (email: string, password: string, name: string): Promise<boolean> => {
    loading.value = true
    error.value = ''

    try {
      const response = await authAPI.register({ email, password, name })
      const { user: userData, access_token, refresh_token } = response.data

      // ユーザー情報を保存
      user.value = userData
      
      // トークンをローカルストレージに保存
      localStorage.setItem('access_token', access_token)
      localStorage.setItem('refresh_token', refresh_token)

      console.log('Registration successful:', userData)
      return true
    } catch (err: any) {
      console.error('Registration error:', err)
      error.value = err.response?.data?.error || '登録に失敗しました'
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    try {
      await authAPI.logout()
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      // ローカル状態をクリア
      user.value = null
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    }
  }

  const checkAuth = async (): Promise<boolean> => {
    const token = localStorage.getItem('access_token')
    if (!token) {
      return false
    }

    try {
      // 簡易実装：トークンがあれば認証済みとみなす
      return true
    } catch (err) {
      // トークンが無効な場合はクリア
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      return false
    }
  }

  // 初期化時に認証状態確認
  const initAuth = () => {
    const token = localStorage.getItem('access_token')
    if (token) {
      console.log('Found existing token, user authenticated')
    }
  }

  return {
    // State
    user,
    loading,
    error,
    
    // Getters
    isAuthenticated,
    accessToken,
    
    // Actions
    login,
    register,
    logout,
    checkAuth,
    initAuth
  }
})