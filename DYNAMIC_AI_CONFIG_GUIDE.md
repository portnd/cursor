# 🤖 Dynamic AI Configuration System

**Status:** ✅ IMPLEMENTED & DEPLOYED

## Overview

The Dynamic AI Configuration system allows CEOs to change AI behavior **at runtime** without restarting the server. Configuration changes take effect **immediately** on the next API call.

---

## 📋 System Configuration

### Configuration Parameters

| Parameter | Type | Range | Description |
|-----------|------|-------|-------------|
| `active_model` | string | See available models | Gemini model to use for all AI operations |
| `temperature` | float32 | 0.0 - 1.0 | AI creativity level (0.0 = deterministic, 1.0 = creative) |
| `cursor_assistance` | int | 0 - 100 | AI assistance level percentage (affects time estimation) |

### Available Gemini Models

```json
[
  "gemini-1.5-flash",
  "gemini-1.5-pro",
  "gemini-2.0-flash-exp",
  "gemini-2.5-flash-lite",
  "gemini-exp-1206"
]
```

### Temperature Guide

- **0.0 - 0.2**: Very deterministic, consistent outputs (recommended for time estimation)
- **0.3 - 0.5**: Balanced, slightly creative (recommended for code review)
- **0.6 - 0.8**: More creative, varied outputs (for brainstorming)
- **0.9 - 1.0**: Highly creative, unpredictable (experimental)

### Cursor Assistance Level

The `cursor_assistance` parameter affects how AI estimates time:

- **0-20%**: Minimal AI assistance → Slower estimates (traditional development)
- **21-50%**: Moderate AI assistance → Moderate time savings
- **51-80%**: Heavy AI assistance → Significant time reduction
- **81-100%**: AI-first workflow → Very aggressive estimates

---

## 🔌 API Endpoints

### 1. Get Current Configuration

**Endpoint:** `GET /api/v1/admin/config`

**Access:** CEO, PM, DEV (all authenticated users can view)

**Response:**
```json
{
  "message": "System configuration retrieved successfully",
  "data": {
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.4,
    "cursor_assistance": 80,
    "updated_at": "2026-01-26T06:51:17Z"
  }
}
```

**cURL Example:**
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/admin/config
```

---

### 2. Update Configuration (CEO Only)

**Endpoint:** `PUT /api/v1/admin/config`

**Access:** CEO only (403 Forbidden for PM/DEV)

**Request Body:**
```json
{
  "active_model": "gemini-2.5-flash-lite",
  "temperature": 0.3,
  "cursor_assistance": 90
}
```

**Validation:**
- `active_model`: Must be from available models list
- `temperature`: Must be 0.0 - 1.0
- `cursor_assistance`: Must be 0 - 100

**Response:**
```json
{
  "message": "System configuration updated successfully. Changes take effect immediately.",
  "data": {
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.3,
    "cursor_assistance": 90,
    "updated_at": "2026-01-26T13:52:00Z"
  }
}
```

**cURL Example:**
```bash
curl -X PUT \
  -H "Authorization: Bearer CEO_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.3,
    "cursor_assistance": 90
  }' \
  http://localhost:8080/api/v1/admin/config
```

**Error Responses:**

```json
// Non-CEO user
{
  "error": "Forbidden",
  "message": "access denied: only CEO can modify system configuration"
}

// Invalid temperature
{
  "error": "Bad Request",
  "message": "temperature must be between 0.0 and 1.0"
}

// Invalid model
{
  "error": "Bad Request",
  "message": "invalid model: gemini-invalid (must be one of: [gemini-1.5-flash gemini-1.5-pro ...])"
}
```

---

### 3. Get Available Models

**Endpoint:** `GET /api/v1/admin/models`

**Access:** All authenticated users

**Response:**
```json
{
  "message": "Available Gemini models",
  "data": [
    "gemini-1.5-flash",
    "gemini-1.5-pro",
    "gemini-2.0-flash-exp",
    "gemini-2.5-flash-lite",
    "gemini-exp-1206"
  ]
}
```

**cURL Example:**
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/admin/models
```

---

## 🎯 How It Works

### 1. Configuration Storage (Singleton Pattern)

- **Database Table:** `system_configs`
- **Single Record:** Always ID = 1
- **Auto-Creation:** If config doesn't exist, creates default on first access

**Default Configuration:**
```go
{
  ID:               1,
  ActiveModel:      "gemini-2.5-flash-lite",
  Temperature:      0.4,
  CursorAssistance: 80,
}
```

### 2. Dynamic AI Behavior

Every time AI is called, it:

1. **Fetches current config** from database
2. **Builds dynamic prompt** with Cursor Assistance context
3. **Uses active model** for API call
4. **Applies temperature** to generation config

**Example Flow:**

```
Task Created → EstimateEffort() Called
  ↓
Fetch SystemConfig from DB (cache potential here)
  ↓
Build Prompt with Cursor Assistance context:
  "The developer uses AI heavily (80% assistance)..."
  ↓
Call Gemini API with:
  - Model: gemini-2.5-flash-lite
  - Temperature: 0.4
  ↓
Return estimation
```

### 3. Cursor Assistance Impact

The system dynamically adjusts prompts based on assistance level:

**20% or less:**
```
"The developer is coding mostly manually with minimal AI assistance. 
Expect slower implementation times similar to traditional development."
```

**~50%:**
```
"The developer uses AI moderately for code suggestions and debugging. 
Estimate time with moderate AI speedup."
```

**~80%:**
```
"The developer heavily relies on AI tools for boilerplate, refactoring, and debugging. 
Expect significant time savings."
```

**90%+:**
```
"The developer works in an AI-first workflow with near-full assistance. 
Expect very aggressive time estimates - AI handles most implementation."
```

---

## 🔧 Implementation Details

### Database Schema

```sql
CREATE TABLE system_configs (
    id SERIAL PRIMARY KEY,
    active_model VARCHAR(255) DEFAULT 'gemini-2.5-flash-lite',
    temperature FLOAT DEFAULT 0.4,
    cursor_assistance INT DEFAULT 80,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Architecture

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                        │
│  GET/PUT /admin/config  |  GET /admin/models        │
└──────────────┬──────────────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────────────┐
│                  Usecase Layer                      │
│  - GetSystemConfig()                                │
│  - UpdateSystemConfig() [CEO Only]                  │
│  - GetAvailableModels()                             │
└──────────────┬──────────────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────────────┐
│                Repository Layer                     │
│  - GetSystemConfig() [Singleton ID=1]               │
│  - UpdateSystemConfig() [Force ID=1]                │
└──────────────┬──────────────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────────────┐
│                  AI Service Layer                   │
│  GeminiService.EstimateEffort():                    │
│    1. Fetch config from repo                        │
│    2. Build dynamic prompt with cursor context      │
│    3. Call Gemini API with config.ActiveModel       │
│    4. Use config.Temperature in generationConfig    │
└─────────────────────────────────────────────────────┘
```

### Key Files Modified

1. **Domain Layer:**
   - `domain/entities.go`: Added `SystemConfig` entity

2. **Repository Layer:**
   - `repository/postgres_repository.go`: Added config CRUD methods
   - `repository/gemini_service.go`: Injected repo, dynamic config fetching

3. **Usecase Layer:**
   - `usecase/sentinel_usecase.go`: Config management with CEO validation

4. **HTTP Layer:**
   - `delivery/http/sentinel_handler.go`: Admin endpoint handlers
   - `delivery/http/route.go`: Registered `/admin/*` routes

5. **Main:**
   - `cmd/server/main.go`: Updated DI, added SystemConfig migration

---

## 🧪 Testing Scenarios

### Scenario 1: Change Model to Test Rate Limits

```bash
# Switch to lighter model to avoid rate limits
curl -X PUT \
  -H "Authorization: Bearer CEO_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-1.5-flash",
    "temperature": 0.2,
    "cursor_assistance": 80
  }' \
  http://localhost:8080/api/v1/admin/config
```

### Scenario 2: Adjust Temperature for Consistency

```bash
# Lower temperature for more deterministic estimates
curl -X PUT \
  -H "Authorization: Bearer CEO_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.1,
    "cursor_assistance": 80
  }' \
  http://localhost:8080/api/v1/admin/config
```

### Scenario 3: Change Cursor Assistance Level

```bash
# Team is now AI-first, increase assistance to 95%
curl -X PUT \
  -H "Authorization: Bearer CEO_JWT" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.4,
    "cursor_assistance": 95
  }' \
  http://localhost:8080/api/v1/admin/config

# Next task creation will use "AI-first workflow" context
# Resulting in more aggressive time estimates
```

---

## 🚀 Production Recommendations

### Optimal Settings by Use Case

**For Consistent Time Estimates:**
```json
{
  "active_model": "gemini-2.5-flash-lite",
  "temperature": 0.1,
  "cursor_assistance": 80
}
```

**For High Throughput (Avoiding Rate Limits):**
```json
{
  "active_model": "gemini-1.5-flash",
  "temperature": 0.3,
  "cursor_assistance": 80
}
```

**For Highest Quality (Cost be damned):**
```json
{
  "active_model": "gemini-1.5-pro",
  "temperature": 0.2,
  "cursor_assistance": 70
}
```

**For Experimental/Latest Features:**
```json
{
  "active_model": "gemini-exp-1206",
  "temperature": 0.4,
  "cursor_assistance": 80
}
```

---

## 📊 Performance Considerations

### Caching Strategy (Future Enhancement)

Currently, config is fetched on **every AI call**. For optimization:

```go
// Option 1: In-memory cache with TTL
var configCache *SystemConfig
var cacheExpiry time.Time

// Option 2: Redis cache
redisClient.Set("ai_config", config, 5*time.Minute)

// Option 3: Event-driven invalidation
// Pub/Sub: When CEO updates config, broadcast invalidation
```

### Database Query Optimization

```sql
-- Index on singleton ID (though not needed for single row)
CREATE INDEX idx_system_configs_id ON system_configs(id);

-- Consider UNLOGGED table for even faster reads (data not critical)
ALTER TABLE system_configs SET UNLOGGED;
```

---

## 🔒 Security Notes

1. **CEO-Only Updates:** Only CEO role can modify configuration
2. **All Users Can View:** Transparency - everyone sees AI settings
3. **Validation:** All inputs validated before DB update
4. **No Secrets Exposed:** API keys remain in environment variables

---

## 📝 Future Enhancements

1. **Per-User AI Settings:** Allow users to override system defaults
2. **Model Cost Tracking:** Log cost per model/request
3. **A/B Testing:** Test different configs on subset of tasks
4. **Config History:** Audit log of all config changes
5. **UI Dashboard:** Visual config editor for CEO
6. **Auto-Tuning:** ML-based optimization of temperature/assistance

---

## ✅ Summary

✅ Dynamic AI Configuration fully implemented
✅ CEO can change model, temperature, cursor assistance at runtime
✅ Changes take effect immediately (no restart needed)
✅ Comprehensive validation and error handling
✅ Access control: CEO-only updates, all can view
✅ Production-ready with sensible defaults

**Default Configuration:**
- Model: `gemini-2.5-flash-lite` (balanced speed/quality)
- Temperature: `0.4` (slightly creative but stable)
- Cursor Assistance: `80%` (heavy AI-assisted workflow)
