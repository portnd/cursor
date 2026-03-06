export interface IRucAccChanceParams {
	number_of_accidents_with_minor_injuries: number
	number_of_accidents_with_property_damaged: number
	number_of_accidents_with_serious_injuries: number
	number_of_fatal_accidents: number
	road_group_id: number
}

export interface IRucAccLossParams {
	value_of_accidents_with_minor_injuries: number
	value_of_accidents_with_property_damaged: number
	value_of_accidents_with_serious_injuries: number
	value_of_fatal_accidents: number
}
