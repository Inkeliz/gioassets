package gioassets

import (
	"errors"
	"gioui.org/font/opentype"
	"gioui.org/op/paint"
	"gioui.org/text"
	"github.com/inkeliz/giosvg"
	"golang.org/x/image/webp"
	"image"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

// NewSharper creates a text.Shaper.
// It can be re-used across all widget.Label or widget.Editor.
//
// The name of the file must be `{name}_{weight}[_{style}]`, for instance
// it can be Montserrat-700.ttf or Montserrat-700-Italic.ttf.
func NewSharper(embed fs.FS) text.Shaper {
	fonts := make([]text.FontFace, 0, 16)

	err := fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(d.Name()) != ".ttf" {
			return nil
		}

		split := strings.Split(strings.Replace(strings.TrimRight(d.Name(), ".ttf"), "-", "_", -1), "_")
		if len(split) < 2 {
			return errors.New("invalid font name, it must be fontName_{weight}_[{style}].ttf, for instance: Montserrat-700.ttf or Montserrat-700-Italic.ttf")
		}

		name := split[0]
		weight, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return errors.New("invalid font name, it must be fontName_{weight}_[{style}].ttf, for instance: Montserrat-700.ttf or Montserrat-700-Italic.ttf")
		}
		style := text.Regular
		if len(split) >= 3 {
			if strings.ToLower(split[2]) == "italic" {
				style = text.Italic
			}
		}

		file, err := fs.ReadFile(embed, path)
		if err != nil {
			return err
		}

		face, err := opentype.Parse(file)
		if err != nil {
			return err
		}

		fonts = append(fonts, text.FontFace{
			Font: text.Font{
				Typeface: text.Typeface(name),
				Weight:   text.Weight(weight),
				Style:    style,
			},
			Face: face,
		})

		return nil
	})

	if err != nil {
		panic(err)
	}

	return text.NewCache(fonts)
}

// NewImages returns a list of `paint.ImageOp`.
// It expects the image to be WEBP, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func NewImages(embed fs.FS) map[string]paint.ImageOp {
	images := make(map[string]paint.ImageOp, 16)

	err := fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(d.Name()) != ".webp" {
			return nil
		}

		fileReader, err := embed.Open(path)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		img, err := webp.Decode(fileReader)
		if err != nil {
			return err
		}

		// Make it faster for WASM
		switch src := img.(type) {
		case *image.NRGBA:
			img = (*image.RGBA)(src)
		}

		name := d.Name()
		name = name[:len(name)-5]
		images[name] = paint.NewImageOp(img)
		return nil
	})

	if err != nil {
		panic(err)
	}

	return images
}

// NewVectors returns a list of `*giosvg.IconOp`.
// It expects the image to be SVG, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func NewVectors(embed fs.FS) map[string]*giosvg.IconOp {
	images := make(map[string]*giosvg.IconOp, 16)

	err := fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(d.Name()) != ".svg" {
			return nil
		}

		fileReader, err := embed.Open(path)
		if err != nil {
			return err
		}
		defer fileReader.Close()

		name := d.Name()
		name = name[:len(name)-4]
		images[name], err = giosvg.NewIconOpReader(fileReader)
		return err
	})

	if err != nil {
		panic(err)
	}

	return images
}
