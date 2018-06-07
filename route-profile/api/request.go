package api

// Request is used throughout by our API handlers. All incoming requests should
// match this struct.
type Request struct {
	Geometry   string  `json:"geometry"`
	Resolution float64 `json:"resolution"`
	Format     string  `json:"format"`
}
