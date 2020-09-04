package api

import (
	"net/http"
	"os"
	"time"

	"github.com/biensupernice/krane/internal/api/routes"

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
	// router.HandleFunc("/", routes.GetServerStatus).Methods(status.MethodGet)
	// router.HandleFunc("/login", routes.RequestLoginPhrase).Methods(status.MethodGet)
	// router.HandleFunc("/auth", routes.AuthenticateClientJWT).Methods(status.MethodPost)
	// router.Handle("/spec", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.CreateSpec))).Methods(status.MethodPost)
	// router.Handle("/spec/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.UpdateSpec))).Methods(status.MethodPut)
	// router.Handle("/deployments", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.GetDeployments))).Methods(status.MethodGet)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.GetDeployment))).Methods(status.MethodGet)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.RunDeployment))).Methods(status.MethodPost)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.DeleteDeployment))).Methods(status.MethodDelete)
	// router.Handle("/deployments/{name}/stop", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.StopDeployment))).Methods(status.MethodPost)
	// router.Handle("/alias/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.UpdateDeploymentAlias))).Methods(status.MethodPost)
	// router.Handle("/alias/{name}", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.DeleteDeploymentAlias))).Methods(status.MethodDelete)
	// router.Handle("/activity", middlewares.AuthSessionMiddleware(status.HandlerFunc(routes.GetRecentActivity))).Methods(status.MethodGet)

	noAuthRouter := router.PathPrefix("/").Subrouter()
	withRoute(noAuthRouter, "/", routes.GetServerStatus).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/login", routes.RequestLoginPhrase).Methods(http.MethodGet)
	withRoute(noAuthRouter, "/auth", routes.AuthenticateClientJWT).Methods(http.MethodPost)

	authRouter := router.PathPrefix("/").Subrouter()
	withRoute(authRouter, "/deployment", routes.CreateDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodPost)
	withRoute(authRouter, "/deployment/{name}", routes.GetDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodGet)
	withRoute(authRouter, "/deployment/{name}", routes.DeleteDeployment, middlewares.AuthSessionMiddleware).Methods(http.MethodDelete)
}

type routeHandler func(http.ResponseWriter, *http.Request)

func withRoute(r *mux.Router, path string, handler routeHandler, middlewares ...mux.MiddlewareFunc) *mux.Route {
	for _, mw := range middlewares {
		r.Use(mw)
	}

	return r.HandleFunc(path, handler)
}
