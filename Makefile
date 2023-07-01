
setup-dev:
	docker-compose up -d
	make up-client-web
	make configure-sns
	make up-app-api

run-api:
	go run -race ./cmd/main.go api

run-worker:
	go run -race ./cmd/main.go worker

configure-sns:
	docker exec localstack_main sh -c "awslocal sns create-topic --name notifications"
	docker exec localstack_main sh -c "awslocal sns subscribe --topic-arn arn:aws:sns:us-east-1:000000000000:notifications --protocol http --notification-endpoint http://172.28.3.10:8083"

up-localstack:
	docker-compose up -d localstack

up-client-web:
	docker build -t client-web-app -f client-web/Dockerfile .
	docker run -d --name client-web-app --net=network-notificator --ip=172.28.3.10 -p 8083:8083 client-web-app 

dependencies:
	go mod download

up-mongo:
	docker-compose up -d mongodb

up-app-api:
	docker build -t ms-notificator .
	docker run -d --env APP_COMMAND=api --name notificator-api --net=network-notificator --ip=172.28.3.8 -p 8081:8081 ms-notificator

	docker run -d --env APP_COMMAND=worker --name notificator-worker --net=network-notificator ms-notificator

down-all:
	docker stop notificator-worker && docker rm notificator-worker || true
	docker stop notificator-api && docker rm notificator-api || true
	docker stop client-web-app && docker rm client-web-app || true
	docker-compose down || true

format:
	go fmt ./...
	
test:
	go test -v -tags testing ./...

test-cov:
	go test -coverprofile=cover.txt ./... && go tool cover -html=cover.txt -o cover.html
