# ติดตั้ง Homebrew และ Node.js (สำหรับ MCP Browser / โปรเจกต์)

ถ้าเครื่องยังไม่มี Homebrew หรือ Node ทำตามนี้ใน **Terminal** (ต้องใช้รหัสผ่านเครื่องหนึ่งครั้ง)

---

## ขั้นที่ 1: ติดตั้ง Homebrew

เปิด **Terminal** แล้วรันคำสั่งเดียวนี้ (กด Enter แล้วใส่รหัสผ่านเมื่อถาม):

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

รอจนติดตั้งเสร็จ แล้วทำขั้นที่ 2 ต่อ (เพิ่ม brew เข้า PATH)

---

## ขั้นที่ 2: เพิ่ม Homebrew เข้า PATH

หลังติดตั้ง Homebrew เสร็จ ต้องเพิ่มเข้า PATH ก่อนใช้คำสั่ง `brew`:

**ถ้าเป็น Mac แบบ Apple Silicon (M1/M2/M3):**
```bash
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile && eval "$(/opt/homebrew/bin/brew shellenv)"
```

**ถ้าเป็น Mac แบบ Intel:**
```bash
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile && eval "$(/usr/local/bin/brew shellenv)"
```

จากนั้นตรวจสอบ:
```bash
brew --version
```

---

## ขั้นที่ 3: ติดตั้ง Node.js

```bash
brew install node
```

ตรวจสอบ:
```bash
node -v
npx -v
```

---

## ขั้นที่ 4: ใช้ MCP Browser ใน Cursor

1. ปิด Cursor แล้วเปิดใหม่ (หรือ Reload Window)
2. MCP browser ควรใช้ได้ เพราะ `.cursor/run-npx.sh` จะหา `npx` จาก `/opt/homebrew/bin` หรือ `/usr/local/bin` ให้เอง

---

## สรุปคำสั่ง ( copy-paste ทั้งบล็อก )

**Apple Silicon:**
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile && eval "$(/opt/homebrew/bin/brew shellenv)"
brew install node
```

**Intel:**
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
echo 'eval "$(/usr/local/bin/brew shellenv)"' >> ~/.zprofile && eval "$(/usr/local/bin/brew shellenv)"
brew install node
```
