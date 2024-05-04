<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { SystemModels } from '@/layouts/SystemModels'

useDark()

const systemStore = useSystemStore()
const { locale } = useI18n()

systemStore.setWebsite()
systemStore.setNavigation(SystemModels.defaultNavigation)

const setLanguage = (value: string = systemStore.language) => {
  locale.value = value
  systemStore.setLanguage(value)
}

setLanguage()
</script>

<template>
  <div>
    <router-view v-slot="{ Component }">
      <transition name="slide-fade" :appear="true">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>

<style scoped>

</style>
