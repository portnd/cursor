<template>
  <div class="min-h-screen bg-gray-900 text-gray-100">
    <!-- Enterprise Header -->
    <header class="sticky top-0 z-10 border-b border-gray-800 bg-gray-900/95 backdrop-blur-sm">
      <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-4 py-6 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
              Execution Dashboard
            </h1>
            <p class="mt-1 text-sm text-gray-400">
              Deep work & task pipeline
            </p>
          </div>
          <div class="flex flex-wrap items-center gap-3">
            <div class="flex items-center gap-4 rounded-xl border border-gray-700/80 bg-gray-800/60 px-4 py-2.5">
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Focus</span>
              <span class="text-lg font-bold text-purple-400">{{ currentFocus.length }}</span>
              <span class="h-4 w-px bg-gray-600" />
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Backlog</span>
              <span class="text-lg font-bold text-amber-400">{{ myBacklog.length }}</span>
              <span class="h-4 w-px bg-gray-600" />
              <span class="text-xs font-medium uppercase tracking-wider text-gray-500">Available</span>
              <span class="text-lg font-bold text-emerald-400">{{ availableMissions.length }}</span>
            </div>
            <button
              type="button"
              @click="fetchTasks"
              class="inline-flex items-center gap-2 rounded-xl border border-gray-600 bg-gray-800 px-4 py-2.5 text-sm font-medium text-gray-200 transition-colors hover:border-gray-500 hover:bg-gray-700 hover:text-white"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refresh
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="isLoading" class="flex flex-col items-center justify-center py-24">
      <div class="h-10 w-10 animate-spin rounded-full border-2 border-purple-500 border-t-transparent" />
      <p class="mt-4 text-sm text-gray-500">Loading your tasks...</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <div class="rounded-xl border border-red-500/50 bg-red-950/20 px-5 py-4 text-red-400">
        <div class="flex items-start gap-3">
          <span class="text-xl">⚠️</span>
          <div>
            <p class="font-medium">Failed to load data</p>
            <p class="mt-1 text-sm text-red-300">{{ error }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <main v-else class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <!-- Your Performance -->
      <section v-if="performanceStore.personal" class="mb-8">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Your performance
          </h2>
          <button
            type="button"
            @click="showKpiPanel = !showKpiPanel"
            class="text-xs font-medium text-gray-500 transition-colors hover:text-white"
          >
            {{ showKpiPanel ? 'Hide' : 'Show' }}
          </button>
        </div>
        <template v-if="showKpiPanel">
          <div class="grid grid-cols-2 gap-4 md:grid-cols-3 lg:grid-cols-6">
            <PerformanceKpiScoreCard
              label="Delivery Rate"
              :value="performanceStore.personal.delivery_rate_pct"
              format="pct"
              :status="deliveryStatus(performanceStore.personal.delivery_rate_pct)"
            />
            <PerformanceKpiScoreCard
              label="Code Quality"
              :value="performanceStore.personal.code_quality_index"
              :status="qualityStatus(performanceStore.personal.code_quality_index)"
            />
            <PerformanceKpiScoreCard
              label="Rework Rate"
              :value="performanceStore.personal.rework_rate_pct"
              format="pct"
              :status="reworkStatus(performanceStore.personal.rework_rate_pct)"
            />
            <PerformanceKpiScoreCard
              label="Time Accuracy"
              :value="performanceStore.personal.time_accuracy_pct"
              format="pct"
              :status="accuracyStatus(performanceStore.personal.time_accuracy_pct)"
            />
            <PerformanceKpiScoreCard
              label="Sprint Velocity"
              :value="performanceStore.personal.sprint_velocity_sp"
              format="number"
              sublabel="SP (last 3 sprints)"
              :trend="performanceStore.personal.velocity_trend === 'up' ? 'up' : performanceStore.personal.velocity_trend === 'down' ? 'down' : 'stable'"
              :trend-label="performanceStore.personal.velocity_trend"
            />
            <PerformanceKpiScoreCard
              label="Score"
              :value="compositeScore"
              format="number"
              sublabel="Composite (0–100)"
              :status="scoreStatus(compositeScore)"
            />
          </div>
        </template>
      </section>

      <!-- Current Focus -->
      <section v-if="currentFocus.length > 0" class="mb-8">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-purple-400">
            Current focus
          </h2>
          <span class="rounded-full bg-purple-500/20 px-3 py-1 text-xs font-medium text-purple-300">
            {{ currentFocus.length }} in progress
          </span>
        </div>
        <div class="grid gap-4">
          <div
            v-for="task in currentFocus"
            :key="task.id"
            :class="[
              'relative cursor-pointer rounded-xl border-2 p-6 shadow-lg transition-all hover:shadow-purple-500/10',
              getDeadlineBorderClass(task) === 'border-red-500' ? 'border-red-500/60 bg-red-950/10' :
              getDeadlineBorderClass(task) === 'border-yellow-500' ? 'border-amber-500/60 bg-amber-950/10' :
              'border-purple-500/50 bg-purple-950/10'
            ]"
            @click="goToTask(task)"
          >
            <div
              v-if="getDeadlineUrgency(task) === 'overdue'"
              class="absolute -top-2 right-4 rounded-md bg-red-600 px-2 py-1 text-xs font-bold text-white"
            >
              Overdue
            </div>
            <div
              v-else-if="getDeadlineUrgency(task) === 'urgent'"
              class="absolute -top-2 right-4 rounded-md bg-amber-500 px-2 py-1 text-xs font-bold text-black"
            >
              Urgent
            </div>
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <h3 class="text-xl font-bold text-white">{{ task.title }}</h3>
                <p class="mt-1 line-clamp-2 text-sm text-gray-400">{{ task.description || 'No description' }}</p>
              </div>
              <span class="shrink-0 rounded-md bg-purple-600/80 px-3 py-1 text-xs font-semibold text-purple-100">
                In progress
              </span>
            </div>
            <div class="mt-4 flex flex-wrap items-center gap-4 text-sm">
              <span class="text-gray-500">Est:</span>
              <span class="font-semibold text-purple-400">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}h</span>
              <template v-if="task.started_at">
                <span class="text-gray-500">Started:</span>
                <span class="text-gray-400">{{ formatDate(task.started_at) }}</span>
              </template>
              <template v-if="task.due_at">
                <span
                  :class="[
                    'text-xs font-semibold',
                    getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' :
                    getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400' : 'text-gray-400'
                  ]"
                >
                  Due {{ formatDeadline(task.due_at) }} · {{ getDeadlineCountdown(task.due_at) }}
                </span>
              </template>
            </div>
            <div class="mt-4 flex gap-3">
              <button
                type="button"
                @click.stop="openSubmitModal(task)"
                class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 px-4 py-2 text-sm font-semibold text-white shadow-md transition-all hover:from-purple-500 hover:to-pink-500"
              >
                Submit work
              </button>
              <NuxtLink
                :to="`/task/${task.code || task.id}`"
                @click.stop
                class="rounded-lg border border-gray-600 bg-gray-700 px-4 py-2 text-sm font-medium text-gray-200 transition-colors hover:bg-gray-600 hover:text-white"
              >
                View details
              </NuxtLink>
            </div>
          </div>
        </div>
      </section>

      <!-- My Backlog -->
      <section v-if="myBacklog.length > 0" class="mb-8">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500">
            My backlog
          </h2>
          <span class="rounded-full bg-amber-500/20 px-3 py-1 text-xs font-medium text-amber-400">
            {{ myBacklog.length }} pending
          </span>
        </div>
        <div class="overflow-hidden rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/80 text-left text-xs font-medium uppercase tracking-wider text-gray-400">
                  <th class="px-5 py-4">Task</th>
                  <th class="px-5 py-4 text-right">Est. hours</th>
                  <th class="px-5 py-4 text-right">Deadline</th>
                  <th class="px-5 py-4 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/80">
                <tr
                  v-for="task in myBacklog"
                  :key="task.id"
                  :class="[
                    'cursor-pointer transition-colors hover:bg-gray-700/30',
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-950/20' :
                    getDeadlineUrgency(task) === 'urgent' ? 'bg-amber-950/20' : ''
                  ]"
                  @click="goToTask(task)"
                >
                  <td class="px-5 py-4">
                    <div class="font-medium text-white">
                      {{ task.title }}
                      <span
                        v-if="getDeadlineUrgency(task) === 'overdue'"
                        class="ml-2 rounded bg-red-600 px-2 py-0.5 text-xs font-semibold text-white"
                      >Overdue</span>
                      <span
                        v-else-if="getDeadlineUrgency(task) === 'urgent'"
                        class="ml-2 rounded bg-amber-500 px-2 py-0.5 text-xs font-semibold text-black"
                      >Urgent</span>
                    </div>
                  </td>
                  <td class="px-5 py-4 text-right text-gray-400">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}h</td>
                  <td class="px-5 py-4 text-right text-xs" :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-500'">
                    <template v-if="task.due_at">
                      {{ formatDeadline(task.due_at) }}
                      <div class="mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</div>
                    </template>
                    <span v-else class="text-gray-600">No deadline</span>
                  </td>
                  <td class="px-5 py-4 text-right">
                    <button
                      type="button"
                      @click.stop="goToTask(task)"
                      class="rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:from-purple-500 hover:to-pink-500"
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

      <!-- Available Missions -->
      <section v-if="availableMissions.length > 0" class="mb-8">
        <div class="mb-4 flex items-center justify-between">
          <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500">
            Available missions
          </h2>
          <span class="rounded-full bg-emerald-500/20 px-3 py-1 text-xs font-medium text-emerald-400">
            {{ availableMissions.length }} unassigned
          </span>
        </div>
        <div class="overflow-hidden rounded-xl border border-gray-700/80 bg-gray-800/60 shadow-lg">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b border-gray-700 bg-gray-900/80 text-left text-xs font-medium uppercase tracking-wider text-gray-400">
                  <th class="px-5 py-4">Task</th>
                  <th class="hidden px-5 py-4 lg:table-cell">Description</th>
                  <th class="px-5 py-4 text-right">Est. hours</th>
                  <th class="px-5 py-4 text-right">Deadline</th>
                  <th class="px-5 py-4 text-right">Action</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-700/80">
                <tr
                  v-for="task in availableMissions"
                  :key="task.id"
                  :class="[
                    'cursor-pointer transition-colors hover:bg-gray-700/30',
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-950/20' :
                    getDeadlineUrgency(task) === 'urgent' ? 'bg-amber-950/20' : ''
                  ]"
                  @click="goToTask(task)"
                >
                  <td class="px-5 py-4">
                    <div class="font-medium text-white">
                      {{ task.title }}
                      <span v-if="getDeadlineUrgency(task) === 'overdue'" class="ml-2 rounded bg-red-600 px-2 py-0.5 text-xs font-semibold text-white">Overdue</span>
                      <span v-else-if="getDeadlineUrgency(task) === 'urgent'" class="ml-2 rounded bg-amber-500 px-2 py-0.5 text-xs font-semibold text-black">Urgent</span>
                    </div>
                  </td>
                  <td class="hidden max-w-xs truncate px-5 py-4 text-gray-500 lg:table-cell">{{ task.description || 'No description' }}</td>
                  <td class="px-5 py-4 text-right text-gray-400">{{ (task.ai_estimated_minutes / 60).toFixed(1) }}h</td>
                  <td class="px-5 py-4 text-right text-xs" :class="getDeadlineUrgency(task) === 'overdue' ? 'text-red-400 font-semibold' : getDeadlineUrgency(task) === 'urgent' ? 'text-amber-400 font-semibold' : 'text-gray-500'">
                    <template v-if="task.due_at">
                      {{ formatDeadline(task.due_at) }}
                      <div class="mt-0.5">{{ getDeadlineCountdown(task.due_at) }}</div>
                    </template>
                    <span v-else class="text-gray-600">No deadline</span>
                  </td>
                  <td class="px-5 py-4 text-right">
                    <button
                      type="button"
                      @click.stop="claimTask(task.id)"
                      :disabled="claiming === task.id"
                      class="rounded-lg bg-emerald-600 px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-emerald-500 disabled:pointer-events-none disabled:opacity-50"
                    >
                      {{ claiming === task.id ? 'Claiming…' : 'Claim' }}
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Empty State -->
      <section v-if="currentFocus.length === 0 && myBacklog.length === 0 && availableMissions.length === 0">
        <div class="flex flex-col items-center justify-center rounded-xl border border-gray-700/80 bg-gray-800/60 py-16 text-center">
          <div class="rounded-full bg-gray-700/80 p-5 text-gray-500">
            <svg class="h-10 w-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <p class="mt-4 text-lg font-medium text-gray-400">No tasks available</p>
          <p class="mt-1 text-sm text-gray-500">Create a new task or wait for assignments</p>
          <NuxtLink
            to="/create"
            class="mt-6 inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-purple-600 to-pink-600 px-5 py-2.5 text-sm font-semibold text-white shadow-lg transition-all hover:from-purple-500 hover:to-pink-500"
          >
            Create new task
          </NuxtLink>
        </div>
      </section>
    </main>

    <!-- Success Toast -->
    <Teleport to="body">
      <Transition name="toast">
        <div
          v-if="showSuccess"
          class="fixed bottom-6 right-6 z-50 flex items-center gap-3 rounded-xl border border-emerald-500/40 bg-gray-800 px-5 py-4 shadow-xl shadow-black/30"
        >
          <span class="flex h-9 w-9 items-center justify-center rounded-full bg-emerald-600 text-white">✓</span>
          <p class="text-sm font-medium text-white">{{ successMessage }}</p>
        </div>
      </Transition>
    </Teleport>

    <!-- Submit Work Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showSubmitModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="submit-modal-title"
          @keydown.escape="closeSubmitModal"
        >
          <div class="fixed inset-0 bg-black/70 backdrop-blur-sm" @click="closeSubmitModal" />
          <div
            class="relative max-h-[90vh] w-full max-w-3xl overflow-y-auto rounded-2xl border-2 border-purple-500/50 bg-gray-800 shadow-2xl"
            @click.stop
          >
            <div class="sticky top-0 z-10 flex items-start justify-between border-b border-gray-700 bg-gray-800 px-6 py-4">
              <div>
                <h2 id="submit-modal-title" class="text-xl font-semibold text-white">Submit mission</h2>
                <p class="mt-1 line-clamp-2 text-sm text-gray-400">{{ selectedTask?.title }}</p>
              </div>
              <button
                type="button"
                @click="closeSubmitModal"
                class="rounded-lg p-2 text-gray-400 transition-colors hover:bg-gray-700 hover:text-white"
                title="Close"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
            </div>
            <form class="p-6" @submit.prevent="submitWork">
              <div v-if="error" class="mb-4 rounded-xl border border-red-500/50 bg-red-950/20 px-4 py-3 text-sm text-red-400">
                {{ error }}
              </div>
              <div class="space-y-4">
                <div>
                  <label for="commit-hash" class="mb-2 block text-sm font-medium text-gray-300">Commit hash <span class="text-red-400">*</span></label>
                  <input
                    id="commit-hash"
                    v-model="submitForm.commitHash"
                    type="text"
                    placeholder="e.g. a3f5d9c7b2e8..."
                    class="w-full rounded-xl border border-gray-600 bg-gray-900 px-4 py-3 font-mono text-sm text-white placeholder-gray-500 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/30 focus:outline-none"
                  />
                </div>
                <div>
                  <label for="code-diff" class="mb-2 block text-sm font-medium text-gray-300">Code diff / changes <span class="text-red-400">*</span></label>
                  <textarea
                    id="code-diff"
                    v-model="submitForm.diff"
                    rows="12"
                    placeholder="Paste your code changes here..."
                    class="w-full resize-none rounded-xl border border-gray-600 bg-gray-900 px-4 py-3 font-mono text-sm text-white placeholder-gray-500 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/30 focus:outline-none"
                    style="min-height: 250px"
                  />
                  <p class="mt-1 text-xs text-gray-500">Include context lines for better AI analysis</p>
                </div>
              </div>
              <div class="mt-6 flex gap-3">
                <button
                  type="button"
                  @click="closeSubmitModal"
                  :disabled="isSubmitting"
                  class="flex-1 rounded-xl border border-gray-600 bg-gray-700 px-4 py-3 text-sm font-medium text-gray-200 transition-colors hover:bg-gray-600 disabled:opacity-50"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  :disabled="isSubmitting || !submitForm.commitHash || !submitForm.diff"
                  class="flex-1 rounded-xl bg-gradient-to-r from-purple-600 to-pink-600 px-4 py-3 text-sm font-semibold text-white shadow-lg transition-all hover:from-purple-500 hover:to-pink-500 disabled:pointer-events-none disabled:opacity-50"
                >
                  {{ isSubmitting ? 'Analyzing…' : 'Submit to AI auditor' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Result Modal (Verdict) -->
    <Teleport to="body">
      <Transition name="modal">
        <div
          v-if="showResultModal"
          class="fixed inset-0 z-50 flex items-center justify-center p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="result-modal-title"
          @keydown.escape="closeResultModal"
        >
          <div class="fixed inset-0 bg-black/80 backdrop-blur-sm" @click="closeResultModal" />
          <div
            :class="[
              'relative max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-2xl border-2 shadow-2xl',
              submissionResult?.ai_verdict === 'PASS'
                ? 'border-emerald-500/50 bg-gray-800'
                : 'border-red-500/50 bg-gray-800'
            ]"
            @click.stop
          >
            <div class="p-8">
              <button
                type="button"
                @click="closeResultModal"
                :class="[
                  'absolute right-4 top-4 rounded-lg p-2 transition-colors',
                  submissionResult?.ai_verdict === 'PASS' ? 'text-emerald-400 hover:bg-emerald-500/20 hover:text-white' : 'text-red-400 hover:bg-red-500/20 hover:text-white'
                ]"
                title="Close"
              >
                <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
              </button>
              <div class="text-center">
                <div
                  :class="[
                    'text-5xl font-black sm:text-6xl',
                    submissionResult?.ai_verdict === 'PASS' ? 'text-emerald-400' : 'text-red-400'
                  ]"
                >
                  {{ submissionResult?.ai_verdict === 'PASS' ? 'Mission accomplished' : 'Mission failed' }}
                </div>
                <div class="mt-4 text-2xl font-bold text-gray-300 sm:text-3xl">
                  AI score: {{ submissionResult?.ai_score }}/100
                </div>
              </div>
              <div
                :class="[
                  'mt-6 rounded-xl border-2 p-6',
                  submissionResult?.ai_verdict === 'PASS'
                    ? 'border-emerald-500/40 bg-emerald-950/30'
                    : 'border-red-500/40 bg-red-950/30'
                ]"
              >
                <div class="mb-3 flex items-center gap-2 font-semibold text-white">
                  <span>AI auditor feedback</span>
                </div>
                <div class="whitespace-pre-wrap text-sm leading-relaxed text-gray-200">
                  {{ getFeedbackText() }}
                </div>
              </div>
              <div class="mt-6 space-y-3">
                <button
                  type="button"
                  @click="closeResultModal"
                  :class="[
                    'w-full rounded-xl px-6 py-4 text-lg font-bold text-white transition-colors',
                    submissionResult?.ai_verdict === 'PASS'
                      ? 'bg-emerald-600 hover:bg-emerald-500'
                      : 'bg-red-600 hover:bg-red-500'
                  ]"
                >
                  {{ submissionResult?.ai_verdict === 'PASS' ? 'Continue working' : 'Go to task & appeal' }}
                </button>
                <div v-if="submissionResult?.ai_verdict !== 'PASS'" class="text-center text-sm text-gray-400">
                  On the task page you can appeal, review feedback, or submit revised code.
                </div>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import PerformanceKpiScoreCard from '~/components/performance/KpiScoreCard.vue'
import { usePerformanceStore } from '~/core/modules/performance/performance-store'

interface Task {
  id: string
  title: string
  description: string
  status: string
  ai_estimated_minutes: number
  assigned_to?: number
  created_at: string
  due_at?: string
  started_at?: string
  completed_at?: string
}

interface Submission {
  id: string
  task_id: string
  dev_id: number
  commit_hash: string
  ai_verdict: string
  ai_score: number
  ai_feedback: any
  created_at: string
}

const { fetchWithAuth, currentUser } = useAuth()
const performanceStore = usePerformanceStore()

const myTasks = ref<Task[]>([])
const unassignedTasks = ref<Task[]>([])
const isLoading = ref(true)
const error = ref('')
const claiming = ref<string | null>(null)
const showKpiPanel = ref(true)

const showSuccess = ref(false)
const successMessage = ref('')

const showSubmitModal = ref(false)
const showResultModal = ref(false)
const isSubmitting = ref(false)
const selectedTask = ref<Task | null>(null)
const submitForm = ref({ commitHash: '', diff: '' })
const submissionResult = ref<Submission | null>(null)

function deliveryStatus(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct >= 85) return 'good'
  if (pct >= 70) return 'warn'
  return 'bad'
}
function qualityStatus(q: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (q >= 75) return 'good'
  if (q >= 60) return 'warn'
  return 'bad'
}
function reworkStatus(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct <= 15) return 'good'
  if (pct <= 25) return 'warn'
  return 'bad'
}
function accuracyStatus(pct: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (pct >= 70) return 'good'
  if (pct >= 50) return 'warn'
  return 'bad'
}
function scoreStatus(score: number): 'good' | 'warn' | 'bad' | 'neutral' {
  if (score >= 80) return 'good'
  if (score >= 60) return 'warn'
  return 'neutral'
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
    .sort((a, b) => {
      if (!a.due_at) return 1
      if (!b.due_at) return -1
      return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
    })
)

const myBacklog = computed(() =>
  myTasks.value
    .filter(t => t.status === 'PENDING')
    .sort((a, b) => {
      if (!a.due_at) return 1
      if (!b.due_at) return -1
      return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
    })
)

const availableMissions = computed(() =>
  unassignedTasks.value.sort((a, b) => {
    if (!a.due_at) return 1
    if (!b.due_at) return -1
    return new Date(a.due_at).getTime() - new Date(b.due_at).getTime()
  })
)

const fetchTasks = async () => {
  try {
    isLoading.value = true
    const [myTasksResponse, unassignedResponse] = await Promise.all([
      fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/my'),
      fetchWithAuth<{ data: Task[] }>('/sentinel/tasks/unassigned')
    ])
    myTasks.value = myTasksResponse.data || []
    unassignedTasks.value = unassignedResponse.data || []
    error.value = ''
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to load tasks'
    console.error('Failed to fetch tasks:', err)
  } finally {
    isLoading.value = false
  }
}

const claimTask = async (taskId: string) => {
  if (!currentUser.value?.user_id) {
    error.value = 'User ID not found. Please log in again.'
    return
  }
  try {
    claiming.value = taskId
    await fetchWithAuth(`/sentinel/tasks/${taskId}/assign`, {
      method: 'POST',
      body: { dev_id: currentUser.value.user_id }
    })
    const claimedTask = unassignedTasks.value.find(t => t.id === taskId)
    if (claimedTask) {
      claimedTask.status = 'IN_PROGRESS'
      claimedTask.assigned_to = currentUser.value.user_id
      unassignedTasks.value = unassignedTasks.value.filter(t => t.id !== taskId)
      myTasks.value.unshift(claimedTask)
    }
    successMessage.value = 'Task claimed successfully'
    showSuccess.value = true
    setTimeout(() => { showSuccess.value = false }, 3000)
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to claim task'
    console.error('Failed to claim task:', err)
  } finally {
    claiming.value = null
  }
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

const getDeadlineUrgency = (task: Task) => {
  if (!task.due_at || task.status === 'COMPLETED') return 'none'
  const now = new Date().getTime()
  const dueDate = new Date(task.due_at).getTime()
  const hoursUntilDue = (dueDate - now) / (1000 * 60 * 60)
  if (hoursUntilDue < 0) return 'overdue'
  if (hoursUntilDue < 24) return 'urgent'
  return 'normal'
}

const getDeadlineBorderClass = (task: Task) => {
  const urgency = getDeadlineUrgency(task)
  if (urgency === 'overdue') return 'border-red-500'
  if (urgency === 'urgent') return 'border-yellow-500'
  return 'border-purple-500'
}

const formatDeadline = (dueAt: string) => {
  const date = new Date(dueAt)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

const getDeadlineCountdown = (dueAt: string) => {
  const now = new Date().getTime()
  const due = new Date(dueAt).getTime()
  const diff = due - now
  if (diff < 0) {
    const hours = Math.abs(Math.floor(diff / (1000 * 60 * 60)))
    const days = Math.floor(hours / 24)
    if (days > 0) return `Overdue by ${days}d`
    return `Overdue by ${hours}h`
  }
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(hours / 24)
  if (days > 0) return `${days}d left`
  if (hours > 0) return `${hours}h left`
  return 'Due soon!'
}

const openSubmitModal = (task: Task) => {
  selectedTask.value = task
  submitForm.value = { commitHash: '', diff: '' }
  showSubmitModal.value = true
}

const closeSubmitModal = () => {
  if (!isSubmitting.value) {
    showSubmitModal.value = false
    selectedTask.value = null
    submitForm.value = { commitHash: '', diff: '' }
  }
}

const submitWork = async () => {
  if (!selectedTask.value || !submitForm.value.commitHash || !submitForm.value.diff) return
  try {
    isSubmitting.value = true
    error.value = ''
    const response = await fetchWithAuth<{ data: Submission }>(`/sentinel/tasks/${selectedTask.value.id}/submit`, {
      method: 'POST',
      body: { commit_hash: submitForm.value.commitHash, diff: submitForm.value.diff }
    })
    submissionResult.value = response.data
    showSubmitModal.value = false
    showResultModal.value = true
    await fetchTasks()
  } catch (err: any) {
    error.value = err.data?.message || err.message || 'Failed to submit work'
    console.error('Failed to submit work:', err)
  } finally {
    isSubmitting.value = false
  }
}

const closeResultModal = async () => {
  if (submissionResult.value?.ai_verdict !== 'PASS' && selectedTask.value?.id) {
    const task = selectedTask.value
    showResultModal.value = false
    submissionResult.value = null
    selectedTask.value = null
    await navigateTo(`/task/${task?.code || task?.id}`)
    return
  }
  showResultModal.value = false
  submissionResult.value = null
  selectedTask.value = null
  fetchTasks()
}

const getFeedbackText = () => {
  if (!submissionResult.value?.ai_feedback) return 'No feedback available'
  try {
    const feedback = submissionResult.value.ai_feedback
    if (typeof feedback === 'string') return feedback
    if (feedback.error) {
      const errorMsg = feedback.error
      if (errorMsg.includes('429') || errorMsg.includes('quota') || errorMsg.includes('RESOURCE_EXHAUSTED')) {
        return 'AI review service temporarily unavailable (quota). Your submission is saved as PENDING. You can appeal or resubmit later.'
      }
      return `AI review error: ${errorMsg}. Your submission is saved as PENDING. You can appeal or resubmit.`
    }
    if (feedback.feedback) return feedback.feedback
    return JSON.stringify(feedback, null, 2)
  } catch {
    return 'Unable to parse feedback'
  }
}

const goToTask = (task: { id: string; code?: string }) => {
  navigateTo(`/task/${task?.code || task?.id}`)
}

onMounted(() => {
  fetchTasks()
  performanceStore.fetchAll('DEV')
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(0.5rem);
}
</style>
