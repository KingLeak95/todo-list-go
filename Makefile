BINARY_NAME=to-do-list.out
LOCAL_PG_CONTAINERNAME=todolist
LOCAL_PG_USER=postgres
LOCAL_PG_PASS=postgres
LOCAL_PG_DBNAME=todolist
LOCAL_PG_PORT=5432

all: run

run: build
	./${BINARY_NAME}

test:
	find . -name go.mod -execdir go test ./... \;

build:
	go build -o ${BINARY_NAME} main.go

docker:
	docker build -t to-do-list .

postgres-start:
	docker run -d --name "postgres-${LOCAL_PG_CONTAINERNAME}" \
        -e POSTGRES_USER="${LOCAL_PG_USER}" \
        -e POSTGRES_PASSWORD="${LOCAL_PG_PASS}" \
        -e POSTGRES_DB="${LOCAL_PG_DBNAME}" \
        -p "${LOCAL_PG_PORT}":5432 \
        postgres >&2

postgres-stop:
	docker stop postgres-${LOCAL_PG_CONTAINERNAME} || true

clean:
	docker rm postgres-${LOCAL_PG_CONTAINERNAME} || true 
	go clean
	rm ${BINARY_NAME}
