/**
 * API Plugin - Configure $fetch baseURL for SSR
 */

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()

  const api = $fetch.create({
    baseURL: config.public.apiBase
  })

  return {
    provide: {
      api
    }
  }
})
