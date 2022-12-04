.SILENT:
-include .env

DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

run:
	go run cmd/main.go

print:
	echo $(DB_URL)

swag-init:
	swag init -g api/server.go -o api/docs

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-up1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migrate-down1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1