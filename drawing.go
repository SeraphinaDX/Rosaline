// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// Rect describes a rectangle used for clipping and other drawing operations.
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// TextStyle controls text drawn on a Canvas.
type TextStyle struct {
	Color Color
	Size  int
}

type drawingState struct {
	matrix gg.Matrix
	clips  [][]gg.Point
}

// DrawingCanvas provides Rosaline's beginner-friendly 2D drawing operations.
// The same API draws both visible widgets and off-screen images.
type DrawingCanvas struct {
	context    *gg.Context
	background Color
	matrix     gg.Matrix
	clips      [][]gg.Point
	stack      []drawingState
}

func newDrawingCanvas(width, height int, background Color) *DrawingCanvas {
	context := gg.NewContext(width, height)
	canvas := &DrawingCanvas{
		context:    context,
		background: background,
		matrix:     gg.Identity(),
	}
	canvas.Clear(background)
	return canvas
}

func renderDrawing(width, height int, background Color, draw func(*DrawingCanvas)) (*image.RGBA, Color) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if draw == nil {
		draw = func(*DrawingCanvas) {}
	}

	canvas := newDrawingCanvas(width, height, background)
	draw(canvas)
	return canvas.context.Image().(*image.RGBA), canvas.background
}

func nativeColor(value Color) color.NRGBA {
	return color.NRGBA{R: value.R, G: value.G, B: value.B, A: value.A}
}

// Clear removes existing drawing and fills the entire canvas with color.
// Clear is unaffected by the current transform or clipping region.
func (c *DrawingCanvas) Clear(value Color) {
	if c == nil || c.context == nil {
		return
	}
	c.background = value
	c.context.SetColor(nativeColor(value))
	c.context.Clear()
}

// FillRect draws a filled rectangle.
func (c *DrawingCanvas) FillRect(x, y, width, height float64, value Color) {
	if c == nil || c.context == nil {
		return
	}
	c.context.SetColor(nativeColor(value))
	c.context.DrawRectangle(x, y, width, height)
	c.context.Fill()
}

// Rect draws the outline of a rectangle.
func (c *DrawingCanvas) Rect(x, y, width, height, stroke float64, value Color) {
	if c == nil || c.context == nil {
		return
	}
	c.prepareStroke(stroke, value)
	c.context.DrawRectangle(x, y, width, height)
	c.context.Stroke()
}

// Line draws a line.
func (c *DrawingCanvas) Line(x1, y1, x2, y2, stroke float64, value Color) {
	if c == nil || c.context == nil {
		return
	}
	c.prepareStroke(stroke, value)
	c.context.DrawLine(x1, y1, x2, y2)
	c.context.Stroke()
}

// FillCircle draws a filled circle.
func (c *DrawingCanvas) FillCircle(x, y, radius float64, value Color) {
	if c == nil || c.context == nil || radius < 0 {
		return
	}
	c.context.SetColor(nativeColor(value))
	c.context.DrawCircle(x, y, radius)
	c.context.Fill()
}

// Circle draws the outline of a circle.
func (c *DrawingCanvas) Circle(x, y, radius, stroke float64, value Color) {
	if c == nil || c.context == nil || radius < 0 {
		return
	}
	c.prepareStroke(stroke, value)
	c.context.DrawCircle(x, y, radius)
	c.context.Stroke()
}

// FillPath fills a reusable Path.
func (c *DrawingCanvas) FillPath(path *Path, value Color) {
	if c == nil || c.context == nil || path == nil {
		return
	}
	c.context.SetColor(nativeColor(value))
	path.replay(c.context)
	c.context.Fill()
}

// StrokePath draws the outline of a reusable Path.
func (c *DrawingCanvas) StrokePath(path *Path, stroke float64, value Color) {
	if c == nil || c.context == nil || path == nil {
		return
	}
	c.prepareStroke(stroke, value)
	path.replay(c.context)
	c.context.Stroke()
}

func (c *DrawingCanvas) prepareStroke(stroke float64, value Color) {
	if stroke < 0 {
		stroke = 0
	}
	c.context.SetColor(nativeColor(value))
	c.context.SetLineWidth(stroke * c.strokeScale())
	c.context.SetLineCapRound()
	c.context.SetLineJoinRound()
}

func (c *DrawingCanvas) strokeScale() float64 {
	xScale := math.Hypot(c.matrix.XX, c.matrix.YX)
	yScale := math.Hypot(c.matrix.XY, c.matrix.YY)
	scale := (xScale + yScale) / 2
	if scale <= 0 {
		return 1
	}
	return scale
}

// Text draws text from its top-left corner.
func (c *DrawingCanvas) Text(text string, x, y float64, style TextStyle) {
	if c == nil || c.context == nil {
		return
	}
	if style.Color.A == 0 {
		style.Color = Black
	}
	if style.Size <= 0 {
		style.Size = 14
	}
	c.context.SetColor(nativeColor(style.Color))
	c.context.SetFontFace(drawingFont(style.Size))
	c.context.DrawStringAnchored(text, x, y, 0, 1)
}

var parsedDrawingFont, drawingFontError = opentype.Parse(goregular.TTF)

func drawingFont(size int) font.Face {
	if drawingFontError != nil {
		return basicfont.Face7x13
	}
	face, err := opentype.NewFace(parsedDrawingFont, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return basicfont.Face7x13
	}
	return face
}

// Push saves the current transform and clipping region. Push calls may be
// nested and paired with Pop.
func (c *DrawingCanvas) Push() {
	if c == nil || c.context == nil {
		return
	}
	c.stack = append(c.stack, drawingState{
		matrix: c.matrix,
		clips:  copyClips(c.clips),
	})
	c.context.Push()
}

// Pop restores the most recently saved transform and clipping region. Calling
// Pop without a matching Push has no effect.
func (c *DrawingCanvas) Pop() {
	if c == nil || c.context == nil || len(c.stack) == 0 {
		return
	}
	last := len(c.stack) - 1
	state := c.stack[last]
	c.stack = c.stack[:last]
	c.context.Pop()
	c.matrix = state.matrix
	c.clips = state.clips
	c.applyClips()
}

// Translate moves subsequent drawing by x and y pixels.
func (c *DrawingCanvas) Translate(x, y float64) {
	if c == nil || c.context == nil {
		return
	}
	c.context.Translate(x, y)
	c.matrix = c.matrix.Translate(x, y)
}

// Scale scales subsequent drawing around the current origin.
func (c *DrawingCanvas) Scale(x, y float64) {
	if c == nil || c.context == nil {
		return
	}
	c.context.Scale(x, y)
	c.matrix = c.matrix.Scale(x, y)
}

// Rotate rotates subsequent drawing clockwise by degrees around the current
// origin.
func (c *DrawingCanvas) Rotate(degrees float64) {
	if c == nil || c.context == nil {
		return
	}
	radians := degrees * math.Pi / 180
	c.context.Rotate(radians)
	c.matrix = c.matrix.Rotate(radians)
}

// ResetTransform restores the normal untransformed coordinate system.
func (c *DrawingCanvas) ResetTransform() {
	if c == nil || c.context == nil {
		return
	}
	c.context.Identity()
	c.matrix = gg.Identity()
}

// Clip restricts subsequent drawing to rect. The rectangle follows the
// current transform and intersects any existing clipping region.
func (c *DrawingCanvas) Clip(rect Rect) {
	if c == nil || c.context == nil || rect.Width <= 0 || rect.Height <= 0 {
		return
	}
	points := make([]gg.Point, 4)
	points[0].X, points[0].Y = c.matrix.TransformPoint(rect.X, rect.Y)
	points[1].X, points[1].Y = c.matrix.TransformPoint(rect.X+rect.Width, rect.Y)
	points[2].X, points[2].Y = c.matrix.TransformPoint(rect.X+rect.Width, rect.Y+rect.Height)
	points[3].X, points[3].Y = c.matrix.TransformPoint(rect.X, rect.Y+rect.Height)
	c.clips = append(c.clips, points)
	c.applyClip(points)
}

// ClipRect is a convenient form of Clip using separate coordinates.
func (c *DrawingCanvas) ClipRect(x, y, width, height float64) {
	c.Clip(Rect{X: x, Y: y, Width: width, Height: height})
}

// ResetClip removes every active clipping region.
func (c *DrawingCanvas) ResetClip() {
	if c == nil || c.context == nil {
		return
	}
	c.clips = nil
	c.context.ResetClip()
}

func copyClips(clips [][]gg.Point) [][]gg.Point {
	result := make([][]gg.Point, len(clips))
	for index, clip := range clips {
		result[index] = append([]gg.Point(nil), clip...)
	}
	return result
}

func (c *DrawingCanvas) applyClips() {
	c.context.ResetClip()
	for _, clip := range c.clips {
		c.applyClip(clip)
	}
}

func (c *DrawingCanvas) applyClip(points []gg.Point) {
	if len(points) < 3 {
		return
	}
	c.context.Push()
	c.context.Identity()
	c.context.MoveTo(points[0].X, points[0].Y)
	for _, point := range points[1:] {
		c.context.LineTo(point.X, point.Y)
	}
	c.context.ClosePath()
	c.context.Clip()
	c.context.Pop()
}
