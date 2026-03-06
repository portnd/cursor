import { useActionStore } from "~~/core/shared/store/ActionStore"

const actionStore = useActionStore()

const useAction = (actions: any) => {
	return actionStore.setActions(actions)
}

export default useAction
