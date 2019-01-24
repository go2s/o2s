// authors: wangoo
// created: 2018-06-05
// oauth2 valid

package o2x

type ValidResponse struct {
	ClientID string `json:"client_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Scope    string `json:"scope,omitempty"`
}
