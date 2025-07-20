<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- ヘッダー -->
    <header class="ios-top-safe px-4 py-4 flex items-center justify-between bg-white/80 backdrop-blur-md">
      <h1 class="text-xl font-bold text-gray-900">瞑想履歴</h1>
      <button
        @click="showStats = !showStats"
        class="p-2 rounded-full bg-gray-100 touch-manipulation"
      >
        <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      </button>
    </header>

    <!-- 統計サマリー -->
    <div v-if="showStats" class="px-4 py-4 space-y-4">
      <div class="grid grid-cols-2 gap-4">
        <div class="card text-center">
          <div class="text-2xl font-bold text-meditation-main">{{ totalSessions }}</div>
          <div class="text-sm text-gray-600">総セッション数</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-meditation-main">{{ totalMinutes }}</div>
          <div class="text-sm text-gray-600">総瞑想時間（分）</div>
        </div>
      </div>
      
      <div class="grid grid-cols-2 gap-4">
        <div class="card text-center">
          <div class="text-2xl font-bold text-meditation-main">{{ currentStreak }}</div>
          <div class="text-sm text-gray-600">現在の連続日数</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-meditation-main">{{ longestStreak }}</div>
          <div class="text-sm text-gray-600">最長連続日数</div>
        </div>
      </div>

      <!-- 月間統計グラフ -->
      <div class="card">
        <h3 class="text-lg font-medium text-gray-900 mb-4">月間セッション数</h3>
        <div class="flex items-end justify-between h-32 mb-4">
          <div
            v-for="(week, index) in monthlyData"
            :key="index"
            class="flex flex-col items-center space-y-2"
          >
            <div
              class="w-8 bg-meditation-main rounded-t-sm transition-all duration-300"
              :style="{ height: `${(week.sessions / maxMonthlySessions) * 100}px` }"
            ></div>
            <span class="text-xs text-gray-600">{{ week.label }}</span>
          </div>
        </div>
        <div class="text-center text-sm text-gray-600">
          今月の合計: {{ monthlyTotal }}セッション
        </div>
      </div>
    </div>

    <!-- フィルターとソート -->
    <div class="px-4 py-4 flex space-x-3">
      <select
        v-model="filterType"
        class="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-white text-sm"
      >
        <option value="">すべてのタイプ</option>
        <option value="timer">静寂瞑想</option>
        <option value="breathing">呼吸瞑想</option>
        <option value="guided">ガイド瞑想</option>
      </select>
      
      <select
        v-model="sortBy"
        class="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-white text-sm"
      >
        <option value="date">日付順</option>
        <option value="duration">時間順</option>
        <option value="mood">気分順</option>
      </select>
    </div>

    <!-- セッション履歴リスト -->
    <div class="px-4 pb-20">
      <div v-if="filteredSessions.length === 0" class="text-center py-12">
        <div class="w-20 h-20 bg-gray-200 rounded-full mx-auto mb-4 flex items-center justify-center">
          <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <p class="text-gray-500">まだ瞑想セッションがありません</p>
        <button
          @click="$router.push('/timer')"
          class="mt-4 btn-meditation touch-manipulation"
        >
          最初の瞑想を始める
        </button>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="session in paginatedSessions"
          :key="session.id"
          class="card hover:shadow-md transition-shadow duration-200"
        >
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center space-x-3">
              <div class="w-12 h-12 bg-meditation-main/10 rounded-full flex items-center justify-center">
                <svg class="w-6 h-6 text-meditation-main" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <div class="font-medium text-gray-900">{{ getTypeLabel(session.type) }}</div>
                <div class="text-sm text-gray-600">{{ formatDate(session.startTime) }}</div>
              </div>
            </div>
            <div class="text-right">
              <div class="text-lg font-semibold text-meditation-main">{{ session.duration }}分</div>
              <div v-if="session.mood" class="text-sm">
                {{ getMoodEmoji(session.mood) }} {{ getMoodLabel(session.mood) }}
              </div>
            </div>
          </div>
          
          <div v-if="session.note" class="bg-gray-50 rounded-lg p-3 mt-3">
            <p class="text-sm text-gray-700">{{ session.note }}</p>
          </div>
          
          <div class="flex justify-between items-center mt-3 text-xs text-gray-500">
            <span>{{ formatTime(session.startTime) }} - {{ formatTime(session.endTime) }}</span>
            <button
              @click="deleteSession(session.id)"
              class="text-red-500 hover:text-red-700 touch-manipulation"
            >
              削除
            </button>
          </div>
        </div>
      </div>

      <!-- ページネーション -->
      <div v-if="totalPages > 1" class="flex justify-center mt-8 space-x-2">
        <button
          v-for="page in totalPages"
          :key="page"
          @click="currentPage = page"
          :class="[
            'w-10 h-10 rounded-full touch-manipulation transition-colors',
            currentPage === page
              ? 'bg-meditation-main text-white'
              : 'bg-white text-gray-700 hover:bg-gray-100'
          ]"
        >
          {{ page }}
        </button>
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
import BaseIcon from '@/components/BaseIcon.vue'

// 状態管理
const showStats = ref(false)
const filterType = ref('')
const sortBy = ref('date')
const currentPage = ref(1)
const sessionsPerPage = 10

// 統計データ（仮データ）
const totalSessions = ref(47)
const totalMinutes = ref(423)
const currentStreak = ref(7)
const longestStreak = ref(21)

// 月間統計（仮データ）
const monthlyData = ref([
  { label: '1週', sessions: 5 },
  { label: '2週', sessions: 8 },
  { label: '3週', sessions: 6 },
  { label: '4週', sessions: 9 }
])

// セッション履歴（仮データ）
const sessions = ref([
  {
    id: 1,
    type: 'timer',
    duration: 15,
    mood: 'calm',
    note: '今日は特に集中できた。呼吸が深くなり、心が落ち着いた。',
    startTime: '2025-07-20T08:00:00Z',
    endTime: '2025-07-20T08:15:00Z'
  },
  {
    id: 2,
    type: 'breathing',
    duration: 10,
    mood: 'focused',
    note: '',
    startTime: '2025-07-19T19:30:00Z',
    endTime: '2025-07-19T19:40:00Z'
  },
  {
    id: 3,
    type: 'timer',
    duration: 20,
    mood: 'relaxed',
    note: '長めの瞑想で深いリラックス状態に入れた。',
    startTime: '2025-07-18T07:00:00Z',
    endTime: '2025-07-18T07:20:00Z'
  }
])

// 計算プロパティ
const maxMonthlySessions = computed(() => 
  Math.max(...monthlyData.value.map(d => d.sessions), 1)
)

const monthlyTotal = computed(() =>
  monthlyData.value.reduce((total, week) => total + week.sessions, 0)
)

const filteredSessions = computed(() => {
  let filtered = sessions.value

  // タイプフィルター
  if (filterType.value) {
    filtered = filtered.filter(session => session.type === filterType.value)
  }

  // ソート
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'date':
        return new Date(b.startTime).getTime() - new Date(a.startTime).getTime()
      case 'duration':
        return b.duration - a.duration
      case 'mood':
        return (a.mood || '').localeCompare(b.mood || '')
      default:
        return 0
    }
  })

  return filtered
})

const totalPages = computed(() => 
  Math.ceil(filteredSessions.value.length / sessionsPerPage)
)

const paginatedSessions = computed(() => {
  const start = (currentPage.value - 1) * sessionsPerPage
  const end = start + sessionsPerPage
  return filteredSessions.value.slice(start, end)
})

// メソッド
const getTypeLabel = (type: string) => {
  const types: Record<string, string> = {
    timer: '静寂瞑想',
    breathing: '呼吸瞑想',
    guided: 'ガイド瞑想'
  }
  return types[type] || type
}

const getMoodLabel = (mood: string) => {
  const moods: Record<string, string> = {
    calm: '穏やか',
    focused: '集中',
    relaxed: 'リラックス',
    energized: 'エネルギッシュ'
  }
  return moods[mood] || mood
}

const getMoodEmoji = (mood: string) => {
  const emojis: Record<string, string> = {
    calm: '😌',
    focused: '🎯',
    relaxed: '😊',
    energized: '✨'
  }
  return emojis[mood] || ''
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const yesterday = new Date(now.getTime() - 24 * 60 * 60 * 1000)
  
  if (date.toDateString() === now.toDateString()) {
    return '今日'
  } else if (date.toDateString() === yesterday.toDateString()) {
    return '昨日'
  } else {
    return date.toLocaleDateString('ja-JP', { month: 'short', day: 'numeric' })
  }
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString('ja-JP', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const deleteSession = async (sessionId: number) => {
  if (confirm('このセッションを削除しますか？')) {
    // TODO: API呼び出しでセッション削除
    sessions.value = sessions.value.filter(s => s.id !== sessionId)
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
  // セッション履歴の取得
  // TODO: API呼び出し
})
</script>