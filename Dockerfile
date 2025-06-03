############# builder
FROM golang:1.24 AS builder

WORKDIR /go/src/github.com/schrodit/gardener-extension-provider-dns-cloudflare
COPY . .
RUN make install

############# base image
FROM alpine:3.22.0 AS base

############# gardener-extension-provider-dns-cloudflare
FROM base AS gardener-extension-provider-dns-cloudflare

COPY charts /charts
COPY --from=builder /go/bin/gardener-extension-provider-dns-cloudflare /gardener-extension-provider-dns-cloudflare
ENTRYPOINT ["/gardener-extension-provider-dns-cloudflare"]
