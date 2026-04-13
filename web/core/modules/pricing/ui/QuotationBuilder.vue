<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-white">Fully Loaded Cost Quotation</h2>
        <p class="mt-1 text-sm text-gray-400">
          Calculate project cost using the company-wide fully loaded cost model.
        </p>
      </div>
      <div v-if="store.hasResult" class="flex items-center gap-2">
        <span class="rounded-full bg-amber-500/10 px-3 py-1 text-xs font-semibold text-amber-400 border border-amber-500/20">
          Grand Total: {{ formatTHB(store.grandTotal) }}
        </span>
      </div>
    </div>

    <!-- Error banner -->
    <div
      v-if="store.error"
      class="flex items-start gap-3 rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400"
    >
      <svg class="mt-0.5 h-4 w-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <span>{{ store.error }}</span>
    </div>

    <!-- ── Cost Parameters (from Admin Cost Config) ──────────────────────── -->
    <div class="rounded-xl border border-gray-700 bg-gray-800/60 p-5">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-sm font-semibold uppercase tracking-widest text-amber-400">
          Cost Parameters
        </h3>
        <div class="flex items-center gap-2">
          <span v-if="loadingConfig" class="flex items-center gap-1.5 text-xs text-gray-500">
            <svg class="h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            Loading…
          </span>
          <span v-else class="rounded-full bg-emerald-500/10 px-2.5 py-0.5 text-xs font-medium text-emerald-400 border border-emerald-500/20">
            Synced from Cost Config
          </span>
          <NuxtLink
            to="/admin/cost-config"
            class="rounded-lg border border-gray-600 px-2.5 py-1 text-xs text-gray-400 hover:border-gray-500 hover:text-white transition-colors"
          >
            ⚙️ Edit
          </NuxtLink>
        </div>
      </div>

      <div v-if="configError" class="mb-4 flex items-start gap-3 rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400">
        <svg class="mt-0.5 h-4 w-4 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <span>{{ configError }}</span>
      </div>

      <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
        <!-- Cost / Manday -->
        <div class="col-span-2 sm:col-span-3 lg:col-span-2 rounded-xl border border-purple-500/25 bg-gradient-to-br from-purple-900/40 to-pink-900/30 px-4 py-3">
          <p class="text-xs font-semibold uppercase tracking-widest text-purple-300 mb-1">Cost / Manday</p>
          <p class="text-2xl font-black text-white tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-base">—</span>
            <span v-else>{{ formatTHB(costPerManday) }}</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">Fully loaded ÷ billable days</p>
        </div>
        <!-- Cost / Hour -->
        <div class="rounded-xl border border-cyan-500/20 bg-gray-900/40 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-widest text-cyan-300 mb-1">Cost / Hour</p>
          <p class="text-lg font-extrabold text-white tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-sm">—</span>
            <span v-else>{{ formatTHB(costPerHour) }}</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">Manday ÷ {{ workingHoursPerDay }}h</p>
        </div>
        <!-- Billable Days -->
        <div class="rounded-xl border border-amber-500/20 bg-gray-900/40 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-widest text-amber-300 mb-1">Billable Days</p>
          <p class="text-lg font-extrabold text-amber-400 tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-sm">—</span>
            <span v-else>{{ billableDays.toFixed(1) }}</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">of {{ workingDaysPerMonth }} days/mo</p>
        </div>
        <!-- Utilisation -->
        <div class="rounded-xl border border-emerald-500/20 bg-gray-900/40 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-widest text-emerald-300 mb-1">Utilisation</p>
          <p class="text-lg font-extrabold text-emerald-400 tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-sm">—</span>
            <span v-else>{{ (utilizationRate * 100).toFixed(0) }}%</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">1 ÷ {{ overheadMultiplier }}×</p>
        </div>
        <!-- Risk Buffer -->
        <div class="rounded-xl border border-amber-500/20 bg-amber-500/5 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-widest text-amber-500/70 mb-1">Risk Buffer</p>
          <p class="text-lg font-extrabold text-amber-400 tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-sm">—</span>
            <span v-else>{{ (form.risk_margin_pct * 100).toFixed(0) }}%</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">default margin</p>
        </div>
        <!-- Profit Margin -->
        <div class="rounded-xl border border-purple-500/20 bg-purple-500/5 px-3 py-2.5">
          <p class="text-xs font-semibold uppercase tracking-widest text-purple-500/70 mb-1">Profit Margin</p>
          <p class="text-lg font-extrabold text-purple-400 tabular-nums">
            <span v-if="loadingConfig" class="text-gray-500 text-sm">—</span>
            <span v-else>{{ (form.profit_margin_pct * 100).toFixed(0) }}%</span>
          </p>
          <p class="text-xs text-gray-500 mt-1">target margin</p>
        </div>
      </div>
    </div>

    <!-- ── Actions ──────────────────────────────────────────────────────── -->
    <div class="flex flex-wrap items-center gap-3">
      <button
        :disabled="store.loading"
        class="btn-primary"
        @click="openTaskSelectionModal"
      >
        <svg v-if="store.loading" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        <svg v-else class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 11h.01M12 11h.01M15 11h.01M4 19h16a2 2 0 002-2V7a2 2 0 00-2-2H4a2 2 0 00-2 2v10a2 2 0 002 2z"/>
        </svg>
        {{ store.loading ? 'Calculating…' : 'Calculate Cost' }}
      </button>

      <button
        v-if="store.hasResult"
        :disabled="store.exporting"
        class="btn-secondary"
        @click="exportPDF"
      >
        <svg v-if="store.exporting" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        <svg v-else class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
        {{ store.exporting ? 'Generating PDF…' : 'Export PDF Quotation' }}
      </button>

      <button
        v-if="store.hasResult"
        :disabled="exportingCustomer"
        class="btn-primary"
        @click="exportCustomerPDF"
      >
        <svg v-if="exportingCustomer" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
        <svg v-else class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
        </svg>
        {{ exportingCustomer ? 'Generating PDF…' : 'Export PDF for Customer' }}
      </button>

      <div
        v-if="store.hasResult"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-600 bg-gray-900/40 px-2 py-1"
      >
        <span class="text-xs text-gray-400">Customer Line Pricing</span>
        <button
          class="rounded-md px-2.5 py-1 text-xs font-semibold transition-colors"
          :class="customerLinePricingMode === 'base' ? 'bg-slate-700 text-white' : 'text-gray-400 hover:text-white'"
          @click="customerLinePricingMode = 'base'"
        >
          แบบเดิม
        </button>
        <button
          class="rounded-md px-2.5 py-1 text-xs font-semibold transition-colors"
          :class="customerLinePricingMode === 'absorbed' ? 'bg-purple-700/80 text-white' : 'text-gray-400 hover:text-white'"
          @click="customerLinePricingMode = 'absorbed'"
        >
          รวมความเสี่ยง+กำไร
        </button>
      </div>
    </div>

    <!-- ── MA Quotation Calculator ───────────────────────────────────────── -->
    <div class="rounded-2xl border border-slate-700/80 bg-gradient-to-br from-slate-900 via-slate-900 to-indigo-950/40 p-5 shadow-2xl shadow-indigo-500/10">
      <div class="flex flex-wrap items-start justify-between gap-3 border-b border-slate-700/70 pb-4">
        <div>
          <h3 class="text-base font-bold text-white">MA Quotation Builder</h3>
          <p class="mt-1 text-xs text-slate-400">
            คำนวณใบเสนอราคา MA จากราคาโครงการและเปอร์เซ็นต์ MA เพื่อใช้เสนอราคาอย่างมืออาชีพ
          </p>
        </div>
        <span class="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-xs font-semibold text-emerald-300">
          Premium Proposal Style
        </span>
      </div>

      <div class="mt-5 grid gap-4 md:grid-cols-3">
        <label class="space-y-2">
          <span class="text-xs font-semibold uppercase tracking-widest text-slate-400">MA Price (THB / ปี)</span>
          <input
            v-model.number="maForm.maPrice"
            type="number"
            min="0"
            step="1000"
            class="input-field"
            placeholder="เช่น 300000"
          >
        </label>
        <label class="space-y-2">
          <span class="text-xs font-semibold uppercase tracking-widest text-slate-400">ระยะเวลา MA (ปี)</span>
          <input
            v-model.number="maForm.maDurationYears"
            type="number"
            min="1"
            step="1"
            class="input-field"
            placeholder="เช่น 1"
          >
        </label>
        <div class="flex items-end">
          <button class="btn-primary w-full justify-center" @click="calculateMAQuotation">
            Generate MA Quotation
          </button>
        </div>
      </div>

      <!-- MA Task Scope -->
      <div class="mt-4 rounded-xl border border-slate-700/60 bg-slate-800/40 px-4 py-3">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <svg class="h-4 w-4 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01"/>
            </svg>
            <span class="text-xs font-semibold uppercase tracking-widest text-slate-400">MA Task Scope</span>
            <span v-if="maSelectedTaskIds.size > 0" class="rounded-full border border-indigo-500/30 bg-indigo-500/20 px-2 py-0.5 text-xs font-semibold text-indigo-300">
              {{ maSelectedTaskIds.size }} tasks
            </span>
          </div>
          <button
            class="flex items-center gap-1.5 rounded-lg border border-slate-600 px-2.5 py-1.5 text-xs text-slate-300 transition-colors hover:border-indigo-500/60 hover:text-white"
            @click="openMATaskModal"
          >
            <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            Manage Tasks
          </button>
        </div>
        <div v-if="maSelectedTaskIds.size === 0" class="mt-2.5 text-xs italic text-slate-500">
          ยังไม่ได้เลือก task — กด Manage Tasks เพื่อเลือก task ที่อยู่ใน scope ของ MA
        </div>
        <div v-else class="mt-3 max-h-56 space-y-3 overflow-y-auto pr-1">
          <div v-for="epic in epicsWithTasks.filter(e => maEpicSelectedCount(e) > 0)" :key="epic.id">
            <p class="mb-1 flex items-center gap-1.5 text-xs font-medium text-indigo-300/80">
              <span class="h-1.5 w-1.5 flex-shrink-0 rounded-full" :style="{ backgroundColor: epic.color || '#6366f1' }"/>
              {{ epic.title }}
            </p>
            <div class="ml-3 space-y-0.5">
              <div
                v-for="task in epic.tasks.filter(t => maSelectedTaskIds.has(t.id))"
                :key="task.id"
                class="flex items-center gap-2 rounded px-2 py-1 bg-slate-900/50"
              >
                <span class="flex-1 truncate text-xs text-slate-300">{{ task.title }}</span>
                <span class="flex-shrink-0 text-xs text-slate-600">{{ task.code }}</span>
              </div>
            </div>
          </div>
          <div v-if="maNoEpicSelectedCount > 0">
            <p class="mb-1 flex items-center gap-1.5 text-xs font-medium text-slate-500">
              <span class="h-1.5 w-1.5 flex-shrink-0 rounded-full bg-slate-500"/>
              No Epic
            </p>
            <div class="ml-3 space-y-0.5">
              <div
                v-for="task in tasksWithoutEpic.filter(t => maSelectedTaskIds.has(t.id))"
                :key="task.id"
                class="flex items-center gap-2 rounded px-2 py-1 bg-slate-900/50"
              >
                <span class="flex-1 truncate text-xs text-slate-300">{{ task.title }}</span>
                <span class="flex-shrink-0 text-xs text-slate-600">{{ task.code }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="maResult" class="mt-5 overflow-hidden rounded-2xl border border-slate-600/70 bg-slate-950/60">
        <div class="bg-gradient-to-r from-indigo-600/20 via-cyan-500/10 to-emerald-500/20 px-5 py-4">
          <div class="flex flex-wrap items-start justify-between gap-3">
            <div>
              <h4 class="text-sm font-bold uppercase tracking-[0.18em] text-slate-200">Maintenance Agreement Quotation</h4>
              <p class="mt-1 text-xs text-slate-400">สรุปข้อมูล MA — กด Export เพื่อออกเอกสารให้ลูกค้า</p>
            </div>
            <button
              :disabled="exportingMA"
              class="btn-secondary"
              @click="exportMAPDF"
            >
              <svg v-if="exportingMA" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              <svg v-else class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
              </svg>
              {{ exportingMA ? 'Generating MA PDF…' : 'Export MA PDF' }}
            </button>
          </div>
        </div>

        <div class="grid gap-3 border-b border-slate-700/70 px-5 py-4 sm:grid-cols-3">
          <div class="premium-kpi-card">
            <p class="premium-kpi-label">MA Annual Price</p>
            <p class="premium-kpi-value text-cyan-300">{{ formatTHB(maResult.maPrice) }}</p>
          </div>
          <div class="premium-kpi-card">
            <p class="premium-kpi-label">ระยะเวลา MA</p>
            <p class="premium-kpi-value text-indigo-300">{{ maResult.maDurationYears }} ปี</p>
          </div>
          <div class="premium-kpi-card">
            <p class="premium-kpi-label">Monthly Fee</p>
            <p class="premium-kpi-value text-emerald-300">{{ formatTHB(maResult.monthlyFee) }}</p>
          </div>
        </div>

        <div class="divide-y divide-slate-700/60 px-5 py-2 text-sm">
          <div class="flex items-center justify-between py-3">
            <span class="text-slate-400">MA Annual Price</span>
            <span class="font-semibold text-slate-100">{{ formatTHB(maResult.maPrice) }}</span>
          </div>
          <div class="flex items-center justify-between py-3">
            <span class="text-slate-400">Monthly Fee (Annual ÷ 12)</span>
            <span class="font-semibold text-slate-200">{{ formatTHB(maResult.monthlyFee) }}</span>
          </div>
          <div class="flex items-center justify-between py-3">
            <span class="text-sm font-bold uppercase tracking-wide text-white">ระยะเวลา MA</span>
            <span class="text-lg font-black text-white">{{ maResult.maDurationYears }} ปี</span>
          </div>
        </div>
      </div>

      <!-- Delivery Milestones -->
      <div class="mt-4 rounded-xl border border-slate-700/60 bg-slate-800/40 px-4 py-4">
        <div class="mb-3 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <svg class="h-4 w-4 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
            </svg>
            <span class="text-xs font-semibold uppercase tracking-widest text-slate-400">งวดส่งมอบ & เบิกเงิน</span>
          </div>
          <button
            class="flex items-center gap-1 rounded-lg border border-slate-600 px-2.5 py-1.5 text-xs text-slate-300 transition-colors hover:border-cyan-500/60 hover:text-white"
            @click="addMilestone"
          >
            <svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            Add งวด
          </button>
        </div>

        <!-- Column headers -->
        <div class="mb-1.5 grid grid-cols-[24px_1fr_120px_140px_28px] items-center gap-2 px-3 text-xs font-semibold uppercase tracking-wider text-slate-600">
          <span></span>
          <span>ชื่องวด</span>
          <span class="text-center">จำนวน Task</span>
          <span class="text-right">จำนวนเงิน (THB)</span>
          <span></span>
        </div>

        <div v-if="maDeliveryMilestones.length === 0" class="py-4 text-center text-xs italic text-slate-500">
          ยังไม่มีงวดส่งมอบ — กด Add งวด เพื่อเพิ่ม
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="(milestone, idx) in maDeliveryMilestones"
            :key="milestone.id"
            class="grid grid-cols-[24px_1fr_120px_140px_28px] items-center gap-2 rounded-lg border px-3 py-2.5"
            :class="milestone.isEndOfMA
              ? 'border-cyan-700/40 bg-cyan-950/30'
              : 'border-slate-700/40 bg-slate-900/50'"
          >
            <!-- Index badge -->
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full text-xs font-bold"
              :class="milestone.isEndOfMA ? 'bg-cyan-600/30 text-cyan-300' : 'bg-indigo-600/30 text-indigo-300'"
            >
              {{ idx + 1 }}
            </span>

            <!-- Label -->
            <input
              v-model="milestone.label"
              type="text"
              class="min-w-0 rounded-lg border border-slate-700 bg-slate-800 px-2.5 py-1.5 text-xs text-white placeholder-slate-500 focus:border-indigo-500 focus:outline-none"
              placeholder="ชื่องวด"
            />

            <!-- Task count (or end-of-MA badge) -->
            <div class="flex items-center justify-center gap-1">
              <template v-if="!milestone.isEndOfMA">
                <input
                  v-model.number="milestone.taskCount"
                  type="number"
                  min="0"
                  class="w-16 rounded-lg border border-slate-700 bg-slate-800 px-2 py-1.5 text-center text-xs text-white focus:border-indigo-500 focus:outline-none"
                  placeholder="0"
                />
                <span class="text-xs text-slate-500">tasks</span>
              </template>
              <span v-else class="rounded-full border border-cyan-700/50 bg-cyan-900/30 px-2 py-0.5 text-xs text-cyan-400">
                สิ้นสุด MA
              </span>
            </div>

            <!-- Amount -->
            <div class="flex items-center gap-1">
              <input
                v-model.number="milestone.amount"
                type="number"
                min="0"
                step="1000"
                class="min-w-0 flex-1 rounded-lg border border-slate-700 bg-slate-800 px-2 py-1.5 text-right text-xs text-white focus:border-indigo-500 focus:outline-none"
                placeholder="0"
              />
            </div>

            <!-- Remove (only non-end-of-MA) -->
            <button
              v-if="!milestone.isEndOfMA"
              class="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded text-slate-600 transition-colors hover:bg-red-500/10 hover:text-red-400"
              @click="removeMilestone(milestone.id)"
            >
              <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
            <span v-else class="h-6 w-6" />
          </div>

          <!-- Summary row -->
          <div class="mt-1 flex items-center justify-between rounded-lg bg-slate-900/30 px-3 py-2">
            <div class="flex flex-col gap-0.5">
              <span class="text-xs text-slate-500">รวมทั้งหมด (TOTAL)</span>
              <span class="text-[10px] text-slate-500">*ไม่รวม VAT</span>
            </div>
            <div class="flex items-center gap-3">
              <span class="text-xs text-slate-500">
                {{ milestoneTotalAmount > 0 ? formatTHB(milestoneTotalAmount) : '—' }}
              </span>
              <span
                v-if="maResult"
                class="text-xs font-semibold"
                :class="milestoneAmountDiff === 0 ? 'text-emerald-400' : 'text-amber-400'"
              >
                {{ milestoneAmountDiff === 0 ? '✓ ครบ' : (milestoneAmountDiff > 0 ? `ขาด ${formatTHB(milestoneAmountDiff)}` : `เกิน ${formatTHB(Math.abs(milestoneAmountDiff))}`) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ── Results ───────────────────────────────────────────────────────── -->
    <template v-if="store.hasResult && store.result">
      <!-- Model Metrics -->
      <div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
        <div class="metric-card">
          <p class="metric-label">Cost / Manday</p>
          <p class="metric-value text-amber-400">{{ formatTHB(store.result.cost_per_manday) }}</p>
        </div>
        <div class="metric-card">
          <p class="metric-label">Total Mandays</p>
          <p class="metric-value">{{ store.result.total_mandays.toFixed(2) }}</p>
        </div>
        <div class="metric-card">
          <p class="metric-label">Tasks Costed</p>
          <p class="metric-value">{{ store.result.tasks.length }}</p>
        </div>
      </div>

      <!-- Task Breakdown Table -->
      <div class="rounded-xl border border-gray-700 bg-gray-800/60 overflow-hidden">
        <div class="flex items-center justify-between border-b border-gray-700 px-5 py-3">
          <div>
            <h3 class="text-sm font-semibold text-white">Task Cost Breakdown</h3>
            <p class="mt-0.5 text-xs text-gray-500">
              {{ customerLinePricingMode === 'absorbed' ? 'แสดงแบบรวมความเสี่ยงและกำไรในแต่ละรายการ (สำหรับ Customer PDF)' : 'แสดงต้นทุนฐานของแต่ละรายการ (แบบเดิม)' }}
            </p>
          </div>
          <span class="text-xs text-gray-500">{{ store.result.currency }}</span>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-700 bg-gray-900/50 text-xs uppercase tracking-wider text-gray-500">
                <th class="px-4 py-3 text-left">Epic</th>
                <th class="px-4 py-3 text-left">Task</th>
                <th class="px-4 py-3 text-right">Mandays</th>
                <th class="px-4 py-3 text-right">Cost (THB)</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(task, idx) in store.result.tasks"
                :key="task.task_id"
                class="border-b border-gray-700/50 transition-colors hover:bg-gray-700/20"
                :class="idx % 2 === 0 ? '' : 'bg-gray-900/20'"
              >
                <td class="px-4 py-3 text-gray-400 text-xs">{{ task.epic_title || '—' }}</td>
                <td class="px-4 py-3 text-gray-200 font-medium">{{ task.title }}</td>
                <td class="px-4 py-3 text-right text-gray-300">{{ task.mandays.toFixed(2) }}</td>
                <td class="px-4 py-3 text-right font-semibold text-amber-400">
                  {{ formatNumber(getTaskDisplayAmount(task, idx)) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Financial Summary -->
      <div class="ml-auto w-full max-w-sm rounded-xl border border-gray-700 bg-gray-800/60 overflow-hidden">
        <div class="border-b border-gray-700 px-5 py-3">
          <h3 class="text-sm font-semibold text-white">Quotation Summary</h3>
        </div>
        <div class="divide-y divide-gray-700/50">
          <div class="summary-row">
            <span class="summary-label">Labor Subtotal</span>
            <span class="summary-amount">{{ formatTHB(store.result.subtotal) }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">
              Risk Buffer
              <span class="ml-1 text-xs text-amber-500/80">({{ (form.risk_margin_pct * 100).toFixed(0) }}%)</span>
            </span>
            <span class="summary-amount text-amber-400">+ {{ formatTHB(store.result.risk_amount) }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">
              Profit Margin
              <span class="ml-1 text-xs text-purple-400/80">({{ (form.profit_margin_pct * 100).toFixed(0) }}%)</span>
            </span>
            <span class="summary-amount text-purple-400">+ {{ formatTHB(store.result.profit_amount) }}</span>
          </div>
          <div class="summary-row border-t border-gray-600">
            <span class="summary-label font-medium text-gray-300">Total (before VAT)</span>
            <span class="summary-amount font-semibold text-white">{{ formatTHB(store.result.subtotal + store.result.risk_amount + store.result.profit_amount) }}</span>
          </div>
          <div class="summary-row">
            <span class="summary-label">VAT (7%)</span>
            <span class="summary-amount text-gray-300">+ {{ formatTHB(store.result.vat) }}</span>
          </div>
          <div class="summary-row bg-gradient-to-r from-amber-500/10 to-purple-500/10">
            <span class="text-sm font-bold text-white">Grand Total</span>
            <span class="text-lg font-extrabold text-white">{{ formatTHB(store.result.grand_total) }}</span>
          </div>
        </div>
      </div>
    </template>
  </div>

  <!-- ── Task Selection Modal ───────────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="showTaskModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="showTaskModal = false"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="showTaskModal = false" />

        <!-- Modal Panel -->
        <div class="relative z-10 flex w-full max-w-2xl flex-col rounded-2xl border border-gray-700 bg-gray-900 shadow-2xl max-h-[85vh]">
          <!-- Modal Header -->
          <div class="flex items-center justify-between border-b border-gray-700 px-6 py-4">
            <div>
              <h3 class="text-base font-bold text-white">Select Scope for Calculation</h3>
              <p class="mt-0.5 text-xs text-gray-400">เลือก Epic และ Task ที่ต้องการนำมาคำนวณต้นทุน</p>
            </div>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-700 hover:text-gray-900 dark:text-white transition-colors"
              @click="showTaskModal = false"
            >
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Loading state -->
          <div v-if="loadingEpics" class="flex flex-1 items-center justify-center py-16">
            <svg class="h-6 w-6 animate-spin text-amber-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <span class="ml-3 text-sm text-gray-400">Loading tasks…</span>
          </div>

          <!-- Empty state -->
          <div v-else-if="epicsWithTasks.length === 0 && tasksWithoutEpic.length === 0" class="flex flex-1 flex-col items-center justify-center py-16 text-center">
            <svg class="h-10 w-10 text-gray-600 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
            <p class="text-sm text-gray-400">No tasks found for this project.</p>
            <p class="text-xs text-gray-500 mt-1">Create tasks with start/end dates to enable cost calculation.</p>
          </div>

          <!-- Epic/Task List -->
          <div v-else class="flex-1 overflow-y-auto px-6 py-4 space-y-4">
            <!-- Select All / Deselect All row -->
            <div class="flex items-center justify-between rounded-lg border border-gray-700/50 bg-gray-800/40 px-4 py-2.5">
              <span class="text-xs font-medium text-gray-400">
                {{ selectedTaskIds.size }} / {{ totalTaskCount }} tasks selected
              </span>
              <div class="flex items-center gap-3">
                <button
                  class="text-xs text-amber-400 hover:text-amber-300 transition-colors"
                  @click="selectAll"
                >
                  Select All
                </button>
                <span class="text-gray-600">|</span>
                <button
                  class="text-xs text-gray-400 hover:text-gray-300 transition-colors"
                  @click="deselectAll"
                >
                  Deselect All
                </button>
              </div>
            </div>

            <!-- Epics -->
            <div
              v-for="epic in epicsWithTasks"
              :key="epic.id"
              class="rounded-xl border border-gray-700/60 bg-gray-800/30 overflow-hidden"
            >
              <!-- Epic Header -->
              <div
                class="flex items-center gap-3 px-4 py-3 cursor-pointer select-none hover:bg-gray-700/30 transition-colors"
                @click="toggleEpicExpand(epic.id)"
              >
                <!-- Epic color dot -->
                <span
                  class="h-2.5 w-2.5 flex-shrink-0 rounded-full"
                  :style="{ backgroundColor: epic.color || '#6366f1' }"
                />
                <!-- Epic checkbox -->
                <input
                  type="checkbox"
                  :checked="isEpicFullySelected(epic)"
                  :indeterminate="isEpicPartiallySelected(epic)"
                  class="h-4 w-4 rounded accent-amber-500 flex-shrink-0"
                  @click.stop
                  @change="toggleEpicSelection(epic)"
                />
                <span class="flex-1 text-sm font-semibold text-white">{{ epic.title }}</span>
                <span class="rounded-full bg-gray-700 px-2 py-0.5 text-xs text-gray-400">
                  {{ epicSelectedCount(epic) }}/{{ epic.tasks.length }}
                </span>
                <!-- Expand chevron -->
                <svg
                  class="h-4 w-4 text-gray-500 transition-transform"
                  :class="expandedEpics.has(epic.id) ? 'rotate-180' : ''"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </div>

              <!-- Task list (collapsible) -->
              <div v-show="expandedEpics.has(epic.id)" class="divide-y divide-gray-700/40">
                <label
                  v-for="task in epic.tasks"
                  :key="task.id"
                  class="flex cursor-pointer items-start gap-3 px-5 py-2.5 hover:bg-gray-700/20 transition-colors"
                  :class="selectedTaskIds.has(task.id) ? 'bg-gray-700/10' : ''"
                >
                  <input
                    v-model="selectedTaskIdsArr"
                    :value="task.id"
                    type="checkbox"
                    class="mt-0.5 h-3.5 w-3.5 rounded accent-amber-500 flex-shrink-0"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="text-sm text-gray-200 leading-snug">{{ task.title }}</p>
                    <div class="mt-1 flex flex-wrap items-center gap-2">
                      <span class="text-xs text-gray-500">{{ task.code }}</span>
                      <span v-if="task.start_date && task.end_date" class="text-xs text-gray-500">
                        {{ formatDate(task.start_date) }} → {{ formatDate(task.end_date) }}
                      </span>
                      <span v-else class="text-xs text-red-400/70">No dates set</span>
                      <span
                        class="rounded px-1.5 py-0.5 text-xs font-medium"
                        :class="priorityClass(task.priority)"
                      >
                        {{ task.priority }}
                      </span>
                    </div>
                  </div>
                </label>
              </div>
            </div>

            <!-- Tasks without epic -->
            <div
              v-if="tasksWithoutEpic.length > 0"
              class="rounded-xl border border-gray-700/60 bg-gray-800/30 overflow-hidden"
            >
              <div
                class="flex items-center gap-3 px-4 py-3 cursor-pointer select-none hover:bg-gray-700/30 transition-colors"
                @click="toggleEpicExpand('__no_epic__')"
              >
                <span class="h-2.5 w-2.5 flex-shrink-0 rounded-full bg-gray-500" />
                <input
                  type="checkbox"
                  :checked="isNoEpicFullySelected"
                  :indeterminate="isNoEpicPartiallySelected"
                  class="h-4 w-4 rounded accent-amber-500 flex-shrink-0"
                  @click.stop
                  @change="toggleNoEpicSelection"
                />
                <span class="flex-1 text-sm font-semibold text-gray-400">No Epic</span>
                <span class="rounded-full bg-gray-700 px-2 py-0.5 text-xs text-gray-400">
                  {{ noEpicSelectedCount }}/{{ tasksWithoutEpic.length }}
                </span>
                <svg
                  class="h-4 w-4 text-gray-500 transition-transform"
                  :class="expandedEpics.has('__no_epic__') ? 'rotate-180' : ''"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </div>
              <div v-show="expandedEpics.has('__no_epic__')" class="divide-y divide-gray-700/40">
                <label
                  v-for="task in tasksWithoutEpic"
                  :key="task.id"
                  class="flex cursor-pointer items-start gap-3 px-5 py-2.5 hover:bg-gray-700/20 transition-colors"
                  :class="selectedTaskIds.has(task.id) ? 'bg-gray-700/10' : ''"
                >
                  <input
                    v-model="selectedTaskIdsArr"
                    :value="task.id"
                    type="checkbox"
                    class="mt-0.5 h-3.5 w-3.5 rounded accent-amber-500 flex-shrink-0"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="text-sm text-gray-200 leading-snug">{{ task.title }}</p>
                    <div class="mt-1 flex flex-wrap items-center gap-2">
                      <span class="text-xs text-gray-500">{{ task.code }}</span>
                      <span v-if="task.start_date && task.end_date" class="text-xs text-gray-500">
                        {{ formatDate(task.start_date) }} → {{ formatDate(task.end_date) }}
                      </span>
                      <span v-else class="text-xs text-red-400/70">No dates set</span>
                      <span
                        class="rounded px-1.5 py-0.5 text-xs font-medium"
                        :class="priorityClass(task.priority)"
                      >
                        {{ task.priority }}
                      </span>
                    </div>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Modal Footer -->
          <div class="flex items-center justify-between border-t border-gray-700 px-6 py-4">
            <p class="text-xs text-gray-500">
              {{ selectedTaskIds.size === 0 ? 'Select tasks to include in the calculation' : `${selectedTaskIds.size} task(s) will be costed` }}
            </p>
            <div class="flex items-center gap-3">
              <button
                class="rounded-lg border border-gray-600 px-4 py-2 text-sm text-gray-400 hover:border-gray-500 hover:text-white transition-colors"
                @click="showTaskModal = false"
              >
                Cancel
              </button>
              <button
                :disabled="selectedTaskIds.size === 0 || store.loading"
                class="btn-primary"
                @click="confirmAndCalculate"
              >
                <svg v-if="store.loading" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                {{ store.loading ? 'Calculating…' : 'Calculate' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>

  <!-- ── MA Task Selection Modal ─────────────────────────────────────────── -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div
        v-if="showMATaskModal"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click.self="showMATaskModal = false"
      >
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="showMATaskModal = false" />

        <div class="relative z-10 flex max-h-[85vh] w-full max-w-2xl flex-col rounded-2xl border border-slate-700 bg-slate-900 shadow-2xl">
          <!-- Header -->
          <div class="flex items-center justify-between border-b border-slate-700 px-6 py-4">
            <div>
              <h3 class="text-base font-bold text-white">MA Task Scope</h3>
              <p class="mt-0.5 text-xs text-slate-400">เลือก Task ที่อยู่ใน scope ของ MA (ไม่คำนวณราคา — เพื่อระบุขอบเขตงานเท่านั้น)</p>
            </div>
            <button
              class="flex h-8 w-8 items-center justify-center rounded-lg text-slate-500 transition-colors hover:bg-slate-700 hover:text-white"
              @click="showMATaskModal = false"
            >
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- Loading -->
          <div v-if="loadingEpics" class="flex flex-1 items-center justify-center py-16">
            <svg class="h-6 w-6 animate-spin text-indigo-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <span class="ml-3 text-sm text-slate-400">Loading tasks…</span>
          </div>

          <!-- Empty -->
          <div v-else-if="epicsWithTasks.length === 0 && tasksWithoutEpic.length === 0" class="flex flex-1 flex-col items-center justify-center py-16 text-center">
            <svg class="mb-3 h-10 w-10 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
            <p class="text-sm text-slate-400">No tasks found for this project.</p>
          </div>

          <!-- Task List -->
          <div v-else class="flex-1 space-y-4 overflow-y-auto px-6 py-4">
            <div class="flex items-center justify-between rounded-lg border border-slate-700/50 bg-slate-800/40 px-4 py-2.5">
              <span class="text-xs font-medium text-slate-400">
                {{ maSelectedTaskIds.size }} / {{ totalTaskCount }} tasks selected
              </span>
              <div class="flex items-center gap-3">
                <button class="text-xs text-indigo-400 transition-colors hover:text-indigo-300" @click="maSelectAll">Select All</button>
                <span class="text-slate-600">|</span>
                <button class="text-xs text-slate-400 transition-colors hover:text-slate-300" @click="maDeselectAll">Deselect All</button>
              </div>
            </div>

            <!-- Epics -->
            <div
              v-for="epic in epicsWithTasks"
              :key="epic.id"
              class="overflow-hidden rounded-xl border border-slate-700/60 bg-slate-800/30"
            >
              <div
                class="flex cursor-pointer select-none items-center gap-3 px-4 py-3 transition-colors hover:bg-slate-700/30"
                @click="toggleEpicExpand(epic.id)"
              >
                <span class="h-2.5 w-2.5 flex-shrink-0 rounded-full" :style="{ backgroundColor: epic.color || '#6366f1' }"/>
                <input
                  type="checkbox"
                  :checked="isMAEpicFullySelected(epic)"
                  :indeterminate="isMAEpicPartiallySelected(epic)"
                  class="h-4 w-4 flex-shrink-0 rounded accent-indigo-500"
                  @click.stop
                  @change="toggleMAEpicSelection(epic)"
                />
                <span class="flex-1 text-sm font-semibold text-white">{{ epic.title }}</span>
                <span class="rounded-full bg-slate-700 px-2 py-0.5 text-xs text-slate-400">
                  {{ maEpicSelectedCount(epic) }}/{{ epic.tasks.length }}
                </span>
                <svg
                  class="h-4 w-4 text-slate-500 transition-transform"
                  :class="expandedEpics.has(epic.id) ? 'rotate-180' : ''"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </div>
              <div v-show="expandedEpics.has(epic.id)" class="divide-y divide-slate-700/40">
                <label
                  v-for="task in epic.tasks"
                  :key="task.id"
                  class="flex cursor-pointer items-start gap-3 px-5 py-2.5 transition-colors hover:bg-slate-700/20"
                  :class="maSelectedTaskIds.has(task.id) ? 'bg-slate-700/10' : ''"
                >
                  <input
                    v-model="maSelectedTaskIdsArr"
                    :value="task.id"
                    type="checkbox"
                    class="mt-0.5 h-3.5 w-3.5 flex-shrink-0 rounded accent-indigo-500"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="text-sm leading-snug text-slate-200">{{ task.title }}</p>
                    <div class="mt-1 flex flex-wrap items-center gap-2">
                      <span class="text-xs text-slate-500">{{ task.code }}</span>
                      <span class="rounded px-1.5 py-0.5 text-xs font-medium" :class="priorityClass(task.priority)">
                        {{ task.priority }}
                      </span>
                    </div>
                  </div>
                </label>
              </div>
            </div>

            <!-- Tasks without epic -->
            <div
              v-if="tasksWithoutEpic.length > 0"
              class="overflow-hidden rounded-xl border border-slate-700/60 bg-slate-800/30"
            >
              <div
                class="flex cursor-pointer select-none items-center gap-3 px-4 py-3 transition-colors hover:bg-slate-700/30"
                @click="toggleEpicExpand('__no_epic__')"
              >
                <span class="h-2.5 w-2.5 flex-shrink-0 rounded-full bg-slate-500" />
                <input
                  type="checkbox"
                  :checked="isMANoEpicFullySelected"
                  :indeterminate="isMANoEpicPartiallySelected"
                  class="h-4 w-4 flex-shrink-0 rounded accent-indigo-500"
                  @click.stop
                  @change="toggleMANoEpicSelection"
                />
                <span class="flex-1 text-sm font-semibold text-slate-400">No Epic</span>
                <span class="rounded-full bg-slate-700 px-2 py-0.5 text-xs text-slate-400">
                  {{ maNoEpicSelectedCount }}/{{ tasksWithoutEpic.length }}
                </span>
                <svg
                  class="h-4 w-4 text-slate-500 transition-transform"
                  :class="expandedEpics.has('__no_epic__') ? 'rotate-180' : ''"
                  fill="none" viewBox="0 0 24 24" stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
                </svg>
              </div>
              <div v-show="expandedEpics.has('__no_epic__')" class="divide-y divide-slate-700/40">
                <label
                  v-for="task in tasksWithoutEpic"
                  :key="task.id"
                  class="flex cursor-pointer items-start gap-3 px-5 py-2.5 transition-colors hover:bg-slate-700/20"
                  :class="maSelectedTaskIds.has(task.id) ? 'bg-slate-700/10' : ''"
                >
                  <input
                    v-model="maSelectedTaskIdsArr"
                    :value="task.id"
                    type="checkbox"
                    class="mt-0.5 h-3.5 w-3.5 flex-shrink-0 rounded accent-indigo-500"
                  />
                  <div class="min-w-0 flex-1">
                    <p class="text-sm leading-snug text-slate-200">{{ task.title }}</p>
                    <div class="mt-1 flex flex-wrap items-center gap-2">
                      <span class="text-xs text-slate-500">{{ task.code }}</span>
                      <span class="rounded px-1.5 py-0.5 text-xs font-medium" :class="priorityClass(task.priority)">
                        {{ task.priority }}
                      </span>
                    </div>
                  </div>
                </label>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="flex items-center justify-between border-t border-slate-700 px-6 py-4">
            <p class="text-xs text-slate-500">
              {{ maSelectedTaskIds.size === 0 ? 'เลือก task ที่ต้องการให้อยู่ใน MA scope' : `${maSelectedTaskIds.size} task(s) ใน MA scope` }}
            </p>
            <button
              class="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-indigo-500"
              @click="showMATaskModal = false"
            >
              Done
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { useCostingStore } from '../store/costing-store'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { usePricingApi } from '../infrastructure/pricing-api'
import type { QuotationRequest, CustomerLinePricingMode } from '../infrastructure/pricing-api'
import type { Epic, Task } from '~/core/modules/projects/infrastructure/projects-api'
import { isEngineerLikeRole } from '~/utils/roles'

const props = defineProps<{
  projectId: string
  projectName?: string
}>()

const { } = useAuth()
const store = useCostingStore()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()
const pricingApi = usePricingApi()

const exportingCustomer = ref(false)
const exportingMA = ref(false)
const customerLinePricingMode = ref<CustomerLinePricingMode>('absorbed')

const maForm = reactive({
  maPrice: 0,
  maDurationYears: 1,
})

const maResult = ref<{
  maPrice: number
  maDurationYears: number
  monthlyFee: number
} | null>(null)

// ── MA Task Scope & Delivery Milestones ──────────────────────────────────────

const showMATaskModal = ref(false)
const maSelectedTaskIdsArr = ref<string[]>([])
const maSelectedTaskIds = computed(() => new Set(maSelectedTaskIdsArr.value))

function isMAEpicFullySelected(epic: EpicWithTasks): boolean {
  return epic.tasks.length > 0 && epic.tasks.every(t => maSelectedTaskIds.value.has(t.id))
}

function isMAEpicPartiallySelected(epic: EpicWithTasks): boolean {
  const count = epic.tasks.filter(t => maSelectedTaskIds.value.has(t.id)).length
  return count > 0 && count < epic.tasks.length
}

function maEpicSelectedCount(epic: EpicWithTasks): number {
  return epic.tasks.filter(t => maSelectedTaskIds.value.has(t.id)).length
}

const isMANoEpicFullySelected = computed(() =>
  tasksWithoutEpic.value.length > 0 && tasksWithoutEpic.value.every(t => maSelectedTaskIds.value.has(t.id))
)

const isMANoEpicPartiallySelected = computed(() => {
  const count = tasksWithoutEpic.value.filter(t => maSelectedTaskIds.value.has(t.id)).length
  return count > 0 && count < tasksWithoutEpic.value.length
})

const maNoEpicSelectedCount = computed(() =>
  tasksWithoutEpic.value.filter(t => maSelectedTaskIds.value.has(t.id)).length
)

function toggleMAEpicSelection(epic: EpicWithTasks) {
  if (isMAEpicFullySelected(epic)) {
    maSelectedTaskIdsArr.value = maSelectedTaskIdsArr.value.filter(id => !epic.tasks.some(t => t.id === id))
  } else {
    const existing = new Set(maSelectedTaskIdsArr.value)
    epic.tasks.forEach(t => existing.add(t.id))
    maSelectedTaskIdsArr.value = [...existing]
  }
}

function toggleMANoEpicSelection() {
  if (isMANoEpicFullySelected.value) {
    maSelectedTaskIdsArr.value = maSelectedTaskIdsArr.value.filter(id => !tasksWithoutEpic.value.some(t => t.id === id))
  } else {
    const existing = new Set(maSelectedTaskIdsArr.value)
    tasksWithoutEpic.value.forEach(t => existing.add(t.id))
    maSelectedTaskIdsArr.value = [...existing]
  }
}

function maSelectAll() {
  const allIds: string[] = []
  epicsWithTasks.value.forEach(e => e.tasks.forEach(t => allIds.push(t.id)))
  tasksWithoutEpic.value.forEach(t => allIds.push(t.id))
  maSelectedTaskIdsArr.value = allIds
}

function maDeselectAll() {
  maSelectedTaskIdsArr.value = []
}

async function openMATaskModal() {
  showMATaskModal.value = true
  if (epicsWithTasks.value.length === 0 && tasksWithoutEpic.value.length === 0) {
    await loadEpicsAndTasks()
    maSelectAll()
  } else if (maSelectedTaskIdsArr.value.length === 0) {
    maSelectAll()
  }
}

interface DeliveryMilestone {
  id: string
  label: string
  taskCount: number | null  // null = end-of-MA milestone (no task delivery)
  amount: number
  isEndOfMA: boolean
}

const maDeliveryMilestones = ref<DeliveryMilestone[]>([
  { id: 'm1', label: 'งวดที่ 1', taskCount: 0, amount: 0, isEndOfMA: false },
  { id: 'm2', label: 'งวดที่ 2', taskCount: 0, amount: 0, isEndOfMA: false },
  { id: 'm3', label: 'งวดสุดท้าย (สิ้นสุด MA)', taskCount: null, amount: 0, isEndOfMA: true },
])

const milestoneTotalAmount = computed(() =>
  maDeliveryMilestones.value.reduce((sum, m) => sum + (Number(m.amount) || 0), 0)
)

const milestoneAmountDiff = computed(() => {
  if (!maResult.value) return 0
  return maResult.value.maPrice - milestoneTotalAmount.value
})

function addMilestone() {
  const n = maDeliveryMilestones.value.length + 1
  const newMilestone: DeliveryMilestone = { id: `m${Date.now()}`, label: `งวดที่ ${n}`, taskCount: 0, amount: 0, isEndOfMA: false }
  const lastIdx = maDeliveryMilestones.value.findIndex(m => m.isEndOfMA)
  if (lastIdx >= 0) {
    maDeliveryMilestones.value.splice(lastIdx, 0, newMilestone)
  } else {
    maDeliveryMilestones.value.push(newMilestone)
  }
}

function removeMilestone(id: string) {
  maDeliveryMilestones.value = maDeliveryMilestones.value.filter(m => m.id !== id)
}

// ── Form state ──────────────────────────────────────────────────────────────

const form = reactive<QuotationRequest>({
  dev_user_ids: [],
  risk_margin_pct: 0.10,
  profit_margin_pct: 0.25,
})

// ── Cost Config (from admin/cost-config) ─────────────────────────────────────

const loadingConfig = ref(false)
const configError = ref<string | null>(null)
const pmCount = ref(0)
const devCount = ref(0)

// Cost model metrics: from config + company manday-rate API (backend computes overhead from config + DB)
const overheadMultiplier = ref(1.3)
const workingDaysPerMonth = ref(22)
const workingHoursPerDay = ref(8)
const totalDevSalary = ref(0)
const totalDevSS = ref(0)
const costPerManday = ref(0)
const costPerHour = ref(0)
const billableDays = ref(0)

const utilizationRate = computed(() =>
  overheadMultiplier.value > 0 ? 1 / overheadMultiplier.value : 0
)

async function loadCostConfig() {
  loadingConfig.value = true
  configError.value = null
  try {
    const [config, salaries, mandayRate] = await Promise.all([
      pricingApi.getCostConfig(),
      pricingApi.listSalaries(),
      pricingApi.getCompanyMandayRate(),
    ])

    form.risk_margin_pct = config.default_risk_buffer ?? 0.10
    form.profit_margin_pct = config.default_profit_margin ?? 0.25

    overheadMultiplier.value = config.overhead_multiplier ?? 1.3
    workingDaysPerMonth.value = config.working_days_per_month ?? 22
    workingHoursPerDay.value = config.working_hours_per_day ?? 8

    costPerManday.value = mandayRate.cost_per_manday ?? 0
    costPerHour.value = mandayRate.cost_per_hour ?? 0
    billableDays.value = mandayRate.billable_days ?? 0

    // Product Owner count for display only (overhead is from backend)
    const activePMSalaries = salaries.filter(s => (s.user_role === 'PRODUCT_OWNER' || s.user_role === 'PM') && !s.effective_to)
    pmCount.value = activePMSalaries.length

    // DEV user IDs & salary totals: all active DEV-role salary records
    const activeDevSalaries = salaries.filter(s => isEngineerLikeRole(s.user_role) && !s.effective_to)
    devCount.value = activeDevSalaries.length
    form.dev_user_ids = activeDevSalaries.map(s => s.user_id)
    totalDevSalary.value = activeDevSalaries.reduce((sum, s) => sum + s.monthly_salary, 0)
    totalDevSS.value = activeDevSalaries.reduce((sum, s) => sum + (s.ss_cost ?? 0), 0)
  } catch (err: unknown) {
    configError.value = err instanceof Error ? err.message : 'Failed to load cost configuration.'
  } finally {
    loadingConfig.value = false
  }
}

onMounted(async () => {
  await loadCostConfig()
  store.reset()
})

onUnmounted(() => {
  store.reset()
})

// ── Task Selection Modal ──────────────────────────────────────────────────────

const showTaskModal = ref(false)
const loadingEpics = ref(false)

interface EpicWithTasks extends Epic {
  tasks: Task[]
}

const epicsWithTasks = ref<EpicWithTasks[]>([])
const tasksWithoutEpic = ref<Task[]>([])
const selectedTaskIdsArr = ref<string[]>([])
const expandedEpics = ref<Set<string>>(new Set())

const selectedTaskIds = computed(() => new Set(selectedTaskIdsArr.value))

const totalTaskCount = computed(() => {
  return epicsWithTasks.value.reduce((sum, e) => sum + e.tasks.length, 0) + tasksWithoutEpic.value.length
})

function isEpicFullySelected(epic: EpicWithTasks): boolean {
  return epic.tasks.length > 0 && epic.tasks.every(t => selectedTaskIds.value.has(t.id))
}

function isEpicPartiallySelected(epic: EpicWithTasks): boolean {
  const count = epic.tasks.filter(t => selectedTaskIds.value.has(t.id)).length
  return count > 0 && count < epic.tasks.length
}

function epicSelectedCount(epic: EpicWithTasks): number {
  return epic.tasks.filter(t => selectedTaskIds.value.has(t.id)).length
}

const isNoEpicFullySelected = computed(() =>
  tasksWithoutEpic.value.length > 0 && tasksWithoutEpic.value.every(t => selectedTaskIds.value.has(t.id))
)

const isNoEpicPartiallySelected = computed(() => {
  const count = tasksWithoutEpic.value.filter(t => selectedTaskIds.value.has(t.id)).length
  return count > 0 && count < tasksWithoutEpic.value.length
})

const noEpicSelectedCount = computed(() =>
  tasksWithoutEpic.value.filter(t => selectedTaskIds.value.has(t.id)).length
)

function toggleEpicExpand(epicId: string) {
  if (expandedEpics.value.has(epicId)) {
    expandedEpics.value.delete(epicId)
  } else {
    expandedEpics.value.add(epicId)
  }
}

function toggleEpicSelection(epic: EpicWithTasks) {
  const allSelected = isEpicFullySelected(epic)
  if (allSelected) {
    selectedTaskIdsArr.value = selectedTaskIdsArr.value.filter(id => !epic.tasks.some(t => t.id === id))
  } else {
    const ids = epic.tasks.map(t => t.id)
    const existing = new Set(selectedTaskIdsArr.value)
    ids.forEach(id => existing.add(id))
    selectedTaskIdsArr.value = [...existing]
  }
}

function toggleNoEpicSelection() {
  if (isNoEpicFullySelected.value) {
    selectedTaskIdsArr.value = selectedTaskIdsArr.value.filter(id => !tasksWithoutEpic.value.some(t => t.id === id))
  } else {
    const ids = tasksWithoutEpic.value.map(t => t.id)
    const existing = new Set(selectedTaskIdsArr.value)
    ids.forEach(id => existing.add(id))
    selectedTaskIdsArr.value = [...existing]
  }
}

function selectAll() {
  const allIds: string[] = []
  epicsWithTasks.value.forEach(e => e.tasks.forEach(t => allIds.push(t.id)))
  tasksWithoutEpic.value.forEach(t => allIds.push(t.id))
  selectedTaskIdsArr.value = allIds
}

function deselectAll() {
  selectedTaskIdsArr.value = []
}

async function openTaskSelectionModal() {
  if (form.dev_user_ids.length === 0) {
    store.error = 'ไม่พบ Engineer ใน Cost Config / No engineers found in cost configuration. Please add ENGINEER or CHIEF_ENGINEER salary records at Admin → Cost Config.'
    return
  }
  store.error = null
  showTaskModal.value = true

  if (epicsWithTasks.value.length === 0 && tasksWithoutEpic.value.length === 0) {
    await loadEpicsAndTasks()
  }
}

async function loadEpicsAndTasks() {
  loadingEpics.value = true
  try {
    const project = await projectsApi.getProject(props.projectId)
    const resolvedProjectId = project.id

    const [epics, allTasks] = await Promise.all([
      projectsApi.getEpics(resolvedProjectId),
      tasksApi.getTasksByProject(resolvedProjectId),
    ])

    const epicMap = new Map(epics.map(e => [e.id, { ...e, tasks: [] as Task[] }]))

    const withoutEpic: Task[] = []
    for (const task of allTasks) {
      if (task.parent_id) continue // skip subtasks
      if (task.epic_id && epicMap.has(task.epic_id)) {
        epicMap.get(task.epic_id)!.tasks.push(task)
      } else {
        withoutEpic.push(task)
      }
    }

    epicsWithTasks.value = [...epicMap.values()].filter(e => e.tasks.length > 0)
    tasksWithoutEpic.value = withoutEpic

    // Expand all epics by default
    expandedEpics.value = new Set([
      ...epicsWithTasks.value.map(e => e.id),
      ...(withoutEpic.length > 0 ? ['__no_epic__'] : []),
    ])

    // Pre-select all tasks
    selectAll()
  } catch {
    // silently fail, user sees empty state
  } finally {
    loadingEpics.value = false
  }
}

async function confirmAndCalculate() {
  showTaskModal.value = false
  await store.calculateQuotation(props.projectId, {
    ...form,
    task_ids: [...selectedTaskIds.value],
  })
}

// ── Actions ──────────────────────────────────────────────────────────────────

async function exportPDF() {
  const { token, apiBase } = useAuth()
  if (!token.value) { store.error = 'Not authenticated'; return }

  // Open blank tab immediately inside the user gesture — before any await
  // so browser popup blocker does not block it (same pattern as timeline export)
  const tab = window.open('', '_blank')

  store.exporting = true
  store.error = null
  try {
    const url = `${apiBase.value}/sentinel/projects/${props.projectId}/quotation/export`
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({ ...form, task_ids: [...selectedTaskIds.value] }),
      signal: AbortSignal.timeout(120_000),
    })

    if (!response.ok) {
      tab?.close()
      let msg = `PDF generation failed (${response.status})`
      try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
      throw new Error(msg)
    }

    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)

    if (tab) {
      tab.location.href = objectUrl
    } else {
      // Fallback if popup was blocked
      const link = document.createElement('a')
      link.href = objectUrl
      link.target = '_blank'
      link.rel = 'noopener noreferrer'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
  } catch (e: any) {
    tab?.close()
    store.error = e?.message ?? 'PDF export failed'
  } finally {
    store.exporting = false
  }
}

async function exportCustomerPDF() {
  const { token, apiBase } = useAuth()
  if (!token.value) { store.error = 'Not authenticated'; return }

  const tab = window.open('', '_blank')

  exportingCustomer.value = true
  store.error = null
  try {
    const url = `${apiBase.value}/sentinel/projects/${props.projectId}/quotation/export`
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({
        ...form,
        task_ids: [...selectedTaskIds.value],
        customer_view: true,
        project_name: props.projectName ?? '',
        customer_line_pricing_mode: customerLinePricingMode.value,
      }),
      signal: AbortSignal.timeout(120_000),
    })

    if (!response.ok) {
      tab?.close()
      let msg = `PDF generation failed (${response.status})`
      try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
      throw new Error(msg)
    }

    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)

    if (tab) {
      tab.location.href = objectUrl
    } else {
      const link = document.createElement('a')
      link.href = objectUrl
      link.target = '_blank'
      link.rel = 'noopener noreferrer'
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
    }
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
  } catch (e: any) {
    tab?.close()
    store.error = e?.message ?? 'Customer PDF export failed'
  } finally {
    exportingCustomer.value = false
  }
}

async function exportMAPDF() {
  if (!maResult.value) {
    store.error = 'กรุณากด Generate MA Quotation ก่อน export PDF'
    return
  }

  const { token, apiBase } = useAuth()
  if (!token.value) { store.error = 'Not authenticated'; return }

  // Open blank tab immediately inside the user gesture to avoid popup blocker
  const previewTab = window.open('', '_blank')
  if (!previewTab) {
    store.error = 'Popup blocked. Please allow popups and try again.'
    return
  }

  exportingMA.value = true
  store.error = null

  try {
    const now = new Date()
    const buddhistYear = now.getFullYear() + 543
    const quoteNo = `MA-${String(buddhistYear)}${String(now.getMonth() + 1).padStart(2, '0')}${String(now.getDate()).padStart(2, '0')}-${props.projectId.slice(0, 6).toUpperCase()}`
    const issueDate = `${String(now.getDate()).padStart(2, '0')}/${String(now.getMonth() + 1).padStart(2, '0')}/${buddhistYear}`

    // Build MA task list (preserve epic grouping order)
    const tasks: Array<{ code: string; title: string }> = []
    for (const epic of epicsWithTasks.value) {
      for (const task of epic.tasks) {
        if (maSelectedTaskIds.value.has(task.id)) {
          tasks.push({ code: task.code, title: task.title })
        }
      }
    }
    for (const task of tasksWithoutEpic.value) {
      if (maSelectedTaskIds.value.has(task.id)) {
        tasks.push({ code: task.code, title: task.title })
      }
    }

    const milestones = maDeliveryMilestones.value.map(m => ({
      label: m.label,
      task_count: m.isEndOfMA ? null : (Number(m.taskCount) || 0),
      amount: Number(m.amount) || 0,
      is_end_of_ma: m.isEndOfMA,
    }))

    const url = `${apiBase.value}/sentinel/projects/${props.projectId}/ma-quotation/export`
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({
        project_name: props.projectName ?? 'Enterprise Project',
        quote_no: quoteNo,
        issue_date: issueDate,
        ma_price: maResult.value.maPrice,
        ma_duration_years: maResult.value.maDurationYears,
        tasks,
        milestones,
      }),
      signal: AbortSignal.timeout(120_000),
    })

    if (!response.ok) {
      previewTab.close()
      let msg = `MA PDF generation failed (${response.status})`
      try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
      throw new Error(msg)
    }

    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)
    previewTab.location.href = objectUrl
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
  } catch (e: any) {
    previewTab.close()
    store.error = e?.message ?? 'MA PDF export failed'
  } finally {
    exportingMA.value = false
  }
}

// exportMAPdfBlob removed — MA PDF is now generated server-side via HTML+chromedp.
// Keeping the stub signature commented out to avoid import errors during transition.
function _exportMAPdfBlobUnused(payload: {
  quoteNo: string
  issueDate: string
  projectName: string
  maPrice: number
  maDurationYears: number
  monthlyFee: number
  tasks: Array<{ code: string; title: string; epicTitle?: string; epicColor?: string }>
  milestones: Array<{ label: string; taskCount: number | null; amount: number; isEndOfMA: boolean }>
  thaiFont?: { regular: string; bold: string }
}): Blob {
  const doc = new jsPDF({ orientation: 'p', unit: 'mm', format: 'a4' })

  // ── Register Sarabun font if available ────────────────────────────────────
  const hasThai = !!payload.thaiFont?.regular
  if (hasThai && payload.thaiFont) {
    doc.addFileToVFS('Sarabun-Regular.ttf', payload.thaiFont.regular)
    doc.addFont('Sarabun-Regular.ttf', 'Sarabun', 'normal')
    doc.addFileToVFS('Sarabun-Bold.ttf', payload.thaiFont.bold)
    doc.addFont('Sarabun-Bold.ttf', 'Sarabun', 'bold')
  }
  // F / FB are the active font name — Sarabun when loaded, helvetica as fallback
  const F = hasThai ? 'Sarabun' : 'helvetica'
  const FB = hasThai ? 'Sarabun' : 'helvetica'

  const MARGIN = 14
  const PAGE_W = 210
  const CONTENT_W = PAGE_W - MARGIN * 2
  const FOOTER_Y = 278
  const ROW_H = 8

  const companyName = 'Komgrip Technologies'
  const companyTagline = 'Software Engineering & Digital Solutions'

  let currentPage = 1

  function drawFooter() {
    doc.setDrawColor(199, 216, 237)
    doc.setLineWidth(0.3)
    doc.line(MARGIN, FOOTER_Y, PAGE_W - MARGIN, FOOTER_Y)
    doc.setTextColor(156, 163, 175)
    doc.setFont(F, 'normal')
    doc.setFontSize(7.5)
    doc.text(`Page ${currentPage}`, PAGE_W - MARGIN, FOOTER_Y + 4, { align: 'right' })
  }

  function checkBreak(y: number, needed = ROW_H): number {
    if (y + needed > FOOTER_Y - 8) {
      drawFooter()
      doc.addPage()
      currentPage++
      doc.setDrawColor(30, 58, 95)
      doc.setLineWidth(0.5)
      doc.line(MARGIN, 12, PAGE_W - MARGIN, 12)
      doc.setTextColor(30, 58, 95)
      doc.setFont(FB, 'bold')
      doc.setFontSize(8)
      doc.text(companyName.toUpperCase(), MARGIN, 10)
      doc.setFont(F, 'normal')
      doc.setTextColor(75, 85, 99)
      doc.text(`MA Quotation · ${payload.quoteNo} (cont.)`, PAGE_W - MARGIN, 10, { align: 'right' })
      return 20
    }
    return y
  }

  function sectionTitle(y: number, title: string): number {
    doc.setDrawColor(199, 216, 237)
    doc.setLineWidth(0.3)
    doc.line(MARGIN, y + 1, PAGE_W - MARGIN, y + 1)
    doc.setTextColor(30, 58, 95)
    doc.setFont(FB, 'bold')
    doc.setFontSize(10)
    doc.text(title, MARGIN, y)
    return y + 6
  }

  // ── HEADER ─────────────────────────────────────────────────────────────────
  doc.setDrawColor(30, 58, 95)
  doc.setLineWidth(1.2)
  doc.line(MARGIN, 14, PAGE_W - MARGIN, 14)

  doc.setTextColor(30, 58, 95)
  doc.setFont(FB, 'bold')
  doc.setFontSize(18)
  doc.text(companyName.toUpperCase(), MARGIN, 23)

  doc.setFont(F, 'normal')
  doc.setFontSize(8)
  doc.setTextColor(74, 111, 165)
  doc.text(companyTagline, MARGIN, 28)

  doc.setTextColor(75, 85, 99)
  doc.setFontSize(9)
  doc.setFont(FB, 'bold')
  doc.text('Maintenance Agreement Quotation', PAGE_W - MARGIN, 22, { align: 'right' })
  doc.setFont(F, 'normal')
  doc.text(`Date: ${payload.issueDate}`, PAGE_W - MARGIN, 27, { align: 'right' })
  doc.text(`Project: ${payload.projectName}`, PAGE_W - MARGIN, 32, { align: 'right' })
  doc.text(`Quote No: ${payload.quoteNo}`, PAGE_W - MARGIN, 37, { align: 'right' })

  let y = 46

  // ── SECTION 1: MA OVERVIEW ─────────────────────────────────────────────────
  y = sectionTitle(y, 'MAINTENANCE AGREEMENT — OVERVIEW')

  const boxW = (CONTENT_W - 6) / 3
  const boxes = [
    { label: 'Annual MA Price', value: formatTHB(payload.maPrice), sub: 'incl. all services' },
    { label: 'MA Duration', value: `${payload.maDurationYears} ปี (Year${payload.maDurationYears > 1 ? 's' : ''})`, sub: 'ระยะเวลาสัญญา' },
    { label: 'Monthly Fee', value: formatTHB(payload.monthlyFee), sub: 'Annual ÷ 12' },
  ]

  boxes.forEach((box, i) => {
    const bx = MARGIN + i * (boxW + 3)
    const isPrice = i === 0

    doc.setFillColor(isPrice ? 30 : 245, isPrice ? 58 : 249, isPrice ? 95 : 255)
    doc.roundedRect(bx, y, boxW, 22, 2, 2, 'F')

    doc.setTextColor(isPrice ? 200 : 100, isPrice ? 210 : 116, isPrice ? 255 : 139)
    doc.setFont(F, 'normal')
    doc.setFontSize(6.5)
    doc.text(box.label.toUpperCase(), bx + 4, y + 5.5)

    doc.setFont(FB, 'bold')
    doc.setFontSize(isPrice ? 10 : 9)
    doc.setTextColor(isPrice ? 255 : 30, isPrice ? 255 : 58, isPrice ? 255 : 95)
    doc.text(box.value, bx + 4, y + 14)

    doc.setFont(F, 'normal')
    doc.setFontSize(7)
    doc.setTextColor(isPrice ? 180 : 120, isPrice ? 200 : 130, isPrice ? 240 : 150)
    doc.text(box.sub, bx + 4, y + 19.5)
  })

  y += 28

  // ── SECTION 2: TASK SCOPE ──────────────────────────────────────────────────
  if (payload.tasks.length > 0) {
    y = checkBreak(y, 22)
    y = sectionTitle(y, `SCOPE OF WORK — TASK LIST  (${payload.tasks.length} tasks)`)

    const xNum = MARGIN + 3
    const xCode = MARGIN + 11
    const xTitle = MARGIN + 32

    doc.setFillColor(30, 58, 95)
    doc.rect(MARGIN, y, CONTENT_W, ROW_H, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFont(FB, 'bold')
    doc.setFontSize(7.5)
    doc.text('#', xNum, y + 5.3)
    doc.text('Code', xCode, y + 5.3)
    doc.text('Task Title', xTitle, y + 5.3)
    y += ROW_H

    payload.tasks.forEach((task, idx) => {
      y = checkBreak(y)
      if (idx % 2 === 0) {
        doc.setFillColor(248, 250, 255)
        doc.rect(MARGIN, y, CONTENT_W, ROW_H, 'F')
      }
      doc.setDrawColor(230, 237, 245)
      doc.line(MARGIN, y + ROW_H, PAGE_W - MARGIN, y + ROW_H)

      doc.setFont(F, 'normal')
      doc.setFontSize(8)

      doc.setTextColor(100, 116, 139)
      doc.text(String(idx + 1), xNum, y + 5.3)

      doc.setTextColor(30, 58, 95)
      doc.setFont(FB, 'bold')
      doc.text(task.code ?? '', xCode, y + 5.3)

      const maxLen = 80
      const titleText = task.title.length > maxLen ? task.title.slice(0, maxLen - 1) + '…' : task.title
      doc.setFont(F, 'normal')
      doc.setTextColor(26, 32, 53)
      doc.text(titleText, xTitle, y + 5.3)

      y += ROW_H
    })

    y += 10
  }

  // ── SECTION 3: DELIVERY SCHEDULE ──────────────────────────────────────────
  if (payload.milestones.length > 0) {
    y = checkBreak(y, payload.milestones.length * ROW_H + 40)
    y = sectionTitle(y, 'DELIVERY SCHEDULE & PAYMENT MILESTONES')

    const xMilestone = MARGIN + 4
    const xTasks = MARGIN + 70
    const xPayment = PAGE_W - MARGIN - 3

    doc.setFillColor(30, 58, 95)
    doc.rect(MARGIN, y, CONTENT_W, ROW_H, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFont(FB, 'bold')
    doc.setFontSize(8)
    doc.text('งวดส่งมอบ', xMilestone, y + 5.3)
    doc.text('การส่งมอบงาน', xTasks, y + 5.3)
    doc.text('จำนวนเงิน (THB)', xPayment, y + 5.3, { align: 'right' })
    y += ROW_H

    payload.milestones.forEach((m, idx) => {
      const rowBg = m.isEndOfMA ? [232, 248, 255] : idx % 2 === 0 ? [245, 249, 255] : [255, 255, 255]
      doc.setFillColor(rowBg[0], rowBg[1], rowBg[2])
      doc.rect(MARGIN, y, CONTENT_W, ROW_H, 'F')
      doc.setDrawColor(199, 216, 237)
      doc.line(MARGIN, y + ROW_H, PAGE_W - MARGIN, y + ROW_H)

      doc.setTextColor(m.isEndOfMA ? 0 : 26, m.isEndOfMA ? 100 : 32, m.isEndOfMA ? 130 : 53)
      doc.setFont(m.isEndOfMA ? FB : F, m.isEndOfMA ? 'bold' : 'normal')
      doc.setFontSize(9)
      doc.text(m.label, xMilestone, y + 5.3)

      if (m.isEndOfMA) {
        doc.setTextColor(0, 100, 140)
        doc.setFont(F, 'normal')
        doc.setFontSize(8.5)
        doc.text('เบิกเมื่อสิ้นสุดระยะเวลา MA', xTasks, y + 5.3)
      } else {
        doc.setTextColor(74, 111, 165)
        doc.setFont(FB, 'bold')
        doc.setFontSize(9)
        doc.text(`ส่งมอบ ${m.taskCount ?? 0} tasks`, xTasks, y + 5.3)
      }

      doc.setTextColor(30, 58, 95)
      doc.setFont(FB, 'bold')
      doc.setFontSize(9)
      doc.text(m.amount > 0 ? formatTHB(m.amount) : '—', xPayment, y + 5.3, { align: 'right' })

      y += ROW_H
    })

    const totalAmount = payload.milestones.reduce((s, m) => s + (m.amount || 0), 0)
    const totalRowH = 12
    doc.setFillColor(30, 58, 95)
    doc.rect(MARGIN, y, CONTENT_W, totalRowH, 'F')
    doc.setTextColor(255, 255, 255)
    doc.setFont(FB, 'bold')
    doc.setFontSize(9)
    doc.text('รวมทั้งหมด (TOTAL)', xMilestone, y + 4.8)
    doc.setFont(F, 'normal')
    doc.setFontSize(7.5)
    doc.setTextColor(210, 220, 235)
    doc.text('*ไม่รวม VAT', xMilestone, y + 9.2)
    doc.setTextColor(255, 255, 255)
    doc.setFont(FB, 'bold')
    doc.setFontSize(9)
    doc.text(formatTHB(totalAmount), xPayment, y + 7.2, { align: 'right' })
    y += totalRowH

    y += 10
  }

  // ── VALIDITY NOTE ──────────────────────────────────────────────────────────
  y = checkBreak(y, 22)
  doc.setDrawColor(199, 216, 237)
  doc.setFillColor(240, 245, 251)
  doc.roundedRect(MARGIN, y, CONTENT_W, 16, 1.5, 1.5, 'FD')
  doc.setTextColor(75, 85, 99)
  doc.setFont(F, 'normal')
  doc.setFontSize(8)
  doc.text('ใบเสนอราคาฉบับนี้มีผลภายใน 30 วัน นับจากวันที่ออกเอกสาร', MARGIN + 4, y + 5.5)
  doc.text('ราคาเป็นสกุลเงินบาท (THB) · เงื่อนไขตามสัญญาที่ลงนามร่วมกัน', MARGIN + 4, y + 11)

  drawFooter()

  return doc.output('blob')
}

// ── Formatting ───────────────────────────────────────────────────────────────

function formatTHB(val: number): string {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    currencyDisplay: 'code',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(val)
}

function formatNumber(val: number): string {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(val)
}

function getTaskDisplayAmount(task: { cost: number }, idx: number): number {
  if (!store.result) return task.cost
  if (customerLinePricingMode.value === 'base') return task.cost

  const tasks = store.result.tasks
  if (!tasks.length) return task.cost

  const totalBeforeVAT = store.result.subtotal + store.result.risk_amount + store.result.profit_amount

  // Same allocation logic as backend customer PDF (absorbed mode)
  if (store.result.subtotal > 0) {
    if (idx === tasks.length - 1) {
      const prior = tasks.slice(0, idx).reduce((sum, t) => sum + round2((t.cost / store.result!.subtotal) * totalBeforeVAT), 0)
      return round2(totalBeforeVAT - prior)
    }
    return round2((task.cost / store.result.subtotal) * totalBeforeVAT)
  }

  const equalShare = round2(totalBeforeVAT / tasks.length)
  if (idx === tasks.length - 1) {
    return round2(totalBeforeVAT - equalShare * (tasks.length - 1))
  }
  return equalShare
}

function formatDate(d: string | null): string {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB', { day: '2-digit', month: 'short' })
}

function round2(v: number): number {
  return Math.round(v * 100) / 100
}

function calculateMAQuotation() {
  const maPrice = Number(maForm.maPrice) || 0
  const maDurationYears = Math.max(1, Number(maForm.maDurationYears) || 1)

  if (maPrice <= 0) {
    store.error = 'กรุณาระบุราคา MA ให้มากกว่า 0'
    maResult.value = null
    return
  }

  store.error = null
  maResult.value = {
    maPrice,
    maDurationYears,
    monthlyFee: maPrice / 12,
  }
}

function formatPercent(val: number): string {
  return `${new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 1,
    maximumFractionDigits: 2,
  }).format(val)}%`
}

function priorityClass(priority: string): string {
  const map: Record<string, string> = {
    CRITICAL: 'bg-red-500/15 text-red-400',
    HIGH: 'bg-orange-500/15 text-orange-400',
    MEDIUM: 'bg-yellow-500/15 text-yellow-400',
    LOW: 'bg-gray-500/15 text-gray-400',
  }
  return map[priority] ?? 'bg-gray-500/15 text-gray-400'
}
</script>

<style scoped>
.input-field {
  @apply w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2 text-sm text-white
         placeholder-gray-500 transition-colors
         focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500/50;
}

.btn-primary {
  @apply inline-flex items-center rounded-lg bg-gradient-to-r from-purple-600 to-pink-600
         px-5 py-2.5 text-sm font-semibold text-white shadow-lg transition-all
         hover:from-purple-500 hover:to-pink-500 hover:shadow-purple-500/25
         disabled:cursor-not-allowed disabled:opacity-50;
}

.btn-secondary {
  @apply inline-flex items-center rounded-lg border border-amber-500/40 bg-amber-500/10
         px-5 py-2.5 text-sm font-semibold text-amber-400 transition-all
         hover:bg-amber-500/20 hover:border-amber-500/60
         disabled:cursor-not-allowed disabled:opacity-50;
}

.metric-card {
  @apply rounded-xl border border-gray-700 bg-gray-800/60 px-4 py-3;
}
.metric-label {
  @apply text-xs font-medium uppercase tracking-wide text-gray-500;
}
.metric-value {
  @apply mt-1 text-xl font-bold text-white;
}

.summary-row {
  @apply flex items-center justify-between px-5 py-3;
}
.summary-label {
  @apply text-sm text-gray-400;
}
.summary-amount {
  @apply text-sm font-semibold text-gray-200 tabular-nums;
}

.premium-kpi-card {
  @apply rounded-xl border border-slate-700 bg-slate-900/60 px-4 py-3;
}

.premium-kpi-label {
  @apply text-[11px] font-semibold uppercase tracking-widest text-slate-400;
}

.premium-kpi-value {
  @apply mt-1 text-lg font-black text-white tabular-nums;
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
