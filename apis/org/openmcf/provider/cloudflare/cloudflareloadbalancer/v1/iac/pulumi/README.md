# Pulumi Module: Cloudflare Load Balancer

This directory contains the Pulumi module for deploying Cloudflare Load Balancers using OpenMCF.

## Overview

The Pulumi module automatically provisions and links three Cloudflare resources:

1. **LoadBalancerMonitor** (account-level): Health check configuration
2. **LoadBalancerPool** (account-level): Group of origin servers
3. **LoadBalancer** (zone-level): DNS hostname and routing policy

You provide a simplified `CloudflareLoadBalancerStackInput` proto message - the module handles the resource dependency graph and ID wiring.

## Directory Structure

```
iac/pulumi/
├── README.md              # This file - deployment guide
├── overview.md            # Architecture and design decisions
├── main.go                # Pulumi entrypoint
├── Pulumi.yaml            # Pulumi project configuration
├── Makefile               # Build and deployment helpers
├── debug.sh               # Debug script for local testing
└── module/
    ├── main.go            # Module entrypoint (Resources function)
    ├── locals.go          # Local variables and helpers
    ├── load_balancer.go   # Core resource provisioning logic
    └── outputs.go         # Output constant definitions
```

## Prerequisites

1. **Cloudflare Account**:
   - Active Cloudflare account with Load Balancing add-on enabled
   - Cloudflare API token with permissions:
     - Load Balancers: Edit
     - Zones: Read
     - Account Settings: Read

2. **Required Environment Variables**:
   ```bash
   export CLOUDFLARE_API_TOKEN="your-api-token-here"
   export CLOUDFLARE_ACCOUNT_ID="your-account-id"       # Account-level resources
   export CLOUDFLARE_ZONE_ID="your-zone-id"             # Zone-level resource
   ```

3. **Pulumi CLI**:
   ```bash
   # macOS
   brew install pulumi

   # Linux
   curl -fsSL https://get.pulumi.com | sh

   # Verify installation
   pulumi version
   ```

4. **Go SDK**:
   - Go 1.21 or later
   - Pulumi Cloudflare provider SDK (auto-installed)

## Deployment

### Step 1: Create Stack Input

Create a YAML file defining your load balancer configuration:

```yaml
# cloudflare-lb-stack-input.yaml
target:
  metadata:
    name: api-lb
  spec:
    hostname: api.example.com
    zone_id:
      value: "abc123def456"
    origins:
      - name: primary
        address: "203.0.113.10"
        weight: 1
      - name: secondary
        address: "198.51.100.20"
        weight: 1
    proxied: true
    health_probe_path: "/health"
    session_affinity: 1  # SESSION_AFFINITY_COOKIE
    steering_policy: 0   # STEERING_OFF (failover)

provider_config:
  # Cloudflare credentials are provided via environment variables
```

### Step 2: Initialize Pulumi Stack

```bash
cd iac/pulumi

# Initialize a new stack (first time only)
pulumi stack init dev

# Set Pulumi config (optional - uses env vars by default)
pulumi config set cloudflare:apiToken --secret "${CLOUDFLARE_API_TOKEN}"
```

### Step 3: Deploy

```bash
# Preview changes
pulumi preview --stack dev

# Apply changes
pulumi up --stack dev
```

**Expected Output**:

```
Updating (dev)

View in Browser (Ctrl+O): https://app.pulumi.com/...

     Type                                   Name                    Status
 +   pulumi:pulumi:Stack                    cloudflare-lb-dev       created
 +   ├─ cloudflare:index:LoadBalancerMonitor  monitor                created
 +   ├─ cloudflare:index:LoadBalancerPool     pool                   created
 +   └─ cloudflare:index:LoadBalancer         load_balancer          created

Outputs:
    load_balancer_cname_target: "api.example.com"
    load_balancer_dns_record_name: "api.example.com"
    load_balancer_id: "a1b2c3d4e5f6"

Resources:
    + 3 created

Duration: 12s
```

### Step 4: Verify Deployment

```bash
# View stack outputs
pulumi stack output

# Test the load balancer
curl https://api.example.com/health
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `CLOUDFLARE_API_TOKEN` | Yes | Cloudflare API token with Load Balancer permissions |
| `CLOUDFLARE_ACCOUNT_ID` | Yes | Cloudflare account ID (for pool and monitor) |
| `CLOUDFLARE_ZONE_ID` | Optional | Can be provided in stack input instead |
| `PULUMI_ACCESS_TOKEN` | Optional | Required for Pulumi Cloud backend |

## Stack Outputs

After deployment, the following outputs are available:

- **`load_balancer_id`**: Cloudflare Load Balancer resource ID
- **`load_balancer_dns_record_name`**: The hostname (e.g., `api.example.com`)
- **`load_balancer_cname_target`**: CNAME target for DNS configuration

Access outputs:

```bash
# View all outputs
pulumi stack output

# Get specific output
pulumi stack output load_balancer_id
```

## Updating the Load Balancer

Modify your stack input YAML and re-run:

```bash
pulumi up --stack dev
```

Pulumi will show a diff of changes before applying.

### Common Updates

**Add a new origin**:
```yaml
origins:
  - name: primary
    address: "203.0.113.10"
  - name: secondary
    address: "198.51.100.20"
  - name: tertiary  # New origin
    address: "192.0.2.30"
```

**Change health check path**:
```yaml
health_probe_path: "/api/v1/healthz"  # Changed from /health
```

**Enable session affinity**:
```yaml
session_affinity: 1  # SESSION_AFFINITY_COOKIE
```

## Destroying the Load Balancer

```bash
# Preview what will be deleted
pulumi destroy --stack dev --preview

# Confirm and delete all resources
pulumi destroy --stack dev
```

**Warning**: This will delete the load balancer, pool, and monitor. Ensure no production traffic is routing through it.

## Debugging

### Enable Debug Mode

```bash
# Run with verbose logging
pulumi up --stack dev --logtostderr -v=9

# Or use the debug script
./debug.sh
```

### Debug Script

The `debug.sh` script provides detailed logging:

```bash
#!/bin/bash
set -x  # Print commands
export PULUMI_DEBUG_COMMANDS=true
export PULUMI_DEBUG_GRPC=true

pulumi up --stack dev --logtostderr -v=9
```

### Common Issues

**Issue**: `Error: could not create load balancer: zone_id is required`

**Solution**: Ensure `zone_id` is set in stack input or via environment variable

---

**Issue**: `Error: authentication error - invalid API token`

**Solution**: Verify `CLOUDFLARE_API_TOKEN` environment variable is set correctly

---

**Issue**: `Error: pool not found`

**Solution**: The pool resource may not have been created yet. This is a timing issue - Pulumi handles dependencies automatically, but check if the pool exists in Cloudflare dashboard.

---

**Issue**: Origins show as "Unhealthy"

**Solution**: 
1. Check origin servers are running and accessible
2. Verify `health_probe_path` returns HTTP 200
3. Ensure origin firewalls allow Cloudflare health check IPs

## Multi-Environment Deployment

Use separate Pulumi stacks for dev, staging, and production:

```bash
# Development
pulumi stack init dev
pulumi up --stack dev

# Staging
pulumi stack init staging
pulumi up --stack staging

# Production
pulumi stack init prod
pulumi up --stack prod
```

Each stack maintains independent state.

## State Management

By default, Pulumi stores state in Pulumi Cloud (free for individuals). For self-hosted state:

```bash
# Use S3 backend
pulumi login s3://my-pulumi-state-bucket

# Use local filesystem (not recommended for teams)
pulumi login file://~/.pulumi

# Use Azure Blob
pulumi login azblob://my-container

# Use Google Cloud Storage
pulumi login gs://my-pulumi-state-bucket
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Deploy Cloudflare Load Balancer

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - uses: pulumi/actions@v4
        with:
          command: up
          stack-name: prod
          work-dir: iac/pulumi
        env:
          PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          CLOUDFLARE_ACCOUNT_ID: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
```

## Testing

Run the module tests:

```bash
cd module
go test -v ./...
```

## Troubleshooting

### View Resource Details

```bash
# List all resources in stack
pulumi stack --show-urns

# View resource details
pulumi stack export
```

### Refresh State

If Cloudflare resources were modified outside of Pulumi:

```bash
pulumi refresh --stack dev
```

### Import Existing Resources

To import an existing Cloudflare Load Balancer:

```bash
pulumi import cloudflare:index/loadBalancer:LoadBalancer main <load-balancer-id>
```

## Best Practices

1. **Use environment-specific stacks**: Separate dev, staging, prod
2. **Store secrets securely**: Use Pulumi secrets or external vaults
3. **Enable stack tags**: Tag resources for cost allocation
4. **Use remote state backend**: Don't rely on local state files
5. **Implement CI/CD**: Automate deployments via GitHub Actions or similar
6. **Test in dev first**: Always validate changes in non-prod before prod deployment

## Additional Resources

- [Pulumi Cloudflare Provider Docs](https://www.pulumi.com/registry/packages/cloudflare/)
- [Cloudflare Load Balancer API Docs](https://developers.cloudflare.com/api/operations/load-balancers-create-load-balancer)
- [overview.md](./overview.md) - Architecture and design decisions
- [Component README](../../README.md) - User-facing component documentation

## Support

For issues or questions:
1. Check [Common Issues](#common-issues) above
2. Review [overview.md](./overview.md) for architectural context
3. Consult Cloudflare and Pulumi official documentation

---

**Ready to deploy?** Run `pulumi up --stack dev` to get started!

