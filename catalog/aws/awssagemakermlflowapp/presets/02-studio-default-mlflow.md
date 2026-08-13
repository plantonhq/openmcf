# Studio Default MLflow

This preset makes the app the default MLflow for a SageMaker domain —
Studio users in the domain track to it automatically, every logged
model lands in the SageMaker Model Registry, and maintenance stays in
a declared quiet window.

## When to Use

- A platform team providing MLflow to Studio users without circulating
  tracking URIs
- Workflows where model lineage must reach the Model Registry
  automatically

## What You Get

- The domain's Studio users tracking to this app by default
  (`defaultDomainIds`)
- Automatic model registration into the SageMaker Model Registry —
  and unlike the tracking server's boolean, this mode toggles both
  ways
- Maintenance held to Sundays 03:00 UTC

## Customize

- Replace the literal domain ID with a `valueFrom` reference to a
  composed `AwsSagemakerDomain`, or list more domains — associations
  update in place
- Add `accountDefaultStatus: ENABLED` to make this the account-wide
  default MLflow — but only ONE app per account can hold it, so decide
  which app owns the default before two manifests fight over it
- Everything here except `roleArn` updates in place — the role is the
  one replace-on-change field
