package db

import (
	"math"
	"route-profile/geometry"
)

// SplitSegments takes a geometry and a goal and splits the geometry so that
// each segment has the goal length.
func SplitSegments(geom string, goal float64) []geometry.SubGeometry {
	var result []geometry.SubGeometry

	// Max error in meters
	err := 1.0

	length := 0.0
	subGeom := Geometry(geom)
	subGeom.StartPosition = 0
	subGeom.EndPosition = 1
	out := make(chan geometry.SubGeometry, 32)
	Split(subGeom, goal, out)

	for segment := range out {
		result = append(result, segment)
		length += segment.Length
		if math.Abs(length-subGeom.Length) < err {
			close(out)
		}
	}

	return result
}

// Split will recursively split a SubGeometry until the goal length is reached.
func Split(geom geometry.SubGeometry, goal float64, out chan geometry.SubGeometry) <-chan geometry.SubGeometry {
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
