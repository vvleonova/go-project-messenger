lint:
	golangci-lint run ./... --config=./.golangci.yml

run:
	go run cmd/main.go

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down --volumes --remove-orphans

logs:
	docker-compose logs -f

ps:
	docker-compose ps

tools.install: ## Install tools
	cd tools; \
 	$(GO) mod download; \
 	GOBIN=$(GOBIN) $(GO) generate -tags tools tools.go