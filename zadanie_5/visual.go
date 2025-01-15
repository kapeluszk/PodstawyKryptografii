package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

// Funkcja do wczytywania obrazu
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// Konwersja obrazu na czarno-biały (tablica 2D int)
func convertImageToBW(img image.Image) [][]int {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	bw := make([][]int, width)
	for x := 0; x < width; x++ {
		bw[x] = make([]int, height)
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if (r+g+b)/3 > 32767 { // Próg (połowa max wartości uint32)
				bw[x][y] = 0 // Biały
			} else {
				bw[x][y] = 1 // Czarny
			}
		}
	}
	return bw
}

func generateShares(original [][]int) ([][]int, [][]int) {
	width := len(original)
	height := len(original[0])
	share1 := make([][]int, width)
	share2 := make([][]int, width)

	for x := 0; x < width; x++ {
		share1[x] = make([]int, height)
		share2[x] = make([]int, height)
		for y := 0; y < height; y++ {
			if original[x][y] == 0 { // Biały piksel
				if rand.Intn(2) == 0 {
					share1[x][y] = 0
					share2[x][y] = 1
				} else {
					share1[x][y] = 1
					share2[x][y] = 0
				}
			} else { // Czarny piksel
				if rand.Intn(2) == 0 {
					share1[x][y] = 0
					share2[x][y] = 0
				} else {
					share1[x][y] = 1
					share2[x][y] = 1
				}
			}
		}
	}
	return share1, share2
}

// Tworzenie obrazu z udziału
func createImageFromShare(share [][]int) image.Image {
	width := len(share)
	height := len(share[0])
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := color.White
			if share[x][y] == 1 {
				c = color.Black
			}
			img.Set(x, y, c)
		}
	}
	return img
}

func saveImage(img image.Image, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func combineShares(share1, share2 [][]int) [][]int {
	width := len(share1)
	height := len(share1[0])
	combined := make([][]int, width)
	for x := 0; x < width; x++ {
		combined[x] = make([]int, height)
		for y := 0; y < height; y++ {
			if share1[x][y] == 1 && share2[x][y] == 1 {
				combined[x][y] = 1 // czarny
			} else {
				combined[x][y] = 0 // biały
			}
		}
	}
	return combined
}

func main() {
	rand.Seed(time.Now().UnixNano())

	img, err := loadImage("input.png")
	if err != nil {
		log.Fatal(err)
	}

	bwImage := convertImageToBW(img)
	share1, share2 := generateShares(bwImage)

	img1 := createImageFromShare(share1)
	img2 := createImageFromShare(share2)
	saveImage(img1, "share1.png")
	saveImage(img2, "share2.png")

	combined := combineShares(share1, share2)
	combinedImg := createImageFromShare(combined)
	saveImage(combinedImg, "combined.png")

	fmt.Println("Udziały i połączony obraz zostały zapisane.")
}
