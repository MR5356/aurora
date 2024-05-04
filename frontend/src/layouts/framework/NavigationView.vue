<script setup lang="ts">
import { useSystemStore } from '@/stores/system'
import NavigationItemView from '@/layouts/framework/NavigationItemView.vue'

const systemStore = useSystemStore()

</script>

<template>
  <div class="flex flex-col gap-2">
    <div v-for="item in systemStore.navigation" :key="item.title">
      <NavigationItemView :item="item" />
    </div>
  </div>
</template>

<style scoped>

</style>