build:
	docker build -t omarkhd/memkv:latest .

down:
	docker-compose down --remove-orphans

up: down build
	docker-compose up
