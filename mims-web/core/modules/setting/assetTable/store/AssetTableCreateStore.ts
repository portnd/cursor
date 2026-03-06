import { IRequestAssetTable, AssetTableService, IRequestDataAssetTable } from "../infrastructure"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	asset: IRequestAssetTable
	imageAsset: IFile
	loading: boolean
}

export const useAssetTableCreateStore = defineStore("setting/asset-table/create", {
	state: (): IState => ({
		imageAsset: {} as IFile,
		asset: {
			icon_filepath: null,
			icon_filepath_status: "no_file",
			data: {
				table_name: "",
				table_label: "",
				delete_columns: [],
				asset_type: "",
				asset_group: "",
				geom_type: "",
				line_color_code: "",
				approver_id: [],
				viewer_id: [],
				columns: [],
			},
		},
		loading: false,
	}),
	actions: {
		reset(type: string) {
			this.imageAsset = {} as IFile

			this.asset = {
				icon_filepath: null,
				icon_filepath_status: "no_file",
				data: {
					table_name: "",
					table_label: "",
					delete_columns: [],
					asset_type: "",
					asset_group: "",
					geom_type: type,
					line_color_code: "",
					approver_id: [],
					viewer_id: [],
					columns: [],
				},
			}
			this.loading = false
			this.storeAsset()
		},
		addColumnAsset(item: any) {
			item.column_id = this.asset.data.columns.length + 1
			this.asset.data.columns.push(item)
		},
		deleteColumnAsset(item: any) {
			const idToDelete = item.column_id

			const indexToDelete = this.asset.data.columns.findIndex((element) => element.column_id === idToDelete)
			if (indexToDelete !== -1) {
				this.asset.data.columns.splice(indexToDelete, 1)
			}
		},
		editColumnAsset(item: any) {
			const idToEdit = item.column_id

			const indexToEdit = this.asset.data.columns.findIndex((element) => element.column_id === idToEdit)
			if (indexToEdit !== -1) {
				this.asset.data.columns[indexToEdit] = item
			}
		},
		generateParams(assetType: string) {
			let dataColumn = this.asset.data.columns
			dataColumn.forEach((e) => {
				e.column_id = 0
			})
			dataColumn = dataColumn.filter((e) => e.is_ban !== true)

			const data: IRequestDataAssetTable = {
				table_name: this.asset.data.table_name,
				table_label: this.asset.data.table_label,
				delete_columns: this.asset.data.delete_columns,
				asset_type: assetType,
				asset_group: this.asset.data.asset_group,
				geom_type: this.asset.data.geom_type,
				line_color_code: this.asset.data.line_color_code,
				approver_id: this.asset.data.approver_id,
				viewer_id: this.asset.data.viewer_id,
				columns: dataColumn,
			}

			const params: IRequestAssetTable = {
				// อัปโหลดไฟล์
				icon_filepath: this.imageAsset.data?.file,
				icon_filepath_status: this.imageAsset.status,
				// ข้อมูล
				data,
			}

			return params
		},
		storeAsset() {
			switch (this.asset.data.geom_type) {
				case "km":
					return (this.asset.data.columns = [
						...[
							{
								column_id: 1,
								column_name: "the_geom",
								table_name_ref: "",
								component_title: "พิกัด",
								component_type: "geom-point",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},

							{
								column_id: 2,
								column_name: "km",
								table_name_ref: "",
								component_title: "กม. ที่ตั้ง",
								component_type: "text",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
						],
						...this.asset.data.columns.filter((item: any) => item.is_ban !== true),
					])
				case "km_range":
					return (this.asset.data.columns = [
						...[
							{
								column_id: 1,
								column_name: "the_geom",
								table_name_ref: "",
								component_title: "พิกัด",
								component_type: "geom-line",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},

							{
								column_id: 2,
								column_name: "km_start",
								table_name_ref: "",
								component_title: "กม. เริ่มต้น",
								component_type: "text",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
							{
								column_id: 3,
								column_name: "km_end",
								table_name_ref: "",
								component_title: "กม. สิ้นสุด",
								component_type: "text",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
						],
						...this.asset.data.columns.filter((item: any) => item.is_ban !== true),
					])
				case "point":
					return (this.asset.data.columns = [
						...[
							{
								column_id: 1,
								column_name: "the_geom",
								table_name_ref: "",
								component_title: "พิกัด",
								component_type: "geom-point",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
							{
								column_id: 2,
								column_name: "latitude",
								table_name_ref: "",
								component_title: "Latitude",
								component_type: "text-number",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
							{
								column_id: 3,
								column_name: "longitude",
								table_name_ref: "",
								component_title: "Longitude",
								component_type: "text-number",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
							{
								column_id: 4,
								column_name: "altitude",
								table_name_ref: "",
								component_title: "Altitude",
								component_type: "text-number",
								is_required: true,
								is_visible_view: true,
								is_visible_edit: true,
								is_visible_report: true,
								is_ban: true,
							},
						],
						...this.asset.data.columns.filter((item: any) => item.is_ban !== true),
					])
				default:
					break
			}
		},
		async assetCreate(assetType: string) {
			// Loading
			this.loading = true
			const params = this.generateParams(assetType)

			const assetTableService = new AssetTableService()
			const res = await assetTableService.post(params)
			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
				return res
			} else {
				return res
			}
		},
		generatePaste(event: ClipboardEvent) {
			event.preventDefault()
			let pastedText = event.clipboardData ? event.clipboardData.getData("text") : ""

			pastedText = pastedText.replace(/[^a-zA-Z0-9_]/g, "").toLowerCase()

			this.asset.data.table_name += pastedText
		},
	},
	getters: {},
})
