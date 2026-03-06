<script setup lang="ts">
import { IDefaultUsersData } from "../infrastructure"
import { useUserListStore } from "../store/UserListStore"
import { THeader } from "~/core/shared/types/Datatable"
import ServerSideDataTable from "~~/core/modules/common/datatable/ui/ServerSideDataTable.vue"

const store = useUserListStore()
useStoreLifecycle(store)

onMounted(() => {})

const headers: THeader[] = [
	{ text: "ลำดับ", value: "no", width: 10 },
	{ text: "ชื่อ-นามสกุล", value: "full_name", width: 200 },
	{ text: "หน่วยงานที่รับผิดชอบ", value: "ref_user_owner_id", width: 150 },
	{ text: "หมวด", value: "ref_depot_id", width: 150 },
	{ text: "เบอร์โทรศัพท์", value: "tel", width: 150 },
	{ text: "อีเมล", value: "email", width: 200 },
	{ text: "สิทธิ์การใช้งาน", value: "roles", width: 200 },
	{ text: "สถานะการใช้งาน", value: "status", width: 150 },
	{ text: "จัดการ", value: "operation", width: 150 },
]

const toggleStatus = (status: boolean) => {
	return status ? "เปิดการใช้งาน" : "ปิดการใช้งาน"
}

const toggleColor = (status: boolean) => {
	if (status) {
		return "badge badge-light-success text-success"
	} else {
		return " badge badge-light-danger text-danger"
	}
}

// เพิ่มข้อมูล
// const createItem = () => {
// 	return navigateTo(`/users/users/create`)
// }

// แก้ไขข้อมูล
const editItem = (id: number) => {
	return navigateTo(`/users/users/${id}/edit`)
}

// ลบข้อมูล
const dataTable: Ref = ref()
const deleteItem = (item: IDefaultUsersData) => {
	useDeleteItem({
		name: item.firstname,
		url: `/users/${item.id}`,
		callBack: function () {
			dataTable.value.loadData()
		},
	})
}

const onSearch = async () => {
	await dataTable.value.searchData(store.updateParams())
}

watch(
	() => store.params.ref_user_owner_id,
	(value) => {
		if (value !== 3) {
			store.params.ref_depot_id = null
		}
	}
)

onUnmounted(() => {
	store.$reset()
})
</script>

<template>
	<div class="row">
		<div class="col-md-4 col-12 mb-2">
			<VTextInput v-model="store.params.fullname" label="ชื่อ-นามสกุล" name="name" />
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VTextInput v-model="store.params.username" label="ชื่อผู้ใช้งาน" name="email" />
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params.status"
				label="สถานะการใช้งาน"
				:options="[
					{ label: 'เปิดใช้งาน', value: 1 },
					{ label: 'ปิดใช้งาน', value: 2 },
				]"
				name="status"
			/>
		</div>
	</div>
	<div class="row mb-3">
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params.permission"
				label="สิทธิ์การใช้งาน"
				:options="toOptions(store.roles)"
				name="permission"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params.ref_user_owner_id"
				:options="store.getUserOwnersOption"
				:label="'หน่วยงานที่รับผิดชอบ'"
				:name="'ref_user_owner_id'"
				:close-on-select="true"
			/>
		</div>
		<div class="col-md-4 col-12 mb-2">
			<VSelect
				v-model="store.params.ref_depot_id"
				:options="store.getDepotOption"
				:label="'หมวด'"
				:name="'ref_depot_id'"
				:disabled="store.params.ref_user_owner_id !== 3"
				:close-on-select="true"
			/>
		</div>
	</div>
	<div class="row mb-3">
		<div class="col-12 mb-2 text-end align-self-end align-items-end">
			<BtnSearch @click="onSearch" />

			<NuxtLink :to="'/users/users/create'" class="btn btn-outline btn-outline-primary ms-4"> เพิ่มข้อมูล </NuxtLink>
		</div>
	</div>
	<div class="row">
		<div class="col-xl-12">
			<ServerSideDataTable ref="dataTable" :headers="headers" url="users">
				<template #item-no="{ item }">
					<div class="text-center">{{ item.no }}</div>
				</template>
				<template #item-full_name="{ item }">
					<div class="text-start">
						{{ item.firstname || item.lastname !== "" ? `${item.firstname} ${item.lastname}` : "-" }}
					</div>
				</template>
				<template #item-ref_user_owner_id="{ item }">
					<div class="text-center">
						{{ item.ref_user_owner.email === "" ? "-" : item.ref_user_owner.email }}
					</div>
				</template>
				<template #item-ref_depot_id="{ item }">
					<div class="text-center">
						{{ item.ref_depot.name === "" ? "-" : item.ref_depot.name }}
					</div>
				</template>
				<template #item-email="{ item }">
					<div class="text-right">{{ item.email || "-" }}</div>
				</template>
				<template #item-tel="{ item }">
					<div class="text-right">{{ item.tel === "" ? "-" : item.tel }}</div>
				</template>
				<template #item-roles="{ item }">
					<ul v-if="item.roles.length > 0" class="d-flex flex-column justify-content-center my-2">
						<li v-for="(role, index) of item.roles" :key="index" class="text-start">{{ role.name }}</li>
					</ul>
					<div v-else class="text-center">-</div>
				</template>
				<template #item-status="{ item }">
					<div class="text-center" :class="toggleColor(item.status)">{{ toggleStatus(item.status) }}</div>
				</template>
				<template #item-operation="{ item }">
					<BtnEdit @click="editItem(item.id)" />
					<BtnDelete @click="deleteItem(item)" />
				</template>
			</ServerSideDataTable>
		</div>
	</div>

	<!-- Modal -->
	<InAssetCreate ref="modalCreate" :data-table="dataTable" />
	<InAssetEdit ref="modalEdit" :data-table="dataTable" />
</template>

<style scoped></style>
