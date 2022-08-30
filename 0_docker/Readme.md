#  Робота з *Docker*

## Початковий варіант

Розглянемо приклад простого *Dockerfile* для створення сервісу написаного на *Golang*.

<details>
  <summary>Початковий <i>Dockerfile</i></summary>

```Dockerfile
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
```
</details>

Спершу спробуємо зібрати початковий варіант, виконавши команду:

`docker build -t service1:01 -f services/service1/v1.Dockerfile .`

Запустимо щойно створений образ:

`docker run -p 8080:8080 service1:01`

Та перевіримо роботу, команда `curl http://localhost:8080/ping` має вивести `OK`

Також перевіримо створений образ за допомогою `docker image ls`

```
REPOSITORY    TAG   IMAGE ID       CREATED         SIZE
service1      01    39c89d785a2d   4 minutes ago   812MB
```

Отже, *Dockerfile* написаний правильно, сервіс працює, як і очікується. Але є ряд проблем і даний варіант неоптимальний.

## Оптимізація *Dockerfile*

### 1. Базовий образ

У початковому *Dockerfile* базовий образ беремо з `golang`, це 

Тут є 2 проблеми: 
1) По перше образ без версії, тобто завжди береться останній `latest`
Це означає, що при випуску нової версії буде очищуватись кеш, і в гіршому випадку код може працювати некоректно.

2) `golang` використовує стандартний дистрибутив `linux` у якій встановлено багато додаткових інструментів,
що збільшують розмір образу та збільшують ризики кібератак. Як правило,
якщо немає потреби в різних додаткових інструментах використовується образ з суфіксом `apline`, мінімальну збірку `linux`.

Змінимо базовий образ на `golang:1.19-alpine3.16` та перестворимо образ:

`docker build -t service:02 -f services/service1/v2.Dockerfile .`

Запустимо другу версію `docker run -p 8080:8080 service1:02`, і перевіримо роботу сервісу:

Вивід `docker image ls` буде наступним:

```
REPOSITORY TAG IMAGE ID       CREATED          SIZE
service1   02  58932afd900f   8 seconds ago    320MB
```

Зверніть увагу, як змінився розмір образу, просто використавши образ `apline` вдалось зекономити майже 500МБ! 

### 2. Кешування

У Докер кешування працює по рівнях, наприклад, якщо код сервісу змінився (рядок `COPY services/service1 services/service1`),
всі наступні кроки будуть виконуватись заново, зокрема встановлення залежностей (`RUN go mod download`),
навіть якщо залежності не змінились.

Для оптимізації кешування на початку мають іти такі кроки, які будуть змінюватись рідше.
В нашому випадку спочатку (як і в більшості загальних випадках) це буде встановлення системних бібліотек:

```Dockerfile
RUN apk add gcc libc-dev
```

Потім встановлення залежностей:

```Dockerfile
COPY go.mod go.sum ./
RUN go mod download
```


<details>
<summary><i>Dockerfile</i> з оптимізованим кешуванням:</summary>

```Dockerfile
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
```
</details>

### 3 .dockeringore

При виконанні кроку копіювання `COPY services/service1 services/service1` буде скопійовано не лише код сервісу,
а і різні службові файли, такі як Readme.md, Dockerfile та ін. які не беруть участі у збиранні проекту.
Для оптимізації стадії копіювання можна використати файл `.dockeringore` в директорії, з контексту якої буде відбуватись збирання образу Докер.
і прописати всі виключення так само як і в `.gitignore`.

<details>
    <summary>Приклад:</summary>

```.dockerignore
*.md
services/*/Dockerfile
```
</details>

### 4. Багатоетапна збірка *(multi-stage build)*

Як правило при збиранні проекту можуть використовуватись інструменти, які не потрібні для запуску додатку.

Наприклад, для збирання проекту написаного на `go` потрібен `go`, але результат збірки - це бінарний файл,
який можна запускати без сторонніх платформ. Для запуску такого файлу не потрібні ніякі бібліотеки чи фреймфорки.

В такому випадку, для того, щоб мінімізувати вихідний образ можна розбити *Dockerfile* на дві частини:

В першій частині буде відбувається збирання сервісу.
Ця частина вимагає додаткових залежностей (в даному прикладі це `golang` + та інструменти `gcc` та `libc-dev`).
В цілому це майже попередній варіант, але при використанні багатоетапної збірки потрібно робити мітки (`AS [label]`),
щоб потім в наступних кроках можна було використати результат попередніх частин

```Dockerfile
FROM golang:1.19-alpine3.16 AS server_builder
# server_builder - мітка, яку можна використати в наступних кроках як посилання на першу частину образу

RUN apk add gcc libc-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY services/service1 services/service1

WORKDIR services/service1
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go
```

Друга частина має лише запускати файл, який отримали в результаті збирання.
Більшість залежностей, які використовуються для збирання, для запуску не потрібні, це дозволяє значно зменшити об'єм образу. 
У випадку, наприклад, з сервером на JavaScript мінімальний образ має містити втановлений nodejs для запуску серверу.

```Dockerfile
FROM scratch
EXPOSE 8080
## копіювання даних з попереднього образу за міткою
COPY --from=server_builder /services/service1/main . 
CMD ["./main"]
```

Таким чином ми можемо оптимізувати розмір образу ще більше, і отримати образ розміром всього 4.43МБ:

```
REPOSITORY TAG    IMAGE ID      CREATED         SIZE
service1   final  775b134077e7  6 seconds ago   4.43MB
```

<details>
<summary>Кінцевий <i>Dockerfile</i></summary>

```Dockerfile
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
```
</details>

## Docker для клієнта

Так само як і в прикладі з сервером, *Dockerfile* для клієнта має бути розділений на дві частини:

```Dockerfile
#збірка статичних файлів
FROM node:18.8.0-alpine as build
WORKDIR /app
COPY client/package.json ./
COPY client/package-lock.json ./
RUN npm install
COPY client/. ./
RUN npm run build

#запуск серверу
FROM nginx:stable-alpine
COPY --from=build /app/build /usr/share/nginx/html
# COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

В першій частині виконується збирання проекту. Для цього необхідно використати образ node, також встановити всі залежності та зібрати проект.
Клієнт не виконує ніякої серверної роботи, а просто роздає статичні файли (js, html, css, ...),
тому як сервер ми можемо використати `nginx` і просто скопіювати, отримані на етапі збирання, файли в директорію `/usr/share/nginx/html`.
Це коренева директорія `nginx` з якої будуть роздаватись файли.
За замовчуванням `nginx` буде використовувати свій стандартний файл конфігурацій.
Якщо його потрібно налаштувати, можна в директорії з клієнтським кодом створити файл `/client/nginx/nginx.conf` та задати налаштування,
які необхідні для вашого додатка. І розкоментувати строку в Dockerfile `COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf`

`docker build -t client:0.1 -f client/Dockerfile .`

Тепер запустивши контейнер, `docker run -p 8080:80 -t client:0.1` можна побачити стартову сторінку з написом `Client is up and running`
