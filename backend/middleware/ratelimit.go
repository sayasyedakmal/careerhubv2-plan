package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type rateLimiterStore struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	r        rate.Limit
	b        int
}

func newStore(r rate.Limit, b int) *rateLimiterStore {
	s := &rateLimiterStore{
		limiters: make(map[string]*ipLimiter),
		r:        r,
		b:        b,
	}
	// Clean up stale entries every minute
	go func() {
		for {
			time.Sleep(time.Minute)
			s.mu.Lock()
			for ip, l := range s.limiters {
				if time.Since(l.lastSeen) > 3*time.Minute {
					delete(s.limiters, ip)
				}
			}
			s.mu.Unlock()
		}
	}()
	return s
}

func (s *rateLimiterStore) get(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	l, exists := s.limiters[ip]
	if !exists {
		l = &ipLimiter{limiter: rate.NewLimiter(s.r, s.b)}
		s.limiters[ip] = l
	}
	l.lastSeen = time.Now()
	return l.limiter
}

// microsoftLoginStore — 30 req/min
var microsoftLoginStore = newStore(rate.Every(2*time.Second), 5)

// refreshStore — 60 req/min
var refreshStore = newStore(rate.Every(time.Second), 5)

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.FullPath()

		var limiter *rate.Limiter
		switch path {
		case "/api/v1/auth/login/microsoft":
			limiter = microsoftLoginStore.get(ip)
		case "/api/v1/auth/refresh":
			limiter = refreshStore.get(ip)
		default:
			c.Next()
			return
		}

		if !limiter.Allow() {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please slow down and try again later.",
				"code":  "RATE_LIMIT_EXCEEDED",
			})
			return
		}

		c.Next()
	}
}
