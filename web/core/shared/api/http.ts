import type { UseFetchOptions } from '#app'

export interface ApiResponse<T = any> {
  data: T
  error: any
  pending: Ref<boolean>
  refresh: () => Promise<void>
}

export const useHttp = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  const request = async <T = any>(
    endpoint: string,
    options: UseFetchOptions<T> = {}
  ): Promise<ApiResponse<T>> => {
    const defaultOptions: UseFetchOptions<T> = {
      baseURL,
      ...options,
      onRequest({ options }) {
        const token = useCookie('auth_token').value
        if (token) {
          options.headers = {
            ...options.headers,
            Authorization: `Bearer ${token}`
          }
        }
      },
      onResponseError({ response }) {
        if (response.status === 401) {
          useCookie('auth_token').value = null
          navigateTo('/login')
        }
      }
    }

    const { data, error, pending, refresh } = await useFetch<T>(endpoint, defaultOptions)

    return {
      data: data.value as T,
      error: error.value,
      pending,
      refresh
    }
  }

  const get = <T = any>(endpoint: string, options: UseFetchOptions<T> = {}) => {
    return request<T>(endpoint, { ...options, method: 'GET' })
  }

  const post = <T = any>(endpoint: string, body: any, options: UseFetchOptions<T> = {}) => {
    return request<T>(endpoint, { ...options, method: 'POST', body })
  }

  const put = <T = any>(endpoint: string, body: any, options: UseFetchOptions<T> = {}) => {
    return request<T>(endpoint, { ...options, method: 'PUT', body })
  }

  const patch = <T = any>(endpoint: string, body: any, options: UseFetchOptions<T> = {}) => {
    return request<T>(endpoint, { ...options, method: 'PATCH', body })
  }

  const del = <T = any>(endpoint: string, options: UseFetchOptions<T> = {}) => {
    return request<T>(endpoint, { ...options, method: 'DELETE' })
  }

  return {
    request,
    get,
    post,
    put,
    patch,
    delete: del
  }
}
