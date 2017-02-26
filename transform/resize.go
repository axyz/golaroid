package transform

import (
	"fmt"

	"gopkg.in/gographics/imagick.v2/imagick"
)

var filters = map[string]imagick.FilterType{
	"lanczos": imagick.FILTER_LANCZOS,
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func getOutputDimensions(iw, ih, ow, oh float64) (uint, uint) {
	scaleW := ow / iw
	scaleH := oh / ih

	if oh == 0 {
		return uint(iw * scaleW), uint(ih * scaleW)
	} else if ow == 0 {
		return uint(iw * scaleH), uint(ih * scaleH)
	}

	scale := min(scaleW, scaleH)

	return uint(iw * scale), uint(ih * scale)
}

func Resize(data []byte, opts map[string]interface{}) ([]byte, error) {
	filter := opts["filter"].(string)
	quality := opts["quality"].(int)
	width, ok := opts["width"].(int)
	if !ok {
		width = 0
	}
	height, ok := opts["height"].(int)
	if !ok {
		height = 0
	}
	fmt.Println("TODO: resizing...")
	fmt.Printf("%+v\n", opts)
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	mw.ReadImageBlob(data)

	iw := mw.GetImageWidth()
	ih := mw.GetImageHeight()

	ow, oh := getOutputDimensions(float64(iw), float64(ih), float64(width), float64(height))
	mw.ResizeImage(ow, oh, filters[filter], 1)
	mw.SetImageCompressionQuality(uint(quality))
	output := mw.GetImageBlob()

	return output, nil
}
