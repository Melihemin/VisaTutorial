# Hızlı Başlangıç Kılavuzu / Quick Start Guide

## 🚀 Türkçe / Turkish

### 5 Dakikalık Kurulum

1. **VFS Global'a Giriş Yap**
   - https://visa.vfsglobal.com adresine git
   - Hesabınızla giriş yap

2. **Gerekli Token'ları Topla**
   ```
   a) Browser DevTools aç (F12)
   b) Application → Cookies → "sessionCookie" değerini kopyala
   c) Network tab'ında API isteğine bak
   d) Headers bölümünden:
      - "authorize" header'ını kopyala
      - "clientsource" header'ını kopyala
   ```

3. **.env Dosyasını Düzenle**
   ```bash
   # .env dosyasını aç ve aşağıdakileri doldur:
   SESSION_COOKIE=<kopyaladığın cookie>
   AUTHORIZE_TOKEN=<authorize header>
   CLIENT_SOURCE_TOKEN=<clientsource header>
   ```

4. **Programı Çalıştır**
   ```bash
   go run main.go proxy.go config.go notifier.go utils.go
   ```

### Telegram Bildirimi Ekle (İsteğe Bağlı)

1. BotFather'dan bot oluştur: @BotFather
2. Aldığın token'ı `.env`e ekle: `TELEGRAM_BOT_TOKEN=...`
3. Chat ID'ni bularak ekle: `TELEGRAM_CHAT_ID=...`

---

## 🚀 English / İngilizce

### 5-Minute Setup

1. **Login to VFS Global**
   - Visit https://visa.vfsglobal.com
   - Login with your account

2. **Collect Required Tokens**
   ```
   a) Open Browser DevTools (F12)
   b) Application → Cookies → Copy "sessionCookie" value
   c) Check Network tab for API request
   d) Copy from Headers:
      - "authorize" header value
      - "clientsource" header value
   ```

3. **Edit .env File**
   ```bash
   # Open .env and fill in:
   SESSION_COOKIE=<your_copied_cookie>
   AUTHORIZE_TOKEN=<authorize_header>
   CLIENT_SOURCE_TOKEN=<clientsource_header>
   ```

4. **Run the Program**
   ```bash
   go run main.go proxy.go config.go notifier.go utils.go
   ```

### Add Telegram Notifications (Optional)

1. Create a bot with BotFather: @BotFather
2. Add bot token to `.env`: `TELEGRAM_BOT_TOKEN=...`
3. Add Chat ID to `.env`: `TELEGRAM_CHAT_ID=...`

---

## 📁 Dosya Yapısı / File Structure

```
randevu_tracker/
├── main.go              # Ana slot kontrol mantığı / Main slot checking logic
├── proxy.go             # Proxy ve worker yönetimi / Proxy & worker management
├── config.go            # Konfigürasyon yükleme / Configuration loading
├── notifier.go          # Telegram bildirimleri / Telegram notifications
├── utils.go             # Logger ve istatistikler / Logger & statistics
├── .env                 # Kişisel ayarlar / Personal settings
├── .env.example         # Örnek ayarlar / Example settings
└── go.mod              # Go bağımlılıkları / Go dependencies
```

---

## ⚙️ Temel Komutlar / Basic Commands

```bash
# Bağımlılıkları indir / Download dependencies
go mod download

# Programı çalıştır / Run program
go run *.go

# Derle / Build executable
go build -o randevu_tracker.exe

# Çalıştırılabilir dosyayı çalıştır / Run executable
./randevu_tracker.exe
```

---

## 🔍 Debugging

### Log Dosyalarını Kontrol Et / Check Log Files

```bash
# Günlük log / Daily log
type randevu_tracker_2026-09-01.log

# Slot kontrol log / Slot check log
type slot_checks_2026-09-01.log
```

### Yaygın Sorunlar / Common Issues

| Sorun / Problem | Çözüm / Solution |
|---|---|
| "401 Unauthorized" | Token'ları yenile / Refresh tokens |
| "429 Too Many Requests" | Kontrol aralığını artır / Increase interval |
| "Proxy Error" | Farklı proxy kullan / Use different proxy |
| "Telegram not working" | Bot token'ı kontrol et / Check bot token |

---

## 💡 İpuçları / Tips

- 🔐 Kimlik bilgilerini asla paylaşma / Never share credentials
- ⏱️ Minimum 45 saniye aralık bırak / Keep minimum 45 sec interval
- 📊 Log dosyalarını izle / Monitor log files
- 🔄 Proxy kullan (isteğe bağlı) / Use proxy (optional)

---

## 📞 Yardım / Help

Sorular için log dosyalarını kontrol edin ve emin olun:
- SESSION_COOKIE doğru ve güncel
- AUTHORIZE_TOKEN geçerli
- CLIENT_SOURCE_TOKEN eksik değil

For questions, check log files and ensure:
- SESSION_COOKIE is correct and up to date
- AUTHORIZE_TOKEN is valid
- CLIENT_SOURCE_TOKEN is not missing
