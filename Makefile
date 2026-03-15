DB_URL=postgres://autoedgeuser:autoedgetest1234@localhost:5432/autoedge_db?sslmode=disable

migrate-up:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down

.PHONY: migrate-up migrate-down