<template>
  <div style="min-height: 100vh; display: flex; flex-direction: column; background: linear-gradient(135deg, #f8fafc 0%, #e0f2fe 100%); padding-top: env(safe-area-inset-top); padding-bottom: env(safe-area-inset-bottom);">
    <!-- ヘッダー -->
    <header style="padding: 24px;">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="font-size: 24px; font-weight: 700; color: #1f2937;">
          統計・履歴
        </h1>
        <button @click="goBack" style="padding: 8px 16px; background: rgba(255, 255, 255, 0.1); border: none; border-radius: 12px; color: #374151; cursor: pointer; display: flex; align-items: center; transition: all 0.2s;">
          <svg style="width: 20px; height: 20px; margin-right: 4px;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          戻る
        </button>
      </div>
    </header>

    <!-- メインコンテンツ -->
    <div style="flex: 1; padding: 0 24px; overflow-y: auto;">
      <div style="max-width: 800px; margin: 0 auto;">
        <!-- ローディング表示 -->
        <div v-if="loading" style="display: flex; justify-content: center; align-items: center; padding: 48px;">
          <div style="text-align: center;">
            <div style="width: 40px; height: 40px; border: 3px solid #93c5fd; border-top: 3px solid #3b82f6; border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto 16px;"></div>
            <p style="color: #6b7280;">統計データを読み込み中...</p>
          </div>
        </div>

        <!-- 統計サマリー -->
        <div v-else style="margin-bottom: 24px;">
          <!-- メイン統計 -->
          <div style="background: white; border-radius: 16px; padding: 24px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); margin-bottom: 24px;">
            <h2 style="font-size: 20px; font-weight: 700; color: #1f2937; margin-bottom: 24px;">
              全体統計
            </h2>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px;">
              <div style="text-align: center; padding: 16px; background: #eff6ff; border-radius: 12px;">
                <div style="font-size: 30px; font-weight: 700; color: #2563eb;">
                  {{ stats.total_sessions || 0 }}
                </div>
                <div style="font-size: 14px; color: #6b7280;">総セッション数</div>
              </div>
              <div style="text-align: center; padding: 16px; background: #f3e8ff; border-radius: 12px;">
                <div style="font-size: 30px; font-weight: 700; color: #7c3aed;">
                  {{ stats.total_minutes || 0 }}
                </div>
                <div style="font-size: 14px; color: #6b7280;">総瞑想時間（分）</div>
              </div>
              <div style="text-align: center; padding: 16px; background: #ecfdf5; border-radius: 12px;">
                <div style="font-size: 30px; font-weight: 700; color: #059669;">
                  {{ stats.current_streak || 0 }}
                </div>
                <div style="font-size: 14px; color: #6b7280;">現在の連続日数</div>
              </div>
              <div style="text-align: center; padding: 16px; background: #fef3c7; border-radius: 12px;">
                <div style="font-size: 30px; font-weight: 700; color: #d97706;">
                  {{ stats.longest_streak || 0 }}
                </div>
                <div style="font-size: 14px; color: #6b7280;">最長連続日数</div>
              </div>
            </div>
          </div>

          <!-- 今月の統計 -->
          <div style="background: white; border-radius: 16px; padding: 24px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); margin-bottom: 24px;">
            <h2 style="font-size: 20px; font-weight: 700; color: #1f2937; margin-bottom: 24px;">
              今月の活動
            </h2>
            <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px;">
              <div style="text-align: center; padding: 12px;">
                <div style="font-size: 24px; font-weight: 700; color: #2563eb;">
                  {{ stats.monthly_sessions || 0 }}
                </div>
                <div style="font-size: 12px; color: #6b7280;">月間セッション</div>
              </div>
              <div style="text-align: center; padding: 12px;">
                <div style="font-size: 24px; font-weight: 700; color: #7c3aed;">
                  {{ stats.weekly_sessions || 0 }}
                </div>
                <div style="font-size: 12px; color: #6b7280;">週間セッション</div>
              </div>
              <div style="text-align: center; padding: 12px;">
                <div style="font-size: 24px; font-weight: 700; color: #059669;">
                  {{ Math.round(stats.average_duration || 0) }}
                </div>
                <div style="font-size: 12px; color: #6b7280;">平均時間（分）</div>
              </div>
            </div>
          </div>

          <!-- 最近のセッション履歴 -->
          <div style="background: white; border-radius: 16px; padding: 24px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); margin-bottom: 24px;">
            <h2 style="font-size: 20px; font-weight: 700; color: #1f2937; margin-bottom: 24px;">
              最近のセッション
            </h2>
            <div v-if="sessions.length === 0" style="text-align: center; padding: 32px; color: #6b7280;">
              まだセッションがありません<br>
              瞑想を始めてみましょう
            </div>
            <div v-else style="display: flex; flex-direction: column; gap: 16px;">
              <div 
                v-for="session in sessions.slice(0, 5)" 
                :key="session.id"
                style="display: flex; justify-content: space-between; align-items: center; padding: 16px; background: #f9fafb; border-radius: 12px;"
              >
                <div>
                  <div style="font-size: 16px; font-weight: 500; color: #1f2937;">
                    {{ session.duration }}分間の瞑想
                  </div>
                  <div style="font-size: 14px; color: #6b7280;">
                    {{ formatDate(session.created_at) }}
                  </div>
                </div>
                <div style="text-align: right;">
                  <div v-if="session.mood" style="font-size: 14px; color: #2563eb; font-weight: 500;">
                    {{ getMoodEmoji(session.mood) }} {{ getMoodText(session.mood) }}
                  </div>
                  <div style="font-size: 12px; color: #9ca3af;">
                    {{ session.type === 'timer' ? 'タイマー' : 'ガイド' }}
                  </div>
                </div>
              </div>
            </div>
            
            <!-- すべて見るボタン -->
            <div v-if="sessions.length > 5" style="text-align: center; margin-top: 24px;">
              <button @click="showAllSessions" style="padding: 12px 24px; background: #f3f4f6; border: none; border-radius: 12px; color: #374151; cursor: pointer; font-weight: 500;">
                すべてのセッションを見る（{{ sessions.length }}）
              </button>
            </div>
          </div>

          <!-- 気分の分析 -->
          <div v-if="stats.most_common_mood" style="background: white; border-radius: 16px; padding: 24px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); margin-bottom: 24px;">
            <h2 style="font-size: 20px; font-weight: 700; color: #1f2937; margin-bottom: 24px;">
              気分の傾向
            </h2>
            <div style="text-align: center; padding: 24px;">
              <div style="font-size: 4rem; margin-bottom: 16px;">
                {{ getMoodEmoji(stats.most_common_mood) }}
              </div>
              <div style="font-size: 18px; font-weight: 500; color: #1f2937; margin-bottom: 8px;">
                よく感じる気分: {{ getMoodText(stats.most_common_mood) }}
              </div>
              <div style="font-size: 14px; color: #6b7280;">
                瞑想後に最も多く選ばれている気分です
              </div>
            </div>
          </div>
        </div>

        <!-- エラー表示 -->
        <div v-if="error" style="background: #fef2f2; border: 1px solid #fecaca; border-radius: 16px; padding: 24px; margin-bottom: 24px;">
          <div style="text-align: center; padding: 16px;">
            <div style="color: #dc2626; font-weight: 500; margin-bottom: 8px;">
              データの読み込みに失敗しました
            </div>
            <div style="color: #ef4444; font-size: 14px; margin-bottom: 16px;">
              {{ error }}
            </div>
            <button @click="loadData" style="padding: 12px 24px; background: #3b82f6; color: white; border: none; border-radius: 12px; cursor: pointer; font-weight: 500;">
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