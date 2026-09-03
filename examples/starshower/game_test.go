// SPDX-License-Identifier: LGPL-3.0-or-later

package main

import (
	"math"
	"testing"
	"time"
)

func TestWrapPointWrapsEveryEdge(t *testing.T) {
	got := wrapPoint(vector{x: -4, y: 103}, 100, 100)
	if got != (vector{x: 96, y: 3}) {
		t.Fatalf("wrapPoint() = %#v, want {96 3}", got)
	}
}

func TestWrappedDistanceSeesAcrossScreenEdge(t *testing.T) {
	distance := wrappedDistanceSquared(vector{x: 2, y: 50}, vector{x: 98, y: 50}, 100, 100)
	if distance != 16 {
		t.Fatalf("wrapped distance = %v, want 16", distance)
	}
}

func TestFixedStepDoesNotAdvanceWhilePaused(t *testing.T) {
	game := newGame(840, 520, 7)
	before := game.ship.position
	game.input.thrust = true
	game.togglePause()
	game.advance(time.Second)
	if game.ship.position != before {
		t.Fatalf("paused ship moved from %#v to %#v", before, game.ship.position)
	}
	if game.input != (inputState{}) {
		t.Fatal("pausing should release held controls")
	}
}

func TestHeldThrustMovesShip(t *testing.T) {
	game := newGame(840, 520, 11)
	before := game.ship.position
	game.input.thrust = true
	game.advance(100 * time.Millisecond)
	if !(game.ship.position.y < before.y) {
		t.Fatalf("ship did not thrust upward: before %#v, after %#v", before, game.ship.position)
	}
}

func TestHeldFireCreatesShots(t *testing.T) {
	game := newGame(840, 520, 13)
	game.input.fire = true
	game.advance(20 * time.Millisecond)
	if len(game.shots) != 1 {
		t.Fatalf("shots = %d, want 1", len(game.shots))
	}
}

func TestShotSplitsLargeAsteroidAndScores(t *testing.T) {
	game := newGame(840, 520, 17)
	game.asteroids = []asteroid{game.makeAsteroid(vector{x: 120, y: 120}, 40, 0)}
	game.asteroids[0].velocity = vector{}
	game.shots = []shot{{position: vector{x: 120, y: 120}, lifetime: 1}}

	game.resolveShotHits()
	if game.score != 20 {
		t.Fatalf("score = %d, want 20", game.score)
	}
	if len(game.asteroids) != 2 {
		t.Fatalf("asteroids = %d, want two smaller pieces", len(game.asteroids))
	}
	if !(game.asteroids[0].radius < 40) {
		t.Fatalf("child radius = %v, want less than 40", game.asteroids[0].radius)
	}
}

func TestShipCollisionUsesLifeAndEndsGame(t *testing.T) {
	game := newGame(840, 520, 19)
	game.ship.invulnerable = 0
	game.lives = 1
	game.asteroids = []asteroid{{position: game.ship.position, radius: 40}}

	game.resolveShipHit()
	if game.lives != 0 || game.mode != modeGameOver {
		t.Fatalf("collision produced lives=%d mode=%d, want game over", game.lives, game.mode)
	}
}

func TestClearingWaveStartsNextWave(t *testing.T) {
	game := newGame(840, 520, 23)
	game.asteroids = nil
	game.update(fixedStepSeconds)
	if game.wave != 2 {
		t.Fatalf("wave = %d, want 2", game.wave)
	}
	if len(game.asteroids) == 0 {
		t.Fatal("next wave should contain asteroids")
	}
}

func TestWrappedPositionsDuplicatesCornerObject(t *testing.T) {
	positions := wrappedPositions(vector{x: 2, y: 3}, 10, 100, 80)
	if len(positions) != 4 {
		t.Fatalf("positions = %d, want 4", len(positions))
	}
	foundCorner := false
	for _, position := range positions {
		if math.Abs(position.x-102) < 0.001 && math.Abs(position.y-83) < 0.001 {
			foundCorner = true
		}
	}
	if !foundCorner {
		t.Fatal("wrapped corner copy was not generated")
	}
}
