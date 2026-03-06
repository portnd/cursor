import { IFile } from "~~/core/shared/types/File"

export interface IRoadLaneList {
	status: boolean
	code: number
	data: ILane[]
}

export interface ILane {
	lane_no: number
	lane_name: string
}

export interface IRoadDamageList {
	status: boolean
	code: number
	data: IDatum[]
}

export interface IDatum {
	year: number
	items: IYearItem[]
}

export interface IYearItem {
	id: number
	id_parent: number
	direction: Direction[]
	lane_no: number
	surveyed_date: Date
	revision: number
}

export interface Direction {
	id: number
	name: string
}

export interface IRoadDamage {
	status: boolean
	code: number
	data: IRoadDamageData
}

export interface IRoadDamageData {
	id: number
	id_parent: number
	permissions: IPermissions
	road_damage: IRoadDamageClass
	the_geom: string
	status: string
	updated_by: IUpdatedBy
	updated_date: string
}

export interface IPermissions {
	can_approve: boolean
	can_delete: boolean
	can_edit: boolean
	can_reject: boolean
	can_send: boolean
}

export interface IRoadDamageClass {
	ac_bleeding: number
	ac_icrack: number
	ac_patching: number
	ac_pothole: number
	ac_ravelling: number
	ac_pothole_area: number
	ac_pothole_count: number
	ac_ucrack: number
	cc_faulting: number
	cc_joint_seal_damage: number
	cc_non_transverse_crack: number
	cc_patching: number
	cc_spalling: number
	cc_transverse_crack: number
	cc_scaling: number
	cc_corner_break: number
	cc_scaling: number
	created_by: number
	created_date: string
	damage_input_filepath: string
	id: number
	id_parent: number
	img_filepath: string
	km_end: number
	km_start: number
	lane_no: number
	reject_reason: string
	revision: number
	road_damage_range: IRoadDamageRange[]
	road_damage_status: IRoadDamageStatus
	road_id: number
	status: string
	surveyed_date: string
	updated_by: number
	updated_date: string
	year: number
}

export interface IRoadDamageRange {
	ac_bleeding: number
	ac_cracks: number
	ac_icrack: number
	ac_patching: number
	ac_pothole: number
	ac_ravelling: number
	ac_pothole_area: number
	ac_pothole_count: number
	ac_ucrack: number
	cc_corner_breaks: number
	cc_faulting: number
	cc_joint_seal_damage: number
	cc_non_transverse_crack: number
	cc_patching: number
	cc_spalling: number
	cc_corner_break: number
	cc_scaling: number
	cc_transverse_crack: number
	id: number
	km_end: number
	km_start: number
	road_damage_id: number
	road_damage_m: IRoadDamageM[]
	the_geom: string
}

export interface IRoadDamageM {
	ac_bleeding: number
	ac_cracks: number
	ac_icrack: number
	ac_patching: number
	ac_pothole: number
	ac_ravelling: number
	ac_surface_deform: number
	ac_pothole_area: number
	ac_pothole_count: number
	ac_ucrack: number
	cc_cornerbreaks: number
	cc_faulting: number
	cc_joint_seal_damage: number
	cc_non_transverse_crack: number
	cc_patching: number
	cc_spalling: number
	cc_transverse_crack: number
	cc_corner_break: number
	cc_scaling: number
	created_by: number
	created_date: string
	hash_data: string
	id: number
	img_filepath: string
	km: number
	road_damage_range_id: number
	the_geom: string
	updated_by: number
	updated_date: string
}

export interface IRoadDamageStatus {
	id: number
	name: string
	status_code: string
}

export interface IUpdatedBy {
	created_by: number
	ref_user_owner: IDepartment
	ref_user_owner_id: number
	ref_depot_id: number
	ref_depot: IRefDepot
	email: string
	firstname: string
	id: number
	lastname: string
	profile_img_path: string
	status: string
	tel: string
	updated_by: number
}

export interface IRefDepot {
	id: number
	name: string
	depot_code: string
}

export interface IDepartment {
	id: number
	email: string
}

export interface IRoadDamageImport {
	status: boolean
	code: number
	data: IDataImport
}

export interface IDataImport {
	id: number
	id_parent: number
	lane_no: number
	surveyed_date: string
	damage_filename: string
	img_filepath: string
	direction: IDirection
}

export interface IDirection {
	id: number
	name: string
}
