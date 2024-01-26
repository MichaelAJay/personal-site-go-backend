package custom_utility

import (
	"fmt"
	"image/color"
	"strconv"
)

func ParseSixDigitHexCodeToColor(s string) (color.Color, error) {
	if s[0] != '#' {
		return nil, fmt.Errorf("invalid color string: missing '#'")
	}

	if len(s) != 7 {
		return nil, fmt.Errorf("invalid color string:  should be in format '#rrggbb'")
	}

	var r, g, b uint8
	var err error

	hex := s[1:]
	r, err = parseHexByte(hex[0:2])
	if err != nil {
		return nil, err
	}
	g, err = parseHexByte(hex[2:4])
	if err != nil {
		return nil, err
	}
	b, err = parseHexByte(hex[4:6])
	if err != nil {
		return nil, err
	}

	return color.RGBA{R: r, G: g, B: b, A: 0xff}, nil
}

func parseHexByte(s string) (uint8, error) {
	b, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		return 0, err
	}
	return uint8(b), nil
}
