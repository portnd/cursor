export interface IHRISItem {
	id: number
	road_number: string
	office_of_highways_code: string
	section_road_number: string
	status: boolean
}

export interface IHRISPreview {
	road_group: IHRISRoadGroup[]
	road_section: IHRISRoadSection[]
}

export interface IHRISRoadSection {
	km_end: string
	km_start: string
	road_group_number: string
	section_road_eng_name: string
	section_road_number: string
	section_road_th_name: string
}

export interface IHRISRoadGroup {
	name: string
	number: string
}
