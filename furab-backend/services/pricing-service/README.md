# Pricing Service -join ' ')

Dynamic pricing and surge pricing service

## Deskripsi

TODO: Tambahkan deskripsi lengkap service ini.

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Language | Go 1.22+ |
| HTTP Router | chi |
| Database | PostgreSQL |
| Testing | gomock, go test |

## Struktur Folder

```
pricing-service/
+-- cmd/main.go
+-- internal/
    +-- handler/price_handler.go
    +-- service/price_service.go
    +-- repository/price_repository.go
    +-- model/price.go
+-- test/
    +-- unit/price_service_test.go
    +-- unit/mock/
    +-- functional/price_functional_test.go
+-- go.mod
+-- Dockerfile
+-- README.md
```

## Cara Menjalankan

```bash
# Set environment variables
export SERVER_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=furab
export DB_PASSWORD=furab_secret
export DB_NAME=pricing-service

# Jalankan service
go run cmd/main.go
```

## Menjalankan Tests

### Unit Tests (Tanpa Database)
```bash
go test ./test/unit/... -v
```

**Test BERHASIL jika output:**
```
--- PASS: TestNewPriceService_Creation
PASS
```

**Test GAGAL jika output:**
```
--- FAIL: TestNewPriceService_Creation
FAIL
```

### Functional Tests (Dengan Database)
```bash
# Pastikan PostgreSQL berjalan
go test ./test/functional/... -v -tags=functional
```

## Docker

```bash
# Build (dari root project)
docker build -t furab/pricing-service:latest -f services/pricing-service/Dockerfile .

# Run
docker run -p 8080:8080 furab/pricing-service:latest
```
