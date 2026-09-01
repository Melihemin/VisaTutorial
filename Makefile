.PHONY: help build run docker-build docker-up docker-down clean test

help:
	@echo "Randevu Uygunluk Takip Sistemi - Build Hedefleri"
	@echo "================================================"
	@echo ""
	@echo "Go Geliştirme:"
	@echo "  make build              - Go uygulamasını derle"
	@echo "  make run                - Uygulamayı çalıştır"
	@echo "  make clean              - Build dosyalarını temizle"
	@echo "  make test               - Testleri çalıştır"
	@echo "  make fmt                - Kodları formatla"
	@echo "  make lint               - Linter çalıştır"
	@echo ""
	@echo "Docker Komutları:"
	@echo "  make docker-build       - Docker image'ları derle"
	@echo "  make docker-up          - Kontainer'ları başlat"
	@echo "  make docker-down        - Kontainer'ları durdur"
	@echo "  make docker-logs        - Kontainer loglarını göster"
	@echo "  make docker-clean       - Kontainer'ları temizle"
	@echo ""
	@echo "Diğer:"
	@echo "  make env                - .env.example'dan .env oluştur"
	@echo "  make help               - Bu yardımı göster"

# Build targets
build:
	@echo "Go uygulaması derleniyorum (src/ klasöründen)..."
	cd src && go build -o ../build/randevu_tracker.exe *.go
	@echo "✓ Derleme tamamlandı: build/randevu_tracker.exe"

run: build
	@echo "Uygulama çalıştırılıyor..."
	./build/randevu_tracker.exe

clean:
	@echo "Temizleme yapılıyor..."
	rm -f build/randevu_tracker.exe
	rm -f *.log
	cd src && go clean -cache -testcache
	@echo "✓ Temizleme tamamlandı"

test:
	@echo "Testler çalıştırılıyor..."
	cd src && go test -v -cover ./...

fmt:
	@echo "Kod formatlanıyor..."
	cd src && go fmt ./...
	@echo "✓ Formatla tamamlandı"

lint:
	@echo "Linter çalıştırılıyor..."
	cd src && go vet ./...
	@echo "✓ Lint tamamlandı"

# Docker targets
docker-build:
	@echo "Docker image'ları derleniyorum..."
	docker build -f docker/Dockerfile.api -t randevu/api:latest .
	docker build -f docker/Dockerfile.worker -t randevu/worker:latest .
	docker build -f docker/Dockerfile.session -t randevu/session:latest .
	docker build -f docker/Dockerfile.proxy -t randevu/proxy:latest .
	@echo "✓ Docker image'ları derlenmiş"

docker-up:
	@echo "Kontainer'lar başlatılıyor..."
	docker-compose -f docker/docker-compose.yml up -d
	@echo "✓ Kontainer'lar başlatıldı"
	@echo ""
	@echo "Servisler:"
	@echo "  API:      http://localhost:8080"
	@echo "  Session:  http://localhost:8081"
	@echo "  Proxy:    http://localhost:8082"

docker-down:
	@echo "Kontainer'lar durduruluyor..."
	docker-compose -f docker/docker-compose.yml down
	@echo "✓ Kontainer'lar durduruldu"

docker-logs:
	docker-compose -f docker/docker-compose.yml logs -f

docker-clean:
	@echo "Docker kaynakları temizleniyor..."
	docker-compose -f docker/docker-compose.yml down -v
	docker system prune -f
	@echo "✓ Docker temizlemesi tamamlandı"

# Environment setup
env:
	@if [ ! -f config/.env ]; then \
		cp config/.env.example config/.env; \
		echo "✓ .env dosyası oluşturuldu (config/ klasöründe)"; \
		echo "⚠️  config/.env dosyasını düzenleyin ve kimlik bilgilerini girin"; \
	else \
		echo "✓ config/.env dosyası zaten var"; \
	fi

# Help
help:
	@echo "Randevu Uygunluk Takip Sistemi - Build Hedefleri"
	@echo "================================================"
	@echo ""
	@echo "Go Geliştirme:"
	@echo "  make build              - Go uygulamasını derle"
	@echo "  make run                - Uygulamayı çalıştır"
	@echo "  make clean              - Build dosyalarını temizle"
	@echo "  make test               - Testleri çalıştır"
	@echo "  make fmt                - Kodları formatla"
	@echo "  make lint               - Linter çalıştır"
	@echo ""
	@echo "Docker Komutları:"
	@echo "  make docker-build       - Docker image'ları derle"
	@echo "  make docker-up          - Kontainer'ları başlat"
	@echo "  make docker-down        - Kontainer'ları durdur"
	@echo "  make docker-logs        - Kontainer loglarını göster"
	@echo "  make docker-clean       - Kontainer'ları temizle"
	@echo ""
	@echo "Diğer:"
	@echo "  make env                - .env.example'dan .env oluştur"
	@echo "  make help               - Bu yardımı göster"
