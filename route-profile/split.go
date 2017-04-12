package main

import (
    "math"
)

type SubGeometry struct {
	Geometry         string
	StartPosition    float64
    Length           float64
}

type RasterSegment struct {
    Values         []float64
    StartPosition  float64
}

type RasterSegments []RasterSegment

func (slice RasterSegments) Len() int {
    return len(slice)
}

func (slice RasterSegments) Less(i, j int) bool {
    return slice[i].StartPosition < slice[j].StartPosition;
}

func (slice RasterSegments) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}

func SplitSegments(geometry string, goal float64) []SubGeometry {
    var result []SubGeometry

    // Max error in meters
    err := 1.0

    length := 0.0
    geom := Geometry(geometry)
    out := make(chan SubGeometry, 32)
    Split(geom, goal, out)

    for segment := range(out) {
        result = append(result, segment)
        length += segment.Length
        if (math.Abs(length - geom.Length) < err) {
            close(out)
        }
    }

    return result
}

func Split(geom SubGeometry, goal float64, out chan SubGeometry) <-chan SubGeometry {
    // Max split per iteration
    max := 10

    if (geom.Length < goal) {
        out <- geom
    } else {
        split := int(math.Ceil(geom.Length/goal))
        if (split > max) {
            split = max
        }

        for i, segment := range(Segments(geom, split)) {
            start := geom.StartPosition
            segment.StartPosition = start +  (1.0 - start) * (float64(i) * (1.0/float64(split)))
            go Split(segment, goal, out)
        }
    }

    return out
}
