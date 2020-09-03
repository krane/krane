package api

import (
	"net/http"
	"os"
	"time"

	"github.com/biensupernice/krane/api/routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/api/middlewares"
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

	logrus.Infof("Krane server listening on %s", srv.Addr)
	logrus.Fatal(srv.ListenAndServe())
}

func withBaseMiddlewares(router *mux.Router) {
	router.Use(middlewares.Logging)
	router.Use(handlers.RecoveryHandler())
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost}),
		handlers.AllowedOrigins([]string{"localhost", "*"})))
}

func withRoutes(router *mux.Router) {
	// router.HandleFunc("/", routes.GetServerStatus).Methods(http.MethodGet)
	// router.HandleFunc("/login", routes.RequestLoginPhrase).Methods(http.MethodGet)
	// router.HandleFunc("/auth", routes.AuthenticateClientJWT).Methods(http.MethodPost)
	// router.Handle("/spec", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.CreateSpec))).Methods(http.MethodPost)
	// router.Handle("/spec/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.UpdateSpec))).Methods(http.MethodPut)
	// router.Handle("/deployments", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.GetDeployments))).Methods(http.MethodGet)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.GetDeployment))).Methods(http.MethodGet)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.RunDeployment))).Methods(http.MethodPost)
	// router.Handle("/deployments/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.DeleteDeployment))).Methods(http.MethodDelete)
	// router.Handle("/deployments/{name}/stop", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.StopDeployment))).Methods(http.MethodPost)
	// router.Handle("/alias/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.UpdateDeploymentAlias))).Methods(http.MethodPost)
	// router.Handle("/alias/{name}", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.DeleteDeploymentAlias))).Methods(http.MethodDelete)
	// router.Handle("/activity", middlewares.AuthSessionMiddleware(http.HandlerFunc(routes.GetRecentActivity))).Methods(http.MethodGet)

	// Open endpoints

	noAuthRouter := router.PathPrefix("/").Subrouter()
	withR(noAuthRouter, "/", routes.GetServerStatus).Methods(http.MethodGet)
	withR(noAuthRouter, "/login", routes.RequestLoginPhrase).Methods(http.MethodGet)
	withR(noAuthRouter, "/auth", routes.AuthenticateClientJWT).Methods(http.MethodPost)

	// // Spec
	// specRouter := router.PathPrefix("/spec").Subrouter()
	// withR(specRouter, "/", routes.CreateSpec, middlewares.AuthSessionMiddleware).Methods(http.MethodPost)
	// withR(specRouter, "/{name}", routes.CreateSpec, middlewares.AuthSessionMiddleware).Methods(http.MethodGet)

}

type routeHandler func(http.ResponseWriter, *http.Request)

func withR(r *mux.Router, path string, handler routeHandler, mwf ...mux.MiddlewareFunc) *mux.Route {
	for _, mw := range mwf {
		r.Use(mw)
	}

	return r.HandleFunc(path, handler)
}
