# Terraform Docs Authoring Guide

Purpose: author `README.md` and `examples.md` for the Terraform module under `iac/tf/`.

## Target directory
- `apis/project/planton/provider/<provider>/<kindfolder>/v1/iac/tf/`

## Files
- `README.md` — CLI flows using OpenMCF tofu with default local backend
- `examples.md` — runnable examples (manifests and CLI flows)

## Notes
- Do not lint here; focus on content.
- Align examples with the resource’s `iac/hack/manifest.yaml` and CLI flows.
