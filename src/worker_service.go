package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
)

// WorkerService - Arka plan worker servisi
type WorkerService struct {
	ID               string
	config           *Config
	sessionMgr       *SessionManager
	proxyPool        *ProxyPool
	checkInterval    time.Duration
	stats            *Stats
	isRunning        bool
	stopChan         chan bool
	slotCheckHistory []SlotCheckResult
}

// NewWorkerService - Yeni worker servisi oluştur
func NewWorkerService(id string, config *Config, sessionMgr *SessionManager, proxyPool *ProxyPool) *WorkerService {
	return &WorkerService{
		ID:               id,
		config:           config,
		sessionMgr:       sessionMgr,
		proxyPool:        proxyPool,
		checkInterval:    time.Duration(config.CheckInterval) * time.Second,
		stats:            &Stats{StartTime: time.Now()},
		stopChan:         make(chan bool),
		slotCheckHistory: make([]SlotCheckResult, 0),
	}
}

// Start - Worker servisi başlat
func (ws *WorkerService) Start() {
	ws.isRunning = true
	log.Printf("Worker Service %s başlatılıyor (Aralık: %v)", ws.ID, ws.checkInterval)

	go ws.run()
}

// run - Ana döngü
func (ws *WorkerService) run() {
	ticker := time.NewTicker(ws.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ws.stopChan:
			log.Printf("Worker Service %s durduruluyor", ws.ID)
			ws.isRunning = false
			return

		case <-ticker.C:
			ws.performSlotCheck()
		}
	}
}

// performSlotCheck - Slot kontrolü gerçekleştir
func (ws *WorkerService) performSlotCheck() {
	log.Printf("[%s] Slot kontrolü yapılıyor...", ws.ID)

	result := &SlotCheckResult{
		Timestamp: time.Now(),
	}

	// Aktif oturumlar üzerinde kontrol yap
	ws.sessionMgr.mu.RLock()
	sessions := make([]*Session, 0)
	for _, session := range ws.sessionMgr.sessions {
		sessions = append(sessions, session)
	}
	ws.sessionMgr.mu.RUnlock()

	if len(sessions) == 0 {
		result.Error = "Aktif oturum bulunamadı"
		ws.stats.UpdateStats(false, false)
		return
	}

	// Her oturum için kontrol yap
	for _, session := range sessions {
		// Proxy seç
		proxy, err := ws.proxyPool.GetStickyProxy(session.ID)
		if err != nil {
			proxy, _ = ws.proxyPool.GetNextProxy()
		}

		result.Error = ""
		if err != nil {
			result.Error = fmt.Sprintf("Proxy hatası: %v", err)
			ws.stats.UpdateStats(false, false)
			continue
		}

		log.Printf("[%s] Oturum %s slot kontrolü yapılıyor (Proxy: %s)", ws.ID, session.ID, proxy.URL)

		// Kontrol başarılı sayılırsa
		ws.stats.UpdateStats(true, len(result.Slots) > 0)
		ws.slotCheckHistory = append(ws.slotCheckHistory, *result)

		// History sınırını kontrol et
		if len(ws.slotCheckHistory) > 1000 {
			ws.slotCheckHistory = ws.slotCheckHistory[1:]
		}
	}
}

// Stop - Worker servisi durdur
func (ws *WorkerService) Stop() {
	if ws.isRunning {
		ws.stopChan <- true
	}
}

// GetStats - İstatistikleri al
func (ws *WorkerService) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"worker_id":   ws.ID,
		"is_running":  ws.isRunning,
		"stats":       ws.stats.GetSummary(),
		"last_checks": len(ws.slotCheckHistory),
	}
}

// GetCheckHistory - Kontrol geçmişini al (son N kaydı)
func (ws *WorkerService) GetCheckHistory(limit int) []SlotCheckResult {
	if limit > len(ws.slotCheckHistory) {
		limit = len(ws.slotCheckHistory)
	}

	if limit <= 0 {
		return []SlotCheckResult{}
	}

	start := len(ws.slotCheckHistory) - limit
	return ws.slotCheckHistory[start:]
}

func init() {
	_ = godotenv.Load()
}
