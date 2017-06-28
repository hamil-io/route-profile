package api

import "net/http"

// Route is used by our router to assign handlers
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is a slice of all available routes
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Segments",
		"POST",
		"/segments",
		GetSegments,
	},
	Route{
		"Wind",
		"POST",
		"/wind",
		GetWind,
	},
	Route{
		"Elevation",
		"POST",
		"/elevation",
		GetElevation,
	},
}
