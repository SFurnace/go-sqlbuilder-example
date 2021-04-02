FROM alpine:latest
ARG BASE_DIR=/root/workspace
ENV GOPATH=/root/.go

RUN apk update && \
  apk add iproute2-ss go mysql mysql-client

COPY root /root
WORKDIR ${BASE_DIR}
RUN mysql_install_db --user=mysql && \
  sed -i -e 's/skip-networking//g' /etc/my.cnf.d/mariadb-server.cnf && \
  sh -c '/usr/bin/mysqld_safe &' && \
  sleep 2 && \
  mysql -u root -e "GRANT ALL PRIVILEGES ON *.* TO 'tester'@'%' IDENTIFIED BY 'tester123'" && \
  mysql -u root -e "GRANT ALL PRIVILEGES ON *.* TO 'tester'@'localhost' IDENTIFIED BY 'tester123'" && \
  go get -u ./... && \
  go test -c -o ./main ./tests && \
  ./main -test.run '^TestCreateDB$' && \
  ./main -test.run '^TestGenerateData$'

CMD ["./start.sh"]
