package pkg

import (
	"context"
	"net/http"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/golang-jwt/jwt"
	"github.com/golang/glog"
)

const (
	LoginCookieName = "X-Gateway-User"
	LoginHeaderName = "X-Gateway-User"
)

type LoginService interface {
	Middleware(handler http.Handler) http.Handler
	libhttp.WithError
}

// NewLoginService for validating request against a jwt secret
func NewLoginService(
	cookieGenerator CookieGenerator,
	stateGenerator StateGenerator,
	googleOAuth GoogleOAuth,
	callbackPath string,
) LoginService {
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
		if err := l.expectedAuthenicated(ctx, req); err != nil {
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

func (l *loginMiddleware) expectedAuthenicated(ctx context.Context, req *http.Request) error {
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

func (l *loginMiddleware) handleCallback(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("logging in via google")
	user, origin, err := l.validateGoogleAuthCallback(ctx, req)
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) && validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			glog.V(1).Infof("callback state expired, reinitializing auth flow: %v", err)
			return errors.Wrapf(ctx, err, "callback expired")
		}
		glog.V(1).Infof("callback validation failed: %v", err)
		return errors.Wrapf(ctx, err, "callback validation failed")
	}

	cookie, err := l.cookieGenerator.Generate(ctx, user)
	if err != nil {
		glog.V(1).Infof("generate cookie for %s failed", user)
		return errors.Wrapf(ctx, err, "generating cookie failed")
	}

	glog.V(2).Infof("set X-Gateway-User to %s", user)
	http.SetCookie(resp, cookie.HTTPCookie())
	glog.V(2).Infof("redirect to %s", origin)
	http.Redirect(resp, req, origin, http.StatusTemporaryRedirect)
	return nil
}

func (l *loginMiddleware) validateGoogleAuthCallback(ctx context.Context, req *http.Request) (string, string, error) {
	if err := req.ParseForm(); err != nil {
		return "", "", errors.Wrapf(ctx, err, "parse form failed")
	}
	state, err := l.stateGenerator.Decode(ctx, req.Form.Get("state"))
	if err != nil {
		return "", "", errors.Wrapf(ctx, err, "invalid oauth state")
	}
	info, err := l.googleOAuth.UserInfo(ctx, Code(req.Form.Get("code")))
	if err != nil {
		return "", "", errors.Wrapf(ctx, err, "get user info failed")
	}
	return info.Email, state.Origin, nil
}

// ServeHTTP callback
func (l *loginMiddleware) ServeHTTP(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
	glog.V(2).Infof("logging in via google")
	user, origin, err := l.validateGoogleAuthCallback(ctx, req)
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) && validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			glog.V(1).Infof("callback state expired, reinitializing auth flow: %v", err)
			return errors.Wrapf(ctx, err, "callback expired")
		}
		glog.V(1).Infof("callback validation failed: %v", err)
		return errors.Wrapf(ctx, err, "callback validation failed")
	}

	cookie, err := l.cookieGenerator.Generate(ctx, user)
	if err != nil {
		glog.V(1).Infof("generate cookie for %s failed", user)
		return errors.Wrapf(ctx, err, "generating cookie failed")
	}

	glog.V(2).Infof("set X-Gateway-User to %s", user)
	http.SetCookie(resp, cookie.HTTPCookie())
	glog.V(2).Infof("redirect to %s", origin)
	http.Redirect(resp, req, origin, http.StatusTemporaryRedirect)
	return nil
}
