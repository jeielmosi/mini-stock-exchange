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
	python3 ./scripts/test.py

clean:
	python3 ./scripts/erase_db.py
