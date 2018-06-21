package img2ansi

import (
	"fmt"
	"image"
	"image/color"
	"io"
)

const (
	ansiColorBase  int     = 16
	ansiColorSteps float64 = 6
	rgbaColorSpace float64 = float64(1 << 16)
)

// ImageToAnsi writes the converted ANSI image to the writer
func ImageToAnsi(img image.Image, w io.Writer) (err error) {
	//var previous string
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		fmt.Fprint(w, "console.log('")

		lineColorStr := ""
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			color := img.At(x, y)
			colorStr := toDevtoolsConsoleStyle(color)
			if isBg(color) {
				fmt.Printf("%%c ")
			} else {
				fmt.Printf("%%c0")
			}
			lineColorStr += fmt.Sprintf(", '%s'", colorStr)

			//if current != previous {
			//	_, err = fmt.Fprint(w, current)
			//	if err != nil {
			//		return err
			//	}
			//	_, err = fmt.Fprint(w, "0</css>")
			//}
			//_, err = fmt.Fprint(w, "0</css>")
			//if err != nil {
			//	return err
			//}
			//previous = current
		}
		_, err = fmt.Printf("'%s)\n", lineColorStr)
		if err != nil {
			return err
		}
	}
	//_, err = fmt.Fprint(w, "\x1b[0m")
	if err != nil {
		return err
	}
	return nil
}

func isBg(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	r2 := uint8(r)
	g2 := uint8(g)
	b2 := uint8(b)
	if r2 == 0 && g2 == 0 && b2 == 0 {
		return true
	}
	return false
}

func toDevtoolsConsoleStyle(c color.Color) string {
	r, g, b, _ := c.RGBA()
	r2 := uint8(r)
	g2 := uint8(g)
	b2 := uint8(b)
	// note: there are 256 colors in total in the xterm color space
	//       the first 16 colors are 8 base colors and 8 bright colors (skipped with `ansiColorBase`)
	//       the next 216 colors are all colors in steps of 6
	//       the last 24 are a grayscale from black to white
	//code := ansiColorBase + toAnsiSpace(r)*36 + toAnsiSpace(g)*6 + toAnsiSpace(b)
	//if r != 0 && g != 0 && b != 0 {
	//	fmt.Printf("r: %d=%X\ng: %d=%X\nb: %d=%X\n", r2, r2, g2, g2, b2, b2)
	//}
	return fmt.Sprintf("color: #%x%x%x", r2, g2, b2)
}

func IntToHex(n int64) string {
	//return []byte(strconv.FormatInt(n, 16))
	return fmt.Sprintf("%x", n)
}


func toAnsiCode(c color.Color) string {
	r, g, b, _ := c.RGBA()
	// note: there are 256 colors in total in the xterm color space
	//       the first 16 colors are 8 base colors and 8 bright colors (skipped with `ansiColorBase`)
	//       the next 216 colors are all colors in steps of 6
	//       the last 24 are a grayscale from black to white
	code := ansiColorBase + toAnsiSpace(r)*36 + toAnsiSpace(g)*6 + toAnsiSpace(b)
	return fmt.Sprintf("\x1b[38;5;%dm", code)
}

func toAnsiSpace(val uint32) int {
	return int(ansiColorSteps * (float64(val) / rgbaColorSpace))
}
