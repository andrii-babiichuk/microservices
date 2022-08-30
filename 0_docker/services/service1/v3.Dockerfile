# Початковий образ
FROM golang:1.19-alpine3.16

# Встановлення додаткових системних бібліотек та інструментів
RUN apk add gcc libc-dev

# Встановлення робочої директорії
WORKDIR /app

# Встановлення залежностей
COPY go.mod go.sum ./
RUN go mod download

# Копіювання коду
COPY services/service1 services/service1

# Збирання коду
WORKDIR services/service1
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

# Відкиття порту
EXPOSE 8080

# Точка входу в додаток
CMD ["./main"]
