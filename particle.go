package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type particle interface {
	Draw()
	Animate()
	canDelete() bool
}

type sparks struct {
	timer, timerMax        int
	positions, speeds      []V2
	lives, maxlives, seeds []uint8
	life                   int
	sparksNo               int
}

func newSparks(pos, mspeed V2, maxradius, duration float64) *sparks {
	s := new(sparks)
	s.sparksNo = 50 + rand.Intn(20)
	speed := 0.5 + rnd()*1.5
	s.positions = make([]V2, s.sparksNo)
	s.speeds = make([]V2, s.sparksNo)
	s.lives = make([]uint8, s.sparksNo)
	s.maxlives = make([]uint8, s.sparksNo)
	s.seeds = make([]uint8, s.sparksNo)

	angle := 0.0
	frames := int(duration * FPS)

	s.life = frames
	for i := 0; i < s.sparksNo; i++ {
		angle += (360 / float64(s.sparksNo)) + rndSym(15)
		s.positions[i] = pos
		s.speeds[i] = mspeed.Add(rotV(angle).MulA(5 * speed * (0.5 + rnd())))
		s.maxlives[i] = uint8(frames/2 + rand.Intn(frames/2))
		s.seeds[i] = uint8(rand.Intn(256))
	}

	return s
}

func (s *sparks) canDelete() bool {
	if s.life > 0 {
		return false
	} else {
		return true
	}
}

func (s *sparks) Animate() {
	for i := 0; i < s.sparksNo; i++ {
		age := float64(s.lives[i]) / 10
		disturb := _noise2D(int(s.lives[i] + s.seeds[i])).MulA(age).SubA(age / 2)
		s.positions[i] = s.positions[i].Add(s.speeds[i].Add(disturb))
		s.speeds[i] = s.speeds[i].MulA(0.996)
		if s.lives[i] < s.maxlives[i] {
			s.lives[i]++
		}
	}
	if s.life > 0 {
		s.life--

	}
}
func (s *sparks) Draw() {
	for i := 0; i < s.sparksNo; i++ {
		if s.lives[i] < s.maxlives[i] {
			if s.lives[i] < s.maxlives[i]/3 {
				c := _colorBlend(s.lives[i], s.maxlives[i]/3, rl.Orange, rl.Red)
				_square(s.positions[i], 2, c)
			} else {
				t := float32(s.lives[i]-s.maxlives[i]/3) / (float32(s.maxlives[i] / 3 * 2))
				v := float32(rand.Intn(2))
				c := rl.ColorFromHSV(t*33, 1.0, v)

				//_square(s.positions[i], 2, rl.ColorAlpha(c, 1-t))
				_square(s.positions[i], 2, c)

			}
		}
	}
}

type explosion struct {
	timer, timerMax       int
	position, speed, offs V2
	r, rstep              float64
	maxr, dur, t          float64
}

func newExplosion(pos, speed V2, maxradius, duration float64) *explosion {
	e := new(explosion)
	e.position = pos
	e.speed = speed
	e.offs = e.position.Add(V2{rndSym(maxradius / 10), rndSym(maxradius / 10)})
	e.maxr, e.dur = maxradius, duration
	e.rstep = maxradius / (duration * FPS)
	e.timerMax = int(duration * FPS)
	return e
}

func (e *explosion) Animate() {
	if e.timer < e.timerMax {
		e.timer++
	}
	e.position = e.position.Add(e.speed)
	e.offs = e.offs.Add(e.speed)
}

func (e *explosion) canDelete() bool {
	return e.timer >= e.timerMax
}

func (e *explosion) Draw() {
	t := 1 - e.timer/(e.timerMax/2)
	if e.timer < e.timerMax/3 {
		_gradientdisc(e.position, e.r*e.r*e.r/5, rl.ColorAlpha(rl.Yellow, float32(t)*0.3), rl.Black)
	}
	if e.timer < e.timerMax/2 {
		_disc(e.position, e.r, rl.Yellow)

	} else {
		t := e.r*2 - e.r/2
		_gradientdisc(e.position, e.maxr, rl.ColorAlpha(rl.Yellow, float32(t)), rl.ColorAlpha(rl.Orange, float32(t)))
		_disc(e.offs, t, rl.Black)

	}
	e.r += e.rstep
}
