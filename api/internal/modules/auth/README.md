# 🔐 Authentication Module

## Overview

The Authentication Module implements user registration and login functionality using JWT (JSON Web Tokens) for session management. This module follows the **Hexagonal Architecture** pattern with clear separation of concerns.

## Architecture

```
auth/
├── domain/
│   └── entity.go           # Core entities, DTOs, and interfaces (PORT)
├── repository/
│   └── postgres_repo.go    # PostgreSQL implementation (ADAPTER)
├── usecase/
│   └── auth_usecase.go     # Business logic (CORE)
└── delivery/
    └── http/
        ├── auth_handler.go # HTTP handlers (ADAPTER)
        └── route.go        # Route registration
```

## Features

- ✅ User Registration with email/password
- ✅ User Login with JWT token generation
- ✅ Password hashing using bcrypt (cost factor: 12)
- ✅ JWT token generation (HS256, 24-hour expiration)
- ✅ Email uniqueness validation
- ✅ Input validation using Gin binding

## Environment Variables

Add this to your `.env` file:

```env
# JWT Authentication
# Generate a strong random secret using: openssl rand -base64 32
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

**⚠️ CRITICAL:** Never commit your actual `JWT_SECRET` to version control!

## API Endpoints

### 1. Register User

**POST** `/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "confirm_password": "password123"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "created_at": "2026-01-22T10:00:00Z",
      "updated_at": "2026-01-22T10:00:00Z"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation errors
- `409 Conflict` - Email already registered
- `500 Internal Server Error` - Server error

### 2. Login User

**POST** `/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "created_at": "2026-01-22T10:00:00Z",
      "updated_at": "2026-01-22T10:00:00Z"
    }
  }
}
```

**Error Responses:**
- `400 Bad Request` - Validation errors
- `401 Unauthorized` - Invalid credentials
- `500 Internal Server Error` - Server error

## Security Features

### Password Hashing
- Uses `bcrypt` with cost factor 12
- Passwords are never stored in plain text
- Hash comparison happens in constant time to prevent timing attacks

### JWT Tokens
- Algorithm: HS256 (HMAC with SHA-256)
- Expiration: 24 hours from issuance
- Claims include: `user_id`, `email`, `iat`, `exp`
- Secret key loaded from environment variable

### Input Validation
- Email format validation
- Password minimum length (8 characters)
- Password confirmation matching
- SQL injection prevention (GORM parameterized queries)

## Database Schema

### `users` table

| Column | Type | Constraints |
|:-------|:-----|:------------|
| `id` | SERIAL | PRIMARY KEY |
| `email` | VARCHAR | UNIQUE, NOT NULL |
| `password` | VARCHAR | NOT NULL (hashed) |
| `created_at` | TIMESTAMP | AUTO |
| `updated_at` | TIMESTAMP | AUTO |

## Testing

### Manual Testing with cURL

**Register:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Using the Token

Save the token from the response and use it in subsequent requests:

```bash
TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8080/protected-endpoint \
  -H "Authorization: Bearer $TOKEN"
```

## Dependencies

- `github.com/golang-jwt/jwt/v5` - JWT token generation and validation
- `golang.org/x/crypto/bcrypt` - Password hashing
- `gorm.io/gorm` - Database ORM
- `github.com/gin-gonic/gin` - HTTP framework

## Future Enhancements

- [ ] Email verification
- [ ] Password reset functionality
- [ ] Refresh token mechanism
- [ ] OAuth2 integration (Google, Facebook, etc.)
- [ ] Two-factor authentication (2FA)
- [ ] Rate limiting for login attempts
- [ ] Account lockout after failed attempts
- [ ] Password strength validation
- [ ] Session management (logout all devices)

## Troubleshooting

### "email already registered" error
- The email is already in the database
- Try logging in instead, or use a different email

### "invalid email or password" error
- Check email spelling
- Verify password is correct
- Ensure user exists (register first)

### "failed to generate token" error
- Check `JWT_SECRET` environment variable is set
- Ensure secret is not empty

## Contributing

When adding new authentication features:
1. Update domain interfaces first
2. Implement in usecase layer (business logic)
3. Update repository if data access changes
4. Update handlers for new endpoints
5. Update this README with new endpoints

---

**Created by:** Thanandorn (Komgrip CEO)  
**Last Updated:** 2026-01-22  
**Version:** 1.0.0
