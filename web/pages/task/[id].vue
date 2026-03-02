<template>
  <div class="min-h-screen bg-gray-900 p-6">
    <!-- Loading State -->
    <div v-if="isLoading" class="text-center py-20">
      <div class="animate-spin text-6xl mb-4">⚙️</div>
      <div class="text-gray-400">Loading mission intel...</div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-4xl mx-auto">
      <div class="bg-red-900/20 border border-red-500 rounded p-6 text-red-400">
        <h2 class="text-xl font-bold mb-2">Failed to Load Mission</h2>
        <p>{{ error }}</p>
        <button 
          @click="goToDashboard"
          class="mt-4 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded transition-colors"
        >
          ← Back to Dashboard
        </button>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="task">
      <!-- Header -->
      <div class="mb-6 border-b border-gray-700 pb-4">
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1">
            <h1 class="text-3xl font-bold text-white mb-2">{{ task.title }}</h1>
            <div class="flex items-center gap-4 text-sm text-gray-400">
              <span class="flex items-center gap-1">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
                <code class="text-xs">{{ task.id }}</code>
              </span>
              <span>•</span>
              <span>Created {{ formatDate(task.created_at) }}</span>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <span 
              :class="getStatusClass(task.status)"
              class="px-4 py-2 text-sm font-bold rounded"
            >
              {{ getStatusLabel(task.status) }}
            </span>
            
            <!-- 🚦 Approve & Complete Button (CEO/PM Only, REVIEW_PENDING tasks) -->
            <button
              v-if="canApproveTask"
              @click="approveTask"
              :disabled="isApprovingTask"
              class="px-6 py-3 bg-green-600 hover:bg-green-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-all flex items-center gap-2"
              title="Approve this task and mark as completed"
            >
              <span v-if="isApprovingTask" class="animate-spin">⚙️</span>
              <span v-else>✅</span>
              <span>{{ isApprovingTask ? 'Approving...' : 'Approve & Complete' }}</span>
            </button>
            
            <!-- Edit & Delete Buttons (CEO or Creator Only) -->
            <button
              v-if="canEditOrDelete"
              @click="openEditModal"
              class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors flex items-center gap-2"
              title="Edit Mission"
            >
              <span>✏️</span>
              <span>Edit</span>
            </button>
            
            <button
              v-if="canEditOrDelete"
              @click="openDeleteConfirmation"
              class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded transition-colors flex items-center gap-2"
              title="Delete Mission"
            >
              <span>🗑️</span>
              <span>Abort</span>
            </button>
            
            <button 
              @click="goToDashboard"
              class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded transition-colors flex items-center gap-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Back
            </button>
          </div>
        </div>
      </div>

      <!-- 🎉 Developer Feedback: Task Awaiting Approval -->
      <div 
        v-if="task.status === 'REVIEW_PENDING' && !isCeoOrPm"
        class="mb-6 bg-blue-900/20 border border-blue-500 rounded p-4"
      >
        <div class="flex items-center gap-3">
          <span class="text-2xl">🎉</span>
          <div class="flex-1">
            <div class="text-sm font-bold text-blue-400 mb-1">
              AI Security Checks Passed!
            </div>
            <div class="text-sm text-gray-300">
              Your code passed all AI security audits. The task is now awaiting PM/CEO verification for functionality review before being marked as complete. 
              <span class="text-blue-400 font-medium">Score: {{ getLatestSubmissionScore() }}/100</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 🔒 PM/CEO Info Banner: Tasks Ready for Review -->
      <div 
        v-if="task.status === 'REVIEW_PENDING' && isCeoOrPm"
        class="mb-6 bg-indigo-900/20 border border-indigo-500 rounded p-4"
      >
        <div class="flex items-center gap-3">
          <span class="text-2xl">🔍</span>
          <div class="flex-1">
            <div class="text-sm font-bold text-indigo-400 mb-1">
              Quality Gate: Awaiting Your Approval
            </div>
            <div class="text-sm text-gray-300">
              AI has approved this code (Score: <strong class="text-green-400">{{ getLatestSubmissionScore() }}/100</strong>). 
              Please verify functionality and approve to mark as COMPLETED.
            </div>
          </div>
        </div>
      </div>

      <!-- Grid Layout -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- LEFT COLUMN: Mission Intel (1/3) -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Card 1: Mission Brief -->
          <div class="bg-gray-800 border border-gray-700 rounded p-4">
            <h2 class="text-sm font-bold text-gray-400 uppercase mb-3">📋 Mission Brief</h2>
            <p class="text-gray-300 text-sm whitespace-pre-wrap leading-relaxed">{{ task.description || 'No description provided' }}</p>
          </div>

          <!-- Card 2: Personnel & Stats -->
          <div class="bg-gray-800 border border-gray-700 rounded p-4">
            <h2 class="text-sm font-bold text-gray-400 uppercase mb-4">📊 Mission Stats</h2>
            
            <!-- Assignee -->
            <div class="mb-4 pb-4 border-b border-gray-700">
              <div class="text-xs text-gray-500 mb-2">ASSIGNED TO</div>
              <div v-if="task.assigned_to" class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold">
                  {{ getDevInitial(task.assigned_to) }}
                </div>
                <div>
                  <div class="text-white font-medium">Dev #{{ task.assigned_to }}</div>
                  <div class="text-xs text-gray-500">Active</div>
                </div>
              </div>
              <div v-else class="text-orange-400 font-medium">
                🚨 UNASSIGNED
              </div>
            </div>

            <!-- AI Estimation & Time Negotiation -->
            <div class="mb-4">
              <div class="flex items-center justify-between mb-2">
                <div class="text-xs text-gray-500">AI ESTIMATED TIME</div>
                
                <!-- Negotiate Button (NONE status) -->
                <button
                  v-if="task.negotiation_status === 'NONE'"
                  @click="openNegotiateModal"
                  class="px-3 py-1 text-xs font-medium bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors"
                >
                  Request More Time
                </button>
              </div>
              
              <div class="text-2xl font-bold text-white">
                {{ (task.ai_estimated_minutes / 60).toFixed(1) }}h
              </div>
              <div class="text-xs text-gray-500 mb-2">{{ task.ai_estimated_minutes }} minutes</div>
              
              <!-- Negotiation Status Badges -->
              <!-- PENDING -->
              <div v-if="task.negotiation_status === 'PENDING'" 
                   class="mt-3 px-3 py-2 bg-yellow-900/30 border border-yellow-600 rounded flex items-center gap-2">
                <span class="text-lg">⏳</span>
                <div class="flex-1">
                  <div class="text-xs font-bold text-yellow-400">Asking for {{ (task.proposed_minutes / 60).toFixed(1) }}h</div>
                  <div class="text-xs text-yellow-200/70">Pending PM/CEO review</div>
                </div>
              </div>
              
              <!-- APPROVED -->
              <div v-if="task.negotiation_status === 'APPROVED'" 
                   class="mt-3 px-3 py-2 bg-green-900/30 border border-green-600 rounded flex items-center gap-2">
                <span class="text-lg">✅</span>
                <div class="flex-1">
                  <div class="text-xs font-bold text-green-400">Time Updated to {{ (task.proposed_minutes / 60).toFixed(1) }}h</div>
                  <div class="text-xs text-green-200/70">Request approved</div>
                </div>
              </div>
              
              <!-- REJECTED -->
              <div v-if="task.negotiation_status === 'REJECTED'" 
                   class="mt-3 px-3 py-2 bg-red-900/30 border border-red-600 rounded flex items-center gap-2">
                <span class="text-lg">❌</span>
                <div class="flex-1">
                  <div class="text-xs font-bold text-red-400">Request Denied</div>
                  <div class="text-xs text-red-200/70">AI estimate stands</div>
                </div>
              </div>
            </div>

            <!-- Created By -->
            <div>
              <div class="text-xs text-gray-500 mb-2">CREATED BY</div>
              <div class="text-white font-medium">Dev #{{ task.created_by }}</div>
            </div>
          </div>

          <!-- Card 3: Schedule Breakdown -->
          <div v-if="task.due_at || task.started_at || task.completed_at" class="bg-gray-800 border border-gray-700 rounded p-4">
            <h2 class="text-sm font-bold text-gray-400 uppercase mb-4">⏱️ Schedule Breakdown</h2>
            
            <div class="space-y-3">
              <!-- Assigned Date -->
              <div v-if="task.started_at" class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center text-blue-400 shrink-0">
                  📅
                </div>
                <div class="flex-1">
                  <div class="text-xs text-gray-500">ASSIGNED</div>
                  <div class="text-white font-medium text-sm">{{ formatDateTime(task.started_at) }}</div>
                </div>
              </div>

              <!-- Deadline -->
              <div v-if="task.due_at" class="flex items-start gap-3">
                <div 
                  :class="[
                    'w-8 h-8 rounded-full flex items-center justify-center shrink-0',
                    getDeadlineUrgency(task) === 'overdue' ? 'bg-red-900/30 text-red-400' :
                    getDeadlineUrgency(task) === 'urgent' ? 'bg-yellow-900/30 text-yellow-400' :
                    'bg-gray-700 text-gray-400'
                  ]"
                >
                  ⏳
                </div>
                <div class="flex-1">
                  <div class="text-xs text-gray-500">DEADLINE</div>
                  <div 
                    :class="[
                      'font-medium text-sm',
                      getDeadlineUrgency(task) === 'overdue' ? 'text-red-400' :
                      getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-400' :
                      'text-white'
                    ]"
                  >
                    {{ formatDateTime(task.due_at) }}
                  </div>
                  <div 
                    v-if="task.status !== 'COMPLETED'"
                    :class="[
                      'text-xs font-bold mt-1',
                      getDeadlineUrgency(task) === 'overdue' ? 'text-red-300' :
                      getDeadlineUrgency(task) === 'urgent' ? 'text-yellow-300' :
                      'text-gray-400'
                    ]"
                  >
                    {{ getDeadlineCountdown(task.due_at) }}
                  </div>
                </div>
              </div>

              <!-- Finished Date -->
              <div v-if="task.completed_at" class="flex items-start gap-3">
                <div class="w-8 h-8 rounded-full bg-green-500/20 flex items-center justify-center text-green-400 shrink-0">
                  🏁
                </div>
                <div class="flex-1">
                  <div class="text-xs text-gray-500">FINISHED</div>
                  <div class="text-green-400 font-medium text-sm">{{ formatDateTime(task.completed_at) }}</div>
                  <div v-if="task.started_at" class="text-xs text-gray-500 mt-1">
                    Duration: {{ calculateDuration(task.started_at, task.completed_at) }}
                  </div>
                </div>
              </div>

              <!-- Performance Badge -->
              <div 
                v-if="task.completed_at && task.started_at && task.due_at"
                :class="[
                  'mt-4 px-3 py-2 rounded border flex items-center gap-2',
                  new Date(task.completed_at) <= new Date(task.due_at)
                    ? 'bg-green-900/30 border-green-600'
                    : 'bg-red-900/30 border-red-600'
                ]"
              >
                <span class="text-lg">
                  {{ new Date(task.completed_at) <= new Date(task.due_at) ? '✅' : '⚠️' }}
                </span>
                <div class="flex-1">
                  <div 
                    :class="[
                      'text-xs font-bold',
                      new Date(task.completed_at) <= new Date(task.due_at)
                        ? 'text-green-400'
                        : 'text-red-400'
                    ]"
                  >
                    {{ new Date(task.completed_at) <= new Date(task.due_at)
                      ? 'Completed On Time'
                      : 'Completed Late' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- RIGHT COLUMN: The Black Box / Timeline (2/3) -->
        <div class="lg:col-span-2">
          <div class="bg-gray-800 border border-gray-700 rounded p-6">
            <h2 class="text-lg font-bold text-white mb-6">📜 Mission Log & Audit History</h2>

            <!-- Timeline -->
            <div class="relative">
              <!-- Vertical Line -->
              <div class="absolute left-6 top-0 bottom-0 w-0.5 bg-gray-700"></div>

              <!-- Events -->
              <div class="space-y-6">
                <!-- Submissions (Newest First) -->
                <div 
                  v-for="(submission, index) in sortedSubmissions" 
                  :key="submission.id"
                  class="relative pl-16"
                >
                  <!-- Timeline Dot (Gold if overridden) -->
                  <div 
                    :class="[
                      'absolute left-4 w-5 h-5 rounded-full border-4 border-gray-800',
                      getTimelineDotClass(submission)
                    ]"
                  ></div>

                  <!-- Event Card (Gold border if overridden) -->
                  <div 
                    :class="[
                      'border-l-4 rounded-lg p-5',
                      getSubmissionCardClass(submission)
                    ]"
                  >
                    <!-- Overridden Badge (Top Right) -->
                    <div v-if="submission.is_overridden" class="mb-3">
                      <span class="inline-flex items-center gap-2 px-3 py-1 bg-blue-900/50 border border-blue-500 rounded text-blue-400 text-xs font-bold">
                        👑 OVERRIDDEN BY HUMAN
                      </span>
                    </div>

                    <!-- Header -->
                    <div class="flex items-start justify-between mb-3">
                      <div>
                        <h3 
                          :class="[
                            'text-xl font-bold mb-1',
                            submission.is_overridden ? 'text-blue-400' : 
                            submission.ai_verdict === 'PASS' ? 'text-green-400' : 'text-red-400'
                          ]"
                        >
                          {{ submission.is_overridden ? '👑 VERDICT OVERRIDDEN' :
                             submission.ai_verdict === 'PASS' ? '✅ MISSION ACCOMPLISHED' : '❌ MISSION FAILED' }}
                        </h3>
                        <div class="text-sm text-gray-400">
                          {{ formatDateTime(submission.created_at) }}
                        </div>
                      </div>
                      <div 
                        :class="[
                          'text-4xl font-black',
                          submission.is_overridden ? 'text-blue-400' :
                          submission.ai_verdict === 'PASS' ? 'text-green-400' : 'text-red-400'
                        ]"
                      >
                        {{ submission.ai_score }}
                      </div>
                    </div>

                    <!-- Metrics -->
                    <div class="flex items-center gap-4 text-sm mb-4">
                      <div class="flex items-center gap-2">
                        <span class="text-gray-500">Commit:</span>
                        <code class="px-2 py-1 bg-gray-900 rounded text-gray-300 text-xs">
                          {{ submission.commit_hash.substring(0, 12) }}
                        </code>
                      </div>
                      <div class="flex items-center gap-2">
                        <span class="text-gray-500">Developer:</span>
                        <span class="text-gray-300">Dev #{{ submission.dev_id }}</span>
                      </div>
                    </div>

                    <!-- Appeal System UI -->
                    <div class="mt-4 space-y-3">
                      <!-- Appeal Pending Banner + Review Button (CEO/PM) -->
                      <div v-if="submission.appeal && submission.appeal.status === 'PENDING'">
                        <!-- Pending Banner -->
                        <div class="px-4 py-3 bg-yellow-900/30 border border-yellow-600 rounded flex items-center gap-3">
                          <span class="text-2xl">⚖️</span>
                          <div class="flex-1">
                            <div class="text-sm font-bold text-yellow-400">Appeal Under Review</div>
                            <div class="text-xs text-yellow-200/70 mt-1">
                              {{ isCeoOrPm ? 'This appeal requires your judgment.' : 'Your appeal is being reviewed by management. Decision pending.' }}
                            </div>
                          </div>
                        </div>
                        
                        <!-- Review Button (CEO/PM Only) -->
                        <button
                          v-if="canReviewAppeal(submission)"
                          @click="openAdjudicationModal(submission)"
                          class="w-full px-4 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
                        >
                          <span>⚖️</span>
                          <span>Review Appeal</span>
                        </button>
                      </div>

                      <!-- Appeal Rejected Banner -->
                      <div v-if="submission.appeal && submission.appeal.status === 'REJECTED'" 
                           class="px-4 py-3 bg-red-900/30 border border-red-600 rounded">
                        <div class="flex items-start gap-3">
                          <span class="text-2xl">❌</span>
                          <div class="flex-1">
                            <div class="text-sm font-bold text-red-400">Appeal Rejected</div>
                            <div v-if="submission.appeal.resolver_note" class="text-xs text-red-200/70 mt-1">
                              {{ submission.appeal.resolver_note }}
                            </div>
                          </div>
                        </div>
                      </div>

                      <!-- Appeal Button (For FAIL verdicts without appeal) -->
                      <button
                        v-if="canAppeal(submission)"
                        @click="openAppealModal(submission.id)"
                        class="px-4 py-2 bg-amber-600 hover:bg-amber-700 text-white text-sm font-medium rounded transition-colors flex items-center gap-2"
                      >
                        <span>📢</span>
                        <span>Appeal Verdict</span>
                      </button>

                      <!-- AI Feedback (Expandable) -->
                      <button
                        @click="toggleFeedback(submission.id)"
                        class="w-full flex items-center justify-between px-4 py-2 bg-gray-900/50 hover:bg-gray-900 border border-gray-700 rounded transition-colors"
                      >
                        <span class="text-sm font-medium text-gray-300">
                          🤖 AI Auditor Feedback
                        </span>
                        <svg 
                          :class="[
                            'w-4 h-4 text-gray-400 transition-transform',
                            expandedFeedback[submission.id] ? 'rotate-180' : ''
                          ]"
                          fill="none" 
                          stroke="currentColor" 
                          viewBox="0 0 24 24"
                        >
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>

                      <!-- Expanded Feedback -->
                      <div 
                        v-if="expandedFeedback[submission.id]"
                        class="mt-3 p-4 bg-gray-950 border border-gray-700 rounded text-xs text-gray-300 leading-relaxed whitespace-pre-wrap overflow-x-auto"
                      >
                        {{ getFeedbackText(submission.ai_feedback) }}
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Initialization Event (Bottom) -->
                <div class="relative pl-16">
                  <!-- Timeline Dot -->
                  <div class="absolute left-4 w-5 h-5 rounded-full border-4 border-gray-800 bg-purple-500"></div>

                  <!-- Event Card -->
                  <div class="bg-gray-900 border border-gray-700 rounded-lg p-5">
                    <h3 class="text-lg font-bold text-purple-400 mb-1">
                      🚀 Mission Initialized
                    </h3>
                    <div class="text-sm text-gray-400">
                      {{ formatDateTime(task.created_at) }}
                    </div>
                    <div class="mt-3 text-sm text-gray-500">
                      Task created by Dev #{{ task.created_by }} • AI estimated {{ (task.ai_estimated_minutes / 60).toFixed(1) }} hours
                    </div>
                  </div>
                </div>

                <!-- No Submissions Message -->
                <div 
                  v-if="!task.submissions || task.submissions.length === 0"
                  class="relative pl-16"
                >
                  <div class="absolute left-4 w-5 h-5 rounded-full border-4 border-gray-800 bg-gray-600"></div>
                  <div class="bg-gray-900 border border-gray-700 rounded-lg p-5 text-center">
                    <div class="text-gray-500 text-sm">
                      💤 No submissions yet. Awaiting developer action.
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Appeal Modal -->
    <div 
      v-if="showAppealModal"
      class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
      @click.self="closeAppealModal"
    >
      <div class="bg-gray-800 border border-gray-700 rounded max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="border-b border-gray-700 p-6">
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold text-white mb-1">📢 Appeal AI Verdict</h2>
              <p class="text-gray-400 text-sm">Challenge the AI Auditor's decision</p>
            </div>
            <button 
              @click="closeAppealModal"
              class="text-white hover:text-amber-200 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="p-6">
          <!-- Warning -->
          <div class="mb-6 p-4 bg-yellow-900/30 border border-yellow-600 rounded flex items-start gap-3">
            <span class="text-2xl">⚠️</span>
            <div class="flex-1 text-sm text-yellow-100">
              <div class="font-bold mb-1">Important Notice</div>
              <div class="text-yellow-200/80">
                Your appeal will be reviewed by management. Please provide a clear, professional explanation 
                of why you believe the AI verdict is incorrect. Be specific about technical context that the AI may have missed.
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="appealError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
            {{ appealError }}
          </div>

          <!-- Form -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Reason for Appeal <span class="text-red-400">*</span>
              </label>
              <textarea
                v-model="appealForm.reason"
                rows="8"
                placeholder="Explain why you believe the AI verdict is incorrect. Include:
• Technical context the AI may have missed
• Justification for your approach
• Evidence supporting your implementation
• Specific points where the AI's assessment was wrong"
                class="w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:border-transparent"
                :disabled="isSubmittingAppeal"
              ></textarea>
              <div class="mt-2 text-xs text-gray-500">
                Minimum 50 characters recommended for a strong appeal
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="mt-6 flex gap-3">
            <button
              @click="submitAppeal"
              :disabled="isSubmittingAppeal || !appealForm.reason.trim()"
              class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
            >
              <span v-if="isSubmittingAppeal" class="animate-spin">⚙️</span>
              <span v-else>⚖️</span>
              <span>{{ isSubmittingAppeal ? 'Submitting Appeal...' : 'Submit Appeal' }}</span>
            </button>
            <button
              @click="closeAppealModal"
              :disabled="isSubmittingAppeal"
              class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Negotiate Time Modal -->
    <div 
      v-if="showNegotiateModal"
      class="fixed inset-0 bg-black/80 flex items-center justify-center z-50 p-4"
      @click.self="closeNegotiateModal"
    >
      <div class="bg-gray-800 border border-gray-700 rounded max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="border-b border-gray-700 p-6">
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold text-white mb-1">⏱️ Dispute AI Estimate</h2>
              <p class="text-gray-400 text-sm">Request more realistic time allocation</p>
            </div>
            <button 
              @click="closeNegotiateModal"
              class="text-white hover:text-blue-200 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="p-6">
          <!-- Info Banner -->
          <div class="mb-6 p-4 bg-blue-900/30 border border-blue-600 rounded flex items-start gap-3">
            <span class="text-2xl">💡</span>
            <div class="flex-1 text-sm text-blue-100">
              <div class="font-bold mb-1">Reality Check</div>
              <div class="text-blue-200/80">
                The AI estimated <strong>{{ (task.ai_estimated_minutes / 60).toFixed(1) }} hours</strong> for this task. 
                If this doesn't match reality (legacy code, unclear requirements, tech debt), explain why you need more time. 
                PM/CEO will review your request.
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="negotiateError" class="mb-4 p-3 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
            {{ negotiateError }}
          </div>

          <!-- Form -->
          <div class="space-y-4">
            <!-- Proposed Time -->
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Proposed Time (Minutes) <span class="text-red-400">*</span>
              </label>
              <input
                v-model.number="negotiateForm.minutes"
                type="number"
                min="1"
                placeholder="e.g., 120 (2 hours)"
                class="w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :disabled="isSubmittingNegotiation"
              />
              <div class="mt-2 text-xs text-gray-500">
                Current AI estimate: {{ task.ai_estimated_minutes }} minutes ({{ (task.ai_estimated_minutes / 60).toFixed(1) }}h)
              </div>
              <div v-if="negotiateForm.minutes > 0" class="mt-1 text-xs text-blue-400">
                Your request: {{ negotiateForm.minutes }} minutes ({{ (negotiateForm.minutes / 60).toFixed(1) }}h)
              </div>
            </div>

            <!-- Reason -->
            <div>
              <label class="block text-sm font-medium text-gray-300 mb-2">
                Reason <span class="text-red-400">*</span>
              </label>
              <textarea
                v-model="negotiateForm.reason"
                rows="6"
                minlength="20"
                placeholder="Explain why the AI estimate is unrealistic:
• Legacy code complexity
• Unclear requirements
• Technical debt
• Security considerations
• Testing requirements
• Documentation needs"
                :class="[
                  'w-full px-4 py-3 bg-gray-800 border rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:border-transparent',
                  negotiateForm.reason.length > 0 && negotiateForm.reason.length < 20
                    ? 'border-red-500 focus:ring-red-500'
                    : negotiateForm.reason.length >= 20
                    ? 'border-green-500 focus:ring-green-500'
                    : 'border-gray-700 focus:ring-blue-500'
                ]"
                :disabled="isSubmittingNegotiation"
              ></textarea>
              <div class="mt-2 flex items-center justify-between text-xs">
                <span :class="negotiateForm.reason.length < 20 ? 'text-red-400' : 'text-green-400'">
                  {{ negotiateForm.reason.length }} / 20 characters minimum
                </span>
                <span v-if="negotiateForm.reason.length < 20 && negotiateForm.reason.length > 0" class="text-red-400">
                  {{ 20 - negotiateForm.reason.length }} more needed
                </span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="mt-6 flex gap-3">
            <button
              @click="submitNegotiation"
              :disabled="isSubmittingNegotiation || negotiateForm.reason.trim().length < 20 || negotiateForm.minutes <= 0"
              class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
            >
              <span v-if="isSubmittingNegotiation" class="animate-spin">⚙️</span>
              <span v-else>⏱️</span>
              <span>{{ isSubmittingNegotiation ? 'Submitting Request...' : 'Submit Negotiation' }}</span>
            </button>
            <button
              @click="closeNegotiateModal"
              :disabled="isSubmittingNegotiation"
              class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded transition-colors"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Adjudication Modal (CEO/PM Only) -->
    <div 
      v-if="showAdjudicationModal"
      class="fixed inset-0 bg-black/90 flex items-center justify-center z-50 p-4"
      @click.self="closeAdjudicationModal"
    >
      <div class="bg-gray-800 border border-gray-700 rounded max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="border-b border-gray-700 p-6">
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-2xl font-bold text-white mb-2 flex items-center gap-3">
                <span>⚖️</span>
                <span>Appeal Review</span>
              </h2>
              <p class="text-gray-400 text-sm">
                {{ authStore.user?.role }} {{ authStore.user?.email }} presiding
              </p>
            </div>
            <button 
              @click="closeAdjudicationModal"
              class="text-white hover:text-amber-200 transition-colors"
              :disabled="isResolvingAppeal"
            >
              <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="p-6 space-y-6">
          <!-- Case Header -->
          <div class="bg-gray-800 border-2 border-gray-700 rounded-lg p-5">
            <div class="flex items-start justify-between mb-4">
              <div>
                <div class="text-xs text-gray-500 uppercase tracking-wide mb-1">Case #</div>
                <code class="text-sm text-gray-300">{{ adjudicationForm.appealId }}</code>
              </div>
              <div class="text-right">
                <div class="text-xs text-gray-500 uppercase tracking-wide mb-1">Submission</div>
                <code class="text-sm text-gray-300">{{ adjudicationForm.submissionId.substring(0, 12) }}</code>
              </div>
            </div>
            
            <!-- Appellant Info -->
            <div class="pt-4 border-t border-gray-700">
              <div class="text-xs text-gray-500 uppercase tracking-wide mb-2">Appellant</div>
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center text-white font-bold">
                  {{ adjudicationForm.appellantName.substring(4, 5) }}
                </div>
                <div>
                  <div class="text-white font-bold">{{ adjudicationForm.appellantName }}</div>
                  <div class="text-xs text-gray-500">Developer</div>
                </div>
              </div>
            </div>
          </div>

          <!-- The Plea -->
          <div class="bg-gray-900 border border-gray-700 rounded p-4">
            <div class="flex items-center gap-2 mb-3">
              <span class="text-xl">📜</span>
              <h3 class="text-sm font-bold text-gray-400 uppercase">Developer's Plea</h3>
            </div>
            <blockquote class="border-l-4 border-gray-600 pl-4 text-gray-300 text-sm leading-relaxed">
              "{{ adjudicationForm.appealReason }}"
            </blockquote>
          </div>

          <!-- AI Advisor Opinion -->
          <div class="bg-gray-900 border border-gray-700 rounded p-4">
            <div class="flex items-center gap-3 mb-4">
              <span class="text-2xl">🤖</span>
              <div>
                <h3 class="text-sm font-bold text-gray-400 uppercase">AI Advisor Opinion</h3>
              </div>
            </div>

            <!-- Recommendation -->
            <div class="mb-4 flex items-center justify-between">
              <span class="text-sm text-gray-400">AI Suggests:</span>
              <span 
                :class="[
                  'px-3 py-1 text-sm font-bold rounded',
                  adjudicationForm.aiRecommendation === 'OVERTURN'
                    ? 'bg-green-700 text-green-100'
                    : 'bg-red-700 text-red-100'
                ]"
              >
                {{ adjudicationForm.aiRecommendation === 'OVERTURN' ? '✅ APPROVE' : '❌ REJECT' }}
              </span>
            </div>

            <!-- Confidence -->
            <div class="mb-4 flex items-center justify-between">
              <span class="text-sm text-gray-400">Confidence:</span>
              <span class="text-sm font-bold text-white">{{ adjudicationForm.aiConfidence }}%</span>
            </div>

            <!-- AI Reasoning -->
            <div class="bg-gray-950 border border-gray-700 rounded p-3">
              <div class="text-xs text-gray-500 uppercase mb-2">Analysis</div>
              <div class="text-gray-300 text-sm leading-relaxed">
                {{ adjudicationForm.aiReasoning }}
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="adjudicationError" class="p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
            {{ adjudicationError }}
          </div>

          <!-- Resolver Note (Optional) -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📝</span>
                <span>Your Verdict Statement</span>
                <span class="text-gray-500 text-xs">(Optional)</span>
              </span>
            </label>
            <textarea
              v-model="adjudicationForm.resolverNote"
              rows="4"
              placeholder="Optional: Explain your decision. This will be visible to the developer.
Example: 'After careful review, the AI was correct about the security risk.' or 'The developer's argument is valid - this is a false positive.'"
              class="w-full px-4 py-3 bg-gray-800 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              :disabled="isResolvingAppeal"
            ></textarea>
          </div>

          <!-- Verdict Actions -->
          <div class="grid grid-cols-2 gap-4 pt-4 border-t border-gray-700">
            <!-- APPROVE -->
            <button
              @click="resolveAppeal('APPROVED')"
              :disabled="isResolvingAppeal"
              class="px-6 py-3 bg-green-600 hover:bg-green-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors"
            >
              <div class="flex flex-col items-center gap-1">
                <span v-if="isResolvingAppeal" class="animate-spin">⚙️</span>
                <span v-else>✅</span>
                <span class="text-sm">Approve Appeal</span>
              </div>
            </button>

            <!-- REJECT -->
            <button
              @click="resolveAppeal('REJECTED')"
              :disabled="isResolvingAppeal"
              class="px-6 py-3 bg-red-600 hover:bg-red-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors"
            >
              <div class="flex flex-col items-center gap-1">
                <span v-if="isResolvingAppeal" class="animate-spin">⚙️</span>
                <span v-else>❌</span>
                <span class="text-sm">Reject Appeal</span>
              </div>
            </button>
          </div>

          <!-- Cancel Button -->
          <button
            @click="closeAdjudicationModal"
            :disabled="isResolvingAppeal"
            class="w-full px-4 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded transition-colors"
          >
            Cancel Review
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Mission Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeEditModal"
    >
      <div class="bg-gray-800 border-2 border-blue-600 rounded-lg p-6 max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="flex items-center justify-between mb-6 pb-4 border-b-2 border-gray-700">
          <div>
            <h2 class="text-2xl font-bold text-white flex items-center gap-3">
              <span class="text-3xl">✏️</span>
              <span>Edit Mission</span>
            </h2>
            <div class="text-sm text-gray-400 mt-1">Update mission parameters</div>
          </div>
          <button
            @click="closeEditModal"
            class="text-gray-400 hover:text-white transition-colors text-2xl"
            :disabled="isUpdatingTask"
          >
            ✕
          </button>
        </div>

        <!-- AI Re-estimation Warning -->
        <div class="mb-6 p-4 bg-yellow-900/30 border-2 border-yellow-600/50 rounded-lg">
          <div class="flex items-start gap-3">
            <span class="text-2xl">⚠️</span>
            <div>
              <div class="font-bold text-yellow-400 mb-1">AI Re-estimation Alert</div>
              <div class="text-sm text-yellow-300/90">
                Changing the title or description will trigger <strong>automatic AI re-estimation</strong>. 
                Any pending time negotiation will be reset.
              </div>
            </div>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="editError" class="mb-4 p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
          {{ editError }}
        </div>

        <!-- Form -->
        <div class="space-y-5">
          <!-- Title -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📋</span>
                <span>Mission Title</span>
              </span>
            </label>
            <input
              v-model="editForm.title"
              type="text"
              placeholder="Enter mission title..."
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            />
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📝</span>
                <span>Mission Description</span>
              </span>
            </label>
            <textarea
              v-model="editForm.description"
              rows="6"
              placeholder="Describe the mission objectives and requirements..."
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            ></textarea>
          </div>

          <!-- Deadline (Optional) -->
          <div>
            <label class="block text-sm font-medium text-gray-300 mb-2">
              <span class="flex items-center gap-2">
                <span>📅</span>
                <span>Deadline</span>
                <span class="text-gray-500 text-xs">(Optional)</span>
              </span>
            </label>
            <input
              v-model="editForm.deadline"
              type="datetime-local"
              class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="isUpdatingTask"
            />
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4 mt-6 pt-4 border-t-2 border-gray-700">
          <button
            @click="submitEdit"
            :disabled="isUpdatingTask || !editForm.title.trim()"
            class="flex-1 px-6 py-3 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
          >
            <span v-if="isUpdatingTask" class="animate-spin">⚙️</span>
            <span v-else>💾</span>
            <span>{{ isUpdatingTask ? 'Updating...' : 'Update Mission' }}</span>
          </button>
          
          <button
            @click="closeEditModal"
            :disabled="isUpdatingTask"
            class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded-lg transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteModal"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4"
      @click.self="closeDeleteModal"
    >
      <div class="bg-gray-800 border-2 border-red-600 rounded-lg p-6 max-w-lg w-full relative">
        <!-- Close Button (X) -->
        <button
          @click="closeDeleteModal"
          class="absolute top-4 right-4 text-gray-400 hover:text-white transition-colors"
          :disabled="isDeletingTask"
          title="Close"
        >
          <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>

        <!-- Header -->
        <div class="flex items-center gap-3 mb-6">
          <span class="text-4xl">🗑️</span>
          <div>
            <h2 class="text-2xl font-bold text-white">Abort Mission?</h2>
            <div class="text-sm text-gray-400 mt-1">This action cannot be undone</div>
          </div>
        </div>

        <!-- Warning -->
        <div class="mb-6 p-4 bg-red-900/30 border-2 border-red-600/50 rounded-lg">
          <div class="text-sm text-red-300 leading-relaxed">
            <p class="font-bold mb-2">⚠️ Critical Operation</p>
            <p>Are you sure you want to <strong>permanently delete</strong> this mission?</p>
            <p class="mt-2">This will remove:</p>
            <ul class="list-disc list-inside ml-2 mt-1">
              <li>Mission data</li>
              <li>All submissions</li>
              <li>All appeals</li>
              <li>Complete audit trail</li>
            </ul>
          </div>
        </div>

        <!-- Task Info -->
        <div class="mb-6 p-4 bg-gray-900 border border-gray-700 rounded">
          <div class="text-xs text-gray-500 uppercase mb-1">Mission to Delete:</div>
          <div class="font-bold text-white">{{ task?.title }}</div>
          <div class="text-xs text-gray-400 mt-1">ID: {{ task?.id }}</div>
        </div>

        <!-- Error Message -->
        <div v-if="deleteError" class="mb-4 p-4 bg-red-900/30 border border-red-600 rounded text-red-400 text-sm">
          {{ deleteError }}
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4">
          <button
            @click="confirmDelete"
            :disabled="isDeletingTask"
            class="flex-1 px-6 py-3 bg-red-600 hover:bg-red-700 disabled:bg-gray-600 disabled:cursor-not-allowed text-white font-bold rounded transition-colors flex items-center justify-center gap-2"
          >
            <span v-if="isDeletingTask" class="animate-spin">⚙️</span>
            <span v-else>💥</span>
            <span>{{ isDeletingTask ? 'Deleting...' : 'Yes, Delete Forever' }}</span>
          </button>
          
          <button
            @click="closeDeleteModal"
            :disabled="isDeletingTask"
            class="px-6 py-3 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:cursor-not-allowed text-gray-300 font-medium rounded-lg transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

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
  ai_reasoning: string      // Advice for CEO/PM
  
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

interface Task {
  id: string
  title: string
  description: string
  resource_urls: any
  ai_estimated_minutes: number
  
  // Time Negotiation
  negotiation_status: string // NONE, PENDING, APPROVED, REJECTED
  proposed_minutes: number
  negotiation_reason: string
  
  due_at: string | null
  started_at: string | null
  completed_at: string | null
  status: string
  assigned_to: number | null
  created_by: number
  created_at: string
  updated_at: string
  submissions?: Submission[]
}

const route = useRoute()
const { fetchWithAuth } = useAuth()
const authStore = useAuthStore()

// State
const task = ref<Task | null>(null)
const isLoading = ref(true)
const error = ref('')
const expandedFeedback = ref<Record<string, boolean>>({})

// Appeal System State
const showAppealModal = ref(false)
const appealForm = ref({
  submissionId: '',
  reason: ''
})
const isSubmittingAppeal = ref(false)
const appealError = ref('')

// Time Negotiation State
const showNegotiateModal = ref(false)
const negotiateForm = ref({
  minutes: 0,
  reason: ''
})
const isSubmittingNegotiation = ref(false)
const negotiateError = ref('')

// Appeal Review State (CEO/PM Only)
const showAdjudicationModal = ref(false)
const adjudicationForm = ref({
  appealId: '',
  submissionId: '',
  appellantName: '',
  appealReason: '',
  resolverNote: '',
  
  // AI Advisory
  aiRecommendation: '',  // OVERTURN or UPHOLD
  aiConfidence: 0,       // 0-100
  aiReasoning: ''        // AI advice text
})
const isResolvingAppeal = ref(false)
const adjudicationError = ref('')

// Edit Task State
const showEditModal = ref(false)
const editForm = ref({
  title: '',
  description: '',
  deadline: ''
})
const isUpdatingTask = ref(false)
const editError = ref('')

// Delete Task State
const showDeleteModal = ref(false)
const isDeletingTask = ref(false)
const deleteError = ref('')

// Task Approval State (Quality Gate)
const isApprovingTask = ref(false)
const approvalError = ref('')

// Computed
const sortedSubmissions = computed(() => {
  if (!task.value?.submissions) return []
  // Sort by created_at descending (newest first)
  return [...task.value.submissions].sort((a, b) => 
    new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  )
})

const isCeoOrPm = computed(() => {
  return authStore.user?.role === 'CEO' || authStore.user?.role === 'PM'
})

const canEditOrDelete = computed(() => {
  if (!task.value || !authStore.user) return false
  
  // CEO can edit/delete any task
  if (authStore.user.role === 'CEO') return true
  
  // Creator can edit/delete their own task
  return task.value.created_by === authStore.user.id
})

const canApproveTask = computed(() => {
  if (!task.value || !authStore.user) return false
  
  // Only CEO or PM can approve
  if (authStore.user.role !== 'CEO' && authStore.user.role !== 'PM') return false
  
  // Task must be in REVIEW_PENDING status
  return task.value.status === 'REVIEW_PENDING'
})

// Methods
const fetchTask = async () => {
  const taskId = route.params.id as string
  
  if (!taskId) {
    error.value = 'Invalid task ID'
    isLoading.value = false
    return
  }

  try {
    isLoading.value = true
    error.value = ''
    
    const response = await fetchWithAuth<{ data: Task }>(`/sentinel/tasks/${taskId}`)
    task.value = response.data
  } catch (err: any) {
    console.error('Failed to fetch task:', err)
    error.value = err.data?.message || err.message || 'Failed to load task'
  } finally {
    isLoading.value = false
  }
}

const toggleFeedback = (submissionId: string) => {
  expandedFeedback.value[submissionId] = !expandedFeedback.value[submissionId]
}

const getFeedbackText = (feedback: any): string => {
  if (!feedback) return 'No feedback available'
  
  // If feedback is a string, return it directly
  if (typeof feedback === 'string') {
    try {
      const parsed = JSON.parse(feedback)
      return parsed.feedback || feedback
    } catch {
      return feedback
    }
  }
  
  // If feedback is an object
  if (typeof feedback === 'object') {
    return feedback.feedback || JSON.stringify(feedback, null, 2)
  }
  
  return 'No feedback available'
}

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    'COMPLETED': 'bg-green-700 text-green-100 border border-green-500',
    'IN_PROGRESS': 'bg-blue-700 text-blue-100 border border-blue-500',
    'PENDING': 'bg-yellow-700 text-yellow-100 border border-yellow-500',
    'BLOCKED': 'bg-red-700 text-red-100 border border-red-500',
    'REVIEW_PENDING': 'bg-indigo-900 text-indigo-200 border border-indigo-600' // 🚦 Quality Gate
  }
  return classes[status] || 'bg-gray-700 text-gray-100 border border-gray-500'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'COMPLETED': '✅ COMPLETED',
    'IN_PROGRESS': '🔄 IN PROGRESS',
    'PENDING': '⏳ PENDING',
    'BLOCKED': '🚫 BLOCKED',
    'REVIEW_PENDING': '⏳ WAITING FOR APPROVAL', // 🚦 Quality Gate
    'ASSIGNED': '📌 ASSIGNED'
  }
  return labels[status] || status
}

const getDevInitial = (devId: number) => {
  return `D${devId}`
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

// Appeal System Methods
const openAppealModal = (submissionId: string) => {
  appealForm.value = {
    submissionId,
    reason: ''
  }
  appealError.value = ''
  showAppealModal.value = true
}

const closeAppealModal = () => {
  showAppealModal.value = false
  appealForm.value = {
    submissionId: '',
    reason: ''
  }
  appealError.value = ''
}

const submitAppeal = async () => {
  if (!appealForm.value.reason.trim()) {
    appealError.value = 'Please provide a reason for your appeal'
    return
  }

  try {
    isSubmittingAppeal.value = true
    appealError.value = ''

    await fetchWithAuth(`/sentinel/submissions/${appealForm.value.submissionId}/appeal`, {
      method: 'POST',
      body: JSON.stringify({
        reason: appealForm.value.reason
      })
    })

    // Success! Close modal and refresh task
    closeAppealModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to submit appeal:', err)
    appealError.value = err.data?.message || err.message || 'Failed to submit appeal'
  } finally {
    isSubmittingAppeal.value = false
  }
}

const canAppeal = (submission: Submission): boolean => {
  // Can appeal if: NOT PASS (FAIL or PENDING) AND no existing appeal
  // PENDING = AI service unavailable (quota/error), FAIL = AI rejected code
  return submission.ai_verdict !== 'PASS' && !submission.appeal
}

// Time Negotiation Methods
const openNegotiateModal = () => {
  if (!task.value) return
  
  negotiateForm.value = {
    minutes: task.value.ai_estimated_minutes * 2, // Default: 2x AI estimate
    reason: ''
  }
  negotiateError.value = ''
  showNegotiateModal.value = true
}

const closeNegotiateModal = () => {
  showNegotiateModal.value = false
  negotiateForm.value = {
    minutes: 0,
    reason: ''
  }
  negotiateError.value = ''
}

const submitNegotiation = async () => {
  if (!task.value) return
  
  const reason = negotiateForm.value.reason.trim()
  
  if (!reason) {
    negotiateError.value = 'Please provide a reason for your negotiation'
    return
  }

  if (reason.length < 20) {
    negotiateError.value = 'Reason must be at least 20 characters long. Please provide more details.'
    return
  }

  if (negotiateForm.value.minutes <= 0) {
    negotiateError.value = 'Proposed time must be greater than 0'
    return
  }

  if (negotiateForm.value.minutes <= task.value.ai_estimated_minutes) {
    negotiateError.value = `Proposed time must be greater than AI estimate (${task.value.ai_estimated_minutes} minutes)`
    return
  }

  try {
    isSubmittingNegotiation.value = true
    negotiateError.value = ''

    await fetchWithAuth(`/sentinel/tasks/${task.value.id}/negotiate`, {
      method: 'POST',
      body: JSON.stringify({
        minutes: negotiateForm.value.minutes,
        reason: negotiateForm.value.reason
      })
    })

    // Success! Close modal and refresh task
    closeNegotiateModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to submit negotiation:', err)
    negotiateError.value = err.data?.message || err.message || 'Failed to submit negotiation'
  } finally {
    isSubmittingNegotiation.value = false
  }
}

const getSubmissionCardClass = (submission: Submission): string => {
  // Blue border if overridden by human
  if (submission.is_overridden) {
    return 'bg-blue-900/20 border-blue-500'
  }
  
  // Default: Green for PASS, Red for FAIL
  return submission.ai_verdict === 'PASS' 
    ? 'bg-green-900/20 border-green-500' 
    : 'bg-red-900/20 border-red-500'
}

const getTimelineDotClass = (submission: Submission): string => {
  if (submission.is_overridden) {
    return 'bg-blue-500'
  }
  return submission.ai_verdict === 'PASS' ? 'bg-green-500' : 'bg-red-500'
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

// Appeal Review Methods (CEO/PM Only)
const canReviewAppeal = (submission: Submission): boolean => {
  return isCeoOrPm.value && 
         submission.appeal?.status === 'PENDING'
}

const openAdjudicationModal = (submission: Submission) => {
  if (!submission.appeal) return
  
  adjudicationForm.value = {
    appealId: submission.appeal.id,
    submissionId: submission.id,
    appellantName: `Dev #${submission.dev_id}`,
    appealReason: submission.appeal.reason,
    resolverNote: '',
    
    // AI Advisory
    aiRecommendation: submission.appeal.ai_recommendation || 'UPHOLD',
    aiConfidence: submission.appeal.ai_confidence || 0,
    aiReasoning: submission.appeal.ai_reasoning || 'AI analysis unavailable'
  }
  adjudicationError.value = ''
  showAdjudicationModal.value = true
}

const closeAdjudicationModal = () => {
  showAdjudicationModal.value = false
  adjudicationForm.value = {
    appealId: '',
    submissionId: '',
    appellantName: '',
    appealReason: '',
    resolverNote: '',
    aiRecommendation: '',
    aiConfidence: 0,
    aiReasoning: ''
  }
  adjudicationError.value = ''
}

const applyAIRecommendation = () => {
  // Apply AI recommendation automatically
  if (adjudicationForm.value.aiRecommendation === 'OVERTURN') {
    // AI suggests approving the appeal
    if (!adjudicationForm.value.resolverNote) {
      adjudicationForm.value.resolverNote = `Following AI recommendation (${adjudicationForm.value.aiConfidence}% confidence): ${adjudicationForm.value.aiReasoning}`
    }
    resolveAppeal('APPROVED')
  } else {
    // AI suggests rejecting the appeal (UPHOLD)
    if (!adjudicationForm.value.resolverNote) {
      adjudicationForm.value.resolverNote = `Following AI recommendation (${adjudicationForm.value.aiConfidence}% confidence): ${adjudicationForm.value.aiReasoning}`
    }
    resolveAppeal('REJECTED')
  }
}

const resolveAppeal = async (status: 'APPROVED' | 'REJECTED') => {
  if (!adjudicationForm.value.appealId) return
  
  try {
    isResolvingAppeal.value = true
    adjudicationError.value = ''

    await fetchWithAuth(`/sentinel/appeals/${adjudicationForm.value.appealId}/resolve`, {
      method: 'POST',
      body: JSON.stringify({
        status: status,
        note: adjudicationForm.value.resolverNote || `Appeal ${status.toLowerCase()} by ${authStore.user?.role}`
      })
    })

    // Success! Close modal and refresh task
    closeAdjudicationModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to resolve appeal:', err)
    adjudicationError.value = err.data?.message || err.message || 'Failed to resolve appeal'
  } finally {
    isResolvingAppeal.value = false
  }
}

// Edit Task Methods
const openEditModal = () => {
  if (!task.value) return
  
  // Pre-fill form with current values
  editForm.value.title = task.value.title
  editForm.value.description = task.value.description || ''
  
  // Convert due_at to datetime-local format if exists
  if (task.value.due_at) {
    const date = new Date(task.value.due_at)
    // Format: YYYY-MM-DDTHH:mm
    editForm.value.deadline = date.toISOString().slice(0, 16)
  } else {
    editForm.value.deadline = ''
  }
  
  editError.value = ''
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.value = {
    title: '',
    description: '',
    deadline: ''
  }
  editError.value = ''
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
    
    // Prepare request body (only send changed fields)
    const body: any = {}
    
    if (editForm.value.title !== task.value.title) {
      body.title = editForm.value.title
    }
    
    if (editForm.value.description !== (task.value.description || '')) {
      body.description = editForm.value.description
    }
    
    // Check if at least one field is being updated
    if (Object.keys(body).length === 0) {
      editError.value = 'No changes detected. Please modify at least one field.'
      isUpdatingTask.value = false
      return
    }
    
    // Call PATCH API
    await fetchWithAuth(`/sentinel/tasks/${taskId}`, {
      method: 'PATCH',
      body: JSON.stringify(body)
    })
    
    // Success! Show toast notification
    // Note: You can add a toast library like vue-toastification for better UX
    alert('✅ Mission Updated!\n\n🤖 AI is recalculating time estimate...')
    
    // Close modal and refresh task data
    closeEditModal()
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to update task:', err)
    editError.value = err.data?.message || err.message || 'Failed to update mission'
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
    
    // Success! Show notification and redirect
    alert('💥 Mission Deleted Successfully!\n\nReturning to dashboard...')
    
    // Redirect to dashboard
    goToDashboard()
  } catch (err: any) {
    console.error('Failed to delete task:', err)
    deleteError.value = err.data?.message || err.message || 'Failed to delete mission'
  } finally {
    isDeletingTask.value = false
  }
}

// Navigation helper
const goToDashboard = () => {
  navigateTo('/dashboard')
}

// Helper: Get latest submission score for banners
const getLatestSubmissionScore = (): number => {
  if (!task.value?.submissions || task.value.submissions.length === 0) return 0
  
  // Get the most recent submission (already sorted)
  const latest = sortedSubmissions.value[0]
  return latest?.ai_score || 0
}

// 🚦 Quality Gate: Approve Task (PM/CEO only)
const approveTask = async () => {
  if (!task.value) return
  
  try {
    isApprovingTask.value = true
    approvalError.value = ''
    
    const taskId = route.params.id as string
    
    // Call approval endpoint
    await fetchWithAuth(`/sentinel/tasks/${taskId}/approve`, {
      method: 'POST'
    })
    
    // Success! Show celebration and refresh
    alert('🎉 Task Approved & Completed!\n\n✅ Mission marked as COMPLETED successfully.')
    
    // Refresh task data to show new status
    await fetchTask()
  } catch (err: any) {
    console.error('Failed to approve task:', err)
    approvalError.value = err.data?.message || err.message || 'Failed to approve task'
    
    // Show error to user
    alert(`❌ Approval Failed\n\n${approvalError.value}`)
  } finally {
    isApprovingTask.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchTask()
})
</script>
