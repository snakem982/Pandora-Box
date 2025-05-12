<template>
  <div class="vdc-out-container" :style="`width:${props.width}`">
    <TransitionGroup name="fade" tag="div" class="vdc-trans-group-container"
                     :style="{ '--vdc-gap': gap + 'px',
                     '--vdc-top': top + 'px'
    }"
    >
      <div class="vdc-item-container"
           :draggable="draggable"
           v-for="(item, index) in items"
           :key="item"
           @dragstart="drag($event, index)"
           @dragover="over"
           @drop="drop($event, index)">
        <slot
            name="VDC"
            :data="item"
            :index="index"></slot>
      </div>
    </TransitionGroup>
  </div>
</template>
<script setup lang="ts">
import {reactive, ref} from 'vue'
import {eType, istate} from './VDContainer'

// eslint-disable-next-line no-undef, no-unused-vars
// 使用 defineProps 和 withDefaults
const props = withDefaults(defineProps<{
  width?: any,
  data: any[],
  row?: number,
  gap?: number,
  top?: number,
  type?: string,
  draggable?: boolean,
}>(), {
  width: '100%',         // 默认值：宽度
  gap: 10,               // 默认值：间距为 10
  top: 0,                // 默认值：顶部偏移为 0
  type: 'sort',       // 默认值：类型为 "default"
  draggable: false       // 默认值：不可拖拽
});


const state: istate = reactive({
  ...props,
  target: 0
})
const getItems = () => props.data
const items = ref(getItems())

// eslint-disable-next-line no-undef
const emit = defineEmits([
  'getData'
])

// while target start dragged
const drag = (event: DragEvent, index: number) => {
  state.target = index
}

// while target is on the  drop point
const over = (event: DragEvent) => {
  event.preventDefault()
}

// while drop the object into target
const drop = (event: DragEvent, index: number) => {
  if (props.type === eType.SORT) {
    items.value.splice(index, 0, items.value.splice(state.target, 1)[0])
  } else if (props.type === eType.SWITCH) {
    items.value[index] = items.value.splice(state.target, 1, items.value[index])[0]
  } else {
    window.console.error("wrong type name,check <VDContainer></VDContainer>element's [type] modal")
  }
  emit('getData', items.value)
}

</script>

<style>
.vdc-trans-group-container {
  display: flex;
  flex-wrap: inherit;
  -ms-flex-wrap: inherit;
  gap: var(--vdc-gap);
  margin-top: var(--vdc-top);
  margin-bottom: var(--vdc-top);
  width: 100%;
}

.vdc-out-container {
  display: flex;
  flex-wrap: wrap;
  -ms-flex-wrap: wrap;
}

/* 动画类名统一，避免冗余 */
.fade-move,
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s cubic-bezier(0.55, 0, 0.1, 1),
  transform 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

/* 进入和离开状态声明*/
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(30px);
}


/* 3. 确保离开的项目被移除出了布局流
      以便正确地计算移动时的动画效果。 */
.fade-leave-active {
  position: absolute;
}
</style>
