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
  PM: PmView,
  DEV: DevView
}

// Determine current view based on user role
const currentView = computed(() => {
  const userRole = currentUser.value?.role?.toUpperCase() || 'DEV'
  
  // Log for debugging
  console.log('🔍 User Role:', userRole)
  console.log('📊 Current User:', currentUser.value)
  
  // Return appropriate component (default to DevView)
  return roleViewMap[userRole] || DevView
})

// Watch for role changes (useful for debugging)
watch(currentUser, (newUser) => {
  console.log('👤 User changed:', newUser)
}, { immediate: true })
</script>
