import { IRucListData, IRucDataTable } from "../infrastructure/RoadUserCostRucModel"
import { RoadUserCostRucService } from "../infrastructure/RoadUserCostRucService"
import { IRucParentParams } from "../infrastructure/RoadUserCostRucRequest"

interface IState {
	loading: boolean
	rucList: IRucListData[]
	rucListId: number
	rucData: IRucDataTable
}

export const useRoadUserCostRucStore = defineStore("ruc/ruc", {
	state: (): IState => ({
		loading: false,
		rucList: [
			{
				id: 1,
				name: "ข้อมูลตั้งต้น",
				name_en: "default_data",
			},
			{
				id: 2,
				name: "กำลังขับเคลื่อน",
				name_en: "driving",
			},
			{
				id: 3,
				name: "ความเร็วเครื่องยนต์",
				name_en: "engine_speed",
			},
			{
				id: 4,
				name: "การสิ้นเปลืองน้ำมันเชื้อเพลิง",
				name_en: "fuel_consumption",
			},
			{
				id: 5,
				name: "การสิ้นเปลืองน้ำมันหล่อลื่น",
				name_en: "lubricant_consumption",
			},
			{
				id: 6,
				name: "การสิ้นเปลืองล้อยาง",
				name_en: "waste_of_consumption",
			},
			{
				id: 7,
				name: "การซ่อมบำรุง",
				name_en: "maintenance",
			},
			{
				id: 8,
				name: "เวลาการเดินทาง",
				name_en: "travel_time",
			},
			{
				id: 9,
				name: "การคำนวณความเร็วยานพาหนะ",
				name_en: "vehicle_speed_calculation",
			},
			{
				id: 10,
				name: "ข้อมูลปริมาณจราจร",
				name_en: "traffic_data",
			},
		],
		rucListId: 1,
		rucData: {} as IRucDataTable,
	}),
	actions: {
		async getRucData(id: number) {
			this.loading = true

			const dataName = this.rucList.find((item) => item.id === id)?.name_en

			const service = new RoadUserCostRucService()
			const res = await service.getRucData(dataName || "")

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showToast: true })
			} else {
				this.rucData = res.data
			}
		},
		async postRucParams() {
			this.loading = true
			const params = this.checkParams(this.rucData)
			const dataName = this.rucList.find((item) => item.id === this.rucListId)?.name_en

			const service = new RoadUserCostRucService()
			const res = await service.postRucParams(dataName || "", params)

			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
		switchPath(rucId: number) {
			switch (rucId) {
				case 1:
					navigateTo("/settings/models/road-user-cost/ruc/default-data")
					break
				case 2:
					navigateTo("/settings/models/road-user-cost/ruc/driving")
					break
				case 3:
					navigateTo("/settings/models/road-user-cost/ruc/engine-speed")
					break
				case 4:
					navigateTo("/settings/models/road-user-cost/ruc/fuel-consumption")
					break
				case 5:
					navigateTo("/settings/models/road-user-cost/ruc/wasted-lubricant")
					break
				case 6:
					navigateTo("/settings/models/road-user-cost/ruc/wasted-tires")
					break
				case 7:
					navigateTo("/settings/models/road-user-cost/ruc/maintenance")
					break
				case 8:
					navigateTo("/settings/models/road-user-cost/ruc/travel-time")
					break
				case 9:
					navigateTo("/settings/models/road-user-cost/ruc/vehicle-speed-calc")
					break
				case 10:
					navigateTo("/settings/models/road-user-cost/ruc/traffic-data")
					break
			}
		},
		generateName(name: string) {
			switch (name) {
				case "car_less_than_equal_seven":
					return "Car <= 7"
				case "car_over_than_seven":
					return "Car > 7"
				case "light_bus":
					return "Light Bus"
				case "medium_bus":
					return "Medium Bus"
				case "heavy_bus":
					return "Heavy Bus"
				case "light_truck":
					return "Light Truck"
				case "medium_truck":
					return "Medium Truck"
				case "heavy_truck":
					return "Heavy Truck"
				case "full_trailor":
					return "Full Trailor"
				case "semi_trailor":
					return "Semi Trailor"
			}
		},
		checkParams(params: IRucParentParams) {
			const result = {} as IRucParentParams

			for (const key in params) {
				result[key as keyof IRucParentParams] = {}

				for (const childKey in params[key as keyof IRucParentParams]) {
					const value = params[key as keyof IRucParentParams][childKey]

					if (childKey === "vehicle_name") {
						result[key as keyof IRucParentParams][childKey] = value
					} else if (typeof value === "string") {
						const numberValue = Number(value)
						result[key as keyof IRucParentParams][childKey] = isNaN(numberValue) ? value : numberValue
					} else {
						result[key as keyof IRucParentParams][childKey] = value
					}
				}
			}

			return result
		},
		generatePropHover(name: string) {
			switch (name) {
				case "car_less_than_equal_seven":
					return "รถยนต์ส่วนบุคคลไม่เกิน 7 ที่นั่ง"
				case "car_over_than_seven":
					return "รถยนต์ส่วนบุคคลมากกว่า 7 ที่นั่ง"
				case "light_bus":
					return "รถโดยสารขนาดเล็ก"
				case "medium_bus":
					return "รถโดยสารขนาดกลาง"
				case "heavy_bus":
					return "รถโดยสารขนาดใหญ่"
				case "light_truck":
					return "รถบรรทุกเล็กเกิน 4 ล้อ"
				case "medium_track":
					return "รถบรรทุกขนาดกลาง 6-10 ล้อ"
				case "heavy_truck":
					return "รถบรรทุกขนาดใหญ่มากกว่า 10 ล้อ"
				case "full_trailor":
					return "รถพ่วง"
				case "semi_trailor":
					return "รถกึ่งพ่วง"
			}
		},
	},
	getters: {
		getRucListOptions(state) {
			if (state.rucList?.length === 0) {
				return []
			}

			const options = state.rucList.map((item) => {
				return { label: item.name, value: item.id }
			})

			return options
		},
		getDataArray(state) {
			if (!state.rucData) {
				return {}
			}

			const key = Object.keys(state.rucData)
			const result = key.map((key) => ({ [key]: state.rucData[key as keyof IRucDataTable] }))

			return result
		},
	},
})
