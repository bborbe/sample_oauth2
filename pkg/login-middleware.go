package pkg

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/golang/glog"
)

const (
	LoginCookieName = "X-Gateway-User"
	LoginHeaderName = "X-Gateway-User"
)

type LoginMiddleware interface {
	Middleware(handler http.Handler) http.Handler
}

// NewLoginMiddleware for validating request against a jwt secret
func NewLoginMiddleware(
	cookieGenerator CookieGenerator,
	stateGenerator StateGenerator,
	googleOAuth GoogleOAuth,
	callbackPath string,
) LoginMiddleware {
	return &loginMiddleware{
		stateGenerator:  stateGenerator,
		cookieGenerator: cookieGenerator,
		googleOAuth:     googleOAuth,
		callbackPath:    callbackPath,
	}
}

type loginMiddleware struct {
	cookieGenerator CookieGenerator
	stateGenerator  StateGenerator
	googleOAuth     GoogleOAuth
	callbackPath    string
}

func (l *loginMiddleware) Middleware(handler http.Handler) http.Handler {
	return libhttp.NewErrorHandler(libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
		glog.V(2).Infof("login middleware started with url %s", req.URL.String())
		if err := l.authenticate(ctx, req); err != nil {
			if err := l.login(ctx, resp, req); err != nil {
				return errors.Wrapf(ctx, err, "redirect to google login failed")
			}
			glog.V(2).Infof("login redirect completed")
			return nil
		}
		glog.V(2).Infof("user is authenticated")
		handler.ServeHTTP(resp, req)
		glog.V(2).Infof("login middleware completed")
		return nil
	}))
}

func (l *loginMiddleware) authenticate(ctx context.Context, req *http.Request) error {
	if req.URL.Path == l.callbackPath {
		glog.V(2).Info("skip auth for callback")
		return nil
	}

	cookie, err := req.Cookie(LoginCookieName)
	if err != nil || cookie.Value == "" {
		return errors.Wrap(ctx, err, "invalid auth cookie")
	}

	secureCookie, err := l.cookieGenerator.Decode(ctx, cookie.Value)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid auth cookie")
	}
	req.Header.Set(LoginHeaderName, secureCookie.Subject)

	glog.V(2).Infof("user %s is authenticated", secureCookie.Subject)
	return nil
}

func (l *loginMiddleware) login(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
	state, err := l.stateGenerator.Generate(ctx, req.URL.String())
	if err != nil {
		return errors.Wrapf(ctx, err, "generate state failed")
	}
	url := l.googleOAuth.AuthCodeURL(state)
	glog.V(3).Infof("redirect url '%s'", url)
	http.Redirect(resp, req, url, http.StatusTemporaryRedirect)
	return nil
}
