package traefik

var (
	// Rules are a set of matchers configured with values, that determine if a particular
	// request matches specific criteria. If the rule is verified, the router becomes active, calls
	// middlewares, and then forwards the request to the service.
	RuleRoutingConfigName  = "traefik.status.routers.%s.rule" // "traefik.status.routers.any_name.rule"
	RuleRoutingConfigValue = "Host(`%s`)"                     // "Host(`example.com`)"

	// Registers a port. Useful when the container exposes multiples ports.
	ServiceRouterConfig = "traefik.status.services.%s.loadbalancer.server.port" // traefik.status.services.myservice.loadbalancer.server.port=80

	// TODO: Complete adding all possible labels available
)
