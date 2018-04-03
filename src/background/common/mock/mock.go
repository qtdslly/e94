package mock

import (
	"common/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// this handler allows to intercept the api request and return mocked response.
func MockHandler(c *gin.Context) {
	url := c.Request.URL.Path

	logger.Error(url)

	switch url {
	case "/xxx":
		c.String(http.StatusOK, "{}")
		c.AbortWithStatus(http.StatusOK)
		return
	}

}
