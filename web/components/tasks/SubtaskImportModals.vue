<template>
  <Teleport to="body">
    <!-- Google Slides -->
    <div
      v-if="showSlides"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto"
      @click.self="closeSlides"
    >
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="slidesStep === 'select' ? 'max-w-5xl' : 'max-w-xl'"
      >
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-purple-600/20 border border-purple-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-purple-400" fill="currentColor" viewBox="0 0 20 20"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import Slides → Sub-tasks</h2>
              <p class="text-xs text-gray-400">Under: {{ parentTitle || 'this task' }}</p>
            </div>
          </div>
          <button type="button" class="text-gray-500 hover:text-white shrink-0 ml-4" @click="closeSlides">✕</button>
        </div>
        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <template v-if="slidesStep === 'result' && slidesResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ slidesResult.presentation_title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ slidesResult.created_count }} sub-tasks จาก {{ slidesResult.slide_count }} slides</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div v-for="t in slidesResult.tasks" :key="t.id" class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(t.code) }}</span>
                <span class="text-gray-200 truncate">{{ t.title }}</span>
              </div>
            </div>
            <button type="button" class="w-full btn-primary py-2.5" @click="closeSlidesAfterSuccess">Done</button>
          </template>
          <template v-else-if="slidesStep === 'select' && slidesPreview">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ slidesPreview.presentation_title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">{{ slidesSelected.length }} / {{ slidesPreview.slides.length }} slides selected</p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" class="btn-ghost-sm" @click="slidesSelectAll">ทั้งหมด</button>
                <button type="button" class="btn-ghost-sm" @click="slidesDeselectAll">ยกเลิก</button>
                <button type="button" class="btn-ghost-sm text-purple-400" @click="slidesSelectOnlyNew">เฉพาะใหม่</button>
              </div>
            </div>
            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-3 text-left w-8" />
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-10">#</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[200px]">Task Title</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Assignee</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Est. min</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="s in slidesPreview.slides"
                    :key="s.index"
                    class="border-b border-gray-700/30 transition-colors"
                    :class="slidesSelected.includes(s.index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-3">
                      <input v-model="slidesSelected" type="checkbox" :value="s.index" class="rounded border-gray-500 bg-gray-700 text-purple-500 focus:ring-purple-500">
                    </td>
                    <td class="py-2 px-3 text-xs text-gray-400 font-mono">
                      {{ s.index }}
                      <span v-if="s.hidden" class="text-amber-400 ml-1 text-[10px]">ซ่อน</span>
                      <span v-else-if="(slidesPreview.already_imported_slide_indices || []).includes(s.index)" class="text-gray-500 ml-1 text-[10px]">นำเข้าแล้ว</span>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="slidesTriaged[s.index]"
                        v-model="slidesTriaged[s.index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-purple-500/60"
                        :disabled="!slidesSelected.includes(s.index)"
                      >
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="slidesTriaged[s.index]"
                        v-model="slidesTriaged[s.index].assignee_id"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!slidesSelected.includes(s.index)"
                      >
                        <option :value="null">— Unassigned —</option>
                        <option v-for="u in importAssignees" :key="u.id" :value="u.id">{{ u.display_name || u.email }}</option>
                      </select>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="slidesTriaged[s.index]"
                        v-model.number="slidesTriaged[s.index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!slidesSelected.includes(s.index)"
                        placeholder="0"
                      >
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="slidesTriaged[s.index]"
                        v-model="slidesTriaged[s.index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-purple-500/60"
                        :disabled="!slidesSelected.includes(s.index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="slidesErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ slidesErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="slidesImporting || slidesSelected.length === 0"
                @click="submitSlides"
              >
                <svg v-if="slidesImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ slidesImporting ? 'กำลัง import...' : `Import ${slidesSelected.length} Slides` }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="slidesStep = 'form'">กลับ</button>
            </div>
          </template>
          <template v-else>
            <div>
              <label class="label">Google Slides URL *</label>
              <input v-model="slidesUrl" type="url" class="input-field w-full" placeholder="https://docs.google.com/presentation/d/..." :disabled="slidesLoading">
              <p class="text-xs text-gray-500 mt-1">ต้องเปิดสิทธิ์ "Anyone with the link can view"</p>
            </div>
            <div>
              <label class="label">Default priority (new slides)</label>
              <select v-model="slidesFormPriority" class="input-field w-full" :disabled="slidesLoading">
                <option value="CRITICAL">CRITICAL</option>
                <option value="HIGH">HIGH</option>
                <option value="MEDIUM">MEDIUM</option>
                <option value="LOW">LOW</option>
              </select>
            </div>
            <div v-if="slidesErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ slidesErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="slidesLoading || !slidesUrl.trim()"
                @click="loadSlidesPreview"
              >
                <svg v-if="slidesLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ slidesLoading ? 'กำลังโหลด...' : 'โหลดรายการ slide' }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="closeSlides">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- PPTX -->
    <div
      v-if="showPptx"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto"
      @click.self="closePptx"
    >
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="pptxStep === 'select' ? 'max-w-5xl' : 'max-w-xl'"
      >
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-orange-600/20 border border-orange-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-orange-400" fill="currentColor" viewBox="0 0 24 24"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm4 18H6V4h7v5h5v11z"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import PPTX → Sub-tasks</h2>
              <p class="text-xs text-gray-400">Under: {{ parentTitle || 'this task' }}</p>
            </div>
          </div>
          <button type="button" class="text-gray-500 hover:text-white shrink-0 ml-4" @click="closePptx">✕</button>
        </div>
        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <template v-if="pptxStep === 'result' && pptxResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ pptxResult.title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ pptxResult.created_count }} sub-tasks จาก {{ pptxResult.page_count }} slides</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div v-for="t in pptxResult.tasks" :key="t.id" class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(t.code) }}</span>
                <span class="text-gray-200 truncate">{{ t.title }}</span>
              </div>
            </div>
            <button type="button" class="w-full btn-primary py-2.5" @click="closePptxAfterSuccess">Done</button>
          </template>
          <template v-else-if="pptxStep === 'select' && pptxPreview">
            <div class="flex items-center justify-between gap-3">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ pptxPreview.title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">{{ pptxSelected.length }} / {{ pptxPreview.slides.length }} slides selected</p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" class="btn-ghost-sm" @click="pptxSelectAll">ทั้งหมด</button>
                <button type="button" class="btn-ghost-sm" @click="pptxDeselectAll">ยกเลิก</button>
                <button type="button" class="btn-ghost-sm text-orange-400" @click="pptxSelectVisible">เฉพาะที่แสดง</button>
              </div>
            </div>
            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-3 text-left w-8" />
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-10">#</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[200px]">Task Title</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Assignee</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Est. min</th>
                    <th class="py-2 px-3 text-left text-xs text-gray-400 font-semibold w-32">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="s in pptxPreview.slides"
                    :key="s.index"
                    class="border-b border-gray-700/30 transition-colors"
                    :class="pptxSelected.includes(s.index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-3">
                      <input v-model="pptxSelected" type="checkbox" :value="s.index" class="rounded border-gray-500 bg-gray-700 text-orange-500 focus:ring-orange-500">
                    </td>
                    <td class="py-2 px-3 text-xs text-gray-400 font-mono">
                      {{ s.index }}
                      <span v-if="s.hidden" class="text-amber-400 ml-1 text-[10px]">ซ่อน</span>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="pptxTriaged[s.index]"
                        v-model="pptxTriaged[s.index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxSelected.includes(s.index)"
                      >
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="pptxTriaged[s.index]"
                        v-model="pptxTriaged[s.index].assignee_id"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxSelected.includes(s.index)"
                      >
                        <option :value="null">— Unassigned —</option>
                        <option v-for="u in importAssignees" :key="u.id" :value="u.id">{{ u.display_name || u.email }}</option>
                      </select>
                    </td>
                    <td class="py-2 px-2">
                      <input
                        v-if="pptxTriaged[s.index]"
                        v-model.number="pptxTriaged[s.index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxSelected.includes(s.index)"
                        placeholder="0"
                      >
                    </td>
                    <td class="py-2 px-2">
                      <select
                        v-if="pptxTriaged[s.index]"
                        v-model="pptxTriaged[s.index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-orange-500/60"
                        :disabled="!pptxSelected.includes(s.index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="pptxErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ pptxErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="pptxImporting || pptxSelected.length === 0"
                @click="submitPptx"
              >
                <svg v-if="pptxImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ pptxImporting ? 'กำลัง import...' : `Import ${pptxSelected.length} Slides` }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="pptxStep = 'form'">กลับ</button>
            </div>
          </template>
          <template v-else>
            <div
              class="relative border-2 border-dashed rounded-xl p-8 text-center transition-colors"
              :class="pptxDrag ? 'border-orange-500/60 bg-orange-500/5' : 'border-gray-600/60 hover:border-gray-500/60'"
              @dragover.prevent="pptxDrag = true"
              @dragleave="pptxDrag = false"
              @drop.prevent="onPptxDrop"
            >
              <p class="text-sm font-medium text-gray-300 mb-1">{{ pptxFile ? pptxFile.name : 'ลากไฟล์ .pptx มาวาง หรือ' }}</p>
              <p v-if="!pptxFile" class="text-xs text-gray-500 mb-3">รองรับไฟล์ .pptx</p>
              <p v-else class="text-xs text-green-400 mb-3">{{ (pptxFile.size / 1024 / 1024).toFixed(1) }} MB</p>
              <label class="btn-import-sm cursor-pointer inline-flex items-center gap-1.5">
                {{ pptxFile ? 'เปลี่ยนไฟล์' : 'เลือกไฟล์' }}
                <input type="file" accept=".pptx" class="sr-only" @change="onPptxFile">
              </label>
            </div>
            <div v-if="pptxErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ pptxErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="pptxLoading || !pptxFile"
                @click="loadPptxPreview"
              >
                <svg v-if="pptxLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ pptxLoading ? 'กำลังโหลด...' : 'โหลดรายการ slide' }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="closePptx">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- Google Sheets -->
    <div
      v-if="showSheets"
      class="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 p-4 overflow-y-auto"
      @click.self="closeSheets"
    >
      <div
        class="bg-gray-800 border border-gray-700 rounded-2xl shadow-2xl w-full my-auto flex flex-col max-h-[90vh]"
        :class="sheetsStep === 'select' ? 'max-w-6xl' : 'max-w-xl'"
      >
        <div class="flex items-center justify-between px-6 pt-5 pb-4 shrink-0 border-b border-gray-700/60">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-emerald-600/20 border border-emerald-500/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-emerald-400" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clip-rule="evenodd"/></svg>
            </div>
            <div>
              <h2 class="text-lg font-bold text-white">Import Sheets → Sub-tasks</h2>
              <p class="text-xs text-gray-400">Under: {{ parentTitle || 'this task' }}</p>
            </div>
          </div>
          <button type="button" class="text-gray-500 hover:text-white shrink-0 ml-4" @click="closeSheets">✕</button>
        </div>
        <div class="overflow-y-auto flex-1 px-6 py-5 space-y-4">
          <template v-if="sheetsStep === 'result' && sheetsResult">
            <div class="p-4 bg-green-900/20 border border-green-600/40 rounded-xl">
              <div class="flex items-center gap-2 mb-2">
                <svg class="w-5 h-5 text-green-400 shrink-0" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/></svg>
                <span class="text-green-400 font-semibold text-sm">Import สำเร็จ!</span>
              </div>
              <p class="text-gray-300 text-sm font-medium mb-1">{{ sheetsResult.sheet_title }}</p>
              <p class="text-gray-400 text-xs">สร้าง {{ sheetsResult.created_count }} sub-tasks</p>
            </div>
            <div class="space-y-1.5 max-h-48 overflow-y-auto pr-1">
              <div v-for="t in sheetsResult.tasks" :key="t.id" class="flex items-center gap-2 py-2 px-3 bg-gray-700/40 rounded-lg text-sm">
                <span class="text-xs font-mono text-gray-500 shrink-0">{{ taskCodeSuffix(t.code) }}</span>
                <span class="text-gray-200 truncate">{{ t.title }}</span>
              </div>
            </div>
            <button type="button" class="w-full btn-primary py-2.5" @click="closeSheetsAfterSuccess">Done</button>
          </template>
          <template v-else-if="sheetsStep === 'select' && sheetsPreview">
            <div class="flex items-center justify-between gap-3 flex-wrap">
              <div class="p-3 bg-gray-700/40 rounded-xl flex-1 min-w-0">
                <p class="text-sm font-medium text-white truncate">{{ sheetsPreview.sheet_title }}</p>
                <p class="text-xs text-gray-500 mt-0.5">{{ sheetsRowsSelected.length }} / {{ sheetsPreview.rows.length }} rows selected</p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <button type="button" class="btn-ghost-sm" @click="sheetsSelectAll">ทั้งหมด</button>
                <button type="button" class="btn-ghost-sm" @click="sheetsDeselectAll">ยกเลิก</button>
              </div>
            </div>
            <div class="overflow-x-auto border border-gray-700/60 rounded-xl">
              <table class="w-full text-sm">
                <thead>
                  <tr class="border-b border-gray-700/60 bg-gray-900/60">
                    <th class="py-2 px-2 text-left w-8" />
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-12">#</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[180px]">Title</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-36">Due</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[120px]">Status</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold min-w-[140px]">Notes</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-28">Est. min</th>
                    <th class="py-2 px-2 text-left text-xs text-gray-400 font-semibold w-28">Priority</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="r in sheetsPreview.rows"
                    :key="r.row_index"
                    class="border-b border-gray-700/30 transition-colors align-top"
                    :class="sheetsRowsSelected.includes(r.row_index) ? 'bg-gray-800/80' : 'bg-gray-900/40 opacity-50'"
                  >
                    <td class="py-2 px-2">
                      <input v-model="sheetsRowsSelected" type="checkbox" :value="r.row_index" class="rounded border-gray-500 bg-gray-700 text-emerald-500 focus:ring-emerald-500">
                    </td>
                    <td class="py-2 px-2 text-xs text-gray-400 font-mono">{{ r.row_index }}</td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsTriaged[r.row_index]"
                        v-model="sheetsTriaged[r.row_index].title"
                        type="text"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                      >
                    </td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsTriaged[r.row_index]"
                        v-model="sheetsTriaged[r.row_index].due_date"
                        type="date"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60 max-w-[9.5rem]"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                      >
                    </td>
                    <td class="py-2 px-1">
                      <select
                        v-if="sheetsTriaged[r.row_index]"
                        v-model="sheetsTriaged[r.row_index].status"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                      >
                        <option value="PENDING">PENDING</option>
                        <option value="IN_PROGRESS">IN_PROGRESS</option>
                        <option value="READY_FOR_TEST">READY_FOR_TEST</option>
                        <option value="READY_FOR_UAT">READY_FOR_UAT</option>
                        <option value="COMPLETED">COMPLETED</option>
                        <option value="CANCELLED">CANCELLED</option>
                      </select>
                    </td>
                    <td class="py-2 px-1">
                      <textarea
                        v-if="sheetsTriaged[r.row_index]"
                        v-model="sheetsTriaged[r.row_index].notes"
                        rows="2"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white placeholder-gray-500 focus:outline-none focus:border-emerald-500/60 resize-y min-h-[2.25rem]"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                      />
                    </td>
                    <td class="py-2 px-1">
                      <input
                        v-if="sheetsTriaged[r.row_index]"
                        v-model.number="sheetsTriaged[r.row_index].estimated_minutes"
                        type="number"
                        min="0"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-2 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                        placeholder="—"
                      >
                    </td>
                    <td class="py-2 px-1">
                      <select
                        v-if="sheetsTriaged[r.row_index]"
                        v-model="sheetsTriaged[r.row_index].priority"
                        class="w-full bg-gray-700/60 border border-gray-600/60 rounded-lg px-1 py-1 text-xs text-white focus:outline-none focus:border-emerald-500/60"
                        :disabled="!sheetsRowsSelected.includes(r.row_index)"
                      >
                        <option value="CRITICAL">CRITICAL</option>
                        <option value="HIGH">HIGH</option>
                        <option value="MEDIUM">MEDIUM</option>
                        <option value="LOW">LOW</option>
                      </select>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div v-if="sheetsErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ sheetsErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="sheetsImporting || sheetsRowsSelected.length === 0"
                @click="submitSheets"
              >
                <svg v-if="sheetsImporting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ sheetsImporting ? 'กำลัง import...' : `Import ${sheetsRowsSelected.length} rows` }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="sheetsStep = 'form'">กลับ</button>
            </div>
          </template>
          <template v-else>
            <div>
              <label class="label">Google Sheets URL *</label>
              <input v-model="sheetsUrl" type="url" class="input-field w-full" placeholder="https://docs.google.com/spreadsheets/d/..." :disabled="sheetsLoading">
              <p class="text-xs text-gray-500 mt-1">ต้องเปิดสิทธิ์ "Anyone with the link can view"</p>
            </div>
            <div v-if="sheetsErr" class="p-3 bg-red-900/30 border border-red-600 rounded-lg text-red-400 text-sm">{{ sheetsErr }}</div>
            <div class="flex gap-3">
              <button
                type="button"
                class="flex-1 btn-primary py-2.5 disabled:opacity-40 flex items-center justify-center gap-2"
                :disabled="sheetsLoading || !sheetsUrl.trim()"
                @click="loadSheetsPreview"
              >
                <svg v-if="sheetsLoading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/></svg>
                {{ sheetsLoading ? 'กำลังโหลด...' : 'โหลด preview' }}
              </button>
              <button type="button" class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-xl transition-colors" @click="closeSheets">Cancel</button>
            </div>
          </template>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useTasksApi } from '~/core/modules/tasks/infrastructure/tasks-api'
import { useAuth } from '~/composables/useAuth'
import { useTeamsApi } from '~/core/modules/teams/infrastructure/teams-api'
import { useTeamsStore } from '~/core/modules/teams/store/teams-store'

interface ImportAssignee { id: number; email: string; display_name: string; role: string }
interface TriagedSlide { title: string; assignee_id: number | null; estimated_minutes: number; priority: string }
interface SheetsTriagedRow {
  title: string
  priority: string
  estimated_minutes: number
  due_date: string
  status: string
  notes: string
}

const props = defineProps<{
  projectId: string
  parentTaskId: string
  parentTitle?: string
  epicId?: string | null
}>()

const emit = defineEmits<{ imported: [] }>()

const { fetchWithAuth, currentUser } = useAuth()
const tasksApi = useTasksApi()
const { getTeams } = useTeamsApi()
const teamsStore = useTeamsStore()

const importAssignees = ref<ImportAssignee[]>([])

async function ensureAssignees() {
  if (importAssignees.value.length > 0) return
  try {
    const role = (currentUser.value?.role || '').toUpperCase()
    if (role === 'PM') {
      await teamsStore.fetchTeamsFeatureEnabled()
      if (teamsStore.teamsFeatureEnabled) {
        const userId = currentUser.value?.user_id
        const teams = await getTeams()
        const myTeam = teams.find((t) => t.users?.some((u) => u.id === userId))
        importAssignees.value = (myTeam?.users ?? [])
          .filter((u) => ['DEV', 'PM', 'MANAGER', 'SUPPORT'].includes(u.role))
          .map((u) => ({ id: u.id, email: u.email, display_name: u.display_name, role: u.role }))
      } else {
        const res = await fetchWithAuth<{ data: ImportAssignee[] }>('/auth/users')
        importAssignees.value = (res.data ?? []).filter((u) => ['DEV', 'PM', 'MANAGER', 'SUPPORT'].includes(u.role))
      }
    } else {
      const res = await fetchWithAuth<{ data: ImportAssignee[] }>('/auth/users')
      importAssignees.value = (res.data ?? []).filter((u) => ['DEV', 'PM', 'MANAGER', 'SUPPORT'].includes(u.role))
    }
  } catch { /* optional */ }
}

function taskCodeSuffix(code?: string): string {
  if (!code) return '—'
  const parts = code.split('-')
  const n = Number(parts[parts.length - 1])
  if (!Number.isFinite(n)) return code
  return String(n).padStart(4, '0')
}

// ── Slides ────────────────────────────────────────────────
const showSlides = ref(false)
const slidesStep = ref<'form' | 'select' | 'result'>('form')
const slidesUrl = ref('')
const slidesFormPriority = ref('MEDIUM')
const slidesLoading = ref(false)
const slidesImporting = ref(false)
const slidesErr = ref('')
const slidesPreview = ref<{
  presentation_title: string
  slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[]
  /** Present when API returns dedupe hints (extends tasks-api preview type) */
  already_imported_slide_indices?: number[]
} | null>(null)
const slidesSelected = ref<number[]>([])
const slidesTriaged = ref<Record<number, TriagedSlide>>({})
const slidesResult = ref<{ created_count: number; slide_count: number; presentation_title: string; tasks: { id: string; title: string; code?: string }[] } | null>(null)

function openSlides() {
  slidesStep.value = 'form'
  slidesUrl.value = ''
  slidesFormPriority.value = 'MEDIUM'
  slidesErr.value = ''
  slidesPreview.value = null
  slidesSelected.value = []
  slidesTriaged.value = {}
  slidesResult.value = null
  showSlides.value = true
  ensureAssignees()
}

function closeSlides() {
  showSlides.value = false
}

function closeSlidesAfterSuccess() {
  const had = !!slidesResult.value
  showSlides.value = false
  slidesResult.value = null
  if (had) emit('imported')
}

async function loadSlidesPreview() {
  if (!slidesUrl.value.trim()) return
  slidesLoading.value = true
  slidesErr.value = ''
  try {
    const data = await tasksApi.previewGoogleSlides({ presentation_url: slidesUrl.value.trim() })
    const ext = data as typeof data & { already_imported_slide_indices?: number[] }
    slidesPreview.value = ext
    const already = new Set(ext.already_imported_slide_indices ?? [])
    slidesSelected.value = data.slides.filter((s) => !s.hidden && !already.has(s.index)).map((s) => s.index)
    const triaged: Record<number, TriagedSlide> = {}
    for (const s of data.slides) {
      const st = s.suggested_task_title?.trim()
      triaged[s.index] = {
        title: st || `Slide ${s.index}`,
        assignee_id: null,
        estimated_minutes: 0,
        priority: slidesFormPriority.value || 'MEDIUM',
      }
    }
    slidesTriaged.value = triaged
    slidesStep.value = 'select'
  } catch (e: any) {
    slidesErr.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
  } finally {
    slidesLoading.value = false
  }
}

function slidesSelectAll() {
  if (slidesPreview.value) slidesSelected.value = slidesPreview.value.slides.map((s) => s.index)
}
function slidesDeselectAll() {
  slidesSelected.value = []
}
function slidesSelectOnlyNew() {
  if (!slidesPreview.value) return
  const already = new Set(slidesPreview.value.already_imported_slide_indices ?? [])
  slidesSelected.value = slidesPreview.value.slides.filter((s) => !s.hidden && !already.has(s.index)).map((s) => s.index)
}

async function submitSlides() {
  slidesImporting.value = true
  slidesErr.value = ''
  try {
    const slides = slidesSelected.value.map((idx) => {
      const t = slidesTriaged.value[idx]
      return {
        slide_index: idx,
        title: t?.title || `Slide ${idx}`,
        assignee_id: t?.assignee_id ?? null,
        estimated_minutes: t?.estimated_minutes ?? 0,
        priority: t?.priority || 'MEDIUM',
      }
    })
    const payload: Record<string, unknown> = {
      presentation_url: slidesUrl.value.trim(),
      project_id: props.projectId,
      parent_id: props.parentTaskId,
      slides,
    }
    if (props.epicId) payload.epic_id = props.epicId
    slidesResult.value = await tasksApi.importGoogleSlides(payload as Parameters<typeof tasksApi.importGoogleSlides>[0])
    slidesStep.value = 'result'
  } catch (e: any) {
    slidesErr.value = e?.data?.message ?? e?.message ?? 'Import failed'
  } finally {
    slidesImporting.value = false
  }
}

// ── PPTX ───────────────────────────────────────────────────
const showPptx = ref(false)
const pptxStep = ref<'form' | 'select' | 'result'>('form')
const pptxFile = ref<File | null>(null)
const pptxDrag = ref(false)
const pptxLoading = ref(false)
const pptxImporting = ref(false)
const pptxErr = ref('')
const pptxPreview = ref<{ title: string; slides: { index: number; title: string; hidden?: boolean; suggested_task_title?: string }[] } | null>(null)
const pptxSelected = ref<number[]>([])
const pptxTriaged = ref<Record<number, TriagedSlide>>({})
const pptxResult = ref<{ created_count: number; page_count: number; title: string; tasks: { id: string; title: string; code?: string }[] } | null>(null)

function openPptx() {
  pptxStep.value = 'form'
  pptxFile.value = null
  pptxErr.value = ''
  pptxPreview.value = null
  pptxSelected.value = []
  pptxTriaged.value = {}
  pptxResult.value = null
  showPptx.value = true
  ensureAssignees()
}

function closePptx() {
  showPptx.value = false
}

function closePptxAfterSuccess() {
  const had = !!pptxResult.value
  showPptx.value = false
  pptxResult.value = null
  if (had) emit('imported')
}

function onPptxFile(e: Event) {
  const input = e.target as HTMLInputElement
  if (input.files?.[0]) {
    pptxFile.value = input.files[0]
    pptxErr.value = ''
  }
}

function onPptxDrop(e: DragEvent) {
  pptxDrag.value = false
  const file = e.dataTransfer?.files[0]
  if (file && (file.name.endsWith('.pptx') || file.type.includes('presentationml'))) {
    pptxFile.value = file
    pptxErr.value = ''
  } else {
    pptxErr.value = 'กรุณาเลือกไฟล์ .pptx เท่านั้น'
  }
}

function pptxSelectAll() {
  if (pptxPreview.value) pptxSelected.value = pptxPreview.value.slides.map((s) => s.index)
}
function pptxDeselectAll() {
  pptxSelected.value = []
}
function pptxSelectVisible() {
  if (pptxPreview.value) pptxSelected.value = pptxPreview.value.slides.filter((s) => !s.hidden).map((s) => s.index)
}

async function loadPptxPreview() {
  if (!pptxFile.value) return
  pptxLoading.value = true
  pptxErr.value = ''
  try {
    const data = await tasksApi.previewPPTXUpload(pptxFile.value)
    pptxPreview.value = data
    pptxSelected.value = data.slides.filter((s) => !s.hidden).map((s) => s.index)
    const triaged: Record<number, TriagedSlide> = {}
    for (const s of data.slides) {
      triaged[s.index] = {
        title: s.suggested_task_title?.trim() || `Slide ${s.index}`,
        assignee_id: null,
        estimated_minutes: 0,
        priority: 'MEDIUM',
      }
    }
    pptxTriaged.value = triaged
    pptxStep.value = 'select'
  } catch (e: any) {
    pptxErr.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
  } finally {
    pptxLoading.value = false
  }
}

async function submitPptx() {
  if (!pptxFile.value) return
  pptxImporting.value = true
  pptxErr.value = ''
  try {
    const pages = pptxSelected.value.map((idx) => {
      const t = pptxTriaged.value[idx]
      return {
        slide_index: idx,
        title: t?.title || `Slide ${idx}`,
        assignee_id: t?.assignee_id ?? null,
        estimated_minutes: t?.estimated_minutes ?? 0,
        priority: t?.priority || 'MEDIUM',
      }
    })
    const payload: Parameters<typeof tasksApi.importPPTXUpload>[1] = {
      project_id: props.projectId,
      parent_id: props.parentTaskId,
      priority: 'MEDIUM',
      story_points: 1,
      pages,
    }
    if (props.epicId) payload.epic_id = props.epicId
    pptxResult.value = await tasksApi.importPPTXUpload(pptxFile.value, payload)
    pptxStep.value = 'result'
  } catch (e: any) {
    pptxErr.value = e?.data?.message ?? e?.message ?? 'Import failed'
  } finally {
    pptxImporting.value = false
  }
}

// ── Sheets ────────────────────────────────────────────────
const showSheets = ref(false)
const sheetsStep = ref<'form' | 'select' | 'result'>('form')
const sheetsUrl = ref('')
const sheetsLoading = ref(false)
const sheetsImporting = ref(false)
const sheetsErr = ref('')
const sheetsPreview = ref<{
  sheet_title: string
  sheet_id: string
  rows: { row_index: number; title: string; due_date: string; status: string; raw_status: string; notes: string }[]
} | null>(null)
const sheetsRowsSelected = ref<number[]>([])
const sheetsTriaged = ref<Record<number, SheetsTriagedRow>>({})
const sheetsResult = ref<{ created_count: number; sheet_title: string; tasks: { id: string; title: string; code?: string }[] } | null>(null)

function openSheets() {
  sheetsStep.value = 'form'
  sheetsUrl.value = ''
  sheetsErr.value = ''
  sheetsPreview.value = null
  sheetsRowsSelected.value = []
  sheetsTriaged.value = {}
  sheetsResult.value = null
  showSheets.value = true
}

function closeSheets() {
  showSheets.value = false
}

function closeSheetsAfterSuccess() {
  const had = !!sheetsResult.value
  showSheets.value = false
  sheetsResult.value = null
  if (had) emit('imported')
}

function sheetsSelectAll() {
  if (sheetsPreview.value) sheetsRowsSelected.value = sheetsPreview.value.rows.map((r) => r.row_index)
}
function sheetsDeselectAll() {
  sheetsRowsSelected.value = []
}

async function loadSheetsPreview() {
  if (!sheetsUrl.value.trim()) return
  sheetsLoading.value = true
  sheetsErr.value = ''
  try {
    const data = await tasksApi.previewGoogleSheets({ sheet_url: sheetsUrl.value.trim() })
    sheetsPreview.value = data
    const triaged: Record<number, SheetsTriagedRow> = {}
    const selected: number[] = []
    for (const r of data.rows) {
      selected.push(r.row_index)
      triaged[r.row_index] = {
        title: r.title,
        priority: 'MEDIUM',
        estimated_minutes: 0,
        due_date: r.due_date || '',
        status: r.status,
        notes: r.notes,
      }
    }
    sheetsTriaged.value = triaged
    sheetsRowsSelected.value = selected
    sheetsStep.value = 'select'
  } catch (e: any) {
    sheetsErr.value = e?.data?.message ?? e?.message ?? 'โหลดรายการไม่สำเร็จ'
  } finally {
    sheetsLoading.value = false
  }
}

async function submitSheets() {
  if (!sheetsPreview.value) return
  sheetsImporting.value = true
  sheetsErr.value = ''
  try {
    const rows = sheetsRowsSelected.value.map((rowIndex) => {
      const t = sheetsTriaged.value[rowIndex]
      const rawEst = Number(t?.estimated_minutes)
      const estimatedMinutes = Number.isFinite(rawEst) && rawEst >= 0 ? Math.floor(rawEst) : 0
      return {
        row_index: rowIndex,
        title: t?.title?.trim() || '',
        priority: t?.priority || 'MEDIUM',
        estimated_minutes: estimatedMinutes,
        due_date: t?.due_date?.trim() || '',
        status: t?.status || 'PENDING',
        notes: t?.notes?.trim() || '',
      }
    })
    const payload: Record<string, unknown> = {
      sheet_url: sheetsUrl.value.trim(),
      sheet_title: sheetsPreview.value.sheet_title,
      project_id: props.projectId,
      parent_id: props.parentTaskId,
      rows,
    }
    if (props.epicId) payload.epic_id = props.epicId
    sheetsResult.value = await tasksApi.importGoogleSheets(payload as Parameters<typeof tasksApi.importGoogleSheets>[0])
    sheetsStep.value = 'result'
  } catch (e: any) {
    sheetsErr.value = e?.data?.message ?? e?.message ?? 'Import failed'
  } finally {
    sheetsImporting.value = false
  }
}

defineExpose({ openSlides, openPptx, openSheets })
</script>

<style scoped>
.label {
  @apply block text-xs text-gray-400 mb-1.5 font-medium;
}
.input-field {
  @apply bg-gray-700 border border-gray-600 rounded-xl px-4 py-2.5 text-sm text-gray-200 placeholder-gray-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/50 transition-colors;
}
.btn-primary {
  @apply bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 text-white font-semibold rounded-xl transition-colors;
}
.btn-import-sm {
  @apply px-3 py-1.5 text-xs bg-purple-900/50 hover:bg-purple-800/60 border border-purple-700/50 text-purple-300 font-medium rounded-lg transition-colors flex items-center gap-1.5;
}
.btn-ghost-sm {
  @apply px-3 py-1.5 text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 font-medium rounded-lg transition-colors;
}
</style>
