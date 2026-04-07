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
          <h3 class="text-sm font-semibold text-white">Task Cost Breakdown</h3>
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
                  {{ formatNumber(task.cost) }}
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
              class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-700 hover:text-white transition-colors"
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
</template>

<script setup lang="ts">
import { useCostingStore } from '../store/costing-store'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { usePricingApi } from '../infrastructure/pricing-api'
import type { QuotationRequest } from '../infrastructure/pricing-api'
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
    const [epics, allTasks] = await Promise.all([
      projectsApi.getEpics(props.projectId),
      tasksApi.getTasksByProject(props.projectId),
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

function formatDate(d: string | null): string {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB', { day: '2-digit', month: 'short' })
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

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
