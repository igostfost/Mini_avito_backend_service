build:
	docker-compose build avito-app

run:
	docker-compose up -d avito-app

migrate:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

stop:
	docker-compose stop avito-app
	docker-compose down