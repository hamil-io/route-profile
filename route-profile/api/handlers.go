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
	GetProfile("elevation", w, r)
}

// GetWind takes a request and returns the headwind in m/s along the
// supplied geometry.
func GetWind(w http.ResponseWriter, r *http.Request) {
	GetProfile("wind", w, r)
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

// GetProfile is used by GetWind and GetElevations, it handles splitting the
// request geometry and parallelizing the underlying raster calculations.
func GetProfile(raster string, w http.ResponseWriter, r *http.Request) {
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	geom := req.Geometry
	reqRes := req.Resolution
	format := req.Format
	length := db.GeometryLength(geom)
	resolution := length / 100.0

	if format == "" {
		format = "array"
	}

	fmt.Printf("Format: %v\n", format)

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

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if format == "array" {
		result := GetRaster(raster, geom, resolution, goal)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
	} else {
		result := GetPolyline(geom, resolution)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
	}
}

// GetPolyline returns a 3d polyline for a particular geometry
func GetPolyline(geom string, resolution float64) string {
	length := db.GeometryLength(geom)
	samples := 1000.0
	goal := samples / (length / resolution)
	goal *= db.GeographyLength(geom)

	rasterFunc := func(geom geometry.SubGeometry, res float64, out chan geometry.SubGeometry) <-chan geometry.SubGeometry {
		out <- db.ElevationGeometry(geom, resolution)
		return out
	}

	segments := db.SplitSegments(geom, goal)

	total := len(segments)
	out := make(chan geometry.SubGeometry)
	var geometrySegments geometry.SubGeometries

	for _, segment := range segments {
		go rasterFunc(segment, resolution, out)
	}

	count := 0
	for segment := range out {
		geometrySegments = append(geometrySegments, segment)
		count++
		if count == total {
			close(out)
		}
	}

	sort.Sort(geometrySegments)

	line := db.Combine(geometrySegments)
	result := db.EncodePolyline(line.Geometry)
	return result
}

// GetRaster returns the raster values for a particular raster layer and polyline
func GetRaster(raster string, geom string, resolution float64, goal float64) []float64 {
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

	var result []float64
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

	return result
}
