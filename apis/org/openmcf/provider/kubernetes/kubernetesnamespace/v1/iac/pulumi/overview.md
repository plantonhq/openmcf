# Kubernetes Namespace - Pulumi Module Architecture

## Purpose

This Pulumi module implements the "Namespace-as-a-Service" pattern, automating the creation of production-ready Kubernetes namespaces with complete resource quotas, network policies, access controls, and optional service mesh integration.

## Architecture

### Component Hierarchy

```
KubernetesNamespace (API Resource)
    │
    ├── Kubernetes Namespace
    │   ├── Labels: managed-by, resource, resource-kind, team, environment, cost-center
    │   ├── Annotations: service mesh injection, TTL, custom
    │   └── Pod Security Standards: enforcement label
    │
    ├── ResourceQuota (if enabled)
    │   ├── CPU: requests + limits
    │   ├── Memory: requests + limits
    │   └── Object Counts: pods, services, configmaps, secrets, PVCs, load balancers
    │
    ├── LimitRange (if custom limits specified)
    │   └── Default Container Resources: CPU/memory requests and limits
    │
    └── NetworkPolicies (if isolation enabled)
        ├── Ingress Policy: default-deny with explicit allows
        └── Egress Policy: DNS + API + specified CIDRs/domains
```

### Data Flow

1. **Input**: KubernetesNamespace API resource with spec configuration
2. **Transformation**: Module computes resource quotas, labels, annotations, policies
3. **Deployment**: Creates namespace and associated Kubernetes resources
4. **Output**: Exports namespace name, quota status, network policy status, mesh status

## Design Decisions

### Why Abstract ResourceQuota into Profiles?

**Problem**: ResourceQuota configuration is complex with many fields (requests.cpu, limits.cpu, requests.memory, limits.memory, count/pods, count/services, etc.). Most users make the same mistakes (forgetting object counts, wrong CPU/memory ratios).

**Solution**: T-shirt sizes (SMALL/MEDIUM/LARGE/XLARGE) with opinionated defaults based on best practices. Advanced users can still specify custom quotas.

### Why Default-Deny Network Policies?

**Problem**: Without NetworkPolicies, all pods in a namespace can communicate with any other pod in the cluster (flat network), violating zero-trust principles.

**Solution**: Implement "Default Deny" pattern with explicit allow lists. Always allow DNS (kube-system) and intra-namespace communication. Users specify additional allowed namespaces/CIDRs.

### Why Service Mesh Abstraction?

**Problem**: Each service mesh (Istio/Linkerd/Consul) uses different labels/annotations for sidecar injection. Users need to remember mesh-specific syntax.

**Solution**: Unified API with `mesh_type` enum. Module handles mesh-specific label/annotation injection. Supports Istio revision tags for safe canary upgrades.

### Why Pod Security Standards?

**Problem**: PodSecurityPolicy (deprecated) was complex and cluster-wide. Kubernetes 1.23+ uses Pod Security Standards with namespace-level enforcement.

**Solution**: Simple enum (PRIVILEGED/BASELINE/RESTRICTED) applied as namespace labels. Kubernetes admission controller enforces automatically.

## Module Implementation

### Locals Computation (locals.go)

The `initializeLocals` function:
1. Extracts spec configuration from stack input
2. Builds combined labels (spec labels + standard labels + PSS labels)
3. Builds combined annotations (spec annotations + mesh annotations)
4. Computes resource quota values from preset or custom config
5. Computes limit range values if specified
6. Extracts network policy settings
7. Extracts service mesh settings

**Key Functions**:
- `computeResourceQuota`: Maps preset profiles to actual quota values
- `computeLimitRange`: Extracts custom default limits
- `buildAnnotations`: Adds mesh-specific annotations based on mesh type

### Resource Creation Flow (main.go)

1. **Initialize Locals**: Compute all derived values
2. **Create Provider**: Setup Kubernetes provider from credentials
3. **Create Namespace**: Base namespace resource with labels/annotations
4. **Create ResourceQuota**: If profile is configured
5. **Create LimitRange**: If default limits are specified
6. **Create NetworkPolicies**: If isolation is enabled
7. **Export Outputs**: Observable values for service discovery

### Resource Quota Logic (resource_quota.go)

**Preset Profiles**:
- SMALL: Dev/test environments (2-4 CPU, 4-8Gi RAM)
- MEDIUM: Staging (4-8 CPU, 8-16Gi RAM)
- LARGE: Production (8-16 CPU, 16-32Gi RAM)
- XLARGE: High-scale production (16-32 CPU, 32-64Gi RAM)

**Custom Quotas**: Direct mapping from spec to ResourceQuota hard limits

### Network Policy Logic (network_policies.go)

**Ingress Policy**:
- Default: Deny all ingress
- Allow: Traffic from allowed_ingress_namespaces
- Allow: Intra-namespace traffic

**Egress Policy**:
- Default: Deny all egress
- Allow: DNS to kube-system (UDP/TCP 53)
- Allow: Traffic to allowed_egress_cidrs
- Allow: Intra-namespace traffic

## Comparison to Alternative Approaches

### vs. Manual YAML

**Manual**: Write namespace.yaml, resourcequota.yaml, limitrange.yaml, networkpolicy.yaml separately. Risk of inconsistency.

**Module**: Single API resource. All resources created atomically. Consistent labeling. DRY principle.

### vs. Helm Chart

**Helm**: Templating with values.yaml. Difficult to encode complex logic (e.g., quota calculations). No type safety.

**Module**: Full programming language (Go). Type-safe protobuf API. Complex logic (e.g., mesh-specific annotations) trivial to implement.

### vs. Capsule Operator

**Capsule**: Cluster-level operator. Requires operator installation. Tenants submit regular namespace creation requests; operator intercepts and applies policies.

**Module**: Declarative IaC. No operator dependency. Works on any Kubernetes cluster. Policy is explicit in the manifest.

## Production Considerations

### Resource Quota Sizing

**Request-Based Billing**: Set `requests` to what teams reliably use. `limits` can be 2x for burst capacity.

**Object Count Safety**: Always set object count quotas. Prevents "1000 ConfigMaps" attacks that exhaust etcd.

### Network Policy Edge Cases

**DNS Resolution**: Egress policies must always allow port 53 to kube-system. Otherwise pods fail DNS lookups.

**Kubernetes API**: Some workloads need API access (e.g., controllers). Allow egress to Kubernetes API CIDR if needed.

### Service Mesh Upgrades

**Istio Revision Tags**: Use revision tags (e.g., "prod-stable") instead of hardcoded versions. When upgrading Istio, move the tag to the new version. Namespace config doesn't change; pods pick up new sidecar on rollout.

### Pod Security Standards

**Start Baseline**: RESTRICTED blocks many common patterns (host ports, privilege escalation). Start with BASELINE. Audit violations. Migrate to RESTRICTED incrementally.

## Testing Strategy

1. **Unit Tests**: spec_test.go validates protobuf rules
2. **Local Testing**: Deploy to kind/minikube with test manifest
3. **Integration Testing**: Verify ResourceQuota is enforced (try exceeding quota)
4. **Network Testing**: Verify NetworkPolicies block unauthorized traffic
5. **Production Validation**: Monitor namespace creation time, quota enforcement, policy violations

## Future Enhancements

1. **Hierarchical Namespaces**: Support HNC (Hierarchical Namespace Controller) for parent/child relationships
2. **DNS Policy**: Calico/Cilium DNS-based egress filtering for `allowed_egress_domains`
3. **Admission Webhooks**: Custom validation beyond protobuf rules (e.g., org-specific naming conventions)
4. **Cost Tracking**: Integration with Kubecost for real-time cost visibility per namespace
5. **Automatic Cleanup**: TTL-based garbage collection for ephemeral namespaces

## References

- [Kubernetes Namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/)
- [ResourceQuotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/)
- [LimitRanges](https://kubernetes.io/docs/concepts/policy/limit-range/)
- [NetworkPolicies](https://kubernetes.io/docs/concepts/services-networking/network-policies/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Multi-Tenancy Guide](https://kubernetes.io/docs/concepts/security/multi-tenancy/)
- [Namespace-as-a-Service Pattern](https://docs.rafay.co/template_catalog/get_started/namespace_asaservice/)


