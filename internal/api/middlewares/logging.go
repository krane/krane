package middlewares

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Logging : custom middleware for logging HTTP requests wrapping the gorrila combined logging
// handler which provides Apache Combined Log Format commonly used by both Apache and nginx.
func Logging(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Logic before - reading request values, putting things into the
			// request context, performing authentication

			// Important that we call the 'next' handler in the chain. If we don't,
			// then request handling will stop here.
			next.ServeHTTP(w, r)

			// Logic after - useful for logging, metrics, etc.

			// NOTE: It's important that we don't use the ResponseWriter after we've called the
			// next handler: we may cause conflicts when trying to write the response
		},
	)))
}
