// authors: wangoo
// created: 2018-05-30
// gin engine

package engine

import (
	"github.com/go2s/o2s/o2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func mapping(engine *gin.Engine, uri string, handler func(w http.ResponseWriter, r *http.Request)) {
	engine.GET(uri, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
	engine.POST(uri, func(c *gin.Context) {
		handler(c.Writer, c.Request)
	})
}

func MappingHandlers(engine *gin.Engine, uriPrefix string) {
	mapping(engine, uriPrefix+o2.Oauth2UriLogin, o2.LoginHandler)
	mapping(engine, uriPrefix+o2.Oauth2UriAuth, o2.AuthHandler)
	mapping(engine, uriPrefix+o2.Oauth2UriAuthorize, o2.AuthorizeRequestHandler)
	mapping(engine, uriPrefix+o2.Oauth2UriToken, o2.TokenRequestHandler)
	mapping(engine, uriPrefix+o2.Oauth2UriValid, o2.BearerTokenValidator)
}

func NewEngine() *gin.Engine {
	engine := gin.Default()
	MappingHandlers(engine, "")
	return engine
}
