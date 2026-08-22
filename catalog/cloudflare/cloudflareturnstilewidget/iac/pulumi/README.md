# CloudflareTurnstileWidget Pulumi Module

Pulumi IaC module for provisioning a Cloudflare Turnstile widget — a privacy-preserving CAPTCHA alternative that yields a public site key and a server-side secret.

## Architecture

```
main.go (entrypoint)
  └── module/
        ├── main.go              — Resources() orchestrator
        ├── locals.go            — Locals struct and initialization
        ├── outputs.go           — Stack output key constants
        └── turnstile_widget.go  — widget creation
```

## How It Works

1. `main.go` loads the `CloudflareTurnstileWidgetStackInput` from the `STACK_INPUT` environment variable (base64-encoded YAML).
2. `module.Resources()` initializes locals, creates a Cloudflare provider, and provisions the widget. Optional flags are omitted when unset so the provider applies its own defaults.
3. Stack outputs are exported matching `CloudflareTurnstileWidgetStackOutputs`. `secret` is a secret; `sitekey` is the API identity.

## Local Development

```bash
# Build the binary
make build

# Preview with test manifest
make test

# Or use debug.sh for a specific manifest
./debug.sh ../../e2e/manifest.yaml
```

## Dependencies

- `github.com/pulumi/pulumi-cloudflare/sdk/v6` — Cloudflare Pulumi provider
- `github.com/pulumi/pulumi/sdk/v3` — Pulumi SDK
- `github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule` — Shared stack input loading and provider wiring
