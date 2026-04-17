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
        <p v-if="isAdminLeaveManager" class="text-xs text-gray-500 mb-2">role `SUPPORT` และ `CEO` จัดการได้เท่ากัน</p>
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
            <UiDatePicker v-model="holidayForm.date" placeholder="เลือกวันที่…" />
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
        <p class="text-xs text-gray-400 mt-0.5">แสดงประเภทลาและช่วงเวลาที่พนักงานขอ</p>
        <div v-if="pending.length === 0" class="text-sm text-gray-500 italic">ไม่มีคำขอค้างอนุมัติ</div>
        <div v-else class="space-y-3">
          <article v-for="r in pending" :key="r.id" class="rounded-xl border border-gray-700 bg-gray-900/40 p-4">
            <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-3">
              <div>
                <p class="text-sm font-semibold">
                  {{ r.user_email || ('User #' + r.user_id) }}
                </p>
                <p class="text-xs text-gray-400 mt-0.5">{{ leaveTypeLabel(r.leave_type) }} · {{ requestedDurationLabel(r) }}</p>
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
              <UiDatePicker v-model="backfillForm.start_date" placeholder="วันที่เริ่ม…" />
            </div>
            <div>
              <label class="label">วันที่สิ้นสุด</label>
              <UiDatePicker v-model="backfillForm.end_date" placeholder="วันที่สิ้นสุด…" :disabled="backfillForm.is_half_day" />
            </div>
            <div>
              <label class="label">รูปแบบวันลา</label>
              <select v-model="backfillDurationMode" class="input">
                <option value="FULL">เต็มวัน</option>
                <option value="HALF_AM">ครึ่งวันเช้า</option>
                <option value="HALF_PM">ครึ่งวันบ่าย</option>
              </select>
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
            Header ที่รองรับ: employee_email,start_date,end_date,leave_type,is_half_day,half_day_session,status,reason,comment
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
        <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
          <div>
            <h2 class="section mb-1">ข้อมูลการลา</h2>
            <p class="text-xs text-gray-500">แสดงรายการคำขอลาทั้งหมดของพนักงาน</p>
            <p v-if="isAdminLeaveManager" class="text-xs text-gray-500 mt-1">role `SUPPORT` ใช้งานได้เหมือน `CEO` ทุกอย่างในหน้านี้</p>
          </div>
        </div>

        <form class="form-search mt-4 grid gap-3 md:grid-cols-4" @submit.prevent>
          <div>
            <label class="label">วันที่</label>
            <UiDatePicker v-model="leaveFilter.from" placeholder="จากวันที่…" />
          </div>
          <div>
            <label class="label">ถึง</label>
            <UiDatePicker v-model="leaveFilter.to" placeholder="ถึงวันที่…" />
          </div>
          <div>
            <label class="label">ประเภทลา</label>
            <select v-model="leaveFilter.leave_type" class="input">
              <option v-for="t in leaveTypeOptions" :key="t.value" :value="t.value">{{ t.label }}</option>
            </select>
          </div>
          <div>
            <label class="label">พนักงาน</label>
            <select v-model="leaveFilter.employee_key" class="input">
              <option value="ALL">ทั้งหมด</option>
              <option v-for="u in leaveEmployeeOptions" :key="u.key" :value="u.key">{{ u.label }}</option>
            </select>
          </div>
        </form>

        <div v-if="leaveRecords.length === 0" class="mt-4 text-sm text-gray-500 italic">ยังไม่มีข้อมูลการลา</div>
        <div v-else-if="filteredLeaveRecords.length === 0" class="mt-4 text-sm text-gray-500 italic">ไม่พบข้อมูลตามเงื่อนไขที่เลือก</div>
        <div v-else class="table-responsive mt-4">
          <div id="datatable_wrapper">
            <div class="dt-row">
              <div class="dt-col">
                <table class="leave-table w-full text-sm">
                  <thead>
                    <tr class="text-xs text-gray-400">
                      <th class="text-left">ลำดับ</th>
                      <th class="text-left">พนักงาน</th>
                      <th class="text-left">ประเภท</th>
                      <th class="text-left">วันที่ลา</th>
                      <th class="text-left">จำนวนวัน</th>
                      <th class="text-left">รายละเอียด</th>
                      <th class="text-left">สถานะ</th>
                      <th class="text-left">ประวัติการแก้ไข</th>
                      <th v-if="canManageLeaveRecords" class="text-left">จัดการ</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(r, idx) in filteredLeaveRecords" :key="r.id">
                      <td>{{ idx + 1 }}</td>
                      <td>{{ r.user_email || ('User #' + r.user_id) }}</td>
                      <template v-if="editingLeaveId === r.id">
                        <td>
                          <select v-model="editLeaveForm.leave_type" class="input text-xs">
                            <option value="ANNUAL">ลาพักร้อน</option>
                            <option value="SICK">ลาป่วย</option>
                            <option value="PERSONAL">ลากิจ</option>
                            <option value="UNPAID">ลาไม่รับค่าจ้าง</option>
                          </select>
                        </td>
                        <td>
                          <div class="space-y-2">
                            <UiDatePicker v-model="editLeaveForm.start_date" placeholder="วันที่เริ่ม…" />
                            <UiDatePicker v-model="editLeaveForm.end_date" placeholder="วันที่สิ้นสุด…" :disabled="editLeaveForm.is_half_day" />
                          </div>
                        </td>
                        <td>
                          <div class="space-y-2 text-xs">
                            <label class="inline-flex items-center gap-2"><input v-model="editLeaveForm.is_half_day" type="checkbox" /> ครึ่งวัน</label>
                            <select v-if="editLeaveForm.is_half_day" v-model="editLeaveForm.half_day_session" class="input text-xs">
                              <option value="AM">เช้า</option>
                              <option value="PM">บ่าย</option>
                            </select>
                          </div>
                        </td>
                        <td><input v-model="editLeaveForm.reason" class="input text-xs" /></td>
                      </template>
                      <template v-else>
                        <td>{{ leaveTypeLabel(r.leave_type) }} <span class="text-xs text-gray-500">({{ r.is_half_day ? requestedDurationLabel(r) : 'เต็มวัน' }})</span></td>
                        <td>{{ leaveDateRangeLabel(r) }}</td>
                        <td>{{ formatRequestedDays(r) }}</td>
                        <td>{{ r.reason }}</td>
                      </template>
                      <td :class="leaveStatusClass(r.status)">{{ leaveStatusLabel(r.status) }}</td>
                      <td>
                        <button class="btn" :disabled="actionLoadingId === r.id" @click="openLeaveHistory(r.id)">ดูประวัติ</button>
                      </td>
                      <td v-if="canManageLeaveRecords" class="whitespace-nowrap">
                        <div class="flex flex-wrap gap-2">
                          <button v-if="editingLeaveId !== r.id" class="btn" :disabled="actionLoadingId === r.id" @click="startEditLeave(r)">แก้ไข</button>
                          <button v-if="editingLeaveId === r.id" class="btn-primary" :disabled="actionLoadingId === r.id" @click="saveEditLeave(r.id)">บันทึก</button>
                          <button v-if="editingLeaveId === r.id" class="btn" :disabled="actionLoadingId === r.id" @click="cancelEditLeave">ยกเลิกแก้ไข</button>
                          <button class="btn-danger" :disabled="actionLoadingId === r.id" @click="cancelLeaveByAdmin(r)">ยกเลิกคำขอ</button>
                          <button class="btn-danger" :disabled="actionLoadingId === r.id" @click="deleteLeaveByAdmin(r)">ลบ</button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div v-if="historyLeaveId !== null" class="fixed inset-0 z-40 flex items-center justify-center bg-black/60 p-4" @click.self="closeLeaveHistory">
          <div class="w-full max-w-3xl rounded-2xl border border-gray-700 bg-gray-900 p-5">
            <div class="mb-4 flex items-center justify-between gap-3">
              <h3 class="text-sm font-semibold">ประวัติการแก้ไขคำขอลา #{{ historyLeaveId }}</h3>
              <button class="btn" @click="closeLeaveHistory">ปิด</button>
            </div>

            <div v-if="historyLoading" class="text-sm text-gray-400">กำลังโหลดประวัติ...</div>
            <div v-else-if="historyLogs.length === 0" class="text-sm text-gray-500 italic">ไม่พบประวัติการแก้ไข</div>
            <div v-else class="space-y-2 max-h-[60vh] overflow-y-auto">
              <div v-for="item in historyLogs" :key="item.id" class="row">
                <div class="space-y-0.5">
                  <p class="text-xs text-gray-200">{{ auditActionLabel(item.action) }}</p>
                  <p class="text-xs text-gray-500">สถานะ: {{ auditStatusLabel(item.old_status) }} → {{ auditStatusLabel(item.new_status) }}</p>
                  <p v-if="item.comment" class="text-xs text-gray-400">หมายเหตุ: {{ item.comment }}</p>
                </div>
                <span class="text-xs text-gray-500">{{ fmtDateTime(item.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type LeavePolicy, type HolidayCalendar, type LeaveRequest, type LeaveAuditLog, type LeaveTrendPoint, type LeaveBackfillBulkResultItem } from '~/core/modules/attendance/infrastructure/attendance-api'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useNotification } from '~/composables/useNotification'
import { useAuth } from '~/composables/useAuth'

definePageMeta({ layout: 'default', middleware: 'auth' })

const api = useAttendanceApi()
const teamsApi = useTeamsApi()
const { showSuccess, showError } = useNotification()
const { currentUser } = useAuth()

const loading = ref(false)
const savingPolicy = ref(false)
const savingHoliday = ref(false)
const reviewingId = ref<number | null>(null)

const tabs = [
  { id: 'policies', label: 'นโยบายวันลา' },
  { id: 'holidays', label: 'วันหยุดประจำปี' },
  { id: 'approvals', label: 'อนุมัติวันลา' },
  { id: 'backfill', label: 'กรอกย้อนหลัง' },
  { id: 'trend', label: 'ข้อมูลการลา' },
]
const activeTab = ref('policies')

const policies = ref<LeavePolicy[]>([])
const holidays = ref<HolidayCalendar[]>([])
const pending = ref<LeaveRequest[]>([])
const leaveRecords = ref<LeaveRequest[]>([])
const trend = ref<LeaveTrendPoint[]>([])
const trendMonthsBack = ref(6)
const auditLogs = ref<LeaveAuditLog[]>([])
const selectedAuditId = ref<number | null>(null)
const reviewComment = ref<Record<number, string>>({})
const backfillSubmitting = ref(false)
const backfillBulkSubmitting = ref(false)
const backfillBulkResults = ref<LeaveBackfillBulkResultItem[]>([])
const backfillBulkSummary = ref<{ total: number; succeeded: number; failed: number } | null>(null)
const editingLeaveId = ref<number | null>(null)
const actionLoadingId = ref<number | null>(null)
const historyLoading = ref(false)
const historyLogs = ref<LeaveAuditLog[]>([])
const historyLeaveId = ref<number | null>(null)
const editLeaveForm = reactive({
  start_date: '',
  end_date: '',
  leave_type: 'ANNUAL' as LeaveRequest['leave_type'],
  is_half_day: false,
  half_day_session: 'AM' as 'AM' | 'PM',
  reason: '',
})

const leaveFilter = reactive({
  from: '',
  to: '',
  leave_type: 'ALL',
  employee_key: 'ALL',
})
const allEmployees = ref<Array<{ key: string; label: string }>>([])
const leaveTypeOptions = [
  { value: 'ALL', label: 'ทั้งหมด' },
  { value: 'ANNUAL', label: 'ลาพักร้อน' },
  { value: 'SICK', label: 'ลาป่วย' },
  { value: 'PERSONAL', label: 'ลากิจ' },
  { value: 'ORDINATION', label: 'ลาบวช' },
  { value: 'MATERNITY', label: 'ลาคลอด' },
]

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
  is_half_day: false,
  half_day_session: 'AM' as 'AM' | 'PM',
  status: 'APPROVED' as 'PENDING' | 'APPROVED' | 'REJECTED',
  reason: '',
  comment: '',
})

const backfillDurationMode = computed({
  get: () => {
    if (!backfillForm.is_half_day) return 'FULL'
    return backfillForm.half_day_session === 'PM' ? 'HALF_PM' : 'HALF_AM'
  },
  set: (v: 'FULL' | 'HALF_AM' | 'HALF_PM') => {
    if (v === 'FULL') {
      backfillForm.is_half_day = false
      return
    }
    backfillForm.is_half_day = true
    backfillForm.half_day_session = v === 'HALF_PM' ? 'PM' : 'AM'
  }
})

watch(() => [backfillForm.is_half_day, backfillForm.start_date] as const, ([isHalf, start]) => {
  if (isHalf && start) {
    backfillForm.end_date = start
  }
})

function leaveTypeLabel(t: string) {
  if (t === 'ANNUAL') return 'ลาพักร้อน'
  if (t === 'SICK') return 'ลาป่วย'
  if (t === 'PERSONAL') return 'ลากิจ'
  if (t === 'UNPAID') return 'ลาไม่รับค่าจ้าง'
  return t
}

const fmt = (s: string) => new Date(s).toLocaleDateString('th-TH', { day: '2-digit', month: 'short', year: 'numeric' })
const fmtDateTime = (s: string) => new Date(s).toLocaleString('th-TH', { dateStyle: 'medium', timeStyle: 'short' })

function fallbackRequestedDays(r: Pick<LeaveRequest, 'start_date' | 'end_date' | 'is_half_day'>): number {
  if (r.is_half_day) return 0.5
  const start = new Date(r.start_date)
  const end = new Date(r.end_date)
  if (Number.isNaN(start.getTime()) || Number.isNaN(end.getTime())) return 0
  const startUtc = Date.UTC(start.getUTCFullYear(), start.getUTCMonth(), start.getUTCDate())
  const endUtc = Date.UTC(end.getUTCFullYear(), end.getUTCMonth(), end.getUTCDate())
  const diff = Math.floor((endUtc - startUtc) / (1000 * 60 * 60 * 24)) + 1
  return diff > 0 ? diff : 0
}

function formatRequestedDays(r: LeaveRequest): string {
  if (r.is_half_day) return '0.5'
  const raw = Number(r.days_requested)
  const days = raw > 0 ? raw : fallbackRequestedDays(r)
  return Number.isInteger(days) ? String(days) : days.toFixed(1)
}

function requestedDurationLabel(r: LeaveRequest): string {
  if (r.is_half_day) {
    const session = (r.half_day_session || '').toUpperCase()
    if (session === 'PM') return 'ครึ่งวันบ่าย'
    return 'ครึ่งวันเช้า'
  }
  return `${formatRequestedDays(r)} วัน`
}

function trendOwnerLabel(t: LeaveTrendPoint): string {
  if (t.user_name || t.user_email) {
    const name = t.user_name || (t.user_id ? `User #${t.user_id}` : 'Unknown user')
    if (t.user_email) return `${name} (${t.user_email})`
    return name
  }
  return t.team_name || 'Unassigned'
}

function trendRowKey(t: LeaveTrendPoint): string {
  if (t.user_id) return `${t.month}-u-${t.user_id}`
  return `${t.month}-t-${t.team_id || 0}-${t.team_name || 'unassigned'}`
}

function leaveDateRangeLabel(r: LeaveRequest): string {
  const start = fmtDate(r.start_date)
  const end = fmtDate(r.end_date)
  return `${start} ถึง ${end}`
}

function fmtDate(s: string): string {
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return s
  return d.toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })
}

function leaveStatusLabel(status: LeaveRequest['status']): string {
  if (status === 'PENDING') return 'กำลังพิจารณา'
  if (status === 'APPROVED') return 'อนุมัติ'
  if (status === 'REJECTED') return 'ไม่อนุมัติ'
  return status
}

function leaveStatusClass(status: LeaveRequest['status']): string {
  if (status === 'APPROVED') return 'text-emerald-300'
  if (status === 'REJECTED') return 'text-red-300'
  return 'text-amber-300'
}

function auditActionLabel(action: string): string {
  if (action === 'UPDATED') return 'แก้ไขคำขอ'
  if (action === 'REVIEWED') return 'พิจารณาคำขอ'
  if (action === 'CANCELLED') return 'ยกเลิกคำขอ'
  if (action === 'DELETED') return 'ลบคำขอ'
  if (action === 'CREATED') return 'สร้างคำขอ'
  return action
}

function auditStatusLabel(status: string | null | undefined): string {
  if (!status) return '-'
  if (status === 'PENDING' || status === 'APPROVED' || status === 'REJECTED') {
    return leaveStatusLabel(status)
  }
  return status
}

const leaveEmployeeOptions = computed<Array<{ key: string; label: string }>>(() => {
  const map = new Map<string, string>()

  for (const u of allEmployees.value) {
    map.set(u.key, u.label)
  }

  for (const r of leaveRecords.value) {
    const key = String(r.user_id || '')
    if (!key) continue
    if (!map.has(key)) {
      map.set(key, r.user_email || `User #${r.user_id}`)
    }
  }

  return Array.from(map.entries()).map(([key, label]) => ({ key, label }))
})

const filteredLeaveRecords = computed<LeaveRequest[]>(() => {
  const fromTs = leaveFilter.from ? new Date(leaveFilter.from).setHours(0, 0, 0, 0) : null
  const toTs = leaveFilter.to ? new Date(leaveFilter.to).setHours(23, 59, 59, 999) : null

  return leaveRecords.value.filter((r) => {
    if (leaveFilter.leave_type !== 'ALL' && r.leave_type !== leaveFilter.leave_type) return false
    if (leaveFilter.employee_key !== 'ALL' && String(r.user_id) !== leaveFilter.employee_key) return false

    const start = new Date(r.start_date).getTime()
    const end = new Date(r.end_date).getTime()
    if (fromTs != null && end < fromTs) return false
    if (toTs != null && start > toTs) return false

    return true
  })
})

const isAdminLeaveManager = computed(() => {
  const role = String(currentUser.value?.role || '').toUpperCase()
  return role === 'CEO' || role === 'SUPPORT'
})

const canManageLeaveRecords = computed(() => isAdminLeaveManager.value)

function toLocalDateString(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function trendRangeFromMonthsBack(monthsBack: number): { from: string; to: string } {
  const now = new Date()
  const to = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const from = new Date(now.getFullYear(), now.getMonth() - (monthsBack - 1), 1)
  return {
    from: toLocalDateString(from),
    to: toLocalDateString(to),
  }
}

async function refreshTrend() {
  const { from, to } = trendRangeFromMonthsBack(trendMonthsBack.value)
  const tr = await api.getLeaveTrend(from, to)
  trend.value = tr.items || []
}

async function refreshAll() {
  loading.value = true
  try {
    const now = new Date()
    const from = new Date(now.getFullYear(), 0, 1).toISOString().slice(0, 10)
    const to = new Date(now.getFullYear(), 11, 31).toISOString().slice(0, 10)
    const trendRange = trendRangeFromMonthsBack(trendMonthsBack.value)

    const [p, h, pend, allLeaves, tr, teams] = await Promise.all([
      api.getLeavePolicies(),
      api.getHolidays(from, to),
      api.getPendingLeaveRequests(),
      api.getAdminLeaveRequests().catch(() => ({ items: [] })),
      api.getLeaveTrend(trendRange.from, trendRange.to),
      teamsApi.getTeams().catch(() => []),
    ])
    policies.value = p.items || []
    holidays.value = h.items || []
    pending.value = pend.items || []
    leaveRecords.value = (allLeaves.items && allLeaves.items.length > 0) ? allLeaves.items : (pend.items || [])
    trend.value = tr.items || []

    const employeeMap = new Map<string, string>()
    for (const team of teams) {
      for (const user of (team.users || [])) {
        const key = String(user.id || '')
        if (!key) continue
        if (!employeeMap.has(key)) {
          employeeMap.set(key, user.email || `User #${user.id}`)
        }
      }
    }
    allEmployees.value = Array.from(employeeMap.entries()).map(([key, label]) => ({ key, label }))
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

async function openLeaveHistory(leaveId: number) {
  historyLeaveId.value = leaveId
  historyLoading.value = true
  try {
    const res = await api.getLeaveAudit(leaveId)
    historyLogs.value = res.items || []
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'โหลดประวัติการแก้ไขไม่สำเร็จ')
  } finally {
    historyLoading.value = false
  }
}

function closeLeaveHistory() {
  historyLeaveId.value = null
  historyLogs.value = []
}

async function submitBackfillSingle() {
  backfillSubmitting.value = true
  try {
    const item = {
      employee_email: backfillForm.employee_email.trim(),
      start_date: backfillForm.start_date,
      end_date: backfillForm.is_half_day ? backfillForm.start_date : backfillForm.end_date,
      leave_type: backfillForm.leave_type,
      is_half_day: backfillForm.is_half_day,
      half_day_session: backfillForm.is_half_day ? backfillForm.half_day_session : undefined,
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
  const required = ['employee_email', 'start_date', 'end_date', 'leave_type', 'is_half_day', 'half_day_session', 'status', 'reason', 'comment']
  const missing = required.filter(k => !header.includes(k))
  if (missing.length) {
    throw new Error(`CSV header missing: ${missing.join(', ')}`)
  }
  const idx = (k: string) => header.indexOf(k)
  const out = [] as Array<{ employee_email: string; start_date: string; end_date: string; leave_type: 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID'; is_half_day: boolean; half_day_session?: 'AM' | 'PM'; status: 'PENDING' | 'APPROVED' | 'REJECTED'; reason: string; comment: string }>
  for (let i = 1; i < rows.length; i++) {
    const cols = rows[i].split(',').map(c => c.trim())
    const isHalfRaw = (cols[idx('is_half_day')] || '').toLowerCase()
    const isHalfDay = ['true', '1', 'yes', 'y'].includes(isHalfRaw)
    const halfDaySessionRaw = (cols[idx('half_day_session')] || '').toUpperCase()
    const halfDaySession = halfDaySessionRaw === 'PM' ? 'PM' : 'AM'
    const startDate = cols[idx('start_date')] || ''
    out.push({
      employee_email: cols[idx('employee_email')] || '',
      start_date: startDate,
      end_date: isHalfDay ? startDate : (cols[idx('end_date')] || ''),
      leave_type: ((cols[idx('leave_type')] || 'ANNUAL').toUpperCase() as any),
      is_half_day: isHalfDay,
      half_day_session: isHalfDay ? halfDaySession : undefined,
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
    'employee_email,start_date,end_date,leave_type,is_half_day,half_day_session,status,reason,comment',
    'alice@company.com,2026-01-15,2026-01-16,ANNUAL,false,,APPROVED,ลาพักร้อนต้นปี,อนุมัติย้อนหลังโดย HR',
    'bob@company.com,2026-02-03,2026-02-03,SICK,true,AM,APPROVED,ป่วยไข้หวัดครึ่งวันเช้า,มีใบรับรองแพทย์',
    'carol@company.com,2026-03-10,2026-03-10,PERSONAL,true,PM,PENDING,ธุระส่วนตัวครึ่งวันบ่าย,รอหัวหน้าอนุมัติ',
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

function startEditLeave(r: LeaveRequest) {
  editingLeaveId.value = r.id
  editLeaveForm.start_date = String(r.start_date || '').slice(0, 10)
  editLeaveForm.end_date = String(r.end_date || '').slice(0, 10)
  editLeaveForm.leave_type = r.leave_type
  editLeaveForm.is_half_day = !!r.is_half_day
  editLeaveForm.half_day_session = (r.half_day_session === 'PM' ? 'PM' : 'AM')
  editLeaveForm.reason = r.reason || ''
}

function cancelEditLeave() {
  editingLeaveId.value = null
}

async function saveEditLeave(leaveID: number) {
  actionLoadingId.value = leaveID
  try {
    const payload = {
      start_date: editLeaveForm.start_date,
      end_date: editLeaveForm.is_half_day ? editLeaveForm.start_date : editLeaveForm.end_date,
      leave_type: editLeaveForm.leave_type,
      is_half_day: editLeaveForm.is_half_day,
      half_day_session: editLeaveForm.is_half_day ? editLeaveForm.half_day_session : undefined,
      reason: editLeaveForm.reason,
    }
    await api.updateAdminLeaveRequest(leaveID, payload)
    showSuccess('แก้ไขคำขอลาเรียบร้อย')
    editingLeaveId.value = null
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'แก้ไขคำขอลาไม่สำเร็จ')
  } finally {
    actionLoadingId.value = null
  }
}

async function cancelLeaveByAdmin(r: LeaveRequest) {
  const ok = window.confirm(`ยืนยันยกเลิกคำขอลา #${r.id} ?`)
  if (!ok) return
  actionLoadingId.value = r.id
  try {
    await api.cancelAdminLeaveRequest(r.id, { comment: 'Cancelled by admin' })
    showSuccess('ยกเลิกคำขอเรียบร้อย')
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'ยกเลิกคำขอไม่สำเร็จ')
  } finally {
    actionLoadingId.value = null
  }
}

async function deleteLeaveByAdmin(r: LeaveRequest) {
  const ok = window.confirm(`ยืนยันลบคำขอลา #${r.id} ? การลบไม่สามารถย้อนกลับได้`)
  if (!ok) return
  actionLoadingId.value = r.id
  try {
    await api.deleteAdminLeaveRequest(r.id)
    showSuccess('ลบคำขอเรียบร้อย')
    if (editingLeaveId.value === r.id) editingLeaveId.value = null
    await refreshAll()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'ลบคำขอไม่สำเร็จ')
  } finally {
    actionLoadingId.value = null
  }
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

.dt-row { @apply rounded-xl border border-gray-700/80 bg-gray-900/30 overflow-hidden; }
.dt-col { @apply w-full overflow-x-auto; }
.leave-table { @apply min-w-[980px] border-separate border-spacing-0; }
.leave-table thead tr { @apply bg-gray-800/70; }
.leave-table th { @apply px-4 py-3 font-semibold uppercase tracking-wide whitespace-nowrap border-b border-gray-700; }
.leave-table td { @apply px-4 py-3 align-top border-b border-gray-800/70; }
.leave-table tbody tr:nth-child(even) { @apply bg-gray-800/20; }
.leave-table tbody tr:hover { @apply bg-cyan-900/10; }
.leave-table tbody tr:last-child td { @apply border-b-0; }
.leave-table td:nth-child(1) { @apply w-14 text-gray-400; }
.leave-table td:nth-child(5) { @apply w-24 font-semibold text-gray-200; }
.leave-table td:nth-child(7) { @apply w-36 font-semibold; }
</style>
