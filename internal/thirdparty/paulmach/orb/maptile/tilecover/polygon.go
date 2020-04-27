package tilecover

import (
	"log"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// Ring creates a tile cover for the ring.
func Ring(r orb.Ring, z maptile.Zoom) maptile.Set {
	if len(r) == 0 {
		return make(maptile.Set)
	}
	return Polygon(orb.Polygon{r}, z)
}

// RingCount creates a tile cover for the ring.
func RingCount(r orb.Ring, z maptile.Zoom) int64 {
	if len(r) == 0 {
		return 0
	}
	return PolygonCount(orb.Polygon{r}, z)
}

// RingChannel creates a tile cover for the ring.
func RingChannel(r orb.Ring, z maptile.Zoom, ch chan<- maptile.Tile) {
	if len(r) == 0 {
		return
	}
	PolygonChannel(orb.Polygon{r}, z, ch)
}

// Polygon creates a tile cover for the polygon.
func Polygon(p orb.Polygon, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	polygon(set, p, z)

	return set
}

// PolygonCount creates a tile cover for the polygon.
func PolygonCount(p orb.Polygon, z maptile.Zoom) int64 {
	return polygonCount(p, z)
}

// PolygonChannel creates a tile cover for the polygon.
func PolygonChannel(p orb.Polygon, z maptile.Zoom, ch chan<- maptile.Tile) {
	polygonChannel(p, z, ch)
}

// MultiPolygon creates a tile cover for the multi-polygon.
func MultiPolygon(mp orb.MultiPolygon, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, p := range mp {
		polygon(set, p, z)
	}
	return set
}

// MultiPolygonCount creates a tile cover for the multi-polygon.
func MultiPolygonCount(mp orb.MultiPolygon, z maptile.Zoom) int64 {
	var cnt int64
	for _, p := range mp {
		cnt += polygonCount(p, z)
	}
	return cnt
}

// MultiPolygonChannel creates a tile cover for the multi-polygon.
func MultiPolygonChannel(mp orb.MultiPolygon, z maptile.Zoom, ch chan<- maptile.Tile) {
	for _, p := range mp {
		polygonChannel(p, z, ch)
	}
}

func polygon(set maptile.Set, p orb.Polygon, zoom maptile.Zoom) {
	intersections := make([][2]uint32, 0)

	for _, r := range p {
		ring := line(set, orb.LineString(r), zoom, make([][2]uint32, 0))

		pi := len(ring) - 2
		for i := range ring {
			pi = (pi + 1) % len(ring)
			ni := (i + 1) % len(ring)
			y := ring[i][1]

			// add interesction if it's not local extremum or duplicate
			if (y > ring[pi][1] || y > ring[ni][1]) && // not local minimum
				(y < ring[pi][1] || y < ring[ni][1]) && // not local maximum
				y != ring[ni][1] {

				intersections = append(intersections, ring[i])
			}
		}
	}

	// sort by y, then x
	sort.Slice(intersections, func(i, j int) bool {
		it := intersections[i]
		jt := intersections[j]

		if it[1] != jt[1] {
			return it[1] < jt[1]
		}

		return it[0] < jt[0]
	})

	for i := 0; i < len(intersections); i += 2 {
		// fill tiles between pairs of intersections
		y := intersections[i][1]
		for x := intersections[i][0] + 1; x < intersections[i+1][0]; x++ {
			set[maptile.New(x, y, zoom)] = true
		}
	}
}

func polygonCount(p orb.Polygon, zoom maptile.Zoom) int64 {
	set := make(maptile.Set)
	var cnt int64
	intersections := make([][2]uint32, 0)
	for _, r := range p {
		ring := line(set, orb.LineString(r), zoom, make([][2]uint32, 0))
		// ring, n := lineCount(orb.LineString(r), zoom, make([][2]uint32, 0))
		// cnt += n
		pi := len(ring) - 2
		for i := range ring {
			pi = (pi + 1) % len(ring)
			ni := (i + 1) % len(ring)
			y := ring[i][1]

			// add interesction if it's not local extremum or duplicate
			if (y > ring[pi][1] || y > ring[ni][1]) && // not local minimum
				(y < ring[pi][1] || y < ring[ni][1]) && // not local maximum
				y != ring[ni][1] {

				intersections = append(intersections, ring[i])
			}
		}
	}

	// sort by y, then x
	sort.Slice(intersections, func(i, j int) bool {
		it := intersections[i]
		jt := intersections[j]

		if it[1] != jt[1] {
			return it[1] < jt[1]
		}

		return it[0] < jt[0]
	})

	for _, v := range set {
		if v {
			cnt++
		}
	}

	for i := 0; i < len(intersections); i += 2 {
		// fill tiles between pairs of intersections
		// y := intersections[i][1]
		for x := intersections[i][0] + 1; x < intersections[i+1][0]; x++ {
			// set[maptile.New(x, y, zoom)] = true
			cnt++
		}
	}
	return cnt
}

func polygonChannel(p orb.Polygon, zoom maptile.Zoom, ch chan<- maptile.Tile) {
	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()
	set := make(maptile.Set)
	intersections := make([][2]uint32, 0)
	for _, r := range p {
		ring := line(set, orb.LineString(r), zoom, make([][2]uint32, 0))
		// ring := lineChannel(orb.LineString(r), zoom, make([][2]uint32, 0), ch)
		pi := len(ring) - 2
		for i := range ring {
			pi = (pi + 1) % len(ring)
			ni := (i + 1) % len(ring)
			y := ring[i][1]

			// add interesction if it's not local extremum or duplicate
			if (y > ring[pi][1] || y > ring[ni][1]) && // not local minimum
				(y < ring[pi][1] || y < ring[ni][1]) && // not local maximum
				y != ring[ni][1] {

				intersections = append(intersections, ring[i])
			}
		}
	}

	// sort by y, then x
	sort.Slice(intersections, func(i, j int) bool {
		it := intersections[i]
		jt := intersections[j]

		if it[1] != jt[1] {
			return it[1] < jt[1]
		}

		return it[0] < jt[0]
	})

	for t, v := range set {
		if v {
			ch <- t
		}
	}

	for i := 0; i < len(intersections); i += 2 {
		// fill tiles between pairs of intersections
		y := intersections[i][1]
		for x := intersections[i][0] + 1; x < intersections[i+1][0]; x++ {
			// set[maptile.New(x, y, zoom)] = true
			ch <- maptile.New(x, y, zoom)
		}
	}
}
