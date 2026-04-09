<template>
  <div class="min-h-screen bg-gray-50 text-gray-900 dark:bg-gray-900 dark:text-white">
    <div class="border-b border-gray-200 dark:border-gray-800 bg-white/95 dark:bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="max-w-screen-xl mx-auto flex items-center justify-between gap-3">
        <div>
          <h1 class="text-base font-bold text-gray-900 dark:text-white">Leave Self-Service</h1>
          <p class="text-xs text-gray-600 dark:text-gray-500">Request leave, check balance, holidays, and approval updates</p>
        </div>
        <button
          class="inline-flex items-center gap-2 px-3 py-2 rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800/60 text-xs font-medium text-gray-700 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700"
          :disabled="loading"
          @click="refresh"
        >
          Refresh
        </button>
      </div>
    </div>

    <main class="max-w-screen-xl mx-auto px-6 py-8 space-y-6">
      <div class="grid gap-6 lg:grid-cols-3">
        <section class="lg:col-span-2 rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800/50 p-5">
          <h2 class="text-sm font-semibold mb-4">ส่งคำขออนุมัติลา</h2>
          <form class="space-y-4" @submit.prevent="submitLeave">
            <div class="grid sm:grid-cols-4 gap-3">
              <div>
                <label class="label">ประเภทการลา</label>
                <select v-model="form.leave_type" class="input">
                  <option value="ANNUAL">ลาพักร้อน</option>
                  <option value="SICK">ลาป่วย</option>
                  <option value="PERSONAL">ลากิจ</option>
                </select>
              </div>
              <div>
                <label class="label">วันที่เริ่มลา</label>
                <UiDatePicker v-model="form.start_date" placeholder="เลือกวันที่เริ่มลา…" />
              </div>
              <div>
                <label class="label">วันที่สิ้นสุดลา</label>
                <UiDatePicker v-model="form.end_date" placeholder="เลือกวันที่สิ้นสุดลา…" :disabled="form.is_half_day" />
              </div>
              <div>
                <label class="label">รูปแบบวันลา</label>
                <select v-model="leaveDurationMode" class="input">
                  <option value="FULL">เต็มวัน</option>
                  <option value="HALF_AM">ครึ่งวันเช้า</option>
                  <option value="HALF_PM">ครึ่งวันบ่าย</option>
                </select>
              </div>
            </div>
            <div>
              <label class="label">เหตุผลการลา</label>
              <textarea v-model="form.reason" rows="3" class="input" required placeholder="ระบุเหตุผลเพื่อประกอบการอนุมัติ" />
            </div>
            <div class="flex justify-end">
              <button
                type="submit"
                :disabled="submitting"
                class="rounded-lg border border-blue-500/40 bg-blue-900/20 px-4 py-2 text-sm font-semibold text-blue-300 hover:bg-blue-900/40 disabled:opacity-50"
              >
                {{ submitting ? 'กำลังส่ง...' : 'ส่งคำขอ' }}
              </button>
            </div>
          </form>
        </section>

        <section class="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800/50 p-5">
          <h2 class="text-sm font-semibold mb-4">คงเหลือวันลา ปี {{ year }}</h2>
          <div class="space-y-3">
            <div v-for="item in balance" :key="item.leave_type" class="rounded-lg border border-gray-200 dark:border-gray-700/70 bg-gray-50 dark:bg-gray-900/40 p-3">
              <div class="flex justify-between text-xs">
                <span class="text-gray-700 dark:text-gray-400">{{ item.leave_type }}</span>
                <span class="font-bold text-emerald-300">{{ item.remaining_days }} days</span>
              </div>
              <p class="text-[11px] text-gray-600 dark:text-gray-500 mt-1">Quota {{ item.annual_quota_days }} + Carry {{ item.carry_forward_days }} · Used {{ item.approved_days_taken }}</p>
            </div>
          </div>
        </section>
      </div>

      <section class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
        <h2 class="text-sm font-semibold mb-4">My Leave Requests</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-xs text-gray-600 dark:text-gray-500 border-b border-gray-200 dark:border-gray-700">
                <th class="text-left py-2">Date Range</th>
                <th class="text-left py-2">Type</th>
                <th class="text-left py-2">Days</th>
                <th class="text-left py-2">Status</th>
                <th class="text-left py-2">Manager Comment</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in requests" :key="item.id" class="border-b border-gray-200 dark:border-gray-800/70">
                <td class="py-2 text-gray-700 dark:text-gray-300">{{ fmt(item.start_date) }} → {{ fmt(item.end_date) }}</td>
                <td class="py-2 text-gray-700 dark:text-gray-300">{{ item.leave_type }}<span v-if="item.is_half_day" class="ml-1 text-xs text-cyan-700 dark:text-cyan-300">({{ item.half_day_session === 'AM' ? 'ครึ่งเช้า' : 'ครึ่งบ่าย' }})</span></td>
                <td class="py-2 text-gray-700 dark:text-gray-300">{{ formatDays(item.days_requested) }}</td>
                <td class="py-2">
                  <span class="text-xs px-2 py-0.5 rounded-full" :class="statusCls(item.status)">{{ item.status }}</span>
                </td>
                <td class="py-2 text-gray-600 dark:text-gray-400">{{ item.manager_comment || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <div class="grid gap-6 lg:grid-cols-2">
        <section class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h2 class="text-sm font-semibold mb-4">Holiday Calendar</h2>
          <div class="space-y-2">
            <div v-for="h in holidays" :key="h.id" class="flex items-center justify-between text-sm rounded-lg border border-gray-200 dark:border-gray-700/70 px-3 py-2 bg-gray-50 dark:bg-gray-900/30">
              <span class="text-gray-700 dark:text-gray-300">{{ fmt(h.date) }}</span>
              <span class="text-gray-600 dark:text-gray-400">{{ h.name }}</span>
            </div>
          </div>
        </section>

        <section class="rounded-2xl border border-gray-700 bg-gray-800/50 p-5">
          <h2 class="text-sm font-semibold mb-4">Notifications</h2>
          <div class="space-y-2">
            <div v-for="n in notifications" :key="n.id" class="rounded-lg border border-gray-200 dark:border-gray-700/70 p-3 bg-gray-50 dark:bg-gray-900/30">
              <div class="flex items-center justify-between gap-2">
                <p class="text-sm font-medium text-white">{{ n.title }}</p>
                <button v-if="!n.is_read" class="text-xs text-blue-300 hover:text-blue-200" @click="markRead(n.id)">Mark read</button>
              </div>
              <p class="text-xs text-gray-400 mt-1">{{ n.message }}</p>
              <p class="text-[11px] text-gray-500 mt-1">{{ n.channel }} · {{ fmtDateTime(n.created_at) }}</p>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAttendanceApi, type LeaveRequest, type LeaveBalanceSummary, type HolidayCalendar, type LeaveNotification } from '~/core/modules/attendance/infrastructure/attendance-api'
import { useNotification } from '~/composables/useNotification'

definePageMeta({ layout: 'default', middleware: 'auth' })

const api = useAttendanceApi()
const { showSuccess, showError } = useNotification()

const loading = ref(false)
const submitting = ref(false)
const year = new Date().getFullYear()

const form = reactive({
  leave_type: 'ANNUAL' as 'ANNUAL' | 'SICK' | 'PERSONAL' | 'UNPAID',
  start_date: '',
  end_date: '',
  is_half_day: false,
  half_day_session: 'AM' as 'AM' | 'PM',
  reason: '',
})

const leaveDurationMode = computed({
  get: () => {
    if (!form.is_half_day) return 'FULL'
    return form.half_day_session === 'PM' ? 'HALF_PM' : 'HALF_AM'
  },
  set: (v: 'FULL' | 'HALF_AM' | 'HALF_PM') => {
    if (v === 'FULL') {
      form.is_half_day = false
      return
    }
    form.is_half_day = true
    form.half_day_session = v === 'HALF_PM' ? 'PM' : 'AM'
  }
})

const balance = ref<LeaveBalanceSummary[]>([])
const requests = ref<LeaveRequest[]>([])
const holidays = ref<HolidayCalendar[]>([])
const notifications = ref<LeaveNotification[]>([])

watch(() => [form.is_half_day, form.start_date] as const, ([isHalf, start]) => {
  if (isHalf && start) {
    form.end_date = start
  }
})

const fmt = (s: string) => new Date(s).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
const fmtDateTime = (s: string) => new Date(s).toLocaleString('en-GB', { dateStyle: 'medium', timeStyle: 'short' })
const formatDays = (days: number) => Number.isInteger(days) ? `${days}` : `${days.toFixed(1)}`

function statusCls(status: string) {
  if (status === 'APPROVED') return 'bg-emerald-500/15 text-emerald-300 border border-emerald-500/30'
  if (status === 'REJECTED') return 'bg-red-500/15 text-red-300 border border-red-500/30'
  return 'bg-amber-500/15 text-amber-300 border border-amber-500/30'
}

async function refresh() {
  loading.value = true
  try {
    const [bal, reqs, hol, notif] = await Promise.all([
      api.getMyLeaveBalance(year),
      api.getMyLeaveRequests(),
      api.getHolidays(),
      api.getLeaveNotifications(false),
    ])
    balance.value = bal.items || []
    requests.value = reqs.items || []
    holidays.value = hol.items || []
    notifications.value = notif.items || []
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'Failed to load leave self-service data')
  } finally {
    loading.value = false
  }
}

async function submitLeave() {
  if (!form.start_date || !form.end_date || !form.reason.trim()) return
  submitting.value = true
  try {
    const payload = {
      leave_type: form.leave_type,
      start_date: form.start_date,
      end_date: form.is_half_day ? form.start_date : form.end_date,
      is_half_day: form.is_half_day,
      half_day_session: form.is_half_day ? form.half_day_session : undefined,
      reason: form.reason.trim(),
    }
    await api.createLeaveRequest(payload)
    showSuccess('Leave request submitted successfully')
    form.reason = ''
    await refresh()
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'Failed to submit leave request')
  } finally {
    submitting.value = false
  }
}

async function markRead(id: number) {
  try {
    await api.markLeaveNotificationRead(id)
    notifications.value = notifications.value.map(n => n.id === id ? { ...n, is_read: true } : n)
  } catch (e: any) {
    showError(e?.data?.error || e?.message || 'Failed to mark notification as read')
  }
}

onMounted(refresh)
</script>

<style scoped>
.label { @apply block text-xs font-semibold uppercase tracking-widest text-gray-600 dark:text-gray-500 mb-1.5; }
.input { @apply w-full rounded-lg border border-gray-300 dark:border-gray-700 bg-white dark:bg-gray-800/80 px-3 py-2 text-sm text-gray-900 dark:text-white focus:border-blue-500 focus:outline-none; }
</style>
