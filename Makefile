## connect: connect into the database docker container
.PHONY: connect
connect:
	@docker exec -it $$(docker ps -f name=lens-db -q) /usr/bin/psql -U baloo -d lenslocked