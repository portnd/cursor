# 🤖 AI Advisory System - Complete Guide

## 📖 Overview
The AI Advisory System provides intelligent recommendations to CEO/PM when reviewing developer appeals. It analyzes code diffs, original AI feedback, and developer arguments to suggest whether to approve or reject an appeal.

---

## 🔄 Complete Flow

```
1. DEVELOPER SUBMITS CODE
   POST /tasks/:id/submit
   Body: { commit_hash, diff }
   ↓
   Submission created with:
   - ai_verdict: FAIL
   - ai_score: 25
   - ai_feedback: "SQL injection vulnerability..."
   - diff: "SELECT * FROM users WHERE id = " + userID

2. DEVELOPER APPEALS
   POST /submissions/:id/appeal
   Body: { reason: "I am using GORM which prevents SQL injection..." }
   ↓
   🤖 AI ADVISORY ANALYSIS TRIGGERED
   ├─ Input: diff + originalFeedback + appealReason
   ├─ Gemini AI analyzes validity
   └─ Output: recommendation + confidence + reasoning
   ↓
   Appeal created with:
   - status: PENDING
   - ai_recommendation: "OVERTURN" or "UPHOLD"
   - ai_confidence: 85
   - ai_reasoning: "Developer is correct. GORM prevents..."

3. CEO/PM REVIEWS APPEAL
   Navigate to: /task/:id
   Click: "⚖️ Review Appeal"
   ↓
   🎨 ADJUDICATION MODAL OPENS
   Shows:
   ├─ Developer's Plea
   ├─ 🤖 AI Advisor Opinion (NEW!)
   │   ├─ Recommendation: ✅ APPROVE or ❌ REJECT
   │   ├─ Confidence: 85% (progress bar)
   │   └─ Reasoning: Full AI analysis
   └─ Quick Action: Apply AI Recommendation

4. CEO/PM MAKES DECISION
   Option A: Click "⚡ Apply AI Recommendation"
   ↓
   - Auto-fills resolver note with AI reasoning
   - Submits decision (APPROVED or REJECTED)
   
   Option B: Manual decision (override AI)
   ↓
   - Click Approve or Reject button manually
   - Can add own reasoning

5. APPEAL RESOLVED
   POST /appeals/:id/resolve
   Body: { status: "APPROVED", note: "..." }
   ↓
   If APPROVED:
   - Submission verdict → PASS
   - is_overridden → true
   - Task → COMPLETED
   
   If REJECTED:
   - Submission remains FAIL
   - Appeal status → REJECTED
```

---

## 🗄️ Database Schema

### **Appeals Table**
```sql
CREATE TABLE appeals (
  id UUID PRIMARY KEY,
  submission_id UUID NOT NULL,
  developer_id INTEGER NOT NULL,
  reason TEXT NOT NULL,
  status VARCHAR DEFAULT 'PENDING', -- PENDING, APPROVED, REJECTED
  
  -- AI Advisory System (NEW)
  ai_recommendation TEXT,  -- OVERTURN or UPHOLD
  ai_confidence INTEGER,   -- 0-100
  ai_reasoning TEXT,       -- Advice for CEO/PM
  
  resolver_id INTEGER,
  resolver_note TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### **Submissions Table**
```sql
CREATE TABLE submissions (
  id UUID PRIMARY KEY,
  task_id UUID NOT NULL,
  dev_id INTEGER NOT NULL,
  commit_hash VARCHAR NOT NULL,
  
  diff TEXT,  -- NEW: Code diff for appeal analysis
  
  ai_verdict VARCHAR,   -- PASS, FAIL, PENDING
  ai_score INTEGER,
  ai_feedback JSONB,
  
  is_overridden BOOLEAN DEFAULT false,
  created_at TIMESTAMP
);
```

---

## 🔌 API Endpoints

### **Submit Work**
```http
POST /api/v1/sentinel/tasks/:id/submit
Authorization: Bearer <token>
Content-Type: application/json

{
  "commit_hash": "abc123",
  "diff": "func Login() {\n  query := \"SELECT * FROM users WHERE name=\" + username\n}"
}
```

**Response:**
```json
{
  "data": {
    "id": "...",
    "ai_verdict": "FAIL",
    "ai_score": 25,
    "ai_feedback": {...},
    "diff": "func Login() {...}"
  }
}
```

---

### **Submit Appeal (with AI Analysis)**
```http
POST /api/v1/sentinel/submissions/:id/appeal
Authorization: Bearer <token>
Content-Type: application/json

{
  "reason": "The AI review misunderstood the context. I am using parameterized queries."
}
```

**Response:**
```json
{
  "data": {
    "id": "appeal-uuid",
    "submission_id": "sub-uuid",
    "developer_id": 3,
    "reason": "The AI review misunderstood...",
    "status": "PENDING",
    
    "ai_recommendation": "OVERTURN",
    "ai_confidence": 85,
    "ai_reasoning": "Developer is correct. GORM uses parameterized queries. False positive.",
    
    "resolver_id": null,
    "resolver_note": "",
    "created_at": "...",
    "updated_at": "..."
  },
  "message": "Appeal submitted successfully"
}
```

---

### **Resolve Appeal**
```http
POST /api/v1/sentinel/appeals/:id/resolve
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "APPROVED",
  "note": "Following AI recommendation (85% confidence): Developer is correct..."
}
```

**Response:**
```json
{
  "message": "Appeal resolved successfully"
}
```

---

## 🎨 Frontend UI Components

### **Task Detail Page**
**File:** `web/pages/task/[id].vue`

#### **Submission Card (Timeline)**
Shows appeal status:
- **PENDING:** Yellow banner "⚖️ Appeal Under Review"
- **CEO/PM:** "Review Appeal" button (purple-gold)

#### **Adjudication Modal**
Opens when CEO/PM clicks "Review Appeal"

**Structure:**
```vue
<div class="adjudication-modal">
  <!-- Header: High Court of Sentinel -->
  <header>CEO presiding</header>
  
  <!-- Case Information -->
  <section>Case #, Submission ID, Appellant</section>
  
  <!-- The Plea -->
  <section>Developer's appeal reason</section>
  
  <!-- 🤖 AI Advisor Opinion (NEW!) -->
  <section class="ai-advisor">
    <!-- Recommendation Badge -->
    <div :class="recommendation-theme">
      ✅ APPROVE APPEAL (OVERTURN)
      or
      ❌ REJECT APPEAL (UPHOLD)
    </div>
    
    <!-- Confidence Meter -->
    <div class="confidence-meter">
      Progress bar (0-100%)
      High/Moderate/Low interpretation
    </div>
    
    <!-- Analysis Report -->
    <div class="reasoning">
      {{ ai_reasoning }}
    </div>
    
    <!-- Quick Action -->
    <button @click="applyAIRecommendation">
      ⚡ Apply AI Recommendation
    </button>
  </section>
  
  <!-- Resolver Note -->
  <textarea v-model="resolverNote" />
  
  <!-- Verdict Buttons -->
  <div class="verdict-actions">
    <button @click="resolveAppeal('APPROVED')">✅ Sustain Appeal</button>
    <button @click="resolveAppeal('REJECTED')">❌ Dismiss Appeal</button>
  </div>
</div>
```

---

## 🧪 Complete Test Scenario

### **Scenario: False Positive SQL Injection**

#### **Step 1: Developer Submits Code**
```bash
DEV_TOKEN="..."

curl -X POST "http://localhost:8080/api/v1/sentinel/tasks/TASK_ID/submit" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "commit_hash": "test123",
    "diff": "func GetUser(id string) (*User, error) {\n  return db.Where(\"id = ?\", id).First(&user)\n}"
  }'
```

**Expected:** AI returns `FAIL` (overzealous review)

---

#### **Step 2: Developer Appeals**
```bash
curl -X POST "http://localhost:8080/api/v1/sentinel/submissions/SUB_ID/appeal" \
  -H "Authorization: Bearer $DEV_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "This is a false positive. I am using GORM with parameterized queries. The ? placeholder prevents SQL injection. The AI did not recognize the ORM pattern."
  }'
```

**Expected Response:**
```json
{
  "ai_recommendation": "OVERTURN",
  "ai_confidence": 90,
  "ai_reasoning": "Developer is correct. GORM uses parameterized queries via ? placeholders which prevent SQL injection. Original review was a false positive."
}
```

---

#### **Step 3: CEO Reviews in Browser**
1. Login as CEO: `ceo@sentinel.com`
2. Navigate to: `http://localhost:3000/task/TASK_ID`
3. Click: "⚖️ Review Appeal" button

**Expected UI:**
```
┌──────────────────────────────────────────────┐
│ 🤖 AI Advisor Opinion                        │
│ Advanced Legal Analysis System               │
├──────────────────────────────────────────────┤
│ ✅  AI Suggests: APPROVE APPEAL              │
│                                               │
│ Confidence Level                       90%   │
│ ██████████████████░░                          │
│ High confidence - Strong recommendation      │
│                                               │
│ 📊 Analysis Report                           │
│ Developer is correct. GORM uses              │
│ parameterized queries via ? placeholders     │
│ which prevent SQL injection. Original        │
│ review was a false positive.                 │
│                                               │
│ [ ⚡ Apply AI Recommendation ]               │
└──────────────────────────────────────────────┘
```

---

#### **Step 4: CEO Applies AI Recommendation**
1. Click: "⚡ Apply AI Recommendation"
2. Resolver note auto-fills: "Following AI recommendation (90% confidence): Developer is correct..."
3. Appeal resolves as APPROVED
4. Modal closes, page refreshes

**Expected Results:**
- Submission card turns **GOLD**
- Badge: "👑 OVERRIDDEN BY HUMAN"
- Verdict: "👑 VERDICT OVERRIDDEN"
- Task status: **COMPLETED** (green)

---

#### **Step 5: Verify Database**
```sql
SELECT 
  a.status as appeal_status,
  a.ai_recommendation,
  a.ai_confidence,
  s.ai_verdict,
  s.is_overridden,
  t.status as task_status
FROM appeals a
JOIN submissions s ON s.id = a.submission_id
JOIN tasks t ON t.id = s.task_id
WHERE a.id = 'APPEAL_ID';
```

**Expected:**
```
appeal_status: APPROVED
ai_recommendation: OVERTURN
ai_confidence: 90
ai_verdict: PASS
is_overridden: true
task_status: COMPLETED
```

---

## 📊 AI Recommendation Logic

### **When AI Suggests OVERTURN (Approve)**
- Developer's defense is technically sound
- Original review was likely a false positive
- Code does not have the claimed vulnerability
- Context was missed by original AI review

**Example Reasons:**
- "Developer uses ORM which prevents SQL injection"
- "Code has proper input validation that was overlooked"
- "Security concern is mitigated by framework protections"

### **When AI Suggests UPHOLD (Reject)**
- Original review was correct
- Vulnerability exists as described
- Developer's defense is weak or incorrect
- Risk is real and should not be ignored

**Example Reasons:**
- "Direct string concatenation creates SQL injection risk"
- "Developer's claim about ORM usage is not reflected in code"
- "Vulnerability is confirmed, requires immediate fix"

---

## 🔍 Troubleshooting

### **Problem: AI Advisory Shows "UPHOLD 0% - AI analysis unavailable"**

**Cause:** Submission has no `diff` field (old submission before feature was implemented)

**Solution:** 
- This is expected for legacy submissions
- CEO/PM should make manual decision
- For new submissions, ensure diff is submitted

---

### **Problem: Gemini API Returns 429 Quota Exceeded**

**Cause:** Free tier API quota exhausted (20 requests/minute)

**Solution:**
- Wait ~60 seconds for quota reset
- Or upgrade to paid tier for higher limits
- System uses fallback values (UPHOLD, 0%, "AI unavailable")

---

### **Problem: AI Recommendation Seems Wrong**

**Cause:** AI analysis is probabilistic, not perfect

**Solution:**
- Check confidence score (low = uncertain)
- Read AI reasoning carefully
- CEO/PM should use judgment to override if needed
- This is why human is the final decision-maker

---

## 📈 Success Metrics

Track these KPIs:

1. **AI Accuracy:** How often does CEO/PM agree with AI?
2. **Confidence Correlation:** Do higher confidence scores = higher agreement?
3. **Time Savings:** How much faster are decisions with AI advice?
4. **Override Rate:** How often does CEO/PM override AI?
5. **Appeal Resolution Rate:** % of appeals resolved within 24 hours

---

## 🚀 Future Enhancements

### **Phase 2 Ideas**
- [ ] Multi-language support for reasoning
- [ ] Historical accuracy display ("AI has been 85% accurate")
- [ ] Explanation of confidence score calculation
- [ ] Link to similar past cases
- [ ] CEO/PM feedback loop (was AI right?)

### **Phase 3 Ideas**
- [ ] Multiple AI models comparison
- [ ] Vote-based ensemble AI recommendations
- [ ] Appeal category classification
- [ ] Automated low-confidence routing to specific reviewers

---

## 📚 Related Documentation

- **Backend Implementation:** `AI_ADVISORY_IMPLEMENTATION.md`
- **Frontend UI:** `AI_ADVISOR_UI_IMPLEMENTATION.md`
- **Quick Summary:** `AI_ADVISORY_SUMMARY.md`
- **Appeal Review (without AI):** `APPEAL_REVIEW_IMPLEMENTATION.md`

---

**🎉 AI Advisory System is fully implemented end-to-end!**

From code submission → AI analysis → visual UI → human decision, the complete flow is operational and ready for production use.
