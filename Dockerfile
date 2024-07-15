FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd/api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

ENV POSTGRES_HOST=localhost
ENV POSTGRES_DB=swibly-db
ENV POSTGRES_USER=swibly-db_owner
ENV POSTGRES_PASSWORD=swibly
ENV POSTGRES_SSLMODE=disable
ENV POSTGRES_PORT=5432
ENV JWT_SECRET=jwtsecret

WORKDIR /root/

COPY --from=builder /app/cmd/api/main .
COPY config /root/config
COPY translations /root/translations

EXPOSE 8080

CMD ["./main"]
