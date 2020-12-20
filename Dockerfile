FROM golang:1.15-buster

ENV DISCORD_BOT_TOKEN null

WORKDIR /go/src/app
COPY . .

RUN apt-get install ca-certificates -y

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["derpbot"]