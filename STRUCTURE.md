# Randevu Uygunluk Takip Sistemi - Proje Yapısı

Temiz ve organize edilmiş Go microservices projesi.

## 📁 Dizin Yapısı

```
randevu_tracker/
│
├── 📄 Kök Dosyalar (Root Level)
│   ├── go.mod                 # Go modülü tanımı
│   ├── go.sum                 # Bağımlılıkları kontrol
│   ├── Makefile               # Build hedefleri ve komutları
│   ├── build.sh               # Linux/macOS build betiği
│   ├── build.bat              # Windows build betiği
│   └── .gitignore             # Git ignore kuralları
│
├── 📂 src/                    # Go Kaynak Kodları
│   ├── main.go                # Ana uygulama, slot tracking
│   ├── proxy.go               # Proxy yönetimi
│   ├── config.go              # Konfigürasyon yükleme
│   ├── notifier.go            # Telegram bildirimleri
│   ├── utils.go               # Logger, istatistikler
│   ├── orchestrator.go        # Servis orkestrasyonu
│   ├── api_service.go         # REST API servisi (8080)
│   ├── worker_service.go      # Arka plan işçi
│   ├── session_service.go     # Oturum yönetimi (8081)
│   └── proxy_service.go       # Proxy servisi (8082)
│
├── 📂 config/                 # Konfigürasyon Dosyaları
│   ├── .env                   # Gizli kimlik bilgileri (local)
│   └── .env.example           # Örnek konfigürasyon (public)
│
├── 📂 docker/                 # Docker Dosyaları
│   ├── docker-compose.yml     # Multi-container orchestration
│   ├── Dockerfile.api         # API servisi container
│   ├── Dockerfile.worker      # Worker servisi container
│   ├── Dockerfile.session     # Session servisi container
│   ├── Dockerfile.proxy       # Proxy servisi container
│   └── .dockerignore          # Docker ignore kuralları
│
├── 📂 docs/                   # Belgelendirme
│   ├── README.md              # Asıl proje dokümantasyonu
│   ├── QUICKSTART.md          # Hızlı başlangıç kılavuzu
│   ├── MICROSERVICES.md       # Mikroservis mimarisi
│   ├── ARCHITECTURE.md        # Sistem mimarisi detayları
│   ├── DEPLOYMENT.md          # Dağıtım kılavuzu
│   └── INDEX.md               # Belgelendirme harita
│
└── 📂 build/                  # Derlenmiş Çıktılar
    └── randevu_tracker.exe    # Ana binary (Windows)
    
    📝 Not: Linux/macOS için 'randevu_tracker' (executable)
```

## 🚀 Hızlı Başlangıç

### 1. Kurulum

```bash
# Windows PowerShell
./build.bat

# Linux/macOS Bash
chmod +x build.sh
./build.sh

# Veya Make ile (tüm platformlar)
make build
```

### 2. Konfigürasyon

```bash
# .env dosyasını oluştur
cp config/.env.example config/.env

# Metin editörü ile aç ve kimlik bilgilerini gir
nano config/.env
```

### 3. Çalıştırma

```bash
# Windows
./build/randevu_tracker.exe

# Linux/macOS
./build/randevu_tracker

# Veya Make ile
make run
```

## 📦 Klasör Açıklamaları

### `src/` - Go Kaynak Kodları
- Tüm Go uygulama kodları burada
- 10 Go modülü (main, proxy, config, vb.)
- Microservices mimarisi

### `config/` - Konfigürasyon
- `.env` - Gizli kimlik bilgileri (git'e eklenmez)
- `.env.example` - Örnek config (public)

### `docker/` - Containerization
- Docker Compose ile 4 servis orkestrasyonu
- Multi-stage builds ile minimal image
- Health checks ve networking

### `docs/` - Belgelendirme
- **README.md** - Proje tanıtımı ve teknik mimari
- **QUICKSTART.md** - 5 dakika kurulum
- **MICROSERVICES.md** - Servis detayları
- **ARCHITECTURE.md** - Sistem tasarımı diyagramları
- **DEPLOYMENT.md** - Production dağıtım
- **INDEX.md** - Belgelendirme haritası

### `build/` - Çıktılar
- Derlenmiş binary dosyaları
- Loglar (runtime)

## 🔧 Build Komutları

### Make ile (Tavsiye Edilen)

```bash
make help              # Tüm hedefleri göster
make build             # Derle
make run               # Çalıştır
make fmt               # Formatla
make lint              # Lint
make clean             # Temizle
make docker-build      # Docker image derle
make docker-up         # Docker başlat
make docker-down       # Docker durdur
```

### PowerShell (Windows)

```powershell
# Derle
./build.bat

# Çalıştır
./build/randevu_tracker.exe

# Docker
docker-compose -f docker/docker-compose.yml up -d
```

### Bash (Linux/macOS)

```bash
# Derle
chmod +x build.sh
./build.sh

# Çalıştır
./build/randevu_tracker

# Docker
cd docker && docker-compose up -d
```

## 🌐 Servisler

- **API Service** - Port 8080 - REST API Gateway
- **Worker Service** - Background - Slot kontrolleri
- **Session Service** - Port 8081 - Oturum yönetimi
- **Proxy Service** - Port 8082 - Proxy havuz

## 📊 Proje İstatistikleri

- **Go Dosyaları**: 10 + test (1500+ satır)
- **Docker Containers**: 4 (api, worker, session, proxy)
- **Belgelendirme**: 6 dosya (1400+ satır)
- **API Endpoints**: 15+
- **Microservices**: 4 farklı servis

## ✨ Özellikler

✅ Microservices mimarisi  
✅ Docker containerization  
✅ REST API  
✅ Health checks  
✅ Modüler tasarım  
✅ Comprehensive documentation  
✅ Build automation  
✅ Production ready  

## 📖 Dokümantasyon Haritası

1. **İlk Kez Bakıyor?** → `docs/README.md`
2. **Hızlı Kur?** → `docs/QUICKSTART.md`
3. **Servisler?** → `docs/MICROSERVICES.md`
4. **Mimari?** → `docs/ARCHITECTURE.md`
5. **Deploy?** → `docs/DEPLOYMENT.md`
6. **Harita?** → `docs/INDEX.md`

## 🔐 Güvenlik

- `.env` dosyası .gitignore'da
- Kimlik bilgileri environment variables'de
- Container isolation (Docker)
- TLS maskeleme (Chrome profili)

## 🚀 Next Steps

1. `config/.env` dosyasını güncelleyin
2. VFS Global'den token'ları alın
3. `make build` ile derleyin
4. `make run` ile başlatın

---

**Yapı Tarihı**: 2026-09-01  
**Versiyon**: 1.0  
**Durum**: ✅ Production Ready
