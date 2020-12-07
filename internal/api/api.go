package api

import (
	"net/http"
	"os"
	"time"

	"github.com/biensupernice/krane/internal/api/controllers"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/biensupernice/krane/internal/api/middlewares"
)

// Run : start the Krane api
func Run() {
	logger.Debugf("Starting Krane API on pid: %d", os.Getpid())
	router := mux.NewRouter()

	withBaseMiddlewares(router)
	withRoutes(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         os.Getenv(constants.EnvListenAddress),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("API on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}

// withBaseMiddlewares : configure api middlewares
func withBaseMiddlewares(router *mux.Router) {
	router.Use(middlewares.Logging)
	router.Use(handlers.RecoveryHandler())
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
		handlers.AllowedOrigins([]string{"*"}))) // TODO: Not allowing wild card origins (*) use envar LISTEN_ADDRESS
}

// withRoutes : configure api routes
func withRoutes(router *mux.Router) {
	noAuthRouter := router.PathPrefix("/").Subrouter()
	withRoute(noAuthRouter, "/", controllers.GetServerStatus).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/login", controllers.RequestLoginPhrase).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/auth", controllers.AuthenticateClientJWT).Methods(http.MethodPost)

	authRouter := router.PathPrefix("/").Subrouter()
	// deployments
	withRoute(authRouter, "/deployments", controllers.GetAllDeployments, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/deployments", controllers.ApplyDeployment, middlewares.ValidateSessionMiddleware).Methods(http.MethodPost)
	withRoute(authRouter, "/deployments/{name}", controllers.GetDeployment, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/deployments/{name}", controllers.DeleteDeployment, middlewares.ValidateSessionMiddleware).Methods(http.MethodDelete)
	withRoute(authRouter, "/deployments/{name}/containers", controllers.GetContainers, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	// secrets
	withRoute(authRouter, "/secrets/{name}", controllers.GetSecrets, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/secrets/{name}", controllers.CreateSecret, middlewares.ValidateSessionMiddleware).Methods(http.MethodPost)
	withRoute(authRouter, "/secrets/{name}/{key}", controllers.DeleteSecret, middlewares.ValidateSessionMiddleware).Methods(http.MethodDelete)
	// jobs
	withRoute(authRouter, "/jobs", controllers.GetRecentJobs, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/jobs/{namespace}", controllers.GetJobsByNamespace, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/jobs/{namespace}/{id}", controllers.GetJobByID, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
	// session
	withRoute(authRouter, "/sessions", controllers.GetSessions, middlewares.ValidateSessionMiddleware).Methods(http.MethodGet)
}

type routeHandler func(http.ResponseWriter, *http.Request)

func withRoute(r *mux.Router, path string, handler routeHandler, middlewares ...mux.MiddlewareFunc) *mux.Route {
	for _, mw := range middlewares {
		r.Use(mw)
	}
	return r.HandleFunc(path, handler)
}
