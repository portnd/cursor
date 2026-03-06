import { IRequestAssetTable, IAssetTable, AssetTableService, IRequestDataAssetTable } from "../infrastructure"
import { IFile } from "~~/core/shared/types/File"

interface IState {
	id: string
	imageAsset: IFile
	asset: IAssetTable
	loading: boolean
}

export const useAssetTableEditStore = defineStore("setting/asset-table/edit", {
	state: (): IState => ({
		id: "",
		imageAsset: {} as IFile,
		asset: {
			icon_filepath: "",
			delete_columns: [],
			asset_id: 0,
			table_name: "",
			table_label: "",
			asset_type: "",
			asset_group: { id: "", can_delete: false, name: "" },
			geom_type: "km",
			line_color_code: "",
			approver: [],
			viewer: [],
			columns: [],
		},
		loading: false,
	}),
	actions: {
		addColumnAsset(item: any) {
			item.column_id = this.asset.columns.length + 1
			item.is_new = true
			item.is_edit = false
			this.asset.columns.push(item)
		},
		deleteColumnAsset(item: any) {
			if (this.asset.delete_columns === undefined) {
				this.asset.delete_columns = []
			}

			if (item.column_id !== 0) {
				this.asset.delete_columns.push(item.column_id)
			}

			const idToDelete = item.column_id

			const indexToDelete = this.asset.columns.findIndex((element) => element.column_id === idToDelete)
			if (indexToDelete !== -1) {
				this.asset.columns.splice(indexToDelete, 1)
			}
		},
		editColumnAsset(item: any) {
			const idToEdit = item.column_id
			item.is_new = true
			item.is_edit = true

			const indexToEdit = this.asset.columns.findIndex((element) => element.column_id === idToEdit)
			if (indexToEdit !== -1) {
				this.asset.columns[indexToEdit] = item
			}
		},
		generateParams(assetType: string) {
			let dataColumn = this.asset.columns
			dataColumn.forEach((e) => {
				if (e.is_new === true && e.is_edit === false) {
					e.column_id = 0
				}
			})
			dataColumn = dataColumn.filter((e) => e.is_ban !== true)
			dataColumn = dataColumn.filter((e) => e.is_new !== false)

			const data: IRequestDataAssetTable = {
				table_name: this.asset.table_name,
				table_label: this.asset.table_label,
				delete_columns: this.asset.delete_columns,
				asset_type: assetType,
				asset_group: this.asset.asset_group.id,
				geom_type: this.asset.geom_type,
				line_color_code: this.asset.line_color_code,
				approver_id: this.asset.approver,
				viewer_id: this.asset.viewer,
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
		async getAsset(id: string) {
			this.id = id
			// Loading
			this.loading = true
			const assetTableService = new AssetTableService()
			const res = await assetTableService.get(this.id)

			// Loading
			this.loading = false

			if (res.status === false) {
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				this.asset = res.data as IAssetTable
				this.asset.columns = this.asset.columns.filter((item) => item.component_type !== "hidden")
				this.asset.columns = this.asset.columns.filter((item) => item.column_name !== "the_geom_camera")
				for (let i = 0; i < this.asset.columns.length; i++) {
					if (
						this.asset.columns[i].column_name === "the_geom" ||
						this.asset.columns[i].column_name === "km" ||
						this.asset.columns[i].column_name === "km_start" ||
						this.asset.columns[i].column_name === "km_end" ||
						this.asset.columns[i].column_name === "latitude" ||
						this.asset.columns[i].column_name === "longitude" ||
						this.asset.columns[i].column_name === "altitude"
					) {
						this.asset.columns[i].is_ban = true
					} else {
						this.asset.columns[i].is_new = false
					}
				}
				if (res.data.geom_type === 1) {
					this.asset.geom_type = "km"
				} else if (res.data.geom_type === 2) {
					this.asset.geom_type = "km_range"
				} else if (res.data.geom_type === 3) {
					this.asset.geom_type = "point"
				}
				this.asset.asset_group.id = res.data.asset_group.id as string
				this.asset.viewer = res.data.viewer?.map((e: any) => e.id)
				this.asset.approver = res.data.approver?.map((e: any) => e.id)
			}
		},
		async editAsset(assetType: string) {
			const columns = this.asset.columns
			// Loading
			this.loading = true

			const params = this.generateParams(assetType)
			const assetTableService = new AssetTableService()
			const res = await assetTableService.put(this.id, params)
			// Loading
			this.loading = false

			if (res.status === false) {
				this.asset.columns = columns
				useHandlerError(res.code, res.error, { showAlert: true })
			} else {
				return res
			}
		},
	},
	getters: {},
})
