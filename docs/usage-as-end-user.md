# Using the Cloudflare dns provider extension with Gardener as end-user

The [`core.gardener.cloud/v1beta1.Shoot` resource](https://github.com/gardener/gardener/blob/master/example/90-shoot.yaml) declares a few fields that are meant to contain provider-specific configuration.

This document describes the configurable options for GCP and provides an example `Shoot` manifest with minimal configuration that can be used to create a GCP cluster (modulo the landscape-specific information like cloud profile names, secret binding names, etc.).

## Cloudflare Provider Credentials

In order for Gardener to create dns records using CloudFlare, a Shoot has to provide credentials with sufficient permissions to the desired cloudflare zones.

Every shoot can either reference these credentials 
- in the shoot manifest using a custom domain (see [Example shoot manifest](#example-shoot-manifest))
- or use the default domain (see [Usage as Operator for details](./usage-as-operator.md))

Create a api token for your cloudflare account as described in https://developers.cloudflare.com/api/tokens/create/ 

This `Secret` must look as follows:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: domain-cloudflare
  namespace: garden-dev
type: Opaque
data:
  apiToken: base64(api token)
```

## Example `Shoot` manifest

Please find below an example `Shoot` manifest:

```yaml
apiVersion: core.gardener.cloud/v1alpha1
kind: Shoot
metadata:
  name: johndoe-cloudflare
  namespace: garden-dev
spec:

  dns:
    domain: my-domain.example.com
    providers:
    - type: cloudflare-dns
      secretName: domain-cloudflare

  cloudProfileName: xxx
  region: xxx
  secretBindingName: xxx
  provider:
    type: xxx
    infrastructureConfig:
      ...
    controlPlaneConfig:
      ...
    workers:
    ...
  networking:
    ...
  kubernetes:
    version: 1.16.1
```
