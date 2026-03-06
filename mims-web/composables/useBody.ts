import { DrawerComponent } from "~~/assets/themes/ts/components/_DrawerComponent"

const useBody = () => {
	const addBodyClassname = (className: string) => {
		document.body.classList.add(className)
	}

	const removeBodyClassName = (className: string) => {
		document.body.classList.remove(className)
	}

	const addBodyAttribute = (payload: { qualifiedName: string; value: string }) => {
		const { qualifiedName, value } = payload
		document.body.setAttribute(qualifiedName, value)
	}

	const removeBodyAttribute = (payload: { qualifiedName: string }) => {
		const { qualifiedName } = payload
		document.body.removeAttribute(qualifiedName)
	}

	const removeOverlay = () => {
		setTimeout(() => {
			DrawerComponent.hideAll()
		}, 250)

		// กรณีหน้าจอ Mobile ตอนกดเมนูย่อย บางครั้งไม่ยอมซ่อน
		const subMenus = document.querySelectorAll("#kt_aside_menu .show")
		subMenus.forEach((subMenu) => {
			subMenu.classList.remove("show")
		})
	}

	return {
		addBodyClassname,
		removeBodyClassName,
		addBodyAttribute,
		removeBodyAttribute,
		removeOverlay,
	}
}

export default useBody
