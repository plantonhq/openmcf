# Terraform Module: DigitalOcean Droplet Autoscale Pool

Provisions a pool of identical droplets with static or utilization-driven scaling -- the complete `digitalocean_droplet_autoscale` resource surface.

## Resources

| Resource | Purpose |
|---|---|
| `digitalocean_droplet_autoscale.pool` | The pool: name + scaling config + the member droplet template |

## Inputs

Generated `variables.tf` mirrors the `DigitalOceanDropletAutoscalePoolSpec` proto: `pool_name`, the scaling oneof as two optional objects (`static` / `dynamic`), and `droplet_template` with flattened reference strings for `ssh_keys` / `vpc` / `project_id`. Authentication uses `digitalocean_token` (sensitive).

## Outputs

Exactly the `DigitalOceanDropletAutoscalePoolStackOutputs` contract: `pool_id`, `status`.

## Behavior notes

- The config block renders ONLY the chosen scaling branch's leaves (null attributes are omitted, matching the API's zero-means-unset wire behavior); the spec oneof makes mixed shapes unrepresentable.
- Member tags are spec.droplet_template.tags ∪ the standard Planton labels -- the exact set the Pulumi module applies.
- `public_networking` is never rendered: dead on write at the pinned provider (declared but never sent).
- Create waits for the pool AND every member to reach active (up to 15 minutes upstream); delete is `DeleteDangerous` -- it destroys the member droplets and polls up to 1 minute for the 404.
- Import: `terraform import ... <pool_id>` (see `iac/import-map.yaml`; the template image reads back as a numeric id -- expect an image diff on the first post-import plan when the config uses a slug).
