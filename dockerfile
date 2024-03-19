FROM golang:1.22

WORKDIR /app

COPY go.mod .
COPY cmd/main.go .

RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]