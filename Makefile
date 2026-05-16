# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@if command -v column >/dev/null 2>&1; then \
		sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'; \
	else \
		sed -n 's/^##//p' ${MAKEFILE_LIST} | sed -e 's/^/ /'; \
	fi

# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker/up: build images and start all services (detached)
.PHONY: docker/up
docker/up:
	docker compose up --build -d

## docker/up/attached: build images and start all services (follow logs)
.PHONY: docker/up/attached
docker/up/attached:
	docker compose up --build

## docker/down: stop all services
.PHONY: docker/down
docker/down:
	docker compose down

## docker/down/volumes: stop all services AND delete volumes (wipes DB data)
.PHONY: docker/down/volumes
docker/down/volumes:
	docker compose down --volumes

## docker/nuke: stop all services, delete volumes, images and orphan containers
.PHONY: docker/nuke
docker/nuke:
	docker compose down --volumes --rmi all --remove-orphans

## docker/logs: follow logs for all services
.PHONY: docker/logs
docker/logs:
	docker compose logs -f

## docker/logs/api: follow logs for the api service only
.PHONY: docker/logs/api
docker/logs/api:
	docker compose logs -f api

## docker/psql: connect to the containerised database using psql
.PHONY: docker/psql
docker/psql:
	docker compose exec db psql -U $${POSTGRES_USER} -d $${POSTGRES_DB}

## docker/migrate/new name=$1: create a new database migration
.PHONY: docker/migrate/new
docker/migrate/new:
	@echo 'Creating migration files for ${name}...'
	migrate create --seq --ext=.sql --dir=./migrations ${name}

## docker/migrate/up: run migrations against the containerised database
.PHONY: docker/migrate/up
docker/migrate/up:
	docker compose run --rm migrate

## docker/migrate/down: roll back all migrations on the containerised database
.PHONY: docker/migrate/down
docker/migrate/down:
	$(eval include .env)
	docker compose run --rm migrate \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable" \
		down

# ==================================================================================== #
# BUILD
# ==================================================================================== #

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build/api: build the cmd/api application locally
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
