export interface IRucAccMaster {
	status: boolean
	code: number
	data: IRucAccMasterData[]
}

export interface IRucAccMasterData {
	id: number
	name: string
	name_en: string
}

export interface IRucAccChance {
	status: boolean
	code: number
	data: IRucAccChanceData
}

export interface IRucAccChanceData {
	road_group_id: number
	number_of_fatal_accidents: number
	number_of_accidents_with_serious_injuries: number
	number_of_accidents_with_minor_injuries: number
	number_of_accidents_with_property_damaged: number
}

export interface IRucAccLoss {
	status: boolean
	code: number
	data: IRucAccLossData
}

export interface IRucAccLossData {
	value_of_accidents_with_minor_injuries: number
	value_of_accidents_with_property_damaged: number
	value_of_accidents_with_serious_injuries: number
	value_of_fatal_accidents: number
}
