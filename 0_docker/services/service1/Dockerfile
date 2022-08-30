FROM golang:1.19-alpine3.16 AS server_builder

RUN apk add gcc libc-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY services/service1 services/service1

WORKDIR services/service1
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

FROM scratch
EXPOSE 8080
COPY --from=server_builder /services/service1/main .
CMD ["./main"]
