package pkg

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/golang/glog"
)

func NewLoginCallbackHandler(
	cookieGenerator CookieGenerator,
	stateGenerator StateGenerator,
	googleOAuth GoogleOAuth,
) libhttp.WithError {
	return libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return errors.Wrapf(ctx, err, "parse form failed")
		}
		state, err := stateGenerator.Decode(ctx, req.Form.Get("state"))
		if err != nil {
			return errors.Wrapf(ctx, err, "invalid oauth state")
		}
		info, err := googleOAuth.UserInfo(ctx, Code(req.Form.Get("code")))
		if err != nil {
			return errors.Wrapf(ctx, err, "get user info failed")
		}
		user := info.Email
		origin := state.Origin

		cookie, err := cookieGenerator.Generate(ctx, user)
		if err != nil {
			glog.V(1).Infof("generate cookie for %s failed", user)
			return errors.Wrapf(ctx, err, "generating cookie failed")
		}

		glog.V(2).Infof("set X-Gateway-User to %s", user)
		http.SetCookie(resp, cookie.HTTPCookie())
		glog.V(2).Infof("redirect to %s", origin)
		http.Redirect(resp, req, origin, http.StatusTemporaryRedirect)
		return nil
	})

}
