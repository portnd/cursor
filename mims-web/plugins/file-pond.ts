import vueFilePond from "vue-filepond"
import "filepond/dist/filepond.min.css"

import "filepond-plugin-image-preview/dist/filepond-plugin-image-preview.min.css"
import "filepond-plugin-image-overlay/dist/filepond-plugin-image-overlay.min.css"
import "filepond-plugin-pdf-preview/dist/filepond-plugin-pdf-preview.min.css"
import "filepond-plugin-get-file/dist/filepond-plugin-get-file.min.css"
import "filepond-plugin-media-preview/dist/filepond-plugin-media-preview.min.css"
import "filepond-plugin-file-poster/dist/filepond-plugin-file-poster.css"
import "filepond-plugin-image-edit/dist/filepond-plugin-image-edit.css"

import FilePondPluginFileValidateType from "filepond-plugin-file-validate-type"
import FilePondPluginImagePreview from "filepond-plugin-image-preview"
import FilePondPluginImageExifOrientation from "filepond-plugin-image-exif-orientation"
import FilePondPluginFileValidateSize from "filepond-plugin-file-validate-size"
// @ts-ignore
import FilePondPluginImageOverlay from "filepond-plugin-image-overlay"
import FilePondPluginGetFile from "filepond-plugin-get-file"
// @ts-ignore
import FilePondPluginPdfPreview from "filepond-plugin-pdf-preview"
// @ts-ignore
import FilePondPluginMediaPreview from "filepond-plugin-media-preview"
import FilePondPluginImageValidateSize from "filepond-plugin-image-validate-size"
import FilePondPluginFilePoster from "filepond-plugin-file-poster"
import FilePondPluginImageEdit from "filepond-plugin-image-edit"
import FilePondPluginFileEncode from "filepond-plugin-file-encode"

const FilePond: any = vueFilePond(
	FilePondPluginFileValidateType,
	FilePondPluginImagePreview,
	FilePondPluginImageExifOrientation,
	FilePondPluginFileValidateSize,
	FilePondPluginImageOverlay,
	FilePondPluginGetFile,
	FilePondPluginPdfPreview,
	FilePondPluginMediaPreview,
	FilePondPluginImageValidateSize,
	FilePondPluginFilePoster,
	FilePondPluginImageEdit,
	FilePondPluginFileEncode
)

export default defineNuxtPlugin((nuxtApp) => {
	nuxtApp.vueApp.component("file-pond", FilePond)
})
