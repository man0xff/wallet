.DEFAULT_GOAL = run

.PHONY: PHONY
PHONY:

UID := $(shell id -u)
GID := $(shell id -g)
export UID GID

$(shell mkdir -p .tmp.home)

docker_exec := docker exec -e HOME=/x/.tmp.home -it -u $(UID):$(GID) \
	-e APP_DB_ADDR=database \
	-e APP_DB_USER=app \
	-e APP_DB_PASSWORD=gfhjkm-app \
	-e APP_DB_NAME=example

env: PHONY
	-mkdir -p .tmp.home/db
	docker-compose -f docker/docker-compose.yaml --project-directory docker up

run: PHONY
	$(docker_exec) example_app_1 go run .

shell: PHONY
	$(docker_exec) example_app_1 bash

db:
	$(docker_exec) example_database_1 mysql -uroot -pgfhjkm-root example

test: PHONY
	$(docker_exec) example_app_1 go test -v -count 1 -tags test_env ./http

vendor: PHONY
	go mod vendor
	go mod tidy
