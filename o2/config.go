// authors: wangoo
// created: 2018-06-04
// oauth2 server config

package o2

import jwt "github.com/dgrijalva/jwt-go"

//ServerConfig oauth2 server config
type ServerConfig struct {
	// oauth2 server name, will be show in login and authorize page
	ServerName string

	// favicon url
	Favicon string

	// logo url
	Logo string

	// uri context
	URIContext string

	// uri prefix to add before authRedirect uri
	URIPrefix string

	//JWT config
	JWT JWTConfig
}

// JWTConfig jwt config
type JWTConfig struct {
	// Support jwt token
	Support bool

	// SignKey jwt sign key
	SignKey []byte

	// SignMethod jwt sign method
	SignMethod jwt.SigningMethod
}

// DefaultServerConfig default server config
func DefaultServerConfig() *ServerConfig {
	if defaultOauth2Cfg == nil {
		defaultOauth2Cfg = &ServerConfig{
			URIPrefix:  "",
			URIContext: "/oauth2",
			ServerName: "Oauth2 Server",
			Logo:       "https://oauth.net/images/oauth-2-sm.png",
			Favicon:    "https://oauth.net/images/oauth-logo-square.png",
			JWT: JWTConfig{
				Support: false,
			},
		}
	}
	return defaultOauth2Cfg
}
