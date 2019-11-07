// authors: wangoo
// created: 2018-05-30
// gin engine

package engine

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var ginEngine *gin.Engine

func GetGinEngine() *gin.Engine {
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	return ginEngine
}

func GinMap(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	GetGinEngine().Handle(method, pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
