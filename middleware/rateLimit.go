package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"voucher-platform/util"
)

type RateLimitConfig struct {
	MaxRequests    int
	WindowSeconds  int64
}

type visitor struct {
	count     int
	lastTime  time.Time
	operation map[string]*visitorOperation
}

type visitorOperation struct {
	count    int
	lastTime time.Time
}

var (
	visitors = sync.Map{}
	mu       sync.RWMutex

	defaultConfigs = map[string]RateLimitConfig{
		"auth": {
			MaxRequests:   5,
			WindowSeconds: 60,
		},
		"send_code": {
			MaxRequests:   3,
			WindowSeconds: 60,
		},
		"register": {
			MaxRequests:   3,
			WindowSeconds: 300,
		},
		"login": {
			MaxRequests:   5,
			WindowSeconds: 60,
		},
		"trade": {
			MaxRequests:   10,
			WindowSeconds: 60,
		},
		"order": {
			MaxRequests:   20,
			WindowSeconds: 60,
		},
		"upload": {
			MaxRequests:   30,
			WindowSeconds: 60,
		},
		"admin": {
			MaxRequests:   300,
			WindowSeconds: 60,
		},
		"admin_login": {
			MaxRequests:   10,
			WindowSeconds: 60,
		},
		"default": {
			MaxRequests:   30,
			WindowSeconds: 60,
		},
	}
)

func SetRateLimitConfig(operation string, config RateLimitConfig) {
	mu.Lock()
	defer mu.Unlock()
	defaultConfigs[operation] = config
}

func RateLimitMiddleware(maxRequests int, windowSeconds int64) gin.HandlerFunc {
	return RateLimitWithOperation("default", maxRequests, windowSeconds)
}

func RateLimitWithOperation(operation string, maxRequests int, windowSeconds int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := ip + ":" + operation

		mu.Lock()
		defer mu.Unlock()

		v, exists := visitors.Load(key)
		now := time.Now()

		if !exists {
			visitors.Store(key, &visitor{
				count:    1,
				lastTime: now,
				operation: map[string]*visitorOperation{
					operation: {count: 1, lastTime: now},
				},
			})
			c.Next()
			return
		}

		visitorData := v.(*visitor)
		opData, opExists := visitorData.operation[operation]

		if !opExists {
			visitorData.operation[operation] = &visitorOperation{count: 1, lastTime: now}
			c.Next()
			return
		}

		if now.Sub(opData.lastTime).Seconds() > float64(windowSeconds) {
			opData.count = 1
			opData.lastTime = now
			c.Next()
			return
		}

		if opData.count >= maxRequests {
			util.GetLogger().Warn("rate_limit_exceeded",
				util.StringField("ip", ip),
				util.StringField("operation", operation),
				util.IntField("count", opData.count),
			)
			util.ResponseError(c, http.StatusTooManyRequests, "请求过于频繁，请稍后重试")
			c.Abort()
			return
		}

		opData.count++
		c.Next()
	}
}

func AuthRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("auth", 5, 60)
}

func SendCodeRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("send_code", 3, 60)
}

func RegisterRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("register", 3, 300)
}

func LoginRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("login", 5, 60)
}

func TradeRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("trade", 10, 60)
}

func OrderRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("order", 20, 60)
}

func UploadRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("upload", 30, 60)
}

func AdminRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("admin", 300, 60)
}

func AdminLoginRateLimit() gin.HandlerFunc {
	return RateLimitWithOperation("admin_login", 10, 60)
}

func CleanExpiredVisitors() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	visitors.Range(func(key, value interface{}) bool {
		v := value.(*visitor)
		if now.Sub(v.lastTime).Seconds() > 3600 {
			visitors.Delete(key)
		}
		return true
	})
}
