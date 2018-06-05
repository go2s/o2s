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

	// uri prefix to add before redirect uri
	UriPrefix string

	// template path prefix
	TemplatePrefix string
}
