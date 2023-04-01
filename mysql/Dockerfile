FROM mysql:latest

#RUN ln -snf /usr/share/zoneinfo/UTD /etc/localtime && echo UTC > /etc/timezone && chown -R mysql:root /var/lib/mysql/

COPY my.cnf /etc/mysql/conf.d/my.cnf

RUN chmod 0444 /etc/mysql/conf.d/my.cnf

CMD ["mysqld"]

EXPOSE 3306
