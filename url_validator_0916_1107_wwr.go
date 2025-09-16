// 代码生成时间: 2025-09-16 11:07:19
package main

import (
    "fmt"
    "net/url"
    "strings"

    "github.com/kataras/iris/v12"
)

// URLValidator checks if a URL is valid
type URLValidator struct{}

// Validate checks the validity of the given URL
func (v *URLValidator) Validate(ctx iris.Context) {
    // Extracting URL from the request
    inputURL := ctx.URLParam("url")

    // Check if URL is empty
    if inputURL == "" {
        ctx.StatusCode(iris.StatusBadRequest)
        ctx.JSON(iris.Map{
            "error": "URL parameter is missing or empty",
        })
        return
    }

    // Parse the URL
    parsedURL, err := url.ParseRequestURI(inputURL)
    if err != nil {
        ctx.StatusCode(iris.StatusBadRequest)
        ctx.JSON(iris.Map{
            "error": fmt.Sprintf("Invalid URL format: %s", err.Error()),
        })
        return
    }

    // Check if the scheme is either http or https
    if !strings.HasPrefix(parsedURL.Scheme, "http") {
        ctx.StatusCode(iris.StatusBadRequest)
        ctx.JSON(iris.Map{
            "error": "URL scheme must be either http or https",
        })
        return
    }

    // If all checks pass, return a valid response
    ctx.JSON(iris.Map{
        "message": "URL is valid",
        "url": inputURL,
    })
}

func main() {
    app := iris.New()

    // Setup a route for URL validation
    app.Get("/validate", func(ctx iris.Context) {
        validationService := &URLValidator{}
        validationService.Validate(ctx)
    })

    // Start the server
    app.Listen(":8080")
}