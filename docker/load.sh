#!/bin/bash

BASE="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
URL="http://www.ftp.ncep.noaa.gov/data/nccf/com/hrrr/prod"

YEAR=`date -u +"%Y"`
MONTH=`date -u +"%m"`
DAY=`date -u +"%d"`
DAY=$(( ${DAY} - 1 ))
FILE="hrrr.${YEAR}${MONTH}${DAY}/hrrr.t00z.wrfsubhf01.grib2"

cd $BASE
rm -f hrrr.t* 

echo "fetching $URL/$FILE"
wget --quiet "$URL/$FILE"

echo "Download sucessful!"
echo "Exporting wind rasters"

raster2pgsql -c -I -M -Y -b 15,16,56,57,97,98,138,139 -s 98411 -t 10x10 hrrr.t* wind >> init.sql

echo "Extracting elevation table..."
tar -xzvf elevation.tar.gz
cat docker.sql >> init.sql
rm docker.sql

echo "Installing crontab..."
echo "56 * * * * /usr/local/bin/load-wind" >> /tmp/crontab-profile
crontab -u postgres /tmp/crontab-profile
rm /tmp/crontab-profile
echo "Finished!"

