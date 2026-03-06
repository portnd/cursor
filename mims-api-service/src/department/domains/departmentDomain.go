package domains

// business logic
type DepartmentUseCase interface {
	// User(param interface{}) ([]models.User, error)
	// Getss() Param
	// GenerateToken(email string, isUser bool) string
	// ValidateToken(encodedToken string) (*jwt.Token, error)
}

// อะไรเชื่อมต่อกับ DB
type DepartmentRepository interface {
	// GetUserAll(user *[]models.User, offset, limit int, column, dir string, search, columnOrder map[string]string) (err error)
	// GetUserByEmail(user *models.User, email, password string) (err error)
}
