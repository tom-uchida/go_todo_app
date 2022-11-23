package config

import (
	"fmt"
	"testing"
)

// デフォルト値が設定されたか、設定した環境変数が期待通り設定されているのかを確認する。
func TestNew(t *testing.T) {
	wantPort := 3333
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if got.Port != wantPort {
		t.Errorf("want %d, but %d", wantPort, got.Port)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
}
