export const useActionStore = defineStore("action", {
	state: (): any => ({
		actions: null,
	}),
	actions: {
		setActions(actions: any) {
			this.actions = actions
		},
	},
	getters: {},
})
