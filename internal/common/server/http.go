package server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"github.com/weicheng95/go-mongo-template/internal/common/server/logging"
	"net/http"
)

func RunHTTPServer(addr string, handler http.Handler) {
	rootRouter := chi.NewRouter()
	setMiddlewares(rootRouter)
	// we are mounting all APIs under /api path
	rootRouter.Mount("/api", handler)

	logrus.Info("Starting HTTP server")

	http.ListenAndServe(addr, rootRouter)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(logs.NewStructuredLogger(logrus.StandardLogger()))
	router.Use(middleware.Recoverer)

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}
	router.Use(logging.NewStructuredLogger(logger))

	//addCorsMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosniff"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func AddAuthMiddleware(router *chi.Mux) {
	router.Use(AuthMiddleware{}.Middleware)
}

//
//func addCorsMiddleware(router *chi.Mux) {
//	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
//	if len(allowedOrigins) == 0 {
//		return
//	}
//
//	corsMiddleware := cors.New(cors.Options{
//		AllowedOrigins:   allowedOrigins,
//		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
//		ExposedHeaders:   []string{"Link"},
//		AllowCredentials: true,
//		MaxAge:           300,
//	})
//	router.Use(corsMiddleware.Handler)
//}
