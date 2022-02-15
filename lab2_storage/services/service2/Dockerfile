FROM golang:1.17.5-alpine3.15 AS service_builder

WORKDIR /

# встановлення додаткових інструментів та бібліотек
RUN apk add gcc libc-dev

# встановлення залежностей
COPY go.mod go.sum ./
RUN go mod download

# копіювання основного коду сервісу
COPY services/service2 services/service2

# збарання сервісу
WORKDIR /services/service2
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

# підготовка фінального образу
FROM scratch
EXPOSE 8080
COPY --from=service_builder /services/service2/main .
CMD ["./main"]
