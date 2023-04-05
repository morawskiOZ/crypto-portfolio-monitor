FROM golang:1.20-alpine AS build_base

WORKDIR /app

RUN apk add bash ca-certificates gcc g++ libc-dev

COPY ./go.mod ./go.sum ./

RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN CGO_ENABLED=0 go build -o "/go/bin/service" -a -tags netgo -ldflags '-w -extldflags "-static"' "./cmd/main.go";

FROM alpine AS runner

RUN apk add ca-certificates

RUN addgroup -S appme && adduser -S appme -G appme
USER appme

COPY --from=server_builder "/go/bin/service" "/bin/service"
COPY --from=server_builder /app/config/ /config/

EXPOSE 8080

ENTRYPOINT ["/bin/service"]