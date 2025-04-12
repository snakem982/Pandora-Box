<script setup lang="ts">
import { nextTick, onMounted, ref, watch } from 'vue'

const props = defineProps<FlatSortableContentProps>()
const emits = defineEmits(['update:modelValue'])

interface FlatSortableContentProps {
  vl: number
  class?: string
  direction?: 'row' | 'row-reverse' | 'column' | 'column-reverse'
  gap?: number
}

interface INodeType {
  top: number
  left: number
  bottom: number
  right: number
  width: number
  height: number
  el: HTMLElement | Element | null
}

interface IDraggedNodeType {
  top: number
  left: number
  width: number
  height: number
  el: HTMLElement | Element | null
  shadowEl: HTMLElement | Element | null
  offsetXFromMouse: number
  offsetYFromMouse: number
}

const containerRef = ref<HTMLElement | null>(null)
const draggedNode = ref<IDraggedNodeType>({
  top: 0,
  left: 0,
  width: 0,
  height: 0,
  offsetXFromMouse: 0,
  offsetYFromMouse: 0,
  el: null,
  shadowEl: null,
})
const isDragging = ref<boolean>(false)

onMounted(() => {
  initFlatDom()
})
watch(() => props.vl, initFlatDom)

function initFlatDom() {
  nextTick(() => {
    if (!containerRef.value) return
    if (!props.vl) {
      console.warn('FlatSortableContent modelValue is required')
      return
    }

    const flatItems = Array.from(containerRef.value.children)?.filter(el =>
      el.classList.contains('flat-sortable-item')
    )
    flatItems?.forEach((el, index) =>
      el.classList.add(`flat-sortable-content-${index}`)
    )
  })
}

watch(isDragging, (v) => {
  if (!containerRef.value) return
  const flatItems = Array.from(containerRef.value.children)?.filter(el =>
    el.classList.contains('flat-sortable-item')
  )
  flatItems.forEach((element) => {
    Array.from(element.children).forEach((child) => {
      (child as HTMLElement).style.pointerEvents = v ? 'none' : 'auto'
    })
  })
})

function handleDragstart(e: DragEvent) {
  if (!isFlatSortableItem(e.target as HTMLElement)) return
  draggedNode.value.el = e.target as HTMLElement

  setTimeout(() => {
    if (draggedNode.value && draggedNode.value.el)
      draggedNode.value?.el.classList.add('sortable-chosen')
  })

  e.dataTransfer!.effectAllowed = 'move'
  isDragging.value = true
}

function handleDragEnter(e: DragEvent) {
  e.preventDefault()
  if (!isFlatSortableItem(e.target as HTMLElement)) return
  if (!draggedNode.value.el || draggedNode.value.el === e.target || e.target === containerRef.value || !isFlatSortableItem(e.target as HTMLElement)) return

  const allNodes = (Array.from(containerRef.value!.children) as HTMLElement[]).filter(node => isFlatSortableItem(node))
  hitTest(draggedNode.value.el as HTMLElement, e.target as HTMLElement, allNodes)
}

function handleDragOver(e: DragEvent) {
  e.preventDefault()
}

function handleDragEnd(_e: DragEvent) {
  if (draggedNode.value && draggedNode.value.el) {
    draggedNode.value?.el.classList.remove('sortable-chosen')
    draggedNode.value.el = null
    isDragging.value = false
  }
}

function isFlatSortableItem(el: HTMLElement) {
  return el.classList.contains('flat-sortable-item')
}

function recordSingle(el: HTMLElement | Element): INodeType {
  const { top, left, width, height, right, bottom } = el.getBoundingClientRect()
  return { top, left, width, height, el, right, bottom }
}

async function hitTest(originNode: HTMLElement, targetNode: HTMLElement, allNodes: HTMLElement[]) {
  const targetIsAnimating = targetNode.getAttribute('data-animating')
  if (targetIsAnimating === 'true') return

  const currentIndex = allNodes.indexOf(originNode)
  const targetIndex = allNodes.indexOf(targetNode)

  allNodes.filter((node, index) => {
    return index >= Math.min(currentIndex, targetIndex) && index <= Math.max(currentIndex, targetIndex)
  }).forEach((node) => {
    node.setAttribute('data-animating', 'true')
  })

  const filterNodes = allNodes.filter((node, index) => {
    return index >= Math.min(currentIndex, targetIndex)
  })

  const firsts = filterNodes.map((node) => {
    const last = recordSingle(node)
    const animation = node.getAnimations()[0]
    if (animation) animation.cancel()
    return last
  })

  if (currentIndex < targetIndex)
    targetNode.parentElement?.insertBefore(originNode, targetNode.nextSibling)
  else
    targetNode.parentElement?.insertBefore(originNode, targetNode)

  nextTick(async () => {
    const updatedAllNodes = (Array.from(containerRef.value!.children) as HTMLElement[]).filter(node => isFlatSortableItem(node))
    const updatedFlatSortableContent = updatedAllNodes.map((el) => {
      const classes = Array.from(el.classList)
      const matchingClasses = classes.filter(className => className.startsWith('flat-sortable-content-'))
      return matchingClasses[0].split('flat-sortable-content-')[1]
    })
    emits('update:modelValue', updatedFlatSortableContent)

    const lasts = filterNodes.map((node) => {
      return recordSingle(node)
    })

    if (currentIndex > targetIndex) {
      for (let i = filterNodes.length - 1; i >= 0; i--) {
        const node = filterNodes[i]
        const first = firsts[i]
        const last = lasts[i]
        const diff = {
          top: last.top - first.top,
          left: last.left - first.left,
        }
        animateElement(node, diff)
      }
    } else {
      for (let i = 0; i < filterNodes.length; i++) {
        const node = filterNodes[i]
        const first = firsts[i]
        const last = lasts[i]
        const diff = {
          top: last.top - first.top,
          left: last.left - first.left,
        }
        animateElement(node, diff)
      }
    }
  })
}

async function animateElement(
  element: HTMLElement,
  diff: { top: number, left: number },
  options: {
    reverse?: boolean
    duration?: number
    easing?: 'linear' | 'ease' | 'ease-in' | 'ease-out' | 'ease-in-out' | 'step-start' | 'step-end' | string
    delay?: number
  } = {
    reverse: true,
    duration: 300,
    easing: 'cubic-bezier(0.25, 0.8, 0.25, 1)', // 优化缓动效果
    delay: 0,
  },
) {
  return new Promise<void>((resolve) => {
    const animates = [
      `translate3d(${-diff.left}px, ${-diff.top}px,0px)`,
      'translate3d(0px, 0px, 0px)',
    ]
    const animation = element.animate(
      {
        transform: options.reverse
          ? animates
          : [...animates].reverse(),
      },
      { duration: options.duration, easing: options.easing, delay: options.delay, fill: 'backwards' },
    )
    animation.onfinish = () => {
      element.removeAttribute('data-animating')
      resolve()
    }
  })
}
</script>

<template>
  <div
    ref="containerRef" :class="props.class" class="translate-x-0" :style="{
      display: 'flex',
      flexDirection: props.direction || 'column',
      gap: `${props.gap || 0}px`,
      transform: 'skew(0)',
    }"
    @dragstart="handleDragstart"
    @dragenter="handleDragEnter"
    @dragover="handleDragOver"
    @dragend="handleDragEnd"
  >
    <slot/>
  </div>
</template>

<style scoped>
.sortable-chosen {
  will-change: transform;
  pointer-events: none !important;
  opacity: 0.2;
  transition: opacity 0.2s ease-in-out; /* 优化透明度过渡 */
}
</style>
