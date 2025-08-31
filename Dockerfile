############# base image
FROM alpine:3.22.0 AS base

############# gardener-extension-provider-dns-cloudflare
FROM base AS gardener-extension-provider-dns-cloudflare

COPY charts /charts
COPY ./dist/gardener-extension-provider-dns-cloudflare_linux_amd64_v1/gardener-extension-provider-dns-cloudflare /gardener-extension-provider-dns-cloudflare
ENTRYPOINT ["/gardener-extension-provider-dns-cloudflare"]
