# Stack Input Authoring Guide

Purpose: author `stack_input.proto` for a resource kind, defining how the CLI passes inputs to IaC modules.

## Location and Package
- Path: `apis/project/planton/provider/<provider>/<kindfolder>/v1/stack_input.proto`
- `syntax = "proto3";`
- `package org.openmcf.provider.<provider>.<kindfolder>.v1;`
- Do NOT include `go_package`.

## Imports
- `project/planton/shared/iac.proto`
- `project/planton/shared/iac/pulumi/pulumi.proto`
- `project/planton/shared/iac/terraform/terraform.proto`
- `project/planton/provider/<provider>/<kindfolder>/v1/api.proto`

## Message
- Define `<Kind>StackInput` with fields (in this exact order/numbering):
  1. `org.openmcf.shared.IacProvisioner provisioner = 1;`
  2. `org.openmcf.shared.iac.pulumi.PulumiStackInfo pulumi = 2;`
  3. `org.openmcf.shared.iac.terraform.TerraformStackInfo terraform = 3;`
  4. `<Kind> target = 1;` (from `api.proto`)
  5. `provider_credential = 2;` (provider-specific credential spec type per import above)

## Notes
- Keep field ordering and numbers stable.
- No validation options/imports here; validations belong in `api.proto`/`spec.proto` rules.
