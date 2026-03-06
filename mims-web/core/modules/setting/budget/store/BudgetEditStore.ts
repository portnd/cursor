import { IRequestBudget, BudgetService, IItemBudget } from "../infrastructure"

interface IBudget {
	id: number
	costPerUnit: number | null
	methodName: string
	is_show_method?: boolean
}
interface IState {
	id: number
	name: string
	budget: IBudget[]
	loading: boolean
}

export const useBudgetEditStore = defineStore("setting/budget/edit", {
	state: (): IState => ({
		id: 0,
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
				id: 0,
				costPerUnit: null,
				methodName: "",
			})
		},
		deleteItemBudget(id: number) {
			const idToDelete = id
			const indexToDelete = this.budget.findIndex((_, key) => key === idToDelete)
			if (indexToDelete !== -1) {
				this.budget.splice(indexToDelete, 1)
			}
		},
		async get(id: number) {
			this.id = id
			// Loading
			this.loading = true

			const budgetService = new BudgetService()
			const res = await budgetService.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.name = res.data.name
				const budgetArr: IBudget[] = []

				for (const item of res.data.budget as IItemBudget[]) {
					const budget: IBudget = {
						id: item.id,
						costPerUnit: item.cost_per_unit,
						methodName: item.method_name,
						is_show_method: item.is_show_method,
					}
					budgetArr.push(budget)
				}

				this.budget = budgetArr
				return res
			}
		},
		async edit() {
			// Loading
			this.loading = true

			const budget: IItemBudget[] = this.budget.map((e) => {
				return { id: e.id, cost_per_unit: e.costPerUnit, method_name: e.methodName }
			})
			const params: IRequestBudget = {
				id: this.id,
				name: this.name,
				budget,
			}

			const budgetService = new BudgetService()
			const res = await budgetService.put(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
				return res
			} else {
				// ล้างค่า
				return res
			}
		},
	},
	getters: {},
})
