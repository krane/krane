package utils

import "net/http"

// QueryParamOrDefault : get a query param from http request, returns a default value if query param not set
func QueryParamOrDefault(r *http.Request, param, fallback string) string {
	query := r.URL.Query()
	value := query.Get(param)
	if value == "" {
		return fallback
	}
	return value
}
