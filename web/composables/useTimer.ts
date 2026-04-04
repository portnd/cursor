/**
 * useTimer — Global task timer with localStorage persistence.
 * Survives page navigation and browser refresh.
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

export function useTimer() {
  const timerState = ref<TimerState | null>(null)
  const now = ref(Date.now())
  let ticker: ReturnType<typeof setInterval> | null = null

  function loadFromStorage() {
    if (!import.meta.client) return
    const raw = localStorage.getItem(TIMER_KEY)
    if (raw) {
      try {
        timerState.value = JSON.parse(raw) as TimerState
      }
      catch {
        localStorage.removeItem(TIMER_KEY)
      }
    }
  }

  function startTicker() {
    if (ticker) clearInterval(ticker)
    ticker = setInterval(() => { now.value = Date.now() }, 1000)
  }

  function stopTicker() {
    if (ticker) { clearInterval(ticker); ticker = null }
  }

  /** Elapsed seconds since timer started (live) */
  const elapsedSeconds = computed(() => {
    if (!timerState.value) return 0
    return Math.floor((now.value - timerState.value.startedAt) / 1000)
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

  const isRunning = computed(() => !!timerState.value)

  /** Start timer for a task. Replaces any existing timer. */
  function start(taskId: string, taskTitle: string, taskCode = '') {
    if (!import.meta.client) return
    const state: TimerState = { taskId, taskTitle, taskCode, startedAt: Date.now() }
    timerState.value = state
    localStorage.setItem(TIMER_KEY, JSON.stringify(state))
    startTicker()
  }

  /**
   * Stop timer and return elapsed minutes.
   * Caller should open the log form pre-filled with this value.
   */
  function stop(): { minutes: number; taskId: string; taskTitle: string; taskCode: string } | null {
    if (!timerState.value) return null
    const mins = elapsedMinutes.value
    const result = {
      minutes: Math.max(mins, 1), // at least 1 minute
      taskId: timerState.value.taskId,
      taskTitle: timerState.value.taskTitle,
      taskCode: timerState.value.taskCode,
    }
    clear()
    return result
  }

  function clear() {
    if (!import.meta.client) return
    timerState.value = null
    localStorage.removeItem(TIMER_KEY)
    stopTicker()
  }

  onMounted(() => {
    loadFromStorage()
    if (timerState.value) startTicker()
  })

  onUnmounted(() => stopTicker())

  return { timerState, elapsedMinutes, elapsedDisplay, isRunning, start, stop, clear }
}
