export default defineNuxtConfig({
  ssr: true,

  devtools: { enabled: true },

  modules: [
    '@pinia/nuxt'
  ],
  
  css: ['~/assets/css/tailwind.css'],
  
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  typescript: {
    strict: true,
    typeCheck: false
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080/api/v1'
    }
  },

  devServer: {
    port: 3000,
    host: '0.0.0.0'
  },

  app: {
    head: {
      title: 'The Sentinel - AI-Powered Task Management',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'AI-powered task management system with code review and intelligent time estimation' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },

  compatibilityDate: '2026-01-22'
})
