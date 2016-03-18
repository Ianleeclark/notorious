FROM golang

RUN apt-get update && apt-get install -y redis-server
RUN apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor

ADD . /go/src/github.com/GrappigPanda/notorious
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN go get gopkg.in/redis.v3
RUN go install github.com/GrappigPanda/notorious

CMD ["/usr/bin/supervisord"]

EXPOSE 3000
