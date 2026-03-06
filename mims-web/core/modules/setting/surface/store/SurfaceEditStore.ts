import { IRequestSurface, SurfaceService } from "../infrastructure"
interface DataInput {
	surfaceType: number
	surfaceName: string
	surfaceDrainage: number | null
	surfaceLayer: number | null
	a: number | null
	b: number | null
	c: number | null
	subC: number | null
	rrf: number | null
	crt: number | null
}
interface IState {
	id: number
	data: DataInput
	validates: boolean
	loading: boolean
}

export const useSurfaceEditStore = defineStore("setting/surface/edit", {
	state: (): IState => ({
		id: 0,
		data: {
			surfaceType: 0,
			surfaceName: "",
			surfaceDrainage: null,
			surfaceLayer: null,

			subC: null,
			a: null,
			b: null,
			c: null,
			crt: null,
			rrf: null,
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
		async get(id: number) {
			// ล้างค่า
			this.reset()

			this.id = id

			// Loading
			this.loading = true

			const surfaceService = new SurfaceService()
			const res = await surfaceService.get(this.id)

			// Loading
			this.loading = false
			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				let type = 0
				const dataRef = useInitData()
					.refSurfaceType()
					?.map((e) => {
						return { id: e.id, name: e.name }
					})
				const foundItem = dataRef?.find((item) => item.name === res.data.type)
				if (foundItem) {
					type = foundItem.id
				} else {
					type = 0
				}
				const cSub = res.data.c ? res.data.c.split("^") : ["", ""]
				this.data.a = res.data.a !== 0 ? res.data.a : 0
				this.data.b = res.data.b !== 0 ? res.data.b : 0
				this.data.c = Number(cSub[0])
				this.data.subC = Number(cSub[1])
				this.data.surfaceType = type
				this.data.surfaceName = res.data.name ? res.data.name : ""
				this.data.surfaceDrainage = res.data.drainage !== 0 ? res.data.drainage : 0
				this.data.surfaceLayer = res.data.layer_coefficient !== 0 ? res.data.layer_coefficient : 0
				this.data.crt = res.data.crt !== 0 ? res.data.crt : 0
				this.data.rrf = res.data.rrf !== 0 ? res.data.rrf : 0
				this.checkValidate(this.data.surfaceType)
				return res
			}
		},
		async edit() {
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
				type: type === "" ? null : type,
				name: this.data.surfaceName === "" ? null : this.data.surfaceName,
				drainage:
					this.data.surfaceDrainage || this.data.surfaceDrainage === 0 ? Number(this.data.surfaceDrainage) : null,
				layer_coefficient:
					this.data.surfaceLayer || this.data.surfaceLayer === 0 ? Number(this.data.surfaceLayer) : null,
				a: this.data.a || this.data.a === 0 ? Number(this.data.a) : null,
				b: this.data.b || this.data.b === 0 ? Number(this.data.b) : null,
				c1: this.data.c || this.data.c === 0 ? Number(this.data.c) : null,
				c2: this.data.subC || this.data.subC === 0 ? Number(this.data.subC) : null,
				surface_group: group === "" ? null : group,
				crt: this.data.crt || this.data.crt === 0 ? Number(this.data.crt) : null,
				rrf: this.data.rrf || this.data.rrf === 0 ? Number(this.data.rrf) : null,
			}

			const surfaceService = new SurfaceService()
			const res = await surfaceService.put(this.id, params)

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
			this.id = 0
			this.data.a = null
			this.data.b = null
			this.data.c = null
			this.data.subC = null
			this.data.surfaceDrainage = null
			this.data.surfaceLayer = null
			this.data.surfaceName = ""
			this.data.surfaceType = 0
		},
	},
	getters: {},
})
