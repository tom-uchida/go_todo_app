package main

import "net/http"

// 戻り値を*http.ServeMux型の値ではなく、
// http.Handlerインターフェースにしておくことで、
// 内部実装に依存しない関数シグネチャになる。
func NewMux() http.Handler {
	mux := http.NewServeMux()
	// HTTPサーバーが稼働中か確認するための/healthエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析のエラーを回避するため明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
	return mux
}
