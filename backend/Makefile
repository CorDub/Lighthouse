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