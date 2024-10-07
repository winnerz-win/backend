package jimg

import (
	"encoding/base64"
	"errors"
	"image"
	"io"
	"jtools/jcrypt/uuid"
	"jtools/jparallel"
	"os"
	"strings"
)

type iJImg interface {
	IJImgInfo
	iJImgMaker
}
type IJImgInfo interface {
	GetTag() string
	GetExtension() string
}
type iJImgMaker interface {
	decode(io.Reader) (image.Image, error)
	encode(io.Writer, image.Image) error
}

var jImgs = map[string]iJImg{
	JPngInfo{}.GetTag(): jPng{},
	JJpgInfo{}.GetTag(): jJpg{},
	JBmpInfo{}.GetTag(): jBmp{},
}
var jImgInfos map[string]IJImgInfo

func GetImgInfo(imgTag string) IJImgInfo {
	return jImgs[imgTag]
}
func GetImgInfoAll() map[string]IJImgInfo {
	if jImgInfos != nil {
		return jImgInfos
	}
	jImgInfos = map[string]IJImgInfo{}
	for tag, img := range jImgs {
		jImgInfos[tag] = img
	}
	return jImgInfos
}

type Base64ToImgParams []Base64ToImgParam
type Base64ToImgParam struct {
	ImgBase64 string `json:"img_base64"`
	ImgPath   string `json:"img_path"`
	ImgName   string `json:"img_name"`
	IJImgInfo
}

func MakeBase64ToImgParam(imgBase64, imgPath, imgName string) Base64ToImgParam {
	imgTag := GetTag(imgBase64)
	imgInfo := GetImgInfo(imgTag)
	return Base64ToImgParam{
		imgBase64,
		imgPath,
		imgName,
		imgInfo,
	}
}
func (my Base64ToImgParam) GetImgPathAndNameAndExtension() string {
	return my.ImgPath + my.GetImgNameAndExtension()
}
func (my Base64ToImgParam) GetImgNameAndExtension() string {
	return my.ImgName + my.GetExtension()
}

func Base64ToImgFile(param Base64ToImgParam) error {
	imgMaker := jImgs[param.GetTag()]
	if imgMaker == nil {
		return errors.New("param.ImgBase64 tag is invalid")
	}

	imgBase64 := param.ImgBase64
	if strings.Contains(param.ImgBase64, imgMaker.GetTag()) {
		imgBase64 = strings.Split(param.ImgBase64, imgMaker.GetTag())[1]
	}

	decoder := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imgBase64))
	decodeImage, err := imgMaker.decode(decoder)
	if err != nil {
		return err
	}
	err = os.MkdirAll(param.ImgPath, 0755)
	if err != nil {
		return err
	}
	f, err := os.Create(param.GetImgPathAndNameAndExtension())
	if err != nil {
		return err
	}
	defer f.Close()
	err = imgMaker.encode(f, decodeImage)
	if err != nil {
		return err
	}

	return nil
}

func Base64ToImgFiles(params Base64ToImgParams) error {
	_, errors, _ := jparallel.Foreach(
		params,
		func(i int, param Base64ToImgParam) (string, error) {
			err := Base64ToImgFile(param)
			if err != nil {
				return "", err
			}
			return "", nil
		},
		200,
	)
	return errors.Error()
}

func MakeImgName(separator string, name string, isAppendUuid bool, appendEtcs ...string) (string, error) {
	b := strings.Builder{}

	b.WriteString(name)

	for _, appendEtc := range appendEtcs {
		b.WriteString(separator)
		b.WriteString(appendEtc)
	}

	if isAppendUuid {
		uuid, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		b.WriteString(separator)
		b.WriteString(uuid.String())
	}

	return b.String(), nil
}

func GetTag(imgBase64 string) string {
	for tag := range jImgs {
		if strings.Contains(imgBase64, tag) {
			return tag
		}
	}
	return ""
}

func HasImg(imgBase64 string) bool {
	return GetTag(imgBase64) != ""
}
