<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 ios-safe-area">
    <div class="text-center px-6">
      <div class="w-32 h-32 bg-gray-200 rounded-full mx-auto mb-8 flex items-center justify-center">
        <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6-4h6m2 5.291A7.962 7.962 0 0112 15c-2.25 0-4.27.69-5.936 1.872M3 21v-1a13.036 13.036 0 0118 0v1" />
        </svg>
      </div>
      
      <h1 class="text-6xl font-bold text-gray-900 mb-4">404</h1>
      <h2 class="text-2xl font-semibold text-gray-700 mb-4">ページが見つかりません</h2>
      <p class="text-gray-600 mb-8 max-w-md mx-auto">
        お探しのページは存在しないか、移動された可能性があります。
      </p>
      
      <div class="space-y-4">
        <button
          @click="$router.push('/')"
          class="w-full max-w-xs btn-meditation touch-manipulation"
        >
          ホームに戻る
        </button>
        
        <button
          @click="$router.go(-1)"
          class="w-full max-w-xs btn-secondary touch-manipulation"
        >
          前のページに戻る
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 404ページは静的なので特別なロジックは不要
</script>