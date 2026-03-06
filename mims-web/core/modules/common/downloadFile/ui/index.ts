import { DownloadFileService } from "../infrastructure"

/** fileType: 'html' | 'pdf' = open in new tab; 'excel' | undefined = download */
export const downloadFile = (title = "ดาวน์โหลด", url: string, fileType?: "html" | "pdf" | "excel"): void => {
	const { $swal }: any = useNuxtApp()
	const openInNewTab = fileType === "html" || fileType === "pdf"

	$swal.fire({
		html: `
        <div class="text-center mt-8 loading">
	      <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24"><g><circle cx="3" cy="12" r="2" fill="currentColor"/><circle cx="21" cy="12" r="2" fill="currentColor"/><circle cx="12" cy="21" r="2" fill="currentColor"/><circle cx="12" cy="3" r="2" fill="currentColor"/><circle cx="5.64" cy="5.64" r="2" fill="currentColor"/><circle cx="18.36" cy="18.36" r="2" fill="currentColor"/><circle cx="5.64" cy="18.36" r="2" fill="currentColor"/><circle cx="18.36" cy="5.64" r="2" fill="currentColor"/><animateTransform attributeName="transform" dur="1.5s" repeatCount="indefinite" type="rotate" values="0 12 12;360 12 12"/></g></svg>
	      <h3 class="swal2-title mb-14 mt-3">กำลังโหลดข้อมูล</h3>
	    </div>
        <div class="text-center mt-8 download" style="display: none">
            <a href target="_blank">
                <i class="fi fi-rr-download lh-0" style="font-size: 2.75rem"></i>
                <h3 class="swal2-title mb-14  mt-5">${title}</h3>
            </a>
	    </div>
        `,
		showCloseButton: true,
		showCancelButton: false,
		showConfirmButton: false,
		allowOutsideClick: false,
		didOpen: async () => {
			const downloadFileService = new DownloadFileService()
			const res = await downloadFileService.download(url)
			if (res.status === false) {
				$swal.close()
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				const config = useRuntimeConfig()
				const apiUrl = (config.public && (config.public as any).apiUrl) || (config as any).apiUrl
				let href = typeof res.data === "string" ? res.data : ""
				if (href && apiUrl) {
					try {
						const dataOrigin = new URL(href).origin
						const apiOrigin = new URL(apiUrl).origin
						const isApiLocal = /^https?:\/\/localhost(:\d+)?$/i.test(apiOrigin) || /^https?:\/\/127\.0\.0\.1(:\d+)?$/i.test(apiOrigin)
						if (isApiLocal && dataOrigin !== apiOrigin) {
							href = apiOrigin + new URL(href).pathname
						}
					} catch (_) {}
				}
				const container = $swal.getHtmlContainer()
				if (container) {
					const loading = container.querySelector(".loading")
					if (loading) {
						loading.style.display = "none"
					}

					const download = container.querySelector(".download")
					if (download) {
						download.style.display = ""

						const link = download.querySelector("a")
						if (link) {
							link.setAttribute("href", href)
							if (!openInNewTab) {
								link.setAttribute("download", "")
							} else {
								link.removeAttribute("download")
							}
						}
					}
				}
			}
		},
	})
}
