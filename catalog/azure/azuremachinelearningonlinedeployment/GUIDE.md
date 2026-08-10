# Azure Machine Learning Online Deployment -- Operational Guide

Judgment that saves real time when running managed online deployments. The field reference lives in the API Explorer; this is the operational layer above it.

## Deployments are cattle; name them for the rollout

The endpoint is the stable contract; deployments exist to be replaced. Name them for their rollout role (`blue`/`green`, or a model-version suffix), ship every model change as a NEW deployment, and shift the endpoint's traffic map -- never mutate a serving deployment's model in place if you can avoid it. An in-place model change goes through a full PUT that rolls every instance while it serves.

## Budget the quota BEFORE the first apply

Managed online endpoints draw from their OWN VM quota, separate from the regular compute quota you already requested for training. A deployment that exceeds it fails at provisioning with a quota error that looks like a capacity problem. Check the `Machine Learning managed online endpoint` quota for the instance family and region first -- and remember mirrored traffic and blue/green overlaps mean BOTH deployments' instances count at once.

## There is no scale-to-zero -- size the floor honestly

Every instance bills around the clock while the deployment exists. `instanceCount` updates in place (it rides the ARM SKU capacity -- the one change the service applies without rolling containers), so scale down to one instance in quiet hours if you must, but the floor is one. If a model is called rarely, batch scoring or an app-embedded model is the cheaper architecture -- do not park a VM behind an endpoint nobody calls.

## The endpoint's identity does the pulling, not the deployment

Provisioning failures that read like image or artifact errors (`unable to pull image`, model download failures) are almost always missing grants on the ENDPOINT's identity -- AcrPull on the registry, Storage Blob Data Reader on the model store. Fix the grants and retry; nothing in the deployment spec is wrong. This is why charts prefer a pre-granted user-assigned identity on the endpoint.

## Probes: trust the defaults until your container proves otherwise

The service's probe defaults suit the standard scoring images. Set them when reality demands it: a model that takes minutes to load wants a `startupProbe` with a generous `failureThreshold` (the startup probe gates the other two), and a request that legitimately runs long wants `requestSettings.requestTimeout` raised ABOVE the client's timeout, not just above the default 5 seconds. A deployment that flaps between healthy and unhealthy under load usually needs `maxConcurrentRequestsPerInstance` aligned with what the container can actually parallelize -- the default is 1 for a reason.

## Data collection is cheap insurance, but it lands in YOUR storage

`dataCollector` writes scoring payloads to workspace blob storage: drift detection needs it, debugging loves it, and storage bills grow with it. Start with `samplingRate` well below 1.0 on high-traffic deployments and let the rolling rate (`Hour` default) shard the paths. Payloads may carry sensitive data -- the collection destination inherits the workspace storage's access controls, so review who can read it.

## Secure egress changes the failure mode

`egressPublicNetworkAccessEnabled: false` forces image and model pulls through the workspace's managed network. It is the right hardened posture -- and it means a missing private endpoint on the registry or storage surfaces as a deployment provisioning failure, not a network error you can see. Prove the workspace's outbound path (registry, storage) before flipping it on a serving deployment.
