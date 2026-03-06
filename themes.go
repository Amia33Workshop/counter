package main

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ThemeImage struct {
	Width  int
	Height int
	Data   string
}

var themeList = make(map[string]map[string]ThemeImage)

func loadThemes() {
	themePath := "./assets/theme"
	files, err := os.ReadDir(themePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			themeName := file.Name()
			themeList[themeName] = make(map[string]ThemeImage)
			themePath := filepath.Join(themePath, themeName)
			imgFiles, err := os.ReadDir(themePath)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, imgFile := range imgFiles {
				ext := filepath.Ext(imgFile.Name())
				if ext == ".png" || ext == ".gif" || ext == ".jpg" || ext == ".jpeg" {
					char := strings.TrimSuffix(imgFile.Name(), ext)
					fullPath := filepath.Join(themePath, imgFile.Name())
					img, err := loadImage(fullPath)
					if err != nil {
						log.Printf("Error loading image %s: %s\n", fullPath, err)
						continue
					}
					themeList[themeName][char] = img
				}
			}
		}
	}
}
func loadImage(path string) (ThemeImage, error) {
	file, err := os.Open(path)
	if err != nil {
		return ThemeImage{}, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return ThemeImage{}, err
	}
	file.Seek(0, 0)
	bytes, err := io.ReadAll(file)
	if err != nil {
		return ThemeImage{}, err
	}
	return ThemeImage{
		Width:  img.Bounds().Dx(),
		Height: img.Bounds().Dy(),
		Data:   fmt.Sprintf("data:%s;base64,%s", http.DetectContentType(bytes), base64.StdEncoding.EncodeToString(bytes)),
	}, nil
}
