# ⚡ AI Configuration Quick Reference

## 🔌 Endpoints

```
GET  /api/v1/admin/config   → Get current config (all users)
PUT  /api/v1/admin/config   → Update config (CEO only)
GET  /api/v1/admin/models   → List available models (all users)
```

## 📋 Parameters

| Field | Type | Range | Default |
|-------|------|-------|---------|
| `active_model` | string | See models | `gemini-2.5-flash-lite` |
| `temperature` | float32 | 0.0 - 1.0 | `0.4` |
| `cursor_assistance` | int | 0 - 100 | `80` |

## 🤖 Available Models

```
gemini-1.5-flash
gemini-1.5-pro
gemini-2.0-flash-exp
gemini-2.5-flash-lite
gemini-exp-1206
```

## 🎯 Quick Examples

### Get Config
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/admin/config
```

### Update Config (CEO)
```bash
curl -X PUT \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.3,
    "cursor_assistance": 90
  }' \
  http://localhost:8080/api/v1/admin/config
```

## 🎨 Temperature Guide

| Value | Behavior |
|-------|----------|
| 0.0-0.2 | Very deterministic (recommended) |
| 0.3-0.5 | Balanced creativity |
| 0.6-1.0 | High creativity (experimental) |

## 🚀 Cursor Assistance Impact

| Level | AI Workflow | Time Estimates |
|-------|-------------|----------------|
| 0-20% | Manual coding | Slower (traditional) |
| 21-50% | Moderate AI use | Moderate savings |
| 51-80% | Heavy AI reliance | Significant reduction |
| 81-100% | AI-first workflow | Very aggressive |

## 🔍 Check Current Config

```bash
# Via API
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/admin/config | jq '.data'

# Via Database
docker-compose exec postgres psql -U komgrip -d komgrip_db \
  -c "SELECT * FROM system_configs;"
```

## 📊 Monitor in Real-Time

```bash
# Watch AI config being used
docker-compose logs -f api | grep "AI Config"

# Output:
# ⚙️  AI Config: Model=gemini-2.5-flash-lite, Temp=0.40, Cursor=80%
# 📡 Calling Gemini API (model: gemini-2.5-flash-lite, temp: 0.40, cursor: 80%)
```

## ⚠️ Common Errors

```json
// Non-CEO trying to update
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
  "message": "invalid model: xyz (must be one of: [...])"
}
```

## 🧪 Test Script

```bash
./test_dynamic_config.sh
```

## 📚 Full Documentation

- **Complete Guide:** `DYNAMIC_AI_CONFIG_GUIDE.md`
- **Summary:** `DYNAMIC_AI_CONFIG_SUMMARY.md`
- **This Quick Ref:** `AI_CONFIG_QUICK_REF.md`
