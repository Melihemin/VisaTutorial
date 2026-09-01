package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger - Özel logger yapısı
type Logger struct {
	LogToFile bool
	LogFile   *os.File
}

// NewLogger - Yeni logger oluştur
func NewLogger(logToFile bool) *Logger {
	logger := &Logger{
		LogToFile: logToFile,
	}

	if logToFile {
		filename := fmt.Sprintf("randevu_tracker_%s.log", time.Now().Format("2006-01-02"))
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Printf("Log dosyası açma hatası: %v", err)
			return logger
		}
		logger.LogFile = file
	}

	return logger
}

// Info - Bilgi mesajı
func (l *Logger) Info(format string, v ...interface{}) {
	msg := fmt.Sprintf("[INFO] %s - %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, v...))
	fmt.Println(msg)

	if l.LogFile != nil {
		l.LogFile.WriteString(msg + "\n")
	}
}

// Error - Hata mesajı
func (l *Logger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[ERROR] %s - %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, v...))
	fmt.Println(msg)

	if l.LogFile != nil {
		l.LogFile.WriteString(msg + "\n")
	}
}

// Success - Başarı mesajı
func (l *Logger) Success(format string, v ...interface{}) {
	msg := fmt.Sprintf("[✓ SUCCESS] %s - %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, v...))
	fmt.Println(msg)

	if l.LogFile != nil {
		l.LogFile.WriteString(msg + "\n")
	}
}

// Warning - Uyarı mesajı
func (l *Logger) Warning(format string, v ...interface{}) {
	msg := fmt.Sprintf("[WARNING] %s - %s", time.Now().Format("15:04:05"), fmt.Sprintf(format, v...))
	fmt.Println(msg)

	if l.LogFile != nil {
		l.LogFile.WriteString(msg + "\n")
	}
}

// Close - Logger'ı kapat
func (l *Logger) Close() {
	if l.LogFile != nil {
		l.LogFile.Close()
	}
}

// Stats - İstatistik tutma
type Stats struct {
	TotalChecks      int
	SuccessfulChecks int
	FailedChecks     int
	SlotsFound       int
	StartTime        time.Time
	LastCheckTime    time.Time
}

// UpdateStats - İstatistikleri güncelle
func (s *Stats) UpdateStats(success bool, slotsFound bool) {
	s.TotalChecks++
	s.LastCheckTime = time.Now()

	if success {
		s.SuccessfulChecks++
	} else {
		s.FailedChecks++
	}

	if slotsFound {
		s.SlotsFound++
	}
}

// GetSummary - İstatistik özeti al
func (s *Stats) GetSummary() map[string]interface{} {
	uptime := time.Since(s.StartTime)
	successRate := 0.0
	if s.TotalChecks > 0 {
		successRate = float64(s.SuccessfulChecks) / float64(s.TotalChecks) * 100
	}

	return map[string]interface{}{
		"total_checks":      s.TotalChecks,
		"successful_checks": s.SuccessfulChecks,
		"failed_checks":     s.FailedChecks,
		"success_rate":      fmt.Sprintf("%.2f%%", successRate),
		"slots_found":       s.SlotsFound,
		"uptime":            uptime.String(),
		"last_check":        s.LastCheckTime.Format("2006-01-02 15:04:05"),
	}
}
