package middlewares

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

// Logging : custom middleware for logging http requests
func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// logic before reading request values, putting things into the request context, performing authentication

		// important that we call the 'next' handler in the chain. If we don't, then request handling will stop here.
		next.ServeHTTP(w, r)

		// logic after; useful for logger, metrics, etc.
		// NOTE: It's important that we don't use the ResponseWriter after we've called the
		// next handler: we may cause conflicts when trying to write the response
	}

	return handlers.CombinedLoggingHandler(os.Stdout, http.HandlerFunc(fn))
}
