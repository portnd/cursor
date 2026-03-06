import { IVehicleType, IVehicleTypeRequest } from "."
import { IResponse } from "~~/core/shared/http"
import { HttpService } from "~~/core/shared/http/HttpService"

export class VehicleTypeService extends HttpService {
	public constructor() {
		super("/")
	}

	public post(params: IVehicleTypeRequest): Promise<IResponse<{}>> {
		return this.http
			.post(`/settings/aadt/percentage_vehicle_type`, params)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public get(id: number): Promise<IResponse<IVehicleType>> {
		return this.http
			.get(`/settings/aadt/percentage_vehicle_type/${id}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
