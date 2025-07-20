import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3003,
    host: '0.0.0.0', // ネットワーク経由でのアクセスを許可
    strictPort: true
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  }
})