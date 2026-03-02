# ⚡ CEO AI Control Panel - Quick Start

## 🚀 Access in 3 Steps

### 1. Login as CEO
```
URL: http://localhost:3000/login
Email: ceo@komgrip.com
Password: [your password]
```

### 2. Find the Golden Link
Look for: **⚙️ AI Control Tower** in the sidebar (gold accent, CEO-only)

### 3. Start Controlling AI
You'll see a premium control panel with 3 main sections.

---

## 🎮 Quick Controls

### Change AI Model
```
1. Open "AI Model Selection" panel
2. Click dropdown → Select model
3. Click "SAVE CONFIGURATION"
✅ Done! Next task uses new model
```

### Adjust Consistency
```
Temperature Slider:
← 0.0 (Very consistent estimates)
→ 1.0 (Creative, varied estimates)

Recommended: 0.2-0.4
```

### Set Team AI Level
```
Cursor Assistance Slider:
← 20%  = Manual coding (slower estimates)
→ 50%  = Hybrid (balanced)
→ 80%  = AI-assisted (faster)
→ 100% = Ultra AI mode (very fast!)

Watch the color change: Red → Blue → Purple → Gold
```

---

## 🎯 Common Scenarios

### "Team uses Cursor heavily now"
```
Set Cursor Assistance to 90-95%
Result: AI estimates much faster times
```

### "Estimates are too inconsistent"
```
Lower Temperature to 0.1-0.2
Result: More predictable estimates
```

### "Rate limiting issues"
```
Switch to gemini-1.5-flash
Result: Higher quota, faster responses
```

---

## 📊 What Each Setting Does

| Setting | What It Controls | Impact |
|---------|-----------------|--------|
| **Model** | Which Gemini AI to use | Speed, quality, quota |
| **Temperature** | Response consistency | 0.0=same, 1.0=varied |
| **Cursor %** | Expected AI assistance | Directly affects time estimates |

---

## 💡 Pro Tips

1. **Start Conservative**
   - Model: gemini-2.5-flash-lite
   - Temperature: 0.4
   - Cursor: 80%

2. **Experiment Safely**
   - Changes apply immediately
   - Can reset anytime
   - No server restart needed

3. **Monitor Results**
   - Create test task
   - Check AI estimate
   - Adjust if needed

---

## 🔍 Where to See Results

**Immediate Impact:**
- Create new task → See AI estimate
- Submit code → AI review uses new model
- Submit appeal → AI analysis uses new settings

**Check Logs:**
```bash
docker-compose logs -f api | grep "AI Config"
```

You'll see:
```
⚙️  AI Config: Model=gemini-2.5-flash-lite, Temp=0.30, Cursor=95%
```

---

## 🆘 Troubleshooting

### Can't see AI Control Tower link?
- Must be logged in as CEO
- Refresh page after login

### Changes not saving?
- Check browser console for errors
- Verify API is running: http://localhost:8080/health

### Want to reset everything?
- Click "Reset to Defaults" button
- Confirms before resetting

---

## 📱 Access URLs

**Control Panel:** http://localhost:3000/admin/ai-settings
**Dashboard:** http://localhost:3000/dashboard
**API Health:** http://localhost:8080/health

---

## 🎨 Visual Guide

### Cursor Assistance Colors
```
🔴 Red (0-20%)    = Manual coding
🔵 Blue (21-50%)  = Hybrid workflow  
🟣 Purple (51-80%) = AI-assisted
🟡 Gold (81-100%) = Ultra AI mode
```

### What You'll See
```
Header: "AI CONTROL TOWER" (gold gradient)
Status: Green "LIVE" badge
Panels: 3 glassmorphism cards
Sliders: Smooth, colorful, responsive
Button: Gold gradient with glow
Toast: Green success notification
```

---

## ⏱️ Quick Reference

| Action | Time | Effect |
|--------|------|--------|
| Page Load | 1s | Fetches config |
| Change Setting | Instant | Visual update |
| Save Config | 0.5s | Persists changes |
| Toast Display | 3s | Auto-dismisses |

---

## 🎯 Default Settings

```json
{
  "active_model": "gemini-2.5-flash-lite",
  "temperature": 0.4,
  "cursor_assistance": 80
}
```

These are balanced for most teams. Adjust based on your actual workflow!

---

## 📚 More Info

- **Full Guide:** `CEO_AI_CONTROL_PANEL.md`
- **Summary:** `CEO_AI_PANEL_SUMMARY.md`
- **API Docs:** `DYNAMIC_AI_CONFIG_GUIDE.md`

---

**That's it! You're ready to control AI like a boss.** 👑
