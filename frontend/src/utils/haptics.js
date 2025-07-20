/**
 * iOS風ハプティックフィードバック・タッチフィードバックユーティリティ
 * Web Vibration API を使用してiOSライクな触覚フィードバックを提供
 */

/**
 * ハプティックフィードバックの種類
 */
export const HapticType = {
  // iOS UIImpactFeedbackGenerator equivalent
  LIGHT: 'light',      // 軽いタップ（ボタンタップ、選択）
  MEDIUM: 'medium',    // 中程度のタップ（切り替え、確認）
  HEAVY: 'heavy',      // 重いタップ（重要なアクション、エラー）
  
  // iOS UINotificationFeedbackGenerator equivalent
  SUCCESS: 'success',  // 成功（瞑想完了、タスク完了）
  WARNING: 'warning',  // 警告（注意が必要）
  ERROR: 'error',      // エラー（失敗、問題発生）
  
  // iOS UISelectionFeedbackGenerator equivalent
  SELECTION: 'selection' // 選択変更（リスト項目、値調整）
}

/**
 * 振動パターンの定義（ミリ秒）
 * iPhone実機での体感に近づけるよう調整
 */
const vibrationPatterns = {
  [HapticType.LIGHT]: [10],
  [HapticType.MEDIUM]: [20],
  [HapticType.HEAVY]: [50],
  [HapticType.SUCCESS]: [10, 50, 10, 50],
  [HapticType.WARNING]: [30, 20, 30],
  [HapticType.ERROR]: [50, 30, 50, 30, 50],
  [HapticType.SELECTION]: [5]
}

/**
 * ハプティックフィードバックが利用可能かチェック
 */
export function isHapticsSupported() {
  return 'vibrate' in navigator && navigator.vibrate
}

/**
 * ハプティックフィードバックを実行
 * @param {string} type - HapticType の値
 * @param {boolean} force - フォーカス状態でも強制実行するか
 */
export function triggerHaptic(type = HapticType.LIGHT, force = false) {
  // PWAとしてホーム画面から起動されている場合のみ実行（通常のブラウザでは無効）
  if (!isHapticsSupported()) {
    return false
  }
  
  // ページがフォーカスされていない場合はスキップ（バックグラウンド時の無駄な振動を防ぐ）
  if (!force && document.hidden) {
    return false
  }
  
  const pattern = vibrationPatterns[type] || vibrationPatterns[HapticType.LIGHT]
  
  try {
    navigator.vibrate(pattern)
    return true
  } catch (error) {
    console.warn('Haptic feedback failed:', error)
    return false
  }
}

/**
 * ボタンクリック用の軽いハプティック
 */
export function hapticButtonTap() {
  return triggerHaptic(HapticType.LIGHT)
}

/**
 * 重要なアクション用の中程度ハプティック
 */
export function hapticMediumImpact() {
  return triggerHaptic(HapticType.MEDIUM)
}

/**
 * 瞑想開始/停止などの重要アクション用
 */
export function hapticHeavyImpact() {
  return triggerHaptic(HapticType.HEAVY)
}

/**
 * 成功時のハプティック（瞑想完了など）
 */
export function hapticSuccess() {
  return triggerHaptic(HapticType.SUCCESS)
}

/**
 * エラー時のハプティック
 */
export function hapticError() {
  return triggerHaptic(HapticType.ERROR)
}

/**
 * 選択変更時のハプティック（タイマー時間選択など）
 */
export function hapticSelection() {
  return triggerHaptic(HapticType.SELECTION)
}

/**
 * Vue.js ディレクティブ用のハプティック関数
 * v-haptic="'light'" のように使用
 */
export function createHapticDirective() {
  return {
    mounted(el, binding) {
      const hapticType = binding.value || HapticType.LIGHT
      
      const handleClick = (event) => {
        triggerHaptic(hapticType)
      }
      
      el.addEventListener('click', handleClick)
      el._hapticHandler = handleClick
    },
    
    beforeUnmount(el) {
      if (el._hapticHandler) {
        el.removeEventListener('click', el._hapticHandler)
        delete el._hapticHandler
      }
    }
  }
}

/**
 * iOS設定から振動設定が無効化されているかの推定
 * 実際の設定は取得できないため、テスト振動での推定
 */
export async function checkHapticsEnabled() {
  if (!isHapticsSupported()) {
    return false
  }
  
  // 短時間の軽い振動をテスト実行
  // 設定で無効化されている場合、navigator.vibrate()は成功するが実際には振動しない
  // このため完全な検出は不可能だが、API利用可能性はチェックできる
  try {
    navigator.vibrate(1) // 1msの短い振動
    return true
  } catch (error) {
    return false
  }
}

export default {
  HapticType,
  isHapticsSupported,
  triggerHaptic,
  hapticButtonTap,
  hapticMediumImpact,
  hapticHeavyImpact,
  hapticSuccess,
  hapticError,
  hapticSelection,
  createHapticDirective,
  checkHapticsEnabled
}