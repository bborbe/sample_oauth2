// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/log"
	"github.com/bborbe/sample_oauth2/pkg"
	libsentry "github.com/bborbe/sentry"
	"github.com/bborbe/service"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	app := &application{}
	os.Exit(service.Main(context.Background(), app, &app.SentryDSN, &app.SentryProxy))
}

type application struct {
	SentryDSN          string `required:"true" arg:"sentry-dsn" env:"SENTRY_DSN" usage:"SentryDSN" display:"length"`
	SentryProxy        string `required:"false" arg:"sentry-proxy" env:"SENTRY_PROXY" usage:"Sentry Proxy"`
	Listen             string `required:"true" arg:"listen" env:"LISTEN" usage:"address to listen to"`
	GoogleClientID     string `required:"false" arg:"google-client-id" env:"GOOGLE_CLIENT_ID" usage:"Google client id"`
	GoogleClientSecret string `required:"false" arg:"google-client-secret" env:"GOOGLE_CLIENT_SECRET" usage:"Google client secret:" display:"length"`
	GoogleHostedDomain string `required:"false" arg:"google-hosted-domain" env:"GOOGLE_HOSTED_DOMAIN" usage:"Domain name of the Google Instance (G Suite)"`
	GoogleRedirectURL  string `required:"false" arg:"google-redirect-url" env:"GOOGLE_REDIRECT_URL" usage:"Google redirect url"`
	JWTSigningKey      string `required:"false" arg:"jwt-signing-key" env:"JWT_SIGNING_KEY" usage:"Key to use for signing jwts" display:"length"`
}

func (a *application) Run(ctx context.Context, sentryClient libsentry.Client) error {

	router := mux.NewRouter()
	router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
	router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
	router.Path("/metrics").Handler(promhttp.Handler())
	router.Path("/setloglevel/{level}").Handler(log.NewSetLoglevelHandler(ctx, log.NewLogLevelSetter(2, 5*time.Minute)))

	callbackUrl, err := url.Parse(a.GoogleRedirectURL)
	if err != nil {
		return errors.Wrapf(ctx, err, "parse callback url failed")
	}

	loginService := pkg.NewLoginService(
		pkg.NewCookieGenerator([]byte(a.JWTSigningKey)),
		pkg.NewStateGenerator([]byte(a.JWTSigningKey)),
		pkg.NewGoogleOAuth(
			a.GoogleClientID,
			a.GoogleClientSecret,
			a.GoogleRedirectURL,
			a.GoogleHostedDomain,
		),
		callbackUrl.Path,
	)
	router.Use(loginService.Middleware)
	router.Path(callbackUrl.Path).Handler(libhttp.NewErrorHandler(loginService))

	router.Path("/").Handler(libhttp.NewErrorHandler(libhttp.WithErrorFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) error {
		libhttp.WriteAndGlog(resp, "login success")
		return nil
	})))

	glog.V(2).Infof("starting http server listen on %s", a.Listen)
	return libhttp.NewServer(
		a.Listen,
		router,
	).Run(ctx)
}
