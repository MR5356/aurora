<script setup lang="ts">
import { ref } from 'vue'
import { RouterView, RouterLink, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import ColorView from '@/layouts/framework/ColorView.vue'
import NavigationView from '@/layouts/framework/NavigationView.vue'
import { useDark, useToggle } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { locales } from '@/lang/i18n'
import { User } from '@/views/LoginView'
import { Logout, Avatar } from '@icon-park/vue-next'
import { ElMessageBox } from 'element-plus'

const { locale, t } = useI18n()
const systemStore = useSystemStore()
const isDark = useDark()
const toggleDark = useToggle(isDark)
const router = useRouter()

const userInfo = ref<User.UserInfo | null>(null)

const getUserInfo = async () => {
  userInfo.value = await User.getUserInfo()
}

const setLanguage = (value: string = systemStore.language) => {
  systemStore.setLanguage(value)
  locale.value = value
}

const onClickLogout = () => {
  ElMessageBox.confirm(t('confirmLogout'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning'
  }).then(async () => {
    await User.logout()
    router.go(0)
  })
}

getUserInfo()
setLanguage()
</script>

<template>
  <el-container
    class="fixed inset-0 overflow-hidden transition duration-300 ease-in-out bg-gradient-to-bl from-indigo-100 dark:from-slate-900 to-slate-50 dark:to-slate-700">
    <el-aside
      class="transition duration-300 ease-in-out text-slate-600 dark:text-slate-400 bg-opacity-20 bg-white dark:bg-slate-900 relative">
      <div class="px-6 py-4 flex flex-col justify-between h-full">
        <!-- start -->
        <div class="flex flex-col gap-8 flex-grow">
          <!-- website title -->
          <div class="flex items-center justify-between">
            <router-link to="/">
              <div class="flex gap-2 items-center">
                <img class="w-8 h-8" :src="systemStore.website.logo" alt="logo">
                <div class="text-2xl font-bold text-teal-400">{{ systemStore.website.title }}.</div>
              </div>
            </router-link>
            <div class="flex items-center gap-3">
              <div class="cursor-pointer opacity-80">
                <el-popover trigger="hover">
                  <template #reference>
                    <img
                      src="data:image/svg+xml;charset=utf-8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHdpZHRoPSIyMiIgaGVpZ2h0PSIyMiIgdmlld0JveD0iMCAwIDQ4IDQ4IiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0yOC4yODU3IDM3SDM5LjcxNDNNNDIgNDJMMzkuNzE0MyAzN0w0MiA0MlpNMjYgNDJMMjguMjg1NyAzN0wyNiA0MlpNMjguMjg1NyAzN0wzNCAyNEwzOS43MTQzIDM3SDI4LjI4NTdaIiBzdHJva2U9IiMzMzMiIHN0cm9rZS13aWR0aD0iNCIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+PHBhdGggZD0iTTE2IDZMMTcgOSIgc3Ryb2tlPSIjMzMzIiBzdHJva2Utd2lkdGg9IjQiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjxwYXRoIGQ9Ik02IDExSDI4IiBzdHJva2U9IiMzMzMiIHN0cm9rZS13aWR0aD0iNCIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+PHBhdGggZD0iTTEwIDE2QzEwIDE2IDExLjc4OTUgMjIuMjYwOSAxNi4yNjMyIDI1LjczOTFDMjAuNzM2OCAyOS4yMTc0IDI4IDMyIDI4IDMyIiBzdHJva2U9IiMzMzMiIHN0cm9rZS13aWR0aD0iNCIgc3Ryb2tlLWxpbmVjYXA9InJvdW5kIiBzdHJva2UtbGluZWpvaW49InJvdW5kIi8+PHBhdGggZD0iTTI0IDExQzI0IDExIDIyLjIxMDUgMTkuMjE3NCAxNy43MzY4IDIzLjc4MjZDMTMuMjYzMiAyOC4zNDc4IDYgMzIgNiAzMiIgc3Ryb2tlPSIjMzMzIiBzdHJva2Utd2lkdGg9IjQiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjwvc3ZnPg=="
                      alt="language">
                  </template>
                  <div class="flex flex-col gap-2">
                    <div v-for="item in locales" :key="item.key" class="cursor-pointer" @click="setLanguage(item.key)">
                      <div
                        :class="'transition duration-300 ease-in-out ' + (item.key === systemStore.language ? 'text-blue-500' : '')">
                        {{ item.value }}
                      </div>
                    </div>
                  </div>
                </el-popover>
              </div>
              <div @click="toggleDark()">
                <div v-if="!isDark" class="cursor-pointer opacity-80">
                  <img
                    src="data:image/svg+xml;charset=utf-8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHdpZHRoPSIyMiIgaGVpZ2h0PSIyMiIgdmlld0JveD0iMCAwIDQ4IDQ4IiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0yNC4wMDMzIDRMMjkuMjczNyA5LjI3MDM4SDM4LjcyOTZWMTguNzI2M0w0NCAyMy45OTY3TDM4LjcyOTYgMjkuMjczN1YzOC43Mjk2SDI5LjI3MzdMMjQuMDAzMyA0NEwxOC43MjY0IDM4LjcyOTZIOS4yNzAzNlYyOS4yNzM3TDQgMjMuOTk2N0w5LjI3MDM2IDE4LjcyNjNWOS4yNzAzOEgxOC43MjY0TDI0LjAwMzMgNFoiIGZpbGw9Im5vbmUiIHN0cm9rZT0iIzMzMyIgc3Ryb2tlLXdpZHRoPSI0IiBzdHJva2UtbWl0ZXJsaW1pdD0iMTAiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjxwYXRoIGQ9Ik0yNyAxN0MyNyAyNSAyMiAyNiAxNyAyNkMxNyAzMCAyMy41IDM0IDI5IDMwQzM0LjUgMjYgMzEgMTcgMjcgMTdaIiBmaWxsPSJub25lIiBzdHJva2U9IiMzMzMiIHN0cm9rZS13aWR0aD0iNCIgc3Ryb2tlLW1pdGVybGltaXQ9IjEwIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz48L3N2Zz4="
                    alt="dark" />
                </div>
                <div v-else class="cursor-pointer opacity-80">
                  <img
                    src="data:image/svg+xml;charset=utf-8;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHdpZHRoPSIyMiIgaGVpZ2h0PSIyMiIgdmlld0JveD0iMCAwIDQ4IDQ4IiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxwYXRoIGQ9Ik0yNC4wMDMzIDRMMjkuMjczNyA5LjI3MDM4SDM4LjcyOTZWMTguNzI2M0w0NCAyMy45OTY3TDM4LjcyOTYgMjkuMjczN1YzOC43Mjk2SDI5LjI3MzdMMjQuMDAzMyA0NEwxOC43MjY0IDM4LjcyOTZIOS4yNzAzNlYyOS4yNzM3TDQgMjMuOTk2N0w5LjI3MDM2IDE4LjcyNjNWOS4yNzAzOEgxOC43MjY0TDI0LjAwMzMgNFoiIGZpbGw9IiNmZmZhMmYiIHN0cm9rZT0iI2ZmZmZmZiIgc3Ryb2tlLXdpZHRoPSI0IiBzdHJva2UtbWl0ZXJsaW1pdD0iMTAiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIgc3Ryb2tlLWxpbmVqb2luPSJyb3VuZCIvPjxwYXRoIGQ9Ik0yNyAxN0MyNyAyNSAyMiAyNiAxNyAyNkMxNyAzMCAyMy41IDM0IDI5IDMwQzM0LjUgMjYgMzEgMTcgMjcgMTdaIiBmaWxsPSIjNmM3MWMwIiBzdHJva2U9IiNGRkYiIHN0cm9rZS13aWR0aD0iNCIgc3Ryb2tlLW1pdGVybGltaXQ9IjEwIiBzdHJva2UtbGluZWNhcD0icm91bmQiIHN0cm9rZS1saW5lam9pbj0icm91bmQiLz48L3N2Zz4="
                    alt="light" />
                </div>
              </div>
            </div>
          </div>
          <!-- user info -->
          <ColorView custom-class="transition duration-300 ease-in-out bg-slate-200 dark:bg-zinc-800">
            <div class="flex gap-2 items-center justify-between" v-if="userInfo">
              <div class="flex gap-2 items-center">
                <div>
                  <img :src="userInfo?.avatar" alt="logo"
                       class="bg-orange-200 dark:bg-orange-300 rounded-full w-10 h-10" />
                </div>
                <div>
                  <div class="font-bold transition duration-300 ease-in-out text-black dark:text-white">
                    {{ userInfo?.nickname }}
                  </div>
                  <div
                    class="text-xs transition duration-300 ease-in-out text-gray-500 dark:text-gray-400 truncate max-w-40">
                    {{ userInfo?.username }}
                  </div>
                </div>
              </div>
              <div class="cursor-pointer">
                <logout theme="outline" size="18" fill="#333" @click="onClickLogout" />
              </div>
            </div>
            <div v-else>
              <router-link to="/login">
                <div class="flex items-center px-2">
                  <avatar />
                  <div class="p-2">{{ $t('signIn') }}</div>
                </div>
              </router-link>
            </div>
          </ColorView>

          <!-- navigation -->
          <NavigationView class="flex-grow" />
        </div>

        <!-- end -->
        <ColorView custom-class="bg-none">
          <div class="text-xs text-gray-500" v-html="systemStore.website.copyright" />
        </ColorView>
      </div>
    </el-aside>
    <el-container class="transition duration-300 ease-in-out bg-opacity-20 bg-indigo-200 dark:bg-slate-800">
      <el-main class="relative">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
:deep(.el-header),
:deep(.el-main),
:deep(.el-aside) {
  padding: 0;
}
</style>