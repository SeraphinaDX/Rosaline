# Building Starshower

Starshower is Rosaline v0.15.0's complete canvas-game tutorial. It combines a
fixed-step game model, held keyboard input, vector drawing, transforms, screen
wrapping, collision detection, timers, menus, buttons, score and lives, pause,
restart, and non-visual tests.

Run it from the project root:

```bash
CGO_ENABLED=0 go run ./examples/starshower
```

Click the game canvas if it does not already have keyboard focus. Turn with the
arrow keys or A and D, thrust with Up or W, hold Space to fire, and pause with P
or Escape.

## Separate the game from the window

The game model lives in `game.go`. It stores plain Go values:

```go
type game struct {
	ship      ship
	shots     []shot
	asteroids []asteroid
	input     inputState
	score     int
	lives     int
	wave      int
}
```

It does not create widgets or call the window backend. That keeps movement,
collisions, scoring, and level progression testable without opening a window.
`main.go` handles Rosaline events and draws the current model.

This separation is useful, but it is not a framework requirement. A smaller
game can keep all its variables in `main` just like the animation example.

## Use held-key state for smooth controls

A key press sets a Boolean and its matching release clears it:

```go
canvas.OnKeyDown(func(event rosaline.KeyEvent) {
	switch event.Key {
	case rosaline.KeyLeft, rosaline.Key("a"):
		game.input.left = true
	case rosaline.KeySpace:
		game.input.fire = true
	}
})

canvas.OnKeyUp(func(event rosaline.KeyEvent) {
	switch event.Key {
	case rosaline.KeyLeft, rosaline.Key("a"):
		game.input.left = false
	case rosaline.KeySpace:
		game.input.fire = false
	}
})
```

The update loop reads those values every step. Movement therefore remains
smooth when a key is held instead of depending on operating-system key-repeat
timing. Starshower also remembers which keys are already down so a repeated P
event cannot rapidly toggle pause.

## Keep simulation steps consistent

Rendering timers can arrive a little early or late. Starshower collects the
elapsed time and updates the model using a constant step:

```go
game.accumulator += elapsed.Seconds()
for game.accumulator >= fixedStepSeconds {
	game.update(fixedStepSeconds)
	game.accumulator -= fixedStepSeconds
}
```

The fixed step makes movement and collisions predictable. A maximum elapsed
time prevents a delayed frame or closed dialog from causing a large jump. The
60 FPS Rosaline animation timer redraws after advancing the model:

```go
animation := rosaline.Animate(60, func() {
	now := time.Now()
	game.advance(now.Sub(lastFrame))
	lastFrame = now
	canvas.Redraw()
})
```

`Animate` owns scheduling and automatically stops when the application closes.
The game model owns the simulation rules.

## Draw reusable vector shapes

The ship is one `Path` created at startup:

```go
shipPath := rosaline.NewPath().
	MoveTo(20, 0).
	LineTo(-14, -12).
	LineTo(-8, 0).
	LineTo(-14, 12).
	Close()
```

Drawing moves the origin to the ship and rotates the local shape:

```go
c.Push()
c.Translate(ship.position.x, ship.position.y)
c.Rotate(ship.angle * 180 / math.Pi)
c.FillPath(shipPath, deepPlum)
c.StrokePath(shipPath, 2.6, rosaline.Rose)
c.Pop()
```

`Push` and `Pop` keep that transform from affecting asteroids, stars, or the
score. Asteroids use the same operations with slightly randomized outlines, so
the entire game remains vector-only and needs no asset files.

## Wrap positions and collisions together

Moving objects use `wrapPoint` to reappear on the opposite side. Drawing makes
a second copy near an edge so a large asteroid crosses smoothly instead of
vanishing halfway through.

Collision distance must use the same wrapped world. Two objects at X positions
2 and 838 are close together on an 840-pixel screen even though ordinary
subtraction says they are far apart. `wrappedDistanceSquared` checks the
shortest distance across either edge and avoids an unnecessary square root.

Large asteroid hits create two smaller asteroids. Small pieces disappear and
award more points. Clearing every piece increments the wave and creates a
slightly faster group.

## Treat pause and game over as model states

The game uses three explicit modes:

```go
const (
	modePlaying gameMode = iota
	modePaused
	modeGameOver
)
```

Paused and game-over frames are still drawn, but the simulation does not
advance. The same state controls the centered overlay, keyboard behavior, and
button commands. Keeping this in the model avoids several unrelated Boolean
flags that could contradict one another.

## Restore keyboard focus after buttons

Buttons and menus naturally receive focus when used. Calling `canvas.Focus()`
from their callbacks immediately returns keyboard control to the game:

```go
restart := func() {
	game.restart()
	canvas.Focus()
	canvas.Redraw()
}
```

Before a window opens, `Focus()` still requests initial focus. After mounting,
it focuses the canvas immediately. This makes the same method useful during
both setup and ordinary callbacks.

## Test the rules without opening a window

The example tests edge wrapping, wrapped collision distance, pause behavior,
held thrust and firing, asteroid splitting, scoring, lives, game over, and wave
progression. Because these rules are ordinary Go, the tests stay quick:

```bash
CGO_ENABLED=0 go test ./examples/starshower
```

## Explore the complete source

- [`examples/starshower/main.go`](../examples/starshower/main.go) builds the
  window, connects input, and draws the scene.
- [`examples/starshower/game.go`](../examples/starshower/game.go) contains the
  portable game model.
- [`examples/starshower/game_test.go`](../examples/starshower/game_test.go)
  tests its rules without a graphical desktop.

For smaller introductions, read [Canvas Animation](ANIMATION.md),
[Keyboard Input and Shortcuts](KEYBOARD_INPUT.md), and
[Drawing Paths](DRAWING_PATHS.md) first.

## Go concepts used here

- Structs, slices, and methods
- Maps used as sets
- Durations and fixed-step loops
- Trigonometry and vectors
- Random-number generators
- State machines
- Separating testable logic from presentation
