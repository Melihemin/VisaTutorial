package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/joho/godotenv"
)

// Her bir başvuru grubunun slot tarihini tutacak alt yapı
type SlotItem struct {
	Applicant string `json:"applicant"`
	Date      string `json:"date"`
}

// API'den dönen ana JSON yapısı
type SlotResponse struct {
	EarliestDate      string      `json:"earliestDate"`
	EarliestSlotLists []SlotItem  `json:"earliestSlotLists"`
	Error             interface{} `json:"error"` // Hata yoksa null döner
}

// Session yönetimi için yapı
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

// Oturum bilgilerini tutan yapı
type Session struct {
	ID           string
	Cookies      []*http.Cookie
	Headers      map[string]string
	ProxyURL     string
	LastActivity time.Time
	Client       tls_client.HttpClient
}

// Slot takip sonuçları
type SlotCheckResult struct {
	Timestamp    time.Time
	EarliestDate string
	Slots        []SlotItem
	Error        string
}

func readEnvFileLenient(path string) map[string]string {
	values := make(map[string]string)

	data, err := os.ReadFile(path)
	if err != nil {
		return values
	}

	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(strings.TrimPrefix(parts[0], "\uFEFF"))
		val := strings.TrimSpace(parts[1])

		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		values[key] = val
	}

	return values
}

func main() {
	// 1. TLS İstemcisi (Şifreleme profili 120 olarak kalır, güvenlidir)
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_146),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Fatalf("TLS İstemcisi oluşturulamadı: %v", err)
	}

	// 2. VFS Global API Uç Noktası
	apiEndpoint := "https://lift-api.vfsglobal.com/appointment/CheckIsSlotAvailable"

	dotenvValues, dotenvErr := godotenv.Read(".env")
	if dotenvErr != nil {
		log.Printf(".env parse uyarisi: %v. Esnek okuyucuya geciliyor.", dotenvErr)
		dotenvValues = readEnvFileLenient(".env")
	}

	getCfg := func(key string) string {
		if v := strings.TrimSpace(os.Getenv(key)); v != "" {
			return v
		}
		return strings.TrimSpace(dotenvValues[key])
	}

	sessionCookie := getCfg("SESSION_COOKIE")
	authorizeToken := getCfg("AUTHORIZE_TOKEN")
	clientSourceToken := getCfg("CLIENT_SOURCE_TOKEN")
	userAgent := getCfg("USER_AGENT")
	secChUa := getCfg("SEC_CH_UA")

	if sessionCookie == "" || authorizeToken == "" || clientSourceToken == "" {
		log.Fatal("Eksik ortam değişkeni: SESSION_COOKIE, AUTHORIZE_TOKEN ve CLIENT_SOURCE_TOKEN .env içinde dolu olmalı")
	}

	if userAgent == "" {
		userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"
	}

	if secChUa == "" {
		secChUa = `"Not=A?Brand";v="99", "Google Chrome";v="151", "Chromium";v="151"`
	}

	// 3. Randevu Sorgulama JSON Paketi (Payload)
	requestBody := []byte(`{"countryCode":"usa","missionCode":"che","vacCode":"NYC","visaCategoryCode":"SUSVFF","roleName":"Individual","loginUser":"meliheminbusiness@gmail.com","payCode":""}`)

	// 4. Ana Kontrol Döngüsü
	for {
		req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(requestBody))
		if err != nil {
			log.Println("İstek oluşturulamadı:", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		// Temel Başlıklar
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Cookie", sessionCookie)
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
		req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("Priority", "u=1, i")

		// VFS Global'e Özel Güvenlik ve Yönlendirme Başlıkları
		req.Header.Set("authorize", authorizeToken)
		req.Header.Set("clientsource", clientSourceToken)
		req.Header.Set("route", "usa/en/swe")

		req.Header.Set("Origin", "https://visa.vfsglobal.com")
		req.Header.Set("Referer", "https://visa.vfsglobal.com/")
		req.Header.Set("Sec-Ch-Ua", secChUa)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")

		fmt.Printf("[%s] VACHOUS merkezi için takvim sorgulanıyor...\n", time.Now().Format("15:04:05"))

		// İsteği gönder
		resp, err := client.Do(req)

		if err != nil {
			log.Println("İstek hatası (Ağ veya TLS problemi):", err)
		} else {
			if resp.StatusCode == 200 {
				bodyBytes, _ := io.ReadAll(resp.Body)

				var slotData SlotResponse
				// Gelen JSON'ı Go struct yapımıza dök
				if err := json.Unmarshal(bodyBytes, &slotData); err != nil {
					log.Println("JSON Parse Hatası:", err)
					fmt.Println("Gelen Yanıt (Bozuk/Farklı Format):", string(bodyBytes))
				} else {
					// Slot Var Mı Kontrolü
					if slotData.Error == nil && len(slotData.EarliestSlotLists) > 0 && slotData.EarliestDate != "" {
						fmt.Printf("\n🔥 BOŞ SLOT BULUNDU! En Erken Tarih: %s\n", slotData.EarliestDate)

						// Bulunan tüm slot alternatiflerini yazdır
						for _, slot := range slotData.EarliestSlotLists {
							fmt.Printf("  -> Başvuru Kombinasyonu: %s | Tarih: %s\n", slot.Applicant, slot.Date)
						}

						// TODO: Telegram mesaj fonksiyonunu buraya ekle

					} else {
						fmt.Println("Şu an boş slot yok. Bir sonraki kontrole kadar bekleniyor...")
					}
				}
			} else if resp.StatusCode == 401 || resp.StatusCode == 403 {
				bodyBytes, _ := io.ReadAll(resp.Body)
				fmt.Printf("🚨 UYARI: WAF Engeli (%d)! Oturum süresi doldu, parmak izi uyuşmuyor veya çerezler/tokenler bağlam dışı.\n", resp.StatusCode)
				if len(bodyBytes) > 0 {
					fmt.Println("WAF yanıtı:", string(bodyBytes))
				}
				resp.Body.Close()
				break // Döngüyü kır ve programı bitir
			} else {
				fmt.Printf("Beklenmeyen HTTP Durum Kodu: %d\n", resp.StatusCode)
				bodyBytes, _ := io.ReadAll(resp.Body)
				fmt.Println("Sunucu Yanıtı:", string(bodyBytes))
			}

			resp.Body.Close()
		}

		// WAF'ı şüphelendirmemek için 45 saniye bekle
		time.Sleep(45 * time.Second)
	}
}

// NewSessionManager - Yeni session yöneticisi oluştur
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

// CreateSession - Yeni oturum oluştur
func (sm *SessionManager) CreateSession(id, proxyURL string, headers map[string]string, client tls_client.HttpClient) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		ID:           id,
		ProxyURL:     proxyURL,
		Headers:      headers,
		LastActivity: time.Now(),
		Client:       client,
		Cookies:      make([]*http.Cookie, 0),
	}

	sm.sessions[id] = session
	return session
}

// GetSession - Oturumu al
func (sm *SessionManager) GetSession(id string) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[id]
	return session, exists
}

// UpdateSessionActivity - Oturum aktivitesini güncelle
func (sm *SessionManager) UpdateSessionActivity(id string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[id]; exists {
		session.LastActivity = time.Now()
		return nil
	}
	return fmt.Errorf("oturum bulunamadı: %s", id)
}

// DeleteSession - Oturumu sil
func (sm *SessionManager) DeleteSession(id string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.sessions, id)
}

// CheckSlots - Randevu slotlarını kontrol et
func CheckSlots(client tls_client.HttpClient, apiEndpoint string, requestBody []byte, headers map[string]string) (*SlotCheckResult, error) {
	result := &SlotCheckResult{
		Timestamp: time.Now(),
	}

	req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		result.Error = fmt.Sprintf("İstek oluşturma hatası: %v", err)
		return result, err
	}

	// Başlıkları ekle
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		result.Error = fmt.Sprintf("İstek gönderme hatası: %v", err)
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		result.Error = fmt.Sprintf("HTTP %d hatası", resp.StatusCode)
		return result, fmt.Errorf("HTTP hata kodu: %d", resp.StatusCode)
	}

	var slotData SlotResponse
	if err := json.Unmarshal(bodyBytes, &slotData); err != nil {
		result.Error = fmt.Sprintf("JSON parse hatası: %v", err)
		return result, err
	}

	result.EarliestDate = slotData.EarliestDate
	result.Slots = slotData.EarliestSlotLists

	return result, nil
}

// LogSlotCheck - Slot kontrolünü kaydet
func LogSlotCheck(result *SlotCheckResult) {
	filename := fmt.Sprintf("slot_checks_%s.log", time.Now().Format("2006-01-02"))

	logEntry := fmt.Sprintf("[%s] EarliestDate: %s | Slots: %d | Error: %s\n",
		result.Timestamp.Format("15:04:05"),
		result.EarliestDate,
		len(result.Slots),
		result.Error,
	)

	// Dosyaya ekle
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Log dosyası yazma hatası: %v", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Log yazma hatası: %v", err)
	}
}
