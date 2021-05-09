package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func httpServer() error {
	return http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}))
}

func startHttpServer(g *errgroup.Group, ctx context.Context) {
	g.Go(func() error {
		return httpServer()
	})
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("ctx.Done,http server exited!")
			return
		}
	}
}

func signalNotify(g *errgroup.Group, ctx context.Context) {
	//创建一个os.Signal channel
	sig := make(chan os.Signal)

	//注册要接收的信号，syscall.SIGINT:接收ctrl+c ,syscall.SIGTERM:程序退出
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("ctx.Done")
				return nil
			case s := <-sig:
				fmt.Printf("received sig: %s\n", s.String())
				return fmt.Errorf("received sig %s", s.String())
			}
		}
	})
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	signalNotify(g, ctx)
	startHttpServer(g, ctx)
}
