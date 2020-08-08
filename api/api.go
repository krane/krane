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

	withMiddlewares(router)
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

func withMiddlewares(router *mux.Router) {
	router.Use(middlewares.Logging)
	router.Use(handlers.RecoveryHandler())
	router.Use(handlers.CORS(
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodDelete}),
		handlers.AllowedOrigins([]string{"localhost", "*"})))
	return
}

func withRoutes(router *mux.Router) {
	router.HandleFunc("/", routes.IndexRoute).Methods(http.MethodGet)
	router.HandleFunc("/spec", routes.CreateSpec).Methods(http.MethodPost)                                  // Create a Spec
	router.HandleFunc("/spec/{name}", routes.UpdateSpec).Methods(http.MethodPut)                            // Update Spec
	router.HandleFunc("/deployments", routes.GetDeployments).Methods(http.MethodGet)                        // Get all deployments
	router.HandleFunc("/deployments/{name}", routes.RunDeployment).Methods(http.MethodPost)                 // Run a deployment
	router.HandleFunc("/deployments/{name}", routes.DeleteDeployment).Methods(http.MethodDelete)            // Delete a deployment
	router.HandleFunc("/deployments/{name}", routes.GetDeployment).Methods(http.MethodGet)                  // Get one deployment
	router.HandleFunc("/deployments/{name}/stop", routes.StopDeployment).Methods(http.MethodPost)           // Stop a deployment
	router.HandleFunc("/deployments/alias/{name}", routes.UpdateDeploymentAlias).Methods(http.MethodPost)   // Update deployment alias
	router.HandleFunc("/deployments/alias/{name}", routes.DeleteDeploymentAlias).Methods(http.MethodDelete) // Delete deployment alias
	router.HandleFunc("/jobs", routes.GetRunningJobs).Methods(http.MethodGet)                               // Doesnt currently work because we are not posting Jobs
	router.HandleFunc("/activity", routes.GetRecentActivity).Methods(http.MethodGet)                        // Get recent activity within a date range
	return
}
