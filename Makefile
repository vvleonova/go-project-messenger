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

GOBIN := $(CURDIR)\bin
tools.install:
	cd tools && go mod download
	set GOBIN=$(GOBIN) && go generate -tags tools ./tools