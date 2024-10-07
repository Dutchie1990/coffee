include .env

stop_containers:
	@echo "Stopping all docker containers"
	@if [ $$(docker ps -q) ]; then \
		echo "Containers found"; \
		docker stop $$(docker ps -q); \
	else \
		echo "No containers found"; \
	fi

create_container:
	docker run --name ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=${USER} -e POSTGRES_PASSWORD=${PASSWORD} -d postgres

create_db:
	docker exec ${DB_DOCKER_CONTAINER} createdb --username=${USER} --owner=${USER} ${DB_NAME}

start_container:
	docker start ${DB_DOCKER_CONTAINER} 

create_migrations:
	sqlx migrate add -r init 

migrate_up:
	docker exec coffee-container sqlx migrate run --database-url "postgres://root:secret@localhost:5432/coffee?sslmode=disable"

migrate_down:
	sqlx migrate revert --database-url "postgres://${USER}:${PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

build: 
	@if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
		echo "Deleted ${BINARY}"; \
	fi

	@echo "Building Binary"
	go build -o ${BINARY} cmd/server/main.go


run: build
	@echo "Starting the server\n" \
	./${BINARY} \
	# @echo "Started the server succesfully" \

stop:
	@echo "Stopping the server\n"
	@kill "${BINARY}"
	@echo "Server stopped"