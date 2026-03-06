package usecases_test

// import (
// 	"errors"
// 	"mime/multipart"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jinzhu/copier"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	helpersmocks "gitlab.com/mims-api-service/helpers/file/mocks"
// 	mockdata "gitlab.com/mims-api-service/mock"
// 	"gitlab.com/mims-api-service/models"
// 	requests "gitlab.com/mims-api-service/requests"
// 	responses "gitlab.com/mims-api-service/responses"
// 	"gitlab.com/mims-api-service/src/road/mocks/repomocks"
// 	"gitlab.com/mims-api-service/src/road/usecases"
// 	"gorm.io/gorm"
// )

// // command for generate mock function file if you change RoadRepository
// // mockery --dir=src/road/domains --name=RoadRepository --filename=roadRepository_mock.go --output=src/road/mocks/repomocks --outpkg=repomocks

// // command for generate mock function file if you change HelpersFileDomain
// // mockery --dir=helpers/file --name=HelpersFileDomain --filename=file_mock.go --output=helpers/file/mocks --outpkg=helpersmocks

// func TestGetRoadByID(t *testing.T) {

// 	t.Run("get data road success", func(t *testing.T) {
// 		mockRoadById := mockdata.RoadByIdMock
// 		roadId := 1
// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadByID", roadId).Return(&mockRoadById, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		road, err := roadUsecase.GetRoadByID(roadId)

// 		var expectedRoad responses.RoadById
// 		copier.Copy(&expectedRoad, &mockdata.RoadByIdMock)
// 		expectedRoad.RoadCode = "00070001"
// 		expectedRoad.RoadSectionNameTH = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoad.RoadSectionNameEN = "Origin Name (English) - Destination Name (English)"
// 		expectedRoad.Province = "กรุงเทพมหานคร"
// 		expectedRoad.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoad.OriginToDestination = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoad.KmRange = "0 - 10.5"
// 		expectedRoad.Distance = 0.0105
// 		expectedRoad.RoadInfo.CenterLaneShapeFilePath = "//path/to/center_lane_shape_file"
// 		expectedRoad.RoadInfo.CenterLineShapeFilePath = "//path/to/center_line_shape_file"

// 		assert.Equal(t, expectedRoad, *road)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road case section name destination nil string", func(t *testing.T) {
// 		roadId := 1

// 		roadMock := mockdata.RoadByIdMock
// 		roadMock.RoadSection.NameDestinationTH = ""

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadByID", roadId).Return(&roadMock, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		road, err := roadUsecase.GetRoadByID(roadId)

// 		var expectedRoad responses.RoadById
// 		copier.Copy(&expectedRoad, &mockdata.RoadByIdMock)
// 		expectedRoad.RoadCode = "00070001"
// 		expectedRoad.RoadSectionNameTH = "Road Origin Name (Thai)"
// 		expectedRoad.RoadSectionNameEN = "Origin Name (English)"
// 		expectedRoad.Province = "กรุงเทพมหานคร"
// 		expectedRoad.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoad.OriginToDestination = "Road Origin Name (Thai)"
// 		expectedRoad.KmRange = "0 - 10.5"
// 		expectedRoad.Distance = 0.0105
// 		expectedRoad.RoadInfo.CenterLaneShapeFilePath = "//path/to/center_lane_shape_file"
// 		expectedRoad.RoadInfo.CenterLineShapeFilePath = "//path/to/center_line_shape_file"

// 		assert.Equal(t, expectedRoad, *road)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road error call repository", func(t *testing.T) {
// 		roadId := 1
// 		newError := errors.New("error call database")
// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadByID", roadId).Return(nil, newError)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		road, err := roadUsecase.GetRoadByID(roadId)

// 		expectedError := responses.AppErr(responses.AppErr{StatusCode: 404, Message: "error call database"})

// 		assert.Nil(t, road)
// 		assert.Equal(t, err, &expectedError)
// 	})

// }

// func TestGetRoadGroupList(t *testing.T) {

// 	t.Run("get data road list success", func(t *testing.T) {
// 		params := requests.RoadPrams{
// 			Keyword:       "",
// 			RoadGroupId:   []int{},
// 			RoadSectionId: []int{},
// 			KmStart:       nil,
// 			KmEnd:         nil,
// 			RefSurfaceId:  []int{1, 2},
// 			DepotCode:     []string{},
// 		}

// 		newMockDataRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{
// 									mockdata.ChildRoadDataMaock,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadGroupList", params).Return(newMockDataRoadList, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		roadList, err := roadUsecase.GetRoadGroupList(params)

// 		expectedRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{
// 									mockdata.ChildRoadDataMaock,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		expectedRoadList[0].Number = "7"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.Name = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.OriginToDestination = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.KmRange = "0 - 10.5"

// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.KmRange = "0 - 10.5"

// 		assert.Equal(t, expectedRoadList, roadList)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road list success case surface child not found", func(t *testing.T) {
// 		params := requests.RoadPrams{
// 			Keyword:       "",
// 			RoadGroupId:   []int{},
// 			RoadSectionId: []int{},
// 			KmStart:       nil,
// 			KmEnd:         nil,
// 			RefSurfaceId:  []int{1, 2},
// 			DepotCode:     []string{},
// 		}

// 		newMockDataRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{
// 									mockdata.ChildRoadDataFoundRefSurfaceMaock,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadGroupList", params).Return(newMockDataRoadList, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		roadList, err := roadUsecase.GetRoadGroupList(params)

// 		expectedRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		expectedRoadList[0].Number = "7"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.Name = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.OriginToDestination = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.KmRange = "0 - 10.5"

// 		// expectedRoadList[0].Sections[0].Roads[0].ChildRoads =

// 		assert.Equal(t, expectedRoadList, roadList)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road case section name destination nil string", func(t *testing.T) {
// 		params := requests.RoadPrams{
// 			Keyword:       "",
// 			RoadGroupId:   []int{},
// 			RoadSectionId: []int{},
// 			KmStart:       nil,
// 			KmEnd:         nil,
// 			RefSurfaceId:  []int{1, 2},
// 			DepotCode:     []string{},
// 		}
// 		newMockDataRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionNilDestMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadGroupList", params).Return(newMockDataRoadList, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		roadList, err := roadUsecase.GetRoadGroupList(params)

// 		expectedRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionNilDestMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}
// 		expectedRoadList[0].Number = "7"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.Name = "Road Origin Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.OriginToDestination = "Road Origin Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.KmRange = "0 - 10.5"

// 		assert.Equal(t, roadList, expectedRoadList)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road case params keyword", func(t *testing.T) {
// 		params := requests.RoadPrams{
// 			Keyword:       "Road",
// 			RoadGroupId:   []int{},
// 			RoadSectionId: []int{},
// 			KmStart:       nil,
// 			KmEnd:         nil,
// 			RefSurfaceId:  []int{1, 2},
// 			DepotCode:     []string{},
// 		}
// 		newMockDataRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{
// 									mockdata.ChildRoadDataMaock,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadGroupList", params).Return(newMockDataRoadList, nil)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		roadList, err := roadUsecase.GetRoadGroupList(params)

// 		expectedRoadList := []models.RoadList{
// 			{
// 				RoadGroup: mockdata.RoadGroupMock,
// 				Sections: []models.RoadSectionData{
// 					{
// 						RoadSection: mockdata.RoadSectionMock,
// 						Roads: []models.RoadData{
// 							{
// 								Road:       mockdata.RoadMock,
// 								RoadInfo:   mockdata.RoadInfoAddDataMock,
// 								RefSurface: mockdata.RefSurfaceRoadMock,
// 								RoadSurfaceIcon: []models.RoadSurfaceIcon{
// 									mockdata.RoadSurfaceIconMock,
// 								},
// 								RoadGeom: []models.RoadGeom{
// 									mockdata.RoadGeomMock,
// 								},
// 								ChildRoads: []models.ChildRoadData{
// 									mockdata.ChildRoadDataMaock,
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		expectedRoadList[0].Number = "7"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.Name = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.OriginToDestination = "Road Origin Name (Thai) - Destination Name (Thai)"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].RoadInfo.KmRange = "0 - 10.5"

// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.RoadCode = "00070001"
// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.ResponsibleCode = "ขท.ธนบุรี - หมวดทางหลวงบางขุนเทียน"
// 		expectedRoadList[0].Sections[0].Roads[0].ChildRoads[0].RoadInfo.KmRange = "0 - 10.5"

// 		assert.Equal(t, roadList, expectedRoadList)
// 		assert.Nil(t, err)
// 	})

// 	t.Run("get data road error call repository", func(t *testing.T) {
// 		params := requests.RoadPrams{
// 			Keyword:       "",
// 			RoadGroupId:   []int{},
// 			RoadSectionId: []int{},
// 			KmStart:       nil,
// 			KmEnd:         nil,
// 			RefSurfaceId:  []int{1, 2},
// 			DepotCode:     []string{},
// 		}

// 		newError := errors.New("error call database")
// 		roadRepo := repomocks.NewRoadRepository(t)
// 		roadRepo.On("GetRoadGroupList", params).Return([]models.RoadList{}, newError)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		roadList, err := roadUsecase.GetRoadGroupList(params)

// 		expectedError := responses.AppErr(responses.AppErr{StatusCode: 404, Message: "error call database"})

// 		assert.Nil(t, roadList)
// 		assert.Equal(t, err, &expectedError)
// 	})
// }

// func TestCreateRoad(t *testing.T) {
// 	os.Setenv("ROAD_CENTER_LINE_SHAPE_FILE_DIR", "/path/to/mock/dir")
// 	os.Setenv("ROAD_CENTER_LANE_SHAPE_FILE_DIR", "/path/to/mock/dir")
// 	t.Run("create data road  success", func(t *testing.T) {

// 		// Create a mock context
// 		w := httptest.NewRecorder()
// 		c, _ := gin.CreateTestContext(w)

// 		roadReq := requests.Road{
// 			Name:                "Test Road",
// 			RoadSectionID:       1,
// 			RoadID:              2,
// 			RoadLevel:           1,
// 			KmStart:             0,
// 			KmEnd:               14760,
// 			RampId:              "Ramp123",
// 			RoadColorCode:       "#FFFFFF",
// 			RefRoadTypeID:       4,
// 			RegisterDate:        "2024-02-29 15:04:05",
// 			Remark:              "Test Remark",
// 			CenterLineShapeFile: &multipart.FileHeader{Filename: "LineShape"},
// 			CenterLaneShapeFile: &multipart.FileHeader{Filename: "LaneShape"},
// 		}

// 		roadRepo := repomocks.NewRoadRepository(t)
// 		helperFile := helpersmocks.NewHelpersFileDomain(t)
// 		lineString := "LINESTRING(100.64237403383734 13.741329869405819)"
// 		kmStartStr := "0+000"
// 		kmEndStr := "14+760"
// 		tx := &gorm.DB{} // mocked tx object
// 		helperFile.On("EnsureDir", "/path/to/mock/dir").Return(nil)
// 		helperFile.On("SaveFile", c, roadReq.CenterLineShapeFile, "/path/to/mock/dir").Return("LineShape.zip", nil)
// 		helperFile.On("SaveFile", c, roadReq.CenterLaneShapeFile, "/path/to/mock/dir").Return("LaneShape.zip", nil)
// 		helperFile.On("ReadCenterLineShepesFile", "LineShape.zip").Return(&lineString, &kmStartStr, &kmEndStr, nil)
// 		helperFile.On("ReadCenterLaneShepesFile", "LaneShape.zip").Return([]models.RoadGeom{mockdata.RoadGeomMock}, nil)

// 		roadRepo.On("GetDataById", mock.AnythingOfType("*models.RoadSection"), 1).Return(nil)
// 		roadRepo.On("GetRoadMaxSeq").Return(1, nil)
// 		roadRepo.On("StartTransSection").Return(tx)
// 		roadRepo.On("CreateData", tx, mock.AnythingOfType("*models.Road")).Return(nil)
// 		roadRepo.On("CreateRoadInfo", tx, mock.AnythingOfType("*models.RoadInfo")).Return(nil)
// 		roadRepo.On("CreateRoadGeom", tx, mock.AnythingOfType("[]models.RoadGeom")).Return(nil)
// 		roadRepo.On("Commit", tx).Return(nil)
// 		roadUsecase := usecases.NewRoadUseCase(roadRepo, helperFile)

// 		_, err := roadUsecase.CreateRoad(c, 1, roadReq)
// 		// idInt := id.(int)

// 		assert.Nil(t, err)
// 		// assert.Nil(t, id)
// 	})
// }
