// 代码生成时间: 2025-09-20 21:15:36
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/kataras/iris/v12"
)

// FileBackupSync represents a file backup and sync tool service
type FileBackupSync struct {
    srcPath  string
    destPath string
    backupPath string
}

// NewFileBackupSync creates a new instance of FileBackupSync
func NewFileBackupSync(src, dest, backup string) *FileBackupSync {
    return &FileBackupSync{
        srcPath:  src,
        destPath: dest,
        backupPath: backup,
    }
}

// Sync syncs the source directory to the destination directory
func (fbs *FileBackupSync) Sync() error {
    err := filepath.Walk(fbs.srcPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        destFilePath := filepath.Join(fbs.destPath, filepath.Rel(fbs.srcPath, path))
        if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
            return err
        }

        srcFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer srcFile.Close()

        destFile, err := os.Create(destFilePath)
        if err != nil {
            return err
        }
        defer destFile.Close()

        if _, err := io.Copy(destFile, srcFile); err != nil {
            return err
        }

        return nil
    }); err != nil {
        return err
    }

    return nil
}

// Backup creates a backup of the source directory
func (fbs *FileBackupSync) Backup() error {
    if err := os.MkdirAll(fbs.backupPath, 0755); err != nil {
        return err
    }

    timestamp := time.Now().Format("20060102150405")
    backupDir := filepath.Join(fbs.backupPath, timestamp)
    if err := os.MkdirAll(backupDir, 0755); err !=
        nil {
        return err
    }

    return filepath.Walk(fbs.srcPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        relPath := filepath.Rel(fbs.srcPath, path)
        backupFilePath := filepath.Join(backupDir, relPath)
        return ioutil.WriteFile(backupFilePath, []byte{}, 0644)
    })
}

func main() {
    app := iris.New()

    // Initialize the file backup and sync tool
    fbs := NewFileBackupSync("./src", "./dest", "./backup")

    // Handle sync request
    app.Post("/sync", func(ctx iris.Context) {
        if err := fbs.Sync(); err != nil {
            ctx.StatusCode(http.StatusInternalServerError)
            ctx.JSON(iris.Map{"error": "Sync failed"})
            return
        }

        ctx.JSON(iris.Map{"message": "Sync completed successfully"})
    })

    // Handle backup request
    app.Post("/backup", func(ctx iris.Context) {
        if err := fbs.Backup(); err != nil {
            ctx.StatusCode(http.StatusInternalServerError)
            ctx.JSON(iris.Map{"error": "Backup failed"})
            return
        }

        ctx.JSON(iris.Map{"message": "Backup completed successfully"})
    })

    // Start the Iris HTTP server
    if err := app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8")); err !=
        nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
