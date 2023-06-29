
run-api:
	go run -race ./cmd/main.go api

up-localstack:
	docker-compose up -d localstack
