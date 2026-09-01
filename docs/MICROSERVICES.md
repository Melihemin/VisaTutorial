# Randevu Uygunluk Takip Sistemi - Mikroservis Mimarisi

Bu klasör, Go dilinde geliştirilmiş mikroservis mimarisine sahip bir randevu takip sistemi içerir.

## 🏗️ Sistem Mimarisi

```
┌──────────────────────────────────────────────────────────────┐
│                    API Gateway (8080)                         │
│         - Slot kontrolü, Session yönetimi, Stats             │
└────────────────────┬─────────────────────────────────────────┘
                     │
        ┌────────────┼────────────┐
        │            │            │
        ▼            ▼            ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│   Worker     │ │   Session    │ │    Proxy     │
│  Service     │ │   Service    │ │   Service    │
│ (Slot Check) │ │  (8081)      │ │   (8082)     │
└──────────────┘ └──────────────┘ └──────────────┘
```

## 📋 Servisler

### 1. **API Service** (Port: 8080)
- REST API gatewayi
- Slot kontrolü istekleri
- Oturum yönetimi endpoints
- İstatistikler

**Endpoints:**
- `GET /health` - Sağlık kontrolü
- `GET /api/sessions` - Oturumları listele
- `POST /api/sessions/create` - Oturum oluştur
- `DELETE /api/sessions/delete` - Oturumu sil
- `POST /api/slots/check` - Slot kontrolü yap
- `GET /api/slots/latest` - Son kontrolü al
- `GET /api/proxy/status` - Proxy durumu
- `GET /api/stats` - Genel istatistikler

### 2. **Worker Service**
- Arka plan işçi
- Periyodik slot kontrolleri
- Oturum monitörü
- Sonuç arşivleme

**Özellikler:**
- Yapılandırılabilir kontrol aralığı
- Otomatik retry mekanizması
- Proxy rotasyonu
- İstatistik tutma

### 3. **Session Service** (Port: 8081)
- Oturum yaşam döngüsü yönetimi
- Session doğrulaması
- Oturum yenileme

**Endpoints:**
- `GET /health` - Sağlık kontrolü
- `GET /sessions` - Oturum sayısı
- `GET /sessions/validate?id=X` - Oturum doğrula
- `GET /sessions/renew?id=X` - Oturumu yenile

### 4. **Proxy Service** (Port: 8082)
- Proxy havuzu yönetimi
- Sticky IP yönetimi
- Proxy rotasyonu
- Proxy test ve monitör

**Endpoints:**
- `GET /health` - Sağlık kontrolü
- `GET /proxies` - Proxy listesi
- `GET /proxies/stats` - Proxy istatistikleri
- `GET /proxies/rotate` - Proxy rotasyonu
- `GET /proxies/test?url=X` - Proxy test

## 🚀 Hızlı Başlangıç

### Docker Compose ile

```bash
# Tüm servisleri başlat
docker-compose up -d

# Logları izle
docker-compose logs -f

# Servisleri durdur
docker-compose down
```

### Doğrudan Go ile

```bash
# Bağımlılıkları yükle
go mod download

# Program çalıştır
go run main.go orchestrator.go proxy.go config.go \
         notifier.go utils.go api_service.go \
         worker_service.go session_service.go proxy_service.go
```

## 📊 API Örnekleri

### Sağlık Kontrolü

```bash
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health
```

### Oturum Oluştur

```bash
curl -X POST http://localhost:8080/api/sessions/create \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "session-1",
    "proxy_url": "http://proxy:8080",
    "headers": {}
  }'
```

### Slot Kontrol Et

```bash
curl -X POST http://localhost:8080/api/slots/check \
  -H "Content-Type: application/json" \
  -d '{"session_id": "session-1"}'
```

### Proxy Listesi

```bash
curl http://localhost:8082/proxies
```

### Proxy İstatistikleri

```bash
curl http://localhost:8082/proxies/stats
```

## 🔧 Konfigürasyon

`.env` dosyasında ayarlar yapılır:

```env
# VFS Global
SESSION_COOKIE=...
AUTHORIZE_TOKEN=...
CLIENT_SOURCE_TOKEN=...

# Worker
CHECK_INTERVAL=45
REQUEST_TIMEOUT=30
MAX_RETRIES=3

# Proxy
PROXY_LIST=http://proxy1:8080,http://proxy2:8080

# Telegram (Bildirim)
TELEGRAM_BOT_TOKEN=...
TELEGRAM_CHAT_ID=...
```

## 📈 İzleme ve Logging

Her servis log dosyaları oluşturur:
- `randevu_tracker_YYYY-MM-DD.log` - Genel loglar
- `slot_checks_YYYY-MM-DD.log` - Slot kontrol sonuçları

Docker ile:
```bash
docker-compose logs -f api
docker-compose logs -f worker
docker-compose logs -f session
docker-compose logs -f proxy
```

## 🔐 Güvenlik

- Servisleri ayrı containerlarda çalıştır
- Environment variables kullan (kimlik bilgileri)
- Network isolation (Docker network)
- Health checks ile otomatik restart
- TLS mutual authentication (opsiyonel)

## 🐛 Sorun Giderme

### Servis başlamıyor
```bash
docker-compose logs proxy
```

### Bağlantı hatası
```bash
docker network ls
docker network inspect randevu_randevu-network
```

### Yüksek CPU/Bellek
- Worker interval'ını artır
- Proxy pool boyutunu azalt
- Session geçmişi sınırını kontrol et

## 📄 Dosya Yapısı

```
.
├── main.go                 # Asıl slot takip (README'de tanımlanmış)
├── proxy.go               # Proxy yönetimi
├── config.go              # Konfigürasyon
├── notifier.go            # Bildirim sistemi (README'de tanımlanmış)
├── utils.go               # Yardımcı fonksiyonlar
│
├── orchestrator.go        # Servis orkestrasyonu
├── api_service.go         # REST API servisi
├── worker_service.go      # Arka plan işçi
├── session_service.go     # Oturum servisi
├── proxy_service.go       # Proxy servisi
│
├── docker-compose.yml     # Kontainer orkestrasyonu
├── Dockerfile.api         # API servisi container
├── Dockerfile.worker      # Worker servisi container
├── Dockerfile.session     # Session servisi container
├── Dockerfile.proxy       # Proxy servisi container
│
├── .env                   # Konfigürasyon (gizli)
├── .env.example           # Örnek konfigürasyon
├── go.mod                 # Go modülleri
└── README.md             # Asıl proje dokümanı
```

## 🎯 Gelecek Geliştirmeler

- [ ] gRPC servisleri arası haberleşme
- [ ] Message queue (RabbitMQ) entegrasyonu
- [ ] Database backend (PostgreSQL)
- [ ] Kubernetes deployment
- [ ] Prometheus metrikleri
- [ ] Jaeger trace
- [ ] Circuit breaker pattern

## 📞 İletişim

Sorular ve öneriler için proje sahibine ulaşın.

---

**Versyon**: 1.0 Mikroservisler  
**Tarih**: 2026-09-01  
**Durum**: ✅ Aktif
