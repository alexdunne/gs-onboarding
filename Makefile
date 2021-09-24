.PHONY: start consumer

start:
	docker-compose --profile api up

consumer:
	docker-compose run --rm consumer