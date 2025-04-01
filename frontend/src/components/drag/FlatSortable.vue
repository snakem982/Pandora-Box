<script setup lang="ts">
import {FlatSortableProps} from "@/types/drag";
const props = defineProps<FlatSortableProps & { class?: string }>()
</script>

<template>
  <div :class="props.class">
    <slot />
  </div>
</template>

<style scoped>

</style>
