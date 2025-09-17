// 代码生成时间: 2025-09-17 15:58:57
package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/kataras/iris/v12"
)

// LogEntry represents a single log entry with timestamp and message.
type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message"`
}

// parseLogFile reads a log file and parses its content into a slice of LogEntry.
// It assumes the log entries are in a simple format with a timestamp followed by a message.
func parseLogFile(filePath string) ([]LogEntry, error) {
    fileContent, err := ioutil.ReadFile(filePath)
    if err != nil {
        return nil, err
    }

    var entries []LogEntry
    lines := strings.Split(string(fileContent), "
")

    for _, line := range lines {
        if line == "" {
            continue
        }

        parts := strings.Fields(line)
        if len(parts) < 2 {
            continue // Skip lines that do not have enough parts.
        }

        timestampStr := parts[0] + " " + parts[1]
        timestamp, err := time.Parse(`2006-01-02 15:04:05`, timestampStr)
        if err != nil {
            continue // Skip lines that do not have a valid timestamp.
        }

        entry := LogEntry{
            Timestamp: timestamp,
            Message:   strings.Join(parts[2:], " "),
        }
        entries = append(entries, entry)
    }

    return entries, nil
}

func main() {
    app := iris.New()

    // Define a route to parse a log file and return its entries.
    app.Get("/parse", func(ctx iris.Context) {
        logFilePath := ctx.URLParam("logFile")
        if logFilePath == "" {
            ctx.StatusCode(iris.StatusBadRequest)
            ctx.JSON(iris.Map{
                "error": "Missing log file path in the URL parameter.",
            })
            return
        }

        // Check if the log file exists.
        if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
            ctx.StatusCode(iris.StatusNotFound)
            ctx.JSON(iris.Map{
                "error": fmt.Sprintf("Log file not found: %s", logFilePath),
            })
            return
        }

        entries, err := parseLogFile(logFilePath)
        if err != nil {
            ctx.StatusCode(iris.StatusInternalServerError)
            ctx.JSON(iris.Map{
                "error": fmt.Sprintf("Failed to parse log file: %s", err.Error()),
            })
            return
        }

        ctx.JSON(entries)
    })

    // Start the Iris web server.
    if err := app.Run(iris.Addr(":8080"), iris.WithoutBanner()); err != nil {
        fmt.Printf("Failed to start server: %s
", err)
   }
}
