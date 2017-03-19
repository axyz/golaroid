package transform

import (
	"fmt"
	"math"

	"gopkg.in/gographics/imagick.v2/imagick"
)

var filters = map[string]imagick.FilterType{
	"bartlett":       imagick.FILTER_BARTLETT,
	"bohman":         imagick.FILTER_BOHMAN,
	"blackman":       imagick.FILTER_BLACKMAN,
	"box":            imagick.FILTER_BOX,
	"catrom":         imagick.FILTER_CATROM,
	"cosine":         imagick.FILTER_COSINE,
	"cubic":          imagick.FILTER_CUBIC,
	"gaussian":       imagick.FILTER_GAUSSIAN,
	"hamming":        imagick.FILTER_HAMMING,
	"hanning":        imagick.FILTER_HANNING,
	"hermite":        imagick.FILTER_HERMITE,
	"jinc":           imagick.FILTER_JINC,
	"kaiser":         imagick.FILTER_KAISER,
	"lagrange":       imagick.FILTER_LAGRANGE,
	"lanczos":        imagick.FILTER_LANCZOS,
	"lanczos-radius": imagick.FILTER_LANCZOS_RADIUS,
	"lanczos-sharp":  imagick.FILTER_LANCZOS_SHARP,
	"lanczos2":       imagick.FILTER_LANCZOS2,
	"lanczos2-sharp": imagick.FILTER_LANCZOS2_SHARP,
	"mitchell":       imagick.FILTER_MITCHELL,
	"parzen":         imagick.FILTER_PARZEN,
	"point":          imagick.FILTER_POINT,
	"quadratic":      imagick.FILTER_QUADRATIC,
	"rubidoux":       imagick.FILTER_ROBIDOUX,
	"rubidoux-sharp": imagick.FILTER_ROBIDOUX_SHARP,
	"sinc":           imagick.FILTER_SINC,
	"triangle":       imagick.FILTER_TRIANGLE,
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}

	return b
}

func max(a, b float64) float64 {
	if a >= b {
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

func calculateSteps(iw, ih, ow, oh, ratio float64, maxSteps uint) uint {
	steps := uint(max(iw, ih) / max(ow, oh) * ratio)
	if steps > maxSteps {
		return maxSteps
	}

	return steps
}

func multistepResize(mw *imagick.MagickWand, iw, ih, ow, oh float64, opts map[string]interface{}) {
	filter, ok := opts["filter"].(string)
	if !ok {
		filter = "lanczos"
	}
	ratio, ok := opts["stepRatio"].(float64)
	if !ok {
		ratio = 0.66
	}
	maxSteps, ok := opts["maxSteps"].(uint)
	if !ok {
		maxSteps = 4
	}
	var unsharpStepParams [4]float64
	parseUnsharpConfig(&unsharpStepParams, opts, "unsharpStepParams", [4]float64{2, 0.5, 0.5, 0.008})
	steps := calculateSteps(iw, ih, ow, oh, ratio, maxSteps)

	for i := uint(1); i <= steps; i++ {
		scale := math.Pow(ratio, float64(i))
		mw.ResizeImage(uint(iw*scale), uint(ih*scale), filters[filter], 1)
		mw.UnsharpMaskImage(
			unsharpStepParams[0],
			unsharpStepParams[1],
			unsharpStepParams[2],
			unsharpStepParams[3],
		)
	}

	mw.ResizeImage(uint(ow), uint(oh), filters[filter], 1)
}

func parseUnsharpConfig(dst *[4]float64, opts map[string]interface{}, key string, fallback [4]float64) {
	unsharpParamsArray, ok := opts[key].([]interface{})
	if !ok {
		*dst = fallback
	} else {
		var radius, sigma, gain, threshold float64
		switch p0 := unsharpParamsArray[0].(type) {
		case int:
			radius = float64(p0)
		case float64:
			radius = p0
		}
		switch p1 := unsharpParamsArray[1].(type) {
		case int:
			sigma = float64(p1)
		case float64:
			sigma = p1
		}
		switch p2 := unsharpParamsArray[2].(type) {
		case int:
			gain = float64(p2)
		case float64:
			gain = p2
		}
		switch p3 := unsharpParamsArray[3].(type) {
		case int:
			threshold = float64(p3)
		case float64:
			threshold = p3
		}

		*dst = [4]float64{radius, sigma, gain, threshold}
	}

}

func Resize(data []byte, opts map[string]interface{}) ([]byte, error) {
	filter, ok := opts["filter"].(string)
	if !ok {
		filter = "lanczos"
	}
	width, ok := opts["width"].(int)
	if !ok {
		width = 0
	}
	height, ok := opts["height"].(int)
	if !ok {
		height = 0
	}
	unsharp, ok := opts["height"].(bool)
	if !ok {
		unsharp = true
	}
	var unsharpParams [4]float64
	parseUnsharpConfig(&unsharpParams, opts, "unsharpParams", [4]float64{2, 0.5, 0.5, 0.006})
	multistep, ok := opts["multistep"].(bool)
	if !ok {
		multistep = false
	}

	fmt.Println("TODO: resizing...")
	fmt.Printf("%+v\n", opts)
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()

	mw.ReadImageBlob(data)

	iw := mw.GetImageWidth()
	ih := mw.GetImageHeight()

	ow, oh := getOutputDimensions(
		float64(iw),
		float64(ih),
		float64(width),
		float64(height),
	)

	if multistep {
		multistepResize(
			mw,
			float64(iw),
			float64(ih),
			float64(ow),
			float64(oh),
			opts,
		)
	} else {
		mw.ResizeImage(ow, oh, filters[filter], 1)
	}

	if unsharp {
		mw.UnsharpMaskImage(
			unsharpParams[0],
			unsharpParams[1],
			unsharpParams[2],
			unsharpParams[3],
		)
	}

	output := mw.GetImageBlob()

	return output, nil
}
