package pkg_test

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/sample_oauth2/pkg"
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
		Expect(cookie.ID).NotTo(BeEmpty())
		Expect(cookie.String()).NotTo(BeEmpty())
		Expect(cookie.IssuedAt.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(cookie.NotBefore.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(cookie.ExpiresAt.Time).To(BeTemporally(">=", time.Unix(time.Now().AddDate(0, 0, 1).Unix(), 0)))
	})
	It("generates valid token", func() {
		user := "jdoe@example.com"
		cookie, err := cookieGenerator.Generate(ctx, user)
		Expect(err).To(BeNil())
		Expect(cookie.String()).NotTo(BeEmpty())
		cookie, err = cookieGenerator.Decode(ctx, cookie.String())
		Expect(err).To(BeNil())
		Expect(cookie.Subject).To(BeEquivalentTo(user))
		Expect(cookie.ID).NotTo(BeEmpty())
		Expect(cookie.IssuedAt.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(cookie.NotBefore.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(cookie.ExpiresAt.Time).To(BeTemporally(">=", time.Unix(time.Now().AddDate(0, 0, 1).Unix(), 0)))
	})
	It("returns error when decoding outdated token", func() {
		user := "jdoe@example.com"
		cookie, err := cookieGenerator.Generate(ctx, user)
		Expect(err).To(BeNil())
		Expect(cookie.String()).NotTo(BeEmpty())

		cookie.IssuedAt = jwt.NewNumericDate(cookie.IssuedAt.Add(time.Duration(-25) * time.Hour))
		cookie.NotBefore = jwt.NewNumericDate(cookie.NotBefore.Add(time.Duration(-25) * time.Hour))
		cookie.ExpiresAt = jwt.NewNumericDate(cookie.ExpiresAt.Add(time.Duration(-25) * time.Hour))

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
