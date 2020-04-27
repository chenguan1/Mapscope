package orb

// MultiLineString is a set of polylines.
type MultiLineString []LineString

// GeoJSONType returns the GeoJSON type for the object.
func (mls MultiLineString) GeoJSONType() string {
	return "MultiLineString"
}

// Dimensions returns 1 because a MultiLineString is a 2d object.
func (mls MultiLineString) Dimensions() int {
	return 1
}

// Bound returns a bound around all the line strings.
func (mls MultiLineString) Bound() Bound {
	if len(mls) == 0 {
		return emptyBound
	}

	bound := mls[0].Bound()
	for i := 1; i < len(mls); i++ {
		bound = bound.Union(mls[i].Bound())
	}

	return bound
}

// Equal compares two multi line strings. Returns true if lengths are the same
// and all points are Equal.
func (mls MultiLineString) Equal(multiLineString MultiLineString) bool {
	if len(mls) != len(multiLineString) {
		return false
	}

	for i, ls := range mls {
		if !ls.Equal(multiLineString[i]) {
			return false
		}
	}

	return true
}

// Clone returns a new deep copy of the multi line string.
func (mls MultiLineString) Clone() MultiLineString {
	if mls == nil {
		return nil
	}

	nmls := make(MultiLineString, 0, len(mls))
	for _, ls := range mls {
		nmls = append(nmls, ls.Clone())
	}

	return nmls
}

// GCJ02ToWGS84 GCJ02 to WGS84.
func (mls MultiLineString) GCJ02ToWGS84() {
	if len(mls) == 0 {
		return
	}

	for i := 0; i < len(mls); i++ {
		mls[i].GCJ02ToWGS84()
	}
}

// WGS84ToGCJ02 WGS84 to GCJ02.
func (mls MultiLineString) WGS84ToGCJ02() {
	if len(mls) == 0 {
		return
	}

	for i := 0; i < len(mls); i++ {
		mls[i].WGS84ToGCJ02()
	}
}

// BD09ToWGS84 BD09 to WGS84.
func (mls MultiLineString) BD09ToWGS84() {
	if len(mls) == 0 {
		return
	}

	for i := 0; i < len(mls); i++ {
		mls[i].BD09ToWGS84()
	}
}

// WGS84ToBD09 WGS84 to BD09.
func (mls MultiLineString) WGS84ToBD09() {
	if len(mls) == 0 {
		return
	}

	for i := 0; i < len(mls); i++ {
		mls[i].WGS84ToBD09()
	}
}
