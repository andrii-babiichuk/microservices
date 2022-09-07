FROM golang:1.15 AS server_builder

WORKDIR /
COPY services/root-service/main.go .
RUN go build -ldflags "-w -s -linkmode external -extldflags -static" -a main.go

FROM scratch
EXPOSE 8080
COPY --from=server_builder /main .
CMD ["./main"]