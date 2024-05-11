postgres16:
	docker run -d --name postgres16 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secretpassword -e PGDATA=/mnt/postgresql/data -v /home/pinkpanther/docker/postgresql:/mnt/postgresql/ -p 5432:5432 postgres:16.2-alpine3.19
createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simplebank
dropdb:
	docker exec -it postgres16 dropdb simplebank
migratedb-up:
	migrate -path database/migration -database "postgresql://root:secretpassword@localhost:5432/simplebank?sslmode=disable" -verbose up
migratedb-down:
	migrate -path database/migration -database "postgresql://root:secretpassword@localhost:5432/simplebank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
.PHONY: postgres16, createdb, migratedb-up, migratedb-down, sqlc