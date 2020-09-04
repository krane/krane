# This is a custom image of traefik with managed configurations for krane
FROM traefik:v2.2

# Ref: https://docs.traefik.io/reference/static-configuration/env/
ENV TRAEFIK_PROVIDERS_DOCKER_NETWORK=krane
ENV TRAEFIK_API_INSECURE=true
ENV TRAEFIK_PROVIDERS_DOCKER=true
ENV TRAEFIK_PROVIDERS_DOCKER_EXPOSEDBYDEFAULT=true
ENV TRAEFIK_API_DASHBOARD=true
ENV TRAEFIK_LOG_LEVEL=DEBUG

EXPOSE 80
EXPOSE 8080

VOLUME /var/run/docker.sock