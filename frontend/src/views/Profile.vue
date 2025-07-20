<template>
  <div class="screen flex-col safe-area-all" 
       style="background: linear-gradient(135deg, var(--color-neutral-50) 0%, var(--color-primary-50) 100%);">
    <!-- ヘッダー -->
    <header style="padding: var(--space-6);">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h1 style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-neutral-800);">
          プロフィール
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
            <p style="color: var(--color-neutral-600);">プロフィールを読み込み中...</p>
          </div>
        </div>

        <!-- プロフィール情報 -->
        <div v-else class="space-y-6">
          <!-- ユーザー基本情報 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              ユーザー情報
            </h2>
            
            <!-- アバター -->
            <div style="text-align: center; margin-bottom: var(--space-6);">
              <div class="center" 
                   style="width: 100px; height: 100px; background: var(--color-primary-500); border-radius: var(--radius-full); margin: 0 auto var(--space-4);">
                <svg style="width: 50px; height: 50px; color: var(--color-neutral-0);" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h3 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800);">
                {{ user.name || 'ユーザー' }}
              </h3>
              <p style="color: var(--color-neutral-600); font-size: var(--text-base);">
                {{ user.email || 'user@example.com' }}
              </p>
            </div>

            <!-- 基本情報フォーム -->
            <div class="space-y-4">
              <div>
                <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-2);">
                  名前
                </label>
                <input
                  v-model="editUser.name"
                  type="text"
                  class="input-field"
                  placeholder="お名前を入力"
                  :disabled="updating"
                />
              </div>

              <div>
                <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-2);">
                  メールアドレス
                </label>
                <input
                  v-model="editUser.email"
                  type="email"
                  class="input-field"
                  placeholder="your@example.com"
                  :disabled="updating"
                />
              </div>
            </div>
          </div>

          <!-- 瞑想設定 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              瞑想設定
            </h2>
            
            <div class="space-y-4">
              <div>
                <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-2);">
                  デフォルト瞑想時間
                </label>
                <select
                  v-model="editUser.preferred_time"
                  class="input-field"
                  :disabled="updating"
                >
                  <option value="5">5分</option>
                  <option value="10">10分</option>
                  <option value="15">15分</option>
                  <option value="20">20分</option>
                  <option value="30">30分</option>
                </select>
              </div>

              <div>
                <label style="display: block; font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--color-neutral-700); margin-bottom: var(--space-2);">
                  リマインダー時刻
                </label>
                <input
                  v-model="editUser.reminder_time"
                  type="time"
                  class="input-field"
                  :disabled="updating"
                />
              </div>
            </div>
          </div>

          <!-- 瞑想統計 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              あなたの瞑想統計
            </h2>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: var(--space-4);">
              <div style="text-align: center; padding: var(--space-4); background: var(--color-primary-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-primary-600);">
                  {{ user.total_sessions || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">総セッション数</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-secondary-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-secondary-600);">
                  {{ user.total_minutes || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">総瞑想時間（分）</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-success-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-success-600);">
                  {{ user.current_streak || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">現在の連続日数</div>
              </div>
              <div style="text-align: center; padding: var(--space-4); background: var(--color-warning-50); border-radius: var(--radius-lg);">
                <div style="font-size: var(--text-2xl); font-weight: var(--font-bold); color: var(--color-warning-600);">
                  {{ user.longest_streak || 0 }}
                </div>
                <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">最長連続日数</div>
              </div>
            </div>
          </div>

          <!-- アプリ設定 -->
          <div class="card-solid">
            <h2 style="font-size: var(--text-xl); font-weight: var(--font-bold); color: var(--color-neutral-800); margin-bottom: var(--space-6);">
              アプリ設定
            </h2>
            
            <div class="space-y-4">
              <!-- 通知設定 -->
              <div style="display: flex; align-items: center; justify-content: space-between; padding: var(--space-4); background: var(--color-neutral-50); border-radius: var(--radius-lg);">
                <div>
                  <div style="font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-800);">
                    プッシュ通知
                  </div>
                  <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">
                    リマインダー通知を受け取る
                  </div>
                </div>
                <button
                  @click="toggleNotifications"
                  :class="[
                    notificationsEnabled ? 'bg-primary-500' : 'bg-neutral-300',
                    'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none'
                  ]"
                  style="transition: background-color var(--transition-fast);"
                >
                  <span
                    :class="[
                      notificationsEnabled ? 'translate-x-5' : 'translate-x-0',
                      'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                    ]"
                    style="transition: transform var(--transition-fast);"
                  ></span>
                </button>
              </div>

              <!-- ダークモード -->
              <div style="display: flex; align-items: center; justify-content: space-between; padding: var(--space-4); background: var(--color-neutral-50); border-radius: var(--radius-lg);">
                <div>
                  <div style="font-size: var(--text-base); font-weight: var(--font-medium); color: var(--color-neutral-800);">
                    ダークモード
                  </div>
                  <div style="font-size: var(--text-sm); color: var(--color-neutral-600);">
                    暗いテーマを使用する
                  </div>
                </div>
                <button
                  @click="toggleDarkMode"
                  :class="[
                    darkModeEnabled ? 'bg-primary-500' : 'bg-neutral-300',
                    'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none'
                  ]"
                  style="transition: background-color var(--transition-fast);"
                >
                  <span
                    :class="[
                      darkModeEnabled ? 'translate-x-5' : 'translate-x-0',
                      'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out'
                    ]"
                    style="transition: transform var(--transition-fast);"
                  ></span>
                </button>
              </div>
            </div>
          </div>

          <!-- 保存ボタン -->
          <div class="space-y-4">
            <button
              @click="saveProfile"
              :disabled="updating"
              class="btn-primary"
              style="width: 100%;"
            >
              <span v-if="updating" class="center">
                <div style="width: 16px; height: 16px; border: 2px solid transparent; border-top: 2px solid white; border-radius: var(--radius-full); animation: spin 1s linear infinite; margin-right: var(--space-2);"></div>
                保存中...
              </span>
              <span v-else>プロフィールを保存</span>
            </button>

            <!-- ログアウトボタン -->
            <button
              @click="handleLogout"
              class="btn-ghost"
              style="width: 100%; color: var(--color-error-600);"
            >
              ログアウト
            </button>
          </div>
        </div>

        <!-- エラー表示 -->
        <div v-if="error" class="card-solid" style="border: 1px solid var(--color-error-200); background: var(--color-error-50);">
          <div style="text-align: center; padding: var(--space-4);">
            <div style="color: var(--color-error-600); font-weight: var(--font-medium); margin-bottom: var(--space-2);">
              エラーが発生しました
            </div>
            <div style="color: var(--color-error-500); font-size: var(--text-sm); margin-bottom: var(--space-4);">
              {{ error }}
            </div>
            <button @click="loadProfile" class="btn-primary">
              再試行
            </button>
          </div>
        </div>

        <!-- 成功メッセージ -->
        <div v-if="successMessage" class="card-solid" style="border: 1px solid var(--color-success-200); background: var(--color-success-50);">
          <div style="text-align: center; padding: var(--space-4); color: var(--color-success-600); font-weight: var(--font-medium);">
            {{ successMessage }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { profileAPI, type User } from '@/utils/api'
import { hapticButtonTap, hapticSuccess } from '@/utils/haptics.js'

const router = useRouter()

const loading = ref(true)
const updating = ref(false)
const error = ref('')
const successMessage = ref('')

const user = ref<User>({
  id: '',
  email: '',
  name: '',
  total_sessions: 0,
  total_minutes: 0,
  current_streak: 0,
  longest_streak: 0,
  preferred_time: 10,
  reminder_time: '08:00',
  created_at: '',
  updated_at: ''
})

const editUser = reactive({
  name: '',
  email: '',
  preferred_time: 10,
  reminder_time: '08:00'
})

const notificationsEnabled = ref(true)
const darkModeEnabled = ref(false)

const loadProfile = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await profileAPI.getProfile()
    user.value = response.data
    
    // 編集フォームに反映
    editUser.name = user.value.name || ''
    editUser.email = user.value.email || ''
    editUser.preferred_time = user.value.preferred_time || 10
    editUser.reminder_time = user.value.reminder_time || '08:00'
    
    console.log('プロフィール取得:', user.value)
  } catch (err: any) {
    console.error('プロフィール取得エラー:', err)
    error.value = err.response?.data?.error || 'プロフィールの読み込みに失敗しました'
    
    // デモデータで代替
    user.value = {
      id: 'demo-user',
      email: 'demo@mindfulme.app',
      name: 'デモユーザー',
      total_sessions: 12,
      total_minutes: 135,
      current_streak: 5,
      longest_streak: 14,
      preferred_time: 10,
      reminder_time: '08:00',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    editUser.name = user.value.name
    editUser.email = user.value.email
    editUser.preferred_time = user.value.preferred_time
    editUser.reminder_time = user.value.reminder_time
  } finally {
    loading.value = false
  }
}

const saveProfile = async () => {
  updating.value = true
  error.value = ''
  successMessage.value = ''
  
  try {
    await profileAPI.updateProfile(editUser)
    user.value = { ...user.value, ...editUser }
    
    successMessage.value = 'プロフィールを更新しました'
    hapticSuccess()
    
    // 3秒後にメッセージを消す
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
    
    console.log('プロフィール更新完了')
  } catch (err: any) {
    console.error('プロフィール更新エラー:', err)
    error.value = err.response?.data?.error || 'プロフィールの更新に失敗しました'
  } finally {
    updating.value = false
  }
}

const toggleNotifications = () => {
  notificationsEnabled.value = !notificationsEnabled.value
  hapticButtonTap()
  console.log(`通知設定: ${notificationsEnabled.value ? 'ON' : 'OFF'}`)
}

const toggleDarkMode = () => {
  darkModeEnabled.value = !darkModeEnabled.value
  hapticButtonTap()
  console.log(`ダークモード: ${darkModeEnabled.value ? 'ON' : 'OFF'}`)
}

const handleLogout = () => {
  hapticButtonTap()
  // 将来的にログアウト処理を実装
  router.push('/')
  console.log('ログアウト')
}

const goBack = () => {
  hapticButtonTap()
  router.push('/')
}

onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.input-field {
  width: 100%;
  padding: var(--space-3);
  border: 1px solid var(--color-neutral-300);
  border-radius: var(--radius-lg);
  font-size: var(--text-base);
  background: var(--color-neutral-0);
  transition: border-color var(--transition-fast);
}

.input-field:focus {
  outline: none;
  border-color: var(--color-primary-500);
  box-shadow: 0 0 0 2px var(--color-primary-100);
}

.input-field:disabled {
  background: var(--color-neutral-100);
  color: var(--color-neutral-500);
}
</style>