# Forge Workflow Analysis

## Comparison: Forge Output vs Ideal State

### What Forge Currently Creates

Based on rules 001-019:

#### Proto Definitions
- ✅ `spec.proto` (rule 001)
- ✅ `spec.proto` with validations (rule 002)
- ✅ `spec_test.go` (rule 003)
- ✅ `stack_outputs.proto` (rule 004)
- ✅ `api.proto` (rule 005)
- ✅ `stack_input.proto` (rule 006)
- ✅ Generated `.pb.go` stubs (rule 017)

#### Documentation (v1 level)
- ✅ `README.md` (rule 007)
- ✅ `examples.md` (rule 007)

#### IaC - Pulumi
- ✅ Module files: `main.go`, `locals.go`, `outputs.go`, resource files (rule 009)
- ✅ Entrypoint files: `main.go`, `Pulumi.yaml`, `Makefile` (rule 010)
- ✅ E2E testing (rule 011)
- ✅ `README.md` (rule 012)
- ✅ `examples.md` (rule 012)
- ✅ `debug.sh` (rule 012)

#### IaC - Terraform
- ✅ Module files: `variables.tf`, `provider.tf`, `locals.tf`, `main.tf`, `outputs.tf` (rule 013)
- ✅ E2E testing (rule 014)
- ✅ `README.md` (rule 015)
- ✅ `examples.md` (rule 015)

#### Supporting Files
- ✅ `iac/hack/manifest.yaml` (rule 008)

#### Registry
- ✅ Enum entry in `cloud_resource_kind.proto` (rule 016)

#### Validation
- ✅ Build validation (rule 018)
- ✅ Test validation (rule 019)

### What Ideal State Requires But Forge Doesn't Create

#### Critical Gap
1. ❌ **`v1/docs/README.md`** - Comprehensive research document
   - Landscape analysis
   - Deployment methods comparison
   - Best practices
   - 80/20 scoping decisions
   - This is a MAJOR documentation piece (typically 300-1000+ lines)

#### Nice to Have Gaps
2. ❌ **`iac/pulumi/overview.md`** - Module architecture document
   - High-level architecture
   - Design decisions
   - Resource relationships

### What Forge Creates But Ideal State Says Is Optional

Based on updated requirements from user:

1. ⚠️ **`iac/pulumi/examples.md`** - Not needed (examples live in `v1/examples.md`)
2. ⚠️ **`iac/tf/examples.md`** - Not needed (examples live in `v1/examples.md`)

**Rationale:** Same examples work for both Pulumi and Terraform, so no need for IaC-specific examples.

### Recommendations

#### Must Fix
1. **Add Rule 020: Generate `v1/docs/README.md`**
   - This is the research document
   - Major documentation effort
   - Requires LLM to generate comprehensive content
   - Should be informed by proto schema and component purpose

#### Should Consider
2. **Add Rule 021: Generate `iac/pulumi/overview.md`** (optional)
   - Module architecture document
   - Could be generated after rule 009 (pulumi module)
   - Less critical than research doc

#### Should Update
3. **Update Rule 012** - Make `examples.md` generation optional or remove it
4. **Update Rule 015** - Make `examples.md` generation optional or remove it

#### Must Update
5. **Update `forge-project-planton-component.mdc`**
   - Include rules 016-019 in the sequence
   - Add rule 020 (research doc)
   - Optionally add rule 021 (overview.md)
   - Reference ideal state document

### Updated Forge Sequence (Proposed)

```
1. 001-spec-proto.mdc - Generate spec.proto (minimal)
2. 002-spec-validate.mdc - Add validations
3. 003-spec-tests.mdc - Add unit tests
4. 004-stack-outputs.mdc - Generate stack_outputs.proto
5. 005-api.mdc - Generate api.proto
6. 006-stack-input.mdc - Generate stack_input.proto
7. 016-cloud-resource-kind.mdc - Register in enum
8. 017-generate-proto-stubs.mdc - Generate .pb.go files
9. 007-docs.mdc - Generate v1/README.md and v1/examples.md
10. 020-research-docs.mdc - Generate v1/docs/README.md (NEW)
11. 008-hack-manifest.mdc - Generate test manifest
12. 009-pulumi-module.mdc - Generate Pulumi module
13. 010-pulumi-entrypoint.mdc - Generate Pulumi entrypoint
14. 021-pulumi-overview.mdc - Generate iac/pulumi/overview.md (NEW, optional)
15. 011-pulumi-e2e.mdc - Pulumi E2E test
16. 012-pulumi-docs.mdc - Generate iac/pulumi/README.md, debug.sh (no examples.md)
17. 013-terraform-module.mdc - Generate Terraform module
18. 014-terraform-e2e.mdc - Terraform E2E test
19. 015-terraform-docs.mdc - Generate iac/tf/README.md (no examples.md)
20. 018-build-validation.mdc - Validate build
21. 019-test-validation.mdc - Validate tests
```

### Alignment with Ideal State

After implementing the recommended changes:

**Critical Items (40% - Must Have)**
- ✅ All 9 critical checklist items will be covered

**Important Items (40% - Should Have)**
- ✅ All 6 important checklist items will be covered
- ✅ Includes comprehensive research doc (v1/docs/README.md)

**Nice to Have (20% - Polish)**
- ⚠️ Partially covered (overview.md is optional)

**Result:** Forge will create 95-100% complete components matching ideal state!

## Next Steps

1. Create rule 020 for research documentation generation
2. (Optional) Create rule 021 for Pulumi overview.md
3. Update rules 012 and 015 to remove examples.md generation
4. Update main forge orchestrator with correct sequence
5. Write comprehensive README for forge workflow

