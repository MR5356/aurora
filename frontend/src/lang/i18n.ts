import { createI18n } from 'vue-i18n'
import { ref } from 'vue'
import en from './locales/en.json'
import zh_CN from './locales/zh_CN.json'

export interface Locale {
  key: string
  value: string
  component: any
}

export const locales = ref<Locale[]>([
  {
    key: 'en',
    value: 'English',
    component: en
  },
  {
    key: 'zh_CN',
    value: '简体中文',
    component: zh_CN
  }
])

const localeMap = locales.value.reduce((acc, curr) => {
  acc[curr.key] = curr.component
  return acc
}, {} as {[key: string]: any})

const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: 'en',
  messages: localeMap
})

export default i18n