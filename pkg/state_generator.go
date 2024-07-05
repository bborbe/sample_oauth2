package pkg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// State stores a requests state for passing through the oauth2 flow,
// ensuring CSRF protection and a fluent experience by passing the origin url.
type State struct {
	Origin string `json:"origin"`
	jwt.StandardClaims

	token string
}

func (s State) String() string {
	return s.token
}

// StateGenerator generates and decodes secure states
//
//counterfeiter:generate -o ../mocks/state-generator.go --fake-name StateGenerator . StateGenerator
type StateGenerator interface {
	Generate(ctx context.Context, originURL string) (State, error)
	Decode(ctx context.Context, token string) (State, error)
}

// NewStateGenerator using key to sign state tokens
func NewStateGenerator(key []byte) StateGenerator {
	return &stateGenerator{key: key}
}

type stateGenerator struct {
	key []byte
}

// Generate a signed state
func (s *stateGenerator) Generate(ctx context.Context, originURL string) (State, error) {
	issuedAt := time.Now().UTC()
	generateUUID, err := uuid.NewUUID()
	if err != nil {
		return State{}, err
	}

	state := State{
		Origin: originURL,
		StandardClaims: jwt.StandardClaims{
			Subject:   generateUUID.String(),
			IssuedAt:  issuedAt.Unix(),
			NotBefore: issuedAt.Unix(),
			ExpiresAt: issuedAt.Add(1 * time.Minute).Unix(),
		},
	}

	return s.sign(state)
}

func (s *stateGenerator) sign(state State) (State, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, state)
	signed, err := token.SignedString(s.key)
	if err != nil {
		return State{}, err
	}

	state.token = signed
	return state, nil
}

// Decode a state string and validate it
func (s *stateGenerator) Decode(ctx context.Context, state string) (State, error) {
	token, err := jwt.ParseWithClaims(state, &State{}, s.keyFunc)
	if err != nil {
		return State{}, err
	}

	if claims, ok := token.Claims.(*State); ok && token.Valid {
		if len(claims.Subject) < 1 {
			return State{}, errors.New("token invalid")
		}
		return *claims, nil
	}

	return State{}, errors.New("token invalid")
}

func (s *stateGenerator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}

	return s.key, nil
}
