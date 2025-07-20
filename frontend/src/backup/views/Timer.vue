<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <!-- ヘッダー -->
    <header class="ios-top-safe px-4 py-4 flex items-center justify-between">
      <button @click="goBack" class="p-2 rounded-full bg-white/80 touch-manipulation">
        <svg class="w-6 h-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h1 class="text-lg font-semibold text-gray-900">瞑想タイマー</h1>
      <div class="w-10"></div>
    </header>

    <!-- タイマー設定画面 -->
    <div v-if="!isActive && !isPaused" class="px-4 py-8 space-y-8">
      <!-- 時間選択 -->
      <div class="text-center">
        <h2 class="text-2xl font-bold text-gray-900 mb-8">瞑想時間を選択</h2>
        
        <!-- プリセット時間 -->
        <div class="grid grid-cols-3 gap-4 mb-8">
          <button
            v-for="preset in presetTimes"
            :key="preset"
            @click="setDuration(preset)"
            :class="[
              'py-6 px-4 rounded-2xl font-semibold text-lg touch-manipulation transition-all duration-200',
              duration === preset * 60 
                ? 'bg-meditation-main text-white shadow-lg' 
                : 'bg-white text-gray-700 hover:bg-gray-50'
            ]"
          >
            {{ preset }}分
          </button>
        </div>

        <!-- カスタム時間スライダー -->
        <div class="card">
          <h3 class="text-lg font-medium text-gray-900 mb-6">カスタム時間</h3>
          <div class="space-y-6">
            <div class="text-center">
              <div class="text-4xl font-bold text-meditation-main mb-2">
                {{ Math.floor(duration / 60) }}分{{ duration % 60 !== 0 ? ` ${duration % 60}秒` : '' }}
              </div>
            </div>
            <input
              v-model="duration"
              type="range"
              min="60"
              max="3600"
              step="30"
              class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
            <div class="flex justify-between text-sm text-gray-500">
              <span>1分</span>
              <span>60分</span>
            </div>
          </div>
        </div>

        <!-- 瞑想タイプ選択 -->
        <div class="card">
          <h3 class="text-lg font-medium text-gray-900 mb-4">瞑想タイプ</h3>
          <div class="grid grid-cols-1 gap-3">
            <button
              v-for="type in meditationTypes"
              :key="type.value"
              @click="selectedType = type.value"
              :class="[
                'p-4 rounded-xl text-left touch-manipulation transition-all duration-200',
                selectedType === type.value
                  ? 'bg-meditation-main text-white'
                  : 'bg-gray-50 hover:bg-gray-100'
              ]"
            >
              <div class="font-medium">{{ type.name }}</div>
              <div class="text-sm opacity-80">{{ type.description }}</div>
            </button>
          </div>
        </div>

        <!-- 開始ボタン -->
        <button
          @click="startMeditation"
          class="w-full btn-meditation touch-manipulation"
        >
          瞑想を開始
        </button>
      </div>
    </div>

    <!-- タイマー実行画面 -->
    <div v-else class="flex flex-col items-center justify-center min-h-screen p-8">
      <!-- タイマー円 -->
      <div class="relative mb-12">
        <svg class="w-80 h-80 -rotate-90" viewBox="0 0 120 120">
          <!-- 背景円 -->
          <circle
            cx="60"
            cy="60"
            r="54"
            fill="none"
            stroke="#E5E7EB"
            stroke-width="4"
          />
          <!-- 進捗円 -->
          <circle
            cx="60"
            cy="60"
            r="54"
            fill="none"
            stroke="url(#gradient)"
            stroke-width="4"
            stroke-linecap="round"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="strokeDashoffset"
            class="transition-all duration-1000"
          />
          <!-- グラデーション定義 -->
          <defs>
            <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" style="stop-color:#667EEA;stop-opacity:1" />
              <stop offset="100%" style="stop-color:#5A67D8;stop-opacity:1" />
            </linearGradient>
          </defs>
        </svg>
        
        <!-- 中央の時間表示 -->
        <div class="absolute inset-0 flex flex-col items-center justify-center">
          <div class="text-5xl font-bold text-gray-900 mb-2 font-mono">
            {{ formatTime(timeLeft) }}
          </div>
          <div class="text-lg text-gray-600">{{ currentTypeLabel }}</div>
        </div>

        <!-- 呼吸ガイド -->
        <div
          v-if="showBreathingGuide"
          :class="[
            'absolute inset-0 rounded-full border-4 border-meditation-main/30 transition-transform duration-4000',
            isInhaling ? 'scale-110' : 'scale-100'
          ]"
        ></div>
      </div>

      <!-- コントロールボタン -->
      <div class="flex space-x-6">
        <button
          @click="pauseResume"
          class="w-16 h-16 rounded-full bg-white shadow-lg flex items-center justify-center touch-manipulation"
        >
          <svg v-if="!isPaused" class="w-8 h-8 text-gray-700" fill="currentColor" viewBox="0 0 24 24">
            <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/>
          </svg>
          <svg v-else class="w-8 h-8 text-gray-700" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </button>
        
        <button
          @click="stopMeditation"
          class="w-16 h-16 rounded-full bg-red-500 shadow-lg flex items-center justify-center touch-manipulation"
        >
          <svg class="w-8 h-8 text-white" fill="currentColor" viewBox="0 0 24 24">
            <path d="M6 6h12v12H6z"/>
          </svg>
        </button>
      </div>

      <!-- 完了画面 -->
      <div
        v-if="isCompleted"
        class="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50"
      >
        <div class="bg-white rounded-2xl p-8 w-full max-w-md text-center">
          <div class="w-20 h-20 bg-green-500 rounded-full mx-auto mb-6 flex items-center justify-center">
            <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-4">瞑想完了！</h3>
          <p class="text-gray-600 mb-6">
            {{ Math.floor(originalDuration / 60) }}分間の瞑想お疲れさまでした
          </p>
          
          <!-- 気分選択 -->
          <div class="mb-6">
            <h4 class="text-lg font-medium text-gray-900 mb-3">今の気分は？</h4>
            <div class="grid grid-cols-2 gap-3">
              <button
                v-for="mood in moods"
                :key="mood.value"
                @click="selectedMood = mood.value"
                :class="[
                  'p-3 rounded-xl touch-manipulation transition-all duration-200',
                  selectedMood === mood.value
                    ? 'bg-meditation-main text-white'
                    : 'bg-gray-100 hover:bg-gray-200'
                ]"
              >
                {{ mood.emoji }} {{ mood.label }}
              </button>
            </div>
          </div>

          <!-- メモ入力 -->
          <div class="mb-6 text-left">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              メモ（任意）
            </label>
            <textarea
              v-model="sessionNote"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg resize-none"
              rows="3"
              placeholder="今回の瞑想についてメモを残しましょう"
            ></textarea>
          </div>

          <div class="flex space-x-3">
            <button
              @click="saveSession"
              class="flex-1 btn-meditation touch-manipulation"
            >
              保存
            </button>
            <button
              @click="skipSave"
              class="flex-1 btn-secondary touch-manipulation"
            >
              スキップ
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

// 状態管理
const duration = ref(600) // 秒
const timeLeft = ref(0)
const originalDuration = ref(0)
const isActive = ref(false)
const isPaused = ref(false)
const isCompleted = ref(false)
const selectedType = ref('timer')
const selectedMood = ref('')
const sessionNote = ref('')
const showBreathingGuide = ref(false)
const isInhaling = ref(true)

let intervalId: number | null = null
let breathingIntervalId: number | null = null

// プリセット時間
const presetTimes = [5, 10, 15, 20, 30]

// 瞑想タイプ
const meditationTypes = [
  {
    value: 'timer',
    name: '静寂瞑想',
    description: '静かな環境で呼吸に集中'
  },
  {
    value: 'breathing',
    name: '呼吸瞑想',
    description: '呼吸リズムをガイド'
  },
  {
    value: 'guided',
    name: 'ガイド瞑想',
    description: '音声ガイド付き（準備中）'
  }
]

// 気分選択肢
const moods = [
  { value: 'calm', label: '穏やか', emoji: '😌' },
  { value: 'focused', label: '集中', emoji: '🎯' },
  { value: 'relaxed', label: 'リラックス', emoji: '😊' },
  { value: 'energized', label: 'エネルギッシュ', emoji: '✨' }
]

// 計算プロパティ
const currentTypeLabel = computed(() => {
  const type = meditationTypes.find(t => t.value === selectedType.value)
  return type?.name || ''
})

const circumference = computed(() => 2 * Math.PI * 54)

const strokeDashoffset = computed(() => {
  const progress = timeLeft.value / originalDuration.value
  return circumference.value * (1 - progress)
})

// メソッド
const setDuration = (minutes: number) => {
  duration.value = minutes * 60
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

const startMeditation = () => {
  originalDuration.value = duration.value
  timeLeft.value = duration.value
  isActive.value = true
  isPaused.value = false
  
  // 呼吸ガイドを有効にする
  if (selectedType.value === 'breathing') {
    showBreathingGuide.value = true
    startBreathingGuide()
  }
  
  startTimer()
}

const startTimer = () => {
  intervalId = setInterval(() => {
    if (timeLeft.value > 0) {
      timeLeft.value--
    } else {
      completeMeditation()
    }
  }, 1000)
}

const pauseResume = () => {
  if (isPaused.value) {
    isPaused.value = false
    startTimer()
    if (showBreathingGuide.value) {
      startBreathingGuide()
    }
  } else {
    isPaused.value = true
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
    if (breathingIntervalId) {
      clearInterval(breathingIntervalId)
      breathingIntervalId = null
    }
  }
}

const stopMeditation = () => {
  isActive.value = false
  isPaused.value = false
  isCompleted.value = false
  timeLeft.value = 0
  showBreathingGuide.value = false
  
  if (intervalId) {
    clearInterval(intervalId)
    intervalId = null
  }
  if (breathingIntervalId) {
    clearInterval(breathingIntervalId)
    breathingIntervalId = null
  }
}

const completeMeditation = () => {
  isActive.value = false
  isCompleted.value = true
  showBreathingGuide.value = false
  
  if (intervalId) {
    clearInterval(intervalId)
    intervalId = null
  }
  if (breathingIntervalId) {
    clearInterval(breathingIntervalId)
    breathingIntervalId = null
  }
  
  // バイブレーション（対応デバイスのみ）
  if ('vibrate' in navigator) {
    navigator.vibrate([500, 200, 500])
  }
}

const startBreathingGuide = () => {
  breathingIntervalId = setInterval(() => {
    isInhaling.value = !isInhaling.value
  }, 4000) // 4秒周期
}

const saveSession = async () => {
  // TODO: API呼び出しでセッションを保存
  const sessionData = {
    duration: Math.floor(originalDuration.value / 60),
    type: selectedType.value,
    mood: selectedMood.value,
    note: sessionNote.value,
    startTime: new Date(Date.now() - originalDuration.value * 1000).toISOString(),
    endTime: new Date().toISOString()
  }
  
  console.log('セッション保存:', sessionData)
  
  isCompleted.value = false
  router.push('/')
}

const skipSave = () => {
  isCompleted.value = false
  router.push('/')
}

const goBack = () => {
  if (isActive.value || isPaused.value) {
    if (confirm('瞑想を中断しますか？')) {
      stopMeditation()
      router.push('/')
    }
  } else {
    router.push('/')
  }
}

// ライフサイクル
onMounted(() => {
  // URLパラメータから初期時間を設定
  const urlDuration = route.query.duration as string
  if (urlDuration) {
    const minutes = parseInt(urlDuration)
    if (!isNaN(minutes) && minutes > 0) {
      setDuration(minutes)
    }
  }
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
  if (breathingIntervalId) {
    clearInterval(breathingIntervalId)
  }
})
</script>

<style scoped>
.slider::-webkit-slider-thumb {
  appearance: none;
  height: 24px;
  width: 24px;
  border-radius: 50%;
  background: #667EEA;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
}

.slider::-moz-range-thumb {
  height: 24px;
  width: 24px;
  border-radius: 50%;
  background: #667EEA;
  cursor: pointer;
  border: none;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
}
</style>