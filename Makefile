start:
	docker-compose -f buildings_config/docker-compose.yml up

stop:
	docker-compose -f buildings_config/docker-compose.yml down

build:
	docker-compose -f buildings_config/docker-compose.yml up --build