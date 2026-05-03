module furab-backend/services/user-service

go 1.23.0

require (
	furab-backend/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.12
	go.uber.org/mock v0.6.0
)

replace furab-backend/shared => ../../shared
