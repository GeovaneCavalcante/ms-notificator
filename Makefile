
run-api:
	go run -race ./cmd/main.go api

run-worker:
	go run -race ./cmd/main.go worker

up-localstack:
	docker-compose up -d localstack

dependencies:
	go mod download

up-mongo:
	docker-compose up -d mongodb