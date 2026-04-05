<template>
  <Teleport to="body">
    <Transition name="slide-up">
      <div
        v-if="show"
        class="fixed inset-0 z-[90] flex items-end sm:items-center justify-center p-0 sm:p-4"
        role="dialog"
        aria-modal="true"
        @keydown.escape="$emit('close')"
      >
        <!-- Backdrop -->
        <div
          class="fixed inset-0 bg-black/75 backdrop-blur-sm"
          aria-hidden="true"
          @click="$emit('close')"
        />

        <!-- Panel -->
        <div
          class="relative w-full sm:max-w-2xl max-h-[92dvh] sm:max-h-[85vh] bg-gray-900 border border-gray-700 sm:rounded-2xl rounded-t-2xl shadow-2xl flex flex-col overflow-hidden"
          @click.stop
        >
          <!-- Header -->
          <div class="flex items-start justify-between gap-3 px-5 pt-5 pb-4 border-b border-gray-700 shrink-0">
            <div class="flex items-center gap-3 min-w-0">
              <!-- Avatar -->
              <div
                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold shrink-0"
                :style="{ background: avatarColor(userEmail) }"
              >
                {{ initial }}
              </div>
              <div class="min-w-0">
                <div class="text-white font-semibold text-base truncate">
                  {{ displayName || userEmail.split('@')[0] }}
                </div>
                <div class="text-sm text-gray-400">{{ formattedDate }}</div>
              </div>
            </div>
            <button
              type="button"
              @click="$emit('close')"
              class="shrink-0 p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-800 transition-colors"
              aria-label="Close"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Loading -->
          <div v-if="store.dayDetailLoading" class="flex-1 flex items-center justify-center py-16">
            <div class="flex flex-col items-center gap-3 text-gray-500">
              <svg class="animate-spin w-8 h-8" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
              </svg>
              <span class="text-sm">กำลังโหลด...</span>
            </div>
          </div>

          <!-- Error -->
          <div v-else-if="store.dayDetailError" class="flex-1 p-5 text-red-400 text-sm">
            {{ store.dayDetailError }}
          </div>

          <!-- Content -->
          <div v-else-if="detail" class="flex-1 overflow-y-auto p-5 space-y-5">

            <!-- Attendance Banner -->
            <div
              v-if="detail.attendance"
              class="rounded-xl px-4 py-3 flex flex-wrap items-center gap-3 border"
              :class="detail.attendance.is_late
                ? 'bg-rose-950/30 border-rose-700/40'
                : 'bg-teal-950/20 border-teal-700/30'"
            >
              <!-- Status badge -->
              <div
                class="flex items-center gap-1.5 text-sm font-semibold shrink-0"
                :class="detail.attendance.is_late ? 'text-rose-300' : 'text-teal-300'"
              >
                <span>{{ detail.attendance.is_late ? '🕐' : '✅' }}</span>
                <span>{{ detail.attendance.is_late ? 'เข้างานสาย' : 'เข้างานตรงเวลา' }}</span>
              </div>
              <!-- Times -->
              <div class="flex items-center gap-4 text-sm flex-wrap">
                <div v-if="detail.attendance.check_in_at" class="flex items-center gap-1.5 text-gray-300">
                  <span class="text-gray-500 text-xs">เข้า</span>
                  <span class="font-bold font-mono">{{ detail.attendance.check_in_at }}</span>
                </div>
                <div v-if="detail.attendance.check_out_at" class="flex items-center gap-1.5" :class="detail.attendance.early_checkout ? 'text-amber-300' : 'text-gray-300'">
                  <span class="text-gray-500 text-xs">ออก</span>
                  <span class="font-bold font-mono">{{ detail.attendance.check_out_at }}</span>
                  <span v-if="detail.attendance.early_checkout" class="text-[10px] px-1.5 py-0.5 rounded bg-amber-900/60 text-amber-300 font-semibold">🚪 กลับก่อนเวลา</span>
                </div>
                <div v-if="!detail.attendance.check_out_at" class="text-gray-600 text-xs italic">ยังไม่ได้ check-out</div>
              </div>
              <!-- Method badge -->
              <div v-if="detail.attendance.check_in_method" class="ml-auto shrink-0">
                <span class="text-[10px] px-2 py-1 rounded-full bg-gray-700/70 text-gray-400 font-medium uppercase">
                  {{ detail.attendance.check_in_method }}
                </span>
              </div>
            </div>
            <!-- No attendance record -->
            <div
              v-else
              class="rounded-xl px-4 py-3 bg-gray-800/40 border border-gray-700/30 flex items-center gap-2 text-gray-500 text-sm"
            >
              <span>📭</span>
              <span>ไม่มีข้อมูล check-in วันนี้</span>
            </div>

            <!-- Summary row -->
            <div class="grid gap-3" :class="(detail.deployed_requests?.length ?? 0) > 0 ? 'grid-cols-5' : 'grid-cols-4'">
              <div class="bg-gray-800 rounded-xl p-3 text-center border"
                :class="detail.has_daily_pulse ? 'border-violet-700/40' : 'border-red-700/40'">
                <div class="text-lg font-bold" :class="detail.has_daily_pulse ? 'text-violet-400' : 'text-red-400'">
                  {{ detail.has_daily_pulse ? '✓' : '✗' }}
                </div>
                <div class="text-[10px] text-gray-500 mt-0.5">Daily Pulse</div>
              </div>
              <div class="bg-gray-800 rounded-xl p-3 text-center border border-emerald-700/40">
                <div class="text-lg font-bold text-emerald-400">{{ detail.completed_tasks.length }}</div>
                <div class="text-[10px] text-gray-500 mt-0.5">Tasks ปิด</div>
              </div>
              <div v-if="(detail.deployed_requests?.length ?? 0) > 0" class="bg-gray-800 rounded-xl p-3 text-center border border-orange-700/40">
                <div class="text-lg font-bold text-orange-400">{{ detail.deployed_requests.length }}</div>
                <div class="text-[10px] text-gray-500 mt-0.5">🚀 Deploy</div>
              </div>
              <div class="bg-gray-800 rounded-xl p-3 text-center border"
                :class="detail.reworks.length > 0 ? 'border-red-700/40' : 'border-gray-700'">
                <div class="text-lg font-bold" :class="detail.reworks.length > 0 ? 'text-red-400' : 'text-gray-600'">
                  {{ detail.reworks.length }}
                </div>
                <div class="text-[10px] text-gray-500 mt-0.5">Rework</div>
              </div>
              <div class="bg-gray-800 rounded-xl p-3 text-center border border-blue-700/40">
                <div class="text-lg font-bold text-blue-400">{{ (detail.total_logged_minutes / 60).toFixed(1) }}h</div>
                <div class="text-[10px] text-gray-500 mt-0.5">Logged</div>
              </div>
            </div>

            <!-- Time Logs section -->
            <section>
              <h3 class="text-sm font-semibold text-gray-300 mb-2 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-blue-500 inline-block"></span>
                Time Logs
                <span class="text-gray-600 font-normal">({{ detail.time_logs.length }} รายการ)</span>
              </h3>
              <div v-if="detail.time_logs.length === 0" class="text-sm text-gray-600 italic px-2">
                ไม่มี log time วันนี้
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="log in detail.time_logs"
                  :key="log.task_id + log.minutes"
                  class="flex items-start gap-3 bg-gray-800/60 rounded-lg px-3 py-2.5 border border-gray-700/50 hover:border-blue-700/30 transition-colors"
                >
                  <div class="shrink-0 mt-0.5">
                    <div class="w-7 h-7 rounded-lg bg-blue-900/50 border border-blue-700/30 flex items-center justify-center text-blue-300 text-xs font-bold">
                      {{ log.hours.toFixed(1) }}h
                    </div>
                  </div>
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center gap-2 flex-wrap">
                      <span v-if="log.task_code" class="text-[10px] text-gray-500 font-mono">{{ log.task_code }}</span>
                      <span class="text-sm text-white font-medium truncate">{{ log.task_title }}</span>
                      <span
                        v-if="log.work_type"
                        class="text-[9px] px-1.5 py-0.5 rounded font-semibold uppercase tracking-wide shrink-0"
                        :class="workTypeBadge(log.work_type)"
                      >{{ log.work_type }}</span>
                      <span v-if="log.is_timer" class="text-[9px] px-1.5 py-0.5 rounded bg-violet-900/50 text-violet-300 font-medium shrink-0">⏱ Timer</span>
                    </div>
                    <div v-if="log.description" class="text-xs text-gray-400 mt-0.5 line-clamp-2">
                      {{ log.description }}
                    </div>
                  </div>
                  <div class="shrink-0 text-right">
                    <div class="text-sm font-bold text-blue-300">{{ log.minutes }}m</div>
                  </div>
                </div>
                <!-- Total bar + work_type breakdown -->
                <div class="bg-blue-950/30 rounded-lg border border-blue-800/30 overflow-hidden">
                  <div class="flex items-center justify-between px-3 py-2">
                    <span class="text-xs text-blue-400 font-medium">รวมทั้งหมด</span>
                    <span class="text-sm font-bold text-blue-300">
                      {{ detail.total_logged_minutes }}m = {{ (detail.total_logged_minutes / 60).toFixed(2) }} ชม.
                    </span>
                  </div>
                  <!-- Work type breakdown mini-bar -->
                  <div v-if="workTypeBreakdown.length > 0" class="px-3 pb-2.5 flex flex-wrap gap-1.5">
                    <div
                      v-for="wt in workTypeBreakdown"
                      :key="wt.type"
                      class="flex items-center gap-1 text-[10px]"
                    >
                      <span class="px-1.5 py-0.5 rounded font-semibold uppercase" :class="workTypeBadge(wt.type)">{{ wt.type }}</span>
                      <span class="text-gray-400">{{ wt.minutes }}m</span>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- Completed Tasks section -->
            <section>
              <h3 class="text-sm font-semibold text-gray-300 mb-2 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-emerald-500 inline-block"></span>
                Tasks ที่ปิดวันนี้
                <span class="text-gray-600 font-normal">({{ detail.completed_tasks.length }} งาน)</span>
              </h3>
              <div v-if="detail.completed_tasks.length === 0" class="text-sm text-gray-600 italic px-2">
                ยังไม่มี task ที่ปิดวันนี้
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="task in detail.completed_tasks"
                  :key="task.task_id"
                  class="flex items-center gap-3 bg-gray-800/60 rounded-lg px-3 py-2.5 border border-gray-700/50 hover:border-emerald-700/30 transition-colors"
                >
                  <div class="w-6 h-6 rounded-full bg-emerald-900/60 border border-emerald-700/40 flex items-center justify-center shrink-0">
                    <svg class="w-3.5 h-3.5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                    </svg>
                  </div>
                  <div class="min-w-0 flex-1">
                    <div class="flex items-center gap-2 flex-wrap">
                      <span v-if="task.task_code" class="text-[10px] text-gray-500 font-mono">{{ task.task_code }}</span>
                      <span class="text-sm text-white truncate">{{ task.task_title }}</span>
                    </div>
                  </div>
                  <div class="flex items-center gap-2 shrink-0">
                    <span
                      class="text-[10px] px-1.5 py-0.5 rounded font-medium"
                      :class="taskTypeBadge(task.task_type)"
                    >{{ task.task_type }}</span>
                    <span v-if="task.story_points > 0" class="text-xs text-gray-400">{{ task.story_points }} SP</span>
                  </div>
                </div>
              </div>
            </section>

            <!-- Deployed Requests section (Chief Engineer) -->
            <section v-if="(detail.deployed_requests?.length ?? 0) > 0">
              <h3 class="text-sm font-semibold text-gray-300 mb-2 flex items-center gap-2">
                <span class="text-orange-400">🚀</span>
                Deployment ที่ Deploy แล้ว
                <span class="text-gray-600 font-normal">({{ detail.deployed_requests.length }} รายการ)</span>
              </h3>
              <div class="space-y-2">
                <div
                  v-for="dr in detail.deployed_requests"
                  :key="dr.id"
                  class="flex items-center gap-3 bg-orange-950/20 rounded-lg px-3 py-2.5 border border-orange-800/30"
                >
                  <div class="w-6 h-6 rounded-full bg-orange-900/60 border border-orange-700/40 flex items-center justify-center shrink-0 text-sm">
                    🚀
                  </div>
                  <div class="min-w-0 flex-1">
                    <span class="text-sm text-white font-medium truncate block">{{ dr.title }}</span>
                    <div class="flex items-center gap-2 mt-0.5">
                      <code class="text-[10px] text-cyan-400 font-mono">⎇ {{ dr.branch }}</code>
                      <span class="text-[10px] px-1.5 py-0.5 rounded font-semibold uppercase tracking-wide"
                        :class="{
                          'bg-orange-500/15 text-orange-300': dr.environment === 'PRODUCTION',
                          'bg-amber-500/15 text-amber-300': dr.environment === 'PRE-PROD',
                          'bg-violet-500/15 text-violet-300': dr.environment === 'STAGING',
                        }">{{ dr.environment }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- Reworks section -->
            <section>
              <h3 class="text-sm font-semibold text-gray-300 mb-2 flex items-center gap-2">
                <span class="w-2 h-2 rounded-full bg-red-500 inline-block"></span>
                Rework / ถูก Reject
                <span class="text-gray-600 font-normal">({{ detail.reworks.length }} ครั้ง)</span>
              </h3>
              <div v-if="detail.reworks.length === 0" class="text-sm text-gray-600 italic px-2">
                ไม่มี rework วันนี้ 🎉
              </div>
              <div v-else class="space-y-2">
                <div
                  v-for="(rw, idx) in detail.reworks"
                  :key="rw.task_id + idx"
                  class="bg-red-950/20 rounded-lg px-3 py-2.5 border border-red-800/30"
                >
                  <div class="flex items-center gap-2 flex-wrap mb-1">
                    <span class="text-red-400 text-xs font-medium">✗ REJECTED</span>
                    <span v-if="rw.task_code" class="text-[10px] text-gray-500 font-mono">{{ rw.task_code }}</span>
                    <span class="text-sm text-white truncate font-medium">{{ rw.task_title }}</span>
                  </div>
                  <div class="text-xs text-red-300/70 leading-relaxed">
                    {{ cleanComment(rw.rejected_comment) }}
                  </div>
                </div>
              </div>
            </section>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { usePerformanceStore } from '~/core/modules/performance/performance-store'

const props = defineProps<{
  show: boolean
  userId: number
  userEmail: string
  displayName?: string
  date: string
}>()

defineEmits<{ close: [] }>()

const store = usePerformanceStore()
const detail = computed(() => store.dayDetail)

watch(
  () => [props.show, props.userId, props.date] as const,
  ([show]) => {
    if (show && props.userId && props.date) {
      store.fetchDayDetail(props.userId, props.date)
    }
  },
  { immediate: true }
)

const formattedDate = computed(() => {
  if (!props.date) return ''
  const d = new Date(props.date + 'T00:00:00')
  return d.toLocaleDateString('th-TH', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
})

const initial = computed(() => {
  const name = props.displayName || props.userEmail
  return name.charAt(0).toUpperCase()
})

const avatarColors = [
  'linear-gradient(135deg,#6366f1,#8b5cf6)',
  'linear-gradient(135deg,#ec4899,#f43f5e)',
  'linear-gradient(135deg,#14b8a6,#06b6d4)',
  'linear-gradient(135deg,#f59e0b,#f97316)',
  'linear-gradient(135deg,#10b981,#059669)',
  'linear-gradient(135deg,#3b82f6,#6366f1)',
]
function avatarColor(email: string): string {
  let hash = 0
  for (const c of email) hash = (hash * 31 + c.charCodeAt(0)) & 0xffffffff
  return avatarColors[Math.abs(hash) % avatarColors.length]
}

function taskTypeBadge(type: string): string {
  const map: Record<string, string> = {
    FEATURE: 'bg-purple-900/60 text-purple-300',
    BUG: 'bg-red-900/60 text-red-300',
    TASK: 'bg-gray-700 text-gray-400',
  }
  return map[type] || 'bg-gray-700 text-gray-400'
}

function cleanComment(comment: string): string {
  return comment.replace(/^\[REJECTED\]\s*/i, '').trim() || '(ไม่มีคำอธิบาย)'
}

const workTypeBadge = (type: string): string => {
  const map: Record<string, string> = {
    DEV: 'bg-blue-900/60 text-blue-300',
    REVIEW: 'bg-cyan-900/60 text-cyan-300',
    TESTING: 'bg-green-900/60 text-green-300',
    MEETING: 'bg-orange-900/60 text-orange-300',
    RESEARCH: 'bg-purple-900/60 text-purple-300',
    OTHER: 'bg-gray-700/80 text-gray-400',
  }
  return map[type?.toUpperCase()] || 'bg-gray-700/80 text-gray-400'
}

const workTypeBreakdown = computed(() => {
  if (!detail.value?.time_logs?.length) return []
  const map: Record<string, number> = {}
  for (const log of detail.value.time_logs) {
    const wt = log.work_type || 'DEV'
    map[wt] = (map[wt] || 0) + log.minutes
  }
  return Object.entries(map)
    .map(([type, minutes]) => ({ type, minutes }))
    .sort((a, b) => b.minutes - a.minutes)
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: opacity 0.2s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
}
.slide-up-enter-active > div:last-child,
.slide-up-leave-active > div:last-child {
  transition: transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.slide-up-enter-from > div:last-child {
  transform: translateY(40px) scale(0.97);
}
.slide-up-leave-to > div:last-child {
  transform: translateY(20px) scale(0.97);
}
</style>
