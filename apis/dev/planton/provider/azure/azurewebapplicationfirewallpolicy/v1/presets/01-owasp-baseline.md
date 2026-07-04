# OWASP 3.2 Baseline

This preset creates the policy almost every gateway should start with: the
OWASP 3.2 core rule set in Prevention mode (Azure's default), no custom
rules, no overrides. It blocks SQL injection, cross-site scripting, remote
and local file inclusion, remote code execution, and protocol violations
out of the box.

## When to Use

- The first WAF policy for any Application Gateway on the WAF_v2 SKU
- Organizations standardizing one org-wide baseline policy that individual
  listeners or routes override where needed

## Key Configuration Choices

- **OWASP 3.2** -- the current production standard; the anomaly-scoring
  model reduces false positives compared with 3.0/3.1, and several policy
  dials (file-upload enforcement) only work with 3.2
- **Prevention mode (default)** -- matching requests are blocked; switch a
  new policy to `DETECTION` first if you want to watch real traffic before
  enforcing
- **No tuning yet** -- add `ruleGroupOverrides` and `exclusions` as real
  traffic surfaces false positives (see the detection-tuning preset)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the policy in | The resource group's `status.outputs.resource_group_name` |
| `<policy-name>` | The policy's name, unique within the resource group | Your naming convention |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |

## Downstream Wiring

Attach the policy gateway-wide:

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureApplicationGateway
metadata:
  name: web-gateway
spec:
  sku: WAF_V2
  firewallPolicyId:
    valueFrom:
      kind: AzureWebApplicationFirewallPolicy
      name: waf-baseline
      fieldPath: status.outputs.policy_id
```
