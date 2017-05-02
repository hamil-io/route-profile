FROM postgres:9.6
 
RUN apt-get -y update && \
    apt-get install -y wget sudo git python-pip libpq-dev python-dev \
                       postgresql-9.6-postgis-2.3 postgresql-9.6-postgis-scripts postgis golang && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV POSTGRES_DB route-profile
ENV DB_NAME route-profile
ENV DB_USER postgres

ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin

ADD . /go/src/github.com/hamil-io
WORKDIR /go/src/github.com/hamil-io

# Build Service
WORKDIR route-profile
RUN go get ./...
RUN go build
RUN pwd && ls -alh
WORKDIR ..

RUN pip install -r docker/requirements.txt

EXPOSE 8080

# Create Postgres functions
RUN cat db/projection.sql >> docker/init.sql
RUN cat db/interpolate.sql >> docker/init.sql
RUN cat db/drape.sql >> docker/init.sql
RUN cat db/profile.sql >> docker/init.sql
RUN cat db/wind.sql >> docker/init.sql

# Load Data
WORKDIR utils/elevation/srtm2postgis
RUN python download.py North_America

WORKDIR /go/src/github.com/hamil-io
RUN ./docker/load.sh

COPY docker/init.sql /docker-entrypoint-initdb.d/

CMD /go/src/github.com/hamil-io/route-profile/route-profile
