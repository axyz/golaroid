package transform

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

func Jpeg(data []byte, opts map[string]interface{}) ([]byte, error) {
	quality, ok := opts["quality"].(int)
	if !ok {
		quality = 85
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	mw.ReadImageBlob(data)

	mw.SetImageCompressionQuality(uint(quality))
	output := mw.GetImageBlob()

	return output, nil
}
