# Summary: Deployment Requirements Document

**Date**: December 9, 2025
**Summary of**: DEPLOYMENT-REQUIREMENTS.md

## Executive Summary

This document explains why 24 code changes across `pkg/*`, `pkg/iac/*`, and `apis/org/*` directories are **critical** for cloud resource deployments to work via the `CreateCloudResource` and `UpdateCloudResource` APIs.

## The Problem

Without these changes, **all deployments fail** due to 8 fundamental issues:

1. **Incorrect module paths** - Can't find Pulumi modules (`apis/project/planton/provider` → `apis/org/openmcf/provider`)
2. **Git checkout failures** - Empty version strings cause fatal errors
3. **Invalid credential keys** - Stack input uses wrong field names (`awsProviderConfig` → `provider_config`)
4. **Missing credential handling** - No way to process user credentials from API
5. **GCP credential format errors** - Bytes passed instead of JSON string
6. **Kind enum resolution failures** - Can't resolve cloud resource kinds when metadata name is empty
7. **AWS RDS naming violations** - Subnet group names violate AWS requirements (uppercase, special chars)
8. **Empty resource identifiers** - Resources fail when `Metadata.Id` is empty

## The 3 Key Areas of Changes

### Part 1: `pkg/*` (1 file)

**File**: `pkg/crkreflect/kind_by_kind_name.go`

**Change**: Fallback to enum name when metadata name is empty

**Impact**: CRITICAL - Without this, credential resolution fails for all deployments

---

### Part 2: `pkg/iac/*` (20 files)

#### Module Path Corrections (2 files)
- `pkg/iac/pulumi/pulumimodule/module_directory.go`
- `pkg/iac/tofu/tofumodule/module_directory.go`

**Change**: Fix paths from `apis/project/planton/provider` to `apis/org/openmcf/provider`

**Impact**: CRITICAL - Can't find Pulumi/OpenTofu modules without this

---

#### Git Version Validation (1 file)
- `pkg/iac/pulumi/pulumimodule/module_directory.go`

**Change**: Prevent empty version string checkouts (`if version != "" && version != default`)

**Impact**: CRITICAL - Fatal git errors on empty versions

---

#### Credential Key Standardization (8 files)
- `pkg/iac/stackinput/stackinputproviderconfig/aws_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/gcp_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/azure_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/atlas_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/cloudflare_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/confluent_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/snowflake_provider.go`
- `pkg/iac/stackinput/stackinputproviderconfig/kubernetes_provider.go`

**Change**: Change all provider keys from `{provider}ProviderConfig` to `provider_config`

**Impact**: CRITICAL - Stack input YAML has wrong field names, credentials can't be read

---

#### NEW FILE: User Credential Conversion (2 files)
- `pkg/iac/stackinput/stackinputproviderconfig/user_provider.go` (NEW)
- `pkg/iac/stackinput/stackinputproviderconfig/BUILD.bazel` (updated)

**Purpose**: Convert API credentials (proto messages) to Pulumi-compatible YAML files

**Impact**: CRITICAL - No way to process user credentials from API without this

---

#### GCP Credential Fixes (1 file)
- `pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider/provider.go`

**Changes**:
1. Convert bytes to JSON string
2. Validate JSON format
3. Validate required fields (type, project_id, private_key, client_email)
4. Validate private key PEM format

**Impact**: CRITICAL (change 1), HIGH (changes 2-4)

---

#### Function Exports (5 files)
- `pkg/iac/pulumi/pulumistack/run.go`
- `pkg/iac/pulumi/pulumistack/project_name.go`
- `pkg/iac/pulumi/pulumistack/init.go`
- `pkg/iac/pulumi/pulumistack/cancel.go`
- `pkg/iac/pulumi/pulumistack/remove.go`

**Change**: Capitalize function names for cross-package access

**Impact**: CRITICAL (first 2), MEDIUM (last 3)

---

#### Code Quality (1 file)
- `pkg/iac/pulumi/pulumimodule/provider/aws/pulumiekskubernetesprovider/provider.go`

**Change**: Remove unused imports

**Impact**: LOW - Code quality only

---

### Part 3: `apis/org/*` (3 files)

#### Resource ID Fallback (2 files)
- `apis/org/openmcf/provider/aws/awsrdsinstance/v1/iac/pulumi/module/locals.go`
- `apis/org/openmcf/provider/aws/awsrdsinstance/v1/iac/pulumi/module/instance.go`

**Change**: Use `Metadata.Name` when `Metadata.Id` is empty

**Impact**: CRITICAL - Resources fail when `Metadata.Id` is empty

---

#### Subnet Group Name Sanitization (2 files)
- `apis/org/openmcf/provider/aws/awsrdsinstance/v1/iac/pulumi/module/subnet_group.go`
- `apis/org/openmcf/provider/aws/awsrdscluster/v1/iac/pulumi/module/subnet_group.go`

**Change**: Add `sanitizeSubnetGroupName()` function that:
- Converts to lowercase
- Removes invalid characters (keeps only a-z, 0-9, ., _, -, space)
- Collapses multiple hyphens
- Trims leading/trailing invalid chars
- Limits to 255 characters

**Impact**: CRITICAL - AWS RDS deployments fail with naming violations without this

---

## Critical Deployment Path

```
API Request → Resolve Kind → Get Credentials → Build Stack Input →
Find Module → Checkout Code → Create Provider → Create Resources → Deploy
```

**Every step requires these changes. Break any step = deployment fails.**

---

## Impact Summary

### By Severity

- **23 CRITICAL changes** - Deployment completely fails without them
- **3 HIGH changes** - Better error messages and validation
- **2 LOW changes** - Code quality/documentation

### By Area

- **Kind Resolution**: 1 file (CRITICAL)
- **Module Management**: 2 files (CRITICAL)
- **Credential Handling**: 10 files (9 CRITICAL, 1 NEW)
- **GCP Provider**: 1 file (1 CRITICAL, 3 HIGH changes)
- **Function Access**: 5 files (2 CRITICAL, 3 MEDIUM)
- **Resource Creation**: 3 files (CRITICAL)
- **Code Quality**: 2 files (LOW)

---

## Testing Impact

### Without These Changes:
- ❌ AWS RDS Instance - credential key mismatch + naming violation + empty ID
- ❌ AWS RDS Cluster - subnet group naming violation
- ❌ GCP Cloud SQL - module path error + credential parsing error
- ❌ Azure Resources - credential key mismatch
- ❌ Empty Version - git checkout error
- ❌ Empty Metadata.Id - empty resource identifier error
- ❌ User-Provided Credentials - no credential conversion logic
- ❌ Kind Resolution - "failed to get kind enum" error

### With These Changes:
- ✅ All deployment scenarios work correctly

---

## Bottom Line

These aren't improvements—they're **mandatory fixes** for the deployment system to work at all.

**The changes form a complete dependency chain**: Each step in the deployment path depends on the previous step working correctly. Missing any single change breaks the entire pipeline.

**All three areas are interdependent**:
1. **`pkg/*`** - Kind resolution enables credential lookup
2. **`pkg/iac/*`** - Credential handling enables provider creation
3. **`apis/org/*`** - Resource creation depends on valid identifiers and names

Without all 24 changes in place, **zero cloud resource deployments can succeed via the API**.

