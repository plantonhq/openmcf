# Container Image App

This preset deploys a DigitalOcean App Platform application from a public Docker Hub image, with professional-tier compute and CPU autoscaling. It does not build from source.

The sample uses `library/nginx` so the preset applies without placeholders. Change `registry`, `repository`, and `tag` to your image. For DigitalOcean Container Registry set `registryType: docr` and leave `registry` empty (DOCR does not take a hostname here). For GitHub Container Registry set `registryType: ghcr` and `registry` to `ghcr.io`.

Do not set `instanceCount` while `autoscaling` is on -- App Platform ignores a fixed count in that case, and the spec rejects the combination.

## When to Use

- You already build images in CI and want App Platform to run them
- Production HTTP services that need more than one instance
- Public images on Docker Hub, or private images on DOCR / GHCR

## Key Configuration Choices

- **Registry type** (`registryType: docker_hub`) -- `docker_hub`, `docr`, or `ghcr`. Not the string `"docker-hub"` as a DOCR hostname.
- **Professional XS** (`instanceSizeSlug: professional-xs`) -- dedicated CPU. Slugs are free-form; check current App Platform sizes in DigitalOcean's docs.
- **Autoscaling** (`minInstanceCount: 2`, `maxInstanceCount: 5`, `cpuPercent: 80`) -- App Platform scales on average CPU. Leave `instanceCount` unset.

## Related Presets

- **01-git-source-web** -- use when App Platform should build from a Git repository instead of running a pre-built image
