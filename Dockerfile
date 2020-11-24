FROM golang:1.15

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/gocql/gocql

ENV GO111MODULE=on
ENV APP_USER app
ENV APP_HOME /go/src/github.com/prateekgupta3991/refresher

ARG GROUP_ID
ARG USER_ID

RUN groupadd --gid $GROUP_ID app && useradd -m -l --uid $USER_ID --gid $GROUP_ID $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

USER $APP_USER
WORKDIR $APP_HOME
EXPOSE 8080
CMD ["go", "run","hello.go"]