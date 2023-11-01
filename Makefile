BINARY_NAME=to-do-list.out

all: run

run: build
	./${BINARY_NAME}

test:
	find . -name go.mod -execdir go test ./... \;

build:
	go build -o ${BINARY_NAME} main.go

docker:
	docker build -t to-do-list .

clean:
	go clean
	rm ${BINARY_NAME}
