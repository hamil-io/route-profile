package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"route-profile/db"
	"route-profile/geometry"
	"sort"
)

// Index is a stub handler and used only to verify that the service is running
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// GetElevation takes an API request and returns the elevation along the
// supplied geometry.
func GetElevation(w http.ResponseWriter, r *http.Request) {
	GetRaster("elevation", w, r)
}

// GetWind takes a request and returns the headwind in m/s along the
// supplied geometry.
func GetWind(w http.ResponseWriter, r *http.Request) {
	GetRaster("wind", w, r)
}

// GetSegments takes a geometry and a resolution and returns an array
// of encoded polylines split from the specified geometry at the specified
// resolution.
func GetSegments(w http.ResponseWriter, r *http.Request) {
	var response []string
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	geom := req.Geometry
	reqRes := req.Resolution
	goal := 4000.0
	if reqRes != 0.0 {
		goal = reqRes
	}

	segments := db.SplitSegments(geom, goal)

	for _, segment := range segments {
		response = append(response, segment.Geometry)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// GetRaster is used by GetWind and GetElevations, it handles splitting the
// request geometry and parallelizing the underlying raster calculations.
func GetRaster(raster string, w http.ResponseWriter, r *http.Request) {
	var result []float64
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	geom := req.Geometry
	reqRes := req.Resolution
	length := db.GeometryLength(geom)
	resolution := length / 100.0

	// Prevent oversampling
	if resolution < 0.00025 {
		resolution = 0.00025
	}

	if reqRes != 0.0 {
		resolution = reqRes
	}

	// The number of samples we want to process per thread
	samples := 1000.0
	goal := samples / (length / resolution)
	goal *= db.GeographyLength(geom)

	rasterFunc := func(raster string, geom geometry.SubGeometry, res float64, out chan geometry.RasterSegment) <-chan geometry.RasterSegment {
		if raster == "wind" {
			out <- db.Wind(geom, resolution)
		} else if raster == "elevation" {
			out <- db.Elevation(geom, resolution)
		}
		return out
	}

	segments := db.SplitSegments(geom, goal)
	total := len(segments)
	out := make(chan geometry.RasterSegment)
	var rasterSegments geometry.RasterSegments

	for _, segment := range segments {
		go rasterFunc(raster, segment, resolution, out)
	}

	count := 0
	for segment := range out {
		rasterSegments = append(rasterSegments, segment)
		count++
		if count == total {
			close(out)
		}
	}

	sort.Sort(rasterSegments)

	for _, segment := range rasterSegments {
		result = append(result, segment.Values...)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}
