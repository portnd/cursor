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
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080'
    }
  },

  devServer: {
    port: 3000,
    host: '0.0.0.0'
  },

  app: {
    head: {
      title: 'Komgrip - God-Tier Starter Kit',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Production-ready, scalable monorepo starter kit built for national-level applications' }
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ]
    }
  },

  compatibilityDate: '2026-01-22'
})
