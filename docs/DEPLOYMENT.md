# Randevu Uygunluk Takip Sistemi - Dağıtım Kılavuzu

## 📦 Kurulum Gereksinimleri

### Doğrudan Go ile

- **Go 1.27+**
- **Windows/Linux/macOS**
- **İnternet bağlantısı**

```bash
# Kurulum kontrol et
go version
```

### Docker ile

- **Docker 20.10+**
- **Docker Compose 2.0+**

```bash
docker --version
docker-compose --version
```

## 🚀 Dağıtım Yöntemleri

### 1. Doğrudan Go ile (Geliştirme)

#### Step 1: Repository'yi Klonla

```bash
cd c:\Users\Melih\OneDrive\Desktop\Projects\go
```

#### Step 2: Bağımlılıkları Yükle

```bash
go mod download
go mod tidy
```

#### Step 3: Derle

```bash
# Tek komut ile
make build

# Veya manuel
go build -o randevu_tracker.exe \
  main.go proxy.go config.go notifier.go utils.go \
  orchestrator.go api_service.go worker_service.go \
  session_service.go proxy_service.go
```

#### Step 4: Çalıştır

```bash
# Make ile
make run

# Veya direkt
./randevu_tracker.exe
```

### 2. Docker Compose ile (Üretim)

#### Step 1: .env Dosyasını Hazırla

```bash
make env
# Sonra .env dosyasını düzenle
```

#### Step 2: Docker Image'ları Derle

```bash
make docker-build
```

#### Step 3: Kontainer'ları Başlat

```bash
make docker-up
```

Servisler başlatıldıktan sonra:
- **API**: http://localhost:8080
- **Session**: http://localhost:8081
- **Proxy**: http://localhost:8082

#### Step 4: Logları İzle

```bash
make docker-logs

# Veya belirli servisler
docker-compose logs -f api
docker-compose logs -f worker
```

#### Step 5: Durdur

```bash
make docker-down
```

## 📝 Konfigürasyon

### .env Dosyası

```env
# VFS Global API Kimlik Bilgileri
SESSION_COOKIE=abc123xyz...
AUTHORIZE_TOKEN=token_...
CLIENT_SOURCE_TOKEN=source_...

# HTTP Başlıkları
USER_AGENT=Mozilla/5.0 (Windows NT 10.0; Win64; x64)...
SEC_CH_UA="Not=A?Brand";v="99", "Google Chrome";v="151"

# Proxy
PROXY_LIST=http://proxy1:8080,http://proxy2:8080
USE_PROXY=false

# Kontrol Parametreleri
CHECK_INTERVAL=45
REQUEST_TIMEOUT=30
MAX_RETRIES=3
LOG_TO_FILE=true

# Payload
COUNTRY=usa
MISSION=che
VACANCY_CODE=NYC
VISA_CATEGORY_CODE=SUSVFF
LOGIN_USER=email@example.com

# Telegram Bildirimleri
TELEGRAM_BOT_TOKEN=123456:ABC...
TELEGRAM_CHAT_ID=987654321
```

## 🔧 Kalite Kontrol

### Kodları Formatla

```bash
make fmt
```

### Linting

```bash
make lint
```

### Testleri Çalıştır

```bash
make test
```

## 🐛 Sorun Giderme

### Build Hatası: "go: command not found"

**Çözüm**: Go'yu yükleyin
```bash
# Windows: https://golang.org/dl indirin
# Linux: sudo apt-get install golang-go
```

### Docker Hatası: "Cannot connect to Docker daemon"

**Çözüm**: Docker'ı başlat
```bash
# Windows: Docker Desktop'u başlat
# Linux: sudo systemctl start docker
```

### Port Zaten Kullanımda

**Çözüm**: Port değiştir veya konflikt eden işlem durdur
```bash
# Hangi işlem port 8080'i kullanıyor?
netstat -ano | findstr :8080

# İşlemi durdur (Windows)
taskkill /PID <PID> /F
```

### Oturum Süresi Doldu

**Çözüm**: VFS Global'de yeni oturum aç ve token'ları güncelle
```bash
# .env dosyasını edit et ve yeni değerleri gir
SESSION_COOKIE=new_value...
```

## 📊 Sistem Monitoring

### Sağlık Kontrolü

```bash
# Tüm servisler
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

### İstatistikler

```bash
# API istatistikleri
curl http://localhost:8080/api/stats

# Proxy istatistikleri
curl http://localhost:8082/proxies/stats
```

### Log Dosyalarını İzle

```bash
# Günlük log
type randevu_tracker_2026-09-01.log

# Slot kontrol log
type slot_checks_2026-09-01.log

# Real-time (Linux)
tail -f randevu_tracker_*.log
```

## 🔐 Güvenlik

### .env Dosyası

- ⚠️ `.env`'i asla Git'e commit etme
- `.gitignore`'da `.env` olmalı
- Kimlik bilgilerini güvenli tut

### Docker

- Official base images kullan
- Image'ları düzenli güncelle
- Network isolation kullan
- Resource limits belirle

### API

- Rate limiting ekle
- Authentication implement et
- HTTPS kullan (production)
- Input validation

## 📈 Performans İyileştirmeleri

### Memory

- Worker kontrol aralığını artır (daha az sık kontrol)
- Session history sınırını azalt
- Proxy pool boyutunu azalt

### CPU

- Goroutine sayısını sınırla
- Batch işlemleri kullan
- Caching implement et

### Network

- Connection pooling
- Keep-alive bağlantılar
- Gzip compression

## 🎯 Best Practices

1. **Üretim Konfigürasyonu**
   - `.env` dosyası için güvenli depo kullan (vault, secrets manager)
   - Resource limits belirle
   - Logging ve monitoring konfigüre et

2. **Backup ve Recovery**
   - Session data'yı yedekle
   - Log dosyalarını arşivle
   - Disaster recovery planı yap

3. **Updates**
   - Go versiyonunu güncel tut
   - Dependencies'i regular olarak güncelle
   - Breaking changes'ler için test et

## 📞 Destek

Daha fazla bilgi için:
- README.md - Genel bilgi
- QUICKSTART.md - Hızlı başlangıç
- MICROSERVICES.md - Mikroservis mimarisi

---

**Tarih**: 2026-09-01  
**Versiyon**: 1.0
