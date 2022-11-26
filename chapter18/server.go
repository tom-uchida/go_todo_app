package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { // 別ゴルーチンでHTTPサーバーを起動する
		// ListenAndServeメソッドではなく、Serveメソッドに変更する。
		if err := s.srv.Serve(s.l); err != nil &&
			// http.ErrServerClosedは
			// http.Server.Shutdown() が正常に終了したことを示すので異常ではない。
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done() // チャネルからの通知(終了通知)を待機する
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait() // グレースフルシャットダウンの終了を待つ
}
