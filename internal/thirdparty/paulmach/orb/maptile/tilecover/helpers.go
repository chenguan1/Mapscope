// Package tilecover computes the covering set of tiles for an orb.Geometry.
package tilecover

import (
	"fmt"
	"log"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// Geometry returns the covering set of tiles for the given geometry.
func Geometry(g orb.Geometry, z maptile.Zoom) maptile.Set {
	if g == nil {
		return nil
	}

	switch g := g.(type) {
	case orb.Point:
		return Point(g, z)
	case orb.MultiPoint:
		return MultiPoint(g, z)
	case orb.LineString:
		return LineString(g, z)
	case orb.MultiLineString:
		return MultiLineString(g, z)
	case orb.Ring:
		return Ring(g, z)
	case orb.Polygon:
		return Polygon(g, z)
	case orb.MultiPolygon:
		return MultiPolygon(g, z)
	case orb.Collection:
		return Collection(g, z)
	case orb.Bound:
		return Bound(g, z)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// GeometryCount returns the covering set of tiles for the given geometry.
func GeometryCount(g orb.Geometry, z maptile.Zoom) int64 {
	if g == nil {
		return 0
	}

	switch g := g.(type) {
	case orb.Point:
		return int64(len(Point(g, z)))
	case orb.MultiPoint:
		return int64(len(MultiPoint(g, z)))
	case orb.LineString:
		return LineStringCount(g, z)
	case orb.MultiLineString:
		return MultiLineStringCount(g, z)
	case orb.Ring:
		return RingCount(g, z)
	case orb.Polygon:
		return PolygonCount(g, z)
	case orb.MultiPolygon:
		return MultiPolygonCount(g, z)
	case orb.Collection:
		return CollectionCount(g, z)
	case orb.Bound:
		return BoundCount(g, z)
	}

	panic(fmt.Sprintf("geometry type not supported: %T", g))
}

// GeometryChannel returns the covering set of tiles for the given geometry.
func GeometryChannel(g orb.Geometry, z maptile.Zoom, ch chan<- maptile.Tile) {
	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()
	if g == nil {
		return
	}

	switch g := g.(type) {
	case orb.Point:
		PointChannel(g, z, ch)
	case orb.MultiPoint:
		MultiPointChannel(g, z, ch)
	case orb.LineString:
		LineStringChannel(g, z, ch)
	case orb.MultiLineString:
		MultiLineStringChannel(g, z, ch)
	case orb.Ring:
		RingChannel(g, z, ch)
	case orb.Polygon:
		PolygonChannel(g, z, ch)
	case orb.MultiPolygon:
		MultiPolygonChannel(g, z, ch)
	case orb.Collection:
		CollectionChannel(g, z, ch)
	case orb.Bound:
		BoundChannel(g, z, ch)
	}
}

// Point creates a tile cover for the point, i.e. just the tile
// containing the point.
func Point(ll orb.Point, z maptile.Zoom) maptile.Set {
	return maptile.Set{
		maptile.At(ll, z): true,
	}
}

// PointChannel creates a tile cover for the point, i.e. just the tile
// containing the point.
func PointChannel(ll orb.Point, z maptile.Zoom, ch chan<- maptile.Tile) {
	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()
	ch <- maptile.At(ll, z)
}

// MultiPoint creates a tile cover for the set of points,
func MultiPoint(mp orb.MultiPoint, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, p := range mp {
		set[maptile.At(p, z)] = true
	}

	return set
}

// MultiPointChannel creates a tile cover for the point, i.e. just the tile
// containing the point.
func MultiPointChannel(mp orb.MultiPoint, z maptile.Zoom, ch chan<- maptile.Tile) {
	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()

	for _, p := range mp {
		ch <- maptile.At(p, z)
	}
}

// Bound creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func Bound(b orb.Bound, z maptile.Zoom) maptile.Set {
	lo := maptile.At(b.Min, z)
	hi := maptile.At(b.Max, z)

	result := make(maptile.Set, (hi.X-lo.X+1)*(lo.Y-hi.Y+1))

	for x := lo.X; x <= hi.X; x++ {
		for y := hi.Y; y <= lo.Y; y++ {
			result[maptile.Tile{X: x, Y: y, Z: z}] = true
		}
	}

	return result
}

// BoundCount creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func BoundCount(b orb.Bound, z maptile.Zoom) int64 {
	lo := maptile.At(b.Min, z)
	hi := maptile.At(b.Max, z)

	cnt := int64(hi.X-lo.X+1) * int64(lo.Y-hi.Y+1)

	if cnt == 0 {
		return 1
	}
	return cnt
}

// BoundChannel creates a tile cover for the bound. i.e. all the tiles
// that intersect the bound.
func BoundChannel(b orb.Bound, z maptile.Zoom, ch chan<- maptile.Tile) {
	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()

	lo := maptile.At(b.Min, z)
	hi := maptile.At(b.Max, z)
	for x := lo.X; x <= hi.X; x++ {
		for y := hi.Y; y <= lo.Y; y++ {
			ch <- maptile.Tile{X: x, Y: y, Z: z}
		}
	}
}

// Collection returns the covering set of tiles for the
// geoemtry collection.
func Collection(c orb.Collection, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, g := range c {
		set.Merge(Geometry(g, z))
	}

	return set
}

// CollectionCount returns the covering set of tiles for the
// geoemtry collection.
func CollectionCount(c orb.Collection, z maptile.Zoom) int64 {
	var cnt int64
	for _, g := range c {
		cnt += GeometryCount(g, z)
	}
	return cnt
}

// CollectionChannel returns the covering set of tiles for the
// geoemtry collection.
func CollectionChannel(c orb.Collection, z maptile.Zoom, ch chan<- maptile.Tile) {
	defer close(ch)
	for _, g := range c {
		GeometryChannel(g, z, ch)
	}
}
