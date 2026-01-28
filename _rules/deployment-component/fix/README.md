# Fix: Targeted Fixes with Cascading Updates

## Overview

**Fix** is the operation for making targeted improvements to deployment components with automatic propagation to all related artifacts. It ensures that when you fix the source code (proto, IaC), all documentation, examples, and tests are automatically updated to match.

**Core Philosophy:** Source code is the ultimate source of truth. Documentation describes code, not the other way around.

## Why Fix Exists

Components need targeted fixes:
- Bugs in validation logic
- Incorrect IaC implementation
- Missing proto fields
- Documentation drift (examples out of sync)
- Test failures
- Feature parity issues (Pulumi ≠ Terraform)

**Fix makes targeted improvements while maintaining consistency across all artifacts.**

## The Source Code Truth Principle

### The Right Order

```
1. Fix Source Code (Proto, Pulumi, Terraform)
        ↓
2. Update Documentation (README, docs/README.md)
        ↓
3. Update Examples (examples.md)
        ↓
4. Update Tests (spec_test.go)
        ↓
5. Validate Everything Matches
```

### The Wrong Order

```
1. Update documentation to describe desired behavior
        ↓
2. Try to make code match documentation
        ↓
3. Documentation and code drift apart
        ↓
4. Examples stop working
        ↓
5. Chaos ensues
```

**Fix enforces the right order automatically.**

## When to Use Fix

### ✅ Use Fix When

- **Specific bug to fix** - Clear issue with known solution
- **Documentation out of sync** - Examples/docs don't match code
- **Validation logic wrong** - buf.validate rules incorrect
- **IaC implementation incorrect** - Module doesn't deploy properly
- **Missing proto field** - Need to add essential field
- **Test failures** - Tests don't match reality
- **Feature parity broken** - Pulumi and Terraform differ

### ❌ Don't Use Fix When

- Component doesn't exist → Use **forge**
- Want to fill missing files → Use **complete** or **update --fill-gaps**
- General improvements → Use **update** with scenario
- Just checking status → Use **audit**
- Want to remove → Use **delete**

### Fix vs Update

| Aspect | Fix | Update |
|--------|-----|--------|
| **Purpose** | Targeted fix with propagation | General improvements |
| **Input** | Specific explanation of fix | Scenario or explanation |
| **Scope** | Narrow (specific issue) | Broad (fill gaps, refresh all) |
| **Propagation** | Automatic to all artifacts | Depends on scenario |
| **Consistency** | Enforces actively | Trusts existing state |
| **Best For** | Specific bugs/issues | Systematic improvements |

**Fix = Surgical, Update = Systematic**

## How Fix Works

### The Five Consistency Checks

Fix automatically verifies and enforces consistency:

#### 1. Proto ↔ Terraform Variables

**Check:** Every field in spec.proto has matching variable in variables.tf

```
spec.proto:
  int64 disk_size_gb = 5;

variables.tf:
  variable "disk_size_gb" {
    type = number
  }
```

**If mismatch:** Updates variables.tf to match proto

#### 2. Proto ↔ Examples

**Check:** All examples use current field names and meet validation rules

```
examples.md:
  diskSizeGb: 100  # ✓ Matches spec.proto field name
  
Not:
  disk_size: 100   # ❌ Wrong field name
```

**If mismatch:** Updates examples to use correct fields

#### 3. Pulumi ↔ Terraform

**Check:** Both modules create same resources with same behavior

```
Pulumi creates:
  - RDS Instance with backup
  - Security group
  - Subnet group

Terraform creates:
  - RDS Instance with backup  ✓
  - Security group  ✓
  - Subnet group  ✓
```

**If mismatch:** Updates module to restore parity

#### 4. Validations ↔ Tests

**Check:** Every buf.validate rule has test in spec_test.go

```
spec.proto:
  string region = 3 [(buf.validate.field).required = true];

spec_test.go:
  TestMissingRegion(t *testing.T) { ... }  ✓
```

**If mismatch:** Adds missing tests

#### 5. Documentation ↔ Implementation

**Check:** Documentation describes actual behavior, not wishful thinking

```
README.md:
  "Supports PostgreSQL 11-15"

spec.proto:
  enum Version {
    V11 = 1; V12 = 2; V13 = 3; V14 = 4; V15 = 5;
  }
  ✓ Matches
```

**If mismatch:** Updates docs to match code

## Usage

### Basic Syntax

```bash
@fix-openmcf-component <ComponentName> --explain "<detailed fix description>"
```

### The --explain Flag

**Critical:** You must explain what needs fixing. Fix uses this to:
1. Understand the problem
2. Determine which files to change
3. Know what the correct behavior should be
4. Decide what documentation needs updating

**Good Explanations:**
```bash
--explain "primaryDomainName validation rejects *.example.com wildcards, should accept them"

--explain "Pulumi module hardcodes backup_retention_period to 7 days instead of using spec field"

--explain "examples.md uses deprecated 'database_name' field, should use 'db_identifier' from current spec"

--explain "spec_test.go expects required validation on 'region' field but spec.proto doesn't have that validation"
```

**Bad Explanations:**
```bash
--explain "fix it"  # Too vague

--explain "make it work"  # Not specific

--explain "update docs"  # What specifically?
```

## Execution Workflow

### Phase 1: Analysis (30 seconds)

```
📋 Analyzing Fix Request

Component: GcpCertManagerCert
Fix: "Allow wildcard domains"

Reading source code:
  ✓ spec.proto (current validation pattern)
  ✓ Pulumi module (field usage)
  ✓ Terraform module (field usage)
  ✓ spec_test.go (current tests)
  ✓ examples.md (current examples)
  ✓ v1/docs/README.md (design context)

Current Behavior:
  - Pattern: ^[a-z0-9-]+\.[a-z]{2,}$
  - Rejects: *.example.com
  - Tests: No wildcard tests
  - Examples: No wildcard examples

Fix Plan:
  1. Update spec.proto validation pattern
  2. Add wildcard tests
  3. Add wildcard examples
  4. Update documentation
```

### Phase 2: Source Code Fix (2-5 minutes)

```
🔧 Fixing Source Code

[1/5] ✅ spec.proto
  - Updated pattern: ^(\*\.)?[a-z0-9-]+\.[a-z]{2,}$
  - Now accepts wildcards

[2/5] ✅ Proto stubs
  - Regenerated: make protos
  - Status: Success

[3/5] ✅ spec_test.go
  - Added: TestWildcardDomainValidation
  - Added: TestNonWildcardDomainValidation (ensure still works)
  - Result: 18/18 tests pass

[4/5] ✅ Pulumi module
  - Verified: Uses spec field directly (no changes needed)

[5/5] ✅ Terraform module
  - Verified: Uses variable directly (no changes needed)
  - variables.tf already compatible

Source code fix complete!
Component tests: ✅ PASS (18/18)
```

### Phase 3: Documentation Propagation (3-5 minutes)

```
📝 Updating Documentation to Match Source

[1/5] ✅ examples.md
  - Added: Wildcard Domain Example
  - Added: Multi-Domain Wildcard Example
  - Validated: All examples pass schema validation

[2/5] ✅ README.md (user-facing)
  - Updated: Multi-Domain Support section
  - Added: Wildcard support to Key Features
  - Updated: Example usage shows wildcard

[3/5] ✅ docs/README.md (research)
  - Updated: Comparison table (wildcard support)
  - Updated: Best Practices (wildcard patterns)
  - Updated: OpenMCF Approach section

[4/5] ✅ iac/pulumi/README.md
  - Added: Wildcard example to usage section

[5/5] ✅ iac/tf/README.md
  - Added: Wildcard example to usage section

Documentation sync complete!
All docs now match source code!
```

### Phase 4: Consistency Validation (1-2 minutes)

```
🔍 Validating Consistency

Check 1: Proto ↔ Terraform
  ✅ All 17 spec.proto fields in variables.tf
  ✅ Types match (string→string, int64→number, etc.)
  ✅ Required fields match

Check 2: Proto ↔ Examples
  ✅ All 7 examples validate against schema
  ✅ No deprecated field names
  ✅ All values meet validation rules

Check 3: Pulumi ↔ Terraform
  ✅ Same resources created
  ✅ Same spec fields used
  ✅ Same outputs returned
  ✅ Feature parity maintained

Check 4: Validations ↔ Tests
  ✅ All 12 validation rules tested
  ✅ 100% coverage of validation logic

Check 5: Docs ↔ Implementation
  ✅ README describes actual behavior
  ✅ Examples demonstrate actual API
  ✅ Research docs explain actual approach

All consistency checks passed!
```

### Phase 5: Validation (2-4 minutes)

```
🧪 Running Validation

Component Tests:
  go test ./apis/org/openmcf/provider/gcp/gcpcertmanagercert/v1/
  ✅ PASS (18/18 tests, +2 new tests)

Build Validation:
  make build
  ✅ Success (all packages compiled)

Full Test Suite:
  make test
  ✅ PASS (156/156 tests)

All validation passed!
```

### Phase 6: Summary Report

```
✅ Fix Complete: GcpCertManagerCert

Fix Applied:
  "Allow wildcard domains in primaryDomainName validation"

Changes Summary:

SOURCE CODE:
  📝 spec.proto (1 validation pattern updated)
  📝 spec_test.go (2 tests added)
  ℹ️  Pulumi module (verified, no changes needed)
  ℹ️  Terraform module (verified, no changes needed)

DOCUMENTATION:
  📝 examples.md (2 new examples added)
  📝 README.md (features updated, example updated)
  📝 docs/README.md (comparison updated, best practices updated)
  📝 iac/pulumi/README.md (example added)
  📝 iac/tf/README.md (example added)

CONSISTENCY:
  ✅ Proto ↔ Terraform: 17/17 fields match
  ✅ Proto ↔ Examples: 7/7 examples validate
  ✅ Pulumi ↔ Terraform: Feature parity maintained
  ✅ Validations ↔ Tests: 12/12 rules tested
  ✅ Docs ↔ Implementation: Fully synchronized

VALIDATION:
  ✅ Component tests: 18/18 passed (+2 new)
  ✅ Build: Success
  ✅ Full test suite: 156/156 passed

Files Modified: 7
Lines Changed: +120, -2
Duration: 8 minutes

Next Steps:
  1. Review changes (git diff)
  2. Test manually (deploy with hack manifest)
  3. Commit:
     git add -A
     git commit -m "fix(gcp-cert): allow wildcard domains in validation"
```

## Common Fix Scenarios

### Scenario: Validation Too Strict

**Problem:** Validation rejects valid values

**Fix:**
```bash
@fix-openmcf-component AwsVpc --explain "CIDR validation rejects 10.0.0.0/8 which is valid private range"
```

**Actions:**
- Update spec.proto pattern to allow 10.x.x.x
- Add test for 10.0.0.0/8
- Add example using 10.0.0.0/8
- Update docs mentioning valid private ranges

### Scenario: Missing Required Field

**Problem:** Essential field not in proto

**Fix:**
```bash
@fix-openmcf-component MongodbAtlas --explain "spec.proto missing 'region' field which is essential for cluster deployment"
```

**Actions:**
- Add region field to spec.proto
- Add validation for region
- Update Terraform variables.tf (add region variable)
- Update Pulumi module (use spec.Region)
- Update Terraform module (use var.region)
- Add tests for region validation
- Update examples (show region usage)
- Update README (document region field)

### Scenario: IaC Hardcoded Value

**Problem:** Module doesn't use spec field

**Fix:**
```bash
@fix-openmcf-component AwsRdsInstance --explain "Pulumi hardcodes backup_retention_period=7, should use spec.backupRetentionDays"
```

**Actions:**
- Update Pulumi module (use spec field)
- Verify Terraform already uses spec field
- Add test validating different retention values
- Add example showing custom retention
- Update overview (document backup behavior)

### Scenario: Documentation Drift

**Problem:** Docs don't match reality

**Fix:**
```bash
@fix-openmcf-component PostgresKubernetes --explain "examples use 'database_name' but spec.proto uses 'db_identifier'"
```

**Actions:**
- **Read spec.proto** (confirm db_identifier is correct)
- Update examples.md (database_name → db_identifier everywhere)
- Update README (fix any references)
- Validate examples (ensure they work)
- **No code changes** (code is already correct)

### Scenario: Test Failure

**Problem:** Tests failing or incorrect

**Fix:**
```bash
@fix-openmcf-component GcpCloudSql --explain "spec_test.go fails because it expects required validation on project_id but spec.proto doesn't have that rule"
```

**Actions:**
- Read spec.proto (check if validation exists)
- **Decision:** Should project_id be required?
- If YES: Add validation to spec.proto, keep test
- If NO: Fix test to not expect error
- Regenerate stubs if proto changed
- Run tests to verify
- Update examples if validation behavior changed

### Scenario: Feature Parity Broken

**Problem:** Pulumi and Terraform behave differently

**Fix:**
```bash
@fix-openmcf-component GcpGkeCluster --explain "Terraform doesn't create node pool autoscaling but Pulumi does"
```

**Actions:**
- Read Pulumi module (understand correct behavior)
- Update Terraform module (add autoscaling)
- Ensure both use same spec fields
- Add tests for autoscaling
- Update examples (show autoscaling config)
- Run both E2E tests
- Verify feature parity restored

## Consistency Enforcement

Fix actively enforces consistency with automated checks:

### Check 1: Proto ↔ Terraform Variables

**What:** After proto changes, ensures variables.tf matches

**How:**
- Parses spec.proto for all fields
- Parses variables.tf for all variables
- Compares field names, types, validations
- Updates variables.tf if mismatch

**Example:**
```
spec.proto added:
  int64 max_connections = 10;

Fix automatically adds to variables.tf:
  variable "max_connections" {
    type = number
    description = "Maximum database connections"
  }
```

### Check 2: Examples Must Validate

**What:** After any change, ensures examples work

**How:**
- Extracts YAML from examples.md
- Validates each against current spec.proto
- Updates examples if validation fails
- Adds examples if new fields added

**Example:**
```
After adding 'region' field to spec.proto:
  ✓ Updates all existing examples to include region
  ✓ Adds new example showing different regions
  ✓ Validates all examples pass schema validation
```

### Check 3: Feature Parity

**What:** After IaC changes, ensures Pulumi = Terraform

**How:**
- Lists resources created by Pulumi
- Lists resources created by Terraform
- Compares resource types and configurations
- Updates lagging module to match

**Example:**
```
Pulumi creates monitoring dashboard:
  ✓ Detects Terraform doesn't
  ✓ Adds monitoring dashboard to Terraform
  ✓ Verifies same configuration
```

### Check 4: Tests Cover Validations

**What:** After validation changes, ensures tests exist

**How:**
- Lists all buf.validate rules in spec.proto
- Lists all validation tests in spec_test.go
- Identifies untested rules
- Adds missing tests

**Example:**
```
Added validation:
  [(buf.validate.field).string.min_len = 3]

Fix automatically adds test:
  TestFieldMinLength(t *testing.T) {
    // Test with 2 chars (should fail)
    // Test with 3 chars (should pass)
  }
```

### Check 5: Docs Describe Reality

**What:** After code changes, ensures docs match

**How:**
- Reads actual code behavior
- Compares to documentation claims
- Updates docs to match reality
- Removes outdated information

**Example:**
```
Code now supports:
  - PostgreSQL 11, 12, 13, 14, 15

README still says:
  - "Supports PostgreSQL 11-13"

Fix updates README:
  - "Supports PostgreSQL 11-15"
```

## Execution Examples

### Example 1: Fix Proto Validation

```bash
@fix-openmcf-component GcpCertManagerCert --explain "primaryDomainName should allow wildcards *.example.com"
```

**What happens:**
```
1. Analysis (30 sec)
   - Reads spec.proto: Pattern rejects wildcards
   - Reads tests: No wildcard tests
   - Reads examples: No wildcard examples
   - Plan: Update pattern, add tests, add examples

2. Source Code Fix (2 min)
   - Update spec.proto pattern: ^(\*\.)?[a-z0-9-]+\.[a-z]{2,}$
   - Regenerate stubs: make protos ✓
   - Add 2 wildcard tests to spec_test.go
   - Run tests: 18/18 pass ✓

3. Documentation Propagation (3 min)
   - Add wildcard examples (2 new)
   - Update README (add wildcard to features)
   - Update docs/README.md (update comparison)
   - Update IaC READMEs (add examples)

4. Consistency Validation (1 min)
   - Proto ↔ TF: ✓ Match
   - Examples ↔ Schema: ✓ All validate
   - Pulumi ↔ TF: ✓ Parity
   - Tests ↔ Validations: ✓ Complete

5. Final Validation (2 min)
   - Component tests: ✓ 18/18 pass
   - Build: ✓ Success
   - Full tests: ✓ 156/156 pass

Total: 8 minutes
```

### Example 2: Fix IaC Implementation

```bash
@fix-openmcf-component AwsRdsInstance --explain "backup_retention_period hardcoded to 7, should use spec.backupRetentionDays"
```

**What happens:**
```
1. Analysis
   - Reads Pulumi: BackupRetentionPeriod: pulumi.Int(7)
   - Reads Terraform: backup_retention_period = 7
   - Reads spec.proto: int32 backup_retention_days = 8;
   - Issue: Both hardcoded instead of using spec

2. Source Code Fix
   - Pulumi: BackupRetentionPeriod: pulumi.Int(int(spec.BackupRetentionDays))
   - Terraform: backup_retention_period = var.backup_retention_days
   - Tests: Add test for different retention values
   - Run: ✓ Tests pass

3. Documentation Propagation
   - Examples: Show various retention periods (7, 14, 30 days)
   - README: Document backup_retention_days field
   - Overview: Explain backup behavior
   - Validate: Examples work

4. Consistency
   - Pulumi ↔ TF: ✓ Both use spec field now
   - Tests validate various values: ✓

5. Validation
   - Tests: ✓ Pass
   - Build: ✓ Success
```

### Example 3: Fix Documentation Only

```bash
@fix-openmcf-component PostgresKubernetes --explain "examples.md uses deprecated 'database_name', should be 'db_identifier'"
```

**What happens:**
```
1. Analysis
   - Reads spec.proto: Confirms field is 'db_identifier'
   - Reads examples: Uses 'database_name' (wrong!)
   - Code is correct, docs are wrong

2. Source Code Fix
   - No changes needed (code is already correct)

3. Documentation Update
   - Update examples.md: database_name → db_identifier
   - Update README: Fix any references
   - Validate examples: ✓ All pass

4. Consistency
   - Examples ↔ Proto: ✓ Now match

5. Validation
   - No code changed, tests still pass: ✓
```

### Example 4: Fix Test Logic

```bash
@fix-openmcf-component MongodbAtlas --explain "test expects error for empty cluster_tier but proto has no validation rule"
```

**What happens:**
```
1. Analysis
   - Reads spec.proto: No required rule on cluster_tier
   - Reads test: Expects error for empty cluster_tier
   - Decision needed: Should it be required?

2. Source Code Fix (assuming should be required)
   - Add validation to spec.proto:
     string cluster_tier = 3 [(buf.validate.field).required = true];
   - Regenerate stubs
   - Test now correct (no changes needed)

3. Documentation Update
   - Examples: Ensure all have cluster_tier
   - README: Note cluster_tier is required
   - Docs: Update if relevant

4. Validation
   - Tests: ✓ Now pass
   - Build: ✓ Success
```

## Best Practices

### Writing Good Fix Explanations

**Be specific:**
- ✅ "Field X validation rejects Y which should be valid"
- ✅ "Module hardcodes Z instead of using spec field"
- ✅ "Examples use old field name, update to current"

**Not vague:**
- ❌ "Fix validation"
- ❌ "Update module"
- ❌ "Sync docs"

### Source Code First

**Always:**
1. Fix code first (proto, IaC)
2. Validate code works
3. Update docs to match code
4. Validate docs match code

**Never:**
1. Update docs with desired behavior
2. Try to make code match docs

### Verify Consistency

After fix:
- [ ] Run component tests
- [ ] Validate examples
- [ ] Check Pulumi ↔ Terraform parity
- [ ] Review documentation accuracy
- [ ] Run full test suite

### Test the Fix

```bash
# After fix completes
cd apis/org/openmcf/provider/<provider>/<component>/v1/iac/hack/

# Test with Pulumi
cd pulumi && make local && pulumi up

# Test with Terraform
cd ../tf && terraform init && terraform plan
```

## Troubleshooting

### Fix Breaks Tests

```
❌ Tests failed after fix

Failed: TestDomainValidation
Error: expected error, got nil

Analysis:
  - Fix relaxed validation (now allows more)
  - Test expects strict validation (old behavior)
  
Auto-Fix:
  - Update test to expect new behavior
  - Retry: ✓ Pass
```

### Fix Creates Inconsistency

```
⚠️  Inconsistency detected after fix

Proto has new field but:
  ❌ Terraform variables.tf missing field
  ❌ Examples don't show field

Auto-Fix:
  - Add to variables.tf
  - Add to examples
  - Verify: ✓ Consistent
```

### Examples Don't Validate

```
❌ Examples validation failed

Example 3: Field 'new_field' not found in schema

Analysis:
  - Example uses field that doesn't exist
  - Likely typo or incomplete fix

Auto-Fix:
  - Check spec.proto for correct field name
  - Update example
  - Validate: ✓ Pass
```

## Success Criteria

After fix completes:

✅ Fix applied to source code
✅ All related artifacts updated
✅ Consistency verified (5 checks pass)
✅ Component tests pass
✅ Build succeeds
✅ Full test suite passes
✅ Examples validate
✅ Documentation accurate
✅ No regressions
✅ Ready to commit

## Integration

### With Audit

```bash
# Audit identifies issue
@audit-openmcf-component MyComponent
# Report: "Examples use deprecated fields"

# Fix it
@fix-openmcf-component MyComponent --explain "update examples to use current field names"

# Verify fix
@audit-openmcf-component MyComponent
# Score maintained or improved
```

### With Complete

```bash
# Complete fills gaps
@complete-openmcf-component MyComponent

# Then fix specific issue
@fix-openmcf-component MyComponent --explain "validation logic has bug"

# Result: 95%+ with bug fixed
```

### With Update

```bash
# Update for general improvements
@update-openmcf-component MyComponent --scenario refresh-docs

# Fix specific bug discovered
@fix-openmcf-component MyComponent --explain "discovered validation bug while refreshing docs"
```

## Tips

### Effective Fixes

1. **Be surgical** - Fix one thing well
2. **Think cascading** - Consider what else needs updating
3. **Validate rigorously** - Run all tests
4. **Document changes** - Update all relevant docs
5. **Maintain parity** - Keep Pulumi = Terraform

### Avoiding Pitfalls

1. ❌ Don't fix docs without fixing code
2. ❌ Don't skip test updates
3. ❌ Don't break feature parity
4. ❌ Don't leave examples broken
5. ❌ Don't commit without validation

## Related Commands

- `@audit-openmcf-component` - Check status, identify issues
- `@update-openmcf-component` - General improvements
- `@complete-openmcf-component` - Fill all gaps
- `@forge-openmcf-component` - Create new component

## Reference

- **Ideal State:** `architecture/deployment-component.md`
- **Fix Rule:** `_rules/deployment-component/fix/fix-openmcf-component.mdc`
- **Master README:** `_rules/deployment-component/README.md`

---

**Remember:** Source code is truth, documentation describes truth. Fix the code first, then sync everything else!

**Ready to fix?** Run `@fix-openmcf-component <ComponentName> --explain "<what needs fixing>"` for targeted fixes with automatic propagation!

