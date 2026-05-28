
up:
	docker compose up --build -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

restart:
	docker compose restart

test:
	docker compose run --rm app go test ./... -v

run:
	docker compose up --build

clean:
	docker compose down -v
