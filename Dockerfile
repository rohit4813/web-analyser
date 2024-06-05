FROM golang:1.22

RUN mkdir /opt/web-analyser
COPY . /opt/web-analyser

WORKDIR /opt/web-analyser

RUN go build -o web-analyser cmd/web/main.go
