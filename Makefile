DB_URL=postgres://wms:wms@localhost:5432/wms?sslmode=disable

build-api:
	cd apps/api && go build ./...

api:
	cd apps/api && air

db:
	docker compose up -d

migration:
	cd apps/api && migrate create -ext sql -dir migrations $(name)

migrate-up:
	cd apps/api && migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	cd apps/api && migrate -path migrations -database "$(DB_URL)" down 1

sql-generate:
	cd apps/api && sqlc generate	