package jimg

import (
	"image"
	"image/png"
	"io"
)

type jPng struct {
	JPngInfo
	yPngMaker
}
type JPngInfo struct{}

func (JPngInfo) GetTag() string {
	return "data:image/png;base64,"
}
func (JPngInfo) GetExtension() string {
	return ".png"
}

type yPngMaker struct{}

func (yPngMaker) decode(decoder io.Reader) (image.Image, error) {
	return png.Decode(decoder)
}
func (yPngMaker) encode(f io.Writer, decodeImage image.Image) error {
	return png.Encode(f, decodeImage)
}
