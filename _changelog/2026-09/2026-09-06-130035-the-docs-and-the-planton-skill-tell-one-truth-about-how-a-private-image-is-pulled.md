# The docs and the Planton skill tell one truth about how a private image is pulled

## What changed

- **The site stops saying the deployment stage pulls with the registry connection.** The registry docs gain **Pulling Private Images**: the three ways a Kubernetes workload pulls (nothing, when the cluster's own identity reaches the registry; `pod.imageRegistries`, the login declared on the workload; `pod.imagePullSecrets`, a Secret declared beside it), what the service deploy fills onto the workload from the registry connection and exactly when it fills nothing (a trusted arm without a GHCR pull token; ECR on any arm), the GHCR read-only pull token, the runtimes that pull only from their own cloud's registry (Cloud Run, App Runner, Lambda) with the pull-through kinds that bridge them, and the refusal and rollout-remedy sentences. The Kubernetes-clusters page says how each cluster kind's identity pulls from its cloud's registry; the deployment-stage page gains **How the Image Is Pulled**; the registry tutorial's excerpt and overview stop promising pulls through the connection, and its GHCR path ends with the pull token.
- **Five catalog kinds state their own pull posture in their spec docs.** Cloud Run and Cloud Run Job image fields say private images come only from Artifact Registry and name the remote-repository route; Lambda's `image_uri` says ECR-only in the same Region and names the pull-through cache; the GKE node-pool `oauth_scopes` comment says `devstorage.read_only` is what lets the kubelet pull; the EKS node role comment says `AmazonEC2ContainerRegistryReadOnly` is the only way a cluster ever pulls from ECR. The ExternalSecret catalog page links its docker-registry template to `pod.imagePullSecrets` with the explicit kind and output path. Stubs, reference pages, and the proto-docs index regenerated.
- **The Planton skill gains `service.pulling-private-images.md`** -- the whole pull story as the Assistant must act on it: the three ways and when to pick each, the fill table by connection arm, where the run records what happened, the GHCR pull token and every place it is authored, the own-registry runtimes, the wizard's **Private image?** dial, the sentences to quote back with the fix for each (the unresolved-reference refusal, the literal password, `ImagePullBackOff` with its diagnosis order, the offline preflight, the duplicate server), previews and rotation, and what never to suggest. Seven references gain the pull beat where their topic meets it (configuring deployments, detection-first registration, build failures, offline deploy, previews, config references, connecting registries). `SKILL.md` stays at its ceiling.

## Why

Code that declares a pull login on the workload and fills it in the open was only half the promise. Every document a person or an agent would read still described a platform that derived pull credentials out of sight, and the skill said nothing about pulls at all -- so an Assistant asked "why won't my pod pull?" improvised from training: a docker-config on a machine, a `kubectl create secret`, a pasted password. The repositories now read as if the registry connection was always how builds push, and pulls were always declared in the manifest.

## How to check

```bash
go run ./pkg/skills/defspack && go test ./pkg/skills/...                                          # the tree packages; the new reference is cited and SKILL.md is within the 500-line ceiling
go test ./pkg/explain/refgen/ ./pkg/protodocs/                                                     # the regenerated reference pages and docs index agree with the protos
for k in gcp/gcpcloudrun gcp/gcpcloudrunjob aws/awslambda gcp/gcpgkenodepool aws/awseksnodegroup; do go build ./catalog/$k/v1alpha1/... && go build -o /dev/null ./catalog/$k/iac/pulumi; done
rg -n "pull them during deployment|pulls the image from the registry|clusters pull images from for deployment" site/   # zero hits
```
