# ⚡ LocalBeam

> انتقال فایل و متن بین دستگاه‌ها در شبکه محلی — بدون اینترنت، بدون کابل، بدون درد سر.

![LocalBeam Banner](https://img.shields.io/badge/LocalBeam-v1.0.0-6c63ff?style=for-the-badge&logo=lightning)
![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

---

## ✨ ویژگی‌ها

- 📁 **انتقال فایل** — هر فایلی تا 500MB (قابل تنظیم)
- 📝 **انتقال متن** — متن، لینک، کد، کلیپ‌بورد
- 📷 **QR Code** — دستگاه گیرنده را با اسکن QR وصل کن
- 🔢 **PIN Code** — بدون QR هم میشه با PIN 6 رقمی وصل شد
- 🔒 **امن** — هرچیزی روی شبکه محلی میمونه، سرور خارجی وجود نداره
- 🔄 **دو‌طرفه** — هر دستگاهی می‌تونه بفرسته یا دریافت کنه
- 📱 **ریسپانسیو** — UI روی موبایل، تبلت و دسکتاپ کار می‌کنه
- 🎨 **UI زیبا** — رابط کاربری تمیز و مدرن
- ⚙️ **قابل تنظیم** — از طریق فایل JSON یا پرچم‌های CLI

---

## 🚀 نصب سریع

### پیش‌نیاز
- Go 1.22 یا بالاتر

### ساخت از سورس

```bash
git clone https://github.com/localbeam/localbeam
cd localbeam
make build
```

یا مستقیم:

```bash
go build -o localbeam ./cmd/localbeam/
```

### اجرا

```bash
./localbeam
```

---

## 📖 راهنمای استفاده

### ۱. سرور را روی دستگاه اصلی اجرا کن

```bash
./localbeam
```

چیزی شبیه به این می‌بینی:

```
 _                    _ ____
| |     ___  ___ __ _| | __ )  ___  __ _ _ __ ___
...

🚀 LocalBeam server running at http://0.0.0.0:8765
📱 Local network: http://192.168.1.105:8765

📷 Scan to open on your device:
[QR Code displayed here]

   URL: http://192.168.1.105:8765
```

### ۲. دستگاه دیگر را وصل کن

**روش A — اسکن QR کد:**
- QR کدی که توی ترمینال نمایش داده می‌شه رو اسکن کن
- مرورگر دستگاه دیگه باز میشه

**روش B — آدرس دستی:**
- آدرس نمایش‌داده شده رو توی مرورگر دستگاه دیگه وارد کن
- مثلاً: `http://192.168.1.105:8765`

### ۳. ارسال فایل یا متن

1. روی **Send** کلیک کن
2. نوع محتوا انتخاب کن: **Files** یا **Text**
3. فایل بکش و بنداز (یا کلیک کن) / متن وارد کن
4. **Create Beam** رو بزن
5. QR Code و PIN Code 6 رقمی نمایش داده میشه

### ۴. دریافت

1. دستگاه گیرنده آدرس سرور رو باز کنه
2. روی **Receive** کلیک کنه
3. PIN 6 رقمی رو وارد کنه (یا QR اسکن کنه)
4. فایل‌ها رو دانلود یا متن رو کپی کنه

---

## ⚙️ تنظیمات

### ساخت فایل پیش‌فرض

```bash
./localbeam --init-config
```

فایل در `~/.localbeam/config.json` ساخته می‌شه.

### نمونه config.json

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": 8765,
    "tls": false
  },
  "security": {
    "session_timeout_minutes": 10,
    "max_sessions": 10,
    "pin_length": 6,
    "allowed_origins": ["*"]
  },
  "transfer": {
    "max_file_size_mb": 500,
    "allowed_types": [],
    "chunk_size_kb": 256,
    "upload_dir": "/tmp",
    "auto_cleanup_minutes": 30
  },
  "ui": {
    "theme": "dark",
    "app_name": "LocalBeam"
  }
}
```

### پرچم‌های CLI

| پرچم | توضیح | پیش‌فرض |
|------|--------|---------|
| `--port` | پورت سرور | `8765` |
| `--host` | هاست سرور | `0.0.0.0` |
| `--config` | مسیر فایل تنظیمات | `~/.localbeam/config.json` |
| `--init-config` | ساخت فایل تنظیمات پیش‌فرض | — |
| `--version` | نمایش نسخه | — |

### مثال‌ها

```bash
# پورت دیگه
./localbeam --port 9000

# فقط روی IP خاص گوش بده
./localbeam --host 192.168.1.100

# با فایل تنظیمات دلخواه
./localbeam --config /etc/localbeam/config.json
```

---

## 🔒 امنیت

- **شبکه محلی only** — تمام داده‌ها فقط روی شبکه محلی منتقل می‌شن
- **Session timeout** — جلسه‌ها به‌طور خودکار منقضی می‌شن (پیش‌فرض: 10 دقیقه)
- **File hash** — هر فایل SHA-256 hash داره برای تأیید صحت
- **Auto cleanup** — فایل‌های موقت بعد از انقضا پاک می‌شن
- **PIN authentication** — دسترسی به session نیاز به PIN داره

> ⚠️ **توجه**: LocalBeam برای شبکه‌های محلی امن طراحی شده. برای استفاده روی شبکه‌های عمومی، TLS فعال کن.

---

## 🏗️ ساختار پروژه

```
localbeam/
├── cmd/
│   └── localbeam/
│       └── main.go          # نقطه ورود برنامه
├── internal/
│   ├── config/
│   │   └── config.go        # مدیریت تنظیمات
│   ├── qr/
│   │   └── qr.go           # تولید QR Code (بدون dependency خارجی)
│   ├── server/
│   │   ├── server.go        # HTTP server و API routes
│   │   └── html.go         # UI (Embedded HTML/CSS/JS)
│   └── transfer/
│       └── session.go       # مدیریت جلسه‌ها و فایل‌ها
├── go.mod
├── Makefile
├── config.example.json
└── README.md
```

---

## 🔌 API

پس از راه‌اندازی سرور:

| Method | Endpoint | توضیح |
|--------|----------|--------|
| `POST` | `/api/session/create` | ایجاد جلسه جدید |
| `GET` | `/api/session/{id}` | اطلاعات جلسه |
| `POST` | `/api/join` | پیوستن با PIN |
| `POST` | `/api/upload/{id}` | آپلود فایل |
| `GET` | `/api/download/{id}/{fileID}` | دانلود فایل |
| `POST` | `/api/text/{id}` | ارسال متن |
| `GET` | `/api/text/{id}` | دریافت متن |
| `GET` | `/api/qr/{id}` | دریافت QR Code (PNG) |
| `GET` | `/api/info` | اطلاعات سرور |

---

## 🛠️ ساخت

```bash
# ساخت binary
make build

# ساخت برای همه پلتفرم‌ها
make build-all

# اجرا مستقیم
make run

# نصب در PATH
make install

# پاکسازی
make clean
```

---

## 📋 سناریوهای استفاده

### گوشی → لپ‌تاپ
1. لپ‌تاپ: `./localbeam` اجرا کن
2. گوشی: مرورگر باز کن، QR اسکن کن یا آدرس وارد کن
3. گوشی: Send → انتخاب فایل → Create Beam
4. لپ‌تاپ: مرورگر باز → Receive → PIN وارد کن → دانلود

### لپ‌تاپ → گوشی
1. لپ‌تاپ: `./localbeam` اجرا کن
2. لپ‌تاپ: Send → فایل انتخاب کن → Create Beam
3. گوشی: آدرس سرور باز کن → Receive → PIN وارد کن → دانلود

### گوشی → گوشی
1. یکی از گوشی‌ها باید سرور را روی یک لپ‌تاپ یا کامپیوتر اجرا کنه
2. هر دو گوشی از طریق مرورگر وصل می‌شن

---

## 📄 لایسنس

MIT License — آزادانه استفاده کنید.
