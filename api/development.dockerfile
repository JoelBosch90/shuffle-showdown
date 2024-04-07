FROM golang:1.22.2-alpine3.18

WORKDIR /api
COPY . .

ENV CGO_ENABLED=1
ENV PORT=80

RUN apk add build-base
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# CompileDaemon doesn't seem to always start the binary initially, so we add
# a manual build and run command here to start up the server on container start.
RUN go build -o /go/bin/api
RUN /go/bin/api &

ENTRYPOINT CompileDaemon --build="go build -o /go/bin/api" --command="/go/bin/api"