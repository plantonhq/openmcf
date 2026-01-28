# Pulumi E2E Run Guide

Purpose: run an end-to-end Pulumi update using the local CLI and fix issues as needed.

## Prereqs
- Protos generated and Pulumi module + entrypoint created.
- A minimal manifest at `iac/hack/manifest.yaml`.

## Steps
1. Install local CLI: `make local`
2. pulumi login (local backend): `pulumi login file://$HOME/.pulumi`
3. pulumi stack select (create if needed): `pulumi stack select <org/project/stack> --create`
4. Run update via CLI:
   - `openmcf pulumi update --manifest <path> --stack <org/project/stack> --module-dir <module_dir>`

## Notes
- Organization for local backend is the `organization` part of the stack name.
- Iterate on module code if errors occur; re-run update until successful.
