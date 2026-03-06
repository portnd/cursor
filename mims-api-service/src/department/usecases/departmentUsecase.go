package usecases

import (
	"gitlab.com/mims-api-service/src/department/domains"
)

type departmentUseCase struct {
	departmentUseCaseRepo domains.DepartmentRepository
}

// init usecase
func NewDepartmentUseCase(repo domains.DepartmentRepository) domains.DepartmentUseCase {
	return &departmentUseCase{
		departmentUseCaseRepo: repo,
	}
}

// =========================================================
