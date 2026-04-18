<template>
  <div class="min-h-screen task-enterprise-bg text-gray-800 dark:text-gray-100">

    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center min-h-screen">
      <div class="w-12 h-12 rounded-2xl bg-purple-600/20 border border-purple-500/30 flex items-center justify-center mb-4">
        <svg class="w-6 h-6 text-purple-400 animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
        </svg>
      </div>
      <p class="text-sm text-gray-500">กำลังโหลด task...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex flex-col items-center justify-center min-h-screen px-6">
      <div class="max-w-md w-full bg-red-900/20 border border-red-500/40 rounded-2xl p-8 text-center">
        <div class="w-14 h-14 rounded-2xl bg-red-900/40 flex items-center justify-center mx-auto mb-4">
          <svg class="w-7 h-7 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
        </div>
        <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-2">Failed to load task</h2>
        <p class="text-sm text-red-300 mb-6">{{ error }}</p>
        <button @click="goToDashboard" class="btn-primary px-6 py-2.5 text-sm">← Back</button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="task" class="relative mx-auto w-full max-w-[1600px] px-4 sm:px-6 lg:px-8 py-6 lg:py-8">

      <!-- ══ TOP BAR ══ -->
      <div class="flex flex-wrap items-center justify-between gap-3 mb-6">
        <!-- Breadcrumb + nav -->
        <div class="flex items-center gap-2 min-w-0">
          <NuxtLink :to="backTarget" class="flex items-center gap-1.5 text-sm text-gray-500 hover:text-gray-700 dark:text-gray-200 transition-colors group">
            <svg class="w-4 h-4 group-hover:-translate-x-0.5 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
            Back
          </NuxtLink>
          <span class="text-gray-700">/</span>
          <code class="text-xs font-mono px-2 py-0.5 bg-gray-800 border border-gray-700 rounded-md text-purple-400">{{ taskCodeDisplay(task) }}</code>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2 shrink-0">
          <!-- Prev / Next (sprint/backlog context) -->
          <template v-if="showPrevNext">
            <button
              type="button"
              @click="goToPrevTask"
              :disabled="!activePrevTaskLink"
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-lg border transition-colors"
              :class="activePrevTaskLink ? 'border-gray-600 text-gray-600 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-500 hover:text-white' : 'border-gray-700 text-gray-600 cursor-default'"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
              Prev
            </button>
            <button
              type="button"
              @click="goToNextTask"
              :disabled="!activeNextTaskLink"
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-lg border transition-colors"
              :class="activeNextTaskLink ? 'border-gray-600 text-gray-600 dark:text-gray-300 hover:border-gray-400 dark:hover:border-gray-500 hover:text-white' : 'border-gray-700 text-gray-600 cursor-default'"
            >
              Next
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
            </button>
            <span v-if="sprintNavMessage" class="text-xs text-amber-400">{{ sprintNavMessage }}</span>
          </template>

          <template v-if="canEditOrDelete">
            <button
              v-if="showCEOQuickFinishAction"
              type="button"
              :disabled="ceoQuickFinishSubmitting"
              @click="markTaskFinishedByCEO"
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-semibold rounded-lg border border-emerald-300 dark:border-emerald-500/45 bg-emerald-100 dark:bg-emerald-500/15 text-emerald-300 hover:bg-emerald-100 dark:bg-emerald-500/25 hover:border-emerald-400/60 transition-colors disabled:opacity-50"
            >
              <span>✅</span>
              {{ ceoQuickFinishSubmitting ? 'Finishing…' : 'Finished' }}
            </button>
            <button
              @click="openEditModal"
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 rounded-lg transition-colors border border-gray-600 hover:border-gray-400 dark:hover:border-gray-500"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
              Edit
            </button>
            <button
              @click="openDeleteConfirmation"
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-400 bg-red-900/20 hover:bg-red-900/40 border border-red-800/60 hover:border-red-300 dark:border-red-700 rounded-lg transition-colors"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
              Delete
            </button>
          </template>
        </div>
      </div>

      <!-- ══ PARENT CONTEXT ══ -->
      <div
        v-if="task.parent_task"
        class="flex items-center gap-2 px-4 py-2.5 mb-4 rounded-xl border"
        :class="{
          'bg-purple-900/20 border-purple-700/40': task.parent_task.task_type === 'FEATURE',
          'bg-blue-900/20 border-blue-700/40': task.parent_task.task_type === 'TASK',
          'bg-red-900/20 border-red-700/40': task.parent_task.task_type === 'BUG',
          'bg-gray-800/40 border-gray-200 dark:border-gray-700/40': !task.parent_task.task_type,
        }"
      >
        <NuxtLink
          :to="`/task/${task.parent_task.id}`"
          class="flex items-center gap-2 shrink-0 transition-opacity hover:opacity-100 opacity-90"
        >
          <svg class="w-3.5 h-3.5 shrink-0" :class="{
          'text-purple-400': task.parent_task.task_type === 'FEATURE',
          'text-blue-400': task.parent_task.task_type === 'TASK',
          'text-red-400': task.parent_task.task_type === 'BUG',
          'text-gray-500 dark:text-gray-400': !task.parent_task.task_type,
        }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
        </svg>
        <span class="text-xs cursor-pointer hover:underline underline-offset-1" :class="{
          'text-purple-400/70 hover:text-purple-400': task.parent_task.task_type === 'FEATURE',
          'text-blue-400/70 hover:text-blue-400': task.parent_task.task_type === 'TASK',
          'text-red-400/70 hover:text-red-400': task.parent_task.task_type === 'BUG',
          'text-gray-400/70 hover:text-gray-500 dark:text-gray-400': !task.parent_task.task_type,
        }">Part of {{ task.parent_task.parent_id ? 'Sub-task' : (task.parent_task.task_type === 'FEATURE' ? 'Feature' : task.parent_task.task_type === 'BUG' ? 'Bug' : 'Task') }}</span>
        <svg class="w-3 h-3 shrink-0" :class="{
          'text-purple-600': task.parent_task.task_type === 'FEATURE',
          'text-blue-600': task.parent_task.task_type === 'TASK',
          'text-red-600': task.parent_task.task_type === 'BUG',
          'text-gray-600': !task.parent_task.task_type,
        }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>
        </NuxtLink>
        <NuxtLink
          :to="`/task/${task.parent_task.id}`"
          class="flex items-center gap-1.5 text-xs font-semibold transition-colors group"
          :class="{
            'text-purple-300 hover:text-purple-100': task.parent_task.task_type === 'FEATURE',
            'text-blue-300 hover:text-blue-100': task.parent_task.task_type === 'TASK',
            'text-red-300 hover:text-red-100': task.parent_task.task_type === 'BUG',
            'text-gray-600 dark:text-gray-300 hover:text-gray-100': !task.parent_task.task_type,
          }"
        >
          <span v-if="task.parent_task.task_type === 'FEATURE'" class="text-purple-500">★</span>
          <span v-else-if="task.parent_task.task_type === 'BUG'" class="text-red-400">⚠</span>
          <span v-else class="text-blue-400">📋</span>
          <span class="group-hover:underline underline-offset-2">{{ task.parent_task.title }}</span>
          <code
            v-if="task.parent_task.code"
            class="ml-1 font-mono text-[10px] px-1.5 py-0.5 rounded border"
            :class="{
              'bg-purple-900/40 border-purple-700/50 text-purple-400': task.parent_task.task_type === 'FEATURE',
              'bg-blue-900/40 border-blue-700/50 text-blue-400': task.parent_task.task_type === 'TASK',
              'bg-red-900/40 border-red-700/50 text-red-400': task.parent_task.task_type === 'BUG',
              'bg-slate-800/70 border-gray-200 dark:border-gray-700/50 text-gray-500 dark:text-gray-400': !task.parent_task.task_type,
            }"
          >
            {{ String(Number(task.parent_task.code.split('-').pop()) || 0).padStart(4, '0') }}
          </code>
        </NuxtLink>
      </div>

      <!-- ══ HERO TITLE SECTION ══ -->
      <div class="enterprise-hero-panel rounded-2xl p-6 mb-6 backdrop-blur-sm">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <!-- Type + Status badges -->
            <div class="flex flex-wrap items-center gap-2 mb-3">
              <span
                v-if="task.task_type"
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg text-xs font-semibold border"
                :class="{
                  'bg-purple-900/40 border-purple-600/50 text-purple-300': task.task_type === 'FEATURE',
                  'bg-blue-900/40 border-blue-600/50 text-blue-300': task.task_type === 'TASK',
                  'bg-red-900/40 border-red-600/50 text-red-300': task.task_type === 'BUG',
                }"
              >
                <span v-if="task.task_type === 'FEATURE'" class="text-purple-300">★</span>
                <span v-else-if="task.task_type === 'BUG'" class="text-red-300">⚠</span>
                <span v-else class="text-blue-300">📋</span>
                {{ task.task_type }}
              </span>
              <span :class="getStatusBadgeClass(task.status)" class="inline-flex items-center px-2.5 py-1 rounded-lg text-xs font-semibold border">
                {{ getStatusLabel(task.status) }}
              </span>
              <span
                v-if="task.priority"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg text-xs font-semibold border"
                :class="{
                  'bg-red-900/30 border-red-700/50 text-red-300': task.priority === 'CRITICAL',
                  'bg-orange-900/30 border-orange-700/50 text-orange-300': task.priority === 'HIGH',
                  'bg-yellow-900/30 border-yellow-700/50 text-yellow-300': task.priority === 'MEDIUM',
                  'bg-green-900/30 border-green-700/50 text-green-300': task.priority === 'LOW',
                }"
              >
                <span v-if="task.priority === 'CRITICAL'">🔴</span>
                <span v-else-if="task.priority === 'HIGH'">🟠</span>
                <span v-else-if="task.priority === 'MEDIUM'">🟡</span>
                <span v-else>🟢</span>
                {{ task.priority }}
              </span>
            </div>

            <div v-if="availableStatusTransitions.length > 0" class="flex flex-wrap items-center gap-2 mb-4">
              <span class="text-[11px] uppercase tracking-wider text-gray-500 dark:text-gray-400 font-semibold">Quick status</span>
              <button
                v-for="action in availableStatusTransitions"
                :key="action.status"
                type="button"
                :disabled="statusChangeLoading"
                @click="changeTaskStatus(action.status)"
                class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-sm font-semibold border transition-all shadow-sm disabled:opacity-50"
                :class="action.className"
              >
                <span>{{ action.icon }}</span>
                {{ statusChangeLoading && pendingStatusChange === action.status ? 'Updating…' : action.label }}
              </button>
            </div>

            <!-- Title -->
            <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white leading-tight tracking-tight mb-2">{{ task.title }}</h1>

            <!-- Meta row -->
            <div class="flex flex-wrap items-center gap-x-4 gap-y-1 text-xs text-gray-500">
              <span>Created {{ formatDate(task.created_at) }}</span>
              <span v-if="creatorLabel">by <span class="text-gray-500 dark:text-gray-400">{{ creatorLabel }}</span></span>
              <span v-if="task.story_points" class="flex items-center gap-1">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/></svg>
                {{ task.story_points }} SP
              </span>
            </div>
          </div>

          <!-- Due date pill -->
          <div v-if="task.due_at" class="shrink-0">
            <div
              class="px-4 py-3 rounded-xl border text-center min-w-[110px]"
              :class="{
                'bg-red-900/30 border-red-700/60 text-red-300': getDeadlineUrgency(task) === 'overdue',
                'bg-amber-900/30 border-amber-700/60 text-amber-300': getDeadlineUrgency(task) === 'urgent',
                'bg-gray-50 dark:bg-gray-800/60 border-gray-200 dark:border-gray-700/60 text-gray-600 dark:text-gray-300': getDeadlineUrgency(task) === 'normal',
                'bg-green-900/30 border-green-700/60 text-green-300': task.status === 'COMPLETED',
              }"
            >
              <p class="text-[10px] uppercase tracking-wider opacity-70 mb-0.5">Due</p>
              <p class="text-sm font-semibold">{{ formatDate(task.due_at) }}</p>
              <p v-if="task.status !== 'COMPLETED'" class="text-[11px] mt-0.5 opacity-80">{{ getDeadlineCountdown(task.due_at) }}</p>
            </div>
          </div>
        </div>

        <!-- Progress bar (sub-tasks aggregated) -->
        <div v-if="isParentTask" class="mt-5 pt-4 border-t border-gray-200 dark:border-gray-700/40">
          <div class="flex items-center justify-between text-xs text-gray-500 mb-2">
            <span>Sub-task progress</span>
            <span class="text-gray-500 dark:text-gray-400 font-medium">{{ subtaskAggregateProgress }}%</span>
          </div>
          <div class="h-1.5 bg-gray-100 dark:bg-gray-700/60 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="subtaskAggregateProgress === 100 ? 'bg-green-500' : 'bg-gradient-to-r from-purple-500 to-pink-500'"
              :style="{ width: subtaskAggregateProgress + '%' }"
            />
          </div>
        </div>
      </div>

      <!-- Product Owner step: READY_FOR_TEST — test & approve to send to deploy queue -->
      <section
        v-if="showPMUATActions"
        class="mb-6 rounded-2xl border border-cyan-500/35 bg-cyan-950/20 px-5 py-4"
      >
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div class="min-w-0">
            <p class="text-xs font-semibold uppercase tracking-widest text-cyan-400 mb-1">Awaiting your test approval</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              This sub-task is in <span class="text-cyan-300 font-medium">Ready for Test</span>. Approve with test evidence to move it to <span class="text-orange-300 font-medium">Wait for Deploy</span>, or reject to send it back to the engineer.
            </p>
          </div>
          <div class="flex flex-col sm:flex-row gap-2 shrink-0">
            <button
              type="button"
              :disabled="uatActionLoading || uatApproveSubmitting"
              class="flex items-center justify-center gap-1.5 px-4 py-2.5 rounded-xl bg-emerald-100 dark:bg-emerald-500/15 border border-emerald-300 dark:border-emerald-500/35 text-emerald-400 text-sm font-semibold hover:bg-emerald-100 dark:bg-emerald-500/25 hover:border-emerald-400/50 transition-colors disabled:opacity-50"
              @click="openUATApproveConfirm"
            >
              <span>✅</span>
              APPROVE
            </button>
            <button
              type="button"
              :disabled="uatActionLoading"
              class="flex items-center justify-center gap-1.5 px-4 py-2.5 rounded-xl bg-red-100 dark:bg-red-500/15 border border-red-300 dark:border-red-500/35 text-red-400 text-sm font-semibold hover:bg-red-100 dark:bg-red-500/25 hover:border-red-400/50 transition-colors disabled:opacity-50"
              @click="openUATRejectModal"
            >
              <span>❌</span>
              REJECT
            </button>
          </div>
        </div>
      </section>

      <!-- WAIT_FOR_DEPLOY step: awaiting Chief Engineer deployment -->
      <section
        v-else-if="showWaitForDeploySection"
        class="mb-6 rounded-2xl border border-orange-500/35 bg-orange-950/20 px-5 py-4"
      >
        <div class="flex flex-col gap-4">
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <p class="text-xs font-semibold uppercase tracking-widest text-orange-400 mb-1">🚀 Waiting for Deployment</p>
              <p class="text-sm text-gray-500 dark:text-gray-400">
                Test approved by Product Owner. The Chief Engineer must create a deployment request, get it reviewed, and mark it as deployed before this task advances to
                <span class="text-amber-300 font-medium">Ready for UAT</span>.
              </p>
            </div>
          </div>

          <!-- No deployment request yet -->
          <div v-if="!deploymentLoading && !deploymentForTask" class="flex flex-col sm:flex-row items-start sm:items-center gap-3 px-4 py-3 rounded-xl bg-orange-500/10 border border-orange-500/25">
            <div class="flex items-center gap-2 flex-1 min-w-0">
              <span class="text-orange-400 text-lg shrink-0">⚠️</span>
              <p class="text-sm text-orange-300 font-medium">No deployment request created yet.</p>
            </div>
            <button
              @click="openDeploymentModal()"
              class="shrink-0 flex items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-orange-100 dark:from-orange-600 to-amber-100 dark:to-amber-600 hover:from-orange-200 dark:hover:from-orange-500 hover:to-amber-200 dark:hover:to-amber-500 text-gray-900 dark:text-white text-sm font-semibold transition-all shadow-lg"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
              </svg>
              New Deployment Request
            </button>
          </div>

          <!-- Deployment request exists — show status -->
          <div v-else-if="deploymentForTask" class="px-4 py-3 rounded-xl bg-gray-50 dark:bg-gray-800/60 border border-gray-200 dark:border-gray-700/50 space-y-2">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="flex items-center gap-2 min-w-0">
                <span class="text-sm font-semibold text-gray-900 dark:text-white truncate">{{ deploymentForTask.title }}</span>
                <span class="text-[10px] font-bold px-2 py-0.5 rounded-md border uppercase"
                  :class="{
                    'bg-yellow-500/10 text-yellow-400 border-yellow-500/30': deploymentForTask.status === 'PENDING',
                    'bg-blue-500/10 text-blue-400 border-blue-500/30': deploymentForTask.status === 'REVIEWING',
                    'bg-green-500/10 text-green-400 border-green-500/30': deploymentForTask.status === 'APPROVED',
                    'bg-red-500/10 text-red-400 border-red-500/30': deploymentForTask.status === 'REJECTED',
                    'bg-cyan-500/10 text-cyan-400 border-cyan-500/30': deploymentForTask.status === 'DEPLOYED',
                  }">{{ deploymentForTask.status }}</span>
              </div>
              <NuxtLink to="/deployment" class="text-xs text-blue-400 hover:text-blue-300 hover:underline shrink-0">View in Deployment →</NuxtLink>
            </div>
            <div class="flex items-center gap-2 text-xs text-gray-500">
              <code class="text-cyan-400 font-mono">⎇ {{ deploymentForTask.branch }}</code>
              <span>{{ deploymentForTask.environment }}</span>
            </div>
            <p v-if="deploymentForTask.status === 'REJECTED'" class="text-xs text-red-300 mt-1">
              ✗ Rejected: {{ deploymentForTask.rejection_reason }}
            </p>
            <p v-if="deploymentForTask.status === 'DEPLOYED'" class="text-xs text-green-400 mt-1">
              ✓ Deployed — task will advance to Ready for UAT automatically.
            </p>
          </div>

          <!-- Loading -->
          <div v-else class="flex items-center gap-2 text-xs text-gray-500 px-1">
            <div class="w-3 h-3 rounded-full border-2 border-gray-500 border-t-transparent animate-spin" />
            Checking deployment status…
          </div>
        </div>
      </section>

      <!-- CEO step: READY_FOR_UAT — final approval after deployment -->
      <section
        v-else-if="showCEOUATActions"
        class="mb-6 rounded-2xl border border-amber-500/35 bg-amber-950/20 px-5 py-4"
      >
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div class="min-w-0">
            <p class="text-xs font-semibold uppercase tracking-widest text-amber-400 mb-1">CEO Final Approval Required</p>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              The Product Owner has tested this task and submitted evidence. Review the test details below and give your <span class="text-amber-300 font-medium">final approval</span> or reject it back to the Product Owner.
            </p>
            <!-- Display Product Owner test evidence from uat_payload -->
            <template v-if="uatPayloadData">
              <div class="mt-3 space-y-2">
                <div v-if="uatPayloadData.test_url" class="flex items-center gap-2">
                  <span class="text-[10px] font-semibold uppercase tracking-wider text-gray-500 w-20 shrink-0">Test URL</span>
                  <a :href="uatPayloadData.test_url" target="_blank" rel="noopener" class="text-xs text-blue-400 hover:text-blue-300 underline truncate max-w-xs">{{ uatPayloadData.test_url }}</a>
                </div>
                <div v-if="uatPayloadData.test_steps" class="flex items-start gap-2">
                  <span class="text-[10px] font-semibold uppercase tracking-wider text-gray-500 w-20 shrink-0 mt-0.5">Test Steps</span>
                  <pre class="text-xs text-gray-600 dark:text-gray-300 whitespace-pre-wrap break-words max-w-lg">{{ uatPayloadData.test_steps }}</pre>
                </div>
              </div>
            </template>
          </div>
          <div class="flex flex-col sm:flex-row gap-2 shrink-0">
            <button
              type="button"
              :disabled="uatActionLoading || uatApproveSubmitting"
              class="flex items-center justify-center gap-1.5 px-4 py-2.5 rounded-xl bg-emerald-100 dark:bg-emerald-500/15 border border-emerald-300 dark:border-emerald-500/35 text-emerald-400 text-sm font-semibold hover:bg-emerald-100 dark:bg-emerald-500/25 hover:border-emerald-400/50 transition-colors disabled:opacity-50"
              @click="openUATApproveConfirm"
            >
              <span>✅</span>
              FINAL APPROVE
            </button>
            <button
              type="button"
              :disabled="uatActionLoading"
              class="flex items-center justify-center gap-1.5 px-4 py-2.5 rounded-xl bg-red-100 dark:bg-red-500/15 border border-red-300 dark:border-red-500/35 text-red-400 text-sm font-semibold hover:bg-red-100 dark:bg-red-500/25 hover:border-red-400/50 transition-colors disabled:opacity-50"
              @click="openUATRejectModal"
            >
              <span>❌</span>
              REJECT
            </button>
          </div>
        </div>
      </section>

      <!-- ══ TWO-COLUMN LAYOUT ══ -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

        <!-- ── LEFT: Description ── -->
        <div class="lg:col-span-2 space-y-6">

          <!-- Description Card -->
          <div class="enterprise-card rounded-2xl overflow-hidden">
            <div class="flex items-center justify-between px-5 py-3.5 border-b border-gray-200 dark:border-gray-700/60 bg-gray-50 dark:bg-gray-800/60">
              <div class="flex items-center gap-2">
                <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/></svg>
                <h2 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Description</h2>
              </div>
              <div class="flex items-center gap-2">
                <a
                  v-if="slideOpenInSlidesURL"
                  :href="slideOpenInSlidesURL"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-1.5 text-xs text-purple-400 hover:text-purple-300 px-2.5 py-1 rounded-lg bg-purple-900/20 border border-purple-700/40 hover:border-purple-600/60 transition-colors"
                >
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/></svg>
                  Open in Slides
                </a>
                <template v-if="canEditDescription">
                  <template v-if="!isEditingDescription">
                    <button @click="startInlineEdit" class="flex items-center gap-1 text-xs text-gray-500 hover:text-blue-400 px-2 py-1 rounded-lg hover:bg-blue-900/20 transition-colors">
                      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                      Edit
                    </button>
                  </template>
                  <template v-else>
                    <button @click="saveInlineDescription" :disabled="isSavingDescription" class="text-xs px-3 py-1 bg-blue-100 dark:bg-blue-600 hover:bg-blue-200 dark:bg-blue-700 disabled:opacity-50 text-gray-900 dark:text-white rounded-lg transition-colors">
                      {{ isSavingDescription ? 'Saving…' : 'Save' }}
                    </button>
                    <button @click="cancelInlineEdit" class="text-xs px-3 py-1 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-600 dark:text-gray-300 rounded-lg transition-colors">Cancel</button>
                  </template>
                </template>
              </div>
            </div>
            <div class="p-5">
              <RichTextEditor v-if="isEditingDescription" v-model="inlineDescriptionHtml" placeholder="Describe what needs to be done… (paste images with ⌘V)" />
              <div v-else>
                <!-- Skeleton shimmer: only shown when we know there IS rich content to load -->
                <div v-if="richDetailLoading && hasRichContent" class="space-y-3 py-2 animate-pulse">
                  <div class="h-4 bg-gray-700/60 rounded-md w-full" />
                  <div class="h-4 bg-gray-700/60 rounded-md w-5/6" />
                  <div class="h-4 bg-gray-700/60 rounded-md w-4/6" />
                  <div class="h-32 bg-gray-700/40 rounded-lg w-full mt-4" />
                  <div class="h-4 bg-gray-700/60 rounded-md w-3/4" />
                  <div class="h-4 bg-gray-700/60 rounded-md w-2/3" />
                </div>
                <RichTextEditor v-else-if="task.description && task.description.trim()" :model-value="task.description" :readonly="true" />
                <p v-else-if="!richDetailLoading" class="text-gray-500 text-sm italic py-4 text-center">No description yet. Click Edit to add one.</p>
              </div>
            </div>
          </div>

          <!-- Sub-tasks -->
          <SubtaskList
            :parent-task-id="task.id"
            :project-id="task.project_id"
            :parent-task="task"
            :subtasks="subtasks"
            :can-edit="canManageSubtasks"
            :is-max-depth="!!(task.parent_task?.parent_id)"
            @refresh="fetchTask"
          />

          <!-- Discussion & Time Tracking -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Comments -->
            <div class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700/60 rounded-2xl overflow-hidden">
              <div class="flex items-center justify-between px-5 py-3.5 border-b border-gray-200 dark:border-gray-700/60 bg-gray-50 dark:bg-gray-800/60">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/></svg>
                  <h2 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Discussion</h2>
                </div>
                <span class="text-xs text-gray-600 bg-gray-100 dark:bg-gray-700/60 px-2 py-0.5 rounded-full">{{ comments.length }}</span>
              </div>
              <div class="p-5">
                <TaskComments
                  :comments="comments"
                  :loading="commentsLoading"
                  :current-user-avatar="currentUserAvatarURL"
                  :current-user-initial="currentUserInitial"
                  :current-user-id="Number(effectiveUser?.id || authStore.userId || 0)"
                  @add-comment="handleAddComment"
                  @edit-comment="handleEditComment"
                />
              </div>
            </div>

            <!-- Time Tracking -->
            <div class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700/60 rounded-2xl overflow-hidden">
              <div class="flex items-center gap-2 px-5 py-3.5 border-b border-gray-200 dark:border-gray-700/60 bg-gray-50 dark:bg-gray-800/60">
                <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                <h2 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Time Tracking</h2>
              </div>
              <div class="p-5">
                <TimeLogger
                  :time-logs="timeLogs"
                  :estimated-minutes="task.estimated_minutes || 0"
                  :task-id="route.params.id as string"
                  :loading="timeLogsLoading"
                  @log-time="handleLogTime"
                  @refresh="fetchCommentsAndLogs"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- ── RIGHT: Details Sidebar ── -->
        <div class="lg:col-span-1 space-y-4">

          <!-- Parent-task notice -->
          <div v-if="isParentTask" class="flex items-start gap-3 p-4 bg-amber-900/20 border border-amber-700/40 rounded-2xl">
            <svg class="w-4 h-4 text-amber-400 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            <p class="text-xs text-amber-300 leading-relaxed">
              Parent Task — you can estimate time at this level (sub-tasks may also exist).
            </p>
          </div>

          <!-- Details card -->
          <div class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700/60 rounded-2xl overflow-hidden">
            <div class="px-5 py-3.5 border-b border-gray-200 dark:border-gray-700/60 bg-gray-50 dark:bg-gray-800/60">
              <h2 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Details</h2>
            </div>
            <div class="divide-y divide-gray-200 dark:divide-gray-700/40">

              <!-- Project board link (not shown for Komgrip tasks) -->
              <div v-if="task.project_id && !task.is_komgrip" class="px-5 py-3.5">
                <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1.5">Project</p>
                <NuxtLink
                  :to="`/projects/${task.project_id}?tab=board`"
                  class="inline-flex items-center gap-1.5 text-sm text-violet-400 hover:text-violet-300 transition-colors group"
                >
                  <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2h-2a2 2 0 00-2 2"/></svg>
                  <span class="group-hover:underline">View on project board</span>
                  <svg class="w-3 h-3 opacity-60" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
                </NuxtLink>
              </div>

              <!-- Komgrip badge -->
              <div v-if="task.is_komgrip" class="px-5 py-3.5">
                <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1.5">Workspace</p>
                <NuxtLink to="/komgrip" class="inline-flex items-center gap-1.5 text-sm text-violet-400 hover:text-violet-300 transition-colors group">
                  <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/></svg>
                  <span class="group-hover:underline">Komgrip</span>
                </NuxtLink>
              </div>

              <!-- Assignee -->
              <div class="px-5 py-3.5">
                <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1.5">Assignee</p>
                <div class="flex items-center gap-2 flex-wrap">
                  <div v-if="task.assigned_to" class="flex items-center gap-2">
                    <div class="w-7 h-7 rounded-full overflow-hidden bg-gradient-to-br from-purple-600 to-pink-600 flex items-center justify-center text-xs font-bold text-gray-900 dark:text-white shrink-0">
                      <img
                        v-if="task.assigned_to_avatar_url"
                        :src="task.assigned_to_avatar_url"
                        alt="Assignee avatar"
                        class="w-full h-full object-cover"
                      />
                      <span v-else>
                        {{ (task.assigned_to_display_name || task.assigned_to_email || 'U').charAt(0).toUpperCase() }}
                      </span>
                    </div>
                    <span class="text-sm text-gray-900 dark:text-white">{{ task.assigned_to_display_name || task.assigned_to_email || `Dev #${task.assigned_to}` }}</span>
                  </div>
                  <span v-else class="text-sm text-gray-500">Unassigned</span>
                  <button
                    v-if="canClaimTask"
                    type="button"
                    @click="claimTask"
                    :disabled="assignLoading"
                    class="text-[11px] text-emerald-300 hover:text-emerald-200 px-2 py-0.5 rounded-md bg-emerald-900/20 border border-emerald-800/40 hover:border-emerald-300 dark:border-emerald-700/60 transition-colors disabled:opacity-50"
                  >
                    {{ assignLoading ? 'Claiming…' : 'Claim task' }}
                  </button>
                  <button
                    v-if="canEditOrDelete && !showAssignDropdown"
                    type="button"
                    @click="openAssignDropdown"
                    class="text-[11px] px-2 py-0.5 rounded-md border shadow-sm transition-all bg-blue-50 hover:bg-blue-100 text-blue-700 hover:text-blue-800 border-blue-200 hover:border-blue-300 dark:bg-blue-900/20 dark:border-blue-700/60 dark:text-blue-400 dark:hover:text-blue-300"
                  >
                    Change
                  </button>
                </div>
                <template v-if="canEditOrDelete && showAssignDropdown">
                  <select
                    v-model="assignSelectedId"
                    @change="confirmChangeAssignee"
                    class="mt-2 block w-full rounded-xl border border-gray-600 bg-gray-800 px-3 py-2 text-sm text-gray-900 dark:text-white focus:border-purple-500 focus:outline-none focus:ring-1 focus:ring-purple-500"
                  >
                    <option value="">— Select —</option>
                    <option value="0">— Unassign —</option>
                    <option v-for="u in visibleAssigneeUsers" :key="u.id" :value="u.id">{{ u.display_name ? `${u.display_name} (${u.email})` : u.email }} ({{ u.role }})</option>
                  </select>
                  <div class="flex items-center gap-2 mt-1.5">
                    <button type="button" @click="showAssignDropdown = false" class="text-xs text-gray-500 hover:text-gray-600 dark:text-gray-300">Cancel</button>
                    <p v-if="assignError" class="text-xs text-red-400">{{ assignError }}</p>
                  </div>
                </template>
              </div>

              <!-- Estimated Effort (hours, 1 decimal; stored as minutes API-side) -->
              <div class="px-5 py-3.5">
                <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1.5">Estimated Effort</p>
                <div v-if="!canEditOrDelete" class="text-sm font-semibold text-gray-900 dark:text-white">
                  {{ formatMinutesAsHours(task.estimated_minutes ?? 0) }}
                  <span class="text-gray-500 dark:text-gray-400 font-normal text-xs">h</span>
                  <span class="text-gray-500 font-normal text-xs ml-1">({{ task.estimated_minutes ?? 0 }} min)</span>
                </div>
                <div v-else class="flex items-center gap-2 flex-wrap">
                  <input
                    v-model.number="estimatedHoursLocal"
                    type="number" min="0" step="0.5"
                    class="w-28 px-2.5 py-1.5 bg-gray-900 border border-gray-600 rounded-xl text-gray-100 text-sm focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="0"
                    @blur="saveEstimatedMinutes"
                  />
                  <span class="text-xs text-gray-500">hours</span>
                  <button
                    v-if="estimatedMinutesDirty"
                    type="button"
                    class="text-xs px-2.5 py-1 bg-blue-100 dark:bg-blue-600 hover:bg-blue-200 dark:bg-blue-700 text-gray-900 dark:text-white rounded-lg disabled:opacity-50 transition-colors"
                    :disabled="isSavingEstimate"
                    @click="saveEstimatedMinutes"
                  >
                    {{ isSavingEstimate ? '…' : 'Save' }}
                  </button>
                </div>
              </div>

              <!-- Dates group -->
              <div v-if="task.start_date || task.end_date || task.completed_at" class="px-5 py-3.5 space-y-3">
                <div v-if="task.start_date">
                  <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1">Start</p>
                  <p class="text-sm text-gray-600 dark:text-gray-300">{{ formatDateTime(task.start_date) }}</p>
                </div>
                <div v-if="task.end_date">
                  <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1">End</p>
                  <p class="text-sm text-gray-600 dark:text-gray-300">{{ formatDateTime(task.end_date) }}</p>
                </div>
                <div v-if="task.completed_at">
                  <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1">Completed</p>
                  <p class="text-sm text-green-400">{{ formatDateTime(task.completed_at) }}</p>
                  <p v-if="task.started_at" class="text-xs text-gray-500 mt-0.5">Duration: {{ calculateDuration(task.started_at, task.completed_at) }}</p>
                </div>
              </div>

              <!-- Outsource to Team (Product Owner / CEO only) -->
              <div v-if="canEditOrDelete" class="px-5 py-3.5">
                <p class="text-[11px] text-gray-500 uppercase tracking-wider mb-1.5">B2B Outsource</p>
                <button
                  @click="showOutsourceModal = true"
                  class="w-full flex items-center justify-center gap-2 rounded-xl border border-blue-300 dark:border-blue-500/30 bg-blue-100 dark:bg-blue-500/10 hover:bg-blue-100 dark:bg-blue-500/20 px-3 py-2 text-xs font-semibold text-blue-300 transition-colors"
                >
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 4H6a2 2 0 00-2 2v12a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-2m-4-1v8m0 0l3-3m-3 3L9 8m-5 5h2.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293h3.172a1 1 0 00.707-.293l2.414-2.414A1 1 0 0014.414 13H17"/>
                  </svg>
                  Outsource to Team
                </button>
              </div>

            </div>
          </div>

          <!-- Activity timeline (audit trail) -->
          <div class="bg-white dark:bg-gray-800/50 border border-gray-200 dark:border-gray-700/60 rounded-2xl overflow-hidden">
            <div class="px-5 py-3.5 border-b border-gray-200 dark:border-gray-700/60 bg-gray-50 dark:bg-gray-800/60 flex items-center justify-between gap-2">
              <div class="flex items-center gap-2 min-w-0">
                <svg class="w-4 h-4 text-violet-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h2 class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">Activity</h2>
              </div>
              <button
                type="button"
                class="text-[11px] text-violet-400 hover:text-violet-300 px-2 py-1 rounded-lg border border-violet-300 dark:border-violet-500/25 hover:border-violet-400/40 bg-violet-100 dark:bg-violet-500/5 transition-colors disabled:opacity-40 shrink-0"
                :disabled="activityLoading"
                @click="fetchTaskActivity"
              >
                Refresh
              </button>
            </div>
            <div class="p-4 max-h-[min(52vh,520px)] overflow-y-auto custom-scrollbar">
              <div v-if="activityLoading" class="flex flex-col items-center justify-center py-10 gap-2">
                <div class="w-8 h-8 rounded-xl bg-violet-600/20 border border-violet-500/30 flex items-center justify-center">
                  <svg class="w-4 h-4 text-violet-400 animate-spin" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                  </svg>
                </div>
                <p class="text-xs text-gray-500">Loading activity…</p>
              </div>
              <p v-else-if="activityError" class="text-xs text-red-400 py-6 text-center px-2">{{ activityError }}</p>
              <p v-else-if="taskActivity.length === 0" class="text-xs text-gray-500 py-8 text-center leading-relaxed px-2">
                No activity entries yet.<br />
                <span class="text-gray-600">New events will appear here automatically.</span>
              </p>
              <ol v-else class="relative ms-2 border-s border-gray-200 dark:border-gray-700/70 ps-0 list-none space-y-0">
                <li
                  v-for="(item, idx) in taskActivity"
                  :key="item.id + '-' + idx"
                  class="relative pb-7 last:pb-1 ps-6"
                >
                  <span class="absolute start-0 top-1 -translate-x-1/2 flex h-7 w-7 rounded-full bg-white dark:bg-gray-900 shadow-sm">
                    <span
                      class="flex h-full w-full items-center justify-center rounded-full border text-[11px] font-bold"
                      :class="activityDotClass(item.action)"
                    >
                      {{ activityIconGlyph(item.action) }}
                    </span>
                  </span>
                  <div class="pt-0.5 space-y-1">
                    <div class="flex flex-wrap items-center gap-2">
                      <span class="text-sm font-semibold text-gray-900 dark:text-white leading-snug">{{ activityTitle(item) }}</span>
                      <span
                        v-if="item.inferred"
                        class="text-[10px] font-semibold uppercase tracking-wide px-2 py-0.5 rounded-md bg-amber-500/10 text-amber-400 border border-amber-500/25"
                      >Inferred from record</span>
                    </div>
                    <p class="text-[11px] text-gray-500 font-mono tabular-nums">{{ formatActivityDateTime(item.at) }}</p>
                    <p class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed break-words">{{ activityDetailLine(item) }}</p>
                  </div>
                </li>
              </ol>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ══ EDIT MODAL ══ -->
    <div v-if="showEditModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-start justify-center z-50 p-3 sm:p-6 overflow-y-auto" @click.self="closeEditModal">
      <div class="edit-task-modal edit-task-modal-enterprise rounded-2xl w-full max-w-7xl my-4 sm:my-8 flex flex-col max-h-[calc(100dvh-2rem)] min-h-0">
        <div class="flex items-center justify-between px-6 sm:px-8 pt-6 sm:pt-8 pb-4 shrink-0 border-b border-gray-200 dark:border-gray-700/80">
          <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white tracking-tight">Edit Task</h2>
          <button type="button" @click="closeEditModal" class="shrink-0 w-11 h-11 flex items-center justify-center rounded-xl text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:text-white hover:bg-gray-700 transition-colors disabled:opacity-40" :disabled="isUpdatingTask" aria-label="Close">✕</button>
        </div>
        <div class="overflow-y-auto flex-1 px-6 sm:px-8 py-6 sm:py-8 space-y-6 sm:space-y-7 min-h-0 overscroll-contain">
          <div v-if="editError" class="p-4 md:p-5 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-base">{{ editError }}</div>
          <div>
            <label class="label">Type *</label>
            <div class="grid grid-cols-3 gap-3 sm:gap-4">
              <button type="button" @click="editForm.task_type = 'FEATURE'" :class="editForm.task_type === 'FEATURE' ? 'border-purple-300 dark:border-purple-500 bg-purple-100 dark:bg-purple-500/20 text-purple-300' : 'border-gray-600 bg-gray-50 dark:bg-gray-900/50 text-gray-500 dark:text-gray-400 hover:border-purple-300 dark:border-purple-500/50'" class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]" :disabled="isUpdatingTask"><span class="text-xl sm:text-2xl leading-none">★</span> Feature</button>
              <button type="button" @click="editForm.task_type = 'TASK'" :class="editForm.task_type === 'TASK' ? 'border-blue-300 dark:border-blue-500 bg-blue-100 dark:bg-blue-500/20 text-blue-300' : 'border-gray-600 bg-gray-50 dark:bg-gray-900/50 text-gray-500 dark:text-gray-400 hover:border-blue-300 dark:border-blue-500/50'" class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]" :disabled="isUpdatingTask"><span class="text-xl sm:text-2xl leading-none">📋</span> Task</button>
              <button type="button" @click="editForm.task_type = 'BUG'" :class="editForm.task_type === 'BUG' ? 'border-red-300 dark:border-red-500 bg-red-100 dark:bg-red-500/20 text-red-300' : 'border-gray-600 bg-gray-50 dark:bg-gray-900/50 text-gray-500 dark:text-gray-400 hover:border-red-300 dark:border-red-500/50'" class="flex flex-col items-center justify-center gap-1.5 px-4 py-4 sm:py-5 rounded-xl border text-sm sm:text-base font-semibold transition-all min-h-[4.5rem]" :disabled="isUpdatingTask"><span class="text-xl sm:text-2xl leading-none">⚠</span> Bug</button>
            </div>
            <div v-if="editForm.task_type === 'FEATURE'" class="mt-3 flex items-start gap-3 p-4 bg-purple-900/20 border border-purple-500/30 rounded-xl text-sm sm:text-base text-purple-300 leading-relaxed">
              <span class="shrink-0 mt-0.5">★</span>
              <span><strong>Feature mode:</strong> Acts as a parent container. Assignee and estimated effort are disabled.</span>
            </div>
          </div>
          <div>
            <label class="label">Title *</label>
            <input v-model="editForm.title" type="text" placeholder="Task title…" class="input-field w-full" :disabled="isUpdatingTask" />
          </div>
          <div>
            <label class="label">Description</label>
            <RichTextEditor v-model="editForm.description" placeholder="Describe the task objectives… (paste images with ⌘V)" />
          </div>
          <div>
            <label class="label" :class="editForm.task_type === 'FEATURE' ? 'text-gray-500' : ''">
              Estimated Effort (hours)
              <span v-if="editForm.task_type === 'FEATURE'" class="text-gray-600 font-normal">(disabled for Features)</span>
            </label>
            <template v-if="isParentTask && editForm.task_type !== 'FEATURE'">
              <div class="flex items-center gap-3 px-4 py-4 bg-gray-100 dark:bg-gray-900/60 border border-amber-700/30 rounded-xl text-amber-200 text-base font-medium">{{ formatMinutesAsHours(subtaskTotalEstimatedMinutes) }} h (roll-up, {{ subtaskTotalEstimatedMinutes }} min)</div>
            </template>
            <template v-else>
              <input v-model.number="editForm.estimated_hours" type="number" min="0" step="0.5" class="input-field w-full" :class="editForm.task_type === 'FEATURE' ? 'opacity-40 cursor-not-allowed' : ''" :disabled="isUpdatingTask || editForm.task_type === 'FEATURE'" placeholder="e.g. 1.5" />
            </template>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Priority</label>
              <select v-model="editForm.priority" class="input-field w-full" :disabled="isUpdatingTask">
                <option value="CRITICAL">🔴 Critical</option>
                <option value="HIGH">🟠 High</option>
                <option value="MEDIUM">🟡 Medium</option>
                <option value="LOW">🟢 Low</option>
              </select>
            </div>
            <div>
              <label class="label">Story Points</label>
              <input v-model.number="editForm.story_points" type="number" min="0" step="0.5" class="input-field w-full" :disabled="isUpdatingTask" />
            </div>
          </div>
          <div v-if="editSprints.length > 0">
            <label class="label">Sprint</label>
            <select v-model="editForm.sprint_id" class="input-field w-full" :disabled="isUpdatingTask">
              <option value="">Backlog</option>
              <option v-for="s in editSprints" :key="s.id" :value="s.id">{{ s.name }}</option>
            </select>
          </div>
          <div>
            <label class="label">Due Date</label>
            <UiDatePicker v-model="editForm.deadline" placeholder="Select due date…" :disabled="isUpdatingTask" @update:modelValue="closeDateTimePicker" />
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-5">
            <div>
              <label class="label">Start Date</label>
              <UiDatePicker v-model="editForm.start_date" placeholder="Select start date…" :disabled="isUpdatingTask" @update:modelValue="closeDateTimePicker" />
            </div>
            <div>
              <label class="label">End Date</label>
              <UiDatePicker v-model="editForm.end_date" placeholder="Select end date…" :disabled="isUpdatingTask" @update:modelValue="closeDateTimePicker" />
            </div>
          </div>
        </div>
        <div class="flex flex-col-reverse sm:flex-row gap-3 sm:gap-4 px-6 sm:px-8 py-5 sm:py-6 border-t border-gray-200 dark:border-gray-700/60 shrink-0">
          <button type="button" @click="submitEdit" :disabled="isUpdatingTask || !editForm.title.trim()" class="flex-1 btn-primary py-4 text-base sm:text-lg font-semibold rounded-xl disabled:opacity-40 flex items-center justify-center gap-2 min-h-[3.25rem]">
            <svg v-if="isUpdatingTask" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
            {{ isUpdatingTask ? 'Saving…' : 'Save Changes' }}
          </button>
          <button type="button" @click="closeEditModal" :disabled="isUpdatingTask" class="sm:shrink-0 px-6 py-4 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-40 text-gray-700 dark:text-gray-200 rounded-xl transition-colors text-base font-medium min-h-[3.25rem]">Cancel</button>
        </div>
      </div>
    </div>

    <!-- ══ DELETE MODAL ══ -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeDeleteModal">
      <div class="bg-gray-800 border border-red-900/60 rounded-2xl p-6 max-w-md w-full shadow-2xl">
        <div class="flex items-center gap-3 mb-5">
          <div class="w-10 h-10 rounded-xl bg-red-900/40 border border-red-700/50 flex items-center justify-center shrink-0">
            <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
          </div>
          <div>
            <h2 class="text-base font-bold text-gray-900 dark:text-white">Delete Task?</h2>
            <p class="text-xs text-gray-500 mt-0.5">This action cannot be undone</p>
          </div>
          <button @click="closeDeleteModal" :disabled="isDeletingTask" class="ml-auto text-gray-500 hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
        </div>
        <div class="p-4 bg-gray-100 dark:bg-gray-900/60 border border-gray-200 dark:border-gray-700/60 rounded-xl mb-4">
          <p class="text-xs text-gray-500 mb-1">Task to delete</p>
          <p class="text-sm font-semibold text-gray-900 dark:text-white">{{ task?.title }}</p>
          <p class="text-xs text-gray-600 mt-0.5 font-mono">{{ task?.id }}</p>
        </div>
        <div v-if="deleteError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded-xl text-red-400 text-sm">{{ deleteError }}</div>
        <div class="flex gap-3">
          <button @click="confirmDelete" :disabled="isDeletingTask" class="flex-1 px-4 py-2.5 bg-red-100 dark:bg-red-600 hover:bg-red-200 dark:bg-red-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-gray-900 dark:text-white text-sm font-semibold rounded-xl transition-colors flex items-center justify-center gap-2">
            <svg v-if="isDeletingTask" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
            {{ isDeletingTask ? 'Deleting…' : 'Yes, Delete Forever' }}
          </button>
          <button @click="closeDeleteModal" :disabled="isDeletingTask" class="px-5 py-2.5 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 disabled:opacity-40 text-gray-600 dark:text-gray-300 text-sm font-medium rounded-xl transition-colors">Cancel</button>
        </div>
      </div>
    </div>

    <!-- ══ Product Owner approve: Test Evidence Form (READY_FOR_TEST → READY_FOR_UAT) ══ -->
    <Teleport to="body">
      <div
        v-if="uatApproveConfirmOpen && showPMUATActions"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click.self="!uatApproveSubmitting && closeUATApproveConfirm()"
      >
        <div class="w-full max-w-lg rounded-2xl border border-emerald-500/30 bg-gray-900 shadow-2xl p-6">
          <div class="flex items-start gap-3 mb-5">
            <div class="w-8 h-8 rounded-lg bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center flex-shrink-0 mt-0.5">
              <span class="text-lg" aria-hidden="true">🧪</span>
            </div>
            <div class="min-w-0">
              <h3 class="text-sm font-bold text-gray-900 dark:text-white">Submit Test Evidence to CEO</h3>
              <p class="text-xs text-gray-500 truncate max-w-[320px]">{{ task?.title }}</p>
              <p class="text-xs text-amber-400/80 mt-1">Task will be forwarded to CEO for final approval — not marked as Done yet.</p>
            </div>
          </div>
          <!-- Test URL -->
          <div class="mb-4">
            <label class="block text-[11px] font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1.5">
              Test / Staging URL <span class="text-red-400">*</span>
            </label>
            <input
              ref="uatTestUrlRef"
              v-model="uatTestUrl"
              type="url"
              placeholder="https://staging.example.com/feature-xyz"
              :disabled="uatApproveSubmitting"
              class="w-full rounded-xl border border-gray-700 bg-gray-50 dark:bg-gray-800/60 px-3 py-2.5 text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-600 focus:border-emerald-500/50 focus:outline-none focus:ring-1 focus:ring-emerald-500/30 disabled:opacity-50"
            />
            <p v-if="uatTestUrl.length > 0 && !uatTestUrl.startsWith('http')" class="text-[11px] text-red-400 mt-1">
              URL must start with http:// or https://
            </p>
          </div>
          <!-- Test Steps -->
          <div class="mb-5">
            <label class="block text-[11px] font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1.5">
              Test Steps for CEO <span class="text-red-400">*</span>
            </label>
            <textarea
              v-model="uatTestSteps"
              rows="6"
              placeholder="Describe step-by-step how the CEO should test this feature:&#10;1. Navigate to...&#10;2. Click on...&#10;3. Verify that..."
              :disabled="uatApproveSubmitting"
              class="w-full rounded-xl border border-gray-700 bg-gray-50 dark:bg-gray-800/60 px-3 py-2.5 text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-600 focus:border-emerald-500/50 focus:outline-none focus:ring-1 focus:ring-emerald-500/30 resize-none disabled:opacity-50"
            />
            <div class="flex items-center justify-between mt-1">
              <p v-if="uatTestSteps.length > 0 && uatTestSteps.length < 20" class="text-[11px] text-red-400">
                At least {{ 20 - uatTestSteps.length }} more character(s) required
              </p>
              <span class="text-[11px] text-gray-600 ml-auto">{{ uatTestSteps.length }} chars</span>
            </div>
          </div>
          <div class="flex items-center justify-end gap-3">
            <button
              type="button"
              :disabled="uatApproveSubmitting"
              class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-500 dark:text-gray-400 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-700 dark:text-gray-200 transition-colors disabled:opacity-50"
              @click="closeUATApproveConfirm"
            >Cancel</button>
            <button
              type="button"
              :disabled="uatApproveSubmitting || !isUATApproveFormValid"
              class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-emerald-100 dark:bg-emerald-500/15 border border-emerald-300 dark:border-emerald-500/40 text-emerald-400 text-xs font-bold hover:bg-emerald-100 dark:bg-emerald-500/25 transition-colors disabled:opacity-50"
              @click="submitUATApprove"
            >
              <svg v-if="uatApproveSubmitting" class="w-3.5 h-3.5 animate-spin shrink-0" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              <span v-if="uatApproveSubmitting">Submitting…</span>
              <span v-else>✅ Submit to CEO</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ══ CEO FINAL APPROVE CONFIRM (READY_FOR_UAT → COMPLETED) ══ -->
    <Teleport to="body">
      <div
        v-if="uatApproveConfirmOpen && showCEOUATActions"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click.self="closeUATApproveConfirm"
      >
        <div class="w-full max-w-md rounded-2xl border border-amber-500/30 bg-gray-900 shadow-2xl p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-8 h-8 rounded-lg bg-amber-500/15 border border-amber-500/30 flex items-center justify-center flex-shrink-0">
              <span class="text-lg" aria-hidden="true">👑</span>
            </div>
            <div>
              <h3 class="text-sm font-bold text-gray-900 dark:text-white">Final Approval — Mark as Done?</h3>
              <p class="text-xs text-gray-500 truncate max-w-[280px]">{{ task?.title }}</p>
            </div>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-4">
            This marks the task as <span class="text-emerald-400 font-medium">COMPLETED</span>. Continue only if you have verified the test evidence.
          </p>
          <div class="flex items-center justify-end gap-3">
            <button
              type="button"
              :disabled="uatApproveSubmitting"
              class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-500 dark:text-gray-400 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-700 dark:text-gray-200 transition-colors disabled:opacity-50"
              @click="closeUATApproveConfirm"
            >Cancel</button>
            <button
              type="button"
              :disabled="uatApproveSubmitting"
              class="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-lg bg-amber-100 dark:bg-amber-500/15 border border-amber-300 dark:border-amber-500/40 text-amber-400 text-xs font-bold hover:bg-amber-100 dark:bg-amber-500/25 transition-colors disabled:opacity-50"
              @click="submitUATApprove"
            >
              <svg v-if="uatApproveSubmitting" class="w-3.5 h-3.5 animate-spin shrink-0" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              <span v-if="uatApproveSubmitting">Approving…</span>
              <span v-else>👑 Yes, Final Approve</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ══ UAT REJECT MODAL (ready-for-test → IN_PROGRESS) ══ -->
    <Teleport to="body">
      <div
        v-if="uatRejectModalOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
        @click.self="closeUATRejectModal"
      >
        <div class="w-full max-w-md rounded-2xl border border-red-500/30 bg-gray-900 shadow-2xl p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-8 h-8 rounded-lg bg-red-500/15 border border-red-500/30 flex items-center justify-center flex-shrink-0">
              <svg class="w-4 h-4 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </div>
            <div>
              <h3 class="text-sm font-bold text-gray-900 dark:text-white">Reject sub-task</h3>
              <p class="text-xs text-gray-500 truncate max-w-[280px]">{{ task?.title }}</p>
            </div>
          </div>
          <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">Explain what failed so the developer can fix it. This will be logged as a comment.</p>
          <textarea
            ref="uatRejectTextareaRef"
            v-model="uatRejectReason"
            rows="4"
            placeholder="Describe the issue (min. 10 characters)…"
            class="w-full rounded-xl border border-gray-700 bg-gray-50 dark:bg-gray-800/60 px-3 py-2.5 text-sm text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-600 focus:border-red-500/50 focus:outline-none focus:ring-1 focus:ring-red-500/30 resize-none"
          />
          <p v-if="uatRejectReason.length > 0 && uatRejectReason.length < 10" class="text-[11px] text-red-400 mt-1">
            At least {{ 10 - uatRejectReason.length }} more character(s) required
          </p>
          <div class="flex items-center justify-end gap-3 mt-4">
            <button
              type="button"
              class="px-4 py-2 rounded-lg border border-gray-700 text-xs font-medium text-gray-500 dark:text-gray-400 hover:border-gray-400 dark:hover:border-gray-600 hover:text-gray-700 dark:text-gray-200 transition-colors"
              @click="closeUATRejectModal"
            >Cancel</button>
            <button
              type="button"
              :disabled="uatRejectReason.length < 10 || uatRejectSubmitting"
              class="px-4 py-2 rounded-lg bg-red-100 dark:bg-red-500/15 border border-red-300 dark:border-red-500/40 text-red-400 text-xs font-bold hover:bg-red-100 dark:bg-red-500/25 transition-colors disabled:opacity-50"
              @click="submitUATReject"
            >
              <span v-if="uatRejectSubmitting">Rejecting…</span>
              <span v-else>Confirm reject</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- ══ OUTSOURCE MODAL ══ -->
    <OutsourceRequestModal
      v-model="showOutsourceModal"
      :prefill-title="task?.title"
      :prefill-description="task?.description"
      :prefill-minutes="task?.estimated_minutes"
      @created="onOutsourceCreated"
    />

    <!-- ══ NEW DEPLOYMENT REQUEST MODAL ══ -->
    <Teleport to="body">
      <div v-if="showCreateDeploymentModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm" @click.self="showCreateDeploymentModal = false">
        <div class="w-full max-w-lg bg-gray-900 border border-gray-200 dark:border-gray-700/60 rounded-2xl shadow-2xl shadow-black/40 overflow-hidden">
          <div class="flex items-center justify-between px-6 py-4 border-b border-gray-200 dark:border-gray-700/50 bg-gray-50 dark:bg-gray-800/60">
            <div class="flex items-center gap-2">
              <span class="text-orange-400">🚀</span>
              <h2 class="text-sm font-bold text-gray-900 dark:text-white">New Deployment Request</h2>
            </div>
            <button @click="showCreateDeploymentModal = false" class="text-gray-500 hover:text-gray-600 dark:text-gray-300 transition-colors">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
          <form @submit.prevent="submitCreateDeployment" class="px-6 py-5 space-y-4">
            <p class="text-xs text-gray-500 dark:text-gray-400">This deployment request is linked to task <span class="text-orange-300 font-mono">{{ task?.code }}</span> and will automatically advance it to <span class="text-amber-300 font-medium">Ready for UAT</span> once deployed.</p>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">Title <span class="text-red-400">*</span></label>
              <input v-model="deployForm.title" type="text" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-none focus:border-orange-500 transition-colors" :placeholder="`Deploy: ${task?.title}`" required />
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">Branch <span class="text-red-400">*</span></label>
              <input v-model="deployForm.branch" type="text" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm font-mono text-cyan-300 focus:outline-none focus:border-orange-500 transition-colors" placeholder="feature/my-branch" required />
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">Environment</label>
                <select v-model="deployForm.environment" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-none focus:border-orange-500 transition-colors">
                  <option value="STAGING">STAGING</option>
                  <option value="PRE-PROD">PRE-PROD</option>
                  <option value="PRODUCTION">PRODUCTION</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">PR URL</label>
                <input v-model="deployForm.pr_url" type="url" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-none focus:border-orange-500 transition-colors" placeholder="https://github.com/…" />
              </div>
            </div>
            <!-- Assignee selector -->
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">Assign to Chief Engineer</label>
              <select v-model="deployForm.reviewer_id" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-none focus:border-orange-500 transition-colors">
                <option value="">— Unassigned (anyone can pick up) —</option>
                <option v-for="ce in deployChiefEngineers" :key="ce.id" :value="String(ce.id)">
                  {{ ce.display_name || ce.email }}
                </option>
              </select>
              <p v-if="deployChiefEngineers.length === 0" class="text-[10px] text-gray-500 mt-1">No Chief Engineers found</p>
            </div>
            <div>
              <label class="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">Release Notes</label>
              <textarea v-model="deployForm.description" rows="3" class="w-full bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-gray-900 dark:text-white focus:outline-none focus:border-orange-500 transition-colors resize-none" placeholder="What changed, what to watch out for…" />
            </div>
            <div class="flex justify-end gap-2 pt-1">
              <button type="button" @click="showCreateDeploymentModal = false" class="px-4 py-2 rounded-lg text-sm text-gray-500 dark:text-gray-400 hover:text-white border border-gray-700 hover:border-gray-400 dark:hover:border-gray-600 transition-colors">Cancel</button>
              <button type="submit" :disabled="deployFormSubmitting || !deployForm.title || !deployForm.branch" class="flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-semibold bg-gradient-to-r from-orange-100 dark:from-orange-600 to-amber-100 dark:to-amber-600 hover:from-orange-200 dark:hover:from-orange-500 hover:to-amber-200 dark:hover:to-amber-500 text-gray-900 dark:text-white transition-all disabled:opacity-50">
                <span v-if="deployFormSubmitting" class="w-3.5 h-3.5 rounded-full border-2 border-white border-t-transparent animate-spin" />
                {{ deployFormSubmitting ? 'Creating…' : 'Create Request' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useAuthStore } from '~/core/modules/auth/store/auth-store'
import TaskComments from '~/components/tasks/TaskComments.vue'
import TimeLogger from '~/components/tasks/TimeLogger.vue'
import RichTextEditor from '~/components/editor/RichTextEditor.vue'
import SubtaskList from '~/components/tasks/SubtaskList.vue'
import { useProjectsApi } from '~/core/modules/projects/infrastructure/projects-api'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import type { TaskComment, TimeLog } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'
import OutsourceRequestModal from '~/components/b2b/OutsourceRequestModal.vue'
import { minutesToEffortHours, effortHoursToMinutes, formatMinutesAsHours } from '~/utils/effortHours'
import { canSeeCeoAssigneeOption, isTaskAssigneeRole } from '~/utils/roles'
import { useDeploymentApi } from '~/core/modules/deployment/infrastructure/deployment-api'
import { sortBacklogTasks } from '~/utils/backlog-task-utils'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

interface Appeal {
  id: string
  submission_id: string
  developer_id: number
  reason: string
  status: string // PENDING, APPROVED, REJECTED
  
  // AI Advisory System
  ai_recommendation: string // OVERTURN or UPHOLD
  ai_confidence: number     // 0-100
  ai_reasoning: string      // Advice for CEO / Product Owner
  
  resolver_id: number | null
  resolver_note: string
  created_at: string
  updated_at: string
}

interface Submission {
  id: string
  task_id: string
  dev_id: number
  commit_hash: string
  ai_verdict: string
  ai_score: number
  ai_feedback: any
  is_overridden: boolean
  appeal?: Appeal
  created_at: string
}

interface SubTask {
  id: string
  title: string
  status: string
  assigned_to: number | null
  assigned_to_display_name?: string
  assigned_to_email?: string
  assigned_to_avatar_url?: string
  estimated_minutes: number
  progress: number
}

interface Task {
  id: string
  code?: string // e.g. mims-hdmap-main-001
  project_id?: string | null // UUID of project (for Back to project when task has no code)
  epic_id?: string | null
  // Backlog ordering (used for Prev/Next when opened from backlog board)
  sort_order?: number | null
  title: string
  description: string
  resource_urls: any
  estimated_minutes: number
  task_type?: string // FEATURE, TASK, BUG

  // Parent-child hierarchy
  parent_id?: string | null
  parent_task?: {
    id: string
    title: string
    code?: string
    task_type?: string
    status?: string
    parent_id?: string | null  // set when parent is itself a sub-task
  } | null
  sub_tasks?: SubTask[]

  // Time Negotiation
  negotiation_status: string // NONE, PENDING, APPROVED, REJECTED
  proposed_minutes: number
  negotiation_reason: string
  
  due_at: string | null
  start_date: string | null
  end_date: string | null
  priority?: string
  story_points?: number
  sprint_id?: string | null
  started_at: string | null
  completed_at: string | null
  status: string
  assigned_to: number | null
  assigned_to_display_name?: string
  assigned_to_email?: string
  assigned_to_avatar_url?: string
  created_by: number
  created_by_role?: string
  created_by_email?: string
  created_by_display_name?: string
  is_komgrip?: boolean
  created_at: string
  updated_at: string
  submissions?: Submission[]
}

interface TaskActivityItem {
  id: string
  action: string
  at: string
  actor_user_id?: number
  actor_email?: string
  actor_display_name?: string
  payload?: Record<string, unknown> | null
  inferred?: boolean
}

const route = useRoute()
const { fetchWithAuth, currentUser: authCurrentUser } = useAuth()
const authStore = useAuthStore()
const projectsApi = useProjectsApi()
const tasksApi = useTasksApi()
const { getTeams } = useTeamsApi()
const teamsStore = useTeamsStore()
const { showSuccess, showError } = useNotification()

// Current user profile (for comment avatar)
const currentUserAvatarURL = ref('')
const currentUserInitial = computed(() => {
  const u = authCurrentUser.value
  if (!u) return ''
  return (u.display_name || u.email || '').charAt(0).toUpperCase()
})

onMounted(async () => {
  try {
    const { authApi } = await import('~/core/modules/auth/infrastructure/auth-api')
    const me = await authApi.getMe()
    currentUserAvatarURL.value = me.avatar_url || ''
  } catch { /* silent */ }
})

// State
const task = ref<Task | null>(null)
const taskSummary = ref<Task | null>(null)
const taskDetailLoaded = ref(false)
const richDetailLoading = ref(false)
/** True when the summary response indicates there is a rich description / attachments to lazy-load */
const hasRichContent = ref(false)
const isLoading = ref(true)
const error = ref('')

// Sprint context: task IDs in order (for Prev/Next when from_sprint + from_project)
const sprintTaskIds = ref<string[]>([])
// Backlog context: task IDs in display order (for Prev/Next when from_project + from_tab=backlog)
const backlogTaskIds = ref<string[]>([])
const sprintNavMessage = ref('')

// Assignee change state
const showAssignDropdown = ref(false)
const assigneeUsers = ref<{ id: number; email: string; display_name: string; role: string }[]>([])
const assignSelectedId = ref<number | ''>('')
const assignLoading = ref(false)
const assignError = ref('')

const visibleAssigneeUsers = computed(() => {
  const viewerCanSeeCeo = canSeeCeoAssigneeOption(currentRole.value)
  return assigneeUsers.value.filter((u) => viewerCanSeeCeo || u.role?.toUpperCase() !== 'CEO')
})

const canLoadCeoAssignee = computed(() => canSeeCeoAssigneeOption(currentRole.value))

// Continuous UAT: approve / reject on this page
const uatActionLoading = ref(false)
const uatActionType = ref<'approve' | 'reject' | null>(null)
const uatApproveConfirmOpen = ref(false)
const uatApproveSubmitting = ref(false)
const uatRejectModalOpen = ref(false)
const uatRejectReason = ref('')
const uatRejectSubmitting = ref(false)
const uatRejectTextareaRef = ref<HTMLTextAreaElement | null>(null)
// Product Owner test evidence form
const uatTestUrl = ref('')
const uatTestSteps = ref('')
const uatTestUrlRef = ref<HTMLInputElement | null>(null)

// Outsource Modal State
const showOutsourceModal = ref(false)
const onOutsourceCreated = () => {
  showSuccess('Outsource request sent to the target team!', 'B2B Request Sent')
}

// Edit Task State
const showEditModal = ref(false)
const editForm = ref({
  title: '',
  description: '',
  task_type: 'TASK',
  deadline: '',
  priority: 'MEDIUM',
  story_points: 0,
  sprint_id: '',
  start_date: '',
  end_date: '',
  estimated_hours: 0
})
const editSprints = ref<{ id: string; name: string }[]>([])
const isUpdatingTask = ref(false)
const editError = ref('')

// Delete Task State
const showDeleteModal = ref(false)
const isDeletingTask = ref(false)
const deleteError = ref('')

// Inline Description Edit State
const isEditingDescription = ref(false)
const isSavingDescription = ref(false)
const inlineDescriptionHtml = ref('')

// Comments & Time Logs State
const comments = ref<TaskComment[]>([])
const timeLogs = ref<TimeLog[]>([])
const commentsLoading = ref(false)
const timeLogsLoading = ref(false)

const taskActivity = ref<TaskActivityItem[]>([])
const activityLoading = ref(false)
const activityError = ref('')

// Estimated Effort (UI: hours, 1 decimal → API: integer minutes)
const estimatedHoursLocal = ref(0)
const isSavingEstimate = ref(false)
const estimatedMinutesDirty = computed(() =>
  task.value != null &&
  effortHoursToMinutes(Number(estimatedHoursLocal.value)) !== (task.value.estimated_minutes ?? 0)
)

// Sub-tasks (child tasks for parent-child hierarchy)
const subtasks = ref<SubTask[]>([])

/** True when this task is a Parent (has sub-tasks) */
const isParentTask = computed(() => subtasks.value.length > 0)

/** Roll-up: sum of all sub-task estimated_minutes */
const subtaskTotalEstimatedMinutes = computed(() =>
  subtasks.value.reduce((sum, s) => sum + (s.estimated_minutes || 0), 0)
)

/** Roll-up: aggregate progress from sub-tasks */
const subtaskAggregateProgress = computed(() => {
  if (subtasks.value.length === 0) return task.value?.progress ?? 0
  const total = subtasks.value.reduce((sum, s) => {
    if (s.status === 'COMPLETED') return sum + 100
    return sum + (s.progress || 0)
  }, 0)
  return Math.round(total / subtasks.value.length)
})

// Parsed slide resource URLs (for Google Slides imported tasks; images are now in description, only metadata for "Open in Slides")
const slideResourceURLs = computed(() => {
  if (!task.value?.resource_urls) return null
  const ru = typeof task.value.resource_urls === 'string'
    ? JSON.parse(task.value.resource_urls)
    : task.value.resource_urls
  if (ru?.source !== 'google_slides') return null
  return ru as {
    thumbnail_url: string
    images: string[]
    slide_url: string
    source: string
    slide_index: number
    presentation_id: string
    comments: Array<{ content: string; author: string; resolved: boolean }>
  }
})

// Open-in-Slides URL: use stored URL if it has #slide= fragment, else build with slide_index so we open to the right slide
const slideOpenInSlidesURL = computed(() => {
  const ru = slideResourceURLs.value
  if (!ru?.slide_url) return ''
  if (ru.slide_url.includes('#slide=')) return ru.slide_url
  if (ru.presentation_id && ru.slide_index != null)
    return `https://docs.google.com/presentation/d/${ru.presentation_id}/edit#slide=${ru.slide_index}`
  return ru.slide_url
})

// Single source: store first, then JWT payload (so Edit works right after login)
const effectiveUser = computed(() => {
  if (authStore.user) return authStore.user
  const payload = authCurrentUser.value
  if (!payload) return null
  const id = payload.user_id ?? (payload as any).userId
  return id != null || payload.role ? { id: Number(id) || 0, role: payload.role || '', email: payload.email || '' } : null
})

const canEditOrDelete = computed(() => {
  if (!task.value || !effectiveUser.value) return false
  const user = effectiveUser.value
  const role = (user.role || '').trim().toUpperCase()
  if (role === 'CEO' || role === 'PRODUCT_OWNER' || role === 'PM' || role === 'MANAGER') return true
  const creatorId = Number(task.value.created_by)
  const userId = Number(user.id ?? authStore.userId ?? 0)
  return creatorId === userId && !Number.isNaN(userId)
})

/** Sub-tasks: Product Owner / CEO / creator (canEditOrDelete) or the user assigned to this task */
const isCurrentUserAssignee = computed(() => {
  if (!task.value || !effectiveUser.value) return false
  const aid = task.value.assigned_to
  if (aid == null || aid === undefined) return false
  const userId = Number(effectiveUser.value.id ?? authStore.userId ?? 0)
  return Number(aid) === userId && !Number.isNaN(userId)
})

const canManageSubtasks = computed(() => canEditOrDelete.value || isCurrentUserAssignee.value)
const canEditDescription = computed(() => canEditOrDelete.value || isCurrentUserAssignee.value)

const currentRole = computed(() => (effectiveUser.value?.role || '').trim().toUpperCase())
const ceoQuickFinishSubmitting = ref(false)

const canClaimTask = computed(() => {
  if (!task.value || !effectiveUser.value) return false
  if (task.value.assigned_to != null) return false
  const role = currentRole.value
  return role === 'ENGINEER' || role === 'CHIEF_ENGINEER' || role === 'CHIEF'
})

const statusChangeLoading = ref(false)
const pendingStatusChange = ref('')

interface StatusTransitionAction {
  status: string
  label: string
  icon: string
  className: string
  roles?: string[]
  states?: string[]
}

const STATUS_TRANSITIONS: StatusTransitionAction[] = [
  {
    status: 'IN_PROGRESS',
    label: 'Start',
    icon: '▶',
    className: 'bg-blue-50 border-blue-200 text-blue-700 hover:bg-blue-100 hover:border-blue-300 dark:bg-blue-900/30 dark:border-blue-600/50 dark:text-blue-300 dark:hover:bg-blue-900/50 dark:hover:border-blue-500/60',
    roles: ['ENGINEER', 'CHIEF_ENGINEER', 'CHIEF', 'MANAGER', 'CEO', 'PRODUCT_OWNER', 'PM'],
    states: ['PENDING', 'ASSIGNED', 'BLOCKED'],
  },
  {
    status: 'READY_FOR_TEST',
    label: 'Ready to test',
    icon: '🧪',
    className: 'bg-cyan-50 border-cyan-200 text-cyan-700 hover:bg-cyan-100 hover:border-cyan-300 dark:bg-cyan-900/30 dark:border-cyan-600/50 dark:text-cyan-300 dark:hover:bg-cyan-900/50 dark:hover:border-cyan-500/60',
    roles: ['ENGINEER', 'CHIEF_ENGINEER', 'CHIEF', 'MANAGER', 'CEO', 'PRODUCT_OWNER', 'PM'],
    states: ['IN_PROGRESS'],
  },
  {
    status: 'COMPLETED',
    label: 'Done',
    icon: '✅',
    className: 'bg-emerald-50 border-emerald-200 text-emerald-700 hover:bg-emerald-100 hover:border-emerald-300 dark:bg-emerald-900/30 dark:border-emerald-600/50 dark:text-emerald-300 dark:hover:bg-emerald-900/50 dark:hover:border-emerald-500/60',
    roles: ['CEO', 'MANAGER'],
    states: ['READY_FOR_UAT'],
  },
]

const availableStatusTransitions = computed(() => {
  if (!task.value) return []
  const role = currentRole.value
  return STATUS_TRANSITIONS.filter((action) => {
    const roleAllowed = !action.roles || action.roles.includes(role)
    const stateAllowed = !action.states || action.states.includes(task.value!.status)
    return roleAllowed && stateAllowed && action.status !== task.value!.status
  })
})

async function changeTaskStatus(status: string) {
  if (!task.value || statusChangeLoading.value) return
  statusChangeLoading.value = true
  pendingStatusChange.value = status
  try {
    await tasksApi.updateTask(task.value.id, { status })
    showSuccess(`Task status updated to ${getStatusLabel(status)}`, 'Updated')
    await fetchTask()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to update task status')
  } finally {
    statusChangeLoading.value = false
    pendingStatusChange.value = ''
  }
}

const showCEOQuickFinishAction = computed(() => {
  if (!task.value) return false
  const role = currentRole.value
  if (role !== 'CEO' && role !== 'MANAGER') return false
  return task.value.status !== 'COMPLETED'
})

/** Product Owner / MANAGER step: task is READY_FOR_TEST and viewer is Product Owner or MANAGER */
const showPMUATActions = computed(() => {
  if (!task.value) return false
  const role = currentRole.value
  if (role !== 'PRODUCT_OWNER' && role !== 'PM' && role !== 'MANAGER') return false
  if (task.value.status !== 'READY_FOR_TEST') return false
  const t = task.value.task_type
  return t === 'TASK' || t === 'BUG'
})

/** WAIT_FOR_DEPLOY section: task is waiting for Chief Engineer to deploy */
const showWaitForDeploySection = computed(() => {
  if (!task.value) return false
  return task.value.status === 'WAIT_FOR_DEPLOY' && (task.value.task_type === 'TASK' || task.value.task_type === 'BUG')
})

/** CEO/MANAGER step: task is READY_FOR_UAT (Chief Engineer deployed) and viewer is CEO or MANAGER */
const showCEOUATActions = computed(() => {
  if (!task.value) return false
  const role = currentRole.value
  if (role !== 'CEO' && role !== 'MANAGER') return false
  if (task.value.status !== 'READY_FOR_UAT') return false
  const t = task.value.task_type
  return t === 'TASK' || t === 'BUG'
})

// Deployment request linked to this task (fetched when status is WAIT_FOR_DEPLOY or READY_FOR_UAT)
const deploymentForTask = ref<import('~/core/modules/deployment/infrastructure/deployment-api').DeploymentRequest | null>(null)
const deploymentLoading = ref(false)
const showCreateDeploymentModal = ref(false)
const deployFormSubmitting = ref(false)
const deployForm = reactive({
  title: '',
  branch: '',
  environment: 'STAGING' as 'STAGING' | 'PRE-PROD' | 'PRODUCTION',
  pr_url: '',
  description: '',
  reviewer_id: '' as string,
})
const deployChiefEngineers = ref<import('~/core/modules/deployment/infrastructure/deployment-api').DeploymentUser[]>([])

async function openDeploymentModal() {
  showCreateDeploymentModal.value = true
  if (deployChiefEngineers.value.length === 0) {
    const depApi = useDeploymentApi()
    deployChiefEngineers.value = await depApi.fetchChiefEngineers()
  }
}

async function fetchDeploymentForTask() {
  if (!task.value) return
  if (!['WAIT_FOR_DEPLOY', 'READY_FOR_UAT'].includes(task.value.status)) return
  deploymentLoading.value = true
  try {
    const depApi = useDeploymentApi()
    deploymentForTask.value = await depApi.getByTaskId(task.value.id)
  } catch { /* not found is fine */ } finally {
    deploymentLoading.value = false
  }
}

async function submitCreateDeployment() {
  if (!task.value || !deployForm.title || !deployForm.branch) return
  deployFormSubmitting.value = true
  try {
    const depApi = useDeploymentApi()
    deploymentForTask.value = await depApi.createRequest({
      title: deployForm.title,
      branch: deployForm.branch,
      environment: deployForm.environment,
      pr_url: deployForm.pr_url || undefined,
      description: deployForm.description || undefined,
      task_id: task.value.id,
      task_ref: task.value.code || task.value.title,
      reviewer_id: deployForm.reviewer_id ? Number(deployForm.reviewer_id) : undefined,
    })
    showCreateDeploymentModal.value = false
    deployForm.title = ''
    deployForm.branch = ''
    deployForm.environment = 'STAGING'
    deployForm.pr_url = ''
    deployForm.description = ''
    deployForm.reviewer_id = ''
    showSuccess('Deployment request created. Chief Engineer will review it.', 'Created 🚀')
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to create deployment request')
  } finally {
    deployFormSubmitting.value = false
  }
}

/** Parsed UAT payload (test URL + steps stored by Product Owner when they submitted) */
const uatPayloadData = computed<{ test_url?: string; test_steps?: string } | null>(() => {
  const raw = task.value?.uat_payload
  if (!raw) return null
  try {
    const parsed = typeof raw === 'string' ? JSON.parse(raw) : raw
    return parsed as { test_url?: string; test_steps?: string }
  } catch {
    return null
  }
})

/** Validation for Product Owner test evidence form */
const isUATApproveFormValid = computed(() =>
  uatTestUrl.value.startsWith('http') && uatTestSteps.value.trim().length >= 20
)

const creatorLabel = computed(() => {
  if (!task.value) return ''
  const role = task.value.created_by_role
  const email = task.value.created_by_email
  const displayName = task.value.created_by_display_name
  if (role === 'CEO') return displayName || email || 'CEO'
  if (role === 'PRODUCT_OWNER' || role === 'PM') return displayName || email || 'Product Owner'
  return displayName || email || `Dev #${task.value.created_by}`
})

// Methods
const fetchTask = async () => {
  const taskId = (route.params.id as string)?.trim?.() || ''

  if (!taskId || taskId === 'undefined' || taskId === 'null') {
    error.value = 'Invalid or missing task ID. Check the URL.'
    isLoading.value = false
    return
  }

  try {
    isLoading.value = true
    error.value = ''
    taskDetailLoaded.value = false
    richDetailLoading.value = false
    hasRichContent.value = false

    const response = await tasksApi.getTaskSummary(taskId)
    taskSummary.value = response.summary as unknown as Task
    task.value = response.summary as unknown as Task
    hasRichContent.value = response.has_rich_content
    estimatedHoursLocal.value = minutesToEffortHours(task.value.estimated_minutes ?? 0)
    subtasks.value = (task.value.sub_tasks ?? []) as SubTask[]
    void fetchTaskActivity()
    // Always hydrate full detail: populates description, sub_tasks, parent_task, submissions
    void fetchRichTaskDetail(taskId)
    sprintTaskIds.value = []
    backlogTaskIds.value = []
  } catch (err: any) {
    console.error('Failed to fetch task:', err)
    const apiMsg = err?.data?.message ?? err?.data?.error
    const status = err?.statusCode ?? err?.status
    error.value = apiMsg || (status === 404 ? 'Task not found.' : err?.message || 'Failed to load task.')
  } finally {
    isLoading.value = false
  }
}

const fetchRichTaskDetail = async (taskId: string) => {
  if (!taskId) return
  richDetailLoading.value = true
  try {
    const response = await tasksApi.getTaskDetail(taskId)
    task.value = response.task
    taskDetailLoaded.value = true
    subtasks.value = (response.task.sub_tasks ?? []) as SubTask[]
  } catch (err) {
    console.error('Failed to fetch rich task detail:', err)
  } finally {
    richDetailLoading.value = false
  }
}

const fetchTaskActivity = async () => {
  const id = task.value?.id || (route.params.id as string)?.trim?.() || ''
  if (!id || id === 'undefined' || id === 'null') {
    taskActivity.value = []
    return
  }
  activityLoading.value = true
  activityError.value = ''
  try {
    const response = await fetchWithAuth<{ data: TaskActivityItem[] }>(`/sentinel/tasks/${id}/activity`)
    taskActivity.value = Array.isArray(response.data) ? response.data : []
  } catch (err: any) {
    console.error('Failed to fetch task activity:', err)
    const apiMsg = err?.data?.message ?? err?.data?.error
    activityError.value =
      apiMsg || (err?.statusCode === 404 || err?.status === 404 ? 'Task not found.' : err?.message || 'Failed to load activity.')
    taskActivity.value = []
  } finally {
    activityLoading.value = false
  }
}

function activityPayload(item: TaskActivityItem): Record<string, unknown> {
  const p = item.payload
  if (p != null && typeof p === 'object' && !Array.isArray(p)) {
    return p as Record<string, unknown>
  }
  return {}
}

/** Legacy Thai notes from earlier API responses — show English in UI */
const legacyActivityNoteEn: Record<string, string> = {
  'สร้างก่อนระบบบันทึกกิจกรรม — เวลาจากเรคคอร์ด task':
    'Created before activity logging — timestamp from task record',
  'เวลาเริ่มงานจาก started_at — ผู้เริ่มไม่ได้ถูกบันทึกในประวัติเดิม':
    'Start time from started_at — actor was not recorded in legacy history',
  'เวลาปิดงานจาก completed_at — รายละเอียดการอนุมัติอาจอยู่ก่อนระบบบันทึกกิจกรรม':
    'Completion time from completed_at — approval trail may predate activity logging',
}

function displayActivityNote(note: string): string {
  const t = note.trim()
  return legacyActivityNoteEn[t] ?? note
}

function activityActorLabel(item: TaskActivityItem): string {
  const name = (item.actor_display_name || '').trim()
  const email = (item.actor_email || '').trim()
  if (name) return name
  if (email) return email
  if (item.actor_user_id != null && item.actor_user_id > 0) return `User #${item.actor_user_id}`
  return 'Unknown actor'
}

function activityStatusLabel(s: string): string {
  const labels: Record<string, string> = {
    PENDING: 'Pending',
    IN_PROGRESS: 'In progress',
    READY_FOR_TEST: 'Ready for test',
    WAIT_FOR_DEPLOY: 'Wait for deploy',
    READY_FOR_UAT: 'Ready for UAT',
    REVIEW_PENDING: 'Review pending',
    COMPLETED: 'Completed',
    BLOCKED: 'Blocked',
    CANCELLED: 'Cancelled',
  }
  return labels[s] || s
}

function activityTitle(item: TaskActivityItem): string {
  switch (item.action) {
    case 'CREATED':
      return 'Task created'
    case 'ASSIGNED':
      return 'Assignee set'
    case 'UNASSIGNED':
      return 'Unassigned'
    case 'STATUS_CHANGED':
      return 'Status changed'
    case 'SUBMITTED_REVIEW':
      return 'Submitted for review (handover)'
    case 'APPROVED_REVIEW':
      return 'Approved after review'
    case 'REJECTED_REVIEW':
      return 'Rejected — sent back for rework'
    case 'READY_FOR_TEST':
      return 'Marked ready for test'
    case 'PM_APPROVED_TEST':
      return 'PO approved test → wait for deploy'
    case 'DEPLOYED':
      return 'Deployed → ready for UAT'
    case 'CEO_FINAL_APPROVED':
      return 'CEO / Manager final approval'
    case 'WORKFLOW_REJECT':
      return 'Workflow rejection (sent back)'
    case 'SUBMIT_UAT':
      return 'Feature UAT submitted'
    case 'NEGOTIATION_SUBMITTED':
      return 'Time negotiation submitted'
    case 'APPEAL_APPROVED_COMPLETE':
      return 'Appeal approved — task completed'
    case 'PARENT_ROLLUP_STATUS':
      return 'Parent feature status updated (roll-up)'
    default:
      return item.action.replaceAll('_', ' ')
  }
}

function activityIconGlyph(action: string): string {
  const m: Record<string, string> = {
    CREATED: '＋',
    ASSIGNED: '👤',
    UNASSIGNED: '◌',
    STATUS_CHANGED: '⇄',
    SUBMITTED_REVIEW: '⎇',
    APPROVED_REVIEW: '✓',
    REJECTED_REVIEW: '↩',
    READY_FOR_TEST: '🧪',
    PM_APPROVED_TEST: '✅',
    DEPLOYED: '🚀',
    CEO_FINAL_APPROVED: '👑',
    WORKFLOW_REJECT: '✗',
    SUBMIT_UAT: '📦',
    NEGOTIATION_SUBMITTED: '⏱',
    APPEAL_APPROVED_COMPLETE: '⚖',
    PARENT_ROLLUP_STATUS: '🔗',
  }
  return m[action] || '•'
}

function activityDotClass(action: string): string {
  const base = 'border-gray-600/80 bg-gray-900/90 text-gray-200'
  const accent: Record<string, string> = {
    CREATED: 'border-emerald-500/40 bg-emerald-500/15 text-emerald-300',
    ASSIGNED: 'border-blue-500/40 bg-blue-500/15 text-blue-300',
    UNASSIGNED: 'border-gray-500/50 bg-gray-800/80 text-gray-400',
    STATUS_CHANGED: 'border-violet-500/40 bg-violet-500/15 text-violet-300',
    SUBMITTED_REVIEW: 'border-cyan-500/40 bg-cyan-500/15 text-cyan-300',
    APPROVED_REVIEW: 'border-green-500/40 bg-green-500/15 text-green-300',
    REJECTED_REVIEW: 'border-red-500/40 bg-red-500/15 text-red-300',
    READY_FOR_TEST: 'border-cyan-500/40 bg-cyan-500/15 text-cyan-300',
    PM_APPROVED_TEST: 'border-teal-500/40 bg-teal-500/15 text-teal-300',
    DEPLOYED: 'border-orange-500/40 bg-orange-500/15 text-orange-300',
    CEO_FINAL_APPROVED: 'border-amber-500/40 bg-amber-500/15 text-amber-300',
    WORKFLOW_REJECT: 'border-red-500/40 bg-red-500/15 text-red-300',
    SUBMIT_UAT: 'border-purple-500/40 bg-purple-500/15 text-purple-300',
    NEGOTIATION_SUBMITTED: 'border-yellow-500/40 bg-yellow-500/15 text-yellow-300',
    APPEAL_APPROVED_COMPLETE: 'border-lime-500/40 bg-lime-500/15 text-lime-300',
    PARENT_ROLLUP_STATUS: 'border-indigo-500/40 bg-indigo-500/15 text-indigo-300',
  }
  return accent[action] || base
}

function activityDetailLine(item: TaskActivityItem): string {
  const actor = activityActorLabel(item)
  const p = activityPayload(item)
  const parts: string[] = []

  const rawNote = typeof p.note === 'string' ? p.note.trim() : ''
  const fromS = typeof p.from_status === 'string' ? p.from_status : ''
  const toS = typeof p.to_status === 'string' ? p.to_status : ''
  const src = typeof p.source === 'string' ? p.source : ''

  switch (item.action) {
    case 'CREATED':
      parts.push(`By ${actor}`)
      if (typeof p.title === 'string' && p.title) parts.push(`Title: ${p.title}`)
      break
    case 'ASSIGNED': {
      const assignee =
        (typeof p.assignee_display_name === 'string' && p.assignee_display_name) ||
        (typeof p.assignee_user_id === 'number' && p.assignee_user_id > 0 ? `User #${p.assignee_user_id}` : '—')
      parts.push(`Assigned to ${assignee} by ${actor}`)
      break
    }
    case 'UNASSIGNED': {
      const prev =
        (typeof p.previous_assignee_display_name === 'string' && p.previous_assignee_display_name) ||
        (typeof p.previous_assignee_user_id === 'number' && p.previous_assignee_user_id > 0
          ? `User #${p.previous_assignee_user_id}`
          : 'previous assignee')
      parts.push(`Unassigned ${prev} by ${actor}`)
      if (toS) parts.push(`Status → ${activityStatusLabel(toS)}`)
      break
    }
    case 'STATUS_CHANGED':
      if (fromS && toS) parts.push(`${activityStatusLabel(fromS)} → ${activityStatusLabel(toS)} by ${actor}`)
      else if (toS) parts.push(`→ ${activityStatusLabel(toS)} by ${actor}`)
      else parts.push(`By ${actor}`)
      if (src === 'kanban_bulk') parts.push('via Kanban')
      break
    case 'SUBMITTED_REVIEW': {
      const url = typeof p.reference_url === 'string' ? p.reference_url : ''
      parts.push(`By ${actor}`)
      if (url) parts.push(`Reference: ${url}`)
      break
    }
    case 'PM_APPROVED_TEST': {
      const u = typeof p.test_url === 'string' ? p.test_url : ''
      parts.push(`By ${actor}`)
      if (u) parts.push(`Test URL: ${u}`)
      break
    }
    case 'DEPLOYED':
      parts.push(`Chief Engineer / deployer: ${actor}`)
      if (fromS && toS) parts.push(`${activityStatusLabel(fromS)} → ${activityStatusLabel(toS)}`)
      break
    case 'NEGOTIATION_SUBMITTED':
      if (typeof p.proposed_minutes === 'number') parts.push(`Proposed ${p.proposed_minutes} minutes by ${actor}`)
      else parts.push(`By ${actor}`)
      break
    case 'WORKFLOW_REJECT': {
      const ex = typeof p.reason_excerpt === 'string' ? p.reason_excerpt : ''
      parts.push(`By ${actor}`)
      if (fromS && toS) parts.push(`${activityStatusLabel(fromS)} → ${activityStatusLabel(toS)}`)
      if (ex) parts.push(ex)
      break
    }
    case 'PARENT_ROLLUP_STATUS': {
      const child = typeof p.child_task_id === 'string' ? p.child_task_id : ''
      parts.push(`All sub-tasks completed — by ${actor}`)
      if (child) parts.push(`Child task: ${child.slice(0, 8)}…`)
      if (toS) parts.push(`Parent status → ${activityStatusLabel(toS)}`)
      break
    }
    case 'SUBMIT_UAT': {
      const st = typeof p.staging_url === 'string' ? p.staging_url : ''
      parts.push(`By ${actor}`)
      if (st) parts.push(`Staging: ${st}`)
      break
    }
    case 'APPROVED_REVIEW':
    case 'CEO_FINAL_APPROVED':
    case 'READY_FOR_TEST':
    case 'REJECTED_REVIEW':
    case 'APPEAL_APPROVED_COMPLETE':
      if (fromS && toS) parts.push(`${activityStatusLabel(fromS)} → ${activityStatusLabel(toS)} · ${actor}`)
      else parts.push(`By ${actor}`)
      break
    default:
      parts.push(`By ${actor}`)
  }

  if (rawNote) parts.push(displayActivityNote(rawNote))

  return parts.filter(Boolean).join(' · ')
}

const formatActivityDateTime = (iso: string) => {
  const d = new Date(iso)
  if (isNaN(d.getTime())) return iso
  return d.toLocaleString('en-US', { dateStyle: 'medium', timeStyle: 'short' })
}

// Load Prev/Next ordering (sprint/backlog context) without blocking initial task render.
const fetchPrevNextOrdering = async () => {
  const fromSprint = route.query.from_sprint as string | undefined
  const fromProject = route.query.from_project as string | undefined
  const fromTab = route.query.from_tab as string | undefined

  const projectId = task.value?.project_id
  const isSprintContext = !!(fromSprint && fromProject)
  const isBacklogContext = !!(fromProject && fromTab === 'backlog')

  if (!projectId) return

  try {
    if (isSprintContext && fromSprint) {
      const projectTasks = await tasksApi.getTasksByProject(projectId)
      const inSprint = projectTasks.filter((t) => t.sprint_id === fromSprint)
      inSprint.sort((a: any, b: any) => (a.code || a.id).localeCompare(b.code || b.id))
      sprintTaskIds.value = inSprint.map((t) => t.id)
      backlogTaskIds.value = []
      return
    }

    if (isBacklogContext) {
      const [allProjectTasks, epics, sprints] = await Promise.all([
        tasksApi.getTasksByProject(projectId),
        projectsApi.getEpics(projectId),
        projectsApi.getSprints(projectId),
      ])

      // Match projects backlog page: sprint order uses sort_order (fallback created_at)
      const sprintOrderIds = [...sprints]
        .sort((a: any, b: any) => (a.sort_order ?? 0) - (b.sort_order ?? 0) || new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
        .map((s: any) => s.id)

      // Keep epic order as returned in `epics`
      const orderedEpics = epics

      // Prev/Next should move between top-level backlog rows (A items), excluding sub-tasks.
      const flatRoots: any[] = []
      for (const epic of orderedEpics) {
        const roots = sortBacklogTasks(
          (allProjectTasks as any[]).filter((t) => t.epic_id === epic.id && !t.parent_id),
          sprintOrderIds,
        )
        roots.forEach((r) => flatRoots.push(r))
      }

      const unassignedRoots = sortBacklogTasks(
        (allProjectTasks as any[]).filter((t) => !t.epic_id && !t.parent_id),
        sprintOrderIds,
      )
      unassignedRoots.forEach((r) => flatRoots.push(r))

      backlogTaskIds.value = flatRoots.map((t) => t.id)
      sprintTaskIds.value = []
    }
  } catch (e) {
    sprintTaskIds.value = []
    backlogTaskIds.value = []
  }
}

async function openAssignDropdown() {
  showAssignDropdown.value = true
  assignError.value = ''
  assignSelectedId.value = ''
  if (assigneeUsers.value.length === 0) {
    try {
      const role = (authCurrentUser.value?.role || '').toUpperCase()
      if (role === 'PRODUCT_OWNER' || role === 'PM') {
        await teamsStore.fetchTeamsFeatureEnabled()
        if (teamsStore.teamsFeatureEnabled) {
          const userId = authCurrentUser.value?.user_id
          const teams = await getTeams()
          const myTeam = teams.find(t => t.users?.some(u => u.id === userId))
          assigneeUsers.value = myTeam?.users?.filter(u => isTaskAssigneeRole(u.role)) ?? []
        } else {
          const res = await fetchWithAuth<{ data: { id: number; email: string; display_name: string; role: string }[] }>('/auth/users')
          const users = res.data ?? []
          assigneeUsers.value = users.filter((u) => isTaskAssigneeRole(u.role))
        }
      } else {
        const res = await fetchWithAuth<{ data: { id: number; email: string; display_name: string; role: string }[] }>('/auth/users')
        const users = res.data ?? []
        assigneeUsers.value = users.filter((u) => isTaskAssigneeRole(u.role))
      }
      if (!canLoadCeoAssignee.value) {
        assigneeUsers.value = assigneeUsers.value.filter((u) => u.role?.toUpperCase() !== 'CEO')
      }
    } catch {
      assignError.value = 'Failed to load users'
    }
  }
}

async function confirmChangeAssignee() {
  const devId = assignSelectedId.value
  if (devId === '' || !task.value) return
  assignLoading.value = true
  assignError.value = ''
  try {
    const taskId = (route.params.id as string)?.trim?.() || task.value.id
    await tasksApi.assignTask(taskId, Number(devId))
    showAssignDropdown.value = false
    assignSelectedId.value = ''
    showSuccess(devId === '0' ? 'Assignee removed.' : 'Assignee updated.', 'Done')
    await fetchTask()
  } catch (err: any) {
    assignError.value = err?.data?.message ?? err?.message ?? 'Failed to assign'
  } finally {
    assignLoading.value = false
  }
}

async function claimTask() {
  if (!task.value || !effectiveUser.value || !canClaimTask.value) return
  assignLoading.value = true
  assignError.value = ''
  try {
    const taskId = (route.params.id as string)?.trim?.() || task.value.id
    const myId = Number(effectiveUser.value.id)
    await tasksApi.assignTask(taskId, myId)
    showSuccess('Task claimed successfully.', 'Done')
    await fetchTask()
  } catch (err: any) {
    assignError.value = err?.data?.message ?? err?.message ?? 'Failed to claim task'
    showError(assignError.value)
  } finally {
    assignLoading.value = false
  }
}

async function markTaskFinishedByCEO() {
  if (!task.value || ceoQuickFinishSubmitting.value || !showCEOQuickFinishAction.value) return

  ceoQuickFinishSubmitting.value = true
  try {
    await tasksApi.updateTask(task.value.id, { status: 'COMPLETED' })
    showSuccess('Task marked as done.', 'Done ✅')
    await fetchTask()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to mark task as done')
  } finally {
    ceoQuickFinishSubmitting.value = false
  }
}

function openUATApproveConfirm() {
  if (!task.value || uatApproveSubmitting.value) return
  if (uatRejectModalOpen.value) closeUATRejectModal()
  uatTestUrl.value = ''
  uatTestSteps.value = ''
  uatApproveConfirmOpen.value = true
  if (showPMUATActions.value) {
    nextTick(() => uatTestUrlRef.value?.focus())
  }
}

function closeUATApproveConfirm() {
  if (uatApproveSubmitting.value) return
  uatApproveConfirmOpen.value = false
  uatTestUrl.value = ''
  uatTestSteps.value = ''
}

async function submitUATApprove() {
  if (!task.value) return
  uatApproveSubmitting.value = true
  try {
    if (showPMUATActions.value) {
      // Product Owner: submit test evidence → WAIT_FOR_DEPLOY
      await tasksApi.pmApproveSubTask(task.value.id, uatTestUrl.value.trim(), uatTestSteps.value.trim())
      uatApproveConfirmOpen.value = false
      showSuccess('Test approved. Task is now waiting for Chief Engineer deployment.', 'Submitted')
    } else {
      // CEO: final approval → COMPLETED
      await tasksApi.approveSubTask(task.value.id)
      uatApproveConfirmOpen.value = false
      showSuccess('Task approved and marked as completed.', 'Done ✅')
    }
    await fetchTask()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to approve')
  } finally {
    uatApproveSubmitting.value = false
  }
}

function openUATRejectModal() {
  if (uatApproveConfirmOpen.value) closeUATApproveConfirm()
  uatRejectReason.value = ''
  uatRejectModalOpen.value = true
  nextTick(() => uatRejectTextareaRef.value?.focus())
}



function closeUATRejectModal() {
  uatRejectModalOpen.value = false
  uatRejectReason.value = ''
}

async function submitUATReject() {
  if (!task.value || uatRejectReason.value.length < 10) return
  uatRejectSubmitting.value = true
  uatActionLoading.value = true
  uatActionType.value = 'reject'
  try {
    await tasksApi.rejectSubTask(task.value.id, uatRejectReason.value)
    closeUATRejectModal()
    showSuccess('Sub-task rejected and returned to in progress.', 'Done')
    await fetchTask()
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to reject')
  } finally {
    uatRejectSubmitting.value = false
    uatActionLoading.value = false
    uatActionType.value = null
  }
}

const saveEstimatedMinutes = async () => {
  if (!task.value) return
  const mins = effortHoursToMinutes(Number(estimatedHoursLocal.value))
  if (mins < 0 || Number.isNaN(mins)) return
  const taskId = (route.params.id as string)?.trim?.() || task.value.id
  if (!taskId) return
  try {
    isSavingEstimate.value = true
    const updated = await tasksApi.updateTask(taskId, { estimated_minutes: mins })
    task.value = { ...task.value, ...updated }
    estimatedHoursLocal.value = minutesToEffortHours(updated.estimated_minutes ?? 0)
    showSuccess('Estimated effort updated.', 'Done')
  } catch (err: any) {
    showError(err?.data?.message ?? err?.message ?? 'Failed to update estimate')
  } finally {
    isSavingEstimate.value = false
  }
}

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    'COMPLETED':       'bg-green-900/40 border-green-600/50 text-green-300',
    'IN_PROGRESS':     'bg-blue-900/40 border-blue-600/50 text-blue-300',
    'PENDING':         'bg-yellow-900/40 border-yellow-600/50 text-yellow-300',
    'REVIEW_PENDING':  'bg-purple-900/40 border-purple-600/50 text-purple-300',
    'READY_FOR_TEST':  'bg-cyan-900/40 border-cyan-600/50 text-cyan-300',
    'WAIT_FOR_DEPLOY': 'bg-orange-900/40 border-orange-600/50 text-orange-300',
    'READY_FOR_UAT':   'bg-amber-900/40 border-amber-600/50 text-amber-300',
    'ASSIGNED':        'bg-gray-700/60 border-gray-600/50 text-gray-300',
  }
  return classes[status] || 'bg-gray-700/60 border-gray-600/50 text-gray-300'
}

const getStatusClass = (status: string) => getStatusBadgeClass(status)

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'COMPLETED':       '✅ COMPLETED',
    'IN_PROGRESS':     '🔄 IN PROGRESS',
    'PENDING':         '⏳ PENDING',
    'REVIEW_PENDING':  '⏳ WAITING FOR APPROVAL',
    'READY_FOR_TEST':  '🧪 READY FOR TEST',
    'WAIT_FOR_DEPLOY': '🚀 WAIT FOR DEPLOY',
    'READY_FOR_UAT':   '🔬 READY FOR UAT',
    'ASSIGNED':        '📌 ASSIGNED',
  }
  return labels[status] || status
}

/** Display task id as A001 / B001 / C001: top-level = A, sub-task = B, sub-task of sub-task = C */
function taskCodeDisplay(t: Task | null): string {
  if (!t) return '–'
  const suffix = t.code ? t.code.split('-').pop() : ''
  const num = (suffix && /^\d+$/.test(suffix)) ? String(Number(suffix)).padStart(3, '0') : (suffix || t.id.slice(0, 4))
  if (!t.parent_id) return 'A' + num
  if (t.parent_task?.parent_id) return 'C' + num
  return 'B' + num
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Deadline Helpers
const getDeadlineUrgency = (task: Task) => {
  if (!task.due_at || task.status === 'COMPLETED') return 'none'
  
  const now = new Date().getTime()
  const dueDate = new Date(task.due_at).getTime()
  const hoursUntilDue = (dueDate - now) / (1000 * 60 * 60)
  
  if (hoursUntilDue < 0) return 'overdue'
  if (hoursUntilDue < 24) return 'urgent'
  return 'normal'
}

const getDeadlineCountdown = (dueAt: string) => {
  const now = new Date().getTime()
  const due = new Date(dueAt).getTime()
  const diff = due - now
  
  if (diff < 0) {
    // Overdue
    const hours = Math.abs(Math.floor(diff / (1000 * 60 * 60)))
    const days = Math.floor(hours / 24)
    if (days > 0) return `Overdue by ${days} days`
    return `Overdue by ${hours} hours`
  }
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)
  
  if (days > 0) return `${days} days left`
  if (hours > 0) return `${hours} hours left`
  return 'Due very soon!'
}

const calculateDuration = (startAt: string, completedAt: string) => {
  const start = new Date(startAt).getTime()
  const end = new Date(completedAt).getTime()
  const diff = end - start
  
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const days = Math.floor(hours / 24)
  
  if (days > 0) {
    const remainingHours = hours % 24
    return `${days}d ${remainingHours}h`
  }
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  }
  return `${minutes}m`
}

/** Convert ISO (UTC) to "YYYY-MM-DD" in local time for date input */
const toDateOnly = (iso: string | null | undefined) => {
  if (!iso) return ''
  const d = new Date(iso)
  if (isNaN(d.getTime())) return ''
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const dateOnlyToISO = (ymd: string) => new Date(`${ymd}T00:00:00`).toISOString()

// Edit Task Methods
const openEditModal = async () => {
  if (!task.value) return

  editForm.value.title = task.value.title
  editForm.value.description = task.value.description || ''
  editForm.value.task_type = task.value.task_type || 'TASK'
  editForm.value.deadline = toDateOnly(task.value.due_at)
  editForm.value.priority = task.value.priority || 'MEDIUM'
  editForm.value.story_points = task.value.story_points ?? 0
  editForm.value.sprint_id = task.value.sprint_id ?? ''
  editForm.value.start_date = toDateOnly(task.value.start_date)
  editForm.value.end_date = toDateOnly(task.value.end_date)
  editForm.value.estimated_hours = minutesToEffortHours(task.value.estimated_minutes ?? 0)

  editSprints.value = []
  if (task.value.project_id) {
    try {
      const list = await projectsApi.getSprints(task.value.project_id)
      editSprints.value = list.map((s) => ({ id: s.id, name: s.name }))
    } catch {
      // ignore
    }
  }

  editError.value = ''
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.value = {
    title: '',
    description: '',
    task_type: 'TASK',
    deadline: '',
    priority: 'MEDIUM',
    story_points: 0,
    sprint_id: '',
    start_date: '',
    end_date: '',
    estimated_hours: 0
  }
  editError.value = ''
}

const closeDateTimePicker = (event: Event) => {
  const target = event.target as HTMLInputElement | null
  if (!target) return
  // Native datetime picker behaves differently across browsers; do both immediate and delayed blur.
  target.blur()
  window.requestAnimationFrame(() => target.blur())
  window.setTimeout(() => {
    target.blur()
    ;(document.activeElement as HTMLElement | null)?.blur?.()
    const fallback = document.querySelector('.edit-task-modal') as HTMLElement | null
    fallback?.focus?.()
  }, 10)
}

const submitEdit = async () => {
  if (!task.value) return
  
  // Validation
  if (!editForm.value.title.trim()) {
    editError.value = 'Title is required'
    return
  }
  
  try {
    isUpdatingTask.value = true
    editError.value = ''
    
    const taskId = route.params.id as string
    
    // Prepare request body (only send changed fields) — same fields as Create Task
    const body: Record<string, string | number> = {}

    if (editForm.value.title !== task.value.title) {
      body.title = editForm.value.title
    }
    if (editForm.value.description !== (task.value.description || '')) {
      body.description = editForm.value.description
    }
    const currentType = task.value.task_type || 'TASK'
    if (editForm.value.task_type && editForm.value.task_type !== currentType) {
      body.task_type = editForm.value.task_type
    }
    if (editForm.value.priority && editForm.value.priority !== (task.value.priority || 'MEDIUM')) {
      body.priority = editForm.value.priority
    }
    const currentSp = task.value.story_points ?? 0
    if (Number(editForm.value.story_points) !== currentSp) {
      body.story_points = Number(editForm.value.story_points) || 0
    }
    const currentSprint = task.value.sprint_id ?? ''
    if (editForm.value.sprint_id !== currentSprint) {
      body.sprint_id = editForm.value.sprint_id || ''
    }
    const currentStart = toDateOnly(task.value.start_date)
    if (editForm.value.start_date !== currentStart && editForm.value.start_date) {
      body.start_date = dateOnlyToISO(editForm.value.start_date)
    }
    const currentDue = toDateOnly(task.value.due_at)
    if (editForm.value.deadline !== currentDue && editForm.value.deadline) {
      body.due_at = dateOnlyToISO(editForm.value.deadline)
    }
    const currentEnd = toDateOnly(task.value.end_date)
    if (editForm.value.end_date !== currentEnd && editForm.value.end_date) {
      body.end_date = dateOnlyToISO(editForm.value.end_date)
    }
    const currentEst = task.value.estimated_minutes ?? 0
    const newEst = effortHoursToMinutes(Number(editForm.value.estimated_hours) || 0)
    if (newEst !== currentEst) {
      body.estimated_minutes = newEst
    }
    
    if (Object.keys(body).length === 0) {
      editError.value = 'No changes detected. Please modify at least one field.'
      isUpdatingTask.value = false
      return
    }
    
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'PATCH',
      body
    })
    
    showSuccess('Task updated.', 'Done')
    closeEditModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to update task:', err)
    showError(err.data?.message || err.message || 'Failed to update task')
    editError.value = err.data?.message || err.message || 'Failed to update task'
  } finally {
    isUpdatingTask.value = false
  }
}

// Delete Task Methods
const openDeleteConfirmation = () => {
  deleteError.value = ''
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  deleteError.value = ''
}

const confirmDelete = async () => {
  if (!task.value) return
  
  try {
    isDeletingTask.value = true
    deleteError.value = ''
    
    const taskId = route.params.id as string
    
    // Call DELETE API
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'DELETE'
    })
    
    showSuccess('Task deleted.', 'Done')
    goToDashboard()
  } catch (err: any) {
    console.error('Failed to delete task:', err)
    const raw = err.data?.message || err.message || ''
    const isSubTasksBlock = raw.includes('fk_tasks_sub_tasks') || raw.includes('23503') || raw.includes('sub-tasks')
    const message = isSubTasksBlock
      ? 'มี sub task ไม่สามารถลบได้ หากต้องการลบ ต้องลบ sub task ก่อน'
      : (raw || 'Failed to delete task')
    showError(message)
    deleteError.value = message
  } finally {
    isDeletingTask.value = false
  }
}

// Inline Description Edit
const startInlineEdit = () => {
  inlineDescriptionHtml.value = task.value?.description || ''
  isEditingDescription.value = true
}

const cancelInlineEdit = () => {
  isEditingDescription.value = false
  inlineDescriptionHtml.value = ''
}

const saveInlineDescription = async () => {
  if (!task.value) return
  isSavingDescription.value = true
  try {
    const taskId = route.params.id as string
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'PATCH',
      body: { description: inlineDescriptionHtml.value }
    })
    task.value.description = inlineDescriptionHtml.value
    isEditingDescription.value = false
    inlineDescriptionHtml.value = ''
  } catch (err: any) {
    console.error('Failed to save description:', err)
  } finally {
    isSavingDescription.value = false
  }
}

// Back: if came from project page (from_project + from_tab) or sprint page (from_sprint + from_project), return there; else project or dashboard
const backTarget = computed(() => {
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  const fromTab = route.query.from_tab as string
  const from = route.query.from as string
  if (fromSprint && fromProject) {
    return `/projects/sprint/${fromSprint}?project=${encodeURIComponent(fromProject)}`
  }
  if (fromProject) {
    const tab = fromTab || 'backlog'
    return `/projects/${fromProject}?tab=${tab}`
  }
  if (from === 'dashboard') {
    const tab = fromTab || 'board'
    return `/dashboard?tab=${tab}`
  }
  if (from === 'workload') {
    return '/workload'
  }
  if (from === 'discipline') {
    return '/discipline'
  }
  const t = task.value
  if (!t) return '/dashboard'
  // Sub-tasks: navigate back to parent task page after deletion
  if (t.parent_id) return `/task/${t.parent_id}`
  // Komgrip tasks have no project — go back to Komgrip page
  if (t.is_komgrip) return '/komgrip'
  // Use project_id (UUID) first — always reliable
  if (t.project_id) return `/projects/${t.project_id}`
  // Fallback: try deriving project code from task code (only reliable for top-level tasks)
  if (t.code) {
    const projectCode = t.code.replace(/-[0-9]+$/, '')
    if (projectCode !== t.code) return `/projects/${projectCode}`
  }
  return '/dashboard'
})

// Whether we're in sprint context (opened from sprint page) — show Prev/Next area
const inSprintContext = computed(() => {
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  return !!(fromSprint && fromProject)
})

// Whether we're in backlog context (opened from project backlog board)
const inBacklogContext = computed(() => {
  const fromProject = route.query.from_project as string | undefined
  const fromTab = route.query.from_tab as string | undefined
  return !!(fromProject && fromTab === 'backlog')
})

const showPrevNext = computed(() => inSprintContext.value || inBacklogContext.value)

// Prev/Next task links within same sprint (only when from_sprint + from_project)
// Normalize id comparison (UUID may differ in casing between getTask and getTasksByProject)
const currentSprintIndex = computed(() => {
  const t = task.value
  if (!t?.id || !sprintTaskIds.value.length) return -1
  const needle = String(t.id).toLowerCase()
  return sprintTaskIds.value.findIndex((id) => String(id).toLowerCase() === needle)
})
const prevTaskLink = computed(() => {
  if (currentSprintIndex.value <= 0) return null
  const id = sprintTaskIds.value[currentSprintIndex.value - 1]
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  if (!id || !fromSprint || !fromProject) return null
  return `/task/${id}?from_sprint=${encodeURIComponent(fromSprint)}&from_project=${encodeURIComponent(fromProject)}`
})
const nextTaskLink = computed(() => {
  const idx = currentSprintIndex.value
  if (idx < 0 || idx >= sprintTaskIds.value.length - 1) return null
  const id = sprintTaskIds.value[idx + 1]
  const fromSprint = route.query.from_sprint as string
  const fromProject = route.query.from_project as string
  if (!id || !fromSprint || !fromProject) return null
  return `/task/${id}?from_sprint=${encodeURIComponent(fromSprint)}&from_project=${encodeURIComponent(fromProject)}`
})

// Prev/Next task links within same backlog display order (only when from_project + from_tab=backlog)
const currentBacklogIndex = computed(() => {
  const t = task.value
  if (!t?.id || !backlogTaskIds.value.length) return -1
  const needle = String(t.id).toLowerCase()
  return backlogTaskIds.value.findIndex((id) => String(id).toLowerCase() === needle)
})

function buildBacklogTaskLink(id: string) {
  const fromProject = route.query.from_project as string | undefined
  const fromTab = route.query.from_tab as string | undefined
  if (!fromProject || !fromTab) return null
  return `/task/${id}?from_project=${encodeURIComponent(fromProject)}&from_tab=${encodeURIComponent(fromTab)}`
}

const prevBacklogTaskLink = computed(() => {
  const idx = currentBacklogIndex.value
  if (idx <= 0) return null
  const id = backlogTaskIds.value[idx - 1]
  return id ? buildBacklogTaskLink(id) : null
})

const nextBacklogTaskLink = computed(() => {
  const idx = currentBacklogIndex.value
  if (idx < 0 || idx >= backlogTaskIds.value.length - 1) return null
  const id = backlogTaskIds.value[idx + 1]
  return id ? buildBacklogTaskLink(id) : null
})

const activePrevTaskLink = computed(() => (inSprintContext.value ? prevTaskLink.value : prevBacklogTaskLink.value))
const activeNextTaskLink = computed(() => (inSprintContext.value ? nextTaskLink.value : nextBacklogTaskLink.value))

function showSprintNavFeedback(msg: string) {
  sprintNavMessage.value = msg
  setTimeout(() => { sprintNavMessage.value = '' }, 2500)
}

function goToPrevTask() {
  const link = activePrevTaskLink.value
  if (link) {
    navigateTo(link)
  } else {
    showSprintNavFeedback('ไม่มี task ก่อนหน้า')
  }
}

function goToNextTask() {
  const link = activeNextTaskLink.value
  if (link) {
    navigateTo(link)
  } else {
    showSprintNavFeedback('ไม่มี task ถัดไป')
  }
}

const goToDashboard = () => {
  navigateTo(backTarget.value)
}

async function fetchCommentsAndLogs() {
  const taskId = route.params.id as string
  if (!taskId) return
  const tasksApi = useTasksApi()
  commentsLoading.value = true
  timeLogsLoading.value = true
  try {
    const [c, l] = await Promise.all([
      tasksApi.getComments(taskId),
      tasksApi.getTimeLogs(taskId),
    ])
    comments.value = c
    timeLogs.value = l
  } catch {
    // non-critical
  } finally {
    commentsLoading.value = false
    timeLogsLoading.value = false
  }
}

async function handleAddComment(payload: { content: string; attachments: File[] }) {
  const taskId = route.params.id as string
  const tasksApi = useTasksApi()
  commentsLoading.value = true
  try {
    const comment = await tasksApi.addComment(taskId, payload.content, payload.attachments)
    comments.value.push(comment)
  } catch (e: any) {
    console.error('Failed to add comment:', e)
  } finally {
    commentsLoading.value = false
  }
}

async function handleEditComment(payload: { commentId: string; content: string }) {
  const tasksApi = useTasksApi()
  commentsLoading.value = true
  try {
    const updated = await tasksApi.editComment(payload.commentId, payload.content)
    const idx = comments.value.findIndex(c => c.id === payload.commentId)
    if (idx >= 0) comments.value[idx] = updated
  } catch (e: any) {
    console.error('Failed to edit comment:', e)
  } finally {
    commentsLoading.value = false
  }
}

async function handleLogTime(minutes: number, description: string, workType: string, loggedDate: string, isTimer: boolean) {
  const taskId = route.params.id as string
  const tasksApi = useTasksApi()
  timeLogsLoading.value = true
  try {
    const log = await tasksApi.logTime(taskId, minutes, description, workType, loggedDate, isTimer)
    timeLogs.value.unshift(log)
  }
  catch (e: any) {
    console.error('Failed to log time:', e)
  }
  finally {
    timeLogsLoading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await fetchTask()
  fetchCommentsAndLogs()
  fetchDeploymentForTask()
  fetchPrevNextOrdering()
})

watch(
  () => route.params.id,
  async (newId, oldId) => {
    const id = typeof newId === 'string' ? newId.trim() : ''
    if (!id || id === oldId) return
    await fetchTask()
    fetchCommentsAndLogs()
    fetchDeploymentForTask()
    fetchPrevNextOrdering()
  }
)
</script>

<style scoped>
.task-enterprise-bg {
  background:
    radial-gradient(1200px 620px at 82% -15%, rgba(139, 92, 246, 0.18), transparent 60%),
    radial-gradient(900px 460px at -8% 2%, rgba(59, 130, 246, 0.16), transparent 56%),
    linear-gradient(180deg, #070b16 0%, #0b1220 55%, #090f1a 100%);
}

.enterprise-hero-panel {
  @apply border border-white/10 bg-gradient-to-br from-slate-900/85 to-slate-900/65 shadow-[0_20px_50px_rgba(2,6,23,0.45)];
}

.enterprise-card {
  @apply bg-slate-900/70 border border-white/10 shadow-[0_14px_34px_rgba(2,6,23,0.4)] backdrop-blur-sm;
}

.label {
  @apply block text-xs text-slate-300 mb-1.5 font-semibold tracking-wide;
}
.input-field {
  @apply bg-slate-800/95 border border-slate-600/80 rounded-xl px-4 py-2.5 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-violet-400 focus:ring-2 focus:ring-violet-500/35 transition-all;
}
/* Edit Task modal — wide layout, larger fields (rest of page unchanged) */
.edit-task-modal-enterprise {
  @apply border border-white/10 bg-slate-900/75 shadow-[0_18px_42px_rgba(2,6,23,0.46)] backdrop-blur-sm;
}

.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: rgba(100, 116, 139, 0.5) transparent;
}
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(100, 116, 139, 0.45);
  border-radius: 999px;
}

.edit-task-modal .label {
  @apply block text-sm sm:text-base text-slate-200 mb-2 font-semibold tracking-wide;
}
.edit-task-modal .input-field {
  @apply bg-slate-800/95 border border-slate-500/80 rounded-xl px-4 py-3.5 text-base text-slate-100 placeholder-slate-500 focus:outline-none focus:border-violet-400 focus:ring-2 focus:ring-violet-500/35 transition-all;
}
.edit-task-modal :deep(.rich-editor .editor-content) {
  min-height: 18rem;
  max-height: min(58vh, 800px);
}
.edit-task-modal :deep(.rich-editor .ProseMirror) {
  font-size: 1rem;
  line-height: 1.65;
}
.edit-task-modal :deep(.editor-toolbar .toolbar-btn) {
  @apply text-sm px-2.5 py-1.5 min-h-[2.25rem];
}
.btn-primary {
  @apply bg-gradient-to-r from-violet-600 via-fuchsia-600 to-indigo-600 hover:from-violet-500 hover:via-fuchsia-500 hover:to-indigo-500 text-white font-semibold rounded-xl transition-all shadow-[0_12px_24px_rgba(124,58,237,0.32)];
}
</style>
