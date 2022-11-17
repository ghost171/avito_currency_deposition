FROM golang:1.16.3-alpine3.13


WORKDIR /app

COPY go.mod go.sum ./

RUN go get -d -v ./...
RUN go install -v ./...

COPY . .

CMD cd $(mkdtemp -d); go mod init tmp; go get -v github.com/cespare/reflex; cd /app; reflex -r '(\.go|go\.mod)' -s go run cmd/balance/main.go;