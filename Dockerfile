FROM golang

RUN apt-get update && apt-get install -y redis-server

# Install and configure mysql
RUN apt-get install -y mysql-server
ADD build/my.cnf /etc/mysql/my.cnf
RUN mkdir -p /var/lib/mysql
RUN chmod -R 755 /var/lib/mysql

RUN apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor

ADD . /go/src/github.com/GrappigPanda/notorious
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN go get gopkg.in/redis.v3
RUN go get github.com/spf13/viper
RUN go install github.com/GrappigPanda/notorious

CMD ["/usr/bin/supervisord"]

EXPOSE 3000
