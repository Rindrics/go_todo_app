package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	_ "embed"

	"github.com/Rindrics/go_todo_app/clock"
	"github.com/Rindrics/go_todo_app/entity"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	Store                 Store
	Clocker               clock.Clocker
}

func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{
		Store:   s,
		Clocker: c,
	}

	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}

	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: publice key: %w", err)
	}

	j.PrivateKey = privkey
	j.PublicKey = pubkey

	return j, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, userID entity.UserID) ([]byte, error) {
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer("github.com/Rindrics/go_todo_app").
		Subject("access token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30 * time.Minute)).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GetToken: failed to build token: %w", err)
	}

	if err := j.Store.Save(ctx, tok.JwtID(), userID); err != nil {
		return nil, err
	}

	return jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
}

func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.RS256, j.PublicKey),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}

	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}

	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}

	return token, nil
}
