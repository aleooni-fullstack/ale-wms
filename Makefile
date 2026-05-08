api:
	cd apps/api && air

db:
	docker compose up -d

migration:
	cd apps/api && migrate create -ext sql -dir migrations $(name)