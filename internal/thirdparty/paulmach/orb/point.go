package orb

// A Point is a Lon/Lat 2d point.
type Point [2]float64

var _ Pointer = Point{}

// GeoJSONType returns the GeoJSON type for the object.
func (p Point) GeoJSONType() string {
	return "Point"
}

// Dimensions returns 0 because a point is a 0d object.
func (p Point) Dimensions() int {
	return 0
}

// Bound returns a single point bound of the point.
func (p Point) Bound() Bound {
	return Bound{p, p}
}

// Point returns itself so it implements the Pointer interface.
func (p Point) Point() Point {
	return p
}

// Y returns the vertical coordinate of the point.
func (p Point) Y() float64 {
	return p[1]
}

// X returns the horizontal coordinate of the point.
func (p Point) X() float64 {
	return p[0]
}

// Lat returns the vertical, latitude coordinate of the point.
func (p Point) Lat() float64 {
	return p[1]
}

// Lon returns the horizontal, longitude coordinate of the point.
func (p Point) Lon() float64 {
	return p[0]
}

// Equal checks if the point represents the same point or vector.
func (p Point) Equal(point Point) bool {
	return p[0] == point[0] && p[1] == point[1]
}

// GCJ02ToWGS84 GCJ02 to WGS84.
func (p Point) GCJ02ToWGS84() {
	p[0], p[1] = Gcj02ToWgs84(p.X(), p.Y())
}

// WGS84ToGCJ02  WGS84 to GCJ02.
func (p Point) WGS84ToGCJ02() {
	p[0], p[1] = Wgs84ToGcj02(p.X(), p.Y())
}

// BD09ToWGS84 BD09 to WGS84.
func (p Point) BD09ToWGS84() {
	p[0], p[1] = Bd09ToGcj02(p.X(), p.Y())
}

// WGS84ToBD09  WGS84 to BD09.
func (p Point) WGS84ToBD09() {
	p[0], p[1] = Wgs84ToBd09(p.X(), p.Y())
}
