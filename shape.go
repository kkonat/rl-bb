package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type motion struct {
	pos, speed    V2
	rot, rotSpeed float64
	rotM          M22
}

func (m *motion) Move(dt float64) {
	dv := V2MulA(m.speed, dt)
	m.pos.Incr(dv)
	m.rot += m.rotSpeed * dt
	m.rotM = newM22rot(m.rot)
}

func newMotion() *motion {
	m := new(motion)
	m.rotM = newM22Id()
	return m
}

type shape struct {
	points []V2
}

func newShape(p []V2) *shape {
	s := new(shape)
	s.points = p
	return s
}

func (s *shape) Draw(m *motion, colFill, colLine rl.Color) {
	var veryfirst, pp V2

	for i, p := range s.points {
		np := m.rotM.pMulV(p)
		np.Incr(m.pos)

		if i > 0 {
			_triangle(np, pp, m.pos, colFill) // sequence of vertices matters must be counter clockwise, otherwise nothing is drawn
			_line(pp, np, colLine)
		} else {
			veryfirst = np
		}
		pp = np
	}

	_triangle(veryfirst, pp, m.pos, colFill)
	_line(pp, veryfirst, colLine)
}
