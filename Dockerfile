FROM golang:1.21-alpine3.18 as builder
WORKDIR /go/src

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./

COPY . ./
RUN go build -o /app/app ./main.go

COPY ./config.json /app/

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app ./
ENTRYPOINT ["./app"]
