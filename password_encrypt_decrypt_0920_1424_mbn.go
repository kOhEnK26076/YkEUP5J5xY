// 代码生成时间: 2025-09-20 14:24:48
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "strings"

    "github.com/kataras/iris/v12"
)

// Key represents the AES key used for encryption and decryption.
var Key = []byte("your-256-bit-key-here")

// Encrypt encrypts a plaintext message using AES.
func Encrypt(text string) (string, error) {
    block, err := aes.NewCipher(Key)
    if err != nil {
        return "", err
    }

    textBytes := []byte(text)
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    encrypted := gcm.Seal(nonce, nonce, textBytes, nil)
    return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts an encrypted message using AES.
func Decrypt(encrypted string) (string, error) {
    encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(Key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(encryptedBytes) < nonceSize {
        return "", fmt.Errorf("malformed ciphertext")
    }

    nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

func main() {
    app := iris.New()
    app.Handle("GET", "/encrypt", func(ctx iris.Context) {
        text := ctx.URLParam("text")
        encrypted, err := Encrypt(text)
        if err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.WriteString(err.Error())
            return
        }
        ctx.WriteString(fmt.Sprintf("Encrypted: %s
", encrypted))
    })

    app.Handle("GET", "/decrypt", func(ctx iris.Context) {
        encrypted := ctx.URLParam("encrypted")
        decrypted, err := Decrypt(encrypted)
        if err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.WriteString(err.Error())
            return
        }
        ctx.WriteString(fmt.Sprintf("Decrypted: %s
", decrypted))
    })

    // Run the IRIS server on port 8080
    app.Listen(":8080")
}
