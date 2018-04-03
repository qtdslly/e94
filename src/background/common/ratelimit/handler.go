package ratelimit

import "github.com/gin-gonic/gin"

func RatelimitHandler(opPerSeconds int) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
