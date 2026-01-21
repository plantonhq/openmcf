---
title: "CLI Reference"
description: "Complete command-line reference for project-planton CLI - all commands, flags, and options"
icon: "terminal"
order: 1
---

# CLI Reference

Complete command-line reference for the `project-planton` CLI.

---

## Command Tree

```
project-planton
├── apply               Deploy infrastructure (unified, auto-detects provisioner)
├── destroy             Teardown infrastructure (or 'delete')
├── init                Initialize backend/stack (unified, auto-detects provisioner)
├── plan                Preview changes (or 'preview', unified, auto-detects provisioner)
├── refresh             Sync state with reality (unified, auto-detects provisioner)
├── pulumi              Manage infrastructure with Pulumi
│   ├── init           Initialize Pulumi stack
│   ├── preview        Preview infrastructure changes
│   ├── up             Deploy infrastructure (or 'update')
│   ├── refresh        Sync state with cloud reality
│   ├── destroy        Teardown infrastructure
│   ├── delete         Remove stack metadata (or 'rm')
│   └── cancel         Cancel ongoing Pulumi operation
├── tofu                Manage infrastructure with OpenTofu
│   ├── init           Initialize backend and providers
│   ├── plan           Preview infrastructure changes
│   ├── apply          Deploy infrastructure
│   ├── refresh        Sync state with cloud reality
│   └── destroy        Teardown infrastructure
├── validate            Validate manifest against schema
├── load-manifest       Load and display manifest with defaults
└── version             Show CLI version
```

---

## Top-Level Commands

### apply

**NEW!** Unified kubectl-style command to deploy infrastructure by automatically detecting the provisioner from the manifest label `project-planton.org/provisioner`.

**Usage**:

```bash
project-planton apply -f <file> [flags]
# or
project-planton apply -f <file> [flags]
```

**Example**:

```bash
# Auto-detect provisioner from manifest
project-planton apply -f database.yaml

# With kustomize
project-planton apply --kustomize-dir services/api --overlay prod

# With overrides
project-planton apply -f api.yaml --set spec.replicas=5
```

**How it works**:
1. Reads the `project-planton.org/provisioner` label from your manifest
2. Automatically routes to the appropriate provisioner (pulumi/tofu/terraform)
3. If label is missing, prompts you to select a provisioner interactively (defaults to Pulumi)

**Supported provisioners**: `pulumi`, `tofu`, `terraform` (case-insensitive)

### destroy

**NEW!** Unified kubectl-style command to destroy infrastructure. Works exactly like `apply` but tears down resources instead.

**Aliases**: `delete` (for kubectl compatibility)

**Usage**:

```bash
project-planton destroy -f <file> [flags]
project-planton delete -f <file> [flags]
```

**Example**:

```bash
# Auto-detect provisioner from manifest
project-planton destroy -f database.yaml

# Using kubectl-style delete alias
project-planton delete -f database.yaml

# With auto-approve (skips confirmation)
project-planton destroy -f api.yaml --auto-approve
```

### init

**NEW!** Unified command to initialize infrastructure backend or stack by automatically detecting the provisioner.

**Usage**:

```bash
project-planton init -f <file> [flags]
```

**Example**:

```bash
# Auto-detect provisioner from manifest
project-planton init -f database.yaml

# With kustomize
project-planton init --kustomize-dir services/api --overlay prod

# Reconfigure after backend changes
project-planton init -f app.yaml --reconfigure

# With tofu-specific backend config
project-planton init -f app.yaml --backend-type s3 --backend-config bucket=my-bucket
```

**How it works**:
1. Reads the `project-planton.org/provisioner` label from your manifest
2. Routes to appropriate initialization:
   - **Pulumi**: Creates stack if it doesn't exist
   - **Tofu**: Initializes backend and downloads providers
   - **Terraform**: Initializes backend and downloads providers

### plan

**NEW!** Unified command to preview infrastructure changes without applying them.

**Aliases**: `preview` (for Pulumi-style experience)

**Usage**:

```bash
project-planton plan -f <file> [flags]
project-planton preview -f <file> [flags]
```

**Example**:

```bash
# Auto-detect provisioner and preview changes
project-planton plan -f database.yaml

# Using preview alias (Pulumi-style)
project-planton preview -f database.yaml

# With kustomize
project-planton plan --kustomize-dir services/api --overlay staging

# Preview destroy plan (Tofu)
project-planton plan -f app.yaml --destroy
```

**How it works**:
1. Reads the `project-planton.org/provisioner` label from your manifest
2. Routes to appropriate preview operation:
   - **Pulumi**: Runs `pulumi preview`
   - **Tofu**: Runs `tofu plan`
   - **Terraform**: Not yet implemented

### refresh

**NEW!** Unified command to sync state with cloud reality without modifying resources.

**Usage**:

```bash
project-planton refresh -f <file> [flags]
```

**Example**:

```bash
# Auto-detect provisioner and refresh state
project-planton refresh -f database.yaml

# With kustomize
project-planton refresh --kustomize-dir services/api --overlay prod

# Show detailed diffs (Pulumi)
project-planton refresh -f app.yaml --diff
```

**How it works**:
1. Queries your cloud provider for current resource state
2. Updates state file to reflect reality
3. Does NOT modify any cloud resources (read-only operation)
4. Routes based on provisioner:
   - **Pulumi**: Runs `pulumi refresh`
   - **Tofu**: Runs `tofu refresh`
   - **Terraform**: Not yet implemented

### pulumi

Manage infrastructure using Pulumi as the IaC engine.

**Subcommands**: `init`, `preview`, `up`/`update`, `refresh`, `destroy`, `delete`/`rm`, `cancel`

**Documentation**: See [Pulumi Commands Reference](/docs/cli/pulumi-commands)

**Example**:

```bash
project-planton pulumi up -f database.yaml
```

### tofu

Manage infrastructure using OpenTofu/Terraform as the IaC engine.

**Subcommands**: `init`, `plan`, `apply`, `refresh`, `destroy`

**Documentation**: See [OpenTofu Commands Reference](/docs/cli/tofu-commands)

**Example**:

```bash
project-planton tofu apply -f database.yaml
```

### validate

Validate a manifest against its Protocol Buffer schema without deploying.

**Usage**:

```bash
project-planton validate -f <file> [flags]
```

**Example**:

```bash
# Validate single manifest
project-planton validate -f ops/resources/database.yaml

# With kustomize
project-planton validate \
  --kustomize-dir services/api/kustomize \
  --overlay prod

# If valid: exits with code 0, no output
# If invalid: shows detailed errors, exits with code 1
```

**Flags**:
- `-f <file>`: Path to manifest file
- `--kustomize-dir <dir>`: Kustomize base directory
- `--overlay <name>`: Kustomize overlay name

### load-manifest

Load a manifest and display it with defaults applied and overrides resolved.

**Usage**:

```bash
project-planton load-manifest -f <file> [flags]
```

**Example**:

```bash
# Load manifest and see defaults
project-planton load-manifest -f database.yaml

# Load with overrides
project-planton load-manifest \
  -f api.yaml \
  --set spec.replicas=5

# Load kustomize-built manifest
project-planton load-manifest \
  --kustomize-dir services/api/kustomize \
  --overlay prod
```

**Flags**: Same as `validate`

**Output**: YAML manifest with defaults filled in and overrides applied

### version

Show Project Planton CLI version information.

**Usage**:

```bash
project-planton version
```

**Example Output**:

```
project-planton version: v0.1.0
git commit: a1b2c3d
built: 2025-11-11T10:30:00Z
```

---

## Common Flags

These flags are available across multiple commands:

### Manifest Input

**`-f, -f <path>`**  
Path to manifest YAML file (local or URL). The `-f` shorthand is available for kubectl-style experience.

```bash
# Local file (kubectl-style)
-f ops/resources/database.yaml

# Local file (traditional)
-f ops/resources/database.yaml

# URL
-f https://raw.githubusercontent.com/myorg/manifests/main/db.yaml
```

**`--kustomize-dir <directory>`**  
Base directory containing kustomize structure.

```bash
--kustomize-dir services/api/kustomize
```

**`--overlay <name>`**  
Kustomize overlay environment to build (must be used with `--kustomize-dir`).

```bash
--overlay prod
```

**Priority**: `-f` > `--kustomize-dir` + `--overlay`

### Execution Control

**`--module-dir <path>`**  
Override IaC module directory (defaults to current directory).

```bash
--module-dir ~/projects/custom-modules/my-module
```

**`--set <key>=<value>`**  
Override manifest field values at runtime (repeatable).

```bash
--set spec.replicas=5 \
--set spec.container.image.tag=v2.0.0
```

### Pulumi-Specific Flags

**`--stack <org>/<project>/<stack>`**  
Override stack FQDN (instead of using manifest label).

```bash
--stack my-org/my-project/dev-stack
```

**`--yes`**  
Auto-approve operations without confirmation (Pulumi commands).

```bash
--yes
```

**`--force`**  
Force stack removal even if resources exist (`delete`/`rm` only).

```bash
--force
```

### OpenTofu/Terraform-Specific Flags

**`--auto-approve`**  
Skip interactive approval (`apply` and `destroy` commands).

```bash
--auto-approve
```

**`--reconfigure`**  
Reconfigure backend, ignoring any saved configuration. Use when backend configuration changes.

```bash
--reconfigure
```

**`--destroy`**  
Create destroy plan (`plan` command).

```bash
--destroy
```

### Provider Config Flag

**`-p, --provider-config <file>`**  
Path to provider credentials file. The provider type is automatically detected from the manifest's `apiVersion` and `kind`.

```bash
# Single unified flag works for all providers
-p ~/.config/aws-creds.yaml      # For AWS resources
-p ~/.config/gcp-creds.yaml      # For GCP resources
-p ~/.config/azure-creds.yaml    # For Azure resources
-p ~/.kube/config                # For Kubernetes resources
```

**How it works**: When you run a command with a manifest, the CLI:
1. Parses the manifest's `apiVersion` (e.g., `aws.project-planton.org/v1`)
2. Determines the required provider (e.g., AWS)
3. Uses the credentials from the file you specified with `-p`

This means you don't need to remember provider-specific flags - just use `-p` with the appropriate credentials file.

---

## Environment Variables

### Respected by CLI

**`TF_LOG`** / **`PULUMI_LOG_LEVEL`**  
Enable verbose logging for debugging.

```bash
export TF_LOG=DEBUG
export PULUMI_LOG_LEVEL=3
```

### Provider Credentials

See [Credentials Guide](/docs/guides/credentials) for complete list of provider-specific environment variables.

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (validation failed, deployment failed, etc.) |

---

## Configuration Files

### Manifest Download Directory

Downloaded URL manifests are cached in:

```
~/.project-planton/manifests/downloaded/
```

### Module Cache

Cloned IaC modules are cached in:

```
~/.project-planton/modules/
```

---

## Examples by Use Case

### First Deployment

```bash
# Unified kubectl-style (recommended)
project-planton validate -f database.yaml
project-planton apply -f database.yaml

# Using Pulumi directly
project-planton validate -f database.yaml
project-planton pulumi up -f database.yaml

# Using OpenTofu directly
project-planton validate -f database.yaml
project-planton tofu init -f database.yaml
project-planton tofu plan -f database.yaml
project-planton tofu apply -f database.yaml
```

### Multi-Environment Deployment

```bash
# Deploy across environments with unified command
for env in dev staging prod; do
    project-planton apply \
        --kustomize-dir services/api/kustomize \
        --overlay $env \
        --yes
done
```

### CI/CD Deployment

```bash
# Non-interactive with dynamic values (unified command)
project-planton apply \
  -f deployment.yaml \
  --set spec.container.image.tag=$CI_COMMIT_SHA \
  --yes
```

### Testing Local Module Changes

```bash
# Point to local module during development
project-planton pulumi preview \
  -f test.yaml \
  --module-dir ~/dev/my-module
```

---

## Related Documentation

- [Pulumi Commands](/docs/cli/pulumi-commands) - Detailed Pulumi command guide
- [OpenTofu Commands](/docs/cli/tofu-commands) - Detailed OpenTofu command guide
- [Manifest Structure](/docs/guides/manifests) - Understanding manifests
- [Credentials Guide](/docs/guides/credentials) - Setting up cloud credentials
- [Advanced Usage](/docs/guides/advanced-usage) - Power user techniques

---

## Getting Help

**Command help**:

```bash
# General help
project-planton --help

# Command-specific help
project-planton pulumi --help
project-planton tofu apply --help
```

**Found an issue?** [Open an issue](https://github.com/plantonhq/project-planton/issues)

**Need support?** Check existing issues or discussions

