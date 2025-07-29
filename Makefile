DB_URL="postgresql://postgres:123456qwerty@postgres:5432/currency_db?sslmode=disable"

gen-proto:
	protoc --go_out=./proto --go-grpc_out=./proto proto/rate.proto
migrate-create-%:
	goose -dir migrations create $(subst migrate-create-,,$@) sql	
migrate-up:
	goose -dir migrations postgres "$(DB_URL)" up
migrate-down:
	goose -dir migrations postgres "$(DB_URL)" down	
run:
	go run cmd/app/main.go --config=configs/local.yml
build:
	go build -o ./cmd/bin ./cmd/app
test:
	go test ./internal/service -v -count=1 
lint:
	golangci-lint run ./...
docker-up:
	docker compose up -d	
docker-build:
	docker compose build
docker-down:
	docker compose down 	

