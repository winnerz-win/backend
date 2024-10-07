package excel

import "encoding/hex"

const (
	RED      = "ff0000"
	RED_DARK = "CC3D3D"

	GREEN      = "00ff00"
	GREEN_DARK = "47C83E"

	LEAF      = "ABF200"
	LEAF_DARK = "9FC93C"

	BLUE      = "0000ff"
	BLUE_DARK = "4641D9"

	PINK      = "FF00FF"
	PINK_DARK = "D9418C"

	YELLOW      = "FFFF00"
	YELLOW_DAKR = "C4B73B"

	CYAN      = "00D8FF"
	CYAN_DARK = "3DB7CC"

	ORANGE      = "FFBB00"
	ORANGE_DARK = "CCA63D"

	GRAY       = "747474"
	BLACK_DARK = "191919"
)

func RGB(r, g, b byte) string {
	return hex.EncodeToString([]byte{r, g, b})
}
