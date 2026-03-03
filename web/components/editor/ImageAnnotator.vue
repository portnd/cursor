<template>
  <Teleport to="body">
    <div
      class="annotator-overlay"
      @mouseup.self="handleBackdropClick"
    >
      <div class="annotator-modal">
        <!-- Header -->
        <div class="annotator-header">
          <div class="flex items-center gap-3">
            <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
            </svg>
            <span class="text-white font-semibold">Annotate Image</span>
            <span class="text-xs text-gray-500">Draw boxes, add comments</span>
          </div>
          <button @click="$emit('close')" class="annotator-close-btn" type="button">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Toolbar -->
        <div class="annotator-toolbar">
          <!-- Tools -->
          <div class="flex items-center gap-1">
            <button
              v-for="tool in tools"
              :key="tool.id"
              @click="activeTool = tool.id"
              :class="{ 'is-active': activeTool === tool.id }"
              class="tool-btn"
              :title="tool.label"
              type="button"
            >
              <component :is="tool.icon" class="w-4 h-4" />
            </button>
          </div>

          <div class="toolbar-sep" />

          <!-- Color -->
          <div class="flex items-center gap-1.5">
            <span class="text-xs text-gray-500">Color:</span>
            <button
              v-for="color in colors"
              :key="color.value"
              @click="activeColor = color.value"
              :style="{ background: color.value }"
              :class="{ 'ring-2 ring-white ring-offset-1 ring-offset-gray-800': activeColor === color.value }"
              class="w-5 h-5 rounded-full border border-white/20 transition-all"
              :title="color.label"
              type="button"
            />
          </div>

          <div class="toolbar-sep" />

          <!-- Stroke width -->
          <div class="flex items-center gap-1.5">
            <span class="text-xs text-gray-500">Size:</span>
            <button
              v-for="sw in strokeWidths"
              :key="sw"
              @click="strokeWidth = sw"
              :class="{ 'bg-blue-600': strokeWidth === sw, 'bg-gray-700 hover:bg-gray-600': strokeWidth !== sw }"
              class="px-2 py-0.5 rounded text-xs text-white transition-colors"
              type="button"
            >{{ sw }}px</button>
          </div>

          <div class="toolbar-sep" />

          <!-- Actions -->
          <div class="flex items-center gap-1">
            <button
              @click="undo"
              :disabled="annotations.length === 0"
              class="tool-btn"
              title="Undo last annotation"
              type="button"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"/>
              </svg>
            </button>
            <button
              @click="clearAll"
              :disabled="annotations.length === 0"
              class="tool-btn text-red-400 hover:text-red-300"
              title="Clear all annotations"
              type="button"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>

          <div class="ml-auto flex items-center gap-2">
            <button @click="$emit('close')" class="px-3 py-1.5 text-sm bg-gray-700 hover:bg-gray-600 text-gray-300 rounded-lg transition-colors" type="button">
              Cancel
            </button>
            <button
              @click="saveAnnotated"
              class="px-4 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors flex items-center gap-1.5"
              type="button"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              Save
            </button>
          </div>
        </div>

        <!-- Main: canvas + comment panel -->
        <div class="annotator-body">
          <!-- Canvas area -->
          <div class="canvas-wrapper" ref="canvasWrapperRef">
            <div
              class="canvas-container"
              :style="{ width: displayWidth + 'px', height: displayHeight + 'px', cursor: canvasCursor }"
              @mousedown="onMouseDown"
              @mousemove="onMouseMove"
              @mouseup="onMouseUp"
              @mouseleave="onMouseLeave"
            >
              <!-- Base image -->
              <img
                ref="imgRef"
                :src="imageSrc"
                class="canvas-image"
                :style="{ width: displayWidth + 'px', height: displayHeight + 'px' }"
                draggable="false"
                @load="onImageLoad"
              />

              <!-- SVG annotation overlay -->
              <svg
                class="annotation-svg"
                :width="displayWidth"
                :height="displayHeight"
                :viewBox="`0 0 ${displayWidth} ${displayHeight}`"
              >
                <!-- All arrow markers in one defs (fixes multiple annotations) -->
                <defs>
                  <marker
                    v-for="ann in arrowAnnotations"
                    :key="'m-' + ann.id"
                    :id="'arrow-' + ann.id"
                    markerWidth="10"
                    markerHeight="7"
                    refX="9"
                    refY="3.5"
                    orient="auto"
                  >
                    <polygon points="0 0, 10 3.5, 0 7" :fill="ann.color"/>
                  </marker>
                </defs>

                <!-- Saved annotations (multiple supported) -->
                <g
                  v-for="ann in annotations"
                  :key="ann.id"
                  class="annotation-group"
                  @click="selectAnnotation(ann.id)"
                >
                  <!-- Rectangle -->
                  <rect
                    v-if="ann.type === 'rect'"
                    :x="ann.x" :y="ann.y"
                    :width="ann.w" :height="ann.h"
                    :stroke="ann.color"
                    :stroke-width="ann.strokeWidth"
                    fill="none"
                    stroke-dasharray="none"
                    class="annotation-shape"
                    :class="{ 'selected-shape': selectedId === ann.id }"
                  />
                  <!-- Circle -->
                  <ellipse
                    v-if="ann.type === 'circle'"
                    :cx="ann.x + ann.w / 2"
                    :cy="ann.y + ann.h / 2"
                    :rx="Math.abs(ann.w / 2)"
                    :ry="Math.abs(ann.h / 2)"
                    :stroke="ann.color"
                    :stroke-width="ann.strokeWidth"
                    fill="none"
                    class="annotation-shape"
                    :class="{ 'selected-shape': selectedId === ann.id }"
                  />
                  <!-- Arrow (marker referenced from defs above) -->
                  <line
                    v-if="ann.type === 'arrow'"
                    :x1="ann.x" :y1="ann.y"
                    :x2="ann.x + ann.w" :y2="ann.y + ann.h"
                    :stroke="ann.color"
                    :stroke-width="ann.strokeWidth"
                    :marker-end="'url(#arrow-' + ann.id + ')'"
                    class="annotation-shape"
                    :class="{ 'selected-shape': selectedId === ann.id }"
                  />
                  <!-- Highlight -->
                  <rect
                    v-if="ann.type === 'highlight'"
                    :x="ann.x" :y="ann.y"
                    :width="ann.w" :height="ann.h"
                    :fill="ann.color"
                    fill-opacity="0.3"
                    stroke="none"
                  />
                  <!-- Freehand -->
                  <polyline
                    v-if="ann.type === 'pen' && ann.points"
                    :points="ann.points.map((p: Point) => `${p.x},${p.y}`).join(' ')"
                    :stroke="ann.color"
                    :stroke-width="ann.strokeWidth"
                    fill="none"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  />

                  <!-- Comment badge -->
                  <g v-if="ann.comment || ann.label" class="comment-badge" @click.stop="selectAnnotation(ann.id)">
                    <circle
                      :cx="badgeX(ann)"
                      :cy="badgeY(ann)"
                      r="11"
                      :fill="ann.color"
                    />
                    <text
                      :x="badgeX(ann)"
                      :y="badgeY(ann)"
                      text-anchor="middle"
                      dominant-baseline="central"
                      fill="white"
                      font-size="11"
                      font-weight="bold"
                      font-family="system-ui"
                    >{{ ann.label }}</text>
                  </g>
                </g>

                <!-- Active drawing shape preview -->
                <g v-if="isDrawing && drawPreview">
                  <rect
                    v-if="activeTool === 'rect'"
                    :x="drawPreview.x" :y="drawPreview.y"
                    :width="drawPreview.w" :height="drawPreview.h"
                    :stroke="activeColor"
                    :stroke-width="strokeWidth"
                    fill="none"
                    stroke-dasharray="5,3"
                  />
                  <ellipse
                    v-if="activeTool === 'circle'"
                    :cx="drawPreview.x + drawPreview.w / 2"
                    :cy="drawPreview.y + drawPreview.h / 2"
                    :rx="Math.abs(drawPreview.w / 2)"
                    :ry="Math.abs(drawPreview.h / 2)"
                    :stroke="activeColor"
                    :stroke-width="strokeWidth"
                    fill="none"
                    stroke-dasharray="5,3"
                  />
                  <line
                    v-if="activeTool === 'arrow'"
                    :x1="drawPreview.x" :y1="drawPreview.y"
                    :x2="drawPreview.x + drawPreview.w" :y2="drawPreview.y + drawPreview.h"
                    :stroke="activeColor"
                    :stroke-width="strokeWidth"
                    marker-end="url(#arrow-preview)"
                  />
                  <rect
                    v-if="activeTool === 'highlight'"
                    :x="drawPreview.x" :y="drawPreview.y"
                    :width="drawPreview.w" :height="drawPreview.h"
                    :fill="activeColor"
                    fill-opacity="0.3"
                    stroke="none"
                  />
                  <defs>
                    <marker id="arrow-preview" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
                      <polygon points="0 0, 10 3.5, 0 7" :fill="activeColor"/>
                    </marker>
                  </defs>
                </g>
              </svg>
            </div>
          </div>

          <!-- Comment Panel -->
          <div class="comment-panel">
            <div class="comment-panel-header">
              <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
              </svg>
              <span class="text-sm font-medium text-gray-300">Annotations</span>
              <span class="text-xs text-gray-500 ml-1">({{ annotations.length }})</span>
            </div>

            <div class="comment-list">
              <div
                v-if="annotations.length === 0"
                class="text-center py-8 text-gray-600 text-xs"
              >
                <svg class="w-8 h-8 mx-auto mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536M9 11l6.586-6.586a2 2 0 112.828 2.828L11.828 13.828a2 2 0 01-.828.485l-3 .75a.75.75 0 01-.918-.919l.75-3a2 2 0 01.485-.828z"/>
                </svg>
                Draw on the image to add annotations.<br>
                <span class="text-gray-500">You can add as many as you need.</span>
              </div>

              <div
                v-for="ann in annotations"
                :key="ann.id"
                @click="selectAnnotation(ann.id)"
                :class="{ 'selected-comment': selectedId === ann.id }"
                class="comment-item"
              >
                <div class="flex items-start gap-2">
                  <div
                    class="comment-badge-small shrink-0"
                    :style="{ background: ann.color }"
                  >{{ ann.label }}</div>
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center justify-between gap-1 mb-1">
                      <span class="text-xs text-gray-400 capitalize">{{ toolLabel(ann.type) }}</span>
                      <button
                        @click.stop="removeAnnotation(ann.id)"
                        class="text-gray-600 hover:text-red-400 transition-colors"
                        type="button"
                      >
                        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                        </svg>
                      </button>
                    </div>
                    <textarea
                      v-model="ann.comment"
                      @click.stop
                      placeholder="Add comment..."
                      rows="2"
                      class="comment-textarea"
                    />
                  </div>
                </div>
              </div>
            </div>

            <!-- Tip -->
            <div class="comment-tip">
              <svg class="w-3.5 h-3.5 shrink-0 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <span>Add multiple annotations, then click a shape to edit its comment</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
interface Point { x: number; y: number }

interface Annotation {
  id: string
  type: 'rect' | 'circle' | 'arrow' | 'highlight' | 'pen'
  x: number
  y: number
  w: number
  h: number
  color: string
  strokeWidth: number
  label: string
  comment: string
  points?: Point[]
}

interface DrawPreview {
  x: number
  y: number
  w: number
  h: number
}

const props = defineProps<{
  imageSrc: string
}>()

const emit = defineEmits<{
  close: []
  save: [annotatedSrc: string]
}>()

// --- State ---
const canvasWrapperRef = ref<HTMLDivElement | null>(null)
const imgRef = ref<HTMLImageElement | null>(null)

const naturalWidth = ref(0)
const naturalHeight = ref(0)
const displayWidth = ref(800)
const displayHeight = ref(500)

const activeTool = ref<'rect' | 'circle' | 'arrow' | 'highlight' | 'pen'>('rect')
const activeColor = ref('#ef4444')
const strokeWidth = ref(2)

const annotations = ref<Annotation[]>([])
const selectedId = ref<string | null>(null)

const isDrawing = ref(false)
const startPoint = ref<Point>({ x: 0, y: 0 })
const drawPreview = ref<DrawPreview | null>(null)
const penPoints = ref<Point[]>([])

// --- Tools config ---
const RectIcon = defineComponent({
  render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' },
    [h('rect', { x: '3', y: '3', width: '18', height: '18', rx: '2', 'stroke-width': '2', 'stroke-linecap': 'round' })])
})
const CircleIcon = defineComponent({
  render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' },
    [h('circle', { cx: '12', cy: '12', r: '9', 'stroke-width': '2' })])
})
const ArrowIcon = defineComponent({
  render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' },
    [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M17 8l4 4m0 0l-4 4m4-4H3' })])
})
const HighlightIcon = defineComponent({
  render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' },
    [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01' })])
})
const PenIcon = defineComponent({
  render: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' },
    [h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z' })])
})

const tools = [
  { id: 'rect' as const, label: 'Rectangle (box)', icon: RectIcon },
  { id: 'circle' as const, label: 'Circle / Ellipse', icon: CircleIcon },
  { id: 'arrow' as const, label: 'Arrow', icon: ArrowIcon },
  { id: 'highlight' as const, label: 'Highlight area', icon: HighlightIcon },
  { id: 'pen' as const, label: 'Freehand draw', icon: PenIcon },
]

const colors = [
  { value: '#ef4444', label: 'Red' },
  { value: '#f97316', label: 'Orange' },
  { value: '#eab308', label: 'Yellow' },
  { value: '#22c55e', label: 'Green' },
  { value: '#3b82f6', label: 'Blue' },
  { value: '#a855f7', label: 'Purple' },
  { value: '#ffffff', label: 'White' },
]

const strokeWidths = [1, 2, 3, 5]

// --- Computed ---
const arrowAnnotations = computed(() =>
  annotations.value.filter((a) => a.type === 'arrow')
)

const canvasCursor = computed(() => {
  if (activeTool.value === 'pen') return 'crosshair'
  if (activeTool.value === 'arrow') return 'crosshair'
  return 'crosshair'
})

// --- Image load ---
const onImageLoad = () => {
  const img = imgRef.value
  if (!img) return
  naturalWidth.value = img.naturalWidth
  naturalHeight.value = img.naturalHeight
  fitToWrapper()
}

const fitToWrapper = () => {
  const wrapper = canvasWrapperRef.value
  if (!wrapper || !naturalWidth.value) return

  const maxW = wrapper.clientWidth - 8
  const maxH = wrapper.clientHeight - 8
  const ratio = naturalWidth.value / naturalHeight.value

  let w = naturalWidth.value
  let h = naturalHeight.value

  if (w > maxW) { w = maxW; h = w / ratio }
  if (h > maxH) { h = maxH; w = h * ratio }

  displayWidth.value = Math.round(w)
  displayHeight.value = Math.round(h)
}

// --- Mouse events ---
const getRelativePos = (e: MouseEvent): Point => {
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  return {
    x: Math.round(e.clientX - rect.left),
    y: Math.round(e.clientY - rect.top),
  }
}

const onMouseDown = (e: MouseEvent) => {
  if (e.button !== 0) return
  const pos = getRelativePos(e)
  startPoint.value = pos
  isDrawing.value = true
  selectedId.value = null

  if (activeTool.value === 'pen') {
    penPoints.value = [pos]
  }
}

const onMouseMove = (e: MouseEvent) => {
  if (!isDrawing.value) return
  const pos = getRelativePos(e)

  if (activeTool.value === 'pen') {
    penPoints.value.push(pos)
    return
  }

  const x = Math.min(startPoint.value.x, pos.x)
  const y = Math.min(startPoint.value.y, pos.y)
  const w = Math.abs(pos.x - startPoint.value.x)
  const h = Math.abs(pos.y - startPoint.value.y)

  if (activeTool.value === 'arrow') {
    drawPreview.value = {
      x: startPoint.value.x,
      y: startPoint.value.y,
      w: pos.x - startPoint.value.x,
      h: pos.y - startPoint.value.y,
    }
  } else {
    drawPreview.value = { x, y, w, h }
  }
}

const onMouseUp = (e: MouseEvent) => {
  if (!isDrawing.value) return
  const pos = getRelativePos(e)
  isDrawing.value = false

  if (activeTool.value === 'pen') {
    if (penPoints.value.length > 2) {
      commitAnnotation({
        type: 'pen',
        x: 0, y: 0, w: 0, h: 0,
        points: [...penPoints.value],
      })
    }
    penPoints.value = []
    return
  }

  let x: number, y: number, w: number, h: number

  if (activeTool.value === 'arrow') {
    x = startPoint.value.x
    y = startPoint.value.y
    w = pos.x - startPoint.value.x
    h = pos.y - startPoint.value.y
  } else {
    x = Math.min(startPoint.value.x, pos.x)
    y = Math.min(startPoint.value.y, pos.y)
    w = Math.abs(pos.x - startPoint.value.x)
    h = Math.abs(pos.y - startPoint.value.y)
  }

  const minSize = activeTool.value === 'arrow' ? 10 : 8
  if (Math.abs(w) < minSize && Math.abs(h) < minSize) {
    drawPreview.value = null
    return
  }

  commitAnnotation({ type: activeTool.value, x, y, w, h })
  drawPreview.value = null
}

const onMouseLeave = () => {
  if (isDrawing.value) {
    isDrawing.value = false
    drawPreview.value = null
    penPoints.value = []
  }
}

const commitAnnotation = (partial: Partial<Annotation> & Pick<Annotation, 'type' | 'x' | 'y' | 'w' | 'h'>) => {
  // Use array length + 1 so each new annotation gets a unique label (1, 2, 3...) regardless of ref timing
  const nextLabel = String(annotations.value.length + 1)
  const ann: Annotation = {
    id: `ann-${Date.now()}-${Math.random().toString(36).slice(2, 7)}`,
    type: partial.type,
    x: partial.x,
    y: partial.y,
    w: partial.w,
    h: partial.h,
    color: activeColor.value,
    strokeWidth: strokeWidth.value,
    label: nextLabel,
    comment: '',
    points: partial.points,
  }
  annotations.value.push(ann)
  selectedId.value = ann.id
}

// --- Badge position helpers ---
const badgeX = (ann: Annotation): number => {
  if (ann.type === 'pen' && ann.points?.length) return ann.points[0].x
  if (ann.type === 'arrow') return ann.x
  return ann.x + ann.w + 2
}

const badgeY = (ann: Annotation): number => {
  if (ann.type === 'pen' && ann.points?.length) return ann.points[0].y
  if (ann.type === 'arrow') return ann.y
  return ann.y - 2
}

// --- Annotation management ---
const selectAnnotation = (id: string) => {
  selectedId.value = selectedId.value === id ? null : id
}

const removeAnnotation = (id: string) => {
  annotations.value = annotations.value.filter(a => a.id !== id)
  if (selectedId.value === id) selectedId.value = null
}

const undo = () => {
  annotations.value.pop()
}

const clearAll = () => {
  annotations.value = []
  selectedId.value = null
}

const handleBackdropClick = () => {
  // Don't close on backdrop click (user might misclick while drawing)
}

const toolLabel = (type: string) => {
  const labels: Record<string, string> = {
    rect: 'Rectangle', circle: 'Circle', arrow: 'Arrow',
    highlight: 'Highlight', pen: 'Freehand',
  }
  return labels[type] || type
}

// --- Save: burn annotations onto canvas & export base64 ---
const saveAnnotated = async () => {
  const canvas = document.createElement('canvas')
  canvas.width = naturalWidth.value
  canvas.height = naturalHeight.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  // Scale factor between display and natural
  const scaleX = naturalWidth.value / displayWidth.value
  const scaleY = naturalHeight.value / displayHeight.value

  // Draw base image
  const img = new Image()
  img.crossOrigin = 'anonymous'
  await new Promise<void>((resolve) => {
    img.onload = () => resolve()
    img.src = props.imageSrc
  })
  ctx.drawImage(img, 0, 0, naturalWidth.value, naturalHeight.value)

  // Draw each annotation
  for (const ann of annotations.value) {
    ctx.save()
    ctx.strokeStyle = ann.color
    ctx.fillStyle = ann.color
    ctx.lineWidth = ann.strokeWidth * Math.max(scaleX, scaleY)
    ctx.lineCap = 'round'
    ctx.lineJoin = 'round'

    const x = ann.x * scaleX
    const y = ann.y * scaleY
    const w = ann.w * scaleX
    const h = ann.h * scaleY

    if (ann.type === 'rect') {
      ctx.strokeRect(x, y, w, h)
    } else if (ann.type === 'circle') {
      ctx.beginPath()
      ctx.ellipse(x + w / 2, y + h / 2, Math.abs(w / 2), Math.abs(h / 2), 0, 0, Math.PI * 2)
      ctx.stroke()
    } else if (ann.type === 'arrow') {
      drawArrowOnCanvas(ctx, x, y, x + w, y + h, ann.strokeWidth * Math.max(scaleX, scaleY))
    } else if (ann.type === 'highlight') {
      ctx.globalAlpha = 0.3
      ctx.fillRect(x, y, w, h)
      ctx.globalAlpha = 1
    } else if (ann.type === 'pen' && ann.points) {
      ctx.beginPath()
      ann.points.forEach((p, i) => {
        const px = p.x * scaleX
        const py = p.y * scaleY
        if (i === 0) ctx.moveTo(px, py)
        else ctx.lineTo(px, py)
      })
      ctx.stroke()
    }

    // Draw number badge
    const bx = badgeX(ann) * scaleX
    const by = badgeY(ann) * scaleY
    const radius = 14 * Math.max(scaleX, scaleY)
    ctx.beginPath()
    ctx.arc(bx, by, radius, 0, Math.PI * 2)
    ctx.fill()
    ctx.fillStyle = 'white'
    ctx.font = `bold ${14 * Math.max(scaleX, scaleY)}px system-ui, sans-serif`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(ann.label, bx, by)

    ctx.restore()
  }

  // Draw comment legend at bottom if any annotations have comments
  const withComments = annotations.value.filter(a => a.comment.trim())
  if (withComments.length > 0) {
    const padding = 20
    const lineH = 30 * scaleY
    const legendH = (withComments.length + 1) * lineH + padding * 2
    const legendY = naturalHeight.value - legendH

    ctx.fillStyle = 'rgba(0,0,0,0.7)'
    ctx.fillRect(0, legendY, naturalWidth.value, legendH)

    ctx.fillStyle = 'white'
    ctx.font = `bold ${14 * scaleY}px system-ui`
    ctx.textAlign = 'left'
    ctx.textBaseline = 'top'
    ctx.fillText('Annotations:', padding, legendY + padding)

    withComments.forEach((ann, i) => {
      const ty = legendY + padding + lineH * (i + 1)
      // Badge
      const bsize = 12 * scaleY
      ctx.fillStyle = ann.color
      ctx.beginPath()
      ctx.arc(padding + bsize, ty + bsize / 2, bsize, 0, Math.PI * 2)
      ctx.fill()
      ctx.fillStyle = 'white'
      ctx.font = `bold ${11 * scaleY}px system-ui`
      ctx.textAlign = 'center'
      ctx.fillText(ann.label, padding + bsize, ty + bsize / 2)
      // Comment text
      ctx.fillStyle = 'rgb(209, 213, 219)'
      ctx.font = `${13 * scaleY}px system-ui`
      ctx.textAlign = 'left'
      ctx.fillText(ann.comment, padding + bsize * 2 + 8, ty)
    })
  }

  const dataUrl = canvas.toDataURL('image/png')
  emit('save', dataUrl)
}

const drawArrowOnCanvas = (ctx: CanvasRenderingContext2D, x1: number, y1: number, x2: number, y2: number, lw: number) => {
  const angle = Math.atan2(y2 - y1, x2 - x1)
  const headLen = Math.min(20 * lw, 30)

  ctx.beginPath()
  ctx.moveTo(x1, y1)
  ctx.lineTo(x2, y2)
  ctx.stroke()

  ctx.beginPath()
  ctx.moveTo(x2, y2)
  ctx.lineTo(x2 - headLen * Math.cos(angle - Math.PI / 6), y2 - headLen * Math.sin(angle - Math.PI / 6))
  ctx.lineTo(x2 - headLen * Math.cos(angle + Math.PI / 6), y2 - headLen * Math.sin(angle + Math.PI / 6))
  ctx.closePath()
  ctx.fill()
}

// Resize observer
onMounted(() => {
  if (canvasWrapperRef.value) {
    const ro = new ResizeObserver(() => fitToWrapper())
    ro.observe(canvasWrapperRef.value)
    onBeforeUnmount(() => ro.disconnect())
  }
})
</script>

<style scoped>
.annotator-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.annotator-modal {
  background: rgb(17 24 39);
  border: 1px solid rgb(55 65 81);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  width: min(1200px, 100%);
  height: min(90vh, 800px);
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0,0,0,0.6);
}

.annotator-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid rgb(31 41 55);
  background: rgb(15 23 42);
}

.annotator-close-btn {
  color: rgb(107 114 128);
  padding: 4px;
  border-radius: 6px;
  transition: all 0.15s;
}
.annotator-close-btn:hover {
  color: white;
  background: rgb(31 41 55);
}

.annotator-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: rgb(15 23 42);
  border-bottom: 1px solid rgb(31 41 55);
}

.toolbar-sep {
  width: 1px;
  height: 22px;
  background: rgb(55 65 81);
  margin: 0 4px;
}

.tool-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 6px;
  color: rgb(156 163 175);
  background: transparent;
  transition: all 0.1s;
  border: none;
  cursor: pointer;
}
.tool-btn:hover:not(:disabled) {
  background: rgb(55 65 81);
  color: white;
}
.tool-btn.is-active {
  background: rgb(37 99 235);
  color: white;
}
.tool-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.annotator-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.canvas-wrapper {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgb(9 14 26);
  padding: 16px;
}

.canvas-container {
  position: relative;
  flex-shrink: 0;
  user-select: none;
}

.canvas-image {
  display: block;
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  user-select: none;
}

.annotation-svg {
  position: absolute;
  top: 0;
  left: 0;
}

.annotation-shape {
  cursor: pointer;
  transition: opacity 0.1s;
}
.annotation-shape:hover {
  opacity: 0.8;
}
.selected-shape {
  filter: drop-shadow(0 0 4px rgba(255,255,255,0.5));
}

.comment-badge {
  cursor: pointer;
}

/* Comment panel */
.comment-panel {
  width: 260px;
  display: flex;
  flex-direction: column;
  border-left: 1px solid rgb(31 41 55);
  background: rgb(15 23 42);
  overflow: hidden;
}

.comment-panel-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px 14px;
  border-bottom: 1px solid rgb(31 41 55);
  font-weight: 600;
}

.comment-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.comment-list::-webkit-scrollbar { width: 4px; }
.comment-list::-webkit-scrollbar-track { background: transparent; }
.comment-list::-webkit-scrollbar-thumb { background: rgb(55 65 81); border-radius: 2px; }

.comment-item {
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  border: 1px solid transparent;
  margin-bottom: 6px;
  transition: all 0.1s;
  background: rgb(17 24 39);
}
.comment-item:hover {
  border-color: rgb(55 65 81);
}
.selected-comment {
  border-color: rgb(59 130 246) !important;
  background: rgb(23 37 84) !important;
}

.comment-badge-small {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  color: white;
  margin-top: 2px;
}

.comment-textarea {
  width: 100%;
  background: rgb(9 14 26);
  border: 1px solid rgb(55 65 81);
  border-radius: 4px;
  color: rgb(209 213 219);
  font-size: 12px;
  line-height: 1.5;
  padding: 5px 7px;
  resize: none;
  outline: none;
  transition: border-color 0.15s;
}
.comment-textarea:focus {
  border-color: rgb(59 130 246);
}
.comment-textarea::placeholder {
  color: rgb(75 85 99);
}

.comment-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  font-size: 11px;
  color: rgb(75 85 99);
  border-top: 1px solid rgb(31 41 55);
}
</style>
