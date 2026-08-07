# Experimental Gateway API CRDs

This preset installs the Kubernetes Gateway API CRDs in the experimental channel. It contains everything the standard channel has, plus resources and fields that are still under active development.

## When to Use

- You need an experimental resource that has not graduated to the standard channel (as of v1.6: XBackend, XBackendTrafficPolicy, XMesh)
- You need an experimental field on a standard resource that only ships in the experimental-channel CRDs
- You accept that experimental APIs may change or be removed between releases

## Key Configuration Choices

- **Experimental channel** -- all standard resources (with additional experimental fields) plus experimental resources. Stability caveat: future Gateway API releases may break or remove experimental resources without a deprecation cycle, and upgrading can require migrating or deleting resources built on them. Prefer the standard channel (`01-standard`) unless a specific experimental feature is required
- **Version** (`v1.6.1`) -- the spec default; check [Gateway API releases](https://github.com/kubernetes-sigs/gateway-api/releases) for updates
- **No namespace** -- CRDs are cluster-scoped; no namespace is needed

## Placeholders to Replace

No placeholders -- this preset is directly deployable.
