.PHONY: up down migrate

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

migrate:
	docker compose run --rm migration

test: build
	go test -count=1 -v -cover ./...  
	python3 ./scripts/test.py

clean:
	python3 ./scripts/erase_db.py
