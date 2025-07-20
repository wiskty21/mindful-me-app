<template>
  <teleport to="body">
    <!-- バックドロップ -->
    <transition name="backdrop" appear>
      <div
        v-if="modelValue"
        class="bottom-sheet-backdrop"
        @click="handleBackdropClick"
      ></div>
    </transition>

    <!-- Bottom Sheet -->
    <transition name="bottom-sheet" appear>
      <div
        v-if="modelValue"
        ref="sheetElement"
        class="bottom-sheet"
        :class="sheetClasses"
        v-swipe:down="onSwipeDown"
        @touchstart="onTouchStart"
        @touchmove="onTouchMove"
        @touchend="onTouchEnd"
      >
        <!-- Handle（上部のつまみ） -->
        <div class="bottom-sheet-handle" @click="toggle">
          <div class="bottom-sheet-handle-bar"></div>
        </div>

        <!-- Header -->
        <header v-if="title || $slots.header" class="bottom-sheet-header">
          <slot name="header">
            <h3 class="bottom-sheet-title">{{ title }}</h3>
            <button
              v-if="showCloseButton"
              @click="close"
              class="bottom-sheet-close"
              v-haptic="'light'"
            >
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </slot>
        </header>

        <!-- Content -->
        <main class="bottom-sheet-content">
          <slot></slot>
        </main>

        <!-- Footer -->
        <footer v-if="$slots.footer" class="bottom-sheet-footer">
          <slot name="footer"></slot>
        </footer>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import { hapticButtonTap, hapticMediumImpact } from '@/utils/haptics.js'

interface Props {
  modelValue: boolean
  title?: string
  height?: 'small' | 'medium' | 'large' | 'full'
  dismissible?: boolean
  showCloseButton?: boolean
  swipeToClose?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  height: 'medium',
  dismissible: true,
  showCloseButton: true,
  swipeToClose: true
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'close': []
  'open': []
}>()

const sheetElement = ref<HTMLElement>()

// タッチ関連の状態
const touchStartY = ref(0)
const currentTranslateY = ref(0)
const isDragging = ref(false)

const sheetClasses = computed(() => [
  `bottom-sheet--${props.height}`,
  {
    'bottom-sheet--dragging': isDragging.value
  }
])

const open = () => {
  emit('update:modelValue', true)
  emit('open')
  hapticMediumImpact()
}

const close = () => {
  emit('update:modelValue', false)
  emit('close')
  hapticButtonTap()
}

const toggle = () => {
  if (props.modelValue) {
    close()
  } else {
    open()
  }
}

const handleBackdropClick = () => {
  if (props.dismissible) {
    close()
  }
}

// スワイプ処理
const onSwipeDown = (data: any) => {
  if (props.swipeToClose && data.deltaY > 50) {
    close()
  }
}

// タッチドラッグ処理（より細かい制御）
const onTouchStart = (event: TouchEvent) => {
  if (!props.swipeToClose) return
  
  touchStartY.value = event.touches[0].clientY
  isDragging.value = true
  currentTranslateY.value = 0
}

const onTouchMove = (event: TouchEvent) => {
  if (!isDragging.value || !props.swipeToClose) return
  
  const deltaY = event.touches[0].clientY - touchStartY.value
  
  // 下方向のドラッグのみ許可
  if (deltaY > 0) {
    currentTranslateY.value = deltaY
    
    if (sheetElement.value) {
      sheetElement.value.style.transform = `translateY(${deltaY}px)`
      sheetElement.value.style.transition = 'none'
    }
  }
}

const onTouchEnd = (event: TouchEvent) => {
  if (!isDragging.value || !props.swipeToClose) return
  
  isDragging.value = false
  
  if (sheetElement.value) {
    sheetElement.value.style.transition = ''
    sheetElement.value.style.transform = ''
  }
  
  // 150px以上ドラッグした場合は閉じる
  if (currentTranslateY.value > 150) {
    close()
  }
  
  currentTranslateY.value = 0
}

// 外部からの操作用
defineExpose({
  open,
  close,
  toggle
})
</script>

<style scoped>
/* バックドロップ */
.bottom-sheet-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  z-index: var(--z-modal);
}

/* Bottom Sheet */
.bottom-sheet {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--color-neutral-0);
  border-radius: var(--radius-2xl) var(--radius-2xl) 0 0;
  box-shadow: var(--shadow-xl);
  z-index: calc(var(--z-modal) + 1);
  
  /* iOS Safe Area 対応 */
  padding-bottom: max(var(--safe-area-inset-bottom), var(--space-6));
  
  display: flex;
  flex-direction: column;
  
  /* パフォーマンス最適化 */
  will-change: transform;
  backface-visibility: hidden;
  -webkit-backface-visibility: hidden;
}

/* サイズバリエーション */
.bottom-sheet--small {
  max-height: 30vh;
}

.bottom-sheet--medium {
  max-height: 50vh;
}

.bottom-sheet--large {
  max-height: 80vh;
}

.bottom-sheet--full {
  max-height: 95vh;
}

/* ドラッグ中のスタイル */
.bottom-sheet--dragging {
  transition: none !important;
}

/* Handle（つまみ） */
.bottom-sheet-handle {
  display: flex;
  justify-content: center;
  padding: var(--space-3) 0;
  cursor: grab;
  -webkit-tap-highlight-color: transparent;
}

.bottom-sheet-handle:active {
  cursor: grabbing;
}

.bottom-sheet-handle-bar {
  width: 36px;
  height: 4px;
  background: var(--color-neutral-300);
  border-radius: var(--radius-full);
  transition: background-color var(--transition-fast);
}

.bottom-sheet-handle:hover .bottom-sheet-handle-bar {
  background: var(--color-neutral-400);
}

/* Header */
.bottom-sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-6) var(--space-2);
  border-bottom: 1px solid var(--color-neutral-200);
  flex-shrink: 0;
}

.bottom-sheet-title {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--color-neutral-800);
  margin: 0;
}

.bottom-sheet-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--color-neutral-100);
  border: none;
  border-radius: var(--radius-full);
  color: var(--color-neutral-600);
  cursor: pointer;
  transition: all var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
}

.bottom-sheet-close:hover {
  background: var(--color-neutral-200);
  color: var(--color-neutral-700);
}

.bottom-sheet-close:active {
  transform: scale(0.95);
}

/* Content */
.bottom-sheet-content {
  flex: 1;
  padding: var(--space-6);
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

/* Footer */
.bottom-sheet-footer {
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--color-neutral-200);
  flex-shrink: 0;
}

/* アニメーション */
.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity var(--transition-base);
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}

.bottom-sheet-enter-active,
.bottom-sheet-leave-active {
  transition: transform var(--transition-base) cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

.bottom-sheet-enter-from,
.bottom-sheet-leave-to {
  transform: translateY(100%);
}

/* ダークモード対応 */
@media (prefers-color-scheme: dark) {
  .bottom-sheet {
    background: var(--color-neutral-100);
  }
  
  .bottom-sheet-title {
    color: var(--color-neutral-900);
  }
  
  .bottom-sheet-header,
  .bottom-sheet-footer {
    border-color: var(--color-neutral-200);
  }
  
  .bottom-sheet-handle-bar {
    background: var(--color-neutral-400);
  }
  
  .bottom-sheet-close {
    background: var(--color-neutral-200);
    color: var(--color-neutral-700);
  }
  
  .bottom-sheet-close:hover {
    background: var(--color-neutral-300);
    color: var(--color-neutral-800);
  }
}

/* タブレット・デスクトップでの調整 */
@media (min-width: 768px) {
  .bottom-sheet {
    left: 50%;
    right: auto;
    transform: translateX(-50%);
    max-width: 500px;
    border-radius: var(--radius-2xl);
    margin-bottom: var(--space-8);
  }
  
  .bottom-sheet-enter-from,
  .bottom-sheet-leave-to {
    transform: translateX(-50%) translateY(100%);
  }
}
</style>