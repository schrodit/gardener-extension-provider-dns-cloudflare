# Release Flow

1. Increase the version in [./VERSION](./VERSION)
    - Use semver versioning
2. Regenerate all resources by running `make generate`
3. Commit the the updated file to master with the title "Release <new version>"
4. Trigger release by running [Release GitHub Action Job](https://github.com/schrodit/gardener-extension-provider-dns-cloudflare/actions/workflows/release.yaml)
