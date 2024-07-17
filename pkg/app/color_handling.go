package app

import (
	"image/color"
	"math"
)

func normalizeC(c uint32) float64 {
	rc := float64(c) / float64(0xffff)
	if rc <= 0.04045 {
		rc = rc / 12.92
	} else {
		rc = math.Exp(math.Log((rc+0.055)/1.055) * 2.4)
		// rc = math.Pow(, 2.4)
	}

	return rc
}

// https://stackoverflow.com/a/3943023
func UseDarkText(bgc color.Color) bool {
	r, g, b, _ := bgc.RGBA()

	rc, gc, bc := normalizeC(r), normalizeC(g), normalizeC(b)
	l := 0.2126*rc + 0.7152*gc + 0.0722*bc

	return l > 0.179
}
