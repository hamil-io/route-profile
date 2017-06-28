package geometry

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
