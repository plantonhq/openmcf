# Forge: Create Components

## Overview

**Forge** is the rule system for bootstrapping **complete, production-ready components** in Planton. It orchestrates 19 atomic rules that create everything from proto definitions to IaC modules to documentation and presets.

**Key principle:** Forge creates components that match **95-100% of the ideal state** defined in `architecture/component.md`, and whose file layout passes the machine-enforced anatomy gate (`pkg/anatomy`, CI lane `lint.component-anatomy.yaml`).

## Component Anatomy

A component lives at `catalog/{provider}/{kind}/`:

- **Component root (the living component):** `README.md` (GitHub-facing page), `catalog.md` (the catalog page), `logo.svg`, optional `GUIDE.md` (authored operational judgment), `iac/pulumi/` + `iac/tf/` (each with a README.md; optional `iac/import-map.yaml`, optional staged `iac/crds/`), `presets/` (`.yaml` + `.md` sidecar pairs), `e2e/` (manifest.yaml, profile, scenarios), optional `conversions/`.
- **Version dir `v1alpha1/` (the versioned contract ONLY):** `api.proto`, `spec.proto`, `input.proto`, `outputs.proto`, their `.pb.go` stubs, `BUILD.bazel`, `spec_test.go`, `reference.md`.

The proto FILE names are `input.proto` / `outputs.proto`; the MESSAGE names are `{Kind}StackInput` / `{Kind}StackOutputs` — the message names are the identity every downstream consumer keys on.

## What Forge Creates

When you run forge, you get a fully-implemented component:

### Proto API Definitions (v1alpha1/)
- ✅ `spec.proto` - Configuration schema with field validations
- ✅ `input.proto` - Inputs to IaC modules (`{Kind}StackInput`: target + provider config)
- ✅ `outputs.proto` - Deployment outputs (`{Kind}StackOutputs`)
- ✅ `api.proto` - KRM wiring (apiVersion, kind, metadata, spec, status)
- ✅ Generated `.pb.go` stubs for all proto files
- ✅ `spec_test.go` - Unit tests for ALL validation rules
- ✅ **Tests execute and pass** - Validates buf.validate rules work correctly

### IaC Modules - Pulumi (iac/pulumi/)
- ✅ Module files: `module/main.go`, `module/locals.go`, `module/outputs.go`, resource-specific files
- ✅ Entrypoint: `main.go` (`package main`) and `Pulumi.yaml` at the `iac/pulumi/` root — nothing else
- ✅ Documentation: `README.md`

### IaC Modules - Terraform (iac/tf/)
- ✅ Module files: `variables.tf` (generated), `provider.tf`, `locals.tf`, `main.tf`, `outputs.tf`
- ✅ Documentation: `README.md`
- ✅ **100% behavioral parity** with the Pulumi module (Parity Mandate)

### Documentation
- ✅ Component-root `README.md` - User-facing overview (the GitHub-facing component page)
- ✅ Component-root `GUIDE.md` - Authored operational judgment: recipes, design rationale, parity accounting
- ✅ Generated `reference.md` - Regenerated via `make generate-reference`, embedding the e2e manifest as its Example

### Supporting Files
- ✅ Component-root `e2e/manifest.yaml` - The complete, protovalidate-valid example manifest; the reference page's Example block AND the E2E framework's testability marker and deployed fixture
- ✅ Component-root `presets/` - 2-3 ready-to-deploy configuration templates with `.md` sidecars
- ✅ Enum entry in `cloud_resource_kind.proto`
- ✅ Build validation passed
- ✅ Test validation passed

Layout conformance is machine-checked: run the anatomy gate (`go test ./pkg/anatomy/...`) instead of auditing file inventories by hand.

## When to Use Forge

Use forge when you need to:
- ✅ **Bootstrap a new component from scratch**
- ✅ Add support for a new cloud provider resource
- ✅ Add support for a new SaaS platform resource
- ✅ Add a new Kubernetes workload or addon

**Don't use forge when:**
- ❌ Component already exists (use **update** instead)
- ❌ You only need to fix/enhance existing component (use **update**)
- ❌ You want to remove a component (use **delete**)
- ❌ You want to check completion status (use **audit**)

## How to Use Forge

### Basic Usage

```
@forge-planton-component <ComponentName> --provider <provider>
```

### Examples

**Create a SaaS platform resource:**
```
@forge-planton-component MongodbAtlas --provider atlas
```

**Create a GCP resource:**
```
@forge-planton-component GcpStorageBucket --provider gcp
```

**Create an AWS resource:**
```
@forge-planton-component AwsSqsQueue --provider aws
```

**Create a Kubernetes workload:**
```
@forge-planton-component PostgresKubernetes --provider kubernetes --category workload
```

**Create a Kubernetes addon:**
```
@forge-planton-component CertManagerKubernetes --provider kubernetes --category addon
```

### Required Information

Before running forge, have ready:
1. **Component Name** - PascalCase (e.g., `GcpCertManagerCert`)
2. **Provider** - One of: aws, gcp, azure, kubernetes, atlas, snowflake, confluent, digitalocean, civo, cloudflare
3. **Category** - Only for Kubernetes: addon, workload, or config

### What Forge Asks You

Forge will interview you to gather:
- Component purpose and use case
- Key configuration fields (for spec.proto)
- Expected outputs (for outputs.proto)
- Provider-specific details (project IDs, regions, etc.)
- Credential requirements
- Best practices and gotchas

## The 19-Rule Workflow

Forge orchestrates 19 rules in 9 phases. (E2E execution is handled separately via
component e2e profiles, not the forge pipeline.)

### Phase 1: Proto API Definitions
1. `001-spec-proto` - Generate spec.proto
2. `002-spec-validate` - Add validations
3. `003-spec-tests` - Generate tests
4. `004-outputs` - Generate outputs.proto
5. `005-api` - Generate api.proto
6. `006-input` - Generate input.proto

### Phase 2: Registration
7. `014-cloud-resource-kind` - Register enum
8. `015-generate-proto-stubs` - Generate .pb.go files

### Phase 3: Documentation
9. `007-readme` - Generate the component-root README.md

### Phase 4: Test Infrastructure
10. `008-e2e-manifest` - Generate the component-root e2e/manifest.yaml

### Phase 5: Pulumi Implementation
11. `009-pulumi-module` - Generate module
12. `010-pulumi-entrypoint` - Generate entrypoint (main.go, Pulumi.yaml)
13. `011-pulumi-readme` - Generate docs

### Phase 6: Terraform Implementation
14. `012-terraform-module` - Generate module
15. `013-terraform-readme` - Generate docs

### Phase 7: Presets
16. `018-presets` - Generate initial presets (2-3 common configuration templates)

### Phase 8: Final Validation
17. `016-build-validation` - Compile all Go code (recursive component build + release-equivalent entrypoint build)
18. `017-test-validation` - Run all tests

### Phase 9: Guide + Reference Regeneration
19. `019-guide` - Seed the component-root `GUIDE.md` (authored wisdom), then run `make generate-reference` so the reference page, guide head link, and catalog indexes materialize together

## Progress Tracking

Forge provides real-time progress updates:

```
🔨 Forge: Creating MongodbAtlas

Phase 1: Proto API Definitions
[1/19] ✅ Generated spec.proto
[2/19] ✅ Added buf.validate rules
[3/19] ✅ Generated and ran spec tests
[4/19] ✅ Generated outputs.proto
[5/19] ✅ Generated api.proto
[6/19] ✅ Generated input.proto

Phase 2: Registration
[7/19] ✅ Registered MongodbAtlas = 51 in cloud_resource_kind.proto
[8/19] ✅ Generated proto stubs (make protos)

Phase 3: Documentation
[9/19] ✅ Generated component-root README.md

Phase 4: Test Infrastructure
[10/19] ✅ Generated e2e/manifest.yaml

Phase 5: Pulumi Implementation
[11/19] ✅ Generated Pulumi module
[12/19] ✅ Generated Pulumi entrypoint
[13/19] ✅ Generated Pulumi docs

Phase 6: Terraform Implementation
[14/19] ✅ Generated Terraform module
[15/19] ✅ Generated Terraform docs

Phase 7: Presets
[16/19] ✅ Generated initial presets

Phase 8: Final Validation
[17/19] ✅ Build validation passed (go build ./catalog/<provider>/<kind>/...)
[18/19] ✅ Component tests passed (go test -v ./catalog/<provider>/<kind>/v1alpha1/)

Phase 9: Guide + Reference
[19/19] ✅ Seeded GUIDE.md and regenerated the reference

🎉 Component creation complete!

📍 Location: catalog/atlas/mongodbatlas/
📊 Expected Audit Score: 95-100%

Next steps:
1. Review generated files
2. Run: @audit-planton-component MongodbAtlas
3. Make any custom modifications
4. Commit and push
```

## Error Handling

### Automatic Retries
- Each rule retries up to 3 times on fixable errors
- Build errors are fixed automatically when possible
- Test failures trigger fixes and retries

### Manual Intervention
If a rule fails after 3 attempts:
1. Forge stops and shows the error
2. Fix the issue manually
3. Resume from the failed rule:
   ```
   @forge-planton-component MongodbAtlas --resume-from 012
   ```

### Common Issues

**Issue: Proto build fails**
- **Cause:** Invalid protobuf syntax
- **Fix:** Forge auto-fixes and retries
- **If persists:** Check .proto file manually

**Issue: Pulumi/Terraform E2E fails**
- **Cause:** Missing credentials or invalid config
- **Fix:** Check manifest values, update and retry

**Issue: Tests fail**
- **Cause:** Validation rules too strict or test logic error
- **Fix:** Forge analyzes and fixes tests automatically

## Post-Forge Validation

After forge completes, validate the component:

**Option 1: Manual Audit**
```bash
@audit-planton-component <ComponentName>
```
**Expected Result:** 95-100% completion score

If score is lower, the audit report shows what's missing, why it matters, and how to fix it.

**Option 2: Auto-Complete (Recommended)**
```bash
@complete-planton-component <ComponentName>
```
Automatically audits and fills any remaining gaps to reach 95%+. Useful if forge had partial failures.

**Always:** the anatomy gate (`go test ./pkg/anatomy/...`) must be green — it machine-checks the component's file layout against the canonical anatomy.

## Customization After Forge

Forge creates a **production-ready baseline**. Common customizations:

### Add More Fields to Proto
1. Edit `spec.proto` to add fields
2. Update validations in `spec.proto`
3. Update tests in `spec_test.go`
4. Run `make protos` to regenerate stubs
5. Update Pulumi module to use new fields
6. Regenerate Terraform `variables.tf` and update the module to match
7. Update the e2e manifest and presets if the new fields belong in examples
8. Run `go build ./catalog/<provider>/<kind>/... && go test -v ./catalog/<provider>/<kind>/v1alpha1/`

### Modify IaC Implementation
1. Update Pulumi module files (`iac/pulumi/module/*.go`)
2. Update Terraform module files (`iac/tf/*.tf`)
3. Keep both engines at behavioral parity (fix whichever engine is wrong)
4. Update documentation if behavior changes

### Enhance Documentation
1. Expand the component-root `README.md` and `catalog.md`
2. Capture new operational judgment in `GUIDE.md`
3. Add troubleshooting to `iac/pulumi/README.md` or `iac/tf/README.md`

## Comparison to Manual Creation

| Aspect | Manual Creation | Forge |
|--------|----------------|-------|
| Time | 8-16 hours | 15-30 minutes |
| Completeness | 60-80% typical | 95-100% |
| Documentation | Often skipped | Comprehensive |
| Validation | Manual | Automated |
| Consistency | Varies | Standardized |
| Best Practices | Hit or miss | Built-in |
| Error-Prone | Yes | Auto-fixed |

## Reference Documents

- **Ideal State Definition:** `architecture/component.md`
- **Anatomy Conformance Gate:** `pkg/anatomy` (CI lane `.github/workflows/lint.component-anatomy.yaml`)
- **Individual Flow Rules:** `_rules/component/forge/flow/`
- **Main Orchestrator:** `_rules/component/forge/forge-planton-component.mdc`

## Tips and Best Practices

### Before Running Forge

1. **Research the resource** - Understand what you're creating
2. **Check if it exists** - Run `@audit-planton-component` first
3. **Plan your API** - Read the provider schema at the pinned version fully and model it fully (100% of the configurable surface)
4. **Gather examples** - Have reference configurations ready

### During Forge

1. **Be specific** - Provide detailed answers to interview questions
2. **Think production** - Consider real-world use cases
3. **Include validation** - Think about what constraints make sense
4. **Document gotchas** - Share known issues and workarounds

### After Forge

1. **Review everything** - Don't blindly trust generated code
2. **Test locally** - Deploy with the e2e manifest
3. **Enhance docs** - Add your learnings to the surfaces their readers use (spec comments, GUIDE.md, presets)
4. **Run audit** - Verify 100% ideal state compliance

## Troubleshooting

### "Component already exists"
**Error:** `Component MongodbAtlas already exists at ...`

**Solution:** Use `@update-planton-component` instead, or delete first with `@delete-planton-component`.

### "Provider not recognized"
**Error:** `Provider 'xyz' is not valid`

**Valid providers:** aws, gcp, azure, kubernetes, atlas, snowflake, confluent, digitalocean, civo, cloudflare

### "Build failed after 3 attempts"
**Check:**
1. Proto syntax in generated files
2. Go code compiles: `go build ./catalog/<provider>/<kind>/...` (from the repo root)
3. Import paths are correct
4. Manual fix may be needed

## Next Steps

After reading this README:
1. Review the ideal state document: `architecture/component.md`
2. Try forge on a test component
3. Inspect the generated code
4. Run audit to verify completion
5. Use forge for real components!

---

**Questions?** Check the troubleshooting section or run `@audit-planton-component` to see examples of complete components.

**Ready to create?** Run `@forge-planton-component <YourComponentName> --provider <provider>`
