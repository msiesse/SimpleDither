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
	//"fmt"
)

func simpleDither(image *image.Gray16) {
	var err, t	float64
	var width, height	int = image.Bounds().Size().X, image.Bounds().Size().Y

	for i := 0; i < width; i++ {
		if i % 2 == 0 {
			for j := 0; j < height; j++ {
				err = float64(image.Gray16At(i, j).Y) / 65535.0
				if err >= 0.5 {
					t = 1
				} else {
					t = 0
				}
				image.SetGray16(i, j, color.Gray16{uint16(math.Round(err) * 65535.0)})
				err = err - t
				if j + 1 < height {
					image.SetGray16(i, j + 1, 
						color.Gray16{image.Gray16At(i, j + 1).Y + uint16(math.Round((7. / 16.) * err * 65535.0))})
					if i + 1 < width {
						image.SetGray16(i + 1, j + 1, 
							color.Gray16{image.Gray16At(i + 1, j + 1).Y + uint16(math.Round((3. / 16.) * err * 65535.0))})
					}
				} else if i + 1 < width {
					image.SetGray16(i + 1, j, 
						color.Gray16{image.Gray16At(i + 1, j).Y + uint16(math.Round((5. / 16.) * err * 65535.0))})
					if j - 1 >= 0 {
						image.SetGray16(i + 1, j - 1, 
							color.Gray16{image.Gray16At(i + 1, j - 1).Y + uint16(math.Round((1. / 16.) * err * 65535.0))})
					}
				}
			}
		} else {
			for j := height - 1; j >= 0; j-- {
				err = float64(image.Gray16At(i, j).Y) / 65535.0
				if err >= 0.5 {
					t = 1
				} else {
					t = 0
				}
				image.SetGray16(i, j, color.Gray16{uint16(math.Round(err) * 65535.0)})
				err = err - t
				if j - 1 >= 0 {
					image.SetGray16(i, j - 1, 
						color.Gray16{image.Gray16At(i, j - 1).Y + uint16(math.Round((7. / 16.) * err * 65535.0))})
					if i + 1 < width {
						image.SetGray16(i + 1, j - 1, 
							color.Gray16{image.Gray16At(i + 1, j - 1).Y + uint16(math.Round((3. / 16.) * err * 65535.0))})
					}
				} else if i + 1 < width {
					image.SetGray16(i + 1, j, 
						color.Gray16{image.Gray16At(i + 1, j).Y + uint16(math.Round((5. / 16.) * err * 65535.0))})
					if j + 1 < height {
						image.SetGray16(i + 1, j + 1, 
							color.Gray16{image.Gray16At(i + 1, j + 1).Y + uint16(math.Round((1. / 16.) * err * 65535.0))})
					}
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