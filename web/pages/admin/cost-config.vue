<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Page Header ──────────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-4">
          <NuxtLink to="/dashboard" class="text-gray-500 hover:text-gray-300 transition-colors text-sm">
            ← Dashboard
          </NuxtLink>
          <span class="text-gray-700">/</span>
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center">
              <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div>
              <h1 class="text-base font-bold text-white">Cost Configuration</h1>
              <p class="text-xs text-gray-500">Salary registry & company overhead settings</p>
            </div>
          </div>
        </div>

        <!-- Tab nav -->
        <div class="flex items-center gap-1">
          <button
            v-for="tab in pageTabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="px-4 py-2 text-xs font-semibold rounded-lg transition-colors"
            :class="activeTab === tab.id
              ? 'bg-amber-500/15 text-amber-400 border border-amber-500/30'
              : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800'"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Access denied -->
    <div v-if="!isCEO" class="flex items-center justify-center min-h-[60vh]">
      <div class="text-center max-w-sm">
        <div class="text-5xl mb-4">🔒</div>
        <h2 class="text-xl font-bold text-white mb-2">Access Restricted</h2>
        <p class="text-gray-400 text-sm">Only CEO can access cost configuration.</p>
      </div>
    </div>

    <div v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">

      <!-- ── TAB: Dashboard ─────────────────────────────────────────────── -->
      <div v-if="activeTab === 'dashboard'">

        <!-- Loading state -->
        <div v-if="loadingSalaries || loadingConfig || loadingMandayRate" class="flex items-center justify-center py-32">
          <div class="flex flex-col items-center gap-3">
            <svg class="h-8 w-8 animate-spin text-amber-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <p class="text-sm text-gray-500">Computing cost model…</p>
          </div>
        </div>

        <template v-else>
          <!-- No data warning -->
          <div v-if="devCount === 0" class="rounded-xl border border-dashed border-amber-500/30 bg-amber-500/5 p-6 mb-8 flex items-start gap-4">
            <svg class="h-5 w-5 text-amber-400 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
            </svg>
            <div>
              <p class="text-sm font-semibold text-amber-400">No engineering salaries found</p>
              <p class="text-xs text-gray-500 mt-0.5">Add salary records for users with role ENGINEER or CHIEF ENGINEER to see cost calculations. Product Owner & overhead data can still be reviewed.</p>
            </div>
          </div>

          <!-- ── Hero metrics row ── -->
          <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
            <div class="col-span-2 lg:col-span-1 rounded-2xl border border-purple-500/25 bg-gradient-to-br from-purple-900/40 to-pink-900/30 p-5">
              <p class="text-xs font-semibold uppercase tracking-widest text-purple-300 mb-1">Cost / Manday</p>
              <p class="text-3xl font-black text-white tabular-nums">{{ formatTHB(costPerManday) }}</p>
              <p class="text-xs text-gray-500 mt-1.5">Fully loaded ÷ billable days</p>
            </div>
            <div class="rounded-2xl border border-cyan-500/20 bg-gray-800/60 p-5">
              <p class="text-xs font-semibold uppercase tracking-widest text-cyan-300 mb-1">Cost / Hour</p>
              <p class="text-2xl font-extrabold text-white tabular-nums">{{ formatTHB(costPerHour) }}</p>
              <p class="text-xs text-gray-500 mt-1.5">Manday ÷ {{ configForm.working_hours_per_day }}h</p>
            </div>
            <div class="rounded-2xl border border-amber-500/20 bg-gray-800/60 p-5">
              <p class="text-xs font-semibold uppercase tracking-widest text-amber-300 mb-1">Billable Days</p>
              <p class="text-2xl font-extrabold text-amber-400 tabular-nums">{{ billableDays.toFixed(1) }}</p>
              <p class="text-xs text-gray-500 mt-1.5">of {{ configForm.working_days_per_month }} work days/mo</p>
            </div>
            <div class="rounded-2xl border border-emerald-500/20 bg-gray-800/60 p-5">
              <p class="text-xs font-semibold uppercase tracking-widest text-emerald-300 mb-1">Utilisation</p>
              <p class="text-2xl font-extrabold text-emerald-400 tabular-nums">{{ (utilizationRate * 100).toFixed(0) }}%</p>
              <p class="text-xs text-gray-500 mt-1.5">1 ÷ {{ configForm.overhead_multiplier }}×</p>
            </div>
          </div>

          <!-- ── Cost breakdown ── -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">

            <!-- Formula card -->
            <div class="lg:col-span-2 rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden">
              <div class="px-5 py-4 border-b border-gray-700 flex items-center gap-3">
                <div class="w-7 h-7 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center flex-shrink-0">
                  <svg class="w-3.5 h-3.5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 11h.01M12 11h.01M15 11h.01M4 7h16a1 1 0 010 2H4a1 1 0 010-2zm0 4h16a1 1 0 010 2H4a1 1 0 010-2z"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-sm font-bold text-white">Fully Loaded Cost Breakdown</h3>
                  <p class="text-xs text-gray-500">Per developer · per month</p>
                </div>
              </div>

              <div class="divide-y divide-gray-700/60">
                <!-- Row: overhead -->
                <div class="px-5 py-4 flex items-start justify-between gap-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs px-1.5 py-0.5 rounded bg-cyan-900/40 border border-cyan-500/30 text-cyan-400 font-mono">1</span>
                      <p class="text-sm font-semibold text-white">Overhead per Dev</p>
                      <span v-if="mandayRate" class="text-xs px-1.5 py-0.5 rounded bg-emerald-900/30 border border-emerald-500/30 text-emerald-400">live API</span>
                    </div>
                    <p class="text-xs text-gray-500 font-mono">
                      (Executive + Company + Product Owner + Manager + Support + Total SS) ÷ {{ devCount }} devs
                    </p>
                    <div class="flex flex-wrap gap-2 mt-1.5">
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-cyan-500"></span>
                        Executive ฿{{ fmtN(mandayRate?.executive_expense ?? configForm.executive_expense ?? 0) }}
                      </span>
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-blue-500"></span>
                        Company ฿{{ fmtN(mandayRate?.company_expense ?? configForm.company_expense ?? 0) }}
                      </span>
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-purple-500"></span>
                        {{ pmCount }} PO ฿{{ fmtN(totalPMSalary) }}
                      </span>
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
                        {{ managerCount }} Manager ฿{{ fmtN(totalManagerSalary) }}
                      </span>
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-cyan-400"></span>
                        {{ supportCount }} Support ฿{{ fmtN(totalSupportSalary) }}
                      </span>
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                        Total SS ฿{{ fmtN(mandayRate?.total_monthly_ss ?? totalSS) }}
                      </span>
                    </div>
                    <p v-if="mandayRate" class="text-xs text-gray-600 mt-1 font-mono">
                      Company Expense Total: ฿{{ fmtN(mandayRate.company_expense_total) }} ÷ {{ devCount }} devs
                    </p>
                  </div>
                  <div class="text-right flex-shrink-0">
                    <p class="text-lg font-extrabold text-cyan-400 tabular-nums">{{ formatTHB(overheadPerDev) }}</p>
                    <p class="text-xs text-gray-500">/mo per dev</p>
                  </div>
                </div>

                <!-- Row: dev cost -->
                <div class="px-5 py-4 flex items-start justify-between gap-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs px-1.5 py-0.5 rounded bg-amber-900/40 border border-amber-500/30 text-amber-400 font-mono">2</span>
                      <p class="text-sm font-semibold text-white">Cost per Dev</p>
                    </div>
                    <p class="text-xs text-gray-500 font-mono">
                      Dev Salaries ฿{{ fmtN(totalDevSalary) }} ÷ {{ devCount }} devs (SS is in overhead)
                    </p>
                    <div class="flex flex-wrap gap-2 mt-1.5">
                      <span class="inline-flex items-center gap-1 text-xs text-gray-400">
                        <span class="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
                        Avg ฿{{ fmtN(devCount > 0 ? totalDevSalary / devCount : 0) }} /dev
                      </span>
                    </div>
                  </div>
                  <div class="text-right flex-shrink-0">
                    <p class="text-lg font-extrabold text-amber-400 tabular-nums">{{ formatTHB(costPerDev) }}</p>
                    <p class="text-xs text-gray-500">/mo per dev</p>
                  </div>
                </div>

                <!-- Row: fully loaded -->
                <div class="px-5 py-4 flex items-start justify-between gap-4 bg-gray-700/20">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs px-1.5 py-0.5 rounded bg-purple-900/40 border border-purple-500/30 text-purple-400 font-mono">Σ</span>
                      <p class="text-sm font-bold text-white">Fully Loaded Cost</p>
                    </div>
                    <p class="text-xs text-gray-500 font-mono">Overhead + Dev Cost</p>
                  </div>
                  <div class="text-right flex-shrink-0">
                    <p class="text-xl font-black text-purple-400 tabular-nums">{{ formatTHB(fullyLoadedCost) }}</p>
                    <p class="text-xs text-gray-500">/mo per dev</p>
                  </div>
                </div>

                <!-- Row: billable days -->
                <div class="px-5 py-4 flex items-start justify-between gap-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs px-1.5 py-0.5 rounded bg-amber-900/30 border border-amber-600/30 text-amber-300 font-mono">÷</span>
                      <p class="text-sm font-semibold text-white">Billable Days</p>
                    </div>
                    <p class="text-xs text-gray-500 font-mono">
                      {{ configForm.working_days_per_month }} work days ÷ {{ configForm.overhead_multiplier }}× multiplier
                    </p>
                  </div>
                  <div class="text-right flex-shrink-0">
                    <p class="text-lg font-extrabold text-amber-400 tabular-nums">{{ billableDays.toFixed(2) }} days</p>
                    <p class="text-xs text-gray-500">{{ (utilizationRate * 100).toFixed(1) }}% utilisation</p>
                  </div>
                </div>

                <!-- Row: final result -->
                <div class="px-5 py-4 flex items-start justify-between gap-4 bg-gradient-to-r from-purple-900/30 to-pink-900/20">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-0.5">
                      <span class="text-xs px-1.5 py-0.5 rounded bg-gradient-to-r from-purple-600 to-pink-600 text-white font-bold font-mono">=</span>
                      <p class="text-sm font-bold text-white">Cost per Manday</p>
                    </div>
                    <p class="text-xs text-gray-500 font-mono">Fully Loaded ÷ Billable Days</p>
                  </div>
                  <div class="text-right flex-shrink-0">
                    <p class="text-2xl font-black text-white tabular-nums">{{ formatTHB(costPerManday) }}</p>
                    <p class="text-xs text-gray-500">per manday ({{ configForm.working_hours_per_day }}h)</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Side panel: headcount & totals -->
            <div class="space-y-4">
              <!-- Headcount -->
              <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-5">
                <h3 class="text-sm font-bold text-white mb-3 flex items-center gap-2">
                  <svg class="h-4 w-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
                  </svg>
                  Headcount
                </h3>
                <div class="space-y-2">
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-yellow-500/10 text-yellow-400 border border-yellow-500/30">CEO</span>
                    <span class="text-sm font-semibold text-white">{{ ceoSalaries.length }} person{{ ceoSalaries.length !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-orange-500/10 text-orange-400 border border-orange-500/30">MANAGER</span>
                    <span class="text-sm font-semibold text-white">{{ managerCount }} person{{ managerCount !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/30">PO</span>
                    <span class="text-sm font-semibold text-white">{{ pmCount }} person{{ pmCount !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-purple-500/10 text-purple-400 border border-purple-500/30">ENGINEER</span>
                    <span class="text-sm font-semibold text-white">{{ engineerSalaryCount }} person{{ engineerSalaryCount !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-violet-500/10 text-violet-300 border border-violet-500/30">CHIEF ENGINEER</span>
                    <span class="text-sm font-semibold text-white">{{ chiefEngineerSalaryCount }} person{{ chiefEngineerSalaryCount !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="flex justify-between items-center">
                    <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/30">SUPPORT</span>
                    <span class="text-sm font-semibold text-white">{{ supportCount }} person{{ supportCount !== 1 ? 's' : '' }}</span>
                  </div>
                  <div class="border-t border-gray-700 pt-2 mt-2 flex justify-between items-center">
                    <span class="text-xs text-gray-500">Total (with salary)</span>
                    <span class="text-sm font-bold text-white">{{ mandayRate?.active_headcount ?? activeSalaries.length }}</span>
                  </div>
                </div>
              </div>

              <!-- Monthly totals -->
              <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-5">
                <h3 class="text-sm font-bold text-white mb-3 flex items-center gap-2">
                  <svg class="h-4 w-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  Monthly Cost Pool
                  <span v-if="mandayRate" class="text-xs px-1.5 py-0.5 rounded bg-emerald-900/30 border border-emerald-500/30 text-emerald-400 ml-auto">live API</span>
                </h3>
                <div class="space-y-2 text-sm">
                  <div class="flex justify-between">
                    <span class="text-gray-400">Product Owner salaries</span>
                    <span class="text-white tabular-nums">{{ formatTHB(totalPMSalary) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-400">Manager Salaries</span>
                    <span class="text-white tabular-nums">{{ formatTHB(totalManagerSalary) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-400">Support Salaries</span>
                    <span class="text-white tabular-nums">{{ formatTHB(totalSupportSalary) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-400">Engineering salaries</span>
                    <span class="text-white tabular-nums">{{ formatTHB(totalDevSalary) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-emerald-400/90">SS (all staff)</span>
                    <span class="text-emerald-400/90 tabular-nums">{{ formatTHB(mandayRate?.total_monthly_ss ?? totalSS) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-cyan-400/90">Executive Exp.</span>
                    <span class="text-cyan-400/90 tabular-nums">{{ formatTHB(mandayRate?.executive_expense ?? configForm.executive_expense ?? 0) }}</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-blue-400/90">Company Exp.</span>
                    <span class="text-blue-400/90 tabular-nums">{{ formatTHB(mandayRate?.company_expense ?? configForm.company_expense ?? 0) }}</span>
                  </div>
                  <div class="border-t border-gray-700 pt-2 mt-2 flex justify-between font-bold">
                    <span class="text-gray-300">Grand Total</span>
                    <span class="text-purple-400 tabular-nums">{{ formatTHB(mandayRate?.total_monthly_burn_rate ?? (totalPayroll + totalSS + (configForm.executive_expense ?? 0) + (configForm.company_expense ?? 0))) }}</span>
                  </div>
                </div>
              </div>

              <!-- Quick links -->
              <div class="rounded-2xl border border-gray-700/60 bg-gray-800/30 p-4 flex gap-2">
                <button class="flex-1 text-xs py-2 rounded-lg text-gray-400 border border-gray-700 hover:bg-gray-700/50 hover:text-gray-900 dark:text-white transition-colors" @click="activeTab = 'salaries'">
                  Edit Salaries →
                </button>
                <button class="flex-1 text-xs py-2 rounded-lg text-gray-400 border border-gray-700 hover:bg-gray-700/50 hover:text-gray-900 dark:text-white transition-colors" @click="activeTab = 'config'">
                  Edit Config →
                </button>
              </div>
            </div>
          </div>

          <!-- ── Formula reference ── -->
          <div class="rounded-2xl border border-gray-700/60 bg-gray-800/40 p-5">
            <h3 class="text-sm font-bold text-white mb-4 flex items-center gap-2">
              <svg class="h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
              </svg>
              Formula Reference
            </h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-cyan-400 mb-0.5">Overhead per Dev</p>
                <p class="text-xs text-gray-400 font-mono">(Exec + Company + PO + Mgr + Support + Total SS) ÷ #Devs</p>
              </div>
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-amber-400 mb-0.5">Cost per Dev</p>
                <p class="text-xs text-gray-400 font-mono">Dev Salaries ÷ #Devs (SS in overhead)</p>
              </div>
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-purple-400 mb-0.5">Fully Loaded Cost</p>
                <p class="text-xs text-gray-400 font-mono">Overhead + Dev Cost</p>
              </div>
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-emerald-400 mb-0.5">SS Cost</p>
                <p class="text-xs text-gray-400 font-mono">Min(Salary × 5%, ฿875)</p>
              </div>
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-amber-300 mb-0.5">Billable Days</p>
                <p class="text-xs text-gray-400 font-mono">WorkDays ÷ Overhead Multiplier</p>
              </div>
              <div class="rounded-lg bg-gray-900/60 border border-gray-700/50 px-3 py-2">
                <p class="text-xs font-bold text-white mb-0.5">Cost per Manday</p>
                <p class="text-xs text-gray-400 font-mono">Fully Loaded ÷ Billable Days</p>
              </div>
            </div>
          </div>

          <!-- ── VC Finance Summary (hidden when teams feature disabled) ── -->
          <div v-if="teamsStore.teamsFeatureEnabled" class="rounded-2xl border border-yellow-700/30 bg-gradient-to-br from-yellow-900/15 to-amber-900/10 overflow-hidden">
            <div class="px-5 py-4 border-b border-yellow-700/20 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="w-7 h-7 rounded-lg bg-yellow-500/15 border border-yellow-500/30 flex items-center justify-center flex-shrink-0">
                  <svg class="w-3.5 h-3.5 text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div>
                  <h3 class="text-sm font-bold text-white">Internal VC — Team Capital Overview</h3>
                  <p class="text-xs text-gray-500">ยอด capital, burn rate, runway ของแต่ละทีม</p>
                </div>
              </div>
              <NuxtLink
                to="/admin/teams"
                class="text-xs text-yellow-400 hover:text-yellow-300 transition-colors border border-yellow-500/30 rounded-lg px-3 py-1.5"
              >
                Manage Teams →
              </NuxtLink>
            </div>

            <!-- No teams -->
            <div v-if="teamsStore.teams.length === 0" class="px-5 py-8 text-center text-gray-600 text-sm">
              No teams found. Create teams first in the Teams Management page.
            </div>

            <!-- Team rows -->
            <div v-else class="divide-y divide-yellow-700/10">
              <div
                v-for="team in teamsStore.teams"
                :key="team.id"
                class="px-5 py-4 flex flex-wrap items-center gap-4"
              >
                <!-- Team name -->
                <div class="min-w-[140px] flex-shrink-0">
                  <p class="text-sm font-semibold text-white">{{ team.name }}</p>
                  <p class="text-xs text-gray-500">{{ team.users?.length ?? 0 }} members</p>
                </div>

                <!-- Loading -->
                <div v-if="teamsStore.costLoading[team.id]" class="flex-1 text-xs text-gray-600">Computing…</div>

                <template v-else-if="teamsStore.teamCosts[team.id]">
                  <!-- Burn -->
                  <div class="text-center min-w-[110px]">
                    <p class="text-xs text-gray-500 mb-0.5">Burn / Month</p>
                    <p class="text-sm font-bold text-orange-400 tabular-nums">{{ formatTHB(teamsStore.teamCosts[team.id]!.total_monthly_cost) }}</p>
                  </div>

                  <!-- Capital -->
                  <div class="text-center min-w-[110px]">
                    <p class="text-xs text-gray-500 mb-0.5">Capital</p>
                    <p class="text-sm font-bold text-emerald-400 tabular-nums">{{ formatTHB(teamsStore.teamCosts[team.id]!.capital_balance) }}</p>
                  </div>

                  <!-- Bonus % -->
                  <div class="text-center min-w-[80px]">
                    <p class="text-xs text-gray-500 mb-0.5">Bonus Target</p>
                    <p class="text-sm font-bold text-yellow-400">{{ teamsStore.teamCosts[team.id]!.bonus_percentage }}%</p>
                  </div>

                  <!-- Runway bar -->
                  <div class="flex-1 min-w-[160px]">
                    <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
                      <span>Runway</span>
                      <span :class="vcRunwayColor(teamsStore.teamCosts[team.id]!.runway_months)">
                        {{ teamsStore.teamCosts[team.id]!.runway_months > 0 ? teamsStore.teamCosts[team.id]!.runway_months.toFixed(1) + ' mo' : '—' }}
                      </span>
                    </div>
                    <div class="h-1.5 bg-gray-700 rounded-full overflow-hidden">
                      <div
                        :class="['h-full rounded-full transition-all', vcRunwayBarColor(teamsStore.teamCosts[team.id]!.runway_months)]"
                        :style="{ width: vcRunwayBarWidth(teamsStore.teamCosts[team.id]!.runway_months) }"
                      />
                    </div>
                  </div>
                </template>

                <!-- Not loaded -->
                <div v-else class="flex-1 text-xs text-gray-600">
                  <button @click="teamsStore.fetchTeamCost(team.id)" class="text-gray-500 hover:text-gray-300 transition-colors">Load →</button>
                </div>
              </div>
            </div>
          </div>

          <!-- ── Export Cost Analysis Report ── -->
          <div class="rounded-2xl border border-blue-500/25 bg-gradient-to-br from-blue-900/30 to-indigo-900/20 p-6">
            <div class="flex items-start justify-between gap-6 flex-wrap">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-3 mb-2">
                  <div class="w-9 h-9 rounded-xl bg-blue-500/15 border border-blue-500/30 flex items-center justify-center flex-shrink-0">
                    <svg class="w-4.5 h-4.5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
                    </svg>
                  </div>
                  <div>
                    <h3 class="text-sm font-bold text-white">Export Cost Analysis Report</h3>
                    <p class="text-xs text-gray-400">รายงานวิเคราะห์ต้นทุนรายเดือน/รายปีสำหรับผู้บริหาร — 5 sections · PDF · Confidential</p>
                  </div>
                </div>
                <p class="text-xs text-gray-500 leading-relaxed">
                  รายงานโครงสร้างต้นทุนพื้นฐานบริษัท: Monthly Burn Rate · Annual Projection · Fully Loaded Cost/Dev · Headcount Salary Analysis · 12-Month Cash Flow · Pricing Sensitivity Matrix
                </p>
              </div>

              <!-- Controls -->
              <div class="flex flex-col gap-3 min-w-[200px]">
                <!-- Error message -->
                <p v-if="reportExportError" class="text-xs text-red-400">{{ reportExportError }}</p>

                <!-- Export button -->
                <button
                  @click="handleExportCostReport"
                  :disabled="reportExporting"
                  class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl font-bold text-sm transition-all"
                  :class="reportExporting
                    ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-500/20'"
                >
                  <svg v-if="reportExporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  {{ reportExporting ? 'Generating Report…' : 'Preview Report' }}
                </button>
              </div>
            </div>
          </div>

          <!-- ── PDF Preview Modal ── -->
          <Teleport to="body">
            <Transition name="fade">
              <div
                v-if="reportPreviewUrl"
                class="fixed inset-0 z-50 flex flex-col bg-gray-950/95 backdrop-blur-sm"
                @keydown.esc="closeReportPreview"
                tabindex="-1"
              >
                <!-- Modal toolbar -->
                <div class="flex items-center justify-between px-5 py-3 bg-gray-900 border-b border-gray-800 flex-shrink-0">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-lg bg-blue-500/15 border border-blue-500/30 flex items-center justify-center">
                      <svg class="w-4 h-4 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414A1 1 0 0119 9.414V19a2 2 0 01-2 2z"/>
                      </svg>
                    </div>
                    <div>
                      <p class="text-sm font-bold text-white">Cost Analysis Report Preview</p>
                      <p class="text-xs text-gray-500">รายงานวิเคราะห์ต้นทุน · CONFIDENTIAL</p>
                    </div>
                  </div>
                  <div class="flex items-center gap-2">
                    <!-- Download button -->
                    <button
                      @click="downloadReport"
                      class="flex items-center gap-2 px-4 py-2 rounded-lg bg-blue-100 dark:bg-blue-600 hover:bg-blue-100 dark:bg-blue-500 text-gray-900 dark:text-white text-sm font-bold transition-colors shadow-lg shadow-blue-500/20"
                    >
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
                      </svg>
                      Download PDF
                    </button>
                    <!-- Close button -->
                    <button
                      @click="closeReportPreview"
                      class="flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-800 hover:bg-gray-700 text-gray-300 hover:text-gray-900 dark:text-white text-sm font-semibold transition-colors border border-gray-700"
                    >
                      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                      </svg>
                      Close
                    </button>
                  </div>
                </div>

                <!-- PDF iframe -->
                <div class="flex-1 overflow-hidden bg-gray-800">
                  <iframe
                    :src="reportPreviewUrl"
                    class="w-full h-full border-0"
                    type="application/pdf"
                  />
                </div>
              </div>
            </Transition>
          </Teleport>

        </template>
      </div>

      <!-- ── TAB: Salaries ───────────────────────────────────────────────── -->
      <div v-if="activeTab === 'salaries'">

        <!-- Summary metrics -->
        <div class="grid grid-cols-2 sm:grid-cols-5 gap-4 mb-8">
          <div class="metric-card">
            <p class="metric-label">Total Staff</p>
            <p class="metric-value">{{ uniqueUsers }}</p>
          </div>
          <div class="metric-card">
            <p class="metric-label">Salary Records</p>
            <p class="metric-value">{{ salaries.length }}</p>
          </div>
          <div class="metric-card">
            <p class="metric-label">Avg Salary (base)</p>
            <p class="metric-value text-amber-400">{{ formatTHB(avgSalary) }}</p>
          </div>
          <div class="metric-card">
            <p class="metric-label">Total SS/mo</p>
            <p class="metric-value text-emerald-400">{{ formatTHB(totalSS) }}</p>
            <p class="text-xs text-gray-500 mt-0.5">Min(5%, ฿875) per person</p>
          </div>
          <div class="metric-card">
            <p class="metric-label">Total Monthly Payroll</p>
            <p class="metric-value text-purple-400">{{ formatTHB(totalPayrollPlusSS) }}</p>
            <p class="text-xs text-gray-500 mt-0.5">Base + SS</p>
          </div>
        </div>

        <!-- Header row -->
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-lg font-bold text-white">Employee Salary Registry</h2>
            <p class="text-sm text-gray-500 mt-0.5">Set effective-dated monthly salaries for each team member</p>
          </div>
          <button class="btn-primary" @click="openAddSalaryModal()">
            <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            Add Salary Record
          </button>
        </div>

        <!-- Error -->
        <div v-if="salaryError" class="flex items-start gap-3 rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400 mb-4">
          <span class="flex-1">{{ salaryError }}</span>
          <button @click="salaryError = ''" class="text-red-400 hover:text-red-300">✕</button>
        </div>

        <!-- Loading -->
        <div v-if="loadingSalaries" class="flex items-center justify-center py-24">
          <div class="flex flex-col items-center gap-3">
            <svg class="h-8 w-8 animate-spin text-amber-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
            </svg>
            <p class="text-sm text-gray-500">Loading salary data…</p>
          </div>
        </div>

        <!-- Salary table (grouped by role) -->
        <template v-else>
          <div v-for="group in groupedSalaries" :key="group.role" class="mb-6">
            <!-- Group header -->
            <div class="flex items-center gap-3 mb-3">
              <span :class="roleChip(group.role)">{{ displayCostRoleLabel(group.role) }}</span>
              <span class="text-xs text-gray-500">{{ group.members.length }} member{{ group.members.length !== 1 ? 's' : '' }}</span>
              <div class="flex-1 h-px bg-gray-800"/>
              <span class="text-xs text-gray-500">
                Total: {{ formatTHB(group.members.reduce((s, m) => s + m.monthly_salary, 0)) }} base
                + {{ formatTHB(group.members.reduce((s, m) => s + (m.ss_cost ?? 0), 0)) }} SS/mo
              </span>
            </div>

            <!-- Table -->
            <div class="rounded-xl border border-gray-700 bg-gray-800/40 overflow-hidden">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700 text-xs uppercase tracking-wider text-gray-500 bg-gray-900/40">
                    <th class="px-4 py-3 text-left">Employee</th>
                    <th class="px-4 py-3 text-left">Role</th>
                    <th class="px-4 py-3 text-left">Type</th>
                    <th class="px-4 py-3 text-right">Monthly Salary</th>
                    <th class="px-4 py-3 text-right">SS (5%, max ฿875)</th>
                    <th class="px-4 py-3 text-left">Effective From</th>
                    <th class="px-4 py-3 text-left">Effective To</th>
                    <th class="px-4 py-3 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(sal, idx) in group.members"
                    :key="sal.id"
                    class="border-b border-gray-700/50 transition-colors hover:bg-gray-700/20"
                    :class="idx % 2 === 1 ? 'bg-gray-900/20' : ''"
                  >
                    <td class="px-4 py-3">
                      <div class="flex items-center gap-3">
                        <div class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-600 to-pink-600 flex items-center justify-center text-xs font-bold text-white flex-shrink-0">
                          {{ initials(sal.user_display_name || sal.user_email) }}
                        </div>
                        <div>
                          <p class="font-medium text-white">{{ sal.user_display_name || sal.user_email }}</p>
                          <p class="text-xs text-gray-500">{{ sal.user_email }}</p>
                        </div>
                      </div>
                    </td>
                    <td class="px-4 py-3">
                      <span class="text-xs font-medium text-gray-300">{{ displayCostRoleLabel(sal.user_role) }}</span>
                    </td>
                    <td class="px-4 py-3">
                      <span :class="empTypeChip(sal.employment_type)">{{ sal.employment_type }}</span>
                    </td>
                    <td class="px-4 py-3 text-right font-semibold text-amber-400 tabular-nums">
                      {{ formatTHB(sal.monthly_salary) }}
                    </td>
                    <td class="px-4 py-3 text-right text-emerald-400/90 tabular-nums text-xs" title="Min(เงินเดือน × 5%, 875)">
                      {{ formatTHB(sal.ss_cost ?? 0) }}
                    </td>
                    <td class="px-4 py-3 text-gray-300 text-xs">{{ sal.effective_from?.substring(0, 10) }}</td>
                    <td class="px-4 py-3 text-gray-500 text-xs">{{ sal.effective_to?.substring(0, 10) || '—' }}</td>
                    <td class="px-4 py-3 text-right">
                      <div class="flex items-center justify-end gap-2">
                        <button
                          class="p-1.5 rounded-lg text-gray-400 hover:text-gray-900 dark:text-white hover:bg-gray-700 transition-colors"
                          title="Edit"
                          @click="openEditSalaryModal(sal)"
                        >
                          <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                          </svg>
                        </button>
                        <button
                          :disabled="deletingId === sal.id"
                          class="p-1.5 rounded-lg text-gray-400 hover:text-red-400 hover:bg-red-900/20 transition-colors disabled:opacity-40"
                          title="Delete"
                          @click="confirmDelete(sal)"
                        >
                          <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div v-if="salaries.length === 0" class="rounded-xl border border-dashed border-gray-700 py-20 text-center">
            <svg class="mx-auto h-12 w-12 text-gray-600 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
            <p class="text-gray-400 font-medium mb-1">No salary records yet</p>
            <p class="text-gray-600 text-sm mb-4">Add salary records to enable cost calculations</p>
            <button class="btn-primary" @click="openAddSalaryModal()">Add First Record</button>
          </div>
        </template>
      </div>

      <!-- ── TAB: Company Config ─────────────────────────────────────────── -->
      <div v-if="activeTab === 'config'">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-lg font-bold text-white">Company Cost Configuration</h2>
            <p class="text-sm text-gray-500 mt-0.5">Working days, overhead multiplier and default margins</p>
          </div>
        </div>

        <div v-if="configError" class="flex items-start gap-3 rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400 mb-4">
          <span class="flex-1">{{ configError }}</span>
          <button @click="configError = ''" class="text-red-400 hover:text-red-300">✕</button>
        </div>

        <div v-if="loadingConfig" class="flex items-center justify-center py-24">
          <svg class="h-8 w-8 animate-spin text-amber-400" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
          </svg>
        </div>

        <form v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6" @submit.prevent="saveConfig">
          <!-- Workdays -->
          <div class="config-card">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <h3>Working Schedule</h3>
            </div>
            <p class="config-card-desc">Calendar basis for billable day calculations</p>
            <div class="grid grid-cols-2 gap-4 mt-4">
              <div>
                <label class="field-label">Working Days / Month</label>
                <input v-model.number="configForm.working_days_per_month" type="number" min="1" max="31" class="input-field" />
                <p class="field-hint">Typically 20–23 days</p>
              </div>
              <div>
                <label class="field-label">Working Hours / Day</label>
                <input v-model.number="configForm.working_hours_per_day" type="number" min="1" max="24" class="input-field" />
                <p class="field-hint">Standard is 8 hours</p>
              </div>
            </div>
          </div>

          <!-- Overhead multiplier -->
          <div class="config-card">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-2 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              <h3>Overhead Multiplier</h3>
            </div>
            <p class="config-card-desc">Fraction of working days that are billable (1.30 = 30% overhead → ~76.9% billable)</p>
            <div class="mt-4">
              <label class="field-label">
                Multiplier
                <span class="ml-2 text-amber-400 font-semibold">{{ configForm.overhead_multiplier }}×</span>
                <span class="ml-2 text-gray-500 text-xs">({{ billableDaysPreview }} billable days/mo)</span>
              </label>
              <input
                v-model.number="configForm.overhead_multiplier"
                type="range"
                min="1.0"
                max="2.0"
                step="0.05"
                class="mt-2 w-full accent-amber-500"
              />
              <div class="flex justify-between text-xs text-gray-500 mt-1">
                <span>1.00× (100% billable)</span>
                <span>2.00× (50% billable)</span>
              </div>
            </div>
          </div>

          <!-- Profit margin -->
          <div class="config-card">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
              </svg>
              <h3>Default Profit Margin</h3>
            </div>
            <p class="config-card-desc">Applied automatically when creating new quotations</p>
            <div class="mt-4">
              <label class="field-label">
                Profit Margin
                <span class="ml-2 text-green-400 font-semibold">{{ pct(configForm.default_profit_margin) }}</span>
              </label>
              <input
                v-model.number="configForm.default_profit_margin"
                type="range"
                min="0"
                max="0.6"
                step="0.01"
                class="mt-2 w-full accent-green-500"
              />
              <div class="flex justify-between text-xs text-gray-500 mt-1">
                <span>0%</span><span>30%</span><span>60%</span>
              </div>
            </div>
          </div>

          <!-- Risk buffer -->
          <div class="config-card">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
              </svg>
              <h3>Default Risk Buffer</h3>
            </div>
            <p class="config-card-desc">Contingency reserve applied to all project quotations</p>
            <div class="mt-4">
              <label class="field-label">
                Risk Buffer
                <span class="ml-2 text-red-400 font-semibold">{{ pct(configForm.default_risk_buffer) }}</span>
              </label>
              <input
                v-model.number="configForm.default_risk_buffer"
                type="range"
                min="0"
                max="0.5"
                step="0.01"
                class="mt-2 w-full accent-red-500"
              />
              <div class="flex justify-between text-xs text-gray-500 mt-1">
                <span>0%</span><span>25%</span><span>50%</span>
              </div>
            </div>
          </div>

          <!-- Executive & Company expense -->
          <div class="config-card col-span-full">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-cyan-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-2 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              <h3>Default Monthly Expenses</h3>
            </div>
            <p class="config-card-desc">Used as defaults when building project quotations (overhead / executive allocation)</p>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
              <div>
                <label class="field-label">Executive expense (THB/mo)</label>
                <input
                  v-model.number="configForm.executive_expense"
                  type="number"
                  min="0"
                  step="1000"
                  class="input-field"
                  placeholder="e.g. 0"
                />
                <p class="field-hint">ค่าใช้จ่ายฝ่ายบริหารต่อเดือน</p>
              </div>
              <div>
                <label class="field-label">Company expense (THB/mo)</label>
                <input
                  v-model.number="configForm.company_expense"
                  type="number"
                  min="0"
                  step="1000"
                  class="input-field"
                  placeholder="e.g. 0"
                />
                <p class="field-hint">ต้นทุนค่าบริหาร/ค่าใช้จ่ายบริษัทต่อเดือน (ค่าเช่า สาธารณูปโภค ฯลฯ)</p>
              </div>
            </div>
          </div>

          <!-- Currency -->
          <div class="config-card col-span-full">
            <div class="config-card-header">
              <svg class="h-5 w-5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16H9m3 0h3"/>
              </svg>
              <h3>Currency</h3>
            </div>
            <div class="mt-4 flex items-center gap-4">
              <div class="w-40">
                <label class="field-label">ISO 4217 Code</label>
                <input v-model="configForm.currency" type="text" maxlength="3" placeholder="THB" class="input-field uppercase" />
              </div>
              <p class="text-sm text-gray-500 mt-4">Used as the display currency across all quotations and reports.</p>
            </div>
          </div>

          <!-- Save -->
          <div class="col-span-full flex items-center justify-end gap-3 pt-2">
            <p v-if="configSaved" class="text-sm text-green-400 flex items-center gap-1.5">
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              Saved successfully
            </p>
            <button type="submit" :disabled="savingConfig" class="btn-primary">
              <svg v-if="savingConfig" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              {{ savingConfig ? 'Saving…' : 'Save Configuration' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ── Add / Edit Salary Modal ─────────────────────────────────────────── -->
    <Teleport to="body">
      <div
        v-if="showSalaryModal"
        class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeSalaryModal"
      >
        <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 w-full max-w-lg shadow-2xl">
          <div class="flex items-center justify-between mb-5">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center">
                <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <h2 class="text-lg font-bold text-white">{{ editingSalary ? 'Edit Salary Record' : 'Add Salary Record' }}</h2>
            </div>
            <button @click="closeSalaryModal" class="text-gray-500 hover:text-white transition-colors p-1">
              <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>

        <div v-if="modalError" class="flex items-start gap-3 rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400 mb-4">
          <span class="flex-1">{{ modalError }}</span>
          <button @click="modalError = ''" class="text-red-400 hover:text-red-300">✕</button>
        </div>

          <form class="space-y-4" @submit.prevent="submitSalaryModal">
            <!-- Employee -->
            <div>
              <label class="field-label">Employee *</label>
              <select v-model.number="salaryForm.user_id" class="input-field" :disabled="!!editingSalary">
                <option value="" disabled>Select employee…</option>
                <option v-for="m in allMembers" :key="m.id" :value="m.id">
                  [{{ displayCostRoleLabel(m.role) }}] {{ m.display_name || m.email }}
                </option>
              </select>
            </div>

            <!-- Monthly Salary -->
            <div>
              <label class="field-label">Monthly Salary (THB) *</label>
              <div class="relative">
                <input
                  v-model.number="salaryForm.monthly_salary"
                  type="number"
                  min="0"
                  step="500"
                  class="input-field"
                  placeholder="e.g. 45000"
                />
              </div>
              <p class="field-hint text-emerald-500/90 mt-1">
                Costing includes SS: Min(5% of salary, ฿875)/mo · ค่าประกันสังคม
              </p>
            </div>

            <!-- Employment Type -->
            <div>
              <label class="field-label">Employment Type</label>
              <div class="flex gap-2">
                <button
                  v-for="t in empTypes"
                  :key="t.value"
                  type="button"
                  @click="salaryForm.employment_type = t.value"
                  class="flex-1 py-2 px-3 rounded-lg border text-xs font-semibold transition-colors"
                  :class="salaryForm.employment_type === t.value
                    ? 'border-amber-500/50 bg-amber-500/10 text-amber-300'
                    : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-gray-500'"
                >
                  {{ t.label }}
                </button>
              </div>
            </div>

            <!-- Dates -->
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="field-label">Effective From *</label>
                <UiDatePicker v-model="salaryForm.effective_from" placeholder="Select date…" />
              </div>
              <div>
                <label class="field-label">Effective To <span class="text-gray-500">(optional)</span></label>
                <UiDatePicker v-model="salaryForm.effective_to" placeholder="Select date…" />
              </div>
            </div>

            <div class="flex gap-3 pt-2">
              <button type="button" class="btn-ghost flex-1" @click="closeSalaryModal">Cancel</button>
              <button type="submit" :disabled="modalSubmitting" class="btn-primary flex-1">
                <svg v-if="modalSubmitting" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                {{ editingSalary ? 'Save Changes' : 'Add Record' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

    <!-- ── Delete Confirm Modal ────────────────────────────────────────────── -->
    <Teleport to="body">
      <div
        v-if="salaryToDelete"
        class="fixed inset-0 bg-black/70 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="salaryToDelete = null"
      >
        <div class="bg-gray-800 border border-red-500/30 rounded-2xl p-6 w-full max-w-sm shadow-2xl">
          <div class="text-center">
            <div class="w-12 h-12 rounded-full bg-red-900/30 border border-red-500/30 flex items-center justify-center mx-auto mb-4">
              <svg class="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-bold text-white mb-1">Delete Salary Record?</h3>
            <p class="text-sm text-gray-400 mb-1">
              <span class="text-white font-medium">{{ salaryToDelete.user_display_name || salaryToDelete.user_email }}</span>
            </p>
            <p class="text-sm text-gray-500">
              {{ formatTHB(salaryToDelete.monthly_salary) }}/mo · From {{ salaryToDelete.effective_from?.substring(0, 10) }}
            </p>
            <p class="text-xs text-red-400 mt-3">This action cannot be undone.</p>
          </div>
          <div class="flex gap-3 mt-5">
            <button class="btn-ghost flex-1" @click="salaryToDelete = null">Cancel</button>
            <button :disabled="deletingId === salaryToDelete.id" class="btn-danger flex-1" @click="doDelete">
              <svg v-if="deletingId === salaryToDelete?.id" class="mr-2 h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              Delete
            </button>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { usePricingApi } from '~/core/modules/pricing/infrastructure/pricing-api'
import type { EmployeeSalaryWithUser, UpdateCostConfigPayload, CostReportRequest, CostReportPeriod, CompanyMandayRateResponse } from '~/core/modules/pricing/infrastructure/pricing-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { isEngineerLikeRole } from '~/utils/roles'

definePageMeta({ layout: 'default', middleware: 'auth' })

const { currentUser, fetchWithAuth, token, apiBase } = useAuth()

const isCEO = computed(() => currentUser.value?.role === 'CEO')

// ── Page tabs ─────────────────────────────────────────────────────────────────

const pageTabs = [
  { id: 'dashboard', label: '📊 Cost Dashboard' },
  { id: 'salaries',  label: '👥 Employee Salaries' },
  { id: 'config',    label: '⚙️ Company Config' },
]
const activeTab = ref('dashboard')

const api = usePricingApi()
const teamsStore = useTeamsStore()

// ── Company Manday Rate (authoritative backend calculation) ───────────────────

const mandayRate = ref<CompanyMandayRateResponse | null>(null)
const loadingMandayRate = ref(false)

async function loadMandayRate() {
  loadingMandayRate.value = true
  try {
    mandayRate.value = await api.getCompanyMandayRate()
  } catch { /* non-critical; fallback to local computation */ }
  finally { loadingMandayRate.value = false }
}

// Computed aliases that prefer the API values, falling back to local formula
const costPerManday   = computed(() => mandayRate.value?.cost_per_manday   ?? localCostPerManday.value)
const costPerHour     = computed(() => mandayRate.value?.cost_per_hour      ?? localCostPerHour.value)
const billableDays    = computed(() => mandayRate.value?.billable_days       ?? localBillableDays.value)
const utilizationRate = computed(() =>
  mandayRate.value
    ? 1 / (mandayRate.value.overhead_multiplier || 1)
    : localUtilizationRate.value
)
// Use API-authoritative per-dev breakdown so all three numbers are consistent:
// fullyLoadedCost ÷ billableDays === costPerManday
const overheadPerDev  = computed(() =>
  mandayRate.value?.overhead_per_dev ?? localOverheadPerDev.value
)
const fullyLoadedCost = computed(() =>
  mandayRate.value?.fully_loaded_monthly_per_dev ?? (overheadPerDev.value + costPerDev.value)
)

// ── Cost Analysis Report Export ───────────────────────────────────────────────

const reportExporting   = ref(false)
const reportExportError = ref('')
const reportPreviewUrl  = ref<string | null>(null)
const reportFilename    = ref('')

async function handleExportCostReport() {
  if (!token.value) { reportExportError.value = 'Not authenticated'; return }
  reportExporting.value = true
  reportExportError.value = ''
  try {
    const url = `${apiBase.value}/pricing/report/export`
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token.value}`,
      },
      body: JSON.stringify({ period: 'all' }),
      signal: AbortSignal.timeout(180_000),
    })

    if (!response.ok) {
      let msg = `Export failed (${response.status})`
      try { const j = await response.json(); msg = j.error || j.message || msg } catch {}
      throw new Error(msg)
    }

    const blob = await response.blob()
    // Revoke previous URL if any
    if (reportPreviewUrl.value) URL.revokeObjectURL(reportPreviewUrl.value)
    reportPreviewUrl.value = URL.createObjectURL(blob)
    reportFilename.value = `cost-analysis-report-${new Date().toISOString().slice(0, 10)}.pdf`
  } catch (e: any) {
    reportExportError.value = e?.message ?? 'Export failed'
  } finally {
    reportExporting.value = false
  }
}

function downloadReport() {
  if (!reportPreviewUrl.value) return
  const link = document.createElement('a')
  link.href = reportPreviewUrl.value
  link.download = reportFilename.value
  link.rel = 'noopener noreferrer'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

function closeReportPreview() {
  if (reportPreviewUrl.value) {
    URL.revokeObjectURL(reportPreviewUrl.value)
    reportPreviewUrl.value = null
  }
}

// ── Team members (for modal employee dropdown) ────────────────────────────────

interface TeamMember {
  id: number
  email: string
  role: string
  display_name?: string
}
const allMembers = ref<TeamMember[]>([])

async function fetchTeam() {
  try {
    const res = await fetchWithAuth<{ data: TeamMember[] }>('/auth/users')
    allMembers.value = res?.data ?? []
  } catch { /* non-critical */ }
}

// ── Salaries ──────────────────────────────────────────────────────────────────

const salaries = ref<EmployeeSalaryWithUser[]>([])
const loadingSalaries = ref(false)
const salaryError = ref('')
const deletingId = ref<number | null>(null)
const salaryToDelete = ref<EmployeeSalaryWithUser | null>(null)

const uniqueUsers = computed(() => new Set(salaries.value.map(s => s.user_id)).size)
const totalPayroll = computed(() => {
  const seen = new Set<number>()
  return salaries.value
    .filter(s => { if (seen.has(s.user_id)) return false; seen.add(s.user_id); return true })
    .reduce((sum, s) => sum + s.monthly_salary, 0)
})
const totalSS = computed(() => {
  const seen = new Set<number>()
  return salaries.value
    .filter(s => { if (seen.has(s.user_id)) return false; seen.add(s.user_id); return true })
    .reduce((sum, s) => sum + (s.ss_cost ?? 0), 0)
})
const totalPayrollPlusSS = computed(() => totalPayroll.value + totalSS.value)
const avgSalary = computed(() => uniqueUsers.value ? totalPayroll.value / uniqueUsers.value : 0)

// ── Dashboard cost model computeds ──────────────────────────────────────────

// Deduplicate salary records per user — keep latest effective record per user_id
const activeSalaries = computed(() => {
  const seen = new Set<number>()
  const out: EmployeeSalaryWithUser[] = []
  for (const s of salaries.value) {
    if (!seen.has(s.user_id)) { seen.add(s.user_id); out.push(s) }
  }
  return out
})
const devSalaries     = computed(() => activeSalaries.value.filter(s => isEngineerLikeRole(s.user_role)))
const pmSalaries      = computed(() => activeSalaries.value.filter(s => s.user_role === 'PRODUCT_OWNER' || s.user_role === 'PM'))
const managerSalaries = computed(() => activeSalaries.value.filter(s => s.user_role === 'MANAGER'))
const supportSalaries = computed(() => activeSalaries.value.filter(s => s.user_role === 'SUPPORT'))
const ceoSalaries     = computed(() => activeSalaries.value.filter(s => s.user_role === 'CEO'))

const devCount     = computed(() => devSalaries.value.length)
/** Headcount rows: Engineer = ENGINEER + legacy DEV; Chief Engineer counted separately. */
const engineerSalaryCount = computed(() =>
  activeSalaries.value.filter((s) => {
    const r = (s.user_role || '').toUpperCase()
    return r === 'ENGINEER' || r === 'DEV'
  }).length
)
const chiefEngineerSalaryCount = computed(() =>
  activeSalaries.value.filter((s) => (s.user_role || '').toUpperCase() === 'CHIEF_ENGINEER').length
)
const pmCount      = computed(() => pmSalaries.value.length)
const managerCount = computed(() => managerSalaries.value.length)
const supportCount = computed(() => supportSalaries.value.length)

// Overhead roles (Product Owner + MANAGER + SUPPORT) — not direct labour; same as backend
const totalPMSalary      = computed(() => pmSalaries.value.reduce((s, m) => s + m.monthly_salary, 0))
const totalManagerSalary = computed(() => managerSalaries.value.reduce((s, m) => s + m.monthly_salary, 0))
const totalSupportSalary = computed(() => supportSalaries.value.reduce((s, m) => s + m.monthly_salary, 0))
const totalOverheadRoleSalary = computed(() => totalPMSalary.value + totalManagerSalary.value + totalSupportSalary.value)

// Total dev salary only (SS is company overhead, not per-dev)
const totalDevSalary = computed(() => devSalaries.value.reduce((s, m) => s + m.monthly_salary, 0))
// totalSS is already computed above: SS for all employees (deduplicated)

// Default Monthly Expenses from config: executive_expense + company_expense
const defaultMonthlyExpenses = computed(() => (configForm.executive_expense ?? 0) + (configForm.company_expense ?? 0))

// Local fallback formulas — used when mandayRate API is unavailable
const localOverheadPerDev = computed(() => {
  if (devCount.value === 0) return 0
  return (defaultMonthlyExpenses.value + totalOverheadRoleSalary.value + totalSS.value) / devCount.value
})

const costPerDev = computed(() => {
  if (devCount.value === 0) return 0
  return totalDevSalary.value / devCount.value
})

const localUtilizationRate = computed(() => configForm.overhead_multiplier > 0 ? 1 / configForm.overhead_multiplier : 0)
const localBillableDays    = computed(() =>
  configForm.overhead_multiplier > 0
    ? configForm.working_days_per_month / configForm.overhead_multiplier
    : 0
)
const localCostPerManday = computed(() => localBillableDays.value > 0 ? (localOverheadPerDev.value + costPerDev.value) / localBillableDays.value : 0)
const localCostPerHour   = computed(() => configForm.working_hours_per_day > 0 ? localCostPerManday.value / configForm.working_hours_per_day : 0)

/** Salary table grouping: legacy DEV → Engineer; legacy PM → Product Owner. */
function canonicalCostRoleGroup(userRole: string): string {
  const u = (userRole || '').trim().toUpperCase()
  if (u === 'DEV' || u === 'ENGINEER') return 'ENGINEER'
  if (u === 'PM') return 'PRODUCT_OWNER'
  return u || 'OTHER'
}

/** User-visible role title (ALL CAPS; never show raw DEV in UI). */
function displayCostRoleLabel(role: string): string {
  const r = (role || '').trim().toUpperCase()
  switch (r) {
    case 'ENGINEER':
    case 'DEV':
      return 'ENGINEER'
    case 'CHIEF_ENGINEER':
      return 'CHIEF ENGINEER'
    case 'PRODUCT_OWNER':
    case 'PM':
      return 'PRODUCT OWNER'
    case 'CEO':
      return 'CEO'
    case 'MANAGER':
      return 'MANAGER'
    case 'SUPPORT':
      return 'SUPPORT'
    default:
      return r || '—'
  }
}

const groupedSalaries = computed(() => {
  const groups: Record<string, { role: string; members: EmployeeSalaryWithUser[] }> = {}
  const order = ['CEO', 'MANAGER', 'PRODUCT_OWNER', 'ENGINEER', 'CHIEF_ENGINEER', 'SUPPORT']
  for (const s of salaries.value) {
    const r = canonicalCostRoleGroup(s.user_role || '')
    if (!groups[r]) groups[r] = { role: r, members: [] }
    groups[r].members.push(s)
  }
  return order
    .filter(r => groups[r])
    .map(r => groups[r])
    .concat(Object.values(groups).filter(g => !order.includes(g.role)))
})

async function loadSalaries() {
  loadingSalaries.value = true
  salaryError.value = ''
  try {
    salaries.value = await api.listSalaries()
  } catch (e: any) {
    salaryError.value = e?.message ?? 'Failed to load salaries'
  } finally {
    loadingSalaries.value = false
  }
}

function confirmDelete(sal: EmployeeSalaryWithUser) {
  salaryToDelete.value = sal
}

async function doDelete() {
  if (!salaryToDelete.value) return
  deletingId.value = salaryToDelete.value.id
  try {
    await api.deleteSalary(salaryToDelete.value.id)
    salaries.value = salaries.value.filter(s => s.id !== salaryToDelete.value!.id)
    salaryToDelete.value = null
  } catch (e: any) {
    salaryError.value = e?.message ?? 'Delete failed'
    salaryToDelete.value = null
  } finally {
    deletingId.value = null
  }
}

// ── Salary Modal ──────────────────────────────────────────────────────────────

const showSalaryModal = ref(false)
const editingSalary = ref<EmployeeSalaryWithUser | null>(null)
const modalError = ref('')
const modalSubmitting = ref(false)

const empTypes = [
  { value: 'FULLTIME', label: 'Full-time' },
  { value: 'PARTTIME', label: 'Part-time' },
  { value: 'CONTRACTOR', label: 'Contractor' },
] as const

const defaultSalaryForm = () => ({
  user_id: 0,
  monthly_salary: 0,
  effective_from: new Date().toISOString().substring(0, 10),
  effective_to: '',
  employment_type: 'FULLTIME' as 'FULLTIME' | 'PARTTIME' | 'CONTRACTOR',
})
const salaryForm = reactive(defaultSalaryForm())

function openAddSalaryModal() {
  editingSalary.value = null
  modalError.value = ''
  Object.assign(salaryForm, defaultSalaryForm())
  showSalaryModal.value = true
}

function openEditSalaryModal(sal: EmployeeSalaryWithUser) {
  editingSalary.value = sal
  modalError.value = ''
  Object.assign(salaryForm, {
    user_id: sal.user_id,
    monthly_salary: sal.monthly_salary,
    effective_from: sal.effective_from?.substring(0, 10) ?? '',
    effective_to: sal.effective_to?.substring(0, 10) ?? '',
    employment_type: sal.employment_type,
  })
  showSalaryModal.value = true
}

function closeSalaryModal() {
  showSalaryModal.value = false
  editingSalary.value = null
}

async function submitSalaryModal() {
  if (!salaryForm.user_id) { modalError.value = 'Please select an employee.'; return }
  if (!salaryForm.monthly_salary || salaryForm.monthly_salary <= 0) { modalError.value = 'Monthly salary must be > 0.'; return }
  if (!salaryForm.effective_from) { modalError.value = 'Effective From date is required.'; return }

  modalSubmitting.value = true
  modalError.value = ''
  try {
    await api.upsertSalary({
      user_id: salaryForm.user_id,
      monthly_salary: salaryForm.monthly_salary,
      effective_from: salaryForm.effective_from,
      effective_to: salaryForm.effective_to || undefined,
      employment_type: salaryForm.employment_type,
    })
    closeSalaryModal()
    await loadSalaries()
  } catch (e: any) {
    modalError.value = e?.data?.error ?? e?.message ?? 'Save failed'
  } finally {
    modalSubmitting.value = false
  }
}

// ── Company Config ────────────────────────────────────────────────────────────

const loadingConfig = ref(false)
const savingConfig = ref(false)
const configSaved = ref(false)
const configError = ref('')

const configForm = reactive<UpdateCostConfigPayload>({
  working_days_per_month: 22,
  working_hours_per_day: 8,
  overhead_multiplier: 1.30,
  default_profit_margin: 0.25,
  default_risk_buffer: 0.10,
  currency: 'THB',
  executive_expense: 0,
  company_expense: 0,
})

const billableDaysPreview = computed(() => billableDays.value.toFixed(1))

async function loadConfig() {
  loadingConfig.value = true
  configError.value = ''
  try {
    const cfg = await api.getCostConfig()
    Object.assign(configForm, {
      working_days_per_month: cfg.working_days_per_month,
      working_hours_per_day: cfg.working_hours_per_day,
      overhead_multiplier: cfg.overhead_multiplier,
      default_profit_margin: cfg.default_profit_margin,
      default_risk_buffer: cfg.default_risk_buffer,
      currency: cfg.currency,
      executive_expense: cfg.executive_expense ?? 0,
      company_expense: cfg.company_expense ?? 0,
    })
  } catch (e: any) {
    configError.value = e?.message ?? 'Failed to load config'
  } finally {
    loadingConfig.value = false
  }
}

async function saveConfig() {
  savingConfig.value = true
  configSaved.value = false
  configError.value = ''
  try {
    await api.updateCostConfig({ ...configForm })
    configSaved.value = true
    setTimeout(() => { configSaved.value = false }, 3000)
  } catch (e: any) {
    configError.value = e?.data?.error ?? e?.message ?? 'Save failed'
  } finally {
    savingConfig.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────────

function formatTHB(val: number): string {
  return new Intl.NumberFormat('th-TH', {
    style: 'currency',
    currency: 'THB',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val)
}

function fmtN(val: number): string {
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 0 }).format(val)
}

function pct(val: number): string {
  return `${(val * 100).toFixed(0)}%`
}

function initials(name: string): string {
  return name.split(/[\s@.]+/).filter(Boolean).slice(0, 2).map(p => p[0].toUpperCase()).join('')
}

function roleChip(role: string): string {
  const base = 'px-2 py-0.5 rounded-full text-xs font-bold tracking-wide border'
  const map: Record<string, string> = {
    CEO:  `${base} bg-yellow-500/10 text-yellow-400 border-yellow-500/30`,
    MANAGER: `${base} bg-orange-500/10 text-orange-400 border-orange-500/30`,
    PRODUCT_OWNER: `${base} bg-blue-500/10 text-blue-400 border-blue-500/30`,
    PM:   `${base} bg-blue-500/10 text-blue-400 border-blue-500/30`,
    ENGINEER: `${base} bg-purple-500/10 text-purple-400 border-purple-500/30`,
    CHIEF_ENGINEER: `${base} bg-violet-500/10 text-violet-300 border-violet-500/30`,
    DEV:  `${base} bg-purple-500/10 text-purple-400 border-purple-500/30`,
    SUPPORT: `${base} bg-cyan-500/10 text-cyan-400 border-cyan-500/30`,
  }
  return map[role] ?? `${base} bg-gray-700 text-gray-300 border-gray-600`
}

function empTypeChip(type: string): string {
  const base = 'px-2 py-0.5 rounded text-xs font-medium'
  const map: Record<string, string> = {
    FULLTIME:   `${base} bg-green-900/30 text-green-400`,
    PARTTIME:   `${base} bg-blue-900/30 text-blue-400`,
    CONTRACTOR: `${base} bg-orange-900/30 text-orange-400`,
  }
  return map[type] ?? `${base} bg-gray-700 text-gray-300`
}

function vcRunwayColor(months: number): string {
  if (months <= 0) return 'text-gray-600'
  if (months > 2) return 'text-emerald-400'
  if (months > 1) return 'text-yellow-400'
  return 'text-red-400'
}

function vcRunwayBarColor(months: number): string {
  if (months <= 0) return 'bg-gray-600'
  if (months > 2) return 'bg-emerald-500'
  if (months > 1) return 'bg-yellow-500'
  return 'bg-red-500'
}

function vcRunwayBarWidth(months: number): string {
  return Math.min((months / 3) * 100, 100) + '%'
}

// ── Init ──────────────────────────────────────────────────────────────────────

onMounted(async () => {
  if (!isCEO.value) return
  await teamsStore.fetchTeamsFeatureEnabled()
  await Promise.all([fetchTeam(), loadSalaries(), loadConfig(), loadMandayRate()])
  if (teamsStore.teamsFeatureEnabled) {
    if (teamsStore.teams.length === 0) await teamsStore.fetchTeams()
    for (const team of teamsStore.teams) {
      teamsStore.fetchTeamCost(team.id)
    }
  }
})
</script>

<!-- Inline error banner component -->
<script lang="ts">
export default { name: 'CostConfigPage' }
</script>

<style scoped>
.input-field {
  @apply w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white
         placeholder-gray-500 transition-colors
         focus:border-amber-500 focus:outline-none focus:ring-1 focus:ring-amber-500/50;
}

.field-label {
  @apply block text-xs font-medium text-gray-400 mb-1.5;
}

.field-hint {
  @apply text-xs text-gray-600 mt-1;
}

.btn-primary {
  @apply inline-flex items-center justify-center rounded-lg
         bg-gradient-to-r from-purple-600 to-pink-600
         px-4 py-2.5 text-sm font-semibold text-white shadow-lg transition-all
         hover:from-purple-500 hover:to-pink-500 hover:shadow-purple-500/25
         disabled:cursor-not-allowed disabled:opacity-50;
}

.btn-ghost {
  @apply inline-flex items-center justify-center rounded-lg
         border border-gray-600 bg-transparent
         px-4 py-2.5 text-sm font-semibold text-gray-300 transition-colors
         hover:bg-gray-700 hover:text-white
         disabled:cursor-not-allowed disabled:opacity-50;
}

.btn-danger {
  @apply inline-flex items-center justify-center rounded-lg
         bg-red-600 px-4 py-2.5 text-sm font-semibold text-white transition-colors
         hover:bg-red-500
         disabled:cursor-not-allowed disabled:opacity-50;
}

.metric-card {
  @apply rounded-xl border border-gray-700 bg-gray-800/60 px-5 py-4;
}
.metric-label {
  @apply text-xs font-medium uppercase tracking-wider text-gray-500;
}
.metric-value {
  @apply mt-1 text-2xl font-extrabold text-white tabular-nums;
}

.config-card {
  @apply rounded-xl border border-gray-700 bg-gray-800/60 p-5;
}
.config-card-header {
  @apply flex items-center gap-2.5 mb-1;
}
.config-card-header h3 {
  @apply text-sm font-semibold text-white;
}
.config-card-desc {
  @apply text-xs text-gray-500 leading-relaxed;
}

/* PDF Preview modal fade transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
