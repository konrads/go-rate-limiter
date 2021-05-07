package decorator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-rate-limiter/pkg/leakybucket"
)

func Decorate(lb *leakybucket.SafeLeakyBucket, handler *gin.HandlerFunc) gin.HandlerFunc {
	var decorated gin.HandlerFunc = func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rejectRule := lb.GetRejectionRule(ip, now)
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
