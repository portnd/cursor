import { IError } from "./index"

export interface IResponse<T> {
	code: number
	status: boolean
	data: T
	error: IError
}
