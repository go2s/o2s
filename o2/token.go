// authors: wangoo
// created: 2018-05-29
// token config

package o2

import (
	"time"

	"github.com/go2s/oauth2/manage"
)

func defaultTokenConfig(manager *manage.Manager) {
	// ------------------------------
	// SetImplicitTokenCfg set the implicit grant token config
	cfg := &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 1,
	}
	manager.SetImplicitTokenCfg(cfg)

	// ------------------------------
	// SetAuthorizeCodeTokenCfg set the authorization code grant token config
	cfg = &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
		// refresh token expiration time
		RefreshTokenExp: time.Hour * 24 * 3,
		// whether to generate the refreshing token
		IsGenerateRefresh: true,
	}
	manager.SetAuthorizeCodeTokenCfg(cfg)

	// ------------------------------
	// SetPasswordTokenCfg set the password grant token config
	cfg = &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
		// refresh token expiration time
		RefreshTokenExp: time.Hour * 24 * 7,
		// whether to generate the refreshing token
		IsGenerateRefresh: true,
	}
	manager.SetPasswordTokenCfg(cfg)

	// ------------------------------
	// SetClientTokenCfg set the client grant token config
	cfg = &manage.Config{
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
	}
	manager.SetClientTokenCfg(cfg)

	// ------------------------------
	// SetRefreshTokenCfg set the refreshing token config
	refCfg := &manage.RefreshingConfig{
		// whether to generate the refreshing token
		IsGenerateRefresh: true,
		// access token expiration time
		AccessTokenExp: time.Hour * 2,
		// refresh token expiration time
		RefreshTokenExp: time.Hour * 24 * 3,
		// whether to reset the refreshing create time
		IsResetRefreshTime: false,
		// whether to remove access token
		IsRemoveAccess: true,
		// whether to remove refreshing token
		IsRemoveRefreshing: true,
	}
	manager.SetRefreshTokenCfg(refCfg)
}
