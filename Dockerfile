# Dockerfile.deploy

FROM golang:1.14 as builder

ENV APP_USER app
ENV APP_HOME /go/src/github.com/prateekgupta3991/refresher

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER
COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o refresher

FROM debian:buster

ENV APP_USER app
ENV APP_HOME /go/src/github.com/prateekgupta3991/refresher

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

# Add the following command to download CA certificates
RUN apt-get update && apt-get install -y ca-certificates

COPY --chown=0:0 --from=builder $APP_HOME/refresher $APP_HOME
COPY --chown=0:0 --from=builder $APP_HOME/configs $APP_HOME/configs

EXPOSE 8080
USER $APP_USER
CMD ["./refresher"]