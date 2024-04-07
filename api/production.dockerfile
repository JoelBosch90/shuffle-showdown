FROM golang:1.22.2-alpine3.18 AS build

WORKDIR /api
COPY . .

ENV CGO_ENABLED=1
ENV GIN_MODE=release
ENV GOOS=linux GOARCH=amd64

RUN apk add build-base
RUN go mod download
RUN go build -o /go/bin/api

FROM alpine:3.18

ENV PORT=80

WORKDIR /api
COPY --from=build /go/bin/api /api

CMD ["/api/api"]