package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// APIService - REST API servisi
type APIService struct {
	port          string
	sessionMgr    *SessionManager
	proxyPool     *ProxyPool
	slotChecker   *SlotChecker
	router        *http.ServeMux
	mu            sync.RWMutex
	lastSlotCheck *SlotCheckResult
}

// SlotChecker - Slot kontrol işlemleri
type SlotChecker struct {
	config *Config
	client interface{}
}

// NewAPIService - Yeni API servisi oluştur
func NewAPIService(port string) *APIService {
	svc := &APIService{
		port:       port,
		sessionMgr: NewSessionManager(),
		router:     http.NewServeMux(),
	}

	svc.setupRoutes()
	return svc
}

// setupRoutes - API rotalarını konfigüre et
func (as *APIService) setupRoutes() {
	// Health check
	as.router.HandleFunc("/health", as.healthHandler)

	// Session endpoints
	as.router.HandleFunc("/api/sessions", as.listSessionsHandler)
	as.router.HandleFunc("/api/sessions/create", as.createSessionHandler)
	as.router.HandleFunc("/api/sessions/delete", as.deleteSessionHandler)

	// Slot checking endpoints
	as.router.HandleFunc("/api/slots/check", as.checkSlotsHandler)
	as.router.HandleFunc("/api/slots/latest", as.latestSlotHandler)

	// Proxy endpoints
	as.router.HandleFunc("/api/proxy/status", as.proxyStatusHandler)

	// Stats endpoints
	as.router.HandleFunc("/api/stats", as.statsHandler)
}

// healthHandler - Sağlık kontrolü
func (as *APIService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"service":   "api-service",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// listSessionsHandler - Oturumları listele
func (as *APIService) listSessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	as.sessionMgr.mu.RLock()
	defer as.sessionMgr.mu.RUnlock()

	sessions := make([]map[string]interface{}, 0)
	for id, session := range as.sessionMgr.sessions {
		sessions = append(sessions, map[string]interface{}{
			"id":            id,
			"proxy_url":     session.ProxyURL,
			"last_activity": session.LastActivity.Format("2006-01-02 15:04:05"),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total":    len(sessions),
		"sessions": sessions,
	})
}

// createSessionHandler - Oturum oluştur
func (as *APIService) createSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string            `json:"session_id"`
		ProxyURL  string            `json:"proxy_url"`
		Headers   map[string]string `json:"headers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}

	session := as.sessionMgr.CreateSession(req.SessionID, req.ProxyURL, req.Headers, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"session_id": session.ID,
		"proxy_url":  session.ProxyURL,
		"created_at": session.LastActivity.Format("2006-01-02 15:04:05"),
	})
}

// deleteSessionHandler - Oturumu sil
func (as *APIService) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	as.sessionMgr.DeleteSession(sessionID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted":    true,
		"session_id": sessionID,
	})
}

// checkSlotsHandler - Slot kontrolü yap
func (as *APIService) checkSlotsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionID string `json:"session_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Slot kontrol yapacağı işi worker service'e gönder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "checking",
		"session_id": req.SessionID,
		"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	})
}

// latestSlotHandler - Son slot kontrolü sonucunu al
func (as *APIService) latestSlotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	as.mu.RLock()
	defer as.mu.RUnlock()

	if as.lastSlotCheck == nil {
		http.Error(w, "No slot check data available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(as.lastSlotCheck)
}

// proxyStatusHandler - Proxy durumu
func (as *APIService) proxyStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "proxy-service-status",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}

// statsHandler - İstatistikler
func (as *APIService) statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_sessions": len(as.sessionMgr.sessions),
		"uptime":         time.Now().Format("2006-01-02 15:04:05"),
	})
}

// Start - API servisi başlat
func (as *APIService) Start() error {
	log.Printf("API Service başlatılıyor: :%s", as.port)
	return http.ListenAndServe(":"+as.port, as.router)
}

func init() {
	_ = godotenv.Load()
}
