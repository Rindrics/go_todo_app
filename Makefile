.PHONY: build push up down logs ps test help
VERSION := $$(git describe --tags --always --dirty)
BRANCH := $$(git branch --show-current)
.DEFAULT_GOAL := help

build: build-application
build-application: build-%: FORCE ## Build application image to deploy
	docker buildx create --use
	docker buildx build --push \
		--platform=linux/amd64,linux/arm64/v8 \
		-t rindrics/gotodo-$*:${BRANCH} --target deploy \
		-t rindrics/gotodo-$*:latest --target deploy \
		-t rindrics/gotodo-$*:${VERSION} --target deploy -f dockerfiles/$* .

build: build-database
build-database: build-%: FORCE ## Build database image to deploy
	docker buildx create --use
	docker buildx build --push \
		--platform=linux/amd64,linux/arm64/v8 \
		-t rindrics/gotodo-$*:${BRANCH} \
		-t rindrics/gotodo-$*:latest\
		-t rindrics/gotodo-$*:${VERSION} -f dockerfiles/$* .

build: build-compose
build-compose: ## Build for docker compose
	docker compose build --no-cache

push: ## Push build images to DockerHub
	docker push rindrics/gotodo:${BRANCH}
	docker push rindrics/gotodo:latest
	docker push rindrics/gotodo:${VERSION}

up: ## docker compose up with hot reload
	docker compose up -d

down: ## Stop svc deployed by docker compose
	docker compose down

lint: ## Lint codes
	golangci-lint run

generate: ## Generate codes
	go generate ./...

logs: ## Tail docker compose logs
	docker compose logs -f

dry-migrate: ## Try migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < ./_tools/mysql/schema.sql

ps: ## Check container status
	docker compose ps

test: ## Execute tests
	go test -race -v -shuffle=on ./...

deps: ## Install dependencies
	go install github.com/k0kubun/sqldef/cmd/mysqldef@latest

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)| \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

FORCE:
