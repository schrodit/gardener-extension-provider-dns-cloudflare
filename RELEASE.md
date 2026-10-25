# Release Flow

1. Merge the release pull request created by Release Please. It updates the [release manifest](./.release-please-manifest.json) and changelog.
2. The [Release workflow](https://github.com/schrodit/gardener-extension-provider-dns-cloudflare/actions/workflows/release.yaml) creates the tag and publishes the binaries, container image, and Helm chart.
