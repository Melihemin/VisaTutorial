package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Notifier - Bildirim göndericisi
type Notifier struct {
	TelegramToken  string
	TelegramChatID string
	EmailEnabled   bool
	EmailAddress   string
}

// NewNotifier - Yeni notifier oluştur
func NewNotifier(telegramToken, telegramChatID string, emailEnabled bool, emailAddress string) *Notifier {
	return &Notifier{
		TelegramToken:  telegramToken,
		TelegramChatID: telegramChatID,
		EmailEnabled:   emailEnabled,
		EmailAddress:   emailAddress,
	}
}

// SendTelegramNotification - Telegram bildirim gönder
func (n *Notifier) SendTelegramNotification(title, message string) error {
	if n.TelegramToken == "" || n.TelegramChatID == "" {
		return fmt.Errorf("Telegram ayarları eksik")
	}

	// Telegram Bot API URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.TelegramToken)

	payload := map[string]interface{}{
		"chat_id":    n.TelegramChatID,
		"text":       fmt.Sprintf("🔔 *%s*\n\n%s\n\nZaman: %s", title, message, time.Now().Format("2006-01-02 15:04:05")),
		"parse_mode": "Markdown",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram API hatası (%d): %s", resp.StatusCode, string(body))
	}

	log.Printf("Telegram bildirimi gönderildi: %s", title)
	return nil
}

// SendSlotFoundNotification - Slot bulundu bildirimi gönder
func (n *Notifier) SendSlotFoundNotification(slotResult *SlotCheckResult) error {
	if slotResult == nil || len(slotResult.Slots) == 0 {
		return fmt.Errorf("slot verileri boş")
	}

	message := fmt.Sprintf(
		"🎉 *BOŞ SLOT BULUNDU!*\n\n"+
			"En Erken Tarih: *%s*\n\n"+
			"Slot Detayları:\n",
		slotResult.EarliestDate,
	)

	for i, slot := range slotResult.Slots {
		message += fmt.Sprintf("%d. %s - %s\n", i+1, slot.Applicant, slot.Date)
	}

	return n.SendTelegramNotification("SLOT BULUNDU", message)
}

// SendErrorNotification - Hata bildirimi gönder
func (n *Notifier) SendErrorNotification(errorMsg string) error {
	message := fmt.Sprintf("❌ Hata oluştu:\n\n%s", errorMsg)
	return n.SendTelegramNotification("HATA", message)
}

// SendStatusNotification - Durum bildirimi gönder
func (n *Notifier) SendStatusNotification(status string, details map[string]interface{}) error {
	message := fmt.Sprintf("📊 Durum: %s\n\n", status)

	for key, value := range details {
		message += fmt.Sprintf("• %s: %v\n", key, value)
	}

	return n.SendTelegramNotification("DURUM RAPORU", message)
}

// TestConnection - Bildirim bağlantısını test et
func (n *Notifier) TestConnection() error {
	return n.SendTelegramNotification("TEST", "Bildirim sistemi çalışıyor!")
}
