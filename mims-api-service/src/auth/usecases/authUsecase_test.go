package usecases_test

import (
	"fmt"
	"reflect"
	"testing"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/auth/usecases"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func isEqual(t *testing.T, want interface{}, got interface{}) {
	wantValues := reflect.ValueOf(want)
	gotValues := reflect.ValueOf(got)

	for i := 0; i < wantValues.NumField(); i++ {
		assert.Equal(t, wantValues.Field(i).Interface(), gotValues.Field(i).Interface())
	}
}

// func TestLoginUseCase(t *testing.T) {

// 	t.Run("login success case", func(t *testing.T) {
// 		authRepo := repositories.NewAuthRepositoryMock()
// 		authRepoUseCase := usecases.NewAuthUseCase(authRepo)

// 		accessToken, refreshToken, err := authRepoUseCase.Login("userhidden1010@gmail.com", "qwer1234")

// 		fmt.Println(accessToken, refreshToken, err)
// 	})

// }

func Test_GenerateUUIDV4(t *testing.T) {
	//arrange
	expectLength := 36

	//act
	got := usecases.GenerateUUIDV4()

	//assert
	assert.Len(t, got, expectLength)
}

func Test_GenerateAuthDetail(t *testing.T) {
	type args struct {
		userId uint
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{fmt.Sprintf("when user id is %d", 1), args{userId: 1}},
		{fmt.Sprintf("when user id is %d", 2), args{userId: 2}},
		{fmt.Sprintf("when user id is %d", 3), args{userId: 3}},
		{fmt.Sprintf("when user id is %d", 4), args{userId: 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usecases.GenerateAuthDetail(tt.args.userId)

			assert.Equal(t, tt.args.userId, got.UserID)
			assert.Len(t, got.AccessUUID, 36)
			assert.Len(t, got.RefreshUUID, 36)
		})
	}
}

func TestCreateToken(t *testing.T) {
	var accessControl []string
	accessControl = append(accessControl, "dashboard_road_condition_access")
	accessControl = append(accessControl, "dashboard_road_surface_access")
	tests := []struct {
		name          string
		authDetail    models.Auth
		accessControl []string
		departmentId  int
	}{
		// TODO: Add test cases.
		{
			name: fmt.Sprintf("create token when user id is %d", 1),
			authDetail: models.Auth{
				UserID:      1,
				AccessUUID:  usecases.GenerateUUIDV4(),
				RefreshUUID: usecases.GenerateUUIDV4(),
			},
			accessControl: accessControl,
			departmentId:  1,
		},
		{
			name: fmt.Sprintf("create token when user id is %d", 2),
			authDetail: models.Auth{
				UserID:      2,
				AccessUUID:  usecases.GenerateUUIDV4(),
				RefreshUUID: usecases.GenerateUUIDV4(),
			},
			accessControl: accessControl,
			departmentId:  1,
		},
		{
			name: fmt.Sprintf("create token when user id is %d", 3),
			authDetail: models.Auth{
				UserID:      3,
				AccessUUID:  usecases.GenerateUUIDV4(),
				RefreshUUID: usecases.GenerateUUIDV4(),
			},
			accessControl: accessControl,
			departmentId:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAccessToken, gotRefreshToken, err := usecases.CreateToken(tt.authDetail, tt.accessControl)

			access, _ := helpers.ValidateToken(gotAccessToken)
			refresh, _ := helpers.ValidateToken(gotRefreshToken)

			accessPayload, _ := access.Claims.(jwt.MapClaims)
			refreshPayload, _ := refresh.Claims.(jwt.MapClaims)

			assert.NoError(t, err)
			assert.Equal(t, tt.authDetail.UserID, uint(accessPayload["user_id"].(float64)))
			assert.Equal(t, tt.authDetail.UserID, uint(refreshPayload["user_id"].(float64)))
			assert.Equal(t, tt.authDetail.AccessUUID, accessPayload["access_uuid"])
			assert.Equal(t, tt.authDetail.RefreshUUID, refreshPayload["refresh_uuid"])
		})
	}
}
