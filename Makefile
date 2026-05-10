.PHONY: up down migrate

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

migrate:
	docker compose run --rm migration

script: build
	python3 ./scripts/test_order.py
