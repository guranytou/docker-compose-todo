FROM golang:1.16.0-alpine3.13

WORKDIR /go/src/github.com/guranytou/docker-compose-todo

COPY ./backend/app .
RUN go mod download

RUN go get -u github.com/cosmtrek/air

CMD ["air", "-c", ".air.toml"]