package decorator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-rate-limiter/pkg/limiter"
)

func Decorate(l *limiter.Limiter, handler *gin.HandlerFunc) gin.HandlerFunc {
	var decorated gin.HandlerFunc = func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rejectRule := l.GetRejectionRule(ip, now)
		if rejectRule != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit breached due to rule: %v", *rejectRule),
			})
		} else {
			(*handler)(c)
		}
	}
	return decorated
}
