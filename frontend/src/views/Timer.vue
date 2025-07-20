<template>
  <div class="screen flex-col safe-area-all" 
       style="background: linear-gradient(135deg, var(--color-neutral-50) 0%, var(--color-primary-50) 100%);">
    <!-- ヘッダー -->
    <header style="padding: var(--space-6);">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-neutral-800);">
          瞑想タイマー
        </h1>
        <div style="display: flex; gap: var(--space-2);">
          <button @click="showSettings" class="btn-ghost">
            <svg style="width: 20px; height: 20px; margin-right: var(--space-1);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            設定
          </button>
          <button @click="goBack" class="btn-ghost">
            <svg style="width: 20px; height: 20px; margin-right: var(--space-1);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            戻る
          </button>
        </div>
      </div>
    </header>

    <!-- メインコンテンツ -->
    <div class="flex-1 center" style="padding: 0 var(--space-6);">
      <div class="container" style="text-align: center;">
        <!-- タイマー表示 -->
        <div style="position: relative; display: inline-block; margin-bottom: var(--space-12);">
          <!-- メインタイマーサークル -->
          <div class="center" 
               style="width: 280px; height: 280px; background: var(--color-neutral-0); border-radius: var(--radius-full); box-shadow: var(--shadow-xl); position: relative;">
            <!-- 内側のタイマー表示 -->
            <div style="text-align: center; z-index: 2; position: relative;">
              <div style="font-size: 3.5rem; font-family: var(--font-mono); font-weight: var(--font-bold); color: var(--color-neutral-800); line-height: 1;">
                {{ formatTime(currentTime) }}
              </div>
              <div style="font-size: var(--text-base); color: var(--color-neutral-500); margin-top: var(--space-2); font-weight: var(--font-medium);">
                {{ getTimerStatus() }}
              </div>
            </div>
            
            <!-- 呼吸ガイド（瞑想中のみ表示） -->
            <div v-if="isRunning" 
                 class="breathing-guide"
                 style="position: absolute; inset: 20px; border-radius: var(--radius-full); background: linear-gradient(135deg, var(--color-primary-100), var(--color-secondary-100)); opacity: 0.7;">
            </div>
          </div>
          
          <!-- 進行状況サークル -->
          <svg style="position: absolute; top: 0; left: 0; width: 280px; height: 280px; transform: rotate(-90deg);" v-if="selectedTime > 0">
            <!-- 背景サークル -->
            <circle
              cx="140"
              cy="140"
              r="130"
              stroke="var(--color-neutral-200)"
              stroke-width="12"
              fill="none"
            />
            <!-- プログレスサークル -->
            <circle
              cx="140"
              cy="140"
              r="130"
              stroke="url(#progressGradient)"
              stroke-width="12"
              fill="none"
              stroke-linecap="round"
              :stroke-dasharray="circumference"
              :stroke-dashoffset="strokeDashoffset"
              style="transition: stroke-dashoffset var(--transition-slow);"
            />
            <!-- グラデーション定義 -->
            <defs>
              <linearGradient id="progressGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" style="stop-color: var(--color-primary-500);" />
                <stop offset="100%" style="stop-color: var(--color-secondary-500);" />
              </linearGradient>
            </defs>
          </svg>
        </div>

        <!-- 時間選択 -->
        <div v-if="!isRunning && currentTime === selectedTime * 60" style="margin-bottom: var(--space-12);">
          <div style="text-align: center; margin-bottom: var(--space-6);">
            <label style="display: block; font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-700);">
              瞑想時間を選択
            </label>
          </div>
          <div style="display: flex; gap: var(--space-3); justify-content: center; flex-wrap: wrap;">
            <button
              v-for="time in timeOptions"
              :key="time"
              @click="setTime(time)"
              :class="[
                selectedTime === time ? 'btn-primary' : 'btn-secondary-solid'
              ]"
              style="min-width: 60px;"
            >
              {{ time }}分
            </button>
          </div>
        </div>

        <!-- コントロールボタン -->
        <div style="display: flex; gap: var(--space-4); justify-content: center; flex-wrap: wrap;">
          <button
            v-if="!isRunning"
            @click="startTimer"
            class="btn-success"
            style="min-width: 120px;"
          >
            <svg style="width: 20px; height: 20px; margin-right: var(--space-2);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ currentTime === selectedTime * 60 ? '開始' : '再開' }}
          </button>
          
          <button
            v-if="isRunning"
            @click="pauseTimer"
            class="btn-warning"
            style="min-width: 120px;"
          >
            <svg style="width: 20px; height: 20px; margin-right: var(--space-2);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            一時停止
          </button>
          
          <button
            v-if="currentTime !== selectedTime * 60"
            @click="resetTimer"
            class="btn-ghost"
            style="min-width: 100px;"
          >
            <svg style="width: 20px; height: 20px; margin-right: var(--space-2);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            リセット
          </button>
        </div>
      </div>
    </div>

    <!-- 完了モーダル -->
    <div v-if="showCompleteModal" 
         style="position: fixed; inset: 0; background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(4px); z-index: var(--z-modal);"
         class="center">
      <div class="card-solid" style="max-width: 400px; margin: var(--space-4); text-align: center;">
        <!-- 成功アイコン -->
        <div class="center" 
             style="width: 80px; height: 80px; background: var(--color-success-50); border-radius: var(--radius-full); margin: 0 auto var(--space-6);">
          <svg style="width: 40px; height: 40px; color: var(--color-success-500);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        
        <!-- メッセージ -->
        <h3 style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-4);">
          瞑想完了！
        </h3>
        <p style="color: var(--color-neutral-600); font-size: var(--text-base); margin-bottom: var(--space-8);">
          {{ selectedTime }}分間の瞑想<br>お疲れさまでした
        </p>
        
        <!-- ボタン -->
        <div class="space-y-4">
          <button
            @click="closeModal"
            class="btn-primary"
            style="width: 100%;"
          >
            閉じる
          </button>
        </div>
      </div>
    </div>

    <!-- 設定 Bottom Sheet -->
    <BottomSheet 
      v-model="showSettingsSheet" 
      title="タイマー設定"
      height="medium"
    >
      <div class="space-y-6">
        <!-- 環境音設定 -->
        <div>
          <label style="display: block; font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-3);">
            環境音
          </label>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-3);">
            <button
              v-for="sound in soundOptions"
              :key="sound.id"
              @click="selectSound(sound)"
              :class="[
                selectedSound === sound.id ? 'btn-primary' : 'btn-secondary-solid'
              ]"
              style="padding: var(--space-3);"
            >
              <div style="text-align: center;">
                <div style="font-size: var(--text-lg); margin-bottom: var(--space-1);">{{ sound.icon }}</div>
                <div style="font-size: var(--text-sm);">{{ sound.name }}</div>
              </div>
            </button>
          </div>
        </div>

        <!-- 通知設定 -->
        <div>
          <label style="display: block; font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-3);">
            通知設定
          </label>
          <div style="display: flex; align-items: center; justify-content: space-between; padding: var(--space-4); background: var(--color-neutral-50); border-radius: var(--radius-lg);">
            <span style="font-size: var(--text-base); color: var(--color-neutral-700);">
              瞑想完了時に振動
            </span>
            <button
              @click="toggleVibration"
              :class="[
                vibrationEnabled ? 'bg-primary-500' : 'bg-neutral-300',
                'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none'
              ]"
              style="transition: background-color var(--transition-fast);"
            >
              <span
                :class="[
                  vibrationEnabled ? 'translate-x-5' : 'translate-x-0',
                  'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                ]"
                style="transition: transform var(--transition-fast);"
              ></span>
            </button>
          </div>
        </div>

        <!-- 統計情報 -->
        <div>
          <h4 style="font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-3);">
            今日の統計
          </h4>
          <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-4);">
            <div style="text-align: center; padding: var(--space-4); background: var(--color-primary-50); border-radius: var(--radius-lg);">
              <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-primary-600);">3</div>
              <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">セッション</div>
            </div>
            <div style="text-align: center; padding: var(--space-4); background: var(--color-secondary-50); border-radius: var(--radius-lg);">
              <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-secondary-600);">45</div>
              <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">分</div>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <button
          @click="closeSettings"
          class="btn-primary"
          style="width: 100%;"
        >
          閉じる
        </button>
      </template>
    </BottomSheet>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { hapticButtonTap, hapticMediumImpact, hapticHeavyImpact, hapticSuccess, hapticSelection } from '@/utils/haptics.js'
import BottomSheet from '@/components/BottomSheet.vue'

const router = useRouter()

const timeOptions = [5, 10, 15, 20, 30]
const selectedTime = ref(10) // デフォルト10分
const currentTime = ref(selectedTime.value * 60) // 秒単位
const isRunning = ref(false)
const showCompleteModal = ref(false)

// Bottom Sheet 設定
const showSettingsSheet = ref(false)
const vibrationEnabled = ref(true)
const selectedSound = ref('none')

// 環境音オプション
const soundOptions = [
  { id: 'none', name: '無音', icon: '🔇' },
  { id: 'rain', name: '雨音', icon: '🌧️' },
  { id: 'forest', name: '森林', icon: '🌲' },
  { id: 'ocean', name: '海波', icon: '🌊' }
]

let timer: number | null = null

// 新しいサークルサイズに合わせて更新
const circumference = 2 * Math.PI * 130

const strokeDashoffset = computed(() => {
  if (selectedTime.value === 0) return circumference
  const progress = (selectedTime.value * 60 - currentTime.value) / (selectedTime.value * 60)
  return circumference * (1 - progress)
})

const getTimerStatus = () => {
  if (isRunning.value) return '瞑想中'
  if (currentTime.value === selectedTime.value * 60) return '準備完了'
  return '一時停止'
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

const setTime = (minutes: number) => {
  selectedTime.value = minutes
  currentTime.value = minutes * 60
  hapticSelection() // 時間選択時の軽い振動
}

const startTimer = () => {
  isRunning.value = true
  hapticHeavyImpact() // 瞑想開始の重要アクション
  
  timer = setInterval(() => {
    if (currentTime.value > 0) {
      currentTime.value--
    } else {
      completeTimer()
    }
  }, 1000)
}

const pauseTimer = () => {
  isRunning.value = false
  hapticMediumImpact() // 一時停止の中程度フィードバック
  
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

const resetTimer = () => {
  pauseTimer()
  currentTime.value = selectedTime.value * 60
  hapticButtonTap() // リセットボタンの軽いフィードバック
}

const completeTimer = () => {
  pauseTimer()
  showCompleteModal.value = true
  hapticSuccess() // 瞑想完了の成功フィードバック
  // 将来的にセッション記録をバックエンドに送信
}

const closeModal = () => {
  showCompleteModal.value = false
  resetTimer()
  hapticButtonTap() // モーダル閉じるボタン
}

const goBack = () => {
  hapticButtonTap() // 戻るボタン
  router.push('/')
}

// Bottom Sheet 関連
const showSettings = () => {
  showSettingsSheet.value = true
  hapticButtonTap()
}

const closeSettings = () => {
  showSettingsSheet.value = false
  hapticButtonTap()
}

const selectSound = (sound: any) => {
  selectedSound.value = sound.id
  hapticSelection()
  // 将来的に音声ファイルをプリロードまたは再生テスト
  console.log(`環境音選択: ${sound.name}`)
}

const toggleVibration = () => {
  vibrationEnabled.value = !vibrationEnabled.value
  hapticButtonTap()
  console.log(`振動設定: ${vibrationEnabled.value ? 'ON' : 'OFF'}`)
}

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>