<template>
  <div class="min-h-screen bg-gray-900 text-white">
    <!-- Loading -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-[60vh]">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <p class="text-sm text-gray-500">กำลังโหลดโปรเจกต์...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="p-8 max-w-2xl mx-auto">
      <div class="bg-red-900/20 border border-red-500 rounded-xl p-6 text-red-400">
        <h2 class="font-bold text-lg mb-1">Failed to load project</h2>
        <p class="text-sm">{{ error }}</p>
        <NuxtLink to="/projects" class="mt-4 inline-block text-sm text-gray-400 hover:text-white">← Back to Projects</NuxtLink>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="project">
      <!-- Project Header (responsive) -->
      <div class="border-b border-gray-800 bg-gray-900/95 sticky top-0 z-20 px-3 sm:px-6 py-3 sm:py-4">
        <div class="flex flex-wrap items-center gap-2 sm:gap-4">
          <NuxtLink to="/projects" class="text-gray-500 hover:text-gray-300 transition-colors text-sm shrink-0">
            ← Projects
          </NuxtLink>
          <span class="text-gray-700 hidden sm:inline">/</span>
          <div class="flex-1 min-w-0 flex items-center gap-2 sm:gap-3">
            <h1 class="text-base sm:text-lg font-bold text-white truncate">{{ project.name }}</h1>
            <span
              class="px-2 py-0.5 text-xs font-semibold rounded-full border shrink-0"
              :class="statusClass(project.status)"
            >
              {{ project.status.replace('_', ' ') }}
            </span>
            <code class="text-xs text-gray-500 font-mono hidden md:inline shrink-0">{{ project.code }}</code>
          </div>
        </div>

        <!-- Tabs (horizontal scroll on small screens) -->
        <div class="flex gap-1 mt-3 sm:mt-4 overflow-x-auto pb-1 -mx-1 px-1 scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-transparent">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="px-3 sm:px-4 py-2 text-xs sm:text-sm rounded-lg transition-colors font-medium whitespace-nowrap shrink-0"
            :class="activeTab === tab.id
              ? 'bg-indigo-600 text-white'
              : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800'"
          >
            {{ tab.icon }} {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Tab Content -->
      <div class="p-3 sm:p-6">
        <!-- TAB 1: Overview -->
        <div v-if="activeTab === 'overview'" class="space-y-6">
          <!-- Key Metrics -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="metric-card">
              <div class="text-2xl font-bold" :class="activeSprint ? 'text-indigo-400' : 'text-gray-500'">
                {{ activeSprint ? activeSprint.name : 'No sprint' }}
              </div>
              <div class="metric-label">Active Sprint</div>
              <div v-if="activeSprint" class="mt-2 text-xs text-gray-500">
                {{ activeSprint.end_date ? `Ends ${formatDate(activeSprint.end_date)}` : 'No end date' }}
              </div>
            </div>
            <div class="metric-card">
              <div class="text-2xl font-bold text-green-400">{{ completedCount }}/{{ totalTasks }}</div>
              <div class="metric-label">Tasks Done</div>
              <div class="mt-2 h-1 bg-gray-700 rounded-full overflow-hidden">
                <div class="h-full bg-green-500 rounded-full" :style="{ width: completionPct + '%' }"></div>
              </div>
            </div>
            <div class="metric-card">
              <div class="text-2xl font-bold text-yellow-400">{{ inProgressCount }}</div>
              <div class="metric-label">In Progress</div>
            </div>
            <div class="metric-card">
              <div class="text-2xl font-bold" :class="overdueCount > 0 ? 'text-red-400' : 'text-green-400'">
                {{ overdueCount }}
              </div>
              <div class="metric-label">Overdue</div>
            </div>
          </div>

          <!-- Row 2: Active Sprint + Milestones -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Active Sprint Card -->
            <div class="card">
              <div class="flex items-center justify-between mb-4">
                <h3 class="section-title">Current Sprint</h3>
                <div class="flex gap-2">
                  <button v-if="activeSprint" @click="openCompleteSprintModal(activeSprint)" class="btn-ghost-sm text-yellow-400 hover:text-yellow-300">
                    Complete Sprint
                  </button>
                  <button @click="openSprintModal()" class="btn-primary-sm shrink-0">+ Sprint</button>
                </div>
              </div>
              <div v-if="activeSprint" class="space-y-3">
                <div>
                  <p class="text-sm font-semibold text-white mb-1">{{ activeSprint.name }}</p>
                  <p v-if="activeSprint.goal" class="text-xs text-gray-400">{{ activeSprint.goal }}</p>
                </div>
                <div class="grid grid-cols-3 gap-3 text-center">
                  <div class="p-2 bg-gray-700/50 rounded-lg">
                    <div class="text-sm font-bold text-gray-200">{{ sprintTaskCount('total') }}</div>
                    <div class="text-[10px] text-gray-500">Total</div>
                  </div>
                  <div class="p-2 bg-gray-700/50 rounded-lg">
                    <div class="text-sm font-bold text-green-400">{{ sprintTaskCount('done') }}</div>
                    <div class="text-[10px] text-gray-500">Done</div>
                  </div>
                  <div class="p-2 bg-gray-700/50 rounded-lg">
                    <div class="text-sm font-bold text-indigo-400">{{ sprintTaskCount('sp') }}</div>
                    <div class="text-[10px] text-gray-500">Story Pts</div>
                  </div>
                </div>
                <div v-if="sprintTaskCount('total') > 0">
                  <div class="flex justify-between text-xs text-gray-500 mb-1">
                    <span>Sprint Progress</span>
                    <span>{{ Math.round(sprintTaskCount('done') / sprintTaskCount('total') * 100) }}%</span>
                  </div>
                  <div class="h-2 bg-gray-700 rounded-full overflow-hidden">
                    <div
                      class="h-full bg-indigo-500 rounded-full"
                      :style="{ width: Math.round(sprintTaskCount('done') / sprintTaskCount('total') * 100) + '%' }"
                    ></div>
                  </div>
                </div>
              </div>
              <div v-else class="text-center py-8 text-gray-500 text-sm">
                No active sprint. Plan and start a sprint to begin tracking.
              </div>
              <!-- List of all sprints (so user sees where created sprints are) -->
              <div v-if="sprints.length > 0" class="mt-4 pt-4 border-t border-gray-700">
                <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">All sprints</h4>
                <ul class="space-y-2 max-h-40 overflow-y-auto">
                  <li
                    v-for="s in sprintsWithActiveFirst"
                    :key="s.id"
                    class="flex items-center justify-between py-1.5 px-2 rounded-lg group"
                    :class="s.status === 'ACTIVE' ? 'bg-indigo-500/10' : 'hover:bg-gray-700/40'"
                  >
                    <NuxtLink
                      :to="`/projects/sprint/${s.id}?project=${route.params.id}`"
                      class="text-sm text-gray-200 hover:text-indigo-300 transition-colors truncate flex-1 min-w-0 mr-2"
                    >
                      {{ s.name }}
                    </NuxtLink>
                    <span class="flex items-center gap-2">
                      <span class="text-[10px] px-1.5 py-0.5 rounded font-medium" :class="s.status === 'ACTIVE' ? 'bg-indigo-500/20 text-indigo-400' : s.status === 'COMPLETED' ? 'bg-gray-600 text-gray-400' : 'bg-yellow-500/20 text-yellow-400'">
                        {{ s.status }}
                      </span>
                      <button
                        v-if="s.status === 'PLANNING'"
                        type="button"
                        :disabled="!!activeSprint"
                        @click.stop="!activeSprint && handleStartSprint(s.id)"
                        class="text-xs font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        :class="activeSprint ? 'text-gray-500' : 'text-indigo-400 hover:text-indigo-300'"
                        :title="activeSprint ? `มี sprint ที่ Active อยู่แล้ว (${activeSprint.name}) — ปิดหรือ Reopen ก่อน` : 'Start sprint'"
                      >
                        Start
                      </button>
                      <button
                        v-if="s.status === 'COMPLETED'"
                        type="button"
                        @click.stop="openReopenSprintModal(s)"
                        class="text-xs text-amber-400 hover:text-amber-300 font-medium"
                        title="Reopen sprint (set back to Active)"
                      >
                        Reopen
                      </button>
                      <button
                        type="button"
                        @click.stop="openAddTasksToSprintModal(s)"
                        class="text-xs text-green-400 hover:text-green-300"
                        title="Add tasks to this sprint"
                      >
                        + Tasks
                      </button>
                      <button
                        type="button"
                        @click.stop="openEditSprintModal(s)"
                        class="text-xs text-gray-400 hover:text-gray-200"
                        title="Edit sprint"
                      >
                        Edit
                      </button>
                      <button
                        type="button"
                        @click.stop="openDeleteSprintModal(s)"
                        class="text-xs text-red-400 hover:text-red-300"
                        title="Delete sprint"
                      >
                        Delete
                      </button>
                    </span>
                  </li>
                </ul>
                <p class="text-[10px] text-gray-500 mt-2">Sprints also appear in Backlog (Sprint column) and Board (Sprint filter).</p>
                <p class="text-[10px] text-amber-500/90 mt-1">หนึ่งโปรเจกต์มีได้แค่ 1 sprint ที่ Active — ต้องปิด (Complete) หรือ Reopen ก่อนจึงจะ Start อีก sprint ได้</p>
              </div>
            </div>

            <!-- Milestones -->
            <div class="card">
              <h3 class="section-title mb-4">Milestone Tracker</h3>
              <MilestoneTimeline
                :milestones="milestones"
                @add-milestone="openMilestoneModal"
                @milestone-click="openEditMilestoneModal"
              />
            </div>
          </div>

          <!-- Recent Activity -->
          <div class="card">
            <h3 class="section-title mb-4">Recent Activity</h3>
            <div v-if="recentTasks.length" class="space-y-2">
              <div
                v-for="t in recentTasks"
                :key="t.id"
                class="flex items-center justify-between py-2.5 px-3 rounded-lg hover:bg-gray-700/40 transition-colors cursor-pointer"
                @click="navigateToTask(t.id)"
              >
                <div class="flex items-center gap-3">
                  <span class="text-xs font-mono text-gray-500">{{ t.code }}</span>
                  <span class="text-sm text-gray-300 truncate max-w-xs">{{ t.title }}</span>
                  <span class="px-1.5 py-0.5 text-[10px] rounded font-medium" :class="priorityBadge(t.priority)">{{ t.priority }}</span>
                </div>
                <div class="flex items-center gap-3">
                  <span class="text-xs px-2 py-0.5 rounded-full" :class="taskStatusBadge(t.status)">{{ t.status.replace('_', ' ') }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center py-8 text-gray-500 text-sm">No tasks yet.</div>
          </div>
        </div>

        <!-- TAB 2: Board (Kanban) -->
        <div v-if="activeTab === 'board'">
          <KanbanBoard
            :tasks="allTasks"
            :sprints="sprints"
            @task-click="(t) => navigateToTask(t.id)"
            @status-change="handleStatusChange"
          />
        </div>

        <!-- TAB 3: Timeline (Gantt) - Enterprise design -->
        <div v-if="activeTab === 'timeline'" class="timeline-tab space-y-5">
          <!-- Toolbar: enterprise card -->
          <div class="timeline-toolbar rounded-xl border border-slate-600/60 bg-slate-800/80 px-4 py-3 shadow-lg shadow-black/20">
            <div class="flex flex-wrap items-center justify-between gap-4">
              <div class="flex flex-wrap items-center gap-4">
                <!-- Matrix Dimension Toggle: Epic Roadmap | Sprint Execution (both as Gantt) -->
                <div class="flex items-center gap-2">
                  <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Mode</span>
                  <div class="flex rounded-lg bg-slate-900/80 p-0.5">
                    <button @click="timelineMode = 'epic'" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="timelineMode === 'epic' ? 'bg-purple-600 text-white shadow-sm' : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'">
                      Epic Roadmap
                    </button>
                    <button @click="timelineMode = 'sprint'" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="timelineMode === 'sprint' ? 'bg-emerald-600 text-white shadow-sm' : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'">
                      Sprint Execution
                    </button>
                  </div>
                </div>

                <!-- View (Day/Week/Month) for both modes -->
                <div class="h-4 w-px bg-slate-600" />
                <div class="flex items-center gap-2">
                  <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">View</span>
                  <div class="flex rounded-lg bg-slate-900/80 p-0.5">
                    <button v-for="v in ['Day', 'Week', 'Month']" :key="v" @click="ganttView = v.toLowerCase()" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="ganttView === v.toLowerCase() ? (timelineMode === 'epic' ? 'bg-purple-600 text-white shadow-sm' : 'bg-emerald-600 text-white shadow-sm') : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'">
                      {{ v }}
                    </button>
                  </div>
                </div>
              </div>
              <button type="button" class="flex items-center gap-2 rounded-lg border border-indigo-500/50 bg-indigo-600/20 px-3 py-1.5 text-xs font-medium text-indigo-300 transition-colors hover:bg-indigo-600/40 hover:text-indigo-200" @click="scrollTimelineToToday">
                <span aria-hidden="true">◉</span> Today
              </button>
            </div>
          </div>

          <!-- Matrix Gantt: Epic Roadmap or Sprint Execution (Y = Epics/Sprints, X = timeline) -->
          <div v-if="matrixTimelineLoading" class="flex flex-col items-center justify-center py-20 text-slate-500">
            <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-600 mb-3" :class="timelineMode === 'epic' ? 'border-t-purple-400' : 'border-t-emerald-400'"></div>
            <p class="text-xs">{{ timelineMode === 'epic' ? 'Loading Epic Roadmap…' : 'Loading Sprint Execution…' }}</p>
          </div>

          <div v-else-if="timelineMode === 'epic' && (!epicTimelineData || !epicTimelineData.epics.length)" class="flex flex-col items-center justify-center rounded-xl border border-slate-600/50 bg-slate-800/50 py-20 text-center">
            <div class="mb-3 text-4xl opacity-50">🗺</div>
            <p class="text-sm font-medium text-slate-400">No epics yet</p>
            <p class="mt-1 text-xs text-slate-500">Create Epics in the Backlog tab to see the Roadmap</p>
          </div>

          <div v-else-if="timelineMode === 'sprint' && (!sprintTimelineData || !sprintTimelineData.sprints.length)" class="flex flex-col items-center justify-center rounded-xl border border-slate-600/50 bg-slate-800/50 py-20 text-center">
            <div class="mb-3 text-4xl opacity-50">🏃</div>
            <p class="text-sm font-medium text-slate-400">No sprints yet</p>
            <p class="mt-1 text-xs text-slate-500">Create Sprints to see the Execution View</p>
          </div>

          <ClientOnly v-else-if="matrixGanttRows.length > 0">
            <div ref="timelineScrollWrapperRef" class="timeline-scroll-wrapper rounded-xl border border-slate-600/60 bg-slate-800/60 shadow-xl shadow-black/25 overflow-x-auto overflow-y-hidden">
              <div class="timeline-inner relative flex flex-col" :style="matrixChartWidth > 0 ? { width: (220 + matrixChartWidth) + 'px', minWidth: (220 + matrixChartWidth) + 'px' } : { minWidth: '100%' }">
                <GanttMilestoneRow v-if="matrixDateRangeStart && matrixDateRangeEnd && matrixChartWidth > 0" :milestones="milestones" :date-range-start="matrixDateRangeStart" :date-range-end="matrixDateRangeEnd" :grid-width="matrixChartWidth" :grid-offset="220" @milestone-click="openEditMilestoneModal" />
                <g-gantt-chart :chart-start="matrixChartStart" :chart-end="matrixChartEnd" :precision="matrixGanttPrecision" bar-start="barStart" bar-end="barEnd" date-format="YYYY-MM-DD" :width="matrixChartWidth + 'px'" :row-height="40" :grid="true" :current-time="true" current-time-label="Now" color-scheme="dark" :label-column-title="timelineMode === 'epic' ? 'Epic / Task' : 'Sprint / Task'" label-column-width="220px" class="gantt-chart-vue gantt-enterprise" @click-bar="onMatrixGanttClickBar" @dragend-bar="onGanttDragEnd">
                  <template #label-column-row="{ label }">
                    <span class="cursor-pointer w-full block min-w-0 truncate" @click.stop="onMatrixLabelClickByLabel(label)">{{ label }}</span>
                  </template>
                  <g-gantt-row v-for="row in matrixGanttRows" :key="row.taskId" :label="row.label" :bars="row.bars" />
                </g-gantt-chart>
                <div v-if="matrixMilestoneLinePositions.length > 0" class="pointer-events-none absolute inset-0 z-[5]" aria-hidden="true">
                  <div v-for="{ id, left } in matrixMilestoneLinePositions" :key="id" class="absolute top-0 bottom-0 w-px bg-indigo-500/50" :style="{ left: left + 'px' }" />
                </div>
              </div>
            </div>
            <template #fallback>
              <div class="flex min-h-[420px] flex-col items-center justify-center rounded-xl border border-slate-600/50 bg-slate-800/50">
                <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-500 border-t-indigo-400" />
                <p class="mt-3 text-xs font-medium text-slate-400">Loading timeline…</p>
              </div>
            </template>
          </ClientOnly>

          <!-- Milestone legend (when chart is shown and we have milestones) -->
          <div v-if="milestones.length && matrixGanttRows.length > 0" class="rounded-xl border border-slate-600/40 bg-slate-800/40 px-4 py-3 mt-4">
            <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-500">Milestones</p>
            <div class="flex flex-wrap gap-x-6 gap-y-2">
              <div v-for="m in milestones" :key="m.id" class="flex items-center gap-2">
                <span class="milestone-legend-diamond rotate-45 border-2 border-indigo-400/80 bg-slate-800" />
                <span class="text-xs text-slate-300">{{ m.title }}</span>
                <span class="text-xs text-slate-500">{{ m.due_date ? formatDate(m.due_date) : '' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 4: Backlog (WBS + Sprint Planning) -->
        <div v-if="activeTab === 'backlog'" class="space-y-5">
          <!-- Epics Management Section -->
          <div class="bg-gray-800/60 border border-gray-700 rounded-xl p-4">
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-2">
                <span class="text-sm font-semibold text-gray-200">Epics</span>
                <span class="text-xs text-gray-500 bg-gray-700 px-1.5 py-0.5 rounded">{{ epics.length }}</span>
              </div>
              <button @click="openCreateEpicModal()" class="btn-primary-sm">+ Epic</button>
            </div>
            <div v-if="epics.length" class="flex flex-wrap gap-2">
              <div
                v-for="ep in epics"
                :key="ep.id"
                class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg border border-gray-600/50 bg-gray-700/40 group"
              >
                <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ background: ep.color }"></span>
                <span class="text-xs text-gray-200">{{ ep.title }}</span>
                <span v-if="ep.status !== 'PLANNING'" class="text-xs px-1 rounded" :class="ep.status === 'DONE' ? 'text-green-400' : 'text-blue-400'">{{ ep.status }}</span>
                <div class="hidden group-hover:flex items-center gap-1 ml-1">
                  <button @click="openEditEpicModal(ep)" class="text-gray-500 hover:text-indigo-400 text-xs">✎</button>
                  <button @click="deleteEpic(ep)" class="text-gray-500 hover:text-red-400 text-xs">✕</button>
                </div>
              </div>
            </div>
            <div v-else class="text-xs text-gray-500 italic">No epics yet. Create one to start organizing your backlog.</div>
          </div>

          <!-- Backlog Table Header + Add Task + Import Slides -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <h3 class="text-base font-semibold text-gray-200">Product Backlog</h3>
              <span class="text-xs text-gray-500">{{ allTasks.filter(t => !t.parent_id).length }} tasks</span>
            </div>
            <div class="flex items-center gap-2">
              <button @click="openBacklogImportModal()" class="btn-import-sm">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
                Import Slides
              </button>
              <button @click="openCreateTaskModal()" class="btn-primary-sm">+ Task</button>
            </div>
          </div>

          <!-- Backlog Table (horizontal scroll on small screens) -->
          <div class="bg-gray-800 border border-gray-700 rounded-xl overflow-x-auto overflow-y-hidden min-w-0">
            <div class="min-w-[640px]">
              <!-- Epic Groups: header = Task, SP, Priority, Sprint, Status (no Epic) -->
              <template v-for="ep in epics" :key="ep.id">
                <!-- Epic Group Header -->
                <div
                  class="flex items-center gap-2 px-3 sm:px-4 py-2 border-b border-gray-700/60 bg-gray-900/40 cursor-pointer hover:bg-gray-900/60 group"
                  @click="toggleEpicGroup(ep.id)"
                >
                  <span class="text-gray-500 text-xs w-4">{{ expandedEpicGroups[ep.id] ? '▼' : '▶' }}</span>
                  <span class="w-3 h-3 rounded-full shrink-0" :style="{ background: ep.color }"></span>
                  <span class="text-sm font-semibold text-gray-200">{{ ep.title }}</span>
                  <span class="text-xs text-gray-500">({{ getTasksForEpic(ep.id).length }} tasks)</span>
                  <div class="ml-auto hidden group-hover:flex items-center gap-2">
                    <button @click.stop="openCreateTaskModal(undefined, ep.id)" class="text-xs text-indigo-400 hover:text-indigo-300">+ Task</button>
                  </div>
                </div>

                <!-- Section header: Epic = Task | SP | Priority | Sprint | Status -->
                <template v-if="expandedEpicGroups[ep.id]">
                  <div class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-2 border-b border-gray-700 text-xs text-gray-500 uppercase tracking-wide">
                    <div class="col-span-1"></div>
                    <div class="col-span-4">Task</div>
                    <div class="col-span-1 text-center">SP</div>
                    <div class="col-span-2">Priority</div>
                    <div class="col-span-2">Sprint</div>
                    <div class="col-span-1">Status</div>
                    <div class="col-span-1"></div>
                  </div>
                  <div v-for="task in getTasksForEpic(ep.id)" :key="task.id">
                    <!-- Task Row (Epic: Sprint only, no Epic dropdown) -->
                    <div class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-3 border-b border-gray-700/50 hover:bg-gray-700/30 transition-colors group">
                      <div class="col-span-1 flex items-center pl-2">
                        <button @click="toggleEpic(task.id)" class="text-gray-500 hover:text-gray-300 text-xs">
                          {{ expandedEpics[task.id] ? '▼' : '▶' }}
                        </button>
                      </div>
                      <div class="col-span-4 flex items-center gap-2 min-w-0">
                        <span class="text-xs font-mono text-gray-600 shrink-0">{{ task.code }}</span>
                        <span class="text-sm font-medium text-gray-200 cursor-pointer hover:text-indigo-300 truncate" @click="navigateToTask(task.id)">{{ task.title }}</span>
                      </div>
                      <div class="col-span-1 flex items-center justify-center">
                        <span class="text-sm font-mono text-purple-400 cursor-pointer hover:text-purple-300" @click="openEditSpField(task)">{{ task.story_points || '–' }}</span>
                      </div>
                      <div class="col-span-2 flex items-center">
                        <select :value="task.priority" @change="updateTaskField(task.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer" :class="priorityTextClass(task.priority)">
                          <option value="CRITICAL">🔴 CRITICAL</option>
                          <option value="HIGH">🟠 HIGH</option>
                          <option value="MEDIUM">🟡 MEDIUM</option>
                          <option value="LOW">🟢 LOW</option>
                        </select>
                      </div>
                      <div class="col-span-2 flex items-center min-w-0">
                        <select :value="task.sprint_id || ''" @change="updateTaskField(task.id, 'sprint_id', ($event.target as HTMLSelectElement).value || null)" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full">
                          <option value="">Backlog</option>
                          <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
                        </select>
                      </div>
                      <div class="col-span-1 flex items-center">
                        <span class="text-xs px-1.5 py-0.5 rounded" :class="taskStatusBadge(task.status)">{{ task.status.replace('_',' ').substring(0,6) }}</span>
                      </div>
                      <div class="col-span-1 flex items-center justify-end opacity-0 group-hover:opacity-100">
                        <button @click="openCreateTaskModal(task.id)" class="text-xs text-indigo-400 hover:text-indigo-300 px-2">+ Sub</button>
                      </div>
                    </div>
                    <!-- Sub-tasks (Epic: inherit Sprint from parent) -->
                    <template v-if="expandedEpics[task.id]">
                      <div v-for="sub in getSubTasks(task.id)" :key="sub.id" class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-2.5 border-b border-gray-700/30 bg-gray-900/30 hover:bg-gray-700/20 transition-colors">
                        <div class="col-span-1"></div>
                        <div class="col-span-4 flex items-center gap-2 pl-6">
                          <span class="text-gray-600">↳</span>
                          <span class="text-xs font-mono text-gray-600">{{ sub.code }}</span>
                          <span class="text-sm text-gray-300 cursor-pointer hover:text-indigo-300" @click="navigateToTask(sub.id)">{{ sub.title }}</span>
                        </div>
                        <div class="col-span-1 flex items-center justify-center">
                          <span class="text-xs font-mono text-purple-400">{{ sub.story_points || '–' }}</span>
                        </div>
                        <div class="col-span-2 flex items-center">
                          <select :value="sub.priority" @change="updateTaskField(sub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer" :class="priorityTextClass(sub.priority)">
                            <option value="CRITICAL">🔴 CRITICAL</option>
                            <option value="HIGH">🟠 HIGH</option>
                            <option value="MEDIUM">🟡 MEDIUM</option>
                            <option value="LOW">🟢 LOW</option>
                          </select>
                        </div>
                        <div class="col-span-2 flex items-center">
                          <span class="text-xs text-gray-500 italic">Inherits from parent</span>
                        </div>
                        <div class="col-span-1 flex items-center">
                          <span class="text-xs px-1.5 py-0.5 rounded" :class="taskStatusBadge(sub.status)">{{ sub.status.replace('_',' ').substring(0,6) }}</span>
                        </div>
                        <div class="col-span-1"></div>
                      </div>
                    </template>
                  </div>
                  <div v-if="!getTasksForEpic(ep.id).length" class="px-8 py-3 text-xs text-gray-500 italic border-b border-gray-700/30 bg-gray-900/20">
                    No tasks in this epic yet.
                    <button @click="openCreateTaskModal(undefined, ep.id)" class="ml-2 text-indigo-400 hover:text-indigo-300">+ Add Task</button>
                  </div>
                </template>
              </template>

              <!-- Unassigned: header = Task, SP, Priority, Epic, Status (no Sprint) -->
              <template v-if="getUnassignedTasks().length">
                <div
                  class="flex items-center gap-2 px-3 sm:px-4 py-2 border-b border-gray-700/60 bg-gray-900/40 cursor-pointer hover:bg-gray-900/60 group"
                  @click="toggleEpicGroup('__unassigned__')"
                >
                  <span class="text-gray-500 text-xs w-4">{{ expandedEpicGroups['__unassigned__'] !== false ? '▼' : '▶' }}</span>
                  <span class="w-3 h-3 rounded-full shrink-0 bg-gray-600"></span>
                  <span class="text-sm font-semibold text-gray-200">Unassigned</span>
                  <span class="text-xs text-gray-500">({{ getUnassignedTasks().length }} tasks)</span>
                </div>
                <template v-if="expandedEpicGroups['__unassigned__'] !== false">
                  <div class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-2 border-b border-gray-700 text-xs text-gray-500 uppercase tracking-wide">
                    <div class="col-span-1"></div>
                    <div class="col-span-4">Task</div>
                    <div class="col-span-1 text-center">SP</div>
                    <div class="col-span-2">Priority</div>
                    <div class="col-span-2">Epic</div>
                    <div class="col-span-1">Status</div>
                    <div class="col-span-1"></div>
                  </div>
                  <div v-for="task in getUnassignedTasks()" :key="task.id">
                    <div class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-3 border-b border-gray-700/50 hover:bg-gray-700/30 transition-colors group">
                      <div class="col-span-1 flex items-center">
                        <button @click="toggleEpic(task.id)" class="text-gray-500 hover:text-gray-300 text-xs">
                          {{ expandedEpics[task.id] ? '▼' : '▶' }}
                        </button>
                      </div>
                      <div class="col-span-4 flex items-center gap-2 min-w-0">
                        <span class="text-xs font-mono text-gray-600 shrink-0">{{ task.code }}</span>
                        <span class="text-sm font-medium text-gray-200 cursor-pointer hover:text-indigo-300 truncate" @click="navigateToTask(task.id)">{{ task.title }}</span>
                      </div>
                      <div class="col-span-1 flex items-center justify-center">
                        <span class="text-sm font-mono text-purple-400 cursor-pointer hover:text-purple-300" @click="openEditSpField(task)">{{ task.story_points || '–' }}</span>
                      </div>
                      <div class="col-span-2 flex items-center">
                        <select :value="task.priority" @change="updateTaskField(task.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer" :class="priorityTextClass(task.priority)">
                          <option value="CRITICAL">🔴 CRITICAL</option>
                          <option value="HIGH">🟠 HIGH</option>
                          <option value="MEDIUM">🟡 MEDIUM</option>
                          <option value="LOW">🟢 LOW</option>
                        </select>
                      </div>
                      <div class="col-span-2 flex items-center min-w-0">
                        <select :value="task.epic_id || ''" @change="updateTaskField(task.id, 'epic_id', ($event.target as HTMLSelectElement).value || '')" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full">
                          <option value="">No Epic</option>
                          <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
                        </select>
                      </div>
                      <div class="col-span-1 flex items-center">
                        <span class="text-xs px-1.5 py-0.5 rounded" :class="taskStatusBadge(task.status)">{{ task.status.replace('_',' ').substring(0,6) }}</span>
                      </div>
                      <div class="col-span-1 flex items-center justify-end opacity-0 group-hover:opacity-100">
                        <button @click="openCreateTaskModal(task.id)" class="text-xs text-indigo-400 hover:text-indigo-300 px-2">+ Sub</button>
                      </div>
                    </div>
                    <template v-if="expandedEpics[task.id]">
                      <div v-for="sub in getSubTasks(task.id)" :key="sub.id" class="grid grid-cols-12 gap-2 px-3 sm:px-4 py-2.5 border-b border-gray-700/30 bg-gray-900/30 hover:bg-gray-700/20 transition-colors">
                        <div class="col-span-1"></div>
                        <div class="col-span-4 flex items-center gap-2 pl-6">
                          <span class="text-gray-600">↳</span>
                          <span class="text-xs font-mono text-gray-600">{{ sub.code }}</span>
                          <span class="text-sm text-gray-300 cursor-pointer hover:text-indigo-300" @click="navigateToTask(sub.id)">{{ sub.title }}</span>
                        </div>
                        <div class="col-span-1 flex items-center justify-center">
                          <span class="text-xs font-mono text-purple-400">{{ sub.story_points || '–' }}</span>
                        </div>
                        <div class="col-span-2 flex items-center">
                          <select :value="sub.priority" @change="updateTaskField(sub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer" :class="priorityTextClass(sub.priority)">
                            <option value="CRITICAL">🔴 CRITICAL</option>
                            <option value="HIGH">🟠 HIGH</option>
                            <option value="MEDIUM">🟡 MEDIUM</option>
                            <option value="LOW">🟢 LOW</option>
                          </select>
                        </div>
                        <div class="col-span-2 flex items-center">
                          <span class="text-xs text-gray-500 italic">Inherits from parent</span>
                        </div>
                        <div class="col-span-1 flex items-center">
                          <span class="text-xs px-1.5 py-0.5 rounded" :class="taskStatusBadge(sub.status)">{{ sub.status.replace('_',' ').substring(0,6) }}</span>
                        </div>
                        <div class="col-span-1"></div>
                      </div>
                    </template>
                  </div>
                </template>
              </template>

              <!-- Empty State -->
              <div v-if="!allTasks.filter(t => !t.parent_id).length" class="py-16 text-center text-gray-500">
                <p class="text-sm mb-3">No tasks in backlog yet.</p>
                <button @click="openCreateTaskModal()" class="btn-primary px-5 py-2 text-sm">Add First Task</button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB 5: Analytics -->
        <div v-if="activeTab === 'analytics'">
          <div v-if="analyticsLoading" class="flex flex-col items-center justify-center py-20">
            <div class="animate-spin text-6xl mb-4">⚙️</div>
            <p class="text-sm text-gray-500">กำลังโหลด analytics...</p>
          </div>
          <ProjectAnalytics v-else-if="analytics" :analytics="analytics" />
          <div v-else class="text-center py-20 text-gray-500 text-sm">Failed to load analytics.</div>
        </div>
      </div>
    </div>

    <!-- Import from Google Slides Modal (Backlog) -->
    <div v-if="showBacklogImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeBacklogImportModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-xl w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-indigo-600/20 border border-indigo-500/30 flex items-center justify-center">
              <svg class="w-4 h-4 text-indigo-400" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import from Google Slides</h2>
              <p class="text-xs text-gray-400">สร้าง task อัตโนมัติจากแต่ละ slide — เลือก Epic ได้</p>
            </div>
          </div>
          <button @click="closeBacklogImportModal" class="text-gray-500 hover:text-white transition-colors">✕</button>
        </div>

        <!-- Result state -->
        <div v-if="backlogImportStep === 'result' && backlogImportResult" class="space-y-4">
          <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
            <div class="flex items-center gap-2 mb-2">
              <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
              <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
            </div>
            <p class="text-gray-300 text-sm font-medium mb-1">{{ backlogImportResult.presentation_title }}</p>
            <p class="text-gray-400 text-xs">สร้าง {{ backlogImportResult.created_count }} tasks จาก {{ backlogImportResult.slide_count }} slides</p>
          </div>
          <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
            <div
              v-for="task in backlogImportResult.tasks"
              :key="task.id"
              class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm"
            >
              <span class="text-xs font-mono text-gray-500 shrink-0">{{ task.code }}</span>
              <span class="text-gray-200 truncate">{{ task.title }}</span>
            </div>
          </div>
          <button @click="closeBacklogImportModal" class="w-full btn-primary py-2.5">Done</button>
        </div>

        <!-- Step 2: Select slides + Epic -->
        <div v-else-if="backlogImportStep === 'select' && backlogImportPreview" class="space-y-4">
          <div class="p-3 bg-gray-700/40 rounded-xl">
            <p class="text-sm font-medium text-white">{{ backlogImportPreview.presentation_title }}</p>
            <p class="text-xs text-gray-500 mt-0.5">
              เลือก slide ที่จะ import ({{ backlogImportSelectedIndices.length }} / {{ backlogImportPreview.slides.length }})
              <span v-if="(backlogImportPreview.already_imported_slide_indices?.length ?? 0) > 0">— หน้าที่นำเข้าแล้วจะถูก uncheck ไว้</span>
            </p>
          </div>
          <div class="flex items-center gap-2 flex-wrap">
            <button type="button" @click="backlogImportSelectAll" class="btn-ghost-sm">เลือกทั้งหมด</button>
            <button type="button" @click="backlogImportDeselectAll" class="btn-ghost-sm">ยกเลิกทั้งหมด</button>
            <button type="button" @click="backlogImportSelectOnlyNew" class="btn-ghost-sm text-indigo-400">เลือกเฉพาะที่ยังไม่เคยนำเข้า</button>
          </div>
          <div class="max-h-56 overflow-y-auto space-y-1.5 pr-1 border border-gray-700/60 rounded-xl p-2">
            <label
              v-for="s in backlogImportPreview.slides"
              :key="s.index"
              class="flex items-center gap-3 py-2 px-2 rounded-lg hover:bg-gray-700/40 cursor-pointer"
              :class="{ 'opacity-70': s.hidden }"
            >
              <input
                v-model="backlogImportSelectedIndices"
                type="checkbox"
                :value="s.index"
                class="rounded border-gray-500 bg-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-xs text-gray-400 w-8 shrink-0">#{{ s.index }}</span>
              <span class="text-sm text-gray-200 truncate flex-1">{{ s.title || '(ไม่มีชื่อ)' }}</span>
              <span v-if="s.hidden" class="text-xs text-amber-400/90 shrink-0">ซ่อน</span>
              <span v-else-if="(backlogImportPreview.already_imported_slide_indices || []).includes(s.index)" class="text-xs text-gray-500 shrink-0">นำเข้าแล้ว</span>
            </label>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Priority ของทุก task</label>
              <select v-model="backlogImportForm.priority" class="input-field w-full" :disabled="isBacklogImporting">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points (ต่อ task)</label>
              <input v-model.number="backlogImportForm.story_points" type="number" min="0" class="input-field w-full" placeholder="1" :disabled="isBacklogImporting" />
            </div>
          </div>
          <div v-if="epics.length">
            <label class="label">นำเข้า Epic ไหน</label>
            <select v-model="backlogImportForm.epic_id" class="input-field w-full" :disabled="isBacklogImporting">
              <option value="">Unassigned (ไม่ใส่ Epic)</option>
              <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
            </select>
          </div>
          <div v-if="backlogImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ backlogImportError }}</div>
          <div class="flex gap-3">
            <button
              @click="submitBacklogImport"
              :disabled="isBacklogImporting || backlogImportSelectedIndices.length === 0"
              class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
            >
              <svg v-if="isBacklogImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
              {{ isBacklogImporting ? 'กำลัง import...' : `Import ${backlogImportSelectedIndices.length} Slides` }}
            </button>
            <button type="button" @click="backlogImportStep = 'form'" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">กลับ</button>
          </div>
        </div>

        <!-- Step 1: Form (URL) -->
        <div v-else class="space-y-4">
          <div>
            <label class="label">Google Slides URL *</label>
            <input
              v-model="backlogImportForm.presentation_url"
              type="url"
              class="input-field w-full"
              placeholder="https://docs.google.com/presentation/d/..."
              :disabled="isBacklogLoadingPreview"
            />
            <p class="text-xs text-gray-500 mt-1">ต้องเปิดสิทธิ์ "Anyone with the link can view"</p>
          </div>
          <div v-if="backlogImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ backlogImportError }}</div>
          <div class="flex gap-3 mt-1">
            <button
              type="button"
              @click="loadBacklogImportPreview"
              :disabled="isBacklogLoadingPreview || !backlogImportForm.presentation_url.trim()"
              class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
            >
              <svg v-if="isBacklogLoadingPreview" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
              {{ isBacklogLoadingPreview ? 'กำลังโหลด...' : 'โหลดรายการ slide' }}
            </button>
            <button @click="closeBacklogImportModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Task Modal -->
    <div v-if="showCreateTaskModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeCreateTaskModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">{{ createTaskForm.parent_id ? 'Add Sub-task' : 'Add Task' }}</h2>
          <button @click="closeCreateTaskModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="label">Title *</label>
            <input v-model="createTaskForm.title" type="text" class="input-field w-full" placeholder="Task title..." />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="createTaskForm.description" rows="3" class="input-field w-full resize-none" placeholder="Describe the task..."></textarea>
          </div>
          <!-- Sub-task hint -->
          <div v-if="createTaskForm.parent_id" class="p-2.5 bg-indigo-900/20 border border-indigo-500/30 rounded-lg text-xs text-indigo-300">
            This is a sub-task. Dates are inherited from the parent task.
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Priority</label>
              <select v-model="createTaskForm.priority" class="input-field w-full">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points</label>
              <input v-model.number="createTaskForm.story_points" type="number" min="0" class="input-field w-full" placeholder="0" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Sprint</label>
              <select v-model="createTaskForm.sprint_id" class="input-field w-full">
                <option value="">Backlog</option>
                <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
              </select>
            </div>
            <div v-if="!createTaskForm.parent_id && epics.length">
              <label class="label">Epic</label>
              <select v-model="createTaskForm.epic_id" class="input-field w-full">
                <option value="">No Epic</option>
                <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
              </select>
            </div>
            <div v-else-if="!createTaskForm.parent_id">
              <label class="label">Due Date</label>
              <input v-model="createTaskForm.due_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <!-- Dates: only shown for top-level tasks (not sub-tasks) -->
          <template v-if="!createTaskForm.parent_id">
            <div v-if="epics.length" class="grid grid-cols-1 gap-3">
              <div>
                <label class="label">Due Date</label>
                <input v-model="createTaskForm.due_date" type="datetime-local" class="input-field w-full" />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="label">Start Date</label>
                <input v-model="createTaskForm.start_date" type="datetime-local" class="input-field w-full" />
              </div>
              <div>
                <label class="label">End Date</label>
                <input v-model="createTaskForm.end_date" type="datetime-local" class="input-field w-full" />
              </div>
            </div>
          </template>
          <div v-if="createTaskError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ createTaskError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="submitCreateTask" :disabled="isCreatingTask || !createTaskForm.title.trim()" class="flex-1 btn-primary py-2.5 disabled:opacity-40">
            {{ isCreatingTask ? 'Creating...' : 'Create Task' }}
          </button>
          <button @click="closeCreateTaskModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Epic Modal (Create / Edit) -->
    <div v-if="showEpicModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeEpicModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">{{ editingEpic ? 'Edit Epic' : 'Create Epic' }}</h2>
          <button @click="closeEpicModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="label">Title *</label>
            <input v-model="epicForm.title" type="text" class="input-field w-full" placeholder="Epic title..." />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="epicForm.description" rows="2" class="input-field w-full resize-none" placeholder="Epic goal or description..."></textarea>
          </div>
          <div>
            <label class="label">Color</label>
            <div class="flex items-center gap-2 flex-wrap mt-1">
              <button
                v-for="c in EPIC_COLORS"
                :key="c"
                @click="epicForm.color = c"
                class="w-7 h-7 rounded-full border-2 transition-all"
                :style="{ background: c }"
                :class="epicForm.color === c ? 'border-white scale-110' : 'border-transparent'"
              ></button>
              <input v-model="epicForm.color" type="color" class="w-7 h-7 rounded cursor-pointer bg-transparent border border-gray-600" title="Custom color" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Start Date</label>
              <input v-model="epicForm.start_date" type="datetime-local" class="input-field w-full" />
            </div>
            <div>
              <label class="label">End Date</label>
              <input v-model="epicForm.end_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <div v-if="epicError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ epicError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="submitEpic" :disabled="isSavingEpic || !epicForm.title.trim()" class="flex-1 btn-primary py-2.5 disabled:opacity-40">
            {{ isSavingEpic ? 'Saving...' : (editingEpic ? 'Update Epic' : 'Create Epic') }}
          </button>
          <button @click="closeEpicModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Sprint Modal (Create / Edit) -->
    <div v-if="showSprintModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeSprintModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">{{ editingSprint ? 'Edit Sprint' : 'Create Sprint' }}</h2>
          <button @click="closeSprintModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="label">Sprint Name *</label>
            <input v-model="sprintForm.name" type="text" class="input-field w-full" placeholder="e.g. Sprint 1 — Foundation" />
          </div>
          <div>
            <label class="label">Sprint Goal</label>
            <textarea v-model="sprintForm.goal" rows="2" class="input-field w-full resize-none" placeholder="What will this sprint achieve?"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Start Date</label>
              <input v-model="sprintForm.start_date" type="datetime-local" class="input-field w-full" />
            </div>
            <div>
              <label class="label">End Date</label>
              <input v-model="sprintForm.end_date" type="datetime-local" class="input-field w-full" />
            </div>
          </div>
          <div v-if="sprintError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ sprintError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="submitSprint" :disabled="isCreatingSprint || !sprintForm.name.trim()" class="flex-1 btn-primary py-2.5 disabled:opacity-40">
            {{ isCreatingSprint ? (editingSprint ? 'Saving...' : 'Creating...') : (editingSprint ? 'Update Sprint' : 'Create Sprint') }}
          </button>
          <button @click="closeSprintModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Add tasks to Sprint Modal -->
    <div v-if="showAddTasksToSprintModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeAddTasksToSprintModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl max-h-[80vh] flex flex-col">
        <h2 class="text-lg font-bold text-white mb-1">เพิ่มงานเข้า Sprint</h2>
        <p class="text-sm text-gray-400 mb-4">เลือกงานที่ต้องการเพิ่มเข้า sprint <strong class="text-white">"{{ sprintForAddTasks?.name }}"</strong></p>
        <div class="flex-1 overflow-y-auto border border-gray-700 rounded-lg p-3 mb-4 min-h-[200px]">
          <div v-if="tasksNotInSprint.length === 0" class="text-center py-8 text-gray-500 text-sm">ไม่มีงานที่ยังไม่อยู่ใน sprint นี้</div>
          <label v-for="t in tasksNotInSprint" :key="t.id" class="flex items-center gap-3 py-2 px-2 rounded-lg hover:bg-gray-700/50 cursor-pointer">
            <input type="checkbox" :value="t.id" v-model="selectedTaskIdsForSprint" class="rounded border-gray-600 bg-gray-700 text-indigo-500 focus:ring-indigo-500" />
            <span class="text-xs font-mono text-gray-500">{{ t.code }}</span>
            <span class="text-sm text-gray-200 truncate flex-1">{{ t.title }}</span>
          </label>
        </div>
        <div v-if="addTasksToSprintError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ addTasksToSprintError }}</div>
        <div class="flex gap-3">
          <button @click="confirmAddTasksToSprint" :disabled="isAddingTasksToSprint || selectedTaskIdsForSprint.length === 0" class="flex-1 btn-primary py-2.5 disabled:opacity-50">
            {{ isAddingTasksToSprint ? 'กำลังเพิ่ม...' : `เพิ่ม ${selectedTaskIdsForSprint.length} งาน` }}
          </button>
          <button @click="closeAddTasksToSprintModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
        </div>
      </div>
    </div>

    <!-- Delete Sprint Confirmation Modal -->
    <div v-if="showDeleteSprintModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeDeleteSprintModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-lg font-bold text-white mb-2">ยืนยันการลบ Sprint</h2>
        <p class="text-sm text-gray-300 mb-2">
          คุณต้องการลบ sprint <strong class="text-white">"{{ sprintToDelete?.name }}"</strong> ใช่หรือไม่?
        </p>
        <p class="text-xs text-gray-500 mb-4">
          งานใน sprint นี้จะถูกย้ายกลับไป Backlog
        </p>
        <div v-if="deleteSprintError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ deleteSprintError }}</div>
        <div class="flex gap-3">
          <button @click="confirmDeleteSprint" :disabled="isDeletingSprint" class="flex-1 px-4 py-2.5 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white font-medium rounded-xl transition-colors">
            {{ isDeletingSprint ? 'กำลังลบ...' : 'ยืนยันการลบ' }}
          </button>
          <button @click="closeDeleteSprintModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
        </div>
      </div>
    </div>

    <!-- Complete Sprint Confirmation Modal -->
    <div v-if="showCompleteSprintModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeCompleteSprintModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-lg font-bold text-white mb-2">ยืนยันการปิด Sprint</h2>
        <p class="text-sm text-gray-300 mb-2">
          คุณต้องการปิด sprint <strong class="text-white">"{{ sprintToComplete?.name }}"</strong> ใช่หรือไม่?
        </p>
        <p class="text-xs text-gray-500 mb-4">
          Sprint จะเปลี่ยนเป็นสถานะ Completed งานที่ยังไม่เสร็จจะยังอยู่ในการอ้างอิงของ sprint นี้
        </p>
        <div v-if="completeSprintError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ completeSprintError }}</div>
        <div class="flex gap-3">
          <button @click="confirmCompleteSprint" :disabled="isCompletingSprint" class="flex-1 px-4 py-2.5 bg-yellow-600 hover:bg-yellow-700 disabled:opacity-50 text-white font-medium rounded-xl transition-colors">
            {{ isCompletingSprint ? 'กำลังปิด...' : 'ยืนยันการปิด' }}
          </button>
          <button @click="closeCompleteSprintModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
        </div>
      </div>
    </div>

    <!-- Reopen Sprint Confirmation Modal -->
    <div v-if="showReopenSprintModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeReopenSprintModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <h2 class="text-lg font-bold text-white mb-2">ยืนยันการเปิด Sprint กลับ</h2>
        <p class="text-sm text-gray-300 mb-2">
          คุณต้องการเปิด sprint <strong class="text-white">"{{ sprintToReopen?.name }}"</strong> กลับเป็น Active ใช่หรือไม่?
        </p>
        <p class="text-xs text-gray-500 mb-4">
          Sprint นี้จะกลายเป็น Current Sprint อีกครั้ง (ถ้ามี sprint อื่นที่ Active อยู่จะเปิดไม่ได้)
        </p>
        <div v-if="reopenSprintError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ reopenSprintError }}</div>
        <div class="flex gap-3">
          <button @click="confirmReopenSprint" :disabled="isReopeningSprint" class="flex-1 px-4 py-2.5 bg-amber-600 hover:bg-amber-700 disabled:opacity-50 text-white font-medium rounded-xl transition-colors">
            {{ isReopeningSprint ? 'กำลังเปิด...' : 'ยืนยันการเปิด' }}
          </button>
          <button @click="closeReopenSprintModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
        </div>
      </div>
    </div>

    <!-- Milestone Modal -->
    <div v-if="showMilestoneModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeMilestoneModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">{{ editingMilestone ? 'Edit Milestone' : 'Add Milestone' }}</h2>
          <button @click="closeMilestoneModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="label">Title *</label>
            <input v-model="milestoneForm.title" type="text" class="input-field w-full" placeholder="e.g. MVP Release" />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="milestoneForm.description" rows="2" class="input-field w-full resize-none"></textarea>
          </div>
          <div>
            <label class="label">Due Date</label>
            <input v-model="milestoneForm.due_date" type="datetime-local" class="input-field w-full" />
          </div>
          <div v-if="editingMilestone">
            <label class="label">Status</label>
            <select v-model="milestoneForm.status" class="input-field w-full">
              <option value="PENDING">Pending</option>
              <option value="REACHED">Reached</option>
              <option value="MISSED">Missed</option>
            </select>
          </div>
          <div v-if="milestoneError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ milestoneError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="submitMilestone" :disabled="isSubmittingMilestone || !milestoneForm.title.trim()" class="flex-1 btn-primary py-2.5 disabled:opacity-40">
            {{ isSubmittingMilestone ? 'Saving...' : editingMilestone ? 'Update' : 'Add Milestone' }}
          </button>
          <button
            v-if="editingMilestone"
            @click="deleteMilestone"
            class="px-4 py-2.5 bg-red-700 hover:bg-red-600 text-white rounded-xl transition-colors text-sm"
          >Delete</button>
          <button @click="closeMilestoneModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import KanbanBoard from '~/components/projects/KanbanBoard.vue'
import GanttMilestoneRow from '~/components/projects/GanttMilestoneRow.vue'
import MilestoneTimeline from '~/components/projects/MilestoneTimeline.vue'
import ProjectAnalytics from '~/components/projects/ProjectAnalytics.vue'
import type { Project, Sprint, Milestone, ProjectAnalytics as AnalyticsType, Task, Epic } from '~/core/modules/projects/infrastructure/projects-api'

definePageMeta({ layout: 'default', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()

const tabs = [
  { id: 'overview', label: 'Overview', icon: '📊' },
  { id: 'board', label: 'Board', icon: '🗂' },
  { id: 'timeline', label: 'Timeline', icon: '📅' },
  { id: 'backlog', label: 'Backlog', icon: '📋' },
  { id: 'analytics', label: 'Analytics', icon: '📈' },
]

const activeTab = ref((route.query.tab as string) || 'overview')

watch(activeTab, async (tab) => {
  router.replace({ query: { tab } })
  if (tab === 'timeline' && project.value) {
    if (timelineMode.value === 'epic' && !epicTimelineData.value) await loadEpicTimeline()
    else if (timelineMode.value === 'sprint' && !sprintTimelineData.value) await loadSprintTimeline()
    nextTick(() => setTimeout(scrollTimelineToToday, 200))
  }
  if (tab === 'analytics' && !analytics.value && project.value) {
    loadAnalytics()
  }
})

// Gantt (Vue-Ganttastic)
interface GanttDependency { id: string; predecessor_id: string; successor_id: string; type: string }
const ganttDependencies = ref<GanttDependency[]>([])
const timelineScrollWrapperRef = ref<HTMLElement | null>(null)
const timelineFilterSprint = ref<string | null>(null)
const timelineFilterMilestone = ref<string | null>(null)

// Matrix Dimension: Timeline view mode — only Epic Roadmap and Sprint Execution (both as Gantt)
type TimelineMode = 'epic' | 'sprint'
const timelineMode = ref<TimelineMode>('epic')
const epicTimelineData = ref<import('~/core/modules/projects/infrastructure/projects-api').EpicTimelineData | null>(null)
const sprintTimelineData = ref<import('~/core/modules/projects/infrastructure/projects-api').SprintTimelineData | null>(null)
const matrixTimelineLoading = ref(false)

watch(timelineMode, async (mode) => {
  if (!project.value) return
  if (mode === 'epic' && !epicTimelineData.value) {
    await loadEpicTimeline()
  } else if (mode === 'sprint' && !sprintTimelineData.value) {
    await loadSprintTimeline()
  }
  nextTick(() => setTimeout(scrollTimelineToToday, 200))
})

const timelineFilteredTasks = computed(() => {
  return allTasks.value.filter((t) => {
    if (timelineFilterSprint.value != null && t.sprint_id !== timelineFilterSprint.value) return false
    if (timelineFilterMilestone.value != null && t.milestone_id !== timelineFilterMilestone.value) return false
    return true
  })
})

function toYMD(d: string) {
  return d.split('T')[0]
}

const ganttChartStart = computed(() => {
  const tasks = timelineFilteredTasks.value.filter((t) => t.start_date || t.end_date || t.due_at)
  const view = ganttView.value
  const padStartDays = 7
  const padStartMonths = 1
  if (!tasks.length) {
    const d = new Date()
    const start = new Date(d)
    if (view === 'month') start.setMonth(start.getMonth() - padStartMonths)
    else start.setDate(start.getDate() - padStartDays)
    return toYMD(start.toISOString())
  }
  let min = Infinity
  for (const t of tasks) {
    const start = t.start_date ? new Date(t.start_date).getTime() : (t.due_at ? new Date(t.due_at).getTime() : null)
    if (start != null) min = Math.min(min, start)
  }
  const minDate = new Date(min === Infinity ? Date.now() : min)
  if (view === 'month') minDate.setMonth(minDate.getMonth() - padStartMonths)
  else minDate.setDate(minDate.getDate() - padStartDays)
  return toYMD(minDate.toISOString())
})

const ganttChartEnd = computed(() => {
  const tasks = timelineFilteredTasks.value.filter((t) => t.start_date || t.end_date || t.due_at)
  const padDays = 30
  const padMonths = ganttView.value === 'month' ? 12 : 3
  if (!tasks.length) {
    const d = new Date()
    d.setDate(d.getDate() + padDays)
    return toYMD(d.toISOString())
  }
  let max = -Infinity
  for (const t of tasks) {
    const end = t.end_date ? new Date(t.end_date).getTime() : (t.due_at ? new Date(t.due_at).getTime() : null)
    if (end != null) max = Math.max(max, end)
  }
  const maxDate = new Date(max === -Infinity ? Date.now() : max)
  const view = ganttView.value
  if (view === 'month') maxDate.setMonth(maxDate.getMonth() + padMonths)
  else maxDate.setDate(maxDate.getDate() + padDays)
  return toYMD(maxDate.toISOString())
})

const ganttPrecision = computed(() => {
  const v = ganttView.value
  if (v === 'month') return 'month'
  if (v === 'week') return 'week'
  return 'day'
})

const ganttRows = computed(() => {
  const today = toYMD(new Date().toISOString())
  const tomorrow = toYMD(new Date(Date.now() + 86400000).toISOString())
  const addDays = (ymd: string, days: number) => {
    const d = new Date(ymd + 'T12:00:00Z')
    d.setUTCDate(d.getUTCDate() + days)
    return toYMD(d.toISOString())
  }
  return timelineFilteredTasks.value.map((t) => {
    let start = t.start_date ? toYMD(t.start_date) : (t.due_at ? toYMD(t.due_at) : today)
    let end = t.end_date ? toYMD(t.end_date) : (t.due_at ? toYMD(t.due_at) : tomorrow)
    if (start === end) end = addDays(start, 1)
    if (end < start) end = addDays(start, 1)
    const label = `${t.code || ''} ${t.title}`.trim() || t.title
    return {
      taskId: t.id,
      label,
      bars: [
        {
          barStart: start,
          barEnd: end,
          ganttBarConfig: {
            id: t.id,
            label: label.length > 40 ? label.slice(0, 37) + '...' : label,
            hasHandles: true,
          },
        },
      ],
    }
  })
})

const ganttChartWidth = computed(() => {
  const start = new Date(ganttChartStart.value + 'T00:00:00Z').getTime()
  const end = new Date(ganttChartEnd.value + 'T00:00:00Z').getTime()
  const days = Math.max(1, (end - start) / 86400000)
  const v = ganttView.value
  const pxPerDay = v === 'month' ? 4 : v === 'week' ? 24 : 40
  return Math.max(800, Math.min(6000, Math.round(days * pxPerDay)))
})

const ganttDateRangeStart = computed(() => {
  const s = ganttChartStart.value
  return s ? new Date(s + 'T00:00:00Z').toISOString() : ''
})

const ganttDateRangeEnd = computed(() => {
  const e = ganttChartEnd.value
  return e ? new Date(e + 'T00:00:00Z').toISOString() : ''
})

function toLocalMidnight(isoOrYmd: string): number {
  const s = String(isoOrYmd).split('T')[0]
  const [y, m, d] = s.split('-').map(Number)
  return new Date(y, m - 1, d).getTime()
}

/** Same scale as Vue-Ganttastic: chart uses full chart div width (ganttChartWidth), position = pct * width; we draw at 220 + that. */
const milestoneLinePositions = computed(() => {
  if (!ganttChartStart.value || !ganttChartEnd.value || ganttChartWidth.value <= 0) return []
  const start = toLocalMidnight(ganttChartStart.value)
  const end = toLocalMidnight(ganttChartEnd.value)
  if (end <= start) return []
  const gridOffset = 220
  const chartContentWidth = ganttChartWidth.value
  const list = milestones.value
    .filter((m): m is Milestone & { due_date: string } => !!m.due_date)
    .sort((a, b) => toLocalMidnight(a.due_date) - toLocalMidnight(b.due_date))
  return list.map((m) => {
    const date = toLocalMidnight(m.due_date)
    const pct = Math.max(0, Math.min(1, (date - start) / (end - start)))
    const left = gridOffset + pct * chartContentWidth
    return { id: m.id, left }
  })
})

// --- Matrix Gantt (Epic Roadmap / Sprint Execution): date range and rows ---
function getMatrixChartRange(): { start: string; end: string } {
  const mode = timelineMode.value
  const pad = 7 * 86400000
  const padEnd = 14 * 86400000
  let min = Infinity
  let max = -Infinity
  if (mode === 'epic' && epicTimelineData.value?.epics?.length) {
    epicTimelineData.value.epics.forEach((ep) => {
      if (ep.start_date) min = Math.min(min, new Date(ep.start_date).getTime())
      if (ep.end_date) max = Math.max(max, new Date(ep.end_date).getTime())
      ;(ep.tasks || []).forEach((t) => {
        if (t.start_date) min = Math.min(min, new Date(t.start_date).getTime())
        if (t.end_date) max = Math.max(max, new Date(t.end_date).getTime())
      })
    })
  } else if (mode === 'sprint' && sprintTimelineData.value?.sprints?.length) {
    sprintTimelineData.value.sprints.forEach((sp) => {
      if (sp.start_date) min = Math.min(min, new Date(sp.start_date).getTime())
      if (sp.end_date) max = Math.max(max, new Date(sp.end_date).getTime())
      ;(sp.tasks || []).forEach((t) => {
        if (t.start_date) min = Math.min(min, new Date(t.start_date).getTime())
        if (t.end_date) max = Math.max(max, new Date(t.end_date).getTime())
      })
    })
  }
  if (min === Infinity) min = Date.now() - pad
  if (max === -Infinity) max = Date.now() + padEnd
  return {
    start: new Date(min - pad).toISOString(),
    end: new Date(max + padEnd).toISOString(),
  }
}

const matrixChartStart = computed(() => toYMD(getMatrixChartRange().start))
const matrixChartEnd = computed(() => toYMD(getMatrixChartRange().end))
const matrixGanttPrecision = computed(() => {
  const v = ganttView.value
  if (v === 'month') return 'month'
  if (v === 'week') return 'week'
  return 'day'
})
const matrixChartWidth = computed(() => {
  const start = new Date(matrixChartStart.value + 'T00:00:00Z').getTime()
  const end = new Date(matrixChartEnd.value + 'T00:00:00Z').getTime()
  const days = Math.max(1, (end - start) / 86400000)
  const v = ganttView.value
  const pxPerDay = v === 'month' ? 4 : v === 'week' ? 24 : 40
  return Math.max(800, Math.min(6000, Math.round(days * pxPerDay)))
})
const matrixDateRangeStart = computed(() => (matrixChartStart.value ? new Date(matrixChartStart.value + 'T00:00:00Z').toISOString() : ''))
const matrixDateRangeEnd = computed(() => (matrixChartEnd.value ? new Date(matrixChartEnd.value + 'T00:00:00Z').toISOString() : ''))

// Timeline matrix: tasks collapsed by default inside each sprint/epic; click bar to expand/collapse
const timelineExpandedSprints = ref<Record<string, boolean>>({})
const timelineExpandedEpics = ref<Record<string, boolean>>({})

const matrixGanttRows = computed(() => {
  const today = toYMD(new Date().toISOString())
  const tomorrow = toYMD(new Date(Date.now() + 86400000).toISOString())
  const addDays = (ymd: string, days: number) => {
    const d = new Date(ymd + 'T12:00:00Z')
    d.setUTCDate(d.getUTCDate() + days)
    return toYMD(d.toISOString())
  }
  const rows: { taskId: string; label: string; bars: { barStart: string; barEnd: string; ganttBarConfig: { id: string; label: string; hasHandles: boolean } }[] }[] = []
  if (timelineMode.value === 'epic' && epicTimelineData.value?.epics?.length) {
    for (const ep of epicTimelineData.value.epics) {
      const epStart = ep.start_date ? toYMD(ep.start_date) : today
      let epEnd = ep.end_date ? toYMD(ep.end_date) : tomorrow
      if (epEnd <= epStart) epEnd = addDays(epStart, 1)
      const taskCount = (ep.tasks || []).length
      const expanded = timelineExpandedEpics.value[ep.id]
      const toggle = expanded ? '▼' : '▶'
      rows.push({
        taskId: `epic-${ep.id}`,
        label: `${toggle} 📁 ${ep.title}${taskCount ? ` (${taskCount})` : ''}`,
        bars: [{ barStart: epStart, barEnd: epEnd, ganttBarConfig: { id: `epic-${ep.id}`, label: ep.title, hasHandles: false } }],
      })
      if (expanded) {
        for (const task of ep.tasks || []) {
          let start = task.start_date ? toYMD(task.start_date) : (task.due_at ? toYMD(task.due_at) : today)
          let end = task.end_date ? toYMD(task.end_date) : (task.due_at ? toYMD(task.due_at) : tomorrow)
          if (end <= start) end = addDays(start, 1)
          const label = `${task.code || ''} ${task.title}`.trim() || task.title
          rows.push({
            taskId: task.id,
            label: `  ${label.length > 35 ? label.slice(0, 32) + '…' : label}`,
            bars: [{ barStart: start, barEnd: end, ganttBarConfig: { id: task.id, label, hasHandles: true } }],
          })
        }
      }
    }
  } else if (timelineMode.value === 'sprint' && sprintTimelineData.value?.sprints?.length) {
    for (const sp of sprintTimelineData.value.sprints) {
      const spStart = sp.start_date ? toYMD(sp.start_date) : today
      let spEnd = sp.end_date ? toYMD(sp.end_date) : tomorrow
      if (spEnd <= spStart) spEnd = addDays(spStart, 1)
      const taskCount = (sp.tasks || []).length
      const expanded = timelineExpandedSprints.value[sp.id]
      const toggle = expanded ? '▼' : '▶'
      rows.push({
        taskId: `sprint-${sp.id}`,
        label: `${toggle} 🏃 ${sp.name}${taskCount ? ` (${taskCount})` : ''}`,
        bars: [{ barStart: spStart, barEnd: spEnd, ganttBarConfig: { id: `sprint-${sp.id}`, label: sp.name, hasHandles: false } }],
      })
      if (expanded) {
        for (const task of sp.tasks || []) {
          let start = task.start_date ? toYMD(task.start_date) : (task.due_at ? toYMD(task.due_at) : today)
          let end = task.end_date ? toYMD(task.end_date) : (task.due_at ? toYMD(task.due_at) : tomorrow)
          if (end <= start) end = addDays(start, 1)
          const label = `${task.code || ''} ${task.title}`.trim() || task.title
          rows.push({
            taskId: task.id,
            label: `  ${label.length > 35 ? label.slice(0, 32) + '…' : label}`,
            bars: [{ barStart: start, barEnd: end, ganttBarConfig: { id: task.id, label, hasHandles: true } }],
          })
        }
      }
    }
  }
  return rows
})

const matrixMilestoneLinePositions = computed(() => {
  if (!matrixChartStart.value || !matrixChartEnd.value || matrixChartWidth.value <= 0) return []
  const start = toLocalMidnight(matrixChartStart.value)
  const end = toLocalMidnight(matrixChartEnd.value)
  if (end <= start) return []
  const gridOffset = 220
  const chartContentWidth = matrixChartWidth.value
  const list = milestones.value
    .filter((m): m is Milestone & { due_date: string } => !!m.due_date)
    .sort((a, b) => toLocalMidnight(a.due_date) - toLocalMidnight(b.due_date))
  return list.map((m) => {
    const date = toLocalMidnight(m.due_date)
    const pct = Math.max(0, Math.min(1, (date - start) / (end - start)))
    const left = gridOffset + pct * chartContentWidth
    return { id: m.id, left }
  })
})

function onMatrixGanttClickBar(payload: { bar: { ganttBarConfig: { id: string } } }) {
  const id = payload?.bar?.ganttBarConfig?.id
  if (!id) return
  if (id.startsWith('epic-')) {
    const epicId = id.slice('epic-'.length)
    timelineExpandedEpics.value = { ...timelineExpandedEpics.value, [epicId]: !timelineExpandedEpics.value[epicId] }
    return
  }
  if (id.startsWith('sprint-')) {
    const sprintId = id.slice('sprint-'.length)
    timelineExpandedSprints.value = { ...timelineExpandedSprints.value, [sprintId]: !timelineExpandedSprints.value[sprintId] }
    return
  }
  router.push(taskUrl(id))
}

/** Click on label column: expand/collapse sprint or epic; navigate for task rows */
function onMatrixLabelClickByLabel(label: string) {
  const row = matrixGanttRows.value.find((r) => r.label === label)
  if (!row) return
  const id = row.taskId
  if (id.startsWith('epic-')) {
    const epicId = id.slice('epic-'.length)
    timelineExpandedEpics.value = { ...timelineExpandedEpics.value, [epicId]: !timelineExpandedEpics.value[epicId] }
    return
  }
  if (id.startsWith('sprint-')) {
    const sprintId = id.slice('sprint-'.length)
    timelineExpandedSprints.value = { ...timelineExpandedSprints.value, [sprintId]: !timelineExpandedSprints.value[sprintId] }
    return
  }
  router.push(taskUrl(id))
}

function onGanttClickBar(payload: { bar: { ganttBarConfig: { id: string } } }) {
  const id = payload?.bar?.ganttBarConfig?.id
  if (id) router.push(taskUrl(id))
}

async function onGanttDragEnd(payload: { movedBars?: Map<string, { start: string | Date; end: string | Date }> }) {
  const map = payload.movedBars
  if (!map || map.size === 0) return
  for (const [taskId, range] of map) {
    const start = typeof range.start === 'string' ? range.start : (range.start as Date).toISOString()
    const end = typeof range.end === 'string' ? range.end : (range.end as Date).toISOString()
    try {
      await tasksApi.updateTask(taskId, {
        start_date: start,
        end_date: end,
      })
      const idx = allTasks.value.findIndex((t) => t.id === taskId)
      if (idx !== -1) {
        allTasks.value[idx] = {
          ...allTasks.value[idx],
          start_date: start,
          end_date: end,
        }
      }
    } catch (e) {
      console.error('Failed to update task dates:', e)
    }
  }
}

function scrollTimelineToToday() {
  nextTick(() => {
    const wrapper = timelineScrollWrapperRef.value
    const width = matrixChartWidth.value
    if (!wrapper || width <= 0) return
    const start = toLocalMidnight(matrixChartStart.value)
    const end = toLocalMidnight(matrixChartEnd.value)
    const now = Date.now()
    const pct = end > start ? Math.max(0, Math.min(1, (now - start) / (end - start))) : 0.5
    const todayOffsetFromLeft = 0.18
    const todayLeftPx = 220 + pct * width
    const targetScroll = Math.max(0, todayLeftPx - wrapper.clientWidth * todayOffsetFromLeft)
    const maxScroll = Math.max(0, wrapper.scrollWidth - wrapper.clientWidth)
    wrapper.scrollLeft = Math.min(targetScroll, maxScroll)
  })
}

// State
const project = ref<Project | null>(null)
const allTasks = ref<Task[]>([])
const sprints = ref<Sprint[]>([])
const milestones = ref<Milestone[]>([])
const epics = ref<Epic[]>([])
const analytics = ref<AnalyticsType | null>(null)
const isLoading = ref(true)
const analyticsLoading = ref(false)
const error = ref('')
const ganttView = ref('week')
const expandedEpics = ref<Record<string, boolean>>({})
const expandedEpicGroups = ref<Record<string, boolean>>({})

// Epic Modal State
const showEpicModal = ref(false)
const editingEpic = ref<Epic | null>(null)
const epicForm = ref({ title: '', description: '', color: '#6366f1', start_date: '', end_date: '' })
const epicError = ref('')
const isSavingEpic = ref(false)
const isDeletingEpic = ref(false)

const EPIC_COLORS = ['#6366f1', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#ef4444', '#06b6d4']

// Computed
const activeSprint = computed(() => sprints.value.find((s) => s.status === 'ACTIVE') ?? null)
/** Sprints for "All sprints" list: ACTIVE first, then others in original order */
const sprintsWithActiveFirst = computed(() =>
  [...sprints.value].sort((a, b) => (a.status === 'ACTIVE' ? -1 : b.status === 'ACTIVE' ? 1 : 0))
)
const totalTasks = computed(() => allTasks.value.length)
const completedCount = computed(() => allTasks.value.filter((t) => t.status === 'COMPLETED').length)
const inProgressCount = computed(() => allTasks.value.filter((t) => t.status === 'IN_PROGRESS').length)
const completionPct = computed(() => totalTasks.value ? Math.round((completedCount.value / totalTasks.value) * 100) : 0)
const overdueCount = computed(() => {
  const now = Date.now()
  return allTasks.value.filter((t) => t.status !== 'COMPLETED' && t.due_at && new Date(t.due_at).getTime() < now).length
})
const epicTasks = computed(() => allTasks.value.filter((t) => !t.parent_id))
const recentTasks = computed(() =>
  [...allTasks.value]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 10)
)

function getSubTasks(parentId: string) {
  return allTasks.value.filter((t) => t.parent_id === parentId)
}

function sprintTaskCount(type: 'total' | 'done' | 'sp') {
  if (!activeSprint.value) return 0
  const tasks = allTasks.value.filter((t) => t.sprint_id === activeSprint.value!.id)
  if (type === 'total') return tasks.length
  if (type === 'done') return tasks.filter((t) => t.status === 'COMPLETED').length
  if (type === 'sp') return tasks.reduce((s, t) => s + (t.story_points || 0), 0)
  return 0
}

function statusClass(status: string) {
  if (status === 'ACTIVE') return 'bg-green-500/10 text-green-400 border-green-500/30'
  if (status === 'COMPLETED') return 'bg-blue-500/10 text-blue-400 border-blue-500/30'
  return 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30'
}

function taskStatusBadge(status: string) {
  if (status === 'COMPLETED') return 'bg-green-500/20 text-green-400'
  if (status === 'IN_PROGRESS') return 'bg-blue-500/20 text-blue-400'
  if (status === 'REVIEW_PENDING') return 'bg-yellow-500/20 text-yellow-400'
  if (status === 'BLOCKED') return 'bg-red-500/20 text-red-400'
  return 'bg-gray-700 text-gray-400'
}

function priorityBadge(p: string) {
  if (p === 'CRITICAL') return 'bg-red-500/20 text-red-400'
  if (p === 'HIGH') return 'bg-orange-500/20 text-orange-400'
  if (p === 'MEDIUM') return 'bg-yellow-500/20 text-yellow-400'
  return 'bg-green-500/20 text-green-400'
}

function priorityTextClass(p: string) {
  if (p === 'CRITICAL') return 'text-red-400'
  if (p === 'HIGH') return 'text-orange-400'
  if (p === 'MEDIUM') return 'text-yellow-400'
  return 'text-green-400'
}

function formatDate(d: string | null) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

function taskUrl(taskId: string) {
  const projectId = route.params.id as string
  const tab = activeTab.value
  return { path: `/task/${taskId}`, query: { from_project: projectId, from_tab: tab } }
}

function navigateToTask(id: string) {
  router.push(taskUrl(id))
}

function toggleEpic(id: string) {
  expandedEpics.value[id] = !expandedEpics.value[id]
}

// Load data
async function loadAll() {
  isLoading.value = true
  error.value = ''
  const idOrCode = route.params.id as string
  try {
    const p = await projectsApi.getProject(idOrCode)
    project.value = p
    const [t, s, m, e] = await Promise.all([
      tasksApi.getTasksByProject(p.id),
      projectsApi.getSprints(p.id),
      projectsApi.getMilestones(p.id),
      projectsApi.getEpics(p.id),
    ])
    allTasks.value = t
    sprints.value = s
    milestones.value = m
    epics.value = e
    // Auto-expand all epic groups
    e.forEach((ep) => { expandedEpicGroups.value[ep.id] = true })
  } catch (e: any) {
    error.value = e.message || 'Failed to load project'
  } finally {
    isLoading.value = false
  }
  if (activeTab.value === 'timeline') {
    if (timelineMode.value === 'epic') await loadEpicTimeline()
    else if (timelineMode.value === 'sprint') await loadSprintTimeline()
  }
  if (activeTab.value === 'analytics') loadAnalytics()
}

async function loadGanttDataForProject() {
  if (!project.value) return
  try {
    const { dependencies } = await tasksApi.getGanttData(project.value.id)
    ganttDependencies.value = (dependencies ?? []) as GanttDependency[]
  } catch {
    ganttDependencies.value = []
  }
}

async function loadAnalytics() {
  if (!project.value) return
  analyticsLoading.value = true
  try {
    analytics.value = await projectsApi.getProjectAnalytics(project.value.id)
  } catch {
    // show empty state
  } finally {
    analyticsLoading.value = false
  }
}

// Kanban status change (uses bulk-status API because PATCH /tasks/:id does not accept status)
async function handleStatusChange(taskId: string, status: string) {
  const idx = allTasks.value.findIndex((t) => t.id === taskId)
  if (idx !== -1) allTasks.value[idx].status = status as Task['status']
  try {
    await tasksApi.bulkUpdateStatus([taskId], status)
  } catch {
    // revert
    await loadAll()
  }
}

// Inline field update (priority, sprint, epic_id)
async function updateTaskField(taskId: string, field: string, value: any) {
  const idx = allTasks.value.findIndex((t) => t.id === taskId)
  if (idx !== -1) (allTasks.value[idx] as any)[field] = value || null
  try {
    const payload = field === 'epic_id' ? { [field]: value ?? '' } : { [field]: value || undefined }
    await tasksApi.updateTask(taskId, payload)
  } catch {
    await loadAll()
  }
}

function openEditSpField(task: Task) {
  const sp = prompt(`Story points for "${task.title}":`, String(task.story_points || 0))
  if (sp !== null && !isNaN(Number(sp))) {
    updateTaskField(task.id, 'story_points', Number(sp))
  }
}

// Sprint operations
const showSprintModal = ref(false)
const editingSprint = ref<Sprint | null>(null)
const sprintForm = ref({ name: '', goal: '', start_date: '', end_date: '' })
const isCreatingSprint = ref(false)
const sprintError = ref('')

function openSprintModal() {
  editingSprint.value = null
  sprintForm.value = { name: '', goal: '', start_date: '', end_date: '' }
  sprintError.value = ''
  showSprintModal.value = true
}

function openEditSprintModal(sprint: Sprint) {
  editingSprint.value = sprint
  sprintForm.value = {
    name: sprint.name,
    goal: sprint.goal || '',
    start_date: sprint.start_date ? new Date(sprint.start_date).toISOString().slice(0, 16) : '',
    end_date: sprint.end_date ? new Date(sprint.end_date).toISOString().slice(0, 16) : '',
  }
  sprintError.value = ''
  showSprintModal.value = true
}

function closeSprintModal() {
  showSprintModal.value = false
  editingSprint.value = null
}

async function submitSprint() {
  if (!project.value) {
    sprintError.value = 'Project not loaded. Please refresh the page.'
    return
  }
  isCreatingSprint.value = true
  sprintError.value = ''
  try {
    const name = sprintForm.value.name.trim()
    const goal = sprintForm.value.goal?.trim() || undefined
    let start_date: string | undefined
    let end_date: string | undefined
    if (sprintForm.value.start_date) {
      const d = new Date(sprintForm.value.start_date)
      if (!isNaN(d.getTime())) start_date = d.toISOString()
    }
    if (sprintForm.value.end_date) {
      const d = new Date(sprintForm.value.end_date)
      if (!isNaN(d.getTime())) end_date = d.toISOString()
    }

    if (editingSprint.value) {
      const updated = await projectsApi.updateSprint(editingSprint.value.id, { name, goal, start_date, end_date })
      const idx = sprints.value.findIndex((s) => s.id === editingSprint.value!.id)
      if (idx !== -1) sprints.value[idx] = updated
      closeSprintModal()
    } else {
      const sprint = await projectsApi.createSprint({
        project_id: project.value.id,
        name,
        goal,
        start_date,
        end_date,
      })
      sprints.value.unshift(sprint)
      closeSprintModal()
    }
  } catch (e: any) {
    const msg = e?.data?.message ?? e?.data?.error ?? e?.message ?? (editingSprint.value ? 'Failed to update sprint' : 'Failed to create sprint')
    sprintError.value = typeof msg === 'string' ? msg : 'Failed to save sprint'
  } finally {
    isCreatingSprint.value = false
  }
}

// Add tasks to Sprint
const showAddTasksToSprintModal = ref(false)
const sprintForAddTasks = ref<Sprint | null>(null)
const selectedTaskIdsForSprint = ref<string[]>([])
const addTasksToSprintError = ref('')
const isAddingTasksToSprint = ref(false)

const tasksNotInSprint = computed(() => {
  const sprintId = sprintForAddTasks.value?.id
  if (!sprintId) return []
  return allTasks.value.filter((t) => t.sprint_id !== sprintId)
})

function openAddTasksToSprintModal(sprint: Sprint) {
  sprintForAddTasks.value = sprint
  selectedTaskIdsForSprint.value = []
  addTasksToSprintError.value = ''
  showAddTasksToSprintModal.value = true
}

function closeAddTasksToSprintModal() {
  showAddTasksToSprintModal.value = false
  sprintForAddTasks.value = null
  selectedTaskIdsForSprint.value = []
  addTasksToSprintError.value = ''
}

async function confirmAddTasksToSprint() {
  if (!sprintForAddTasks.value || selectedTaskIdsForSprint.value.length === 0) return
  isAddingTasksToSprint.value = true
  addTasksToSprintError.value = ''
  try {
    await projectsApi.addTasksToSprint(sprintForAddTasks.value.id, selectedTaskIdsForSprint.value)
    for (const id of selectedTaskIdsForSprint.value) {
      const t = allTasks.value.find((x) => x.id === id)
      if (t) t.sprint_id = sprintForAddTasks.value!.id
    }
    closeAddTasksToSprintModal()
  } catch (e: any) {
    const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'เพิ่มงานไม่สำเร็จ'
    addTasksToSprintError.value = typeof err === 'string' ? err : 'เพิ่มงานไม่สำเร็จ'
  } finally {
    isAddingTasksToSprint.value = false
  }
}

// Delete Sprint: use modal so confirmation always shows (browser may block confirm())
const showDeleteSprintModal = ref(false)
const sprintToDelete = ref<Sprint | null>(null)
const deleteSprintError = ref('')
const isDeletingSprint = ref(false)

function openDeleteSprintModal(sprint: Sprint) {
  sprintToDelete.value = sprint
  deleteSprintError.value = ''
  showDeleteSprintModal.value = true
}

function closeDeleteSprintModal() {
  showDeleteSprintModal.value = false
  sprintToDelete.value = null
  deleteSprintError.value = ''
}

async function confirmDeleteSprint() {
  if (!sprintToDelete.value) return
  isDeletingSprint.value = true
  deleteSprintError.value = ''
  try {
    await projectsApi.deleteSprint(sprintToDelete.value.id)
    sprints.value = sprints.value.filter((s) => s.id !== sprintToDelete.value!.id)
    closeDeleteSprintModal()
  } catch (e: any) {
    const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'ลบไม่สำเร็จ'
    deleteSprintError.value = typeof err === 'string' ? err : 'ลบไม่สำเร็จ'
  } finally {
    isDeletingSprint.value = false
  }
}

// Complete Sprint: use modal so confirmation always shows
const showCompleteSprintModal = ref(false)
const sprintToComplete = ref<Sprint | null>(null)
const completeSprintError = ref('')
const isCompletingSprint = ref(false)

function openCompleteSprintModal(sprint: Sprint) {
  sprintToComplete.value = sprint
  completeSprintError.value = ''
  showCompleteSprintModal.value = true
}

function closeCompleteSprintModal() {
  showCompleteSprintModal.value = false
  sprintToComplete.value = null
  completeSprintError.value = ''
}

async function confirmCompleteSprint() {
  if (!sprintToComplete.value) return
  isCompletingSprint.value = true
  completeSprintError.value = ''
  try {
    const updated = await projectsApi.completeSprint(sprintToComplete.value.id)
    const idx = sprints.value.findIndex((s) => s.id === sprintToComplete.value!.id)
    if (idx !== -1) sprints.value[idx] = updated
    closeCompleteSprintModal()
  } catch (e: any) {
    const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'ปิด sprint ไม่สำเร็จ'
    completeSprintError.value = typeof err === 'string' ? err : 'ปิด sprint ไม่สำเร็จ'
  } finally {
    isCompletingSprint.value = false
  }
}

// Reopen Sprint: use modal so confirmation always shows
const showReopenSprintModal = ref(false)
const sprintToReopen = ref<Sprint | null>(null)
const reopenSprintError = ref('')
const isReopeningSprint = ref(false)

function openReopenSprintModal(sprint: Sprint) {
  sprintToReopen.value = sprint
  reopenSprintError.value = ''
  showReopenSprintModal.value = true
}

function closeReopenSprintModal() {
  showReopenSprintModal.value = false
  sprintToReopen.value = null
  reopenSprintError.value = ''
}

async function confirmReopenSprint() {
  if (!sprintToReopen.value) return
  isReopeningSprint.value = true
  reopenSprintError.value = ''
  try {
    const updated = await projectsApi.reopenSprint(sprintToReopen.value.id)
    const idx = sprints.value.findIndex((s) => s.id === sprintToReopen.value!.id)
    if (idx !== -1) sprints.value[idx] = updated
    closeReopenSprintModal()
  } catch (e: any) {
    const err = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'เปิด sprint กลับไม่สำเร็จ'
    reopenSprintError.value = typeof err === 'string' ? err : 'เปิด sprint กลับไม่สำเร็จ'
  } finally {
    isReopeningSprint.value = false
  }
}

async function handleStartSprint(id: string) {
  if (activeSprint.value) return
  try {
    const updated = await projectsApi.startSprint(id)
    const idx = sprints.value.findIndex((s) => s.id === id)
    if (idx !== -1) sprints.value[idx] = updated
  } catch (e: any) {
    const msg = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'Failed to start sprint'
    const text = typeof msg === 'string' ? msg : 'Failed to start sprint'
    alert(text)
  }
}

// Milestone operations
const showMilestoneModal = ref(false)
const editingMilestone = ref<Milestone | null>(null)
const milestoneForm = ref({ title: '', description: '', due_date: '', status: 'PENDING' })
const isSubmittingMilestone = ref(false)
const milestoneError = ref('')

function openMilestoneModal() {
  editingMilestone.value = null
  milestoneForm.value = { title: '', description: '', due_date: '', status: 'PENDING' }
  milestoneError.value = ''
  showMilestoneModal.value = true
}
function openEditMilestoneModal(m: Milestone) {
  editingMilestone.value = m
  milestoneForm.value = {
    title: m.title,
    description: m.description,
    due_date: m.due_date ? new Date(m.due_date).toISOString().slice(0, 16) : '',
    status: m.status,
  }
  milestoneError.value = ''
  showMilestoneModal.value = true
}
function closeMilestoneModal() {
  showMilestoneModal.value = false
  editingMilestone.value = null
}

async function submitMilestone() {
  if (!project.value) return
  isSubmittingMilestone.value = true
  milestoneError.value = ''
  try {
    if (editingMilestone.value) {
      const updated = await projectsApi.updateMilestone(editingMilestone.value.id, {
        title: milestoneForm.value.title,
        description: milestoneForm.value.description,
        status: milestoneForm.value.status as Milestone['status'],
        due_date: milestoneForm.value.due_date ? new Date(milestoneForm.value.due_date).toISOString() : undefined,
      })
      const idx = milestones.value.findIndex((m) => m.id === editingMilestone.value!.id)
      if (idx !== -1) milestones.value[idx] = updated
    } else {
      const m = await projectsApi.createMilestone({
        project_id: project.value.id,
        title: milestoneForm.value.title,
        description: milestoneForm.value.description,
        due_date: milestoneForm.value.due_date ? new Date(milestoneForm.value.due_date).toISOString() : undefined,
      })
      milestones.value.push(m)
    }
    closeMilestoneModal()
  } catch (e: any) {
    milestoneError.value = e.message
  } finally {
    isSubmittingMilestone.value = false
  }
}

async function deleteMilestone() {
  if (!editingMilestone.value) return
  const name = editingMilestone.value.title
  if (!confirm(`ยืนยันการลบ milestone "${name}"?\n\nกด OK เพื่อลบ / Cancel เพื่อยกเลิก`)) return
  try {
    await projectsApi.deleteMilestone(editingMilestone.value.id)
    milestones.value = milestones.value.filter((m) => m.id !== editingMilestone.value!.id)
    closeMilestoneModal()
  } catch (e: any) {
    milestoneError.value = e.message
  }
}

// Create Task Modal
const showCreateTaskModal = ref(false)
const createTaskForm = ref({
  title: '', description: '', priority: 'MEDIUM', story_points: 0,
  sprint_id: '', due_date: '', start_date: '', end_date: '', parent_id: '', epic_id: ''
})
const isCreatingTask = ref(false)
const createTaskError = ref('')

function openCreateTaskModal(parentId?: string, epicId?: string) {
  createTaskForm.value = { title: '', description: '', priority: 'MEDIUM', story_points: 0, sprint_id: '', due_date: '', start_date: '', end_date: '', parent_id: parentId || '', epic_id: epicId || '' }
  createTaskError.value = ''
  showCreateTaskModal.value = true
}

function closeCreateTaskModal() { showCreateTaskModal.value = false }

// Backlog Import from Google Slides
const showBacklogImportModal = ref(false)
const backlogImportStep = ref<'form' | 'select' | 'result'>('form')
const isBacklogImporting = ref(false)
const isBacklogLoadingPreview = ref(false)
const backlogImportError = ref('')
const backlogImportResult = ref<{ created_count: number; slide_count: number; presentation_title: string; tasks: any[] } | null>(null)
const backlogImportPreview = ref<{
  presentation_title: string
  slides: { index: number; title: string; hidden?: boolean }[]
  already_imported_slide_indices?: number[]
} | null>(null)
const backlogImportSelectedIndices = ref<number[]>([])
const backlogImportForm = ref({
  presentation_url: '',
  priority: 'MEDIUM' as const,
  story_points: 1,
  epic_id: '',
})

function openBacklogImportModal() {
  backlogImportForm.value = { presentation_url: '', priority: 'MEDIUM', story_points: 1, epic_id: '' }
  backlogImportStep.value = 'form'
  backlogImportError.value = ''
  backlogImportResult.value = null
  backlogImportPreview.value = null
  backlogImportSelectedIndices.value = []
  showBacklogImportModal.value = true
}

function closeBacklogImportModal() {
  showBacklogImportModal.value = false
  if (backlogImportResult.value) loadAll()
}

async function loadBacklogImportPreview() {
  if (!backlogImportForm.value.presentation_url.trim()) return
  isBacklogLoadingPreview.value = true
  backlogImportError.value = ''
  try {
    const data = await tasksApi.previewGoogleSlides({
      presentation_url: backlogImportForm.value.presentation_url.trim(),
    })
    backlogImportPreview.value = data
    const alreadySet = new Set(data.already_imported_slide_indices ?? [])
    backlogImportSelectedIndices.value = data.slides
      .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
      .map((s: { index: number }) => s.index)
    backlogImportStep.value = 'select'
  } catch (e: any) {
    backlogImportError.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
  } finally {
    isBacklogLoadingPreview.value = false
  }
}

function backlogImportSelectAll() {
  if (backlogImportPreview.value) backlogImportSelectedIndices.value = backlogImportPreview.value.slides.map((s) => s.index)
}

function backlogImportDeselectAll() {
  backlogImportSelectedIndices.value = []
}

function backlogImportSelectOnlyNew() {
  if (!backlogImportPreview.value) return
  const alreadySet = new Set(backlogImportPreview.value.already_imported_slide_indices ?? [])
  backlogImportSelectedIndices.value = backlogImportPreview.value.slides
    .filter((s: { index: number; hidden?: boolean }) => !s.hidden && !alreadySet.has(s.index))
    .map((s: { index: number }) => s.index)
}

async function submitBacklogImport() {
  if (!project.value) return
  isBacklogImporting.value = true
  backlogImportError.value = ''
  try {
    const payload: any = {
      presentation_url: backlogImportForm.value.presentation_url.trim(),
      project_id: project.value.id,
      priority: backlogImportForm.value.priority,
      story_points: backlogImportForm.value.story_points,
    }
    if (backlogImportForm.value.epic_id) payload.epic_id = backlogImportForm.value.epic_id
    if (backlogImportSelectedIndices.value.length > 0) payload.slide_indices = backlogImportSelectedIndices.value
    backlogImportResult.value = await tasksApi.importGoogleSlides(payload)
    backlogImportStep.value = 'result'
  } catch (e: any) {
    backlogImportError.value = e?.data?.message ?? e?.message ?? 'Import failed'
  } finally {
    isBacklogImporting.value = false
  }
}

async function submitCreateTask() {
  if (!project.value) return
  isCreatingTask.value = true
  createTaskError.value = ''
  try {
    const payload: any = {
      title: createTaskForm.value.title,
      description: createTaskForm.value.description,
      priority: createTaskForm.value.priority,
      story_points: createTaskForm.value.story_points,
      project_id: project.value.id,
    }
    if (createTaskForm.value.parent_id) payload.parent_id = createTaskForm.value.parent_id
    if (createTaskForm.value.epic_id) payload.epic_id = createTaskForm.value.epic_id
    if (createTaskForm.value.sprint_id) payload.sprint_id = createTaskForm.value.sprint_id
    if (createTaskForm.value.due_date) payload.due_date = new Date(createTaskForm.value.due_date).toISOString()
    // Sub-tasks inherit dates from parent — don't send dates when parent_id is set
    if (!createTaskForm.value.parent_id) {
      if (createTaskForm.value.start_date) payload.start_date = new Date(createTaskForm.value.start_date).toISOString()
      if (createTaskForm.value.end_date) payload.end_date = new Date(createTaskForm.value.end_date).toISOString()
    }
    const task = await tasksApi.createTask(payload)
    allTasks.value.unshift(task)
    closeCreateTaskModal()
  } catch (e: any) {
    createTaskError.value = e.message
  } finally {
    isCreatingTask.value = false
  }
}

// Matrix Dimension Timeline loaders
async function loadEpicTimeline() {
  if (!project.value) return
  matrixTimelineLoading.value = true
  try {
    epicTimelineData.value = await projectsApi.getEpicTimelineData(project.value.id)
  } catch (e) {
    console.error('Failed to load epic timeline:', e)
  } finally {
    matrixTimelineLoading.value = false
  }
}

async function loadSprintTimeline() {
  if (!project.value) return
  matrixTimelineLoading.value = true
  try {
    sprintTimelineData.value = await projectsApi.getSprintTimelineData(project.value.id)
  } catch (e) {
    console.error('Failed to load sprint timeline:', e)
  } finally {
    matrixTimelineLoading.value = false
  }
}

// Epic Management
function openCreateEpicModal() {
  editingEpic.value = null
  epicForm.value = { title: '', description: '', color: '#6366f1', start_date: '', end_date: '' }
  epicError.value = ''
  showEpicModal.value = true
}

function openEditEpicModal(epic: Epic) {
  editingEpic.value = epic
  epicForm.value = {
    title: epic.title,
    description: epic.description || '',
    color: epic.color || '#6366f1',
    start_date: epic.start_date ? epic.start_date.slice(0, 16) : '',
    end_date: epic.end_date ? epic.end_date.slice(0, 16) : '',
  }
  epicError.value = ''
  showEpicModal.value = true
}

function closeEpicModal() { showEpicModal.value = false }

async function submitEpic() {
  if (!project.value) return
  isSavingEpic.value = true
  epicError.value = ''
  try {
    const payload: any = {
      title: epicForm.value.title,
      description: epicForm.value.description,
      color: epicForm.value.color,
    }
    if (epicForm.value.start_date) payload.start_date = new Date(epicForm.value.start_date).toISOString()
    if (epicForm.value.end_date) payload.end_date = new Date(epicForm.value.end_date).toISOString()

    if (editingEpic.value) {
      const updated = await projectsApi.updateEpic(editingEpic.value.id, payload)
      const idx = epics.value.findIndex((e) => e.id === updated.id)
      if (idx >= 0) epics.value[idx] = updated
    } else {
      payload.project_id = project.value.id
      const created = await projectsApi.createEpic(payload)
      epics.value.push(created)
      expandedEpicGroups.value[created.id] = true
    }
    closeEpicModal()
  } catch (e: any) {
    epicError.value = e.message || 'Failed to save epic'
  } finally {
    isSavingEpic.value = false
  }
}

async function deleteEpic(epic: Epic) {
  if (!confirm(`Delete epic "${epic.title}"? Tasks in this epic will be unlinked.`)) return
  isDeletingEpic.value = true
  try {
    await projectsApi.deleteEpic(epic.id)
    epics.value = epics.value.filter((e) => e.id !== epic.id)
    // Unlink tasks locally
    allTasks.value = allTasks.value.map((t) => t.epic_id === epic.id ? { ...t, epic_id: null } : t)
  } catch (e: any) {
    alert(e.message || 'Failed to delete epic')
  } finally {
    isDeletingEpic.value = false
  }
}

function toggleEpicGroup(id: string) {
  expandedEpicGroups.value[id] = !expandedEpicGroups.value[id]
}

function getTasksForEpic(epicId: string) {
  return allTasks.value.filter((t) => t.epic_id === epicId && !t.parent_id)
}

function getUnassignedTasks() {
  return allTasks.value.filter((t) => !t.epic_id && !t.parent_id)
}

onMounted(loadAll)
</script>

<style scoped>
.card {
  @apply bg-gray-800 border border-gray-700 rounded-xl p-5;
}
.metric-card {
  @apply bg-gray-800/60 border border-gray-700/50 rounded-xl p-4;
}
.metric-label {
  @apply text-xs text-gray-500 mt-1 uppercase tracking-wide;
}
.section-title {
  @apply text-sm font-semibold text-gray-300;
}
.label {
  @apply block text-xs text-gray-400 mb-1.5 font-medium;
}
.input-field {
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-colors;
}
.btn-primary {
  @apply bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-xl transition-colors;
}
.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-indigo-600 hover:bg-indigo-700 text-white font-medium rounded-lg transition-colors;
}
.btn-import-sm {
  @apply px-3 py-1.5 text-xs bg-indigo-900/50 hover:bg-indigo-800/60 border border-indigo-700/50 text-indigo-300 font-medium rounded-lg transition-colors flex items-center gap-1.5;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
/* Timeline tab layout */
.timeline-tab {
  --gantt-bar: 99 102 241; /* indigo-500 */
  --gantt-bar-hover: 129 140 248; /* indigo-400 */
  --gantt-today: 96 165 250; /* blue-400 */
}

.gantt-chart-vue {
  min-height: 420px;
  font-family: ui-sans-serif, system-ui, sans-serif;
}

.timeline-scroll-wrapper {
  overscroll-behavior-x: contain;
}

.milestone-legend-diamond {
  display: inline-block;
  width: 10px;
  height: 10px;
}

/* Enterprise overrides for Vue-Ganttastic (dark theme) */
.gantt-enterprise :deep(.g-gantt-chart) {
  background: rgb(30 41 59);
  border-radius: 0 0 0.75rem 0;
}

.gantt-enterprise :deep(.g-timeaxis) {
  background: rgb(15 23 42) !important;
  border-bottom: 1px solid rgb(51 65 85);
  height: 72px;
}

.gantt-enterprise :deep(.g-timeunits-container) {
  color: rgb(148 163 184);
  font-size: 0.7rem;
  font-weight: 500;
  letter-spacing: 0.02em;
}

.gantt-enterprise :deep(.g-upper-timeunit) {
  color: rgb(203 213 225);
  font-size: 0.75rem;
  font-weight: 600;
}

.gantt-enterprise :deep(.g-label-column) {
  background: rgb(15 23 42);
  color: rgb(203 213 225);
  border-right: 1px solid rgb(51 65 85);
  border-radius: 0.75rem 0 0 0;
}

.gantt-enterprise :deep(.g-label-column-header) {
  background: rgb(15 23 42);
  color: rgb(148 163 184);
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  border-bottom: 1px solid rgb(51 65 85);
  height: 72px;
  min-height: 72px;
}

.gantt-enterprise :deep(.g-label-column-row) {
  color: rgb(226 232 240);
  font-size: 0.8rem;
  font-weight: 500;
  padding: 0 0.75rem;
  border-bottom: 1px solid rgb(51 65 85 / 0.6);
}

.gantt-enterprise :deep(.g-gantt-row > .g-gantt-row-bars-container) {
  border-top: 1px solid rgb(51 65 85 / 0.5);
  border-bottom: 1px solid rgb(51 65 85 / 0.5);
  background: rgb(30 41 59 / 0.5);
}

.gantt-enterprise :deep(.g-gantt-row:nth-child(even) > .g-gantt-row-bars-container) {
  background: rgb(30 41 59 / 0.8);
}

.gantt-enterprise :deep(.g-grid-line) {
  border-left: 1px solid rgb(51 65 85 / 0.7);
}

.gantt-enterprise :deep(.g-gantt-bar) {
  background: linear-gradient(135deg, rgb(99 102 241) 0%, rgb(79 70 229) 100%) !important;
  border-radius: 6px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.2);
  border: 1px solid rgb(99 102 241 / 0.4);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.gantt-enterprise :deep(.g-gantt-bar:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgb(99 102 241 / 0.35);
}

.gantt-enterprise :deep(.g-gantt-bar-label) {
  color: rgb(255 255 255);
  font-size: 0.8rem;
  font-weight: 500;
  text-shadow: 0 1px 1px rgb(0 0 0 / 0.3);
}

.gantt-enterprise :deep(.g-gantt-bar-handle-left),
.gantt-enterprise :deep(.g-gantt-bar-handle-right) {
  background: rgb(255 255 255 / 0.25) !important;
  width: 8px;
  border-radius: 4px 0 0 4px;
}

.gantt-enterprise :deep(.g-gantt-bar-handle-right) {
  border-radius: 0 4px 4px 0;
}

.gantt-enterprise :deep(.g-grid-current-time-marker) {
  width: 2px !important;
  background: rgb(96 165 250) !important;
  box-shadow: 0 0 8px rgb(96 165 250 / 0.8);
}

.gantt-enterprise :deep(.g-grid-current-time-text) {
  color: rgb(96 165 250);
  font-size: 0.65rem;
  font-weight: 600;
}

.gantt-enterprise :deep(.g-gantt-tooltip) {
  background: rgb(30 41 59) !important;
  border: 1px solid rgb(71 85 105);
  border-radius: 8px;
  padding: 0.5rem 0.75rem;
  font-size: 0.8rem;
  box-shadow: 0 10px 25px rgb(0 0 0 / 0.4);
}

.gantt-enterprise :deep(.g-gantt-tooltip:before) {
  border-bottom-color: rgb(30 41 59);
}
</style>
