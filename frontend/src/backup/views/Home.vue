<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- ヘッダー -->
    <header class="ios-top-safe bg-white/80 backdrop-blur-md border-b border-gray-200">
      <div class="px-4 py-4 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 bg-gradient-to-br from-meditation-main to-meditation-dark rounded-full flex items-center justify-center">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
          <div>
            <h1 class="text-lg font-bold text-gray-900">MindfulMe</h1>
            <p class="text-sm text-gray-600">こんにちは、{{ user?.name }}さん</p>
          </div>
        </div>
        <button
          @click="$router.push('/profile')"
          class="p-2 rounded-full bg-gray-100 touch-manipulation"
        >
          <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </button>
      </div>
    </header>

    <!-- メインコンテンツ -->
    <main class="px-4 py-6 space-y-6">
      <!-- 今日の統計 -->
      <div class="card">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">今日の瞑想</h2>
        <div class="grid grid-cols-3 gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ todayStats.sessions }}</div>
            <div class="text-sm text-gray-600">セッション</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ todayStats.minutes }}</div>
            <div class="text-sm text-gray-600">分</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-meditation-main">{{ user?.currentStreak || 0 }}</div>
            <div class="text-sm text-gray-600">連続日</div>
          </div>
        </div>
      </div>

      <!-- クイックスタート -->
      <div class="card">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">瞑想を始める</h2>
        <div class="grid grid-cols-3 gap-3">
          <button
            v-for="time in quickTimes"
            :key="time"
            @click="startMeditation(time)"
            class="bg-gradient-to-br from-meditation-main to-meditation-dark text-white py-4 px-3 rounded-xl font-medium touch-manipulation hover:from-meditation-dark hover:to-meditation-dark transition-all duration-200"
          >
            {{ time }}分
          </button>
        </div>
        <button
          @click="$router.push('/timer')"
          class="w-full mt-4 btn-secondary touch-manipulation"
        >
          カスタム時間を設定
        </button>
      </div>

      <!-- 週間進捗 -->
      <div class="card">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">週間進捗</h2>
        <div class="flex justify-between items-end h-24 mb-4">
          <div
            v-for="(day, index) in weeklyProgress"
            :key="index"
            class="flex flex-col items-center space-y-2"
          >
            <div
              class="w-6 bg-meditation-main rounded-t-sm transition-all duration-300"
              :style="{ height: `${(day.minutes / maxWeeklyMinutes) * 80}px` }"
            ></div>
            <span class="text-xs text-gray-600">{{ day.day }}</span>
          </div>
        </div>
        <div class="text-center text-sm text-gray-600">
          今週の合計: {{ weeklyTotal }}分
        </div>
      </div>

      <!-- 瞑想のヒント -->
      <div class="card bg-gradient-to-r from-meditation-light to-white">
        <h2 class="text-lg font-semibold text-gray-900 mb-3">今日のヒント</h2>
        <p class="text-gray-700 leading-relaxed">
          {{ dailyTip }}
        </p>
      </div>
    </main>

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
const quickTimes = [5, 10, 15]

// 今日の統計（仮データ）
const todayStats = ref({
  sessions: 2,
  minutes: 25
})

// 週間進捗（仮データ）
const weeklyProgress = ref([
  { day: '月', minutes: 15 },
  { day: '火', minutes: 30 },
  { day: '水', minutes: 20 },
  { day: '木', minutes: 25 },
  { day: '金', minutes: 10 },
  { day: '土', minutes: 35 },
  { day: '日', minutes: 0 }
])

const maxWeeklyMinutes = computed(() => 
  Math.max(...weeklyProgress.value.map(d => d.minutes), 1)
)

const weeklyTotal = computed(() =>
  weeklyProgress.value.reduce((total, day) => total + day.minutes, 0)
)

// 瞑想のヒント
const tips = [
  '呼吸に意識を向けることから始めましょう。深くゆっくりとした呼吸がリラックスを促します。',
  '思考が浮かんできても判断せず、優しく呼吸に意識を戻しましょう。',
  '毎日少しずつでも続けることが、瞑想の効果を高める鍵です。',
  '座り心地の良い姿勢を見つけて、背筋を伸ばしましょう。',
  '瞑想に「正しい」や「間違い」はありません。今この瞬間を大切にしましょう。'
]

const dailyTip = computed(() => {
  const today = new Date().getDate()
  return tips[today % tips.length]
})

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

const startMeditation = (minutes: number) => {
  router.push(`/timer?duration=${minutes}`)
}

onMounted(() => {
  // 統計データの取得
  // TODO: API呼び出しで実際のデータを取得
})
</script>