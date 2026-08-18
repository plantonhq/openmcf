# AwsCloudMapNamespace

An AWS Cloud Map namespace — the service-discovery registry ECS services and custom applications look each other up in — with its services and statically registered instances managed in-line. Three namespace types shape the surface: HTTP (API-only discovery), PRIVATE_DNS (a private hosted zone in one VPC), PUBLIC_DNS (an internet-resolvable zone).

## Highlights

- **The type union is CELs**: PRIVATE_DNS requires its VPC and nothing else takes one; HTTP services carry no dns_config (no records exist to publish); Route 53 health checks live only in PUBLIC_DNS namespaces (the checkers are on the public internet).
- **Static registration is a first-class arm**: instances fold under their service with TYPED fields (ip, port, cname, alias to an ELB, EC2 instance) instead of the provider's magic-string attribute map — registering an RDS endpoint into a namespace so services discover it by DNS is a two-line chart story.
- **Contracts taught in place**: the HTTP namespace's description is ForceNew (the DNS namespaces update in place), the private namespace's VPC is never read back (imports carry it), an alias registration stands alone, and `force_destroy` deregisters runtime-registered instances too — never set it on ECS-managed services.

## Both Engines

Both modules render whichever namespace arm the spec's type selects, plus services and instances, and export the same outputs: `namespace_id` (import ID), `namespace_arn`, `hosted_zone_id`, `http_name`, plus the `service_ids`, `service_arns`, and `instance_service_ids` maps keyed like the spec entries.

## Chart Wiring

`vpc_id` → AwsVpc `vpc_id`; `instances.cname` → any endpoint output (an RDS address, an internal hostname); `instances.alias_dns_name` → AwsAlb `load_balancer_dns_name`; `instances.ec2_instance_id` → AwsEc2Instance `instance_id`. The `service_arns` map is what ECS service registries consume.
