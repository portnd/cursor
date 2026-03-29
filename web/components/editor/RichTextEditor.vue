<template>
  <div class="rich-editor" :class="{ 'is-readonly': readonly }">

    <!-- Fullscreen image lightbox (shared component) -->
    <ImageFullscreenOverlay
      :show="!!fullscreenImageSrc"
      :image-src="fullscreenImageSrc"
      :show-actions="!readonly"
      @close="closeFullscreen"
      @annotate="openAnnotatorFromFullscreen"
      @delete="deleteImageFromFullscreen"
    />

    <!-- Image Annotator Modal -->
    <ImageAnnotator
      v-if="annotatorSrc"
      :image-src="annotatorSrc"
      @close="annotatorSrc = null"
      @save="onAnnotationSave"
    />
    <!-- Toolbar -->
    <div v-if="!readonly" class="editor-toolbar">
      <!-- Text Style -->
      <div class="toolbar-group">
        <button
          v-for="heading in [1, 2, 3]"
          :key="heading"
          @click="setHeading(heading)"
          :class="{ 'is-active': editor?.isActive('heading', { level: heading }) }"
          class="toolbar-btn"
          :title="`Heading ${heading}`"
          type="button"
        >
          H{{ heading }}
        </button>
      </div>

      <div class="toolbar-divider" />

      <!-- Inline Formatting -->
      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().toggleBold().run()"
          :class="{ 'is-active': editor?.isActive('bold') }"
          class="toolbar-btn font-bold"
          title="Bold (⌘B)"
          type="button"
        >B</button>
        <button
          @click="editor?.chain().focus().toggleItalic().run()"
          :class="{ 'is-active': editor?.isActive('italic') }"
          class="toolbar-btn italic"
          title="Italic (⌘I)"
          type="button"
        >I</button>
        <button
          @click="editor?.chain().focus().toggleUnderline().run()"
          :class="{ 'is-active': editor?.isActive('underline') }"
          class="toolbar-btn underline"
          title="Underline (⌘U)"
          type="button"
        >U</button>
        <button
          @click="editor?.chain().focus().toggleStrike().run()"
          :class="{ 'is-active': editor?.isActive('strike') }"
          class="toolbar-btn line-through"
          title="Strikethrough"
          type="button"
        >S</button>
        <button
          @click="editor?.chain().focus().toggleCode().run()"
          :class="{ 'is-active': editor?.isActive('code') }"
          class="toolbar-btn font-mono text-xs"
          title="Inline Code"
          type="button"
        >&lt;/&gt;</button>
      </div>

      <div class="toolbar-divider" />

      <!-- Lists -->
      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().toggleBulletList().run()"
          :class="{ 'is-active': editor?.isActive('bulletList') }"
          class="toolbar-btn"
          title="Bullet List"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().toggleOrderedList().run()"
          :class="{ 'is-active': editor?.isActive('orderedList') }"
          class="toolbar-btn"
          title="Numbered List"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 6h13M7 12h13M7 18h13M3 6h.01M3 12h.01M3 18h.01"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().toggleBlockquote().run()"
          :class="{ 'is-active': editor?.isActive('blockquote') }"
          class="toolbar-btn"
          title="Blockquote"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-3 3v-3z"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          :class="{ 'is-active': editor?.isActive('codeBlock') }"
          class="toolbar-btn"
          title="Code Block"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider" />

      <!-- Alignment -->
      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().setTextAlign('left').run()"
          :class="{ 'is-active': editor?.isActive({ textAlign: 'left' }) }"
          class="toolbar-btn"
          title="Align Left"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h16"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().setTextAlign('center').run()"
          :class="{ 'is-active': editor?.isActive({ textAlign: 'center' }) }"
          class="toolbar-btn"
          title="Align Center"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M4 18h16"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().setTextAlign('right').run()"
          :class="{ 'is-active': editor?.isActive({ textAlign: 'right' }) }"
          class="toolbar-btn"
          title="Align Right"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M4 18h16"/>
          </svg>
        </button>
      </div>

      <div class="toolbar-divider" />

      <!-- Image & Link -->
      <div class="toolbar-group">
        <button
          @click="insertImageFromUrl"
          class="toolbar-btn"
          title="Insert Image (or paste with ⌘V)"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </button>
        <button
          @click="setLink"
          :class="{ 'is-active': editor?.isActive('link') }"
          class="toolbar-btn"
          title="Add Link"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().setHorizontalRule().run()"
          class="toolbar-btn"
          title="Horizontal Rule"
          type="button"
        >—</button>
      </div>

      <div class="toolbar-divider" />

      <!-- History -->
      <div class="toolbar-group">
        <button
          @click="editor?.chain().focus().undo().run()"
          :disabled="!editor?.can().undo()"
          class="toolbar-btn"
          title="Undo (⌘Z)"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"/>
          </svg>
        </button>
        <button
          @click="editor?.chain().focus().redo().run()"
          :disabled="!editor?.can().redo()"
          class="toolbar-btn"
          title="Redo (⌘⇧Z)"
          type="button"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10H11a8 8 0 00-8 8v2m18-10l-6 6m6-6l-6-6"/>
          </svg>
        </button>
      </div>

      <!-- Image paste hint -->
      <div class="ml-auto flex items-center gap-1 text-xs text-gray-500 pr-1">
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14"/>
        </svg>
        ⌘V to paste image
      </div>
    </div>

    <!-- Hidden file input -->
    <input
      v-if="!readonly"
      ref="fileInputRef"
      type="file"
      accept="image/*"
      class="hidden"
      @change="handleFileUpload"
    />

    <!-- Editor content -->
    <editor-content :editor="editor" class="editor-content" />

    <!-- Image paste indicator -->
    <div
      v-if="isPastingImage"
      class="absolute inset-0 flex items-center justify-center bg-gray-900/80 rounded-lg z-10 pointer-events-none"
    >
      <div class="flex items-center gap-2 text-blue-400 text-sm font-medium">
        <svg class="w-5 h-5 animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
        </svg>
        Inserting image...
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import { createResizableImageExtension } from '~/components/editor/resizable-image-extension'
import Placeholder from '@tiptap/extension-placeholder'
import Underline from '@tiptap/extension-underline'
import Link from '@tiptap/extension-link'
import TextAlign from '@tiptap/extension-text-align'
import { TextStyle } from '@tiptap/extension-text-style'
import ImageAnnotator from '~/components/editor/ImageAnnotator.vue'
import ImageFullscreenOverlay from '~/components/editor/ImageFullscreenOverlay.vue'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

// Slide/PPTX imports stored <p><img ... /></p>. TipTap Image is a block node and cannot live inside a paragraph, so the image was dropped on parse — lift to top-level blocks.
function liftImportedImagesOutOfParagraphs(html: string): string {
  if (!html || typeof html !== 'string') return html
  let out = html
  let prev = ''
  const re = /<p>\s*(<img\b[^>]*>)\s*<\/p>/gi
  while (out !== prev) {
    prev = out
    out = out.replace(re, '$1')
  }
  return out
}

const fileInputRef = ref<HTMLInputElement | null>(null)
const isPastingImage = ref(false)

// Fullscreen image lightbox
const fullscreenImageSrc = ref<string | null>(null)
const fullscreenImgEl = ref<HTMLImageElement | null>(null)

const openFullscreen = (imgEl: HTMLImageElement) => {
  fullscreenImgEl.value = imgEl
  fullscreenImageSrc.value = imgEl.src
}

const closeFullscreen = () => {
  fullscreenImageSrc.value = null
  fullscreenImgEl.value = null
}

function onFullscreenKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && fullscreenImageSrc.value) {
    closeFullscreen()
  }
}

watch(fullscreenImageSrc, (src) => {
  if (typeof document === 'undefined') return
  if (src) {
    document.addEventListener('keydown', onFullscreenKeydown)
  } else {
    document.removeEventListener('keydown', onFullscreenKeydown)
  }
})

onBeforeUnmount(() => {
  if (typeof document !== 'undefined') {
    document.removeEventListener('keydown', onFullscreenKeydown)
  }
})

const openAnnotatorFromFullscreen = () => {
  const el = fullscreenImgEl.value
  if (el) {
    annotatingImgEl = el
    annotatorSrc.value = el.src
  }
  closeFullscreen()
}

const deleteImageFromFullscreen = () => {
  const src = fullscreenImageSrc.value
  const ed = editor.value
  if (!src || !ed) {
    closeFullscreen()
    return
  }
  const { state, view } = ed
  let from: number | null = null
  state.doc.descendants((node, pos) => {
    if (node.type.name === 'image' && node.attrs.src === src) {
      from = pos
      return false
    }
  })
  if (from != null) {
    const node = state.doc.nodeAt(from)
    if (node) {
      const tr = state.tr.delete(from, from + node.nodeSize)
      view.dispatch(tr)
      const html = ed.getHTML()
      emit('update:modelValue', html === '<p></p>' ? '' : html)
    }
  }
  closeFullscreen()
}

// Image annotation
const annotatorSrc = ref<string | null>(null)
let annotatingImgEl: HTMLImageElement | null = null

const openAnnotator = (imgEl: HTMLImageElement) => {
  annotatingImgEl = imgEl
  annotatorSrc.value = imgEl.src
}

const onAnnotationSave = (annotatedSrc: string) => {
  if (annotatingImgEl && editor.value) {
    // Replace the img src in the editor content
    const html = editor.value.getHTML()
    const escapedSrc = annotatingImgEl.src.replace(/[.*+?^${}()|[\]\\]/g, '\\$&').substring(0, 100)
    // Use DOM manipulation on the editor state instead
    const { state, view } = editor.value
    state.doc.descendants((node, pos) => {
      if (node.type.name === 'image' && node.attrs.src === annotatingImgEl!.src) {
        const tr = state.tr.setNodeMarkup(pos, undefined, {
          ...node.attrs,
          src: annotatedSrc,
        })
        view.dispatch(tr)
        return false
      }
    })
  }
  annotatorSrc.value = null
  annotatingImgEl = null
}

const editor = useEditor({
  content: liftImportedImagesOutOfParagraphs(props.modelValue || ''),
  editable: !props.readonly,
  extensions: [
    StarterKit.configure({
      history: { depth: 100 },
      // Disable extensions we add manually to avoid duplicates
      link: false,
      underline: false,
    }),
    createResizableImageExtension().configure({
      inline: false,
      allowBase64: true,
      HTMLAttributes: {
        class: 'editor-image',
      },
    }),
    Placeholder.configure({
      placeholder: props.placeholder || 'Describe what needs to be done... (paste images with ⌘V)',
    }),
    Underline,
    Link.configure({
      openOnClick: false,
      HTMLAttributes: { class: 'editor-link' },
    }),
    TextAlign.configure({
      types: ['heading', 'paragraph'],
    }),
    TextStyle,
  ],
  onUpdate({ editor }) {
    const html = editor.getHTML()
    emit('update:modelValue', html === '<p></p>' ? '' : html)
  },
})

// Sync external modelValue changes into editor
watch(() => props.modelValue, (newVal) => {
  if (!editor.value) return
  const current = editor.value.getHTML()
  const normalized = liftImportedImagesOutOfParagraphs(newVal || '')
  if (current !== normalized) {
    editor.value.commands.setContent(normalized, false)
  }
})

// Sync readonly changes
watch(() => props.readonly, (isReadonly) => {
  editor.value?.setEditable(!isReadonly)
})

// --- Image Helpers ---

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

async function insertImageFile(file: File) {
  if (!editor.value) return
  isPastingImage.value = true
  try {
    const base64 = await fileToBase64(file)
    editor.value.chain().focus().setImage({ src: base64 }).run()
  } finally {
    isPastingImage.value = false
  }
}

function insertImageFromUrl() {
  if (fileInputRef.value) {
    fileInputRef.value.click()
  }
}

async function handleFileUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) {
    await insertImageFile(file)
    input.value = ''
  }
}

function setHeading(level: 1 | 2 | 3) {
  if (!editor.value) return
  if (editor.value.isActive('heading', { level })) {
    editor.value.chain().focus().setParagraph().run()
  } else {
    editor.value.chain().focus().toggleHeading({ level }).run()
  }
}

function setLink() {
  if (!editor.value) return
  const prev = editor.value.getAttributes('link').href
  const url = window.prompt('Enter URL:', prev || 'https://')
  if (url === null) return
  if (url === '') {
    editor.value.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }
  editor.value.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}

// Handle paste event for images (Cmd+V) + image click to annotate
onMounted(() => {
  nextTick(() => {
    // Find this component's ProseMirror instance (use the component root ref)
    const editorEl = editor.value?.view?.dom
    if (!editorEl) return

    editorEl.addEventListener('paste', async (e: Event) => {
      const pasteEvent = e as ClipboardEvent
      const items = pasteEvent.clipboardData?.items
      if (!items) return
      for (const item of Array.from(items)) {
        if (item.type.startsWith('image/')) {
          pasteEvent.preventDefault()
          const file = item.getAsFile()
          if (file) await insertImageFile(file)
          return
        }
      }
    })

    // Image click → fullscreen (readonly: view only; edit mode: fullscreen + option to annotate)
    editorEl.addEventListener('click', (e: Event) => {
      const target = e.target as HTMLElement
      if (target.tagName === 'IMG') {
        openFullscreen(target as HTMLImageElement)
      }
    })
  })
})

onBeforeUnmount(() => {
  editor.value?.destroy()
})
</script>

<style>
/* Editor wrapper */
.rich-editor {
  position: relative;
  display: flex;
  flex-direction: column;
  border: 1px solid rgb(55 65 81); /* gray-700 */
  border-radius: 0.5rem;
  overflow: hidden;
  background: rgb(17 24 39); /* gray-900 */
}

.rich-editor:focus-within:not(.is-readonly) {
  border-color: rgb(59 130 246); /* blue-500 */
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.rich-editor.is-readonly {
  border-color: transparent;
  background: transparent;
}

/* Toolbar */
.editor-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 2px;
  padding: 6px 8px;
  background: rgb(31 41 55); /* gray-800 */
  border-bottom: 1px solid rgb(55 65 81); /* gray-700 */
  min-height: 40px;
}

.toolbar-group {
  display: flex;
  align-items: center;
  gap: 1px;
}

.toolbar-divider {
  width: 1px;
  height: 20px;
  background: rgb(55 65 81); /* gray-700 */
  margin: 0 4px;
}

.toolbar-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
  padding: 0 5px;
  border-radius: 4px;
  font-size: 13px;
  color: rgb(156 163 175); /* gray-400 */
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.1s;
  user-select: none;
}

.toolbar-btn:hover:not(:disabled) {
  background: rgb(55 65 81); /* gray-700 */
  color: white;
}

.toolbar-btn.is-active {
  background: rgb(37 99 235); /* blue-600 */
  color: white;
}

.toolbar-btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

/* Editor content area */
.editor-content {
  flex: 1;
  padding: 14px 16px;
  min-height: 160px;
  max-height: 600px;
  overflow-y: auto;
}

.is-readonly .editor-content {
  padding: 0;
  min-height: unset;
  max-height: unset;
  overflow-y: visible;
}

/* ProseMirror core */
.ProseMirror {
  outline: none;
  color: rgb(209 213 219); /* gray-300 */
  font-size: 14px;
  line-height: 1.7;
  word-break: break-word;
}

.ProseMirror p.is-editor-empty:first-child::before {
  content: attr(data-placeholder);
  float: left;
  color: rgb(107 114 128); /* gray-500 */
  pointer-events: none;
  height: 0;
}

/* Typography */
.ProseMirror h1 { font-size: 1.6em; font-weight: 700; color: white; margin: 1em 0 0.5em; }
.ProseMirror h2 { font-size: 1.3em; font-weight: 700; color: white; margin: 0.9em 0 0.45em; }
.ProseMirror h3 { font-size: 1.1em; font-weight: 600; color: rgb(229 231 235); margin: 0.8em 0 0.4em; }
.ProseMirror p { margin: 0.4em 0; }
.ProseMirror strong { font-weight: 700; color: white; }
.ProseMirror em { font-style: italic; color: rgb(209 213 219); }
.ProseMirror u { text-decoration: underline; }
.ProseMirror s { text-decoration: line-through; color: rgb(107 114 128); }

/* Code */
.ProseMirror code {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 0.875em;
  background: rgb(31 41 55);
  color: rgb(110 231 183); /* emerald-300 */
  padding: 1px 5px;
  border-radius: 3px;
  border: 1px solid rgb(55 65 81);
}

.ProseMirror pre {
  background: rgb(15 23 42); /* slate-900 */
  border: 1px solid rgb(55 65 81);
  border-radius: 6px;
  padding: 12px 14px;
  margin: 0.75em 0;
  overflow-x: auto;
}

.ProseMirror pre code {
  background: none;
  border: none;
  padding: 0;
  color: rgb(110 231 183);
  font-size: 13px;
}

/* Lists */
.ProseMirror ul, .ProseMirror ol { padding-left: 1.5em; margin: 0.5em 0; }
.ProseMirror ul { list-style-type: disc; }
.ProseMirror ol { list-style-type: decimal; }
.ProseMirror li { margin: 0.2em 0; }
.ProseMirror li p { margin: 0; }

/* Blockquote */
.ProseMirror blockquote {
  border-left: 3px solid rgb(59 130 246); /* blue-500 */
  padding: 4px 0 4px 16px;
  margin: 0.75em 0;
  color: rgb(156 163 175); /* gray-400 */
  font-style: italic;
}

/* Horizontal rule */
.ProseMirror hr {
  border: none;
  border-top: 1px solid rgb(55 65 81);
  margin: 1em 0;
}

/* Links */
.ProseMirror a, .editor-link {
  color: rgb(96 165 250); /* blue-400 */
  text-decoration: underline;
  cursor: pointer;
}

/* Resizable image wrapper: 8 handles (4 corners = proportional, 4 edges = single axis) */
.ProseMirror .resizable-image-wrapper {
  position: relative;
  display: inline-block;
  max-width: 100%;
  margin: 8px 0;
}

.ProseMirror .resizable-image-wrapper .resize-handle {
  position: absolute;
  width: 10px;
  height: 10px;
  background: rgb(59 130 246);
  border: 2px solid white;
  border-radius: 2px;
  opacity: 0;
  transition: opacity 0.15s;
  z-index: 2;
}

/* Corners */
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="nw"] { top: -2px; left: -2px; }
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="ne"] { top: -2px; right: -2px; }
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="sw"] { bottom: -2px; left: -2px; }
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="se"] { bottom: -2px; right: -2px; }

/* Edges: longer hit area for easier grab */
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="n"] {
  top: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 10px;
}
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="s"] {
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 24px;
  height: 10px;
}
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="e"] {
  right: -2px;
  top: 50%;
  transform: translateY(-50%);
  width: 10px;
  height: 24px;
}
.ProseMirror .resizable-image-wrapper .resize-handle[data-resize-handle="w"] {
  left: -2px;
  top: 50%;
  transform: translateY(-50%);
  width: 10px;
  height: 24px;
}

.ProseMirror .resizable-image-wrapper:hover .resize-handle,
.ProseMirror .resizable-image-wrapper.resizing .resize-handle {
  opacity: 1;
}

.ProseMirror .resizable-image-wrapper.resizing {
  user-select: none;
}

.ProseMirror .resizable-image-wrapper .editor-image {
  display: block;
  max-width: 100%;
  height: auto;
}

/* Images (fallback when not in wrapper) */
.ProseMirror img.editor-image,
.ProseMirror img {
  max-width: 100%;
  height: auto;
  border-radius: 6px;
  margin: 8px 0;
  border: 2px solid transparent;
  display: block;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.ProseMirror img.editor-image:hover,
.ProseMirror img:hover {
  border-color: rgb(59 130 246); /* blue-500 */
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.ProseMirror img.ProseMirror-selectednode {
  border-color: rgb(59 130 246); /* blue-500 */
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

/* Annotate hint below image in edit mode */
.rich-editor:not(.is-readonly) .ProseMirror img::after {
  content: 'Click to annotate';
}

/* Text alignment */
.ProseMirror [style*="text-align: center"] { text-align: center; }
.ProseMirror [style*="text-align: right"] { text-align: right; }
.ProseMirror [style*="text-align: left"] { text-align: left; }

/* Selection */
.ProseMirror ::selection {
  background: rgba(59, 130, 246, 0.3);
}

/* Scrollbar */
.editor-content::-webkit-scrollbar { width: 6px; }
.editor-content::-webkit-scrollbar-track { background: transparent; }
.editor-content::-webkit-scrollbar-thumb { background: rgb(55 65 81); border-radius: 3px; }
.editor-content::-webkit-scrollbar-thumb:hover { background: rgb(75 85 99); }
</style>
