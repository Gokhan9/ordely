# Orderly - High Performance E-Commerce API

Orderly, modern yazılım mimarisi prensipleri (Clean Architecture) ile geliştirilmiş, yüksek performanslı ve ölçeklenebilir bir e-ticaret backend uygulamasıdır. **Golang**'in gücüyle, düşük gecikme süresi ve yüksek eşzamanlılık hedeflenerek tasarlanmıştır.

## 🚀 Teknolojiler

- **Backend:** Go (Golang)
- **Web Framework:** [Fiber v2](https://gofiber.io/) (Express-like, zero memory allocation)
- **Database:** PostgreSQL 15 (via [pgx/v5](https://github.com/jackc/pgx))
- **SQL Code Generation:** [sqlc](https://sqlc.dev/) (Type-safe SQL)
- **Authentication:** JWT (JSON Web Tokens)
- **Logging:** [Zerolog](https://github.com/rs/zerolog) (Structured JSON logging)
- **Validation:** [Go-Playground Validator](https://github.com/go-playground/validator)
- **Infrastructure:** Docker & Docker Compose

## 🏗️ Mimari Yapı (Clean Architecture)

Proje, bağımlılıkların içe doğru olduğu katmanlı bir yapıda kurulmuştur:

- `cmd/api/`: Uygulamanın giriş noktası ve dependency injection.
- `internal/domain/`: İş modelleri, interface'ler ve global kurallar.
- `internal/usecase/`: Business logic (İş mantığı) katmanı.
- `internal/repository/`: Veritabanı ve dış kaynak erişim katmanı.
- `internal/delivery/http/`: HTTP handler'lar ve API uç noktaları.
- `pkg/`: Paylaşılan yardımcı kütüphaneler (Logger, Config, Utils).

## 🛠️ Temel Özellikler

- [x] **User Management:** Kayıt, Giriş ve JWT tabanlı yetkilendirme.
- [x] **Product Catalog:** Paging (sayfalama) ve kategori bazlı ürün yönetimi.
- [x] **Order System:** 
  - Atomik (Transactional) sipariş oluşturma.
  - Otomatik stok kontrolü ve güncelleme.
  - Sipariş toplam tutarı hesaplama.
- [x] **Graceful Shutdown:** Kesintisiz sistem durdurma.
- [x] **Security:** Bcrypt şifre hashleme ve JWT middleware.

## 🏁 Başlangıç

### Gereksinimler
- Docker & Docker Compose
- Go 1.21+

### Kurulum
1. Repoyu klonlayın.
2. Servisleri başlatın:
   ```bash
   docker-compose up -d
   ```
3. Uygulamayı çalıştırın:
   ```bash
   go run cmd/api/main.go
   ```

## 📝 Gelecek Planları
1. Unit & Integration Testleri (Mocking ile).
2. Swagger/OpenAPI dökümantasyonu.
3. Next.js 14+ Frontend entegrasyonu.
