package constants

const (
	//User error
	INVALID_ID_PASSWORD     string = "ชื่อผู้ใช้/รหัสผ่านไม่ถูกต้อง"
	USER_NOT_FOUND          string = "ไม่พบผู้ใช้งาน"
	USER_BLOCKED            string = "ผู้ใช้งานถูกระงับการเข้าใช้"
	INVALID_USER_PERMISSION string = "ผู้ใช้งานไม่มีสิทธิ์ในการใช้งานในส่วนนี้"

	//Data error
	EDITED_BY_ANOTHER_USER                      string = "ขณะนี้มีคนอื่นกำลังแก้ไขข้อมูลนี้อยู่"
	DATA_NOT_MATCH_USER                         string = "ข้อมูลนี้ไม่ได้ถูกเชื่อมโยงกับผู้ใช้"
	INVALID_CATEGORY                            string = "ประเภทไม่ถูกต้อง"
	INVALID_STATUS                              string = "สถานะไม่ถูกต้อง"
	INVALID_CONDITION_TYPE                      string = "ประเภทสภาพทางไม่ถูกต้อง"
	DATA_WAITING_APPROVAL                       string = "ขณะนี้ข้อมูลอยู่ในสถานะรอการอนุมัติ"
	DATA_ALREADY_CONFIRMED                      string = "ข้อมูลเคยได้รับการยืนยันไปแล้ว"
	INVALID_GEOM_RANGE                          string = "ข้อมูลไม่อยู่ในช่วงระยะของสายทางในระบบ"
	INVALID_ROAD_CONDITION_GEOM_RANGE           string = "โปรดตรวจสอบไฟล์ข้อมูลสำรวจ IRI/RUT/MPD เนื่องจากข้อมูลไม่อยู่ในช่วงกม. ของสายทาง"
	INVALID_ROAD_RETRO_REFLECTIVITY_GEOM_RANGE  string = "โปรดตรวจสอบไฟล์ข้อมูลสำรวจเนื่องจากข้อมูลไม่อยู่ในช่วงกม. ของสายทาง"
	INVALID_ROAD_RETRO_REFLECTIVITY_LINE        string = "โปรดระบุเส้นจราจรไม่เกิน "
	INVALID_ROAD                                string = "หมายเลขสายทางปัจจุบันไม่ตรงกับหมายเลขสายทางก่อนหน้า โปรดตรวจสอบ"
	INVALID_RETRO_REFLECTIVITY_LINE             string = "เส้นจราจรข้อมูลแถบสะท้อนแสงปัจจุบันไม่ตรงกับเส้นจราจรก่อนหน้า โปรดตรวจสอบ"
	INVALID_ROAD_CONDITION_LANE                 string = "เลนส์ข้อมูลสภาพทางปัจจุบันไม่ตรงกับข้อมูลเลนส์ก่อนหน้า โปรดตรวจสอบ"
	INVALID_ROAD_CONDITION_DATA                 string = "ข้อมูลสภาพทางไม่ถูกต้อง กรุณาตรวจสอบข้อมูล ณ "
	INVALID_ROAD_CONDITION_ROW_DATA             string = "ข้อมูลสภาพทางไม่ถูกต้อง จำนวนแถวของข้อมูลไม่ตรงกับสายทาง"
	INVALID_ROAD_CONDITION_VALUE                string = "โปรดตรวจสอบบรรทัดที่ _ เนื่องจากข้อมูลไม่ถูกต้อง"
	INVALID_ROAD_RETRO_REFLECTIVITY_COLOR       string = "โปรดตรวจสอบบรรทัดที่ _ เนื่องจากข้อมูลสีเส้นจราจรไม่ตรงกับข้อมูลในระบบ"
	INVALID_ROAD_RETRO_REFLECTIVITY_STRIPE_TYPE string = "โปรดตรวจสอบบรรทัดที่ _ เนื่องจากข้อมูลประเภทเส้นจราจรไม่ตรงกับข้อมูลในระบบ"
	OVERLAPPING_RANGE                           string = "ช่วงระยะซ้อนทับกัน"
	INVALID_HEX_FORMAT                          string = "โค้ดสีไม่ถูกต้อง"
	INVALID_GEOM_TYPE                           string = "ประเภทพิกัดไม่ถูกต้อง"
	INVALID_SURFACE_ID                          string = "โปรดระบุชนิดผิวทางอย่างน้อย 1 ช่องจราจร ในช่วงกม."
	//File error
	FAILED_TO_SAVE_FILE   string = "ระบบไม่สามารถจัดเก็บไฟล์ได้"
	FAILED_TO_CREATE_FILE string = "ระบบไม่สามารถสร้างไฟล์ได้"
	UNSUPPORTED_FILE_TYPE string = "ระบบไม่รองรับชนิดไฟล์ที่อัปโหลด"

	INVALID_POLY_LINE string = "ข้อมูล PolyLine ไม่ถูกต้องกรุณาตรวจสอบ"

	FAILED_TO_UPDATE_ROAD            string = "ระบบไม่สามารถปรับปรุงข้อมูลสายทางได้ กรุณาตรวจสอบ"
	FAILED_TO_UPDATE_ROAD_HAVE_CHILD string = "ไม่สามารถบันทึกข้อมูลได้ เนื่องจากสายทางนี้มีข้อมูล ServiceRoad, Ramp, Uturn อยู่แล้ว"
	FAILED_TO_CREATE_ROAD            string = "ระบบไม่สามารถสร้างข้อมูลสายทางได้ กรุณาตรวจสอบ"

	FAILED_TO_DELETE_ROAD string = "ระบบไม่สามารถลบข้อมูลสายทางได้ กรุณาตรวจสอบ"

	INVALID_ROAD_RANGE string = "โปรดตรวจสอบไฟล์ (Center Line) หรือ (Center Lane) เนื่องจากข้อมูล กม.ไม่อยู่ในช่วง กม. ของสายทาง"

	FAILED_TO_CREATE_DIR                string = "ระบบไม่สามารถสร้างโฟลเดอร์เพื่อจัดเก็บข้อมูลได้ กรุณาตรวจสอบ"
	INVALID_CONVERT_LINE_STRING_TO_GEOM string = "ระบบไม่สามารถแปลง LineString เป็น Geom ได้ กรุณาตรวจสอบ"
	HAS_MANY_SHAPE_FILE_IN_ZIP          string = "ระบบตรวจพบ Shape File มากกว่า 1 ไฟล์ใน zip กรุณาตรวจสอบ"
	FAILED_TO_READ_ENV_FILE             string = "ระบบไม่สามารถอ่านไฟล์ข้อมูลจากไฟล์ .env ได้"
	FAILED_TO_UPLOAD_CENTER_LINE_FILE   string = "ไฟล์ข้อมูลสายทางกึ่งกลางถนน (Center Line) ไม่ถูกต้อง กรุณาตรวจสอบ"
	FAILED_TO_UPLOAD_CENTER_LANE_FILE   string = "ไฟล์ข้อมูลสายทางในแต่ละช่องจราจร (Center Lane) ไม่ถูกต้อง กรุณาตรวจสอบ"

	CENTER_LANE_CODE_NOT_MATCH string = "รหัสสายทาง ของข้อมูลสายทางในแต่ละช่องจราจร (Center Lane) ไม่ตรงกัน กรุณาตรวจสอบ"
	CENTER_LINE_CODE_NOT_MATCH string = "รหัสสายทาง ของข้อมูลสายทางกึ่งกลางถนน (Center Line) ไม่ตรงกัน กรุณาตรวจสอบ"

	NOT_IN_RANGE string = "ช่วงกม.ไม่อยู่ในสายทาง"

	//RoadGroup Color
	YELLOW_1        = "#FDB833"
	YELLOW_2        = "#FDC65C"
	YELLOW_3        = "#FED485"
	ORANGE_1        = "#FF8929"
	ORANGE_2        = "#FFA154"
	ORANGE_3        = "#FFB87F"
	BRIGHT_YELLOW_1 = "#FFE433"
	BRIGHT_YELLOW_2 = "#FFE95C"
	BRIGHT_YELLOW_3 = "#FFEF85"
	RED_ORANGE_1    = "#FF5638"
	RED_ORANGE_2    = "#FF7860"
	RED_ORANGE_3    = "#FF9A88"
)
