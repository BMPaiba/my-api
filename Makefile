include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migration
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: docker-build
docker-build:
	docker build -t my-api-v1 .

.PHONY: docker-run
docker-run:
	docker run --rm -p 8080:8080 --env-file .env my-api-v1

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt
