package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/weicheng95/go-mongo-template/internal/envvar"
	"github.com/weicheng95/go-mongo-template/pkg/logger"
)

var (
	log = logger.NewLogger("go-mango-server")
)

func init() {
	// set output level based on environment
	log.SetOutputLevel(logger.InfoLevel)
}

func main() {

	var env, address string

	flag.StringVar(&env, "env", ".env", "Environment Variables filename")
	flag.StringVar(&address, "address", ":8000", "HTTP Server Address")
	flag.Parse()
	log.WithLogFields(logger.Fields{
		"animal": "walrus",
	}).Info("init successful")
	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}

}

func run(env, address string) (<-chan error, error) {
	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.WithLogFields(logger.Fields{
				"method": r.Method,
				"url":    r.URL.String(),
			}).Info()

			h.ServeHTTP(w, r)
		})
	}

	//-

	errC := make(chan error, 1)

	server := newServer(address, logging)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		// this will trigger when Shutdown signal received
		<-ctx.Done()

		log.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			stop()
			cancel()
			close(errC)
		}()

		server.SetKeepAlivesEnabled(false)

		if err := server.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		log.Info("Shutdown completed")
	}()

	go func() {
		log.Info("Listening and serving ", address)

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

func newServer(address string, mws ...mux.MiddlewareFunc) *http.Server {
	r := mux.NewRouter()

	for _, mw := range mws {
		r.Use(mw)
	}

	//-

	// repo := postgresql.NewTask(db) // Task Repository
	// svc := service.NewTask(repo)   // Task Application Service

	// r.Handle("/metrics", metrics)

	//-

	return &http.Server{
		Handler:           r,
		Addr:              address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}
}
