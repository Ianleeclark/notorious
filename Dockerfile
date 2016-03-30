FROM golang

RUN apt-get update && apt-get install -y redis-server

# Install and configure mysql
RUN echo 'mysql-server mysql-server/root_password password secret_password' | debconf-set-selections
RUN echo 'mysql-server mysql-server/root_password_again password secret_password' | debconf-set-selections
RUN apt-get install -y mysql-server
ADD build/my.cnf /etc/mysql/my.cnf
RUN mkdir -p /var/lib/mysql
RUN chmod -R 755 /var/lib/mysql

RUN apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor

ADD . /go/src/github.com/GrappigPanda/notorious
COPY build/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

RUN go get gopkg.in/redis.v3
RUN go get github.com/NotoriousTracker/gorm
RUN go get github.com/NotoriousTracker/viper
RUN go install github.com/GrappigPanda/notorious

CMD ["/usr/bin/supervisord"]

EXPOSE 3000
