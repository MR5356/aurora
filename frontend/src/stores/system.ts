import { ref } from 'vue'
import { defineStore } from 'pinia'
import { SystemModels } from '@/layouts/SystemModels'

export const useSystemStore = defineStore('system', () => {
    const website = ref<SystemModels.Website>(SystemModels.defaultWebsite)
    const navigation = ref<SystemModels.Navigation[]>(SystemModels.defaultNavigation)
    const language = ref<string>('en')

    function setWebsite(value: SystemModels.Website = SystemModels.defaultWebsite) {
      website.value = value
      SystemModels.setWebsite(value)
    }

    function setNavigation(value: SystemModels.Navigation[]) {
      navigation.value = value
    }

    function setLanguage(value: string) {
      language.value = value
    }

    return { website, setWebsite, navigation, setNavigation, language, setLanguage }
  },
  {
    persist: true
  }
)