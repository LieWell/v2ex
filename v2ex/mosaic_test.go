package v2ex

import "testing"

func TestCreate(t *testing.T) {
	// 马赛克图像
	targetImagePath := "/Users/sheng/Documents/Code/liewell/v2ex/target.png"
	tileImagesDir := "/Users/sheng/Documents/Code/liewell/v2ex/tmp/"
	outputImagePath := "mosaic_output.jpg"
	tileSize := 10

	mosaic := CreateMosaic(targetImagePath, tileImagesDir, tileSize)
	SaveImage(mosaic, outputImagePath)
}
