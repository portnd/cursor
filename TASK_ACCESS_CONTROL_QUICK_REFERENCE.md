# 🔒 Task Access Control - Quick Reference

## 📖 Who Can Do What?

| Action | Creator | CEO | PM | Developer (Other) |
|--------|---------|-----|----|--------------------|
| ✅ Create Task | Yes | Yes | Yes | Yes |
| ✅ View Task | Yes | Yes | Yes | Yes |
| ✅ **Update Task** | **Yes** | **Yes** | **No** | **No** |
| ✅ **Delete Task** | **Yes** | **Yes** | **No** | **No** |

---

## 🔧 API Endpoints

### **Update Task**
```bash
PATCH /api/v1/sentinel/tasks/{taskId}
Authorization: Bearer <JWT_TOKEN>

{
  "title": "New title (optional)",
  "description": "New description (optional)"
}
```

**What Happens:**
1. ✅ Checks if you're the creator OR CEO
2. ✅ Compares new vs old title/description
3. ✅ If changed → **AI re-estimates automatically**
4. ✅ Resets negotiation status if AI re-estimates
5. ✅ Saves updated task

**Response:**
```json
{
  "message": "Task updated successfully",
  "data": {
    "ai_estimated_minutes": 240,  // NEW ESTIMATE
    "negotiation_status": "NONE"  // RESET
  }
}
```

---

### **Delete Task**
```bash
DELETE /api/v1/sentinel/tasks/{taskId}
Authorization: Bearer <JWT_TOKEN>
```

**What Happens:**
1. ✅ Checks if you're the creator OR CEO
2. ✅ Deletes the task permanently
3. ✅ Logs who deleted it

**Response:**
```json
{
  "message": "Task deleted successfully"
}
```

---

## 🤖 AI Re-estimation

### **When Does It Happen?**
AI automatically re-estimates when:
- ✅ Title changes
- ✅ Description changes
- ✅ Both change

### **What Gets Updated?**
- `ai_estimated_minutes` → New AI estimate
- `negotiation_status` → Reset to "NONE" (if was PENDING/APPROVED)
- `proposed_minutes` → Reset to 0
- `negotiation_reason` → Cleared

### **Example:**
```
Old Title: "Fix bug"
Old Estimate: 30 minutes

Update to:
New Title: "Implement secure database with encryption"
New Estimate: 240 minutes (AI calculated)
```

---

## ⚠️ Error Codes

### **403 Forbidden**
```json
{
  "error": "Forbidden",
  "message": "unauthorized: only the task creator or CEO can update this task"
}
```

**Reason:** You're not the creator and you're not CEO

**Solution:** Ask the creator or CEO to make the change

---

### **404 Not Found**
```json
{
  "error": "Not Found",
  "message": "task not found"
}
```

**Reason:** Task doesn't exist or was deleted

**Solution:** Check the task ID

---

### **400 Bad Request**
```json
{
  "error": "Bad Request",
  "message": "At least one field (title or description) must be provided"
}
```

**Reason:** You sent an empty update request

**Solution:** Provide at least title OR description

---

## 🧪 Quick Test Commands

### **Test Update (as Creator/CEO)**
```bash
CEO_TOKEN="your-token-here"
TASK_ID="task-uuid-here"

curl -X PATCH "http://localhost:8080/api/v1/sentinel/tasks/$TASK_ID" \
  -H "Authorization: Bearer $CEO_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title",
    "description": "Updated description"
  }'
```

---

### **Test Delete (as Creator/CEO)**
```bash
curl -X DELETE "http://localhost:8080/api/v1/sentinel/tasks/$TASK_ID" \
  -H "Authorization: Bearer $CEO_TOKEN"
```

---

## 💡 Best Practices

### **When Updating Tasks:**
1. ✅ **Only update what changed:** Don't send unchanged fields
2. ✅ **Expect AI re-estimation:** Task estimate will update automatically
3. ✅ **Check logs:** AI re-estimation logs show old vs new values
4. ✅ **Negotiation resets:** Any pending negotiation will be cleared

### **When Deleting Tasks:**
1. ⚠️ **Permanent action:** Cannot be undone (consider soft delete in future)
2. ⚠️ **Check dependencies:** Ensure no critical data is lost
3. ✅ **Logged action:** All deletions are logged for audit trail
4. ✅ **CEO override:** CEO can delete any task if needed

---

## 🔍 Troubleshooting

### **Problem: "Unauthorized" error when updating my own task**
**Solution:** Make sure you're logged in as the creator or CEO

### **Problem: AI re-estimation takes a long time**
**Solution:** This is normal (Gemini API call). Update will complete after AI responds.

### **Problem: AI estimation didn't update**
**Solution:** 
- Check if title/description actually changed
- If both are the same, AI won't re-estimate
- Check API logs for AI errors

---

## 📊 Monitoring

### **Check API Logs**
```bash
# See AI re-estimation activity
docker compose logs api | grep "AI Re-estimation"

# See who updated/deleted tasks
docker compose logs api | grep "Task Updated\|Task Deleted"
```

### **Log Format**
```
🔄 Task content changed. Triggering AI re-estimation...
   Old: [Old Title] Old description
   New: [New Title] New description
✅ AI Re-estimation Complete: 240 minutes (4.0 hours)
🔄 Resetting negotiation status (AI has new estimate)
✅ Task Updated: uuid by CEO (User ID: 1)
```

---

## 🎯 Common Scenarios

### **Scenario 1: CEO Updates Task**
```
1. CEO calls PATCH /tasks/:id
2. System checks: CEO? ✅ Yes
3. Title changed? ✅ Yes
4. AI re-estimates: 30min → 120min
5. Negotiation reset: PENDING → NONE
6. Response: "Task updated successfully"
```

### **Scenario 2: Developer Tries to Update Another's Task**
```
1. DEV calls PATCH /tasks/:id
2. System checks: Creator? ❌ No, CEO? ❌ No
3. Response: 403 Forbidden
```

### **Scenario 3: Creator Deletes Own Task**
```
1. Creator calls DELETE /tasks/:id
2. System checks: Creator? ✅ Yes
3. Task deleted from database
4. Response: "Task deleted successfully"
```

---

**🔒 Access Control is now active! Only authorized users can modify tasks.** 🚀
