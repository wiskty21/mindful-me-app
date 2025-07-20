/**
 * iOS風スワイプジェスチャー検出ユーティリティ
 * Touch Events を使用してiOSライクなスワイプ操作を提供
 */

/**
 * スワイプ方向の定義
 */
export const SwipeDirection = {
  LEFT: 'left',
  RIGHT: 'right',
  UP: 'up',
  DOWN: 'down'
}

/**
 * スワイプ設定のデフォルト値
 */
const defaultOptions = {
  threshold: 50,        // スワイプと判定する最小移動距離（px）
  restraint: 100,       // 垂直方向の最大許容移動距離（px）
  allowedTime: 300,     // スワイプと判定する最大時間（ms）
  velocity: 0.3,        // 最小スワイプ速度（px/ms）
  preventDefault: true,  // デフォルトのタッチ動作を防ぐ
  passive: false        // パッシブリスナーを使用しない
}

/**
 * タッチ位置とタイミングを記録するクラス
 */
class TouchData {
  constructor() {
    this.startX = 0
    this.startY = 0
    this.endX = 0
    this.endY = 0
    this.startTime = 0
    this.endTime = 0
  }

  start(touch) {
    this.startX = touch.clientX
    this.startY = touch.clientY
    this.startTime = Date.now()
  }

  end(touch) {
    this.endX = touch.clientX
    this.endY = touch.clientY
    this.endTime = Date.now()
  }

  get deltaX() {
    return this.endX - this.startX
  }

  get deltaY() {
    return this.endY - this.startY
  }

  get duration() {
    return this.endTime - this.startTime
  }

  get velocity() {
    return Math.abs(this.deltaX) / this.duration
  }
}

/**
 * スワイプジェスチャー検出器クラス
 */
export class SwipeDetector {
  constructor(element, options = {}) {
    this.element = element
    this.options = { ...defaultOptions, ...options }
    this.touchData = new TouchData()
    this.handlers = new Map()
    
    this.onTouchStart = this.onTouchStart.bind(this)
    this.onTouchMove = this.onTouchMove.bind(this)
    this.onTouchEnd = this.onTouchEnd.bind(this)
    
    this.init()
  }

  init() {
    const passive = this.options.passive
    
    this.element.addEventListener('touchstart', this.onTouchStart, { passive })
    this.element.addEventListener('touchmove', this.onTouchMove, { passive })
    this.element.addEventListener('touchend', this.onTouchEnd, { passive })
  }

  destroy() {
    this.element.removeEventListener('touchstart', this.onTouchStart)
    this.element.removeEventListener('touchmove', this.onTouchMove)
    this.element.removeEventListener('touchend', this.onTouchEnd)
    this.handlers.clear()
  }

  onTouchStart(event) {
    if (this.options.preventDefault) {
      event.preventDefault()
    }
    
    const touch = event.touches[0]
    this.touchData.start(touch)
  }

  onTouchMove(event) {
    if (this.options.preventDefault) {
      event.preventDefault()
    }
  }

  onTouchEnd(event) {
    if (this.options.preventDefault) {
      event.preventDefault()
    }
    
    const touch = event.changedTouches[0]
    this.touchData.end(touch)
    
    this.detectSwipe()
  }

  detectSwipe() {
    const { deltaX, deltaY, duration, velocity } = this.touchData
    const { threshold, restraint, allowedTime } = this.options
    
    // 時間制限チェック
    if (duration > allowedTime) {
      return
    }
    
    // 速度チェック
    if (velocity < this.options.velocity) {
      return
    }
    
    // 水平スワイプの検出
    if (Math.abs(deltaX) >= threshold && Math.abs(deltaY) <= restraint) {
      const direction = deltaX > 0 ? SwipeDirection.RIGHT : SwipeDirection.LEFT
      this.triggerHandler(direction, { deltaX, deltaY, duration, velocity })
      return
    }
    
    // 垂直スワイプの検出
    if (Math.abs(deltaY) >= threshold && Math.abs(deltaX) <= restraint) {
      const direction = deltaY > 0 ? SwipeDirection.DOWN : SwipeDirection.UP
      this.triggerHandler(direction, { deltaX, deltaY, duration, velocity })
      return
    }
  }

  triggerHandler(direction, data) {
    const handler = this.handlers.get(direction)
    if (handler) {
      handler(data)
    }
    
    // 汎用ハンドラーも実行
    const anyHandler = this.handlers.get('any')
    if (anyHandler) {
      anyHandler(direction, data)
    }
  }

  /**
   * スワイプハンドラーを登録
   * @param {string} direction - SwipeDirection の値または 'any'
   * @param {function} handler - ハンドラー関数
   */
  on(direction, handler) {
    this.handlers.set(direction, handler)
    return this
  }

  /**
   * スワイプハンドラーを削除
   * @param {string} direction - SwipeDirection の値または 'any'
   */
  off(direction) {
    this.handlers.delete(direction)
    return this
  }
}

/**
 * Vue.js コンポーザブル関数
 * @param {Ref<HTMLElement>} elementRef - 対象要素のref
 * @param {object} options - スワイプ設定
 */
export function useSwipe(elementRef, options = {}) {
  let detector = null
  
  const initSwipe = () => {
    if (elementRef.value && !detector) {
      detector = new SwipeDetector(elementRef.value, options)
    }
  }
  
  const destroySwipe = () => {
    if (detector) {
      detector.destroy()
      detector = null
    }
  }
  
  const on = (direction, handler) => {
    if (detector) {
      detector.on(direction, handler)
    }
    return { on, off }
  }
  
  const off = (direction) => {
    if (detector) {
      detector.off(direction)
    }
    return { on, off }
  }
  
  return {
    initSwipe,
    destroySwipe,
    on,
    off
  }
}

/**
 * Vue.js ディレクティブ用のスワイプ関数
 * v-swipe:left="handler" のように使用
 */
export function createSwipeDirective() {
  return {
    mounted(el, binding) {
      const direction = binding.arg || 'any'
      const handler = binding.value
      
      if (typeof handler !== 'function') {
        console.warn('Swipe directive handler must be a function')
        return
      }
      
      const detector = new SwipeDetector(el, binding.modifiers)
      detector.on(direction, handler)
      
      el._swipeDetector = detector
    },
    
    beforeUnmount(el) {
      if (el._swipeDetector) {
        el._swipeDetector.destroy()
        delete el._swipeDetector
      }
    }
  }
}

/**
 * 簡単なスワイプ検出関数（一回限りの使用向け）
 * @param {HTMLElement} element - 対象要素
 * @param {function} onSwipe - スワイプ時のコールバック (direction, data) => void
 * @param {object} options - オプション設定
 */
export function attachSwipeListener(element, onSwipe, options = {}) {
  const detector = new SwipeDetector(element, options)
  detector.on('any', onSwipe)
  
  return () => detector.destroy()
}

export default {
  SwipeDirection,
  SwipeDetector,
  useSwipe,
  createSwipeDirective,
  attachSwipeListener
}