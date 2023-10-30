FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN ls
COPY . .

RUN go build -o /to-do-list

EXPOSE 8080
CMD ["/to-do-list"]
