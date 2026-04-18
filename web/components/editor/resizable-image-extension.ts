import Image from '@tiptap/extension-image'

const MIN_SIZE = 60
const MAX_SIZE = 1200

type HandlePosition = 'n' | 'ne' | 'e' | 'se' | 's' | 'sw' | 'w' | 'nw'

function parseSize(val: unknown): number | null {
  if (val == null) return null
  const n = typeof val === 'number' ? val : parseInt(String(val), 10)
  return Number.isNaN(n) ? null : n
}

const HANDLES: { position: HandlePosition; cursor: string }[] = [
  { position: 'n', cursor: 'ns-resize' },
  { position: 'ne', cursor: 'nesw-resize' },
  { position: 'e', cursor: 'ew-resize' },
  { position: 'se', cursor: 'nwse-resize' },
  { position: 's', cursor: 'ns-resize' },
  { position: 'sw', cursor: 'nesw-resize' },
  { position: 'w', cursor: 'ew-resize' },
  { position: 'nw', cursor: 'nwse-resize' },
]

function createHandle(
  position: HandlePosition,
  cursor: string,
  editor: { isEditable: boolean },
  img: HTMLImageElement,
  wrapper: HTMLDivElement,
  node: { attrs: Record<string, unknown> },
  getPos: () => number | undefined,
  applyUpdate: (w: number, h: number, tx?: number, ty?: number) => void
) {
  const handle = document.createElement('div')
  handle.className = 'resize-handle'
  handle.setAttribute('data-resize-handle', position)
  handle.setAttribute('aria-label', `Resize ${position}`)
  handle.style.cursor = cursor

  let startX = 0
  let startY = 0
  let startW = 0
  let startH = 0

  const onMouseDown = (e: MouseEvent) => {
    if (e.button !== 0 || !editor.isEditable) return
    e.preventDefault()
    e.stopPropagation()
    const rect = img.getBoundingClientRect()
    startX = e.clientX
    startY = e.clientY
    startW = rect.width
    startH = rect.height
    wrapper.classList.add('resizing')
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
  }

  const onMouseMove = (e: MouseEvent) => {
    const dx = e.clientX - startX
    const dy = e.clientY - startY
    const aspect = startW / startH

    let newW = startW
    let newH = startH
    let translateX: number | undefined
    let translateY: number | undefined

    switch (position) {
      case 'e':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW + dx))
        break
      case 'w':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW - dx))
        translateX = startW - newW
        break
      case 's':
        newH = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startH + dy))
        break
      case 'n':
        newH = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startH - dy))
        translateY = startH - newH
        break
      case 'se':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW + dx))
        newH = newW / aspect
        if (newH < MIN_SIZE) {
          newH = MIN_SIZE
          newW = newH * aspect
        } else if (newH > MAX_SIZE) {
          newH = MAX_SIZE
          newW = newH * aspect
        }
        break
      case 'sw':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW - dx))
        newH = newW / aspect
        if (newH < MIN_SIZE) {
          newH = MIN_SIZE
          newW = newH * aspect
        } else if (newH > MAX_SIZE) {
          newH = MAX_SIZE
          newW = newH * aspect
        }
        translateX = startW - newW
        break
      case 'ne':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW + dx))
        newH = newW / aspect
        if (newH < MIN_SIZE) {
          newH = MIN_SIZE
          newW = newH * aspect
        } else if (newH > MAX_SIZE) {
          newH = MAX_SIZE
          newW = newH * aspect
        }
        translateY = -(newH - startH)
        break
      case 'nw':
        newW = Math.max(MIN_SIZE, Math.min(MAX_SIZE, startW - dx))
        newH = newW / aspect
        if (newH < MIN_SIZE) {
          newH = MIN_SIZE
          newW = newH * aspect
        } else if (newH > MAX_SIZE) {
          newH = MAX_SIZE
          newW = newH * aspect
        }
        translateX = startW - newW
        translateY = -(newH - startH)
        break
    }

    applyUpdate(newW, newH, translateX, translateY)
  }

  const onMouseUp = () => {
    wrapper.classList.remove('resizing')
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    const w = Math.round(parseFloat(img.style.width) || img.offsetWidth)
    const h = Math.round(parseFloat(img.style.height) || img.offsetHeight)
    const pos = typeof getPos === 'function' ? getPos() : undefined
    if (pos !== undefined && editor?.view) {
      editor.view.dispatch(
        editor.view.state.tr.setNodeMarkup(pos, undefined, {
          ...node.attrs,
          width: w,
          height: h,
        })
      )
    }
    img.style.transform = ''
    img.style.transformOrigin = ''
  }

  handle.addEventListener('mousedown', onMouseDown)
  return handle
}

export function createResizableImageExtension() {
  return Image.extend({
    addAttributes() {
      return {
        ...this.parent?.(),
        width: {
          default: null,
          parseHTML: (el) => parseSize((el as HTMLElement).getAttribute('width')),
          renderHTML: (attrs) => (attrs.width != null ? { width: String(attrs.width) } : {}),
        },
        height: {
          default: null,
          parseHTML: (el) => parseSize((el as HTMLElement).getAttribute('height')),
          renderHTML: (attrs) => (attrs.height != null ? { height: String(attrs.height) } : {}),
        },
      }
    },

    addNodeView() {
      return ({ node, getPos, editor }) => {
        const wrapper = document.createElement('div')
        wrapper.className = 'resizable-image-wrapper'

        const img = document.createElement('img')
        img.src = node.attrs.src ?? ''
        img.className = 'editor-image'
        img.loading = 'lazy'
        img.decoding = 'async'
        img.draggable = false
        img.setAttribute('data-drag-handle', 'true')

        const w = parseSize(node.attrs.width)
        const h = parseSize(node.attrs.height)
        if (w != null) img.style.width = `${w}px`
        if (h != null) img.style.height = `${h}px`

        wrapper.appendChild(img)

        const applyUpdate = (newW: number, newH: number, tx?: number, ty?: number) => {
          img.style.width = `${newW}px`
          img.style.height = `${newH}px`
          if (tx !== undefined || ty !== undefined) {
            const x = tx ?? 0
            const y = ty ?? 0
            img.style.transformOrigin = 'top left'
            img.style.transform = `translate(${x}px, ${y}px)`
          } else {
            img.style.transform = ''
            img.style.transformOrigin = ''
          }
        }

        if (editor.isEditable) {
          for (const { position, cursor } of HANDLES) {
            const handle = createHandle(
              position,
              cursor,
              editor,
              img,
              wrapper,
              node,
              getPos,
              applyUpdate
            )
            wrapper.appendChild(handle)
          }
        }

        return {
          dom: wrapper,
          update(updatedNode) {
            if (updatedNode.type !== node.type) return false
            if (updatedNode.attrs.src !== img.src) img.src = updatedNode.attrs.src ?? ''
            const nw = parseSize(updatedNode.attrs.width)
            const nh = parseSize(updatedNode.attrs.height)
            if (nw != null) img.style.width = `${nw}px`
            else img.style.width = ''
            if (nh != null) img.style.height = `${nh}px`
            else img.style.height = ''
            img.style.transform = ''
            img.style.transformOrigin = ''
            return true
          },
        }
      }
    },
  })
}
