import axios from '@/utils/request'

export namespace User {
  export interface AvailableOAuth {
    oauth: string
    type: string
  }

  export interface UserInfo {
    id: string
    username: string
    nickname: string
    password: string
    avatar: string
    email: string
    phone: string
    createdAt: string
    updatedAt: string
  }


  export const getAvailableOAuth = async (): Promise<AvailableOAuth[]> => {
    return axios.get<AvailableOAuth[]>("/user/oauth/all")
  }

  export const getOAuthURL = async (oauth: string, redirectURL: string = '/'): Promise<string> => {
    return axios.get<string>(`/user/oauth?oauth=${oauth}&redirectURL=${redirectURL}`)
  }

  export const getUserInfo = async (): Promise<UserInfo> => {
    return axios.get<UserInfo>("/user/info")
  }

  export const logout = async () => {
    return axios.get("/user/logout")
  }
}
