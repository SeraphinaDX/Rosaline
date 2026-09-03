// SPDX-License-Identifier: LGPL-3.0-or-later

// Starshower is a complete little vector arcade game built with Rosaline.
// Run it with: go run ./examples/starshower
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/SeraphinaDX/Rosaline"
)

const (
	gameWidth  = 840
	gameHeight = 520
)

var (
	night       = rosaline.Hex("#120b20")
	deepPlum    = rosaline.Hex("#261333")
	starlight   = rosaline.Hex("#fff0fa")
	softViolet  = rosaline.Hex("#caa4e8")
	asteroidInk = rosaline.Hex("#b98aca")
	gold        = rosaline.Hex("#ffd37a")
	coral       = rosaline.Hex("#ff9b9b")
)

func main() {
	game := newGame(gameWidth, gameHeight, time.Now().UnixNano())
	held := make(map[rosaline.Key]bool)
	lastFrame := time.Now()

	shipPath := rosaline.NewPath().
		MoveTo(20, 0).
		LineTo(-14, -12).
		LineTo(-8, 0).
		LineTo(-14, 12).
		Close()
	flamePath := rosaline.NewPath().
		MoveTo(-10, -7).
		LineTo(-25, 0).
		LineTo(-10, 7).
		Close()

	var canvas *rosaline.CanvasWidget
	canvas = rosaline.Canvas(func(c *rosaline.DrawingCanvas) {
		drawScene(c, game, shipPath, flamePath)
	}).Size(gameWidth, gameHeight).Background(night).Expand().Focus()

	releaseKeys := func() {
		clear(held)
		game.input = inputState{}
	}
	restart := func() {
		game.restart()
		releaseKeys()
		lastFrame = time.Now()
		canvas.Focus()
		canvas.Redraw()
	}
	togglePause := func() {
		game.togglePause()
		game.input = inputState{}
		lastFrame = time.Now()
		canvas.Focus()
		canvas.Redraw()
	}
	togglePauseFromControl := func() {
		togglePause()
		releaseKeys()
	}

	canvas.OnKeyDown(func(event rosaline.KeyEvent) {
		if held[event.Key] {
			return
		}
		held[event.Key] = true
		switch event.Key {
		case rosaline.KeyLeft, rosaline.Key("a"):
			game.input.left = true
		case rosaline.KeyRight, rosaline.Key("d"):
			game.input.right = true
		case rosaline.KeyUp, rosaline.Key("w"):
			game.input.thrust = true
		case rosaline.KeySpace:
			game.input.fire = true
		case rosaline.Key("p"), rosaline.KeyEscape:
			togglePause()
		case rosaline.KeyEnter:
			if game.mode == modeGameOver {
				restart()
			}
		}
	})

	canvas.OnKeyUp(func(event rosaline.KeyEvent) {
		delete(held, event.Key)
		switch event.Key {
		case rosaline.KeyLeft, rosaline.Key("a"):
			game.input.left = false
		case rosaline.KeyRight, rosaline.Key("d"):
			game.input.right = false
		case rosaline.KeyUp, rosaline.Key("w"):
			game.input.thrust = false
		case rosaline.KeySpace:
			game.input.fire = false
		}
	})

	help := func() {
		wasPlaying := game.mode == modePlaying
		if wasPlaying {
			game.togglePause()
		}
		releaseKeys()
		rosaline.Message(
			"How to play Starshower",
			"Turn with Left/Right or A/D.\nThrust with Up or W.\nHold Space to fire.\nPress P or Escape to pause.\nPress Enter after game over to restart.\n\nClear every asteroid to reach the next wave.",
		)
		if wasPlaying {
			game.togglePause()
		}
		lastFrame = time.Now()
		canvas.Focus()
	}

	animation := rosaline.Animate(60, func() {
		now := time.Now()
		game.advance(now.Sub(lastFrame))
		lastFrame = now
		canvas.Redraw()
	})

	menu := rosaline.MenuBar(
		rosaline.Menu("Game",
			rosaline.MenuItem("New Game", restart).Shortcut("Primary+R"),
			rosaline.MenuItem("Pause / Resume", togglePauseFromControl),
			rosaline.MenuSeparator(),
			rosaline.MenuItem("Quit", rosaline.Quit).Shortcut("Primary+Q"),
		),
		rosaline.Menu("Help",
			rosaline.MenuItem("How to Play", help).Shortcut("F1"),
			rosaline.MenuItem("About", func() {
				rosaline.Message("About Starshower", "A tiny vector arcade game made entirely with Rosaline and ordinary Go.")
				lastFrame = time.Now()
				canvas.Focus()
			}),
		),
	)

	theme := rosaline.DefaultTheme
	theme.Background = rosaline.Hex("#ead5e5")
	theme.Surface = rosaline.Hex("#fff7fc")
	theme.Primary = rosaline.Hex("#c63f82")
	theme.Border = rosaline.Hex("#cfadc5")

	rosaline.RunApp(rosaline.App{
		Title:   "Starshower — Rosaline Arcade",
		Width:   920,
		Height:  700,
		Padding: 14,
		Theme:   theme,
		Menu:    menu,
		Timers:  []*rosaline.Timer{animation},
		Content: rosaline.Column(
			rosaline.Row(
				rosaline.Label("STARSHOWER").Bold().FontSize(20).Color(theme.Primary),
				rosaline.Label("A ROSALINE ARCADE GAME").Color(theme.Muted),
				rosaline.Spring(),
				rosaline.LabelFunc(func() string {
					return fmt.Sprintf("SCORE %06d   LIVES %d   WAVE %d", game.score, game.lives, game.wave)
				}).Bold(),
			).Gap(12),
			rosaline.Card(canvas).Padding(5).Expand(),
			rosaline.Row(
				rosaline.Button("New Game", restart),
				rosaline.Button("Pause / Resume", togglePauseFromControl).Primary(),
				rosaline.Button("How to Play", help),
				rosaline.Spring(),
				rosaline.Label("Turn: ← → / A D   Thrust: ↑ / W   Fire: Space   Pause: P / Esc").Color(theme.Muted),
			).Gap(9),
		).Gap(10).Expand(),
	})
}

func drawScene(c *rosaline.DrawingCanvas, game *game, shipPath, flamePath *rosaline.Path) {
	c.Clear(night)

	for _, star := range game.stars {
		brightness := 0.5 + 0.5*math.Sin(game.clock*1.8+star.twinkle)
		color := softViolet
		if brightness > 0.62 {
			color = starlight
		}
		c.FillCircle(star.position.x, star.position.y, star.radius*(0.7+brightness*0.45), color)
	}

	for _, current := range game.asteroids {
		for _, position := range wrappedPositions(current.position, current.radius+4, game.width, game.height) {
			drawAsteroid(c, current, position)
		}
	}

	for _, current := range game.shots {
		for _, position := range wrappedPositions(current.position, 5, game.width, game.height) {
			c.FillCircle(position.x, position.y, 4.5, rosaline.Hex("#6c3a5a"))
			c.FillCircle(position.x, position.y, 2.2, gold)
		}
	}

	visible := game.ship.invulnerable <= 0 || int(game.ship.invulnerable*10)%2 == 0
	if game.mode != modeGameOver && visible {
		for _, position := range wrappedPositions(game.ship.position, 27, game.width, game.height) {
			drawShip(c, game, position, shipPath, flamePath)
		}
	}

	c.Text(fmt.Sprintf("SCORE %06d", game.score), 18, 16, rosaline.TextStyle{Color: starlight, Size: 17})
	c.Text(fmt.Sprintf("LIVES %d", game.lives), game.width-102, 16, rosaline.TextStyle{Color: starlight, Size: 17})
	c.Text(fmt.Sprintf("WAVE %d", game.wave), game.width/2-31, 16, rosaline.TextStyle{Color: softViolet, Size: 15})

	switch game.mode {
	case modePaused:
		drawOverlay(c, game, "PAUSED", "Press P or Escape to continue")
	case modeGameOver:
		drawOverlay(c, game, "STARSHOWER OVER", "Press Enter or choose New Game")
	}
}

func drawShip(c *rosaline.DrawingCanvas, game *game, position vector, shipPath, flamePath *rosaline.Path) {
	c.Push()
	c.Translate(position.x, position.y)
	c.Rotate(game.ship.angle * 180 / math.Pi)
	if game.input.thrust && game.mode == modePlaying {
		c.FillPath(flamePath, coral)
		c.StrokePath(flamePath, 2, gold)
	}
	c.FillPath(shipPath, deepPlum)
	c.StrokePath(shipPath, 2.6, rosaline.Rose)
	c.FillCircle(2, 0, 3, starlight)
	c.Pop()
}

func drawAsteroid(c *rosaline.DrawingCanvas, asteroid asteroid, position vector) {
	path := rosaline.NewPath()
	for index, scale := range asteroid.shape {
		angle := 2 * math.Pi * float64(index) / float64(len(asteroid.shape))
		x := math.Cos(angle) * asteroid.radius * scale
		y := math.Sin(angle) * asteroid.radius * scale
		if index == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()

	c.Push()
	c.Translate(position.x, position.y)
	c.Rotate(asteroid.angle * 180 / math.Pi)
	c.FillPath(path, deepPlum)
	c.StrokePath(path, 2.3, asteroidInk)
	c.Circle(-asteroid.radius*0.22, -asteroid.radius*0.12, asteroid.radius*0.17, 1.3, rosaline.Hex("#714d7f"))
	c.Circle(asteroid.radius*0.27, asteroid.radius*0.2, asteroid.radius*0.11, 1.2, rosaline.Hex("#714d7f"))
	c.Pop()
}

func drawOverlay(c *rosaline.DrawingCanvas, game *game, title, instruction string) {
	width := 390.0
	height := 112.0
	x := (game.width - width) / 2
	y := (game.height - height) / 2
	c.FillRect(x, y, width, height, rosaline.Hex("#1d1029e8"))
	c.Rect(x, y, width, height, 2, rosaline.Rose)
	c.Text(title, x+28, y+24, rosaline.TextStyle{Color: starlight, Size: 25})
	c.Text(instruction, x+28, y+70, rosaline.TextStyle{Color: softViolet, Size: 15})
}
