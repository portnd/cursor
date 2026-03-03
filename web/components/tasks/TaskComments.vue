<template>
  <div class="task-comments">
    <!-- Comments List -->
    <div class="space-y-4 mb-6 max-h-96 overflow-y-auto">
      <div v-if="!comments.length" class="text-center py-8 text-gray-500 text-sm">
        No comments yet. Start the discussion.
      </div>
      <div
        v-for="comment in comments"
        :key="comment.id"
        class="flex gap-3"
      >
        <div class="flex-shrink-0 w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-white text-sm font-bold">
          {{ (comment.user_email || String(comment.user_id)).charAt(0).toUpperCase() }}
        </div>
        <div class="flex-1 bg-gray-800 rounded-xl px-4 py-3 border border-gray-700/50">
          <div class="flex items-center justify-between mb-1.5">
            <span class="text-sm font-medium text-gray-300">{{ comment.user_email || `User #${comment.user_id}` }}</span>
            <span class="text-xs text-gray-500">{{ formatTime(comment.created_at) }}</span>
          </div>
          <p class="text-sm text-gray-300 leading-relaxed whitespace-pre-wrap">{{ comment.content }}</p>
        </div>
      </div>
    </div>

    <!-- Add Comment Input -->
    <div class="flex gap-3">
      <div class="flex-shrink-0 w-8 h-8 rounded-full bg-gray-700 flex items-center justify-center text-gray-400 text-sm">
        You
      </div>
      <div class="flex-1 flex gap-2">
        <textarea
          v-model="newComment"
          placeholder="Write a comment... (Shift+Enter for new line)"
          rows="2"
          class="flex-1 bg-gray-800 border border-gray-700 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-indigo-500 resize-none"
          @keydown.enter.prevent="handleEnter"
        ></textarea>
        <button
          @click="submitComment"
          :disabled="!newComment.trim() || loading"
          class="btn-primary self-end px-4 py-2 text-sm disabled:opacity-40"
        >
          {{ loading ? '...' : 'Post' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TaskComment } from '~/core/modules/tasks/infrastructure/tasks-api'

const props = defineProps<{
  comments: TaskComment[]
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'add-comment', content: string): void
}>()

const newComment = ref('')

function formatTime(dateStr: string) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  if (hrs < 24) return `${hrs}h ago`
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

function handleEnter(e: KeyboardEvent) {
  if (!e.shiftKey) {
    submitComment()
  } else {
    newComment.value += '\n'
  }
}

function submitComment() {
  const content = newComment.value.trim()
  if (!content) return
  emit('add-comment', content)
  newComment.value = ''
}
</script>

<style scoped>
.btn-primary {
  @apply bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-medium transition-colors;
}
</style>
