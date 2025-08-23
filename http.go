package main

import (
	"context"
	_ "embed"
	"net/http"
	"strconv"
	"time"

	"github.com/siiimooon/go-polar/pkg/h10"
)

//go:embed index.html
var indexHtml []byte
var indexHtmlContentLength = strconv.Itoa(len(indexHtml))

const defaultAddr = "localhost:4880"

func serveIndexHtml(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", indexHtmlContentLength)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexHtml)
}

func serveHttp(hr chan h10.HeartRateMeasurement, addr string, ctx context.Context) error {
	ws := WSServer{hr}
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndexHtml)
	mux.Handle("/hr", &ws)

	s := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = s.Shutdown(shutdownCtx)
	}()

	return s.ListenAndServe()
}
