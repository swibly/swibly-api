FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd/api && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

ENV POSTGRES_HOST=localhost
ENV POSTGRES_DB=arkhon-db
ENV POSTGRES_USER=arkhon-db_owner
ENV POSTGRES_PASSWORD=arkhon
ENV POSTGRES_SSLMODE=disable
ENV JWT_SECRET=jwtsecret

WORKDIR /root/

COPY --from=builder /app/cmd/api/main .

EXPOSE 8080

CMD ["./main"]
