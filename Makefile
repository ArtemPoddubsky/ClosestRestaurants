all:
	@mkdir -p /var/tmp/docker/restaurants/postgresql
	@docker compose --env-file ./.env build
	@docker compose --env-file ./.env up -d

app:
	@docker compose --env-file ./.env up -d app

db:
	@docker compose --env-file ./.env up -d db

start:
	@docker compose --env-file ./.env up -d

stop:
	@docker compose stop

clean: stop
	docker compose down
	@-docker volume rm $$(docker volume ls -q)
	@-docker rmi $$(docker images -q)

fclean: clean
	@-rm -rf /var/tmp/docker/restaurants/postgresql

re: fclean all
