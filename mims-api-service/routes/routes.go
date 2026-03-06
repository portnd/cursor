package routes

import (
	"os"
	"strings"

	"github.com/gin-contrib/gzip"
	// _ "gitlab.com/mims-api-service/docs"
	"gitlab.com/mims-api-service/databases"
	"gitlab.com/mims-api-service/logs"
	middlewares "gitlab.com/mims-api-service/middlewares"
	servicesDB "gitlab.com/mims-api-service/services/database"
	authRepo "gitlab.com/mims-api-service/src/auth/repositories"
	dashboardRepo "gitlab.com/mims-api-service/src/dashboard/repositories"
	discordRepo "gitlab.com/mims-api-service/src/discord/repositories"
	hrisRepo "gitlab.com/mims-api-service/src/hris/repositories"
	hsmsRepo "gitlab.com/mims-api-service/src/hsms/repositories"
	maintenanceRepo "gitlab.com/mims-api-service/src/maintenance/repositories"
	maintenanceAnalysisRepo "gitlab.com/mims-api-service/src/maintenanceAnalysis/repositories"
	menuRepo "gitlab.com/mims-api-service/src/menu/repositories"
	refRepo "gitlab.com/mims-api-service/src/ref/repositories"
	reportRepo "gitlab.com/mims-api-service/src/report/repositories"
	roadRepo "gitlab.com/mims-api-service/src/road/repositories"
	roadAssetRepo "gitlab.com/mims-api-service/src/roadAsset/repositories"
	roadRepoCondition "gitlab.com/mims-api-service/src/roadCondition/repositories"
	roadDamageRepo "gitlab.com/mims-api-service/src/roadDamage/repositories"
	roadGroupRepo "gitlab.com/mims-api-service/src/roadGroup/repositories"
	roadRetroReflectivity "gitlab.com/mims-api-service/src/roadRetroReflectivity/repositories"
	roadSectionRepo "gitlab.com/mims-api-service/src/roadSection/repositories"
	roadSurface "gitlab.com/mims-api-service/src/roadSurface/repositories"
	roleRepo "gitlab.com/mims-api-service/src/role/repositories"
	settingRepo "gitlab.com/mims-api-service/src/setting/repositories"
	userRepo "gitlab.com/mims-api-service/src/user/repositories"
	volumeAadtRepo "gitlab.com/mims-api-service/src/volumeAadt/repositories"
	volumeAccidentRepo "gitlab.com/mims-api-service/src/volumeAccident/repositories"

	// NewRepositoryHandler
	// ""gitlab.com/mims-api-service/src/volumeAadt/usecases

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// _ "gitlab.com/mims-api-service/docs"
)

type usecase struct {
	servicesDB servicesDB.ServicesDatabaseDomain
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(func(c *gin.Context) {
		// add header Access-Control-Allow-Origin
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	})

	r.GET("/public/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		if strings.HasSuffix(filepath, ".png") {
			c.Header("Content-Type", "image/png")
		} else if strings.HasSuffix(filepath, ".jpg") {
			c.Header("Content-Type", "image/jpg")
		} else if strings.HasSuffix(filepath, ".pdf") {
			c.Header("Content-Type", "application/pdf")
		}
		c.File("public" + filepath)
	})

	r.GET("/images/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		if strings.HasSuffix(filepath, ".png") {
			c.Header("Content-Type", "image/png")
		} else if strings.HasSuffix(filepath, ".jpg") {
			c.Header("Content-Type", "image/jpg")
		} else if strings.HasSuffix(filepath, ".pdf") {
			c.Header("Content-Type", "application/pdf")
		}
		c.File("images" + filepath)
	})

	r.GET("/storages/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		if strings.HasSuffix(filepath, ".csv") {
			c.Header("Content-Type", "text/csv")
		} else if strings.HasSuffix(filepath, ".xlsx") {
			c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		} else if strings.HasSuffix(filepath, ".pdf") {
			c.Header("Content-Type", "application/pdf")
			c.Header("Content-Disposition", "inline")
		} else if strings.HasSuffix(filepath, ".html") {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Header("Content-Disposition", "inline")
		} else if strings.HasSuffix(filepath, ".zip") {
			c.Header("Content-Type", "application/zip")
		} else if strings.HasSuffix(filepath, ".png") {
			c.Header("Content-Type", "image/png")
		} else if strings.HasSuffix(filepath, ".jpg") {
			c.Header("Content-Type", "image/jpeg")
		} else if strings.HasSuffix(filepath, ".jpeg") {
			c.Header("Content-Type", "image/jpeg")
		}

		localPath := "storages" + filepath
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			// File not found locally — redirect to production server if configured
			fallbackIP := os.Getenv("STORAGE_FALLBACK_IP")
			if fallbackIP != "" {
				c.Redirect(302, fallbackIP+"/storages"+filepath)
				return
			}
		}
		c.File(localPath)
	})
	r.Use(logs.LoggingMiddleware)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// middlewares.AuthorizeJWT()
	v1 := selectAPIPath(r, os.Getenv("ENV"))
	{

		menu := v1.Group("/menu", middlewares.AuthorizeJWT())
		{
			menuHandler := menuRepo.NewMenuRepositoryHandler(databases.DB)
			// menu.GET("", middlewares.AuthorizeJWT(), menuHandler.GetMenu)
			menu.GET("", menuHandler.GetMenu)

		}

		roleCtrl := v1.Group("/roles", middlewares.AuthorizeJWT(), middlewares.RolePermission())
		{
			handler := roleRepo.NewRoleRepositoryHandler(databases.DB)
			roleCtrl.GET("", handler.GetRole)
			roleCtrl.GET(":id", handler.GetRoleById)
			roleCtrl.POST("", handler.CreateRole)
			roleCtrl.PUT(":id", handler.UpdateRole)
			roleCtrl.DELETE(":id", handler.DeleteRole)
			roleCtrl.GET("access_control", handler.GetRoleAccessControlAll)
		}

		auth := v1.Group("/auth")
		{
			handler := authRepo.NewAuthRepositoryHandler(databases.DB)
			auth.GET("logout", handler.Logout)
			auth.GET("check_reset_password_token/:reset_password_token", handler.CheckResetPasswordToken)
			auth.GET("verify_email/:verify_email_token", handler.VerifyEmail)
			auth.POST("login", handler.Login)
			auth.POST("refresh_token", handler.RefreshToken)
			auth.POST("forgot_password", handler.ForgotPassword)
			auth.POST("reset_password", handler.ResetPassword)

			auth.Use(middlewares.AuthorizeJWT())
			auth.POST("resend_verify_email", handler.ResendVerifyEmail)
			auth.GET("test", handler.Test)
		}

		user := v1.Group("/users", middlewares.AuthorizeJWT(), middlewares.UserPermission())
		{
			handler := userRepo.NewUserRepositoryHandler(databases.DB)
			user.GET("", handler.User)
			user.GET(":id", handler.GetUserById)
			user.POST("", handler.CreateUser)
			user.PUT(":id", handler.UpdateUserById)
			user.DELETE(":id", handler.DeleteUserById)
		}

		userInfo := v1.Group("user_info", middlewares.AuthorizeJWT())
		{
			handler := userRepo.NewUserRepositoryHandler(databases.DB)
			userInfo.GET("", handler.GetUserInfo)
			userInfo.PUT("", handler.UpdateUserInfo)
			userInfo.PUT("change_password", handler.UpdatePasswordUserInfo)
		}

		initData := v1.Group("/initdata")
		initData.Use(middlewares.AuthorizeJWT())
		{
			handler := refRepo.NewRefRepositoryHandler(databases.DB)
			initData.GET("", handler.InitData)
		}

		ref := v1.Group("/ref")
		ref.Use(middlewares.AuthorizeJWT())
		{
			handler := refRepo.NewRefRepositoryHandler(databases.DB)
			ref.GET("asset", handler.GetRefAsset)
			ref.GET("data_status", handler.GetRefDataStatus)
			ref.GET("department", handler.GetRefDepartment)
			ref.GET("direction", handler.GetRefDirection)
			ref.GET("grade", handler.GetRefGrade)
			ref.GET("material_base", handler.GetRefMaterialBase)
			ref.GET("material_subbase", handler.GetRefMaterialSubbase)
			ref.GET("material_subgrade", handler.GetRefMaterialSubgrade)
			ref.GET("owner", handler.GetRefOwner)
			ref.GET("road_type", handler.GetRefRoadType)
			ref.GET("surface", handler.GetRefSurface)
			ref.GET("table_list", handler.GetRefTableList)
			ref.GET("color_list", handler.GetRefColorList)
			ref.GET("structure_surface", handler.GetRefStructureSurface)
			ref.GET("criteria", handler.GetRefCriteriaType)
			ref.GET("aadt_parameter_vehicle_type/:road_group_id", handler.GetParameterVehicleType)
			ref.GET("road_user_cost/acc", handler.GetRoadUserCostAcc)
			ref.GET("road_user_cost/ruc", handler.GetRoadUserCostRuc)
			ref.GET("districts", handler.GetRefDistrictsList)
			/////////
			ref.GET("maintenance_analysis_strategic", handler.GetMaintenanceAnalysisStrategic)
		}

		setting := v1.Group("/settings")
		setting.Use(middlewares.AuthorizeJWT())
		{
			handler := settingRepo.NewSettingRepositoryHandler(databases.DB, databases.MongoDb)
			setting.GET("asset_groups", middlewares.AssetGroupPermission(), handler.GetAssetGroups)
			setting.GET("asset_groups/:id", middlewares.AssetGroupPermission(), handler.GetAssetGroupByID)
			setting.POST("asset_groups", middlewares.AssetGroupPermission(), handler.CreateAssetGroup)
			setting.PUT("asset_groups/:id", middlewares.AssetGroupPermission(), handler.UpdateAssetGroupByID)
			setting.DELETE("asset_groups/:id", middlewares.AssetGroupPermission(), handler.DeleteAssetGroupByID)

			setting.GET("owners", middlewares.RoadConditionCriteriaPermission(), handler.GetOwners)
			setting.GET("owners/:id", middlewares.RoadConditionCriteriaPermission(), handler.GetOwnersByID)
			setting.POST("owners", middlewares.RoadConditionCriteriaPermission(), handler.CreateOwner)
			setting.PUT("owners/:id", middlewares.RoadConditionCriteriaPermission(), handler.UpdateOwnerByID)
			setting.DELETE("owners/:id", middlewares.RoadConditionCriteriaPermission(), handler.DeleteOwnerByID)

			setting.GET("owners_road_line", middlewares.RoadG7CriteriaPermission(), handler.GetOwnersRoadLine)
			setting.GET("owners_road_line/:id", middlewares.RoadG7CriteriaPermission(), handler.GetOwnersRoadLineByID)
			setting.POST("owners_road_line", middlewares.RoadG7CriteriaPermission(), handler.CreateOwnerRoadLine)
			setting.PUT("owners_road_line/:id", middlewares.RoadG7CriteriaPermission(), handler.UpdateOwnerRoadLineByID)
			setting.DELETE("owners_road_line/:id", middlewares.RoadG7CriteriaPermission(), handler.DeleteOwnerRoadLineByID)

			setting.GET("signs", middlewares.SignPermission(), handler.GetSigns)
			setting.GET("signs/:id", middlewares.SignPermission(), handler.GetSignByID)
			setting.POST("signs", middlewares.SignPermission(), handler.CreateSign)
			setting.PUT("signs/:id", middlewares.SignPermission(), handler.UpdateSignByID)
			setting.DELETE("signs/:id", middlewares.SignPermission(), handler.DeleteSignByID)

			setting.GET("asset_tables", middlewares.AssetTablePermission(), handler.GetAssetTables)
			setting.GET("asset_tables/:id", middlewares.AssetTablePermission(), handler.GetAssetTableByID)
			setting.POST("asset_tables", middlewares.AssetTablePermission(), handler.CreateAssetTable)
			setting.PUT("asset_tables/:id", middlewares.AssetTablePermission(), handler.UpdateAssetTableByID)
			setting.DELETE("asset_tables/:id", middlewares.AssetTablePermission(), handler.DeleteAssetTableByID)

			setting.GET("ref/surface", middlewares.SurfacePermission(), handler.GetRefSurface)
			setting.GET("ref/surface_param", middlewares.SurfacePermission(), handler.GetParamRefSurface)
			setting.GET("ref/surface/:id", middlewares.SurfacePermission(), handler.GetRefSurfaceByID)
			setting.POST("ref/surface", middlewares.SurfacePermission(), handler.PostRefSurface)
			setting.PUT("ref/surface/:id", middlewares.SurfacePermission(), handler.PutRefSurface)
			setting.DELETE("ref/surface/:id", middlewares.SurfacePermission(), handler.DeleteSettingRefSurfaceByID)

			setting.GET("/budget", middlewares.BudgetPermission(), handler.GetBudget)
			setting.GET("/budget/:id", middlewares.BudgetPermission(), handler.GetBudgetById)
			setting.POST("/budget", middlewares.BudgetPermission(), handler.CreateBudget)
			setting.PUT("/budget", middlewares.BudgetPermission(), handler.UpdateBudget)
			setting.DELETE("/budget/:id", middlewares.BudgetPermission(), handler.DeleteBudget)

			setting.POST("intervention_criteria", middlewares.ModelPermission(), handler.CreateInterventionCriteria)
			setting.GET("intervention_criteria", middlewares.ModelPermission(), handler.GetInterventionCriteria)
			setting.GET("intervention_criteria/list", middlewares.ModelPermission(), handler.GetInterventionCriteriaList)
			setting.GET("intervention_criteria/sequence", middlewares.ModelPermission(), handler.GetInterventionCriteriaSequence)
			setting.POST("intervention_criteria/sequence", middlewares.ModelPermission(), handler.CreateInterventionCriteriaSequence)
			setting.DELETE("intervention_criteria/:id", middlewares.ModelPermission(), handler.DeleteInterventionCriteria)

			setting.POST("/road_work_effect/asphalt", middlewares.ModelPermission(), handler.CreateRoadWorkEffectAsphalt)
			setting.GET("/road_work_effect/asphalt", middlewares.ModelPermission(), handler.GetRoadWorkEffectAsphalt)
			setting.POST("/road_work_effect/concrete", middlewares.ModelPermission(), handler.CreateRoadWorkEffectConcrete)
			setting.GET("/road_work_effect/concrete", middlewares.ModelPermission(), handler.GetRoadWorkEffectConcrete)

			setting.POST("/aadt/growth_rate", middlewares.ModelPermission(), handler.CreateAadtGrowthRate)
			setting.GET("/aadt/growth_rate", middlewares.ModelPermission(), handler.GetAadtGrowthRate)

			setting.POST("/aadt/percentage_vehicle_type", middlewares.ModelPermission(), handler.CreateAadtPercentageVehicleType)
			setting.GET("/aadt/percentage_vehicle_type/:road_group_id", middlewares.ModelPermission(), handler.GetAadtPercentageVehicleType)

			setting.POST("/aadt/parameter", middlewares.ModelPermission(), handler.CreateAadtParameter)
			setting.GET("/aadt/parameter/:road_group_id", middlewares.ModelPermission(), handler.GetAadtParameter)
			setting.GET("/aadt/parameter/road_group_with_volume_aadt", middlewares.ModelPermission(), handler.GetAadtParameterRoadGroupWithVolumeAadt)

			setting.GET("/road_user_cost/acc/loss_value", middlewares.ModelPermission(), handler.GetRoadUserCostAccLossValue)
			setting.POST("/road_user_cost/acc/loss_value", middlewares.ModelPermission(), handler.CreateRoadUserCostAccLossValue)
			setting.GET("/road_user_cost/acc/chance_of_accident/:road_group_id", middlewares.ModelPermission(), handler.GetRoadUserCostAccChanceOfAccident)
			setting.POST("/road_user_cost/acc/chance_of_accident", middlewares.ModelPermission(), handler.CreateRoadUserCostAccChanceOfAccident)
			setting.POST("/road_user_cost/ruc/default_data", middlewares.ModelPermission(), handler.CreateRoadUserCostRucDefaultData)
			setting.GET("/road_user_cost/ruc/default_data", middlewares.ModelPermission(), handler.GetRoadUserCostRucDefaultData)
			setting.POST("/road_user_cost/ruc/driving", middlewares.ModelPermission(), handler.CreateRoadUserCostRucDriving)
			setting.GET("/road_user_cost/ruc/driving", middlewares.ModelPermission(), handler.GetRoadUserCostRucDriving)
			setting.POST("/road_user_cost/ruc/engine_speed", middlewares.ModelPermission(), handler.CreateRoadUserCostRucEngineSpeed)
			setting.GET("/road_user_cost/ruc/engine_speed", middlewares.ModelPermission(), handler.GetRoadUserCostRucEngineSpeed)
			setting.POST("/road_user_cost/ruc/fuel_consumption", middlewares.ModelPermission(), handler.CreateRoadUserCostRucFuelConsumption)
			setting.GET("/road_user_cost/ruc/fuel_consumption", middlewares.ModelPermission(), handler.GetRoadUserCostRucFuelConsumption)
			setting.POST("/road_user_cost/ruc/lubricant_consumption", middlewares.ModelPermission(), handler.CreateRoadUserCostRucLubricantConsumption)
			setting.GET("/road_user_cost/ruc/lubricant_consumption", middlewares.ModelPermission(), handler.GetRoadUserCostRucLubricantConsumption)
			setting.POST("/road_user_cost/ruc/waste_of_consumption", middlewares.ModelPermission(), handler.CreateRoadUserCostRucWasteOfConsumption)
			setting.GET("/road_user_cost/ruc/waste_of_consumption", middlewares.ModelPermission(), handler.GetRoadUserCostRucWasteOfConsumption)
			setting.POST("/road_user_cost/ruc/maintenance", middlewares.ModelPermission(), handler.CreateRoadUserCostRucMaintenance)
			setting.GET("/road_user_cost/ruc/maintenance", middlewares.ModelPermission(), handler.GetRoadUserCostRucMaintenance)
			setting.POST("/road_user_cost/ruc/travel_time", middlewares.ModelPermission(), handler.CreateRoadUserCostRucTravelTime)
			setting.GET("/road_user_cost/ruc/travel_time", middlewares.ModelPermission(), handler.GetRoadUserCostRucTravelTime)
			setting.POST("/road_user_cost/ruc/vehicle_speed_calculation", middlewares.ModelPermission(), handler.CreateRoadUserCostRucVehicleSpeedCalculation)
			setting.GET("/road_user_cost/ruc/vehicle_speed_calculation", middlewares.ModelPermission(), handler.GetRoadUserCostRucVehicleSpeedCalculation)
			setting.POST("/road_user_cost/ruc/traffic_data", middlewares.ModelPermission(), handler.CreateRoadUserCostRucTrafficData)
			setting.GET("/road_user_cost/ruc/traffic_data", middlewares.ModelPermission(), handler.GetRoadUserCostRucTrafficData)

			setting.POST("/optimization", middlewares.ModelPermission(), handler.CreateOptimization)
			setting.GET("/optimization", middlewares.ModelPermission(), handler.GetOptimization)

			setting.POST("/deterioration/asphalt", middlewares.ModelPermission(), handler.CreateDeteriorationAsphalt)
			setting.GET("/deterioration/asphalt/:road_group_id", middlewares.ModelPermission(), handler.GetDeteriorationAsphalt)
			setting.POST("/deterioration/concrete", middlewares.ModelPermission(), handler.CreateDeteriorationConcrete)
			setting.GET("/deterioration/concrete/:road_group_id", middlewares.ModelPermission(), handler.GetDeteriorationConcrete)

			setting.DELETE("/deterioration/asphalt/:road_group_id", middlewares.ModelPermission(), handler.GetDeteriorationAsphalt)

			setting.GET("/hris", middlewares.HrisPermission(), handler.GetHris)
			setting.GET("/hris/:id", middlewares.HrisPermission(), handler.GetHrisById)
			setting.GET("/hris_preview", middlewares.HrisPermission(), handler.GetHrisPreview)
			setting.POST("/hris", middlewares.HrisPermission(), handler.CreateHris)
			setting.PUT("/hris/:id", middlewares.HrisPermission(), handler.UpdateHris)
			setting.DELETE("/hris/:id", middlewares.HrisPermission(), handler.DeleteHris)
			setting.POST("/hris_import", middlewares.HrisPermission(), handler.ImportHris)

			setting.GET("/hsms", middlewares.HsmsPermission(), handler.GetAllHsms)
			setting.DELETE("/hsms/:type/table/:id", middlewares.HsmsPermission(), handler.DeleteHsmsByTypeAndId)
		}

		roadCtrl := v1.Group("/roads", middlewares.AuthorizeJWT(), middlewares.RoadPermission())
		{
			handler := roadRepo.NewRoadRepositoryHandler(databases.DB)
			roadCtrl.GET("tree", handler.GetRoadTree)
			roadCtrl.GET("", handler.GetRoadList)
			roadCtrl.GET(":id", handler.GetRoadByID)
			roadCtrl.GET("menu", handler.GetRoadDetailMenu)
			roadCtrl.GET("init/:id/:level", handler.GetRoadInit)
			//roadCtrl.GET("init/data", handler.GetRoadInitData)
			roadCtrl.GET("/:id/lanse", handler.GetRoadLanes)
			roadCtrl.POST("", handler.CreateRoad)
			roadCtrl.PUT("/:id", handler.UpdateRoad)
			roadCtrl.DELETE(":id", handler.DeleteRoad)

			rcHandler := roadRepoCondition.NewRoadConditionRepositoryHandler(databases.DB)
			roadCtrl.GET("condition/:id_parent", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadCondition)
			roadCtrl.GET("condition_details/:id_parent", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionDetails)
			roadCtrl.GET(":id/condition_list", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionList)
			roadCtrl.GET(":id/condition_template", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionTemplate)
			roadCtrl.GET(":id/condition_compare_lane", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionCompareLane)
			roadCtrl.GET(":id/condition_compare_year", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionCompareYear)
			roadCtrl.GET(":id/condition_compare_average/:lane", middlewares.RoadConditionAccessPermission(), rcHandler.GetRoadConditionCompareAverage)
			roadCtrl.GET(":id/lane_list", middlewares.RoadConditionAccessPermission(), handler.GetRoadDirectionLaneList)
			roadCtrl.POST(":id/condition", middlewares.RoadConditionManagePermission(), rcHandler.CreateRoadCondition)
			roadCtrl.PUT(":id/condition/:id_parent", middlewares.RoadConditionManagePermission(), rcHandler.UpdateRoadCondition)
			roadCtrl.DELETE("condition/:id_parent", middlewares.RoadConditionManagePermission(), rcHandler.DeleteRoadCondition)

			// road damage
			roadDamageHandler := roadDamageRepo.NewRoadDamageRepositoryHandler(databases.DB)
			roadCtrl.GET(":id/damage_list", middlewares.RoadDamageAccessPermission(), roadDamageHandler.GetRoadDamageList)
			roadCtrl.GET(":id/damage_detail/:id_parent", middlewares.RoadDamageAccessPermission(), roadDamageHandler.GetRoadDamageDetail)
			roadCtrl.GET(":id/damage_template", middlewares.RoadDamageAccessPermission(), roadDamageHandler.GetRoadDamageTemplate)
			damageImp := roadCtrl.Group(":id/damage_import")
			{
				damageImp.GET(":id_parent", middlewares.RoadDamageAccessPermission(), roadDamageHandler.GetRoadDamageForImport)
				damageImp.POST("", middlewares.RoadDamageManagePermission(), roadDamageHandler.SetRoadDamageFromImport)
				damageImp.PUT(":id_parent", middlewares.RoadDamageManagePermission(), roadDamageHandler.UpdateRoadDamageFromImport)
				damageImp.DELETE(":id_parent", middlewares.RoadDamageManagePermission(), roadDamageHandler.DeleteRoadDamageForImport)
			}

			rrfhandler := roadRetroReflectivity.NewRepositoryHandler(databases.DB)

			roadCtrl.GET("retro_reflectivity/details/:id_parent", rrfhandler.GetRoadRetroReflectivityDetails)
			roadCtrl.GET("retro_reflectivity/:id_parent", rrfhandler.GetRoadRetroReflectivity)
			roadCtrl.DELETE("retro_reflectivity/:id_parent", rrfhandler.DeleteRoadRetroReflectivity)

			roadCtrl.GET(":id/retro_reflectivity/line_list", rrfhandler.GetRoadRetroReflectivityLineList)
			roadCtrl.GET(":id/retro_reflectivity/list", rrfhandler.GetRoadRetroReflectivityList)
			roadCtrl.POST(":id/retro_reflectivity", rrfhandler.CreateRoadRetroReflectivity)
			roadCtrl.PUT(":id/retro_reflectivity/:id_parent", rrfhandler.UpdateRoadRetroReflectivity)
			roadCtrl.GET(":id/retro_reflectivity/template", rrfhandler.GetRoadRetroReflectivityTemplate)

			roadCtrl.GET(":id/retro_reflectivity/compare_line", rrfhandler.GetRoadRetroReflectivityCompareLine)
			roadCtrl.GET(":id/retro_reflectivity/compare_year", rrfhandler.GetRoadRetroReflectivityCompareYear)
			roadCtrl.GET(":id/retro_reflectivity/compare_average/:line", rrfhandler.GetRoadRetroReflectivityCompareAverage)

			// road asset
			roadAssetHandler := roadAssetRepo.NewRoadAssetRepositoryHandler(databases.DB)
			roadCtrl.GET(":id/asset_revision_list", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetRoadAssetRevisions)
			roadCtrl.GET(":id/asset_details/:road_asset_id", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetRoadAssetDetail)
			roadCtrl.GET("asset_permission", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetRoadAssetPermission)
			roadCtrl.GET("asset_edit_template", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetRoadAssetTemplate)
			roadCtrl.GET(":id/km", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetRoadKmByGeom)
			roadCtrl.GET("asset_table/:id", middlewares.RoadAssetAccessPermission(), roadAssetHandler.GetAssetRoadType)
			roadCtrl.PUT(":id/asset/:id_parent_asset", middlewares.RoadAssetManagePermission(), roadAssetHandler.UpdateRoadAsset)
			roadCtrl.PUT("asset_cancel/:id_parent", middlewares.RoadAssetManagePermission(), roadAssetHandler.CancelRoadAsset)
			roadCtrl.PUT("asset_confirm/:id_parent", middlewares.RoadAssetManagePermission(), roadAssetHandler.ConfirmRoadAsset)
			roadCtrl.POST(":id/asset", middlewares.RoadAssetManagePermission(), roadAssetHandler.CreateRoadAsset)
			roadCtrl.DELETE("asset_delete/:id_parent", middlewares.RoadAssetManagePermission(), roadAssetHandler.DeleteRoadAsset)
			roadCtrl.DELETE("asset/:id/table/:ref_asset_table_id/asset_object/:asset_object_id", middlewares.RoadAssetManagePermission(), roadAssetHandler.DeleteRoadAssetObject)

			// Maintain History
			handlerM := maintenanceRepo.NewMaintainRepoHandler(databases.DB)
			roadCtrl.GET(":id/maintenance", handlerM.GetMaintenanceByRoadID)
			roadCtrl.GET(":id/maintenance_year", handlerM.GetMaintenanceByRoadYear)

			handlerV := volumeAadtRepo.NewRepositoryHandler(databases.DB)
			roadCtrl.GET(":id/volume_aadt/revision", middlewares.RoadAadtAccessPermission(), handlerV.GetVolumeRevision)
			roadCtrl.GET(":id/volume_aadt/:aadt_id", middlewares.RoadAadtAccessPermission(), handlerV.GetVolume)
			roadCtrl.POST(":id/volume_aadt", middlewares.RoadAadtManagePermission(), handlerV.CreateVolume)
			roadCtrl.PUT(":id/volume_aadt/:id_parent/aadt/:aadt_id", middlewares.RoadAadtManagePermission(), handlerV.UpdateVolume)
			roadCtrl.DELETE(":id/volume_aadt/:aadt_id", middlewares.RoadAadtManagePermission(), handlerV.DeleteVolume)
		}

		handler := roadSurface.NewRoadSurfaceRepositoryHandler(databases.DB)
		roadSurface := v1.Group("/road")
		roadSurface.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			roadSurface.Use(middlewares.AuthorizeJWT(), middlewares.RoadPermission())
			{
				roadSurface.GET("/surface", handler.GetRoadSurface)
				roadSurface.POST("/surface", handler.PostRoadSurface)
				roadSurface.GET("/surface/icon/:road_id", handler.GetIconRoadSurface)
			}
		}

		maintain := v1.Group("/maintenance")
		maintain.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			maintain.Use(middlewares.AuthorizeJWT())
			{
				handler := maintenanceRepo.NewMaintainRepoHandler(databases.DB)
				// master
				maintain.GET("intervention_criteria", middlewares.MaintenanceAccessPermission(), handler.GetInterventionCriteria)
				//maintain.GET("plan_stauts", middlewares.MaintenanceAccessPermission(), handler.GetMaintenancePlanStatus)
				maintain.GET("budgets", middlewares.MaintenanceAccessPermission(), handler.GetMaintenanceBudget)
				//maintain.GET("road_division_filter", middlewares.MaintenanceAccessPermission(), handler.GetRoadDivisionFilter)
				maintain.GET("road_dropdown_list", middlewares.MaintenanceAccessPermission(), handler.GetRoadDropdownList)
				maintain.GET("road_dropdown_list_analyze", middlewares.MaintenanceAccessPermission(), handler.GetRoadDropdownListAnalyze)
				maintain.GET("road_dropdown_list_dashboard", middlewares.MaintenanceAccessPermission(), handler.GetRoadDropdownListDashboard)
				maintain.GET("division", middlewares.MaintenanceAccessPermission(), handler.GetDivisionList)

				//GetMaintenance
				//ประวัติซ่อมบำรุง
				maintain.GET("", middlewares.MaintenanceAccessPermission(), handler.GetMaintenance)
				maintain.GET(":id_parent", middlewares.MaintenanceAccessPermission(), handler.GetMaintenanceByID)
				maintain.GET("years", middlewares.MaintenanceAccessPermission(), handler.GetMaintenanceYear)
				maintain.POST("", middlewares.MaintenanceManagePermission(), handler.CreateMaintenance)
				maintain.PUT(":id_parent", middlewares.MaintenanceManagePermission(), handler.UpdateMaintenance)
				maintain.DELETE(":id_parent", middlewares.MaintenanceManagePermission(), handler.DeleteMaintenance)

				//เพิ่มข้อมูลการซ่อมบำรุง
				maintain.GET(":id_parent/road/:m_road_id", middlewares.MaintenanceAccessPermission(), handler.GetMaintenanceRoadByID)
				maintain.POST(":id_parent/road", middlewares.MaintenanceAccessPermission(), handler.CreateMaintenanceRoad)
				maintain.PUT(":id_parent/road/:m_road_id", middlewares.MaintenanceAccessPermission(), handler.UpdateMaintenanceRoad)
				maintain.DELETE(":id_parent/road/:m_road_id", middlewares.MaintenanceAccessPermission(), handler.DeleteMaintenanceRoad)

				// ข้อมูลประวัติการซ่อมบำรุงในช่วงค้ำประกัน
				maintain.GET("/:id_parent/road_history/:m_road_his_id", middlewares.MaintenanceHisAccessPermission(), handler.GetMaintenanceHistoryByID)
				maintain.POST("/:id_parent/road_history", middlewares.MaintenanceHisManagePermission(), handler.CreateMaintenanceRoadHistory)
				maintain.PUT("/:id_parent/road_history/:m_road_his_id", middlewares.MaintenanceHisManagePermission(), handler.UpdateMaintenanceRoadHistory)
				maintain.DELETE("/:id_parent/road_history/:m_road_his_id", middlewares.MaintenanceHisManagePermission(), handler.DeleteMaintenanceRoadHistory)
			}
		}
		roadGrp := v1.Group("/road_group")
		roadGrp.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			roadGrp.Use(middlewares.AuthorizeJWT())
			{
				handler := roadGroupRepo.NewRepositoryHandler(databases.DB)
				roadGrp.GET("", handler.GetRoadGroup)
				roadGrp.GET(":id", handler.GetRoadGroupByID)
			}

			roadGrp.Use(middlewares.AuthorizeJWT())
			{
				handler := volumeAccidentRepo.NewRepositoryHandler(databases.DB)
				roadGrp.GET(":id/volume_accident/revision", middlewares.RoadAccidentAccessPermission(), handler.GetVolumeRevision)
				roadGrp.GET(":id/volume_accident/:accident_id", middlewares.RoadAccidentAccessPermission(), handler.GetVolume)
				roadGrp.POST(":id/volume_accident", middlewares.RoadAccidentManagePermission(), handler.CreateVolume)
				roadGrp.PUT(":id/volume_accident/:id_parent/accident/:accident_id", middlewares.RoadAccidentManagePermission(), handler.UpdateVolume)
				roadGrp.DELETE(":id/volume_accident/:accident_id", middlewares.RoadAccidentManagePermission(), handler.DeleteVolume)
			}

		}

		roadSection := v1.Group("/road_section")
		roadSection.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			roadSection.Use(middlewares.AuthorizeJWT())
			{
				handler := roadSectionRepo.NewRepositoryHandler(databases.DB)
				roadSection.GET("", handler.GetRoadSection)
				roadSection.GET(":id", handler.GetRoadSectionByID)
			}
		}

		analyze := v1.Group("/analyze")
		analyze.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			analyze.Use(middlewares.AuthorizeJWT())
			{
				handler := maintenanceAnalysisRepo.NewRepositoryHandler(databases.DB, databases.MongoDb)
				analyze.GET("", middlewares.AnalyzeAccessPermission(), handler.GetMaintenanceAnalysis)                           // หน้า listting
				analyze.GET(":id", middlewares.AnalyzeAccessPermission(), handler.GetMaintenanceAnalysisById)                    // step 1
				analyze.POST("prepare_data", middlewares.AnalyzeManagePermission(), handler.CreateMaintenanceAnalysis)           // ค้นหา
				analyze.PUT(":id/prepare_data", middlewares.AnalyzeManagePermission(), handler.UpdateMaintenanceAnalysis)        // ค้นหาเมื่อแก้ไข
				analyze.POST(":id/condition", middlewares.AnalyzeManagePermission(), handler.GetMaintenanceAnalysisCondition)    // step 2
				analyze.POST(":id/analyzing", middlewares.AnalyzeManagePermission(), handler.CreateMaintenanceAnalysisAnalyzing) // เริ้มวิเคราะห์
				analyze.DELETE(":id", middlewares.AnalyzeManagePermission(), handler.DeleteMaintenanceAnalysis)                  // step 1
				analyze.POST(":id/copy", middlewares.AnalyzeManagePermission(), handler.CopyMaintenanceAnalysis)
				analyze.POST(":id/favorite", middlewares.AnalyzeManagePermission(), handler.FavoriteMaintenanceAnalysis) // step 1
				analyze.GET(":id/export_data", middlewares.AnalyzeManagePermission(), handler.ExportData)                // step 1

				analyze.GET(":id/model", middlewares.AnalyzeManagePermission(), handler.GetAnalysisModel)
				analyze.PUT(":id/model", middlewares.AnalyzeManagePermission(), handler.UpdateAnalysisModel)
				analyze.GET("intervention_criterias", middlewares.AnalyzeManagePermission(), handler.GetRefCriteriaMethod)

				analyze.GET(":id/prepare_data_selected", handler.GetPrepareDataAllByAnalysisSelected)

				// report
				analyze.GET(":id/report/report1", middlewares.AnalyzeManagePermission(), handler.GetReport1)
				analyze.GET(":id/report/report2", middlewares.AnalyzeManagePermission(), handler.GetReport2)
				analyze.GET(":id/report/report3", middlewares.AnalyzeManagePermission(), handler.GetReport3)
				analyze.GET(":id/report/report4", middlewares.AnalyzeManagePermission(), handler.GetReport4)
				analyze.GET(":id/report/report5", middlewares.AnalyzeManagePermission(), handler.GetReport5)

				// dashboard
				analyze.GET("dashboard/strategic/:id", middlewares.AnalyzeManagePermission(), handler.DashboardStrategicMaintenanceAnalysis)
				analyze.GET("dashboard/annual/:id", middlewares.AnalyzeManagePermission(), handler.DashboardAnnualMaintenanceAnalysis)

				// New
				analyze.GET("/:id/check_prepare_data", middlewares.AnalyzeManagePermission(), handler.CheckPrepareDataById)
				analyze.GET("/:id/prepare_data_id", middlewares.AnalyzeManagePermission(), handler.GetPrepareDataIdById)
				analyze.GET("/:id/prepare_data", middlewares.AnalyzeManagePermission(), handler.GetPrepareDataById)

				analyze.GET("dashboard-map/:id/filter", middlewares.AnalyzeManagePermission(), handler.DashboardMapFilter)
				analyze.GET("dashboard-map/:id", middlewares.AnalyzeManagePermission(), handler.DashboardMap)
			}
		}
		dashboard := v1.Group("/dashboard")
		handler2 := dashboardRepo.NewRepositoryHandler(databases.DB)
		dashboard.GET("surface/data_mart_sys", handler2.GetDataMartSys)
		dashboard.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			dashboard.Use(middlewares.AuthorizeJWT())
			{
				handler := dashboardRepo.NewRepositoryHandler(databases.DB)
				dashboard.GET("road", handler.GetRoadDashboard)
				dashboard.GET("years", handler.GetDashboardYears)

				// road damage report
				dashboard.GET("asset", handler.GetAsset)
				dashboard.GET("asset_detail", handler.GetAssetDetail)
				dashboard.GET("asset_map", handler.GetAssetMap)
				dashboard.GET("asset_map/:id/detail/:asset_table_id", handler.GetAssetMapDetailByID)

				//  road surface
				dashboard.GET("surface/data_mart_check", handler.GetDataMartCheck)
				dashboard.GET("surface/data_mart", handler.GetDataMart)
				dashboard.GET("surface", handler.GetSurfaceDashboard)
				dashboard.GET("surface_map", handler.GetSurfaceDashboardMap)

				// road condition
				dashboard.GET("condition", handler.GetDashboardCondition)
				dashboard.GET("condition_map", handler.GetDashboardConditionMap)

				handlerMaintenance := dashboardRepo.NewRepositoryHandlerMaintenance(databases.DB)
				dashboard.GET("maintenance", handlerMaintenance.GetMaintenanceDashboard)
				dashboard.GET("maintenance_map", handlerMaintenance.GetMaintenanceMapDashboard)
				dashboard.GET("maintenance_table", handlerMaintenance.GetMaintenanceTableDashboard)

			}
		}

		report := v1.Group("/report")
		report.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			report.Use(middlewares.AuthorizeJWT(), middlewares.ReportPermission())
			{
				handler := reportRepo.NewRepositoryHandler(databases.DB)

				//////////////// NEW MIMS ////////////////
				report.GET("status", handler.ReportStatus)
				report.GET("status/check/:id", handler.CheckReportStatus)

				report.GET("asset/filter/type1", handler.GetAssetFilterType1)
				report.GET("asset/filter/type2", handler.GetAssetFilterType2)
				report.GET("asset/filter/type3", handler.GetAssetFilterType3)

				report.GET("road/filter/type1", handler.GetRoadFilterType1)
				report.GET("road/filter/type2", handler.GetRoadFilterType2)
				report.GET("road/filter/type3", handler.GetRoadFilterType3)
				report.GET("road/filter/type4", handler.GetRoadFilterType4)
				report.GET("road/filter/type5", handler.GetRoadFilterType5)
				report.GET("road/filter/type6", handler.GetRoadFilterType6)

				report.GET("road_damage/filter/type1", handler.GetRoadDamageFilterType1)
				report.GET("road_damage/filter/type2", handler.GetRoadDamageFilterType2)

				report.GET("maintenance/filter/type1", handler.GetMaintenanceFilterType1)

				report.GET("maintenance_kpi/filter/type1", handler.GetMaintenanceKpiFilterType1)

				report.GET("aadt/filter/type1", handler.GetAadtFilterType1)

				report.GET("type1", handler.GetReport1)
				report.GET("type2", handler.GetReport2)
				report.GET("type3", handler.GetReport3)
				report.GET("type4", handler.GetReportRoad)
				report.GET("type5", handler.Report5)
				report.GET("type6", handler.Report6)
				report.GET("type7", handler.Report7)
				report.GET("type8", handler.Report8)
				report.GET("type9", handler.Report9)
				report.GET("type10", handler.Report10)
				report.GET("type11", handler.Report11)
				report.GET("type12", handler.Report12)
				report.GET("type13", handler.Report13)
				report.GET("type14", handler.GetReportTrafficVolume)
				//////////////// NEW MIMS ////////////////
			}
		}

		hris := v1.Group("/hris")
		hris.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
		{
			handler := hrisRepo.NewRepository(databases.DB, databases.MongoDb)
			hris.GET("section_geom", handler.GetSectionGeom)
			hris.GET("road_latest", handler.GetRoadLatest)

			hris.GET("match_data", handler.MatchData)

		}
	}

	hsms := v1.Group("/hsms")
	hsms.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
	{
		handler := hsmsRepo.NewRepositoryHandler(databases.DB)
		hsms.GET("bridge", handler.GetHsmsBridge)
		hsms.GET("guard", handler.GetHsmsGuard)
		hsms.GET("interchange", handler.GetHsmsInterchange)
		hsms.GET("intersection", handler.GetHsmsIntersection)
		hsms.GET("streetlight", handler.GetHsmsStreetlight)
		hsms.GET("railwaycrossing", handler.GetHsmsRailwaycrossing)
		hsms.GET("trafficlight", handler.GetHsmsTrafficlight)
		hsms.GET("uturnbridge", handler.GetHsmsUturnbridge)
	}

	CheckLog := v1.Group("/discord")
	CheckLog.Use(middlewares.RequestLoggingMiddleware(logs.LogRequest))
	{
		handler := discordRepo.NewRepository(databases.DB, databases.MongoDb)
		CheckLog.GET("log", handler.GetHrisAndHsmsLog)
	}

	return r
}

func selectAPIPath(r *gin.Engine, env string) *gin.RouterGroup {
	if env == "dev" {
		return r.Group("api/v1")
	}

	return r.Group("/v1")
}
