package gioassets

import (
	"gioui.org/op/paint"
	"gioui.org/text"
	"github.com/inkeliz/gioassets/fsfont"
	"github.com/inkeliz/gioassets/fsimage"
	"github.com/inkeliz/gioassets/fsvector"
	"github.com/inkeliz/giosvg"
	"io/fs"
)

// NewSharper creates a text.Shaper.
// It can be re-used across all widget.Label or widget.Editor.
//
// The name of the file must be `{name}_{weight}[_{style}]`, for instance
// it can be Montserrat-700.ttf or Montserrat-700-Italic.ttf.
func NewSharper(embed fs.FS) text.Shaper {
	return fsfont.NewMust(embed)
}

// NewImages returns a list of `paint.ImageOp`.
// It expects the image to be WEBP, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func NewImages(embed fs.FS) map[string]paint.ImageOp {
	return fsimage.NewMust(embed)
}

// NewVectors returns a list of `*giosvg.IconOp`.
// It expects the image to be SVG, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func NewVectors(embed fs.FS) map[string]*giosvg.IconOp {
	return fsvector.NewMust(embed)
}
