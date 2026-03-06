import { IRequestSurface, SurfaceService } from "../infrastructure"

interface DataInput {
	surfaceType: number
	surfaceName: string
	surfaceDrainage: string
	surfaceLayer: string
	a: string
	b: string
	c: string
	subC: string
	rrf: string
	crt: string
}
interface IState {
	data: DataInput
	validates: boolean
	loading: boolean
}

export const useSurfaceStore = defineStore("setting/surface/create", {
	state: (): IState => ({
		data: {
			surfaceType: 0,
			surfaceName: "",
			surfaceDrainage: "",
			surfaceLayer: "",

			subC: "",
			a: "",
			b: "",
			c: "",
			crt: "",
			rrf: "",
		},
		validates: false,
		loading: false,
	}),
	actions: {
		checkValidate(id: number) {
			const dataRef = useInitData()
				.refSurfaceType()
				?.map((e) => {
					return { id: e.id, name: e.name, surfaceGroup: e.surface_group }
				})
			const foundItem = dataRef?.find((item) => item.id === id)
			if (foundItem) {
				if (foundItem.surfaceGroup === "Asphalt") {
					this.validates = true
				} else {
					this.validates = false
				}
			}
		},
		async create() {
			// Loading
			this.loading = true
			let type = ""
			let group = ""
			const dataRef = useInitData()
				.refSurfaceType()
				?.map((e) => {
					return { id: e.id, name: e.name, surfaceGroup: e.surface_group }
				})
			const foundItem = dataRef?.find((item) => item.id === this.data.surfaceType)
			if (foundItem) {
				type = foundItem.name
				group = foundItem.surfaceGroup
			} else {
				type = ""
				group = ""
			}

			const params: IRequestSurface = {
				type,
				name: this.data.surfaceName,
				drainage: this.data.surfaceDrainage ? Number(this.data.surfaceDrainage) : null,
				layer_coefficient: this.data.surfaceLayer ? Number(this.data.surfaceLayer) : null,
				a: this.data.a ? Number(this.data.a) : null,
				b: this.data.b ? Number(this.data.b) : null,
				c1: this.data.c ? Number(this.data.c) : null,
				c2: this.data.subC ? Number(this.data.subC) : null,
				surface_group: group,
				crt: this.data.crt ? Number(this.data.crt) : null,
				rrf: this.data.rrf ? Number(this.data.rrf) : null,
			}

			const surfaceService = new SurfaceService()
			const res = await surfaceService.post(params)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				// ล้างค่า
				this.reset()

				return res
			}
		},
		reset() {
			this.data = {
				surfaceType: 0,
				surfaceName: "",
				surfaceDrainage: "",
				surfaceLayer: "",

				subC: "",
				a: "",
				b: "",
				c: "",
				rrf: "",
				crt: "",
			}
		},
	},
	getters: {},
})
