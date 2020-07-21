package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/color"
	_ "image/jpeg"
	"image/png"
)

// HANDLES ERRORS (not really, just quits :p )
func cErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// BIT ENCODER / DECODER
func msgEncoder(msg []uint8) []uint8 {
	ret := make([]uint8, len(msg)*8)
	count, ix := 8, 0
	for i := 0; i < len(ret); i++ {
		if count == 0 {
			ix++
			count = 8
		}
		ret[i] = msg[ix] >> (count - 1) & 0x01
		count--
	}
	return ret
}

func msgDecoder(msg []uint8) []byte {
	dst := make([]byte, len(msg))

	count := 0
	var cur uint8 = 0
	ix := 0
	for _, b := range msg {
		cur |= b
		if count == 7 {
			dst[ix] = byte(cur)
			ix++
			count = 0
			cur = 0
		} else {
			count++
			cur <<= 1
		}
	}

	return dst[:]
}

// IMAGE MANIPULATION
func toGray(img image.Image) *image.Gray {
	b := img.Bounds()
	g := image.NewGray(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			g.Set(x, y, img.At(x, y))
		}
	}

	return g
}

// ARGS
func gray() {
	imgPath, err := os.Open(os.Args[2])
	cErr(err)
	defer imgPath.Close()
	img, _, err := image.Decode(imgPath)
	cErr(err)
	dstPath, err := os.Create("gray_" + os.Args[2]) // e.g: "gray_cat.jpg"
	cErr(err)
	defer dstPath.Close()

	png.Encode(dstPath, toGray(img))
}

func enc() {
	imgPath, err := os.Open(os.Args[2])
	cErr(err)
	defer imgPath.Close()

	dstPath, err := os.Create("enc_" + os.Args[2])
	cErr(err)
	defer dstPath.Close()

	// Read src image
	src, err := png.Decode(imgPath)
	cErr(err)

	// Encode text
	enc := msgEncoder([]uint8(os.Args[3]))

	// Encode image
	b := src.Bounds()
	ix := 0
	dst := image.NewGray(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			a := color.GrayModel.Convert(src.At(x, y)).(color.Gray).Y
			if ix == len(enc) {
				ix = 0
			}
			c := color.Gray{Y: a&0xfe | enc[ix]}
			dst.Set(x, y, c)
			ix++
		}
	}

	// Write to destination
	png.Encode(dstPath, dst)
}

func dec() {
	imgPath, err := os.Open(os.Args[2])
	cErr(err)
	defer imgPath.Close()

	img, err := png.Decode(imgPath)
	cErr(err)

	b := img.Bounds()

	msg := make([]uint8, b.Max.Y*b.Max.X)
	ix := 0
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			p := img.At(x, y)
			c := color.GrayModel.Convert(p).(color.Gray).Y
			msg[ix] = c & 0x01
			ix++
		}
	}

	//fmt.Println(msgDecoder(msg)[:100])
	fmt.Println(string(msgDecoder(msg)))

}

// MAIN
func main() {
	if len(os.Args) < 3 {
		fmt.Printf("USAGE: %v <gray|enc|dec> <image> [text]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gray":
		gray()
	case "enc":
		if len(os.Args) != 4 {
			fmt.Println("Text not provided")
			os.Exit(3)
		}
		enc()
	case "dec":
		dec()
	default:
		fmt.Println("Unknown command", os.Args[1])
		os.Exit(2)
	}
}
