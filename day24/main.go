package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

type bigPoint struct {
	X, Y, Z *big.Float
}
type point struct {
	X, Y, Z float64
}
type bigHailstone [2]bigPoint
type hailstone [2]point

type mat2x2 [2][2]float64
type mat3x3 [3][3]*big.Float

func (m mat2x2) det() float64 {
	return m[0][0]*m[1][1] - m[0][1]*m[1][0]
}
func (m mat3x3) det() *big.Float {
	// a1b2c3 + b1c2a3 + c1a2b3 - c1b2a3 - b1a2c3 - a1c2b3
	a1b2c3 := new(big.Float).Mul(new(big.Float).Mul(m[0][0], m[1][1]), m[2][2])
	b1c2a3 := new(big.Float).Mul(new(big.Float).Mul(m[0][1], m[1][2]), m[2][0])
	c1a2b3 := new(big.Float).Mul(new(big.Float).Mul(m[0][2], m[1][0]), m[2][1])
	c1b2a3 := new(big.Float).Mul(new(big.Float).Mul(m[0][2], m[1][1]), m[2][0])
	b1a2c3 := new(big.Float).Mul(new(big.Float).Mul(m[0][1], m[1][0]), m[2][2])
	a1c2b3 := new(big.Float).Mul(new(big.Float).Mul(m[0][0], m[1][2]), m[2][1])
	ret := new(big.Float).Sub(new(big.Float).Add(new(big.Float).Add(a1b2c3, b1c2a3), c1a2b3),
		new(big.Float).Add(new(big.Float).Add(c1b2a3, b1a2c3), a1c2b3))
	return ret
}

func newFloat(f float64) *big.Float {
	return new(big.Float).SetPrec(256).SetFloat64(f)
}
func solvedets(h1, h2, h3 bigHailstone) map[string]float64 {
	x1, y1, z1 := h1[0].X, h1[0].Y, h1[0].Z
	vx1, vy1, vz1 := h1[1].X, h1[1].Y, h1[1].Z
	x2, y2, z2 := h2[0].X, h2[0].Y, h2[0].Z
	vx2, vy2, vz2 := h2[1].X, h2[1].Y, h2[1].Z
	x3, y3, z3 := h3[0].X, h3[0].Y, h3[0].Z
	vx3, vy3, vz3 := h3[1].X, h3[1].Y, h3[1].Z
	yzMatrix := mat3x3{{newFloat(1), y1, z1}, {newFloat(1), y2, z2}, {newFloat(1), y3, z3}}
	xzMatrix := mat3x3{{x1, newFloat(1), z1}, {x2, newFloat(1), z2}, {x3, newFloat(1), z3}}
	xyMatrix := mat3x3{{x1, y1, newFloat(1)}, {x2, y2, newFloat(1)}, {x3, y3, newFloat(1)}}
	vxvyMatrix := mat3x3{{vx1, vy1, newFloat(1)}, {vx2, vy2, newFloat(1)}, {vx3, vy3, newFloat(1)}}
	vxvzMatrix := mat3x3{{vx1, newFloat(1), vz1}, {vx2, newFloat(1), vz2}, {vx3, newFloat(1), vz3}}
	vyvzMatrix := mat3x3{{newFloat(1), vy1, vz1}, {newFloat(1), vy2, vz2}, {newFloat(1), vy3, vz3}}
	//float64 innacurate for operation but rounding after is fine
	yz, _ := yzMatrix.det().Float64()
	xz, _ := xzMatrix.det().Float64()
	xy, _ := xyMatrix.det().Float64()
	vxvy, _ := vxvyMatrix.det().Float64()
	vxvz, _ := vxvzMatrix.det().Float64()
	vyvz, _ := vyvzMatrix.det().Float64()
	values := map[string]float64{
		"yz":   yz,
		"xz":   xz,
		"xy":   xy,
		"vxvy": vxvy,
		"vxvz": vxvz,
		"vyvz": vyvz,
	}
	return values
}
func solve(h2, h3 hailstone, dets map[string]float64) (float64, error) {
	x2, y2, z2 := h2[0].X, h2[0].Y, h2[0].Z
	vx2, vy2, vz2 := h2[1].X, h2[1].Y, h2[1].Z
	x3, y3, z3 := h3[0].X, h3[0].Y, h3[0].Z
	vx3, vy3, vz3 := h3[1].X, h3[1].Y, h3[1].Z
	yz, xz, xy, vxvy, vxvz, vyvz := dets["yz"], dets["xz"], dets["xy"], dets["vxvy"], dets["vxvz"], dets["vyvz"]
	vx2Vx3 := vx2 - vx3
	vy2Vy3 := vy2 - vy3
	vz2Vz3 := vz2 - vz3
	z2Z3 := z2 - z3
	y2Y3 := y2 - y3
	x2X3 := x2 - x3
	n := vx2Vx3*yz + vy2Vy3*xz + vz2Vz3*xy
	d := z2Z3*vxvy + y2Y3*vxvz + x2X3*vyvz
	// Check for zero denominator and handle it
	if d == 0 {
		return 0, fmt.Errorf("denominator is zero")
	}
	return n / d, nil
}
func part1(h []hailstone) int {
	count := 0
	for i, h1 := range h {
		for _, h2 := range h[i+1:] {
			p1, v1, p2, v2 := h1[0], h1[1], h2[0], h2[1]
			det := mat2x2{{v2.X, v2.Y}, {v1.X, v1.Y}}.det()
			t1 := ((p2.Y-p1.Y)*v2.X - (p2.X-p1.X)*v2.Y) / det
			t2 := ((p2.Y-p1.Y)*v1.X - (p2.X-p1.X)*v1.Y) / det
			px := p1.X + t1*v1.X
			py := p1.Y + t1*v1.Y
			lower := 200000000000000.0
			upper := 400000000000000.0
			if t1 > 0 && t2 > 0 && px >= lower && px <= upper && py >= lower && py <= upper {
				count++
			}
		}
	}
	return count
}

func part2(h []bigHailstone, g []hailstone) int {
	dets1 := solvedets(h[0], h[1], h[2])
	dets2 := solvedets(h[1], h[0], h[2])

	t1, err := solve(g[1], g[2], dets1)
	if err != nil {
		panic(err)
	}
	t2, err := solve(g[0], g[2], dets2)
	if err != nil {
		panic(err)
	}
	p1, v1 := g[0][0], g[0][1]
	p2, v2 := g[1][0], g[1][1]
	c1 := point{p1.X + t1*v1.X, p1.Y + t1*v1.Y, p1.Z + t1*v1.Z}
	c2 := point{p2.X + t2*v2.X, p2.Y + t2*v2.Y, p2.Z + t2*v2.Z}
	v := point{(c2.X - c1.X) / (t2 - t1), (c2.Y - c1.Y) / (t2 - t1), (c2.Z - c1.Z) / (t2 - t1)}
	p := point{p1.X + v1.X*t1 - v.X*t1, p1.Y + v1.Y*t1 - v.Y*t1, p1.Z + v1.Z*t1 - v.Z*t1}
	ans := int(math.Round(p.X + p.Y + p.Z))
	return ans
}
func main() {
	var bigHails []bigHailstone
	var hails3 []hailstone
	var allHails []hailstone

	input, _ := os.ReadFile("input.txt")
	lines := strings.Split(strings.TrimSpace(string(input)), "\n")

	for i, s := range lines {
		var x, y, z, vx, vy, vz float64
		_, err := fmt.Sscanf(s, "%f, %f, %f @ %f, %f, %f", &x, &y, &z, &vx, &vy, &vz)
		if err != nil {
			fmt.Println("Error reading line:", err)
			return
		}

		// Add to allHails as float64
		allHails = append(allHails, hailstone{
			point{x, y, z},
			point{vx, vy, vz},
		})

		// For the first three points, add to bigHails as big.Float, and hails3 as float64
		if i < 3 {
			bigHails = append(bigHails, bigHailstone{
				bigPoint{newFloat(x), newFloat(y), newFloat(z)},
				bigPoint{newFloat(vx), newFloat(vy), newFloat(vz)},
			})
			hails3 = append(hails3, hailstone{
				point{x, y, z},
				point{vx, vy, vz},
			})
		}
	}
	fmt.Println(part1(allHails))
	fmt.Println(part2(bigHails, hails3))
}
