package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Rindrics/go_todo_app/entity"
	"github.com/Rindrics/go_todo_app/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)

	sut := &KVS{Client: cli}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTest(t)
	sut := &KVS{Client: cli}

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})

		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Errorf("want no error, but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		key := "TestKVS_Load_notFound"

		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value = %d)", ErrNotFound, err, got)
		}
	})
}
