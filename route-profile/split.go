package main

import (
	"math"
)

// SubGeometry is the main type used throughout. All calls to the db functions
// will use a SubGeometry. Subgeometry is simply a geometry along with its
// start/end points and its length. Storing this extra information makes it
// easier to parallelize our raster functions.
type SubGeometry struct {
	StartPosition float64
	EndPosition   float64
	Length        float64
	Geometry      string
}

// SubGeometries is a slice of SubGeometry structs.
// This is sortable by StartPosition.
type SubGeometries []SubGeometry

func (slice SubGeometries) Len() int {
	return len(slice)
}

func (slice SubGeometries) Less(i, j int) bool {
	return slice[i].StartPosition < slice[j].StartPosition
}

func (slice SubGeometries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// RasterSegment is the main data type returned by our db functions.
// RasterSegment has two parts. The start position, which is represented
// as a percentage of the incoming SubGeometry's length. And Values, which is
// an array of the calculated raster values for that SubGeometry.
type RasterSegment struct {
	Values        []float64
	StartPosition float64
}

// RasterSegments is a slice of RasterSegment structs.
// This is sortable by StartPosition.
type RasterSegments []RasterSegment

func (slice RasterSegments) Len() int {
	return len(slice)
}

func (slice RasterSegments) Less(i, j int) bool {
	return slice[i].StartPosition < slice[j].StartPosition
}

func (slice RasterSegments) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// SplitSegments takes a geometry and a goal and splits the geometry so that
// each split has the goal length.
func SplitSegments(geometry string, goal float64) []SubGeometry {
	var result []SubGeometry

	// Max error in meters
	err := 1.0

	length := 0.0
	geom := Geometry(geometry)
	geom.StartPosition = 0
	geom.EndPosition = 1
	out := make(chan SubGeometry, 32)
	Split(geom, goal, out)

	for segment := range out {
		result = append(result, segment)
		length += segment.Length
		if math.Abs(length-geom.Length) < err {
			close(out)
		}
	}

	return result
}

// Split will recursively split a SubGeometry until the goal length is reached.
func Split(geom SubGeometry, goal float64, out chan SubGeometry) <-chan SubGeometry {
	// Max split per iteration
	max := 10

	if geom.Length < goal {
		out <- geom
	} else {
		split := int(math.Ceil(geom.Length / goal))
		if split > max {
			split = max
		}

		for i, segment := range Segments(geom, split) {
			start := geom.StartPosition
			span := geom.EndPosition - start
			segment.StartPosition = start + (1.0-start)*(float64(i)*(span/float64(split)))
			segment.EndPosition = start + (1.0-start)*(float64(i+1)*(span/float64(split)))
			go Split(segment, goal, out)
		}
	}

	return out
}
