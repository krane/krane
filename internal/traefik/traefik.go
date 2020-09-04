package traefik

var (
	// Rules are a set of matchers configured with values, that determine if a particular
	// request matches specific criteria. If the rule is verified, the router becomes active, calls
	// middlewares, and then forwards the request to the service.
	RuleRoutingConfigName  = "traefik.http.routers.%s.rule" // "traefik.http.routers.any_name.rule"
	RuleRoutingConfigValue = "Host(`%s`)"                   // "Host(`example.com`)"

	// Registers a port. Useful when the container exposes multiples ports.
	ServiceRouterConfig = "traefik.http.services.%s.loadbalancer.server.port" // traefik.http.services.myservice.loadbalancer.server.port=80
)
