// game/fonts.go
package game

import (
	_ "embed"
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed assets/fonts/RobotoMono-Regular.ttf
var robotoMonoRegular []byte

const (
	dpi = 72

	SmallFontSize  = 12
	MediumFontSize = 16
	LargeFontSize  = 24
)

type Fonts struct {
	Small  font.Face
	Medium font.Face
	Large  font.Face
}

func LoadFonts() (*Fonts, error) {
	tt, err := opentype.Parse(robotoMonoRegular)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	smallFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    SmallFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create small font face: %v", err)
	}

	mediumFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    MediumFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create medium font face: %v", err)
	}

	largeFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    LargeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create large font face: %v", err)
	}

	return &Fonts{
		Small:  smallFont,
		Medium: mediumFont,
		Large:  largeFont,
	}, nil
}
