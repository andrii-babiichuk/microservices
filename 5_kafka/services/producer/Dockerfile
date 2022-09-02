FROM golang:1.19-alpine3.16 AS service_builder

WORKDIR /

# встановлення додаткових інструментів та бібліотек
RUN apk add gcc libc-dev

# встановлення залежностей
COPY go.mod go.sum ./
RUN go mod download

# копіювання основного коду сервісу
COPY services/producer services/producer

# збарання сервісу
WORKDIR /services/producer
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

# підготовка фінального образу
FROM scratch
EXPOSE 8080
COPY --from=service_builder /services/producer/main .
CMD ["./main"]
