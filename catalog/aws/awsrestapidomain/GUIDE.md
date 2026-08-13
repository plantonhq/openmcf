# AwsRestApiDomain — Component Guide

Authored operational judgment for the REST API custom-domain component:
the design decisions behind the spec's shape, and what to know before
fronting APIs with your own hostname.

## Design decisions

- **A domain is not an API field.** One hostname maps many APIs and
  outlives any of them; putting it on AwsRestApiGateway would force
  every mapping to share an API lifecycle.
- **`routing_mode` is modeled; routing rules are not.** Rules are an
  API Gateway v2 surface that also attaches to v1 domains, and they
  already live on AwsHttpApiDomain. This component exposes the v1
  knob that chooses between base-path mappings and those rules.
- **DNS stays outside.** The domain exists independently of any
  record pointing at it. Alias targets and zone IDs are outputs so
  AwsRoute53DnsRecord can compose them.
- **Certificate source is exactly one.** ACM in-region for REGIONAL /
  PRIVATE, ACM in us-east-1 for EDGE — the modules wire the matching
  provider argument.

## Running custom domains in production

- **REGIONAL is the default for new work.** EDGE still exists for
  global CloudFront distribution, but it requires the certificate in
  us-east-1 and a longer create.
- **Keep the root mapping.** An empty `base_path` (the `(root)`
  output key) is how `https://api.example.com/` reaches a stage;
  forgetting it 404s the hostname itself.
- **PRIVATE domains need access associations.** Without them the
  hostname exists and nobody inside the VPC can call it.
- **Do not rotate the hostname in place.** Changing `domain_name`
  replaces the domain; stand up the new name, move DNS, then destroy
  the old one.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
