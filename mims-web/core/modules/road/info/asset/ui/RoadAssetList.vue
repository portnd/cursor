<script setup lang="ts">
import { IRoadAssetItem } from "../infrastructure"
import { useRoadAssetCreateStore, useRoadAssetEditStore, useRoadAssetListStore } from "../store"
import RoadReason from "../../reason/ui"
// import { useRoadListStore } from "../../../roadList/store"
import { RoadAssetCreate, RoadAssetEdit } from "./index"
import { useRoadTitleStore } from "~/core/modules/common/roadTitle/store"
import RoadTitle from "~~/core/modules/common/roadTitle/ui"

import { useInitUserStore } from "~/core/modules/initUser/store/InitUserStore"

const initUserStore = useInitUserStore()

const props = defineProps({
	assetType: {
		type: String,
		default: "",
	},
})

// interface IStateParams {
// 	keyword: string
// 	road_group_id: string[]
// 	road_section_id: number[]
// 	km_start: string
// 	km_end: string
// 	depot_code: number[]
// 	ref_surface_id: number[]
// }

const store = useRoadAssetListStore()
const roadTitleStore = useRoadTitleStore()
const RoadAssetCreateStore = useRoadAssetCreateStore()
const RoadAssetEditStore = useRoadAssetEditStore()
// const RoadListStore = useRoadListStore()

useStoreLifecycle([store, roadTitleStore, RoadAssetCreateStore, RoadAssetEditStore])

const route = useRoute()
const canEdit = ref<boolean>()

onMounted(async () => {
	canEdit.value =
		initUserStore.accessPermissions[IUserRolesAccess.manage_road_in_assets] ||
		initUserStore.getIsOwnerManagePermission(
			initUserStore.accessPermissions[IUserRolesAccess.manage_owner_road_in_assets],
			roadTitleStore.data.ref_depot.id
		)

	store.canEdit = canEdit.value
	const roadId = route.params.roadId
	const refAssetTableId = route.params.assetId
	store.assetType = props.assetType
	store.id = Number(roadId)
	store.refAssetTableId = Number(refAssetTableId)
	const res = await store.getRevision()
	await store.getDataTable(res ?? 1)
})

onBeforeMount(() => {
	window.addEventListener("beforeunload", () => {
		localStorage.removeItem("search")
		// roadTitleStore.params = {} as IStateParams
	})
})

// ซ่อน/แสดง แผนที่
const mapShow = ref({
	collapsed: false,
})

const activeRowNumber = ref(1)
const selectRow = (item: IRoadAssetItem) => {
	store.setLocation(item)

	store.locationCurrent = item.no

	// เพิ่ม class ใส่ใน row นั้น ๆ
	activeRowNumber.value = Number(item.no) ?? 0
}

watch(
	() => store.map,
	(_) => {
		store.lineString = {
			geom_cl: roadTitleStore.data.road_info?.the_geom,
			road_color_code: roadTitleStore.data.road_info?.road_color_code,
		}
		store.draw_asset_geom()
		store.setLocation(store.dataGeom[0])
	}
)

watch(
	() => store.dataAsset.road_assets,
	() => {
		if (store.map) {
			store.draw_asset_geom()
			// store.setLocation(newdata)

			const currentZoom = store.map.zoom()
			store.map.zoom(currentZoom + 1)
			store.map.zoom(currentZoom)
		}
	}
)

const locationMap = ref()

const dataTable = ref()
const modalCreate: Ref = ref()
const modalEdit: Ref = ref()
const modalReason = ref()

const startPoint = ref()
const polyline = ref()
const endPoint = ref()
const markend = ref()

const pushMap = () => {
	store.map.Overlays.add(polyline.value)
}

const createItem = async () => {
	if (store.assetDetail.geom_type_id === 2) {
		if (!startPoint.value) {
			const iconPosition = locationMap.value.location() // ตำแหน่งของไอคอน
			// @ts-ignore
			const longdo = window.longdo
			const marker = new longdo.Marker(
				{ lon: iconPosition.lon, lat: iconPosition.lat },
				{
					icon: {
						html: `<img src="/images/icons/png/location-pin.png" style="width:25px; height:25px"></img>`,
						offset: { x: 12.5, y: 20.5 },
					},
					visible: true,
					visibleRange: { min: 13, max: 20 },
				}
			)
			startPoint.value = { lon: iconPosition.lon, lat: iconPosition.lat }
			store.map.Overlays.add(marker)
			store.map.zoom(15)
			store.map.zoom(16)
		} else {
			const iconPosition = locationMap.value.location() // ตำแหน่งของไอคอน
			// @ts-ignore
			const longdo = window.longdo
			const marker = new longdo.Marker(
				{ lon: iconPosition.lon, lat: iconPosition.lat },
				{
					icon: {
						html: `<img src="/images/icons/png/location-pin.png" style="width:25px; height:25px"></img>`,
						offset: { x: 12.5, y: 20.5 },
					},
					visible: true,
					visibleRange: { min: 13, max: 20 },
				}
			)
			endPoint.value = { lon: iconPosition.lon, lat: iconPosition.lat }

			polyline.value = new longdo.Polyline([startPoint.value, endPoint.value], {
				lineColor: store.assetDetail.color,
			})
			store.map.Overlays.add(marker)
			store.map.Overlays.add(polyline.value)
			const data = []
			data.push(startPoint.value)
			data.push(endPoint.value)
			RoadAssetCreateStore.location = data
			RoadAssetCreateStore.id = store.id
			RoadAssetCreateStore.refAssetTableId = store.refAssetTableId
			RoadAssetCreateStore.revisionSelectId = store.revisionSelectId
			RoadAssetCreateStore.revisionIdParent = store.revisionIdParent

			await RoadAssetCreateStore.getTemplateData(RoadAssetCreateStore.location, store.assetDetail.geom_type_id)
			modalCreate.value.showModal(store.refAssetTableId, RoadAssetCreateStore.template)
			startPoint.value = null
			endPoint.value = null
		}
		store.map.Event.bind("location", function () {
			store.map.Overlays.remove(polyline.value)
			if (startPoint.value) {
				// @ts-ignore
				const longdo = window.longdo
				markend.value = store.map.location()
				const marker = new longdo.Marker(markend.value, {
					icon: {
						html: `<p></p>`,
						offset: { x: 12.5, y: 20.5 },
					},
					visible: true,
					draggable: true,
					visibleRange: { min: 13, max: 20 },
				})
				store.map.Overlays.add(marker)
				startPoint.value = { lon: startPoint.value.lon, lat: startPoint.value.lat }
				endPoint.value = { lon: markend.value.lon, lat: markend.value.lat }
				polyline.value = new longdo.Polyline([startPoint.value, endPoint.value], {
					lineStyle: longdo.LineStyle.Dashed,
					pointer: true,
					lineWidth: 4,
					lineColor: store.assetDetail.color,
				})
				pushMap()
			}
		})
	} else {
		const location = locationMap.value.location()
		const data = []
		data.push(location)
		RoadAssetCreateStore.location = data
		RoadAssetCreateStore.id = store.id
		RoadAssetCreateStore.refAssetTableId = store.refAssetTableId
		RoadAssetCreateStore.revisionSelectId = store.revisionSelectId
		RoadAssetCreateStore.revisionIdParent = store.revisionIdParent

		await RoadAssetCreateStore.getTemplateData(RoadAssetCreateStore.location, store.assetDetail.geom_type_id)
		modalCreate.value.showModal(store.refAssetTableId, RoadAssetCreateStore.template)
	}
}

const editItem = async (item: any) => {
	item.refAssetTableId = Number(route.params.assetId)
	RoadAssetEditStore.refAssetTableId = item.refAssetTableId
	RoadAssetEditStore.revisionSelectId = store.revisionSelectId
	RoadAssetEditStore.revisionIdParent = store.revisionIdParent
	RoadAssetEditStore.assetObjectId = item.raw_data.asset_object_id
	const res = await RoadAssetEditStore.getTemplateData()
	modalEdit.value.showModal(item, res?.data)
}

const showReason = () => {
	modalReason.value.showModal()
}

const showHideIcon = ref("fi fi-sr-eye")
const showHideMarker = () => {
	if (store.showVisible) {
		store.map.Overlays.clear()
		store.showVisible = !store.showVisible
		showHideIcon.value = "fi fi-sr-eye-crossed"
	} else {
		store.draw_asset_geom()
		store.dataGeom.forEach((e) => {
			if (e.no === store.locationCurrent) {
				store.setLocation(e)
			}
		})

		store.showVisible = !store.showVisible
		showHideIcon.value = "fi fi-sr-eye"
	}
}

watch(
	() => store.dataGeom,
	(newValue, oldValue) => {
		const diff = newValue.filter((a1) => !oldValue.find((a2) => a2.km === a1.km))
		store.setLocation(diff[0])
	}
)

const generateOptionTable = () => {
	const data: { id: number; name: string }[] = []
	store.dataRevision?.map((e) =>
		data.push({ id: e.id, name: `${formatDate(e.updated_date, "วันที่ dd mmm yyyy", true)} (${e.status})` })
	)
	return data
}

const deleteItem = (item: any) => {
	useDeleteItem({
		name: "",
		url: `/roads/asset/${store.revisionSelectId}/table/${store.refAssetTableId}/asset_object/${item.raw_data.asset_object_id}`,
		callBack: async () => {
			const res = await store.getRevision()
			await store.getDataTable(res ?? 1)
			store.map.Overlays.clear()
			store.draw_asset_geom()
			store.setLocation(store.dataGeom[0])
		},
	})
}

const deleteRevision = () => {
	useDeleteItem({
		name: `${formatDate(store.dataAsset.updated_date, "วันที่ dd mmm yyyy เวลา HH:ii น.", true)}`,
		url: `/roads/asset_delete/${store.revisionIdParent}`,
		callBack: async () => {
			const res = await store.getRevision()
			await store.getDataTable(res ?? 1)
		},
	})
}

const cancelRevision = () => {
	showAlert({
		title: `ยกเลิกฉบับร่าง`,
		html: `คุณต้องการยกเลิกฉบับร่าง <b class="fw-semibold text-danger">${formatDate(
			store.dataAsset.updated_date,
			"วันที่ dd mmm yyyy เวลา HH:ii น.",
			true
		)}</b> ใช่หรือไม่`,
		type: "question",
		callBack: function () {
			store.putCancel()
			store.resetStore()
		},
	})
}

const saveRevision = () => {
	showAlert({
		title: `ต้องการบันทึกการแก้ไข`,
		html: `คุณต้องการบันทึกการแก้ไข <b class="fw-semibold text-danger">${formatDate(
			store.dataAsset.updated_date,
			"วันที่ dd mmm yyyy เวลา HH:ii น.",
			true
		)}</b> ใช่หรือไม่`,
		type: "question",
		callBack: async () => {
			const res = await store.putConfirm()
			if (res?.status) {
				useHandlerSuccess(res.code, {
					showAlert: true,
					fn: async () => {
						const res = await store.getRevision()
						await store.getDataTable(res ?? 1)
					},
				})
			}
		},
	})
}

const onCancel = () => {
	store.map.Overlays.clear(polyline.value)
	store.draw_asset_geom()
	store.dataGeom.forEach((e) => {
		if (e.no === store.locationCurrent) {
			store.setLocation(e)
		}
	})
	return undefined
}

const noImage = ref<any[]>([])
const handleImageError = (e: any, data: any) => {
	if (e.type === "error") {
		noImage.value.push({ no: data.no, index: data.index })
	}
}

onUnmounted(() => {
	store.resetStore()
})
</script>

<template>
	<div class="row">
		<div class="col-12" :class="!mapShow.collapsed ? 'col-lg-8 col-xl-7' : 'col-lg-12 col-xl-12'">
			<div class="card card-rounded">
				<RoadTitle :road-id="Number(route.params.roadId)" />
			</div>
			<div class="card card-rounded p-5 mt-5">
				<div class="row">
					<div class="col-12 col-md-5">
						<VSelect
							v-model="store.revisionSelectId"
							:options="toOptions(generateOptionTable())"
							label="วันที่ปรับปรุงข้อมูล"
							name="lane"
							:can-clear="false"
							:can-deselect="false"
							placeholder="เลือก"
							@update:model-value="(e:any) => store.getDataTable(e)"
						/>
					</div>

					<div
						v-if="store.dataAsset.status_code === 'T' && store.statusLock && store.revisionSelectId !== 0 && canEdit"
						class="align-self-end text-end col-12 col-md-7 mt-md-0 mt-2 mb-1"
					>
						<button
							type="button"
							class="btn btn-light rounded-4 px-7 py-0 me-3 fw-semibold text-black fs-6"
							@click="cancelRevision"
						>
							ยกเลิกฉบับร่าง
						</button>
						<button type="button" class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6" @click="saveRevision">
							บันทึกการแก้ไข
						</button>
					</div>
					<div
						v-else-if="store.dataAsset.status_code !== 'I' && store.revisionSelectId !== 0 && canEdit"
						class="align-self-end text-end col-12 col-md-7 mt-md-0 mt-2 mb-1"
					>
						<button
							v-show="store.dataAsset.status_code === 'R'"
							type="button"
							class="btn btn-light rounded-4 px-7 py-0 me-3 fw-semibold text-black fs-6"
							@click="showReason"
						>
							เหตุผลการส่งกลับแก้ไข
						</button>
						<button type="button" class="btn btn-primary rounded-4 px-6 py-3 fw-semibold fs-6" @click="deleteRevision">
							ลบข้อมูล
						</button>
					</div>
				</div>
				<div class="row mt-3 mb-3">
					<div v-show="store.dataAsset?.updated_by" class="col-12">
						<span
							class="badge rounded-1 fs-7 fw-semibold px-4 py-2"
							:class="`badge-status-${store.dataAsset?.status_code}`"
						>
							{{ store.dataAsset?.status }} <i v-show="store.statusLock" class="fi fi-rr-lock ms-1"></i>
						</span>

						<label class="text-gray-900 ms-2 me-2 fs-7 mt-2">ปรับปรุงข้อมูลโดย</label>
						<VUser
							:label="`${store.dataAsset?.updated_by?.full_name}`"
							:name="`${store.dataAsset?.updated_by?.full_name}`"
							:role="`${store.dataAsset?.updated_by?.department.name}`"
						/>
						<label class="text-gray-900 ms-2 fs-7 mt-2"
							>เมื่อ{{ formatDate(store.dataAsset.updated_date, "วันที่ dd mmm yyyy เวลา HH:ii น.  ", true) }}</label
						>
					</div>
				</div>

				<VDatatable
					ref="dataTable"
					:headers="store.headers"
					:items="store.dataGeom"
					active-item-class-name="active"
					:active-row-number="activeRowNumber"
				>
					<!-- begin::Items -->
					<template v-for="(header, index) in store.headers" #[`item-${header.value}`]="{ item }" :key="index">
						<div class="cursor-pointer" @click="selectRow(item)">
							<div
								v-if="header.value.includes('_filepath') && item[header.value]"
								v-viewer
								class="symbol symbol-40px cursor-pointer"
							>
								<img :src="`${item[header.value]}`" class="rounded-1" @error="(e) => handleImageError(e, index)" />
							</div>
							<div v-else-if="header.value.includes('action') && canEdit" v-show="store.dataAsset.status_code !== 'I'">
								<BtnEdit @click="editItem(item)" />
								<BtnDelete @click="deleteItem(item)" />
							</div>
							<div
								v-else-if="header.value.includes('_image') && item[header.value]"
								v-viewer
								class="symbol symbol-40px cursor-pointer"
							>
								<img
									v-if="item[header.value].filepath"
									:id="`${item.no}_${index}`"
									:src="`${item[header.value].filepath}`"
									class="rounded-1"
									@error="(e) => handleImageError(e, index)"
								/>
								<span v-else>-</span>
							</div>
							<div v-else-if="header.type === 'datepicker'">
								{{
									/^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)$/.test(item[header.value]) && item[header.value] !== ""
										? buddhistFormatDate(item[header.value], "dd mmm yyyy") === ""
											? "-"
											: buddhistFormatDate(item[header.value], "dd mmm yyyy")
										: item[header.value] === ""
										? "-"
										: item[header.value]
								}}
							</div>
							<div v-else-if="header.type === 'text-year'">
								{{ item[header.value] === "" ? "-" : Number(item[header.value]) + 543 }}
							</div>

							<div v-else-if="header.type === 'text-number'">
								<div v-if="header.value === 'latitude' || header.value === 'longitude' || header.value === 'altitude'">
									{{ item[header.value] }}
								</div>
								<div v-else>
									{{ item[header.value] === "0" ? "-" : toNumber(Number(item[header.value])) }}
								</div>
							</div>
							<div v-else-if="header.type === 'text-km'">
								{{ item[header.value] === "" ? "-" : convertMeterToKm(item[header.value]) }}
							</div>
							<div v-else>
								{{ item[header.value] === 0 || item[header.value] ? item[header.value] : "-" }}
							</div>
						</div>
					</template>
					<!-- end::Items -->
				</VDatatable>
			</div>
		</div>
		<div class="col-12 map-sticky" :class="!mapShow.collapsed ? 'col-lg-4 col-xl-5' : 'map-collapsed'">
			<div class="widget">
				<KeepAlive>
					<VMap
						v-model="mapShow"
						:loading="false"
						height="calc(97vh)"
						:is-sticky="true"
						@map="($event:any) => {locationMap = $event ,store.setMap($event)}"
					>
						<template #bottom-end>
							<div class="d-flex">
								<div class="button-action me-5 cursor-pointer" @click="showHideMarker">
									<div class="rounded-circle text-primary text-center bg-white w-50px h-50px pt-4">
										<i class="lh-0 fs-1" :class="showHideIcon"></i>
									</div>
								</div>
								<div
									v-if="canEdit && store.dataAsset.status_code !== 'I'"
									class="button-action cursor-pointer"
									@click="createItem"
								>
									<div class="rounded-circle text-primary text-center bg-white w-50px h-50px">
										<div class="button">
											<i class="fi fi-sr-marker fs-2"> </i>
											<i class="fi fi-br-plus-small fs-5 icon-plus"></i>
										</div>
									</div>
								</div>
							</div>
						</template>
					</VMap>
				</KeepAlive>
			</div>
		</div>
	</div>

	<!-- Modal -->
	<RoadAssetCreate ref="modalCreate" :props-store="store" :asset-type="store.assetType" :on-cancel="onCancel" />
	<RoadAssetEdit ref="modalEdit" :props-store="store" :asset-type="store.assetType" :on-cancel="onCancel" />
	<RoadReason ref="modalReason" :message="store?.dataAsset.reject_reason" />
</template>

<style scoped></style>
