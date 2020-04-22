package orb

// Polygon is a closed area. The first LineString is the outer ring.
// The others are the holes. Each LineString is expected to be closed
// ie. the first point matches the last.
type Polygon []Ring

// GeoJSONType returns the GeoJSON type for the object.
func (p Polygon) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Polygon is a 2d object.
func (p Polygon) Dimensions() int {
	return 2
}

// Bound returns a bound around the polygon.
func (p Polygon) Bound() Bound {
	if len(p) == 0 {
		return emptyBound
	}
	return p[0].Bound()
}

// Equal compares two polygons. Returns true if lengths are the same
// and all points are Equal.
func (p Polygon) Equal(polygon Polygon) bool {
	if len(p) != len(polygon) {
		return false
	}

	for i := range p {
		if !p[i].Equal(polygon[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the polygon.
// All of the rings are also cloned.
func (p Polygon) Clone() Polygon {
	if p == nil {
		return p
	}

	np := make(Polygon, 0, len(p))
	for _, r := range p {
		np = append(np, r.Clone())
	}

	return np
}

// GCJ02ToWGS84 GCJ02 to WGS84.
func (p Polygon) GCJ02ToWGS84() {
	if len(p) == 0 {
		return
	}
	for i := 0; i < len(p); i++ {
		p[i].GCJ02ToWGS84()
	}
}

// WGS84ToGCJ02 WGS84 to GCJ02.
func (p Polygon) WGS84ToGCJ02() {
	if len(p) == 0 {
		return
	}
	for i := 0; i < len(p); i++ {
		p[i].WGS84ToGCJ02()
	}
}

// BD09ToWGS84 BD09 to WGS84.
func (p Polygon) BD09ToWGS84() {
	if len(p) == 0 {
		return
	}
	for i := 0; i < len(p); i++ {
		p[i].BD09ToWGS84()
	}
}

// WGS84ToBD09 WGS84 to BD09.
func (p Polygon) WGS84ToBD09() {
	if len(p) == 0 {
		return
	}
	for i := 0; i < len(p); i++ {
		p[i].WGS84ToBD09()
	}
}
