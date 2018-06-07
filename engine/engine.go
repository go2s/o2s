// authors: wangoo
// created: 2018-05-30
// gin engine

package engine

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var ginEngine *gin.Engine

func GetGinEngine() *gin.Engine {
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	return ginEngine
}

func GinMap(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	GetGinEngine().GET(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
	GetGinEngine().POST(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
