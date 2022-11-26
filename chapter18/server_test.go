// 期待通りにHTTPサーバーが起動しているか
// テストコードから意図通りに終了するか
package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("fatal to listen port %v", err)
	}

	// キャンセル可能な「context.Context」にオブジェクトを作る。
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})

	// 別ゴルーチンでテスト対象の「run」関数を実行してHTTPサーバーを起動する。
	eg.Go(func() error {
		s := NewServer(l, mux)
		return s.Run(ctx)
	})

	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// ポート番号の確認
	t.Logf("try request to %q", url)
	// エンドポイントに対してGETリクエストを送信する。
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// HTTPサーバーの戻り値を検証する。
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}
	// 「cancel」関数を実行して、run関数に終了通知を送信する。
	cancel()

	// 「*errgroup.Group.Wait」メソッド経由で「run」関数の戻り値を検証する。
	// GETリクエストで取得したレスポンスボディが期待する文字列であることを検証する。
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
