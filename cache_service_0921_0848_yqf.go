// 代码生成时间: 2025-09-21 08:48:33
package main

import (
    "fmt"
    "time"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/cache"
    "github.com/kataras/iris/v12/cache/cfg"
)

// CacheService 结构体用于管理缓存操作
type CacheService struct {
    cache cache.Cache
}

// NewCacheService 初始化缓存服务
func NewCacheService() *CacheService {
    // 使用默认配置创建缓存
    cache := cache.New(cfg.Default())
    return &CacheService{
        cache: cache,
    }
}

// SetCache 设置缓存项
func (s *CacheService) SetCache(key string, value interface{}, expiration time.Duration) error {
    // 检查缓存是否支持设置过期时间
    if _, ok := s.cache.(cache.Expires); ok {
        // 设置带有过期时间的缓存项
        if err := s.cache.Set(key, value, expiration); err != nil {
            return fmt.Errorf("set cache error: %w", err)
        }
    } else {
        return fmt.Errorf("cache does not support expiration")
    }
    return nil
}

// GetCache 获取缓存项
func (s *CacheService) GetCache(key string) (interface{}, error) {
    value, err := s.cache.Get(key)
    if err != nil {
        return nil, fmt.Errorf("get cache error: %w", err)
    }
    return value, nil
}

func main() {
    app := iris.New()
    cacheService := NewCacheService()

    // 设置缓存项
    app.Post("/set_cache", func(ctx iris.Context) {
        key := ctx.URLParam("key")
        value := ctx.URLParam("value")
        expirationStr := ctx.URLParam("expiration")

        expiration, err := time.ParseDuration(expirationStr)
        if err != nil {
            ctx.JSON(iris.StatusBadRequest, iris.Map{
                "error": "invalid expiration format",
            })
            return
        }

        if err := cacheService.SetCache(key, value, expiration); err != nil {
            ctx.JSON(iris.StatusInternalServerError, iris.Map{
                "error": err.Error(),
            })
            return
        }

        ctx.JSON(iris.StatusOK, iris.Map{
            "message": "cache set successfully",
        })
    })

    // 获取缓存项
    app.Get="/get_cache", func(ctx iris.Context) {
        key := ctx.URLParam("key")
        value, err := cacheService.GetCache(key)
        if err != nil {
            ctx.JSON(iris.StatusInternalServerError, iris.Map{
                "error": err.Error(),
            })
            return
        }

        ctx.JSON(iris.StatusOK, iris.Map{
            "value": value,
        })
    })

    // 启动服务器
    if err := app.Run(iris.Addr(":8080")); err != nil {
        fmt.Printf("server error: %s
", err)
    }
}
