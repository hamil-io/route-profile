FROM debian
 
RUN apt-get update -y && apt-get install -y software-properties-common gnupg && \
    apt-key adv --keyserver pgp.mit.edu --recv-keys D0E480B0 && \
    add-apt-repository -y "deb http://repo.hamil.io stretch main" && \
    apt-get -y update && \
    apt-get install -y wget supervisor route-profile && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV POSTGRES_DB route-profile
ENV DB_NAME route-profile
ENV DB_USER postgres
ENV PATH $PATH:/usr/local/bin

ADD . /go/src/github.com/hamil-io
WORKDIR /go/src/github.com/hamil-io

# Create Postgres functions
RUN cat db/projection.sql >> docker/init.sql
RUN cat db/segments.sql >> docker/init.sql
RUN cat db/interpolate.sql >> docker/init.sql
RUN cat db/drape.sql >> docker/init.sql
RUN cat db/profile.sql >> docker/init.sql
RUN cat db/wind.sql >> docker/init.sql

# Load Data
WORKDIR /go/src/github.com/hamil-io
RUN mkdir -p /var/log/route-profile
RUN ./docker/load.sh

# Setup Postgres
RUN mkdir -p /var/run/postgresql/9.6-main.pg_stat_tmp/
RUN chown postgres:postgres /var/run/postgresql/9.6-main.pg_stat_tmp/
USER postgres
RUN /etc/init.d/postgresql start &&\
    psql postgres < docker/init.sql
USER root

EXPOSE 8080

CMD supervisord -n -e debug -c /go/src/github.com/hamil-io/docker/supervisord.conf
