FROM alpine:latest
ARG BASE_DIR=/root/workspace
ENV GOPATH=/root/.go

COPY root /root
RUN apk update && \
  apk add iproute2-ss go mysql mysql-client

WORKDIR ${BASE_DIR}
RUN mysql_install_db --user=mysql && \
  sh -c '/usr/bin/mysqld_safe &' && \
  sleep 2 && \
  mysql -u root -e "GRANT ALL PRIVILEGES ON *.* TO 'tester'@'%' IDENTIFIED BY 'tester123'" && \
  mysql -u root -e "GRANT ALL PRIVILEGES ON *.* TO 'tester'@'localhost' IDENTIFIED BY 'tester123'" && \
  sed -i -e 's/skip-networking//g' /etc/my.cnf.d/mariadb-server.cnf && \
  go get -u ./... && \
  go mod download

WORKDIR ${BASE_DIR}/tests
ENTRYPOINT ["mysqld_safe"]
