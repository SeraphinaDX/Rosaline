// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import "github.com/fogleman/gg"

type pathOperation uint8

const (
	pathMove pathOperation = iota
	pathLine
	pathQuadratic
	pathCubic
	pathClose
)

type pathCommand struct {
	operation pathOperation
	points    [6]float64
}

// Path describes a reusable shape made from straight lines and Bézier curves.
type Path struct {
	commands []pathCommand
}

// NewPath creates an empty path.
func NewPath() *Path { return &Path{} }

// MoveTo begins a new part of the path at x, y.
func (p *Path) MoveTo(x, y float64) *Path {
	if p != nil {
		p.commands = append(p.commands, pathCommand{operation: pathMove, points: [6]float64{x, y}})
	}
	return p
}

// LineTo adds a straight line from the current point to x, y.
func (p *Path) LineTo(x, y float64) *Path {
	if p != nil {
		p.commands = append(p.commands, pathCommand{operation: pathLine, points: [6]float64{x, y}})
	}
	return p
}

// QuadraticTo adds a quadratic Bézier curve with one control point.
func (p *Path) QuadraticTo(controlX, controlY, x, y float64) *Path {
	if p != nil {
		p.commands = append(p.commands, pathCommand{
			operation: pathQuadratic,
			points:    [6]float64{controlX, controlY, x, y},
		})
	}
	return p
}

// CubicTo adds a cubic Bézier curve with two control points.
func (p *Path) CubicTo(control1X, control1Y, control2X, control2Y, x, y float64) *Path {
	if p != nil {
		p.commands = append(p.commands, pathCommand{
			operation: pathCubic,
			points:    [6]float64{control1X, control1Y, control2X, control2Y, x, y},
		})
	}
	return p
}

// Close connects the current point to the beginning of its path section.
func (p *Path) Close() *Path {
	if p != nil {
		p.commands = append(p.commands, pathCommand{operation: pathClose})
	}
	return p
}

func (p *Path) replay(context *gg.Context) {
	if p == nil {
		return
	}
	for _, command := range p.commands {
		switch command.operation {
		case pathMove:
			context.MoveTo(command.points[0], command.points[1])
		case pathLine:
			context.LineTo(command.points[0], command.points[1])
		case pathQuadratic:
			context.QuadraticTo(
				command.points[0], command.points[1],
				command.points[2], command.points[3],
			)
		case pathCubic:
			context.CubicTo(
				command.points[0], command.points[1],
				command.points[2], command.points[3],
				command.points[4], command.points[5],
			)
		case pathClose:
			context.ClosePath()
		}
	}
}
