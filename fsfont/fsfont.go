package fsfont

import (
	"errors"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

var defaultWeights = map[string]font.Weight{
	strings.ToLower("Thin"):       font.Thin,
	strings.ToLower("ExtraLight"): font.ExtraLight,
	strings.ToLower("Light"):      font.Light,
	strings.ToLower("Normal"):     font.Normal,
	strings.ToLower("Medium"):     font.Medium,
	strings.ToLower("SemiBold"):   font.SemiBold,
	strings.ToLower("Bold"):       font.Bold,
	strings.ToLower("ExtraBold"):  font.ExtraBold,
	strings.ToLower("Black"):      font.Black,
}

// New creates a text.Shaper.
// It can be re-used across all widget.Label or widget.Editor.
//
// The name of the file must be `{name}_{weight}[_{style}]`, for instance
// it can be Montserrat-700.ttf or Montserrat-700-Italic.ttf.
func New(embed fs.FS) (*text.Shaper, error) {
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

		weight, ok := defaultWeights[strings.ToLower(split[1])]
		if !ok {
			w, err := strconv.ParseInt(split[1], 10, 64)
			if err != nil {
				return errors.New("invalid font name, it must be fontName_{weight}_[{style}].ttf, for instance: Montserrat-700.ttf or Montserrat-700-Italic.ttf")
			}
			weight = font.Weight(w)
		}

		style := font.Regular
		if len(split) >= 3 {
			if strings.ToLower(split[2]) == "italic" {
				style = font.Italic
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
			Font: font.Font{
				Typeface: font.Typeface(name),
				Weight:   weight,
				Style:    style,
			},
			Face: face,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(fonts); i++ {
		if fonts[i].Font.Weight == 300 {
			fonts[0], fonts[i] = fonts[i], fonts[0]
			break
		}
	}

	return text.NewShaper(text.WithCollection(fonts), text.NoSystemFonts()), nil
}

// NewMust uses New and panic when error is returned.
func NewMust(embed fs.FS) *text.Shaper {
	r, err := New(embed)
	if err != nil {
		panic(err)
	}
	return r
}
