package fsvector

import (
	"github.com/inkeliz/giosvg"
	"io/fs"
	"path/filepath"
)

// New returns a list of `*giosvg.IconOp`.
// It expects the image to be SVG, and the name of the file must
// be unique, even if the file is storage on distinct folders.
//
// The name of the file doesn't contain the extension.
func New(embed fs.FS) (map[string]*giosvg.IconOp, error) {
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
		return nil, err
	}

	return images, nil
}

// NewMust uses New and panic when error is returned.
func NewMust(embed fs.FS) map[string]*giosvg.IconOp {
	r, err := New(embed)
	if err != nil {
		panic(err)
	}
	return r
}
