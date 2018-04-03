package fms

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func buildThumbPath(path string, width int, height int) string {
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(path)
	if ext != "" {
		name = name[:len(name)-len(ext)]
	}
	return filepath.Join(dir, name+"@w"+strconv.Itoa(width)+"_h"+strconv.Itoa(height)+ext)
}

func isThumb(path string) bool {
	n := filepath.Base(path)
	return strings.Contains(n, "@w") && strings.Contains(n, "_h")
}

func createThumb(path string) error {
	ext := filepath.Ext(path)
	/*
		if ext != ".jpg" && ext != ".jpeg" &&
			ext != ".png" && ext != ".bmp" && ext != ".gif" {
			return errors.New("Unsupport file format " + ext)
		}
	*/

	// example path: ../file/o/2016-05-19/91bd48ce-1dae-11e6-9d7f-408d5cdf2c91@w360_h100.jpg
	// take the string after "@" (w360_h100.jpg) to extract the size (360*270)
	sl := strings.Split(path, "@")
	sizeString := sl[1]
	ws := sizeString[strings.Index(sizeString, "@w")+2 : strings.Index(sizeString, "_h")]
	var hs string
	if strings.Index(sizeString, ".") == -1 {
		hs = sizeString[strings.Index(sizeString, "_h")+2:]
	} else {
		hs = sizeString[strings.Index(sizeString, "_h")+2 : strings.Index(sizeString, ".")]
	}
	w, _ := strconv.Atoi(ws)
	h, _ := strconv.Atoi(hs)

	srcPath := sl[0] + ext
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	i, _, err := image.Decode(src)
	if err != nil {
		return err
	}

	m := imaging.Fill(i, w, h, imaging.Center, imaging.Lanczos)

	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, m, nil)
}
