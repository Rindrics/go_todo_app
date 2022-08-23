package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Rindrics/go_todo_app/clock"
	"github.com/Rindrics/go_todo_app/entity"
	"github.com/Rindrics/go_todo_app/store"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JWTer struct {
	PrivateKey string
	Store      store.KVS
	Clocker    clock.Clocker
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
