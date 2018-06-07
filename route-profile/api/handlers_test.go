package api

import (
	"math"
	"route-profile/db"
	"testing"
)

var geometries = []struct {
	Geometry string
	Length   float64
}{
	{"{``mE|h}aOgFEApECtE?p@AbDAfD?\\At@A~DAxK?d@A\\Cl@?nD?P?L?LB`@f@tFPjC@P@Lx@dJ@F?N?VCdJ?R?VChI?P?FeECy@?}FEqFCM?E?CpF?v@ChHBlHAhG@^D?nB@?P@PONWAGF", 0.0293336298528689},
}

func TestGetSegments(t *testing.T) {
}

func TestGetProfile(t *testing.T) {
}

func TestGetRaster(t *testing.T) {
}

func TestGetPolyline(t *testing.T) {
	for _, route := range geometries {
		polyline := GetPolyline(route.Geometry, 0.00005)
		length := db.GeometryLength(polyline)
		if math.Abs(length-route.Length)/route.Length > 0.01 {
			t.Errorf("Length mismatch: computed length is %f, expected %f", length, route.Length)
		}
	}
}
