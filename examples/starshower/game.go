// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"math"
	"math/rand"
	"time"
)

const (
	fixedStepSeconds = 1.0 / 120.0
	maximumFrameTime = 100 * time.Millisecond
)

type vector struct {
	x float64
	y float64
}

func (v vector) add(other vector) vector {
	return vector{x: v.x + other.x, y: v.y + other.y}
}

func (v vector) scale(amount float64) vector {
	return vector{x: v.x * amount, y: v.y * amount}
}

func direction(angle float64) vector {
	return vector{x: math.Cos(angle), y: math.Sin(angle)}
}

type gameMode uint8

const (
	modePlaying gameMode = iota
	modePaused
	modeGameOver
)

type inputState struct {
	left   bool
	right  bool
	thrust bool
	fire   bool
}

type ship struct {
	position     vector
	velocity     vector
	angle        float64
	invulnerable float64
}

type shot struct {
	position vector
	velocity vector
	lifetime float64
}

type asteroid struct {
	position vector
	velocity vector
	radius   float64
	angle    float64
	spin     float64
	shape    []float64
}

type star struct {
	position vector
	radius   float64
	twinkle  float64
}

type game struct {
	width       float64
	height      float64
	random      *rand.Rand
	ship        ship
	shots       []shot
	asteroids   []asteroid
	stars       []star
	input       inputState
	mode        gameMode
	score       int
	lives       int
	wave        int
	fireDelay   float64
	clock       float64
	accumulator float64
}

func newGame(width, height float64, seed int64) *game {
	g := &game{
		width:  width,
		height: height,
		random: rand.New(rand.NewSource(seed)),
	}
	g.makeStars(90)
	g.restart()
	return g
}

func (g *game) restart() {
	g.score = 0
	g.lives = 3
	g.wave = 1
	g.mode = modePlaying
	g.clock = 0
	g.accumulator = 0
	g.fireDelay = 0
	g.input = inputState{}
	g.shots = nil
	g.resetShip(2)
	g.asteroids = nil
	g.spawnWave()
}

func (g *game) resetShip(invulnerable float64) {
	g.ship = ship{
		position:     vector{x: g.width / 2, y: g.height / 2},
		angle:        -math.Pi / 2,
		invulnerable: invulnerable,
	}
}

func (g *game) togglePause() {
	switch g.mode {
	case modePlaying:
		g.mode = modePaused
		g.input = inputState{}
		g.accumulator = 0
	case modePaused:
		g.mode = modePlaying
		g.accumulator = 0
	}
}

func (g *game) advance(elapsed time.Duration) {
	if g.mode != modePlaying {
		g.accumulator = 0
		return
	}
	if elapsed < 0 {
		elapsed = 0
	}
	if elapsed > maximumFrameTime {
		elapsed = maximumFrameTime
	}
	g.accumulator += elapsed.Seconds()
	for g.accumulator >= fixedStepSeconds {
		g.update(fixedStepSeconds)
		g.accumulator -= fixedStepSeconds
	}
}

func (g *game) update(dt float64) {
	g.clock += dt
	g.updateShip(dt)
	g.updateShots(dt)
	g.updateAsteroids(dt)
	g.resolveShotHits()
	g.resolveShipHit()

	if g.mode == modePlaying && len(g.asteroids) == 0 {
		g.wave++
		g.resetShip(1.5)
		g.spawnWave()
	}
}

func (g *game) updateShip(dt float64) {
	const (
		turnSpeed  = 3.7
		thrust     = 250.0
		maximum    = 330.0
		shotDelay  = 0.17
		shotSpeed  = 460.0
		shotLife   = 1.15
		shipRadius = 17.0
	)

	turn := 0.0
	if g.input.left {
		turn--
	}
	if g.input.right {
		turn++
	}
	g.ship.angle += turn * turnSpeed * dt

	if g.input.thrust {
		g.ship.velocity = g.ship.velocity.add(direction(g.ship.angle).scale(thrust * dt))
	}
	speed := math.Hypot(g.ship.velocity.x, g.ship.velocity.y)
	if speed > maximum {
		g.ship.velocity = g.ship.velocity.scale(maximum / speed)
	}
	g.ship.velocity = g.ship.velocity.scale(math.Pow(0.994, dt*60))
	g.ship.position = wrapPoint(g.ship.position.add(g.ship.velocity.scale(dt)), g.width, g.height)

	if g.ship.invulnerable > 0 {
		g.ship.invulnerable = max(0, g.ship.invulnerable-dt)
	}
	g.fireDelay -= dt
	if g.input.fire && g.fireDelay <= 0 {
		forward := direction(g.ship.angle)
		g.shots = append(g.shots, shot{
			position: wrapPoint(g.ship.position.add(forward.scale(shipRadius+4)), g.width, g.height),
			velocity: g.ship.velocity.add(forward.scale(shotSpeed)),
			lifetime: shotLife,
		})
		g.fireDelay = shotDelay
	}
}

func (g *game) updateShots(dt float64) {
	alive := g.shots[:0]
	for _, current := range g.shots {
		current.position = wrapPoint(current.position.add(current.velocity.scale(dt)), g.width, g.height)
		current.lifetime -= dt
		if current.lifetime > 0 {
			alive = append(alive, current)
		}
	}
	g.shots = alive
}

func (g *game) updateAsteroids(dt float64) {
	for index := range g.asteroids {
		current := &g.asteroids[index]
		current.position = wrapPoint(current.position.add(current.velocity.scale(dt)), g.width, g.height)
		current.angle += current.spin * dt
	}
}

func (g *game) resolveShotHits() {
	for shotIndex := len(g.shots) - 1; shotIndex >= 0; shotIndex-- {
		hit := -1
		for asteroidIndex := range g.asteroids {
			current := &g.asteroids[asteroidIndex]
			if wrappedDistanceSquared(g.shots[shotIndex].position, current.position, g.width, g.height) <= current.radius*current.radius {
				hit = asteroidIndex
				break
			}
		}
		if hit < 0 {
			continue
		}

		g.shots = removeShot(g.shots, shotIndex)
		destroyed := g.asteroids[hit]
		g.asteroids = removeAsteroid(g.asteroids, hit)
		g.score += scoreForRadius(destroyed.radius)
		g.splitAsteroid(destroyed)
	}
}

func (g *game) resolveShipHit() {
	if g.ship.invulnerable > 0 || g.mode != modePlaying {
		return
	}
	const shipRadius = 14.0
	for index := range g.asteroids {
		current := &g.asteroids[index]
		distance := shipRadius + current.radius*0.82
		if wrappedDistanceSquared(g.ship.position, current.position, g.width, g.height) > distance*distance {
			continue
		}

		g.lives--
		g.input = inputState{}
		g.shots = nil
		if g.lives <= 0 {
			g.lives = 0
			g.mode = modeGameOver
			return
		}
		g.resetShip(2.25)
		return
	}
}

func (g *game) spawnWave() {
	count := min(3+g.wave, 8)
	for index := 0; index < count; index++ {
		position := g.safeAsteroidPosition()
		radius := 38 + g.random.Float64()*10
		g.asteroids = append(g.asteroids, g.makeAsteroid(position, radius, 32+float64(g.wave)*4))
	}
}

func (g *game) safeAsteroidPosition() vector {
	center := vector{x: g.width / 2, y: g.height / 2}
	for attempt := 0; attempt < 20; attempt++ {
		position := vector{x: g.random.Float64() * g.width, y: g.random.Float64() * g.height}
		if wrappedDistanceSquared(position, center, g.width, g.height) > 190*190 {
			return position
		}
	}
	return vector{x: 30, y: 30}
}

func (g *game) makeAsteroid(position vector, radius, minimumSpeed float64) asteroid {
	travel := g.random.Float64() * 2 * math.Pi
	speed := minimumSpeed + g.random.Float64()*38
	points := 9 + g.random.Intn(4)
	shape := make([]float64, points)
	for index := range shape {
		shape[index] = 0.76 + g.random.Float64()*0.32
	}
	spin := 0.2 + g.random.Float64()*0.55
	if g.random.Intn(2) == 0 {
		spin = -spin
	}
	return asteroid{
		position: position,
		velocity: direction(travel).scale(speed),
		radius:   radius,
		angle:    g.random.Float64() * 2 * math.Pi,
		spin:     spin,
		shape:    shape,
	}
}

func (g *game) splitAsteroid(parent asteroid) {
	if parent.radius < 24 {
		return
	}
	childRadius := parent.radius * 0.61
	for index := 0; index < 2; index++ {
		child := g.makeAsteroid(parent.position, childRadius, 70+float64(g.wave)*5)
		child.position = wrapPoint(parent.position.add(direction(child.angle).scale(childRadius*0.35)), g.width, g.height)
		g.asteroids = append(g.asteroids, child)
	}
}

func (g *game) makeStars(count int) {
	g.stars = make([]star, count)
	for index := range g.stars {
		g.stars[index] = star{
			position: vector{x: g.random.Float64() * g.width, y: g.random.Float64() * g.height},
			radius:   0.5 + g.random.Float64()*1.35,
			twinkle:  g.random.Float64() * 2 * math.Pi,
		}
	}
}

func scoreForRadius(radius float64) int {
	switch {
	case radius >= 35:
		return 20
	case radius >= 22:
		return 50
	default:
		return 100
	}
}

func removeShot(shots []shot, index int) []shot {
	return append(shots[:index], shots[index+1:]...)
}

func removeAsteroid(asteroids []asteroid, index int) []asteroid {
	return append(asteroids[:index], asteroids[index+1:]...)
}

func wrapPoint(point vector, width, height float64) vector {
	if width > 0 {
		point.x = math.Mod(point.x, width)
		if point.x < 0 {
			point.x += width
		}
	}
	if height > 0 {
		point.y = math.Mod(point.y, height)
		if point.y < 0 {
			point.y += height
		}
	}
	return point
}

func wrappedDistanceSquared(a, b vector, width, height float64) float64 {
	dx := math.Abs(a.x - b.x)
	dy := math.Abs(a.y - b.y)
	if width > 0 {
		dx = min(dx, width-dx)
	}
	if height > 0 {
		dy = min(dy, height-dy)
	}
	return dx*dx + dy*dy
}

func wrappedPositions(position vector, radius, width, height float64) []vector {
	xOffsets := []float64{0}
	yOffsets := []float64{0}
	if position.x < radius {
		xOffsets = append(xOffsets, width)
	}
	if position.x > width-radius {
		xOffsets = append(xOffsets, -width)
	}
	if position.y < radius {
		yOffsets = append(yOffsets, height)
	}
	if position.y > height-radius {
		yOffsets = append(yOffsets, -height)
	}

	positions := make([]vector, 0, len(xOffsets)*len(yOffsets))
	for _, xOffset := range xOffsets {
		for _, yOffset := range yOffsets {
			positions = append(positions, vector{x: position.x + xOffset, y: position.y + yOffset})
		}
	}
	return positions
}
