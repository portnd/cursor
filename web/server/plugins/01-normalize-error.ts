/**
 * Ensure errors passed to the Nuxt dev server have both message and statusMessage.
 * Prevents @nuxt/cli createError() from crashing when it receives { status: 404 }
 * (no message/statusMessage) — which results in an empty-string H3Error that logs
 * as a noisy unhandled error in the dev server output.
 */
const STATUS_MESSAGES: Record<number, string> = {
  400: 'Bad Request',
  401: 'Unauthorized',
  403: 'Forbidden',
  404: 'Not Found',
  422: 'Unprocessable Entity',
  500: 'Internal Server Error',
  502: 'Bad Gateway',
  503: 'Service Unavailable',
}

export default defineNitroPlugin((nitroApp) => {
  nitroApp.hooks.hook('error', (error: unknown, _context: unknown) => {
    if (!error || typeof error !== 'object') return
    const e = error as Record<string, unknown>
    if (e.message && e.statusMessage) return
    const status = (e.status ?? e.statusCode) as number | undefined
    const msg = (e.statusMessage as string | undefined)
      ?? (e.message as string | undefined)
      ?? (status ? STATUS_MESSAGES[status] : undefined)
      ?? 'Error'
    try {
      if (!e.statusMessage) e.statusMessage = msg
      if (!e.message) e.message = msg
    } catch {
      // read-only property — ignore
    }
  })
})
