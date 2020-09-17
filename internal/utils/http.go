package utils

import "net/http"

func QueryParamOrDefault(r *http.Request, param, fallback string) string {
	query := r.URL.Query()
	value := query.Get(param)
	if value == "" {
		return fallback
	}
	return value
}
