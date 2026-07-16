package pkg_test

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/sample_oauth2/pkg"
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
		Expect(state.IssuedAt.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(state.NotBefore.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(state.ExpiresAt.Time).To(BeTemporally(">", time.Unix(time.Now().Unix(), 0)))
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
		Expect(state.IssuedAt.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(state.NotBefore.Time).To(BeTemporally(">=", time.Unix(time.Now().Unix(), 0)))
		Expect(state.ExpiresAt.Time).To(BeTemporally(">", time.Unix(time.Now().Unix(), 0)))
	})
	It("returns error when decoding outdated token", func() {
		origin := "https://test.localhost/foo"
		state, err := stateGenerator.Generate(ctx, origin)
		Expect(err).To(BeNil())
		Expect(state.String()).NotTo(BeEmpty())

		state.IssuedAt = jwt.NewNumericDate(state.IssuedAt.Add(time.Duration(-2) * time.Minute))
		state.ExpiresAt = jwt.NewNumericDate(state.ExpiresAt.Add(time.Duration(-2) * time.Minute))
		state.NotBefore = jwt.NewNumericDate(state.NotBefore.Add(time.Duration(-2) * time.Minute))

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
