FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go get -d -v
RUN go build -o server .
# Stage 2
FROM alpine
RUN adduser -S -D -H -h /app satit
USER satit
COPY --from=builder /build/ /app/
WORKDIR /app
