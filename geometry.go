package sketchy

import (
	"fmt"
	"github.com/fogleman/gg"
	"log"
	"math"
)

// Point is a simple point in 2D space
type Point struct {
	X float64
	Y float64
}

// Line is two points that form a line
type Line struct {
	P Point
	Q Point
}

// Curve A curve is a list of points, may be closed
type Curve struct {
	Points []Point
	Closed bool
}

// Rect is a simple rectangle
type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

// Tuple representation of a point, useful for debugging
func (p Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// Lerp is a linear interpolation between two points
func (p Point) Lerp(a Point, i float64) Point {
	return Point{
		X: Lerp(p.X, a.X, i),
		Y: Lerp(p.Y, a.Y, i),
	}
}

// String representation of a line, useful for debugging
func (l Line) String() string {
	return fmt.Sprintf("(%f, %f) -> (%f, %f)", l.P.X, l.P.Y, l.Q.X, l.Q.Y)
}

func (l Line) Angle() float64 {
	dy := l.Q.Y - l.P.Y
	dx := l.Q.X - l.P.X
	angle := math.Atan(dy / dx)
	return angle
}

// Slope computes the slope of the line
func (l Line) Slope() float64 {
	dy := l.Q.Y - l.P.Y
	dx := l.Q.X - l.P.X
	if math.Abs(dx) < Smol {
		if dx < 0 {
			if dy > 0 {
				return math.Inf(-1)
			} else {
				return math.Inf(1)
			}
		} else {
			if dy > 0 {
				return math.Inf(1)
			} else {
				return math.Inf(-1)
			}
		}
	}
	return dy / dx
}

func (l Line) InvertedSlope() float64 {
	slope := l.Slope()
	if math.IsInf(slope, 1) || math.IsInf(slope, -1) {
		return 0
	}
	return -1 / slope
}

func (l Line) PerpendicularAt(percentage float64, length float64) Line {
	angle := l.Angle()
	point := l.P.Lerp(l.Q, percentage)
	sinOffset := 0.5 * length * math.Sin(angle)
	cosOffset := 0.5 * length * math.Cos(angle)
	p := Point{
		X: NoTinyVals(point.X - sinOffset),
		Y: NoTinyVals(point.Y + cosOffset),
	}
	q := Point{
		X: NoTinyVals(point.X + sinOffset),
		Y: NoTinyVals(point.Y - cosOffset),
	}
	return Line{
		P: p,
		Q: q,
	}
}

func (l Line) PerpendicularBisector(length float64) Line {
	return l.PerpendicularAt(0.5, length)
}

// Lerp is an interpolation between the two points of a line
func (l Line) Lerp(i float64) Point {
	return Point{
		X: Lerp(l.P.X, l.Q.X, i),
		Y: Lerp(l.P.Y, l.Q.Y, i),
	}
}

func (l Line) Draw(ctx *gg.Context) {
	ctx.DrawLine(l.P.X, l.P.Y, l.Q.X, l.Q.Y)
}

// Distance between two points
func Distance(p Point, q Point) float64 {
	return math.Sqrt(math.Pow(q.X-p.X, 2) + math.Pow(q.Y-p.Y, 2))
}

// SquaredDistance is the square of the distance between two points
func SquaredDistance(p Point, q Point) float64 {
	return math.Pow(q.X-p.X, 2) + math.Pow(q.Y-p.Y, 2)
}

// Midpoint Calculates the midpoint between two points
func Midpoint(p Point, q Point) Point {
	return Point{X: 0.5 * (p.X + q.X), Y: 0.5 * (p.Y + q.Y)}
}

// Midpoint Calculates the midpoint of a line
func (l Line) Midpoint() Point {
	return Midpoint(l.P, l.Q)
}

// Length Calculates the length of a line
func (l Line) Length() float64 {
	return Distance(l.P, l.Q)
}

// Length Calculates the length of the line segments of a curve
func (c *Curve) Length() float64 {
	result := 0.0
	n := len(c.Points)
	for i := 0; i < n-1; i++ {
		result += Distance(c.Points[i], c.Points[i+1])
	}
	if c.Closed {
		result += Distance(c.Points[0], c.Points[n-1])
	}
	return result
}

// Last returns the last point in a curve
func (c *Curve) Last() Point {
	n := len(c.Points)
	switch n {
	case 0:
		return Point{
			X: 0,
			Y: 0,
		}
	case 1:
		return c.Points[0]
	}
	if c.Closed {
		return c.Points[0]
	}
	return c.Points[n-1]
}

func (c *Curve) LastLine() Line {
	n := len(c.Points)
	switch n {
	case 0:
		return Line{
			P: Point{X: 0, Y: 0},
			Q: Point{X: 0, Y: 0},
		}
	case 1:
		return Line{
			P: c.Points[0],
			Q: c.Points[0],
		}
	}
	if c.Closed {
		return Line{
			P: c.Points[n-1],
			Q: c.Points[0],
		}
	}
	return Line{
		P: c.Points[n-2],
		Q: c.Points[n-1],
	}
}

// Lerp calculates a point a given percentage along a curve
func (c *Curve) Lerp(percentage float64) Point {
	var point Point
	if percentage < 0 || percentage > 1 {
		log.Fatalf("percentage in Lerp not between 0 and 1: %v\n", percentage)
	}
	if NoTinyVals(percentage) == 0 {
		return c.Points[0]
	}
	if math.Abs(percentage-1) < Smol {
		return c.Last()
	}
	totalDist := c.Length()
	targetDist := percentage * totalDist
	partialDist := 0.0
	var foundPoint bool
	n := len(c.Points)
	for i := 0; i < n-1; i++ {
		dist := Distance(c.Points[i], c.Points[i+1])
		if partialDist+dist >= targetDist {
			remainderDist := targetDist - partialDist
			pct := remainderDist / dist
			point = c.Points[i].Lerp(c.Points[i+1], pct)
			foundPoint = true
			break
		}
		partialDist += dist
	}
	if !foundPoint {
		if c.Closed {
			dist := Distance(c.Points[n-1], c.Points[0])
			remainderDist := targetDist - partialDist
			pct := remainderDist / dist
			point = c.Points[n-1].Lerp(c.Points[0], pct)
		} else {
			panic("couldn't find curve lerp point")
		}
	}
	return point
}

func (c *Curve) LineAt(percentage float64) (Line, float64) {
	var line Line
	var linePct float64
	if percentage < 0 || percentage > 1 {
		log.Fatalf("percentage in Lerp not between 0 and 1: %v\n", percentage)
	}
	if NoTinyVals(percentage) == 0 {
		return Line{P: c.Points[0], Q: c.Points[1]}, 0
	}
	if math.Abs(percentage-1) < Smol {
		return c.LastLine(), 1
	}
	totalDist := c.Length()
	targetDist := percentage * totalDist
	partialDist := 0.0
	var foundPoint bool
	n := len(c.Points)
	for i := 0; i < n-1; i++ {
		dist := Distance(c.Points[i], c.Points[i+1])
		if partialDist+dist >= targetDist {
			remainderDist := targetDist - partialDist
			linePct = remainderDist / dist
			line.P = c.Points[i]
			line.Q = c.Points[i+1]
			foundPoint = true
			break
		}
		partialDist += dist
	}
	if !foundPoint {
		if c.Closed {
			dist := Distance(c.Points[n-1], c.Points[0])
			remainderDist := targetDist - partialDist
			linePct = remainderDist / dist
			line.P = c.Points[n-1]
			line.Q = c.Points[0]
		} else {
			panic("couldn't find curve lerp point")
		}
	}
	return line, linePct
}

func (c *Curve) PerpendicularAt(percentage float64, length float64) Line {
	line, linePct := c.LineAt(percentage)
	return line.PerpendicularAt(linePct, length)
}

func (c *Curve) Draw(ctx *gg.Context) {
	n := len(c.Points)
	for i := 0; i < n-1; i++ {
		ctx.DrawLine(c.Points[i].X, c.Points[i].Y, c.Points[i+1].X, c.Points[i+1].Y)
	}
	if c.Closed {
		ctx.DrawLine(c.Points[n-1].X, c.Points[n-1].Y, c.Points[0].X, c.Points[0].Y)
	}
}

// ContainsPoint determines if a point lies within a rectangle
func (r *Rect) ContainsPoint(p Point) bool {
	return p.X >= r.X && p.X <= r.X+r.W && p.Y >= r.Y && p.Y <= r.Y+r.H
}
