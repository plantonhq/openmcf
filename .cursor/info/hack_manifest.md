# Hack Manifest Authoring Guide

Purpose: create a minimal top-level `iac/hack/manifest.yaml` for a resource to test quickly.

## Inputs to read
- `api.proto` — `api_version` and `kind`
- `spec.proto` — key required fields and safe defaults

## Path
- `apis/project/planton/provider/<provider>/<kindfolder>/v1/iac/hack/manifest.yaml`

## Manifest skeleton
```yaml
apiVersion: <provider>.openmcf.org/v1
kind: <Kind>
metadata:
  name: <kindfolder>-demo
spec:
  # Fill essential fields with safe defaults derived from <Kind>Spec
```

## Notes
- Do not include credentials; they are passed via stack input.
- Keep optional sections commented or minimal.
 - For fields using ValueOrRef wrappers (e.g., StringValueOrRef, Int32ValueOrRef), always set the direct value in the manifest (use value, not valueFrom).
