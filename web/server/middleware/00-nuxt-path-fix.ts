/**
 * Avoid unhandled 404 when something requests GET /_nuxt/ or GET /_nuxt (no filename).
 * Nuxt serves assets at /_nuxt/[file]; a request to just /_nuxt/ returns 404 and triggers ERROR in dev.
 */
export default defineEventHandler((event) => {
  const path = getRequestURL(event).pathname
  if (event.method === 'GET' && (path === '/_nuxt' || path === '/_nuxt/')) {
    setResponseStatus(event, 204)
    if (event.node?.res?.end) {
      event.node.res.end()
    }
    return
  }
})
