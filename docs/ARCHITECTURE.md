# Randevu Uygunluk Takip Sistemi - Sistem Mimarisi

## 📐 Genel Mimari Diyagram

```
┌──────────────────────────────────────────────────────────────────────┐
│                        İstemci / Client                              │
│            (curl, Postman, Browser, Web UI)                         │
└──────────────────────────────────────┬───────────────────────────────┘
                                       │
                    ┌──────────────────▼──────────────────┐
                    │      API Gateway Service           │
                    │         (Port: 8080)               │
                    │  - REST API endpoints             │
                    │  - Request routing                │
                    │  - Response aggregation           │
                    └──────────────────┬──────────────────┘
                                       │
              ┌────────────────────────┼────────────────────────┐
              │                        │                        │
              ▼                        ▼                        ▼
    ┌──────────────────┐   ┌──────────────────┐   ┌──────────────────┐
    │  Worker Service  │   │ Session Service  │   │  Proxy Service   │
    │  (Background)    │   │   (Port: 8081)   │   │   (Port: 8082)   │
    │                  │   │                  │   │                  │
    │ - Slot checking  │   │ - Lifecycle mgmt │   │ - Pool management│
    │ - Scheduling     │   │ - Validation     │   │ - Rotation       │
    │ - Notifications  │   │ - Renewal        │   │ - Health check   │
    └────────┬─────────┘   └──────────────────┘   └──────────────────┘
             │
             │ Periyodik kontrol
             ▼
    ┌──────────────────────────┐
    │   VFS Global API         │
    │ https://lift-api...      │
    └──────────────────────────┘
```

## 🏢 Katmanlı Mimari

```
┌─────────────────────────────────────────────────────────┐
│              Presentation Layer                         │
│              (REST API - 8080)                          │
│  - HTTP endpoints                                       │
│  - JSON serialization                                   │
│  - Request validation                                   │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Application Layer                          │
│  - Service orchestration                                │
│  - Business logic                                       │
│  - Event handling                                       │
│  - Error handling                                       │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Service Layer                              │
│  - API Service (8080)                                   │
│  - Worker Service (background)                          │
│  - Session Service (8081)                               │
│  - Proxy Service (8082)                                 │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              Data/Manager Layer                         │
│  - SessionManager                                       │
│  - ProxyPool                                            │
│  - Config management                                    │
│  - Statistics tracking                                  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│              External Services                          │
│  - VFS Global API                                       │
│  - HTTP Proxies                                         │
│  - Telegram Bot API (optional)                          │
│  - File system (logs)                                   │
└─────────────────────────────────────────────────────────┘
```

## 📦 Servis İletişimi

### Service-to-Service Communication

```
┌─────────────┐
│  API SVC    │──────────────────┐
│ (8080)      │                  │
└─────────────┘                  │
       │                         │ REST calls
       │ REST API                │ (via HTTP)
       │ requests                │
       │                         │
       ▼                         ▼
┌─────────────────────────────────────────┐
│         Shared Memory                   │
│  - SessionManager                       │
│  - ProxyPool                            │
│  - Stats                                │
└─────────────────────────────────────────┘
       ▲                         ▲
       │                         │
       │                         │ REST calls
       │ REST API                │
       │ requests                │
       │                         │
┌─────────────┐          ┌──────────────┐
│ Worker SVC  │          │ Session SVC  │
│(background) │          │   (8081)     │
└─────────────┘          └──────────────┘

┌──────────────┐
│ Proxy SVC    │
│  (8082)      │
└──────────────┘
```

## 🔄 Veri Akışı

### Slot Kontrol Süreci

```
1. İstemci İstek Gönder
   │
   ├─→ POST /api/slots/check
   │   └─→ {session_id: "session-1"}
   │
2. API Servisi
   │
   ├─→ SessionManager'dan oturum al
   ├─→ ProxyPool'dan proxy seç
   └─→ Worker Servisi'ne kontrol işi gönder
   │
3. Worker Servisi
   │
   ├─→ Proxy ile TLS client oluştur
   ├─→ VFS Global API'ye istek gönder
   ├─→ Cevabı parse et
   └─→ İstatistikleri güncelle
   │
4. Notification (opsiyonel)
   │
   ├─→ Slot varsa Telegram bildirimi gönder
   └─→ Log dosyasına kaydet
   │
5. Yanıt İstemciye Dön
   │
   └─→ {status: "checking", timestamp: "..."}
```

### Oturum Yönetim Süreci

```
1. Oturum Oluştur
   │
   ├─→ POST /api/sessions/create
   │   └─→ {session_id: "s1", proxy_url: "...", headers: {...}}
   │
2. Session Manager
   │
   ├─→ Oturum nesnesi oluştur
   ├─→ ProxyPool'a sticky proxy ata
   └─→ Cookie'leri sakla
   │
3. Session Service'e Sync
   │
   └─→ Oturum bilgisini sakla
   │
4. Worker, oturumu aktif tutar
   │
   └─→ Periyodik kontrol sırasında last_activity güncelle
```

## 📊 Data Structures

### Session Structure
```go
type Session struct {
    ID           string
    Cookies      []*http.Cookie
    Headers      map[string]string
    ProxyURL     string
    LastActivity time.Time
    Client       tls_client.HttpClient
}
```

### Proxy Entry Structure
```go
type ProxyEntry struct {
    URL        string
    IsActive   bool
    FailCount  int
    LastUsed   time.Time
    Sticky     bool
    SessionID  string
}
```

### SlotCheckResult Structure
```go
type SlotCheckResult struct {
    Timestamp    time.Time
    EarliestDate string
    Slots        []SlotItem
    Error        string
}
```

## 🔌 API Endpoints Haritası

```
API SERVICE (8080)
├── /health
│   └── GET: Sağlık kontrolü
│
├── /api/sessions
│   ├── GET: Oturumları listele
│   ├── POST: Yeni oturum oluştur
│   └── DELETE: Oturumu sil
│
├── /api/slots
│   ├── POST /check: Slot kontrolü yap
│   └── GET /latest: Son kontrol sonucunu al
│
├── /api/proxy
│   └── GET /status: Proxy durumunu al
│
└── /api/stats
    └── GET: Genel istatistikler

SESSION SERVICE (8081)
├── /health
│   └── GET: Sağlık kontrolü
│
├── /sessions
│   └── GET: Oturum sayısı
│
├── /sessions/validate
│   └── GET: Oturum doğrulaması
│
└── /sessions/renew
    └── GET: Oturum yenileme

PROXY SERVICE (8082)
├── /health
│   └── GET: Sağlık kontrolü
│
├── /proxies
│   └── GET: Proxy listesi
│
├── /proxies/stats
│   └── GET: Istatistikler
│
├── /proxies/rotate
│   └── GET: Proxy rotasyonu
│
└── /proxies/test
    └── GET: Proxy test
```

## 🔐 Güvenlik Modeli

```
┌────────────────────────────────────────┐
│   Credential Management                │
│                                        │
│ .env (local, .gitignore)              │
│   ├── SESSION_COOKIE                  │
│   ├── AUTHORIZE_TOKEN                 │
│   └── CLIENT_SOURCE_TOKEN             │
└────────────────────────────────────────┘
           │
           ▼
┌────────────────────────────────────────┐
│  Environment Isolation                 │
│                                        │
│  Docker Container   ≠   Host           │
│  - Isolated ENV       - Protected      │
│  - Network isolated   - Secure storage │
│  - Volume mounted     - .env ignored   │
└────────────────────────────────────────┘
           │
           ▼
┌────────────────────────────────────────┐
│  Service Communication                 │
│                                        │
│  Loopback network:                     │
│  - API Service     ← localhost:8080   │
│  - Session Service ← localhost:8081   │
│  - Proxy Service   ← localhost:8082   │
│  - Worker Service  ← background       │
└────────────────────────────────────────┘
```

## 📈 Ölçeklenebilirlik

### Horizontal Scaling

```
┌─────────────────────────────┐
│   Load Balancer             │
│  (nginx, traefik)           │
└─────────────┬───────────────┘
              │
    ┌─────────┼─────────┐
    │         │         │
    ▼         ▼         ▼
┌───────┐ ┌───────┐ ┌───────┐
│ API 1 │ │ API 2 │ │ API 3 │
└───────┘ └───────┘ └───────┘
    │         │         │
    └─────────┼─────────┘
              │
    ┌─────────┴─────────┐
    │                   │
    ▼                   ▼
┌──────────┐      ┌──────────┐
│ Worker 1 │      │ Worker N │
└──────────┘      └──────────┘
    │                   │
    └────────┬──────────┘
             │
        ┌────▼─────────────┐
        │ Shared Services  │
        │  - Database      │
        │  - Cache         │
        │  - Message Queue │
        └──────────────────┘
```

### Vertical Scaling

- Worker kontrol aralığını artır (daha az sık)
- Batch işlemleri kullan
- Connection pooling
- Caching layer ekle

## 🔄 Deployment Pipeline

```
Source Code
    │
    ▼
┌──────────────┐
│  Git Repo    │
│  (VCS)       │
└──────┬───────┘
       │
       ▼
┌──────────────────────────┐
│  CI/CD Pipeline          │
│  (GitHub Actions, etc)   │
│                          │
│  1. Build               │
│  2. Test                │
│  3. Lint                │
│  4. Security scan       │
└──────┬───────────────────┘
       │
       ▼
┌──────────────────────┐
│  Docker Registry     │
│  (DockerHub, etc)    │
└──────┬───────────────┘
       │
       ▼
┌──────────────────────┐
│  Production Env      │
│  (Kubernetes, etc)   │
│                      │
│  - Orchestration     │
│  - Auto-scaling      │
│  - Self-healing      │
│  - Updates           │
└──────────────────────┘
```

## 🎯 Tasarım Prensipleri

### 1. **Separation of Concerns**
- Her servis tek sorumluluğu
- Loose coupling
- High cohesion

### 2. **Stateless Services**
- State centralized (SessionManager, ProxyPool)
- Horizontal scaling mümkün
- No session affinity needed

### 3. **Resilience**
- Health checks
- Automatic restarts
- Graceful shutdown
- Error recovery

### 4. **Observability**
- Structured logging
- Metrics collection
- Health endpoints
- Error tracking

### 5. **Security**
- Credential isolation
- Network isolation (Docker)
- Input validation
- Rate limiting (future)

---

**Mimari Versiyon**: 1.0  
**Tarih**: 2026-09-01  
**Durum**: Production Ready
