build:
	docker-compose build avito-app

run:
	docker-compose up -d avito-app

migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

migrate_down:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' down

tests:
	go test ./test

stop:
	docker-compose stop avito-app
	docker-compose down