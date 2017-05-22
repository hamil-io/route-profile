#!/bin/bash

# Test for getopt
getopt --test > /dev/null
if [ $? -ne 4 ]; then
    echo "I’m sorry, `getopt --test` failed in this environment."
    exit 1
fi

# Find source directory of this script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
    DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

show_help () {
    usage="load-elevation -- Program to load SRTMGL1 tiles. Requires a NASA Earthdata login.

Usage: load-elevation -u user -p password [OPTION...]

Where:

    -h                                 Show this help text
    -u, --user=USERNAME                Username for NASA Earthdata
    -p, --password=PASSWORD            Password for NASA Earthdata
    -c, --continent=CONTINENT          Continent you would like to load, available values are:
                                           Africa, Australia, Eurasia, North_America, South_America

    Constraints:
        --latitude-min=[-90, 90]       Minimum latitude you would like to load
        --latitude-max=[-90, 90]       Maximum latitude you would like to load
        --longitude-min=[-180, 180]    Minimum longitude you would like to load
        --longitude-max=[-180, 180]    Maximum longitude you would like to load"
        
    echo "$usage"
}

SHORT=cuph
LONG=latitude-min,latitude-max,longitude-min,longitude-max,continent,user,password

PARSED=$(getopt --options $SHORT --longoptions $LONG --name "$0" -- "$@")
if [ $? -ne 0 ]; then
    exit 2
fi

eval set -- "$PARSED"

while true; do
    case "$1" in
        -h|--help)
            show_help
            exit
            ;;
        -c|--continent)
            continent=$2
            shift
            ;;
        -u|--user)
            user=$2
            shift
            ;;
        -p|--password)
            pass=$2
            shift
            ;;
        --latitude-min)
            lat_min=$2
            shift
            ;;
        --latitude-max)
            lat_max=$2
            shift
            ;;
        --longitude-min)
            lon_min=$2
            shift
            ;;
        --longitude-max)
            lon_max=$2
            shift
            ;;
        --)
            shift
            break
            ;;
        *)
            echo "Programming error"
            exit 3
            ;;
    esac
done

if [ -z "$user" ]; then
    echo "Username must be set, use -h for help"
    exit 3;
fi

if [ -z "$pass" ]; then
    echo "Password must be set, use -h for help"
    exit 3;
fi

# initialize tmp directory for tiles
mkdir -p /tmp/route-profile/elevation
cd /tmp/route-profile/elevation

input="$DIR/tile.urls"

# test to see if user/pass is correct
test_tile=$(head -n 1 $input)
test_resp=0
test_resp=$(wget --user $user --password $pass $test_tile  2>&1 | grep -c "401")

if [ "$test_resp" -gt 0 ]; then
    echo "Username or password is incorrect"
    exit 3
fi

# download the specified tiles
files=()
while IFS= read -r var
do
    files+=($var)
done < "$input"
echo "$user"
echo "$pass"
echo ${files[*]} | xargs -n 1 -P 16 wget --user $user --password $pass
