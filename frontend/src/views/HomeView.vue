<script setup lang="ts">
import { Search } from '@icon-park/vue-next'
</script>

<template>
  <div class="bg-white bg-opacity-0 backdrop-blur-lg absolute inset-0 p-6">
    <div class="flex flex-col gap-4">
      <!-- header -->
      <div class="flex items-center justify-between">
        <div class="flex flex-col gap-2">
          <div class="text-3xl font-medium text-blue-500 dark:text-red-200">Hello Aurora {{ $t("language") }}</div>
          <div class="text-xs text-gray-500">Good day for you</div>
        </div>
        <div class="bg-white bg-opacity-30 rounded-xl py-3 px-6 backdrop-blur-lg flex justify-between min-w-[220px]">
          <div class="text-gray-500">{{ $t("search") }}</div>
          <div class="cursor-pointer">
            <search />
          </div>
        </div>
      </div>

      <el-divider />
      <!-- content -->
    </div>
  </div>
</template>
