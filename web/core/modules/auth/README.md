# 🔐 Authentication Module (Frontend)

## Overview

The Authentication Module implements user authentication on the frontend using **Pinia** for state management, **Nuxt 3** composables, and **TailwindCSS** for styling. This module follows the **Feature-Sliced Design (FSD)** architecture.

## Architecture

```
auth/
├── store/
│   └── auth-store.ts        # Pinia store (global state)
├── infrastructure/
│   └── auth-api.ts          # API integration layer
└── ui/
    ├── LoginForm.vue        # Login form component
    └── RegisterForm.vue     # Register form component
```

## Features

- ✅ User registration with email/password
- ✅ User login with JWT token
- ✅ Token storage in HTTP-only cookies
- ✅ Auto-redirect for authenticated users
- ✅ Beautiful Tailwind CSS forms
- ✅ Form validation (client-side)
- ✅ Loading states and error handling
- ✅ Auth middleware for protected routes

## Usage

### 1. Using Auth Store

```typescript
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

// In your component
const authStore = useAuthStore()

// Check if logged in
if (authStore.isLoggedIn) {
  console.log('User is logged in')
}

// Get user info
console.log(authStore.user)
console.log(authStore.userEmail)

// Login
await authStore.login('user@example.com', 'password123')

// Register
await authStore.register('user@example.com', 'password123', 'password123')

// Logout
await authStore.logout()
```

### 2. Protecting Routes with Middleware

Add middleware to any page that requires authentication:

```vue
<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})
</script>
```

Example - Protected Dashboard:

```vue
<!-- pages/dashboard.vue -->
<template>
  <div>
    <h1>Dashboard</h1>
    <p>Welcome, {{ authStore.userEmail }}</p>
    <button @click="authStore.logout()">Logout</button>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '~/core/modules/auth/store/auth-store'

definePageMeta({
  middleware: 'auth' // This page requires authentication
})

const authStore = useAuthStore()
</script>
```

### 3. Using Auth API Directly

If you need to call auth endpoints without using the store:

```typescript
import { authApi } from '~/core/modules/auth/infrastructure/auth-api'

// Login
const response = await authApi.login('user@example.com', 'password123')
console.log(response.data.token)

// Register
const response = await authApi.register('user@example.com', 'password123', 'password123')
console.log(response.data.user)
```

### 4. Using Components

Import and use the form components:

```vue
<template>
  <LoginForm />
  <!-- or -->
  <RegisterForm />
</template>

<script setup lang="ts">
import LoginForm from '~/core/modules/auth/ui/LoginForm.vue'
import RegisterForm from '~/core/modules/auth/ui/RegisterForm.vue'
</script>
```

## Pages

### Login Page (`/login`)

- Beautiful gradient background
- Email and password fields
- Form validation
- Error messages
- Link to register page
- Auto-redirect if already logged in

### Register Page (`/register`)

- Beautiful gradient background
- Email, password, and confirm password fields
- Form validation
- Password strength indicator
- Error messages
- Link to login page
- Auto-redirect if already logged in

## State Management

### Auth Store State

```typescript
interface AuthState {
  user: User | null           // Current user object
  token: string | null        // JWT token
  isAuthenticated: boolean    // Authentication status
  isLoading: boolean          // Loading state
  error: string | null        // Error message
}
```

### Getters

- `isLoggedIn` - Check if user is authenticated
- `userEmail` - Get current user's email
- `userId` - Get current user's ID

### Actions

- `initialize()` - Initialize auth from stored token (call on app start)
- `login(email, password)` - Login user
- `register(email, password, confirmPassword)` - Register new user
- `logout()` - Logout user and clear state
- `clearError()` - Clear error message

## Token Management

Tokens are stored in HTTP cookies with these settings:

```typescript
{
  maxAge: 60 * 60 * 24 * 7,  // 7 days
  path: '/',
  sameSite: 'lax'
}
```

### Why Cookies?

- More secure than localStorage
- Automatic inclusion in requests
- Can be HTTP-only (future enhancement)
- XSS protection

## Form Validation

### Client-Side Validation

Both forms include validation:

**Email:**
- Required
- Must be valid email format

**Password:**
- Required
- Minimum 8 characters

**Confirm Password (Register only):**
- Required
- Must match password

### Server-Side Validation

Backend performs additional validation:
- Email uniqueness
- Password strength
- SQL injection prevention

## Styling

Forms use TailwindCSS with:
- Gradient backgrounds
- Smooth transitions
- Focus states
- Error states
- Loading states
- Responsive design

## Error Handling

Errors are displayed in multiple ways:

1. **Field-level errors** - Red border and message below input
2. **Form-level errors** - Red alert box above submit button
3. **API errors** - Mapped from backend responses

Example error display:

```vue
<div v-if="authStore.error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
  <p class="text-sm">{{ authStore.error }}</p>
</div>
```

## API Integration

The `auth-api.ts` layer uses the shared `useHttp` composable:

```typescript
const { data, error } = await useHttp('/auth/login', {
  method: 'POST',
  body: { email, password }
})
```

This automatically:
- Adds base URL
- Includes auth token (if exists)
- Handles 401 errors
- Parses JSON responses

## Security Features

### Client-Side

- Input sanitization
- XSS prevention (Vue escapes by default)
- Token stored in cookies (not localStorage)
- CSRF protection (future: add CSRF token)

### Network

- HTTPS only in production
- CORS configured
- Token sent in Authorization header

## Testing

### Manual Testing

1. **Register New User**
   - Go to `/register`
   - Fill form with valid email/password
   - Submit and verify redirect to `/`

2. **Login Existing User**
   - Go to `/login`
   - Enter credentials
   - Submit and verify redirect to `/`

3. **Protected Route**
   - Logout
   - Try to access a protected page
   - Verify redirect to `/login`

4. **Logout**
   - Click logout button
   - Verify redirect to `/login`
   - Verify token cookie is cleared

### Example Test

```typescript
// Test login flow
const authStore = useAuthStore()

// Should be logged out initially
expect(authStore.isLoggedIn).toBe(false)

// Login
await authStore.login('test@example.com', 'password123')

// Should be logged in
expect(authStore.isLoggedIn).toBe(true)
expect(authStore.user).toBeTruthy()
expect(authStore.token).toBeTruthy()

// Logout
await authStore.logout()

// Should be logged out
expect(authStore.isLoggedIn).toBe(false)
```

## Future Enhancements

- [ ] Remember me functionality
- [ ] Password reset flow
- [ ] Email verification
- [ ] Social login (Google, Facebook)
- [ ] Two-factor authentication (2FA)
- [ ] Token refresh mechanism
- [ ] Persistent sessions
- [ ] Profile editing
- [ ] Password change
- [ ] Account deletion

## Troubleshooting

### "Failed to login" error

- Check API is running on `http://localhost:8080`
- Verify email and password are correct
- Check network tab for API errors
- Ensure JWT_SECRET is set in backend

### Token not persisting

- Check browser cookies (DevTools → Application → Cookies)
- Verify cookie settings (sameSite, path)
- Clear cookies and try again

### Redirect loop

- Check middleware is not applied to login/register pages
- Verify token validation logic

### Styling not applied

- Ensure TailwindCSS is configured
- Check `tailwind.config.ts` includes auth pages
- Run `npm run dev` with hot reload

## Contributing

When adding new auth features:

1. Update store if state changes needed
2. Update API layer for new endpoints
3. Create/update UI components
4. Update pages if needed
5. Update middleware for new auth logic
6. Update this README

---

**Created by:** Thanandorn (Komgrip CEO)  
**Last Updated:** 2026-01-22  
**Version:** 1.0.0
