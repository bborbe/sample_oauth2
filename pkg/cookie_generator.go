package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Cookie storing the user pass-through information that is passed on to ERPNext for authentication.
type Cookie struct {
	jwt.StandardClaims

	token string
}

func (s Cookie) String() string {
	return s.token
}

// HTTPCookie based on Cookie
func (s Cookie) HTTPCookie() *http.Cookie {
	return &http.Cookie{
		Name:  LoginCookieName,
		Path:  "/",
		Value: s.String(),
	}
}

// CookieGenerator generates and decodes secure cookies
//
//counterfeiter:generate -o ../mocks/cookie-generator.go --fake-name CookieGenerator . CookieGenerator
type CookieGenerator interface {
	Generate(ctx context.Context, user string) (Cookie, error)
	Decode(ctx context.Context, cookie string) (Cookie, error)
}

// NewCookieGenerator using key to sign cookie tokens
func NewCookieGenerator(key []byte) CookieGenerator {
	return &cookieGenerator{key: key}
}

type cookieGenerator struct {
	key []byte
}

// Generate a signed cookie
func (s *cookieGenerator) Generate(ctx context.Context, user string) (Cookie, error) {
	issuedAt := time.Now().UTC()
	generateUUID, err := uuid.NewUUID()
	if err != nil {
		return Cookie{}, err
	}

	cookie := Cookie{
		StandardClaims: jwt.StandardClaims{
			Id:        generateUUID.String(),
			Subject:   user,
			IssuedAt:  issuedAt.Unix(),
			NotBefore: issuedAt.Unix(),
			ExpiresAt: issuedAt.Add(24 * time.Hour).Unix(),
		},
	}

	return s.sign(cookie)
}

func (s *cookieGenerator) sign(cookie Cookie) (Cookie, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cookie)
	signed, err := token.SignedString(s.key)
	if err != nil {
		return Cookie{}, err
	}

	cookie.token = signed
	return cookie, nil
}

// Decode a cookie string and validate it
func (s *cookieGenerator) Decode(ctx context.Context, cookie string) (Cookie, error) {
	token, err := jwt.ParseWithClaims(cookie, &Cookie{}, s.keyFunc)
	if err != nil {
		return Cookie{}, err
	}

	if claims, ok := token.Claims.(*Cookie); ok && token.Valid {
		if len(claims.Subject) < 1 {
			return Cookie{}, errors.New("token invalid")
		}
		return *claims, nil
	}

	return Cookie{}, errors.New("token invalid")
}

func (s *cookieGenerator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	return s.key, nil
}
