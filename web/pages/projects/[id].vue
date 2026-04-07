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
            <button
              type="button"
              @click="openEditProjectModal"
              class="shrink-0 p-1.5 rounded-lg text-gray-500 hover:text-purple-400 hover:bg-gray-700/50 transition-colors"
              title="แก้ไขชื่อโครงการ"
            >
              <span class="text-sm">✏️</span>
            </button>
            <span
              class="px-2 py-0.5 text-xs font-semibold rounded-full border shrink-0"
              :class="statusClass(project.status)"
            >
              {{ project.status.replace('_', ' ') }}
            </span>
            <code class="text-xs text-gray-500 font-mono hidden md:inline shrink-0">{{ project.code }}</code>
          </div>
          <button
            v-if="activeTab === 'timeline'"
            type="button"
            class="flex items-center gap-2 rounded-lg border border-slate-500/50 bg-slate-700/50 px-3 py-1.5 text-xs font-medium text-slate-300 transition-colors hover:bg-slate-600/60 hover:text-white disabled:opacity-60 shrink-0"
            :disabled="timelineRefreshing"
            title="โหลดข้อมูลใหม่ / Refresh"
            @click="refreshTimeline"
          >
            <span v-if="timelineRefreshing" class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-slate-400 border-t-transparent" aria-hidden="true" />
            <span v-else aria-hidden="true">↻</span>
            Refresh
          </button>
        </div>

        <!-- Tabs (horizontal scroll on small screens) -->
        <div class="flex gap-1 mt-3 sm:mt-4 overflow-x-auto pb-1 -mx-1 px-1 scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-transparent">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            class="px-3 sm:px-4 py-2 text-xs sm:text-sm rounded-lg transition-colors font-medium whitespace-nowrap shrink-0"
            :class="activeTab === tab.id
              ? 'bg-gradient-to-r from-purple-600 to-pink-600 text-white'
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
              <div class="text-2xl font-bold" :class="activeSprint ? 'text-purple-400' : 'text-gray-500'">
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
                    <div class="text-sm font-bold text-purple-400">{{ sprintTaskCount('sp') }}</div>
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
                      class="h-full bg-purple-500 rounded-full"
                      :style="{ width: Math.round(sprintTaskCount('done') / sprintTaskCount('total') * 100) + '%' }"
                    ></div>
                  </div>
                </div>
              </div>
              <div v-else class="text-center py-8 text-gray-500 text-sm">
                No active sprint. Plan and start a sprint to begin tracking.
              </div>
              <!-- List of all sprints (so user sees where created sprints are); drag to reorder -->
              <div v-if="sprints.length > 0" class="mt-4 pt-4 border-t border-gray-700">
                <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">All sprints</h4>
                <ul class="space-y-2 max-h-40 overflow-y-auto">
                  <li
                    v-for="(s, sIdx) in sprintsOrdered"
                    :key="s.id"
                    class="flex items-center justify-between py-1.5 px-2 rounded-lg group"
                    :class="[s.status === 'ACTIVE' ? 'bg-purple-500/10' : 'hover:bg-gray-700/40', sprintDragId === s.id && 'opacity-60']"
                    draggable="true"
                    @dragstart="onSprintDragStart($event, s.id)"
                    @dragover="onSprintDragOver"
                    @drop.stop="onSprintDrop($event, sIdx)"
                  >
                    <span class="text-gray-500 text-xs w-4 shrink-0 cursor-grab select-none" title="ลากเพื่อเรียงลำดับ">⋮⋮</span>
                    <NuxtLink
                      :to="`/projects/sprint/${s.id}?project=${route.params.id}`"
                      class="text-sm text-gray-200 hover:text-purple-300 transition-colors truncate flex-1 min-w-0 mr-2"
                    >
                      {{ s.name }}
                    </NuxtLink>
                    <span class="flex items-center gap-2">
                      <span class="text-[10px] px-1.5 py-0.5 rounded font-medium" :class="s.status === 'ACTIVE' ? 'bg-purple-500/20 text-purple-400' : s.status === 'COMPLETED' ? 'bg-gray-600 text-gray-400' : 'bg-yellow-500/20 text-yellow-400'">
                        {{ s.status }}
                      </span>
                      <button
                        v-if="s.status === 'PLANNING'"
                        type="button"
                        :disabled="!!activeSprint"
                        @click.stop="!activeSprint && handleStartSprint(s.id)"
                        class="text-xs font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        :class="activeSprint ? 'text-gray-500' : 'text-purple-400 hover:text-purple-300'"
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
            :task-display-code-map="taskDisplayCodeMap"
            :user-role="currentUser?.role"
            :active-sprint="activeSprint"
            :deployed-task-ids="deployedTaskIds"
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
                    <button v-for="v in ['Day', 'Week', 'Month']" :key="v" type="button" @click="ganttView = v.toLowerCase()" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="ganttView === v.toLowerCase() ? (timelineMode === 'epic' ? 'bg-purple-600 text-white shadow-sm' : 'bg-emerald-600 text-white shadow-sm') : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'" :title="v === 'Week' ? '1 column = 7 days (Mon–Sun)' : v === 'Day' ? '1 column = 1 day' : '1 column = 1 month'">
                      {{ v === 'Week' ? 'Week (7d)' : v }}
                    </button>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <div v-if="matrixGanttRows.length > 0" class="flex items-center rounded-lg border border-slate-600/60 bg-slate-800/50 overflow-hidden">
                  <button type="button" class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-slate-300 hover:text-white hover:bg-slate-700/60 transition-colors border-r border-slate-600/60" title="กางทั้งหมด" @click="expandAllTimelineTasks">
                    Expand all
                  </button>
                  <button type="button" class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-slate-300 hover:text-white hover:bg-slate-700/60 transition-colors" title="ย่อทั้งหมด" @click="collapseAllTimelineTasks">
                    Collapse all
                  </button>
                </div>
                <button v-if="matrixGanttRows.length > 0" type="button" class="flex items-center gap-2 rounded-lg border border-slate-500/50 bg-slate-700/50 px-3 py-1.5 text-xs font-medium text-slate-300 transition-colors hover:bg-slate-600/60 hover:text-white" @click="timelineFullscreen = true" title="ขยายเต็มจอ / Fullscreen">
                  <span aria-hidden="true">⛶</span> Fullscreen
                </button>
                <button v-if="matrixGanttRows.length > 0" type="button" class="flex items-center gap-2 rounded-lg border border-slate-500/50 bg-slate-700/50 px-3 py-1.5 text-xs font-medium text-slate-300 transition-colors hover:bg-slate-600/60 hover:text-white" title="Export timeline as PDF (opens in new tab)" @click="onExportTimelinePdf">
                  <span aria-hidden="true">📄</span> Export PDF
                </button>
                <button type="button" class="flex items-center gap-2 rounded-lg border border-purple-500/50 bg-purple-600/20 px-3 py-1.5 text-xs font-medium text-purple-300 transition-colors hover:bg-purple-600/40 hover:text-purple-200" @click="scrollTimelineToToday">
                  <span aria-hidden="true">◉</span> Today
                </button>
              </div>
            </div>
          </div>

          <!-- Dynamic epic bar colors (so each epic uses its chosen color on Gantt) -->
          <component :is="'style'" v-if="epicBarStyles">{{ epicBarStyles }}</component>

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
            <div ref="timelineScrollWrapperRef" class="timeline-scroll-wrapper cursor-grab rounded-xl border border-slate-600/60 bg-slate-800/60 shadow-xl shadow-black/25 overflow-x-auto overflow-y-hidden active:cursor-grabbing" :class="{ 'timeline-fullscreen': timelineFullscreen, 'fixed inset-0 z-50 m-0 rounded-none flex flex-col overflow-hidden': timelineFullscreen }" @mousedown="onTimelinePanStart" @touchstart.passive="onTimelinePanStartTouch">
              <!-- Exit fullscreen bar (only when expanded) -->
              <div v-if="timelineFullscreen" class="flex shrink-0 items-center justify-between gap-4 border-b border-slate-600/60 bg-slate-800/95 px-4 py-2 shadow-md">
                <span class="text-sm font-medium text-slate-300">Timeline — Fullscreen</span>
                <div class="flex flex-wrap items-center gap-4">
                  <!-- Mode: Epic Roadmap | Sprint Execution -->
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">Mode</span>
                    <div class="flex rounded-lg bg-slate-900/80 p-0.5">
                      <button type="button" @click="timelineMode = 'epic'" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="timelineMode === 'epic' ? 'bg-purple-600 text-white shadow-sm' : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'">
                        Epic Roadmap
                      </button>
                      <button type="button" @click="timelineMode = 'sprint'" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="timelineMode === 'sprint' ? 'bg-emerald-600 text-white shadow-sm' : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'">
                        Sprint Execution
                      </button>
                    </div>
                  </div>
                  <div class="h-4 w-px bg-slate-600" />
                  <!-- View: Day | Week (7 days) | Month -->
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">View</span>
                    <div class="flex rounded-lg bg-slate-900/80 p-0.5">
                      <button v-for="v in ['Day', 'Week', 'Month']" :key="v" type="button" @click="ganttView = v.toLowerCase()" class="rounded-md px-3 py-1.5 text-xs font-medium transition-all duration-200" :class="ganttView === v.toLowerCase() ? (timelineMode === 'epic' ? 'bg-purple-600 text-white shadow-sm' : 'bg-emerald-600 text-white shadow-sm') : 'text-slate-400 hover:bg-slate-700/60 hover:text-slate-200'" :title="v === 'Week' ? '1 column = 7 days (Mon–Sun)' : v === 'Day' ? '1 column = 1 day' : '1 column = 1 month'">
                        {{ v === 'Week' ? 'Week (7d)' : v }}
                      </button>
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <div class="flex items-center rounded-lg border border-slate-600/60 bg-slate-800/50 overflow-hidden">
                    <button type="button" class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-slate-300 hover:text-white hover:bg-slate-700/60 transition-colors border-r border-slate-600/60" title="กางทั้งหมด" @click="expandAllTimelineTasks">
                      Expand all
                    </button>
                    <button type="button" class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-slate-300 hover:text-white hover:bg-slate-700/60 transition-colors" title="ย่อทั้งหมด" @click="collapseAllTimelineTasks">
                      Collapse all
                    </button>
                  </div>
                  <button type="button" class="flex items-center gap-2 rounded-lg border border-purple-500/50 bg-purple-600/20 px-3 py-1.5 text-xs font-medium text-purple-300 transition-colors hover:bg-purple-600/40 hover:text-purple-200" @click="scrollTimelineToToday">
                    <span aria-hidden="true">◉</span> Today
                  </button>
                  <button type="button" class="flex items-center gap-2 rounded-lg border border-slate-500/50 bg-slate-700/50 px-3 py-1.5 text-xs font-medium text-slate-300 transition-colors hover:bg-slate-600/60 hover:text-white" title="Export timeline as PDF (opens in new tab)" @click="onExportTimelinePdf">
                    <span aria-hidden="true">📄</span> Export PDF
                  </button>
                  <button type="button" class="flex items-center gap-2 rounded-lg border border-slate-500 bg-slate-700 px-3 py-1.5 text-xs font-medium text-slate-200 hover:bg-slate-600 hover:text-white" @click="timelineFullscreen = false" title="ย่อกลับ / Exit fullscreen">
                    ✕ ย่อกลับ
                  </button>
                </div>
              </div>
              <!-- พื้นที่เลื่อนแนวตั้ง+แนวนอน เมื่อ fullscreen (min-h-0 ให้ flex ลูก scroll ได้) -->
              <div :class="timelineFullscreen ? 'timeline-fullscreen-scroll min-h-0 flex-1 overflow-auto' : ''">
              <div class="timeline-inner relative flex flex-col" :style="matrixChartWidth > 0 ? { width: (220 + matrixChartWidth) + 'px', minWidth: (220 + matrixChartWidth) + 'px' } : { minWidth: '100%' }">
                <GanttMilestoneRow v-if="matrixDateRangeStart && matrixDateRangeEnd && matrixChartWidth > 0" :milestones="milestones" :date-range-start="matrixDateRangeStart" :date-range-end="matrixDateRangeEnd" :grid-width="matrixChartWidth" :grid-offset="220" @milestone-click="openEditMilestoneModal" @milestone-drag-move="onMilestoneDragMove" @milestone-drag-end="onMilestoneDragEnd" />
                <g-gantt-chart :chart-start="matrixChartStart" :chart-end="matrixChartEnd" :precision="matrixGanttPrecision" bar-start="barStart" bar-end="barEnd" date-format="YYYY-MM-DD" :width="matrixChartWidth + 'px'" :row-height="52" :grid="true" :current-time="true" current-time-label="Now" color-scheme="dark" :label-column-title="timelineMode === 'epic' ? 'Epic / Task' : 'Sprint / Task'" label-column-width="220px" class="gantt-chart-vue gantt-enterprise" @click-bar="onMatrixGanttClickBar" @dragstart-bar="onMatrixGanttDragStart" @dragend-bar="onGanttDragEnd" @mouseenter-bar="onMatrixGanttMouseEnter" @mouseleave-bar="onMatrixGanttMouseLeave">
                  <template #timeunit="{ label, date }">
                    <span v-if="ganttView === 'week' && date" class="whitespace-nowrap">{{ weekRangeLabel(date) }}</span>
                    <span v-else>{{ label }}</span>
                  </template>
                  <template #label-column-row="{ label }">
                    <span class="cursor-pointer w-full block min-w-0 whitespace-normal break-words text-[13px] leading-tight py-0.5" @click.stop="onMatrixLabelClickByLabel(label)">{{ label }}</span>
                  </template>
                  <template #bar-tooltip="{ bar }">
                    <div v-if="bar" class="rounded-md border border-slate-600 bg-slate-900/95 px-3 py-2 text-xs text-slate-200 shadow-lg">
                      <p class="font-semibold text-white mb-1">{{ bar.ganttBarConfig?.label || 'Timeline item' }}</p>
                      <p>
                        {{ formatGanttTooltipDate(bar.barStart) }} → {{ formatGanttTooltipInclusiveEnd(bar.barEnd) }}
                      </p>
                    </div>
                  </template>
                  <g-gantt-row v-for="row in matrixGanttRows" :key="row.taskId" :label="row.label" :bars="row.bars" :class="row.taskId.startsWith('epic-') ? 'gantt-row-epic' : row.taskId.startsWith('sprint-') ? 'gantt-row-sprint' : 'gantt-row-task'" />
                </g-gantt-chart>
                <div v-if="matrixMilestoneLinePositions.length > 0" class="pointer-events-none absolute inset-0 z-[5]" aria-hidden="true">
                  <div v-for="{ id, left } in matrixMilestoneLinePositions" :key="id" class="absolute top-0 bottom-0 w-px bg-purple-500/50" :style="{ left: left + 'px' }" />
                </div>
              </div>
              </div>
              <div
                v-if="matrixHoveredBar"
                class="pointer-events-none fixed z-[80] rounded-md border border-slate-600 bg-slate-900/95 px-3 py-2 text-xs text-slate-200 shadow-xl"
                :style="{ left: matrixTooltipPos.x + 'px', top: matrixTooltipPos.y + 'px' }"
              >
                <p class="font-semibold text-white mb-1">{{ matrixHoveredBar?.ganttBarConfig?.label || 'Timeline item' }}</p>
                <p>{{ formatGanttTooltipDate(matrixHoveredBar?.barStart) }} → {{ formatGanttTooltipInclusiveEnd(matrixHoveredBar?.barEnd) }}</p>
              </div>
            </div>
            <template #fallback>
              <div class="flex min-h-[420px] flex-col items-center justify-center rounded-xl border border-slate-600/50 bg-slate-800/50">
                <div class="h-8 w-8 animate-spin rounded-full border-2 border-slate-500 border-t-purple-400" />
                <p class="mt-3 text-xs font-medium text-slate-400">Loading timeline…</p>
              </div>
            </template>
          </ClientOnly>

          <!-- Milestone legend (when chart is shown and we have milestones) -->
          <div v-if="milestones.length && matrixGanttRows.length > 0" class="rounded-xl border border-slate-600/40 bg-slate-800/40 px-4 py-3 mt-4">
            <p class="mb-2 text-xs font-semibold uppercase tracking-wider text-slate-500">Milestones</p>
            <div class="flex flex-wrap gap-x-6 gap-y-2">
              <div v-for="m in milestones" :key="m.id" class="flex items-center gap-2">
                <span class="milestone-legend-diamond rotate-45 border-2 border-purple-400/80 bg-slate-800" />
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
                v-for="(ep, epIdx) in epics"
                :key="ep.id"
                draggable="true"
                class="flex items-center gap-1.5 px-2.5 py-1 rounded-lg border border-gray-600/50 bg-gray-700/40 group cursor-grab active:cursor-grabbing"
                :class="{ 'opacity-60': backlogDrag?.type === 'epic' && backlogDrag?.id === ep.id }"
                @dragstart="onEpicDragStart($event, ep.id)"
                @dragover="onEpicDragOver"
                @drop="onEpicDrop($event, epIdx)"
              >
                <span class="text-gray-500 cursor-grab shrink-0 select-none" title="ลากเพื่อเรียงลำดับ">⋮⋮</span>
                <span class="w-2.5 h-2.5 rounded-full shrink-0" :style="{ background: ep.color }"></span>
                <span class="text-xs text-gray-200">{{ ep.title }}</span>
                <span v-if="ep.status !== 'PLANNING'" class="text-xs px-1 rounded" :class="ep.status === 'DONE' ? 'text-green-400' : 'text-blue-400'">{{ ep.status }}</span>
                <div class="hidden group-hover:flex items-center gap-1 ml-1">
                  <button type="button" @click.stop="openEditEpicModal(ep)" class="text-gray-500 hover:text-purple-400 text-xs">✎</button>
                  <button type="button" @click.stop="deleteEpic(ep)" class="text-gray-500 hover:text-red-400 text-xs">✕</button>
                </div>
              </div>
            </div>
            <div v-else class="text-xs text-gray-500 italic">No epics yet. Create one to start organizing your backlog.</div>
          </div>

          <!-- Backlog Table Header + Add Task + Import Slides -->
          <div class="flex flex-wrap items-center justify-between gap-3">
            <div class="flex items-center gap-3 flex-wrap">
              <h3 class="text-base font-semibold text-gray-200">Product Backlog</h3>
              <span class="text-xs text-gray-500">{{ allTasks.filter(t => !t.parent_id).length }} tasks</span>
              <span v-if="projectDetailsTasksMeta?.has_more" class="text-[11px] text-amber-300 bg-amber-500/10 border border-amber-500/30 rounded px-2 py-0.5">
                Showing first {{ projectDetailsTasksMeta.returned }} tasks (server-limited)
              </span>
              <span
                v-if="projectDetailsTasksMeta?.has_more && isLoadingMoreProjectTasks"
                class="text-[11px] text-blue-300 bg-blue-500/10 border border-blue-500/30 rounded px-2 py-0.5"
              >
                Auto loading more tasks…
              </span>
            </div>
            <div class="flex items-center gap-2">
              <!-- Selection toolbar when items selected -->
              <template v-if="backlogSelectedCount > 0">
                <span class="text-xs text-gray-400">{{ backlogSelectedCount }} selected</span>
                <button type="button" @click="clearBacklogSelection" class="px-3 py-2 text-xs font-medium text-gray-400 hover:text-white rounded-lg border border-gray-600 hover:bg-gray-700/60 transition-colors">
                  Clear
                </button>
                <button
                  type="button"
                  @click="bulkDeleteSelectedBacklogTasks"
                  :disabled="isBulkDeletingBacklog"
                  class="px-3 py-2 text-xs font-medium text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 rounded-lg transition-colors flex items-center gap-1.5"
                >
                  <span v-if="isBulkDeletingBacklog" class="w-3.5 h-3.5 animate-spin rounded-full border-2 border-white border-t-transparent" />
                  Delete selected
                </button>
              </template>
              <!-- Expand / Collapse all (enterprise toolbar) -->
              <div class="flex items-center rounded-lg border border-gray-600/60 bg-gray-800/50 overflow-hidden">
                <button
                  type="button"
                  @click="expandAllBacklog"
                  class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-gray-300 hover:text-white hover:bg-gray-700/60 transition-colors border-r border-gray-600/60"
                  title="กางทั้งหมด"
                >
                  <svg class="w-3.5 h-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                  </svg>
                  <span class="hidden sm:inline">Expand all</span>
                </button>
                <button
                  type="button"
                  @click="collapseAllBacklog"
                  class="inline-flex items-center gap-1.5 px-3 py-2 text-xs font-medium text-gray-300 hover:text-white hover:bg-gray-700/60 transition-colors"
                  title="ย่อทั้งหมด"
                >
                  <svg class="w-3.5 h-3.5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                  </svg>
                  <span class="hidden sm:inline">Collapse all</span>
                </button>
              </div>
              <button @click="openBacklogImportModal()" class="btn-import-sm">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
                Import Slides
              </button>
              <button @click="openSheetsImportModal()" class="btn-import-sm">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd"/></svg>
                Import EU
              </button>
              <button @click="openIODImportModal()" class="btn-import-sm border border-blue-500/40 text-blue-300 hover:bg-blue-500/10">
                <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
                Import IOD
              </button>
              <button @click="openPPTXImportModal()" class="btn-import-sm">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm4 18H6V4h7v5h5v11zM8 15h8v2H8v-2zm0-4h8v2H8v-2z"/></svg>
                Import PPTX
              </button>
              <button @click="openCreateTaskModal()" class="btn-primary-sm">+ Task</button>
            </div>
          </div>

          <!-- Backlog Table (horizontal scroll on small screens) -->
          <div class="bg-gray-800 border border-gray-700 rounded-xl overflow-x-auto overflow-y-hidden min-w-0">
            <div class="min-w-[640px]">
              <!-- Epic Groups: header = Task, SP, Priority, Sprint, Status (no Epic) -->
              <template v-for="(ep, epIdx) in epics" :key="ep.id">
                <!-- Epic Group Header (draggable to reorder) -->
                <div
                  class="flex items-center gap-2 px-3 sm:px-4 py-2 border-b border-gray-700/60 bg-gray-900/40 cursor-pointer hover:bg-gray-900/60 group"
                  :class="{ 'opacity-60': backlogDrag?.type === 'epic' && backlogDrag?.id === ep.id }"
                  draggable="true"
                  @click="toggleEpicGroup(ep.id)"
                  @dragstart="onEpicDragStart($event, ep.id)"
                  @dragover="onEpicDragOver"
                  @drop.stop="onEpicDrop($event, epIdx)"
                >
                  <span class="text-gray-500 text-xs w-4 shrink-0 cursor-grab select-none" title="ลากเพื่อเรียงลำดับ">⋮⋮</span>
                  <span class="text-gray-500 text-xs w-4">{{ expandedEpicGroups[ep.id] ? '▼' : '▶' }}</span>
                  <span class="w-3 h-3 rounded-full shrink-0" :style="{ background: ep.color }"></span>
                  <span class="text-sm font-semibold text-gray-200">{{ ep.title }}</span>
                  <span class="text-xs text-gray-500">({{ getTasksForEpic(ep.id).length }} tasks)</span>
                  <div class="ml-auto hidden group-hover:flex items-center gap-2">
                    <button type="button" @click.stop="openCreateTaskModal(undefined, ep.id)" class="text-xs text-purple-400 hover:text-purple-300">+ Task</button>
                  </div>
                </div>

                <!-- Section header + rows in one grid so columns align (subgrid) -->
                <template v-if="expandedEpicGroups[ep.id]">
                  <div class="backlog-table-grid">
                    <div class="backlog-table-header backlog-subgrid">
                      <div class="flex items-center justify-center shrink-0">
                        <input
                          type="checkbox"
                          :checked="backlogSectionSelectAllChecked(ep.id)"
                          :indeterminate.prop="backlogSectionSelectAllIndeterminate(ep.id)"
                          class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                          @change="setBacklogSectionSelectAll(ep.id, ($event.target as HTMLInputElement).checked)"
                        />
                      </div>
                      <div class="flex items-center justify-center shrink-0"></div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">ID</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Task</div>
                      <div class="flex items-center justify-center shrink-0 font-semibold text-gray-300">SP</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Priority</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Epic</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Sprint</div>
                      <div class="flex items-center justify-center shrink-0 font-semibold text-gray-300">Status</div>
                      <div class="flex items-center justify-center w-min shrink-0"></div>
                    </div>
                    <template v-for="(task, taskIdx) in getTasksForEpic(ep.id)" :key="task.id">
                      <div
                        class="backlog-row backlog-subgrid group"
                        :class="{ 'opacity-60': backlogDrag?.type === 'task' && backlogDrag?.id === task.id }"
                        @dragover="onTaskDragOver"
                        @drop.stop="onTaskDrop($event, ep.id, taskIdx)"
                      >
                        <div class="flex items-center justify-center shrink-0">
                          <input
                            type="checkbox"
                            :checked="isBacklogTaskSelected(task.id)"
                            class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                            @change="toggleBacklogTaskSelection(task.id)"
                            @click.stop
                          />
                        </div>
                        <div class="flex items-center gap-3 shrink-0">
                          <span
                          class="text-gray-500 cursor-grab select-none text-xs"
                          title="ลากเพื่อเรียงลำดับ"
                          draggable="true"
                          @dragstart="onTaskDragStartSetData($event, task.id, ep.id)"
                        >⋮⋮</span>
                        <button type="button" @click="toggleEpic(task.id)" class="text-gray-500 hover:text-gray-300 text-xs shrink-0">
                          {{ expandedEpics[task.id] ? '▼' : '▶' }}
                        </button>
                      </div>
                      <div class="flex items-center min-w-0">
                        <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(task)">{{ taskDisplayCode(task) }}</span>
                      </div>
                      <div class="flex items-center gap-1 min-w-0">
                        <span
                          class="shrink-0 text-xs font-bold"
                          :class="task.task_type === 'FEATURE' ? 'text-purple-400' : task.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                          :title="task.task_type"
                        >{{ task.task_type === 'FEATURE' ? '★' : task.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                        <span class="text-sm font-medium text-gray-200 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0 flex-1" :title="task.title" @click="navigateToTask(task.id)">{{ task.title }}</span>
                        <span class="shrink-0 flex items-center gap-0.5 ml-auto opacity-0 group-hover:opacity-100 transition-opacity">
                          <button type="button" @click.stop="openEditTaskTitle(task)" class="p-0.5 rounded text-gray-500 hover:text-purple-400 hover:bg-gray-700/50" title="แก้ไขชื่อ task">✎</button>
                          <button type="button" @click.stop="duplicateTask(task)" class="p-0.5 rounded text-gray-500 hover:text-purple-400 hover:bg-gray-700/50" title="Duplicate task">⎘</button>
                        </span>
                      </div>
                      <div class="flex items-center justify-center shrink-0">
                        <span class="text-sm font-mono text-purple-400 cursor-pointer hover:text-purple-300" @click="openEditSpField(task)">{{ task.story_points || '–' }}</span>
                      </div>
                      <div class="flex items-center min-w-0">
                        <select :value="task.priority" @change="updateTaskField(task.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(task.priority)">
                          <option value="CRITICAL">🔴 CRITICAL</option>
                          <option value="HIGH">🟠 HIGH</option>
                          <option value="MEDIUM">🟡 MEDIUM</option>
                          <option value="LOW">🟢 LOW</option>
                        </select>
                      </div>
                      <div class="flex items-center min-w-0">
                        <select :value="task.epic_id || ''" @change="updateTaskField(task.id, 'epic_id', ($event.target as HTMLSelectElement).value || '')" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full" title="ย้ายไป Epic อื่น">
                          <option value="">No Epic</option>
                          <option v-for="e in epics" :key="e.id" :value="e.id">{{ e.title }}</option>
                        </select>
                      </div>
                      <div class="flex items-center min-w-0">
                        <select :value="task.sprint_id || ''" @change="updateTaskField(task.id, 'sprint_id', ($event.target as HTMLSelectElement).value || null)" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full">
                          <option value="">Backlog</option>
                          <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
                        </select>
                      </div>
                      <div class="flex items-center shrink-0">
                        <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(task.status)">{{ task.status.replace('_',' ') }}</span>
                      </div>
                      <div class="flex items-center justify-start w-full min-w-0 shrink-0 opacity-0 group-hover:opacity-100">
                        <button @click="openCreateTaskModal(task.id)" class="text-xs text-purple-400 hover:text-purple-300 shrink-0 py-0.5">+ Sub</button>
                      </div>
                    </div>
                    <!-- Sub-tasks B (inherit Epic from parent); C = sub-tasks of B -->
                    <template v-if="expandedEpics[task.id]">
                      <template v-for="sub in getSubTasks(task.id)" :key="sub.id">
                        <div class="backlog-subgrid backlog-sub-row border-b border-gray-700/40 bg-gray-900/55 hover:bg-gray-700/35 transition-colors group">
                          <div class="flex items-center justify-center shrink-0">
                            <input
                              type="checkbox"
                              :checked="isBacklogTaskSelected(sub.id)"
                              class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                              @change="toggleBacklogTaskSelection(sub.id)"
                              @click.stop
                            />
                          </div>
                          <div class="flex items-center pl-6">
                            <button v-if="getSubTasks(sub.id).length" type="button" @click.stop="toggleEpic(sub.id)" class="text-gray-500 hover:text-gray-300 text-xs shrink-0" :aria-label="expandedEpics[sub.id] ? 'Collapse' : 'Expand'">{{ expandedEpics[sub.id] ? '▼' : '▶' }}</button>
                          </div>
                          <div class="flex items-center min-w-0 pl-6">
                            <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(sub)">{{ taskDisplayCode(sub) }}</span>
                          </div>
                          <div class="flex items-center gap-1 min-w-0">
                            <span class="text-gray-600 shrink-0">↳</span>
                            <span
                              class="shrink-0 text-xs font-bold"
                              :class="sub.task_type === 'FEATURE' ? 'text-purple-400' : sub.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                              :title="sub.task_type"
                            >{{ sub.task_type === 'FEATURE' ? '★' : sub.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                            <span class="text-sm text-gray-300 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0" :title="sub.title" @click="navigateToTask(sub.id)">{{ sub.title }}</span>
                          </div>
                          <div class="flex items-center justify-center shrink-0">
                            <span class="text-xs font-mono text-purple-400">{{ sub.story_points || '–' }}</span>
                          </div>
                          <div class="flex items-center min-w-0">
                            <select :value="sub.priority" @change="updateTaskField(sub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(sub.priority)">
                              <option value="CRITICAL">🔴 CRITICAL</option>
                              <option value="HIGH">🟠 HIGH</option>
                              <option value="MEDIUM">🟡 MEDIUM</option>
                              <option value="LOW">🟢 LOW</option>
                            </select>
                          </div>
                          <div class="flex items-center justify-center min-w-0 w-full">
                            <span class="text-xs text-gray-500 italic">Inherits</span>
                          </div>
                          <div class="flex items-center justify-center min-w-0 w-full">
                            <span class="text-xs text-gray-500 italic">Inherits</span>
                          </div>
                          <div class="flex items-center shrink-0">
                            <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(sub.status)">{{ sub.status.replace('_',' ') }}</span>
                          </div>
                          <div class="flex items-center justify-start w-full min-w-0 shrink-0 opacity-0 group-hover:opacity-100">
                            <button type="button" @click.stop="openCreateTaskModal(sub.id)" class="text-xs text-purple-400 hover:text-purple-300 shrink-0 py-0.5">+ Sub</button>
                          </div>
                        </div>
                        <!-- Level C: sub-tasks of B -->
                        <template v-if="expandedEpics[sub.id]">
                          <div v-for="subsub in getSubTasks(sub.id)" :key="subsub.id" class="backlog-subgrid backlog-sub-row border-b border-gray-700/20 bg-gray-950/50 hover:bg-gray-700/10 transition-colors">
                            <div class="flex items-center justify-center shrink-0">
                              <input
                                type="checkbox"
                                :checked="isBacklogTaskSelected(subsub.id)"
                                class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                                @change="toggleBacklogTaskSelection(subsub.id)"
                                @click.stop
                              />
                            </div>
                            <div class="flex items-center"></div>
                            <div class="flex items-center min-w-0 pl-10">
                              <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(subsub)">{{ taskDisplayCode(subsub) }}</span>
                            </div>
                            <div class="flex items-center gap-1 min-w-0">
                              <span class="text-gray-600 shrink-0">↳↳</span>
                              <span
                                class="shrink-0 text-xs font-bold"
                                :class="subsub.task_type === 'FEATURE' ? 'text-purple-400' : subsub.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                                :title="subsub.task_type"
                              >{{ subsub.task_type === 'FEATURE' ? '★' : subsub.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                              <span class="text-sm text-gray-400 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0" :title="subsub.title" @click="navigateToTask(subsub.id)">{{ subsub.title }}</span>
                            </div>
                            <div class="flex items-center justify-center shrink-0">
                              <span class="text-xs font-mono text-purple-400">{{ subsub.story_points || '–' }}</span>
                            </div>
                            <div class="flex items-center min-w-0">
                              <select :value="subsub.priority" @change="updateTaskField(subsub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(subsub.priority)">
                                <option value="CRITICAL">🔴 CRITICAL</option>
                                <option value="HIGH">🟠 HIGH</option>
                                <option value="MEDIUM">🟡 MEDIUM</option>
                                <option value="LOW">🟢 LOW</option>
                              </select>
                            </div>
                            <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                            <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                            <div class="flex items-center shrink-0">
                              <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(subsub.status)">{{ subsub.status.replace('_',' ') }}</span>
                            </div>
                            <div class="flex items-center w-min shrink-0"></div>
                          </div>
                        </template>
                      </template>
                    </template>
                    </template>
                  </div>
                  <div v-if="!getTasksForEpic(ep.id).length" class="px-8 py-3 text-xs text-gray-500 italic border-b border-gray-700/30 bg-gray-900/20">
                    No tasks in this epic yet.
                    <button @click="openCreateTaskModal(undefined, ep.id)" class="ml-2 text-purple-400 hover:text-purple-300">+ Add Task</button>
                  </div>
                </template>
              </template>

              <!-- Unassigned: header = Task, SP, Priority, Epic, Status (no Sprint) -->
              <template v-if="getUnassignedTasks().length">
                <button
                  type="button"
                  class="relative z-10 flex w-full items-center gap-2 px-3 sm:px-4 py-2 border-b border-gray-700/60 bg-gray-900/40 cursor-pointer hover:bg-gray-900/60 group text-left"
                  @click.stop="toggleEpicGroup('__unassigned__')"
                >
                  <span class="text-gray-500 text-xs w-4">{{ expandedEpicGroups['__unassigned__'] !== false ? '▼' : '▶' }}</span>
                  <span class="w-3 h-3 rounded-full shrink-0 bg-gray-600"></span>
                  <span class="text-sm font-semibold text-gray-200">Unassigned</span>
                  <span class="text-xs text-gray-500">({{ getUnassignedTasks().length }} tasks)</span>
                </button>
                <template v-if="expandedEpicGroups['__unassigned__'] !== false">
                  <div class="backlog-table-grid">
                    <div class="backlog-table-header backlog-subgrid">
                      <div class="flex items-center justify-center shrink-0">
                        <input
                          type="checkbox"
                          :checked="backlogSectionSelectAllChecked('__unassigned__')"
                          :indeterminate.prop="backlogSectionSelectAllIndeterminate('__unassigned__')"
                          class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                          @change="setBacklogSectionSelectAll('__unassigned__', ($event.target as HTMLInputElement).checked)"
                        />
                      </div>
                      <div class="flex items-center justify-center shrink-0"></div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">ID</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Task</div>
                      <div class="flex items-center justify-center shrink-0 font-semibold text-gray-300">SP</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Priority</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Epic</div>
                      <div class="flex items-center justify-center min-w-0 font-semibold text-gray-300">Sprint</div>
                      <div class="flex items-center justify-center shrink-0 font-semibold text-gray-300">Status</div>
                      <div class="flex items-center justify-center w-min shrink-0"></div>
                    </div>
                    <template v-for="(task, taskIdx) in getUnassignedTasks()" :key="task.id">
                      <div
                        class="backlog-row backlog-subgrid group"
                        :class="{ 'opacity-60': backlogDrag?.type === 'task' && backlogDrag?.id === task.id }"
                        @dragover="onTaskDragOver"
                        @drop.stop="onTaskDrop($event, null, taskIdx)"
                      >
                        <div class="flex items-center justify-center shrink-0">
                          <input
                            type="checkbox"
                            :checked="isBacklogTaskSelected(task.id)"
                            class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                            @change="toggleBacklogTaskSelection(task.id)"
                            @click.stop
                          />
                        </div>
                        <div class="flex items-center gap-3 shrink-0">
                          <span
                            class="text-gray-500 cursor-grab select-none text-xs"
                            title="ลากเพื่อเรียงลำดับ"
                            draggable="true"
                            @dragstart="onTaskDragStartSetData($event, task.id, null)"
                          >⋮⋮</span>
                          <button type="button" @click="toggleEpic(task.id)" class="text-gray-500 hover:text-gray-300 text-xs shrink-0">
                            {{ expandedEpics[task.id] ? '▼' : '▶' }}
                          </button>
                        </div>
                        <div class="flex items-center min-w-0">
                          <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(task)">{{ taskDisplayCode(task) }}</span>
                        </div>
                        <div class="flex items-center gap-1 min-w-0">
                          <span
                            class="shrink-0 text-xs font-bold"
                            :class="task.task_type === 'FEATURE' ? 'text-purple-400' : task.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                            :title="task.task_type"
                          >{{ task.task_type === 'FEATURE' ? '★' : task.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                          <span class="text-sm font-medium text-gray-200 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0 flex-1" :title="task.title" @click="navigateToTask(task.id)">{{ task.title }}</span>
                          <span class="shrink-0 flex items-center gap-0.5 ml-auto opacity-0 group-hover:opacity-100 transition-opacity">
                            <button type="button" @click.stop="openEditTaskTitle(task)" class="p-0.5 rounded text-gray-500 hover:text-purple-400 hover:bg-gray-700/50" title="แก้ไขชื่อ task">✎</button>
                            <button type="button" @click.stop="duplicateTask(task)" class="p-0.5 rounded text-gray-500 hover:text-purple-400 hover:bg-gray-700/50" title="Duplicate task">⎘</button>
                          </span>
                        </div>
                        <div class="flex items-center justify-center shrink-0">
                          <span class="text-sm font-mono text-purple-400 cursor-pointer hover:text-purple-300" @click="openEditSpField(task)">{{ task.story_points || '–' }}</span>
                        </div>
                        <div class="flex items-center min-w-0">
                          <select :value="task.priority" @change="updateTaskField(task.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(task.priority)">
                            <option value="CRITICAL">🔴 CRITICAL</option>
                            <option value="HIGH">🟠 HIGH</option>
                            <option value="MEDIUM">🟡 MEDIUM</option>
                            <option value="LOW">🟢 LOW</option>
                          </select>
                        </div>
                        <div class="flex items-center min-w-0">
                          <select :value="task.epic_id || ''" @change="updateTaskField(task.id, 'epic_id', ($event.target as HTMLSelectElement).value || '')" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full">
                            <option value="">No Epic</option>
                            <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
                          </select>
                        </div>
                        <div class="flex items-center min-w-0">
                          <select :value="task.sprint_id || ''" @change="updateTaskField(task.id, 'sprint_id', ($event.target as HTMLSelectElement).value || null)" class="text-xs bg-gray-700 border border-gray-600 rounded px-1.5 py-0.5 text-gray-300 focus:outline-none max-w-full">
                            <option value="">Backlog</option>
                            <option v-for="s in sprints" :key="s.id" :value="s.id">{{ s.name }}</option>
                          </select>
                        </div>
                        <div class="flex items-center shrink-0">
                          <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(task.status)">{{ task.status.replace('_',' ') }}</span>
                        </div>
                        <div class="flex items-center justify-start w-full min-w-0 shrink-0 opacity-0 group-hover:opacity-100">
                          <button @click="openCreateTaskModal(task.id)" class="text-xs text-purple-400 hover:text-purple-300 shrink-0 py-0.5">+ Sub</button>
                        </div>
                      </div>
                      <template v-if="expandedEpics[task.id]">
                        <template v-for="sub in getSubTasks(task.id)" :key="sub.id">
                          <div class="backlog-subgrid backlog-sub-row border-b border-gray-700/40 bg-gray-900/55 hover:bg-gray-700/35 transition-colors group">
                            <div class="flex items-center justify-center shrink-0">
                              <input
                                type="checkbox"
                                :checked="isBacklogTaskSelected(sub.id)"
                                class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                                @change="toggleBacklogTaskSelection(sub.id)"
                                @click.stop
                              />
                            </div>
                            <div class="flex items-center pl-6">
                              <button v-if="getSubTasks(sub.id).length" type="button" @click.stop="toggleEpic(sub.id)" class="text-gray-500 hover:text-gray-300 text-xs shrink-0" :aria-label="expandedEpics[sub.id] ? 'Collapse' : 'Expand'">{{ expandedEpics[sub.id] ? '▼' : '▶' }}</button>
                            </div>
                            <div class="flex items-center min-w-0 pl-6">
                              <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(sub)">{{ taskDisplayCode(sub) }}</span>
                            </div>
                            <div class="flex items-center gap-1 min-w-0">
                              <span class="text-gray-600 shrink-0">↳</span>
                              <span
                                class="shrink-0 text-xs font-bold"
                                :class="sub.task_type === 'FEATURE' ? 'text-purple-400' : sub.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                                :title="sub.task_type"
                              >{{ sub.task_type === 'FEATURE' ? '★' : sub.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                              <span class="text-sm text-gray-300 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0" :title="sub.title" @click="navigateToTask(sub.id)">{{ sub.title }}</span>
                            </div>
                            <div class="flex items-center justify-center shrink-0">
                              <span class="text-xs font-mono text-purple-400">{{ sub.story_points || '–' }}</span>
                            </div>
                            <div class="flex items-center min-w-0">
                              <select :value="sub.priority" @change="updateTaskField(sub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(sub.priority)">
                                <option value="CRITICAL">🔴 CRITICAL</option>
                                <option value="HIGH">🟠 HIGH</option>
                                <option value="MEDIUM">🟡 MEDIUM</option>
                                <option value="LOW">🟢 LOW</option>
                              </select>
                            </div>
                            <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                            <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                            <div class="flex items-center shrink-0">
                              <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(sub.status)">{{ sub.status.replace('_',' ') }}</span>
                            </div>
                            <div class="flex items-center justify-start w-full min-w-0 shrink-0 opacity-0 group-hover:opacity-100">
                              <button type="button" @click.stop="openCreateTaskModal(sub.id)" class="text-xs text-purple-400 hover:text-purple-300 shrink-0 py-0.5">+ Sub</button>
                            </div>
                          </div>
                          <!-- Level C: sub-tasks of B (Unassigned) -->
                          <template v-if="expandedEpics[sub.id]">
                            <div v-for="subsub in getSubTasks(sub.id)" :key="subsub.id" class="backlog-subgrid backlog-sub-row border-b border-gray-700/20 bg-gray-950/50 hover:bg-gray-700/10 transition-colors">
                              <div class="flex items-center justify-center shrink-0">
                                <input
                                  type="checkbox"
                                  :checked="isBacklogTaskSelected(subsub.id)"
                                  class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                                  @change="toggleBacklogTaskSelection(subsub.id)"
                                  @click.stop
                                />
                              </div>
                              <div class="flex items-center"></div>
                              <div class="flex items-center min-w-0 pl-10">
                                <span class="text-xs font-mono text-gray-500 truncate" :title="taskDisplayCode(subsub)">{{ taskDisplayCode(subsub) }}</span>
                              </div>
                              <div class="flex items-center gap-1 min-w-0">
                                <span class="text-gray-600 shrink-0">↳↳</span>
                                <span
                                  class="shrink-0 text-xs font-bold"
                                  :class="subsub.task_type === 'FEATURE' ? 'text-purple-400' : subsub.task_type === 'BUG' ? 'text-red-400' : 'text-blue-400'"
                                  :title="subsub.task_type"
                                >{{ subsub.task_type === 'FEATURE' ? '★' : subsub.task_type === 'BUG' ? '⚠' : '📋' }}</span>
                                <span class="text-sm text-gray-400 cursor-pointer hover:text-purple-300 line-clamp-2 break-words block min-w-0" :title="subsub.title" @click="navigateToTask(subsub.id)">{{ subsub.title }}</span>
                              </div>
                              <div class="flex items-center justify-center shrink-0">
                                <span class="text-xs font-mono text-purple-400">{{ subsub.story_points || '–' }}</span>
                              </div>
                              <div class="flex items-center min-w-0">
                                <select :value="subsub.priority" @change="updateTaskField(subsub.id, 'priority', ($event.target as HTMLSelectElement).value)" class="text-xs bg-transparent border-0 focus:outline-none cursor-pointer w-full min-w-[6.5rem]" :class="priorityTextClass(subsub.priority)">
                                  <option value="CRITICAL">🔴 CRITICAL</option>
                                  <option value="HIGH">🟠 HIGH</option>
                                  <option value="MEDIUM">🟡 MEDIUM</option>
                                  <option value="LOW">🟢 LOW</option>
                                </select>
                              </div>
                              <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                              <div class="flex items-center justify-center min-w-0 w-full"><span class="text-xs text-gray-500 italic">Inherits</span></div>
                              <div class="flex items-center shrink-0">
                                <span class="text-xs px-1.5 py-0.5 rounded whitespace-nowrap" :class="taskStatusBadge(subsub.status)">{{ subsub.status.replace('_',' ') }}</span>
                              </div>
                              <div class="flex items-center w-min shrink-0"></div>
                            </div>
                          </template>
                        </template>
                      </template>
                    </template>
                  </div>
                </template>
              </template>

              <!-- Auto-load sentinel for infinite scroll (backlog tab only) -->
              <div
                v-if="projectDetailsTasksMeta?.has_more"
                ref="backlogAutoLoadSentinelRef"
                class="py-4 text-center text-xs text-gray-500"
              >
                <span v-if="isLoadingMoreProjectTasks">Loading more tasks…</span>
                <span v-else>Scroll down to auto-load more tasks</span>
              </div>

              <!-- Empty State -->
              <div v-if="!allTasks.filter(t => !t.parent_id).length" class="py-16 text-center text-gray-500">
                <p class="text-sm mb-3">No tasks in backlog yet.</p>
                <button @click="openCreateTaskModal()" class="btn-primary px-5 py-2 text-sm">Add First Task</button>
              </div>
            </div>
          </div>
        </div>

        <!-- TAB: Feature Roadmap (FEATURE tasks, scoped to this project) -->
        <div v-if="activeTab === 'feature-roadmap'">
          <FeatureRoadmapBoard
            :scope-project-id="project.id"
            :scope-project-code="project.code"
            :scope-project-name="project.name"
          />
        </div>

        <!-- TAB: Sprints (Sprint Management) -->
        <div v-if="activeTab === 'sprints'" class="space-y-6">
          <!-- Header + CTA -->
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div>
              <h2 class="text-xl font-bold text-white tracking-tight">Sprint Management</h2>
              <p class="text-sm text-gray-400 mt-0.5">Plan iterations, start sprints, and track progress in one place.</p>
            </div>
            <button @click="openSprintModal()" class="btn-primary-sm inline-flex items-center gap-2 shrink-0">
              <span>+</span> Create Sprint
            </button>
          </div>

          <!-- Active Sprint Hero (if any) -->
          <div v-if="activeSprint" class="rounded-2xl border border-purple-500/40 bg-gradient-to-br from-purple-900/30 to-gray-800/80 p-5 sm:p-6 shadow-xl">
            <div class="flex flex-wrap items-start justify-between gap-4">
              <div class="min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="px-2 py-0.5 text-xs font-semibold rounded-full bg-purple-500/30 text-purple-300">Active</span>
                  <h3 class="text-lg font-bold text-white truncate">{{ activeSprint.name }}</h3>
                </div>
                <p v-if="activeSprint.goal" class="text-sm text-gray-400 mt-1 line-clamp-2">{{ activeSprint.goal }}</p>
                <div class="flex flex-wrap gap-3 mt-3 text-xs text-gray-500">
                  <span v-if="activeSprint.start_date">Start: {{ formatDate(activeSprint.start_date) }}</span>
                  <span v-if="activeSprint.end_date">End: {{ formatDate(activeSprint.end_date) }}</span>
                </div>
              </div>
              <div class="flex flex-wrap items-center gap-2">
                <NuxtLink :to="`/projects/sprint/${activeSprint.id}?project=${route.params.id}`" class="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/20 text-white text-sm font-medium transition-colors">
                  Open Sprint →
                </NuxtLink>
                <button type="button" @click="openAddTasksToSprintModal(activeSprint)" class="px-4 py-2 rounded-xl bg-emerald-600/80 hover:bg-emerald-600 text-white text-sm font-medium transition-colors">
                  + Add Tasks
                </button>
                <button type="button" @click="openCompleteSprintModal(activeSprint)" class="px-4 py-2 rounded-xl bg-amber-600/80 hover:bg-amber-600 text-amber-100 text-sm font-medium transition-colors">
                  Complete Sprint
                </button>
                <button type="button" @click="openEditSprintModal(activeSprint)" class="p-2 rounded-lg text-gray-400 hover:text-white hover:bg-gray-700/50 transition-colors" title="Edit sprint">
                  ✎
                </button>
              </div>
            </div>
            <div class="mt-5 grid grid-cols-3 gap-4">
              <div class="bg-gray-800/60 rounded-xl p-4 text-center">
                <div class="text-2xl font-bold text-white">{{ sprintTaskCount('total') }}</div>
                <div class="text-xs text-gray-500 uppercase tracking-wide mt-0.5">Tasks</div>
              </div>
              <div class="bg-gray-800/60 rounded-xl p-4 text-center">
                <div class="text-2xl font-bold text-green-400">{{ sprintTaskCount('done') }}</div>
                <div class="text-xs text-gray-500 uppercase tracking-wide mt-0.5">Done</div>
              </div>
              <div class="bg-gray-800/60 rounded-xl p-4 text-center">
                <div class="text-2xl font-bold text-purple-400">{{ sprintTaskCount('sp') }}</div>
                <div class="text-xs text-gray-500 uppercase tracking-wide mt-0.5">Story Pts</div>
              </div>
            </div>
            <div v-if="sprintTaskCount('total') > 0" class="mt-4">
              <div class="flex justify-between text-xs text-gray-500 mb-1.5">
                <span>Progress</span>
                <span>{{ Math.round((sprintTaskCount('done') / sprintTaskCount('total')) * 100) }}%</span>
              </div>
              <div class="h-2.5 bg-gray-700 rounded-full overflow-hidden">
                <div class="h-full bg-gradient-to-r from-purple-500 to-pink-500 rounded-full transition-all duration-500" :style="{ width: Math.round((sprintTaskCount('done') / sprintTaskCount('total')) * 100) + '%' }" />
              </div>
            </div>
          </div>

          <!-- All Sprints -->
          <div class="rounded-2xl border border-gray-700 bg-gray-800/50 overflow-hidden">
            <div class="px-4 sm:px-5 py-4 border-b border-gray-700">
              <h3 class="text-base font-semibold text-gray-200">All Sprints</h3>
              <p class="text-xs text-gray-500 mt-0.5">Start, add tasks, edit, or reopen from here.</p>
            </div>
            <div v-if="sprints.length === 0" class="py-16 text-center">
              <div class="text-5xl mb-3 opacity-60">🏃</div>
              <p class="text-gray-400 font-medium">No sprints yet</p>
              <p class="text-sm text-gray-500 mt-1 max-w-sm mx-auto">Create a sprint to plan your iterations and move tasks from the backlog.</p>
              <button @click="openSprintModal()" class="mt-4 btn-primary px-5 py-2.5 rounded-xl">Create first sprint</button>
            </div>
            <ul v-else class="divide-y divide-gray-700/80">
              <li
                v-for="(s, sIdx) in sprintsOrdered"
                :key="s.id"
                class="flex flex-col sm:flex-row sm:items-center gap-4 px-4 sm:px-5 py-4 hover:bg-gray-700/30 transition-colors"
                :class="{ 'bg-purple-500/5': s.status === 'ACTIVE', 'opacity-60': sprintDragId === s.id }"
                draggable="true"
                @dragstart="onSprintDragStart($event, s.id)"
                @dragover="onSprintDragOver"
                @drop.stop="onSprintDrop($event, sIdx)"
              >
                <div class="flex items-center gap-2 flex-1 min-w-0">
                  <span class="text-gray-500 cursor-grab shrink-0 select-none text-sm" title="ลากเพื่อเรียงลำดับ">⋮⋮</span>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 flex-wrap">
                    <NuxtLink :to="`/projects/sprint/${s.id}?project=${route.params.id}`" class="font-semibold text-white hover:text-purple-300 truncate transition-colors">
                      {{ s.name }}
                    </NuxtLink>
                    <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full shrink-0" :class="s.status === 'ACTIVE' ? 'bg-purple-500/25 text-purple-300' : s.status === 'COMPLETED' ? 'bg-gray-600 text-gray-400' : 'bg-amber-500/20 text-amber-400'">
                      {{ s.status }}
                    </span>
                  </div>
                  <p v-if="s.goal" class="text-sm text-gray-500 truncate mt-0.5">{{ s.goal }}</p>
                  <div class="flex items-center gap-4 mt-2 text-xs text-gray-500">
                    <span v-if="s.start_date">{{ formatDate(s.start_date) }}</span>
                    <span v-if="s.end_date">– {{ formatDate(s.end_date) }}</span>
                    <span class="text-gray-600">·</span>
                    <span>{{ getSprintStats(s.id).done }}/{{ getSprintStats(s.id).total }} tasks</span>
                    <span v-if="getSprintStats(s.id).total > 0" class="text-purple-400">{{ getSprintStats(s.id).sp }} SP</span>
                  </div>
                </div>
                </div>
                <div class="flex flex-wrap items-center gap-2 shrink-0">
                  <template v-if="s.status === 'PLANNING'">
                    <button type="button" :disabled="!!activeSprint" @click="!activeSprint && handleStartSprint(s.id)" class="px-3 py-1.5 rounded-lg text-xs font-medium bg-purple-600 hover:bg-purple-500 disabled:opacity-50 disabled:cursor-not-allowed text-white transition-colors" :title="activeSprint ? 'Complete or reopen the active sprint first' : 'Start this sprint'">
                      Start
                    </button>
                  </template>
                  <template v-if="s.status === 'COMPLETED'">
                    <button type="button" @click="openReopenSprintModal(s)" class="px-3 py-1.5 rounded-lg text-xs font-medium text-amber-400 hover:bg-amber-500/20 transition-colors" title="Reopen sprint">
                      Reopen
                    </button>
                  </template>
                  <button type="button" @click="openAddTasksToSprintModal(s)" class="px-3 py-1.5 rounded-lg text-xs font-medium text-emerald-400 hover:bg-emerald-500/20 transition-colors" title="Add tasks to sprint">
                    + Tasks
                  </button>
                  <button type="button" @click="openEditSprintModal(s)" class="px-3 py-1.5 rounded-lg text-xs font-medium text-gray-400 hover:bg-gray-600/50 transition-colors" title="Edit sprint">
                    Edit
                  </button>
                  <button type="button" @click="openDeleteSprintModal(s)" class="px-3 py-1.5 rounded-lg text-xs font-medium text-red-400 hover:bg-red-500/20 transition-colors" title="Delete sprint">
                    Delete
                  </button>
                </div>
              </li>
            </ul>
            <p class="px-4 sm:px-5 py-3 text-[10px] text-gray-500 border-t border-gray-700/80 bg-gray-900/40">
              One project can have only one active sprint. Complete or reopen the current sprint before starting another.
            </p>
          </div>
        </div>

        <!-- TAB 6: Analytics -->
        <div v-if="activeTab === 'analytics'">
          <div v-if="analyticsLoading" class="flex flex-col items-center justify-center py-20">
            <div class="animate-spin text-6xl mb-4">⚙️</div>
            <p class="text-sm text-gray-500">กำลังโหลด analytics...</p>
          </div>
          <ProjectAnalytics v-else-if="analytics" :analytics="analytics" />
          <div v-else class="text-center py-20 text-gray-500 text-sm">Failed to load analytics.</div>
        </div>

        <!-- TAB: Capital (Internal VC — per-project capital tracking) — CEO only -->
        <div v-if="activeTab === 'capital' && isProjectCeo">
          <ProjectCapitalPanel :project-id="project.id" :team-id="project.team_id" />
        </div>

        <!-- TAB: Costing & Quotation — CEO only -->
        <div v-if="activeTab === 'costing' && isProjectCeo">
          <QuotationBuilder :project-id="project.id" :project-name="project.name" />
        </div>

        <!-- TAB: Backup & Restore — CEO only -->
        <div v-if="activeTab === 'backup' && isProjectCeo">
          <ProjectBackupPanel :project-id="project.id" @restored="onProjectRestored" />
        </div>
      </div>
    </div>

    <!-- Import from Google Slides Modal (Backlog) -->
    <div v-if="showBacklogImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto" @click.self="closeBacklogImportModal">
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="backlogImportStep === 'select' ? 'max-w-5xl' : 'max-w-xl'"
      >
        <!-- Modal Header (fixed, never scrolls) -->
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-purple-600/20 border border-purple-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-purple-400" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import from Google Slides</h2>
              <p class="text-xs text-gray-400">สร้าง sub-task จากแต่ละ slide — Manual Triage ก่อน import</p>
            </div>
          </div>
          <button @click="closeBacklogImportModal" class="text-gray-500 hover:text-white transition-colors shrink-0 ml-4">✕</button>
        </div>

        <!-- Modal Body (scrollable) -->
        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">

          <!-- Result state -->
          <template v-if="backlogImportStep === 'result' && backlogImportResult">
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
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(task.code) }}</span>
                <span class="text-gray-200 truncate">{{ task.title }}</span>
              </div>
            </div>
            <button @click="closeBacklogImportModal" class="w-full btn-primary py-2.5">Done</button>
          </template>

          <!-- Step 2: Manual Triage Table -->
          <template v-else-if="backlogImportStep === 'select' && backlogImportPreview">
            <div class="flex items-center justify-between gap-3">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ backlogImportPreview.presentation_title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">
                  {{ backlogImportSelectedIndices.length }} / {{ backlogImportPreview.slides.length }} slides selected — Est. minutes ไม่บังคับ
                </p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" @click="backlogImportSelectAll" class="btn-ghost-sm">ทั้งหมด</button>
                <button type="button" @click="backlogImportDeselectAll" class="btn-ghost-sm">ยกเลิก</button>
                <button type="button" @click="backlogImportSelectOnlyNew" class="btn-ghost-sm text-purple-400">เฉพาะใหม่</button>
              </div>
            </div>

            <!-- Triage Table -->
            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-3 text-left w-8"></th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-10">#</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[200px]">Task Title</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Assignee</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Est. min <span class="text-gray-500 font-normal">(ไม่บังคับ)</span></th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="s in backlogImportPreview.slides"
                    :key="s.index"
                    class="border-b border-gray-700/30 transition-colors"
                    :class="backlogImportSelectedIndices.includes(s.index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-3">
                      <input
                        v-model="backlogImportSelectedIndices"
                        type="checkbox"
                        :value="s.index"
                        class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500"
                      />
                    </td>
                    <td class="py-2 px-3 text-xs text-gray-400 font-mono">
                      {{ s.index }}
                      <span v-if="s.hidden" class="text-amber-400 ml-1 text-[10px]">ซ่อน</span>
                      <span v-else-if="(backlogImportPreview.already_imported_slide_indices || []).includes(s.index)" class="text-gray-500 ml-1 text-[10px]">นำเข้าแล้ว</span>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="backlogImportTriagedSlides[s.index]"
                        v-model="backlogImportTriagedSlides[s.index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-purple-500/60"
                        :disabled="!backlogImportSelectedIndices.includes(s.index)"
                      />
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="backlogImportTriagedSlides[s.index]"
                        v-model="backlogImportTriagedSlides[s.index].assignee_id"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!backlogImportSelectedIndices.includes(s.index)"
                      >
                        <option :value="null">— Unassigned —</option>
                        <option v-for="u in backlogImportAssignees" :key="u.id" :value="u.id">{{ u.display_name || u.email }}</option>
                      </select>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="backlogImportTriagedSlides[s.index]"
                        v-model.number="backlogImportTriagedSlides[s.index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!backlogImportSelectedIndices.includes(s.index)"
                        placeholder="0"
                      />
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="backlogImportTriagedSlides[s.index]"
                        v-model="backlogImportTriagedSlides[s.index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!backlogImportSelectedIndices.includes(s.index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
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
          </template>

          <!-- Step 1: Form (URL + Epic + Parent Task) -->
          <template v-else>
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
            <div v-if="epics.length">
              <label class="label">Epic</label>
              <select v-model="backlogImportForm.epic_id" class="input-field w-full" :disabled="isBacklogLoadingPreview" @change="onBacklogImportEpicChange">
                <option value="">— ทุก Epic / Unassigned —</option>
                <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
              </select>
            </div>
            <div>
              <label class="label">Target Parent Task <span class="text-gray-500 font-normal">(Sub-tasks จะถูกสร้างใต้ task นี้)</span></label>
              <select v-model="backlogImportForm.parent_id" class="input-field w-full" :disabled="isBacklogLoadingPreview">
                <option value="">— No Parent (Top-level tasks) —</option>
                <option v-for="t in backlogParentTaskOptions" :key="t.id" :value="t.id">
                  [{{ taskCodeSuffix(t.code) }}] {{ t.title }}
                </option>
              </select>
              <p v-if="backlogImportForm.epic_id && !backlogParentTaskOptions.length" class="text-xs text-amber-400/80 mt-1">ไม่มี task ใน Epic นี้ที่จะเป็น parent ได้</p>
            </div>
            <div v-if="backlogImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ backlogImportError }}</div>
            <div class="flex gap-3">
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
          </template>

        </div>
      </div>
    </div>

    <!-- Import PPTX File Upload Modal -->
    <div v-if="showPPTXImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto" @click.self="closePPTXImportModal">
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="pptxImportStep === 'select' ? 'max-w-5xl' : 'max-w-xl'"
      >
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-orange-600/20 border border-orange-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-orange-400" fill="currentColor" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm4 18H6V4h7v5h5v11z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import PPTX</h2>
              <p class="text-xs text-gray-400">Upload ไฟล์ .pptx (export จาก Canva, PowerPoint, Keynote ฯลฯ) — 1 slide = 1 task</p>
            </div>
          </div>
          <button @click="closePPTXImportModal" class="text-gray-500 hover:text-white transition-colors shrink-0 ml-4">✕</button>
        </div>

        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <!-- Result -->
          <template v-if="pptxImportStep === 'result' && pptxImportResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ pptxImportResult.title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ pptxImportResult.created_count }} tasks จาก {{ pptxImportResult.page_count }} slides</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div
                v-for="task in pptxImportResult.tasks"
                :key="task.id"
                class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm"
              >
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(task.code) }}</span>
                <span class="text-gray-200 truncate">{{ task.title }}</span>
              </div>
            </div>
            <button @click="closePPTXImportModal" class="w-full btn-primary py-2.5">Done</button>
          </template>

          <!-- Step 2: Triage Table -->
          <template v-else-if="pptxImportStep === 'select' && pptxImportPreview">
            <div class="flex items-center justify-between gap-3">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ pptxImportPreview.title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">
                  {{ pptxImportSelectedIndices.length }} / {{ pptxImportPreview.slides.length }} slides selected — Est. minutes ไม่บังคับ
                </p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" @click="pptxImportSelectAll" class="btn-ghost-sm">ทั้งหมด</button>
                <button type="button" @click="pptxImportDeselectAll" class="btn-ghost-sm">ยกเลิก</button>
                <button type="button" @click="pptxImportSelectOnlyVisible" class="btn-ghost-sm text-orange-400">เฉพาะที่แสดง</button>
              </div>
            </div>

            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-3 text-left w-8"></th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-10">#</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[200px]">Task Title</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Assignee</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Est. min <span class="text-gray-500 font-normal">(ไม่บังคับ)</span></th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="s in pptxImportPreview.slides"
                    :key="s.index"
                    class="border-b border-gray-700/30 transition-colors"
                    :class="pptxImportSelectedIndices.includes(s.index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-3">
                      <input
                        v-model="pptxImportSelectedIndices"
                        type="checkbox"
                        :value="s.index"
                        class="rounded border-gray-500 bg-gray-700 text-orange-500 focus:ring-orange-500"
                      />
                    </td>
                    <td class="py-2 px-3 text-xs text-gray-400 font-mono">
                      {{ s.index }}
                      <span v-if="s.hidden" class="text-amber-400 ml-1 text-[10px]">ซ่อน</span>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="pptxImportTriagedSlides[s.index]"
                        v-model="pptxImportTriagedSlides[s.index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxImportSelectedIndices.includes(s.index)"
                      />
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="pptxImportTriagedSlides[s.index]"
                        v-model="pptxImportTriagedSlides[s.index].assignee_id"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxImportSelectedIndices.includes(s.index)"
                      >
                        <option :value="null">— Unassigned —</option>
                        <option v-for="u in backlogImportAssignees" :key="u.id" :value="u.id">{{ u.display_name || u.email }}</option>
                      </select>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="pptxImportTriagedSlides[s.index]"
                        v-model.number="pptxImportTriagedSlides[s.index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxImportSelectedIndices.includes(s.index)"
                        placeholder="0"
                      />
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="pptxImportTriagedSlides[s.index]"
                        v-model="pptxImportTriagedSlides[s.index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxImportSelectedIndices.includes(s.index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="pptxImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ pptxImportError }}</div>
            <div class="flex gap-3">
              <button
                @click="submitPPTXImport"
                :disabled="isPPTXImporting || pptxImportSelectedIndices.length === 0"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
              >
                <svg v-if="isPPTXImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isPPTXImporting ? 'กำลัง import...' : `Import ${pptxImportSelectedIndices.length} Slides` }}
              </button>
              <button type="button" @click="pptxImportStep = 'form'" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">กลับ</button>
            </div>
          </template>

          <!-- Step 1: File Picker + Epic/Parent -->
          <template v-else>
            <!-- File Drop Zone -->
            <div
              class="relative border-2 border-dashed rounded-xl p-8 text-center transition-colors"
              :class="pptxDragOver ? 'border-orange-500/60 bg-orange-500/5' : 'border-gray-600/60 hover:border-gray-500/60'"
              @dragover.prevent="pptxDragOver = true"
              @dragleave="pptxDragOver = false"
              @drop.prevent="onPPTXFileDrop"
            >
              <svg class="w-10 h-10 text-gray-500 mx-auto mb-3" fill="currentColor" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm4 18H6V4h7v5h5v11z"/></svg>
              <p class="text-sm font-medium text-gray-300 mb-1">
                {{ pptxImportFile ? pptxImportFile.name : 'ลากไฟล์ .pptx มาวาง หรือ' }}
              </p>
              <p v-if="!pptxImportFile" class="text-xs text-gray-500 mb-3">รองรับไฟล์ .pptx จาก Canva, PowerPoint, Keynote</p>
              <p v-else class="text-xs text-green-400 mb-3">{{ (pptxImportFile.size / 1024 / 1024).toFixed(1) }} MB</p>
              <label class="btn-import-sm cursor-pointer inline-flex items-center gap-1.5">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
                {{ pptxImportFile ? 'เปลี่ยนไฟล์' : 'เลือกไฟล์' }}
                <input type="file" accept=".pptx" class="sr-only" @change="onPPTXFileChange" />
              </label>
            </div>

            <div v-if="epics.length">
              <label class="label">Epic</label>
              <select v-model="pptxImportForm.epic_id" class="input-field w-full" :disabled="isPPTXLoadingPreview" @change="onPPTXImportEpicChange">
                <option value="">— ทุก Epic / Unassigned —</option>
                <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
              </select>
            </div>
            <div>
              <label class="label">Target Parent Task <span class="text-gray-500 font-normal">(Sub-tasks จะถูกสร้างใต้ task นี้)</span></label>
              <select v-model="pptxImportForm.parent_id" class="input-field w-full" :disabled="isPPTXLoadingPreview">
                <option value="">— No Parent (Top-level tasks) —</option>
                <option v-for="t in pptxParentTaskOptions" :key="t.id" :value="t.id">
                  [{{ taskCodeSuffix(t.code) }}] {{ t.title }}
                </option>
              </select>
            </div>

            <div v-if="pptxImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ pptxImportError }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                @click="loadPPTXImportPreview"
                :disabled="isPPTXLoadingPreview || !pptxImportFile"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
              >
                <svg v-if="isPPTXLoadingPreview" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isPPTXLoadingPreview ? 'กำลังโหลด...' : 'โหลดรายการ slide' }}
              </button>
              <button @click="closePPTXImportModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Import from Google Sheets Modal (Backlog) -->
    <div v-if="showSheetsImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto" @click.self="closeSheetsImportModal">
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="sheetsImportStep === 'select' ? 'max-w-6xl' : 'max-w-xl'"
      >
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-emerald-600/20 border border-emerald-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-emerald-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import from Google Sheets</h2>
              <p class="text-xs text-gray-400">Preview แถวจาก CSV — แก้ไขก่อน import (BUG, unassigned)</p>
            </div>
          </div>
          <button type="button" @click="closeSheetsImportModal" class="text-gray-500 hover:text-white transition-colors shrink-0 ml-4">✕</button>
        </div>

        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <template v-if="sheetsImportStep === 'result' && sheetsImportResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ sheetsImportResult.sheet_title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ sheetsImportResult.created_count }} tasks จาก Google Sheet</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div
                v-for="task in sheetsImportResult.tasks"
                :key="task.id"
                class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm"
              >
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(task.code) }}</span>
                <span class="text-gray-200 truncate">{{ task.title }}</span>
              </div>
            </div>
            <button type="button" @click="closeSheetsImportModal" class="w-full btn-primary py-2.5">Done</button>
          </template>

          <template v-else-if="sheetsImportStep === 'select' && sheetsImportPreview">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ sheetsImportPreview.sheet_title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">
                  {{ sheetsImportSelectedRowIndices.length }} / {{ sheetsImportPreview.rows.length }} rows selected
                </p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" @click="sheetsImportSelectAll" class="btn-ghost-sm">ทั้งหมด</button>
                <button type="button" @click="sheetsImportDeselectAll" class="btn-ghost-sm">ยกเลิก</button>
              </div>
            </div>

            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-2 text-left w-8"></th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-12">#</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[180px]">Title</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-36">Due</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[120px]">Status</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Notes</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-28">Est. min <span class="text-gray-500 font-normal">(ไม่บังคับ)</span></th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-28">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="r in sheetsImportPreview.rows"
                    :key="r.row_index"
                    class="border-b border-gray-700/30 transition-colors align-top"
                    :class="sheetsImportSelectedRowIndices.includes(r.row_index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-2">
                      <input
                        v-model="sheetsImportSelectedRowIndices"
                        type="checkbox"
                        :value="r.row_index"
                        class="rounded border-gray-500 bg-gray-700 text-emerald-500 focus:ring-emerald-500"
                      />
                    </td>
                    <td class="py-2 px-2 text-xs text-gray-400 font-mono">{{ r.row_index }}</td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model="sheetsImportTriagedRows[r.row_index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                      />
                      <p v-if="r.raw_status" class="text-[10px] text-gray-500 mt-0.5 truncate max-w-[200px]" :title="r.raw_status">KG: {{ r.raw_status }}</p>
                    </td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model="sheetsImportTriagedRows[r.row_index].due_date"
                        type="date"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60 max-w-[9.5rem]"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                      />
                    </td>
                    <td class="py-2 px-1">
                      <select
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model="sheetsImportTriagedRows[r.row_index].status"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                      >
                        <option value="PENDING">PENDING</option>
                        <option value="IN_PROGRESS">IN_PROGRESS</option>
                        <option value="READY_FOR_TEST">READY_FOR_TEST</option>
                        <option value="READY_FOR_UAT">READY_FOR_UAT</option>
                        <option value="COMPLETED">COMPLETED</option>
                        <option value="CANCELLED">CANCELLED</option>
                      </select>
                    </td>
                    <td class="py-2 px-1">
                      <textarea
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model="sheetsImportTriagedRows[r.row_index].notes"
                        rows="2"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-emerald-500/60 resize-y min-h-[2.25rem]"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                      />
                    </td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model.number="sheetsImportTriagedRows[r.row_index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                        placeholder="—"
                      />
                    </td>
                    <td class="py-2 px-1">
                      <select
                        v-if="sheetsImportTriagedRows[r.row_index]"
                        v-model="sheetsImportTriagedRows[r.row_index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsImportSelectedRowIndices.includes(r.row_index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="sheetsImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ sheetsImportError }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                @click="submitSheetsImport"
                :disabled="isSheetsImporting || sheetsImportSelectedRowIndices.length === 0"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
              >
                <svg v-if="isSheetsImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isSheetsImporting ? 'กำลัง import...' : `Import ${sheetsImportSelectedRowIndices.length} rows` }}
              </button>
              <button type="button" @click="sheetsImportStep = 'form'" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">กลับ</button>
            </div>
          </template>

          <template v-else>
            <div>
              <label class="label">Google Sheets URL *</label>
              <input
                v-model="sheetsImportForm.sheet_url"
                type="url"
                class="input-field w-full"
                placeholder="https://docs.google.com/spreadsheets/d/..."
                :disabled="isSheetsLoadingPreview"
              />
              <p class="text-xs text-gray-500 mt-1">
                เปิดสิทธิ์ &quot;Anyone with the link can view&quot; — คัดลอก URL จากแท็บที่ต้องการ (เช่น 31/03) ให้มี
                <code class="text-gray-400">gid=...</code> ในลิงก์
                — รองรับทั้งแบบ IFP x KG (คอลัมน์ A–K) และแบบ IOD bug list (ปัญหา / working Status)
              </p>
            </div>
            <div v-if="epics.length">
              <label class="label">Epic</label>
              <select v-model="sheetsImportForm.epic_id" class="input-field w-full" :disabled="isSheetsLoadingPreview" @change="onSheetsImportEpicChange">
                <option value="">— ทุก Epic / Unassigned —</option>
                <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
              </select>
            </div>
            <div>
              <label class="label">Target Parent Task <span class="text-gray-500 font-normal">(optional sub-tasks)</span></label>
              <select v-model="sheetsImportForm.parent_id" class="input-field w-full" :disabled="isSheetsLoadingPreview">
                <option value="">— No Parent (Top-level tasks) —</option>
                <option v-for="t in sheetsParentTaskOptions" :key="t.id" :value="t.id">
                  [{{ taskCodeSuffix(t.code) }}] {{ t.title }}
                </option>
              </select>
              <p v-if="sheetsImportForm.epic_id && !sheetsParentTaskOptions.length" class="text-xs text-amber-400/80 mt-1">ไม่มี task ใน Epic นี้ที่จะเป็น parent ได้</p>
            </div>
            <div v-if="sheetsImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ sheetsImportError }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                @click="loadSheetsImportPreview"
                :disabled="isSheetsLoadingPreview || !sheetsImportForm.sheet_url.trim()"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
              >
                <svg v-if="isSheetsLoadingPreview" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isSheetsLoadingPreview ? 'กำลังโหลด...' : 'โหลด preview' }}
              </button>
              <button type="button" @click="closeSheetsImportModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- ══════ IOD Import Modal ══════ -->
    <div v-if="showIODImportModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto" @click.self="closeIODImportModal">
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[92vh]"
        :class="iodImportStep === 'select' ? 'max-w-7xl' : 'max-w-xl'"
      >
        <!-- Header -->
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-blue-600/20 border border-blue-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import IOD Bug Sheet</h2>
              <p class="text-xs text-gray-400">นำเข้า Bug จาก Google Sheets แบบ IOD — ดึง Header, Detail, URL, Method, Payload ครบ</p>
            </div>
          </div>
          <button type="button" @click="closeIODImportModal" class="text-gray-500 hover:text-white transition-colors shrink-0 ml-4">✕</button>
        </div>

        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <!-- Step: Result -->
          <template v-if="iodImportStep === 'result' && iodImportResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ iodImportResult.sheet_title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ iodImportResult.created_count }} bugs จาก IOD Sheet</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div v-for="task in iodImportResult.tasks" :key="task.id" class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(task.code) }}</span>
                <span class="text-gray-200 truncate">{{ task.title }}</span>
              </div>
            </div>
            <button type="button" @click="closeIODImportModal" class="w-full btn-primary py-2.5">Done</button>
          </template>

          <!-- Step: Preview Table -->
          <template v-else-if="iodImportStep === 'select' && iodImportPreview">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ iodImportPreview.sheet_title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">{{ iodImportSelectedRowIndices.length }} / {{ iodImportPreview.rows.length }} bugs selected</p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" @click="iodImportSelectedRowIndices = iodImportPreview!.rows.map(r => r.row_index)" class="btn-ghost-sm">ทั้งหมด</button>
                <button type="button" @click="iodImportSelectedRowIndices = []" class="btn-ghost-sm">ยกเลิก</button>
              </div>
            </div>

            <!-- Full-detail IOD table -->
            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-xs" style="min-width:900px">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-2 w-8"></th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-10">#</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-36">Section (Header)</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold" style="min-width:200px">Detail (Bug Title)</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-52">Header Link</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-24">Method</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-52">Payload</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-24">Status</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-24">Priority</th>
                    <th class="py-2 px-2 text-left text-gray-400 font-semibold w-20">Est. min</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="r in iodImportPreview.rows"
                    :key="r.row_index"
                    class="border-b border-gray-700/30 transition-colors align-top"
                    :class="iodImportSelectedRowIndices.includes(r.row_index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-40'"
                  >
                    <td class="py-2 px-2">
                      <input v-model="iodImportSelectedRowIndices" type="checkbox" :value="r.row_index"
                        class="rounded border-gray-500 bg-gray-700 text-blue-500 focus:ring-blue-500" />
                    </td>
                    <td class="py-2 px-2 text-gray-500 font-mono">{{ r.row_index }}</td>
                    <!-- Section (readonly info) -->
                    <td class="py-2 px-1">
                      <p class="text-gray-300 text-[11px] leading-snug max-h-16 overflow-y-auto">{{ r.header || '—' }}</p>
                    </td>
                    <!-- Detail / Title (editable) -->
                    <td class="py-2 px-1">
                      <textarea
                        v-if="iodImportTriagedRows[r.row_index]"
                        v-model="iodImportTriagedRows[r.row_index].title"
                        rows="3"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-blue-500/60 resize-y"
                        :disabled="!iodImportSelectedRowIndices.includes(r.row_index)"
                      />
                    </td>
                    <!-- Header Link (readonly) -->
                    <td class="py-2 px-1">
                      <a v-if="r.header_link" :href="r.header_link" target="_blank" rel="noopener"
                        class="text-blue-400 hover:underline break-all text-[10px] leading-snug line-clamp-3">{{ r.header_link }}</a>
                      <span v-else class="text-gray-600">—</span>
                    </td>
                    <!-- Method (readonly badge) -->
                    <td class="py-2 px-1 text-center">
                      <span v-if="r.request_method"
                        class="inline-block px-1.5 py-0.5 rounded text-[10px] font-mono font-bold"
                        :class="{
                          'bg-green-900/50 text-green-300': r.request_method === 'GET',
                          'bg-blue-900/50 text-blue-300': r.request_method === 'POST',
                          'bg-yellow-900/50 text-yellow-300': r.request_method === 'PUT' || r.request_method === 'PATCH',
                          'bg-red-900/50 text-red-300': r.request_method === 'DELETE',
                          'bg-gray-700 text-gray-300': !['GET','POST','PUT','PATCH','DELETE'].includes(r.request_method)
                        }">{{ r.request_method }}</span>
                      <span v-else class="text-gray-600">—</span>
                    </td>
                    <!-- Payload (collapsible readonly) -->
                    <td class="py-2 px-1">
                      <details v-if="r.payload" class="cursor-pointer">
                        <summary class="text-blue-400 text-[10px] hover:text-blue-300 select-none">ดู Payload</summary>
                        <pre class="mt-1 text-[10px] text-gray-300 bg-gray-900/60 rounded p-1 max-h-32 overflow-auto whitespace-pre-wrap break-all">{{ r.payload }}</pre>
                      </details>
                      <span v-else class="text-gray-600">—</span>
                    </td>
                    <!-- Status (editable) -->
                    <td class="py-2 px-1">
                      <select v-if="iodImportTriagedRows[r.row_index]"
                        v-model="iodImportTriagedRows[r.row_index].status"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-blue-500/60"
                        :disabled="!iodImportSelectedRowIndices.includes(r.row_index)">
                        <option value="PENDING">PENDING</option>
                        <option value="IN_PROGRESS">IN_PROGRESS</option>
                        <option value="READY_FOR_TEST">READY_FOR_TEST</option>
                        <option value="READY_FOR_UAT">READY_FOR_UAT</option>
                        <option value="BLOCKED">BLOCKED</option>
                        <option value="COMPLETED">COMPLETED</option>
                        <option value="CANCELLED">CANCELLED</option>
                      </select>
                    </td>
                    <!-- Priority (editable) -->
                    <td class="py-2 px-1">
                      <select v-if="iodImportTriagedRows[r.row_index]"
                        v-model="iodImportTriagedRows[r.row_index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-blue-500/60"
                        :disabled="!iodImportSelectedRowIndices.includes(r.row_index)">
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                    <!-- Est. minutes (editable) -->
                    <td class="py-2 px-1">
                      <input v-if="iodImportTriagedRows[r.row_index]"
                        v-model.number="iodImportTriagedRows[r.row_index].estimated_minutes"
                        type="number" min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-blue-500/60"
                        :disabled="!iodImportSelectedRowIndices.includes(r.row_index)"
                        placeholder="—" />
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="iodImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ iodImportError }}</div>
            <div class="flex gap-3">
              <button type="button" @click="submitIODImport"
                :disabled="isIODImporting || iodImportSelectedRowIndices.length === 0"
                class="flex-1 py-2.5 disabled:opacity-40 flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-500 text-white font-semibold rounded-xl transition-colors">
                <svg v-if="isIODImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isIODImporting ? 'กำลัง import...' : `Import ${iodImportSelectedRowIndices.length} bugs` }}
              </button>
              <button type="button" @click="iodImportStep = 'form'" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">กลับ</button>
            </div>
          </template>

          <!-- Step: Form -->
          <template v-else>
            <div>
              <label class="label">Google Sheets URL (แท็บ 31/03 หรือ แท็บใดก็ได้) *</label>
              <input v-model="iodImportForm.sheet_url" type="url" class="input-field w-full"
                placeholder="https://docs.google.com/spreadsheets/d/...?gid=1794611958"
                :disabled="isIODLoadingPreview" />
              <p class="text-xs text-gray-500 mt-1">
                คัดลอก URL จากแท็บที่ต้องการใน Google Sheets ให้มี <code class="text-gray-400">gid=...</code> —
                เปิดสิทธิ์ "Anyone with the link can view"
              </p>
            </div>
            <div v-if="epics.length">
              <label class="label">Epic</label>
              <select v-model="iodImportForm.epic_id" class="input-field w-full" :disabled="isIODLoadingPreview" @change="onIODImportEpicChange">
                <option value="">— ทุก Epic / Unassigned —</option>
                <option v-for="ep in epics" :key="ep.id" :value="ep.id">{{ ep.title }}</option>
              </select>
            </div>
            <div>
              <label class="label">Target Parent Task <span class="text-gray-500 font-normal">(optional)</span></label>
              <select v-model="iodImportForm.parent_id" class="input-field w-full" :disabled="isIODLoadingPreview">
                <option value="">— No Parent (Top-level tasks) —</option>
                <option v-for="t in iodParentTaskOptions" :key="t.id" :value="t.id">[{{ taskCodeSuffix(t.code) }}] {{ t.title }}</option>
              </select>
            </div>
            <div v-if="iodImportError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ iodImportError }}</div>
            <div class="flex gap-3">
              <button type="button" @click="loadIODImportPreview"
                :disabled="isIODLoadingPreview || !iodImportForm.sheet_url.trim()"
                class="flex-1 py-2.5 disabled:opacity-40 flex items-center justify-center gap-2 bg-blue-600 hover:bg-blue-500 text-white font-semibold rounded-xl transition-colors">
                <svg v-if="isIODLoadingPreview" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ isIODLoadingPreview ? 'กำลังโหลด...' : 'โหลด preview' }}
              </button>
              <button type="button" @click="closeIODImportModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>
    <!-- ══════ End IOD Import Modal ══════ -->

    <!-- Create Task Modal -->
    <div v-if="showCreateTaskModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-start justify-center z-50 p-3 sm:p-6 overflow-y-auto" @click.self="closeCreateTaskModal">
      <div class="create-task-modal bg-gray-800 border border-gray-700 rounded-2xl w-full max-w-7xl shadow-2xl my-4 sm:my-8 flex flex-col max-h-[calc(100dvh-2rem)] min-h-0">
        <div class="flex items-center justify-between px-6 sm:px-8 pt-6 sm:pt-8 pb-4 shrink-0 border-b border-gray-700/80">
          <h2 class="text-2xl sm:text-3xl font-bold text-white tracking-tight">{{ createTaskForm.parent_id ? 'Add Sub-task' : 'Add Task' }}</h2>
          <button type="button" @click="closeCreateTaskModal" class="shrink-0 w-11 h-11 flex items-center justify-center rounded-xl text-gray-400 hover:text-white hover:bg-gray-700 text-xl leading-none" aria-label="Close">✕</button>
        </div>
        <div class="px-6 sm:px-8 py-6 sm:py-8 space-y-6 sm:space-y-7 flex-1 overflow-y-auto overscroll-contain min-h-0">
          <!-- Task Type Selector -->
          <div>
            <label class="label">Type *</label>
            <div class="grid grid-cols-3 gap-3 sm:gap-4">
              <button
                type="button"
                @click="createTaskForm.task_type = 'FEATURE'"
                :class="createTaskForm.task_type === 'FEATURE' ? 'border-purple-500 bg-purple-500/20 text-purple-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-purple-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">★</span> Feature
              </button>
              <button
                type="button"
                @click="createTaskForm.task_type = 'TASK'"
                :class="createTaskForm.task_type === 'TASK' ? 'border-blue-500 bg-blue-500/20 text-blue-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-blue-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">📋</span> Task
              </button>
              <button
                type="button"
                @click="createTaskForm.task_type = 'BUG'"
                :class="createTaskForm.task_type === 'BUG' ? 'border-red-500 bg-red-500/20 text-red-300' : 'border-gray-600 bg-gray-900/50 text-gray-400 hover:border-red-500/50'"
                class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]"
              >
                <span class="text-xl sm:text-2xl leading-none">⚠</span> Bug
              </button>
            </div>
            <!-- Product Owner rule hint for FEATURE type -->
            <div v-if="createTaskForm.task_type === 'FEATURE'" class="mt-3 flex items-start gap-3 p-4 bg-purple-900/20 border border-purple-500/30 rounded-xl text-sm sm:text-base text-purple-300 leading-relaxed">
              <span class="shrink-0 mt-0.5">★</span>
              <span><strong>Feature mode:</strong> Acts as a parent container. Assignee and estimated effort are disabled — add sub-tasks of type Task/Bug to assign work.</span>
            </div>
          </div>

          <div>
            <label class="label">Title *</label>
            <input v-model="createTaskForm.title" type="text" class="input-field w-full" placeholder="Task title..." />
          </div>
          <div>
            <label class="label">Description</label>
            <textarea v-model="createTaskForm.description" rows="6" class="input-field w-full resize-y min-h-[10rem]" placeholder="Describe the task..."></textarea>
          </div>
          <div>
            <label class="label" :class="createTaskForm.task_type === 'FEATURE' ? 'text-gray-500' : ''">
              Estimated Effort (hours)
              <span v-if="createTaskForm.task_type === 'FEATURE'" class="text-gray-600 font-normal">(disabled for Features)</span>
            </label>
            <input
              v-model.number="createTaskForm.estimated_hours"
              type="number"
              min="0"
              step="0.1"
              class="input-field w-full transition-opacity"
              :class="createTaskForm.task_type === 'FEATURE' ? 'opacity-40 cursor-not-allowed' : ''"
              :disabled="createTaskForm.task_type === 'FEATURE'"
              placeholder="e.g. 1.5"
            />
            <p v-if="createTaskForm.task_type !== 'FEATURE'" class="text-sm text-gray-500 mt-2">Hours, up to 1 decimal place (e.g. 1.5). Stored for Manday and Quotation (Costing Engine).</p>
          </div>
          <!-- Sub-task hint -->
          <div v-if="createTaskForm.parent_id" class="p-4 bg-purple-900/20 border border-purple-500/30 rounded-xl text-sm sm:text-base text-purple-300 leading-relaxed">
            This is a sub-task. Dates are inherited from the parent task.
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
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
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
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
            <div v-if="epics.length" class="grid grid-cols-1 gap-4 sm:gap-5">
              <div>
                <label class="label">Due Date</label>
                <input v-model="createTaskForm.due_date" type="datetime-local" class="input-field w-full" />
              </div>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
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
          <div v-if="createTaskError" class="p-4 md:p-5 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-base">{{ createTaskError }}</div>
        </div>
        <div class="flex flex-col-reverse sm:flex-row gap-3 sm:gap-4 px-6 sm:px-8 py-5 sm:py-6 border-t border-gray-700 shrink-0">
          <button
            @click="submitCreateTask"
            :disabled="isCreatingTask || !createTaskForm.title.trim()"
            class="flex-1 btn-primary py-4 text-base sm:text-lg font-semibold rounded-xl disabled:opacity-40 min-h-[3.25rem]"
          >
            {{ isCreatingTask ? 'Creating...' : 'Create Task' }}
          </button>
          <button type="button" @click="closeCreateTaskModal" class="sm:shrink-0 px-6 py-4 bg-gray-700 hover:bg-gray-600 text-gray-200 rounded-xl transition-colors text-base font-medium min-h-[3.25rem]">Cancel</button>
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

    <!-- Edit Project Modal -->
    <div v-if="showEditProjectModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeEditProjectModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">แก้ไขโครงการ</h2>
          <button @click="closeEditProjectModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-400 mb-1.5">ชื่อโครงการ <span class="text-red-400">*</span></label>
            <input v-model="editProjectForm.name" type="text" class="w-full px-4 py-2 bg-gray-900 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500 outline-none" placeholder="Project name (English only)" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-400 mb-1.5">คำอธิบาย</label>
            <textarea v-model="editProjectForm.description" rows="3" class="w-full px-4 py-2 bg-gray-900 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500 outline-none resize-none" placeholder="Description"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-400 mb-1.5">สถานะ</label>
            <select v-model="editProjectForm.status" class="w-full px-4 py-2 bg-gray-900 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500 outline-none">
              <option value="ACTIVE">ACTIVE</option>
              <option value="COMPLETED">COMPLETED</option>
              <option value="ON_HOLD">ON_HOLD</option>
            </select>
          </div>
          <label class="flex items-start gap-2 cursor-pointer">
            <input v-model="editProjectForm.update_code" type="checkbox" class="mt-1 rounded border-gray-600 bg-gray-700 text-purple-500 focus:ring-purple-500" />
            <span class="text-sm text-gray-400">อัปเดตรหัสโครงการ (code) และรหัสงานทั้งหมดตามชื่อใหม่ — ลิงก์เดิม (เช่น /projects/รหัสเก่า) อาจใช้ไม่ได้</span>
          </label>
          <div v-if="editProjectError" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ editProjectError }}</div>
        </div>
        <div class="flex gap-3 mt-5">
          <button @click="saveEditProject" :disabled="isSavingProject || !editProjectForm.name.trim()" class="flex-1 px-5 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-medium rounded-xl transition-colors disabled:opacity-40">
            {{ isSavingProject ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
          <button @click="closeEditProjectModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
        </div>
      </div>
    </div>

    <!-- Edit Task Title Modal -->
    <div v-if="showEditTaskTitleModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeEditTaskTitleModal">
      <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 max-w-lg w-full shadow-2xl">
        <div class="flex items-center justify-between mb-5">
          <h2 class="text-lg font-bold text-white">แก้ไขชื่อ task</h2>
          <button type="button" @click="closeEditTaskTitleModal" class="text-gray-500 hover:text-white">✕</button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-400 mb-1.5">ชื่อ task</label>
            <input v-model="editTaskTitleValue" type="text" class="w-full px-4 py-2 bg-gray-900 border border-gray-600 rounded-lg text-white focus:ring-2 focus:ring-purple-500 focus:border-purple-500 outline-none" placeholder="Task title" @keydown.enter.prevent="saveEditTaskTitle" />
          </div>
        </div>
        <div class="flex gap-3 mt-5">
          <button type="button" @click="saveEditTaskTitle" :disabled="isSavingTaskTitle || !editTaskTitleValue.trim()" class="flex-1 px-5 py-2.5 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-medium rounded-xl transition-colors disabled:opacity-40">
            {{ isSavingTaskTitle ? 'กำลังบันทึก...' : 'บันทึก' }}
          </button>
          <button type="button" @click="closeEditTaskTitleModal" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors">ยกเลิก</button>
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
          <div
            v-if="!editingSprint && project"
            class="rounded-xl border border-purple-500/25 bg-purple-500/5 px-3 py-2.5 text-xs text-gray-400 leading-relaxed"
          >
            <span class="text-gray-200 font-medium">{{ project.name }}</span>
            — New sprint defaults: <span class="text-gray-300">2 weeks</span> (min 1 week), start on a <span class="text-gray-300">Monday 09:00</span>.
            <template v-if="sprints.length">
              If you already have sprints, the start date is placed <span class="text-gray-300">after the latest sprint ends</span> when that is later than “this week’s” Monday.
            </template>
            <template v-else>
              Start is the upcoming Monday (or today if it’s Monday).
            </template>
          </div>
          <div>
            <label class="label">Sprint name *</label>
            <input
              v-model="sprintForm.name"
              type="text"
              class="input-field w-full"
              placeholder="Project — Sprint 1 (Mar 30 – Apr 12, 2026)"
            />
            <div v-if="!editingSprint && project" class="mt-1.5 flex flex-wrap items-center gap-2">
              <button type="button" class="text-xs text-purple-400 hover:text-purple-300" @click="applySuggestedSprintName">
                Use suggested name (project + dates in parentheses)
              </button>
            </div>
          </div>
          <div>
            <label class="label">Sprint goal</label>
            <textarea v-model="sprintForm.goal" rows="2" class="input-field w-full resize-none" placeholder="What will this sprint achieve?"></textarea>
          </div>
          <div v-if="!editingSprint" class="space-y-2">
            <label class="label">Duration</label>
            <select
              :value="sprintDurationWeeks"
              class="input-field w-full"
              @change="onSprintDurationWeeksChange"
            >
              <option v-for="w in sprintDurationOptions" :key="w" :value="w">
                {{ w }} week{{ w === 1 ? '' : 's' }} (Mon–Sun × {{ w }})
              </option>
            </select>
            <p class="text-[11px] text-gray-500">Minimum 1 week. End date updates from start + duration.</p>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="label">Start</label>
              <input
                v-model="sprintForm.start_date"
                type="datetime-local"
                class="input-field w-full"
                @change="syncSprintEndFromStartAndDuration"
              />
            </div>
            <div>
              <label class="label">End</label>
              <input
                v-model="sprintForm.end_date"
                type="datetime-local"
                class="input-field w-full"
              />
            </div>
          </div>
          <div v-if="!editingSprint" class="flex flex-wrap gap-2">
            <button type="button" class="text-xs px-3 py-1.5 rounded-lg border border-gray-600 text-gray-300 hover:bg-gray-700/60" @click="resetSprintModalToDefaults">
              Reset to smart start (after last sprint if any) + 2 weeks + suggested name
            </button>
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
            <input type="checkbox" :value="t.id" v-model="selectedTaskIdsForSprint" class="rounded border-gray-600 bg-gray-700 text-purple-500 focus:ring-purple-500" />
            <span class="text-xs font-mono text-gray-500 shrink-0">{{ t.code }}</span>
            <span class="text-sm text-gray-200 truncate flex-1 min-w-0">{{ t.title }}</span>
            <span class="text-[10px] px-1.5 py-0.5 rounded shrink-0" :class="t.sprint_id ? 'bg-amber-500/20 text-amber-400' : 'bg-gray-600 text-gray-400'">
              {{ t.sprint_id ? (sprints.find(s => s.id === t.sprint_id)?.name ?? 'Sprint อื่น') : 'Backlog' }}
            </span>
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
// Lazy-load heavy tab panels (only load when tab is visited)
const ProjectAnalytics = defineAsyncComponent(() => import('~/components/projects/ProjectAnalytics.vue'))
const QuotationBuilder = defineAsyncComponent(() => import('~/core/modules/pricing/ui/QuotationBuilder.vue'))
const ProjectBackupPanel = defineAsyncComponent(() => import('~/core/modules/projects/ui/ProjectBackupPanel.vue'))
const ProjectCapitalPanel = defineAsyncComponent(() => import('~/core/modules/projects/ui/ProjectCapitalPanel.vue'))
const FeatureRoadmapBoard = defineAsyncComponent(() => import('~/components/tasks/FeatureRoadmapBoard.vue'))
import type { Project, Sprint, Milestone, ProjectAnalytics as AnalyticsType, Task, Epic } from '~/core/modules/projects/infrastructure/projects-api'
import { exportTimelinePdf } from '~/utils/timelinePdfExport'
import { effortHoursToMinutes } from '~/utils/effortHours'
import { buildTaskDisplayCodeMap, sortBacklogTasks, taskCodeSuffix } from '~/utils/backlog-task-utils'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import { useDeploymentApi } from '~/core/modules/deployment/infrastructure/deployment-api'
import { useProjectSprints } from '~/composables/useProjectSprints'
import { useProjectImports } from '~/composables/useProjectImports'

definePageMeta({ layout: 'default', middleware: 'auth' })

const route = useRoute()
const router = useRouter()
const { currentUser } = useAuth()
const projectsApi = useProjectsApi()
const { showError, showSuccess, confirm } = useNotification()
const tasksApi = useTasksApi()
const teamsStore = useTeamsStore()

/** Capital / Costing / Backup — CEO only (not MANAGER / Product Owner / legacy DEV). */
const CEO_ONLY_PROJECT_TAB_IDS = new Set(['capital', 'costing', 'backup'])

const allProjectTabs = [
  { id: 'overview', label: 'Overview', icon: '📊' },
  { id: 'board', label: 'Board', icon: '🗂' },
  { id: 'timeline', label: 'Timeline', icon: '📅' },
  { id: 'backlog', label: 'Backlog', icon: '📋' },
  { id: 'feature-roadmap', label: 'Feature', icon: '🗺️' },
  { id: 'sprints', label: 'Sprints', icon: '🏃' },
  { id: 'analytics', label: 'Analytics', icon: '📈' },
  { id: 'capital', label: 'Capital', icon: '🏦' },
  { id: 'costing', label: 'Costing', icon: '💰' },
  { id: 'backup', label: 'Backup', icon: '🗄' },
]

const isProjectCeo = computed(() => currentUser.value?.role?.toUpperCase() === 'CEO')

const tabs = computed(() => {
  if (isProjectCeo.value) return allProjectTabs
  return allProjectTabs.filter(t => !CEO_ONLY_PROJECT_TAB_IDS.has(t.id))
})

const activeTab = ref('overview')

watch(
  () => [String(route.query.tab || ''), currentUser.value?.role ?? ''] as const,
  ([tabFromRoute, role]) => {
    const raw = tabFromRoute || 'overview'
    const roleUpper = role ? String(role).toUpperCase() : ''
    // Until auth hydrates, do not rewrite CEO-only URLs (would bounce CEOs to overview).
    if (!roleUpper && CEO_ONLY_PROJECT_TAB_IDS.has(raw)) {
      if (activeTab.value !== raw) activeTab.value = raw
      return
    }
    const ceo = roleUpper === 'CEO'
    const next = CEO_ONLY_PROJECT_TAB_IDS.has(raw) && !ceo ? 'overview' : raw
    if (activeTab.value !== next) {
      activeTab.value = next
    }
    if (raw !== next) {
      router.replace({ query: { ...route.query, tab: next } })
    }
  },
  { immediate: true }
)

watch(activeTab, async (tab) => {
  router.replace({ query: { tab } })
  if (tab === 'timeline' && project.value) {
    if (timelineMode.value === 'epic' && !epicTimelineData.value) await loadEpicTimeline()
    else if (timelineMode.value === 'sprint' && !sprintTimelineData.value) await loadSprintTimeline()
    nextTick(() =>
      setTimeout(ganttView.value === 'day' ? scrollTimelineDayFocusRecent : scrollTimelineToToday, 200)
    )
  }
  if (tab === 'analytics' && !analytics.value && project.value) {
    loadAnalytics()
  }
  if (tab === 'backlog') {
    nextTick(() => setupBacklogInfiniteScroll())
  } else if (backlogAutoLoadObserver) {
    backlogAutoLoadObserver.disconnect()
    backlogAutoLoadObserver = null
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
const timelineFullscreen = ref(false)
const timelineRefreshing = ref(false)

watch(timelineMode, async (mode) => {
  if (!project.value) return
  if (mode === 'epic' && !epicTimelineData.value) {
    await loadEpicTimeline()
  } else if (mode === 'sprint' && !sprintTimelineData.value) {
    await loadSprintTimeline()
  }
  nextTick(() =>
    setTimeout(ganttView.value === 'day' ? scrollTimelineDayFocusRecent : scrollTimelineToToday, 200)
  )
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

function addDaysYMD(ymd: string, days: number) {
  const d = new Date(ymd + 'T12:00:00Z')
  d.setUTCDate(d.getUTCDate() + days)
  return toYMD(d.toISOString())
}

/** Persisted end_date is inclusive; Gantt barEnd must be exclusive. */
function inclusiveEndToExclusiveYMD(isoOrYmd: string) {
  return addDaysYMD(toYMD(isoOrYmd), 1)
}

function formatGanttTooltipDate(v: unknown): string {
  const iso = barDateToISO(v)
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

function formatGanttTooltipInclusiveEnd(v: unknown): string {
  const iso = barDateToISO(v)
  if (!iso) return '—'
  const d = new Date(iso)
  d.setUTCDate(d.getUTCDate() - 1)
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

/** Set task start/end to full sprint range so the bar size matches the sprint bar. Returns ISO strings (noon UTC) for API. */
function taskDatesInSprintRange(
  _task: { start_date?: string | null; end_date?: string | null; due_at?: string | null },
  sprint: { start_date: string | null; end_date: string | null }
): { start_date: string; end_date: string } | null {
  if (!sprint?.start_date) return null
  const addDays = (ymd: string, days: number) => {
    const d = new Date(ymd + 'T12:00:00Z')
    d.setUTCDate(d.getUTCDate() + days)
    return toYMD(d.toISOString())
  }
  const spStart = toYMD(sprint.start_date)
  let spEnd = sprint.end_date ? toYMD(sprint.end_date) : addDays(spStart, 14)
  if (spEnd <= spStart) spEnd = addDays(spStart, 1)
  const toNoonUTC = (ymd: string) => ymd + 'T12:00:00.000Z'
  return { start_date: toNoonUTC(spStart), end_date: toNoonUTC(spEnd) }
}

/** For Week view: one column = 7 days (Mon–Sun). Return label "D Mmm – D Mmm" for the timeunit slot. */
function weekRangeLabel(startOfWeek: Date): string {
  const d = new Date(startOfWeek)
  const end = new Date(d)
  end.setDate(end.getDate() + 6)
  const fmt = (x: Date) => x.getDate() + ' ' + x.toLocaleDateString('en-US', { month: 'short' })
  return `${fmt(d)} – ${fmt(end)}`
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
    let end = t.end_date ? inclusiveEndToExclusiveYMD(t.end_date) : (t.due_at ? toYMD(t.due_at) : tomorrow)
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
  const pxPerDay = v === 'month' ? 4 : v === 'week' ? 24 : 72
  const maxChartPx = v === 'day' ? 12000 : 6000
  return Math.max(800, Math.min(maxChartPx, Math.round(days * pxPerDay)))
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
  const pxPerDay = v === 'month' ? 4 : v === 'week' ? 24 : 72
  const maxChartPx = v === 'day' ? 12000 : 6000
  return Math.max(800, Math.min(maxChartPx, Math.round(days * pxPerDay)))
})
const matrixDateRangeStart = computed(() => (matrixChartStart.value ? new Date(matrixChartStart.value + 'T00:00:00Z').toISOString() : ''))
const matrixDateRangeEnd = computed(() => (matrixChartEnd.value ? new Date(matrixChartEnd.value + 'T00:00:00Z').toISOString() : ''))

// Timeline matrix: tasks collapsed by default inside each sprint/epic; click bar to expand/collapse
const timelineExpandedSprints = ref<Record<string, boolean>>({})
const timelineExpandedEpics = ref<Record<string, boolean>>({})

type MatrixGanttRow = { taskId: string; label: string; bars: { barStart: string; barEnd: string; ganttBarConfig: { id: string; label: string; hasHandles: boolean; class?: string } }[] }

/** Safe CSS class suffix from epic id (no leading digit, no special chars) */
function epicBarClassSuffix(epicId: string): string {
  return 'e-' + epicId.replace(/\W/g, '_')
}
const matrixGanttRows = ref<MatrixGanttRow[]>([])

function buildMatrixGanttRows() {
  const today = toYMD(new Date().toISOString())
  const tomorrow = toYMD(new Date(Date.now() + 86400000).toISOString())
  const addDays = (ymd: string, days: number) => {
    const d = new Date(ymd + 'T12:00:00Z')
    d.setUTCDate(d.getUTCDate() + days)
    return toYMD(d.toISOString())
  }
  const rows: MatrixGanttRow[] = []
  if (timelineMode.value === 'epic' && epicTimelineData.value?.epics?.length) {
    const timelineEpicMap = new Map(epicTimelineData.value.epics.map((e) => [e.id, e]))
    const backlogOrder = epics.value
    const epicIdsInOrder = backlogOrder.length
      ? backlogOrder.map((e) => e.id).filter((id) => timelineEpicMap.has(id))
      : epicTimelineData.value.epics.map((e) => e.id)
    type TaskLike = { start_date?: string | null; end_date?: string | null; due_at?: string | null; sub_tasks?: TaskLike[] }
    const taskDateRange = (t: TaskLike): { s: string; e: string } => {
      const s = t.start_date ? toYMD(t.start_date) : (t.due_at ? toYMD(t.due_at) : today)
      let e = t.end_date ? toYMD(t.end_date) : (t.due_at ? toYMD(t.due_at) : tomorrow)
      if (e <= s) e = addDays(s, 1)
      return { s, e }
    }
    const allTaskRangesInEpic = (taskList: TaskLike[]): { s: string; e: string }[] => {
      const out: { s: string; e: string }[] = []
      const walk = (tasks: TaskLike[]) => {
        for (const t of tasks || []) {
          out.push(taskDateRange(t))
          if (t.sub_tasks?.length) walk(t.sub_tasks)
        }
      }
      walk(taskList || [])
      return out
    }
    for (const epicId of epicIdsInOrder) {
      const ep = timelineEpicMap.get(epicId)
      if (!ep) continue
      const tasks = ep.tasks || []
      const ranges = allTaskRangesInEpic(tasks)
      let epStart: string
      let epEnd: string
      if (ranges.length > 0) {
        let minT = ranges[0].s
        let maxT = ranges[0].e
        for (const r of ranges) {
          if (r.s < minT) minT = r.s
          if (r.e > maxT) maxT = r.e
        }
        epStart = minT
        epEnd = maxT <= minT ? addDays(minT, 1) : maxT
      } else {
        // No tasks: epic bar length follows "tasks inside" = minimal one-day bar (not epic's own dates)
        epStart = today
        epEnd = tomorrow
        if (epEnd <= epStart) epEnd = addDays(epStart, 1)
      }
      const taskCount = tasks.length
      const expanded = timelineExpandedEpics.value[ep.id]
      const toggle = expanded ? '▼' : '▶'
      const epicBarClass = `gantt-bar-epic gantt-bar-epic-${epicBarClassSuffix(ep.id)}`
      // Always show epic bar (both when collapsed and when expanded).
      const epicBar = [{ barStart: epStart, barEnd: epEnd, ganttBarConfig: { id: `epic-${ep.id}`, label: ep.title, hasHandles: true, class: epicBarClass } }]
      rows.push({
        taskId: `epic-${ep.id}`,
        label: `${toggle} 📁 ${ep.title}${taskCount ? ` (${taskCount})` : ''}`,
        bars: epicBar,
      })
      if (expanded) {
        const taskBarClass = `gantt-bar-task gantt-bar-task-epic-${epicBarClassSuffix(ep.id)}`
        for (const task of ep.tasks || []) {
          let start = task.start_date ? toYMD(task.start_date) : (task.due_at ? toYMD(task.due_at) : today)
          let end = task.end_date ? inclusiveEndToExclusiveYMD(task.end_date) : (task.due_at ? toYMD(task.due_at) : tomorrow)
          if (end <= start) end = addDays(start, 1)
          const label = task.title || ''
          rows.push({
            taskId: task.id,
            label: `  ${label}`,
            bars: [{ barStart: start, barEnd: end, ganttBarConfig: { id: task.id, label, hasHandles: true, class: taskBarClass } }],
          })
        }
      }
    }
  } else if (timelineMode.value === 'sprint' && sprintTimelineData.value?.sprints?.length) {
    for (const sp of sprintTimelineData.value.sprints) {
      const spStart = sp.start_date ? toYMD(sp.start_date) : today
      let spEnd = sp.end_date ? inclusiveEndToExclusiveYMD(sp.end_date) : tomorrow
      if (spEnd <= spStart) spEnd = addDays(spStart, 1)
      const taskCount = (sp.tasks || []).length
      const expanded = timelineExpandedSprints.value[sp.id]
      const toggle = expanded ? '▼' : '▶'
      rows.push({
        taskId: `sprint-${sp.id}`,
        label: `${toggle} 🏃 ${sp.name}${taskCount ? ` (${taskCount})` : ''}`,
        bars: [{ barStart: spStart, barEnd: spEnd, ganttBarConfig: { id: `sprint-${sp.id}`, label: sp.name, hasHandles: true, class: 'gantt-bar-sprint' } }],
      })
      if (expanded) {
        for (const task of sp.tasks || []) {
          let start = task.start_date ? toYMD(task.start_date) : (task.due_at ? toYMD(task.due_at) : today)
          let end = task.end_date ? inclusiveEndToExclusiveYMD(task.end_date) : (task.due_at ? toYMD(task.due_at) : tomorrow)
          if (end <= start) end = addDays(start, 1)
          const label = task.title || ''
          rows.push({
            taskId: task.id,
            label: `  ${label}`,
            bars: [{ barStart: start, barEnd: end, ganttBarConfig: { id: task.id, label, hasHandles: true, class: 'gantt-bar-task' } }],
          })
        }
      }
    }
  }
  matrixGanttRows.value = rows
}

watch(
  () => [
    timelineMode.value,
    epicTimelineData.value,
    sprintTimelineData.value,
    timelineExpandedEpics.value,
    timelineExpandedSprints.value,
  ],
  () => buildMatrixGanttRows(),
  { deep: true }
)

/** Dynamic CSS for epic (and its tasks) Gantt bar colors – same color for epic and tasks under it. */
const epicBarStyles = computed(() => {
  if (timelineMode.value !== 'epic' || !epicTimelineData.value?.epics?.length) return ''
  const lines: string[] = []
  for (const ep of epicTimelineData.value.epics) {
    const color = ep.color || '#6366f1'
    const safeId = epicBarClassSuffix(ep.id)
    const colorAlpha = color.length === 7 ? `${color}80` : color
    const shadowAlpha = color.length === 7 ? `${color}66` : color
    // Epic bar
    lines.push(
      `.gantt-enterprise .g-gantt-bar.gantt-bar-epic.gantt-bar-epic-${safeId} { background: ${color} !important; border: 1px solid ${colorAlpha} !important; }`,
      `.gantt-enterprise .g-gantt-bar.gantt-bar-epic.gantt-bar-epic-${safeId}:hover { transform: translateY(-1px); box-shadow: 0 4px 12px ${shadowAlpha}; }`
    )
    // Task bars under this epic – same color
    lines.push(
      `.gantt-enterprise .g-gantt-bar.gantt-bar-task.gantt-bar-task-epic-${safeId} { background: ${color} !important; border: 1px solid ${colorAlpha} !important; }`,
      `.gantt-enterprise .g-gantt-bar.gantt-bar-task.gantt-bar-task-epic-${safeId}:hover { transform: translateY(-1px); box-shadow: 0 4px 12px ${shadowAlpha}; }`
    )
  }
  return lines.join('\n')
})

const milestoneDragPosition = ref<{ id: string; left: number } | null>(null)

function onMilestoneDragMove(payload: { milestoneId: string; leftPx: number }) {
  milestoneDragPosition.value = { id: payload.milestoneId, left: payload.leftPx }
}

const matrixMilestoneLinePositions = computed(() => {
  if (!matrixChartStart.value || !matrixChartEnd.value || matrixChartWidth.value <= 0) return []
  const start = toLocalMidnight(matrixChartStart.value)
  const end = toLocalMidnight(matrixChartEnd.value)
  if (end <= start) return []
  const gridOffset = 220
  const chartContentWidth = matrixChartWidth.value
  const drag = milestoneDragPosition.value
  const list = milestones.value
    .filter((m): m is Milestone & { due_date: string } => !!m.due_date)
    .sort((a, b) => toLocalMidnight(a.due_date) - toLocalMidnight(b.due_date))
  return list.map((m) => {
    if (drag?.id === m.id) return { id: m.id, left: drag.left }
    const date = toLocalMidnight(m.due_date)
    const pct = Math.max(0, Math.min(1, (date - start) / (end - start)))
    const left = gridOffset + pct * chartContentWidth
    return { id: m.id, left }
  })
})

const ganttBarJustDragged = ref(false)
const matrixHoveredBar = ref<any | null>(null)
const matrixTooltipPos = ref({ x: 0, y: 0 })

function onMatrixGanttMouseEnter(payload: { bar: any; e: MouseEvent }) {
  matrixHoveredBar.value = payload?.bar || null
  matrixTooltipPos.value = {
    x: (payload?.e?.clientX || 0) + 12,
    y: (payload?.e?.clientY || 0) + 12,
  }
}

function onMatrixGanttMouseLeave() {
  matrixHoveredBar.value = null
}

function onMatrixGanttDragStart() {
  ganttBarJustDragged.value = true
}

function onMatrixGanttClickBar(payload: { bar: { ganttBarConfig: { id: string } } }) {
  if (ganttBarJustDragged.value) {
    ganttBarJustDragged.value = false
    return
  }
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

function expandAllTimelineTasks() {
  if (timelineMode.value === 'epic' && epicTimelineData.value?.epics?.length) {
    const next: Record<string, boolean> = {}
    for (const ep of epicTimelineData.value.epics) next[ep.id] = true
    timelineExpandedEpics.value = next
  } else if (timelineMode.value === 'sprint' && sprintTimelineData.value?.sprints?.length) {
    const next: Record<string, boolean> = {}
    for (const sp of sprintTimelineData.value.sprints) next[sp.id] = true
    timelineExpandedSprints.value = next
  }
}

function collapseAllTimelineTasks() {
  if (timelineMode.value === 'epic' && epicTimelineData.value?.epics?.length) {
    const next: Record<string, boolean> = {}
    for (const ep of epicTimelineData.value.epics) next[ep.id] = false
    timelineExpandedEpics.value = next
  } else if (timelineMode.value === 'sprint' && sprintTimelineData.value?.sprints?.length) {
    const next: Record<string, boolean> = {}
    for (const sp of sprintTimelineData.value.sprints) next[sp.id] = false
    timelineExpandedSprints.value = next
  }
}

function onGanttClickBar(payload: { bar: { ganttBarConfig: { id: string } } }) {
  const id = payload?.bar?.ganttBarConfig?.id
  if (id) router.push(taskUrl(id))
}

function barDateToISO(v: unknown): string {
  if (typeof v === 'string') {
    if (/^\d{4}-\d{2}-\d{2}$/.test(v)) return new Date(v + 'T12:00:00Z').toISOString()
    return new Date(v).toISOString()
  }
  if (v instanceof Date) return v.toISOString()
  return ''
}

/**
 * Vue-Ganttastic emits barEnd as exclusive boundary in day/week views.
 * Persist inclusive end_date by subtracting 1 day from barEnd.
 */
function barEndToInclusiveISO(v: unknown): string {
  const iso = barDateToISO(v)
  if (!iso) return ''
  const d = new Date(iso)
  d.setUTCDate(d.getUTCDate() - 1)
  return d.toISOString()
}

/** อัปเดตข้อมูล timeline ในเครื่อง (ไม่โหลดใหม่) เพื่อไม่ให้ scroll กระโดด */
function updateEpicInTimelineData(epicId: string, payload: { start_date?: string; end_date?: string }) {
  const epics = epicTimelineData.value?.epics
  if (!epics) return
  const ep = epics.find((e) => e.id === epicId)
  if (ep) {
    if (payload.start_date != null) ep.start_date = payload.start_date
    if (payload.end_date != null) ep.end_date = payload.end_date
  }
}

function updateSprintInTimelineData(sprintId: string, payload: { start_date?: string; end_date?: string }) {
  const sprints = sprintTimelineData.value?.sprints
  if (!sprints) return
  const sp = sprints.find((s) => s.id === sprintId)
  if (sp) {
    if (payload.start_date != null) sp.start_date = payload.start_date
    if (payload.end_date != null) sp.end_date = payload.end_date
  }
}

function updateTaskInTimelineData(taskId: string, start_date: string, end_date: string) {
  for (const ep of epicTimelineData.value?.epics ?? []) {
    const t = ep.tasks?.find((x) => x.id === taskId)
    if (t) {
      t.start_date = start_date
      t.end_date = end_date
      break
    }
  }
  for (const sp of sprintTimelineData.value?.sprints ?? []) {
    const t = sp.tasks?.find((x) => x.id === taskId)
    if (t) {
      t.start_date = start_date
      t.end_date = end_date
      break
    }
  }
}

/** คำนวณ start/end ของ epic ให้ครอบคลุมทุก task ใน epic นั้น แล้วอัปเดต API + ข้อมูลในเครื่อง */
async function syncEpicDatesFromTasks(epicId: string) {
  const tasksInEpic = allTasks.value.filter((t) => t.epic_id === epicId && !t.parent_id)
  if (tasksInEpic.length === 0) return
  let minStart: Date | null = null
  let maxEnd: Date | null = null
  for (const t of tasksInEpic) {
    const s = t.start_date ? new Date(t.start_date) : (t.due_at ? new Date(t.due_at) : null)
    const e = t.end_date ? new Date(t.end_date) : (t.due_at ? new Date(t.due_at) : null)
    if (s) minStart = minStart == null ? s : (s < minStart ? s : minStart)
    if (e) maxEnd = maxEnd == null ? e : (e > maxEnd ? e : maxEnd)
  }
  if (minStart == null || maxEnd == null) return
  if (maxEnd <= minStart) maxEnd = new Date(minStart.getTime() + 86400000)
  const startStr = minStart.toISOString()
  const endStr = maxEnd.toISOString()
  try {
    await projectsApi.updateEpic(epicId, { start_date: startStr, end_date: endStr })
    updateEpicInTimelineData(epicId, { start_date: startStr, end_date: endStr })
  } catch (err) {
    console.error('Failed to sync epic dates from tasks:', err)
  }
}

function getTimelineScrollElement(): HTMLElement | null {
  const wrapper = timelineScrollWrapperRef.value
  if (!wrapper) return null
  if (timelineFullscreen.value) {
    const inner = wrapper.querySelector('.timeline-fullscreen-scroll')
    return inner as HTMLElement | null
  }
  return wrapper
}

function onTimelinePanStart(e: MouseEvent) {
  if ((e.target as HTMLElement).closest('.g-gantt-bar, .g-gantt-bar-handle-left, .g-gantt-bar-handle-right, .g-label-column, .milestone-marker, button, a')) return
  const scrollEl = getTimelineScrollElement()
  if (!scrollEl) return
  const startX = e.clientX
  const startY = e.clientY
  const startScrollLeft = scrollEl.scrollLeft
  const startScrollTop = scrollEl.scrollTop
  const onMove = (e2: MouseEvent) => {
    scrollEl.scrollLeft = Math.max(0, startScrollLeft + (startX - e2.clientX))
    if (timelineFullscreen.value) scrollEl.scrollTop = Math.max(0, startScrollTop + (startY - e2.clientY))
  }
  const onUp = () => {
    document.removeEventListener('mousemove', onMove)
    document.removeEventListener('mouseup', onUp)
    document.body.style.removeProperty('cursor')
    document.body.style.removeProperty('user-select')
  }
  document.addEventListener('mousemove', onMove)
  document.addEventListener('mouseup', onUp)
  document.body.style.cursor = 'grabbing'
  document.body.style.userSelect = 'none'
}

function onTimelinePanStartTouch(e: TouchEvent) {
  if ((e.target as HTMLElement).closest('.g-gantt-bar, .g-gantt-bar-handle-left, .g-gantt-bar-handle-right, .g-label-column, .milestone-marker, button, a')) return
  const scrollEl = getTimelineScrollElement()
  if (!scrollEl || e.touches.length !== 1) return
  const startX = e.touches[0].clientX
  const startY = e.touches[0].clientY
  const startScrollLeft = scrollEl.scrollLeft
  const startScrollTop = scrollEl.scrollTop
  const onMove = (e2: TouchEvent) => {
    if (e2.touches.length !== 1) return
    e2.preventDefault()
    scrollEl.scrollLeft = Math.max(0, startScrollLeft + (startX - e2.touches[0].clientX))
    if (timelineFullscreen.value) scrollEl.scrollTop = Math.max(0, startScrollTop + (startY - e2.touches[0].clientY))
  }
  const onEnd = () => {
    document.removeEventListener('touchmove', onMove, { capture: true })
    document.removeEventListener('touchend', onEnd)
  }
  document.addEventListener('touchmove', onMove, { capture: true, passive: false })
  document.addEventListener('touchend', onEnd)
}

async function onGanttDragEnd(payload: { bar?: { barStart?: unknown; barEnd?: unknown; ganttBarConfig?: { id?: string } }; movedBars?: Map<unknown, unknown> }) {
  ganttBarJustDragged.value = true
  const map = payload.movedBars
  if (!map || map.size === 0) {
    nextTick(() => { ganttBarJustDragged.value = false })
    return
  }
  const projectId = project.value?.id
  if (!projectId) {
    nextTick(() => { ganttBarJustDragged.value = false })
    return
  }
  const affectedEpicIds = new Set<string>()
  for (const [barObj, _old] of map) {
    const bar = barObj as { barStart?: unknown; barEnd?: unknown; ganttBarConfig?: { id?: string } }
    const barId = bar?.ganttBarConfig?.id
    if (!barId) continue
    const start = barDateToISO(bar.barStart)
    const end = barEndToInclusiveISO(bar.barEnd)
    if (!start || !end) continue
    try {
      if (barId.startsWith('epic-')) {
        const epicId = barId.slice('epic-'.length)
        await projectsApi.updateEpic(epicId, { start_date: start, end_date: end })
        updateEpicInTimelineData(epicId, { start_date: start, end_date: end })
      } else if (barId.startsWith('sprint-')) {
        const sprintId = barId.slice('sprint-'.length)
        const toNoonUTC = (iso: string) => toYMD(iso) + 'T12:00:00.000Z'
        const startNorm = toNoonUTC(start)
        const endNorm = toNoonUTC(end)
        const newStartMs = new Date(startNorm).getTime()
        const newEndMs = new Date(endNorm).getTime()
        const newDurationMs = Math.max(newEndMs - newStartMs, 86400000)

        const oldSprint = sprintTimelineData.value?.sprints?.find((s) => s.id === sprintId) ?? sprints.value.find((s) => s.id === sprintId)
        const oldStartYMD = oldSprint?.start_date ? toYMD(oldSprint.start_date) : null
        const oldEndYMD = oldSprint?.end_date ? toYMD(oldSprint.end_date) : null
        const oldStartMs = oldStartYMD ? new Date(oldStartYMD + 'T12:00:00.000Z').getTime() : newStartMs
        const oldEndMs = oldEndYMD ? new Date(oldEndYMD + 'T12:00:00.000Z').getTime() : newEndMs
        const oldDurationMs = Math.max(oldEndMs - oldStartMs, 86400000)

        const updatedSprint = await projectsApi.updateSprint(sprintId, { start_date: startNorm, end_date: endNorm })
        const sprintIdx = sprints.value.findIndex((s) => s.id === sprintId)
        if (sprintIdx !== -1) sprints.value[sprintIdx] = { ...sprints.value[sprintIdx], ...updatedSprint }
        updateSprintInTimelineData(sprintId, { start_date: startNorm, end_date: endNorm })

        const tasksInSprint = allTasks.value.filter((t) => t.sprint_id === sprintId)
        for (const t of tasksInSprint) {
          const hasStart = t.start_date != null && t.start_date !== ''
          const hasEnd = t.end_date != null && t.end_date !== ''
          let tStartMs: number
          let tEndMs: number
          if (hasStart && hasEnd) {
            tStartMs = new Date(toYMD(t.start_date!) + 'T12:00:00.000Z').getTime()
            tEndMs = new Date(toYMD(t.end_date!) + 'T12:00:00.000Z').getTime()
          } else {
            tStartMs = oldStartMs
            tEndMs = oldEndMs
          }
          let ratioStart = (tStartMs - oldStartMs) / oldDurationMs
          let ratioEnd = (tEndMs - oldStartMs) / oldDurationMs
          ratioStart = Math.max(0, Math.min(1, ratioStart))
          ratioEnd = Math.max(0, Math.min(1, ratioEnd))
          if (ratioEnd <= ratioStart) ratioEnd = Math.min(1, ratioStart + 1 / 7)
          const newTStartMs = newStartMs + ratioStart * (newEndMs - newStartMs)
          const newTEndMs = newStartMs + ratioEnd * (newEndMs - newStartMs)
          const startVal = toNoonUTC(new Date(newTStartMs).toISOString())
          const endVal = toNoonUTC(new Date(newTEndMs).toISOString())
          try {
            await tasksApi.updateTask(t.id, { start_date: startVal, end_date: endVal })
            const idx = allTasks.value.findIndex((x) => x.id === t.id)
            if (idx !== -1) {
              allTasks.value[idx] = { ...allTasks.value[idx], start_date: startVal, end_date: endVal }
            }
            updateTaskInTimelineData(t.id, startVal, endVal)
          } catch (e) {
            console.error('Failed to scale task dates:', t.id, e)
          }
        }
      } else {
        await tasksApi.updateTask(barId, { start_date: start, end_date: end })
        const idx = allTasks.value.findIndex((t) => t.id === barId)
        if (idx !== -1) {
          allTasks.value[idx] = {
            ...allTasks.value[idx],
            start_date: start,
            end_date: end,
          }
          const epicId = allTasks.value[idx].epic_id
          if (epicId) affectedEpicIds.add(epicId)
        }
        updateTaskInTimelineData(barId, start, end)
      }
    } catch (e) {
      console.error('Failed to update dates after drag/resize:', e)
    }
  }
  for (const epicId of affectedEpicIds) {
    await syncEpicDatesFromTasks(epicId)
  }
  buildMatrixGanttRows()
  nextTick(() => setTimeout(() => { ganttBarJustDragged.value = false }, 0))
}

function scrollTimelineToToday() {
  nextTick(() => {
    const scrollEl = getTimelineScrollElement()
    const width = matrixChartWidth.value
    if (!scrollEl || width <= 0) return
    const start = toLocalMidnight(matrixChartStart.value)
    const end = toLocalMidnight(matrixChartEnd.value)
    const now = Date.now()
    const pct = end > start ? Math.max(0, Math.min(1, (now - start) / (end - start))) : 0.5
    const todayOffsetFromLeft = 0.18
    const todayLeftPx = 220 + pct * width
    const targetScroll = Math.max(0, todayLeftPx - scrollEl.clientWidth * todayOffsetFromLeft)
    const maxScroll = Math.max(0, scrollEl.scrollWidth - scrollEl.clientWidth)
    scrollEl.scrollLeft = Math.min(targetScroll, maxScroll)
  })
}

/** Day view: align scroll so the window shows from (today − 5 days) through today (local dates). */
function scrollTimelineDayFocusRecent() {
  nextTick(() => {
    const scrollEl = getTimelineScrollElement()
    const width = matrixChartWidth.value
    if (!scrollEl || width <= 0) return
    const rangeStart = toLocalMidnight(matrixChartStart.value)
    const rangeEnd = toLocalMidnight(matrixChartEnd.value)
    if (rangeEnd <= rangeStart) return
    const focus = new Date()
    focus.setHours(0, 0, 0, 0)
    focus.setDate(focus.getDate() - 5)
    const focusStart = focus.getTime()
    const clamped = Math.max(rangeStart, Math.min(focusStart, rangeEnd))
    const pct = (clamped - rangeStart) / (rangeEnd - rangeStart)
    const gridLeftPx = 220 + pct * width
    const padding = 8
    const targetScroll = Math.max(0, gridLeftPx - padding)
    const maxScroll = Math.max(0, scrollEl.scrollWidth - scrollEl.clientWidth)
    scrollEl.scrollLeft = Math.min(targetScroll, maxScroll)
  })
}

async function onExportTimelinePdf() {
  if (!project.value) return
  await exportTimelinePdf(
    project.value,
    timelineMode.value,
    epicTimelineData.value,
    sprintTimelineData.value
  )
}

// State
const project = ref<Project | null>(null)
const allTasks = ref<Task[]>([])
const sprints = ref<Sprint[]>([])
const milestones = ref<Milestone[]>([])
const epics = ref<Epic[]>([])
const analytics = ref<AnalyticsType | null>(null)
const isLoading = ref(true)
// Task IDs that already have a deployment request (used to show warning on WAIT_FOR_DEPLOY kanban cards)
const deployedTaskIds = ref<string[]>([])
const projectDetailsTasksMeta = ref<{ limit: number; returned: number; has_more: boolean } | null>(null)
const projectTasksNextCursor = ref<{ created_at: string; id: string } | null>(null)
const projectTasksNextOffset = ref<number>(0)
const isLoadingMoreProjectTasks = ref(false)
const backlogAutoLoadSentinelRef = ref<HTMLElement | null>(null)
let backlogAutoLoadObserver: IntersectionObserver | null = null
const analyticsLoading = ref(false)
const error = ref('')
const ganttView = ref('week')

watch(ganttView, (v, prev) => {
  if (v !== 'day' || prev === 'day') return
  nextTick(() => setTimeout(() => scrollTimelineDayFocusRecent(), 200))
})

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
const totalTasks = computed(() => allTasks.value.length)
const completedCount = computed(() => allTasks.value.filter((t) => t.status === 'COMPLETED').length)
const inProgressCount = computed(() => allTasks.value.filter((t) => t.status === 'IN_PROGRESS').length)
const completionPct = computed(() => totalTasks.value ? Math.round((completedCount.value / totalTasks.value) * 100) : 0)
const overdueCount = computed(() => {
  const now = Date.now()
  return allTasks.value.filter((t) => t.status !== 'COMPLETED' && t.due_at && new Date(t.due_at).getTime() < now).length
})
const epicTasks = computed(() => allTasks.value.filter((t) => !t.parent_id))

/** Backlog multi-select: task IDs selected for bulk delete */
const backlogSelectedTaskIds = ref<Set<string>>(new Set())
const backlogSelectedCount = computed(() => backlogSelectedTaskIds.value.size)
function isBacklogTaskSelected(id: string) {
  return backlogSelectedTaskIds.value.has(id)
}
function toggleBacklogTaskSelection(id: string) {
  const next = new Set(backlogSelectedTaskIds.value)
  if (next.has(id)) next.delete(id)
  else next.add(id)
  backlogSelectedTaskIds.value = next
}
function clearBacklogSelection() {
  backlogSelectedTaskIds.value = new Set()
}

const isBulkDeletingBacklog = ref(false)
async function bulkDeleteSelectedBacklogTasks() {
  const ids = [...backlogSelectedTaskIds.value]
  if (ids.length === 0) return
  const ok = await confirm({
    title: 'ลบ tasks ที่เลือก',
    message: `ต้องการลบ ${ids.length} รายการใช่หรือไม่? (ถ้ามี sub-task ต้องลบ sub-task ก่อน)`,
    confirmLabel: 'ลบ',
    cancelLabel: 'ยกเลิก',
    variant: 'danger',
  })
  if (!ok) return
  isBulkDeletingBacklog.value = true
  const byId = Object.fromEntries(allTasks.value.map((t) => [t.id, t]))
  const remaining = new Set(ids)
  const toDelete: string[] = []
  while (remaining.size > 0) {
    const leaves = [...remaining].filter((id) => ![...remaining].some((oid) => byId[oid]?.parent_id === id))
    if (leaves.length === 0) break
    leaves.forEach((id) => {
      toDelete.push(id)
      remaining.delete(id)
    })
  }
  let deleted = 0
  const errors: string[] = []
  for (const id of toDelete) {
    try {
      await tasksApi.deleteTask(id)
      deleted++
    } catch (e: any) {
      const msg = e?.data?.message || e?.message || 'ลบไม่สำเร็จ'
      errors.push(`${byId[id]?.code || id}: ${msg}`)
    }
  }
  isBulkDeletingBacklog.value = false
  clearBacklogSelection()
  await loadAll()
  if (errors.length > 0) {
    showError(errors.slice(0, 5).join('; ') + (errors.length > 5 ? ` และอีก ${errors.length - 5} รายการ` : ''), 'ลบบางรายการไม่สำเร็จ')
  }
  if (deleted > 0) {
    showSuccess(`ลบแล้ว ${deleted} รายการ`, 'ลบสำเร็จ')
  }
}
const recentTasks = computed(() =>
  [...allTasks.value]
    .sort((a, b) => new Date(b.updated_at).getTime() - new Date(a.updated_at).getTime())
    .slice(0, 10)
)

function getSubTasks(parentId: string) {
  return allTasks.value.filter((t) => t.parent_id === parentId)
}

/** Push task and all descendants in display order (so A001, B001, C001… are assigned correctly). */
function pushTaskAndDescendants(list: Task[], task: Task) {
  list.push(task)
  getSubTasks(task.id).forEach((child) => pushTaskAndDescendants(list, child))
}

/** Backlog display order: epics (with tasks + all subs including C) then unassigned. Used for A001/B001/C001. */
const allTasksInBacklogOrder = computed(() => {
  const list: typeof allTasks.value = []
  epics.value.forEach((ep) => {
    getTasksForEpic(ep.id).forEach((t) => pushTaskAndDescendants(list, t))
  })
  getUnassignedTasks().forEach((t) => pushTaskAndDescendants(list, t))
  return list
})

/** Top-level: A084. Sub-tasks: B084. Sub-tasks of B: C084. Uses task.code suffix so backlog and task page match. */
const taskDisplayCodeMap = computed(() =>
  buildTaskDisplayCodeMap(allTasksInBacklogOrder.value, allTasks.value)
)

function taskDisplayCode(task: { id: string; code?: string }) {
  return taskDisplayCodeMap.value[task.id] ?? taskCodeSuffix(task.code)
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
  if (status === 'READY_FOR_TEST') return 'bg-cyan-500/20 text-cyan-400'
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

/** Convert ISO date string (UTC) to "YYYY-MM-DDTHH:mm" in local time for datetime-local input */
function isoToDatetimeLocal(iso: string | null | undefined): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const h = String(d.getHours()).padStart(2, '0')
  const min = String(d.getMinutes()).padStart(2, '0')
  return `${y}-${m}-${day}T${h}:${min}`
}


const BACKLOG_EXPANDED_STORAGE_KEY = 'sentinel-backlog-expanded'
const BACKLOG_EXPECT_RETURN_KEY = 'sentinel-backlog-expect-return'

function taskUrl(taskId: string) {
  const projectId = route.params.id as string
  const tab = activeTab.value
  return { path: `/task/${taskId}`, query: { from_project: projectId, from_tab: tab } }
}

/** Scroll container is <main> in default layout (overflow-auto), not window. */
function getMainScrollEl(): HTMLElement | null {
  if (typeof document === 'undefined') return null
  return document.querySelector('main')
}

function saveBacklogExpandedState() {
  if (typeof sessionStorage === 'undefined' || !project.value) return
  const key = `${BACKLOG_EXPANDED_STORAGE_KEY}-${project.value.id}`
  try {
    const main = getMainScrollEl()
    const scrollTop = main ? main.scrollTop : 0
    const scrollLeft = main ? main.scrollLeft : 0
    sessionStorage.setItem(key, JSON.stringify({
      expandedEpics: { ...expandedEpics.value },
      expandedEpicGroups: { ...expandedEpicGroups.value },
      scrollTop,
      scrollLeft,
    }))
    sessionStorage.setItem(BACKLOG_EXPECT_RETURN_KEY, project.value.id)
  } catch {
    // ignore quota or parse errors
  }
}

/** Restores expanded state from sessionStorage. Returns saved scroll position if any (for caller to apply after paint). */
function restoreBacklogExpandedState(projectId: string): { scrollTop: number; scrollLeft: number } | null {
  if (typeof sessionStorage === 'undefined') return null
  const key = `${BACKLOG_EXPANDED_STORAGE_KEY}-${projectId}`
  try {
    const raw = sessionStorage.getItem(key)
    if (!raw) return null
    const data = JSON.parse(raw) as {
      expandedEpics?: Record<string, boolean>
      expandedEpicGroups?: Record<string, boolean>
      scrollTop?: number
      scrollLeft?: number
    }
    if (data.expandedEpics && typeof data.expandedEpics === 'object') {
      expandedEpics.value = { ...data.expandedEpics }
    }
    if (data.expandedEpicGroups && typeof data.expandedEpicGroups === 'object') {
      expandedEpicGroups.value = { ...data.expandedEpicGroups }
    }
    const scrollTop = typeof data.scrollTop === 'number' ? data.scrollTop : 0
    const scrollLeft = typeof data.scrollLeft === 'number' ? data.scrollLeft : 0
    if (scrollTop > 0 || scrollLeft > 0) {
      return { scrollTop, scrollLeft }
    }
    return null
  } catch {
    // ignore parse errors
    return null
  }
}

function navigateToTask(id: string) {
  if (activeTab.value === 'backlog') saveBacklogExpandedState()
  router.push(taskUrl(id))
}

function toggleEpic(id: string) {
  expandedEpics.value[id] = !expandedEpics.value[id]
}

// Reload all project data after a backup restore
async function onProjectRestored() {
  await loadAll()
}

// Load data — use combined details endpoint (1 round-trip) for fast initial load.
// When tab=timeline: run details + timeline API in parallel; end full-page spinner as soon as details land (timeline shows its own loader).
async function loadAll() {
  isLoading.value = true
  error.value = ''
  const idOrCode = route.params.id as string
  const isTimelineTab = activeTab.value === 'timeline'
  let timelineOk = false
  try {
    // When timeline tab: fetch details and timeline in parallel (both support id/code)
    const isHeavyBoardTab = activeTab.value === 'backlog' || activeTab.value === 'board' || activeTab.value === 'timeline'
    const detailsPromise = projectsApi.getProjectDetails(idOrCode, { tasksLimit: isHeavyBoardTab ? 1200 : 600 })
    const timelinePromise =
      isTimelineTab
        ? (timelineMode.value === 'epic'
            ? projectsApi.getEpicTimelineData(idOrCode)
            : projectsApi.getSprintTimelineData(idOrCode))
        : null

    const details = await detailsPromise
    project.value = details.project
    allTasks.value = details.tasks
    projectDetailsTasksMeta.value = details.tasks_meta ?? null
    if (details.tasks.length > 0) {
      const last = details.tasks[details.tasks.length - 1]
      projectTasksNextCursor.value = { created_at: last.created_at, id: last.id }
      projectTasksNextOffset.value = details.tasks.length
    } else {
      projectTasksNextCursor.value = null
      projectTasksNextOffset.value = 0
    }
    sprints.value = details.sprints
    milestones.value = details.milestones
    epics.value = details.epics
    // Default: collapse all epic groups and tasks (user expands to see). On refresh we keep collapsed; only restore when returning from task page.
    details.epics.forEach((ep) => { expandedEpicGroups.value[ep.id] = false })
    expandedEpicGroups.value['__unassigned__'] = false
    expandedEpics.value = {}
    const expectReturn = typeof sessionStorage !== 'undefined' ? sessionStorage.getItem(BACKLOG_EXPECT_RETURN_KEY) : null
    const shouldRestore = expectReturn === details.project.id
    if (shouldRestore && typeof sessionStorage !== 'undefined') sessionStorage.removeItem(BACKLOG_EXPECT_RETURN_KEY)
    const savedScroll = shouldRestore ? restoreBacklogExpandedState(details.project.id) : null
    if (savedScroll && activeTab.value === 'backlog') {
      nextTick(() => {
        const apply = () => {
          const main = getMainScrollEl()
          if (main) {
            main.scrollTop = savedScroll!.scrollTop
            main.scrollLeft = savedScroll!.scrollLeft
          }
        }
        requestAnimationFrame(() => {
          requestAnimationFrame(() => {
            apply()
            setTimeout(apply, 80)
          })
        })
      })
    }

    isLoading.value = false

    // Fetch deployment requests for WAIT_FOR_DEPLOY tasks (non-blocking)
    const waitForDeployTaskIds = details.tasks
      .filter((t) => t.status === 'WAIT_FOR_DEPLOY')
      .map((t) => t.id)
    if (waitForDeployTaskIds.length > 0) {
      const depApi = useDeploymentApi()
      const results = await Promise.allSettled(
        waitForDeployTaskIds.map((id) => depApi.getByTaskId(id))
      )
      deployedTaskIds.value = waitForDeployTaskIds.filter((_, i) => {
        const r = results[i]
        return r.status === 'fulfilled' && r.value !== null
      })
    } else {
      deployedTaskIds.value = []
    }

    // Apply timeline data if we fetched in parallel (do not block project shell on this)
    if (timelinePromise) {
      matrixTimelineLoading.value = true
      try {
        const data = await timelinePromise
        if (timelineMode.value === 'epic') epicTimelineData.value = data as typeof epicTimelineData.value
        else sprintTimelineData.value = data as typeof sprintTimelineData.value
        timelineOk = true
        nextTick(() => buildMatrixGanttRows())
        nextTick(() =>
          setTimeout(ganttView.value === 'day' ? scrollTimelineDayFocusRecent : scrollTimelineToToday, 200)
        )
      } catch (e) {
        console.error('Failed to load timeline:', e)
      } finally {
        matrixTimelineLoading.value = false
      }
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load project'
    isLoading.value = false
  }
  // When timeline tab but we didn't get timeline data (no parallel or parallel failed): load now
  if (activeTab.value === 'timeline' && !timelineOk && project.value) {
    if (timelineMode.value === 'epic') await loadEpicTimeline()
    else if (timelineMode.value === 'sprint') await loadSprintTimeline()
  }
  if (activeTab.value === 'analytics') loadAnalytics()
}

async function loadMoreProjectTasks() {
  if (!project.value || isLoadingMoreProjectTasks.value) return
  if (!projectDetailsTasksMeta.value?.has_more) return
  isLoadingMoreProjectTasks.value = true
  try {
    const page = await projectsApi.getProjectTasksPage(project.value.id, {
      limit: projectDetailsTasksMeta.value.limit || 600,
      cursorCreatedAt: projectTasksNextCursor.value?.created_at,
      cursorId: projectTasksNextCursor.value?.id,
      offset: projectTasksNextOffset.value,
    })
    if (page.tasks.length > 0) {
      const merged = [...allTasks.value]
      const seen = new Set(merged.map((t) => t.id))
      for (const t of page.tasks) {
        if (!seen.has(t.id)) {
          merged.push(t)
          seen.add(t.id)
        }
      }
      allTasks.value = merged
      const last = page.tasks[page.tasks.length - 1]
      projectTasksNextCursor.value = { created_at: last.created_at, id: last.id }
    }
    projectTasksNextOffset.value = page.next_offset ?? (projectTasksNextOffset.value + page.returned)
    projectDetailsTasksMeta.value = {
      limit: page.limit,
      returned: allTasks.value.length,
      has_more: page.has_more,
    }
  } catch (e: any) {
    error.value = e?.message || 'Failed to load more tasks'
  } finally {
    isLoadingMoreProjectTasks.value = false
  }
}

function setupBacklogInfiniteScroll() {
  if (typeof window === 'undefined') return
  if (backlogAutoLoadObserver) {
    backlogAutoLoadObserver.disconnect()
    backlogAutoLoadObserver = null
  }
  const el = backlogAutoLoadSentinelRef.value
  if (!el) return
  backlogAutoLoadObserver = new IntersectionObserver((entries) => {
    const entry = entries[0]
    if (!entry?.isIntersecting) return
    if (activeTab.value !== 'backlog') return
    if (!projectDetailsTasksMeta.value?.has_more) return
    void loadMoreProjectTasks()
  }, { root: null, rootMargin: '300px 0px 300px 0px', threshold: 0 })
  backlogAutoLoadObserver.observe(el)
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
  const task = idx !== -1 ? allTasks.value[idx] : null
  if (task) (task as any)[field] = value || null
  try {
    // sprint_id: JSON omits `undefined`; backend needs explicit "" to clear sprint (see UpdateTask hasSprint / empty string branch).
    const payload: Record<string, unknown> =
      field === 'epic_id'
        ? { epic_id: value ?? '' }
        : field === 'sprint_id'
          ? { sprint_id: value && String(value).trim() !== '' ? String(value) : '' }
          : { [field]: value || undefined }
    // When moving task to a sprint, clamp start/end to sprint range so the bar shows within the sprint
    if (field === 'sprint_id' && value && task) {
      let sprint = sprints.value.find((s) => s.id === value) ?? null
      if (!sprint?.start_date && sprintTimelineData.value?.sprints?.length) {
        const fromTimeline = sprintTimelineData.value.sprints.find((s) => s.id === value)
        if (fromTimeline) sprint = fromTimeline as Sprint
      }
      const dates = sprint ? taskDatesInSprintRange(task, sprint) : null
      if (dates) {
        payload.start_date = dates.start_date
        payload.end_date = dates.end_date
        if (idx !== -1) {
          allTasks.value[idx] = { ...allTasks.value[idx], start_date: dates.start_date, end_date: dates.end_date }
        }
      }
    }
    await tasksApi.updateTask(taskId, payload as any)
    // Keep Epic and Sprint timeline in sync: update task dates in both in-memory datasets
    if (field === 'sprint_id' && payload.start_date && payload.end_date) {
      updateTaskInTimelineData(taskId, payload.start_date as string, payload.end_date as string)
    }
    // Refresh timeline data so task appears under the correct sprint/epic row
    if (field === 'sprint_id') await loadSprintTimeline()
    else if (field === 'epic_id') await loadEpicTimeline()
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

// Edit Task Title modal
const showEditTaskTitleModal = ref(false)
const editingTaskForTitle = ref<Task | null>(null)
const editTaskTitleValue = ref('')
const isSavingTaskTitle = ref(false)

function openEditTaskTitle(task: Task) {
  editingTaskForTitle.value = task
  editTaskTitleValue.value = task.title || ''
  showEditTaskTitleModal.value = true
}

function closeEditTaskTitleModal() {
  showEditTaskTitleModal.value = false
  editingTaskForTitle.value = null
  editTaskTitleValue.value = ''
}

async function saveEditTaskTitle() {
  const task = editingTaskForTitle.value
  const trimmed = editTaskTitleValue.value.trim()
  if (!task || !trimmed) return
  isSavingTaskTitle.value = true
  try {
    await updateTaskField(task.id, 'title', trimmed)
    closeEditTaskTitleModal()
  } finally {
    isSavingTaskTitle.value = false
  }
}

const isDuplicatingTask = ref(false)
/** After duplicate: keep new task visually right below original until page refresh. */
const duplicatePlacement = ref<{ newId: string; afterId: string } | null>(null)

async function duplicateTask(task: Task) {
  if (!project.value) return
  isDuplicatingTask.value = true
  duplicatePlacement.value = null
  try {
    let source = task
    try {
      source = await tasksApi.getTask(task.id)
    } catch {
      // list payload may omit description; duplicate still works with empty body
    }
    const payload: any = {
      title: (source.title || '').trim() ? `${(source.title || '').trim()} (copy)` : 'Task (copy)',
      description: source.description || '',
      priority: source.priority || 'MEDIUM',
      story_points: source.story_points ?? 0,
      project_id: project.value.id,
    }
    if (source.epic_id) payload.epic_id = source.epic_id
    if (source.sprint_id != null) payload.sprint_id = source.sprint_id
    let newTask = await tasksApi.createTask(payload)
    const nextOrder = (source.sort_order ?? 0) + 1
    try {
      const updated = await tasksApi.updateTask(newTask.id, { sort_order: nextOrder })
      newTask = updated
    } catch {
      // ignore if backend doesn't support sort_order on update
    }
    const idx = allTasks.value.findIndex((t) => t.id === task.id)
    if (idx !== -1) {
      allTasks.value.splice(idx + 1, 0, newTask)
    } else {
      allTasks.value.unshift(newTask)
    }
    // Keep duplicated task below original in backlog until refresh
    duplicatePlacement.value = { newId: newTask.id, afterId: task.id }
  } catch (e: any) {
    console.error('Duplicate task failed:', e)
  } finally {
    isDuplicatingTask.value = false
  }
}

// Edit project modal
const showEditProjectModal = ref(false)
const editProjectForm = ref({ name: '', description: '', status: 'ACTIVE' as string, update_code: false })
const editProjectError = ref('')
const isSavingProject = ref(false)

function openEditProjectModal() {
  if (!project.value) return
  editProjectForm.value = {
    name: project.value.name,
    description: project.value.description || '',
    status: project.value.status || 'ACTIVE',
    update_code: false,
  }
  editProjectError.value = ''
  showEditProjectModal.value = true
}

function closeEditProjectModal() {
  showEditProjectModal.value = false
}

async function saveEditProject() {
  if (!project.value) return
  const idOrCode = (route.params.id as string) || project.value.id || project.value.code
  if (!idOrCode) return
  isSavingProject.value = true
  editProjectError.value = ''
  try {
    const updated = await projectsApi.updateProject(idOrCode, {
      name: editProjectForm.value.name.trim(),
      description: editProjectForm.value.description,
      status: editProjectForm.value.status,
      update_code: editProjectForm.value.update_code,
    })
    project.value = updated
    if (editProjectForm.value.update_code) {
      await loadAll()
      if (updated.code && updated.code !== route.params.id) {
        await router.replace(`/projects/${updated.code}`)
      }
    }
    closeEditProjectModal()
  } catch (e: any) {
    editProjectError.value = e?.data?.message ?? e?.data?.error ?? e?.message ?? 'บันทึกไม่สำเร็จ'
  } finally {
    isSavingProject.value = false
  }
}

const {
  sprintsOrdered,
  getSprintStats,
  showSprintModal,
  editingSprint,
  sprintForm,
  sprintDurationWeeks,
  sprintDurationOptions,
  isCreatingSprint,
  sprintError,
  syncSprintEndFromStartAndDuration,
  onSprintDurationWeeksChange,
  applySuggestedSprintName,
  resetSprintModalToDefaults,
  openSprintModal,
  openEditSprintModal,
  closeSprintModal,
  submitSprint,
  sprintDragId,
  onSprintDragStart,
  onSprintDragOver,
  onSprintDrop,
  showAddTasksToSprintModal,
  sprintForAddTasks,
  selectedTaskIdsForSprint,
  addTasksToSprintError,
  isAddingTasksToSprint,
  tasksNotInSprint,
  openAddTasksToSprintModal,
  closeAddTasksToSprintModal,
  confirmAddTasksToSprint,
  showDeleteSprintModal,
  sprintToDelete,
  deleteSprintError,
  isDeletingSprint,
  openDeleteSprintModal,
  closeDeleteSprintModal,
  confirmDeleteSprint,
  showCompleteSprintModal,
  sprintToComplete,
  completeSprintError,
  isCompletingSprint,
  openCompleteSprintModal,
  closeCompleteSprintModal,
  confirmCompleteSprint,
  showReopenSprintModal,
  sprintToReopen,
  reopenSprintError,
  isReopeningSprint,
  openReopenSprintModal,
  closeReopenSprintModal,
  confirmReopenSprint,
  handleStartSprint,
} = useProjectSprints({
  project,
  sprints,
  allTasks,
  activeSprint,
  projectsApi,
  tasksApi,
  loadAll,
  loadSprintTimeline,
  updateTaskInTimelineData,
  taskDatesInSprintRange,
  showError,
})

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
    due_date: isoToDatetimeLocal(m.due_date),
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
  const ok = await confirm({
    title: 'ยืนยันการลบ milestone',
    message: `ยืนยันการลบ milestone "${name}"? กด Confirm เพื่อลบ / Cancel เพื่อยกเลิก`,
    confirmLabel: 'ลบ',
    cancelLabel: 'ยกเลิก',
    variant: 'danger'
  })
  if (!ok) return
  try {
    await projectsApi.deleteMilestone(editingMilestone.value.id)
    milestones.value = milestones.value.filter((m) => m.id !== editingMilestone.value!.id)
    closeMilestoneModal()
  } catch (e: any) {
    milestoneError.value = e.message
  }
}

async function onMilestoneDragEnd(payload: { milestone: Milestone; newDueDate: string }) {
  milestoneDragPosition.value = null
  const { milestone, newDueDate } = payload
  const dueDateISO = newDueDate + 'T12:00:00.000Z'
  try {
    const updated = await projectsApi.updateMilestone(milestone.id, { due_date: dueDateISO })
    const idx = milestones.value.findIndex((m) => m.id === milestone.id)
    if (idx !== -1) milestones.value[idx] = updated
  } catch (e) {
    console.error('Failed to update milestone date:', e)
  }
}

// Create Task Modal
const showCreateTaskModal = ref(false)
const createTaskForm = ref({
  title: '', description: '', task_type: 'TASK', priority: 'MEDIUM', story_points: 0,
  sprint_id: '', due_date: '', start_date: '', end_date: '', parent_id: '', epic_id: '',
  estimated_hours: 0
})
const isCreatingTask = ref(false)
const createTaskError = ref('')

function openCreateTaskModal(parentId?: string, epicId?: string) {
  createTaskForm.value = { title: '', description: '', task_type: 'TASK', priority: 'MEDIUM', story_points: 0, sprint_id: '', due_date: '', start_date: '', end_date: '', parent_id: parentId || '', epic_id: epicId || '', estimated_hours: 0 }
  createTaskError.value = ''
  showCreateTaskModal.value = true
}

function closeCreateTaskModal() { showCreateTaskModal.value = false }

const {
  showBacklogImportModal,
  backlogImportStep,
  isBacklogImporting,
  isBacklogLoadingPreview,
  backlogImportError,
  backlogImportResult,
  backlogImportPreview,
  backlogImportSelectedIndices,
  backlogImportTriagedSlides,
  backlogImportAssignees,
  backlogImportForm,
  backlogParentTaskOptions,
  onBacklogImportEpicChange,
  openBacklogImportModal,
  closeBacklogImportModal,
  loadBacklogImportPreview,
  backlogImportSelectAll,
  backlogImportDeselectAll,
  backlogImportSelectOnlyNew,
  submitBacklogImport,
  showPPTXImportModal,
  pptxImportStep,
  isPPTXImporting,
  isPPTXLoadingPreview,
  pptxImportError,
  pptxImportResult,
  pptxImportPreview,
  pptxImportFile,
  pptxDragOver,
  pptxImportSelectedIndices,
  pptxImportTriagedSlides,
  pptxImportForm,
  pptxParentTaskOptions,
  onPPTXImportEpicChange,
  onPPTXFileChange,
  onPPTXFileDrop,
  openPPTXImportModal,
  closePPTXImportModal,
  pptxImportSelectAll,
  pptxImportDeselectAll,
  pptxImportSelectOnlyVisible,
  loadPPTXImportPreview,
  submitPPTXImport,
  showSheetsImportModal,
  sheetsImportStep,
  isSheetsImporting,
  isSheetsLoadingPreview,
  sheetsImportError,
  sheetsImportResult,
  sheetsImportPreview,
  sheetsImportSelectedRowIndices,
  sheetsImportTriagedRows,
  sheetsImportForm,
  sheetsParentTaskOptions,
  onSheetsImportEpicChange,
  openSheetsImportModal,
  closeSheetsImportModal,
  sheetsImportSelectAll,
  sheetsImportDeselectAll,
  loadSheetsImportPreview,
  submitSheetsImport,
  showIODImportModal,
  iodImportStep,
  isIODImporting,
  isIODLoadingPreview,
  iodImportError,
  iodImportResult,
  iodImportPreview,
  iodImportSelectedRowIndices,
  iodImportTriagedRows,
  iodImportForm,
  iodParentTaskOptions,
  onIODImportEpicChange,
  openIODImportModal,
  closeIODImportModal,
  loadIODImportPreview,
  submitIODImport,
} = useProjectImports({
  allTasks,
  project,
  currentUser,
  tasksApi,
  teamsStore,
  loadAll,
})

async function submitCreateTask() {
  if (!project.value) return
  isCreatingTask.value = true
  createTaskError.value = ''
  try {
    const estMins = effortHoursToMinutes(Number(createTaskForm.value.estimated_hours) || 0)
    const payload: any = {
      title: createTaskForm.value.title,
      description: createTaskForm.value.description,
      task_type: createTaskForm.value.task_type || 'TASK',
      priority: createTaskForm.value.priority,
      story_points: createTaskForm.value.story_points,
      project_id: project.value.id,
      estimated_minutes: estMins,
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
    nextTick(() => buildMatrixGanttRows())
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
    nextTick(() => buildMatrixGanttRows())
  } catch (e) {
    console.error('Failed to load sprint timeline:', e)
  } finally {
    matrixTimelineLoading.value = false
  }
}

async function refreshTimeline() {
  if (!project.value || timelineRefreshing.value) return
  timelineRefreshing.value = true
  try {
    await loadAll()
    if (timelineMode.value === 'epic') await loadEpicTimeline()
    else await loadSprintTimeline()
    nextTick(() =>
      setTimeout(ganttView.value === 'day' ? scrollTimelineDayFocusRecent : scrollTimelineToToday, 200)
    )
  } finally {
    timelineRefreshing.value = false
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
    start_date: isoToDatetimeLocal(epic.start_date),
    end_date: isoToDatetimeLocal(epic.end_date),
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
      expandedEpicGroups.value[created.id] = false
    }
    closeEpicModal()
  } catch (e: any) {
    epicError.value = e.message || 'Failed to save epic'
  } finally {
    isSavingEpic.value = false
  }
}

async function deleteEpic(epic: Epic) {
  const ok = await confirm({
    title: 'Delete epic',
    message: `Delete epic "${epic.title}"? Tasks in this epic will be unlinked.`,
    confirmLabel: 'Delete',
    cancelLabel: 'Cancel',
    variant: 'danger'
  })
  if (!ok) return
  isDeletingEpic.value = true
  try {
    await projectsApi.deleteEpic(epic.id)
    epics.value = epics.value.filter((e) => e.id !== epic.id)
    // Unlink tasks locally
    allTasks.value = allTasks.value.map((t) => t.epic_id === epic.id ? { ...t, epic_id: null } : t)
  } catch (e: any) {
    showError(e.message || 'Failed to delete epic', 'Delete epic failed')
  } finally {
    isDeletingEpic.value = false
  }
}

function toggleEpicGroup(id: string) {
  expandedEpicGroups.value[id] = !expandedEpicGroups.value[id]
}

function expandAllBacklog() {
  epics.value.forEach((ep) => { expandedEpicGroups.value[ep.id] = true })
  expandedEpicGroups.value['__unassigned__'] = true
  allTasks.value.forEach((t) => {
    if (!t.parent_id && getSubTasks(t.id).length > 0) expandedEpics.value[t.id] = true
  })
}

function collapseAllBacklog() {
  epics.value.forEach((ep) => { expandedEpicGroups.value[ep.id] = false })
  expandedEpicGroups.value['__unassigned__'] = false
  expandedEpics.value = {}
}


/** After duplicate: place new task right below original until refresh. */
function applyDuplicatePlacement<T extends { id: string }>(tasks: T[]): T[] {
  const placement = duplicatePlacement.value
  if (!placement || placement.newId === placement.afterId) return tasks
  const afterIdx = tasks.findIndex((t) => t.id === placement.afterId)
  const newIdx = tasks.findIndex((t) => t.id === placement.newId)
  if (afterIdx === -1 || newIdx === -1) return tasks
  const list = [...tasks]
  const [item] = list.splice(newIdx, 1)
  const insertAt = afterIdx < newIdx ? afterIdx + 1 : afterIdx
  list.splice(insertAt, 0, item)
  return list
}

function getTasksForEpic(epicId: string) {
  const sprintOrderIds = sprintsOrdered.value.map((s) => s.id)
  const sorted = sortBacklogTasks(
    allTasks.value.filter((t) => t.epic_id === epicId && !t.parent_id),
    sprintOrderIds,
  )
  return applyDuplicatePlacement(sorted)
}

function getUnassignedTasks() {
  const sprintOrderIds = sprintsOrdered.value.map((s) => s.id)
  const sorted = sortBacklogTasks(
    allTasks.value.filter((t) => !t.epic_id && !t.parent_id),
    sprintOrderIds,
  )
  return applyDuplicatePlacement(sorted)
}

/** Row IDs in one backlog section (epic or unassigned), including all sub-tasks in the tree. */
function pushBacklogSectionTaskIds(ids: string[], task: Task) {
  ids.push(task.id)
  getSubTasks(task.id).forEach((child) => pushBacklogSectionTaskIds(ids, child))
}

function backlogSectionRowIds(sectionKey: string): string[] {
  const ids: string[] = []
  const roots =
    sectionKey === '__unassigned__' ? getUnassignedTasks() : getTasksForEpic(sectionKey)
  roots.forEach((t) => pushBacklogSectionTaskIds(ids, t))
  return ids
}

function backlogSectionSelectAllChecked(sectionKey: string): boolean {
  const ids = backlogSectionRowIds(sectionKey)
  if (ids.length === 0) return false
  const sel = backlogSelectedTaskIds.value
  return ids.every((id) => sel.has(id))
}

function backlogSectionSelectAllIndeterminate(sectionKey: string): boolean {
  const ids = backlogSectionRowIds(sectionKey)
  if (ids.length === 0) return false
  const sel = backlogSelectedTaskIds.value
  const n = ids.filter((id) => sel.has(id)).length
  return n > 0 && n < ids.length
}

function setBacklogSectionSelectAll(sectionKey: string, checked: boolean) {
  const rowIds = backlogSectionRowIds(sectionKey)
  const next = new Set(backlogSelectedTaskIds.value)
  if (checked) {
    rowIds.forEach((id) => next.add(id))
  } else {
    rowIds.forEach((id) => next.delete(id))
  }
  backlogSelectedTaskIds.value = next
}

// --- Backlog drag-and-drop ---
async function reorderEpics(newOrder: Epic[]) {
  for (let i = 0; i < newOrder.length; i++) {
    if (newOrder[i].sort_order === i) continue
    try {
      await projectsApi.updateEpic(newOrder[i].id, { sort_order: i })
      const idx = epics.value.findIndex((e) => e.id === newOrder[i].id)
      if (idx >= 0) (epics.value[idx] as { sort_order: number }).sort_order = i
    } catch {
      await loadAll()
      break
    }
  }
}

async function reorderTasksInBacklog(orderedTaskIds: string[]) {
  try {
    for (let i = 0; i < orderedTaskIds.length; i++) {
      await tasksApi.updateTask(orderedTaskIds[i], { sort_order: i })
      const t = allTasks.value.find((x) => x.id === orderedTaskIds[i])
      if (t) (t as Task & { sort_order: number }).sort_order = i
    }
  } catch {
    await loadAll()
  }
}

const backlogDrag = ref<{ type: 'epic' | 'task'; id: string; epicId?: string | null } | null>(null)

function onEpicDragStart(e: DragEvent, epicId: string) {
  backlogDrag.value = { type: 'epic', id: epicId }
  e.dataTransfer?.setData?.('application/json', JSON.stringify({ type: 'epic', id: epicId }))
  e.dataTransfer!.effectAllowed = 'move'
}

function onEpicDragOver(e: DragEvent) {
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
}

function onEpicDrop(e: DragEvent, dropIndex: number) {
  e.preventDefault()
  let dragId: string | null = null
  try {
    const raw = e.dataTransfer?.getData('application/json')
    if (raw) {
      const p = JSON.parse(raw) as { type: string; id: string }
      if (p.type === 'epic') dragId = p.id
    }
  } catch {}
  backlogDrag.value = null
  if (!dragId) return
  const fromIndex = epics.value.findIndex((x) => x.id === dragId)
  if (fromIndex < 0 || fromIndex === dropIndex) return
  const next = [...epics.value]
  const [removed] = next.splice(fromIndex, 1)
  next.splice(dropIndex, 0, removed)
  epics.value = next
  reorderEpics(next)
}

function onTaskDragStartSetData(e: DragEvent, taskId: string, epicId: string | null) {
  backlogDrag.value = { type: 'task', id: taskId, epicId }
  e.dataTransfer?.setData?.('application/json', JSON.stringify({ type: 'task', id: taskId, epicId }))
  e.dataTransfer!.effectAllowed = 'move'
  const row = (e.target as HTMLElement)?.closest?.('.backlog-row')
  if (row && e.dataTransfer?.setDragImage) {
    const rect = row.getBoundingClientRect()
    e.dataTransfer.setDragImage(row, Math.min(20, rect.width / 2), rect.height / 2)
  }
}

function onTaskDragOver(e: DragEvent) {
  e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
}

function onTaskDrop(e: DragEvent, epicId: string | null, dropIndex: number) {
  e.preventDefault()
  let taskId: string | null = null
  let dragEpicId: string | null = null
  try {
    const raw = e.dataTransfer?.getData('application/json')
    if (raw) {
      const p = JSON.parse(raw) as { type: string; id: string; epicId?: string | null }
      if (p.type === 'task') {
        taskId = p.id
        dragEpicId = p.epicId ?? null
      }
    }
  } catch {}
  backlogDrag.value = null
  if (!taskId || dragEpicId !== epicId) return
  const list = epicId ? getTasksForEpic(epicId) : getUnassignedTasks()
  const fromIndex = list.findIndex((t) => t.id === taskId)
  if (fromIndex < 0 || fromIndex === dropIndex) return
  const next = [...list.map((t) => t.id)]
  const [removed] = next.splice(fromIndex, 1)
  next.splice(dropIndex, 0, removed)
  reorderTasksInBacklog(next)
}

watch(
  () => [activeTab.value, projectDetailsTasksMeta.value?.has_more, backlogAutoLoadSentinelRef.value] as const,
  ([tab, hasMore]) => {
    if (tab === 'backlog' && hasMore) {
      nextTick(() => setupBacklogInfiniteScroll())
      return
    }
    if (backlogAutoLoadObserver) {
      backlogAutoLoadObserver.disconnect()
      backlogAutoLoadObserver = null
    }
  }
)

onMounted(loadAll)
onBeforeUnmount(() => {
  if (backlogAutoLoadObserver) {
    backlogAutoLoadObserver.disconnect()
    backlogAutoLoadObserver = null
  }
})
</script>

<style scoped>
/* Backlog table: Task column takes remaining space, + Sub column minimal */
.backlog-grid {
  display: grid;
  grid-template-columns: auto auto minmax(0, 1fr) auto auto auto auto auto auto;
}

/* ตาราง backlog: ใช้ grid คอลัมน์ชัดเจน (ไม่ใช้ subgrid) ให้ทุก browser แสดงตรงกัน */
.backlog-table-grid {
  display: grid;
  grid-template-columns: 2rem 2.5rem 3.25rem minmax(18rem, 4.5fr) 3rem minmax(6rem, 0.8fr) minmax(8rem, 1.2fr) minmax(6rem, 0.9fr) 5.5rem 3.5rem;
  column-gap: 0.75rem;
  row-gap: 0;
  padding: 0;
  align-items: center;
}
@media (min-width: 640px) {
  .backlog-table-grid {
    column-gap: 1rem;
    grid-template-columns: 2rem 2.5rem 3.5rem minmax(20rem, 5fr) 3.5rem minmax(6.5rem, 0.8fr) minmax(10rem, 1.2fr) minmax(7rem, 0.9fr) 6rem 3.5rem;
  }
}
/* แถวหัวตารางและแถวข้อมูลใช้ grid เดียวกับ parent (แต่ละแถวเป็น grid row ใน .backlog-table-grid) */
.backlog-subgrid {
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: 2rem 2.5rem 3.25rem minmax(18rem, 4.5fr) 3rem minmax(6rem, 0.8fr) minmax(8rem, 1.2fr) minmax(6rem, 0.9fr) 5.5rem 3.5rem;
  align-items: center;
  column-gap: 0.75rem;
  row-gap: 0;
}
@media (min-width: 640px) {
  .backlog-subgrid {
    column-gap: 1rem;
    grid-template-columns: 2rem 2.5rem 3.5rem minmax(20rem, 5fr) 3.5rem minmax(6.5rem, 0.8fr) minmax(10rem, 1.2fr) minmax(7rem, 0.9fr) 6rem 3.5rem;
  }
}
/* เซลล์ทุกคอลัมน์: padding สม่ำเสมอ ให้ข้อมูลชิดกับ header และไม่ชิดขอบ */
.backlog-subgrid > div {
  min-width: 0;
  padding-left: 0.5rem;
  padding-right: 0.5rem;
}
.backlog-subgrid > div:first-child {
  padding-left: 0.75rem;
  padding-right: 0.25rem;
}
.backlog-subgrid > div:last-child {
  padding-left: 0.25rem;
  padding-right: 0.5rem;
  margin-left: 0;
}
@media (min-width: 640px) {
  .backlog-subgrid > div:first-child {
    padding-left: 1rem;
    padding-right: 0.375rem;
  }
  .backlog-subgrid > div:last-child {
    padding-right: 0.5rem;
    margin-left: 0;
  }
}

/* หัวตาราง: ความสูงและระยะชิดเนื้อหาให้ตรงกับแถวข้อมูล */
.backlog-table-header {
  font-size: 0.6875rem;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  font-weight: 600;
  background: rgba(30, 41, 59, 0.6);
  color: rgb(148 163 184);
  padding: 0.875rem 0;
  border-bottom: 1px solid rgba(75, 85, 99, 0.9);
}
.backlog-table-header > div {
  padding-top: 0;
  padding-bottom: 0;
}
/* หัวคอลัมน์ ID ชิดขวา */
.backlog-table-header > div:nth-child(3) {
  justify-content: flex-end;
}
/* ช่อง ID ในแถวข้อมูล: จัดตัวเลขชิดขวา */
.backlog-subgrid > div:nth-child(3) {
  justify-content: flex-end;
}
/* ช่องชื่อ Task: ชิดซ้าย (ยกเว้นหัวตาราง) */
.backlog-subgrid > div:nth-child(4) {
  justify-content: flex-start;
}
/* ช่อง Epic ในแถวระดับ A: จัดเนื้อหาชิดขวา */
.backlog-row.backlog-subgrid > div:nth-child(7) {
  justify-content: flex-end;
}
/* แถว sub-row (B/C): คอลัมน์ Epic และ Sprint จัดข้อความ Inherits ตรงกลาง */
.backlog-sub-row > div:nth-child(7),
.backlog-sub-row > div:nth-child(8) {
  justify-content: center;
}
/* คอลัมน์ Status (คอลัมน์ที่ 9): จัด badge ตรงกลาง */
.backlog-subgrid > div:nth-child(9) {
  justify-content: center;
}
/* หัวคอลัมน์ Task และ Epic อยู่ตรงกลาง (ต้องอยู่หลัง rule ทั่วไปเพื่อ override) */
.backlog-table-header.backlog-subgrid > div:nth-child(4),
.backlog-table-header.backlog-subgrid > div:nth-child(7) {
  justify-content: center;
  text-align: center;
}

/* แถวข้อมูล: padding แนวตั้งพอดี อ่านง่าย */
.backlog-row {
  padding: 0.625rem 0;
  border-bottom: 1px solid rgba(55, 65, 81, 0.55);
  min-width: 0;
  background: rgba(31, 41, 55, 0.55);
}
.backlog-row:hover {
  background: rgba(55, 65, 81, 0.6);
}
/* Sub-task rows: ระยะห่างสมดุลกับแถวหลัก */
.backlog-sub-row {
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}

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
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
/* Add Task modal — larger, easier to read (does not affect rest of page) */
.create-task-modal .label {
  @apply block text-sm sm:text-base text-gray-300 mb-2 font-medium;
}
.create-task-modal .input-field {
  @apply bg-gray-700 border border-gray-500 rounded-xl px-4 py-3.5 text-base text-gray-100 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
.btn-primary {
  @apply bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-xl transition-colors;
}
.btn-primary-sm {
  @apply px-3 py-1.5 text-xs bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-medium rounded-lg transition-colors;
}
.btn-import-sm {
  @apply px-3 py-1.5 text-xs bg-purple-900/50 hover:bg-purple-800/60 border border-purple-700/50 text-purple-300 font-medium rounded-lg transition-colors flex items-center gap-1.5;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
/* Timeline tab layout */
.timeline-tab {
  --gantt-bar: 147 51 234; /* purple-600 */
  --gantt-bar-hover: 168 85 247; /* purple-500 */
  --gantt-today: 96 165 250; /* blue-400 */
}

.gantt-chart-vue {
  min-height: 420px;
  font-family: ui-sans-serif, system-ui, sans-serif;
}

.timeline-scroll-wrapper {
  overscroll-behavior-x: contain;
}

/* Fullscreen: พื้นที่ scroll แนวตั้ง+แนวนอน (ลูกของ wrapper ที่มี flex-1 min-h-0) */
.timeline-fullscreen-scroll {
  overscroll-behavior: auto;
  scroll-behavior: smooth;
  -webkit-overflow-scrolling: touch;
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

/* แยก Epic / Sprint / Task ให้เห็นชัด */
.gantt-enterprise :deep(.gantt-row-epic) > div:first-child {
  background: linear-gradient(90deg, rgb(67 56 202 / 0.35) 0%, rgb(30 41 59 / 0.6) 100%);
  border-left: 3px solid rgb(99 102 241);
  font-weight: 600;
  color: rgb(224 231 255);
}

.gantt-enterprise :deep(.gantt-row-sprint) > div:first-child {
  background: linear-gradient(90deg, rgb(6 95 70 / 0.3) 0%, rgb(30 41 59 / 0.6) 100%);
  border-left: 3px solid rgb(16 185 129);
  font-weight: 600;
  color: rgb(167 243 208);
}

.gantt-enterprise :deep(.gantt-row-task) > div:first-child {
  background: rgb(30 41 59 / 0.4);
  border-left: 3px solid transparent;
  font-weight: 400;
  color: rgb(203 213 225);
}

/* สีแท่งตามประเภทแถว (ถ้า library ไม่ใส่ class ที่ bar) */
.gantt-enterprise :deep(.gantt-row-epic .g-gantt-bar) {
  background: linear-gradient(135deg, rgb(99 102 241) 0%, rgb(67 56 202) 100%) !important;
  border: 1px solid rgb(129 140 248 / 0.5);
}

.gantt-enterprise :deep(.gantt-row-sprint .g-gantt-bar) {
  background: linear-gradient(135deg, rgb(16 185 129) 0%, rgb(5 150 105) 100%) !important;
  border: 1px solid rgb(52 211 153 / 0.5);
}

.gantt-enterprise :deep(.gantt-row-task .g-gantt-bar) {
  background: linear-gradient(135deg, rgb(139 92 246) 0%, rgb(124 58 237) 100%) !important;
  border: 1px solid rgb(167 139 250 / 0.5);
}

.gantt-enterprise :deep(.g-grid-line) {
  border-left: 1px solid rgb(51 65 85 / 0.7);
}

.gantt-enterprise :deep(.g-gantt-bar) {
  border-radius: 6px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.2);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

/* Epic: fallback สี indigo เมื่อไม่มีคลาสเฉพาะ (สีจริงมาจาก epicBarStyles ตาม epic.color) */
.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-epic:not([class*="gantt-bar-epic-e-"])) {
  background: linear-gradient(135deg, rgb(99 102 241) 0%, rgb(67 56 202) 100%) !important;
  border: 1px solid rgb(129 140 248 / 0.5);
}

.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-epic:hover) {
  transform: translateY(-1px);
}

/* Sprint: แท่งสี emerald */
.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-sprint) {
  background: linear-gradient(135deg, rgb(16 185 129) 0%, rgb(5 150 105) 100%) !important;
  border: 1px solid rgb(52 211 153 / 0.5);
}

.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-sprint:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgb(16 185 129 / 0.4);
}

/* Task: แท่งสี violet อ่อน แยกจาก Epic */
.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-task) {
  background: linear-gradient(135deg, rgb(139 92 246) 0%, rgb(124 58 237) 100%) !important;
  border: 1px solid rgb(167 139 250 / 0.5);
}

.gantt-enterprise :deep(.g-gantt-bar.gantt-bar-task:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgb(139 92 246 / 0.35);
}

.gantt-enterprise :deep(.g-gantt-bar:hover) {
  transform: translateY(-1px);
}

.gantt-enterprise :deep(.g-gantt-bar-label) {
  color: rgb(255 255 255);
  font-size: 0.8rem;
  font-weight: 500;
  text-shadow: 0 1px 1px rgb(0 0 0 / 0.3);
  white-space: normal;
  word-break: break-word;
  overflow: visible;
}

/* ลากขอบซ้าย/ขวาเพื่อขยาย-ย่อแท่ง (resize); ลากกลางแท่งเพื่อเลื่อน (drag) */
.gantt-enterprise :deep(.g-gantt-bar-handle-left),
.gantt-enterprise :deep(.g-gantt-bar-handle-right) {
  background: rgb(255 255 255 / 0.35) !important;
  width: 12px;
  min-width: 12px;
  border-radius: 4px 0 0 4px;
  cursor: ew-resize;
  transition: background 0.15s ease;
}

.gantt-enterprise :deep(.g-gantt-bar:hover .g-gantt-bar-handle-left),
.gantt-enterprise :deep(.g-gantt-bar:hover .g-gantt-bar-handle-right) {
  background: rgb(255 255 255 / 0.6) !important;
}

.gantt-enterprise :deep(.g-gantt-bar-handle-right) {
  border-radius: 0 4px 4px 0;
}

/* Current time ("Now") – ให้เมาส์ทะลุไปถึงแท่ง ไม่บังการลาก/คลิก */
.gantt-enterprise :deep(.g-grid-current-time),
.gantt-enterprise :deep(.g-grid-current-time-marker),
.gantt-enterprise :deep(.g-grid-current-time-text) {
  pointer-events: none;
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

/* Hide default tooltip color dot from vue-ganttastic */
.gantt-enterprise :deep(.g-gantt-tooltip-color-dot) {
  display: none !important;
}
</style>

<style>
/* Global override because vue-ganttastic tooltip is teleported to <body> */
.g-gantt-tooltip {
  display: none !important;
}
</style>
