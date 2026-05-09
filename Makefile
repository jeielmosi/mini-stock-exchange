.PHONY: up down migrate

up:
	docker compose up -d

down:
	docker compose down

migrate:
	docker compose run --rm migration

script:
	python3 ./scripts/test_order.py
