BINARY_NAME=to-do-list.out
LOCAL_PG_CONTAINERNAME=todolist
LOCAL_PG_USER=postgres
LOCAL_PG_PASS=postgres
LOCAL_PG_DBNAME=todolist
LOCAL_PG_PORT=5432
LOCAL_NETWORK=todolist
LOCAL_APP_CONTAINERNAME=todolist

all: run

.PHONY: all run test build docker-build docker-start postgres-start postgres-stop clean docker-compose-up docker-compose-down
DOCKER_TEST_IMAGE=to-do-list-test:latest

run: build
	./${BINARY_NAME}

test:
	go test ./...

build:
	go build -o ${BINARY_NAME} main.go

docker-build:
	docker build -t to-do-list .

docker-test-build:
	docker build -f Dockerfile.test -t ${DOCKER_TEST_IMAGE} .

docker-test: docker-test-build
	docker run --rm ${DOCKER_TEST_IMAGE}

docker-start: postgres-start 	
	docker run --name "app-${LOCAL_APP_CONTAINERNAME}" \
	--network="${LOCAL_NETWORK}" \
	-p 8080:8080 \
	-e DB_HOST="postgres-${LOCAL_PG_CONTAINERNAME}"
	to-do-list

postgres-start:
	docker run -d --name "postgres-${LOCAL_PG_CONTAINERNAME}" \
        -e POSTGRES_USER="${LOCAL_PG_USER}" \
        -e POSTGRES_PASSWORD="${LOCAL_PG_PASS}" \
        -e POSTGRES_DB="${LOCAL_PG_DBNAME}" \
        -p "${LOCAL_PG_PORT}":5432 \
				--network="${LOCAL_NETWORK}" \
        postgres >&2 || true

postgres-stop:
	docker stop postgres-${LOCAL_PG_CONTAINERNAME} || true

clean: postgres-stop
	docker rm postgres-${LOCAL_PG_CONTAINERNAME} || true 
	go clean
	rm ${BINARY_NAME}

docker-compose-up:
	docker-compose up -d

docker-compose-down:
	docker-compose down

docker-compose-logs:
	docker-compose logs -f
