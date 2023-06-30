FROM golang:1.20-alpine AS build

WORKDIR /src

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/app cmd/main.go

FROM alpine:3.14

EXPOSE 8081

WORKDIR /app

COPY --from=build /src/bin/app .

COPY --from=build /src/docker/docker-entrypoint.sh .

COPY --from=build /src/.env .

ENTRYPOINT ["/bin/sh", "docker-entrypoint.sh"]
