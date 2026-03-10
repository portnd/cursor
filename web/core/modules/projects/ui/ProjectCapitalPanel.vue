<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-white">Project Capital</h2>
        <p class="text-sm text-gray-400 mt-1">Internal VC — track capital injection, burn rate, and milestone payout for this project</p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="openInjectModal"
          class="flex items-center gap-2 px-4 py-2 bg-emerald-600/20 hover:bg-emerald-600/30 border border-emerald-600/40 text-emerald-400 rounded-lg text-sm font-medium transition-all"
        >
          + Inject Capital
        </button>
        <button
          @click="openEditModal"
          class="flex items-center gap-2 px-3 py-2 bg-gray-700/50 hover:bg-gray-700 border border-gray-600 text-gray-300 rounded-lg text-sm font-medium transition-all"
        >
          Edit
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="text-center py-16">
      <div class="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin mx-auto mb-3"></div>
      <p class="text-sm text-gray-400">Loading capital data...</p>
    </div>

    <template v-else-if="capital">
      <!-- KPIs row -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-gray-800/60 border border-gray-700/60 rounded-xl p-4">
          <p class="text-xs text-gray-400 mb-1">Capital Balance</p>
          <p class="text-2xl font-bold text-white">{{ formatMoney(capital.capital_balance) }}</p>
        </div>
        <div class="bg-gray-800/60 border border-gray-700/60 rounded-xl p-4">
          <p class="text-xs text-gray-400 mb-1">Team Monthly Burn</p>
          <p class="text-2xl font-bold text-red-400">{{ formatMoney(capital.team_monthly_cost) }}</p>
          <p class="text-xs text-gray-500 mt-1">Full team cost / month</p>
        </div>
        <div class="bg-gray-800/60 border border-gray-700/60 rounded-xl p-4">
          <p class="text-xs text-gray-400 mb-1">Runway</p>
          <p class="text-2xl font-bold" :class="runwayColor">{{ capital.runway_months.toFixed(1) }} mo</p>
          <p class="text-xs text-gray-500 mt-1">Capital / Monthly burn</p>
        </div>
        <div
          class="bg-gray-800/60 border rounded-xl p-4 transition-colors"
          :class="capital.capital_balance >= capital.team_monthly_cost ? 'border-emerald-500/40' : 'border-red-500/40'"
        >
          <p class="text-xs text-gray-400 mb-1">Gross Margin</p>
          <p
            class="text-2xl font-bold"
            :class="capital.capital_balance >= capital.team_monthly_cost ? 'text-emerald-400' : 'text-red-400'"
          >
            {{ capital.capital_balance >= capital.team_monthly_cost ? '+' : '' }}{{ formatMoney(capital.capital_balance - capital.team_monthly_cost) }}
          </p>
          <p class="text-xs text-gray-500 mt-1">Capital − monthly burn</p>
        </div>
      </div>

      <!-- Runway progress bar -->
      <div class="bg-gray-800/60 border border-gray-700/60 rounded-xl p-4">
        <div class="flex items-center justify-between text-sm mb-2">
          <span class="text-gray-400">Runway</span>
          <span :class="runwayColor" class="font-medium">{{ capital.runway_months.toFixed(2) }} months</span>
        </div>
        <div class="w-full bg-gray-700 rounded-full h-3 overflow-hidden">
          <div
            class="h-3 rounded-full transition-all duration-500"
            :class="runwayBarColor"
            :style="{ width: Math.min(runwayPct, 100) + '%' }"
          ></div>
        </div>
        <p class="text-xs text-gray-500 mt-2">
          Target runway: 3 months &nbsp;·&nbsp;
          {{ capital.runway_months >= 3 ? 'Fully funded ✓' : `${(3 - capital.runway_months).toFixed(1)} months short` }}
        </p>
      </div>

      <!-- Transaction history -->
      <div class="bg-gray-800/60 border border-gray-700/60 rounded-xl overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-700/60">
          <h3 class="text-sm font-semibold text-white">Transaction History</h3>
        </div>
        <div v-if="capital.transactions && capital.transactions.length > 0" class="divide-y divide-gray-700/40">
          <div
            v-for="tx in capital.transactions"
            :key="tx.id"
            class="flex items-center justify-between px-5 py-3 hover:bg-gray-700/20 transition-colors group"
          >
            <div class="flex items-center gap-3">
              <span class="text-lg">{{ txIcon(tx.type) }}</span>
              <div>
                <p class="text-sm text-white font-medium">{{ tx.reference || tx.type }}</p>
                <p class="text-xs text-gray-500">{{ formatDate(tx.created_at) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <div class="text-right">
                <p class="text-sm font-bold" :class="tx.amount >= 0 ? 'text-emerald-400' : 'text-red-400'">
                  {{ tx.amount >= 0 ? '+' : '' }}{{ formatMoney(tx.amount) }}
                </p>
                <span class="text-xs px-2 py-0.5 rounded-full" :class="txBadgeClass(tx.type)">{{ tx.type }}</span>
              </div>
              <!-- Delete: show confirm inline on first click -->
              <div v-if="deletingTxId === tx.id" class="flex items-center gap-1.5">
                <span class="text-xs text-red-400">ลบ?</span>
                <button @click="confirmDeleteTx(tx.id)" class="px-2 py-0.5 text-xs bg-red-600 hover:bg-red-500 text-white rounded font-medium transition-all">ใช่</button>
                <button @click="deletingTxId = null" class="px-2 py-0.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 rounded transition-all">ยกเลิก</button>
              </div>
              <button
                v-else
                @click="deletingTxId = tx.id"
                class="opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-red-400 transition-all"
                title="ลบ transaction"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>
          </div>
        </div>
        <div v-else class="text-center py-8 text-gray-500 text-sm">
          No transactions yet. Inject capital to get started.
        </div>
      </div>
    </template>

    <!-- Inject Capital Modal -->
    <div v-if="showInjectModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl w-full max-w-md shadow-2xl">
        <div class="flex items-center justify-between p-6 border-b border-gray-700">
          <h3 class="text-lg font-bold text-white">Inject Capital</h3>
          <button @click="showInjectModal = false" class="text-gray-400 hover:text-white transition-colors">✕</button>
        </div>
        <div class="p-6 space-y-4">
          <!-- Delivery date -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">วันที่ส่งงาน <span class="text-red-400">*</span></label>
            <input
              v-model="injectModal.deliveryDate"
              type="date"
              :min="injectModalDeliveryDateMin"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
            <p class="text-xs text-gray-500 mt-1">ระบบจะคำนวณยอดเงินให้หมดพอดีวันส่งงาน (จาก Burn Rate ปัจจุบัน)</p>
            <p v-if="injectModal.deliveryDate && injectModalAmountFromDate !== null" class="text-xs text-emerald-400 mt-1">
              คำนวณจาก Burn {{ formatMoney(capital?.team_monthly_cost ?? 0) }}/mo × {{ injectModalMonthsToDelivery.toFixed(1) }} mo
            </p>
          </div>
          <!-- Amount -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">Amount (฿) <span class="text-red-400">*</span></label>
            <input
              v-model.number="injectModal.amount"
              type="number"
              min="1"
              placeholder="คำนวณอัตโนมัติเมื่อเลือกวันที่ส่งงาน"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
          </div>
          <!-- Bonus % -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">Bonus Target (%)</label>
            <input
              v-model.number="injectModal.bonusPct"
              type="number"
              min="0"
              max="100"
              placeholder="e.g. 20"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
          </div>
          <!-- Note -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">Note / Reference</label>
            <input
              v-model="injectModal.note"
              type="text"
              placeholder="e.g. งวดที่ 1 MIMS HD-MAP"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
          </div>
        </div>
        <div class="flex gap-3 p-6 pt-0">
          <button
            @click="showInjectModal = false"
            class="flex-1 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl text-sm font-medium transition-all"
          >
            Cancel
          </button>
          <button
            @click="confirmInject"
            :disabled="!injectModal.amount || injectModal.amount <= 0 || saving"
            class="flex-1 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white rounded-xl text-sm font-bold transition-all"
          >
            {{ saving ? 'Injecting...' : 'Inject Capital' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Capital Modal -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl w-full max-w-md shadow-2xl">
        <div class="flex items-center justify-between p-6 border-b border-gray-700">
          <h3 class="text-lg font-bold text-white">Edit Capital Balance</h3>
          <button @click="showEditModal = false" class="text-gray-400 hover:text-white transition-colors">✕</button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">New Balance (฿)</label>
            <input
              v-model.number="editModal.newBalance"
              type="number"
              min="0"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
            <p v-if="capital && editModal.newBalance !== null" class="text-xs mt-1.5" :class="editDelta >= 0 ? 'text-emerald-400' : 'text-red-400'">
              Delta: {{ editDelta >= 0 ? '+' : '' }}{{ formatMoney(editDelta) }}
            </p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">Bonus Target (%) — optional</label>
            <input
              v-model.number="editModal.bonusPct"
              type="number"
              min="0"
              max="100"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-1.5">Note</label>
            <input
              v-model="editModal.note"
              type="text"
              class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:border-indigo-500"
            />
          </div>
        </div>
        <div class="flex gap-3 p-6 pt-0">
          <button @click="showEditModal = false" class="flex-1 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl text-sm font-medium transition-all">Cancel</button>
          <button
            @click="confirmEdit"
            :disabled="editModal.newBalance === null || editModal.newBalance === ('' as any) || isNaN(Number(editModal.newBalance)) || saving"
            class="flex-1 py-2.5 bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 text-white rounded-xl text-sm font-bold transition-all"
          >
            {{ saving ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useProjectsStore } from '../store/projects-store'
import type { ProjectCapitalResponse } from '../infrastructure/projects-api'

const props = defineProps<{
  projectId: string
  teamId?: number | null
}>()

const store = useProjectsStore()
const capital = computed(() => store.projectCapital)
const loading = ref(false)
const saving = ref(false)

// Cache burn rate locally so modal can use it even before store.projectCapital is loaded
const localTeamMonthlyCost = ref(0)

// Modals
const showInjectModal = ref(false)
const showEditModal = ref(false)

const injectModal = reactive({
  deliveryDate: '' as string,
  amount: null as number | null,
  bonusPct: 0,
  note: '',
})
const editModal = ref({ newBalance: null as number | null, bonusPct: null as number | null, note: '' })
const deletingTxId = ref<number | null>(null)

// ---------- computed ----------

const runwayPct = computed(() => {
  if (!capital.value || capital.value.team_monthly_cost === 0) return 0
  return (capital.value.runway_months / 3) * 100
})

const runwayColor = computed(() => {
  const m = capital.value?.runway_months ?? 0
  if (m > 2) return 'text-emerald-400'
  if (m > 1) return 'text-yellow-400'
  return 'text-red-400'
})

const runwayBarColor = computed(() => {
  const m = capital.value?.runway_months ?? 0
  if (m > 2) return 'bg-emerald-500'
  if (m > 1) return 'bg-yellow-500'
  return 'bg-red-500'
})

const editDelta = computed(() => {
  if (editModal.value.newBalance === null || !capital.value) return 0
  return editModal.value.newBalance - capital.value.capital_balance
})

const injectModalDeliveryDateMin = computed(() => {
  const d = new Date()
  d.setDate(d.getDate() + 1)
  return d.toISOString().split('T')[0]
})

// จำนวนเดือนจากวันนี้ถึงวันส่งงาน
const injectModalMonthsToDelivery = computed(() => {
  if (!injectModal.deliveryDate || !localTeamMonthlyCost.value) return 0
  const start = new Date()
  start.setHours(0, 0, 0, 0)
  const end = new Date(injectModal.deliveryDate)
  end.setHours(0, 0, 0, 0)
  const days = Math.max(0, (end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24))
  return days / (365 / 12) // ~30.44 days per month
})

// ยอดเงินที่ต้อง inject เพื่อให้หมดพอดีวันส่งงาน
const injectModalAmountFromDate = computed(() => {
  if (injectModalMonthsToDelivery.value <= 0 || !localTeamMonthlyCost.value) return null
  return Math.round(localTeamMonthlyCost.value * injectModalMonthsToDelivery.value)
})

watch(
  () => injectModal.deliveryDate,
  () => {
    const suggested = injectModalAmountFromDate.value
    if (suggested != null && suggested > 0) injectModal.amount = suggested
  }
)

// ---------- helpers ----------

function formatMoney(v: number) {
  return '฿' + Math.round(v).toLocaleString('th-TH')
}

function formatDate(s: string) {
  return new Date(s).toLocaleString('th-TH', { dateStyle: 'medium', timeStyle: 'short' })
}

function txIcon(type: string) {
  const icons: Record<string, string> = { INJECTION: '💰', BURN: '🔥', BONUS_PAYOUT: '🎁', ADJUSTMENT: '✏️' }
  return icons[type] ?? '•'
}

function txBadgeClass(type: string) {
  const cls: Record<string, string> = {
    INJECTION: 'bg-emerald-500/20 text-emerald-400',
    BURN: 'bg-red-500/20 text-red-400',
    BONUS_PAYOUT: 'bg-yellow-500/20 text-yellow-400',
    ADJUSTMENT: 'bg-indigo-500/20 text-indigo-400',
  }
  return cls[type] ?? 'bg-gray-500/20 text-gray-400'
}

// ---------- modal openers ----------

async function openInjectModal() {
  if (!capital.value) {
    loading.value = true
    try {
      await store.fetchProjectCapital(props.projectId)
      if (capital.value?.team_monthly_cost) {
        localTeamMonthlyCost.value = capital.value.team_monthly_cost
      }
    } finally {
      loading.value = false
    }
  }
  injectModal.deliveryDate = ''
  injectModal.amount = null
  injectModal.bonusPct = capital.value?.bonus_percentage ?? 0
  injectModal.note = ''
  showInjectModal.value = true
}

function openEditModal() {
  editModal.value = {
    newBalance: capital.value?.capital_balance ?? 0,
    bonusPct: capital.value?.bonus_percentage ?? null,
    note: '',
  }
  showEditModal.value = true
}

// ---------- actions ----------

async function confirmDeleteTx(txId: number) {
  deletingTxId.value = null
  await store.deleteProjectTransaction(props.projectId, txId)
}

async function confirmInject() {
  if (!injectModal.amount || injectModal.amount <= 0) return
  saving.value = true
  try {
    await store.injectProjectCapital(props.projectId, {
      amount: injectModal.amount,
      bonus_percentage: injectModal.bonusPct,
      note: injectModal.note || (injectModal.deliveryDate ? `ส่งงาน ${injectModal.deliveryDate}` : 'Capital injection'),
    })
    await store.fetchProjectCapital(props.projectId)
    showInjectModal.value = false
  } finally {
    saving.value = false
  }
}

async function confirmEdit() {
  if (editModal.value.newBalance === null) return
  const newBalance = Number(editModal.value.newBalance)
  if (isNaN(newBalance)) return
  saving.value = true
  try {
    const bonusPct = editModal.value.bonusPct
    const payload: { new_balance: number; bonus_percentage?: number; note: string } = {
      new_balance: newBalance,
      note: editModal.value.note,
    }
    const parsed = Number(bonusPct)
    if (bonusPct !== null && bonusPct !== undefined && !isNaN(parsed)) {
      payload.bonus_percentage = parsed
    }
    await store.editProjectCapital(props.projectId, payload)
    await store.fetchProjectCapital(props.projectId)
    showEditModal.value = false
  } finally {
    saving.value = false
  }
}

// ---------- lifecycle ----------

// Sync localTeamMonthlyCost whenever capital data arrives/updates
watch(capital, (val) => {
  if (val?.team_monthly_cost) {
    localTeamMonthlyCost.value = val.team_monthly_cost
  }
}, { immediate: true })

onMounted(async () => {
  loading.value = true
  try {
    await store.fetchProjectCapital(props.projectId)
    if (capital.value?.team_monthly_cost) {
      localTeamMonthlyCost.value = capital.value.team_monthly_cost
    }
  } catch (err) {
    console.error('[ProjectCapitalPanel] fetchProjectCapital failed:', err)
  } finally {
    loading.value = false
  }
})
</script>
