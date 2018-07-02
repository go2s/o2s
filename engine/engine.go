// authors: wangoo
// created: 2018-05-30
// gin engine

package engine

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

var ginEngine *gin.Engine

func GetGinEngine() *gin.Engine {
	if ginEngine == nil {
		ginEngine = gin.Default()
	}
	return ginEngine
}

func GinMap(method, pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	fmt.Printf("gin map [%v]%v\n", method, pattern)
	GetGinEngine().Handle(method, pattern, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}
