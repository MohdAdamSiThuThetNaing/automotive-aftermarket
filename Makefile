up:
	docker compose up -d
down:
	docker compose down
test:
	go test ./... -v
run:
	go run ./cmd/app