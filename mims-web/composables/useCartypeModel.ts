export interface ICarTypeModel {
	[key: string]: any
	//  name: string
	//  nameValidate: string
	//  popOver: string
}

const useCarTypeModel = (): ICarTypeModel => {
	return {
		car_less_than_equal_seven: "รถยนต์ส่วนบุคคลไม่เกิน 7 ที่นั่ง",
		car_over_than_seven: "รถยนต์ส่วนบุคคลมากกว่า 7 ที่นั่ง",
		light_bus: "รถโดยสารขนาดเล็ก",
		light_truck: "รถบรรทุกเล็ก 4 ล้อ",
		medium_bus: "รถโดยสารขนาดกลาง",
		medium_truck: "รถบรรทุกขนาดกลาง 6-10 ล้อ",
		heavy_bus: "รถโดยสารขนาดใหญ่",
		heavy_truck: "รถบรรทุกขนาดใหญ่มากกว่า 10 ล้อ",
		full_trailor: "รถพ่วง",
		semi_trailor: "รถกึ่งพ่วง",
	}
	//  [
	// 	 {
	// 	 	name: "รถ 4 ล้อ",
	//      type:'four_wheel',
	// 	 	item: [
	// 			{
	// 				 name: "Car < 7",
	// 				 nameValidate: "car_less_seven",
	// 				popOver: "รถยนต์ส่วนบุคคลไม่เกิน 7 ที่นั่ง",
	// 			},
	// 			{
	// 				 name: "Car > 7",
	// 				 nameValidate: "car_more_seven",
	// 				popOver: "รถยนต์ส่วนบุคคลมากกว่า 7 ที่นั่ง",
	// 			},
	// 			{
	// 				 name: "Light Bus",
	// 				 nameValidate: "light_bus",
	// 				popOver: "รถโดยสารขนาดเล็ก",
	// 			},
	// 			{
	// 				 name: "Light Truck",
	// 				 nameValidate: "light_truck",
	// 				popOver: "รถบรรทุกเล็กเกิน 4 ล้อ",
	// 			},
	// 		 ],
	// 	 },
	// 	 {
	// 	 	name: "รถ 6-10 ล้อ",
	//      type:'six_wheel',
	// 	 	item: [
	// 			{
	// 				 name: "Medium Bus",
	// 				 nameValidate: "medium_bus",
	// 				popOver: "รถโดยสารขนาดกลาง",
	// 			},
	// 			{
	// 				 name: "Medium Truck",
	// 				 nameValidate: "medium_truck",
	// 				popOver: "รถบรรทุกขนาดกลาง 6-10 ล้อ",
	// 			},
	// 		 ],
	// 	 },
	// 	 {
	// 	 	name: "รถ > 10 ล้อ",
	//      type:'Wheel_wheel',
	// 	 	item: [
	// 			{
	// 				 name: "Heavy Bus",
	// 				 nameValidate: "heavy_bus",
	// 				popOver: "รถโดยสารขนาดใหญ่",
	// 			},
	// 			{
	// 				 name: "Heavy Truck",
	// 				 nameValidate: "heavy_truck",
	// 				popOver: "รถบรรทุกขนาดใหญ่มากกว่า 10 ล้อ",
	// 			},
	// 			{
	// 				name: "Full Trailer",
	// 				 nameValidate: "full_trailer",
	// 				popOver: "รถพ่วง",
	// 			},
	// 			{
	// 				name: "Semi-Trailer",
	// 				 nameValidate: "semi_trailer",
	// 				popOver: "รถกึ่งพ่วง",
	// 			},
	// 		 ],
	// 	 },
	//  ]
}

export { useCarTypeModel }
