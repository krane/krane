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
	router := mux.NewRouter()

	withBaseMiddlewares(router)
	withRoutes(router)

	srv := &http.Server{
		Handler:      router,
		Addr:         os.Getenv("LISTEN_ADDRESS"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Krane server on %s", os.Getenv("LISTEN_ADDRESS"))
	logrus.Fatal(srv.ListenAndServe())
}

func withBaseMiddlewares(router *mux.Router) {
	router.Use(middlewares.Logging)
	router.Use(handlers.RecoveryHandler())
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}),
		handlers.AllowedOrigins([]string{"localhost", "*"})))
}

func withRoutes(router *mux.Router) {
	router.HandleFunc("/", routes.GetServerStatus).Methods(http.MethodGet)
	router.HandleFunc("/login", routes.RequestLoginPhrase).Methods(http.MethodGet)
	router.HandleFunc("/auth", routes.AuthenticateClientJWT).Methods(http.MethodPost)

	createSpec := http.HandlerFunc(routes.CreateSpec)
	router.Handle("/spec", middlewares.AuthSessionMiddleware(createSpec)).Methods(http.MethodPost)

	updateSpec := http.HandlerFunc(routes.UpdateSpec)
	router.Handle("/spec/{name}", middlewares.AuthSessionMiddleware(updateSpec)).Methods(http.MethodPut)

	router.HandleFunc("/deployments", routes.GetDeployments).Methods(http.MethodGet)
	router.HandleFunc("/deployments/{name}", routes.GetDeployment).Methods(http.MethodGet)
	router.HandleFunc("/deployments/{name}", routes.RunDeployment).Methods(http.MethodPost)
	router.HandleFunc("/deployments/{name}", routes.DeleteDeployment).Methods(http.MethodDelete)
	router.HandleFunc("/deployments/{name}/stop", routes.StopDeployment).Methods(http.MethodPost)
	router.HandleFunc("/alias/{name}", routes.UpdateDeploymentAlias).Methods(http.MethodPost)
	router.HandleFunc("/alias/{name}", routes.DeleteDeploymentAlias).Methods(http.MethodDelete)
	router.HandleFunc("/activity", routes.GetRecentActivity).Methods(http.MethodGet)
	return
}
