###
# Build
###
FROM golang:1.23-alpine3.19 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o webhook ./cmd/server

###
# Runtime
###
FROM gcr.io/distroless/static
COPY --from=build /app/webhook .
ENTRYPOINT ["/webhook"]
