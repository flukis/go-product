PSQL_DOCKER := postgres14
PSQL_DBNAME := ecommerce

pg-create:
	docker run --name ${PSQL_DOCKER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=wap12345 -d postgres:14-alpine

pg-start:
	docker start ${PSQL_DOCKER}

createdb:
	docker exec -it $(PSQL_DOCKER) createdb --username=root --owner=root $(PSQL_DBNAME)

dropdb:
	docker exec -it $(PSQL_DOCKER) dropdb $(PSQL_DBNAME)

run:
	go run app/main.go

.PHONY: pg-create pg-start createdb dropdb run