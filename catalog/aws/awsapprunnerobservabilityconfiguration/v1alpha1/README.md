# AwsAppRunnerObservabilityConfiguration

Deploy and manage an AWS App Runner observability configuration using Planton -- the reusable, versioned tracing policy that App Runner services reference to enable distributed request tracing.

## Overview

When a service references an observability configuration, each of its instances runs an OpenTelemetry collector sidecar that forwards request spans to the configured vendor (AWS X-Ray is the only vendor App Runner supports today). The application must be instrumented with the AWS Distro for OpenTelemetry (ADOT) SDK for the collector to have anything to forward.

The configuration is deliberately its own resource rather than fields on the service:

- **Shared by many services** -- each service references the configuration by ARN, so a fleet adopts one tracing posture tuned in one place.
- **Versioned by AWS** -- the trace settings are create-time immutable; changing them registers a NEW revision under the same name. The exported ARN carries the revision, so a change rolls referencing services through the resource graph.

## When to Use

- Trace request flows across a fleet of App Runner services and the AWS APIs they call.
- Debug tail latency with X-Ray service maps and trace timelines.

## Minimal Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsAppRunnerObservabilityConfiguration
metadata:
  name: xray-tracing
  org: my-org
  env: prod
  id: xray-tracing-prod
spec:
  region: us-east-1
  traceConfiguration:
    vendor: AWSXRAY
```

## Spec Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | Yes | AWS region for the configuration. |
| `traceConfiguration` | object | No | Tracing settings. Omitted, the configuration is valid but inert. |
| `traceConfiguration.vendor` | string | No (default `AWSXRAY`) | The tracing vendor. `AWSXRAY` is the only supported value today. |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `configuration_arn` | The revision-carrying ARN services reference via `observabilityConfigurationArn`. |
| `configuration_revision` | The revision this deployment registered. |
| `latest` | Whether this revision is the latest under the configuration name. |

## Composing into a Service

```yaml
# In an AwsAppRunnerService spec -- the reference itself is the enable switch:
observabilityConfigurationArn:
  valueFrom:
    kind: AwsAppRunnerObservabilityConfiguration
    name: xray-tracing
    fieldPath: status.outputs.configuration_arn
```

The service's `instanceRoleArn` role needs X-Ray write permissions (the AWS-managed `AWSXRayDaemonWriteAccess` policy covers it).

## Deliberately Omitted

- **Per-kind tags**: identity tags derive from metadata; custom user tags are a platform-wide concern.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
