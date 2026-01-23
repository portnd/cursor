# 🧪 Authentication Module - Testing Guide

Quick guide to test the newly implemented authentication system.

---

## 🚀 Quick Start Testing (5 minutes)

### Step 1: Ensure Services are Running

```bash
# Check if services are running
make ps

# If not running, start them
make up

# Wait 30 seconds for services to initialize
sleep 30

# Check API health
curl http://localhost:8080/health
```

**Expected output:**
```json
{
  "status": "UP",
  "services": {
    "postgres": "UP",
    "mongodb": "UP",
    "redis": "UP"
  }
}
```

---

## 🔧 Backend API Testing

### Test 1: Register a New User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@komgrip.com",
    "password": "password123",
    "confirm_password": "password123"
  }'
```

**Expected output:**
```json
{
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "test@komgrip.com",
      "created_at": "2026-01-22T..."
    }
  }
}
```

**✅ Success Criteria:**
- Status code: 201 Created
- Response contains `token` and `user` object
- User `id` is present
- Password is NOT in response

### Test 2: Login with Existing User

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@komgrip.com",
    "password": "password123"
  }'
```

**Expected output:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "test@komgrip.com",
      "created_at": "2026-01-22T..."
    }
  }
}
```

**✅ Success Criteria:**
- Status code: 200 OK
- Response contains `token` and `user` object
- Same user `id` as registration

### Test 3: Duplicate Email (Should Fail)

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@komgrip.com",
    "password": "another123",
    "confirm_password": "another123"
  }'
```

**Expected output:**
```json
{
  "error": "Conflict",
  "message": "email already registered"
}
```

**✅ Success Criteria:**
- Status code: 409 Conflict
- Error message indicates email is taken

### Test 4: Wrong Password (Should Fail)

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@komgrip.com",
    "password": "wrongpassword"
  }'
```

**Expected output:**
```json
{
  "error": "Authentication failed",
  "message": "Invalid email or password"
}
```

**✅ Success Criteria:**
- Status code: 401 Unauthorized
- Generic error message (doesn't reveal if email exists)

### Test 5: Verify Password is Hashed in Database

```bash
# Connect to PostgreSQL
make shell-db

# In psql shell, run:
# SELECT id, email, LEFT(password, 20) AS password_hash FROM users;
# \q
```

**Expected output:**
```
 id |       email        |   password_hash
----+-------------------+---------------------
  1 | test@komgrip.com  | $2a$12$...
```

**✅ Success Criteria:**
- Password starts with `$2a$12$` (bcrypt hash)
- Password is NOT plain text

---

## 🌐 Frontend Testing

### Test 1: Open Register Page

```bash
# Mac
open http://localhost:3000/register

# Linux
xdg-open http://localhost:3000/register

# Or manually open in browser
```

**✅ Visual Checklist:**
- [ ] Beautiful gradient purple/indigo background
- [ ] "🛡️ KOMGRIP" logo at top
- [ ] "Create Account" heading
- [ ] Email input field
- [ ] Password input field
- [ ] Confirm Password input field
- [ ] "Create Account" button (gradient purple)
- [ ] "Already have an account? Sign in" link
- [ ] Responsive design (try resizing)

### Test 2: Register a New User (Frontend)

1. Fill in the form:
   - Email: `frontend@komgrip.com`
   - Password: `frontend123`
   - Confirm Password: `frontend123`

2. Click "Create Account"

**✅ Success Checklist:**
- [ ] Button shows "Creating account..." with spinner
- [ ] After success, redirected to home page (`http://localhost:3000/`)
- [ ] Home page shows "Welcome, frontend@komgrip.com"
- [ ] "Logout" button is visible (red gradient)

### Test 3: Test Form Validation

1. Clear all fields
2. Try to submit

**✅ Expected Behavior:**
- [ ] Email field shows "Email is required" error
- [ ] Password field shows "Password is required" error
- [ ] Red borders on invalid fields

3. Enter invalid email: `notanemail`

**✅ Expected Behavior:**
- [ ] "Email is invalid" error message

4. Enter password: `123` (too short)

**✅ Expected Behavior:**
- [ ] "Password must be at least 8 characters" error

5. Enter mismatched passwords:
   - Password: `password123`
   - Confirm: `password456`

**✅ Expected Behavior:**
- [ ] "Passwords do not match" error

### Test 4: Login Page

1. Go to http://localhost:3000/login

**✅ Visual Checklist:**
- [ ] Similar design to register page
- [ ] "Welcome Back" heading
- [ ] Email and Password fields only (no confirm)
- [ ] "Sign In" button
- [ ] "Don't have an account? Sign up" link

2. Login with:
   - Email: `frontend@komgrip.com`
   - Password: `frontend123`

**✅ Success Checklist:**
- [ ] Button shows "Signing in..." with spinner
- [ ] Redirected to home page
- [ ] Welcome message displays

### Test 5: Logout

1. On home page, click "Logout" button

**✅ Success Checklist:**
- [ ] Redirected to `/login` page
- [ ] Token cookie is cleared (check DevTools → Application → Cookies)
- [ ] Try going to home page - should redirect to login

### Test 6: Auth Middleware (Protected Route)

1. Logout (if logged in)
2. Try to access home page directly: http://localhost:3000/

**Current Behavior:**
- [ ] Home page loads (not protected yet)

**To protect a route**, add this to any page:

```vue
<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})
</script>
```

Then:
- [ ] Without login, redirects to `/login`
- [ ] With login, page loads normally

### Test 7: Token Persistence

1. Login
2. Close browser completely
3. Open browser again
4. Go to http://localhost:3000/

**✅ Expected Behavior:**
- [ ] Still logged in (token persists in cookie)
- [ ] Welcome message still shows
- [ ] Don't need to login again

---

## 🎨 UI/UX Testing

### Responsive Design

Test on different screen sizes:

**Desktop (1920x1080):**
- [ ] Form is centered
- [ ] Background gradients visible
- [ ] Proper spacing

**Tablet (768px):**
- [ ] Form still readable
- [ ] Buttons still accessible

**Mobile (375px):**
- [ ] Form scales down
- [ ] No horizontal scroll
- [ ] Touch-friendly buttons

### Loading States

1. Slow down network (DevTools → Network → Slow 3G)
2. Try to register/login

**✅ Expected:**
- [ ] Spinner appears in button
- [ ] Button text changes to "Creating account..." or "Signing in..."
- [ ] Button is disabled during loading
- [ ] Form inputs are disabled

### Error Messages

Test various error scenarios:

**Backend Error (e.g., email taken):**
- [ ] Red alert box appears above submit button
- [ ] Error message is clear and helpful

**Network Error:**
- [ ] Error message appears
- [ ] User can retry

---

## 🔒 Security Testing

### Test 1: XSS Protection

Try to register with:
- Email: `<script>alert('xss')</script>@test.com`

**✅ Expected:**
- [ ] Script does NOT execute
- [ ] Email is treated as text
- [ ] Vue auto-escapes the input

### Test 2: SQL Injection

Try to login with:
- Email: `admin' OR '1'='1`
- Password: `anything`

**✅ Expected:**
- [ ] Login fails
- [ ] No database error
- [ ] GORM prevents SQL injection

### Test 3: Token Security

1. Login and get token
2. Open DevTools → Application → Cookies
3. Find `token` cookie

**✅ Checklist:**
- [ ] Cookie exists with JWT value
- [ ] `SameSite: Lax` (prevents CSRF)
- [ ] `Path: /` (available everywhere)
- [ ] Expires in 7 days

4. Try to use token in API call:

```bash
TOKEN="your-token-here"

curl -X GET http://localhost:8080/protected-endpoint \
  -H "Authorization: Bearer $TOKEN"
```

(Note: Protected endpoint doesn't exist yet, but token format is valid)

---

## 🐛 Common Issues & Solutions

### Issue 1: "Failed to connect to API"

**Symptoms:**
- Frontend shows connection error
- Network tab shows failed requests

**Solutions:**
```bash
# Check API is running
curl http://localhost:8080/health

# If not, check logs
make logs

# Restart API
docker-compose restart api
```

### Issue 2: "Database migration failed"

**Symptoms:**
- API logs show migration error
- Users table not created

**Solutions:**
```bash
# Check PostgreSQL is healthy
make ps

# Connect to database
make shell-db

# Check if users table exists
\dt

# If not, check API logs for errors
docker-compose logs api | grep -i migrate
```

### Issue 3: "Token not persisting"

**Symptoms:**
- Login successful but immediately logged out
- Cookie not saved

**Solutions:**
1. Check browser cookies (DevTools → Application)
2. Try incognito mode (no extensions)
3. Check browser console for errors
4. Verify cookie settings in code

### Issue 4: "Passwords don't match" even when they do

**Symptoms:**
- Form shows error incorrectly

**Solutions:**
1. Check for trailing spaces
2. Clear form and try again
3. Check browser console for JS errors

---

## ✅ Complete Test Checklist

### Backend
- [ ] Register endpoint works
- [ ] Login endpoint works
- [ ] Passwords are hashed in database
- [ ] JWT tokens are generated
- [ ] Duplicate email is rejected
- [ ] Wrong password is rejected
- [ ] Email validation works
- [ ] Password minimum length enforced

### Frontend
- [ ] Register page renders correctly
- [ ] Login page renders correctly
- [ ] Form validation works
- [ ] Error messages display
- [ ] Loading states show
- [ ] Success redirects to home
- [ ] Token stored in cookie
- [ ] Logout clears token
- [ ] Token persists across sessions
- [ ] Responsive design works

### Integration
- [ ] Frontend can call backend API
- [ ] CORS allows requests
- [ ] Tokens work across requests
- [ ] Auto-redirect works

### Security
- [ ] XSS protection works
- [ ] SQL injection prevented
- [ ] Passwords never exposed
- [ ] Generic error messages

---

## 📊 Performance Testing

### Load Test (Optional)

Use `wrk` or `ab` (ApacheBench):

```bash
# Install wrk (Mac)
brew install wrk

# Test register endpoint (100 requests, 10 concurrent)
wrk -t10 -c10 -d10s --script register.lua http://localhost:8080/auth/register
```

**Expected:**
- [ ] Response time < 200ms
- [ ] No errors under load
- [ ] Database handles concurrent writes

---

## 🎉 Success Criteria

**You're done when:**

1. ✅ All backend tests pass
2. ✅ All frontend tests pass
3. ✅ All security tests pass
4. ✅ UI looks beautiful
5. ✅ No errors in logs
6. ✅ Token persists correctly
7. ✅ Users can register, login, logout smoothly

---

## 📞 Need Help?

If any tests fail:

1. Check logs: `make logs`
2. Check specific service: `docker-compose logs api` or `docker-compose logs web`
3. Read error messages carefully
4. Check documentation: `AUTH_MODULE.md`
5. Create GitHub issue with:
   - What you were testing
   - Expected vs actual behavior
   - Error logs
   - Screenshots

---

**Happy Testing! 🚀**

Built by Thanandorn (Komgrip CEO)  
Date: 2026-01-22  
Version: 1.0.0
