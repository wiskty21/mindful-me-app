<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- ヘッダー -->
    <header class="ios-top-safe px-4 py-4 flex items-center justify-between bg-white/80 backdrop-blur-md">
      <h1 class="text-xl font-bold text-gray-900">プロフィール</h1>
      <button
        @click="handleLogout"
        class="p-2 rounded-full bg-red-100 text-red-600 touch-manipulation"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
    </header>

    <div class="px-4 py-6 pb-20 space-y-6">
      <!-- ユーザー情報 -->
      <div class="card text-center">
        <div class="w-24 h-24 bg-gradient-to-br from-meditation-main to-meditation-dark rounded-full mx-auto mb-4 flex items-center justify-center">
          <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        <h2 class="text-xl font-bold text-gray-900 mb-1">{{ user?.name }}</h2>
        <p class="text-gray-600">{{ user?.email }}</p>
        <div class="grid grid-cols-3 gap-4 mt-6 pt-6 border-t border-gray-200">
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ user?.totalSessions || 0 }}</div>
            <div class="text-sm text-gray-600">セッション</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ user?.totalMinutes || 0 }}</div>
            <div class="text-sm text-gray-600">分</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ user?.longestStreak || 0 }}</div>
            <div class="text-sm text-gray-600">最長連続</div>
          </div>
        </div>
      </div>

      <!-- 瞑想設定 -->
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">瞑想設定</h3>
        
        <div class="space-y-4">
          <!-- デフォルト瞑想時間 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              デフォルト瞑想時間
            </label>
            <select
              v-model="preferredTime"
              @change="updatePreferences"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-meditation-main focus:border-transparent"
            >
              <option value="5">5分</option>
              <option value="10">10分</option>
              <option value="15">15分</option>
              <option value="20">20分</option>
              <option value="30">30分</option>
            </select>
          </div>

          <!-- リマインダー時刻 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              リマインダー時刻
            </label>
            <input
              v-model="reminderTime"
              @change="updatePreferences"
              type="time"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-meditation-main focus:border-transparent"
            />
          </div>

          <!-- リマインダー有効/無効 -->
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700">リマインダー通知</span>
            <button
              @click="toggleReminder"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none',
                reminderEnabled ? 'bg-meditation-main' : 'bg-gray-200'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  reminderEnabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
        </div>
      </div>

      <!-- アプリ設定 -->
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">アプリ設定</h3>
        
        <div class="space-y-4">
          <!-- ダークモード -->
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700">ダークモード</span>
            <button
              @click="toggleDarkMode"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none',
                darkMode ? 'bg-meditation-main' : 'bg-gray-200'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  darkMode ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>

          <!-- バイブレーション -->
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700">バイブレーション</span>
            <button
              @click="toggleVibration"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none',
                vibrationEnabled ? 'bg-meditation-main' : 'bg-gray-200'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  vibrationEnabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>

          <!-- 音響効果 -->
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700">音響効果</span>
            <button
              @click="toggleSound"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none',
                soundEnabled ? 'bg-meditation-main' : 'bg-gray-200'
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  soundEnabled ? 'translate-x-6' : 'translate-x-1'
                ]"
              />
            </button>
          </div>
        </div>
      </div>

      <!-- データとプライバシー -->
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">データとプライバシー</h3>
        
        <div class="space-y-3">
          <button
            @click="exportData"
            class="w-full flex items-center justify-between p-3 bg-gray-50 rounded-lg touch-manipulation hover:bg-gray-100 transition-colors"
          >
            <span class="text-sm font-medium text-gray-700">データをエクスポート</span>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <button
            @click="deleteAccount"
            class="w-full flex items-center justify-between p-3 bg-red-50 rounded-lg touch-manipulation hover:bg-red-100 transition-colors"
          >
            <span class="text-sm font-medium text-red-600">アカウントを削除</span>
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>

      <!-- サポート -->
      <div class="card">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">サポート</h3>
        
        <div class="space-y-3">
          <button
            @click="openHelp"
            class="w-full flex items-center justify-between p-3 bg-gray-50 rounded-lg touch-manipulation hover:bg-gray-100 transition-colors"
          >
            <span class="text-sm font-medium text-gray-700">ヘルプ</span>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <button
            @click="sendFeedback"
            class="w-full flex items-center justify-between p-3 bg-gray-50 rounded-lg touch-manipulation hover:bg-gray-100 transition-colors"
          >
            <span class="text-sm font-medium text-gray-700">フィードバック</span>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <button
            @click="showAbout"
            class="w-full flex items-center justify-between p-3 bg-gray-50 rounded-lg touch-manipulation hover:bg-gray-100 transition-colors"
          >
            <span class="text-sm font-medium text-gray-700">アプリについて</span>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>

      <!-- バージョン情報 -->
      <div class="text-center text-sm text-gray-500">
        <p>MindfulMe v1.0.0</p>
        <p>© 2025 MindfulMe. All rights reserved.</p>
      </div>
    </div>

    <!-- ナビゲーションバー -->
    <nav class="fixed bottom-0 left-0 right-0 ios-bottom-safe bg-white border-t border-gray-200">
      <div class="grid grid-cols-4 gap-1">
        <button
          v-for="item in navItems"
          :key="item.name"
          @click="$router.push(item.path)"
          :class="[
            'flex flex-col items-center justify-center py-3 px-2 touch-manipulation transition-colors',
            $route.path === item.path ? 'text-meditation-main' : 'text-gray-500'
          ]"
        >
          <BaseIcon :name="item.icon" class="mb-1" />
          <span class="text-xs">{{ item.name }}</span>
        </button>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import BaseIcon from '@/components/BaseIcon.vue'

const router = useRouter()
const authStore = useAuthStore()

const user = computed(() => authStore.user)

// 設定値
const preferredTime = ref(10)
const reminderTime = ref('08:00')
const reminderEnabled = ref(true)
const darkMode = ref(false)
const vibrationEnabled = ref(true)
const soundEnabled = ref(true)

// メソッド
const updatePreferences = async () => {
  // TODO: API呼び出しで設定を更新
  const success = await authStore.updateProfile({
    preferredTime: preferredTime.value,
    reminderTime: reminderTime.value
  })
  
  if (success) {
    console.log('設定を更新しました')
  }
}

const toggleReminder = () => {
  reminderEnabled.value = !reminderEnabled.value
  // TODO: 通知設定の更新
}

const toggleDarkMode = () => {
  darkMode.value = !darkMode.value
  // TODO: ダークモード設定の更新
}

const toggleVibration = () => {
  vibrationEnabled.value = !vibrationEnabled.value
  // TODO: バイブレーション設定の更新
}

const toggleSound = () => {
  soundEnabled.value = !soundEnabled.value
  // TODO: 音響設定の更新
}

const exportData = () => {
  // TODO: データエクスポート機能
  alert('データエクスポート機能は準備中です')
}

const deleteAccount = () => {
  if (confirm('アカウントを削除すると、すべてのデータが失われます。本当に削除しますか？')) {
    if (confirm('この操作は取り消せません。本当によろしいですか？')) {
      // TODO: アカウント削除API呼び出し
      alert('アカウント削除機能は準備中です')
    }
  }
}

const openHelp = () => {
  // TODO: ヘルプページを開く
  alert('ヘルプページは準備中です')
}

const sendFeedback = () => {
  // TODO: フィードバック送信機能
  const email = 'feedback@mindfulme.app'
  const subject = 'MindfulMe フィードバック'
  const body = 'フィードバック内容をここに記入してください。'
  
  window.location.href = `mailto:${email}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`
}

const showAbout = () => {
  alert('MindfulMe v1.0.0\n\n瞑想・マインドフルネスタイマーアプリ\n\n心の平穏と集中力向上のために開発されました。')
}

const handleLogout = async () => {
  if (confirm('ログアウトしますか？')) {
    await authStore.logout()
    router.push('/login')
  }
}

// ナビゲーション
const navItems = [
  {
    name: 'ホーム',
    path: '/',
    icon: 'home' as const
  },
  {
    name: 'タイマー',
    path: '/timer',
    icon: 'timer' as const
  },
  {
    name: '履歴',
    path: '/history',
    icon: 'history' as const
  },
  {
    name: '設定',
    path: '/profile',
    icon: 'profile' as const
  }
]

onMounted(() => {
  // ユーザー設定の取得
  if (user.value) {
    preferredTime.value = user.value.preferredTime || 10
    reminderTime.value = user.value.reminderTime || '08:00'
  }
})
</script>