package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

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
