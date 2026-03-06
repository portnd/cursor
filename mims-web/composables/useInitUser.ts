import { useInitUserStore } from "~~/core/modules/initUser/store"

const initUserStore = useInitUserStore()

const useInitUser = () => {
	return initUserStore.data
}

export default useInitUser
