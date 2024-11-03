postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=4713a4cd628778cd1c37a95518f3eaf3 -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root Instashop_DB

dropdb:
	docker exec -it postgres dropdb Instashop_DB

create-migration:
	@if [ -n "$(name)" ]; then \
		migrate create -ext sql -dir migrations -seq $(name); \
	else \
		echo "Error: Missing 'name' variable" >&2; \
		echo "Usage: make create-migration name=\"<NAME_OF_MIGRATION_FILE>\"" >&2; \
		exit 1; \
	fi

migrate-up:
	migrate -path migrations -database "postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Instashop_DB?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgresql://root:4713a4cd628778cd1c37a95518f3eaf3@localhost:5432/Instashop_DB?sslmode=disable" -verbose down

mock-repo:
	mockgen -package mocked -destination internal/mock/user_repo.go  github.com/zde37/instashop-task/internal/repository Repository

mock-service:
	mockgen -package mocked -destination internal/mock/user_service.go  github.com/zde37/instashop-task/internal/service Service

test:
	go test -v -cover -short -count=1 ./...
	 
run:
	go run cmd/main.go

docs:
	swag init -g cmd/main.go -o docs

build-run:
	go build -o instashop cmd/main.go && ./instashop

.PHONY: postgres createdb dropdb createmigration migrate-up migrate-down mock-repo mock-service test docs run build-run