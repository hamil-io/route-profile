package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("postgres", "postgres://postgres:tribes@127.0.0.1/elevation")
    db.SetMaxOpenConns(32)

	if err != nil {
		panic(err)
	}
}

func Wind(geom SubGeometry, resolution float64) RasterSegment{
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

func Elevation(geom SubGeometry, resolution float64) RasterSegment{
    var altitude float64
    var result []float64
    rows, err := db.Query("select z from profile(drape($1, $2::numeric))", geom.Geometry, resolution)

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

func Geometry(geometry string) SubGeometry {
    var geom string
    var length float64
    rows, err := db.Query("select geom, ST_Length(geom::geography) as length from " +
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

func Segments(geom SubGeometry, pieces int) []SubGeometry{
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

func Length(geometry string) float64{
    var length float64
    rows, err := db.Query("SELECT ST_Length($1::geography) as length", geometry)

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
