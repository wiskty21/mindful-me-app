<template>
  <div class="screen flex-col safe-area-all" 
       style="background: linear-gradient(135deg, var(--color-neutral-50) 0%, var(--color-primary-50) 100%);">
    <!-- ヘッダー -->
    <header style="padding: var(--space-6);">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-neutral-800);">
          統計・履歴
        </h1>
        <button @click="goBack" class="btn-ghost">
          <svg style="width: 20px; height: 20px; margin-right: var(--space-1);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          戻る
        </button>
      </div>
    </header>

    <!-- メインコンテンツ -->
    <div class="flex-1" style="padding: 0 var(--space-6); overflow-y: auto;">
      <div class="container space-y-8">
        <!-- ローディング表示 -->
        <div v-if="loading" class="center" style="padding: var(--space-12);">
          <div style="text-align: center;">
            <div style="width: 40px; height: 40px; border: 3px solid var(--color-primary-200); border-top: 3px solid var(--color-primary-500); border-radius: var(--radius-full); animation: spin 1s linear infinite; margin: 0 auto var(--space-4);"></div>
            <p style="color: var(--color-neutral-600);">統計データを読み込み中...</p>
          </div>
        </div>

        <!-- 統計サマリー -->
        <div v-else class="space-y-6">
          <!-- メイン統計 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              全体統計
            </h2>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-4);">
              <div style="text-align: center; padding: var(--space-4); background: var(--color-primary-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-3xl); font-weight: var(--font-bold); color: var(--color-primary-600);">
                  {{ stats.total_sessions || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">総セッション数</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-secondary-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-3xl); font-weight: var(--font-bold); color: var(--color-secondary-600);">
                  {{ stats.total_minutes || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">総瞑想時間（分）</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-success-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-3xl); font-weight: var(--font-bold); color: var(--color-success-600);">
                  {{ stats.current_streak || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">現在の連続日数</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-warning-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-3xl); font-weight: var(--font-bold); color: var(--color-warning-600);">
                  {{ stats.longest_streak || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">最長連続日数</div>
              </div>
            </div>
          </div>

          <!-- 今月の統計 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              今月の活動
            </h2>
            <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: var(--space-4);">
              <div style="text-align: center; padding: var(--space-3);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-primary-600);">
                  {{ stats.monthly_sessions || 0 }}
                </div>
                <div style="font-size: var(--text-xs); color: var(--color-neutral-600);">月間セッション</div>
              </div>
              <div style="text-align: center; padding: var(--space-3);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-secondary-600);">
                  {{ stats.weekly_sessions || 0 }}
                </div>
                <div style="font-size: var(--text-xs); color: var(--color-neutral-600);">週間セッション</div>
              </div>
              <div style="text-align: center; padding: var(--space-3);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-success-600);">
                  {{ Math.round(stats.average_duration || 0) }}
                </div>
                <div style="font-size: var(--text-xs); color: var(--color-neutral-600);">平均時間（分）</div>
              </div>
            </div>
          </div>

          <!-- 最近のセッション履歴 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              最近のセッション
            </h2>
            <div v-if="sessions.length === 0" style="text-align: center; padding: var(--space-8); color: var(--color-neutral-500);">
              まだセッションがありません<br>
              瞑想を始めてみましょう
            </div>
            <div v-else class="space-y-4">
              <div 
                v-for="session in sessions.slice(0, 5)" 
                :key="session.id"
                style="display: flex; justify-content: space-between; align-items: center; padding: var(--space-4); background: var(--color-neutral-50); border-radius: var(--radius-lg);"
              >
                <div>
                  <div style="font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-800);">
                    {{ session.duration }}分間の瞑想
                  </div>
                  <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">
                    {{ formatDate(session.created_at) }}
                  </div>
                </div>
                <div style="text-align: right;">
                  <div v-if="session.mood" style="font-size: var(--text-sm); color: var(--color-primary-600); font-weight: var(--font-medium);">
                    {{ getMoodEmoji(session.mood) }} {{ getMoodText(session.mood) }}
                  </div>
                  <div style="font-size: var(--text-xs); color: var(--color-neutral-500);">
                    {{ session.type === 'timer' ? 'タイマー' : 'ガイド' }}
                  </div>
                </div>
              </div>
            </div>
            
            <!-- すべて見るボタン -->
            <div v-if="sessions.length > 5" style="text-align: center; margin-top: var(--space-6);">
              <button @click="showAllSessions" class="btn-secondary-solid">
                すべてのセッションを見る（{{ sessions.length }}）
              </button>
            </div>
          </div>

          <!-- 気分の分析 -->
          <div v-if="stats.most_common_mood" class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              気分の傾向
            </h2>
            <div style="text-align: center; padding: var(--space-6);">
              <div style="font-size: 4rem; margin-bottom: var(--space-4);">
                {{ getMoodEmoji(stats.most_common_mood) }}
              </div>
              <div style="font-size: var(--text-lg); font-weight: var(--font-medium); color: var(--color-neutral-800); margin-bottom: var(--space-2);">
                よく感じる気分: {{ getMoodText(stats.most_common_mood) }}
              </div>
              <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">
                瞑想後に最も多く選ばれている気分です
              </div>
            </div>
          </div>
        </div>

        <!-- エラー表示 -->
        <div v-if="error" class="card-solid" style="border: 1px solid var(--color-error-200); background: var(--color-error-50);">
          <div style="text-align: center; padding: var(--space-4);">
            <div style="color: var(--color-error-600); font-weight: var(--font-medium); margin-bottom: var(--space-2);">
              データの読み込みに失敗しました
            </div>
            <div style="color: var(--color-error-500); font-size: var(--text-sm); margin-bottom: var(--space-4);">
              {{ error }}
            </div>
            <button @click="loadData" class="btn-primary">
              再試行
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { sessionAPI, type Session } from '@/utils/api'
import { hapticButtonTap } from '@/utils/haptics.js'

const router = useRouter()

const loading = ref(true)
const error = ref('')
const stats = ref<any>({})
const sessions = ref<Session[]>([])

const loadData = async () => {
  loading.value = true
  error.value = ''
  
  try {
    // 統計データと履歴を並行取得
    const [statsResponse, sessionsResponse] = await Promise.all([
      sessionAPI.getStats(),
      sessionAPI.getSessions()
    ])
    
    stats.value = statsResponse.data
    sessions.value = sessionsResponse.data
    
    console.log('統計データ:', stats.value)
    console.log('セッション履歴:', sessions.value)
  } catch (err: any) {
    console.error('データ読み込みエラー:', err)
    error.value = err.response?.data?.error || 'データの読み込みに失敗しました'
  } finally {
    loading.value = false
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return '今日'
  } else if (diffDays === 1) {
    return '昨日'
  } else if (diffDays < 7) {
    return `${diffDays}日前`
  } else {
    return date.toLocaleDateString('ja-JP', { 
      month: 'short', 
      day: 'numeric' 
    })
  }
}

const getMoodEmoji = (mood: string) => {
  const moodEmojis: { [key: string]: string } = {
    calm: '😌',
    focused: '🧘',
    relaxed: '😊',
    energized: '✨',
    peaceful: '🕊️'
  }
  return moodEmojis[mood] || '😊'
}

const getMoodText = (mood: string) => {
  const moodTexts: { [key: string]: string } = {
    calm: '穏やか',
    focused: '集中',
    relaxed: 'リラックス',
    energized: 'エネルギッシュ',
    peaceful: '平和'
  }
  return moodTexts[mood] || '良好'
}

const showAllSessions = () => {
  hapticButtonTap()
  // 将来的に詳細履歴ページへ遷移
  console.log('全セッション表示')
}

const goBack = () => {
  hapticButtonTap()
  router.push('/')
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>