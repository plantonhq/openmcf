# Pulumi Module: DigitalOcean Droplet Autoscale Pool

Provisions a pool of identical droplets with static or utilization-driven scaling -- the complete `digitalocean_droplet_autoscale` resource surface. Behavioral parity with the Terraform module is the contract.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean.DropletAutoscale` | The pool: name + scaling config + the member droplet template |

## Inputs

`DigitalOceanDropletAutoscalePoolStackInput`: the target `DigitalOceanDropletAutoscalePool` resource and the DigitalOcean provider config (API token).

## Outputs

Exactly the `DigitalOceanDropletAutoscalePoolStackOutputs` contract: `pool_id` (Pulumi's resource id), `status`.

## Behavior notes

- The SDK's `Config` and `DropletTemplate` are singular objects (no one-element-array bridge quirk); only the chosen scaling branch's leaves are set -- nil pointers never reach the API.
- Member tags are spec.droplet_template.tags ∪ the standard Planton labels -- the exact set the Terraform module applies.
- `public_networking` is absent from the SDK's template args, which matches its dead-on-write state at the pinned provider (excluded from the spec).
- SSH key, VPC, and project references resolve to literal ids before the module runs.
