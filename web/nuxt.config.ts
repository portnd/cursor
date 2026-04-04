import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const nuxtConfigDir = dirname(fileURLToPath(import.meta.url))

export default defineNuxtConfig({
  ssr: true,
  // Minimal dark overlay only (no logo) — absolute path so Docker (/app) does not treat ~ literally
  spaLoadingTemplate: join(nuxtConfigDir, 'app/spa-loading-template.html'),

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
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/v1',
      // When set (e.g. in Docker), SSR uses this to reach the API from the web container (e.g. http://api:8080/api/v1)
      apiBaseServer: process.env.NUXT_PUBLIC_API_BASE_SERVER || ''
    }
  },

  nitro: {
    devProxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },

  vite: {
    server: {
      host: '0.0.0.0',
      hmr: {
        protocol: 'ws',
        host: 'localhost',
        clientPort: 3000,
      }
    },
    plugins: [
      {
        name: 'suppress-empty-nuxt-path-404',
        configureServer(server) {
          // Requests to /_nuxt (no filename) hit Vite's asset handler which returns
          // { status: 404 } — an object with no message/statusMessage.
          // @nuxt/cli's createError() then throws because the H3Error constructor
          // receives an empty message string. Intercept early and return 204 to silence it.
          server.middlewares.use((req, res, next) => {
            const url = req.url ?? ''
            if (url === '/_nuxt' || url === '/_nuxt/' || url === '/_nuxt/index.html') {
              res.statusCode = 204
              res.end()
              return
            }
            next()
          })
        }
      }
    ]
  },

  routeRules: {
    '/projects/gantt': { ssr: false },
    '/logtime': { ssr: false },       // Heavy client state (localStorage timer, local timezone dates)
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
        { name: 'description', content: 'AI-powered task management system with code review and intelligent time estimation' },
        { name: 'theme-color', content: '#111827' }
      ],
      htmlAttrs: { class: 'bg-gray-900' },
      bodyAttrs: { class: 'bg-gray-900 text-gray-100' },
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
      ],
      style: [
        { innerHTML: '#__nuxt, #__nuxt > div { background-color: #111827 !important; }', tag: 'style' },
        { innerHTML: '#nuxt-loading svg, #nuxt-loading img, .loader svg, .loader img, [id^="nuxt"] svg, [id^="nuxt"] img { display: none !important; } #nuxt-loading, .loader { background-color: #111827 !important; }', tag: 'style' }
      ]
    }
  },

  compatibilityDate: '2026-01-22'
})
