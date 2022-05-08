FROM golang:1.17 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -o monitor

FROM alpine:3.14 as production
COPY --from=builder /app/cmd/monitor /app/prod.env /
ARG app_env
ENV APP_ENV $app_env
ENV EMAIL_PASS=""
ENV EMAIL_LOGIN=""
ENV EMAIL_SMTP_PORT="587"
ENV EMAIL_SMTP_HOST="smtp.gmail.com"
ENV EMAIL_RECIPIENT=""
# Exec built binary
CMD ["./monitor"]
