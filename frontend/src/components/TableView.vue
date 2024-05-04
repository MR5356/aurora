<script setup lang="ts">

import type { PropType, VNode } from 'vue'

export interface TableColumn {
  field: string
  label: string
  width?: string | number
  fixed?: boolean
  align?: 'left' | 'center' | 'right'
  formatter?: (row: any, column: any, cellValue: any, index: number) => VNode | string
}

defineProps({
  data: {
    type: Array<any>,
    require: true
  },
  columns: {
    type: Array<TableColumn>
  },
  selection: {
    type: Boolean,
    require: false,
    default: false
  },
  loading: {
    type: Boolean,
    require: false,
    default: false
  },
  selectChange: Function as PropType<((val: any[]) => void)>,
  handlerEdit: Function as PropType<((val: any) => void)>,
  handlerDelete: Function as PropType<((val: any) => void)>,
  handlerDetail: Function as PropType<((val: any) => void)>,
  onSearch: Function as PropType<any>
})

const search = defineModel('search', { type: String, required: false })
</script>

<template>
  <el-table
    ref="multipleTableRef"
    :data="data"
    v-loading="loading"
    :empty-text="$t('table.emptyText')"
    style="width: 100%; height: 100%"
    @selection-change="selectChange"
    :scrollbar-always-on="true"
  >
    <template #empty>
      <el-empty :description="$t('table.emptyText')" />
    </template>
    <el-table-column v-if="selection" type="selection" width="55" />
    <el-table-column v-for="column in columns" :key="column.field" :property="column.field" :label="$t(column.label)"
                     :width="column.width" :show-overflow-tooltip="true" :align="column.align"
                     :formatter="column.formatter" />
    <el-table-column
      v-if="handlerDelete || handlerDetail || handlerEdit"
      :label="$t('table.operations')"
      :width="((handlerEdit ? 1 : 0) + (handlerDetail ? 1 : 0) + (handlerDelete ? 1 : 0)) * 86"
    >
      <template #header>
        <el-input
          v-if="onSearch"
          v-model="search"
          size="small"
          :placeholder="$t('table.searchPlaceholder')"
          clearable
          @keyup.enter="onSearch" />
      </template>
      <template #default="scope">
        <el-button size="small" text bg v-if="handlerDetail" @click="handlerDetail(scope.row)">
          {{ $t('table.detail') }}
        </el-button>
        <el-button size="small" text bg v-if="handlerEdit" @click="handlerEdit(scope.row)">
          {{ $t('table.edit') }}
        </el-button>
        <el-button size="small" text bg type="danger" v-if="handlerDelete" @click="handlerDelete(scope.row)">
          {{ $t('table.delete') }}
        </el-button>
      </template>
    </el-table-column>
  </el-table>
  <slot />
</template>

<style lang="scss">
//.el-table {
//  background-color: transparent;
//}

//.el-table .custom-table {
//  background-color: rgba(255, 255, 255, 0.95);
//  border: 0;
//
//  td:first-child {
//    border-radius: 8px 0 0 8px;
//  }
//
//  td:last-child {
//    border-radius: 0 8px 8px 0;
//  }
//}
</style>