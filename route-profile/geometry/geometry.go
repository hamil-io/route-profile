package geometry

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
