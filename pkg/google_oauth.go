package pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bborbe/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// UserInfo returned by the OAuth provider
type UserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}

// Code used for authorization
type Code string

func (c Code) String() string {
	return string(c)
}

// GoogleOAuth defines the interface used for running a Google OAuth flow
//
//counterfeiter:generate -o mocks/login-google-oauth.go --fake-name GoogleOAuth . GoogleOAuth
type GoogleOAuth interface {
	AuthCodeURL(state State) string
	UserInfo(ctx context.Context, code Code) (*UserInfo, error)
}

// NewGoogleOAuth returns an implementation of the Google OAuth flow using the provided credentials
func NewGoogleOAuth(
	clientID string,
	clientSecret string,
	redirectURL string,
	hostedDomain string,
) GoogleOAuth {
	return &googleOAuth{
		config: oauth2.Config{
			RedirectURL:  redirectURL,
			ClientID:     strings.ReplaceAll(clientID, "client_id: ", ""),
			ClientSecret: clientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint: google.Endpoint,
		},
		hostedDomain: hostedDomain,
	}
}

type googleOAuth struct {
	config       oauth2.Config
	hostedDomain string
}

// AuthCodeURL returns the auth code url for the provided state
func (o *googleOAuth) AuthCodeURL(state State) string {
	return o.config.AuthCodeURL(state.String(),
		oauth2.SetAuthURLParam("hd", o.hostedDomain),
	)
}

// UserInfo retrieves the UserInfo for the provided auth code
func (o *googleOAuth) UserInfo(ctx context.Context, code Code) (*UserInfo, error) {

	token, err := o.config.Exchange(ctx, code.String())
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "code exchange failed")
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, errors.Wrapf(ctx, err, "failed getting user info")
	}

	defer response.Body.Close()
	var data UserInfo
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, errors.Wrapf(ctx, err, "decode json failed")
	}
	return &data, nil
}
