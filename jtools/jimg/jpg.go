package jimg

import (
	"image"
	"image/jpeg"
	"io"
)

type jJpg struct {
	JJpgInfo
	yJpgMaker
}
type JJpgInfo struct{}

func (JJpgInfo) GetTag() string {
	return "data:image/jpeg;base64,"
}
func (JJpgInfo) GetExtension() string {
	return ".jpg"
}

type yJpgMaker struct{}

func (yJpgMaker) decode(decoder io.Reader) (image.Image, error) {
	return jpeg.Decode(decoder)
}
func (yJpgMaker) encode(f io.Writer, decodeImage image.Image) error {
	return jpeg.Encode(f, decodeImage, nil)
}
