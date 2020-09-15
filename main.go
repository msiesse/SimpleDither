package main

import (
	"os"
	"log"
	"image/jpeg"
	"image/png"
	"image/color"
	"image/draw"
	"image"
	"math"
)

func simpleDither(image *image.Gray16) {
	var err, t	float64
	var k		int
	var width, height	int = image.Bounds().Size().X, image.Bounds().Size().Y

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			if i % 2 == 0 {
				k = j
			} else {
				k = height - j - 1
			}
			err = float64(image.Gray16At(i, k).Y) / 65535
			if err >= 0.5 {
				t = 1
			} else {
				t = 0
			}
			image.SetGray16(i, k, color.Gray16{uint16(math.Round(err) * 65535.0)})
			err = err - t
			if i + 1 < width {
				image.SetGray16(i + 1, k, 
					color.Gray16{image.Gray16At(i + 1, k).Y + uint16((7 / 16) * err)})
				if k + 1 < height {
					image.SetGray16(i + 1, k + 1, 
						color.Gray16{image.Gray16At(i + 1, k + 1).Y + uint16((3 / 16) * err)})
				}
			} else if k + 1 < height {
				image.SetGray16(i, k + 1, 
					color.Gray16{image.Gray16At(i, k + 1).Y + uint16((5 / 16) * err)})
				if i - 1 >= 0 {
					image.SetGray16(i - 1, k + 1, 
						color.Gray16{image.Gray16At(i - 1, k + 1).Y + uint16((1 / 16) * err)})
				}
			}
		}
	}
}

func main() {

	if (len(os.Args) == 1) {
		os.Exit(1)
	}

	imagePath := os.Args[1]
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	imageData, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	myImg := image.NewGray16(imageData.Bounds())
	draw.Draw(myImg, imageData.Bounds(), imageData, image.Point{}, draw.Over)
	simpleDither(myImg)

	result, err := os.Create("result.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(result, myImg); err != nil {
		result.Close()
		log.Fatal(err)
	}
	if err := result.Close(); err != nil {
		log.Fatal(err)
	}
}