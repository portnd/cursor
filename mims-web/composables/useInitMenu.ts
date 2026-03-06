import { useInitMenuStore } from "~~/core/modules/initMenu/store"

const initMenuStore = useInitMenuStore()

const useInitMenu = () => {
	return initMenuStore.data
}

export default useInitMenu
