# ✏️🗑️ Task Edit & Delete - Quick Start Guide

## 🎯 Quick Access

**Location:** Task Detail Page (`/task/:id`)  
**Buttons:** Top right header (next to status badge)  
**Permissions:** CEO or Task Creator only

---

## 🔑 Who Can Do What?

| Action | CEO | Task Creator | Other Users |
|--------|-----|--------------|-------------|
| **View Task** | ✅ Yes | ✅ Yes | ✅ Yes |
| **Edit Task** | ✅ Yes | ✅ Yes (own) | ❌ No |
| **Delete Task** | ✅ Yes | ✅ Yes (own) | ❌ No |

---

## ✏️ How to Edit a Task

### **Step 1: Open Edit Modal**
1. Navigate to task detail page
2. Click **"✏️ Edit"** button (top right)

### **Step 2: Modify Fields**
- **Title:** Change mission title
- **Description:** Update mission details
- **Deadline:** Adjust due date/time (optional)

### **Step 3: Submit**
- Click **"💾 Update Mission"**
- See alert: "Mission Updated! AI is recalculating..."
- Page refreshes with new AI estimate

---

## 🗑️ How to Delete a Task

### **Step 1: Open Confirmation**
1. Navigate to task detail page
2. Click **"🗑️ Abort"** button (top right)

### **Step 2: Confirm Deletion**
- Read warning: "This will remove all data"
- Review task info
- Click **"💥 Yes, Delete Forever"**

### **Step 3: Auto-Redirect**
- See alert: "Mission Deleted Successfully!"
- Automatically redirected to dashboard
- Task permanently removed

---

## ⚠️ Important Warnings

### **Edit Warning**
```
⚠️ AI Re-estimation Alert

Changing the title or description will trigger 
automatic AI re-estimation. Any pending time 
negotiation will be reset.
```

### **Delete Warning**
```
⚠️ Critical Operation

This will remove:
• Mission data
• All submissions
• All appeals
• Complete audit trail

This action cannot be undone!
```

---

## 🧪 Quick Test (5 Minutes)

### **Test Edit (as CEO)**
```bash
1. Login: ceo@sentinel.com / password123
2. Go to: http://localhost:3000/task/{any-task-id}
3. Click: "✏️ Edit"
4. Change title to: "Updated Secure System"
5. Click: "Update Mission"
6. Result: AI re-estimates time, page refreshes
```

### **Test Delete (as CEO)**
```bash
1. Login: ceo@sentinel.com / password123
2. Go to: http://localhost:3000/task/{any-task-id}
3. Click: "🗑️ Abort"
4. Click: "Yes, Delete Forever"
5. Result: Redirects to dashboard, task gone
```

### **Test Permission Denial (as DEV)**
```bash
1. Login: dev@sentinel.com / password123
2. Go to: http://localhost:3000/task/{ceo-created-task}
3. Look for Edit/Delete buttons
4. Result: Buttons NOT VISIBLE (no permission)
```

---

## 🚨 Troubleshooting

### **Problem: Edit/Delete buttons not visible**
**Solution:**
- Check if you're the CEO
- OR check if you created this task
- Only authorized users see these buttons

### **Problem: "No changes detected" error**
**Solution:**
- Modify at least one field (title or description)
- Current: "Old Title" → New: "Old Title" (no change)
- Fix: Change to different value

### **Problem: "Title is required" error**
**Solution:**
- Title field cannot be empty
- Fill in the title field

### **Problem: API returns 403 Forbidden**
**Solution:**
- You don't have permission to edit/delete this task
- Only CEO or task creator can modify

---

## 🎨 Visual Guide

### **Before (No Permission)**
```
┌────────────────────────────────────────┐
│ Title                    [PENDING] [Back] │
└────────────────────────────────────────┘
```

### **After (With Permission)**
```
┌──────────────────────────────────────────────────────┐
│ Title            [PENDING] [✏️Edit] [🗑️Abort] [Back] │
└──────────────────────────────────────────────────────┘
```

---

## 📊 What Happens After Edit?

### **Immediate Effects:**
1. ✅ Task title/description updated
2. ✅ AI re-estimates time (if content changed)
3. ✅ Negotiation status reset to "NONE" (if was PENDING)
4. ✅ Page refreshes with new data

### **Backend Actions:**
```
🔄 Task content changed. Triggering AI re-estimation...
   Old: [Fix bug] Fix login bug
   New: [Implement secure auth] Add 2FA system
🧠 AI Estimation Request: Implement secure auth
📡 Calling Gemini API (model: gemini-2.5-flash)
✅ AI Re-estimation Complete: 180 minutes (3.0 hours)
🔄 Resetting negotiation status (AI has new estimate)
✅ Task Updated: uuid by CEO (User ID: 1)
```

---

## 📊 What Happens After Delete?

### **Immediate Effects:**
1. ✅ Task permanently deleted from database
2. ✅ All submissions removed
3. ✅ All appeals removed
4. ✅ Audit log created
5. ✅ User redirected to dashboard

### **Backend Actions:**
```
🗑️ Task Deleted: uuid by CEO (User ID: 1)
```

---

## 🔍 Common Scenarios

### **Scenario 1: CEO Updates Any Task**
```
1. CEO navigates to task
2. Sees Edit/Delete buttons
3. Clicks Edit
4. Changes description
5. Submits
6. AI re-estimates (30min → 2h)
7. Success!
```

### **Scenario 2: Developer Updates Own Task**
```
1. Dev navigates to own task
2. Sees Edit/Delete buttons
3. Clicks Edit
4. Adjusts description
5. Submits
6. AI re-estimates
7. Success!
```

### **Scenario 3: Developer Tries to Edit Other's Task**
```
1. Dev navigates to PM's task
2. Edit/Delete buttons HIDDEN
3. Cannot modify (correct behavior)
```

### **Scenario 4: CEO Deletes Obsolete Task**
```
1. CEO finds old task
2. Clicks "Abort Mission"
3. Reads warning
4. Confirms deletion
5. Task removed forever
6. Redirected to dashboard
```

---

## 💡 Pro Tips

### **For Editing:**
✅ **Only change what's needed** - System only sends modified fields  
✅ **Expect AI delay** - Re-estimation takes 2-5 seconds  
✅ **Check negotiation** - Pending negotiations will reset  

### **For Deleting:**
⚠️ **Double-check before deleting** - Action is permanent  
⚠️ **Backup important data** - No restore option  
✅ **Use for cleanup only** - Delete obsolete/duplicate tasks  

---

## 🎯 Best Practices

### **When to Edit:**
- ✅ Fix typos in title/description
- ✅ Adjust scope after requirements change
- ✅ Update deadline
- ✅ Clarify mission objectives

### **When to Delete:**
- ✅ Duplicate task created by mistake
- ✅ Obsolete/cancelled mission
- ✅ Testing/demo task cleanup
- ✅ Incorrect task creation

### **When NOT to Edit:**
- ❌ Task is in progress with submissions
- ❌ Just to "game" the AI estimate
- ❌ After developer already started work

### **When NOT to Delete:**
- ❌ Task has valuable submission history
- ❌ Task is completed (keep for audit)
- ❌ Task has pending appeals
- ❌ Unsure if it's needed

---

## 📞 Support

### **Permission Issues:**
- Contact CEO to modify tasks you didn't create
- CEO has unrestricted access to all tasks

### **Technical Issues:**
- Check browser console for errors
- Check API logs: `docker compose logs api`
- Verify JWT token is valid

### **Data Recovery:**
- Deleted tasks cannot be restored
- Check database backups if critical

---

**🎉 Ready to use! Edit and delete tasks with confidence!** ✏️🗑️🚀
