package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/weicheng95/go-mongo-template/initiator"
	"github.com/weicheng95/go-mongo-template/internal/common/server"
	"github.com/weicheng95/go-mongo-template/pkg/helper"
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

	var env, configFile, address string

	flag.StringVar(&env, "env", "development", "Environment Variables filename")
	flag.StringVar(&configFile, "configFile", "config.yaml", "Environment Variables filename")
	flag.StringVar(&address, "address", ":8000", "HTTP Server Address")
	flag.Parse()

	if err := helper.Load(env, configFile); err != nil {
		log.Fatalf("%w", err)
	}

	// database init
	dbClient := helper.NewMongoDBClient()

	userModule := initiator.UserRestInit(dbClient)

	// auth module (register, login)
	r := chi.NewRouter()
	r.Group(func (r chi.Router) {
		r.Mount("/auth", userModule.AuthRoutes())
	})

	// user module
	r.Group(func (r chi.Router) {
		r.Use(server.AuthMiddleware{}.Middleware)
		r.Mount("/user", userModule.UserRoutes())
	})

	server.RunHTTPServer(address, r)
}

//func run(env, configFile string, address string) (<-chan error, error) {
//
//
//	//logging := func(h http.Handler) http.Handler {
//	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	//		log.WithLogFields(logger.Fields{
//	//			"method": r.Method,
//	//			"url":    r.URL.String(),
//	//		}).Info()
//	//
//	//		h.ServeHTTP(w, r)
//	//	})
//	//}
//
//
//
//	// server init
//	//server := newServer(address, dbClient, logging)
//
//	errC := make(chan error, 1)
//
//	ctx, stop := signal.NotifyContext(context.Background(),
//		os.Interrupt,
//		syscall.SIGTERM,
//		syscall.SIGQUIT)
//
//	go func() {
//		// this will trigger when Shutdown signal received
//		<-ctx.Done()
//
//		log.Info("Shutdown signal received")
//
//		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//
//		defer func() {
//			stop()
//			cancel()
//			close(errC)
//		}()
//
//		server.SetKeepAlivesEnabled(false)
//
//		if err := server.Shutdown(ctxTimeout); err != nil {
//			errC <- err
//		}
//
//		log.Info("Shutdown completed")
//	}()
//
//	go func() {
//		log.Info("Listening and serving ", address)
//
//		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
//		// ErrServerClosed."
//		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
//			errC <- err
//		}
//	}()
//
//	return errC, nil
//}
//
//func newServer(address string, client *mongo.Client, mws ...mux.MiddlewareFunc) *http.Server {
//	r := mux.NewRouter()
//
//	for _, mw := range mws {
//		r.Use(mw)
//	}
//
//	userModule := initiator.UserRestInit(client)
//	userModule.Register(r)
//
//	return &http.Server{
//		Handler:           r,
//		Addr:              address,
//		ReadTimeout:       1 * time.Second,
//		ReadHeaderTimeout: 1 * time.Second,
//		WriteTimeout:      1 * time.Second,
//		IdleTimeout:       1 * time.Second,
//	}
//}
