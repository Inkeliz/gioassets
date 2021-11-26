package fsimage

import (
	"gioui.org/op/paint"
	"golang.org/x/image/webp"
	"image"
	"io/fs"
	"path/filepath"
)

// New returns a list of `paint.ImageOp`.
// It expects the image to be WEBP, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func New(embed fs.FS) (map[string]paint.ImageOp, error) {
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
		return nil, err
	}

	return images, nil
}

// NewMust uses New and panic when error is returned.
func NewMust(embed fs.FS) map[string]paint.ImageOp {
	r, err := New(embed)
	if err != nil {
		panic(err)
	}
	return r
}

