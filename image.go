// SPDX-License-Identifier: LGPL-3.0-or-later

package rosaline

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"math"
	"os"

	_ "github.com/gen2brain/avif"
	_ "golang.org/x/image/bmp"
	xdraw "golang.org/x/image/draw"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	tk "modernc.org/tk9.0"
)

// Picture contains a decoded image that Rosaline can display.
type Picture struct {
	pixels image.Image
	path   string
	format string
}

// NewPicture creates a Rosaline picture from Go's standard image.Image type.
func NewPicture(pixels image.Image) *Picture {
	return &Picture{pixels: pixels}
}

// LoadImage reads and decodes an image file.
// PNG, JPEG, GIF, BMP, TIFF, WebP, and AVIF are supported.
func LoadImage(path string) (*Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("rosaline: could not open image %q: %w", path, err)
	}
	defer file.Close()
	return decodePicture(file, path)
}

// LoadImageFS reads an image from a standard Go filesystem such as embed.FS.
// It is useful for applications that ship pictures inside their executable.
func LoadImageFS(fileSystem fs.FS, path string) (*Picture, error) {
	if fileSystem == nil {
		return nil, fmt.Errorf("rosaline: could not open image %q: filesystem is nil", path)
	}
	file, err := fileSystem.Open(path)
	if err != nil {
		return nil, fmt.Errorf("rosaline: could not open image %q: %w", path, err)
	}
	defer file.Close()
	return decodePicture(file, path)
}

func decodePicture(reader io.Reader, path string) (*Picture, error) {
	pixels, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("rosaline: could not decode image %q: %w", path, err)
	}

	return &Picture{pixels: pixels, path: path, format: format}, nil
}

// Width returns the picture width in pixels.
func (p *Picture) Width() int {
	if p == nil || p.pixels == nil {
		return 0
	}
	return p.pixels.Bounds().Dx()
}

// Height returns the picture height in pixels.
func (p *Picture) Height() int {
	if p == nil || p.pixels == nil {
		return 0
	}
	return p.pixels.Bounds().Dy()
}

// Path returns the filename used by LoadImage. Pictures made with NewPicture
// have an empty path.
func (p *Picture) Path() string {
	if p == nil {
		return ""
	}
	return p.path
}

// Format returns the decoded format name, such as "png" or "jpeg".
func (p *Picture) Format() string {
	if p == nil {
		return ""
	}
	return p.format
}

// Image returns the underlying standard-library image.Image value.
func (p *Picture) Image() image.Image {
	if p == nil {
		return nil
	}
	return p.pixels
}

// ImageWidget displays a Picture.
type ImageWidget struct {
	picture     *Picture
	placeholder string
	fitWidth    int
	fitHeight   int
	expand      bool
	onClick     func()
	label       *tk.LabelWidget
	tkImage     *tk.Img
}

// Image creates a widget that displays picture. A nil picture is allowed and
// shows a friendly placeholder until SetImage is called.
func Image(picture *Picture) *ImageWidget {
	return &ImageWidget{picture: picture, placeholder: "No image loaded."}
}

// Placeholder changes the text shown when no picture is loaded.
func (i *ImageWidget) Placeholder(text string) *ImageWidget {
	i.placeholder = text
	return i
}

// Fit displays the picture inside a fixed pixel area while preserving its
// aspect ratio. Any unused area remains transparent and shows the theme's
// surface color. Non-positive dimensions restore the picture's natural size.
func (i *ImageWidget) Fit(width, height int) *ImageWidget {
	if i == nil {
		return i
	}
	if width <= 0 || height <= 0 {
		i.fitWidth, i.fitHeight = 0, 0
	} else {
		i.fitWidth, i.fitHeight = width, height
	}
	i.apply()
	return i
}

// OnClick runs when the user clicks the image.
func (i *ImageWidget) OnClick(handler func()) *ImageWidget {
	if i != nil {
		i.onClick = handler
	}
	return i
}

// Expand asks the image widget to use available layout space.
func (i *ImageWidget) Expand() *ImageWidget {
	i.expand = true
	return i
}

// SetImage changes the displayed picture. It can be called from Rosaline
// callbacks after the widget has been mounted.
func (i *ImageWidget) SetImage(picture *Picture) {
	i.picture = picture
	i.apply()
}

// Picture returns the picture currently displayed by the widget.
func (i *ImageWidget) Picture() *Picture {
	return i.picture
}

func (i *ImageWidget) apply() {
	if i.label == nil {
		return
	}

	oldImage := i.tkImage
	if i.picture == nil || i.picture.pixels == nil {
		i.tkImage = nil
		i.label.Configure(tk.Image(""), tk.Txt(i.placeholder))
	} else {
		pixels := i.picture.pixels
		if i.fitWidth > 0 && i.fitHeight > 0 {
			pixels = fitImage(pixels, i.fitWidth, i.fitHeight)
		}
		i.tkImage = tk.NewPhoto(tk.Data(pixels))
		i.label.Configure(tk.Image(i.tkImage), tk.Txt(""))
	}
	if oldImage != nil {
		oldImage.Delete()
	}
}

func (i *ImageWidget) mount(ctx *mountContext, parent *tk.Window) mountedWidget {
	i.label = parent.Label(
		tk.Background(ctx.theme.Surface.String()),
		tk.Foreground(ctx.theme.Muted.String()),
		tk.Anchor("center"),
		tk.Borderwidth(0),
	)
	i.apply()
	if i.onClick != nil {
		tk.Bind(i.label.Window, "<Button-1>", tk.Command(func() {
			ctx.flush()
			i.onClick()
			ctx.refresh()
		}))
	}
	ctx.addCleanup(func() {
		if i.tkImage != nil {
			i.tkImage.Delete()
		}
		i.tkImage = nil
		i.label = nil
	})
	return mountedWidget{window: i.label.Window, expandX: i.expand, expandY: i.expand}
}

func fitImage(source image.Image, width, height int) image.Image {
	if source == nil || width <= 0 || height <= 0 {
		return source
	}
	box := source.Bounds()
	if box.Dx() <= 0 || box.Dy() <= 0 {
		return source
	}
	scale := math.Min(float64(width)/float64(box.Dx()), float64(height)/float64(box.Dy()))
	scaledWidth := max(1, int(math.Round(float64(box.Dx())*scale)))
	scaledHeight := max(1, int(math.Round(float64(box.Dy())*scale)))
	x := (width - scaledWidth) / 2
	y := (height - scaledHeight) / 2
	destination := image.NewRGBA(image.Rect(0, 0, width, height))
	xdraw.CatmullRom.Scale(
		destination,
		image.Rect(x, y, x+scaledWidth, y+scaledHeight),
		source,
		box,
		xdraw.Over,
		nil,
	)
	return destination
}
