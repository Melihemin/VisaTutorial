package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// ProxyService - Proxy yönetimi servisi
type ProxyService struct {
	port      string
	proxyPool *ProxyPool
	router    *http.ServeMux
	mu        sync.RWMutex
}

// NewProxyService - Yeni proxy servisi oluştur
func NewProxyService(port string, proxyPool *ProxyPool) *ProxyService {
	svc := &ProxyService{
		port:      port,
		proxyPool: proxyPool,
		router:    http.NewServeMux(),
	}

	svc.setupRoutes()
	return svc
}

// setupRoutes - Rotaları konfigüre et
func (ps *ProxyService) setupRoutes() {
	ps.router.HandleFunc("/health", ps.healthHandler)
	ps.router.HandleFunc("/proxies", ps.listProxiesHandler)
	ps.router.HandleFunc("/proxies/stats", ps.statsHandler)
	ps.router.HandleFunc("/proxies/rotate", ps.rotateHandler)
	ps.router.HandleFunc("/proxies/test", ps.testProxyHandler)
}

// healthHandler - Sağlık kontrolü
func (ps *ProxyService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","service":"proxy-service"}`)
}

// listProxiesHandler - Proxy listele
func (ps *ProxyService) listProxiesHandler(w http.ResponseWriter, r *http.Request) {
	ps.proxyPool.mu.RLock()
	defer ps.proxyPool.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"total_proxies":%d,"proxies":[`, len(ps.proxyPool.proxies))

	for i, proxy := range ps.proxyPool.proxies {
		if i > 0 {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, `{"url":"%s","active":%v,"fail_count":%d}`,
			proxy.URL, proxy.IsActive, proxy.FailCount)
	}

	fmt.Fprintf(w, `]}`)
}

// statsHandler - Proxy istatistikleri
func (ps *ProxyService) statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := ps.proxyPool.GetPoolStats()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"stats":%+v,"timestamp":"%s"}`,
		stats, time.Now().Format("2006-01-02 15:04:05"))
}

// rotateHandler - Proxy rotasyonu
func (ps *ProxyService) rotateHandler(w http.ResponseWriter, r *http.Request) {
	proxy, err := ps.proxyPool.GetNextProxy()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}

	fmt.Fprintf(w, `{"current_proxy":"%s","is_active":%v}`,
		proxy.URL, proxy.IsActive)
}

// testProxyHandler - Proxy test et
func (ps *ProxyService) testProxyHandler(w http.ResponseWriter, r *http.Request) {
	proxyURL := r.URL.Query().Get("url")
	if proxyURL == "" {
		http.Error(w, "url parameter required", http.StatusBadRequest)
		return
	}

	// Simple connectivity test
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://www.google.com")

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Fprintf(w, `{"proxy":"%s","working":false,"error":"%s"}`,
			proxyURL, err.Error())
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, `{"proxy":"%s","working":true,"status_code":%d}`,
		proxyURL, resp.StatusCode)
}

// Start - Servisi başlat
func (ps *ProxyService) Start() error {
	log.Printf("Proxy Service başlatılıyor: :%s", ps.port)
	return http.ListenAndServe(":"+ps.port, ps.router)
}

func init() {
	_ = godotenv.Load()
}
