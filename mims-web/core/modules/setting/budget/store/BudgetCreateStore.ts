import { IRequestBudget, BudgetService, IItemBudget } from "../infrastructure"

interface IBudget {
	id: number
	costPerUnit: number | null
	methodName: string
}

interface IState {
	name: string
	budget: IBudget[]
	loading: boolean
}

export const useBudgetCreateStore = defineStore("setting/budget/create", {
	state: (): IState => ({
		name: "",
		budget: [
			{
				id: 1,
				costPerUnit: null,
				methodName: "",
			},
		],
		loading: false,
	}),
	actions: {
		addItemBudget() {
			this.budget.push({
				id: this.budget.length + 1,
				costPerUnit: null,
				methodName: "",
			})
		},
		deleteItemBudget(id: number) {
			const idToDelete = id
			const indexToDelete = this.budget.findIndex((element) => element.id === idToDelete)
			if (indexToDelete !== -1) {
				this.budget.splice(indexToDelete, 1)
			}
		},
		async create() {
			// Loading
			this.loading = true

			const budget: IItemBudget[] = this.budget.map((e) => {
				return { id: 0, cost_per_unit: e.costPerUnit, method_name: e.methodName }
			})
			const params: IRequestBudget = {
				name: this.name,
				budget,
			}

			const budgetService = new BudgetService()
			const res = await budgetService.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
				return res
			} else {
				return res
			}
		},
	},
	getters: {},
})
