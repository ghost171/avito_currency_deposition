start:
	docker-compose -f buildings_config/docker-compose.yml up

stop:
	docker-compose -f buildings_config/docker-compose.yml down

build:
	docker-compose -f buildings_config/docker-compose.yml up --build

migrateup:
	migrate -path migrations -database "postgresql://root:spartak1@localhost:5432/balance?sslmode=disable" up

migratedown:
	migrate -path migrations -database "postgresql://root:spartak1@localhost:5432/balance?sslmode=disable" down