FROM golang:1.19-alpine

RUN go version
ENV GOPATH=/
ENV APP_PATH=/home/src

EXPOSE ${HTTP_SERVER_PORT}
VOLUME $APP_PATH/logs
WORKDIR $APP_PATH

COPY ./ $APP_PATH

RUN go mod download
RUN go build -o ping-pong $APP_PATH/cmd/main.go

ENTRYPOINT ["./ping-pong"]