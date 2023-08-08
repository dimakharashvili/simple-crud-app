FROM golang:latest as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:latest
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

CMD ["go", "test", "-v", "./internal/repo/it/postgres_it_test.go"]