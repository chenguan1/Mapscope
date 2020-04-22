package orb

// A MultiPoint represents a set of points in the 2D Eucledian or Cartesian plane.
type MultiPoint []Point

// GeoJSONType returns the GeoJSON type for the object.
func (mp MultiPoint) GeoJSONType() string {
	return "MultiPoint"
}

// Dimensions returns 0 because a MultiPoint is a 0d object.
func (mp MultiPoint) Dimensions() int {
	return 0
}

// Clone returns a new copy of the points.
func (mp MultiPoint) Clone() MultiPoint {
	if mp == nil {
		return nil
	}

	points := make([]Point, len(mp))
	copy(points, mp)

	return MultiPoint(points)
}

// Bound returns a bound around the points. Uses rectangular coordinates.
func (mp MultiPoint) Bound() Bound {
	if len(mp) == 0 {
		return emptyBound
	}

	b := Bound{mp[0], mp[0]}
	for _, p := range mp {
		b = b.Extend(p)
	}

	return b
}

// Equal compares two MultiPoint objects. Returns true if lengths are the same
// and all points are Equal, and in the same order.
func (mp MultiPoint) Equal(multiPoint MultiPoint) bool {
	if len(mp) != len(multiPoint) {
		return false
	}

	for i := range mp {
		if !mp[i].Equal(multiPoint[i]) {
			return false
		}
	}

	return true
}

// GCJ02ToWGS84 GCJ02 to WGS84.
func (mp MultiPoint) GCJ02ToWGS84() {
	if len(mp) == 0 {
		return
	}
	for i := 0; i < len(mp); i++ {
		mp[i][0], mp[i][1] = Gcj02ToWgs84(mp[i].X(), mp[i].Y())
	}
}

// WGS84ToGCJ02 WGS84 to GCJ02.
func (mp MultiPoint) WGS84ToGCJ02() {
	if len(mp) == 0 {
		return
	}
	for i := 0; i < len(mp); i++ {
		mp[i][0], mp[i][1] = Wgs84ToGcj02(mp[i].X(), mp[i].Y())
	}
}

// BD09ToWGS84 BD09 to WGS84.
func (mp MultiPoint) BD09ToWGS84() {
	if len(mp) == 0 {
		return
	}
	for i := 0; i < len(mp); i++ {
		mp[i][0], mp[i][1] = Bd09ToWgs84(mp[i].X(), mp[i].Y())
	}
}

// WGS84ToBD09 WGS84 to BD09.
func (mp MultiPoint) WGS84ToBD09() {
	if len(mp) == 0 {
		return
	}
	for i := 0; i < len(mp); i++ {
		mp[i][0], mp[i][1] = Wgs84ToBd09(mp[i].X(), mp[i].Y())
	}
}
