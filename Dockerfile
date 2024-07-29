FROM golang:1.22-alpine AS build
RUN apk update && apk add --no-cache ca-certificates build-base
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun cmd/app/main.go

# FROM scratch
FROM golang:1.22-alpine
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./cmd/app/.env .
COPY --from=build /app/cloudrun .
ENTRYPOINT ["./cloudrun"]
