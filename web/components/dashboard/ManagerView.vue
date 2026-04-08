<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-blue-500/15 border border-blue-500/30 flex items-center justify-center">
            <svg class="w-4 h-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold text-white">Manager Operations Center</h1>
            <p class="text-xs text-gray-500">People operations, project reliability, and leave approvals</p>
          </div>
        </div>
        <button
          @click="refresh"
          :disabled="isLoading"
          class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-700 bg-gray-800/60 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50"
        >
          <svg class="h-3.5 w-3.5" :class="isLoading ? 'animate-spin' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-blue-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading manager dashboard…</p>
    </div>

    <div v-else-if="error" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Failed to load manager data</p>
          <p class="text-xs text-red-300 mt-0.5">{{ error }}</p>
        </div>
      </div>
    </div>

    <main v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">
      <CeoUATApprovalQueue />

      <section>
        <h2 class="section-label">Operations snapshot</h2>
        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-5">
            <p class="card-title">Employees Today</p>
            <p class="card-value">{{ peopleKpis.total }}</p>
            <p class="card-sub">present {{ peopleKpis.present }} · late {{ peopleKpis.late }}</p>
          </div>
          <div class="rounded-2xl border border-red-500/30 bg-red-900/15 p-5">
            <p class="card-title">Absent / Leave</p>
            <p class="card-value text-red-400">{{ peopleKpis.absentOrOnLeave }}</p>
            <p class="card-sub">today</p>
          </div>
          <div class="rounded-2xl border border-amber-500/30 bg-amber-900/15 p-5">
            <p class="card-title">Projects At Risk</p>
            <p class="card-value text-amber-400">{{ projectKpis.atRisk }}</p>
            <p class="card-sub">overdue-heavy projects</p>
          </div>
          <div class="rounded-2xl border border-cyan-500/30 bg-cyan-900/15 p-5">
            <p class="card-title">Leave Approvals</p>
            <p class="card-value text-cyan-300">{{ pendingLeaves.length }}</p>
            <p class="card-sub">waiting your decision</p>
          </div>
        </div>
      </section>

      <PmPerformanceSection :projects="projects" audience="manager" />

      <section class="grid gap-6 lg:grid-cols-3">
        <div class="lg:col-span-2 rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-semibold text-white">Leave Approval Queue</h3>
            <span class="text-xs text-gray-500">manager-owned workflow</span>
          </div>

          <div v-if="pendingLeaves.length === 0" class="text-sm text-gray-500 italic py-6 text-center">
            No pending leave approvals.
          </div>

          <div v-else class="space-y-3">
            <article
              v-for="leave in pendingLeaves"
              :key="leave.id"
              class="rounded-xl border border-gray-700/70 bg-gray-900/40 p-4"
            >
              <div class="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                <div>
                  <p class="text-sm font-semibold text-white">{{ leave.user_display_name || leave.user_email || ('User #' + leave.user_id) }}</p>
                  <p class="text-xs text-gray-400 mt-0.5">{{ leave.leave_type }} · {{ leave.days_requested }} day(s)</p>
                  <p class="text-xs text-gray-500 mt-1">{{ formatDate(leave.start_date) }} → {{ formatDate(leave.end_date) }}</p>
                  <p class="text-xs text-gray-400 mt-2">{{ leave.reason }}</p>
                </div>
                <div class="flex flex-col gap-2 sm:items-end min-w-[220px]">
                  <textarea
                    v-model="reviewCommentById[leave.id]"
                    rows="2"
                    class="w-full rounded-lg border border-gray-700 bg-gray-800/80 px-3 py-2 text-xs text-white focus:border-blue-500 focus:outline-none"
                    placeholder="Comment (optional for approve, recommended for reject)"
                  />
                  <div class="flex items-center gap-2 w-full justify-end">
                    <button
                      class="rounded-lg border border-red-500/40 bg-red-900/20 px-3 py-1.5 text-xs font-semibold text-red-300 hover:bg-red-900/40 disabled:opacity-50"
                      :disabled="reviewingId === leave.id"
                      @click="reviewLeave(leave.id, 'REJECTED')"
                    >
                      Reject
                    </button>
                    <button
                      class="rounded-lg border border-emerald-500/40 bg-emerald-900/20 px-3 py-1.5 text-xs font-semibold text-emerald-300 hover:bg-emerald-900/40 disabled:opacity-50"
                      :disabled="reviewingId === leave.id"
                      @click="reviewLeave(leave.id, 'APPROVED')"
                    >
                      Approve
                    </button>
                  </div>
                </div>
              </div>
            </article>
          </div>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h3 class="text-sm font-semibold text-white mb-4">People Operations Today</h3>
          <div class="space-y-3 text-xs">
            <div class="flex items-center justify-between">
              <span class="text-gray-400">Present</span>
              <span class="font-bold text-emerald-400">{{ peopleKpis.present }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-400">Late</span>
              <span class="font-bold text-amber-400">{{ peopleKpis.late }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-400">Early checkout</span>
              <span class="font-bold text-orange-300">{{ peopleKpis.earlyCheckout }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-gray-400">No check-in</span>
              <span class="font-bold text-red-400">{{ peopleKpis.noCheckIn }}</span>
            </div>
            <div class="pt-3 border-t border-gray-700">
              <p class="text-gray-500">You can manage attendance config and records from <NuxtLink to="/admin/attendance-config" class="text-blue-400 hover:text-blue-300">Attendance Admin</NuxtLink>.</p>
            </div>
          </div>
        </div>
      </section>

      <section class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-white">Project Command Center</h3>
          <NuxtLink to="/projects" class="text-xs text-gray-400 hover:text-white">View all projects</NuxtLink>
        </div>

        <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <NuxtLink
            v-for="project in riskyProjects"
            :key="project.id"
            :to="`/projects/${project.id}`"
            class="block rounded-xl border border-amber-500/25 bg-amber-900/10 p-4 transition hover:border-amber-400/50 hover:bg-amber-900/20 hover:shadow-lg hover:shadow-amber-900/20 focus:outline-none focus:ring-2 focus:ring-amber-400/60"
          >
            <div class="flex items-center justify-between mb-2 gap-2">
              <p class="text-sm font-semibold text-white truncate">{{ project.name }}</p>
              <span class="text-[10px] uppercase px-2 py-0.5 rounded-full border border-amber-500/40 text-amber-300">At Risk</span>
            </div>
            <p class="text-xs text-gray-400">Code: {{ project.code || 'N/A' }}</p>
            <div class="mt-3 space-y-1.5 text-xs">
              <div class="flex justify-between"><span class="text-gray-500">Overdue tasks</span><span class="font-semibold text-red-400">{{ project.task_overdue || 0 }}</span></div>
              <div class="flex justify-between"><span class="text-gray-500">Completed</span><span class="font-semibold text-emerald-400">{{ project.task_completed || 0 }}/{{ project.task_total || 0 }}</span></div>
              <div class="flex justify-between"><span class="text-gray-500">Status</span><span class="font-semibold text-gray-300">{{ project.status }}</span></div>
            </div>
          </NuxtLink>

          <div v-if="riskyProjects.length === 0" class="rounded-xl border border-emerald-500/30 bg-emerald-900/10 p-4 sm:col-span-2 lg:col-span-3">
            <p class="text-sm font-semibold text-emerald-300">All active projects are currently stable.</p>
            <p class="text-xs text-emerald-400/80 mt-1">No overdue-heavy project found in current scope.</p>
          </div>
        </div>
      </section>

      <section class="grid gap-6 lg:grid-cols-2">
        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h3 class="text-sm font-semibold text-white mb-4">Leave Trend by Team / Month</h3>
          <div class="space-y-2">
            <div v-for="m in monthlyTrend" :key="m.month" class="rounded-lg border border-gray-700/70 bg-gray-900/30 p-3 text-xs">
              <div class="flex items-center justify-between">
                <span class="text-gray-300 font-semibold">{{ m.month }}</span>
                <span class="text-gray-400">Days {{ m.totalDays }}</span>
              </div>
              <div class="mt-1 flex gap-4 text-[11px]">
                <span class="text-cyan-300">Requested {{ m.requested }}</span>
                <span class="text-emerald-300">Approved {{ m.approved }}</span>
                <span class="text-red-300">Rejected {{ m.rejected }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h3 class="text-sm font-semibold text-white mb-4">Leave Policy & Holiday Calendar</h3>
          <div class="space-y-3 text-xs">
            <div v-for="p in leavePolicies" :key="p.id" class="rounded-lg border border-gray-700/70 bg-gray-900/30 p-3">
              <div class="flex justify-between"><span class="text-gray-300">{{ p.leave_type }}</span><span class="text-emerald-300">{{ p.annual_quota_days }} days</span></div>
              <p class="text-gray-500 mt-1">Carry forward max: {{ p.max_carry_forward_days }} · {{ p.is_active ? 'Active' : 'Inactive' }}</p>
            </div>
            <div class="pt-2 border-t border-gray-700">
              <p class="text-gray-400 mb-1">Upcoming holidays</p>
              <p v-for="h in holidays.slice(0, 5)" :key="h.id" class="text-gray-500">{{ formatDate(h.date) }} · {{ h.name }}</p>
            </div>
          </div>
        </div>
      </section>

      <section v-if="selectedAuditLeaveId" class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-semibold text-white">Audit Log · Leave #{{ selectedAuditLeaveId }}</h3>
          <button class="text-xs text-gray-400 hover:text-white" @click="selectedAuditLeaveId = null">Close</button>
        </div>
        <div class="space-y-2 text-xs">
          <div v-for="log in selectedLeaveAudit" :key="log.id" class="rounded-lg border border-gray-700/70 bg-gray-900/30 p-3">
            <div class="flex items-center justify-between">
              <span class="text-gray-300 font-semibold">{{ log.action }}</span>
              <span class="text-gray-500">{{ formatDate(log.created_at) }}</span>
            </div>
            <p class="text-gray-500 mt-1">{{ log.old_status || '-' }} → {{ log.new_status || '-' }} · {{ log.actor_name || log.actor_email || 'System' }}</p>
            <p v-if="log.comment" class="text-gray-400 mt-1">{{ log.comment }}</p>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type LeaveRequest, type LeaveTrendPoint, type LeavePolicy, type HolidayCalendar, type LeaveAuditLog } from '~/core/modules/attendance/infrastructure/attendance-api'
import { useProjectsApi, type Project } from '~/core/modules/projects/infrastructure/projects-api'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import PmPerformanceSection from '~/components/dashboard/PmPerformanceSection.vue'
import CeoUATApprovalQueue from '~/components/dashboard/CeoUATApprovalQueue.vue'

const attendanceApi = useAttendanceApi()
const projectsApi = useProjectsApi()
const performanceStore = usePerformanceStore()

const isLoading = ref(true)
const error = ref('')
const reviewingId = ref<number | null>(null)
const attendanceDate = ref(new Date().toISOString().slice(0, 10))

const attendanceRecords = ref<any[]>([])
const projects = ref<Project[]>([])
const pendingLeaves = ref<LeaveRequest[]>([])
const reviewCommentById = ref<Record<number, string>>({})
const leaveTrend = ref<LeaveTrendPoint[]>([])
const leavePolicies = ref<LeavePolicy[]>([])
const holidays = ref<HolidayCalendar[]>([])
const selectedAuditLeaveId = ref<number | null>(null)
const selectedLeaveAudit = ref<LeaveAuditLog[]>([])

const peopleKpis = computed(() => {
  const records = attendanceRecords.value || []
  const present = records.filter((r: any) => r.status === 'present').length
  const late = records.filter((r: any) => r.is_late || r.status === 'late').length
  const earlyCheckout = records.filter((r: any) => r.early_checkout).length
  const noCheckIn = records.filter((r: any) => !r.check_in_at).length
  return {
    total: records.length,
    present,
    late,
    earlyCheckout,
    noCheckIn,
    absentOrOnLeave: noCheckIn,
  }
})

const riskyProjects = computed(() => {
  return (projects.value || [])
    .filter(p => (p.task_overdue || 0) > 0 || p.status === 'ON_HOLD')
    .sort((a, b) => (b.task_overdue || 0) - (a.task_overdue || 0))
    .slice(0, 9)
})

const projectKpis = computed(() => ({
  total: projects.value.length,
  atRisk: riskyProjects.value.length,
}))

function formatDate(input: string) {
  return new Date(input).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
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
})

async function reviewLeave(id: number, status: 'APPROVED' | 'REJECTED') {
  reviewingId.value = id
  error.value = ''
  try {
    await attendanceApi.reviewLeaveRequest(id, {
      status,
      comment: reviewCommentById.value[id] || '',
    })
    pendingLeaves.value = pendingLeaves.value.filter(x => x.id !== id)
    await openAudit(id)
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Failed to review leave request'
  } finally {
    reviewingId.value = null
  }
}

async function openAudit(leaveId: number) {
  selectedAuditLeaveId.value = leaveId
  try {
    const res = await attendanceApi.getLeaveAudit(leaveId)
    selectedLeaveAudit.value = res.items || []
  } catch {
    selectedLeaveAudit.value = []
  }
}

async function refresh() {
  isLoading.value = true
  error.value = ''
  try {
    const [recordsRes, projectsRes, leavesRes, trendRes, policiesRes, holidayRes] = await Promise.all([
      attendanceApi.adminRecords(attendanceDate.value),
      projectsApi.getProjects(),
      attendanceApi.getPendingLeaveRequests(),
      attendanceApi.getLeaveTrend(),
      attendanceApi.getLeavePolicies(),
      attendanceApi.getHolidays(),
    ])
    attendanceRecords.value = recordsRes.records || []
    projects.value = projectsRes || []
    pendingLeaves.value = leavesRes.items || []
    leaveTrend.value = trendRes.items || []
    leavePolicies.value = policiesRes.items || []
    holidays.value = holidayRes.items || []
  } catch (e: any) {
    error.value = e?.data?.error || e?.message || 'Failed to load manager dashboard'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  performanceStore.fetchAll('CEO')
  refresh()
})
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
.card-title {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5;
}
.card-value {
  @apply text-2xl font-black text-white tabular-nums;
}
.card-sub {
  @apply text-xs text-gray-500 mt-1;
}
</style>
