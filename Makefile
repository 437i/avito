# Makefile
DB_URL="postgres://myuser:mypassword@localhost:5432/test_db?sslmode=disable"

migrate-up:
	goose -dir ./migrations postgres $(DB_URL) up

migrate-down:
	goose -dir ./migrations postgres $(DB_URL) down