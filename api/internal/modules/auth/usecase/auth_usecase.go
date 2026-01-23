package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/komgrip/starter-kit/internal/modules/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

// authUsecase implements the business logic for authentication
// This is the CORE of Hexagonal Architecture - pure business logic
type authUsecase struct {
	repo      domain.Repository
	jwtSecret []byte
}

// NewAuthUsecase creates a new authentication usecase instance
// Follows Dependency Injection pattern
func NewAuthUsecase(repo domain.Repository) domain.Usecase {
	// Load JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Fallback to a default secret (NOT recommended for production)
		jwtSecret = "komgrip-default-secret-change-in-production"
	}

	return &authUsecase{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

// Register handles user registration
// Business Rules:
// 1. Check if email already exists
// 2. Hash the password using bcrypt
// 3. Create user in database
// 4. Generate JWT token
func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// Check if user with this email already exists
	existingUser, _ := u.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Hash the password using bcrypt (cost factor: 12)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user entity
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Save to database
	if err := u.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := u.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response with token and user info
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// Login handles user login
// Business Rules:
// 1. Find user by email
// 2. Verify password using bcrypt.CompareHashAndPassword
// 3. Generate JWT token if password is correct
func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.AuthResponse, error) {
	// Find user by email
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := u.generateJWT(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Return response with token and user info
	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

// generateJWT creates a JWT token with user claims
// Token includes: user_id, email, issued_at, expires_at (24 hours)
// Algorithm: HS256 (HMAC with SHA-256)
func (u *authUsecase) generateJWT(userID uint, email string) (string, error) {
	// Define JWT claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iat":     time.Now().Unix(),                        // Issued At
		"exp":     time.Now().Add(24 * time.Hour).Unix(),    // Expires in 24 hours
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
