import { ITrafficParams } from "./ModelTrafficParameterRequest"
import { ITrafficRoadGroupData, ITrafficVehicleData, IGetTrafficData } from "./ModelTrafficParameterModel"
import { IResponse } from "~/core/shared/http"
import { HttpService } from "~/core/shared/http/HttpService"

export class TrafficParameterService extends HttpService {
	public constructor() {
		super("/")
	}

	public postAadtParameter(param: ITrafficParams): Promise<IResponse<ITrafficParams>> {
		return this.http
			.post("settings/aadt/parameter", param)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAadtParameterRoadGroupVolume(): Promise<IResponse<ITrafficRoadGroupData[]>> {
		return this.http
			.get("settings/aadt/parameter/road_group_with_volume_aadt")
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAadtParameter(roadGroupID: number): Promise<IResponse<IGetTrafficData>> {
		return this.http
			.get(`settings/aadt/parameter/${roadGroupID}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}

	public getAadtParametersVehicleType(roadGroupID: number): Promise<IResponse<ITrafficVehicleData[]>> {
		return this.http
			.get(`ref/aadt_parameter_vehicle_type/${roadGroupID}`)
			.then(this.handleResponse.bind(this))
			.catch(this.handleError.bind(this))
	}
}
