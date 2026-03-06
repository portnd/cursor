package repositories

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/mims-api-service/models"
)

type authRepositoryMock struct {
	mock.Mock
}

func NewAuthRepositoryMock() *authRepositoryMock {
	return &authRepositoryMock{}
}

func (m *authRepositoryMock) GetUserByUserName(username string) (models.Users, error) {
	args := m.Called(username)
	return args.Get(0).(models.Users), args.Error(1)
}

func (m *authRepositoryMock) GetUserByID(userId string) (models.Users, error) {
	args := m.Called(userId)
	return args.Get(0).(models.Users), args.Error(1)
}

func (m *authRepositoryMock) GetUserByEmail(email string) (models.Users, error) {
	args := m.Called(email)
	return args.Get(0).(models.Users), args.Error(1)
}

func (m *authRepositoryMock) GetUserByResetPasswordToken(token string) (models.Users, error) {
	args := m.Called(token)
	return args.Get(0).(models.Users), args.Error(1)
}

func (m *authRepositoryMock) UpdateUser(user models.Users) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *authRepositoryMock) CreateAuth(authDetail models.Auth) error {
	args := m.Called(authDetail)
	return args.Error(0)
}

func (m *authRepositoryMock) GetAuthByUserId(userId uint) (models.Auth, error) {
	args := m.Called(userId)
	return args.Get(0).(models.Auth), args.Error(1)
}

func (m *authRepositoryMock) DeleteAuthByUserId(userId uint) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *authRepositoryMock) UpdateAuth(authDetail models.Auth) error {
	args := m.Called(authDetail)
	return args.Error(0)
}
