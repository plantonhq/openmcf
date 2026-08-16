# AWS SSM Document

Your own runbook in AWS's hands: a reusable definition of commands to
run on managed nodes or automation steps to drive against AWS APIs —
versioned, shareable, and executable by State Manager, Run Command,
maintenance windows, and Automation.

## What Gets Managed

- The document content (JSON or YAML) with its declared parameters,
  the document type, and a target-type restriction.
- Versioning: every content change creates a new document version and
  promotes it to the default; version labels pin releases.
- Sharing with other accounts (or publicly) and artifact attachments
  for Package documents.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SSM permissions.

### AWS Account

- Nothing — documents stand alone. Executing them later needs managed
  nodes or automation targets.

## Deploy

### Console

Create the resource from the AWS catalog, paste the document content,
pick the type, and deploy.

### CLI

```bash
planton apply -f ssm-document.yaml
```

## After Deploy

- Bind the document to targets on a schedule with
  [AWS SSM Association](/cloud-catalog/aws-ssm-association), run it in
  a window with
  [AWS SSM Maintenance Window](/cloud-catalog/aws-ssm-maintenance-window),
  or invoke it ad hoc via Run Command.
- Content updates create new versions — associations pinning
  `$DEFAULT` follow automatically.
- Documents are free; only their executions consume resources.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
