package api

import (
	"net/http"
	"os"
	"time"

	"github.com/biensupernice/krane/internal/api/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api/middlewares"
)

func Run() {
	logrus.Infof("Starting Krane API on pid: %d", os.Getpid())
	router := mux.NewRouter()

	withBaseMiddlewares(router)
	withRoutes(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         os.Getenv("LISTEN_ADDRESS"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Krane API listening on %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logrus.Fatal(err.Error())
	}
}

func withBaseMiddlewares(router *mux.Router) {
	router.Use(middlewares.Logging)
	router.Use(handlers.RecoveryHandler())
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
		handlers.AllowedOrigins([]string{"localhost", "*"})))
}

func withRoutes(router *mux.Router) {
	noAuthRouter := router.PathPrefix("/").Subrouter()
	withRoute(noAuthRouter, "/", controllers.GetServerStatus).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/login", controllers.RequestLoginPhrase).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/auth", controllers.AuthenticateClientJWT).Methods(http.MethodPost)
	withRoute(noAuthRouter, "/ping", controllers.PingController).Methods(http.MethodGet)

	authRouter := router.PathPrefix("/").Subrouter()
	withRoute(authRouter, "/deployments", controllers.GetDeployments, middlewares.AuthSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/deployments", controllers.SaveDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodPost)
	withRoute(authRouter, "/deployments/{name}", controllers.GetDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/deployments/{name}", controllers.DeleteDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodDelete)
	withRoute(authRouter, "/sessions", controllers.GetSessions, middlewares.AuthSessionMiddleware).Methods(http.MethodGet)
}

type routeHandler func(http.ResponseWriter, *http.Request)

func withRoute(r *mux.Router, path string, handler routeHandler, middlewares ...mux.MiddlewareFunc) *mux.Route {
	for _, mw := range middlewares {
		r.Use(mw)
	}
	return r.HandleFunc(path, handler)
}
