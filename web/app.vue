<template>
  <div class="min-h-screen bg-gray-900 text-gray-100" style="background-color: #111827;">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <UiAppModal />
  </div>
</template>

<script setup lang="ts">
useHead({
  htmlAttrs: {
    lang: 'en',
    class: 'antialiased bg-gray-900'
  },
  bodyAttrs: {
    class: 'bg-gray-900 text-gray-100'
  },
  /* พื้นหลังเข้มใน head เพื่อไม่ให้เห็น white flash ตอน refresh (ก่อน CSS โหลด) */
  style: [
    {
      textContent: 'html,body,#__nuxt,#__nuxt>div{background-color:#111827!important;margin:0;min-height:100%}#__nuxt,#__nuxt>div{min-height:100vh}'
    }
  ]
})
</script>
