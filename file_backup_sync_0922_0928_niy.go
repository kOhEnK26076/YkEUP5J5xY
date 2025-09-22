// 代码生成时间: 2025-09-22 09:28:10
package main

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

// Config 包含备份和同步的配置信息
type Config struct {
    SourcePath string
    TargetPath string
    BackupPath string
}

// BackupAndSync 执行文件备份和同步操作
func BackupAndSync(config Config) error {
    // 检查源路径是否存在
    if _, err := os.Stat(config.SourcePath); os.IsNotExist(err) {
        return fmt.Errorf("source path does not exist: %s", config.SourcePath)
    }

    // 确保目标路径存在
    if err := os.MkdirAll(config.TargetPath, 0755); err != nil {
        return fmt.Errorf("failed to create target path: %s", err)
    }

    // 确保备份路径存在
    if err := os.MkdirAll(config.BackupPath, 0755); err != nil {
        return fmt.Errorf("failed to create backup path: %s", err)
    }

    // 遍历源路径中的文件和文件夹
    err := filepath.Walk(config.SourcePath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // 跳过源路径本身
        if path == config.SourcePath {
            return nil
        }

        // 创建目标路径中的相对路径
        relativePath := strings.TrimPrefix(path, config.SourcePath)
        targetFilePath := filepath.Join(config.TargetPath, relativePath)

        // 创建备份路径中的相对路径
        backupFilePath := filepath.Join(config.BackupPath, relativePath)

        // 如果是文件，则进行备份和同步
        if !info.IsDir() {
            // 备份文件
            if err := backupFile(path, backupFilePath); err != nil {
                return err
            }

            // 同步文件
            if err := syncFile(path, targetFilePath); err != nil {
                return err
            }
        }

        // 如果是文件夹，则确保目标路径中存在对应的文件夹
        if info.IsDir() {
            if err := os.MkdirAll(targetFilePath, 0755); err != nil {
                return fmt.Errorf("failed to create directory: %s", err)
            }
        }

        return nil
    })

    return err
}

// backupFile 备份文件
func backupFile(sourcePath, backupPath string) error {
    src, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("failed to open source file: %s", err)
    }
    defer src.Close()

    dst, err := os.Create(backupPath)
    if err != nil {
        return fmt.Errorf("failed to create backup file: %s", err)
    }
    defer dst.Close()

    if _, err := io.Copy(dst, src); err != nil {
        return fmt.Errorf("failed to copy file: %s", err)
    }
    return nil
}

// syncFile 同步文件
func syncFile(sourcePath, targetPath string) error {
    src, err := ioutil.ReadFile(sourcePath)
    if err != nil {
        return fmt.Errorf("failed to read source file: %s", err)
    }

    if err := ioutil.WriteFile(targetPath, src, 0644); err != nil {
        return fmt.Errorf("failed to write target file: %s", err)
    }
    return nil
}

func main() {
    // 配置示例
    config := Config{
        SourcePath: "./source",
        TargetPath: "./target",
        BackupPath: "./backup",
    }

    // 执行备份和同步
    if err := BackupAndSync(config); err != nil {
        log.Fatalf("backup and sync failed: %s", err)
    }

    fmt.Println("Backup and sync completed successfully.")
}