FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/backend ./cmd/api

FROM alpine:3.20

WORKDIR /app
COPY --from=build /out/backend /app/backend
EXPOSE 8080
CMD ["/app/backend"]

