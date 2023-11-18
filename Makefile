db_up:
	docker compose up -d postgres

db_down:
	docker compose down -v

migrate_up:
	migrate -path db/migration -database "postgresql://postgres:jaw123@localhost:8892/curso?sslmode=disable" -verbose up
migrate_down:
	migrate -path db/migration -database "postgresql://postgres:jaw123@localhost:8892/curso?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

.PHONY: db_up db_down migrate_up migrate_down test