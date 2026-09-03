postgres:
	podman run --network inventium --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:16-alpine
createdb:
	podman exec -it postgres createdb --username=root --owner=root recipe-service
dropdb:
	podman exec -it postgres dropdb --username=root recipe-service
migrateup:
	migrate -path ./models/migration -database "$(DB_SOURCE)" -verbose up
migratedown:
	migrate -path ./models/migration -database "$(DB_SOURCE)" -verbose down
sqlc:
	sqlc generate --no-remote
loaddata:
	PGPASSWORD="$(DB_PASS)" psql -h "$(DB_HOST)" -U "$(DB_USER)" -p 16677 -d recipe_service -f data/sql/inventium.sql
runcontainer:
	podman run --network inventium --name recipe-service -p 9820:9820 -d -e DB_SOURCE="$(DB_SOURCE)" recipe-service:1.0.0
.PHONY: postgres createdb dropdb migrateup migratedown sqlc loaddata runcontainer