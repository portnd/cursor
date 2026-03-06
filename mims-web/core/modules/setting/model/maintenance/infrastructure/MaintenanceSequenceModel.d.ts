export interface IMaintenanceSequence {
	[key: string]: IMaintenanceSequenceItem[]
}

export interface IMaintenanceSequenceItem {
	id: number
	name?: string
}
