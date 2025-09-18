// 代码生成时间: 2025-09-18 12:21:18
package main

import (
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "log"
    "strings"

    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/middleware/cors"
)

// ImageResizer 结构体，用于批量调整图片尺寸
type ImageResizer struct {
    targetWidth  int
    targetHeight int
}

// NewImageResizer 创建一个新的ImageResizer实例
func NewImageResizer(targetWidth, targetHeight int) *ImageResizer {
    return &ImageResizer{
        targetWidth:  targetWidth,
        targetHeight: targetHeight,
    }
}

// ResizeImage 调整单个图片的尺寸
func (r *ImageResizer) ResizeImage(imgPath, outputPath string) error {
    file, err := os.Open(imgPath)
    if err != nil {
        return err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return err
    }

    // 计算新的尺寸
    imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
    newWidth, newHeight := r.targetWidth, r.targetHeight
    if float64(newWidth)/float64(newHeight) != float64(imgWidth)/float64(imgHeight) {
        // 保持宽高比
        ratio := float64(newWidth) / float64(imgWidth)
        if ratio > float64(newHeight)/float64(imgHeight) {
            newHeight = int(float64(imgHeight) * ratio)
        } else {
            newWidth = int(float64(imgWidth) * ratio)
        }
    }

    // 创建新的图片
    resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
    resize := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
    resizedImg.Pix = resize

    // 保存图片
    fileOutput, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer fileOutput.Close()

    switch filepath.Ext(imgPath) {
    case ".jpg":
        err = jpeg.Encode(fileOutput, resizedImg, nil)
    case ".png":
        err = png.Encode(fileOutput, resizedImg)
    default:
        err = fmt.Errorf("unsupported image format")
    }
    if err != nil {
        return err
    }
    return nil
}

// BatchResize 批量调整图片尺寸
func (r *ImageResizer) BatchResize(sourceDir, outputDir string) error {
    files, err := ioutil.ReadDir(sourceDir)
    if err != nil {
        return err
    }

    for _, file := range files {
        if file.IsDir() {
            continue
        }

        imgPath := filepath.Join(sourceDir, file.Name())
        outputPath := filepath.Join(outputDir, file.Name())
        if err := r.ResizeImage(imgPath, outputPath); err != nil {
            log.Printf("Failed to resize image %s: %v", file.Name(), err)
        }
    }
    return nil
}

func main() {
    app := iris.New()
    app.Use(cors.New(cors.Options{
       AllowedOrigins: []string{irisANY},
       AllowedMethods: []string{irisGET, irisHEAD, irisPOST, irisPUT, irisDELETE},
       AllowCredentials: true,
    }))

    // 创建一个新的ImageResizer实例
    resizer := NewImageResizer(800, 600)

    // 批量调整图片尺寸
    app.Post("/batch-resize", func(ctx iris.Context) {
        sourceDir := ctx.URLParam("sourceDir")
        outputDir := ctx.URLParam("outputDir")

        if err := resizer.BatchResize(sourceDir, outputDir); err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.JSON(iris.Map{
                "error": err.Error(),
            })
        } else {
            ctx.JSON(iris.Map{
                "message": "Images resized successfully",
            })
        }
    })

    app.Listen(":8080")
}
