# --- build stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
# `go mod tidy` resolves the dependency graph and writes go.sum at build time,
# so we don't have to commit a go.sum for this small lab service.
RUN go mod tidy && CGO_ENABLED=0 go build -o /out/service-go .

# --- run stage ---
FROM alpine:3.20
COPY --from=build /out/service-go /service-go
EXPOSE 8080
ENTRYPOINT ["/service-go"]
