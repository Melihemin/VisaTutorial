package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

// ServiceOrchestrator - Tüm servisleri yönet
type ServiceOrchestrator struct {
	apiService     *APIService
	workerService  *WorkerService
	sessionService *SessionService
	proxyService   *ProxyService
	config         *Config
	isRunning      bool
}

// NewServiceOrchestrator - Yeni orchestrator oluştur
func NewServiceOrchestrator(config *Config) *ServiceOrchestrator {
	proxyPool := NewProxyPool(config.ProxyList)
	sessionMgr := NewSessionManager()

	return &ServiceOrchestrator{
		apiService:     NewAPIService("8080"),
		workerService:  NewWorkerService("worker-1", config, sessionMgr, proxyPool),
		sessionService: NewSessionService("8081"),
		proxyService:   NewProxyService("8082", proxyPool),
		config:         config,
	}
}

// Start - Tüm servisleri başlat
func (so *ServiceOrchestrator) Start() {
	so.isRunning = true

	log.Println("╔════════════════════════════════════════════════════════════╗")
	log.Println("║   Randevu Uygunluk Takip Sistemi - Mikroservis Mimarisi   ║")
	log.Println("╚════════════════════════════════════════════════════════════╝")
	log.Println()

	// Graceful shutdown için signal handler
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// API Servisi
	log.Println("📡 API Service başlatılıyor (Port: 8080)")
	log.Println("   - GET  /health")
	log.Println("   - GET  /api/sessions")
	log.Println("   - POST /api/sessions/create")
	log.Println("   - DELETE /api/sessions/delete")
	log.Println("   - POST /api/slots/check")
	log.Println("   - GET  /api/slots/latest")
	log.Println()

	go func() {
		if err := so.apiService.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Printf("API Service hatası: %v", err)
		}
	}()

	// Worker Servisi
	log.Println("👷 Worker Service başlatılıyor")
	log.Printf("   - Kontrol aralığı: %d saniye", so.config.CheckInterval)
	so.workerService.Start()
	log.Println()

	// Session Servisi
	log.Println("🔐 Session Service başlatılıyor (Port: 8081)")
	log.Println("   - GET  /health")
	log.Println("   - GET  /sessions")
	log.Println("   - GET  /sessions/validate")
	log.Println("   - GET  /sessions/renew")
	log.Println()

	go func() {
		if err := so.sessionService.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Printf("Session Service hatası: %v", err)
		}
	}()

	// Proxy Servisi
	log.Println("🌐 Proxy Service başlatılıyor (Port: 8082)")
	log.Println("   - GET  /health")
	log.Println("   - GET  /proxies")
	log.Println("   - GET  /proxies/stats")
	log.Println("   - GET  /proxies/rotate")
	log.Println("   - GET  /proxies/test")
	log.Println()

	go func() {
		if err := so.proxyService.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Printf("Proxy Service hatası: %v", err)
		}
	}()

	log.Println("════════════════════════════════════════════════════════════")
	log.Println("✅ Tüm servisler başarıyla başlatıldı!")
	log.Println("════════════════════════════════════════════════════════════")
	log.Println()
	log.Println("API Endpoints:")
	log.Println("  API Service:     http://localhost:8080/api/slots/latest")
	log.Println("  Session Service: http://localhost:8081/health")
	log.Println("  Proxy Service:   http://localhost:8082/proxies/stats")
	log.Println()
	log.Println("Kapatmak için: Ctrl+C")
	log.Println()

	// Signal bekle
	<-sigChan
	log.Println("\n🛑 Sistemi kapatıyorum...")
	so.Stop()
}

// Stop - Tüm servisleri durdur
func (so *ServiceOrchestrator) Stop() {
	if !so.isRunning {
		return
	}

	so.isRunning = false

	log.Println("Worker Service durduruluyor...")
	so.workerService.Stop()
	time.Sleep(1 * time.Second)

	log.Println("Servisler durduruldu.")
	log.Println("Çıkılıyor...")
	os.Exit(0)
}

// PrintServiceStatus - Servis durumunu yazdır
func (so *ServiceOrchestrator) PrintServiceStatus() {
	fmt.Println("\n📊 Servis Durumları:")
	fmt.Println("────────────────────────────────────")

	fmt.Println("API Service:")
	fmt.Printf("  Port: 8080\n")
	fmt.Printf("  Status: Running\n\n")

	fmt.Println("Worker Service:")
	stats := so.workerService.GetStats()
	fmt.Printf("  %+v\n\n", stats)

	fmt.Println("Session Service:")
	fmt.Printf("  Port: 8081\n")
	fmt.Printf("  Status: Running\n\n")

	fmt.Println("Proxy Service:")
	fmt.Printf("  Port: 8082\n")
	fmt.Printf("  Status: Running\n\n")
}

func init() {
	_ = godotenv.Load()
}

// main - Ana entry point
func init() {
	// Bu init fonksiyonu sadece microservices mimarisini tanıtmak için
	// Ana program main.go'da çalıştırılır
}
