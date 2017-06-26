FROM debian
 
RUN apt-get -y update && \
    apt-get install -y wget git libpq-dev supervisor postgresql-9.6\
                       postgresql-9.6-postgis-2.3 postgresql-9.6-postgis-scripts postgis golang && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV POSTGRES_DB route-profile
ENV DB_NAME route-profile
ENV DB_USER postgres

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
ENV PATH $PATH:/usr/local/bin

ADD . /go/src/github.com/hamil-io
WORKDIR /go/src/github.com/hamil-io

# Build Service
WORKDIR route-profile
RUN go get ./...
RUN go build
WORKDIR ..

# Symlink utils
RUN ln -s /go/src/github.com/hamil-io/utils/wind/load-wind /usr/local/bin/load-wind
RUN ln -s /go/src/github.com/hamil-io/utils/elevation/load-elevation /usr/local/bin/load-elevation

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
