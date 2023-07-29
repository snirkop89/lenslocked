## connect: connect into the database docker container
.PHONY: connect
connect:
	@docker exec -it $$(docker ps -f name=lens-db -q) /usr/bin/psql -U baloo -d lenslocked

compose-dev:
	@docker-compose up -d 

compose-prod:
	@docker-compose -f docker-compose.yaml -f docker-compose.production.yml up -d