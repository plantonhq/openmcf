# CivoComputeInstance Pulumi Module Architecture

## Module Design

The CivoComputeInstance module follows OpenMCF's standard pattern:
1. Protobuf spec → Pulumi resources → Cloud infrastructure
2. Abstraction over Civo provider complexity
3. Consistent multi-cloud API patterns

## Key Design Decisions

### 1. Single SSH Key/Firewall

**Limitation**: Civo provider supports one SSH key and one firewall per instance.

**Module Behavior**: Uses first item from `ssh_key_ids` and `firewall_ids` arrays.

**Rationale**: Protobuf spec allows arrays for future-proofing if Civo adds multi-key/firewall support.

**Workaround**: Add additional SSH keys via user_data cloud-init.

### 2. Foreign Key References

The module uses `StringValueOrRef` for network, firewalls, volumes, and reserved IPs.

**Benefits**:
- Reference other OpenMCF resources by name
- Runtime resolution of resource IDs
- Type-safe cross-resource dependencies

**Example**:
```yaml
network:
  ref:
    kind: CivoVpc
    name: prod-network
    # Resolves to: status.outputs.network_id
```

### 3. User Data vs Configuration Management

The module supports cloud-init scripts (user_data) for bootstrapping.

**Best Practices**:
- Use user_data for initial setup (install packages, enable services)
- Use config management (Ansible, Chef) for ongoing configuration
- Don't embed secrets in user_data (use environment variables)

## Resource Lifecycle

### Creation Flow

```
1. Pulumi parses stack input
2. Module initializes locals
3. Civo provider authenticated
4. Instance resource created:
   a. Civo API call: POST /v2/instances
   b. Instance enters BUILDING state
   c. Cloud-init executes (if user_data provided)
   d. Instance reaches ACTIVE state
5. Outputs exported to stack state
```

### Update Flow

Changing certain fields triggers instance replacement:
- `region` - Must recreate in new region
- `size` - Can resize in-place (but may require reboot)
- `image` - Must recreate with new OS
- `network` - Must recreate in new network

Use `create_before_destroy` for zero-downtime updates.

### Deletion Flow

```
1. Pulumi initiates destroy
2. Module deletes instance
3. Civo API call: DELETE /v2/instances/{id}
4. Instance destroyed (root disk deleted)
5. Volumes detached (persist if not explicitly deleted)
6. Reserved IP released (persists if not explicitly deleted)
```

## State Management

Pulumi tracks:
- Instance ID (UUID)
- All configuration parameters
- Current IPs and status
- Dependencies (network, firewall, volumes)

**State Backend**: Pulumi Service (encrypted) or self-hosted (S3, Azure Blob, GCS).

## Performance Characteristics

- **Cold start**: ~60 seconds (Civo provisioning)
- **With user_data**: +30-180 seconds (script execution)
- **Total**: 1-5 minutes typical

**Scaling**:
- 10 instances (parallel): ~90 seconds
- 100 instances: Limited by Civo API rate limits

## Security Model

### Secrets

- **Civo API Token**: Stored as Pulumi secret (encrypted at rest)
- **SSH Private Keys**: Never stored (only public key IDs referenced)
- **User Data**: Visible in Pulumi state (don't embed secrets)

### Best Practices

1. Separate API tokens per environment
2. Use SSH keys exclusively (no passwords)
3. Custom firewalls for all production instances
4. Private networks for multi-tier architectures

## Testing Strategy

### Unit Tests (`spec_test.go`)

- Validates buf.validate rules
- Tests: instance_name pattern, size/image validation, user_data limits, tag uniqueness

### Integration Tests

- Provision real instance in test account
- Verify SSH connectivity
- Verify cloud-init completion
- Destroy instance

## Extending the Module

### Adding New Fields

1. Update `spec.proto`
2. Regenerate: `make protos`
3. Update `instance.go` to use new field
4. Add validation tests
5. Update documentation

### Custom Instance Types

For Civo-specific features not in spec (e.g., GPU instances):
- Add fields to spec.proto
- Map to Civo provider args in instance.go
- Document limitations

## Troubleshooting

### Module Errors

| Error | Cause | Solution |
|-------|-------|----------|
| "network not found" | Invalid network ID | Verify with `civo network list` |
| "invalid size" | Size not in region | Check `civo size list --region <region>` |
| "image not available" | Image not in region | Verify `civo diskimage list --region <region>` |

### Cloud-Init Debugging

```bash
# Check cloud-init status
ssh root@<ip> "cloud-init status --wait"

# View cloud-init logs
ssh root@<ip> "cat /var/log/cloud-init-output.log"

# Re-run cloud-init (testing only)
ssh root@<ip> "cloud-init clean && cloud-init init"
```

## Related Resources

- **Civo Provider**: [pulumi.com/registry/packages/civo](https://www.pulumi.com/registry/packages/civo/)
- **Civo API**: [civo.com/api/instances](https://www.civo.com/api/instances)
- **Cloud-Init**: [cloudinit.readthedocs.io](https://cloudinit.readthedocs.io/)

## Conclusion

The CivoComputeInstance module provides production-ready instance provisioning on Civo Cloud with proper networking, security, and storage integration. It handles the 80% use case while allowing customization for advanced scenarios.
