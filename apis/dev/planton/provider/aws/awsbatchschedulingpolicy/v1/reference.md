# AwsBatchSchedulingPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsBatchSchedulingPolicySpec defines an AWS Batch fair-share scheduling
policy: the rules that divide a job queue's compute capacity across share
identifiers instead of processing jobs strictly first-in-first-out.

Without a scheduling policy, one team's burst of ten thousand jobs starves
every other submitter on the queue. With one, jobs are submitted with a
share identifier (e.g. per team or per workload class) and Batch balances
dispatch so each share gets capacity proportional to its weight.

One policy is a standalone, shareable object: many queues can reference
the same policy through their scheduling_policy field, so an
organization's fairness rules are defined once and reused. The
weight/decay/reservation dials can all be updated in place on a live
policy.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsBatchSchedulingPolicy
metadata:
  name: test-batch-fair-share
  id: test-batch-fair-share
  org: test-org
  env: dev
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: test
    pulumi.planton.dev/project: test
    pulumi.planton.dev/stack.name: dev.AwsBatchSchedulingPolicy.test-batch-fair-share
spec:
  region: us-west-2
  computeReservation: 10
  shareDecaySeconds: 3600
  shareDistributions:
    - shareIdentifier: teamA
      weightFactor: 0.5
    - shareIdentifier: adhoc*
      weightFactor: 2
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.computeReservation` | `int32` |  |  |  |
| `spec.shareDecaySeconds` | `int32` |  |  |  |
| `spec.shareDistributions` | `[]AwsBatchShareDistribution` |  |  |  |
| `spec.shareDistributions[].shareIdentifier` | `string` | yes |  |  |
| `spec.shareDistributions[].weightFactor` | `double` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the scheduling policy is created. A policy can
only be attached to job queues in the same region.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.computeReservation

`int32` · optional (explicit presence)

The percentage (0-99) of the queue's capacity held back for share
identifiers that are NOT currently represented among running jobs --
headroom so a quiet team's first job does not wait behind a busy
team's backlog. The effective reservation is computed as
(compute_reservation/100)^N where N is the number of active shares,
so the held-back slice shrinks as more shares become active.

- rule: {"int32":{"lte":99,"gte":0}}

### spec.shareDecaySeconds

`int32` · optional (explicit presence)

The sliding window, in seconds (0-604800, up to 7 days), over which
past usage counts against a share's fair allocation. Longer windows
make fairness account for history ("you had the cluster all morning");
0 considers only currently-running jobs.

- rule: {"int32":{"lte":604800,"gte":0}}

### spec.shareDistributions

`[]AwsBatchShareDistribution`

The relative weight of each share identifier. Shares absent from this
list use weight 1.0. AWS allows up to 500 entries per policy.

- rule: {"repeated":{"maxItems":"500"}}
- rule: weight_factor must be between 0.0001 and 999.9999 when set (leave it at 0 to use the AWS default of 1.0)

### spec.shareDistributions[].shareIdentifier

`string` · required

The share identifier jobs carry at submission time (SubmitJob's
shareIdentifier). End with "*" to match a prefix -- e.g. "analytics*"
covers "analyticsDaily" and "analyticsAdhoc" as one share. Up to 255
characters, alphanumeric plus the trailing wildcard.

- rule: {"required":true,"string":{"pattern":"^[0-9A-Za-z]{0,254}[0-9A-Za-z*]?$"}}

### spec.shareDistributions[].weightFactor

`double`

The share's relative weight, 0.0001-999.9999. LOWER weight means MORE
capacity: a share with weight 0.5 receives twice the capacity of a
weight-1.0 share. Defaults to 1.0 when unset.

- rule: {"double":{"lte":999.9999,"gte":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsBatchSchedulingPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.schedulingPolicyArn` | `string` | The Amazon Resource Name (ARN) of the scheduling policy -- what job queues reference through their scheduling_policy field. |
| `status.outputs.schedulingPolicyName` | `string` | The scheduling policy's name (derived from metadata.name). |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
