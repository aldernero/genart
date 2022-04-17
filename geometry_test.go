package sketchy

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestPoint_IsEqual(t *testing.T) {
	assert := assert.New(t)
	p := Point{
		X: 100.123456789,
		Y: 100.123456789,
	}
	q := Point{
		X: 100.123456789,
		Y: 100.123456789,
	}
	assert.Equal(p.IsEqual(q), true)
}

func TestSlope(t *testing.T) {
	assert := assert.New(t)
	line := Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 0, Y: 1},
	}
	assert.Equal(math.Inf(1), line.Slope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 0, Y: -1},
	}
	assert.Equal(math.Inf(-1), line.Slope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 0.0000001, Y: 1},
	}
	assert.Equal(math.Inf(1), line.Slope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: -0.0000001, Y: 1},
	}
	assert.Equal(math.Inf(-1), line.Slope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 1, Y: 2},
	}
	assert.Equal(float64(2), line.Slope())
}

func TestInvertedSlope(t *testing.T) {
	assert := assert.New(t)
	line := Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 0, Y: 1},
	}
	assert.Equal(float64(0), line.InvertedSlope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 0, Y: -1},
	}
	assert.Equal(float64(0), line.InvertedSlope())
	line = Line{
		P: Point{X: 0, Y: 0},
		Q: Point{X: 1, Y: 2},
	}
	assert.Equal(-0.5, line.InvertedSlope())
}

func TestPerpendicularBisector(t *testing.T) {
	assert := assert.New(t)
	line := Line{
		P: Point{X: -1, Y: 0},
		Q: Point{X: 1, Y: 0},
	}
	pb := line.PerpendicularBisector(2)
	assert.Equal(
		Line{P: Point{X: 0, Y: 1}, Q: Point{X: 0, Y: -1}},
		pb,
	)
	line = Line{
		P: Point{X: 0, Y: -1},
		Q: Point{X: 0, Y: 1},
	}
	pb = line.PerpendicularBisector(2)
	assert.Equal(
		Line{P: Point{X: -1, Y: 0}, Q: Point{X: 1, Y: 0}},
		pb,
	)
}

func TestCurveLerp(t *testing.T) {
	assert := assert.New(t)
	var curve Curve
	points := []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 10, Y: 0},
	}
	curve.Points = points
	assert.Equal(Point{X: 0, Y: 0}, curve.Lerp(0.0))
	assert.Equal(Point{X: 10, Y: 0}, curve.Lerp(1.0))
	assert.Equal(Point{X: 2.5, Y: 0}, curve.Lerp(0.25))
	points = []Point{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 1, Y: 1},
		{X: 0, Y: 1},
	}
	curve.Points = points
	curve.Closed = true
	assert.Equal(Point{X: 0, Y: 1}, curve.Lerp(0.75))
	assert.Equal(Point{X: 0, Y: 0.5}, curve.Lerp(0.875))
}
