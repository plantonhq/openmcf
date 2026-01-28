# Build and Docs Authoring Guide

Purpose: produce `README.md` and `examples.md` for a resource after its protos exist.

## Inputs to read
- `api.proto` — envelope (api_version/kind/metadata/spec/status)
- `spec.proto` — fields and brief descriptions
- `stack_input.proto` — provisioner/IaC/target/credential
- `stack_outputs.proto` — outputs fields
- Credential proto (for provider) — name only for documentation

## README.md (outline)
- Title: `<Kind>`
- Intro: what the resource is and what it provisions
- Spec fields (80/20 summary): key fields with one-line descriptions
- Stack outputs summary
- How it works: Pulumi and Terraform modules in this repo
- References: official docs URLs

## examples.md (outline)
- Scenarios
  - `## Minimal manifest (YAML)` — only required fields
  - `## <Scenario name>` — realistic example(s)
- Include complete YAML manifests in fenced blocks
- CLI flows
  - Validate: `openmcf validate --manifest <file>`
  - Pulumi deploy: `openmcf pulumi update --manifest <file> --stack <org/project/stack> --module-dir <path>`
  - Terraform deploy: `openmcf tofu apply --manifest <file> --auto-approve`

## Notes
- Keep concise and practical (80/20). Avoid secrets in examples.
