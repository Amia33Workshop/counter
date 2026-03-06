package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func toFixed(num float64, precision int) string {
	return strconv.FormatFloat(num, 'f', precision, 64)
}
func getCountImage(opts getCountImageOptions) string {
	theme, ok := themeList[opts.Theme]
	if !ok {
		theme = themeList["moebooru"]
	}
	countStr := strconv.Itoa(opts.Count)
	if opts.Padding > len(countStr) {
		countStr = strings.Repeat("0", opts.Padding-len(countStr)) + countStr
	}
	var countArray []string
	countArray = strings.Split(countStr, "")
	if opts.Prefix != -1 {
		prefixArray := strings.Split(strconv.Itoa(opts.Prefix), "")
		countArray = append(prefixArray, countArray...)
	}
	if _, ok := theme["_start"]; ok {
		countArray = append([]string{"_start"}, countArray...)
	}
	if _, ok := theme["_end"]; ok {
		countArray = append(countArray, "_end")
	}
	uniqueChars := make(map[string]bool)
	for _, char := range countArray {
		if _, ok := theme[char]; ok {
			uniqueChars[char] = true
		}
	}
	var x, y float64
	var defsBuilder strings.Builder
	for char := range uniqueChars {
		if img, ok := theme[char]; ok {
			width := float64(img.Width) * opts.Scale
			height := float64(img.Height) * opts.Scale
			y = math.Max(y, height)
			fmt.Fprintf(&defsBuilder, `<image id="%s" width="%s" height="%s" xlink:href="%s" />`, char, toFixed(width, 5), toFixed(height, 5), img.Data)
		}
	}
	defs := defsBuilder.String()
	var partsBuilder strings.Builder
	for _, char := range countArray {
		if img, ok := theme[char]; ok {
			width := float64(img.Width) * opts.Scale
			height := float64(img.Height) * opts.Scale
			yOffset := 0.0
			switch opts.Align {
			case "center":
				yOffset = (y - height) / 2
			case "bottom":
				yOffset = y - height
			}
			yOffsetStr := ""
			if yOffset != 0 {
				yOffsetStr = fmt.Sprintf(` y="%s"`, toFixed(yOffset, 5))
			}
			fmt.Fprintf(&partsBuilder, `<use x="%s"%s xlink:href="#%s" />`, toFixed(x, 5), yOffsetStr, char)
			x += width + float64(opts.Offset)
		}
	}
	parts := partsBuilder.String()
	if len(countArray) > 0 && opts.Offset != 0 {
		x -= float64(opts.Offset)
	}
	var styleBuilder strings.Builder
	hasPixelated := opts.Pixelated
	hasDarkMode := opts.DarkMode == "1"
	if hasPixelated || hasDarkMode {
		styleBuilder.WriteString("svg {")
		if hasPixelated {
			styleBuilder.WriteString(" image-rendering: pixelated;")
		}
		if hasDarkMode {
			styleBuilder.WriteString(" filter: brightness(.6);")
		}
		styleBuilder.WriteString(" }")
	}
	if opts.DarkMode == "auto" {
		if styleBuilder.Len() > 0 {
			styleBuilder.WriteByte(' ')
		}
		styleBuilder.WriteString("@media (prefers-color-scheme: dark) { svg { filter: brightness(.6); } }")
	}
	style := styleBuilder.String()
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg viewBox="0 0 %s %s" width="%s" height="%s" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <title>Counter</title>
  <style>%s</style>
  <defs>%s</defs>
  <g>%s</g>
</svg>`, toFixed(x, 5), toFixed(y, 5), toFixed(x, 5), toFixed(y, 5), style, defs, parts)
}
