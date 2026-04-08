/**
 * useTimer — Global task timer with localStorage persistence.
 * Survives page navigation and browser refresh.
 * Uses module-level singleton state so all components share the same reactive timer.
 *
 * Usage:
 *   const { timerState, elapsedMinutes, elapsedDisplay, start, stop, clear } = useTimer()
 */

export interface TimerState {
  taskId: string
  taskTitle: string
  taskCode: string
  startedAt: number // Date.now() ms
}

const TIMER_KEY = 'sentinel:active-timer'

// Module-level singletons — shared across all component instances
const _timerState = ref<TimerState | null>(null)
const _now = ref(Date.now())
let _ticker: ReturnType<typeof setInterval> | null = null
let _initialized = false

function _loadFromStorage() {
  if (!import.meta.client) return
  const raw = localStorage.getItem(TIMER_KEY)
  if (raw) {
    try {
      _timerState.value = JSON.parse(raw) as TimerState
    }
    catch {
      localStorage.removeItem(TIMER_KEY)
    }
  }
}

function _startTicker() {
  if (_ticker) clearInterval(_ticker)
  _ticker = setInterval(() => { _now.value = Date.now() }, 1000)
}

function _stopTicker() {
  if (_ticker) { clearInterval(_ticker); _ticker = null }
}

export function useTimer() {
  // Initialize once on client
  if (import.meta.client && !_initialized) {
    _initialized = true
    _loadFromStorage()
    if (_timerState.value) _startTicker()
  }

  /** Elapsed seconds since timer started (live) */
  const elapsedSeconds = computed(() => {
    if (!_timerState.value) return 0
    return Math.floor((_now.value - _timerState.value.startedAt) / 1000)
  })

  /** Elapsed whole minutes */
  const elapsedMinutes = computed(() => Math.floor(elapsedSeconds.value / 60))

  /** Human-readable "HH:MM:SS" */
  const elapsedDisplay = computed(() => {
    const s = elapsedSeconds.value
    const h = Math.floor(s / 3600)
    const m = Math.floor((s % 3600) / 60)
    const sec = s % 60
    const pad = (n: number) => String(n).padStart(2, '0')
    return `${pad(h)}:${pad(m)}:${pad(sec)}`
  })

  const isRunning = computed(() => !!_timerState.value)

  /** Start timer for a task. Replaces any existing timer. */
  function start(taskId: string, taskTitle: string, taskCode = '') {
    if (!import.meta.client) return
    const state: TimerState = { taskId, taskTitle, taskCode, startedAt: Date.now() }
    _timerState.value = state
    localStorage.setItem(TIMER_KEY, JSON.stringify(state))
    _startTicker()
  }

  /**
   * Stop timer and return elapsed minutes.
   * Caller should open the log form pre-filled with this value.
   */
  function stop(): { minutes: number; taskId: string; taskTitle: string; taskCode: string } | null {
    if (!_timerState.value) return null
    const mins = elapsedMinutes.value
    const result = {
      minutes: Math.max(mins, 1), // at least 1 minute
      taskId: _timerState.value.taskId,
      taskTitle: _timerState.value.taskTitle,
      taskCode: _timerState.value.taskCode,
    }
    clear()
    return result
  }

  function clear() {
    if (!import.meta.client) return
    _timerState.value = null
    localStorage.removeItem(TIMER_KEY)
    _stopTicker()
  }

  return { timerState: _timerState, elapsedMinutes, elapsedDisplay, isRunning, start, stop, clear }
}
