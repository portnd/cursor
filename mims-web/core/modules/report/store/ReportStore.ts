import {
	IReportAccidentParams,
	IReportAsset,
	IReportAssetImprove,
	IReportAssetSummary,
	IReportConditionParams,
	IReportDamageParams,
	IReportMaintenanceHistoryParams,
	IReportMaintenanceTrackingParams,
	IReportTrackingMaintenance,
	IReportTrafficParams,
	IReportsHistMaintenance,
	IReportSurfaceSummary,
	IReportAccidentData,
	IReportRoadsTree,
	IReportTrafficData,
	IReportAssetData,
	IReportAssetSummaryData,
	IReportAssetImproveData,
	IReportRoadConditionSummaryData,
} from "../infrastructure"
import { ReportSerivce } from "../infrastructure/ReportService"
import { IReportRoadConditionSummary, IReportRoadSurfaceSummary } from "../infrastructure/ReportsRequest"

interface IStateParams {
	road_id: number | null
	year: number | null
	start_year: number | null
	end_year: number | null
	distance: number | null
	road_multiplier_id: number[]
	asset_group_id: number | null
	asset_id: number | null
	month: string | null
	surface_type: string | null
	measure_id: number | null
}

interface IState {
	loading: boolean
	reportName: string
	yearLists: number[]
	histOptions: IReportsHistMaintenance[]
	trackingOptions: IReportTrackingMaintenance[]
	roadTree: IReportRoadsTree[]
	trafficsOptions: IReportTrafficData
	accidentOptions: IReportAccidentData
	assetOptions: IReportAssetData
	assetSummaryOptions: IReportAssetSummaryData[]
	assetImproveOptions: IReportAssetImproveData
	roadConditionOptions: IReportRoadConditionSummaryData
	surfaceSummaryOptions: IReportSurfaceSummary
	params: IStateParams
}

export const useReportStore = defineStore("report", {
	state: (): IState => ({
		loading: false,
		reportName: "รายงานข้อมูลสภาพทาง",
		yearLists: [],
		histOptions: [],
		trackingOptions: [],
		roadTree: [],
		trafficsOptions: {} as IReportTrafficData,
		accidentOptions: {} as IReportAccidentData,
		assetOptions: {} as IReportAssetData,
		assetSummaryOptions: [],
		assetImproveOptions: {} as IReportAssetImproveData,
		roadConditionOptions: {} as IReportRoadConditionSummaryData,
		surfaceSummaryOptions: {} as IReportSurfaceSummary,
		params: {
			road_id: null,
			year: null,
			start_year: null,
			end_year: null,
			distance: null,
			road_multiplier_id: [],
			asset_group_id: null,
			asset_id: null,
			month: null,
			surface_type: null,
			measure_id: null,
		},
	}),
	actions: {
		async getAssetsData() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getAssetOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.assetOptions = res.data
			}
		},
		async getAssetsMapData() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getAssetMapOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.assetOptions = res.data
			}
		},
		async getAssetSummaryData() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getAssetSummaryOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.assetSummaryOptions = res.data
			}
		},
		async getAssetImproveData() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getAssetImproveOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.assetImproveOptions = res.data
			}
		},
		async getConditionSummaryData() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getConditionSummaryOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadConditionOptions = res.data
			}
		},
		async getConditionYearList() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getConditionsYears()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.yearLists = res.data?.year
			}
		},
		async getDamageYearList() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getDamageYears()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.yearLists = res.data?.year
			}
		},
		async getSurfaceSummary() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getSurfaceOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.surfaceSummaryOptions = res.data
			}
		},
		async getMaintenanceHistoryOptions() {
			this.loading = true

			const serivce = new ReportSerivce()
			const res = await serivce.getMaintenanceHistDataOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.histOptions = res.data
			}
		},
		async getMaintenanceTrackingOptions() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getMaintenanceTrackingDataOptions()

			this.loading = false
			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.trackingOptions = res.data
			}
		},
		async getMaintenanceHistDataOptions() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getMaintenanceHistDataOptions()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.histOptions = res.data
			}
		},
		async getRoadTree() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getRoadTrees()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.roadTree = res.data
				this.params.road_id = this.getRoadOptions[1].value
			}
		},
		async getTrafficOptions() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getTrafficData()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.trafficsOptions = res.data
			}
		},
		async getAccidentOptions() {
			this.loading = true

			const service = new ReportSerivce()
			const res = await service.getAccidentData()

			this.loading = false

			if (!res.status) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.accidentOptions = res.data
			}
		},
		async toggleReportName(name: string) {
			this.reportName = name

			switch (this.reportName) {
				case "รายงานสินทรัพย์":
					await this.getAssetsData()
					this.params.asset_group_id = this.getAssestGroupsOptions.length
						? this.getAssestGroupsOptions[0].value
						: this.params.asset_group_id
					this.params.asset_id = this.getAssetListsOptions.length
						? this.getAssetListsOptions[0].value
						: this.params.asset_id
					this.params.road_id = this.getRoadGroupOptions.length ? this.getRoadOptions[1]?.value : this.params.road_id
					break
				case "รายงานแผนที่สินทรัพย์":
					await this.getAssetsMapData()
					this.params.asset_group_id = this.getAssestGroupsOptions.length
						? this.getAssestGroupsOptions[0].value
						: this.params.asset_group_id
					this.params.asset_id = this.getAssetListsOptions.length
						? this.getAssetListsOptions[0].value
						: this.params.asset_id
					this.params.road_id = this.getRoadGroupOptions.length ? this.getRoadOptions[1]?.value : this.params.road_id
					break
				case "รายงานสรุปสินทรัพย์":
					await this.getAssetSummaryData()
					this.params.road_id = this.getRoadGroupOptions.length ? this.getRoadOptions[1]?.value : this.params.road_id
					break
				case "รายงานการปรับแก้สินทรัพย์ประจำเดือน":
					await this.getAssetImproveData()
					this.updateYearList(this.reportName)

					break
				case "รายงานข้อมูลสภาพทาง":
					await this.getConditionYearList()
					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0].value : this.params.year
					this.params.road_id = this.getRoadOptions.length ? this.getRoadOptions[1]?.value : this.params.road_id
					break
				case "รายงานสรุปข้อมูลสภาพทาง":
					await this.getConditionSummaryData()
					this.updateYearList(this.reportName)
					break
				case "รายงานสรุปข้อมูลความเสียหาย":
					await this.getDamageYearList()
					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0].value : this.params.year
					this.params.road_id = this.getRoadOptions.length ? this.getRoadOptions[1]?.value : this.params.road_id
					break
				case "รายงานสรุปรายละเอียดชนิดผิวทาง":
					await this.getSurfaceSummary()
					this.updateYearList(this.reportName)
					break
				case "รายงานติดตามการซ่อมบำรุง":
					await this.getMaintenanceTrackingOptions()
					this.updateYearList(this.reportName)
					break
				case "รายงานประวัติการซ่อมบำรุง":
					await this.getMaintenanceHistDataOptions()
					this.updateYearList(this.reportName)
					break
				case "รายงานปริมาณจราจร":
					await this.getTrafficOptions()
					this.updateYearList(this.reportName)
					break
				case "รายงานอุบัติเหตุ":
					await this.getAccidentOptions()
					this.updateYearList(this.reportName)
					break
			}
		},
		updateYearList(reportName: string) {
			let yearList: number[] = []
			switch (reportName) {
				case "รายงานปริมาณจราจร":
					const trafficYears = this.trafficsOptions.year

					yearList = trafficYears?.sort((a, b) => b - a)

					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year
					this.params.road_id = this.getRoadGroupOptions.length
						? this.getRoadGroupOptions[0]?.value
						: this.params.road_id
					break
				case "รายงานอุบัติเหตุ":
					const accidentYears = this.accidentOptions.year

					yearList = accidentYears?.sort((a, b) => b - a)
					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year
					this.params.road_id = this.getRoadGroupOptions.length
						? this.getRoadGroupOptions[0]?.value
						: this.params.road_id
					break
				case "รายงานการปรับแก้สินทรัพย์ประจำเดือน":
					const assetImproveYear = this.assetImproveOptions.year

					yearList = assetImproveYear?.sort((a, b) => b - a)
					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year
					this.params.month = this.getMonthOptions.length ? this.getMonthOptions[0]?.value : this.params.month
					this.params.road_id = this.getRoadGroupOptions.length
						? this.getRoadGroupOptions[0]?.value
						: this.params.road_id
					break
				case "รายงานสรุปข้อมูลสภาพทาง":
					const roadConditionSummaryYear = this.roadConditionOptions.year

					yearList = roadConditionSummaryYear?.sort((a, b) => b - a)
					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year
					this.params.surface_type = this.getSurfaceTypeOptions.length
						? this.getSurfaceTypeOptions[0]?.value
						: this.params.surface_type

					this.params.road_id = this.getRoadGroupOptions.length
						? this.getRoadGroupOptions[0]?.value
						: this.params.road_id

					this.params.measure_id = this.getOwnerOptions.length ? this.getOwnerOptions[0]?.value : this.params.measure_id
					break
				case "รายงานสรุปรายละเอียดชนิดผิวทาง":
					const roadSurfaceSummary = this.surfaceSummaryOptions.year

					yearList = roadSurfaceSummary?.sort((a, b) => b - a)
					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year
					this.params.road_id = this.getRoadGroupOptions.length
						? this.getRoadGroupOptions[0]?.value
						: this.params.road_id

					break
				case "รายงานติดตามการซ่อมบำรุง":
					yearList = this.trackingOptions.map((item) => item.year).sort((a, b) => b - a)
					this.yearLists = yearList

					this.params.year = this.getYearListOptions.length ? this.getYearListOptions[0]?.value : this.params.year

					this.params.road_multiplier_id =
						this.getRoadsTreeOptions.reduce<number[]>((acc, item) => {
							if (item.children.length) {
								const ids = item.children.map((child) => child.id)
								return acc.concat(ids)
							}

							return acc
						}, []) ?? []
					break
				case "รายงานประวัติการซ่อมบำรุง":
					yearList =
						this.histOptions
							.filter((item) => this.params.road_multiplier_id.includes(item.road_group_id))
							.flatMap((item) => item.year)
							.sort((a, b) => b - a) ?? []
					yearList = [...new Set(yearList)]
					this.yearLists = yearList

					this.params.start_year = this.getYearListOptions.length
						? this.getYearListOptions[1]?.value
						: this.params.start_year

					this.params.end_year = this.getYearListOptions.length
						? this.getYearListOptions[0]?.value
						: this.params.end_year

					this.params.road_multiplier_id =
						this.getRoadsTreeOptions.reduce<number[]>((acc, item) => {
							if (item.children.length) {
								const ids = item.children.map((child) => child.id)
								return acc.concat(ids)
							}

							return acc
						}, []) ?? []
			}
		},
		encodeQuery(data: any) {
			if (data) {
				const result: string[] = []
				const keys = Object.keys(data).filter((key) => data[key as keyof typeof data] !== undefined)
				for (const key of keys) {
					const value = data[key as keyof typeof data]
					if (value !== null && value !== undefined) {
						result.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
					}
				}
				return result.join("&")
			} else {
				return ""
			}
		},
		downloadFile(fileType: string) {
			let params
			const reportName = this.reportName
			let keyword = ""

			switch (reportName) {
				case "รายงานสินทรัพย์":
					params = {} as IReportAsset
					params.road_id = this.params.road_id!
					params.asset_id = this.params.asset_id!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName} `, `report/asset/report?${keyword}`)
					break
				case "รายงานแผนที่สินทรัพย์":
					params = {} as IReportAsset
					params.road_id = this.params.road_id!
					params.asset_id = this.params.asset_id!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName} `, `report/map/report?${keyword}`)
					break
				case "รายงานสรุปสินทรัพย์":
					params = {} as IReportAssetSummary
					params.road_id = this.params.road_id!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName} `, `report/summary_asset/report?${keyword}`)
					break
				case "รายงานการปรับแก้สินทรัพย์ประจำเดือน":
					params = {} as IReportAssetImprove
					params.year = this.params.year!
					params.month = this.params.month!
					params.road_id = this.params.road_id!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName} `, `report/asset_adjustment/report?${keyword}`)
					break
				case "รายงานข้อมูลสภาพทาง":
					if (!this.params.road_id) {
						useHandlerError(0, { message: "โปรดเลือกสายทาง" }, { showAlert: true })
					}
					params = {} as IReportConditionParams
					params.year = this.params.year!
					params.road_id = this.params.road_id!
					params.distance = this.params.distance!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/condition/report?${keyword}`)
					break
				case "รายงานสรุปข้อมูลสภาพทาง":
					if (!this.params.road_id) {
						useHandlerError(0, { message: "โปรดเลือกสายทาง" }, { showAlert: true })
					}
					params = {} as IReportRoadConditionSummary
					params.year = this.params.year!
					params.road_id = this.params.road_id!
					params.factor = this.params.surface_type!
					params.measure_id = this.params.measure_id!
					params.type = fileType

					keyword = this.encodeQuery(params)

					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/summary_condition/report?${keyword}`)
					break
				case "รายงานสรุปข้อมูลความเสียหาย":
					if (!this.params.road_id) {
						useHandlerError(0, { message: "โปรดเลือกสายทาง" }, { showAlert: true })
					}
					params = {} as IReportDamageParams

					params.year = this.params.year!
					params.road_id = this.params.road_id!
					params.type = fileType

					keyword = this.encodeQuery(params)
					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/damage/report?${keyword}`)
					break
				case "รายงานสรุปรายละเอียดชนิดผิวทาง":
					if (!this.params.road_id) {
						useHandlerError(0, { message: "โปรดเลือกสายทาง" }, { showAlert: true })
					}
					params = {} as IReportRoadSurfaceSummary

					params.year = this.params.year!
					params.road_id = this.params.road_id!
					params.type = fileType

					keyword = this.encodeQuery(params)
					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/surface/report?${keyword}`)
					break
				case "รายงานติดตามการซ่อมบำรุง":
					if (this.params.road_multiplier_id.length === 0) {
						useHandlerError(0, { message: "โปรดเลือกสายทาง" }, { showAlert: true })
					} else {
						params = {} as IReportMaintenanceTrackingParams
						params.road_group_id = this.params.road_multiplier_id
						params.type = fileType
						params.year = this.params.year!

						keyword = this.encodeQuery(params)
						useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/maintenance_tracking/report?${keyword}`)
					}

					break
				case "รายงานปริมาณจราจร":
					params = {} as IReportTrafficParams
					params.road_group_id = this.params.road_id!
					params.type = fileType
					params.year = this.params.year!

					keyword = this.encodeQuery(params)
					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/traffic_volume/report?${keyword}`)

					break
				case "รายงานอุบัติเหตุ":
					params = {} as IReportAccidentParams
					params.road_group_id = this.params.road_id!
					params.type = fileType
					params.year = this.params.year!

					keyword = this.encodeQuery(params)
					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/accident_volume/report?${keyword}`)

					break
				case "รายงานประวัติการซ่อมบำรุง":
					params = {} as IReportMaintenanceHistoryParams
					const diff = Math.abs(this.params.end_year! - this.params.start_year!)

					if (diff > 10) {
						useHandlerError(0, { message: "เลือกระยะเวลาได้สูงสุด 10 ปี" }, { showAlert: true })
						break
					}

					if (this.params.start_year! > this.params.end_year!) {
						useHandlerError(0, { message: "ปีเริ่มต้นต้องน้อยกว่าปีสิ้นสุด" }, { showAlert: true })
						break
					}
					params.type = fileType
					params.road_group_id = this.params.road_multiplier_id
					params.year_start = this.params.start_year!
					params.year_end = this.params.end_year!

					keyword = this.encodeQuery(params)
					useDownloadFile(`ดาวน์โหลด : ${reportName}`, `report/maintenance_history/report?${keyword}`)

					break
			}
		},
		updateAssetDefaultParams() {
			this.params.asset_id = this.getAssetListsOptions[0].value
		},
	},
	getters: {
		getRoadsChildTreeOptions(state) {
			const options = state.roadTree.map((parent) => {
				return {
					id: `parent_${parent.id}`,
					label: parent.label,
					children: parent.children
						.sort((a, b) => a.id - b.id)
						.map((child) => {
							return {
								id: child.id,
								label: child.label,
							}
						}),
				}
			})

			return options || []
		},
		getRoadsTreeOptions(state) {
			const options = [
				{
					id: "all_0",
					label: "เลือกทั้งหมด",
					children: state.roadTree.map((item) => {
						return {
							id: item.id,
							label: item.label,
						}
					}),
				},
			]

			return options || []
		},
		getYearListOptions(state) {
			if (!state.yearLists.length) {
				return []
			}

			const options = state.yearLists
				.map((year) => {
					return { label: `${year + 543}`, value: year }
				})
				.sort((a, b) => b.value - a.value)

			return options
		},
		getRoadOptions(state) {
			const road = state.roadTree
			let options = []

			options = road.map((item) => {
				return { label: item.label, value: item.id }
			})

			return options || []
		},
		getRoadGroupOptions(state) {
			let options: { label: string; value: number }[] = []
			let roadGroup = []

			switch (state.reportName) {
				case "รายงานปริมาณจราจร":
					roadGroup = state.trafficsOptions?.road_group
					options = roadGroup?.map((item) => ({ label: item.name, value: item.id }))
					break
				case "รายงานอุบัติเหตุ":
					roadGroup = state.accidentOptions?.road_group
					options = roadGroup?.map((item) => ({ label: item.name, value: item.id }))
					break
				default:
					options = []
					break
			}

			return options
		},
		getAssestGroupsOptions(state) {
			const assetData = state.assetOptions

			const options = assetData.group?.map((asset) => ({ label: asset.name, value: asset.id }))

			return options || []
		},
		getAssetListsOptions(state) {
			const assetData = state.assetOptions

			const assetFilter = assetData.group?.filter((group) => group.id === state.params.asset_group_id)
			const options = assetFilter?.flatMap((group) => {
				return group.asset.map((asset) => ({ label: asset.name, value: asset.id }))
			})

			return options || []
		},
		getMonthOptions(state) {
			const months = state.assetImproveOptions?.month

			const options = months?.map((month) => ({ label: month, value: month }))

			return options || []
		},
		getSurfaceTypeOptions(state) {
			const surfaceType = state.roadConditionOptions?.type

			const options = surfaceType?.map((surface) => ({ label: surface, value: surface }))

			return options || []
		},
		getOwnerOptions(state) {
			const owner = state.roadConditionOptions?.measure

			const options = owner?.map((item) => ({ label: item.name, value: item.id }))

			return options || []
		},
	},
})
