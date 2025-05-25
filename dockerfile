FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main .

FROM alpine:latest  

COPY --from=builder /app/main .

EXPOSE 2112

CMD ["./main"]

