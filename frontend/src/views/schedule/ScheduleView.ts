import axios, { Pager } from '@/utils/request'

export namespace Schedule {
  export interface ScheduleItem {
    createdAt: string
    cronString: string
    desc: string
    enabled: boolean
    executor: string
    id: string
    nextTime: string
    params: string
    status: string
    title: string
    updatedAt: string
  }

  export interface Executor {
    name: string
    displayName: string
  }

  export const pageSchedule = async (page: Number, size: Number): Promise<Pager<ScheduleItem>> => {
    return axios.get<Pager<ScheduleItem>>(`/schedule/page?page=${page}&size=${size}`)
  }

  export const deleteSchedule = async (id: string) => {
    return axios.delete(`/schedule/${id}`)
  }

  export const getExecutors = async (): Promise<Executor[]> => {
    return axios.get<Executor[]>(`/schedule/executors`)
  }

  export const updateSchedule = async (data: ScheduleItem) => {
    return axios.put(`/schedule/${data.id}`, data)
  }
}