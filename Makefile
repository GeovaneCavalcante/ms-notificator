
run-api:
	go run -race ./cmd/main.go api

up-localstack:
	docker-compose up -d localstack

dependencies:
	go mod download

up-mongo:
	docker-compose up -d mongodb