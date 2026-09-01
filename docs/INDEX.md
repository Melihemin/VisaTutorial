# 📚 Proje Belgeleri ve Kılavuzlar

## 📖 Tüm Dokümantasyon Dosyaları

Bu projede aşağıdaki belge dosyaları bulunmaktadır:

### 1. **README.md** (Asıl Proje Dokümanı)
   - **Amaç**: Projenin genel tanıtımı ve teknik mimarisi
   - **İçerik**: 
     - Proje özeti ve problem tanımı
     - Teknoloji seçimi (Go programlama dili)
     - Genel sistem mimarisi
     - Bileşenlerin açıklamaları
     - VFS Global API entegrasyonu
   - **Okuyucular**: Teknik yöneticiler, geliştiriciler
   - **✅ KORUNUYOR**: Orijinal içerik değiştirilmedi

### 2. **QUICKSTART.md** (Hızlı Başlangıç)
   - **Amaç**: 5 dakika içinde sistemi çalıştır
   - **İçerik**:
     - Adım adım kurulum kılavuzu
     - Token toplama yöntemi
     - .env dosyası konfigürasyonu
     - Telegram bildirimi kurulumu
   - **Okuyucular**: Yeni kullanıcılar
   - **Dilleri**: Türkçe + İngilizce

### 3. **MICROSERVICES.md** (Mikroservis Mimarisi)
   - **Amaç**: Servis mimarisinin detaylı anlatımı
   - **İçerik**:
     - Servis bileşenleri (4 farklı servis)
     - API endpoints
     - Docker Compose kullanımı
     - Health checks
   - **Okuyucular**: DevOps, sistem mimarları
   - **Yeni Eklenmiş**: Bu belgeler projeye eklenmiştir

### 4. **ARCHITECTURE.md** (Sistem Mimarisi Detayları)
   - **Amaç**: Derin teknik mimari belgeleme
   - **İçerik**:
     - Sistem diyagramları
     - Katmanlı mimari
     - Service-to-service iletişimi
     - Veri akış diyagramları
     - Güvenlik modeli
     - Ölçeklenebilirlik stratejileri
   - **Okuyucular**: Mimari tasarımcılar, lead developer
   - **Yeni Eklenmiş**: Bu belgeler projeye eklenmiştir

### 5. **DEPLOYMENT.md** (Dağıtım Kılavuzu)
   - **Amaç**: Üretime deployment işlemleri
   - **İçerik**:
     - Gereksinimler (Go, Docker)
     - Doğrudan Go ile dağıtım
     - Docker Compose ile dağıtım
     - Konfigürasyon ve Best Practices
     - Sorun giderme
   - **Okuyucular**: DevOps mühendisleri
   - **Yeni Eklenmiş**: Bu belgeler projeye eklenmiştir

---

## 📁 Proje Dosya Yapısı

```
randevu_tracker/
│
├─ 📄 Konfigürasyon Dosyaları
│  ├── .env                    # Gizli kimlik bilgileri (local)
│  ├── .env.example            # Örnek konfigürasyon (public)
│  ├── go.mod                  # Go bağımlılıkları
│  └── go.sum                  # Bağımlılık kontrol
│
├─ 💻 Go Kaynak Kodu
│  ├── main.go                 # Ana uygulama, slot tracking
│  ├── proxy.go                # Proxy yönetimi
│  ├── config.go               # Konfigürasyon yükleme
│  ├── notifier.go             # Telegram bildirimleri
│  ├── utils.go                # Logger, istatistikler
│  │
│  ├── orchestrator.go         # Servis orkestrasyonu (NEW)
│  ├── api_service.go          # REST API servisi (NEW)
│  ├── worker_service.go       # Arka plan işçi (NEW)
│  ├── session_service.go      # Oturum yönetimi (NEW)
│  ├── proxy_service.go        # Proxy servisi (NEW)
│  │
│  └── randevu_tracker.exe     # Derlenmiş binary
│
├─ 📄 Belgelendirme
│  ├── README.md               # Asıl proje dokü (KORUNUYOR)
│  ├── QUICKSTART.md           # Hızlı başlangıç (NEW)
│  ├── MICROSERVICES.md        # Mikroservis mimarisi (NEW)
│  ├── ARCHITECTURE.md         # Mimari detayları (NEW)
│  ├── DEPLOYMENT.md           # Dağıtım kılavuzu (NEW)
│  └── INDEX.md                # Bu dosya (NEW)
│
├─ 🐳 Docker Dosyaları
│  ├── docker-compose.yml      # Multi-container orchestration (NEW)
│  ├── Dockerfile.api          # API servisi container (NEW)
│  ├── Dockerfile.worker       # Worker servisi container (NEW)
│  ├── Dockerfile.session      # Session servisi container (NEW)
│  ├── Dockerfile.proxy        # Proxy servisi container (NEW)
│  └── .dockerignore           # Docker ignore kuralları (NEW)
│
└─ 🛠️ Build ve Tools
   └── Makefile                # Build hedefleri (NEW)
```

---

## 🎯 Hangi Belgeyi Ne Zaman Oku?

### Yeni başlayanlar için:
1. README.md - Projeyi anla
2. QUICKSTART.md - Hızlı kur ve çalıştır
3. .env.example - Gerekli kimlik bilgilerini hazırla

### Geliştirici için:
1. README.md - Teknik altyapı
2. ARCHITECTURE.md - Sistem tasarımı
3. MICROSERVICES.md - Servis detayları
4. İlgili `.go` dosyası

### DevOps/SysAdmin için:
1. DEPLOYMENT.md - Kurulum adımları
2. MICROSERVICES.md - Servis harita
3. docker-compose.yml - Container yapısı
4. Makefile - Build komutları

### Yönetim/PO için:
1. README.md - Proje özeti (ilk 20 satır)
2. ARCHITECTURE.md - Sistem diyagramları
3. DEPLOYMENT.md - Production requirements

---

## 🔧 Komut Hızlı Referansı

```bash
# Development
make build                    # Derle
make run                      # Çalıştır
make fmt                      # Formatla
make lint                     # Lint çalıştır

# Docker
make docker-build             # Image derle
make docker-up                # Kontainerları başlat
make docker-down              # Kontainerları durdur
make docker-logs              # Logları izle

# Kurulum
make env                      # .env dosyası oluştur
```

---

## 📊 İstatistikler

### Go Kaynak Kodu
| Dosya | Satırlar | Amaç |
|-------|----------|------|
| main.go | ~300 | Ana uygulama |
| proxy.go | ~200 | Proxy yönetimi |
| config.go | ~150 | Konfigürasyon |
| notifier.go | ~100 | Bildirimler |
| utils.go | ~150 | Yardımcı fonksiyonlar |
| api_service.go | ~200 | API servisi |
| worker_service.go | ~150 | Worker servisi |
| session_service.go | ~120 | Session servisi |
| proxy_service.go | ~120 | Proxy servisi |
| orchestrator.go | ~150 | Orkestrasyonu |
| **TOPLAM** | **~1500+** | **10 Go dosyası** |

### Belgelendirme
- README.md: 250+ satır (KORUNUYOR)
- QUICKSTART.md: 150+ satır
- MICROSERVICES.md: 250+ satır
- ARCHITECTURE.md: 400+ satır
- DEPLOYMENT.md: 300+ satır

### Docker Dosyaları
- 1x docker-compose.yml (5 servis)
- 4x Dockerfile (multi-stage builds)
- 1x .dockerignore

---

## 🚀 Microservis Mimarisi Özeti

### Servisler
1. **API Service** (8080) - REST API Gateway
2. **Worker Service** - Periyodik slot kontrolleri
3. **Session Service** (8081) - Oturum yönetimi
4. **Proxy Service** (8082) - Proxy havuz yönetimi

### Özellikler
✅ Horizontal scaling  
✅ Service isolation  
✅ Health checks  
✅ Docker containerization  
✅ Graceful shutdown  
✅ Error handling  

---

## 📌 Önemli Notlar

### README.md Korunması
- ✅ Asıl README.md değiştirilmedi
- ✅ Orijinal içerik tamamen korundu
- ✅ Bildirim kısmı dokunulmadı
- ✅ Yeni belgelendirme ayrı dosyalarda

### Microservis Mimarisi
- ✅ Modüler tasarım
- ✅ Ayrı sorumluluk alanları
- ✅ Kolay bakım ve güncellenme
- ✅ Üretime hazır

### Geliştirme Kolaylıkları
- ✅ Makefile ile hızlı komutlar
- ✅ Docker Compose ile one-command deploy
- ✅ Detaylı belgeler her seviye için
- ✅ Health checks ve monitoring

---

## 📞 Dosya İçeriği Hızlı Bağlantılar

| Belge | Dosya | Amaç |
|-------|-------|------|
| Proje Tanıtımı | [README.md](README.md) | Teknik mimari ve özellikler |
| Hızlı Kurulum | [QUICKSTART.md](QUICKSTART.md) | 5 dakika başlangıç |
| Mikroservisler | [MICROSERVICES.md](MICROSERVICES.md) | Servis detayları |
| Mimari | [ARCHITECTURE.md](ARCHITECTURE.md) | Sistem tasarımı |
| Dağıtım | [DEPLOYMENT.md](DEPLOYMENT.md) | Production deployment |
| Bu Dosya | [INDEX.md](INDEX.md) | Belgelendirme harita |

---

## ✨ Yenilikler

### Eklenen Dosyalar (1.0 Microservices)
- ✨ api_service.go - REST API gateway
- ✨ worker_service.go - Background worker
- ✨ session_service.go - Session management
- ✨ proxy_service.go - Proxy pool manager
- ✨ orchestrator.go - Service orchestration
- ✨ docker-compose.yml - Container orchestration
- ✨ 4x Dockerfile - Container definitions
- ✨ Makefile - Build automation
- ✨ MICROSERVICES.md - Servis belgelendirmesi
- ✨ ARCHITECTURE.md - Mimari belgelendirmesi
- ✨ DEPLOYMENT.md - Dağıtım kılavuzu
- ✨ .dockerignore - Docker ignore rules

### Korunan Dosyalar
- ✅ README.md - Orijinal içerik korundu
- ✅ main.go - Ana uygulama güncellendi
- ✅ proxy.go - Proxy yönetimi
- ✅ config.go - Konfigürasyon
- ✅ notifier.go - Bildirim sistemi
- ✅ utils.go - Yardımcı fonksiyonlar

---

**Belgelendirme Tarihı**: 2026-09-01  
**Versiyon**: 1.0 Microservices  
**Durum**: ✅ Prodüksüyüne Hazır
