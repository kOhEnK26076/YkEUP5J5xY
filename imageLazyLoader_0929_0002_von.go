// 代码生成时间: 2025-09-29 00:02:59
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/kataras/iris/v12"
)

// ImageLazyLoader is a struct that holds configurations for the lazy loader.
type ImageLazyLoader struct {
    Dir string // Directory where images are stored.
}

// NewImageLazyLoader creates a new ImageLazyLoader instance.
func NewImageLazyLoader(dir string) *ImageLazyLoader {
    return &ImageLazyLoader{Dir: dir}
}

// LoadImages lists all images in the directory and serves them with lazy loading.
func (loader *ImageLazyLoader) LoadImages(ctx iris.Context) {
    // Check if directory exists.
    if _, err := os.Stat(loader.Dir); os.IsNotExist(err) {
        ctx.StatusCode(iris.StatusNotFound)
        ctx.WriteString("Image directory does not exist.")
        return
    }

    // List all files in the directory.
    files, err := os.ReadDir(loader.Dir)
    if err != nil {
        ctx.StatusCode(iris.StatusInternalServerError)
        ctx.WriteString("Failed to read image directory.")
        return
    }

    // Create a slice to hold image URLs with lazy loading attributes.
    urls := make([]string, 0, len(files))

    // Iterate through files and prepare the URLs.
    for _, file := range files {
        if !file.IsDir() && isImageFile(file.Name()) {
            // Use a relative path to keep the URL consistent with the server's root.
            relativePath := filepath.Rel(loader.Dir, file.Name())
            // Create a lazy loading URL.
            url := fmt.Sprintf("<img src='" + iris.RelURL("/images/" + relativePath) + "' loading='lazy' alt='Lazy Loaded Image'>")
            urls = append(urls, url)
        }
    }

    // Serve the HTML with lazy loaded images.
    ctx.ViewData("images.html", iris.Map{
        "Images": urls,
    })
}

// isImageFile checks if a file has an image extension.
func isImageFile(filename string) bool {
    extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff"}
    for _, ext := range extensions {
        if strings.HasSuffix(strings.ToLower(filename), ext) {
            return true
        }
    }
    return false
}

func main() {
    app := iris.New()
    app.RegisterView(iris.HTML("./templates", ".html")).Reload(true)

    // Create an instance of ImageLazyLoader.
    loader := NewImageLazyLoader("./images")

    // Define route for serving images with lazy loading.
    app.Get("/images", loader.LoadImages)

    // Start the server.
    app.Listen(":8080")
}

// This program uses the IRIS framework to serve images with lazy loading.
// It scans a directory for image files, generates HTML with img tags
// that have the 'loading=lazy' attribute, and serves this HTML to the client.
// The directory path can be customized when creating an ImageLazyLoader instance.
