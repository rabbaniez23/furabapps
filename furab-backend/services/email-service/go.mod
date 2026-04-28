module furab-backend/services/email-service

go 1.22

require (
	furab-backend/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.12
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.6.0
	go.uber.org/mock v0.4.0
)

replace furab-backend/shared => ../../shared
