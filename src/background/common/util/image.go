package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

// this method allows to build a image with specified width and height
func ResizeImage(url string, width, height int) string {
	ext := filepath.Ext(url)
	fileName := strings.TrimSuffix(url, ext)
	return fmt.Sprintf("%s@w%d_h%d%s", fileName, width, height, ext)
}
