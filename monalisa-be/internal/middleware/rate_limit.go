package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type client struct {
	count int
	last  time.Time
}

var (
	mu      sync.Mutex
	clients = map[string]*client{}
)

func RateLimit(max int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		defer mu.Unlock()

		cl, ok := clients[ip]
		if !ok || time.Since(cl.last) > window {
			clients[ip] = &client{count: 1, last: time.Now()}
			return
		}

		if cl.count >= max {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "too many requests",
			})
			return
		}

		cl.count++
	}
}
