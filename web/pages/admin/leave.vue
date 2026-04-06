<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="max-w-screen-xl mx-auto flex items-center justify-between gap-3">
        <div>
          <h1 class="text-base font-bold">Leave Admin</h1>
          <p class="text-xs text-gray-500">จัดการนโยบายลา วันหยุด และติดตามการอนุมัติ</p>
        </div>
        <button class="btn" :disabled="loading" @click="refreshAll">รีเฟรช</button>
      </div>
    </div>

    <main class="max-w-screen-xl mx-auto px-6 py-8 space-y-6">
      <div class="flex gap-2 border-b border-gray-800 pb-2 overflow-x-auto">
        <button v-for="t in tabs" :key="t.id" class="tab" :class="activeTab === t.id ? 'tab-active' : ''" @click="activeTab = t.id">{{ t.label }}</button>
      </div>

      <section v-if="activeTab === 'policies'" class="card">
        <h2 class="section">นโยบายโควตาวันลา</h2>
        <form class="grid gap-3 md:grid-cols-5 items-end" @submit.prevent="savePolicy">
          <div>
            <label class="label">ประเภทลา</label>
            <select v-model="policyForm.leave_type" class="input">
              <option value="ANNUAL">ลาพักร้อน</option>
              <option value="SICK">ลาป่วย</option>
              <option value="PERSONAL">ลากิจ</option>
            </select>
          </div>
          <div>
            <label class="label">โควต้าต่อปี</label>
            <input v-model.number="policyForm.annual_quota_days" type="number" min="0" class="input" />
          </div>
          <div>
            <label class="label">วันยกยอดสูงสุด</label>
            <input v-model.number="policyForm.max_carry_forward_days" type="number" min="0" class="input" />
          </div>
          <div>
            <label class="label">สถานะ</label>
            <select v-model="policyForm.is_active" class="input">
              <option :value="true">ใช้งาน</option>
              <option :value="false">ปิดใช้งาน</option>
            </select>
          </div>
          <button class="btn-primary" :disabled="savingPolicy">{{ savingPolicy ? 'กำลังบันทึก...' : 'บันทึกนโยบาย' }}</button>
        </form>

        <div class="mt-5 overflow-x-auto">
          <table class="w-full text-sm">
            <thead><tr class="text-xs text-gray-500 border-b border-gray-700"><th class="py-2 text-left">ประเภทลา</th><th class="py-2 text-left">โควตา</th><th class="py-2 text-left">ยกยอด</th><th class="py-2 text-left">สถานะ</th></tr></thead>
            <tbody>
              <tr v-for="p in policies" :key="p.id" class="border-b border-gray-800/60">
                <td class="py-2">{{ leaveTypeLabel(p.leave_type) }}</td>
                <td class="py-2">{{ p.annual_quota_days }} วัน</td>
                <td class="py-2">{{ p.max_carry_forward_days }} วัน</td>
                <td class="py-2"><span class="text-xs" :class="p.is_active ? 'text-emerald-300' : 'text-gray-500'">{{ p.is_active ? 'ใช้งาน' : 'ปิด' }}</span></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section v-if="activeTab === 'holidays'" class="card">
        <h2 class="section">วันหยุดประจำปี</h2>
        <form class="grid gap-3 md:grid-cols-4 items-end" @submit.prevent="saveHoliday">
          <div>
            <label class="label">วันที่</label>
            <input v-model="holidayForm.date" type="date" class="input" required />
          </div>
          <div class="md:col-span-2">
            <label class="label">ชื่อวันหยุด</label>
            <input v-model="holidayForm.name" class="input" placeholder="เช่น วันขึ้นปีใหม่" required />
          </div>
          <button class="btn-primary" :disabled="savingHoliday">{{ savingHoliday ? 'กำลังบันทึก...' : 'เพิ่ม/อัปเดตวันหยุด' }}</button>
        </form>

        <div class="mt-5 space-y-2">
          <div v-for="h in holidays" :key="h.id" class="row">
            <span>{{ fmt(h.date) }}</span>
            <span class="text-gray-400">{{ h.name }}</span>
          </div>
        </div>
      </section>

      <section v-if="activeTab === 'approvals'" class="card">
        <h2 class="section">คำขอลารออนุมัติ</h2>
        <div v-if="pending.length === 0" class="text-sm text-gray-500 italic">ไม่มีคำขอค้างอนุมัติ</div>
        <div v-else class="space-y-3">
          <article v-for="r in pending" :key="r.id" class="rounded-xl border border-gray-700 bg-gray-900/40 p-4">
            <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-3">
              <div>
                <p class="text-sm font-semibold">{{ r.user_display_name || r.user_email || ('User #' + r.user_id) }}</p>
                <p class="text-xs text-gray-400 mt-0.5">{{ leaveTypeLabel(r.leave_type) }} · {{ r.days_requested }} วัน</p>
                <p class="text-xs text-gray-500 mt-1">{{ fmt(r.start_date) }} → {{ fmt(r.end_date) }}</p>
                <p class="text-xs text-gray-400 mt-2">{{ r.reason }}</p>
              </div>
              <div class="w-full md:w-[280px] space-y-2">
                <textarea v-model="reviewComment[r.id]" rows="2" class="input" placeholder="หมายเหตุการอนุมัติ/ไม่อนุมัติ" />
                <div class="flex justify-end gap-2">
                  <button class="btn-danger" :disabled="reviewingId === r.id" @click="review(r.id, 'REJECTED')">ไม่อนุมัติ</button>
                  <button class="btn-success" :disabled="reviewingId === r.id" @click="review(r.id, 'APPROVED')">อนุมัติ</button>
                  <button class="btn" @click="loadAudit(r.id)">Audit</button>
                </div>
              </div>
            </div>
          </article>
        </div>

        <div v-if="selectedAuditId" class="mt-6 rounded-xl border border-gray-700 bg-gray-900/30 p-4">
          <h3 class="text-sm font-semibold mb-2">Audit Log: Leave #{{ selectedAuditId }}</h3>
          <div class="space-y-2 text-xs">
            <div v-for="a in auditLogs" :key="a.id" class="row">
              <span class="text-gray-300">{{ a.action }} · {{ a.old_status || '-' }} → {{ a.new_status || '-' }}</span>
              <span class="text-gray-500">{{ fmtDateTime(a.created_at) }}</span>
            </div>
          </div>
        </div>
      </section>

      <section v-if="activeTab === 'backfill'" class="card space-y-6">
        <div>
          <h2 class="section">กรอกข้อมูลลาย้อนหลัง (รายบุคคล)</h2>
          <form class="grid gap-3 md:grid-cols-3" @submit.prevent="submitBackfillSingle">
            <div>
              <label class="label">อีเมลพนักงาน</label>
              <input v-model="backfillForm.employee_email" type="email" class="input" placeholder="employee@company.com" required />
            </div>
            <div>
              <label class="label">วันที่เริ่ม</label>
              <input v-model="backfillForm.start_date" type="date" class="input" required />
            </div>
            <div>
              <label class="label">วันที่สิ้นสุด</label>
              <input v-model="backfillForm.end_date" type="date" class="input" required />
            </div>
            <div>
              <label class="label">ประเภทลา</label>
              <select v-model="backfillForm.leave_type" class="input">
                <option value="ANNUAL">ลาพักร้อน</option>
                <option value="SICK">ลาป่วย</option>
                <option value="PERSONAL">ลากิจ</option>
              </select>
            </div>
            <div>
              <label class="label">สถานะ</label>
              <select v-model="backfillForm.status" class="input">
                <option value="APPROVED">อนุมัติแล้ว</option>
                <option value="PENDING">รออนุมัติ</option>
                <option value="REJECTED">ไม่อนุมัติ</option>
              </select>
            </div>
            <div>
              <label class="label">หมายเหตุผู้จัดการ</label>
              <input v-model="backfillForm.comment" class="input" placeholder="ใส่หรือเว้นว่างได้" />
            </div>
            <div class="md:col-span-3">
              <label class="label">เหตุผลการลา</label>
              <textarea v-model="backfillForm.reason" rows="2" class="input" placeholder="เช่น ลาพักผ่อนประจำปี" required />
            </div>
            <div class="md:col-span-3 flex justify-end">
              <button class="btn-primary" :disabled="backfillSubmitting">{{ backfillSubmitting ? 'กำลังบันทึก...' : 'บันทึกข้อมูลย้อนหลัง' }}</button>
            </div>
          </form>
        </div>

        <div>
          <div class="flex items-center justify-between gap-3">
            <h2 class="section mb-0">นำเข้าย้อนหลังแบบ Bulk (CSV)</h2>
            <button class="btn" type="button" @click="downloadBackfillTemplate">ดาวน์โหลด CSV Template</button>
          </div>
          <p class="text-xs text-gray-500 mb-2 mt-2">
            Header ที่รองรับ: employee_email,start_date,end_date,leave_type,status,reason,comment
          </p>
          <input class="input" type="file" accept=".csv,text/csv" :disabled="backfillBulkSubmitting" @change="onBulkCsvSelected" />

          <div v-if="backfillBulkSummary" class="mt-4 grid grid-cols-3 gap-3 text-sm">
            <div class="row"><span>ทั้งหมด</span><span>{{ backfillBulkSummary.total }}</span></div>
            <div class="row"><span class="text-emerald-300">สำเร็จ</span><span>{{ backfillBulkSummary.succeeded }}</span></div>
            <div class="row"><span class="text-red-300">ไม่สำเร็จ</span><span>{{ backfillBulkSummary.failed }}</span></div>
          </div>

          <div v-if="backfillBulkResults.length" class="mt-4 overflow-x-auto">
            <table class="w-full text-xs">
              <thead>
                <tr class="text-gray-500 border-b border-gray-700">
                  <th class="py-2 text-left">#</th>
                  <th class="py-2 text-left">Email</th>
                  <th class="py-2 text-left">ผลลัพธ์</th>
                  <th class="py-2 text-left">Leave ID</th>
                  <th class="py-2 text-left">Error</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="r in backfillBulkResults" :key="`${r.index}-${r.email}`" class="border-b border-gray-800/60">
                  <td class="py-2">{{ r.index + 1 }}</td>
                  <td class="py-2">{{ r.email }}</td>
                  <td class="py-2" :class="r.status === 'created' ? 'text-emerald-300' : 'text-red-300'">{{ r.status }}</td>
                  <td class="py-2">{{ r.leave_id || '-' }}</td>
                  <td class="py-2 text-red-300">{{ r.error || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <section v-if="activeTab === 'trend'" class="card">
        <h2 class="section">แนวโน้มการลา (รายเดือน)</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead><tr class="text-xs text-gray-500 border-b border-gray-700"><th class="py-2 text-left">เดือน</th><th class="py-2 text-left">ทีม</th><th class="py-2 text-left">ขอ</th><th class="py-2 text-left">อนุมัติ</th><th class="py-2 text-left">ไม่อนุมัติ</th><th class="py-2 text-left">รวมวันลา</th></tr></thead>
            <tbody>
              <tr v-for="t in trend" :key="`${t.month}-${t.team_name}`" class="border-b border-gray-800/60">
                <td class="py-2">{{ t.month }}</td>
                <td class="py-2">{{ t.team_name || 'Unassigned' }}</td>
                <td class="py-2">{{ t.requested }}</td>
                <td class="py-2 text-emerald-300">{{ t.approved }}</td>
                <td class="py-2 text-red-300">{{ t.rejected }}</td>
                <td class="py-2">{{ t.total_days }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type LeavePolicy, type HolidayCalendar, type LeaveRequest, type LeaveAuditLog, type LeaveTrendPoint, type LeaveBackfillBulkResultItem } from '~/core/modules/attendance/infrastructure/attendance-api'
import { useNotification } from '~/composables/useNotification'

definePageMeta({ layout: 'default', middleware: 'auth' })

const api = useAttendanceApi()
const { showSuccess, showError } = useNotification()

const loading = ref(false)
const savingPolicy = ref(false)
const savingHoliday = ref(false)
const reviewingId = ref<number | null>(null)

const tabs = [
  { id: 'policies', label: 'นโยบายวันลา' },
  { id: 'holidays', label: 'วันหยุดประจำปี' },
  { id: 'approvals', label: 'อนุมัติวันลา' },
  { id: 'backfill', label: 'กรอกย้อนหลัง' },
  { id: 'trend', label: 'รายงานแนวโน้ม' },
]
const activeTab = ref('policies')

const policies = ref<LeavePolicy[]>([])
const holidays = ref<HolidayCalendar[]>([])
const pending = ref<LeaveRequest[]>([])
const trend = ref<LeaveTrendPoint[]>([])
const auditLogs = ref<LeaveAuditLog[]>([])
const selectedAuditId = ref<number | null>(null)
const reviewComment = ref<Record<number, string>>({})
const backfillSubmitting = ref(false)
const backfillBulkSubmitting = ref(false)
const backfillBulkResults = ref<LeaveBackfillBulkResultItem[]>([])
const backfillBulkSummary = ref<{ total: number; succeeded: number; failed: number } | null>(null)

const policyForm = reactive({
  leave_type: 'ANNUAL' as 'ANNUAL' | 'SICK' | 'PERSONAL',
  annual_quota_days: 10,
  max_carry_forward_days: 0,
  is_active: true,
})

const holidayForm = reactive({
  date: '',
  name: '',
})

const backfillForm = reactive({
  employee_email: '',
  start_date: '',
  end_date: '',
  leave_type: 'ANNUAL' as 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID',
  status: 'APPROVED' as 'PENDING' | 'APPROVED' | 'REJECTED',
  reason: '',
  comment: '',
})

function leaveTypeLabel(t: string) {
  if (t === 'ANNUAL') return 'ลาพักร้อน'
  if (t === 'SICK') return 'ลาป่วย'
  if (t === 'PERSONAL') return 'ลากิจ'
  return t
}

const fmt = (s: string) => new Date(s).toLocaleDateString('th-TH', { day: '2-digit', month: 'short', year: 'numeric' })
const fmtDateTime = (s: string) => new Date(s).toLocaleString('th-TH', { dateStyle: 'medium', timeStyle: 'short' })

async function refreshAll() {
  loading.value = true
  try {
    const now = new Date()
    const from = new Date(now.getFullYear(), 0, 1).toISOString().slice(0, 10)
    const to = new Date(now.getFullYear(), 11, 31).toISOString().slice(0, 10)

    const [p, h, pend, tr] = await Promise.all([
      api.getLeavePolicies(),
      api.getHolidays(from, to),
      api.getPendingLeaveRequests(),
      api.getLeaveTrend(from, to),
    ])
    policies.value = p.items || []
    holidays.value = h.items || []
    pending.value = pend.items || []
    trend.value = tr.items || []
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'โหลดข้อมูล Leave Admin ไม่สำเร็จ')
  } finally {
    loading.value = false
  }
}

async function savePolicy() {
  savingPolicy.value = true
  try {
    await api.upsertLeavePolicy(policyForm)
    showSuccess('บันทึกนโยบายวันลาสำเร็จ')
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'บันทึกนโยบายไม่สำเร็จ')
  } finally {
    savingPolicy.value = false
  }
}

async function saveHoliday() {
  savingHoliday.value = true
  try {
    await api.upsertHoliday(holidayForm)
    showSuccess('บันทึกวันหยุดสำเร็จ')
    holidayForm.date = ''
    holidayForm.name = ''
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'บันทึกวันหยุดไม่สำเร็จ')
  } finally {
    savingHoliday.value = false
  }
}

async function review(id: number, status: 'APPROVED' | 'REJECTED') {
  reviewingId.value = id
  try {
    await api.reviewLeaveRequest(id, { status, comment: reviewComment.value[id] || '' })
    showSuccess(status === 'APPROVED' ? 'อนุมัติคำขอลาแล้ว' : 'ไม่อนุมัติคำขอลาแล้ว')
    pending.value = pending.value.filter(x => x.id !== id)
    await loadAudit(id)
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'อัปเดตสถานะคำขอลาไม่สำเร็จ')
  } finally {
    reviewingId.value = null
  }
}

async function loadAudit(leaveId: number) {
  selectedAuditId.value = leaveId
  try {
    const res = await api.getLeaveAudit(leaveId)
    auditLogs.value = res.items || []
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'โหลด Audit log ไม่สำเร็จ')
  }
}

async function submitBackfillSingle() {
  backfillSubmitting.value = true
  try {
    const item = {
      employee_email: backfillForm.employee_email.trim(),
      start_date: backfillForm.start_date,
      end_date: backfillForm.end_date,
      leave_type: backfillForm.leave_type,
      status: backfillForm.status,
      reason: backfillForm.reason.trim(),
      comment: backfillForm.comment.trim(),
    }
    const created = await api.backfillLeave({ item })
    showSuccess(`เพิ่มข้อมูลย้อนหลังสำเร็จ (Leave #${created.id})`)
    await loadAudit(created.id)
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'บันทึกย้อนหลังไม่สำเร็จ')
  } finally {
    backfillSubmitting.value = false
  }
}

function parseCsvRows(csvText: string) {
  const rows = csvText
    .split(/\r?\n/)
    .map(r => r.trim())
    .filter(Boolean)
  if (!rows.length) return []
  const header = rows[0].split(',').map(h => h.trim().toLowerCase())
  const required = ['employee_email', 'start_date', 'end_date', 'leave_type', 'status', 'reason', 'comment']
  const missing = required.filter(k => !header.includes(k))
  if (missing.length) {
    throw new Error(`CSV header missing: ${missing.join(', ')}`)
  }
  const idx = (k: string) => header.indexOf(k)
  const out = [] as Array<{ employee_email: string; start_date: string; end_date: string; leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'; status: 'PENDING' | 'APPROVED' | 'REJECTED'; reason: string; comment: string }>
  for (let i = 1; i < rows.length; i++) {
    const cols = rows[i].split(',').map(c => c.trim())
    out.push({
      employee_email: cols[idx('employee_email')] || '',
      start_date: cols[idx('start_date')] || '',
      end_date: cols[idx('end_date')] || '',
      leave_type: ((cols[idx('leave_type')] || 'ANNUAL').toUpperCase() as any),
      status: ((cols[idx('status')] || 'APPROVED').toUpperCase() as any),
      reason: cols[idx('reason')] || '',
      comment: cols[idx('comment')] || '',
    })
  }
  return out
}

async function onBulkCsvSelected(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  backfillBulkSubmitting.value = true
  try {
    const text = await file.text()
    const items = parseCsvRows(text)
    if (!items.length) throw new Error('CSV ไม่มีข้อมูลแถวรายการ')
    const res = await api.backfillLeaveBulk({ items })
    backfillBulkSummary.value = { total: res.total, succeeded: res.succeeded, failed: res.failed }
    backfillBulkResults.value = res.results || []
    if (res.failed > 0) {
      showError(`นำเข้าเสร็จ: สำเร็จ ${res.succeeded}, ไม่สำเร็จ ${res.failed}`)
    } else {
      showSuccess(`นำเข้า CSV สำเร็จทั้งหมด ${res.succeeded} รายการ`)
    }
    await refreshAll()
  } catch (err: any) {
    showError(err?.message || 'นำเข้า CSV ไม่สำเร็จ')
  } finally {
    backfillBulkSubmitting.value = false
    input.value = ''
  }
}

function downloadBackfillTemplate() {
  const lines = [
    'employee_email,start_date,end_date,leave_type,status,reason,comment',
    'alice@company.com,2026-01-15,2026-01-16,ANNUAL,APPROVED,ลาพักร้อนต้นปี,อนุมัติย้อนหลังโดย HR',
    'bob@company.com,2026-02-03,2026-02-03,SICK,APPROVED,ป่วยไข้หวัด,มีใบรับรองแพทย์',
    'carol@company.com,2026-03-10,2026-03-10,PERSONAL,PENDING,ธุระส่วนตัว,รอหัวหน้าอนุมัติ',
  ]
  const csv = `${lines.join('\n')}\n`
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'leave-backfill-template.csv'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

onMounted(refreshAll)
</script>

<style scoped>
.card { @apply rounded-2xl border border-gray-700 bg-gray-800/50 p-5; }
.section { @apply text-sm font-semibold mb-4; }
.label { @apply block text-xs font-semibold uppercase tracking-widest text-gray-500 mb-1.5; }
.input { @apply w-full rounded-lg border border-gray-700 bg-gray-800/80 px-3 py-2 text-sm text-white focus:border-cyan-500 focus:outline-none; }
.btn { @apply rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700; }
.btn-primary { @apply rounded-lg border border-cyan-500/40 bg-cyan-900/20 px-4 py-2 text-sm font-semibold text-cyan-300 hover:bg-cyan-900/40 disabled:opacity-50; }
.btn-success { @apply rounded-lg border border-emerald-500/40 bg-emerald-900/20 px-3 py-1.5 text-xs font-semibold text-emerald-300 hover:bg-emerald-900/40 disabled:opacity-50; }
.btn-danger { @apply rounded-lg border border-red-500/40 bg-red-900/20 px-3 py-1.5 text-xs font-semibold text-red-300 hover:bg-red-900/40 disabled:opacity-50; }
.tab { @apply px-3 py-1.5 rounded-lg text-xs font-semibold text-gray-400 hover:text-gray-200 whitespace-nowrap; }
.tab-active { @apply bg-gray-700 text-white; }
.row { @apply flex items-center justify-between rounded-lg border border-gray-700/70 px-3 py-2 bg-gray-900/30; }
</style>
