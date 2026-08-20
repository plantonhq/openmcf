# Build Runner (Deploys + Image Builds)

This preset runs everything the cluster runner does and additionally
registers the runner as a build worker: container-image build pipelines
execute through Tekton on this cluster, close to the source of truth and
without shipping build context out of the network.

## When to Use

- Container images should be built inside your own cluster rather than
  on hosted build infrastructure
- The cluster already runs (or will run) Tekton Pipelines and you want
  one runner handling both deploys and builds

## Key Configuration Choices

- **`build.enabled: true`** -- the runner registers as a build worker in
  addition to its deploy duties; build pipelines run as Tekton resources
  on this cluster
- **Tekton Pipelines is a prerequisite** -- the module does not install
  it; the runner expects it present before build operations arrive.
  `tektonNamespace: tekton-pipelines` matches a stock install -- point it
  elsewhere if your Tekton lives in a different namespace (it defaults
  to the runner's own namespace when unset).
- **Token as a managed-secret reference** -- the token authorizes
  joining only; the runner registers itself on first boot and receives
  its own individually revocable identity, and revoking the token never
  touches runners it already admitted
- **Everything else matches the cluster runner** -- same namespace
  handling, same secret handling, chart-default sizing (consider sizing
  up if builds and deploys run concurrently)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |

## Related Presets

- `01-cluster-runner` -- deploys only, the leanest posture
- `03-self-hosted-control-plane` -- joins a self-hosted control plane instead of the hosted one
