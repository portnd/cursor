interface IState {
	loading: boolean
	map: any
	longdo: any
	geom: string
	color: string
	plan: number
	year: number
	result: number
	condition: number
	fullScreen: boolean
}
export const useAnnualAnalyzeSummaryStore = defineStore("annual/summary", {
	state: (): IState => ({
		loading: false,
		map: null,
		longdo: null,
		geom: "",
		color: "",
		plan: 4,
		year: 10,
		result: 1,
		condition: 1,
		fullScreen: false,
	}),
	actions: {
		setMap(map: any) {
			this.map = map
			if (this.map) {
				// @ts-ignore
				this.longdo = window.longdo
			}

			this.createLine()
			this.map.Event.bind("fullscreen", this.fullscreen)
		},
		fullscreen() {
			const element = document.querySelector("#map")
			const map = element?.querySelector(".longdo-map")
			const fullScreenElement = map?.querySelector(".ldmap_placeholder_fullscreen") as HTMLElement
			if (fullScreenElement !== null) {
				this.fullScreen = true
			} else {
				this.fullScreen = false
			}
		},
		getHtmlPopUp() {
			const html = ref("")
			if (this.result !== 2) {
				html.value = `<div class="d-flex gap-2">
        <i class="fi-rr-chart-histogram fs-4 p-0" style="color:${this.color}"></i>
        <h4 class="fw-semibold" style="color:${this.color}">IRI 2.00 ม./กม.</h4>
      </div>
      <div class="row">
        <div class="col-4 text-gray-800 mb-2">สายทาง:</div>
        <div class="col-6 text-gray-700 mb-2">กรุงเทพมหานคร - บ้านฉาง</div>
        <div class="col-4 text-gray-800 mb-2">ปีที่:</div>
        <div class="col-6 text-gray-700 mb-2">1</div>
        <div class="col-4 text-gray-800 mb-2">IRI ก่อนซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">2.00  ม./กม.</div>
        <div class="col-4 text-gray-800 mb-2">IRI หลังซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">1.00  ม./กม.</div>
    </div>`
			} else {
				html.value = `<div class="d-flex gap-2">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <g clip-path="url(#clip0_2983_17320)">
        <path d="M0 9H24V22H16.152C17.586 20.808 18.5 19.011 18.5 17C18.5 14.297 16.849 11.98 14.501 11H9.5C7.152 11.98 5.501 14.297 5.501 17C5.501 19.011 6.415 20.808 7.849 22H0V9ZM14 12.988V15.5C14 16.605 13.105 17.5 12 17.5C10.895 17.5 10 16.605 10 15.5V12.988C8.524 13.726 7.5 15.237 7.5 17C7.5 18.956 8.756 20.605 10.5 21.224V24H13.5V21.224C15.244 20.604 16.5 18.956 16.5 17C16.5 15.237 15.476 13.727 14 12.988ZM24 3V7H0V3C0 1.346 1.346 0 3 0H21C22.654 0 24 1.346 24 3ZM5 3.5C5 2.672 4.328 2 3.5 2C2.672 2 2 2.672 2 3.5C2 4.328 2.672 5 3.5 5C4.328 5 5 4.328 5 3.5ZM9 3.5C9 2.672 8.328 2 7.5 2C6.672 2 6 2.672 6 3.5C6 4.328 6.672 5 7.5 5C8.328 5 9 4.328 9 3.5Z" fill="${this.color}"/>
        </g>
        <defs>
        <clipPath id="clip0_2983_17320">
        <rect width="24" height="24" fill="white"/>
        </clipPath>
        </defs>
        </svg>
        <h4 class="fw-semibold" style="color:${this.color}">SS : Fibro Seal</h4>
      </div>
      <div class="row">
        <div class="col-4 text-gray-800 mb-2">สายทาง:</div>
        <div class="col-6 text-gray-700 mb-2">กรุงเทพมหานคร - บ้านฉาง</div>
        <div class="col-4 text-gray-800 mb-2">ปีที่:</div>
        <div class="col-6 text-gray-700 mb-2">1</div>
        <div class="col-4 text-gray-800 mb-2">IRI ก่อนซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">3.50  ม./กม.</div>
        <div class="col-4 text-gray-800 mb-2">IRI หลังซ่อม:</div>
        <div class="col-6 text-gray-700 mb-2">1.00  ม./กม.</div>
    </div>`
			}

			return html.value
		},
		createLine() {
			console.log(1)
			if (this.map) {
				console.log(2)
				this.map.Overlays.clear()
				this.geom = `LINESTRING(100.60633432409338 13.821522736504463,100.602784216045 13.821653647098415,100.59840724738694 13.821034002969526,100.5962951577875 13.820649997329767,100.59623498763187 13.810615925266177,100.59624888251949 13.808909356040601,100.59677811904618 13.806716534644991,100.59668991295837 13.80593705019568,100.59673401600236 13.803889820386786,100.59668991295837 13.803555751891949,100.59262361230782 13.803615712939205,100.5908360491419 13.803657730944948,100.58982168776913 13.80546510526736)`
				this.color = this.result !== 2 ? "#50CD89" : "#418FFF"
				const coordinator = this.geom.split(",")[0]?.split("(")[1]?.split(" ")
				const lines = this.longdo.Util.overlayFromWkt(this.geom, {
					detail: this.getHtmlPopUp(),
					lineColor: this.color,
				})
				this.map.Overlays.add(lines[0])
				this.map.location({
					lon: coordinator[0],
					lat: coordinator[1],
				})
				this.map.zoom(18)
			}
		},
	},
	getters: {},
})
