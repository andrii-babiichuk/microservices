# Початковий образ
FROM golang

# Встановлення робочої директорії
WORKDIR /app

# Копіювання коду
COPY go.mod go.sum ./
COPY services/service1 services/service1

# Встановлення залежностей
RUN go mod download

# Встановлення додаткових системних бібліотек та інструментів
RUN apk add gcc libc-dev

# Збирання коду
WORKDIR services/service1
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

# Відкиття порту
EXPOSE 8080

# Точка входу в додаток
CMD ["./main"]
