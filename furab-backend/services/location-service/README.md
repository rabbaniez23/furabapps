# Location Service -join ' ')

Real-time location tracking service

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
location-service/
+-- cmd/main.go
+-- internal/
    +-- handler/location_handler.go
    +-- service/location_service.go
    +-- repository/location_repository.go
    +-- model/location.go
+-- test/
    +-- unit/location_service_test.go
    +-- unit/mock/
    +-- functional/location_functional_test.go
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
export DB_NAME=location-service

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
--- PASS: TestNewLocationService_Creation
PASS
```

**Test GAGAL jika output:**
```
--- FAIL: TestNewLocationService_Creation
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
docker build -t furab/location-service:latest -f services/location-service/Dockerfile .

# Run
docker run -p 8080:8080 furab/location-service:latest
```
