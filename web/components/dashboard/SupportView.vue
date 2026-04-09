<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <div class="sticky top-0 z-20 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto flex max-w-screen-xl items-center justify-between px-6 py-4">
        <div class="flex items-center gap-3">
          <div class="flex h-9 w-9 items-center justify-center rounded-lg border border-cyan-500/30 bg-cyan-500/10">
            <svg class="h-5 w-5 text-cyan-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6M7 4h10a2 2 0 012 2v12a2 2 0 01-2 2H7a2 2 0 01-2-2V6a2 2 0 012-2z" />
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold">Support Operations Dashboard</h1>
            <p class="text-xs text-gray-500">ธุรการบุคคล: ขาด · ลา · มาสาย · วินัยการทำงาน</p>
          </div>
        </div>

        <button
          @click="refresh"
          :disabled="isLoading"
          class="inline-flex items-center gap-2 rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2 text-xs font-medium text-gray-300 transition-colors hover:border-gray-600 hover:bg-gray-700 hover:text-gray-900 dark:text-white disabled:opacity-50"
        >
          <svg class="h-3.5 w-3.5" :class="isLoading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
      <svg class="mb-3 h-8 w-8 animate-spin text-cyan-300" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading support operations...</p>
    </div>

    <div v-else-if="error" class="mx-auto max-w-screen-xl px-6 py-8">
      <div class="rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <p class="text-sm font-semibold">โหลดข้อมูลไม่สำเร็จ</p>
        <p class="mt-1 text-xs text-red-300">{{ error }}</p>
      </div>
    </div>

    <main v-else class="mx-auto max-w-screen-xl space-y-8 px-6 py-8">
      <section>
        <h2 class="section-label">Today workforce pulse</h2>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-5">
            <p class="card-title">Employees tracked</p>
            <p class="card-value">{{ attendanceKpis.total }}</p>
            <p class="card-sub">attendance record วันนี้</p>
          </div>
          <div class="rounded-2xl border border-rose-500/30 bg-rose-900/15 p-5">
            <p class="card-title">Late arrivals</p>
            <p class="card-value text-rose-300">{{ attendanceKpis.late }}</p>
            <p class="card-sub">มาสายวันนี้</p>
          </div>
          <div class="rounded-2xl border border-amber-500/30 bg-amber-900/15 p-5">
            <p class="card-title">Early checkout</p>
            <p class="card-value text-amber-300">{{ attendanceKpis.earlyCheckout }}</p>
            <p class="card-sub">กลับก่อนเวลา</p>
          </div>
          <div class="rounded-2xl border border-violet-500/30 bg-violet-900/15 p-5">
            <p class="card-title">Daily Standup missing</p>
            <p class="card-value text-violet-300">{{ noPulseToday }}</p>
            <p class="card-sub">คนที่ไม่ส่ง Daily Standup วันนี้</p>
          </div>
        </div>
      </section>

      <section class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-white">Attendance exceptions (today)</h3>
            <NuxtLink to="/attendance" class="text-xs text-cyan-300 hover:text-cyan-200">เปิดหน้า Attendance</NuxtLink>
          </div>

          <div v-if="exceptionRows.length === 0" class="rounded-xl border border-emerald-500/30 bg-emerald-900/10 p-4 text-sm text-emerald-300">
            ไม่พบรายการผิดปกติในวันนี้
          </div>

          <div v-else class="space-y-2">
            <div
              v-for="row in exceptionRows"
              :key="row.id"
              class="flex items-center justify-between rounded-lg border border-gray-700/70 bg-gray-900/35 px-3 py-2"
            >
              <div class="min-w-0">
                <p class="truncate text-sm font-medium text-white">{{ row.name }}</p>
                <p class="truncate text-xs text-gray-500">{{ row.email }}</p>
              </div>
              <div class="flex items-center gap-2 text-[11px]">
                <span v-if="row.isLate" class="rounded border border-rose-500/40 bg-rose-900/20 px-2 py-0.5 text-rose-300">มาสาย</span>
                <span v-if="row.earlyCheckout" class="rounded border border-amber-500/40 bg-amber-900/20 px-2 py-0.5 text-amber-300">กลับก่อน</span>
                <span v-if="!row.checkInAt" class="rounded border border-red-500/40 bg-red-900/20 px-2 py-0.5 text-red-300">ยังไม่ check-in</span>
              </div>
            </div>
          </div>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h3 class="mb-4 text-sm font-semibold text-white">Leave queue</h3>
          <div v-if="pendingLeaves.length === 0" class="text-xs text-gray-500">ไม่มีรายการลาแบบ pending</div>
          <div v-else class="space-y-2">
            <div
              v-for="leave in pendingLeaves.slice(0, 7)"
              :key="leave.id"
              class="rounded-lg border border-gray-700/70 bg-gray-900/35 p-3"
            >
              <p class="truncate text-xs font-semibold text-white">{{ leave.user_display_name || leave.user_email || ('User #' + leave.user_id) }}</p>
              <p class="mt-0.5 text-[11px] text-gray-400">{{ leave.leave_type }} · {{ leave.days_requested }} วัน</p>
              <p class="mt-1 text-[11px] text-gray-500">{{ formatDate(leave.start_date) }} → {{ formatDate(leave.end_date) }}</p>
            </div>
          </div>
          <div class="mt-3 border-t border-gray-700 pt-3">
            <NuxtLink to="/admin/leave" class="text-xs text-cyan-300 hover:text-cyan-200">ดูหน้าจัดการลาเต็มรูปแบบ</NuxtLink>
          </div>
        </div>
      </section>

      <section class="grid gap-6 lg:grid-cols-2">
        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <div class="mb-4 flex items-center justify-between">
            <h3 class="text-sm font-semibold text-white">Discipline leaderboard (7 days)</h3>
            <NuxtLink to="/discipline" class="text-xs text-orange-300 hover:text-orange-200">เปิดหน้า Discipline</NuxtLink>
          </div>
          <div class="space-y-2">
            <div
              v-for="u in disciplineTop"
              :key="u.user_id"
              class="flex items-center justify-between rounded-lg border border-gray-700/70 bg-gray-900/35 px-3 py-2 text-xs"
            >
              <div class="min-w-0">
                <p class="truncate font-semibold text-white">{{ u.user_display_name || u.user_email }}</p>
                <p class="truncate text-gray-500">{{ u.user_email }}</p>
              </div>
              <div class="text-right">
                <p class="font-bold" :class="scoreColor(disciplineScore(u))">{{ disciplineScore(u) }}%</p>
                <p class="text-[10px] text-gray-500">No Pulse {{ u.missed_pulse_count }} · Late {{ u.total_late_days ?? 0 }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h3 class="mb-4 text-sm font-semibold text-white">Leave trend overview</h3>
          <div class="space-y-2">
            <div v-for="m in monthlyTrend" :key="m.month" class="rounded-lg border border-gray-700/70 bg-gray-900/30 p-3 text-xs">
              <div class="flex items-center justify-between">
                <span class="font-semibold text-gray-300">{{ m.month }}</span>
                <span class="text-gray-500">Days {{ m.totalDays }}</span>
              </div>
              <div class="mt-1 flex gap-4 text-[11px]">
                <span class="text-cyan-300">Req {{ m.requested }}</span>
                <span class="text-emerald-300">Apv {{ m.approved }}</span>
                <span class="text-red-300">Rej {{ m.rejected }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type AttendanceRecord, type LeaveRequest, type LeaveTrendPoint } from '~/core/modules/attendance/infrastructure/attendance-api'
import { usePerformanceApi, type DisciplineUser } from '~/core/modules/performance/performance-api'

const attendanceApi = useAttendanceApi()
const performanceApi = usePerformanceApi()

const isLoading = ref(true)
const error = ref('')

const attendanceRecords = ref<AttendanceRecord[]>([])
const pendingLeaves = ref<LeaveRequest[]>([])
const leaveTrend = ref<LeaveTrendPoint[]>([])
const disciplineUsers = ref<DisciplineUser[]>([])

const today = computed(() => new Date().toISOString().slice(0, 10))

const attendanceKpis = computed(() => {
  const records = attendanceRecords.value
  return {
    total: records.length,
    late: records.filter(r => r.is_late || r.status === 'late').length,
    earlyCheckout: records.filter(r => r.early_checkout).length,
  }
})

const exceptionRows = computed(() => {
  return attendanceRecords.value
    .filter(r => r.is_late || r.early_checkout || !r.check_in_at)
    .map(r => ({
      id: r.id,
      email: r.user_email || `user-${r.user_id}`,
      name: r.user_display_name || r.user_email || `User #${r.user_id}`,
      isLate: Boolean(r.is_late || r.status === 'late'),
      earlyCheckout: Boolean(r.early_checkout),
      checkInAt: r.check_in_at,
    }))
    .sort((a, b) => Number(b.isLate) - Number(a.isLate))
})

function disciplineScore(u: DisciplineUser): number {
  const days = Math.max(1, u.days?.length || 0)
  const pulsePct = Math.max(0, (days - u.missed_pulse_count) / days) * 40
  const logDays = (u.days || []).filter(d => d.logged_minutes > 0).length
  const logPct = (logDays / days) * 40
  const totalSubs = u.total_tasks_closed + u.total_reworks
  const reworkPct = totalSubs > 0 ? (1 - u.total_reworks / totalSubs) * 20 : 20
  return Math.round(pulsePct + logPct + reworkPct)
}

const disciplineTop = computed(() => {
  return [...disciplineUsers.value]
    .sort((a, b) => disciplineScore(b) - disciplineScore(a))
    .slice(0, 8)
})

const noPulseToday = computed(() => {
  return disciplineUsers.value.reduce((sum, user) => {
    const todayDay = user.days?.find(d => d.date === today.value)
    return sum + (todayDay && !todayDay.has_daily_pulse ? 1 : 0)
  }, 0)
})

function scoreColor(score: number): string {
  if (score >= 80) return 'text-emerald-400'
  if (score >= 50) return 'text-yellow-400'
  return 'text-red-400'
}

const monthlyTrend = computed(() => {
  const map = new Map<string, { requested: number; approved: number; rejected: number; totalDays: number }>()
  for (const it of leaveTrend.value) {
    const prev = map.get(it.month) || { requested: 0, approved: 0, rejected: 0, totalDays: 0 }
    prev.requested += it.requested
    prev.approved += it.approved
    prev.rejected += it.rejected
    prev.totalDays += it.total_days
    map.set(it.month, prev)
  }
  return Array.from(map.entries())
    .map(([month, v]) => ({ month, ...v }))
    .sort((a, b) => a.month.localeCompare(b.month))
    .slice(-6)
})

function formatDate(input: string): string {
  return new Date(input).toLocaleDateString('th-TH', { day: '2-digit', month: 'short', year: 'numeric' })
}

async function refresh() {
  isLoading.value = true
  error.value = ''
  try {
    const from = new Date()
    from.setDate(from.getDate() - 6)
    const fromStr = from.toISOString().slice(0, 10)

    const [recordsRes, leavesRes, trendRes, disciplineRes] = await Promise.all([
      attendanceApi.adminRecords(today.value),
      attendanceApi.getPendingLeaveRequests(),
      attendanceApi.getLeaveTrend(),
      performanceApi.getDiscipline(fromStr, today.value),
    ])

    attendanceRecords.value = recordsRes.records || []
    pendingLeaves.value = leavesRes.items || []
    leaveTrend.value = trendRes.items || []
    disciplineUsers.value = disciplineRes.users || []
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Failed to load support dashboard'
  } finally {
    isLoading.value = false
  }
}

onMounted(refresh)
</script>

<style scoped>
.section-label {
  @apply mb-4 text-xs font-semibold uppercase tracking-widest text-gray-500;
}
.card-title {
  @apply mb-1.5 text-xs font-semibold uppercase tracking-widest text-gray-400;
}
.card-value {
  @apply tabular-nums text-2xl font-black text-white;
}
.card-sub {
  @apply mt-1 text-xs text-gray-500;
}
</style>
