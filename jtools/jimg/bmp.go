package jimg

import (
	"image"
	"io"

	"golang.org/x/image/bmp"
)

type jBmp struct {
	JBmpInfo
	yBmpMaker
}
type JBmpInfo struct{}

func (JBmpInfo) GetTag() string {
	return "data:image/bmp;base64,"
}
func (JBmpInfo) GetExtension() string {
	return ".bmp"
}

type yBmpMaker struct{}

func (yBmpMaker) decode(decoder io.Reader) (image.Image, error) {
	return bmp.Decode(decoder)
}
func (yBmpMaker) encode(f io.Writer, decodeImage image.Image) error {
	return bmp.Encode(f, decodeImage)
}
