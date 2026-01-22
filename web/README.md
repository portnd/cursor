# 🛡️ KOMGRIP WEB (Nuxt 3 Frontend)

Production-ready Nuxt 3 application built with Feature-Sliced Design (FSD).

---

## 📁 Project Structure (FSD)

```
web/
├── core/                           # Core business logic
│   ├── modules/                    # Feature modules (FSD)
│   │   └── {feature}/
│   │       ├── infrastructure/     # API clients, external services
│   │       ├── store/              # Pinia stores
│   │       └── ui/                 # Feature-specific components
│   └── shared/                     # Shared utilities
│       ├── api/                    # HTTP client
│       │   └── http.ts             # useHttp composable
│       ├── composables/            # Reusable composables
│       └── ui/                     # Common UI components
├── pages/                          # File-based routing (Nuxt)
│   └── index.vue                   # Landing page
├── layouts/                        # Layout components
├── assets/
│   └── css/
│       └── main.css                # Global styles + Tailwind
├── public/                         # Static assets
├── docker/
│   └── Dockerfile.dev              # Development container
├── app.vue                         # Root component
├── nuxt.config.ts                  # Nuxt configuration
├── tailwind.config.ts              # Tailwind configuration
├── tsconfig.json                   # TypeScript configuration
├── package.json                    # Dependencies
└── README.md                       # This file
```

---

## 🏗️ Architecture: Feature-Sliced Design (FSD)

### Layer Structure

```
core/modules/{feature}/
├── infrastructure/     # External dependencies (API calls, localStorage, etc.)
├── store/              # State management (Pinia)
└── ui/                 # UI components specific to this feature
```

### Dependency Rules
- **Pages** → **Modules** → **Shared**
- Pages are dumb components (no business logic)
- Modules contain all feature logic
- Shared contains reusable utilities

---

## 🚀 Getting Started

### 1. Start Services (From Root Directory)

```bash
# From repository root
make up

# View Web logs
docker-compose logs -f web

# Access Web container shell
make shell-web
```

### 2. Access Application

- **Development:** http://localhost:3000
- **API:** http://localhost:8080

### 3. Hot Module Replacement (HMR)

Edit any `.vue`, `.ts`, or `.css` file and see changes instantly!

---

## 🛠️ Development

### Install Dependencies

```bash
# Inside container
make shell-web
npm install

# Or from host (if Node.js installed)
cd web
npm install
```

### Run Development Server

```bash
npm run dev
```

The app will be available at http://localhost:3000 with hot reload enabled.

### Build for Production

```bash
npm run build
```

### Preview Production Build

```bash
npm run preview
```

### Type Check

```bash
npm run type-check
```

### Lint Code

```bash
npm run lint
npm run lint:fix
```

---

## 📦 Dependencies

| Package | Purpose | Documentation |
| :--- | :--- | :--- |
| `nuxt` | Vue.js framework | https://nuxt.com |
| `@pinia/nuxt` | State management | https://pinia.vuejs.org |
| `@nuxtjs/tailwindcss` | Utility-first CSS | https://tailwindcss.nuxtjs.org |
| `typescript` | Type safety | https://www.typescriptlang.org |
| `sass` | CSS preprocessor | https://sass-lang.com |

---

## 🌐 HTTP Client (`useHttp`)

### Usage

```typescript
// In any component or composable
const { get, post, put, patch, delete: del } = useHttp()

// GET request
const { data, error, pending } = await get<User>('/api/users/1')

// POST request
const { data } = await post<User>('/api/users', {
  name: 'John Doe',
  email: 'john@example.com'
})

// PUT request
await put('/api/users/1', { name: 'Jane Doe' })

// DELETE request
await del('/api/users/1')
```

### Features

- ✅ Automatic base URL from `runtimeConfig`
- ✅ Automatic token injection (from `auth_token` cookie)
- ✅ 401 handling (redirects to `/login`)
- ✅ Type-safe with TypeScript generics
- ✅ Built on Nuxt's `useFetch` (SSR-friendly)

---

## 📝 Adding a New Feature Module

Example: Adding an `auth` module

### 1. Create Directory Structure

```bash
mkdir -p core/modules/auth/{infrastructure,store,ui}
```

### 2. Create API Client (Infrastructure)

**`core/modules/auth/infrastructure/auth-api.ts`**

```typescript
import { useHttp } from '~/core/shared/api/http'

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: number
    email: string
    name: string
  }
}

export const useAuthApi = () => {
  const { post, get } = useHttp()

  const login = async (credentials: LoginRequest) => {
    return await post<LoginResponse>('/api/auth/login', credentials)
  }

  const logout = async () => {
    return await post('/api/auth/logout', {})
  }

  const getProfile = async () => {
    return await get('/api/auth/profile')
  }

  return {
    login,
    logout,
    getProfile
  }
}
```

### 3. Create Pinia Store

**`core/modules/auth/store/auth-store.ts`**

```typescript
import { defineStore } from 'pinia'
import { useAuthApi } from '../infrastructure/auth-api'
import type { LoginRequest } from '../infrastructure/auth-api'

interface User {
  id: number
  email: string
  name: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = useCookie('auth_token')
  const api = useAuthApi()

  const isAuthenticated = computed(() => !!user.value && !!token.value)

  const login = async (credentials: LoginRequest) => {
    const { data, error } = await api.login(credentials)
    
    if (error) {
      throw new Error('Login failed')
    }

    if (data) {
      token.value = data.token
      user.value = data.user
    }
  }

  const logout = async () => {
    await api.logout()
    token.value = null
    user.value = null
    navigateTo('/login')
  }

  const fetchProfile = async () => {
    if (!token.value) return

    const { data } = await api.getProfile()
    if (data) {
      user.value = data
    }
  }

  return {
    user,
    isAuthenticated,
    login,
    logout,
    fetchProfile
  }
})
```

### 4. Create UI Components

**`core/modules/auth/ui/LoginForm.vue`**

```vue
<template>
  <div class="card max-w-md mx-auto">
    <h2 class="text-2xl font-bold mb-6">Login</h2>
    
    <form @submit.prevent="handleSubmit">
      <div class="mb-4">
        <label class="block text-sm font-medium mb-2">Email</label>
        <input 
          v-model="form.email" 
          type="email" 
          class="input" 
          required 
        />
      </div>

      <div class="mb-6">
        <label class="block text-sm font-medium mb-2">Password</label>
        <input 
          v-model="form.password" 
          type="password" 
          class="input" 
          required 
        />
      </div>

      <button type="submit" class="btn-primary w-full" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
    </form>

    <p v-if="error" class="text-red-500 text-sm mt-4">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '../store/auth-store'

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(form)
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
```

### 5. Create Page

**`pages/login.vue`**

```vue
<template>
  <div class="min-h-screen flex items-center justify-center">
    <LoginForm />
  </div>
</template>

<script setup lang="ts">
import LoginForm from '~/core/modules/auth/ui/LoginForm.vue'

definePageMeta({
  layout: false
})
</script>
```

---

## 🎨 Styling with TailwindCSS

### Utility Classes

```vue
<template>
  <!-- Buttons -->
  <button class="btn-primary">Primary Button</button>
  <button class="btn-secondary">Secondary Button</button>

  <!-- Cards -->
  <div class="card">
    <h3 class="text-xl font-bold">Card Title</h3>
    <p>Card content</p>
  </div>

  <!-- Inputs -->
  <input type="text" class="input" placeholder="Enter text" />

  <!-- Gradient Text -->
  <h1 class="text-gradient">Gradient Heading</h1>
</template>
```

### Custom Colors

```typescript
// Available in tailwind.config.ts
'primary-50' to 'primary-950'  // Purple gradient
```

---

## 🚀 Production Deployment

### Build Docker Image

```bash
docker build -t komgrip-web:latest -f Dockerfile .
```

### Run Production Container

```bash
docker run -d \
  -p 3000:3000 \
  -e NUXT_PUBLIC_API_BASE=https://api.yourapp.com \
  komgrip-web:latest
```

### Environment Variables

- `NUXT_PUBLIC_API_BASE`: Backend API URL (default: `http://localhost:8080`)
- `NODE_ENV`: `development` or `production`

---

## 📖 Best Practices

### 1. Component Organization
- Keep pages simple (routing only)
- Move logic to modules
- Use composables for reusable logic

### 2. State Management
- Use Pinia for complex state
- Use `ref`/`reactive` for local state
- Avoid prop drilling (use composables)

### 3. Type Safety
- Always define TypeScript interfaces
- Use generics for API responses
- Enable strict mode

### 4. Performance
- Use `lazy` loading for heavy components
- Implement virtual scrolling for long lists
- Optimize images with Nuxt Image

### 5. SEO
- Use `useHead()` for meta tags
- Implement proper OpenGraph tags
- Enable SSR for public pages

---

## 🧪 Testing

### Unit Tests (Coming Soon)

```bash
npm run test
```

### E2E Tests (Coming Soon)

```bash
npm run test:e2e
```

---

## 📞 Support

For issues or questions:
- **CEO:** Thanandorn
- **Architecture:** Feature-Sliced Design (FSD)
- **Documentation:** See root README.md

---

**Built with 💪 for national-scale applications.**
