.PHONY: run build test migrate-up migrate-down swag-gen generate-sdk

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

migrate-up:
	# Note: This assumes you have golang-migrate CLI installed or use the internal runner
	# Using the internal runner via a small go script or just running the app is easier for a template
	go run cmd/api/main.go --migrate-up

swag-gen:
	export PATH=$$HOME/go/bin:$$PATH && \
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init -g cmd/api/main.go --output docs
	go run cmd/swagger-enhancer/main.go

generate-sdk:
	cd sdk && bun install && bun run generate

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
