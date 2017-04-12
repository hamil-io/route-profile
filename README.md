# Overview

This is the route-profile service. It provides an API for returning the elevation and headwind along a route.

## Data Sources

### Elevation

Elevation data comes from NASA's [SRTM](http://www2.jpl.nasa.gov/srtm/). To load the data we use `utils/elevation/srtm2postgis`. There are three parts to loading the data, first you must download the SRTM tiles, then create your database/tables, and finally import the raster tiles.

#### Download

`python download.py North_America`

#### Create Database

```
CREATE USER username WITH PASSWORD 'password';
CREATE DATABASE route-profile WITH OWNER username;
\connect route-profile
CREATE EXTENSION postgis;```
```

#### Load Data

`sudo -u username python read_data_pg.py North_America`

Note: username should be the Postgres user you created

### Wind

Wind data comes from the NOAA High Resolution Rapid Refresh dataset. Specifically we use the 2d Surface Levels - Sub Hourly product.

The NOAA uses a non-standard projection for this data so you must create the projection in Postgis first:

`psql database < db/projection.sql`

To download and load the data simply run the provided shell script:

`utils/wind/load.sh`

The load script will attempt to load the data for the next hour. If the data is not yet available on the NOAA's FTP it will retry until the data becomes available. The data is updated by the NOAA hourly so we recommend setting up a simple crontab:

`56 * * * * /path/to/project/utils/wind/load.sh`

## Route Profile Service

The route-profile service is written in Go. The service primarily handles the incoming API requests, the geospatial processing is implemented primarily in Postgres via Postgis. To get the service running first load the necessary postgres functions:

`psql database < db/*.sql`

You can then build and run the service:

```
cd route-service
go build
./route-service
```

## API

The API itself is fairly simple. There are two endpoints that each take the same parameters:
* geometry: An encoded polyline representing the route
* resolution: An optional parameter that determines the sampling resolution. This is provided in arcseconds, the default is 0.0025.

### Wind Endpoint

`curl http://localhost:8080/wind -d '{"geometry":"{``mE|h}aOgFEApECtE?p@AbDAfD?\\At@A~DAxK?d@A\\Cl@?nD?P?L?LB`@f@tFPjC@P@Lx@dJ@F?N?VCdJ?R?VChI?P?FeECy@?}FEqFCM?E?CpF?v@ChHBlHAhG@^D?nB@?P@PONWAGF"}'`

### Elevation Endpoint

`curl http://localhost:8080/elevation -d '{"geometry":"{``mE|h}aOgFEApECtE?p@AbDAfD?\\At@A~DAxK?d@A\\Cl@?nD?P?L?LB`@f@tFPjC@P@Lx@dJ@F?N?VCdJ?R?VChI?P?FeECy@?}FEqFCM?E?CpF?v@ChHBlHAhG@^D?nB@?P@PONWAGF"}'`
