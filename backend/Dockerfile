# syntax=docker/dockerfile:1

FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata \
  && adduser -D -H appuser
WORKDIR /app
COPY --from=build /out/api ./api
COPY migrations ./migrations
COPY seeds ./seeds
USER appuser
EXPOSE 8080
# Start only the API. Run migrations as a separate, explicit command/job.
ENTRYPOINT ["./api"]
