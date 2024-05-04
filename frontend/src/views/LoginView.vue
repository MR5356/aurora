<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import GithubIcon from '@/components/icons/GithubIcon.vue'
import GitlabIcon from '@/components/icons/GitlabIcon.vue'
import { User } from '@/views/LoginView'

const systemStore = useSystemStore()
const router = useRouter()

const availableOAuth = ref<User.AvailableOAuth[]>([])

const getAvailableOAuth = async () => {
  let oauths = await User.getAvailableOAuth()
  oauths.sort((a, b) => a.type.localeCompare(b.type))
  availableOAuth.value = oauths
}

const onClickOAuthLogin = async (oauth: string) => {
  let url = await User.getOAuthURL(oauth, router.currentRoute.value.query.redirectURL as string)
  window.open(url, '_self')
}

const init = async () => {
  let userInfo = await User.getUserInfo()
  if (userInfo) {
    await router.replace(router.currentRoute.value.query.redirectURL as string ?? '/')
  } else {
    await getAvailableOAuth()
  }
}

init()

</script>

<template>
  <div
    class="fixed inset-0 select-none overflow-hidden transition duration-300 ease-in-out bg-gradient-to-bl from-indigo-300 dark:from-slate-900 to-slate-200 dark:to-slate-700">
    <div class="absolute top-0 bottom-0 left-0 right-0 bg-opacity-40 bg-white h-fit w-[460px] m-auto p-8 rounded-2xl">
      <div class="flex flex-col gap-16">
        <div class="flex flex-col gap-4">
          <router-link to="/">
            <div class="flex items-center">
              <img class="w-8 h-8" :src="systemStore.website.logo" alt="logo">
              <div class="text-2xl font-bold text-teal-400">{{ systemStore.website.title }}.</div>
            </div>
          </router-link>
          <div class="text-5xl font-bold">{{ $t('signIn') }}</div>
        </div>
        <div class="flex flex-col gap-0 font-mono subpixel-antialiased items-center w-full">
          <div v-for="(item, index) in availableOAuth" :key="item.oauth" class="w-full">
            <div class="flex items-center justify-between w-full" @click="onClickOAuthLogin(item.oauth)">
              <div
                class="cursor-pointer bg-blue-600 text-white rounded-lg p-3 px-6 flex items-center gap-2 w-full justify-center hover:shadow-lg">
                <github-icon v-if="item.type === 'github'" class="w-4 h-4" />
                <gitlab-icon v-if="item.type === 'gitlab'" class="w-4 h-4" />
                <div v-if="item.type === 'github'">{{ $t('signWithGithub').replaceAll('Github', item.oauth) }}</div>
                <div v-if="item.type === 'gitlab'">{{ $t('signWithGitlab').replaceAll('Gitlab', item.oauth) }}</div>
              </div>
            </div>
            <div v-if="index !== availableOAuth.length - 1" class="text-sm text-gray-500 text-center py-1">OR</div>
          </div>
        </div>
        <div class="text-center text-xs text-gray-500" v-html="$t('signInNote')" />
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>