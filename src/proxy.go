package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// ProxyPool - Proxy havuzu yönetimi
type ProxyPool struct {
	mu         sync.RWMutex
	proxies    []ProxyEntry
	currentIdx int
}

// ProxyEntry - Proxy girdisi
type ProxyEntry struct {
	URL       string
	IsActive  bool
	FailCount int
	LastUsed  time.Time
	Sticky    bool
	SessionID string
}

// ProxyWorker - Proxy tabanlı isçi
type ProxyWorker struct {
	ID          string
	ProxyPool   *ProxyPool
	SessionMgr  *SessionManager
	Client      tls_client.HttpClient
	WorkerQueue chan WorkItem
	Done        chan bool
	mu          sync.RWMutex
}

// WorkItem - İş öğesi
type WorkItem struct {
	ID        string
	Endpoint  string
	Payload   []byte
	Headers   map[string]string
	SessionID string
	Retries   int
	MaxRetry  int
}

// WorkResult - İş sonucu
type WorkResult struct {
	WorkID     string
	Success    bool
	Data       interface{}
	Error      string
	ProxyUsed  string
	Timestamp  time.Time
	StatusCode int
}

// NewProxyPool - Yeni proxy havuzu oluştur
func NewProxyPool(proxyList []string) *ProxyPool {
	pool := &ProxyPool{
		proxies:    make([]ProxyEntry, 0),
		currentIdx: 0,
	}

	for _, proxyStr := range proxyList {
		// Proxy URL doğrulaması
		if _, err := url.Parse(proxyStr); err == nil {
			pool.proxies = append(pool.proxies, ProxyEntry{
				URL:      proxyStr,
				IsActive: true,
				Sticky:   true,
			})
		}
	}

	return pool
}

// GetNextProxy - Sonraki aktif proxy'i al
func (pp *ProxyPool) GetNextProxy() (ProxyEntry, error) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	if len(pp.proxies) == 0 {
		return ProxyEntry{}, fmt.Errorf("kullanılabilir proxy yok")
	}

	// Aktif proxy ara
	for i := 0; i < len(pp.proxies); i++ {
		idx := (pp.currentIdx + i) % len(pp.proxies)
		if pp.proxies[idx].IsActive && pp.proxies[idx].FailCount < 5 {
			pp.currentIdx = (idx + 1) % len(pp.proxies)
			return pp.proxies[idx], nil
		}
	}

	return ProxyEntry{}, fmt.Errorf("aktif proxy bulunamadı")
}

// GetStickyProxy - Sticky oturum proxy'i al
func (pp *ProxyPool) GetStickyProxy(sessionID string) (ProxyEntry, error) {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	for _, proxy := range pp.proxies {
		if proxy.SessionID == sessionID && proxy.IsActive {
			return proxy, nil
		}
	}

	return ProxyEntry{}, fmt.Errorf("sticky proxy bulunamadı: %s", sessionID)
}

// SetStickyProxy - Oturum için sticky proxy belirle
func (pp *ProxyPool) SetStickyProxy(sessionID string, proxyURL string) error {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	for i := range pp.proxies {
		if pp.proxies[i].URL == proxyURL {
			pp.proxies[i].SessionID = sessionID
			pp.proxies[i].LastUsed = time.Now()
			return nil
		}
	}

	return fmt.Errorf("proxy bulunamadı: %s", proxyURL)
}

// MarkProxyFailed - Proxy başarısızlığını işaretle
func (pp *ProxyPool) MarkProxyFailed(proxyURL string) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	for i := range pp.proxies {
		if pp.proxies[i].URL == proxyURL {
			pp.proxies[i].FailCount++
			if pp.proxies[i].FailCount >= 5 {
				pp.proxies[i].IsActive = false
				log.Printf("Proxy devre dışı bırakıldı (5 başarısızlık): %s", proxyURL)
			}
			return
		}
	}
}

// ResetProxyFailCount - Proxy başarısızlık sayısını sıfırla
func (pp *ProxyPool) ResetProxyFailCount(proxyURL string) {
	pp.mu.Lock()
	defer pp.mu.Unlock()

	for i := range pp.proxies {
		if pp.proxies[i].URL == proxyURL {
			pp.proxies[i].FailCount = 0
			pp.proxies[i].IsActive = true
			pp.proxies[i].LastUsed = time.Now()
			return
		}
	}
}

// NewProxyWorker - Yeni proxy işçi oluştur
func NewProxyWorker(id string, proxyPool *ProxyPool, sessionMgr *SessionManager) (*ProxyWorker, error) {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_146),
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	worker := &ProxyWorker{
		ID:          id,
		ProxyPool:   proxyPool,
		SessionMgr:  sessionMgr,
		Client:      client,
		WorkerQueue: make(chan WorkItem, 100),
		Done:        make(chan bool),
	}

	return worker, nil
}

// Start - İşçi başlat
func (pw *ProxyWorker) Start() {
	go func() {
		for {
			select {
			case <-pw.Done:
				log.Printf("Worker %s durduruluyor...", pw.ID)
				return
			case workItem := <-pw.WorkerQueue:
				result := pw.ProcessWork(workItem)
				log.Printf("Work %s sonuç: %+v", result.WorkID, result)
			}
		}
	}()
}

// Stop - İşçi durdur
func (pw *ProxyWorker) Stop() {
	pw.Done <- true
}

// ProcessWork - İşi işle
func (pw *ProxyWorker) ProcessWork(work WorkItem) *WorkResult {
	result := &WorkResult{
		WorkID:    work.ID,
		Timestamp: time.Now(),
	}

	// Session yöneticisinden session al
	var session *Session
	if work.SessionID != "" {
		var exists bool
		session, exists = pw.SessionMgr.GetSession(work.SessionID)
		if !exists {
			result.Error = fmt.Sprintf("Oturum bulunamadı: %s", work.SessionID)
			result.Success = false
			return result
		}
	}

	// Proxy seç
	var proxy ProxyEntry
	var err error

	if session != nil && session.ProxyURL != "" {
		// Sticky proxy kullan
		proxy, err = pw.ProxyPool.GetStickyProxy(work.SessionID)
		if err != nil {
			log.Printf("Sticky proxy başarısız, yeni proxy seçiliyor: %v", err)
			proxy, err = pw.ProxyPool.GetNextProxy()
		}
	} else {
		proxy, err = pw.ProxyPool.GetNextProxy()
	}

	if err != nil {
		result.Error = fmt.Sprintf("Proxy seçme hatası: %v", err)
		result.Success = false
		return result
	}

	result.ProxyUsed = proxy.URL

	// İsteği gönder
	req, err := http.NewRequest(http.MethodPost, work.Endpoint, nil)
	if err != nil {
		result.Error = fmt.Sprintf("İstek oluşturma hatası: %v", err)
		result.Success = false
		return result
	}

	// Başlıkları ayarla
	for key, value := range work.Headers {
		req.Header.Set(key, value)
	}

	// Proxy URL'ini ayarla
	pw.Client.SetProxy(proxy.URL)

	// İsteği gönder
	resp, err := pw.Client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("İstek gönderme hatası: %v", err)
		result.Success = false
		pw.ProxyPool.MarkProxyFailed(proxy.URL)

		// Retry mantığı
		if work.Retries < work.MaxRetry {
			work.Retries++
			go func() {
				time.Sleep(time.Duration(work.Retries*5) * time.Second)
				pw.WorkerQueue <- work
			}()
		}
		return result
	}

	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = resp.StatusCode == 200

	if result.Success {
		pw.ProxyPool.ResetProxyFailCount(proxy.URL)

		// Sticky oturum ise proxy'i sakla
		if session != nil {
			pw.ProxyPool.SetStickyProxy(work.SessionID, proxy.URL)
		}
	}

	return result
}

// SubmitWork - İş gönder
func (pw *ProxyWorker) SubmitWork(work WorkItem) error {
	select {
	case pw.WorkerQueue <- work:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("iş gönderme zaman aşımı")
	}
}

// GetPoolStats - Havuz istatistiklerini al
func (pp *ProxyPool) GetPoolStats() map[string]interface{} {
	pp.mu.RLock()
	defer pp.mu.RUnlock()

	active := 0
	inactive := 0

	for _, proxy := range pp.proxies {
		if proxy.IsActive {
			active++
		} else {
			inactive++
		}
	}

	return map[string]interface{}{
		"total_proxies": len(pp.proxies),
		"active":        active,
		"inactive":      inactive,
	}
}
