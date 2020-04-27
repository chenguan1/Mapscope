package tilecover

import (
	"log"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// LineString creates a tile cover for the line string.
func LineString(ls orb.LineString, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	line(set, ls, z, nil)
	return set
}

// LineStringCount creates a tile cover for the line string.
func LineStringCount(ls orb.LineString, z maptile.Zoom) int64 {
	_, cnt := lineCount(ls, z, nil)
	return cnt
}

// LineStringChannel creates a tile cover for the line string.
func LineStringChannel(ls orb.LineString, z maptile.Zoom, ch chan<- maptile.Tile) {
	lineChannel(ls, z, nil, ch)
}

// MultiLineString creates a tile cover for the line strings.
func MultiLineString(mls orb.MultiLineString, z maptile.Zoom) maptile.Set {
	set := make(maptile.Set)
	for _, ls := range mls {
		line(set, ls, z, nil)
	}

	return set
}

// MultiLineStringCount creates a tile cover for the line strings.
func MultiLineStringCount(mls orb.MultiLineString, z maptile.Zoom) int64 {
	set := make(maptile.Set)
	for _, ls := range mls {
		line(set, ls, z, nil)
	}

	var cnt int64
	for _, v := range set {
		if v {
			cnt++
		}
	}
	return cnt
}

// MultiLineStringChannel creates a tile cover for the line string.
func MultiLineStringChannel(mls orb.MultiLineString, z maptile.Zoom, ch chan<- maptile.Tile) {

	set := make(maptile.Set)
	for _, ls := range mls {
		line(set, ls, z, nil)
	}
	for t, v := range set {
		if v {
			ch <- t
		}
	}
}

func line(
	set maptile.Set,
	line orb.LineString,
	zoom maptile.Zoom,
	ring [][2]uint32,
) [][2]uint32 {
	inf := math.Inf(1)

	prevX := -1.0
	prevY := -1.0

	var x, y float64

	for i := 0; i < len(line)-1; i++ {
		start := maptile.Fraction(line[i], zoom)
		stop := maptile.Fraction(line[i+1], zoom)

		dx := stop[0] - start[0]
		dy := stop[1] - start[1]

		if dy == 0 && dx == 0 {
			continue
		}

		sx := -1.0
		if dx > 0 {
			sx = 1.0
		}
		sy := -1.0
		if dy > 0 {
			sy = 1.0
		}

		x = math.Floor(start[0])
		y = math.Floor(start[1])

		tMaxX := inf
		if dx != 0 {
			d := 0.0
			if dx > 0 {
				d = 1.0
			}
			tMaxX = math.Abs((d + x - start[0]) / dx)
		}

		tMaxY := inf
		if dy != 0 {
			d := 0.0
			if dy > 0 {
				d = 1.0
			}
			tMaxY = math.Abs((d + y - start[1]) / dy)
		}

		tdx := math.Abs(sx / dx)
		tdy := math.Abs(sy / dy)

		if x != prevX || y != prevY {
			set[maptile.New(uint32(x), uint32(y), zoom)] = true
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}

		for tMaxX < 1 || tMaxY < 1 {
			if tMaxX < tMaxY {
				tMaxX += tdx
				x += sx
			} else {
				tMaxY += tdy
				y += sy
			}

			set[maptile.New(uint32(x), uint32(y), zoom)] = true
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}
	}

	if ring != nil && uint32(y) == ring[0][1] {
		ring = ring[:len(ring)-1]
	}

	return ring
}

func lineCount(
	line orb.LineString,
	zoom maptile.Zoom,
	ring [][2]uint32,
) ([][2]uint32, int64) {

	var cnt int64
	inf := math.Inf(1)

	prevX := -1.0
	prevY := -1.0

	var x, y float64

	for i := 0; i < len(line)-1; i++ {
		start := maptile.Fraction(line[i], zoom)
		stop := maptile.Fraction(line[i+1], zoom)

		dx := stop[0] - start[0]
		dy := stop[1] - start[1]

		if dy == 0 && dx == 0 {
			continue
		}

		sx := -1.0
		if dx > 0 {
			sx = 1.0
		}
		sy := -1.0
		if dy > 0 {
			sy = 1.0
		}

		x = math.Floor(start[0])
		y = math.Floor(start[1])

		tMaxX := inf
		if dx != 0 {
			d := 0.0
			if dx > 0 {
				d = 1.0
			}
			tMaxX = math.Abs((d + x - start[0]) / dx)
		}

		tMaxY := inf
		if dy != 0 {
			d := 0.0
			if dy > 0 {
				d = 1.0
			}
			tMaxY = math.Abs((d + y - start[1]) / dy)
		}

		tdx := math.Abs(sx / dx)
		tdy := math.Abs(sy / dy)

		if x != prevX || y != prevY {
			// set[maptile.New(uint32(x), uint32(y), zoom)] = true
			cnt++
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}

		for tMaxX < 1 || tMaxY < 1 {
			if tMaxX < tMaxY {
				tMaxX += tdx
				x += sx
			} else {
				tMaxY += tdy
				y += sy
			}

			// set[maptile.New(uint32(x), uint32(y), zoom)] = true
			cnt++
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}
	}

	if ring != nil && uint32(y) == ring[0][1] {
		ring = ring[:len(ring)-1]
	}

	return ring, cnt
}

func lineChannel(
	line orb.LineString,
	zoom maptile.Zoom,
	ring [][2]uint32,
	ch chan<- maptile.Tile,
) [][2]uint32 {

	defer func() {
		if recover() != nil {
			log.Println("buffer got closed...")
		}
	}()

	inf := math.Inf(1)

	prevX := -1.0
	prevY := -1.0

	var x, y float64

	for i := 0; i < len(line)-1; i++ {
		start := maptile.Fraction(line[i], zoom)
		stop := maptile.Fraction(line[i+1], zoom)

		dx := stop[0] - start[0]
		dy := stop[1] - start[1]

		if dy == 0 && dx == 0 {
			continue
		}

		sx := -1.0
		if dx > 0 {
			sx = 1.0
		}
		sy := -1.0
		if dy > 0 {
			sy = 1.0
		}

		x = math.Floor(start[0])
		y = math.Floor(start[1])

		tMaxX := inf
		if dx != 0 {
			d := 0.0
			if dx > 0 {
				d = 1.0
			}
			tMaxX = math.Abs((d + x - start[0]) / dx)
		}

		tMaxY := inf
		if dy != 0 {
			d := 0.0
			if dy > 0 {
				d = 1.0
			}
			tMaxY = math.Abs((d + y - start[1]) / dy)
		}

		tdx := math.Abs(sx / dx)
		tdy := math.Abs(sy / dy)

		if x != prevX || y != prevY {
			// set[maptile.New(uint32(x), uint32(y), zoom)] = true
			ch <- maptile.New(uint32(x), uint32(y), zoom)
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}

		for tMaxX < 1 || tMaxY < 1 {
			if tMaxX < tMaxY {
				tMaxX += tdx
				x += sx
			} else {
				tMaxY += tdy
				y += sy
			}

			// set[maptile.New(uint32(x), uint32(y), zoom)] = true
			ch <- maptile.New(uint32(x), uint32(y), zoom)
			if ring != nil && y != prevY {
				ring = append(ring, [2]uint32{uint32(x), uint32(y)})
			}
			prevX = x
			prevY = y
		}
	}

	if ring != nil && uint32(y) == ring[0][1] {
		ring = ring[:len(ring)-1]
	}

	return ring
}
