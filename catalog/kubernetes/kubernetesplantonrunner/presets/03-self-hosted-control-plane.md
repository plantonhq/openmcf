# Self-Hosted Control Plane

This preset joins the runner to a self-hosted control plane instead of
the hosted one. The endpoint is the ONE bootstrap coordinate the join
cannot deliver -- everything else the runner needs (work queue, tunnel,
API endpoints) arrives in the join response, so a single `host:port`
plus a token from your instance is the entire difference from the
standard cluster runner.

## When to Use

- Your organization runs its own control plane and runners must enroll
  against it
- Air-gapped or compliance-bound environments where the hosted control
  plane is not an option

## Key Configuration Choices

- **`controlPlaneEndpoint: planton.example.com:443`** -- replace with
  your instance's host:port; no scheme prefix. When this field is unset
  the runner joins the hosted control plane, so this one line is what
  redirects enrollment.
- **The token comes from YOUR control plane** -- create it there with
  `planton runner token create` and store it under the managed-secret
  slug this manifest references; a token minted by one control plane is
  meaningless to another
- **Self-configuration still applies** -- the runner presents the token
  at your endpoint, registers itself, receives its own individually
  revocable identity, and picks up everything else from the join
  response; there is nothing else to point at your instance
- **Everything else matches the cluster runner** -- dedicated
  namespace, chart-default sizing, token carried in a Kubernetes Secret

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |
| `planton.example.com:443` | Your control plane's host:port | Your self-hosted instance's public endpoint |

## Related Presets

- `01-cluster-runner` -- the same shape joined to the hosted control plane
- `02-build-runner` -- adds the Tekton build worker for container-image builds
