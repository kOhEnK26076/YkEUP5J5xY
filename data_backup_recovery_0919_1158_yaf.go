// 代码生成时间: 2025-09-19 11:58:10
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/kataras/iris/v12"
)

// BackupData 用于备份数据
func BackupData(w http.ResponseWriter, r *http.Request) {
    backupFilePath := "./backup/" + time.Now().Format("2006-01-02_150405") + ".sql"
    databaseFilePath := "./database/database.sql"
    
    // 检查备份目录是否存在，不存在则创建
    if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
        os.MkdirAll(filepath.Dir(backupFilePath), os.ModePerm)
    }
    
    // 备份数据库
    src, err := os.Open(databaseFilePath)
    if err != nil {
        iris.NewText(w, "Error: Unable to open database file.")
        return
    }
    defer src.Close()

    dst, err := os.Create(backupFilePath)
    if err != nil {
        iris.NewText(w, "Error: Unable to create backup file.")
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        iris.NewText(w, "Error: Failed to backup database.")
        return
    }
    
    iris.NewText(w, "Backup created successfully.")
}

// RestoreData 用于恢复数据
func RestoreData(w http.ResponseWriter, r *http.Request) {
    backupFilePath := "./backup/database_backup.sql"
    databaseFilePath := "./database/database.sql"
    
    // 恢复数据库
    src, err := os.Open(backupFilePath)
    if err != nil {
        iris.NewText(w, "Error: Unable to open backup file.")
        return
    }
    defer src.Close()

    dst, err := os.Create(databaseFilePath)
    if err != nil {
        iris.NewText(w, "Error: Unable to create database file.")
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        iris.NewText(w, "Error: Failed to restore database.")
        return
    }
    
    iris.NewText(w, "Database restored successfully.")
}

func main() {
    app := iris.New()

    // 设置路由
    app.Post("/backup", BackupData)
    app.Post("/restore", RestoreData)

    // 启动服务器
    if err := app.Listen(":8080"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}
