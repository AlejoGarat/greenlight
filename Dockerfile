FROM golang AS builder
WORKDIR /go/src/greenlight

## This will make installing dependencies much faster if using vendoring
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
# Tidy imports, if using vendoring, make sure to use go mod vendor
RUN go mod tidy


################# BUILD BINARY ##################
## You could use the -ldflags="-s" flag, this deletes the symbol table and debug information
## from the binary, making it smaller, but you won't be able to debug it if something goes wrong.
## Also, panic messages won't be shown, so be careful.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/example

################# MINIMAL IMAGE #################
FROM scratch
## Copy binary from builder to next image
## If you changed WORKDIR in builder, make sure to change it here too
COPY --from=builder /go/src/greenlight/main main
## If you have configuration files, you can copy them here
# COPY --from=builder /go/src/greenlight/cmd/example/config.yaml config.yaml

############## ENVIRONMENT VARS #################
# ENV SENTRY_DSN=''
# ENV NO_SENTRY=1
# ENV SENTRY_DEBUG=0
# ENV GIN_MODE=release
# ENV REDIS_URL=127.0.0.1
# ENV REDIS_PORT=6379
# ENV PORT=80

################ EXPOSE PORTS ################
## Not mandatory to expose default port, but recommended
EXPOSE 80

################# RUN BINARY #################
## Entrypoint to the exceutable
ENTRYPOINT ["./main"]
