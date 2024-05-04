<script setup lang="ts">
defineProps({
  customClass: {
    type: String,
    required: false,
    default: 'bg-gray-100'
  }
})
</script>

<template>
  <div :class="'p-2 rounded-xl overflow-hidden ' + customClass">
    <slot />
  </div>
</template>

<style scoped>

</style>