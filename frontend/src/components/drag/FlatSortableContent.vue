<script setup lang="ts">
import {nextTick, onMounted, ref, watch} from 'vue'

/********************************************************/
/*                 dragstart 开始拖拽                    */
/*                 dragend 结束拖拽                      */
/*  来自 https://github.com/pinky-pig/vue-flat-sortable  */
/********************************************************/

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

// 容器 DOM
const containerRef = ref<HTMLElement | null>(null)

// 当前元素的节点
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

// 是否正在拖拽
const isDragging = ref<boolean>(false)

onMounted(() => {
  initFlatDom()
})
watch(() => props.vl, initFlatDom)

function initFlatDom() {
  nextTick(() => {
    // 给拖拽的元素添加类名，这样排序后的顺序可以知道
    if (!containerRef.value)
      return
    if (!props.vl) {
      console.warn('FlatSortableContent modelValue is required')
      return
    }

    const flatItems = Array.from(containerRef.value.children)
        ?.filter(el => el.classList.contains('flat-sortable-item'))

    // if (flatItems.length !== props.modelValue.length) {
    //   console.warn(``)
    //   return
    // }

    flatItems?.forEach((el, index) => el.classList.add(`flat-sortable-content-${index}`))

    // 不是flatItem的元素，放置末尾
    // const nonFlatItem = Array.from(containerRef.value.children)
    //     ?.filter(el => !el.classList.contains('flat-sortable-item')) as HTMLElement[]
    // nonFlatItem.forEach((el) => {
    //   const insertBeforeElement = null // 末尾位置为 null
    //   containerRef.value!.insertBefore(el, insertBeforeElement)
    // })
  })
}

// 给拖拽元素设置 pointerEvents 为 none ，以防 dragenter 触发的子元素
watch(isDragging, (v) => {
  if (!containerRef.value)
    return
  const flatItems = Array.from(containerRef.value.children)
      ?.filter(el => el.classList.contains('flat-sortable-item'))
  flatItems.forEach((element) => {
    Array.from(element.children).forEach((child) => {
      (child as HTMLElement).style.pointerEvents = v ? 'none' : 'auto'
    })
  })
})

/********************************************************/
/*                 1. 占位 DOM                          */
/*                 2. 跟随鼠标 DOM                       */
/*                 2. 拖拽过渡，结束拖拽过渡               */
/********************************************************/

function handleDragstart(e: DragEvent) {
  //  如果拖拽的不是 FlatSortableItem 则不进行拖拽
  if (!isFlatSortableItem(e.target as HTMLElement))
    return

  // 初始化 draggedNode 的状态
  draggedNode.value.el = e.target as HTMLElement

  setTimeout(() => {
    // 添加 draggedNode 样式
    if (draggedNode.value && draggedNode.value.el)
      draggedNode.value?.el.classList.add('sortable-chosen')
  })

  e.dataTransfer!.effectAllowed = 'move'
  isDragging.value = true
}

function handleDragEnter(e: DragEvent) {
  e.preventDefault()
  // 如果拖拽的不是 FlatSortableItem 则不进行拖拽
  if (!isFlatSortableItem(e.target as HTMLElement))
    return
  // 如果没有 el ，也是不进行碰撞检测
  if (!draggedNode.value.el || draggedNode.value.el === e.target || e.target === containerRef.value || !isFlatSortableItem(e.target as HTMLElement))
    return

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

/********************************************************/
/*                       utils                          */
/********************************************************/

function isFlatSortableItem(el: HTMLElement) {
  return el.classList.contains('flat-sortable-item')
}

function recordSingle(el: HTMLElement | Element): INodeType {
  const {top, left, width, height, right, bottom} = el.getBoundingClientRect()
  return {top, left, width, height, el, right, bottom}
}

/**
 * 这里碰撞检测比较简单，但是动画比较繁琐
 * 首先碰撞就是 dragenter 已经拿到了，不需要再操作了
 * 其次是交换位置只需要 insertBefore 插入就行
 * 动画这块，直接使用 animates
 * 但是因为设计到快速多次动画触发，所以造成动画异常
 * @param originNode 拖拽的元素
 * @param targetNode 碰撞的元素
 * @param allNodes 所有的容器内的子元素
 */
async function hitTest(originNode: HTMLElement, targetNode: HTMLElement, allNodes: HTMLElement[]) {
  // 判断当前碰撞的元素是否在动画中，如果是，那么就跳过
  const targetIsAnimating = targetNode.getAttribute('data-animating')
  if (targetIsAnimating === 'true')
    return

  const currentIndex = allNodes.indexOf(originNode)
  const targetIndex = allNodes.indexOf(targetNode)

  // 在中间的元素，添加动画的标志
  allNodes.filter((node, index) => {
    return index >= Math.min(currentIndex, targetIndex) && index <= Math.max(currentIndex, targetIndex)
  }).forEach((node) => {
    node.setAttribute('data-animating', 'true')
  })

  // 过滤出 index 最前面的元素的之后所有的元素，为后面开始动画
  const filterNodes = allNodes.filter((node, index) => {
    return index >= Math.min(currentIndex, targetIndex)
  })

  const firsts = filterNodes.map((node) => {
    // 1. 如果当前的元素有动画效果，那么就要以动画效果的位置为初始
    const last = recordSingle(node)
    const animation = node.getAnimations()[0]
    if (animation)
      animation.cancel()

    return last
  })

  /** 更改 DOM start */
  if (currentIndex < targetIndex)
    targetNode.parentElement?.insertBefore(originNode, targetNode.nextSibling)
  else
    targetNode.parentElement?.insertBefore(originNode, targetNode)
  /** 更改 DOM end */

  nextTick(async () => {
    /** 更改绑定的 class 数组 start */
    const updatedAllNodes = (Array.from(containerRef.value!.children) as HTMLElement[]).filter(node => isFlatSortableItem(node))
    const updatedFlatSortableContent = updatedAllNodes.map((el) => {
      const classes = Array.from(el.classList)
      const matchingClasses = classes.filter(className => className.startsWith('flat-sortable-content-'))
      return matchingClasses[0].split('flat-sortable-content-')[1]
    })
    emits('update:modelValue', updatedFlatSortableContent)
    /** 更改绑定的 class 数组 end */

    const lasts = filterNodes.map((node) => {
      return recordSingle(node)
    })

    if (currentIndex > targetIndex) {
      // 说明拖拽的元素大于碰撞的元素，那么是插入其前面，动画从后面开始播放
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
      easing: 'linear',
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
        {duration: options.duration, easing: options.easing, delay: options.delay, fill: 'backwards'},
    )
    animation.onfinish = () => {
      // 标志位，结束动画
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
}
</style>
