import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080'

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// JWT トークンをリクエストヘッダーに追加
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// レスポンスエラーをハンドリング
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      // トークンが無効な場合、ログアウト処理
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/'
    }
    return Promise.reject(error)
  }
)

export interface User {
  id: string
  email: string
  name: string
  total_sessions: number
  total_minutes: number
  current_streak: number
  longest_streak: number
  last_session_at?: string
  preferred_time: number
  reminder_time: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  user: User
  access_token: string
  refresh_token: string
}

export interface Session {
  id: string
  user_id: string
  start_time: string
  end_time: string
  duration: number
  type: string
  mood?: string
  note?: string
  created_at: string
}

// 認証 API
export const authAPI = {
  login: (data: LoginRequest) => 
    api.post<AuthResponse>('/api/v1/auth/login', data),
  
  register: (data: RegisterRequest) => 
    api.post<AuthResponse>('/api/v1/auth/register', data),
  
  refreshToken: () => {
    const refreshToken = localStorage.getItem('refresh_token')
    return api.post<{ access_token: string; refresh_token: string }>('/api/v1/auth/refresh', {
      refresh_token: refreshToken
    })
  },
  
  logout: () => 
    api.post('/api/v1/auth/logout')
}

// セッション API
export const sessionAPI = {
  createSession: (data: Partial<Session>) => 
    api.post<Session>('/api/v1/sessions', data),
  
  getSessions: () => 
    api.get<Session[]>('/api/v1/sessions'),
  
  updateSession: (id: string, data: Partial<Session>) => 
    api.put<Session>(`/api/v1/sessions/${id}`, data),
  
  getStats: () => 
    api.get('/api/v1/sessions/stats'),
  
  getStreak: () => 
    api.get('/api/v1/sessions/streak')
}

// プロフィール API
export const profileAPI = {
  getProfile: () => 
    api.get<User>('/api/v1/profile'),
  
  updateProfile: (data: Partial<User>) => 
    api.put<User>('/api/v1/profile', data)
}