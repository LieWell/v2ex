// 马赛克图像

package v2ex

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

func averageColor(img image.Image) color.Color {
	var r, g, b, a uint32
	var count uint32
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rr, gg, bb, aa := img.At(x, y).RGBA()
			r += rr
			g += gg
			b += bb
			a += aa
			count++
		}
	}
	if count == 0 {
		return color.RGBA{0, 0, 0, 0}
	}
	return color.RGBA{uint8(r / count / 256), uint8(g / count / 256), uint8(b / count / 256), uint8(a / count / 256)}
}

func colorDistance(c1, c2 color.Color) float64 {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	return math.Sqrt(float64((r1-r2)*(r1-r2) + (g1-g2)*(g1-g2) + (b1-b2)*(b1-b2)))
}

func findBestMatch(avgColor color.Color, tileColors []color.Color) int {
	minDistance := math.MaxFloat64
	bestIndex := 0 // 初始化为0，保证总能返回有效索引
	for i, color := range tileColors {
		distance := colorDistance(avgColor, color)
		if distance < minDistance {
			minDistance = distance
			bestIndex = i
		}
	}
	return bestIndex
}

func CreateMosaic(targetImagePath string, tileImagesDir string, tileSize int) image.Image {
	targetFile, err := os.Open(targetImagePath)
	if err != nil {
		log.Fatalf("failed to open target image: %v", err)
	}
	defer targetFile.Close()

	targetImage, _, err := image.Decode(targetFile)
	if err != nil {
		log.Fatalf("failed to decode target image: %v", err)
	}

	var tileImages []image.Image
	var tileColors []color.Color

	err = filepath.Walk(tileImagesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		log.Printf("Loading tile image: %s", path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}

		resizedImg := resize.Resize(uint(tileSize), uint(tileSize), img, resize.Lanczos3)
		tileImages = append(tileImages, resizedImg)
		tileColors = append(tileColors, averageColor(resizedImg))
		return nil
	})
	if err != nil {
		log.Fatalf("failed to load tile images: %v", err)
	}

	if len(tileImages) == 0 {
		log.Fatalf("no tile images found in directory: %s", tileImagesDir)
	}

	bounds := targetImage.Bounds()
	mosaic := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y += tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x += tileSize {
			rect := image.Rect(x, y, x+tileSize, y+tileSize)
			subImg := targetImage.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(rect)
			avgColor := averageColor(subImg)
			bestMatchIndex := findBestMatch(avgColor, tileColors)
			draw.Draw(mosaic, rect, tileImages[bestMatchIndex], image.Point{0, 0}, draw.Over)
		}
	}

	return mosaic
}

func SaveImage(img image.Image, outputPath string) {
	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, img, nil)
}
