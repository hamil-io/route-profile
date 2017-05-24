# Overview

This is the route-profile service. It provides an API for returning the elevation and headwind along a route.

## Data Sources

Before you can load data from any of the data sources you must be sure to create a spatially enabled database

```
CREATE USER username WITH PASSWORD 'password';
CREATE DATABASE route-profile WITH OWNER username;
\connect route-profile
CREATE EXTENSION postgis;
```

### Elevation

Elevation data comes from NASA's [SRTM](http://www2.jpl.nasa.gov/srtm/). To load the data we use `utils/elevation/load-elevation`. There are three parts to loading the data, first you must download the SRTM tiles, then create your database/tables, and finally import the raster tiles.

#### Download

There is a utility provided for downloading the elevation tiles and getting it ready for import into postgres. Downloading the tiles requires a NASA Earthdata login which can be obtained [here](https://urs.earthdata.nasa.gov/). The `load-elevation` utility will download tiles for the entire globe unless you provide constraints. See `load-elevation --help` for more info.

```
./utils/elevation/load-elevation -u user -p password --latitude-min=20 --latitude-max=30 --out=data.sql
```

#### Load Data

Once you've generated you SQL data using the `load-elevation` utility you can simply import it directly into postgres.

`psql < data.sql`

### Wind

Wind data comes from the NOAA High Resolution Rapid Refresh dataset. Specifically we use the 2d Surface Levels - Sub Hourly product.

The NOAA uses a non-standard projection for this data so you must create the projection in Postgis first:

`psql database < db/projection.sql`

To download and load the data simply run the provided shell script:

`utils/wind/load-wind`

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
* resolution: An optional parameter that determines the sampling resolution. This is provided in degrees, the default is calculated by length(geometry)/100 with a minimum of 0.00025 degrees.

### Wind Endpoint

```curl http://localhost:8080/wind -d '{"geometry":"{``mE|h}aOgFEApECtE?p@AbDAfD?\\At@A~DAxK?d@A\\Cl@?nD?P?L?LB`@f@tFPjC@P@Lx@dJ@F?N?VCdJ?R?VChI?P?FeECy@?}FEqFCM?E?CpF?v@ChHBlHAhG@^D?nB@?P@PONWAGF"}'```

### Elevation Endpoint

```curl http://localhost:8080/elevation -d '{"geometry":"{``mE|h}aOgFEApECtE?p@AbDAfD?\\At@A~DAxK?d@A\\Cl@?nD?P?L?LB`@f@tFPjC@P@Lx@dJ@F?N?VCdJ?R?VChI?P?FeECy@?}FEqFCM?E?CpF?v@ChHBlHAhG@^D?nB@?P@PONWAGF"}'```
