package auth

import (
	"bytes"
	"context"
	"testing"

	"github.com/Rindrics/go_todo_app/clock"
	"github.com/Rindrics/go_todo_app/entity"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN RSA PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPrivKey)
	}
}

func TestJWTer_GenerateToken(t *testing.T) {
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(1234)
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, clock.RealClocker{})
	if err != nil {
		t.Fatal(err)
	}

	got, err := sut.GenerateToken(ctx, wantID)
	if err != nil {
		t.Fatalf("want no error, got %s", err)
	}

	if len(got) == 0 {
		t.Errorf("token is empty")
	}
}
