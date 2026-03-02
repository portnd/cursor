#!/bin/bash

# 🧪 Dynamic AI Configuration Test Script

API_URL="http://localhost:8080"
CEO_TOKEN=""  # Will be set after login

echo "🧪 Testing Dynamic AI Configuration System"
echo "============================================"

# 1. Login as CEO to get token
echo ""
echo "1️⃣ Logging in as CEO..."
LOGIN_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ceo@komgrip.com",
    "password": "password123"
  }' \
  $API_URL/api/v1/auth/login)

CEO_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$CEO_TOKEN" ]; then
  echo "❌ Failed to get CEO token. Please ensure CEO user exists."
  echo "Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✅ CEO token obtained: ${CEO_TOKEN:0:20}..."

# 2. Get available models
echo ""
echo "2️⃣ Getting available Gemini models..."
curl -s -H "Authorization: Bearer $CEO_TOKEN" \
  $API_URL/api/v1/admin/models | jq '.'

# 3. Get current configuration
echo ""
echo "3️⃣ Getting current AI configuration..."
curl -s -H "Authorization: Bearer $CEO_TOKEN" \
  $API_URL/api/v1/admin/config | jq '.'

# 4. Update configuration
echo ""
echo "4️⃣ Updating AI configuration (CEO)..."
curl -s -X PUT \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 0.3,
    "cursor_assistance": 90
  }' \
  $API_URL/api/v1/admin/config | jq '.'

# 5. Verify update
echo ""
echo "5️⃣ Verifying configuration update..."
curl -s -H "Authorization: Bearer $CEO_TOKEN" \
  $API_URL/api/v1/admin/config | jq '.'

# 6. Test creating a task (should use new config)
echo ""
echo "6️⃣ Creating a test task (AI will use new config)..."
TASK_RESPONSE=$(curl -s -X POST \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Dynamic Config Task",
    "description": "This task should be estimated using temperature=0.3 and cursor_assistance=90%"
  }' \
  $API_URL/api/v1/sentinel/tasks)

echo $TASK_RESPONSE | jq '.'

# Extract task details
AI_MINUTES=$(echo $TASK_RESPONSE | grep -o '"ai_estimated_minutes":[0-9]*' | cut -d':' -f2)
echo ""
echo "📊 AI Estimated: $AI_MINUTES minutes ($(echo "scale=2; $AI_MINUTES/60" | bc) hours)"

# 7. Test invalid model
echo ""
echo "7️⃣ Testing validation (invalid model)..."
curl -s -X PUT \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-invalid-model",
    "temperature": 0.3,
    "cursor_assistance": 90
  }' \
  $API_URL/api/v1/admin/config | jq '.'

# 8. Test invalid temperature
echo ""
echo "8️⃣ Testing validation (invalid temperature)..."
curl -s -X PUT \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "active_model": "gemini-2.5-flash-lite",
    "temperature": 2.5,
    "cursor_assistance": 90
  }' \
  $API_URL/api/v1/admin/config | jq '.'

echo ""
echo "✅ All tests completed!"
echo ""
echo "🔍 Check API logs to see dynamic config in action:"
echo "   docker-compose logs -f api | grep 'AI Config'"
