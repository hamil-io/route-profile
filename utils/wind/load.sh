#!/bin/bash

cd ~/projects/wind/
rm hrrr.t* 
URL="http://www.ftp.ncep.noaa.gov/data/nccf/com/hrrr/prod"
FILE=`date -u +"hrrr.%Y%m%d/hrrr.t%Hz.wrfsubhf01.grib2"`

wget "$URL/$FILE"

while [ $? -ne 0 ]
do
    echo "Retrying download in 10 seconds..."
    sleep 10
    wget "$URL/$FILE"
done

echo "Download sucessful!"
echo "Exporting raster bands..."

raster2pgsql -d -I -M -b 15,16,56,57,97,98,138,139 -s 98411 -t 10x10 hrrr.t* wind > wind.sql

echo "Importing into Postgres..."

sudo psql elevation < wind.sql

echo "Finished!"
