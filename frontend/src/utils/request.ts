import axios from 'axios'
import { ElMessage } from 'element-plus'
import type { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'
import router from '@/router'
import i18n from '@/lang/i18n'

export declare interface Pager<T> {
  current: number
  size: number
  total: number
  data?: T[]
}

enum RequestEnums {
  TIMEOUT = 60000,
  SUCCESS = '00000', // 请求成功
  NotLogin = 'B1001',  // 未登录
  BASEURL = '/api/v1'
}
const config = {
  // 设置超时时间
  timeout: RequestEnums.TIMEOUT as number,
  // 跨域时候允许携带凭证
  withCredentials: false,
  baseURL: RequestEnums.BASEURL as string,
}

class RequestHttp {
  // 定义成员变量并指定类型
  service: AxiosInstance;
  public constructor(config: AxiosRequestConfig) {
    // 实例化axios
    this.service = axios.create(config);

    /**
     * 请求拦截器
     * 客户端发送请求 -> [请求拦截器] -> 服务器
     * token校验(JWT) : 接受服务器返回的token,存储到vuex/pinia/本地储存当中
     */
    this.service.interceptors.request.use(
      (config: AxiosRequestConfig) => {
        // const token = localStorage.getItem('token') || '';
        return {
          ...config,
          headers: {
            // 'x-access-token': token, // 请求头中携带token信息
          }
        }
      },
      (error: AxiosError) => {
        // 请求报错
        Promise.reject(error)
      }
    )

    /**
     * 响应拦截器
     * 服务器换返回信息 -> [拦截统一处理] -> 客户端JS获取到信息
     */
    this.service.interceptors.response.use(
      (response: AxiosResponse) => {
        const {data, config} = response; // 解构
        if (data.code === RequestEnums.NotLogin) {
          ElMessage.warning(i18n.global.t('needLogin'))
          router.replace({
            path: `/login`,
            query: {
              redirectURL: router.currentRoute.value.fullPath
            }
          }).then(() => {})
          return Promise.reject(data);
        }
        // 全局错误信息拦截（防止下载文件得时候返回数据流，没有code，直接报错）
        if (data.code && data.code !== RequestEnums.SUCCESS) {
          ElMessage.error(data.message)
          return Promise.reject(data)
        }
        return data.data;
      },
      (error: AxiosError) => {
        const {response} = error;
        if (response) {
          this.handleCode(response.status)
        }
        if (!window.navigator.onLine) {
          ElMessage.error('网络连接失败')
        }
      }
    )
  }
  handleCode(code: number):void {
    switch(code) {
      case 401:
        ElMessage.error('登录失效，请重新登录')
        break;
      default:
        ElMessage.error('数据请求失败')
        break;
    }
  }

  // 常用方法封装
  get<T>(url: string, params?: object): Promise<T> {
    return this.service.get(url, {params});
  }
  post<T>(url: string, params?: object): Promise<T> {
    return this.service.post(url, params);
  }
  put<T>(url: string, params?: object): Promise<T> {
    return this.service.put(url, params);
  }
  delete<T>(url: string, params?: object): Promise<T> {
    return this.service.delete(url, {params});
  }
}

// 导出一个实例对象
export default new RequestHttp(config);