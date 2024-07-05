package pkg_test

import (
	"context"
	"time"

	"github.com/bborbe/sample_oauth2/pkg"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StateGenerator", func() {
	var signingKey = []byte("test-key")
	var stateGenerator = pkg.NewStateGenerator(signingKey)
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()
	})
	It("generates complete token", func() {
		origin := "https://test.localhost/foo"
		state, err := stateGenerator.Generate(ctx, origin)
		Expect(err).To(BeNil())
		Expect(state.Origin).To(BeEquivalentTo(origin))
		Expect(state.Subject).NotTo(BeEmpty())
		Expect(state.String()).NotTo(BeEmpty())
		Expect(time.Unix(state.IssuedAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(state.NotBefore, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(state.ExpiresAt, 0)).To(BeTemporally(">", time.Unix(time.Now().Unix(), 0)))
	})
	It("generates valid token", func() {
		origin := "https://test.localhost/foo"
		state, err := stateGenerator.Generate(ctx, origin)
		Expect(err).To(BeNil())
		Expect(state.String()).NotTo(BeEmpty())
		state, err = stateGenerator.Decode(ctx, state.String())
		Expect(err).To(BeNil())
		Expect(state.Origin).To(BeEquivalentTo(origin))
		Expect(state.Subject).NotTo(BeEmpty())
		Expect(time.Unix(state.IssuedAt, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(state.NotBefore, 0)).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(time.Unix(state.ExpiresAt, 0)).To(BeTemporally(">", time.Unix(time.Now().Unix(), 0)))
	})
	It("returns error when decoding outdated token", func() {
		origin := "https://test.localhost/foo"
		state, err := stateGenerator.Generate(ctx, origin)
		Expect(err).To(BeNil())
		Expect(state.String()).NotTo(BeEmpty())

		state.IssuedAt = time.Unix(state.IssuedAt, 0).Add(time.Duration(-2) * time.Minute).Unix()
		state.ExpiresAt = time.Unix(state.ExpiresAt, 0).Add(time.Duration(-2) * time.Minute).Unix()
		state.NotBefore = time.Unix(state.NotBefore, 0).Add(time.Duration(-2) * time.Minute).Unix()

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, state).
			SignedString(signingKey)
		Expect(err).To(BeNil())
		Expect(token).NotTo(BeEmpty())

		state, err = stateGenerator.Decode(ctx, token)
		Expect(err).NotTo(BeNil())
	})
	It("returns error when decoding invalid string", func() {
		raw := "0123456789"
		state, err := stateGenerator.Decode(ctx, raw)
		Expect(err).NotTo(BeNil())
		Expect(state).To(BeEquivalentTo(pkg.State{}))
	})
})
