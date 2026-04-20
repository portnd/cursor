<template>
  <div class="min-h-full px-4 py-5 text-slate-900 dark:text-white md:px-8 md:py-8">
    <section class="mb-6 rounded-2xl border border-slate-200/80 bg-white/90 p-5 shadow-sm dark:border-white/10 dark:bg-slate-900/70">
      <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
        <div>
          <p class="text-xs uppercase tracking-[0.18em] text-violet-500 dark:text-violet-300">Work Load Intelligence</p>
          <h1 class="mt-1 text-2xl font-bold md:text-3xl">Team Work Load Dashboard</h1>
          <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">
            สร้างเพื่อช่วยตัดสินใจเร็วขึ้น: ใคร overload, ใครยังมี capacity, และควรโยกงานอย่างไรให้ยุติธรรม
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-2">
          <select v-model="selectedRange" class="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-white/15 dark:bg-slate-800">
            <option value="week">สัปดาห์นี้</option>
            <option value="month">เดือนนี้</option>
            <option value="quarter">ไตรมาสนี้</option>
          </select>

          <select v-model="selectedRole" class="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-white/15 dark:bg-slate-800">
            <option value="all">ทุกบทบาท</option>
            <option v-for="role in roleOptions" :key="role" :value="role">{{ role }}</option>
          </select>

          <select v-model="selectedProjectId" class="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm dark:border-white/15 dark:bg-slate-800">
            <option value="all">ทุกโปรเจกต์</option>
            <option v-for="project in projects" :key="project.id" :value="project.id">{{ project.name }}</option>
          </select>

          <button
            class="rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 dark:border-white/20 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700"
            :disabled="loading"
            @click="loadWorkload"
          >
            {{ loading ? 'กำลังโหลด...' : 'Refresh' }}
          </button>

          <button
            class="rounded-lg bg-gradient-to-r from-rose-600 to-orange-600 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-orange-500/30 transition hover:brightness-110 disabled:opacity-60"
            :disabled="recommendations.length === 0"
            @click="openRecommendations = !openRecommendations"
          >
            Rebalance Workload ({{ recommendations.length }})
          </button>
        </div>
      </div>

      <div class="mt-5 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        <article v-for="kpi in kpis" :key="kpi.label" class="rounded-xl border border-slate-200/80 bg-slate-50/80 p-4 dark:border-white/10 dark:bg-slate-800/50">
          <p class="text-xs uppercase tracking-wide text-slate-500 dark:text-slate-400">{{ kpi.label }}</p>
          <p class="mt-2 text-2xl font-bold">{{ kpi.value }}</p>
          <p class="text-xs" :class="kpi.tone">{{ kpi.meta }}</p>
        </article>
      </div>

      <p v-if="error" class="mt-4 rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-sm text-rose-700 dark:border-rose-500/30 dark:bg-rose-500/10 dark:text-rose-300">
        {{ error }}
      </p>

      <div
        v-if="unassignedTasksCount > 0"
        class="mt-4 rounded-xl border border-amber-300/80 bg-amber-50 px-4 py-3 text-amber-800 dark:border-amber-500/40 dark:bg-amber-500/10 dark:text-amber-200"
      >
        <p class="text-sm font-semibold">มี {{ unassignedTasksCount }} งานยังไม่ assign</p>
        <p class="text-xs opacity-90">งานที่ยังไม่ assign จะไม่ถูกรวมในโหลดของพนักงาน เพื่อไม่ให้ข้อมูลบิดเบือน</p>
        <div class="mt-3 flex flex-wrap gap-2">
          <button
            class="rounded-lg border border-amber-400/70 px-3 py-1.5 text-xs font-semibold transition hover:bg-amber-100 dark:hover:bg-amber-500/20"
            @click="showUnassignedOnly = !showUnassignedOnly"
          >
            {{ showUnassignedOnly ? 'กลับมาดูตารางหลัก' : 'ดูเฉพาะงานที่ยังไม่ assign' }}
          </button>
        </div>
      </div>
    </section>

    <section
      v-if="showUnassignedOnly"
      class="mb-6 rounded-2xl border border-amber-200/80 bg-amber-50/80 p-5 shadow-sm dark:border-amber-500/30 dark:bg-amber-500/10"
    >
      <div class="mb-3 flex items-center justify-between gap-3">
        <h2 class="text-lg font-semibold text-amber-700 dark:text-amber-300">Unassigned Tasks</h2>
        <button class="text-xs underline" @click="showUnassignedOnly = false">ปิดมุมมองนี้</button>
      </div>

      <div v-if="unassignedTaskRows.length === 0" class="text-sm text-slate-600 dark:text-slate-300">ไม่พบงาน unassigned ตาม filter ปัจจุบัน</div>
      <div v-else class="space-y-3">
        <article
          v-for="row in unassignedTaskRows"
          :key="row.task.id"
          class="cursor-pointer rounded-xl border border-amber-200 bg-white/90 p-4 transition hover:border-amber-300 hover:shadow-sm dark:border-amber-500/30 dark:bg-slate-900/40"
          @click="goToTask(row.task)"
        >
          <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p class="font-semibold text-slate-900 dark:text-white">{{ row.task.title }}</p>
              <p class="text-xs text-slate-500 dark:text-slate-400">{{ row.task.status }} · {{ (row.hours).toFixed(1) }}h · {{ row.projectName }}</p>
            </div>
            <div class="flex items-center gap-3">
              <div class="text-sm text-slate-700 dark:text-slate-200">
                <p class="font-semibold">แนะนำผู้รับ: {{ row.suggested?.name || 'ไม่พบผู้รับที่เหมาะสม' }}</p>
                <p class="text-xs text-slate-500 dark:text-slate-400" v-if="row.suggested">Load {{ row.suggested.currentLoad }}% → {{ row.suggested.afterLoad }}%</p>
              </div>
              <button
                class="rounded-lg bg-emerald-600 px-3 py-2 text-xs font-semibold text-white transition hover:bg-emerald-500 disabled:opacity-50"
                :disabled="!row.suggested"
                @click.stop="openAutoAssignPreview(row)"
              >
                Apply Auto-assign
              </button>
            </div>
          </div>
        </article>
      </div>
    </section>

    <section
      v-if="openRecommendations && recommendations.length > 0"
      class="mb-6 rounded-2xl border border-orange-200/80 bg-orange-50/80 p-5 shadow-sm dark:border-orange-500/30 dark:bg-orange-500/10"
    >
      <div class="mb-3 flex items-center justify-between gap-3">
        <h2 class="text-lg font-semibold text-orange-700 dark:text-orange-300">Smart Rebalance Suggestions</h2>
        <span class="text-xs text-slate-600 dark:text-slate-300">มี preview ก่อนยืนยันทุกครั้ง</span>
      </div>

      <div class="space-y-3">
        <article
          v-for="rec in recommendations"
          :key="rec.id"
          class="rounded-xl border border-orange-200 bg-white/90 p-4 dark:border-orange-500/30 dark:bg-slate-900/40"
        >
          <div class="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
            <div>
              <p class="font-semibold text-slate-900 dark:text-white">ย้ายงาน <span class="text-orange-600 dark:text-orange-300">{{ rec.taskTitle }}</span></p>
              <p class="text-sm text-slate-600 dark:text-slate-300">{{ rec.fromName }} → {{ rec.toName }}</p>
              <p class="mt-1 text-xs text-slate-500 dark:text-slate-400">{{ rec.reason }}</p>
            </div>
            <div class="flex items-center gap-4">
              <div class="text-xs text-slate-600 dark:text-slate-300">
                <p>From: <span class="font-semibold text-rose-500">{{ rec.fromLoadBefore }}%</span> → <span class="font-semibold">{{ rec.fromLoadAfter }}%</span></p>
                <p>To: <span class="font-semibold text-emerald-500">{{ rec.toLoadBefore }}%</span> → <span class="font-semibold">{{ rec.toLoadAfter }}%</span></p>
              </div>
              <button
                class="rounded-lg bg-orange-600 px-3 py-2 text-xs font-semibold text-white transition hover:bg-orange-500 disabled:opacity-50"
                :disabled="!rec.toUserId"
                @click="openApplyPreview(rec)"
              >
                Apply Rebalance
              </button>
            </div>
          </div>
        </article>
      </div>
    </section>

    <section class="rounded-2xl border border-slate-200/80 bg-white/90 shadow-sm dark:border-white/10 dark:bg-slate-900/70">
      <div class="flex items-center justify-between border-b border-slate-200/80 px-5 py-4 dark:border-white/10">
        <h2 class="text-lg font-semibold">People-first Workload Matrix</h2>
        <p class="text-xs text-slate-500 dark:text-slate-400">เรียงตามคนที่โหลดสูงสุด</p>
      </div>

      <div v-if="loading" class="px-5 py-8 text-sm text-slate-500 dark:text-slate-300">กำลังประมวลผลข้อมูลภาระงาน...</div>
      <div v-else-if="sortedMembers.length === 0" class="px-5 py-8 text-sm text-slate-500 dark:text-slate-300">ยังไม่มีข้อมูลตรงกับ filter ที่เลือก</div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full text-sm">
          <thead class="bg-slate-50/90 text-left text-xs uppercase text-slate-500 dark:bg-slate-800/60 dark:text-slate-400">
            <tr>
              <th class="px-5 py-3 font-medium">Employee</th>
              <th class="px-5 py-3 font-medium">Load %</th>
              <th class="px-5 py-3 font-medium">Planned/Actual</th>
              <th class="px-5 py-3 font-medium">Task Count</th>
              <th class="px-5 py-3 font-medium">Risk</th>
              <th class="px-5 py-3 font-medium">Complexity Mix</th>
              <th class="px-5 py-3 font-medium">Action</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="member in sortedMembers" :key="member.key" class="border-t border-slate-100 dark:border-white/5">
              <td class="px-5 py-4">
                <p class="font-semibold">{{ member.name }}</p>
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ member.role }}</p>
              </td>
              <td class="px-5 py-4">
                <div class="flex items-center gap-3">
                  <div class="h-2 w-28 overflow-hidden rounded-full bg-slate-200 dark:bg-slate-700">
                    <div class="h-full rounded-full" :class="loadTone(member.load).bar" :style="{ width: `${Math.min(member.load, 150) / 1.5}%` }" />
                  </div>
                  <span class="font-semibold" :class="loadTone(member.load).text">{{ member.load }}%</span>
                </div>
              </td>
              <td class="px-5 py-4 text-slate-600 dark:text-slate-300">{{ member.planned }}h / {{ member.actual }}h</td>
              <td class="px-5 py-4">{{ member.inProgress }} in progress / {{ member.pending }} pending</td>
              <td class="px-5 py-4">
                <span class="rounded-full px-2 py-1 text-xs font-semibold" :class="riskTone(member.risk)">{{ member.risk }}</span>
              </td>
              <td class="px-5 py-4 text-xs text-slate-600 dark:text-slate-300">L {{ member.complexity.low }} · M {{ member.complexity.medium }} · H {{ member.complexity.high }}</td>
              <td class="px-5 py-4">
                <button
                  class="rounded-lg border border-slate-300 px-3 py-1.5 text-xs font-semibold transition hover:border-violet-400 hover:text-violet-600 dark:border-white/20"
                  @click="openMemberDetail(member.key)"
                >
                  View
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </section>

    <section class="mt-6 rounded-2xl border border-slate-200/80 bg-white/90 p-4 shadow-sm dark:border-white/10 dark:bg-slate-900/70">
      <div class="mb-3 flex items-center justify-between">
        <h3 class="text-base font-semibold">Audit / Telemetry Log</h3>
        <button class="text-xs text-slate-500 underline" @click="clearAuditLogs">clear</button>
      </div>
      <div v-if="auditLogs.length === 0" class="text-sm text-slate-500 dark:text-slate-400">ยังไม่มีการตัดสินใจรีบาลานซ์</div>
      <ul v-else class="space-y-2 text-sm">
        <li v-for="item in auditLogs" :key="item.id" class="rounded-lg border border-slate-200 px-3 py-2 dark:border-white/10">
          <p class="font-medium">{{ item.action }} <span class="text-xs text-slate-500">({{ item.status }})</span></p>
          <p class="text-xs text-slate-500 dark:text-slate-400">{{ item.detail }} · {{ formatDateTime(item.createdAt) }}</p>
        </li>
      </ul>
    </section>

    <Teleport to="body">
      <div v-if="previewRecommendation" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" @click.self="previewRecommendation = null">
        <div class="w-full max-w-lg rounded-2xl border border-slate-200 bg-white p-5 shadow-2xl dark:border-white/10 dark:bg-slate-900">
          <h3 class="text-lg font-semibold">ยืนยัน Apply Rebalance</h3>
          <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">ระบบจะเปลี่ยน assignee ของงานนี้ผ่าน API จริง</p>

          <div class="mt-4 space-y-2 rounded-xl bg-slate-50 p-4 text-sm dark:bg-slate-800/60">
            <p><span class="font-semibold">Task:</span> {{ previewRecommendation.taskTitle }}</p>
            <p><span class="font-semibold">Move:</span> {{ previewRecommendation.fromName }} → {{ previewRecommendation.toName }}</p>
            <p><span class="font-semibold">Impact:</span> From {{ previewRecommendation.fromLoadBefore }}% → {{ previewRecommendation.fromLoadAfter }}%, To {{ previewRecommendation.toLoadBefore }}% → {{ previewRecommendation.toLoadAfter }}%</p>
          </div>

          <div class="mt-4 flex justify-end gap-2">
            <button class="rounded-lg border border-slate-300 px-3 py-2 text-sm dark:border-white/20" @click="previewRecommendation = null">ยกเลิก</button>
            <button
              class="rounded-lg bg-orange-600 px-3 py-2 text-sm font-semibold text-white hover:bg-orange-500 disabled:opacity-50"
              :disabled="applyingRecommendation"
              @click="confirmApplyRebalance"
            >
              {{ applyingRecommendation ? 'กำลัง apply...' : 'ยืนยันเปลี่ยน assignee' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="previewAutoAssign" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4" @click.self="previewAutoAssign = null">
        <div class="w-full max-w-lg rounded-2xl border border-emerald-200 bg-white p-5 text-slate-900 shadow-2xl dark:border-emerald-500/40 dark:bg-slate-900 dark:text-slate-100">
          <h3 class="text-lg font-semibold text-slate-900 dark:text-slate-100">ยืนยัน Apply Auto-assign</h3>
          <p class="mt-2 text-sm text-slate-600 dark:text-slate-300">ระบบจะ assign งานนี้ให้ผู้รับที่แนะนำผ่าน API จริง</p>

          <div class="mt-4 space-y-2 rounded-xl bg-slate-50 p-4 text-sm text-slate-800 dark:bg-slate-800/60 dark:text-slate-200">
            <p><span class="font-semibold">Task:</span> {{ previewAutoAssign.task.title }}</p>
            <p><span class="font-semibold">Project:</span> {{ previewAutoAssign.projectName }}</p>
            <p><span class="font-semibold">Assign to:</span> {{ previewAutoAssign.suggested?.name || '-' }}</p>
            <p v-if="previewAutoAssign.suggested"><span class="font-semibold">Impact:</span> {{ previewAutoAssign.suggested.currentLoad }}% → {{ previewAutoAssign.suggested.afterLoad }}%</p>
          </div>

          <div class="mt-4 flex justify-end gap-2">
            <button class="rounded-lg border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-100 dark:border-white/20 dark:text-slate-200 dark:hover:bg-slate-800" @click="previewAutoAssign = null">ยกเลิก</button>
            <button
              class="rounded-lg bg-emerald-600 px-3 py-2 text-sm font-semibold text-white hover:bg-emerald-500 disabled:opacity-50"
              :disabled="applyingAutoAssign || !previewAutoAssign.suggested"
              @click="confirmApplyAutoAssign"
            >
              {{ applyingAutoAssign ? 'กำลัง apply...' : 'ยืนยัน assign' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="selectedMember" class="fixed inset-0 z-40 flex justify-end bg-black/40" @click.self="closeMemberDetail">
        <aside class="h-full w-full max-w-lg overflow-y-auto border-l border-slate-200 bg-white p-5 text-slate-900 dark:border-white/10 dark:bg-slate-900 dark:text-slate-100">
          <div class="mb-4 flex items-start justify-between">
            <div>
              <h3 class="text-lg font-semibold">{{ selectedMember.name }}</h3>
              <p class="text-sm text-slate-500 dark:text-slate-400">{{ selectedMember.role }}</p>
            </div>
            <button class="text-sm text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-200" @click="closeMemberDetail">ปิด</button>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div class="rounded-xl border border-slate-200 p-3 dark:border-white/10 dark:bg-slate-800/50">
              <p class="text-xs text-slate-500 dark:text-slate-400">Current Load</p>
              <p class="text-xl font-bold" :class="loadTone(selectedMember.load).text">{{ selectedMember.load }}%</p>
            </div>
            <div class="rounded-xl border border-slate-200 p-3 dark:border-white/10 dark:bg-slate-800/50">
              <p class="text-xs text-slate-500 dark:text-slate-400">Tasks</p>
              <p class="text-xl font-bold">{{ selectedMember.tasks.length }}</p>
            </div>
          </div>

          <div class="mt-4 rounded-xl border border-slate-200 p-3 dark:border-white/10 dark:bg-slate-800/50">
            <p class="mb-2 text-sm font-semibold">Risk Timeline</p>
            <ul class="space-y-1 text-sm text-slate-600 dark:text-slate-300">
              <li>Overdue: {{ memberRiskTimeline.overdue }}</li>
              <li>Today: {{ memberRiskTimeline.today }}</li>
              <li>1-2 Days: {{ memberRiskTimeline.next2Days }}</li>
              <li>3-7 Days: {{ memberRiskTimeline.next7Days }}</li>
            </ul>
          </div>

          <div class="mt-4 rounded-xl border border-slate-200 p-3 dark:border-white/10 dark:bg-slate-800/50">
            <p class="mb-2 text-sm font-semibold">Task Breakdown</p>
            <ul class="space-y-2 text-sm">
              <li
                v-for="task in selectedMember.tasks"
                :key="task.id"
                class="cursor-pointer rounded-lg bg-slate-50 p-2 transition hover:bg-slate-100 dark:bg-slate-800/60 dark:hover:bg-slate-700/70"
                @click="goToTask(task)"
              >
                <p class="font-medium">{{ task.title }}</p>
                <p class="text-xs text-slate-500 dark:text-slate-400">{{ task.status }} · {{ ((task.estimated_minutes || 0) / 60).toFixed(1) }}h</p>
              </li>
            </ul>
          </div>
        </aside>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useTasksApi, type GlobalActiveTask, type TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'
import { isTaskOverdueForMetrics } from '~/utils/task-overdue-metrics'
import { useProjectsApi, type Project } from '~/core/modules/projects/infrastructure/projects-api'
import type { User } from '~/core/modules/auth/infrastructure/auth-api'

definePageMeta({
  layout: 'default',
  middleware: [
    'auth',
    () => {
      const { currentUser } = useAuth()
      const role = currentUser.value?.role
      if (!['CEO', 'MANAGER'].includes(role ?? '')) {
        return navigateTo('/dashboard')
      }
    }
  ]
})

type RiskLevel = 'Low' | 'Medium' | 'High'
type RangeKey = 'week' | 'month' | 'quarter'

interface MemberRow {
  key: string
  userId: number | null
  name: string
  role: string
  load: number
  planned: number
  actual: number
  inProgress: number
  pending: number
  risk: RiskLevel
  complexity: { low: number; medium: number; high: number }
  tasks: GlobalActiveTask[]
}

interface Recommendation {
  id: string
  taskId: string
  toUserId: number | null
  taskTitle: string
  fromName: string
  toName: string
  reason: string
  fromLoadBefore: number
  fromLoadAfter: number
  toLoadBefore: number
  toLoadAfter: number
}

interface AuditLog {
  id: string
  action: string
  detail: string
  status: 'info' | 'success' | 'error'
  createdAt: string
}

interface AutoAssignRow {
  task: GlobalActiveTask
  hours: number
  projectName: string
  suggested: {
    userId: number
    name: string
    currentLoad: number
    afterLoad: number
  } | null
}

const { fetchWithAuth } = useAuth()
const { confirm } = useNotification()
const tasksApi = useTasksApi()
const projectsApi = useProjectsApi()
const route = useRoute()
const router = useRouter()

const selectedRange = ref<RangeKey>('week')
const selectedRole = ref('all')
const selectedProjectId = ref('all')
const selectedMemberKey = ref<string | null>(null)
const showUnassignedOnly = ref(false)
const loading = ref(false)
const error = ref('')
const openRecommendations = ref(true)
const applyingRecommendation = ref(false)

const previewRecommendation = ref<Recommendation | null>(null)
const previewAutoAssign = ref<AutoAssignRow | null>(null)
const applyingAutoAssign = ref(false)
const applyingAllAutoAssign = ref(false)
const auditLogs = ref<AuditLog[]>([])

const allActiveTasks = ref<GlobalActiveTask[]>([])
const allUsers = ref<User[]>([])
const projects = ref<Project[]>([])
const userTimeLogMinutesMap = ref<Record<number, number>>({})

const rangeCapacityHours: Record<RangeKey, number> = { week: 40, month: 160, quarter: 480 }
const AUDIT_STORAGE_KEY = 'sentinel-workload-audit-v1'

const userMap = computed(() => {
  const map = new Map<number, User>()
  for (const u of allUsers.value) map.set(u.id, u)
  return map
})

function getRangeDateWindow(range: RangeKey) {
  const now = new Date()
  const start = new Date(now)
  const end = new Date(now)

  if (range === 'week') {
    start.setHours(0, 0, 0, 0)
    end.setDate(end.getDate() + 7)
  }
  if (range === 'month') {
    start.setDate(1)
    start.setHours(0, 0, 0, 0)
    end.setMonth(end.getMonth() + 1)
    end.setDate(0)
    end.setHours(23, 59, 59, 999)
  }
  if (range === 'quarter') {
    const currentQuarter = Math.floor(now.getMonth() / 3)
    start.setMonth(currentQuarter * 3, 1)
    start.setHours(0, 0, 0, 0)
    end.setMonth(currentQuarter * 3 + 3, 0)
    end.setHours(23, 59, 59, 999)
  }

  return { start, end }
}

function isWithinRange(iso: string | null | undefined, range: RangeKey) {
  if (!iso) return false
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return false
  const { start, end } = getRangeDateWindow(range)
  return d >= start && d <= end
}

async function buildUserTimeLogMinutesMap(tasks: GlobalActiveTask[], range: RangeKey) {
  const map: Record<number, number> = {}

  await Promise.all(
    tasks.map(async (task) => {
      try {
        const logs = await tasksApi.getTimeLogs(task.id)
        for (const log of logs) {
          const userId = Number(log.user_id)
          if (!userId) continue
          const inRange = isWithinRange(log.logged_date || log.logged_at, range)
          if (!inRange) continue
          map[userId] = (map[userId] || 0) + (log.minutes || 0)
        }
      } catch {
        // ignore per-task log fetch failures
      }
    })
  )

  return map
}

const roleOptions = computed(() => {
  const hiddenRoles = new Set(['CEO', 'SUPPORT'])
  const set = new Set<string>()

  for (const u of allUsers.value) {
    const role = (u.role || '').trim()
    if (!role) continue
    if (hiddenRoles.has(role.toUpperCase())) continue
    set.add(role)
  }

  return [...set].sort()
})

const tasksInRange = computed(() => {
  const now = new Date()
  const end = new Date(now)
  if (selectedRange.value === 'week') end.setDate(end.getDate() + 7)
  if (selectedRange.value === 'month') end.setMonth(end.getMonth() + 1)
  if (selectedRange.value === 'quarter') end.setMonth(end.getMonth() + 3)

  return allActiveTasks.value.filter((task) => {
    const sourceDate = task.due_at || task.end_date || task.updated_at || task.created_at
    if (!sourceDate) return true
    const d = new Date(sourceDate)
    if (Number.isNaN(d.getTime())) return true
    return d <= end
  })
})

const filteredTasks = computed(() => {
  return tasksInRange.value.filter((task) => {
    const projectOk = selectedProjectId.value === 'all' || task.project_id === selectedProjectId.value
    return projectOk
  })
})

const sortedMembers = computed<MemberRow[]>(() => {
  const byUser = new Map<string, MemberRow>()
  const hiddenRoles = new Set(['CEO', 'SUPPORT'])

  // Seed all users first so everyone appears even with 0 tasks
  for (const user of allUsers.value) {
    if (hiddenRoles.has((user.role || '').toUpperCase())) continue
    if (selectedRole.value !== 'all' && user.role !== selectedRole.value) continue
    const key = `id:${user.id}`
    byUser.set(key, {
      key,
      userId: user.id,
      name: user.display_name || user.email || `User #${user.id}`,
      role: user.role || 'Team Member',
      load: 0,
      planned: 0,
      actual: 0,
      inProgress: 0,
      pending: 0,
      risk: 'Low',
      complexity: { low: 0, medium: 0, high: 0 },
      tasks: []
    })
  }

  // Aggregate tasks only into assigned users (unassigned shown in alert card)
  for (const task of filteredTasks.value) {
    if (!task.assigned_to) continue
    const key = `id:${task.assigned_to}`
    const user = userMap.value.get(task.assigned_to)

    if (hiddenRoles.has((user?.role || '').toUpperCase())) continue
    if (selectedRole.value !== 'all' && user?.role !== selectedRole.value) continue

    const current = byUser.get(key) || {
      key,
      userId: task.assigned_to,
      name: user?.display_name || task.assigned_to_display_name || task.assigned_to_email || `User #${task.assigned_to}`,
      role: user?.role || 'Team Member',
      load: 0,
      planned: 0,
      actual: 0,
      inProgress: 0,
      pending: 0,
      risk: 'Low' as RiskLevel,
      complexity: { low: 0, medium: 0, high: 0 },
      tasks: []
    }

    const hours = Math.max(0, (task.estimated_minutes || 0) / 60)
    current.planned += hours
    current.tasks.push(task)

    if (task.status === 'IN_PROGRESS' || task.status === 'READY_FOR_TEST' || task.status === 'REVIEW_PENDING') current.inProgress += 1
    if (task.status === 'PENDING' || task.status === 'BLOCKED') current.pending += 1

    if (hours <= 2) current.complexity.low += 1
    else if (hours <= 6) current.complexity.medium += 1
    else current.complexity.high += 1

    const due = task.due_at ? new Date(task.due_at) : null
    if (due && !Number.isNaN(due.getTime())) {
      const days = Math.ceil((due.getTime() - Date.now()) / 86_400_000)
      if (days < 0 && task.status !== 'COMPLETED') current.risk = 'High'
      else if (days <= 2 && current.risk !== 'High') current.risk = 'Medium'
    }

    byUser.set(key, current)
  }

  return [...byUser.values()]
    .map((m) => {
      const loggedHours = (userTimeLogMinutesMap.value[m.userId || 0] || 0) / 60
      const plannedHours = Number(m.planned.toFixed(1))
      const actualHours = Number(loggedHours.toFixed(1))
      const capacity = rangeCapacityHours[selectedRange.value]
      const load = capacity > 0 ? Math.round((plannedHours / capacity) * 100) : 0

      return {
        ...m,
        planned: plannedHours,
        actual: actualHours,
        load
      }
    })
    .sort((a, b) => {
      if (a.load !== b.load) return b.load - a.load

      const aActualZero = a.actual === 0
      const bActualZero = b.actual === 0
      if (aActualZero && bActualZero) return b.tasks.length - a.tasks.length

      return b.actual - a.actual
    })
})

const selectedMember = computed(() => sortedMembers.value.find((m) => m.key === selectedMemberKey.value) || null)

const memberRiskTimeline = computed(() => {
  const m = selectedMember.value
  const output = { overdue: 0, today: 0, next2Days: 0, next7Days: 0 }
  if (!m) return output

  for (const t of m.tasks) {
    if (!t.due_at || t.status === 'COMPLETED') continue
    const days = Math.ceil((new Date(t.due_at).getTime() - Date.now()) / 86_400_000)
    if (days < 0) {
      if (isTaskOverdueForMetrics(t)) output.overdue += 1
      continue
    }
    if (days === 0) output.today += 1
    else if (days <= 2) output.next2Days += 1
    else if (days <= 7) output.next7Days += 1
  }

  return output
})

const kpis = computed(() => {
  const members = sortedMembers.value
  const overloaded = members.filter((m) => m.load > 110).length
  const underutilized = members.filter((m) => m.load < 70).length
  const riskTasks = filteredTasks.value.filter((t) => {
    if (!t.due_at || t.status === 'COMPLETED') return false
    const days = (new Date(t.due_at).getTime() - Date.now()) / 86_400_000
    if (days < 0) return isTaskOverdueForMetrics(t)
    return days <= 2
  }).length
  const totalActual = members.reduce((sum, m) => sum + m.actual, 0)
  const totalPlanned = members.reduce((sum, m) => sum + m.planned, 0)
  const capacityPct = totalPlanned > 0 ? Math.round((totalActual / totalPlanned) * 100) : 0

  return [
    { label: 'Overloaded', value: `${overloaded} คน`, meta: 'Load > 110%', tone: 'text-rose-500' },
    { label: 'Underutilized', value: `${underutilized} คน`, meta: 'Load < 70%', tone: 'text-emerald-500' },
    { label: 'SLA Risk Tasks', value: `${riskTasks} งาน`, meta: 'ใกล้ครบกำหนดหรือเลยกำหนด', tone: 'text-amber-500' },
    { label: 'Team Capacity', value: `${capacityPct}%`, meta: 'Actual เทียบ Planned ตามช่วงเวลา', tone: 'text-sky-500' }
  ]
})

const unassignedTasks = computed(() => filteredTasks.value.filter((task) => !task.assigned_to))
const unassignedTasksCount = computed(() => unassignedTasks.value.length)

const autoAssignRecommendations = computed(() => {
  const receivers = sortedMembers.value
    .filter((m) => m.userId && m.load < 90)
    .sort((a, b) => a.load - b.load)

  return unassignedTasks.value.map((task) => {
    const taskHours = Math.max(0, (task.estimated_minutes || 0) / 60)
    const receiver = receivers.find((r) => {
      const nextLoad = Math.round(((r.planned + taskHours) / rangeCapacityHours[selectedRange.value]) * 100)
      return nextLoad <= 100
    })

    return {
      task,
      hours: taskHours,
      suggested: receiver
        ? {
            userId: receiver.userId as number,
            name: receiver.name,
            currentLoad: receiver.load,
            afterLoad: Math.round(((receiver.planned + taskHours) / rangeCapacityHours[selectedRange.value]) * 100)
          }
        : null
    }
  })
})

const unassignedTaskRows = computed<AutoAssignRow[]>(() => {
  const projectNameMap = new Map(projects.value.map((p) => [p.id, p.name]))
  return autoAssignRecommendations.value.map((r) => ({
    ...r,
    projectName: r.task.project_id ? (projectNameMap.get(r.task.project_id) || 'Unknown Project') : 'No Project'
  }))
})

const applyableAutoAssignRows = computed(() => unassignedTaskRows.value.filter((r) => !!r.suggested))

const recommendations = computed<Recommendation[]>(() => {
  const rows = sortedMembers.value
  const overloaded = rows.filter((r) => r.load > 110)
  const available = rows.filter((r) => r.load < 80 && r.userId)
  if (!overloaded.length || !available.length) return []

  const output: Recommendation[] = []
  for (const from of overloaded) {
    const candidateTask = [...from.tasks].filter((t) => t.status !== 'COMPLETED').sort((a, b) => (b.estimated_minutes || 0) - (a.estimated_minutes || 0))[0]
    if (!candidateTask) continue

    const hours = (candidateTask.estimated_minutes || 0) / 60
    const receiver = [...available].filter((r) => r.userId !== from.userId).sort((a, b) => a.load - b.load)[0]
    if (!receiver) continue

    const planned = rangeCapacityHours[selectedRange.value]
    const fromAfter = Math.round(((from.planned - hours) / planned) * 100)
    const toAfter = Math.round(((receiver.planned + hours) / planned) * 100)
    if (toAfter > 100) continue

    output.push({
      id: `${from.key}-${receiver.key}-${candidateTask.id}`,
      taskId: candidateTask.id,
      toUserId: receiver.userId,
      taskTitle: candidateTask.title,
      fromName: from.name,
      toName: receiver.name,
      reason: `ผู้รับมี load ต่ำกว่าและงานนี้ใช้เวลาประมาณ ${hours.toFixed(1)} ชั่วโมง จึงช่วยลดคอขวดโดยไม่ทำให้ผู้รับ overload`,
      fromLoadBefore: from.load,
      fromLoadAfter: fromAfter,
      toLoadBefore: receiver.load,
      toLoadAfter: toAfter
    })

    if (output.length >= 5) break
  }

  return output
})

function persistAuditLogs() {
  if (!import.meta.client) return
  localStorage.setItem(AUDIT_STORAGE_KEY, JSON.stringify(auditLogs.value.slice(0, 100)))
}

function trackEvent(action: string, detail: string, status: AuditLog['status'] = 'info') {
  if (!action.startsWith('apply_')) return

  auditLogs.value.unshift({
    id: crypto.randomUUID(),
    action,
    detail,
    status,
    createdAt: new Date().toISOString()
  })
  persistAuditLogs()
}

function openApplyPreview(rec: Recommendation) {
  previewRecommendation.value = rec
  trackEvent('preview_rebalance', `${rec.taskTitle}: ${rec.fromName} -> ${rec.toName}`, 'info')
}

function openAutoAssignPreview(row: AutoAssignRow) {
  if (!row.suggested) return
  previewAutoAssign.value = row
  trackEvent('preview_auto_assign', `${row.task.title} -> ${row.suggested.name}`, 'info')
}

async function confirmApplyAutoAssign() {
  const row = previewAutoAssign.value
  if (!row?.suggested) return

  applyingAutoAssign.value = true
  try {
    await tasksApi.assignTask(row.task.id, row.suggested.userId)
    trackEvent('apply_auto_assign_success', `${row.task.title} -> ${row.suggested.name}`, 'success')
    previewAutoAssign.value = null
    await loadWorkload()
  } catch (e: any) {
    trackEvent('apply_auto_assign_failed', e?.message || 'unknown error', 'error')
    error.value = e?.message || 'ไม่สามารถ auto-assign งานได้'
  } finally {
    applyingAutoAssign.value = false
  }
}

async function openApplyAllAutoAssignPreview() {
  const count = applyableAutoAssignRows.value.length
  if (count === 0) return
  const ok = await confirm({
    title: 'ยืนยัน Apply Auto-assign ทั้งหมด',
    message: `จะทำการ assign งานทั้งหมด ${count} งาน ตามผู้รับที่ระบบแนะนำ ดำเนินการต่อหรือไม่?`,
    confirmLabel: 'ยืนยัน',
    cancelLabel: 'ยกเลิก',
    variant: 'primary'
  })
  if (!ok) return

  applyingAllAutoAssign.value = true
  let success = 0
  let failed = 0

  for (const row of applyableAutoAssignRows.value) {
    if (!row.suggested) continue
    try {
      await tasksApi.assignTask(row.task.id, row.suggested.userId)
      success += 1
      trackEvent('apply_auto_assign_success', `${row.task.title} -> ${row.suggested.name}`, 'success')
    } catch (e: any) {
      failed += 1
      trackEvent('apply_auto_assign_failed', `${row.task.title} -> ${e?.message || 'unknown error'}`, 'error')
    }
  }

  applyingAllAutoAssign.value = false
  trackEvent('apply_auto_assign_bulk', `success=${success}, failed=${failed}`, failed > 0 ? 'error' : 'success')
  await loadWorkload()
}

async function confirmApplyRebalance() {
  const rec = previewRecommendation.value
  if (!rec || !rec.toUserId) return

  applyingRecommendation.value = true
  try {
    await tasksApi.assignTask(rec.taskId, rec.toUserId)
    trackEvent('apply_rebalance_success', `${rec.taskTitle}: ${rec.fromName} -> ${rec.toName}`, 'success')
    previewRecommendation.value = null
    await loadWorkload()
  } catch (e: any) {
    trackEvent('apply_rebalance_failed', e?.message || 'unknown error', 'error')
    error.value = e?.message || 'ไม่สามารถเปลี่ยน assignee ได้'
  } finally {
    applyingRecommendation.value = false
  }
}

function goToTask(task: { id: string; code?: string }) {
  navigateTo({
    path: `/task/${task.id}`,
    query: { from: 'workload' }
  })
}

function syncSelectedMemberFromRoute() {
  const member = route.query.member
  selectedMemberKey.value = typeof member === 'string' && member ? member : null
}

function openMemberDetail(memberKey: string) {
  selectedMemberKey.value = memberKey
  const nextQuery = { ...route.query, member: memberKey }
  router.replace({ query: nextQuery })

  const m = sortedMembers.value.find((x) => x.key === memberKey)
  if (m) trackEvent('open_member_detail', `${m.name} (${m.load}%)`, 'info')
}

function closeMemberDetail() {
  selectedMemberKey.value = null
  const nextQuery = { ...route.query }
  delete nextQuery.member
  router.replace({ query: nextQuery })
}

function clearAuditLogs() {
  auditLogs.value = []
  persistAuditLogs()
}

async function loadWorkload() {
  loading.value = true
  error.value = ''
  try {
    const [tasks, usersRes, projectData] = await Promise.all([
      tasksApi.getTeamActiveTasks(),
      fetchWithAuth<{ data: User[] }>('/auth/users'),
      projectsApi.getProjects()
    ])

    allActiveTasks.value = tasks
    allUsers.value = usersRes.data || []
    projects.value = projectData
    userTimeLogMinutesMap.value = await buildUserTimeLogMinutesMap(tasks, selectedRange.value)
    trackEvent('load_workload', `tasks=${tasks.length}, users=${allUsers.value.length}`, 'info')
  } catch (e: any) {
    error.value = e?.message || 'ไม่สามารถโหลดข้อมูล workload ได้'
    trackEvent('load_workload_failed', error.value, 'error')
  } finally {
    loading.value = false
  }
}

function loadTone(load: number) {
  if (load > 110) return { text: 'text-rose-500', bar: 'bg-rose-500' }
  if (load > 90) return { text: 'text-amber-500', bar: 'bg-amber-500' }
  if (load > 70) return { text: 'text-sky-500', bar: 'bg-sky-500' }
  return { text: 'text-emerald-500', bar: 'bg-emerald-500' }
}

function riskTone(risk: RiskLevel) {
  if (risk === 'High') return 'bg-rose-100 text-rose-600 dark:bg-rose-500/20 dark:text-rose-300'
  if (risk === 'Medium') return 'bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300'
  return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300'
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString('th-TH', { dateStyle: 'short', timeStyle: 'short' })
}

onMounted(() => {
  syncSelectedMemberFromRoute()

  if (import.meta.client) {
    try {
      const raw = localStorage.getItem(AUDIT_STORAGE_KEY)
      if (raw) {
        const parsed = JSON.parse(raw) as AuditLog[]
        auditLogs.value = parsed.filter((item) => item.action.startsWith('apply_'))
      }
    } catch {
      auditLogs.value = []
    }
  }
  loadWorkload()
})

watch([selectedRange, selectedRole, selectedProjectId], () => {
  loadWorkload()
})

watch(() => route.query.member, () => {
  syncSelectedMemberFromRoute()
})
</script>
