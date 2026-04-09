<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <header class="sticky top-0 z-10 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-5xl px-4 py-6 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
              บัญชี / Accounting
            </h1>
            <p class="mt-1 text-sm text-gray-400">
              กรอกข้อมูลดิบรายเดือน — ระบบจะคำนวณ Runway, Burn rate, MRR ให้แสดงใน Dashboard
            </p>
          </div>
          <button
            type="button"
            @click="loadEntries"
            class="inline-flex items-center gap-2 rounded-xl border border-gray-600 bg-gray-800 px-4 py-2.5 text-sm font-medium text-gray-200 transition-colors hover:border-gray-500 hover:bg-gray-700 hover:text-gray-900 dark:text-white"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            โหลดใหม่
          </button>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-5xl px-4 py-8 sm:px-6 lg:px-8">
      <!-- Form: กรอกข้อมูลรายเดือน -->
      <section class="mb-10 rounded-xl border border-gray-700/80 bg-gray-800/60 p-6 shadow-lg">
        <h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-gray-400">
          กรอกข้อมูลเดือน
        </h2>
        <form @submit.prevent="submit" class="grid gap-4 sm:grid-cols-2 lg:grid-cols-6">
          <div>
            <label for="year" class="mb-1 block text-xs font-medium text-gray-500">ปี (Year)</label>
            <input
              id="year"
              v-model.number="form.year"
              type="number"
              min="2020"
              max="2100"
              required
              class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <div>
            <label for="month" class="mb-1 block text-xs font-medium text-gray-500">เดือน (Month)</label>
            <select
              id="month"
              v-model.number="form.month"
              required
              class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
            >
              <option v-for="m in 12" :key="m" :value="m">{{ monthName(m) }}</option>
            </select>
          </div>
          <div>
            <label for="revenue" class="mb-1 block text-xs font-medium text-gray-500">รายได้ (Revenue)</label>
            <input
              id="revenue"
              v-model.number="form.revenue"
              type="number"
              min="0"
              step="0.01"
              class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <div>
            <label for="expenses" class="mb-1 block text-xs font-medium text-gray-500">ค่าใช้จ่าย (Expenses)</label>
            <input
              id="expenses"
              v-model.number="form.expenses"
              type="number"
              min="0"
              step="0.01"
              class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <div>
            <label for="cash_balance" class="mb-1 block text-xs font-medium text-gray-500">เงินคงเหลือปลายเดือน (Cash)</label>
            <input
              id="cash_balance"
              v-model.number="form.cash_balance"
              type="number"
              min="0"
              step="0.01"
              class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
            />
          </div>
          <div class="flex items-end">
            <button
              type="submit"
              :disabled="saving"
              class="w-full rounded-xl bg-amber-100 dark:bg-amber-600 px-4 py-2.5 text-sm font-semibold text-gray-900 dark:text-white transition-colors hover:bg-amber-100 dark:bg-amber-500 disabled:opacity-50"
            >
              {{ saving ? 'กำลังบันทึก...' : 'บันทึก' }}
            </button>
          </div>
        </form>
        <div class="mt-4">
          <label for="note" class="mb-1 block text-xs font-medium text-gray-500">หมายเหตุ (Note)</label>
          <input
            id="note"
            v-model="form.note"
            type="text"
            placeholder="เช่น ปิดบัญชี ม.ค. 2567"
            class="w-full rounded-lg border border-gray-600 bg-gray-900 px-3 py-2 text-white placeholder-gray-500 focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500"
          />
        </div>
        <p v-if="saveMessage" class="mt-3 text-sm" :class="saveMessageSuccess ? 'text-emerald-400' : 'text-red-400'">
          {{ saveMessage }}
        </p>
      </section>

      <!-- สรุปที่ระบบคำนวณ (จากข้อมูลล่าสุด) -->
      <section v-if="summary" class="mb-8 rounded-xl border border-gray-700/80 bg-gray-800/60 p-5">
        <h2 class="mb-3 text-xs font-semibold uppercase tracking-wider text-gray-500">
          สรุปที่แสดงใน Dashboard (คำนวณอัตโนมัติ)
        </h2>
        <div class="grid gap-3 sm:grid-cols-4">
          <div>
            <div class="text-xs text-gray-500">Cash Runway</div>
            <div class="text-xl font-bold text-white">{{ formatRunway(summary.runway_months) }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500">Burn rate (เฉลี่ย/เดือน)</div>
            <div class="text-xl font-bold text-white">{{ formatMoney(summary.burn_rate) }} {{ summary.currency }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500">MRR (รายได้เดือนล่าสุด)</div>
            <div class="text-xl font-bold text-white">{{ formatMoney(summary.last_month_mrr) }} {{ summary.currency }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500">Net new ARR</div>
            <div class="text-xl font-bold text-white">{{ formatMoney(summary.net_new_arr) }} {{ summary.currency }}</div>
          </div>
        </div>
      </section>

      <!-- ตารางรายเดือน -->
      <section>
        <h2 class="mb-4 text-xs font-semibold uppercase tracking-wider text-gray-500">
          ข้อมูลที่กรอกแล้ว (เรียงเดือนล่าสุดก่อน)
        </h2>
        <div v-if="loading" class="py-8 text-center text-gray-500">กำลังโหลด...</div>
        <div v-else-if="entries.length === 0" class="rounded-xl border border-gray-700/80 bg-gray-800/40 py-12 text-center text-gray-500">
          ยังไม่มีข้อมูล — กรอกเดือนแรกด้านบนแล้วกดบันทึก
        </div>
        <div v-else class="overflow-hidden rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/80 text-left text-xs font-medium uppercase tracking-wider text-gray-400">
                  <th class="px-5 py-4">ปี / เดือน</th>
                  <th class="px-5 py-4 text-right">รายได้</th>
                  <th class="px-5 py-4 text-right">ค่าใช้จ่าย</th>
                  <th class="px-5 py-4 text-right">เงินคงเหลือ</th>
                  <th class="px-5 py-4">หมายเหตุ</th>
                  <th class="px-5 py-4 text-center">แก้ไข</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/80">
                <tr
                  v-for="e in entries"
                  :key="`${e.year}-${e.month}`"
                  class="transition-colors hover:bg-gray-700/30"
                >
                  <td class="px-5 py-4 font-medium text-white">{{ e.year }} / {{ monthName(e.month) }}</td>
                  <td class="px-5 py-4 text-right text-gray-300">{{ formatMoney(e.revenue) }}</td>
                  <td class="px-5 py-4 text-right text-gray-300">{{ formatMoney(e.expenses) }}</td>
                  <td class="px-5 py-4 text-right text-emerald-400">{{ formatMoney(e.cash_balance) }}</td>
                  <td class="px-5 py-4 text-gray-500">{{ e.note || '—' }}</td>
                  <td class="px-5 py-4 text-center">
                    <button
                      type="button"
                      @click="editEntry(e)"
                      class="rounded-lg border border-gray-600 px-2 py-1 text-xs text-gray-400 hover:bg-gray-700 hover:text-gray-900 dark:text-white"
                    >
                      เลือก
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useFinanceApi } from '~/core/modules/finance/finance-api'
import type { MonthlyEntry, FinanceSummary } from '~/core/modules/finance/finance-api'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const api = useFinanceApi()
const { currentUser } = useAuth()

const entries = ref<MonthlyEntry[]>([])
const summary = ref<FinanceSummary | null>(null)
const loading = ref(true)
const saving = ref(false)
const saveMessage = ref('')
const saveMessageSuccess = ref(false)

const form = reactive({
  year: new Date().getFullYear(),
  month: new Date().getMonth() + 1,
  revenue: 0,
  expenses: 0,
  cash_balance: 0,
  note: ''
})

const monthNames = [
  'ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.',
  'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.'
]
function monthName(m: number) {
  return monthNames[m - 1] ?? String(m)
}

function formatMoney(value: number) {
  return new Intl.NumberFormat('th-TH', { minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(value)
}

function formatRunway(months: number) {
  if (months <= 0) return '—'
  return `${months.toFixed(1)} เดือน`
}

async function loadEntries() {
  loading.value = true
  try {
    entries.value = await api.getEntries(24)
    summary.value = await api.getSummary()
  } catch (e: any) {
    console.error('Failed to load finance data:', e)
  } finally {
    loading.value = false
  }
}

function editEntry(e: MonthlyEntry) {
  form.year = e.year
  form.month = e.month
  form.revenue = e.revenue
  form.expenses = e.expenses
  form.cash_balance = e.cash_balance
  form.note = e.note ?? ''
  saveMessage.value = ''
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function submit() {
  saving.value = true
  saveMessage.value = ''
  try {
    await api.createOrUpdateEntry({
      year: form.year,
      month: form.month,
      revenue: Number(form.revenue) || 0,
      expenses: Number(form.expenses) || 0,
      cash_balance: Number(form.cash_balance) || 0,
      note: form.note?.trim() ?? ''
    })
    saveMessage.value = `บันทึก ${form.year}/${monthName(form.month)} แล้ว`
    saveMessageSuccess.value = true
    await loadEntries()
  } catch (e: any) {
    saveMessage.value = e?.data?.message || e?.message || 'บันทึกไม่สำเร็จ'
    saveMessageSuccess.value = false
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  if (currentUser.value?.role !== 'CEO') {
    navigateTo('/dashboard')
    return
  }
  loadEntries()
})
</script>
