<template>
  <div class="min-h-screen bg-gray-900 text-white">

    <!-- ── Sticky Page Header ─────────────────────────────────────────────────── -->
    <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-6 py-4 backdrop-blur-sm">
      <div class="flex flex-col gap-3 max-w-screen-xl mx-auto sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-3 min-w-0">
          <div class="w-8 h-8 rounded-lg bg-purple-500/15 border border-purple-500/30 flex items-center justify-center shrink-0">
            <svg class="w-4 h-4 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
            </svg>
          </div>
          <div class="min-w-0">
            <h1 class="text-base font-bold text-white">Execution Dashboard</h1>
            <p class="text-xs text-gray-500">
              <template v-if="activeView === 'pipeline'">My Board = your work in <span class="text-gray-400">active sprints</span> · {{ myTeamName ? myTeamName + ' Board' : 'Global Board' }} = TASK/BUG queue across projects</template>
              <template v-else>Deep work & task pipeline</template>
            </p>
          </div>
        </div>
        <div class="flex items-center gap-3 shrink-0">
          <div class="flex rounded-lg border border-gray-700 bg-gray-800/60 overflow-hidden text-xs font-semibold">
            <button
              type="button"
              @click="activeView = 'board'"
              class="px-3 py-2 transition-colors"
              :class="activeView === 'board' ? 'bg-indigo-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-700'"
            >
              {{ myTeamName ? myTeamName + ' Board' : 'Global Board' }}
            </button>
            <button
              type="button"
              @click="activeView = 'pipeline'"
              class="px-3 py-2 transition-colors"
              :class="activeView === 'pipeline' ? 'bg-purple-600 text-white' : 'text-gray-400 hover:text-white hover:bg-gray-700'"
            >
              My Board
            </button>
          </div>
          <button
            v-if="activeView === 'pipeline'"
            type="button"
            :disabled="isLoading"
            class="inline-flex items-center gap-1.5 rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2 text-xs font-medium text-gray-300 hover:border-gray-600 hover:bg-gray-700 hover:text-white transition-colors disabled:opacity-50"
            @click="fetchTasks"
          >
            <svg class="h-3.5 w-3.5" :class="isLoading ? 'animate-spin' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- ── Squad P&L (Global Board only) ── -->
    <div v-if="activeView === 'board' && teamsStore.teamsFeatureEnabled" class="max-w-screen-xl mx-auto px-6 pt-6">
      <TeamFinancialWidget
        :team-member-ids="myTeamMemberIds"
        :completed-task-ids="[]"
        :total-capital="squadTotalCapital"
        :project-capital-count="squadProjectCapitalCount"
        :loaded-monthly-burn-rate="squadLoadedBurnRate"
      />
    </div>

    <main v-if="activeView === 'board'" class="max-w-screen-xl mx-auto px-6 py-8">
      <GlobalKanbanBoard />
    </main>

    <template v-if="activeView === 'pipeline'">

    <div v-if="isLoading" class="flex flex-col items-center justify-center py-32">
      <svg class="h-8 w-8 animate-spin text-purple-400 mb-3" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
      </svg>
      <p class="text-sm text-gray-500">Loading your tasks…</p>
    </div>

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

    <main v-else class="max-w-screen-xl mx-auto px-6 py-8 space-y-8">

      <!-- Your performance — first block; KPI grid always expanded (no hide) -->
      <section>
        <h2 class="section-label">Your performance</h2>
        <div v-if="performanceStore.personal" class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
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
        <p v-else-if="performanceStore.loading" class="text-sm text-gray-500 py-2">Loading metrics…</p>
        <p v-else class="text-sm text-gray-500 py-2">Performance data unavailable.</p>
      </section>

      <!-- Active sprints: multi-chip strip + open sprint -->
      <section v-if="activeSprints.length > 0" class="space-y-3">
        <div class="flex flex-wrap items-center gap-2">
          <span class="text-[10px] font-bold uppercase tracking-widest text-gray-500 w-full sm:w-auto">Active sprints</span>
        </div>
        <div class="flex flex-wrap gap-2">
          <button
            type="button"
            class="rounded-xl border px-3 py-2 text-left text-xs font-semibold transition-colors max-w-full sm:max-w-[280px]"
            :class="selectedSprintId === '' ? 'border-purple-500/50 bg-purple-500/15 text-white' : 'border-gray-700 bg-gray-800/50 text-gray-400 hover:border-gray-600 hover:text-gray-200'"
            @click="selectedSprintId = ''"
          >
            All sprints
            <span class="block font-normal text-[10px] text-gray-500 mt-0.5">Show everything below</span>
          </button>
          <div
            v-for="s in activeSprints"
            :key="s.id"
            class="flex items-stretch rounded-xl border overflow-hidden max-w-full sm:max-w-[320px]"
            :class="selectedSprintId === s.id ? 'border-purple-500/50 ring-1 ring-purple-500/30' : 'border-gray-700'"
          >
            <button
              type="button"
              class="flex-1 text-left px-3 py-2 min-w-0 transition-colors"
              :class="selectedSprintId === s.id ? 'bg-purple-500/10' : 'bg-gray-800/40 hover:bg-gray-800/70'"
              @click="selectedSprintId = s.id"
            >
              <span class="text-xs font-bold text-white truncate block">{{ s.name }}</span>
              <span class="text-[10px] text-gray-500 truncate block">{{ s.project_name || 'Project' }}</span>
              <span v-if="s.end_date" class="text-[10px] text-amber-400/90 mt-0.5 block">{{ getSprintCountdown(s.end_date) }}</span>
            </button>
            <NuxtLink
              :to="`/projects/sprint/${s.id}?project=${encodeURIComponent(s.project_code || s.project_id)}`"
              class="shrink-0 flex items-center px-2.5 border-l border-gray-700/80 text-gray-500 hover:text-purple-300 hover:bg-gray-800/80 text-[10px] font-semibold uppercase tracking-wide"
              title="Open sprint in project"
            >
              Open
            </NuxtLink>
          </div>
        </div>
      </section>

      <section v-else class="rounded-xl border border-amber-500/30 bg-amber-950/20 px-5 py-4 flex items-start gap-3">
        <svg class="h-5 w-5 shrink-0 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
        </svg>
        <p class="text-sm text-amber-300">No active sprint found. Ask your PM to start a sprint and assign tasks to you.</p>
      </section>

      <!-- Filters -->
      <section v-if="myTasks.length > 0" class="flex flex-col sm:flex-row sm:flex-wrap gap-3 items-stretch sm:items-end">
        <div class="flex flex-col gap-1 flex-1 min-w-[140px]">
          <label class="text-[10px] font-bold uppercase tracking-widest text-gray-500">Project</label>
          <select v-model="filterProjectId" class="rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2 text-sm text-white focus:border-purple-500 focus:outline-none focus:ring-1 focus:ring-purple-500/40">
            <option value="">All projects</option>
            <option v-for="p in projectFilterOptions" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>
        <div class="flex flex-col gap-1 flex-1 min-w-[140px]">
          <label class="text-[10px] font-bold uppercase tracking-widest text-gray-500">Status</label>
          <select v-model="filterStatus" class="rounded-lg border border-gray-700 bg-gray-800/60 px-3 py-2 text-sm text-white focus:border-purple-500 focus:outline-none focus:ring-1 focus:ring-purple-500/40">
            <option value="">All statuses</option>
            <option value="IN_PROGRESS">In progress</option>
            <option value="PENDING">Pending</option>
            <option value="READY_FOR_TEST">Ready for test</option>
            <option value="REVIEW_PENDING">Review pending</option>
            <option value="READY_FOR_UAT">Ready for UAT</option>
            <option value="BLOCKED">Blocked</option>
            <option value="COMPLETED">Completed</option>
          </select>
        </div>
        <button
          v-if="selectedSprintId || filterProjectId || filterStatus"
          type="button"
          class="text-xs font-medium text-gray-500 hover:text-white px-3 py-2 rounded-lg border border-transparent hover:border-gray-600 transition-colors"
          @click="clearFilters"
        >
          Clear filters
        </button>
      </section>

      <!-- Filter empty -->
      <section
        v-if="myTasks.length > 0 && filteredMyTasks.length === 0"
        class="rounded-xl border border-gray-700 bg-gray-800/30 px-5 py-8 text-center"
      >
        <p class="text-sm font-medium text-gray-300">No tasks match your filters</p>
        <p class="text-xs text-gray-500 mt-1">Try another sprint, project, or status — or clear filters.</p>
        <button type="button" class="mt-3 text-xs font-semibold text-purple-400 hover:text-purple-300" @click="clearFilters">Clear filters</button>
      </section>

      <!-- Waiting on others -->
      <section v-if="waitingOnReview.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">Waiting on review</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-cyan-500/10 border border-cyan-500/25 text-cyan-400">{{ waitingOnReview.length }}</span>
        </div>
        <div class="grid gap-3">
          <article
            v-for="task in waitingOnReview"
            :key="task.id"
            class="rounded-xl border border-cyan-500/25 bg-cyan-950/10 p-4 cursor-pointer hover:border-cyan-400/40 transition-colors"
            @click="goToTask(task)"
          >
            <DevBoardTaskMeta :task="task" :status-label="statusLabel(task.status)" />
            <p class="font-semibold text-white mt-2">{{ task.title }}</p>
          </article>
        </div>
      </section>

      <!-- Blocked -->
      <section v-if="blockedTasks.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">Blocked</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-red-500/10 border border-red-500/25 text-red-400">{{ blockedTasks.length }}</span>
        </div>
        <div class="grid gap-3">
          <article
            v-for="task in blockedTasks"
            :key="task.id"
            class="rounded-xl border border-red-500/30 bg-red-950/10 p-4 cursor-pointer hover:border-red-400/50 transition-colors"
            @click="goToTask(task)"
          >
            <DevBoardTaskMeta :task="task" status-label="Blocked" />
            <p class="font-semibold text-white mt-2">{{ task.title }}</p>
          </article>
        </div>
      </section>

      <!-- Current focus -->
      <section v-if="currentFocusVisible.length > 0">
        <div class="flex items-center justify-between mb-4 gap-2 flex-wrap">
          <h2 class="section-label mb-0">Current focus</h2>
          <div class="flex items-center gap-2">
            <span class="text-xs font-medium px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/25 text-purple-400">
              {{ currentFocusSorted.length }} in progress
            </span>
            <button
              v-if="currentFocusSorted.length > FOCUS_LIMIT"
              type="button"
              class="text-xs font-medium text-gray-500 hover:text-white"
              @click="showAllFocus = !showAllFocus"
            >
              {{ showAllFocus ? 'Show less' : `Show all (${currentFocusSorted.length})` }}
            </button>
          </div>
        </div>
        <div class="grid gap-4">
          <div
            v-for="task in currentFocusVisible"
            :key="task.id"
            class="relative cursor-pointer rounded-2xl border-2 p-6 shadow-lg transition-all hover:shadow-purple-500/10"
            :class="focusCardClass(task)"
            @click="goToTask(task)"
          >
            <div
              v-if="getDeadlineUrgency(task) === 'overdue'"
              class="absolute -top-2 right-4 rounded-md bg-red-600 px-2.5 py-0.5 text-xs font-bold text-white"
            >Overdue</div>
            <div
              v-else-if="getDeadlineUrgency(task) === 'urgent'"
              class="absolute -top-2 right-4 rounded-md bg-amber-500 px-2.5 py-0.5 text-xs font-bold text-black"
            >Urgent</div>

            <DevBoardTaskMeta :task="task" status-label="In progress" detailed />
            <p v-if="task.description" class="mt-2 line-clamp-2 text-sm text-gray-400">{{ task.description }}</p>

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
                type="button"
                @click.stop="markReadyForTest(task)"
                :disabled="markingReadyId === task.id"
                class="rounded-lg bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 px-4 py-2 text-sm font-semibold text-white shadow-md transition-all disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg v-if="markingReadyId === task.id" class="inline mr-1.5 h-3.5 w-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                {{ markingReadyId === task.id ? 'Submitting…' : 'Ready for Test' }}
              </button>
              <NuxtLink
                :to="`/task/${task.code || task.id}?from=dashboard&from_tab=${activeView}`"
                class="rounded-lg border border-gray-700 bg-gray-800/60 px-4 py-2 text-sm font-medium text-gray-300 hover:bg-gray-700 hover:text-white transition-colors"
                @click.stop
              >
                View details
              </NuxtLink>
            </div>
          </div>
        </div>
      </section>

      <!-- Up next -->
      <section v-if="upNextTasks.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">Up next</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-amber-500/10 border border-amber-500/25 text-amber-400">{{ upNextTasks.length }} pending</span>
        </div>
        <div class="rounded-2xl border border-gray-700 bg-gray-800/60 overflow-hidden shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/60 text-xs uppercase tracking-wider text-gray-500">
                  <th class="px-5 py-3.5 text-left">Task</th>
                  <th class="px-5 py-3.5 text-right hidden sm:table-cell">Est.</th>
                  <th class="px-5 py-3.5 text-right">Deadline</th>
                  <th class="px-5 py-3.5 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/60">
                <tr
                  v-for="task in upNextTasks"
                  :key="task.id"
                  class="cursor-pointer hover:bg-gray-700/20 transition-colors"
                  :class="rowUrgencyClass(task)"
                  @click="goToTask(task)"
                >
                  <td class="px-5 py-3.5">
                    <DevBoardTaskMeta :task="task" :status-label="'Pending'" dense />
                    <p class="font-medium text-white mt-1">{{ task.title }}</p>
                  </td>
                  <td class="px-5 py-3.5 text-right text-gray-400 tabular-nums hidden sm:table-cell">{{ (task.estimated_minutes / 60).toFixed(1) }}h</td>
                  <td class="px-5 py-3.5 text-right text-xs">
                    <template v-if="task.due_at">
                      <p :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-500'">
                        {{ formatDeadline(task.due_at) }}
                      </p>
                      <p class="text-gray-600 mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</p>
                    </template>
                    <span v-else class="text-gray-600">—</span>
                  </td>
                  <td class="px-5 py-3.5 text-right">
                    <button
                      type="button"
                      class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 px-3 py-1.5 text-xs font-semibold text-white transition-all"
                      @click.stop="goToTask(task)"
                    >
                      Open
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Recently done -->
      <section v-if="doneTasks.length > 0">
        <button
          type="button"
          class="flex items-center justify-between w-full text-left mb-4 group"
          @click="doneSectionOpen = !doneSectionOpen"
        >
          <h2 class="section-label mb-0">Recently done</h2>
          <span class="text-xs font-medium text-gray-500 group-hover:text-gray-300">{{ doneSectionOpen ? 'Hide' : `Show (${doneTasks.length})` }}</span>
        </button>
        <div v-if="doneSectionOpen" class="grid gap-2">
          <div
            v-for="task in doneTasks.slice(0, DONE_PREVIEW_LIMIT)"
            :key="task.id"
            class="rounded-xl border border-emerald-500/20 bg-emerald-950/5 px-4 py-3 cursor-pointer hover:border-emerald-500/40 transition-colors flex items-center justify-between gap-3"
            @click="goToTask(task)"
          >
            <div class="min-w-0">
              <DevBoardTaskMeta :task="task" status-label="Done" dense />
              <p class="text-sm font-medium text-gray-300 truncate">{{ task.title }}</p>
            </div>
            <span class="text-[10px] font-semibold uppercase text-emerald-500/80 shrink-0">Done</span>
          </div>
        </div>
      </section>

      <!-- My board (3 columns — respects sprint/project/status filters) -->
      <section v-if="myTasks.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h2 class="section-label mb-0">My board</h2>
          <span class="text-xs font-medium px-3 py-1 rounded-full bg-purple-500/10 border border-purple-500/25 text-purple-400">
            {{ filteredMyTasks.length }} in view
          </span>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
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
                <DevBoardTaskMeta :task="task" status-label="Pending" dense />
                <p class="text-sm font-semibold text-white line-clamp-2 mt-1">{{ task.title }}</p>
                <div class="mt-2 flex items-center justify-between gap-2 text-xs text-gray-500">
                  <span class="tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span>
                  <span
                    v-if="task.due_at"
                    :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-600'"
                  >{{ getDeadlineCountdown(task.due_at) }}</span>
                </div>
              </div>
              <p v-if="boardPending.length === 0" class="text-center text-xs text-gray-600 py-6">No pending</p>
            </div>
          </div>

          <div class="rounded-2xl border border-purple-500/30 bg-purple-950/10 flex flex-col">
            <div class="flex items-center gap-2 px-4 py-3 border-b border-purple-500/20">
              <span class="h-2 w-2 rounded-full bg-purple-400 animate-pulse"></span>
              <span class="text-xs font-bold uppercase tracking-widest text-purple-400">In progress</span>
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
                <DevBoardTaskMeta :task="task" status-label="In progress" dense />
                <p class="text-sm font-semibold text-white line-clamp-2 mt-1">{{ task.title }}</p>
                <div class="mt-2 flex items-center justify-between gap-2 text-xs text-gray-500">
                  <span class="tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</span>
                  <span
                    v-if="task.due_at"
                    :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-600'"
                  >{{ getDeadlineCountdown(task.due_at) }}</span>
                </div>
                <button
                  type="button"
                  class="mt-2 w-full rounded-lg bg-gradient-to-r from-cyan-600 to-teal-600 hover:from-cyan-500 hover:to-teal-500 px-3 py-1 text-xs font-semibold text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                  :disabled="markingReadyId === task.id"
                  @click.stop="markReadyForTest(task)"
                >
                  <svg v-if="markingReadyId === task.id" class="inline mr-1 h-3 w-3 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                  </svg>
                  {{ markingReadyId === task.id ? 'Submitting…' : 'Ready for Test' }}
                </button>
              </div>
              <p v-if="boardInProgress.length === 0" class="text-center text-xs text-gray-600 py-6">Nothing in progress</p>
            </div>
          </div>

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
                <DevBoardTaskMeta :task="task" status-label="Done" dense />
                <p class="text-sm font-semibold text-gray-300 line-clamp-2 mt-1">{{ task.title }}</p>
                <div class="mt-2 text-xs text-gray-600 tabular-nums">{{ (task.estimated_minutes / 60).toFixed(1) }}h</div>
              </div>
              <p v-if="boardDone.length === 0" class="text-center text-xs text-gray-600 py-6">No completed yet</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Empty -->
      <section v-if="myTasks.length === 0">
        <div class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-gray-700 bg-gray-800/30 py-20 text-center">
          <div class="w-14 h-14 rounded-2xl bg-gray-700/50 flex items-center justify-center mb-4">
            <svg class="h-7 w-7 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
          </div>
          <p class="text-base font-semibold text-gray-300">No tasks in active sprints</p>
          <p class="text-sm text-gray-500 mt-1 mb-6">When a PM assigns you work in a started sprint, it will show up here.</p>
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

    </template>

    <DailyCheckinModal :forced="true" />

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
import DevBoardTaskMeta from '~/components/dashboard/DevBoardTaskMeta.vue'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import type { ProjectCapitalResponse } from '~/core/modules/projects/infrastructure/projects-api'

interface Task {
  id: string
  code?: string
  title: string
  description: string
  status: string
  task_type?: string
  priority?: string
  estimated_minutes: number
  assigned_to?: number
  project_id?: string | null
  sprint_id?: string | null
  project_name?: string
  project_color?: string
  sprint_name?: string
  effective_sprint_id?: string | null
  created_at: string
  due_at?: string
  started_at?: string
  completed_at?: string
}

interface ActiveSprint {
  id: string
  project_id: string
  project_name?: string
  project_code?: string
  name: string
  goal: string
  end_date: string | null
  status: string
}

const FOCUS_LIMIT = 5
const DONE_PREVIEW_LIMIT = 12

const PRI_ORDER: Record<string, number> = { CRITICAL: 0, HIGH: 1, MEDIUM: 2, LOW: 3 }

const { fetchWithAuth, currentUser } = useAuth()
const performanceStore = usePerformanceStore()
const teamsStore = useTeamsStore()
const { getTeams } = useTeamsApi()
const { getProjects, getProjectCapital } = useProjectsApi()

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
    // non-fatal
  }
}

const route = useRoute()
const router = useRouter()

const activeView = ref<'board' | 'pipeline'>(
  route.query.tab === 'pipeline' ? 'pipeline' : 'board'
)

watch(activeView, (tab) => {
  router.replace({ query: { ...route.query, tab } })
}, { immediate: true })

const myTasks = ref<Task[]>([])
const activeSprints = ref<ActiveSprint[]>([])
const isLoading = ref(true)
const error = ref('')
const selectedSprintId = ref('')
const filterProjectId = ref('')
const filterStatus = ref('')
const showAllFocus = ref(false)
const doneSectionOpen = ref(false)

const showSuccess = ref(false)
const successMessage = ref('')
const markingReadyId = ref<string | null>(null)

function effectiveSprintKey(task: Task): string {
  return task.effective_sprint_id || task.sprint_id || ''
}

const filteredMyTasks = computed(() => {
  let list = myTasks.value
  if (selectedSprintId.value) {
    const sid = selectedSprintId.value
    list = list.filter(t => effectiveSprintKey(t) === sid)
  }
  if (filterProjectId.value) {
    list = list.filter(t => t.project_id === filterProjectId.value)
  }
  if (filterStatus.value) {
    list = list.filter(t => t.status === filterStatus.value)
  }
  return list
})

const projectFilterOptions = computed(() => {
  const m = new Map<string, string>()
  for (const t of myTasks.value) {
    if (!t.project_id) continue
    const name = t.project_name || t.project_id
    if (!m.has(t.project_id)) m.set(t.project_id, name)
  }
  return [...m.entries()].map(([id, name]) => ({ id, name })).sort((a, b) => a.name.localeCompare(b.name))
})

function sortByDueThenPriority(a: Task, b: Task): number {
  const da = a.due_at ? new Date(a.due_at).getTime() : Number.POSITIVE_INFINITY
  const db = b.due_at ? new Date(b.due_at).getTime() : Number.POSITIVE_INFINITY
  if (da !== db) return da - db
  const pa = a.priority ? (PRI_ORDER[a.priority] ?? 9) : 9
  const pb = b.priority ? (PRI_ORDER[b.priority] ?? 9) : 9
  return pa - pb
}

const waitingOnReview = computed(() =>
  filteredMyTasks.value
    .filter(t => ['READY_FOR_TEST', 'REVIEW_PENDING', 'READY_FOR_UAT'].includes(t.status))
    .sort(sortByDueThenPriority)
)

const blockedTasks = computed(() =>
  filteredMyTasks.value.filter(t => t.status === 'BLOCKED').sort(sortByDueThenPriority)
)

const currentFocusSorted = computed(() =>
  filteredMyTasks.value.filter(t => t.status === 'IN_PROGRESS').sort(sortByDueThenPriority)
)

const currentFocusVisible = computed(() => {
  const all = currentFocusSorted.value
  if (showAllFocus.value || all.length <= FOCUS_LIMIT) return all
  return all.slice(0, FOCUS_LIMIT)
})

const upNextTasks = computed(() =>
  filteredMyTasks.value.filter(t => t.status === 'PENDING').sort(sortByDueThenPriority)
)

const doneTasks = computed(() =>
  filteredMyTasks.value
    .filter(t => t.status === 'COMPLETED' || t.status === 'DONE')
    .sort((a, b) =>
      new Date(b.completed_at || b.created_at || 0).getTime() - new Date(a.completed_at || a.created_at || 0).getTime()
    )
)

const boardPending = computed(() =>
  filteredMyTasks.value.filter(t => t.status === 'PENDING').sort(sortByDueThenPriority)
)

const boardInProgress = computed(() =>
  filteredMyTasks.value.filter(t => t.status === 'IN_PROGRESS').sort(sortByDueThenPriority)
)

const boardDone = computed(() =>
  filteredMyTasks.value
    .filter(t => t.status === 'COMPLETED' || t.status === 'DONE')
    .sort(
      (a, b) =>
        new Date(b.completed_at || b.created_at || 0).getTime() - new Date(a.completed_at || a.created_at || 0).getTime()
    )
)

function clearFilters() {
  selectedSprintId.value = ''
  filterProjectId.value = ''
  filterStatus.value = ''
}

function statusLabel(status: string): string {
  if (status === 'READY_FOR_TEST') return 'Ready for test'
  if (status === 'REVIEW_PENDING') return 'Review pending'
  if (status === 'READY_FOR_UAT') return 'Ready for UAT'
  return status.replace(/_/g, ' ')
}

function focusCardClass(task: Task): string {
  const u = getDeadlineUrgency(task)
  if (u === 'overdue') return 'border-red-500/50 bg-red-950/10'
  if (u === 'urgent') return 'border-amber-500/50 bg-amber-950/10'
  return 'border-purple-500/40 bg-purple-950/10'
}

function rowUrgencyClass(task: Task): string {
  const u = getDeadlineUrgency(task)
  if (u === 'overdue') return 'bg-red-950/15'
  if (u === 'urgent') return 'bg-amber-950/15'
  return ''
}

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

const fetchTasks = async () => {
  isLoading.value = true
  error.value = ''
  try {
    const myRes = await fetchWithAuth<{ data: Task[]; active_sprints: ActiveSprint[] }>('/sentinel/tasks/my')
    myTasks.value = myRes.data || []
    activeSprints.value = myRes.active_sprints || []
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load tasks'
  } finally {
    isLoading.value = false
  }
}

function getDeadlineUrgency(task: Task) {
  if (!task.due_at || task.status === 'COMPLETED' || task.status === 'DONE') return 'none'
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

const markReadyForTest = async (task: Task) => {
  markingReadyId.value = task.id
  try {
    await fetchWithAuth(`/sentinel/tasks/${task.id}/ready-for-test`, { method: 'POST' })
    myTasks.value = myTasks.value.map(t =>
      t.id === task.id ? { ...t, status: 'READY_FOR_TEST' } : t
    )
    successMessage.value = 'Marked as Ready for Test — awaiting PM approval'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)
  } catch (err: any) {
    successMessage.value = err.data?.message || err.message || 'Failed to mark ready for test'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 4000)
  } finally {
    markingReadyId.value = null
  }
}

const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}?from=dashboard&from_tab=${activeView.value}`)
}

onMounted(async () => {
  fetchTasks()
  await teamsStore.fetchTeamsFeatureEnabled()
  if (teamsStore.teamsFeatureEnabled) {
    fetchSquadFinancials()
  }
  performanceStore.fetchAll('DEV')
})
</script>

<style scoped>
.section-label {
  @apply text-xs font-semibold uppercase tracking-widest text-gray-500 mb-4;
}
.toast-enter-active, .toast-leave-active { transition: all 0.25s ease; }
.toast-enter-from, .toast-leave-to { opacity: 0; transform: translateY(0.5rem); }
</style>
