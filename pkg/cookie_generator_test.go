package pkg_test

import (
	"context"
	"time"

	"github.com/bborbe/sample_oauth2/pkg"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CookieGenerator", func() {
	var signingKey = []byte("test-key")
	var cookieGenerator = pkg.NewCookieGenerator(signingKey)
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	It("generates complete token", func() {
		user := "jdoe@example.com"
		cookie, err := cookieGenerator.Generate(ctx, user)
		Expect(err).To(BeNil())
		Expect(cookie.Subject).To(BeEquivalentTo(user))
		Expect(cookie.Id).NotTo(BeEmpty())
		Expect(cookie.String()).NotTo(BeEmpty())
		Expect(time.Unix(cookie.IssuedAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(cookie.NotBefore, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(cookie.ExpiresAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().AddDate(0, 0, 1).Unix(), 0)))
	})
	It("generates valid token", func() {
		user := "jdoe@example.com"
		cookie, err := cookieGenerator.Generate(ctx, user)
		Expect(err).To(BeNil())
		Expect(cookie.String()).NotTo(BeEmpty())
		cookie, err = cookieGenerator.Decode(ctx, cookie.String())
		Expect(err).To(BeNil())
		Expect(cookie.Subject).To(BeEquivalentTo(user))
		Expect(cookie.Id).NotTo(BeEmpty())
		Expect(time.Unix(cookie.IssuedAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(cookie.NotBefore, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(cookie.ExpiresAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().AddDate(0, 0, 1).Unix(), 0)))
	})
	It("returns error when decoding outdated token", func() {
		user := "jdoe@example.com"
		cookie, err := cookieGenerator.Generate(ctx, user)
		Expect(err).To(BeNil())
		Expect(cookie.String()).NotTo(BeEmpty())

		cookie.IssuedAt = time.Unix(cookie.IssuedAt, 0).Add(time.Duration(-25) * time.Hour).Unix()
		cookie.NotBefore = time.Unix(cookie.NotBefore, 0).Add(time.Duration(-25) * time.Hour).Unix()
		cookie.ExpiresAt = time.Unix(cookie.ExpiresAt, 0).Add(time.Duration(-25) * time.Hour).Unix()

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, cookie).SignedString(signingKey)
		Expect(err).To(BeNil())
		Expect(token).NotTo(BeEmpty())

		cookie, err = cookieGenerator.Decode(ctx, token)
		Expect(err).NotTo(BeNil())
	})
	It("returns error when decoding invalid string", func() {
		raw := "0123456789"
		cookie, err := cookieGenerator.Decode(ctx, raw)
		Expect(err).NotTo(BeNil())
		Expect(cookie).To(BeEquivalentTo(pkg.Cookie{}))
	})
})
