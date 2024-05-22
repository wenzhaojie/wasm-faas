package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

// putWatermark 函数将水印添加到图片上
func putWatermark(imageData string, watermarkData string) (string, error) {
	// 解码图片数据
	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", err
	}
	imageReader := bytes.NewReader(imageBytes)
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return "", err
	}

	// 解码水印数据
	watermarkBytes, err := base64.StdEncoding.DecodeString(watermarkData)
	if err != nil {
		return "", err
	}
	watermarkReader := bytes.NewReader(watermarkBytes)
	watermark, err := png.Decode(watermarkReader)
	if err != nil {
		return "", err
	}

	// 在图片上绘制水印
	offset := image.Pt(10, 10)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.Point{}, draw.Src)
	draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.Point{}, draw.Over)

	// 将处理后的图片编码为Base64字符串
	var buf bytes.Buffer
	if err := png.Encode(&buf, m); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// resizeImage 函数调整图片大小
func resizeImage(imageData string, width, height int) (string, error) {
	// 解码图片数据
	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", err
	}
	imageReader := bytes.NewReader(imageBytes)
	img, _, err := image.Decode(imageReader)
	if err != nil {
		return "", err
	}

	// 调整图片大小
	resized := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	// 将处理后的图片编码为Base64字符串
	var buf bytes.Buffer
	if err := png.Encode(&buf, resized); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func main() {
	// 从本地加载图片
	imageFile, err := os.Open("original.png")
	if err != nil {
		log.Fatal(err)
	}
	defer imageFile.Close()

	// 读取图片数据
	var imageBuffer bytes.Buffer
	_, err = imageBuffer.ReadFrom(imageFile)
	if err != nil {
		log.Fatal(err)
	}
	imageData := base64.StdEncoding.EncodeToString(imageBuffer.Bytes())

	// 加载水印图片
	watermarkFile, err := os.Open("watermark.png")
	if err != nil {
		log.Fatal(err)
	}
	defer watermarkFile.Close()

	// 读取水印图片数据
	var watermarkBuffer bytes.Buffer
	_, err = watermarkBuffer.ReadFrom(watermarkFile)
	if err != nil {
		log.Fatal(err)
	}
	watermarkData := base64.StdEncoding.EncodeToString(watermarkBuffer.Bytes())

	// 将水印添加到图片上
	imageWithWatermark, err := putWatermark(imageData, watermarkData)
	if err != nil {
		log.Fatal(err)
	}

	// 调整图片大小
	resizedImage, err := resizeImage(imageWithWatermark, 300, 200)
	if err != nil {
		log.Fatal(err)
	}

	// 输出处理后的图片数据
	fmt.Println("Resized image with watermark (Base64):")
	fmt.Println(strings.TrimSpace(resizedImage))
}
