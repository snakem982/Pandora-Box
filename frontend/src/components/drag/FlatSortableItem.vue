<script setup lang="ts">
import {FlatSortableItemProps} from "@/types/drag";
const props = defineProps<FlatSortableItemProps & { class?: string, inset?: boolean, draggable?: boolean }>()
</script>

<template>
  <div
      class="flat-sortable-item"
      :class="props.class"
      :draggable="props.draggable"
      style="will-change: transform;"
  >
    <slot/>
  </div>
</template>
