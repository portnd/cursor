# 📥 Tasks Inbox - Quick Start Guide

## 🚀 Access the Inbox

**URL:** `http://localhost:3000/tasks`

**Navigation:** Click "📥 Tasks Inbox" in the sidebar

---

## 👥 What You'll See Based on Your Role

### **CEO & PM**
You'll see **TWO sections:**

#### 1️⃣ 🚨 **PENDING APPROVALS** (Top, Gold Border)
- Tasks requiring your decision
- **Time Negotiations:** Developer wants more time than AI estimated
- **Appeals:** Developer disputes AI verdict

**Example:**
```
┌─────────────────────────────────────────┐
│ 🚨 PENDING APPROVALS              [1]   │
├─────────────────────────────────────────┤
│ Implement secure database query         │
│ [⏱️ TIME NEGOTIATION] [⚖️ APPEAL]      │
│                                          │
│ AI: 30min → Dev wants: 120min          │
│ Reason: Legacy code complexity...       │
│                                          │
│ [⚖️ Review Request →]                   │
└─────────────────────────────────────────┘
```

#### 2️⃣ ⚡ **MY ACTIVE MISSIONS**
- Your assigned tasks
- Same as what developers see

---

### **Developers**
You'll see **ONE section:**

#### ⚡ **MY ACTIVE MISSIONS**
- All tasks assigned to you
- Color-coded by urgency:
  - 🚨 **Red Border + OVERDUE Badge:** Past deadline
  - ⚠️ **Yellow Border + URGENT Badge:** < 24 hours left
  - **Gray Border:** Normal

**Example:**
```
┌─────────────────────────────────────────┐
│ ⚡ MY ACTIVE MISSIONS              [3]  │
├─────────────────────────────────────────┤
│ ┌───────────────┐  ┌───────────────┐   │
│ │ 🚨 OVERDUE    │  │ Fix API Bug   │   │
│ │ SQL Injection │  │ [IN_PROGRESS] │   │
│ │ ⏰ 2h overdue │  │ ⏰ 5h left    │   │
│ │ Execute →     │  │ Execute →     │   │
│ └───────────────┘  └───────────────┘   │
└─────────────────────────────────────────┘
```

---

## 🎯 Quick Actions

### **For Approvals (CEO/PM)**
1. Click **"⚖️ Review Request"**
2. Goes to task detail page
3. Make your decision (approve/reject)

### **For Tasks (Everyone)**
1. Click **"Execute →"**
2. Goes to task detail page
3. Work on task / Submit code

---

## 🎨 Visual Guide

### **Urgency Indicators**

| Visual | Meaning | Action Required |
|--------|---------|-----------------|
| 🚨 **OVERDUE** Badge (Red, Pulsing) | Past deadline | **URGENT!** |
| ⚠️ **URGENT** Badge (Yellow) | < 24 hours left | Act soon |
| ⏰ Countdown Timer | Time remaining | Monitor |
| No badge | Normal deadline | Scheduled |

### **Approval Types**

| Badge | Meaning |
|-------|---------|
| ⏱️ **TIME NEGOTIATION** | Dev wants more time |
| ⚖️ **APPEAL PENDING** | Dev disputes AI verdict |

---

## 💡 Tips

1. **Check Daily:** Start your day by reviewing the inbox
2. **Prioritize Overdue:** Handle red-bordered tasks first
3. **Review AI Advice:** For appeals, check AI recommendation and confidence
4. **Use Empty State:** "All Systems Clear" means you're caught up!

---

## 🐛 Troubleshooting

### **Problem: Empty inbox but I have tasks**
- **Solution:** Tasks might not be assigned to you yet
- Check `/dashboard` for unassigned tasks

### **Problem: Don't see Approvals section**
- **Solution:** This section is only for CEO and PM
- If you're a developer, you won't see it

### **Problem: Page won't load**
- **Solution:** Make sure you're logged in
- Try refreshing the page
- Check if API is running: `http://localhost:8080/health`

---

## 📊 Data Refresh

The page loads data when you:
- First open it
- Click "Retry" after an error

To see updated data:
- Refresh the page (F5 or Cmd+R)
- Or navigate away and back

---

## 🎯 Common Workflows

### **CEO Morning Routine**
```
1. Login
2. Open Tasks Inbox
3. Review Pending Approvals first
4. Handle urgent appeals
5. Approve/reject time negotiations
6. Check own tasks
```

### **Developer Daily Check**
```
1. Login
2. Open Tasks Inbox
3. Check for OVERDUE items (act immediately)
4. Review URGENT items (plan for today)
5. Monitor countdown timers
6. Execute on highest priority
```

---

## 🔗 Quick Links

- **API Health:** `http://localhost:8080/health`
- **Tasks Inbox:** `http://localhost:3000/tasks`
- **Dashboard:** `http://localhost:3000/dashboard`
- **Create Task:** `http://localhost:3000/create`

---

**🎉 You're all set! Start managing your tasks like a pro!** 🚀
