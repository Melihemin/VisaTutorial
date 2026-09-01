package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// SessionService - Oturum yönetimi servisi
type SessionService struct {
	port       string
	sessionMgr *SessionManager
	router     *http.ServeMux
	mu         sync.RWMutex
}

// NewSessionService - Yeni session servisi oluştur
func NewSessionService(port string) *SessionService {
	svc := &SessionService{
		port:       port,
		sessionMgr: NewSessionManager(),
		router:     http.NewServeMux(),
	}

	svc.setupRoutes()
	return svc
}

// setupRoutes - Rotaları konfigüre et
func (ss *SessionService) setupRoutes() {
	ss.router.HandleFunc("/health", ss.healthHandler)
	ss.router.HandleFunc("/sessions", ss.sessionsHandler)
	ss.router.HandleFunc("/sessions/validate", ss.validateHandler)
	ss.router.HandleFunc("/sessions/renew", ss.renewHandler)
}

// healthHandler - Sağlık kontrolü
func (ss *SessionService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","service":"session-service"}`)
}

// sessionsHandler - Oturum listele/yönet
func (ss *SessionService) sessionsHandler(w http.ResponseWriter, r *http.Request) {
	ss.sessionMgr.mu.RLock()
	sessionCount := len(ss.sessionMgr.sessions)
	ss.sessionMgr.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"total_sessions":%d,"timestamp":"%s"}`,
		sessionCount, time.Now().Format("2006-01-02 15:04:05"))
}

// validateHandler - Oturum doğrulama
func (ss *SessionService) validateHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	_, exists := ss.sessionMgr.GetSession(sessionID)

	w.Header().Set("Content-Type", "application/json")
	if exists {
		fmt.Fprintf(w, `{"valid":true,"session_id":"%s"}`, sessionID)
	} else {
		fmt.Fprintf(w, `{"valid":false,"session_id":"%s"}`, sessionID)
	}
}

// renewHandler - Oturum yenile
func (ss *SessionService) renewHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("id")
	if sessionID == "" {
		http.Error(w, "id parameter required", http.StatusBadRequest)
		return
	}

	err := ss.sessionMgr.UpdateSessionActivity(sessionID)

	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		fmt.Fprintf(w, `{"renewed":true,"session_id":"%s"}`, sessionID)
	} else {
		fmt.Fprintf(w, `{"renewed":false,"error":"%s"}`, err.Error())
	}
}

// Start - Servisi başlat
func (ss *SessionService) Start() error {
	log.Printf("Session Service başlatılıyor: :%s", ss.port)
	return http.ListenAndServe(":"+ss.port, ss.router)
}

func init() {
	_ = godotenv.Load()
}
