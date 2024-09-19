# Build app
FROM golang:1.22-alpine as build

WORKDIR /code

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -v -o binary

# Run app
FROM alpine:latest

COPY --from=build /code/binary /usr/local/bin/

EXPOSE 8080

WORKDIR /usr/local/bin/

COPY templates ./templates
COPY static ./static

ENTRYPOINT ["/usr/local/bin/binary", "serve"]
CMD ["--config", "config/config.toml"]
