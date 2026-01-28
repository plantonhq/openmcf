# GcpCloudRun API - Examples

This document provides real-world configuration examples for different Cloud Run deployment scenarios.

---

## Create using CLI

Create a YAML file using the examples below, then apply the configuration:

```shell
planton apply -f <yaml-path>
```

---

## Example 1: Basic Cloud Run Service (Direct Values)

A minimal Cloud Run service deployment using direct string values.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: my-api
spec:
  # GCP project (direct value)
  projectId:
    value: my-gcp-project-123
  
  # Region
  region: us-central1
  
  # Container configuration
  container:
    image:
      repo: us-docker.pkg.dev/my-project/my-registry/my-app
      tag: v1.0.0
    cpu: 1
    memory: 512
    port: 8080
    replicas:
      min: 0
      max: 10
  
  # Allow public access
  allowUnauthenticated: true
```

**What this creates:**
- A Cloud Run service in `us-central1`
- Single vCPU container with 512Mi memory
- Scale-to-zero enabled (min: 0)
- Publicly accessible

---

## Example 2: Cloud Run with Cross-Resource References

A production-ready Cloud Run service using references to other OpenMCF resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: prod-api
  org: acme-corp
  env:
    id: production
spec:
  # GCP project (reference to GcpProject resource)
  projectId:
    valueFrom:
      kind: GcpProject
      name: acme-prod-project
      fieldPath: status.outputs.project_id
  
  region: us-central1
  serviceName: production-api
  
  # Container configuration
  container:
    image:
      repo: us-docker.pkg.dev/acme-prod/containers/api
      tag: v2.1.0
    cpu: 2
    memory: 1024
    port: 8080
    replicas:
      min: 1
      max: 100
    env:
      variables:
        NODE_ENV: production
        LOG_LEVEL: info
      secrets:
        DATABASE_URL: projects/acme-prod/secrets/db-connection:latest
        API_KEY: projects/acme-prod/secrets/api-key:latest
  
  # Performance settings
  maxConcurrency: 100
  timeoutSeconds: 60
  
  # Security settings
  allowUnauthenticated: false
  ingress: INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER
  
  # Execution environment
  executionEnvironment: EXECUTION_ENVIRONMENT_GEN2
  
  # Deletion protection
  deleteProtection: true
```

**What this creates:**
- A production Cloud Run service with HA (min: 1)
- References GcpProject for project_id (resolved at deployment time)
- Private service with load balancer ingress
- Environment secrets from Secret Manager
- Deletion protection enabled

---

## Example 3: Cloud Run with VPC Access (Direct Values)

A Cloud Run service with Direct VPC Egress for accessing private resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: internal-api
spec:
  projectId:
    value: my-gcp-project-123
  
  region: us-central1
  
  container:
    image:
      repo: us-docker.pkg.dev/my-project/my-registry/internal-api
      tag: v1.0.0
    cpu: 1
    memory: 512
    replicas:
      min: 0
      max: 10
  
  # VPC access configuration (direct values)
  vpcAccess:
    network:
      value: projects/my-gcp-project-123/global/networks/my-vpc
    subnet:
      value: projects/my-gcp-project-123/regions/us-central1/subnetworks/my-subnet
    egress: PRIVATE_RANGES_ONLY
  
  # Internal only access
  allowUnauthenticated: false
  ingress: INGRESS_TRAFFIC_INTERNAL_ONLY
```

**What this creates:**
- Cloud Run service with Direct VPC Egress
- Only private IP traffic goes through VPC
- Internal-only ingress (no public access)
- Can access private resources in VPC (databases, Redis, etc.)

---

## Example 4: Cloud Run with VPC Access (Cross-Resource References)

A Cloud Run service using references to GcpVpc and GcpSubnetwork resources.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: backend-api
  org: acme-corp
  env:
    id: production
spec:
  # Reference to GcpProject
  projectId:
    valueFrom:
      kind: GcpProject
      name: acme-prod-project
      fieldPath: status.outputs.project_id
  
  region: us-central1
  
  container:
    image:
      repo: us-docker.pkg.dev/acme-prod/containers/backend
      tag: v3.0.0
    cpu: 2
    memory: 2048
    replicas:
      min: 2
      max: 50
  
  # VPC access with cross-resource references
  vpcAccess:
    network:
      valueFrom:
        kind: GcpVpc
        name: acme-prod-vpc
        fieldPath: status.outputs.network_name
    subnet:
      valueFrom:
        kind: GcpSubnetwork
        name: acme-prod-private-subnet
        fieldPath: status.outputs.subnetwork_name
    egress: ALL_TRAFFIC
  
  # Internal with load balancer for secure external access
  allowUnauthenticated: false
  ingress: INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER
```

**What this creates:**
- Backend service with high availability (min: 2)
- VPC network and subnet resolved from GcpVpc/GcpSubnetwork resources
- All egress traffic goes through VPC
- Accessible via internal load balancer

---

## Example 5: Cloud Run with Custom DNS

A Cloud Run service with custom domain mapping.

```yaml
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: public-api
spec:
  projectId:
    value: my-gcp-project-123
  
  region: us-central1
  
  container:
    image:
      repo: us-docker.pkg.dev/my-project/my-registry/api
      tag: v1.0.0
    cpu: 2
    memory: 1024
    replicas:
      min: 1
      max: 100
  
  # Public access
  allowUnauthenticated: true
  ingress: INGRESS_TRAFFIC_ALL
  
  # Custom DNS configuration
  dns:
    enabled: true
    hostnames:
      - api.example.com
      - api.acme.com
    managedZone: example-zone
```

**What this creates:**
- Public Cloud Run service
- Custom domain mapping to `api.example.com`
- TXT verification record created in Cloud DNS
- SSL certificate automatically provisioned

---

## Example 6: Complete Stack with Dependencies

A complete manifest showing a GcpCloudRun with all prerequisite resources.

```yaml
---
# First, create the GCP Project
apiVersion: gcp.openmcf.org/v1
kind: GcpProject
metadata:
  name: my-project
spec:
  projectName: my-cloud-run-project
  billingAccountId: "012345-6789AB-CDEF01"
  organizationId: "123456789012"

---
# Create the VPC network
apiVersion: gcp.openmcf.org/v1
kind: GcpVpc
metadata:
  name: my-vpc
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  autoCreateSubnetworks: false

---
# Create the subnet
apiVersion: gcp.openmcf.org/v1
kind: GcpSubnetwork
metadata:
  name: my-subnet
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  vpcSelfLink:
    valueFrom:
      kind: GcpVpc
      name: my-vpc
      fieldPath: status.outputs.self_link
  region: us-central1
  ipCidrRange: 10.0.0.0/24

---
# Finally, create the Cloud Run service
apiVersion: gcp.openmcf.org/v1
kind: GcpCloudRun
metadata:
  name: my-api
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: my-project
      fieldPath: status.outputs.project_id
  
  region: us-central1
  
  container:
    image:
      repo: us-docker.pkg.dev/my-project/my-registry/my-app
      tag: v1.0.0
    cpu: 1
    memory: 512
    replicas:
      min: 0
      max: 10
  
  vpcAccess:
    network:
      valueFrom:
        kind: GcpVpc
        name: my-vpc
        fieldPath: status.outputs.network_name
    subnet:
      valueFrom:
        kind: GcpSubnetwork
        name: my-subnet
        fieldPath: status.outputs.subnetwork_name
    egress: PRIVATE_RANGES_ONLY
  
  allowUnauthenticated: true
```

**Deploy order:**
1. GcpProject
2. GcpVpc
3. GcpSubnetwork
4. GcpCloudRun

OpenMCF automatically handles dependency ordering based on foreign key references.

---

## Field Reference

### projectId (StringValueOrRef)
- **value**: Direct GCP project ID string
- **valueFrom**: Reference to GcpProject resource
  - **kind**: `GcpProject`
  - **name**: Resource name
  - **fieldPath**: `status.outputs.project_id`

### vpcAccess.network (StringValueOrRef)
- **value**: Direct VPC network name/path
- **valueFrom**: Reference to GcpVpc resource
  - **kind**: `GcpVpc`
  - **name**: Resource name
  - **fieldPath**: `status.outputs.network_name`

### vpcAccess.subnet (StringValueOrRef)
- **value**: Direct subnetwork name/path
- **valueFrom**: Reference to GcpSubnetwork resource
  - **kind**: `GcpSubnetwork`
  - **name**: Resource name
  - **fieldPath**: `status.outputs.subnetwork_name`

---

## Next Steps

After creating a Cloud Run service, you may want to:

1. **Configure IAM**: Set up service account permissions for accessing other GCP services
2. **Set up monitoring**: Configure Cloud Monitoring alerts for the service
3. **Configure CI/CD**: Set up automated deployments from your container registry
4. **Implement traffic splitting**: Use revision-based deployments for canary releases

For more information, see the [research documentation](docs/README.md).
