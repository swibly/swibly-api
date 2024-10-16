FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd/api && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk add --no-cache ca-certificates

ENV POSTGRES_CONNECTION_STRING=str
ENV JWT_SECRET=str
ENV SMTP_HOST=str
ENV SMTP_PORT=str
ENV SMTP_USERNAME=str
ENV SMTP_EMAIL=str
ENV SMTP_PASSWORD=str

WORKDIR /root/

COPY --from=builder /app/cmd/api/main .
COPY config /root/config
COPY translations /root/translations

EXPOSE 8080

CMD ["./main"]
