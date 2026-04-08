import { authApi } from '~/core/modules/auth/infrastructure/auth-api'

const THEME_KEY = 'sentinel-theme'

export const useTheme = () => {
  const isDark = useState<boolean>('theme-dark', () => true)

  function applyTheme(dark: boolean) {
    if (!import.meta.client) return
    const html = document.documentElement
    if (dark) {
      html.classList.add('dark')
      html.classList.remove('light')
    } else {
      html.classList.remove('dark')
      html.classList.add('light')
    }
  }

  function persist(dark: boolean) {
    if (!import.meta.client) return
    localStorage.setItem(THEME_KEY, dark ? 'dark' : 'light')
  }

  // Fire-and-forget: save preference to user account (requires auth token)
  async function syncToBackend(dark: boolean) {
    try {
      const tokenCookie = useCookie('token')
      if (!tokenCookie.value) return
      await authApi.updateTheme(dark ? 'dark' : 'light')
    } catch {
      // silent — localStorage is still the fallback
    }
  }

  function toggle() {
    isDark.value = !isDark.value
    applyTheme(isDark.value)
    persist(isDark.value)
    syncToBackend(isDark.value)
  }

  function setTheme(dark: boolean) {
    isDark.value = dark
    applyTheme(dark)
    persist(dark)
  }

  // Called on page load — resolves theme in priority order:
  //   1. Explicit user account preference (passed in from /auth/me response)
  //   2. localStorage (browser-level fallback)
  //   3. OS dark-mode preference
  function initTheme(accountTheme?: 'dark' | 'light') {
    if (!import.meta.client) return
    const saved = localStorage.getItem(THEME_KEY)
    const prefersDark = window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? true

    let dark: boolean
    if (accountTheme !== undefined) {
      dark = accountTheme === 'dark'
    } else if (saved !== null) {
      dark = saved === 'dark'
    } else {
      dark = prefersDark
    }

    isDark.value = dark
    applyTheme(dark)
    persist(dark)
  }

  return {
    isDark: readonly(isDark),
    toggle,
    setTheme,
    initTheme,
  }
}
