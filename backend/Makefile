include .env
.PHONY: test gooseUp gooseDown gooseStatus
MIGRATIONS_DIR=./sql/schema

test: 
	gotestsum --format testname -- -count=1 ./...

gooseUp:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

gooseDown:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

gooseStatus:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

gooseReset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up
	go run seed/main.go