# OpenMCF CLI - Architecture Reference

Technical reference for the OpenMCF CLI architecture and implementation.

---

## Overview

The OpenMCF CLI is a Go-based command-line tool built with [Cobra](https://github.com/spf13/cobra) that provides a unified interface for deploying infrastructure across multiple cloud providers using either Pulumi or OpenTofu.

**Core Responsibility**: Orchestrate the deployment workflow (load → validate → transform → delegate to IaC engine).

---

## Architecture

### High-Level Flow

```
┌─────────────┐
│  User Input │  (manifest, flags, credentials)
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────────────────┐
│  CLI Layer (cmd/openmcf/)                       │
│  ├── Parse flags (cobra)                                │
│  ├── Resolve manifest path (file, URL, kustomize)       │
│  ├── Apply --set overrides                              │
│  └── Validate manifest (proto-validate)                 │
└────────┬────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────┐
│  Transform Layer (internal/, pkg/)                      │
│  ├── Convert manifest to IaC input                      │
│  │   ├── Pulumi: export as env var                      │
│  │   └── OpenTofu: generate tfvars JSON                 │
│  ├── Extract credentials                                │
│  └── Setup module directory                             │
└────────┬────────────────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────────────────────┐
│  IaC Engine (Pulumi or OpenTofu)                        │
│  ├── Reads manifest from env/tfvars                     │
│  ├── Executes deployment code                           │
│  ├── Manages state                                       │
│  └── Creates/updates cloud resources                    │
└─────────────────────────────────────────────────────────┘
```

---

## Directory Structure

```
cmd/openmcf/
├── root.go                      # Root command definition
├── root/
│   ├── pulumi.go                # Pulumi command group
│   ├── pulumi/
│   │   ├── init.go              # pulumi init
│   │   ├── preview.go           # pulumi preview
│   │   ├── update.go            # pulumi up/update
│   │   ├── refresh.go           # pulumi refresh
│   │   ├── destroy.go           # pulumi destroy
│   │   ├── delete.go            # pulumi delete/rm
│   │   ├── cancel.go            # pulumi cancel
│   │   └── README.md            # Pulumi commands docs
│   ├── tofu.go                  # OpenTofu command group
│   ├── tofu/
│   │   ├── init.go              # tofu init
│   │   ├── plan.go              # tofu plan
│   │   ├── apply.go             # tofu apply
│   │   ├── refresh.go           # tofu refresh
│   │   ├── destroy.go           # tofu destroy
│   │   └── README.md            # OpenTofu commands docs
│   ├── load_manifest.go         # load-manifest command
│   ├── validate_manifest.go     # validate command
│   └── version.go               # version command
```

---

## Command Implementation Pattern

### Standard Handler Pattern

All command handlers follow this pattern:

```go
func commandHandler(cmd *cobra.Command, args []string) {
    // 1. Parse flags
    manifest, err := cmd.Flags().GetString(string(flag.Manifest))
    flag.HandleFlagErrAndValue(err, flag.Manifest, manifest)

    // 2. Load manifest
    cliprint.PrintStep("Loading manifest...")
    targetManifest, isTemp, err := climanifest.ResolveManifestPath(cmd)
    if err != nil {
        log.Fatalf("failed to resolve manifest: %v", err)
    }
    if isTemp {
        defer os.Remove(targetManifest)
    }
    cliprint.PrintSuccess("Manifest loaded")

    // 3. Validate
    cliprint.PrintStep("Validating manifest...")
    if err := manifest.Validate(targetManifest); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    cliprint.PrintSuccess("Manifest validated")

    // 4. Execute IaC operation
    cliprint.PrintHandoff("Pulumi")  // or "OpenTofu"
    err = executeOperation(targetManifest, ...)
    if err != nil {
        log.Fatalf("failed to execute: %v", err)
    }
}
```

### CLI Output Formatting

Uses `internal/cli/cliprint` for consistent output:

```go
cliprint.PrintStep("Loading manifest...")        // ● with spinner
cliprint.PrintSuccess("Manifest loaded")         // ✔ with checkmark
cliprint.PrintHandoff("Pulumi")                  // 🤝 handoff message
```

---

## Flag Handling

### Persistent Flags (Inherited by Subcommands)

Defined in `root/pulumi.go` and `root/tofu.go`:

```go
// Manifest input
Pulumi.PersistentFlags().String(string(flag.Manifest), "", "manifest file path")
Pulumi.PersistentFlags().String(string(flag.KustomizeDir), "", "kustomize directory")
Pulumi.PersistentFlags().String(string(flag.Overlay), "", "kustomize overlay")

// Execution control
Pulumi.PersistentFlags().String(string(flag.ModuleDir), pwd, "module directory")
Pulumi.PersistentFlags().StringToString(string(flag.Set), map[string]string{}, "overrides")

// Credentials
Pulumi.PersistentFlags().String(string(flag.AwsProviderConfig), "", "AWS credential file")
Pulumi.PersistentFlags().String(string(flag.GcpProviderConfig), "", "GCP credential file")
// ... more providers
```

### Command-Specific Flags

Defined in individual command files:

```go
// pulumi/delete.go
func init() {
    Delete.PersistentFlags().Bool(string(flag.Force), false, "force removal")
}

// tofu/apply.go
func init() {
    Apply.PersistentFlags().Bool(string(flag.AutoApprove), false, "skip approval")
}
```

---

## Integration with Internal Packages

### Manifest Resolution

`internal/cli/manifest/resolve.go` handles manifest source priority:

```
--manifest (direct file/URL)
    ↓ (if not provided)
--kustomize-dir + --overlay (build kustomize)
    ↓ (if not provided)
error (no manifest source)
```

### Manifest Operations

`internal/manifest/` provides:

- `LoadManifest()`: Load from file/URL
- `Validate()`: Run proto-validate
- `LoadWithOverrides()`: Apply --set flags
- `ApplyOverridesToFile()`: Create temp file with overrides

### Cloud Resource Kind Reflection

`pkg/crkreflect/` provides:

- Kind string → CloudResourceKind enum mapping
- CloudResourceKind → proto.Message type mapping
- Kind → provider mapping
- Kind metadata (name, group, version)

### IaC Execution

**Pulumi**: `pkg/iac/pulumi/pulumistack/run.go`

- Exports manifest as `PROJECT_PLANTON_MANIFEST` env var
- Runs `pulumi <command>` in module directory
- Streams output to user

**OpenTofu**: `pkg/iac/tofu/tofumodule/run_operation.go`

- Converts manifest to `variables.tf.json`
- Sets provider env vars
- Runs `tofu <command>` in module directory
- Streams output to user

---

## Build System

### Compilation

```bash
# Build CLI
go build -o bin/openmcf ./cmd/openmcf

# Or using Makefile
make build
```

### Dependencies

**Core**:

- `github.com/spf13/cobra`: CLI framework
- `google.golang.org/protobuf`: Proto handling
- `buf.build/go/protovalidate`: Validation

**Cloud Providers**:

- Pulumi AWS, GCP, Azure, etc. SDKs (in modules, not CLI)
- OpenTofu (external binary)

### Bazel Build

Project uses Bazel for builds:

```bash
# Build with Bazel
bazel build //cmd/openmcf

# Run tests
bazel test //...
```

---

## Testing

### Unit Tests

Each command handler should have tests:

```go
func TestCommandHandler(t *testing.T) {
    // Setup
    cmd := &cobra.Command{}
    cmd.Flags().String("manifest", "test.yaml", "")

    // Execute
    err := commandHandler(cmd, []string{})

    // Assert
    if err != nil {
        t.Errorf("expected no error, got: %v", err)
    }
}
```

### Integration Tests

Test full workflows:

```bash
# Create test manifest
cat > test-resource.yaml <<EOF
apiVersion: test.openmcf.org/v1
kind: TestResource
metadata:
  name: test
spec:
  value: test
EOF

# Test validation
openmcf validate --manifest test-resource.yaml

# Test load
openmcf load-manifest --manifest test-resource.yaml
```

---

## Adding New Commands

### Step 1: Create Command File

```go
// cmd/openmcf/root/pulumi/newcommand.go
package pulumi

import (
    "github.com/spf13/cobra"
)

var NewCommand = &cobra.Command{
    Use:   "newcommand",
    Short: "Description of new command",
    Run:   newCommandHandler,
}

func newCommandHandler(cmd *cobra.Command, args []string) {
    // Implementation
}
```

### Step 2: Register Command

```go
// cmd/openmcf/root/pulumi.go
func init() {
    // ...
    Pulumi.AddCommand(
        pulumi.Init,
        pulumi.Preview,
        pulumi.NewCommand,  // Add here
    )
}
```

### Step 3: Add Tests

```go
// cmd/openmcf/root/pulumi/newcommand_test.go
func TestNewCommand(t *testing.T) {
    // Test implementation
}
```

### Step 4: Update Documentation

- Add to `cmd/openmcf/root/pulumi/README.md`
- Add to website docs if user-facing

---

## Development Workflow

### Local Development

```bash
# 1. Make changes
vim cmd/openmcf/root/pulumi/update.go

# 2. Build
go build -o bin/openmcf ./cmd/openmcf

# 3. Test
./bin/openmcf pulumi up --manifest test.yaml

# 4. Run tests
go test ./cmd/openmcf/...

# 5. Commit
git add cmd/openmcf/
git commit -m "feat: improve pulumi up command"
```

### Debugging

```go
// Add debug logging
import log "github.com/sirupsen/logrus"

log.SetLevel(log.DebugLevel)
log.Debugf("Loaded manifest: %+v", manifest)
```

---

## Design Principles

### 1. Consistent User Experience

All commands follow the same output pattern:

- Print steps with spinners
- Show success/failure clearly
- Hand off to IaC engine explicitly
- Stream IaC engine output unmodified

### 2. Fail Fast

Validate early, before expensive operations:

- Manifest loading: Fail if file not found
- Validation: Fail before calling cloud APIs
- Credential check: Fail before deployment starts

### 3. Idempotent Operations

Most operations are idempotent:

- `init`: Safe to run multiple times
- `validate`: Read-only, no side effects
- `up`/`apply`: Creates if missing, updates if exists

### 4. No Surprises

- Auto-approve requires explicit flag (`--yes` or `--auto-approve`)
- Destructive operations ask for confirmation
- Preview/plan operations never modify resources

---

## Related Documentation

- [CLI Reference](/docs/cli/cli-reference) - User-facing CLI reference
- [Manifest Package](../../internal/manifest/README.md) - Manifest loading
- [CRK Reflect Package](../../pkg/crkreflect/README.md) - Kind resolution

---

## Contributing

When modifying the CLI:

- Follow existing command patterns
- Add tests for new functionality
- Update both code README and website docs
- Maintain backwards compatibility
- Use consistent error handling
- Format output with `cliprint` package
- Add flag validation
- Handle temp file cleanup

