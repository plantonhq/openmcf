# Stack Outputs Authoring Guide

Purpose: define `stack_outputs.proto` for a resource with observable, non-secret outputs.

## Folder and Naming
- Path: `apis/project/planton/provider/<provider>/<kindfolder>/v1/stack_outputs.proto`
- Kind (PascalCase), folder `<kindfolder>` lowercased, no underscores.

## Syntax and Package
- `syntax = "proto3";`
- `package org.openmcf.provider.<provider>.<kindfolder>.v1;`
- No `go_package` option.

## Message
- Single message named `<Kind>StackOutputs`.

## Fields
- Include stable, observable outputs only: IDs, ARNs, URLs, hostnames, hosted zone IDs, ports.
- Avoid secrets (passwords, tokens, private keys).
- Prefer `string` for IDs/ARNs/URIs/hostnames; use numeric types only for numeric outputs (e.g., ports).
- Use `repeated` only when the resource inherently returns multiple values.
- Names in snake_case; brief comments.

## Example Skeleton
```proto
syntax = "proto3";
package org.openmcf.provider.aws.awscloudfront.v1;

// AwsCloudFrontStackOutputs captures observable identifiers from CloudFront.
message AwsCloudFrontStackOutputs {
  // distribution id
  string distribution_id = 1;
  // distribution domain name (e.g., d123.cloudfront.net)
  string domain_name = 2;
  // hosted zone id for alias records
  string hosted_zone_id = 3;
}
```

## Notes
- Derive from official provider docs or common IaC outputs.
- Keep schema practical (80/20). Add more outputs later if needed.
