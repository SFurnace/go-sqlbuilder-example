FROM alpine:latest

ENV http_proxy=http://192.168.0.12:7890 https_proxy=http://192.168.0.12:7890
COPY root/* /root
RUN apk update && \
  apk add iproute2-ss go mysql mysql-client && \
  mysql_install_db --user=mysql

WORKDIR /root
ENTRYPOINT ["mysqld_safe"]
