// authors: wangoo
// created: 2018-06-04
// oauth2 server config

package o2

type ServerConfig struct {
	// oauth2 server name, will be show in login and authorize page
	ServerName string

	// favicon url
	Favicon string

	// logo url
	Logo string

	// uri context
	UriContext string

	// uri prefix to add before authRedirect uri
	UriPrefix string

	// template path prefix
	TemplatePrefix string
}


// ---------------------------
func DefaultServerConfig() *ServerConfig {
	if defaultOauth2Cfg == nil {
		defaultOauth2Cfg = &ServerConfig{
			UriPrefix:      "",
			UriContext:     "/oauth2",
			TemplatePrefix: "./",
			ServerName:     "Oauth2 Server",
			Logo:           "https://oauth.net/images/oauth-2-sm.png",
			Favicon:        "https://oauth.net/images/oauth-logo-square.png",
		}
	}
	return defaultOauth2Cfg
}
