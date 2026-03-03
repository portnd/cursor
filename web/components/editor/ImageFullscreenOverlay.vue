<template>
  <Teleport to="body">
    <div
      v-if="show && imageSrc"
      class="image-fullscreen-overlay"
      @click.self="$emit('close')"
    >
      <button
        type="button"
        class="image-fullscreen-close"
        aria-label="Close"
        @click="$emit('close')"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
        </svg>
      </button>
      <img
        :src="imageSrc"
        class="image-fullscreen-img"
        alt="Full screen"
        @click.stop
      />
      <div v-if="showActions" class="image-fullscreen-actions">
        <button
          type="button"
          class="image-fullscreen-annotate-btn"
          @click="$emit('annotate')"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
          </svg>
          Annotate
        </button>
        <button
          type="button"
          class="image-fullscreen-delete-btn"
          @click="$emit('delete')"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
          </svg>
          Delete image
        </button>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
defineProps<{
  show: boolean
  imageSrc: string | null
  showActions: boolean
}>()

defineEmits<{
  close: []
  annotate: []
  delete: []
}>()
</script>

<style scoped>
.image-fullscreen-overlay {
  position: fixed;
  inset: 0;
  z-index: 9998;
  background: rgba(0, 0, 0, 0.92);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 16px 80px;
}

.image-fullscreen-close {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgb(156 163 175);
  background: rgb(31 41 55);
  border: 1px solid rgb(55 65 81);
  border-radius: 10px;
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
}

.image-fullscreen-close:hover {
  color: white;
  background: rgb(55 65 81);
}

.image-fullscreen-img {
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.5);
  user-select: none;
  pointer-events: none;
}

.image-fullscreen-actions {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
}

.image-fullscreen-annotate-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  background: rgb(37 99 235);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}

.image-fullscreen-annotate-btn:hover {
  background: rgb(29 78 216);
}

.image-fullscreen-delete-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 600;
  color: rgb(254 202 202);
  background: rgb(127 29 29);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.image-fullscreen-delete-btn:hover {
  background: rgb(185 28 28);
  color: white;
}
</style>
