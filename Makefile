.PHONY: build push up down logs ps test help
VERSION := $$(git describe --tags --always --dirty)
BRANCH := $$(git symbolic-ref --short HEAD)

build: ## Build docker image to deploy
	docker build \
		-t rindrics/gotodo:${BRANCH} --target deploy \
		-t rindrics/gotodo:latest --target deploy \
		-t rindrics/gotodo:${VERSION} --target deploy ./
	docker compose build --no-cache

push: ## Push build images to DockerHub
	docker push rindrics/gotodo:${BRANCH}
	docker push rindrics/gotodo:latest
	docker push rindrics/gotodo:${VERSION}

up: ## docker compose up with hot reload
	docker compose up -d

down: ## Stop svc deployed by docker compose
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

test: ## Execute tests
	go test -race -v -shuffle=on ./...

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST)| \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n, $$1, $$2}'
