// 代码生成时间: 2025-10-08 22:37:42
package main

import (
    "crypto/md5"
    "encoding/hex"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/kataras/iris/v12"
)

// Goods represents a product in the跨境电商平台
type Goods struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {
    // Initialize the application with Iris
    app := iris.New()

    // Define a route for creating a new product
    app.Post("/goods", createGoodsHandler)

    // Start the Iris server
    log.Fatal(app.Listen(":8080"))
}

// createGoodsHandler handles the POST request for creating a new product
func createGoodsHandler(ctx iris.Context) {
    // Retrieve the product data from the request body
    var goods Goods
    if err := ctx.ReadJSON(&goods); err != nil {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid request body"})
        return
    }

    // Generate a unique ID for the product using MD5 hash of the current timestamp
    id := fmt.Sprintf("%x", md5.Sum([]byte(goods.Name + time.Now().String())))
    goods.ID = id
    goods.CreatedAt = time.Now()

    // Save the product to a database (for demonstration, we'll just print to console)
    fmt.Printf("Created new product: %+v
", goods)

    // Respond with the created product
    ctx.JSON(iris.Map{
        "message": "Product created successfully",
        "product": goods,
    })
}
