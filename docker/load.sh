#!/bin/bash

BASE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
URL="http://www.ftp.ncep.noaa.gov/data/nccf/com/hrrr/prod"
FILE=`date -u +"hrrr.%Y%m%d/hrrr.t01z.wrfsubhf01.grib2"`
DATA="../utils/elevation/srtm2postgis/data/North_America"

cd $BASE
rm hrrr.t* 

wget "$URL/$FILE"

echo "Download sucessful!"
echo "Exporting wind rasters"

raster2pgsql -c -I -M -b 15,16,56,57,97,98,138,139 -s 98411 -t 10x10 hrrr.t* wind >> init.sql

echo "Creating elevation table..."
raster2pgsql -c -I -M N60W174.hgt route-profile.elevation >> init.sql
echo "Exporting elevation rasters..."
for f in "$DATA/*.hgt"; do echo "Processing $f" && raster2pgsql -a -I -M $f route-profile.elevation >> init.sql; done
echo "Finished!"

