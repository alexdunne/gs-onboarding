# user
FROM alpine:3.13.1 as user
ARG uid=10001
ARG gid=10001
RUN echo "scratchuser:x:${uid}:${gid}::/home/scratchuser:/bin/sh" > /scratchpasswd

# certs
FROM alpine:3.13 as certs
RUN apk add -U --no-cache ca-certificates

# builder
FROM golang:1.17 as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# build all of the images in one go to save redownloading dependencies each time
RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o ./bin/api ./cmd/api
RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o ./bin/consumer ./cmd/consumer
RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o ./bin/gateway ./cmd/gateway
RUN GOOS=linux CGO_ENABLED=0 GOGC=off GOARCH=amd64 go build -o ./bin/migrator ./cmd/migrator

# entrypoints
FROM scratch as api
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /app/bin/api .
USER scratchuser
EXPOSE 8001
ENTRYPOINT ["/api"]

FROM scratch as consumer
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /app/bin/consumer .
COPY --from=build /app/migrations/ ./migrations/
USER scratchuser
ENTRYPOINT ["/consumer"]

FROM scratch as gateway
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /app/bin/gateway .
USER scratchuser
EXPOSE 8000
ENTRYPOINT ["/gateway"]

FROM scratch as migrator
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=user /scratchpasswd /etc/passwd
COPY --from=build /app/bin/migrator .
COPY --from=build /app/migrations/ ./migrations/
USER scratchuser
ENTRYPOINT ["/migrator"]