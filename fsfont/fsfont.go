package fsfont

import (
	"errors"
	"gioui.org/font/opentype"
	"gioui.org/text"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
)

var defaultWeights = map[string]text.Weight{
	strings.ToLower("Thin"):       text.Thin,
	strings.ToLower("Hairline"):   text.Hairline,
	strings.ToLower("ExtraLight"): text.ExtraLight,
	strings.ToLower("UltraLight"): text.UltraLight,
	strings.ToLower("Light"):      text.Light,
	strings.ToLower("Normal"):     text.Normal,
	strings.ToLower("Medium"):     text.Medium,
	strings.ToLower("SemiBold"):   text.SemiBold,
	strings.ToLower("DemiBold"):   text.DemiBold,
	strings.ToLower("Bold"):       text.Bold,
	strings.ToLower("ExtraBold"):  text.ExtraBold,
	strings.ToLower("UltraBold"):  text.UltraBold,
	strings.ToLower("Black"):      text.Black,
	strings.ToLower("Heavy"):      text.Heavy,
	strings.ToLower("ExtraBlack"): text.ExtraBlack,
	strings.ToLower("UltraBlack"): text.UltraBlack,
}

// New creates a text.Shaper.
// It can be re-used across all widget.Label or widget.Editor.
//
// The name of the file must be `{name}_{weight}[_{style}]`, for instance
// it can be Montserrat-700.ttf or Montserrat-700-Italic.ttf.
func New(embed fs.FS) (text.Shaper, error) {
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
			weight = text.Weight(w)
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

	return text.NewCache(fonts), nil
}

// NewMust uses New and panic when error is returned.
func NewMust(embed fs.FS) text.Shaper {
	r, err := New(embed)
	if err != nil {
		panic(err)
	}
	return r
}
