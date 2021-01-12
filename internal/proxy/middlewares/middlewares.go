package middlewares

import (
	"fmt"
	"strconv"
)

func RedirectToHTTPSLabels(deployment string) map[string]string {
	labels := make(map[string]string, 0)

	labels[fmt.Sprintf("traefik.http.routers.%s-insecure.middlewares", deployment)] = "redirect-to-https@docker"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.scheme"] = "https"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.port"] = "443"
	labels["traefik.http.middlewares.redirect-to-https.redirectscheme.permanent"] = "true"

	return labels
}

func RateLimitLabels(deployment string, rateLimit uint) map[string]string {
	labels := make(map[string]string, 0)
	labels[fmt.Sprintf("traefik.http.middlewares.%s-ratelimit.ratelimit.average", deployment)] = strconv.FormatUint(uint64(rateLimit), 10)
	return labels
}
