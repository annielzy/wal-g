FROM wal-g/ubuntu:20.04

# mb not another path
ENV MYSQLDATA /var/lib/mysql

RUN apt-get update && \
    apt-get install --yes --no-install-recommends --no-install-suggests software-properties-common dirmngr gnupg2 curl apt-transport-https && \
    curl -sSL https://downloads.mariadb.com/MariaDB/mariadb_repo_setup | /bin/bash -s -- --mariadb-server-version="mariadb-10.11" --skip-maxscale && \
    DEBIAN_FRONTEND=noninteractive apt-get install --yes \
    mariadb-server \
    mariadb-client \
    mariadb-backup

RUN curl -s https://packagecloud.io/install/repositories/akopytov/sysbench/script.deb.sh | bash && apt -y install sysbench
RUN rm -rf $MYSQLDATA

COPY docker/mariadb/client.cnf /root/.my.cnf
COPY docker/mariadb/client.cnf /etc/mysql/debian.cnf
COPY docker/mariadb/init.sql /etc/mysql/init.sql
COPY docker/mariadb/export_common.sh /usr/local/export_common.sh

# append
COPY docker/mariadb/my.cnf /tmp/my.cnf
RUN cat /tmp/my.cnf >> /etc/mysql/my.cnf; rm /tmp/my.cnf
