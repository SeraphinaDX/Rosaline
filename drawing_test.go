// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"image/color"
	"testing"
)

func TestRenderDrawsOffScreen(t *testing.T) {
	picture := Render(80, 60, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.FillRect(10, 12, 20, 18, Rose)
	})

	if picture.Width() != 80 || picture.Height() != 60 {
		t.Fatalf("picture size = %dx%d, want 80x60", picture.Width(), picture.Height())
	}
	assertPixel(t, picture, 15, 18, nativeColor(Rose))
	assertPixel(t, picture, 2, 2, nativeColor(White))
}

func TestPathSupportsBezierCurves(t *testing.T) {
	heart := NewPath().
		MoveTo(40, 58).
		CubicTo(8, 40, 10, 12, 40, 28).
		CubicTo(70, 12, 72, 40, 40, 58).
		Close()

	picture := Render(80, 70, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.FillPath(heart, Rose)
	})

	assertPixel(t, picture, 40, 36, nativeColor(Rose))
	assertPixel(t, picture, 2, 2, nativeColor(White))
}

func TestTransformsMoveDrawing(t *testing.T) {
	picture := Render(80, 60, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.Translate(30, 20)
		canvas.FillRect(0, 0, 12, 10, Rose)
	})

	assertPixel(t, picture, 35, 25, nativeColor(Rose))
	assertPixel(t, picture, 5, 5, nativeColor(White))
}

func TestScaleAndRotationTransformDrawing(t *testing.T) {
	picture := Render(80, 70, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.Translate(40, 20)
		canvas.Rotate(90)
		canvas.Scale(2, 1)
		canvas.FillRect(0, 0, 10, 16, Rose)
	})

	assertPixel(t, picture, 30, 30, nativeColor(Rose))
	assertPixel(t, picture, 50, 30, nativeColor(White))
}

func TestClipFollowsCurrentTransform(t *testing.T) {
	picture := Render(70, 60, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.Translate(20, 12)
		canvas.ClipRect(0, 0, 20, 18)
		canvas.ResetTransform()
		canvas.FillRect(0, 0, 70, 60, Rose)
	})

	assertPixel(t, picture, 25, 18, nativeColor(Rose))
	assertPixel(t, picture, 10, 10, nativeColor(White))
	assertPixel(t, picture, 50, 40, nativeColor(White))
}

func TestClipAndPopRestorePreviousRegion(t *testing.T) {
	picture := Render(80, 60, func(canvas *DrawingCanvas) {
		canvas.Clear(White)
		canvas.ClipRect(0, 0, 40, 60)
		canvas.Push()
		canvas.ClipRect(0, 0, 40, 20)
		canvas.FillRect(0, 0, 80, 60, SoftRose)
		canvas.Pop()
		canvas.FillRect(0, 30, 80, 20, Rose)
	})

	assertPixel(t, picture, 10, 10, nativeColor(SoftRose))
	assertPixel(t, picture, 10, 40, nativeColor(Rose))
	assertPixel(t, picture, 60, 10, nativeColor(White))
	assertPixel(t, picture, 60, 40, nativeColor(White))
}

func TestCanvasPictureWorksBeforeMount(t *testing.T) {
	canvas := Canvas(func(drawing *DrawingCanvas) {
		drawing.Clear(White)
		drawing.FillCircle(20, 20, 10, Rose)
	}).Size(50, 40)

	picture := canvas.Picture()
	if picture.Width() != 50 || picture.Height() != 40 {
		t.Fatalf("picture size = %dx%d, want 50x40", picture.Width(), picture.Height())
	}
	assertPixel(t, picture, 20, 20, nativeColor(Rose))
}

func assertPixel(t *testing.T, picture *Picture, x, y int, want color.NRGBA) {
	t.Helper()
	got := color.NRGBAModel.Convert(picture.Image().At(x, y)).(color.NRGBA)
	if got != want {
		t.Fatalf("pixel (%d,%d) = %#v, want %#v", x, y, got, want)
	}
}
