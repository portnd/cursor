/**
 * Add body class immediately so body:not([class]) { display: none } never hides the page (ป้องกันหน้าขาว)
 */
export default defineNuxtPlugin(() => {
	if (import.meta.client && typeof document !== "undefined") {
		document.body.classList.add("page-loading")
		if (!document.body.style.backgroundColor) {
			document.body.style.backgroundColor = "#f5f8fa"
		}
	}
})
