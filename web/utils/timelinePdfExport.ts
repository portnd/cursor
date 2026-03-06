/**
 * Timeline PDF export — calls the backend API which uses chromedp (Chromium) to render HTML → PDF.
 * Identical pattern to mims-api-service: the API returns raw PDF bytes (application/pdf).
 * We fetch with the JWT Authorization header, receive a Blob, create an object URL, and open in new tab.
 */
import type { Project, EpicTimelineData, SprintTimelineData } from '~/core/modules/projects/infrastructure/projects-api'

export async function exportTimelinePdf(
  project: Project,
  mode: 'epic' | 'sprint',
  _epicData: EpicTimelineData | null,
  _sprintData: SprintTimelineData | null
): Promise<void> {
  if (!project?.id) {
    throw new Error('Project ID is required')
  }

  const { fetchWithAuth, apiBase } = useAuth()
  const url = `${apiBase.value}/sentinel/projects/${project.id}/timeline/export-pdf?mode=${mode}`

  // Fetch PDF bytes with JWT auth — same as mims DownloadFileService.download() but for binary
  const token = useCookie('token')
  if (!token.value) {
    throw new Error('Not authenticated')
  }

  const response = await fetch(url, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token.value}`,
    },
  })

  if (!response.ok) {
    let msg = `PDF generation failed (${response.status})`
    try {
      const json = await response.json()
      msg = json.message || json.error || msg
    } catch {}
    throw new Error(msg)
  }

  const blob = await response.blob()
  const objectUrl = URL.createObjectURL(blob)

  // Open PDF in new tab — browser displays it natively (like mims: openInNewTab = true)
  const link = document.createElement('a')
  link.href = objectUrl
  link.target = '_blank'
  link.rel = 'noopener noreferrer'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  // Do not revoke immediately — let the new tab load the PDF
  setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000)
}
