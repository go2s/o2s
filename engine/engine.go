// authors: wangoo
// created: 2018-05-30
// gin engine

package engine

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var ginEngine *gin.Engine

func init() {
	ginEngine = gin.Default()
}

func GetGinEngine() *gin.Engine {
	return ginEngine
}

func GinMap(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	ginEngine.GET(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
	ginEngine.POST(pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
