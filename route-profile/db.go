package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB

func init() {
	var err error
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		host = "/var/run/postgresql/"
	}
	db, err = sql.Open("postgres", "user="+os.Getenv("DB_USER")+" "+
		"host="+host+" "+
		"dbname="+os.Getenv("DB_NAME"))
	db.SetMaxOpenConns(32)

	if err != nil {
		panic(err)
	}
}

// Wind takes a SubGeometry and returns the headwind in m/s interpolated along
// that geometry at the specified resolution. Resolution is in Degrees.
func Wind(geom SubGeometry, resolution float64) RasterSegment {
	var headwind float64
	var result []float64
	rows, err := db.Query("SELECT headwind FROM wind($1, $2)", geom.Geometry, resolution)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&headwind)
		if err != nil {
			panic(err)
		}
		result = append(result, headwind)
	}

	return RasterSegment{result, geom.StartPosition}
}

// Elevation takes a SubGeometry and returns the elevation in meters interpolated
//along that geometry at the specified resolution. Resolution is in Degrees.
func Elevation(geom SubGeometry, resolution float64) RasterSegment {
	var altitude float64
	var result []float64
	rows, err := db.Query("select z from profile(drape($1, 'elevation', $2::numeric))", geom.Geometry, resolution)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&altitude)
		if err != nil {
			panic(err)
		}
		result = append(result, altitude)
	}

	return RasterSegment{result, geom.StartPosition}
}

// Geometry takes an encoded polyline and returns a SubGeometry consisting of
// the geometry, the start/end points of the geometry, and the length.
func Geometry(geometry string) SubGeometry {
	var geom string
	var length float64
	rows, err := db.Query("select geom, ST_Length(geom::geography) as length from "+
		"ST_LineFromEncodedPolyline($1) geom", geometry)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&geom, &length)
		if err != nil {
			panic(err)
		}
	}

	return SubGeometry{Geometry: geom, Length: length}
}

// Segments takes a SubGeometry and splits it into n SubGeometries, where n is
// the number of pices specified.
func Segments(geom SubGeometry, pieces int) []SubGeometry {
	var segment string
	var length float64
	var result []SubGeometry
	rows, err := db.Query("select geom, length from split_line($1, $2)", geom.Geometry, 1.0/float64(pieces))

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&segment, &length)
		if err != nil {
			panic(err)
		}
		result = append(result, SubGeometry{Geometry: segment, Length: length})
	}

	return result
}

// GeometryLength returns the length of the specified encoded polyline in the
// units of the underlying spatial reference system. By default this will be in
// degrees since we use EPSG 4326 throughout.
func GeometryLength(geometry string) float64 {
	var length float64
	rows, err := db.Query("SELECT ST_Length(ST_LineFromEncodedPolyline($1)) as length", geometry)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&length)
		if err != nil {
			panic(err)
		}
	}

	return length
}

// GeographyLength returns the length of the specified encode polyline in meters.
func GeographyLength(geometry string) float64 {
	var length float64
	rows, err := db.Query("SELECT ST_Length(ST_LineFromEncodedPolyline($1)::geography) as length", geometry)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&length)
		if err != nil {
			panic(err)
		}
	}

	return length
}
