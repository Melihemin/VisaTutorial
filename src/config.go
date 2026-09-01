package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config - Uygulamanın konfigürasyon yapısı
type Config struct {
	// VFS Global API
	APIEndpoint       string
	SessionCookie     string
	AuthorizeToken    string
	ClientSourceToken string

	// HTTP Başlıkları
	UserAgent string
	SecChUa   string

	// Proxy Ayarları
	ProxyList       []string
	UseProxy        bool
	StickyProxyMode bool

	// Kontrol Ayarları
	CheckInterval  int // Saniye
	RequestTimeout int // Saniye
	MaxRetries     int
	LogToFile      bool

	// Bildirim Ayarları
	NotifyOnSlotFound bool
	TelegramToken     string
	TelegramChatID    string
	EmailEnabled      bool
	EmailAddress      string

	// Payload
	Country          string
	Mission          string
	VacancyCode      string
	VisaCategoryCode string
	LoginUser        string
	PayCode          string
}

// LoadConfig - .env dosyasından konfigürasyon yükle
func LoadConfig() *Config {
	// .env dosyasını yükle (başarısız olsa da devam et)
	_ = godotenv.Load()

	cfg := &Config{
		// Varsayılan değerler
		APIEndpoint:       "https://lift-api.vfsglobal.com/appointment/CheckIsSlotAvailable",
		CheckInterval:     45,
		RequestTimeout:    30,
		MaxRetries:        3,
		StickyProxyMode:   true,
		NotifyOnSlotFound: true,
		Country:           "usa",
		Mission:           "che",
		VacancyCode:       "NYC",
		VisaCategoryCode:  "SUSVFF",
		LoginUser:         "meliheminbusiness@gmail.com",
		PayCode:           "",

		// Varsayılan başlıklar
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36",
		SecChUa:   `"Not=A?Brand";v="99", "Google Chrome";v="151", "Chromium";v="151"`,
	}

	// .env değerlerini oku
	cfg.SessionCookie = getEnv("SESSION_COOKIE", "")
	cfg.AuthorizeToken = getEnv("AUTHORIZE_TOKEN", "")
	cfg.ClientSourceToken = getEnv("CLIENT_SOURCE_TOKEN", "")
	cfg.UserAgent = getEnv("USER_AGENT", cfg.UserAgent)
	cfg.SecChUa = getEnv("SEC_CH_UA", cfg.SecChUa)

	// Proxy ayarları
	proxyStr := getEnv("PROXY_LIST", "")
	if proxyStr != "" {
		cfg.ProxyList = strings.Split(proxyStr, ",")
		cfg.UseProxy = true
	}

	// Bildirim ayarları
	cfg.TelegramToken = getEnv("TELEGRAM_BOT_TOKEN", "")
	cfg.TelegramChatID = getEnv("TELEGRAM_CHAT_ID", "")
	if cfg.TelegramToken != "" && cfg.TelegramChatID != "" {
		cfg.NotifyOnSlotFound = true
	}

	cfg.EmailAddress = getEnv("NOTIFICATION_EMAIL", "")
	cfg.EmailEnabled = cfg.EmailAddress != ""

	// Sayısal ayarlar
	if interval := getEnv("CHECK_INTERVAL", ""); interval != "" {
		if val, err := strconv.Atoi(interval); err == nil {
			cfg.CheckInterval = val
		}
	}

	if timeout := getEnv("REQUEST_TIMEOUT", ""); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			cfg.RequestTimeout = val
		}
	}

	if retries := getEnv("MAX_RETRIES", ""); retries != "" {
		if val, err := strconv.Atoi(retries); err == nil {
			cfg.MaxRetries = val
		}
	}

	cfg.LogToFile = getEnvBool("LOG_TO_FILE", true)

	// Payload özelleştirmeleri
	cfg.Country = getEnv("COUNTRY", cfg.Country)
	cfg.Mission = getEnv("MISSION", cfg.Mission)
	cfg.VacancyCode = getEnv("VACANCY_CODE", cfg.VacancyCode)
	cfg.VisaCategoryCode = getEnv("VISA_CATEGORY_CODE", cfg.VisaCategoryCode)
	cfg.LoginUser = getEnv("LOGIN_USER", cfg.LoginUser)

	return cfg
}

// Validate - Konfigürasyonu doğrula
func (c *Config) Validate() error {
	if c.SessionCookie == "" {
		log.Fatal("ERROR: SESSION_COOKIE .env dosyasında tanımlanmalı")
	}
	if c.AuthorizeToken == "" {
		log.Fatal("ERROR: AUTHORIZE_TOKEN .env dosyasında tanımlanmalı")
	}
	if c.ClientSourceToken == "" {
		log.Fatal("ERROR: CLIENT_SOURCE_TOKEN .env dosyasında tanımlanmalı")
	}
	return nil
}

// Helper fonksiyonlar
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val == "true" || val == "1" || val == "yes"
}
