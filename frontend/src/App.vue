<template>
  <div id="app" class="safe-area-all" ref="appElement" v-swipe:left="onSwipeLeft" v-swipe:right="onSwipeRight">
    <router-view v-slot="{ Component, route }">
      <transition 
        :name="getTransitionName(route)"
        mode="out-in"
        @before-enter="onBeforeEnter"
        @after-enter="onAfterEnter"
      >
        <component :is="Component" :key="route.path" />
      </transition>
    </router-view>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      // スワイプ無効化フラグ（アニメーション中などに使用）
      swipeDisabled: false
    }
  },
  mounted() {
    // iOS PWA対応: ステータスバーのスタイルを設定
    this.setIOSMetaTags()
    
    // iOS PWA: viewport-fit=cover の強制適用
    this.setViewportFitCover()
  },
  methods: {
    setIOSMetaTags() {
      // ステータスバーの色を設定（ダークコンテンツ）
      let statusBarMeta = document.querySelector('meta[name="apple-mobile-web-app-status-bar-style"]')
      if (!statusBarMeta) {
        statusBarMeta = document.createElement('meta')
        statusBarMeta.name = 'apple-mobile-web-app-status-bar-style'
        document.head.appendChild(statusBarMeta)
      }
      statusBarMeta.content = 'default'
      
      // ホーム画面に追加時のタイトル
      let titleMeta = document.querySelector('meta[name="apple-mobile-web-app-title"]')
      if (!titleMeta) {
        titleMeta = document.createElement('meta')
        titleMeta.name = 'apple-mobile-web-app-title'
        document.head.appendChild(titleMeta)
      }
      titleMeta.content = 'MindfulMe'
    },
    
    setViewportFitCover() {
      let viewport = document.querySelector('meta[name="viewport"]')
      if (viewport) {
        // viewport-fit=cover を追加してSafe Areaを有効化
        const currentContent = viewport.content
        if (!currentContent.includes('viewport-fit=cover')) {
          viewport.content = currentContent + ', viewport-fit=cover'
        }
      }
    },
    
    // iOS風トランジション名の決定
    getTransitionName(route) {
      const direction = route.meta?.transitionDirection || 'forward'
      return direction === 'forward' ? 'slide-left' : 'slide-right'
    },
    
    // トランジション開始時
    onBeforeEnter() {
      // iOS風ハプティックフィードバック（可能であれば）
      if (navigator.vibrate) {
        navigator.vibrate(10) // 軽い振動
      }
    },
    
    // トランジション完了時
    onAfterEnter() {
      // ページタイトルの更新
      const route = this.$route
      if (route.meta?.title) {
        document.title = `${route.meta.title} - MindfulMe`
      }
      
      // アニメーション完了後はスワイプを再有効化
      this.swipeDisabled = false
    },
    
    // 左スワイプ（次のページへ） - Login → Register → Timer
    onSwipeLeft(data) {
      if (this.swipeDisabled) return
      
      const currentRoute = this.$route
      const currentIndex = currentRoute.meta?.index ?? 0
      
      // 次のページのインデックスを計算
      const nextIndex = currentIndex + 1
      const routes = this.$router.getRoutes()
      const nextRoute = routes.find(route => route.meta?.index === nextIndex)
      
      if (nextRoute) {
        this.swipeDisabled = true // アニメーション中はスワイプ無効
        this.$router.push(nextRoute.path)
      }
    },
    
    // 右スワイプ（前のページへ） - Timer → Register → Login
    onSwipeRight(data) {
      if (this.swipeDisabled) return
      
      const currentRoute = this.$route
      const currentIndex = currentRoute.meta?.index ?? 0
      
      // 前のページのインデックスを計算
      const prevIndex = currentIndex - 1
      const routes = this.$router.getRoutes()
      const prevRoute = routes.find(route => route.meta?.index === prevIndex)
      
      if (prevRoute) {
        this.swipeDisabled = true // アニメーション中はスワイプ無効
        this.$router.push(prevRoute.path)
      }
    }
  }
}
</script>

<style>
/* グローバルスタイルはdesign-tokens.cssで管理しているため最小限に */
#app {
  min-height: 100vh;
  min-height: 100dvh; /* Dynamic viewport height for mobile */
  
  /* iOS PWA: Safe Areaに対応した完全な画面利用 */
  position: relative;
  width: 100%;
  overflow-x: hidden;
  
  /* iOS PWA: ステータスバー領域の背景色（グラデーション背景に合わせる） */
  background: linear-gradient(135deg, var(--color-primary-500), var(--color-secondary-500));
}

/* iOS PWA: ダイナミックアイランド・ノッチ対応 */
@supports (padding: max(0px)) {
  #app {
    padding-top: max(var(--safe-area-inset-top), 0px);
    padding-bottom: max(var(--safe-area-inset-bottom), 0px);
    padding-left: max(var(--safe-area-inset-left), 0px);
    padding-right: max(var(--safe-area-inset-right), 0px);
  }
}

/* iOS PWA: ランドスケープモード（横向き）対応 */
@media (orientation: landscape) and (max-height: 430px) {
  #app {
    /* 横向きの時はtop/bottomのsafe-areaを小さくして画面を有効活用 */
    padding-top: max(var(--safe-area-inset-top), var(--space-2));
    padding-bottom: max(var(--safe-area-inset-bottom), var(--space-2));
  }
}

/* === iOS風ページ遷移アニメーション === */

/* 共通設定 */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94); /* iOS easeInOut */
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

/* 前進（右から左へスライド）Login → Register → Timer */
.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0.9;
}

.slide-left-leave-to {
  transform: translateX(-30%);
  opacity: 0.7;
}

/* 後退（左から右へスライド）Timer → Register → Login */
.slide-right-enter-from {
  transform: translateX(-30%);
  opacity: 0.7;
}

.slide-right-leave-to {
  transform: translateX(100%);
  opacity: 0.9;
}

/* iOS風のz-indexスタッキング */
.slide-left-enter-active {
  z-index: 2;
}

.slide-left-leave-active {
  z-index: 1;
}

.slide-right-enter-active {
  z-index: 1;
}

.slide-right-leave-active {
  z-index: 2;
}

/* パフォーマンス最適化 */
.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  will-change: transform, opacity;
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
  transform-style: preserve-3d;
  -webkit-transform-style: preserve-3d;
}

/* iPhone SEなど小さい画面での調整 */
@media (max-width: 375px) {
  .slide-left-enter-active,
  .slide-left-leave-active,
  .slide-right-enter-active,
  .slide-right-leave-active {
    transition-duration: 0.25s; /* 少し高速化 */
  }
}
</style>