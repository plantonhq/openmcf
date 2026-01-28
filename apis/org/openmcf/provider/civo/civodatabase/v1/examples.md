# Civo Database Examples

This document provides practical examples for deploying Civo managed databases using OpenMCF.

## Table of Contents

- [Minimal Development Database](#minimal-development-database)
- [Development with Custom Storage](#development-with-custom-storage)
- [Staging Environment with Basic HA](#staging-environment-with-basic-ha)
- [Production with Full HA and Security](#production-with-full-ha-and-security)
- [MySQL Database Example](#mysql-database-example)
- [Multi-Environment Setup](#multi-environment-setup)
- [Database with Foreign Key References](#database-with-foreign-key-references)

---

## Minimal Development Database

A minimal PostgreSQL database for local development or testing. Single node, small tier, no replicas.

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: dev-postgres
  env: development
spec:
  dbInstanceName: dev-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.small
  networkId:
    value: "net-12345678-abcd-1234-abcd-1234567890ab"
```

**Use Case**: Quick prototyping, feature development, CI/CD pipelines  
**Cost**: ~$43/month  
**HA**: No (single node)

---

## Development with Custom Storage

Development database with additional storage for larger datasets.

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: dev-postgres-large-storage
  env: development
spec:
  dbInstanceName: dev-large-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.small
  storageGib: 100  # Override default 40 GB with 100 GB
  networkId:
    value: "net-12345678-abcd-1234-abcd-1234567890ab"
  tags:
    - development
    - large-dataset
```

**Use Case**: Development with large test datasets or data migration testing  
**Cost**: ~$43/month + storage overage costs  
**HA**: No (single node)

---

## Staging Environment with Basic HA

Staging database with one replica for basic high availability testing.

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: staging-postgres
  env: staging
spec:
  dbInstanceName: staging-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.medium
  replicas: 1  # 2 total nodes (1 primary + 1 replica)
  networkId:
    value: "net-12345678-abcd-1234-abcd-1234567890ab"
  firewallIds:
    - value: "fw-staging-87654321-dcba-4321-dcba-0987654321ba"
  tags:
    - staging
    - backend
```

**Use Case**: Staging environment, pre-production testing  
**Cost**: ~$174/month (2x medium tier)  
**HA**: Basic (1 primary + 1 replica)

---

## Production with Full HA and Security

Production-grade PostgreSQL database with multiple replicas, custom storage, firewall, and complete tagging.

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: prod-postgres
  env: production
  org: acme-corp
spec:
  dbInstanceName: production-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.large
  replicas: 2  # 3 total nodes (1 primary + 2 replicas)
  networkId:
    value: "net-prod-12345678-abcd-1234-abcd-1234567890ab"
  firewallIds:
    - value: "fw-prod-87654321-dcba-4321-dcba-0987654321ba"
  storageGib: 250
  tags:
    - production
    - backend
    - primary-db
    - critical
```

**Use Case**: Production workloads with high availability requirements  
**Cost**: ~$521/month (3x large tier)  
**HA**: Full (1 primary + 2 replicas, automatic failover)

**Best Practices Applied**:
- ✅ High availability with 2 replicas
- ✅ Network isolation via private network
- ✅ Firewall rules for access control
- ✅ Custom storage for growth
- ✅ Comprehensive tagging for organization

---

## MySQL Database Example

MySQL 8.0 database for applications requiring MySQL compatibility.

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: prod-mysql
  env: production
spec:
  dbInstanceName: production-mysql
  engine: mysql
  engineVersion: "8.0"
  region: nyc1
  sizeSlug: g3.db.medium
  replicas: 1  # 2 total nodes
  networkId:
    value: "net-12345678-abcd-1234-abcd-1234567890ab"
  firewallIds:
    - value: "fw-87654321-dcba-4321-dcba-0987654321ba"
  tags:
    - production
    - mysql
    - backend
```

**Use Case**: Production MySQL workload (WordPress, legacy applications)  
**Cost**: ~$174/month (2x medium tier)  
**HA**: Basic (1 primary + 1 replica)

---

## Multi-Environment Setup

Deploying databases across multiple environments (dev, staging, prod) with consistent configuration.

### Development

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: app-db-dev
  env: development
spec:
  dbInstanceName: app-dev-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.small
  replicas: 0  # Single node
  networkId:
    value: "net-dev-12345678"
  tags:
    - development
    - app
```

### Staging

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: app-db-staging
  env: staging
spec:
  dbInstanceName: app-staging-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.medium
  replicas: 1  # 2 nodes for HA testing
  networkId:
    value: "net-staging-12345678"
  firewallIds:
    - value: "fw-staging-87654321"
  tags:
    - staging
    - app
```

### Production

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: app-db-prod
  env: production
spec:
  dbInstanceName: app-prod-db
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.large
  replicas: 2  # 3 nodes for full HA
  networkId:
    value: "net-prod-12345678"
  firewallIds:
    - value: "fw-prod-87654321"
  storageGib: 200
  tags:
    - production
    - app
    - critical
```

**Pattern**: Consistent naming (`app-db-{env}`), progressive scaling, increasing HA as environment criticality grows.

---

## Database with Foreign Key References

Using OpenMCF's foreign key references to wire dependencies between resources declaratively.

### Step 1: Create Network

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoVpc
metadata:
  name: app-network
spec:
  region: lon1
  label: app-private-network
```

### Step 2: Create Firewall

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoFirewall
metadata:
  name: db-firewall
spec:
  region: lon1
  networkId:
    value_from:
      kind: CivoVpc
      name: app-network
      field_path: status.outputs.network_id
  ingress_rules:
    - label: allow-postgres
      protocol: tcp
      port: "5432"
      cidr: "10.0.0.0/16"
      action: allow
```

### Step 3: Create Database with References

```yaml
apiVersion: code.openmcf.io/v1
kind: CivoDatabase
metadata:
  name: app-db
spec:
  dbInstanceName: app-database
  engine: postgres
  engineVersion: "16"
  region: lon1
  sizeSlug: g3.db.medium
  replicas: 1
  networkId:
    value_from:
      kind: CivoVpc
      name: app-network
      field_path: status.outputs.network_id
  firewallIds:
    - value_from:
        kind: CivoFirewall
        name: db-firewall
        field_path: status.outputs.firewall_id
  tags:
    - production
    - app
```

**Benefit**: OpenMCF automatically resolves dependencies and creates resources in the correct order. No manual copying of IDs required.

---

## Connection String Examples

After the database is provisioned, you can construct connection strings using the outputs.

### PostgreSQL

```bash
# Using DNS endpoint (recommended for HA)
postgresql://username:password@dns-endpoint:5432/database_name

# Example
postgresql://admin:SecurePass123@db-abc123.civo.com:5432/myapp
```

### MySQL

```bash
# Using DNS endpoint (recommended for HA)
mysql://username:password@dns-endpoint:3306/database_name

# Example
mysql://admin:SecurePass123@db-xyz789.civo.com:3306/myapp
```

### Environment Variables (Kubernetes)

```yaml
env:
  - name: DB_HOST
    value: "db-abc123.civo.com"  # Use dns_endpoint from outputs
  - name: DB_PORT
    value: "5432"
  - name: DB_USERNAME
    valueFrom:
      secretKeyRef:
        name: db-credentials
        key: username
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-credentials
        key: password
  - name: DB_NAME
    value: "myapp"
```

---

## Cost Comparison

| Configuration | Tier | Nodes | Monthly Cost |
|---------------|------|-------|--------------|
| Dev (minimal) | g3.db.small | 1 | ~$43 |
| Dev (large storage) | g3.db.small + 100GB | 1 | ~$43 + storage |
| Staging (basic HA) | g3.db.medium | 2 | ~$174 |
| Production (full HA) | g3.db.large | 3 | ~$521 |
| Production (max HA) | g3.db.xlarge | 5 | ~$1,740 |

**Note**: Civo pricing is all-inclusive (no egress or IOPS fees), making costs highly predictable.

---

## Deployment Workflow

1. **Create Network**: Deploy a `CivoVpc` for network isolation
2. **Create Firewall**: Deploy a `CivoFirewall` with least-privilege rules
3. **Create Database**: Deploy the `CivoDatabase` with references to network and firewall
4. **Verify Status**: Check resource status to ensure successful provisioning
5. **Retrieve Credentials**: Extract connection details from outputs
6. **Connect Application**: Update application configuration with connection details

```bash
# Deploy all resources
kubectl apply -f network.yaml
kubectl apply -f firewall.yaml
kubectl apply -f database.yaml

# Check status
kubectl get civodatabase app-db -o yaml

# Extract credentials (from status.outputs)
kubectl get civodatabase app-db -o jsonpath='{.status.outputs.dns_endpoint}'
```

---

## Best Practices Summary

1. **Development**: Use `g3.db.small` with no replicas to minimize costs
2. **Staging**: Use `g3.db.medium` with 1 replica to test HA behavior
3. **Production**: Use `g3.db.large` or higher with 2+ replicas for full HA
4. **Always use private networks**: Never deploy databases without a `network_id`
5. **Always use firewalls**: Attach a firewall with least-privilege rules
6. **Connect via DNS endpoint**: Use `dns_endpoint` for HA configurations
7. **Tag everything**: Use consistent tagging for cost tracking and organization
8. **Test failover**: Regularly test database failover in staging before relying on it in production

---

## Additional Resources

- **Configuration Reference**: See [`README.md`](README.md) for detailed field documentation
- **Research Documentation**: See [`docs/README.md`](docs/README.md) for architectural patterns
- **Pulumi Examples**: See [`iac/pulumi/examples.md`](iac/pulumi/examples.md)
- **Terraform Examples**: See [`iac/tf/examples.md`](iac/tf/examples.md)

