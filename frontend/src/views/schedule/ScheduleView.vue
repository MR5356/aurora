<script setup lang="ts">

import TableView, { type TableColumn } from '@/components/TableView.vue'
import { ref, watch } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useI18n } from 'vue-i18n'
import cronstrue from 'cronstrue/i18n'
import moment from 'moment'
import { Schedule } from '@/views/schedule/ScheduleView'
import type { Pager } from '@/utils/request'
import { ElMessageBox } from 'element-plus'

const systemStore = useSystemStore()
const { t } = useI18n()

const data = ref<Schedule.ScheduleItem[]>([])
const search = ref('')
const currentItem = ref<Schedule.ScheduleItem | null>(null)
const showDrawer = ref(false)
const drawerMode = ref<'edit' | 'detail'>('edit')
const loading = ref(false)

const executors = ref<Schedule.Executor[]>([])
const rawData = ref<Pager<Schedule.ScheduleItem>>({
  current: 1,
  size: 25,
  total: 0
})

const formatData = () => {
  data.value = []
  rawData.value.data.forEach((item) => {
    let newItem = JSON.parse(JSON.stringify(item))
    newItem.cronString = cronstrue.toString(newItem.cronString, { locale: systemStore.language })
    newItem.enabled = t('schedule.' + newItem.enabled)
    newItem.createdAt = moment(newItem.createdAt).format('YYYY-MM-DD HH:mm:ss')
    newItem.executor = executors.value.find((item) => item.name === newItem.executor).displayName
    data.value.push(newItem)
  })
}

const listSchedule = async () => {
  loading.value = true
  executors.value = await Schedule.getExecutors()
  rawData.value = await Schedule.pageSchedule(rawData.value.current, rawData.value.size)
  formatData()
  loading.value = false
}

listSchedule()

watch(() => systemStore.language, () => {
  formatData()
}, {
  deep: true
})

const columns = ref<TableColumn[]>([
  {
    field: 'title',
    label: 'schedule.title',
    fixed: true,
    align: 'left',
    width: 150
  },
  {
    field: 'cronString',
    label: 'schedule.cronString',
    align: 'left',
    width: 200
  },
  {
    field: 'executor',
    label: 'schedule.executor',
    align: 'left',
    width: 150
  },
  {
    field: 'enabled',
    label: 'schedule.enabled',
    align: 'left',
    width: 88
  },
  {
    field: 'desc',
    label: 'schedule.desc',
    align: 'left',
    width: 200
  },
  {
    field: 'params',
    label: 'schedule.params',
    align: 'left'
  },
  {
    field: 'createdAt',
    label: 'schedule.createdAt',
    align: 'left',
    width: 200
  }
])
const multipleSelection = ref<any[]>([])
const handleSelectionChange = (val: any[]) => {
  multipleSelection.value = val
  console.log(multipleSelection.value)
}

const handleDelete = async (val: Schedule.ScheduleItem) => {
  ElMessageBox.confirm(t('confirmDelete'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel')
  }).then(async () => {
    await Schedule.deleteSchedule(val.id)
    await listSchedule()
  })
}

const handleDetail = (val: Schedule.ScheduleItem) => {
  showDrawer.value = true
  drawerMode.value = 'detail'
  setCurrentItem(val)
}

const handleEdit = (val: Schedule.ScheduleItem) => {
  showDrawer.value = true
  drawerMode.value = 'edit'
  setCurrentItem(val)
}

const setCurrentItem = (val: Schedule.ScheduleItem) => {
  currentItem.value = rawData.value.data.find((item) => item.id === val.id)
}

const onPageChange = async (e) => {
  rawData.value.current = e
  await listSchedule()
}

const onUpdateSchedule = async () => {
  loading.value = true
  await Schedule.updateSchedule(currentItem.value)
  loading.value = false
  showDrawer.value = false
  await listSchedule()
}

const onSearch = () => {
  console.log(search.value)
}
</script>

<template>
  <div class="flex gap-0 h-[100vh]">
    <div
      class="w-full h-full flex flex-col rounded-none shadow-2xl shadow-fuchsia-50 dark:shadow-slate-900 overflow-hidden">
      <TableView
        :data="data"
        :columns="columns"
        selection
        :select-change="handleSelectionChange"
        :handler-delete="handleDelete"
        :handler-edit="handleEdit"
        :handler-detail="handleDetail"
        v-model:search="search"
      >
      </TableView>
      <el-pagination
        class="w-full bg-white dark:bg-slate-900"
        layout="prev, pager, next"
        :total="rawData.total"
        :page-size="rawData.size"
        @change="onPageChange"
      />
    </div>
    <transition name="slide-fade" :appear="true">
      <div v-if="showDrawer" class="min-w-[300px] w-1/3 h-full">
        <div
          class="h-full p-4 transition duration-300 ease-in-out text-slate-600 dark:text-slate-400 bg-opacity-40 bg-white dark:bg-slate-900 relative">
          <div class="font-bold select-none">
            <span v-if="drawerMode==='edit'">{{ $t('schedule.editTitle') }}</span>
            <span v-if="drawerMode==='detail'">{{ $t('schedule.detailTitle') }}</span>
            <el-form
              label-position="top"
              label-width="auto"
              :model="currentItem"
            >
              <el-form-item :label="$t('schedule.title')">
                <el-input v-model="currentItem.title" :disabled="drawerMode==='detail'" />
              </el-form-item>
              <el-form-item :label="$t('schedule.desc')">
                <el-input v-model="currentItem.desc" :disabled="drawerMode==='detail'" />
              </el-form-item>
              <el-form-item :label="$t('schedule.cronString')">
                <el-input v-model="currentItem.cronString" :disabled="drawerMode==='detail'" />
              </el-form-item>
              <el-form-item :label="$t('schedule.executor')">
                <el-select
                  v-model="currentItem.executor"
                  :disabled="drawerMode==='detail'"
                  size="large"
                >
                  <el-option
                    v-for="item in executors"
                    :key="item.name"
                    :label="item.displayName"
                    :value="item.name"
                  />
                </el-select>
              </el-form-item>
              <el-form-item :label="$t('schedule.enabled')">
                <el-switch v-model="currentItem.enabled" :disabled="drawerMode==='detail'" />
              </el-form-item>
              <el-form-item :label="$t('schedule.params')">
                <el-input type="textarea" :rows="5" v-model="currentItem.params" :disabled="drawerMode==='detail'" />
              </el-form-item>
            </el-form>
          </div>
          <div class="absolute bottom-4 right-4">
            <el-button size="small" @click="showDrawer=false" :disabled="loading">{{ $t('close') }}</el-button>
            <el-button size="small" type="primary" v-if="drawerMode=='edit'" @click="onUpdateSchedule" :disabled="loading">{{ $t('save') }}</el-button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
</style>