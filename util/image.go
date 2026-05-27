package util

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"voucher-platform/config"

	xdraw "golang.org/x/image/draw"
)

var allowedFormats = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

func ValidateImageFormat(ext string) bool {
	ext = strings.ToLower(ext)
	return allowedFormats[ext]
}

func ValidateImageSize(size int64, maxSize int64) bool {
	return size <= maxSize
}

func SaveImage(file *multipart.FileHeader, uploadPath, filename string) (string, error) {
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(filename)
	ext = strings.ToLower(ext)

	tempPath := filepath.Join(uploadPath, "temp_"+filename)
	dst, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	dst.Close()

	var finalPath string
	if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
		finalPath = filepath.Join(uploadPath, filename)
		if err := CompressImage(tempPath, finalPath); err != nil {
			os.Remove(tempPath)
			return "", err
		}
		os.Remove(tempPath)
	} else {
		finalPath = filepath.Join(uploadPath, filename)
		os.Rename(tempPath, finalPath)
	}

	relativePath := strings.Replace(finalPath, "public/", "/", 1)
	return relativePath, nil
}

func DeleteImage(relativePath string) error {
	if relativePath == "" {
		return nil
	}
	fullPath := "public" + relativePath
	return os.Remove(fullPath)
}

func CompressImage(inputPath, outputPath string) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var img image.Image
	var ext string = strings.ToLower(filepath.Ext(inputPath))

	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return errors.New("unsupported image format for compression")
	}

	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	maxSide := 1920
	if width <= maxSide && height <= maxSide {
		srcFile, err := os.Open(inputPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	}

	var newWidth, newHeight int
	if width > height {
		newWidth = maxSide
		newHeight = int(float64(height) * float64(maxSide) / float64(width))
	} else {
		newHeight = maxSide
		newWidth = int(float64(width) * float64(maxSide) / float64(height))
	}

	newImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	xdraw.BiLinear.Scale(newImg, newImg.Bounds(), img, img.Bounds(), xdraw.Over, nil)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(outFile, newImg, &jpeg.Options{Quality: 85})
	case ".png":
		return png.Encode(outFile, newImg)
	}

	return nil
}

func GetFullImageURL(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	if strings.HasPrefix(path, "/product") {
		return path
	}
	return config.AppConfig.ServerImageURL + path
}
