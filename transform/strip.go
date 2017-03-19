package transform

import (
	"fmt"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func Strip(data []byte, opts map[string]interface{}) ([]byte, error) {
	fmt.Println("TODO: stripping...")

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	mw.ReadImageBlob(data)

	mw.StripImage()
	output := mw.GetImageBlob()

	return output, nil
}
