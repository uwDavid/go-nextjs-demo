server: 
	go run main.go
	
c-m: 
	# create a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

m-up: 
	# migrate psql init
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/testdb?sslmode=disable" up
	migrate -path db/migrations -database "postgres://root:secret@localhost:5433/testdb?sslmode=disable" up

m-down: 
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/testdb?sslmode=disable" down
	migrate -path db/migrations -database "postgres://root:secret@localhost:5433/testdb?sslmode=disable" down

dock-up: 
	docker-compose up -d

dock-down: 
	docker-compose down

db-up: 
	docker exec -it psdb createdb --username=root --owner=root testdb
	docker exec -it psdb_live createdb --username=root --owner=root testdb

db-down: 
	docker exec -it psdb dropdb --username=root testdb
	docker exec -it psdb_live dropdb --username=root testdb

sqlc: 
	sqlc generate

start: 
	CompileDaemon -command="./nextjs"

test: 
	go test -v -failfast ./db/tests
