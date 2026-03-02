# 🚦 Task Approval - Quick Reference Card

## ⚡ Quick Facts

```
New Status:  REVIEW_PENDING 🚦
New Endpoint: POST /tasks/:id/approve
Permission:  PM or CEO only
Impact:      CRITICAL - All tasks need approval
```

---

## 🔄 Flow Change

### **BEFORE:**
```
Submit → AI PASS → ✅ COMPLETED (automatic)
```

### **AFTER:**
```
Submit → AI PASS → 🚦 REVIEW_PENDING → Approve → ✅ COMPLETED
```

---

## 🎯 New Endpoint

### **Approve Task:**
```bash
POST /api/v1/sentinel/tasks/:id/approve
Authorization: Bearer <pm_or_ceo_token>
```

**Success (200):**
```json
{
  "message": "Task approved and marked as COMPLETED"
}
```

**Forbidden (403):**
```json
{
  "error": "Forbidden",
  "message": "access denied: only PM or CEO can approve tasks"
}
```

---

## 📋 Status States

| Status | Meaning | Next Action |
|--------|---------|-------------|
| **PENDING** | AI review failed | Fix code or appeal |
| **REVIEW_PENDING** 🚦 | AI passed, awaiting approval | PM/CEO approves |
| **COMPLETED** ✅ | Approved by PM/CEO | Done |

---

## 👥 Permissions

| Role | Can Approve? |
|------|--------------|
| Developer | ❌ No |
| PM | ✅ Yes |
| CEO | ✅ Yes |

---

## 🧪 Quick Test

### **1. Submit Code (DEV):**
```bash
curl -X POST /api/v1/sentinel/tasks/{id}/submit \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -d '{"commit_hash": "abc123", "diff": "..."}'
```
**Result:** Status = REVIEW_PENDING

---

### **2. Check Approvals (PM/CEO):**
```bash
curl -X GET /api/v1/sentinel/tasks/approvals \
  -H "Authorization: Bearer $CEO_TOKEN"
```
**Result:** List of pending tasks

---

### **3. Approve (PM/CEO):**
```bash
curl -X POST /api/v1/sentinel/tasks/{id}/approve \
  -H "Authorization: Bearer $CEO_TOKEN"
```
**Result:** Status = COMPLETED ✅

---

## 📊 Logs

### **Submit Work:**
```
🚦 Task abc-123 moved to REVIEW_PENDING - awaiting PM/CEO approval
📋 AI Review: PASS (90/100) - Ready for human verification
```

### **Approve:**
```
✅ Task abc-123 APPROVED by CEO (ID: 1)
🎉 Task marked COMPLETED at 2026-01-26T10:46:12Z
📊 Actual Time: 3.5 hours (AI Estimated: 4.0 hours)
```

---

## ⚠️ Common Errors

### **Error 1: Not PM/CEO**
```
403 Forbidden
"access denied: only PM or CEO can approve tasks (your role: DEV)"
```

### **Error 2: Wrong Status**
```
400 Bad Request
"task is not pending review (current status: COMPLETED)"
```

### **Error 3: Not Found**
```
404 Not Found
"task not found"
```

---

## 🎓 Key Points

1. ✅ **AI PASS ≠ Completed** anymore
2. ✅ **PM/CEO must approve** all tasks
3. ✅ **REVIEW_PENDING** is the new waiting state
4. ✅ **Timestamp set on approval**, not on AI pass

---

## 📚 Full Docs

- **Complete Guide:** `HUMAN_QUALITY_GATE.md`
- **Workflows:** `WORKFLOW_DIAGRAM.md`
- **Implementation:** `QUALITY_GATE_IMPLEMENTATION.md`
- **This Card:** `APPROVAL_QUICK_REF.md`

---

**Status: ✅ LIVE**

**Human approval now required! 🚦**
