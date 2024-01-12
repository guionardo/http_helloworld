FROM golang:1.18.3-alpine AS builder

RUN adduser -D -g '' elf

WORKDIR /app

COPY go.mod ./

# fetch dependancies
RUN go mod download && \
    go mod verify

# copy the source code as the last step
COPY . .
# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" .


# build a small image
FROM alpine:3.16

WORKDIR /app
RUN mkdir /app/custom_responses

LABEL language="golang"
LABEL org.opencontainers.image.source https://github.com/guionardo/http_helloworld
# import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
# copy the static executable
COPY --from=builder --chown=elf:1000 /app/http_helloworld /app/http_helloworld
# use a non-root user
USER elf


# run app
ENTRYPOINT ["/app/http_helloworld"]