package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/color"
	"image/jpeg"
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

// MAIN
func main() {
	imgPath, err := os.Open("gray.jpg")
	cErr(err)
	defer imgPath.Close()

	src, err := jpeg.Decode(imgPath)
	cErr(err)
	b := src.Bounds()
	//src := toGray(img) // must be a better way

	// test encoder and decoder
	enc := msgEncoder([]uint8("Cat!"))

	// Image from code
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

	// Write to src file
	dstPath, err := os.Create("grayCode.jpg")
	cErr(err)
	jpeg.Encode(dstPath, dst, nil)

	// get data from encoded image
	msg := make([]uint8, b.Max.Y*b.Max.X)
	ix = 0
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			c := dst.GrayAt(x, y)
			msg[ix] = c.Y & 0x01
			ix++
		}
	}

	fmt.Println(string(msgDecoder(msg)))

	// Make a gray image
	//	dstPath, err := os.Create("gray.jpg")
	//	cErr(err)
	//	defer dstPath.Close()
	//
	//	jpeg.Encode(dstPath, toGray(img), nil)
}
