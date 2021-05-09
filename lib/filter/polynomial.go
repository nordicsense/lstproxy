package filter

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gonum/matrix/mat64"
)

type PolyVariant int

const (
	PolyGoNum PolyVariant = iota
	PolyGoMatrix
)

const timeDenominator = 1.0 / 3600.0 / 24.0

type dataPoint struct {
	t time.Time
	x float64
	y float64
}

type poly struct {
	variant    PolyVariant
	degree     int
	minPoints  int
	maxNaCount int
	lookBack   time.Duration // expected negative
	dps        []dataPoint
	c          *mat64.Dense
	naCount    int
}

func Polynomial(variant PolyVariant, degree, minPoints, maxNaCount int, lookBack time.Duration) *poly {
	return &poly{
		variant:    variant,
		degree:     degree,
		minPoints:  minPoints,
		maxNaCount: maxNaCount,
		lookBack:   lookBack,
	}
}
func (p *poly) Compute(t time.Time, v float64) float64 {
	switch p.variant {
	case PolyGoNum:
		return p.computeGoNum(t, v)
	case PolyGoMatrix:
		return p.computeGoMatrix(t, v)
	default:
		panic("not implemented")
	}
}

// https://rosettacode.org/wiki/Polynomial_regression#Library_go.matrix
func (p *poly) computeGoMatrix(t time.Time, v float64) float64 {
	panic("FIXME: not implemented")
}

// https://rosettacode.org/wiki/Polynomial_regression#Library_gonum.2Fmatrix
func (p *poly) computeGoNum(t time.Time, v float64) float64 {
	x := float64(t.Unix()) * timeDenominator
	if !math.IsNaN(v) {
		if i := sort.Search(len(p.dps), func(i int) bool {
			return p.dps[i].t.After(t.Add(p.lookBack))
		}); i > 0 {
			p.dps = p.dps[i:]
		}
		dp := dataPoint{t: t, x: x, y: v}
		p.dps = append(p.dps, dp)
		p.c = nil
		p.naCount = 0
		if len(p.dps) >= p.minPoints {
			// compute new coefficients
			a := mat64.NewDense(len(p.dps), p.degree+1, nil)
			b := mat64.NewDense(len(p.dps), 1, nil)
			p.c = mat64.NewDense(p.degree+1, 1, nil)
			for i := range p.dps {
				for j, xx := 0, 1.; j <= p.degree; j, xx = j+1, xx*p.dps[i].x {
					a.Set(i, j, xx)
				}
				b.Set(i, 0, p.dps[i].y)
			}
			qr := new(mat64.QR)
			qr.Factorize(a)
			if err := p.c.SolveQR(qr, false, b); err != nil {
				fmt.Println("Failed to solve at", t)
			}
		}
	} else {
		p.naCount++
	}
	if p.c != nil && p.naCount < p.maxNaCount {
		var res float64
		xx := x
		for i := 0; i < p.degree+1; i++ {
			c := p.c.At(i, 0)
			if i == 0 {
				res = c
			} else {
				res += c * xx
				xx *= x
			}
		}
		return res
	}
	return math.NaN()
}
