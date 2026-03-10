<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Sticky Page Header ─────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex items-center justify-between max-w-screen-xl mx-auto">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center">
            <svg class="w-4 h-4 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
            </svg>
          </div>
          <div>
            <h1 class="text-base font-bold text-white">Execution Dashboard</h1>
            <p class="text-xs text-gray-500">Deep work & task pipeline</p>
          </div>
        </div>
        <!-- Tab switcher + inline stats -->
        <div class="flex items-center gap-3">
          <!-- View tabs -->
          <div class="flex rounded-lg border border-gray-700 bg-gray-800/60 overflow-hidden text-xs font-semibold">
            <button
              @click="activeView = 'board'"
              class="px-3 py-2 transition-colors"
              :class="activeView === 'board' ? 'bg-indigo-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-700'"
            >
              {{ myTeamName ? myTeamName + ' Board' : 'Global Board' }}
            </button>
            <button
              @click="activeView = 'pipeline'"
              class="px-3 py-2 transition-colors"
              :class="activeView === 'pipeline' ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-700'"
            >
              My Board
            </button>
          </div>

        </div>
      </div>
    </div>

    <!-- ── Squad P&L (Global Board only) ───────────────────────────────────────── -->
    <div v-if="activeView === 'board'" class="max-w-screen-xl mx-auto px-6 pt-6">
      <TeamFinancialWidget
        :team-member-ids="myTeamMemberIds"
        :completed-task-ids="[]"
        :total-capital="squadTotalCapital"
        :project-capital-count="squadProjectCapitalCount"
        :loaded-monthly-burn-rate="squadLoadedBurnRate"
      />
    </div>

    <!-- ── Global Board View ───────────────────────────────────────────────────── -->
    <main v-if="activeView === 'board'" class="max-w-screen-xl mx-auto px-6 py-8">
      <GlobalKanbanBoard />
    </main>

    <!-- ── Pipeline View ─────────────────────────────────────────────────────── -->
    <template v-if="activeView === 'pipeline'">

    <!-- ── Loading ───────────────────────────────────────────────────────────── -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-purple-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading your tasks…</p>
    </div>

    <!-- ── Error ─────────────────────────────────────────────────────────────── -->
    <div v-else-if="error" class="max-w-screen-xl mx-auto px-6 py-8">
      <div class="flex items-start gap-3 rounded-xl border border-red-500/30 bg-red-900/20 px-5 py-4 text-red-400">
        <svg class="h-5 w-5 flex-shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <div>
          <p class="font-semibold text-sm">Failed to load data</p>
          <p class="text-xs text-red-300 mt-0.5">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- ── Main Content ───────────────────────────────────────────────────────── -->
    <main v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">

      <!-- ── Active Sprint Banner (Blinder) ─────────────────────────────────── -->
      <section v-if="primaryActiveSprint">
        <div class="relative overflow-hidden rounded-2xl border border-purple-500/40 bg-gradient-to-r from-purple-950/60 via-gray-900/80 to-gray-900/60 px-6 py-5 shadow-lg shadow-purple-500/10">
          <!-- Glow accent -->
          <div class="pointer-events-none absolute -top-10 -left-10 h-48 w-48 rounded-full bg-purple-600/10 blur-3xl" />
          <div class="relative flex flex-col sm:flex-row sm:items-center gap-4">
            <div class="flex items-center gap-3 min-w-0 flex-1">
              <div class="shrink-0 w-10 h-10 rounded-xl bg-purple-500/20 border border-purple-500/40 flex items-center justify-center">
                <svg class="w-5 h-5 text-purple-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                </svg>
              </div>
              <div class="min-w-0">
                <div class="flex items-center gap-2 mb-0.5">
                  <span class="text-[10px] font-bold uppercase tracking-widest text-purple-400">Active Sprint</span>
                  <span class="inline-flex h-1.5 w-1.5 rounded-full bg-green-400 animate-pulse" />
                </div>
                <h2 class="text-base font-bold text-white truncate">{{ primaryActiveSprint.name }}</h2>
                <p v-if="primaryActiveSprint.goal" class="text-xs text-gray-400 mt-0.5 line-clamp-1">{{ primaryActiveSprint.goal }}</p>
              </div>
            </div>
            <div v-if="primaryActiveSprint.end_date" class="shrink-0 text-right">
              <p class="text-[10px] font-semibold uppercase tracking-widest text-gray-500">Sprint Ends</p>
              <p class="text-sm font-bold text-amber-400 mt-0.5">{{ formatSprintEndDate(primaryActiveSprint.end_date) }}</p>
              <p class="text-xs text-gray-500 mt-0.5">{{ getSprintCountdown(primaryActiveSprint.end_date) }}</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ── No Active Sprint Warning ───────────────────────────────────────── -->
      <section v-else-if="!isLoading">
        <div class="flex items-center gap-3 rounded-xl border border-amber-500/30 bg-amber-950/20 px-5 py-4">
          <svg class="h-5 w-5 shrink-0 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
          </svg>
          <p class="text-sm text-amber-300">No active sprint found. Ask your PM to start a sprint and assign tasks to you.</p>
        </div>
      </section>

      <!-- ── My Performance ─────────────────────────────────────────────────── -->
      <section v-if="performanceStore.personal">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">Your performance</h2>
          <button
            @click="showKpiPanel = !showKpiPanel"
            class="text-xs font-medium text-gray-500 hover:text-white transition-colors"
          >
            {{ showKpiPanel ? 'Hide' : 'Show' }}
          </button>
        </div>
        <template v-if="showKpiPanel">
          <div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
            <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5">Delivery Rate</p>
              <p class="text-xl font-extrabold tabular-nums" :class="pctColor(performanceStore.personal.delivery_rate_pct)">{{ performanceStore.personal.delivery_rate_pct.toFixed(1) }}%</p>
            </div>
            <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4 flex flex-col justify-between">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-2">Code Quality</p>
              <span class="inline-flex items-center gap-1 self-start rounded-full bg-purple-500/10 border border-purple-500/30 px-2.5 py-0.5 text-xs font-medium text-purple-400">
                <span class="w-1.5 h-1.5 rounded-full bg-purple-400 animate-pulse"></span>
                Coming soon
              </span>
            </div>
            <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5">Rework Rate</p>
              <p class="text-xl font-extrabold tabular-nums" :class="reworkColor(performanceStore.personal.rework_rate_pct)">{{ performanceStore.personal.rework_rate_pct.toFixed(1) }}%</p>
            </div>
            <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5">Time Accuracy</p>
              <p class="text-xl font-extrabold tabular-nums" :class="pctColor(performanceStore.personal.time_accuracy_pct)">{{ performanceStore.personal.time_accuracy_pct.toFixed(1) }}%</p>
            </div>
            <div class="rounded-2xl border border-gray-700 bg-gray-800/60 p-4">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5">Sprint Velocity</p>
              <p class="text-xl font-extrabold text-white tabular-nums">{{ performanceStore.personal.sprint_velocity_sp }}</p>
              <p class="text-xs text-gray-500 mt-0.5">SP (last 3)</p>
            </div>
            <div class="rounded-2xl border border-gray-700 bg-gradient-to-br from-purple-900/30 to-pink-900/20 p-4">
              <p class="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-1.5">Score</p>
              <p class="text-xl font-black tabular-nums" :class="scoreColor(compositeScore)">{{ compositeScore.toFixed(1) }}</p>
              <p class="text-xs text-gray-500 mt-0.5">0–100 composite</p>
            </div>
          </div>
        </template>
      </section>

      <!-- ── Current Focus ──────────────────────────────────────────────────── -->
      <section v-if="currentFocus.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">Current focus</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/25 text-purple-400">
            {{ currentFocus.length }} in progress
          </span>
        </div>
        <div class="grid gap-4">
          <div
            v-for="task in currentFocus"
            :key="task.id"
            class="relative cursor-pointer rounded-2xl border-2 p-6 shadow-lg transition-all hover:shadow-purple-500/10"
            :class="[
              getDeadlineUrgency(task) === 'overdue' ? 'border-red-500/50 bg-red-950/10' :
              getDeadlineUrgency(task) === 'urgent'  ? 'border-amber-500/50 bg-amber-950/10' :
              'border-purple-500/40 bg-purple-950/10'
            ]"
            @click="goToTask(task)"
          >
            <!-- Urgency badge -->
            <div
              v-if="getDeadlineUrgency(task) === 'overdue'"
              class="absolute -top-2 right-4 rounded-md bg-red-600 px-2.5 py-0.5 text-xs font-bold text-white"
            >Overdue</div>
            <div
              v-else-if="getDeadlineUrgency(task) === 'urgent'"
              class="absolute -top-2 right-4 rounded-md bg-amber-500 px-2.5 py-0.5 text-xs font-bold text-black"
            >Urgent</div>

            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <h3 class="text-lg font-bold text-white">{{ task.title }}</h3>
                <p class="mt-1 line-clamp-2 text-sm text-gray-400">{{ task.description || 'No description' }}</p>
              </div>
              <span class="shrink-0 rounded-lg bg-purple-600/70 px-3 py-1 text-xs font-semibold text-purple-100">
                In progress
              </span>
            </div>

            <div class="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1.5 text-xs text-gray-500">
              <span>Est: <span class="font-semibold text-purple-400">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span></span>
              <template v-if="task.started_at">
                <span class="text-gray-700">·</span>
                <span>Started: <span class="text-gray-300">{{ formatDate(task.started_at) }}</span></span>
              </template>
              <template v-if="task.due_at">
                <span class="text-gray-700">·</span>
                <span
                  class="font-semibold"
                  :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400' : 'text-gray-400'"
                >
                  Due {{ formatDeadline(task.due_at) }} · {{ getDeadlineCountdown(task.due_at) }}
                </span>
              </template>
            </div>

            <div class="mt-4 flex gap-3">
              <button
                @click.stop="openSubmitModal(task)"
                class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 px-4 py-2 text-sm font-semibold text-white shadow-md transition-all"
              >
                Submit work
              </button>
              <NuxtLink
                :to="`/task/${task.code || task.id}?from=dashboard&from_tab=${activeView}`"
                @click.stop
                class="rounded-lg border border-gray-700 bg-gray-800/60 px-4 py-2 text-sm font-medium text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
              >
                View details
              </NuxtLink>
            </div>
          </div>
        </div>
      </section>

      <!-- ── My Backlog ─────────────────────────────────────────────────────── -->
      <section v-if="myBacklog.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">My backlog</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-amber-500/10 border border-amber-500/25 text-amber-400">
            {{ myBacklog.length }} pending
          </span>
        </div>
        <div class="rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/60 text-xs uppercase tracking-wider text-gray-500">
                  <th class="px-5 py-3.5 text-left">Task</th>
                  <th class="px-5 py-3.5 text-right">Est. hours</th>
                  <th class="px-5 py-3.5 text-right">Deadline</th>
                  <th class="px-5 py-3.5 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/60">
                <tr
                  v-for="task in myBacklog"
                  :key="task.id"
                  class="cursor-pointer hover:bg-gray-700/20 transition-colors"
                  :class="[
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-950/15' :
                    getDeadlineUrgency(task) === 'urgent'  ? 'bg-amber-950/15' : ''
                  ]"
                  @click="goToTask(task)"
                >
                  <td class="px-5 py-3.5">
                    <p class="font-medium text-white">
                      {{ task.title }}
                      <span v-if="getDeadlineUrgency(task) === 'overdue'" class="ml-2 rounded bg-red-600 px-1.5 py-0.5 text-xs font-semibold text-white">Overdue</span>
                      <span v-else-if="getDeadlineUrgency(task) === 'urgent'" class="ml-2 rounded bg-amber-500 px-1.5 py-0.5 text-xs font-semibold text-black">Urgent</span>
                    </p>
                  </td>
                  <td class="px-5 py-3.5 text-right text-gray-400 tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</td>
                  <td class="px-5 py-3.5 text-right text-xs">
                    <template v-if="task.due_at">
                      <p :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-500'">
                        {{ formatDeadline(task.due_at) }}
                      </p>
                      <p class="text-gray-600 mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</p>
                    </template>
                    <span v-else class="text-gray-600">No deadline</span>
                  </td>
                  <td class="px-5 py-3.5 text-right">
                    <button
                      @click.stop="goToTask(task)"
                      class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 px-3 py-1.5 text-xs font-semibold text-white transition-all"
                    >
                      Start
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- ── My Board (All Assigned Tasks) ────────────────────────────────── -->
      <section v-if="myTasks.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">My board</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/25 text-purple-400">
            {{ myTasks.length }} assigned
          </span>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Pending column -->
          <div class="rounded-2xl border border-gray-700 bg-gray-800/40 flex flex-col">
            <div class="flex items-center gap-2 px-4 py-3 border-b border-gray-700">
              <span class="h-2 w-2 rounded-full bg-amber-400"></span>
              <span class="text-xs font-bold uppercase tracking-widest text-gray-400">Pending</span>
              <span class="ml-auto text-xs font-semibold text-amber-400 tabular-nums">{{ boardPending.length }}</span>
            </div>
            <div class="flex-1 p-3 space-y-2 overflow-y-auto max-h-[480px]">
              <div
                v-for="task in boardPending"
                :key="task.id"
                class="rounded-xl border border-gray-700 bg-gray-800/60 p-3 cursor-pointer hover:border-amber-500/40 hover:bg-amber-950/10 transition-all"
                :class="getDeadlineUrgency(task) === 'overdue' ? 'border-red-500/40 bg-red-950/10' : getDeadlineUrgency(task) === 'urgent' ? 'border-amber-500/40 bg-amber-950/10' : ''"
                @click="goToTask(task)"
              >
                <p class="text-sm font-semibold text-white line-clamp-2">{{ task.title }}</p>
                <div class="mt-2 flex items-center justify-between gap-2 text-xs text-gray-500">
                  <span class="tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span>
                  <span
                    v-if="task.due_at"
                    :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-600'"
                  >{{ getDeadlineCountdown(task.due_at) }}</span>
                </div>
              </div>
              <p v-if="boardPending.length === 0" class="text-center text-xs text-gray-600 py-6">No pending tasks</p>
            </div>
          </div>

          <!-- In Progress column -->
          <div class="rounded-2xl border border-purple-500/30 bg-purple-950/10 flex flex-col">
            <div class="flex items-center gap-2 px-4 py-3 border-b border-purple-500/20">
              <span class="h-2 w-2 rounded-full bg-purple-400 animate-pulse"></span>
              <span class="text-xs font-bold uppercase tracking-widest text-purple-400">In Progress</span>
              <span class="ml-auto text-xs font-semibold text-purple-400 tabular-nums">{{ boardInProgress.length }}</span>
            </div>
            <div class="flex-1 p-3 space-y-2 overflow-y-auto max-h-[480px]">
              <div
                v-for="task in boardInProgress"
                :key="task.id"
                class="rounded-xl border border-purple-500/30 bg-gray-800/60 p-3 cursor-pointer hover:border-purple-400/60 hover:bg-purple-950/20 transition-all"
                :class="getDeadlineUrgency(task) === 'overdue' ? 'border-red-500/40 bg-red-950/10' : getDeadlineUrgency(task) === 'urgent' ? 'border-amber-500/40 bg-amber-950/10' : ''"
                @click="goToTask(task)"
              >
                <p class="text-sm font-semibold text-white line-clamp-2">{{ task.title }}</p>
                <div class="mt-2 flex items-center justify-between gap-2 text-xs text-gray-500">
                  <span class="tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span>
                  <span
                    v-if="task.due_at"
                    :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-600'"
                  >{{ getDeadlineCountdown(task.due_at) }}</span>
                </div>
                <button
                  @click.stop="openSubmitModal(task)"
                  class="mt-2 w-full rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 px-3 py-1 text-xs font-semibold text-white transition-all"
                >
                  Submit work
                </button>
              </div>
              <p v-if="boardInProgress.length === 0" class="text-center text-xs text-gray-600 py-6">Nothing in progress</p>
            </div>
          </div>

          <!-- Completed column -->
          <div class="rounded-2xl border border-emerald-500/20 bg-emerald-950/5 flex flex-col">
            <div class="flex items-center gap-2 px-4 py-3 border-b border-emerald-500/20">
              <span class="h-2 w-2 rounded-full bg-emerald-400"></span>
              <span class="text-xs font-bold uppercase tracking-widest text-emerald-500">Done</span>
              <span class="ml-auto text-xs font-semibold text-emerald-400 tabular-nums">{{ boardDone.length }}</span>
            </div>
            <div class="flex-1 p-3 space-y-2 overflow-y-auto max-h-[480px]">
              <div
                v-for="task in boardDone"
                :key="task.id"
                class="rounded-xl border border-gray-700/60 bg-gray-800/40 p-3 cursor-pointer hover:border-emerald-500/30 hover:bg-emerald-950/10 transition-all opacity-75 hover:opacity-100"
                @click="goToTask(task)"
              >
                <p class="text-sm font-semibold text-gray-300 line-clamp-2">{{ task.title }}</p>
                <div class="mt-2 flex items-center gap-2 text-xs text-gray-600">
                  <span class="tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span>
                </div>
              </div>
              <p v-if="boardDone.length === 0" class="text-center text-xs text-gray-600 py-6">No completed tasks yet</p>
            </div>
          </div>
        </div>
      </section>

      <!-- ── Empty State ────────────────────────────────────────────────────── -->
      <section v-if="myTasks.length === 0">
        <div class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-700 bg-gray-800/30 py-20 text-center">
          <div class="w-14 h-14 rounded-2xl bg-gray-700/50 flex items-center justify-center mb-4">
            <svg class="h-7 w-7 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
          </div>
          <p class="text-base font-semibold text-gray-300">No tasks available</p>
          <p class="text-sm text-gray-500 mt-1 mb-6">Create a new task or wait for assignments</p>
          <NuxtLink
            to="/create"
            class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 px-5 py-2.5 text-sm font-semibold text-white shadow-lg transition-all"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            Create new task
          </NuxtLink>
        </div>
      </section>
    </main>

    </template><!-- end pipeline view -->

    <!-- ── Daily Check-in (forced) ────────────────────────────────────────────── -->
    <DailyCheckinModal :forced="true" />

    <!-- ── Submit Work Modal ──────────────────────────────────────────────────── -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showSubmitModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
          @keydown.escape="closeSubmitModal"
        >
          <div class="fixed inset-0 bg-black/70 backdrop-blur-sm" @click="closeSubmitModal"/>
          <div
            class="relative max-h-[90vh] w-full max-w-3xl overflow-y-auto rounded-2xl border-2 border-purple-500/40 bg-gray-800 shadow-2xl"
            @click.stop
          >
            <div class="sticky top-0 z-10 flex items-start justify-between border-b border-gray-700 bg-gray-800 px-6 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center">
                  <svg class="w-4 h-4 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div>
                  <h2 class="text-sm font-bold text-white">Hand over mission</h2>
                  <p class="text-xs text-gray-500 mt-0.5 line-clamp-1">{{ selectedTask?.title }}</p>
                </div>
              </div>
              <button @click="closeSubmitModal" class="rounded-lg p-2 text-gray-400 hover:bg-gray-700 hover:text-white transition-colors">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>
            <form class="p-6 space-y-4" @submit.prevent="submitWork">
              <div v-if="submitError" class="rounded-lg border border-red-500/30 bg-red-900/20 px-4 py-3 text-sm text-red-400">
                {{ submitError }}
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 mb-1.5">
                  Pull Request / Commit URL <span class="text-red-400">*</span>
                </label>
                <input
                  v-model="submitForm.referenceUrl"
                  type="url"
                  placeholder="https://github.com/org/repo/pull/123"
                  class="input-field"
                />
                <p class="text-xs text-gray-600 mt-1">GitHub PR, GitLab MR, or direct commit URL</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-400 mb-1.5">Note <span class="text-gray-600">(optional)</span></label>
                <textarea
                  v-model="submitForm.note"
                  rows="4"
                  placeholder="What was changed and why? Any known issues?"
                  class="input-field resize-none"
                />
              </div>
              <div class="flex gap-3 pt-2">
                <button type="button" @click="closeSubmitModal" :disabled="isSubmitting" class="btn-ghost flex-1">Cancel</button>
                <button
                  type="submit"
                  :disabled="isSubmitting || !submitForm.referenceUrl"
                  class="btn-primary flex-1"
                >
                  <svg v-if="isSubmitting" class="mr-2 h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                  </svg>
                  {{ isSubmitting ? 'Sending…' : 'Hand Over for Review' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- ── Success Toast ─────────────────────────────────────────────────────── -->
    <Teleport to="body">
      <Transition name="toast">
        <div
          v-if="showSuccess"
          class="fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-xl border border-emerald-500/40 bg-gray-800 px-5 py-3.5 shadow-xl shadow-black/30"
        >
          <span class="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-600 text-white text-sm font-bold">✓</span>
          <p class="text-sm font-medium text-white">{{ successMessage }}</p>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { usePerformanceStore } from '~/core/modules/performance/performance-store'
import DailyCheckinModal from '~/core/modules/pulse/ui/DailyCheckinModal.vue'
import GlobalKanbanBoard from '~/components/tasks/GlobalKanbanBoard.vue'
import TeamFinancialWidget from '~/components/dashboard/TeamFinancialWidget.vue'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { ProjectCapitalResponse } from '~/core/modules/projects/infrastructure/projects-api'

interface Task {
  id: string
  code?: string
  title: string
  description: string
  status: string
  estimated_minutes: number
  assigned_to?: number
  created_at: string
  due_at?: string
  started_at?: string
  completed_at?: string
}

interface ActiveSprint {
  id: string
  name: string
  goal: string
  end_date: string | null
  status: string
}

const { fetchWithAuth, currentUser } = useAuth()
const performanceStore = usePerformanceStore()
const { getTeams } = useTeamsApi()
const { getProjects, getProjectCapital } = useProjectsApi()

// --- Squad P&L state (read-only for DEV) ---
const myTeamMemberIds = ref<number[]>([])
const myTeamName = ref<string>('')
const squadProjectCapitals = ref<Record<string, ProjectCapitalResponse>>({})

const squadTotalCapital = computed(() =>
  Object.values(squadProjectCapitals.value).reduce((s, c) => s + (c.capital_balance ?? 0), 0)
)
const squadLoadedBurnRate = computed(() => {
  const vals = Object.values(squadProjectCapitals.value)
  if (!vals.length) return 0
  return vals.reduce((sum, c) => sum + c.team_monthly_cost, 0)
})
const squadProjectCapitalCount = computed(() => Object.keys(squadProjectCapitals.value).length)

const fetchSquadFinancials = async () => {
  try {
    const userId = currentUser.value?.user_id
    const [teams, projects] = await Promise.all([getTeams(), getProjects()])
    const myTeam = teams.find(t => t.users?.some(u => u.id === userId))
    myTeamMemberIds.value = myTeam?.users?.map(u => u.id) ?? []
    myTeamName.value = myTeam?.name ?? ''

    const capitals = await Promise.allSettled(
      projects.map(p => getProjectCapital(p.id).then(c => ({ id: p.id, capital: c })))
    )
    const map: Record<string, ProjectCapitalResponse> = {}
    for (const r of capitals) {
      if (r.status === 'fulfilled') map[r.value.id] = r.value.capital
    }
    squadProjectCapitals.value = map
  } catch {
    // non-fatal — widget shows skeleton
  }
}

const route = useRoute()
const router = useRouter()

// Restore active tab from query param — persists across refresh and back navigation
const activeView = ref<'board' | 'pipeline'>(
  route.query.tab === 'pipeline' ? 'pipeline' : 'board'
)

// Keep URL in sync so refresh stays on the same tab
watch(activeView, (tab) => {
  router.replace({ query: { ...route.query, tab } })
}, { immediate: true })

const myTasks = ref<Task[]>([])
const unassignedTasks = ref<Task[]>([])
const activeSprints = ref<ActiveSprint[]>([])
const isLoading = ref(true)
const error = ref('')
const claiming = ref<string | null>(null)
const showKpiPanel = ref(true)

const primaryActiveSprint = computed(() => activeSprints.value[0] ?? null)

const showSuccess = ref(false)
const successMessage = ref('')
const submitError = ref('')

const showSubmitModal = ref(false)
const isSubmitting = ref(false)
const selectedTask = ref<Task | null>(null)
const submitForm = ref({ referenceUrl: '', note: '' })

function pctColor(pct: number): string {
  if (pct >= 80) return 'text-emerald-400'
  if (pct >= 60) return 'text-amber-400'
  return 'text-red-400'
}

function scoreColor(score: number): string {
  if (score >= 70) return 'text-emerald-400'
  if (score >= 50) return 'text-amber-400'
  return 'text-red-400'
}

function reworkColor(pct: number): string {
  if (pct <= 10) return 'text-emerald-400'
  if (pct <= 20) return 'text-amber-400'
  return 'text-red-400'
}

const compositeScore = computed(() => {
  const p = performanceStore.personal
  if (!p) return 0
  const qualityNorm = Math.min(100, p.code_quality_index)
  const reworkNorm = Math.max(0, 100 - p.rework_rate_pct)
  const velocityNorm = Math.min(100, p.sprint_velocity_sp * 5)
  return 0.30 * p.delivery_rate_pct + 0.25 * qualityNorm + 0.20 * reworkNorm + 0.15 * velocityNorm + 0.10 * p.time_accuracy_pct
})

const currentFocus = computed(() =>
  myTasks.value
    .filter(t => t.status === 'IN_PROGRESS')
    .sort((a, b) => (!a.due_at ? 1 : !b.due_at ? -1 : new Date(a.due_at).getTime() - new Date(b.due_at).getTime()))
)

const myBacklog = computed(() =>
  myTasks.value
    .filter(t => t.status === 'PENDING')
    .sort((a, b) => (!a.due_at ? 1 : !b.due_at ? -1 : new Date(a.due_at).getTime() - new Date(b.due_at).getTime()))
)

const availableMissions = computed(() =>
  unassignedTasks.value.sort((a, b) => (!a.due_at ? 1 : !b.due_at ? -1 : new Date(a.due_at).getTime() - new Date(b.due_at).getTime()))
)

const boardPending = computed(() =>
  myTasks.value
    .filter(t => t.status === 'PENDING')
    .sort((a, b) => (!a.due_at ? 1 : !b.due_at ? -1 : new Date(a.due_at).getTime() - new Date(b.due_at).getTime()))
)

const boardInProgress = computed(() =>
  myTasks.value
    .filter(t => t.status === 'IN_PROGRESS')
    .sort((a, b) => (!a.due_at ? 1 : !b.due_at ? -1 : new Date(a.due_at).getTime() - new Date(b.due_at).getTime()))
)

const boardDone = computed(() =>
  myTasks.value
    .filter(t => t.status === 'COMPLETED' || t.status === 'DONE')
    .sort((a, b) => new Date(b.completed_at || b.created_at || 0).getTime() - new Date(a.completed_at || a.created_at || 0).getTime())
)

const fetchTasks = async () => {
  isLoading.value = true
  error.value = ''
  try {
    const [myRes, unassignedRes] = await Promise.all([
      fetchWithAuth<{ data: Task[]; active_sprints: ActiveSprint[] }>('/sentinel/tasks/my'),
      fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/unassigned'),
    ])
    myTasks.value = myRes.data || []
    activeSprints.value = myRes.active_sprints || []
    unassignedTasks.value = unassignedRes.data || []
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load tasks'
  } finally {
    isLoading.value = false
  }
}

const claimTask = async (taskId: string) => {
  if (!currentUser.value?.user_id) { error.value = 'User ID not found.'; return }
  claiming.value = taskId
  try {
    await fetchWithAuth(`/sentinel/tasks/${taskId}/assign`, {
      method: 'POST',
      body: { dev_id: currentUser.value.user_id },
    })
    const claimed = unassignedTasks.value.find(t => t.id === taskId)
    if (claimed) {
      claimed.status = 'IN_PROGRESS'
      claimed.assigned_to = currentUser.value.user_id
      unassignedTasks.value = unassignedTasks.value.filter(t => t.id !== taskId)
      myTasks.value.unshift(claimed)
    }
    successMessage.value = 'Task claimed successfully'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to claim task'
  } finally {
    claiming.value = null
  }
}

function getDeadlineUrgency(task: Task) {
  if (!task.due_at || task.status === 'COMPLETED') return 'none'
  const hoursUntilDue = (new Date(task.due_at).getTime() - Date.now()) / 3600000
  if (hoursUntilDue < 0) return 'overdue'
  if (hoursUntilDue < 24) return 'urgent'
  return 'normal'
}

function formatDate(dateString: string) {
  return new Date(dateString).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function formatDeadline(dueAt: string) {
  return new Date(dueAt).toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatSprintEndDate(endDate: string) {
  return new Date(endDate).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function getSprintCountdown(endDate: string) {
  const diff = new Date(endDate).getTime() - Date.now()
  if (diff < 0) return 'Sprint overdue'
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  if (days === 0) return 'Ends today'
  if (days === 1) return '1 day left'
  return `${days} days left`
}

function getDeadlineCountdown(dueAt: string) {
  const diff = new Date(dueAt).getTime() - Date.now()
  if (diff < 0) {
    const hours = Math.abs(Math.floor(diff / 3600000))
    const days = Math.floor(hours / 24)
    return days > 0 ? `Overdue by ${days}d` : `Overdue by ${hours}h`
  }
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(hours / 24)
  if (days > 0) return `${days}d left`
  if (hours > 0) return `${hours}h left`
  return 'Due soon!'
}

const openSubmitModal = (task: Task) => {
  selectedTask.value = task
  submitForm.value = { referenceUrl: '', note: '' }
  submitError.value = ''
  showSubmitModal.value = true
}

const closeSubmitModal = () => {
  if (!isSubmitting.value) {
    showSubmitModal.value = false
    selectedTask.value = null
    submitForm.value = { referenceUrl: '', note: '' }
  }
}

const submitWork = async () => {
  if (!selectedTask.value || !submitForm.value.referenceUrl) return
  isSubmitting.value = true
  submitError.value = ''
  try {
    await fetchWithAuth(`/sentinel/tasks/${selectedTask.value.id}/submit`, {
      method: 'POST',
      body: { reference_url: submitForm.value.referenceUrl, note: submitForm.value.note },
    })
    showSubmitModal.value = false
    selectedTask.value = null
    successMessage.value = 'Handed over — awaiting PM review'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)
    await fetchTasks()
  } catch (err: any) {
    submitError.value = err.data?.message || err.message || 'Failed to submit handover'
  } finally {
    isSubmitting.value = false
  }
}

const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}?from=dashboard&from_tab=${activeView.value}`)
}

onMounted(() => {
  fetchTasks()
  fetchSquadFinancials()
  performanceStore.fetchAll('DEV')
})
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
.input-field {
  @apply w-full rounded-lg border border-gray-600 bg-gray-900/50 px-3 py-2.5 text-sm text-white
         placeholder-gray-500 transition-colors
         focus:border-purple-500 focus:outline-none focus:ring-1 focus:ring-purple-500/50;
}
.btn-primary {
  @apply inline-flex items-center justify-center rounded-lg bg-gradient-to-r from-purple-600 to-pink-600
         px-4 py-2.5 text-sm font-semibold text-white shadow-lg transition-all
         hover:from-purple-500 hover:to-pink-500
         disabled:cursor-not-allowed disabled:opacity-50;
}
.btn-ghost {
  @apply inline-flex items-center justify-center rounded-lg border border-gray-600 bg-transparent
         px-4 py-2.5 text-sm font-semibold text-gray-300 transition-colors
         hover:bg-gray-700 hover:text-white
         disabled:cursor-not-allowed disabled:opacity-50;
}
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
.toast-enter-active, .toast-leave-active { transition: all 0.25s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(0.5rem); }
</style>
