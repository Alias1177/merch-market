# Переменные
DB_URL := "postgres://myuser:mypassword@postgres:5432/mydb?sslmode=disable"

# Команды миграции
migrate-up:
	docker-compose run --rm migrate -path=/migrations -database $(DB_URL) up

migrate-down:
	docker-compose run --rm migrate -path=/migrations -database $(DB_URL) down 1

migrate-force:
	docker-compose run --rm migrate -path=/migrations -database $(DB_URL) force ${version}

migrate-status:
	docker-compose run --rm migrate -path=/migrations -database $(DB_URL) version