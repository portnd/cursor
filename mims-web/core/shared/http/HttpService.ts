import type { AxiosError, AxiosInstance, AxiosResponse } from "axios"
import axios from "axios"
import { IResponse } from "~~/core/shared/http"
import { useAuthStore } from "~~/core/modules/auth/authentication/store"

export function isAxiosError(value: any): value is AxiosError {
	return typeof value?.response === "object"
}

export abstract class HttpService {
	protected readonly http: AxiosInstance

	protected constructor(
		protected readonly path?: string,
		protected readonly baseURL: string = useRuntimeConfig().public?.apiUrl || ""
	) {
		if (path) {
			baseURL += path
		}

		this.http = axios.create({
			baseURL,
		})

		this.http.defaults.headers.common.Accept = "application/json;charset=UTF-8"
		this.http.defaults.headers.common["Content-Type"] = "application/json;charset=UTF-8"

		// JWT Token
		const authStore = useAuthStore()
		if (authStore.isLoggedIn) {
			this.http.defaults.headers.common.Authorization = `Bearer ${authStore.accessToken}`
		}
	}

	protected createParams(record: Record<string, any>): URLSearchParams {
		const params: URLSearchParams = new URLSearchParams()
		for (const key in record) {
			if (Object.prototype.hasOwnProperty.call(record, key)) {
				const value: any = record[key]
				if (value !== null && value !== undefined) {
					params.append(key, value)
				} else {
					console.debug(`Param key '${key}' was null or undefined and will be ignored`)
				}
			}
		}
		return params
	}

	protected handleResponse<T>(response: AxiosResponse<T>): T {
		const result: IResponse<T> = response.data as IResponse<T>
		return {
			code: response.status,
			status: result.status,
			data: result.data,
		} as T
	}

	protected handleError<T>(error: AxiosError): T {
		if (error instanceof Error) {
			if (isAxiosError(error)) {
				if (error.response) {
					const result: IResponse<T> = error.response.data as IResponse<T>
					return {
						code: error.response.status,
						status: result.status ? result.status : false,
						error: result.error,
					} as T
				} else if (error.request) {
					console.log(error.request)
					throw new Error(error as any)
				}
			} else {
				console.log("error", error)
				throw new Error(error)
			}
		}
		throw new Error(error as any)
	}
}
