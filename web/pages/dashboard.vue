<template>
  <div>
    <!-- Role-Based Dashboard Renderer -->
    <component :is="currentView" />
  </div>
</template>

<script setup lang="ts">
import CeoView from '~/components/dashboard/CeoView.vue'
import PmView from '~/components/dashboard/PmView.vue'
import DevView from '~/components/dashboard/DevView.vue'

definePageMeta({
  layout: 'default',
  middleware: 'auth'
})

const { currentUser } = useAuth()

// Role to Component Mapping
const roleViewMap: Record<string, any> = {
  CEO: CeoView,
  MANAGER: CeoView,
  PRODUCT_OWNER: PmView,
  PM: PmView,
  ENGINEER: DevView,
  CHIEF_ENGINEER: DevView,
  DEV: DevView
}

// Determine current view based on user role
const currentView = computed(() => {
  const userRole = currentUser.value?.role?.toUpperCase() || 'ENGINEER'
  return roleViewMap[userRole] || DevView
})
</script>
