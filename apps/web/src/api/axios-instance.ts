import axios, { type AxiosRequestConfig, type AxiosResponse } from 'axios'

declare const API_URL: string
const baseURL = API_URL

// httpOnly Cookie使用時は型ガードは不要
// トークンはCookieで管理され、クライアント側でアクセスしない

const axiosInstance = axios.create({
  baseURL,
  timeout: 30000,
  withCredentials: true, // httpOnly Cookieを含むリクエストを有効化
  headers: {
    'Content-Type': 'application/json',
  },
})

// リクエストインターセプター
axiosInstance.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

// Orval用のカスタムインスタンス関数
export const customAxiosInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const controller = new AbortController()

  const promise = axiosInstance({
    ...config,
    ...options,
    signal: controller.signal,
  }).then((response: AxiosResponse<T>) => response.data)

  // AbortControllerを使用した適切なキャンセル実装
  Object.assign(promise, {
    cancel: () => {
      controller.abort()
    },
  })

  return promise
}

export default axiosInstance
