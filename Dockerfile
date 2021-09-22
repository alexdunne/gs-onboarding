# certs
FROM alpine:3.13 as certs
RUN apk add -U --no-cache ca-certificates

# builder
FROM golang:1.17 as build
ARG cmd
WORKDIR /app
# Download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy over the rest of the code
COPY . .
# Build the specified cmd
RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o "./bin/${cmd}" "./cmd/${cmd}"

FROM scratch as cmd
ARG cmd
COPY .env .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build "/app/bin/${cmd}" .