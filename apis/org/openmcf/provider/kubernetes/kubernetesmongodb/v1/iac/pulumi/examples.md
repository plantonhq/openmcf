# MongoDB Kubernetes API - Example Configurations

## Prerequisites

The **Percona Server for MongoDB Operator** must be installed on your cluster before deploying these examples. See the main [README.md](README.md) for installation instructions.

---

## Example w/ Basic Configuration

### Create Using CLI

Create a YAML file using the example shown below. After the YAML is created, use the following command to apply it:

```shell
planton apply -f <yaml-path>
```

### Basic Example

This basic example demonstrates a minimal configuration for deploying a MongoDB Kubernetes instance using the default settings, including 1 replica and persistence enabled.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: MongodbKubernetes
metadata:
  name: basic-mongodb
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: basic-mongodb
  container:
    replicas: 1
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Note**: The Percona operator automatically configures this as a replica set with 1 member.

---

## Example w/ Persistence Enabled

In this example, MongoDB persistence is enabled, and a persistent volume is created for each MongoDB pod to ensure data durability. The `disk_size` field defines the storage size allocated to the MongoDB pods.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: MongodbKubernetes
metadata:
  name: persistent-mongodb
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: persistent-mongodb
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 10Gi
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
```

**Note**: With 3 replicas, the Percona operator creates a 3-member replica set for high availability.

---

## Example w/ High Availability Configuration

This example demonstrates deploying a highly available MongoDB cluster with multiple replicas. The Percona operator automatically configures replica set relationships and handles failover.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: MongodbKubernetes
metadata:
  name: ha-mongodb
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: ha-mongodb
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 20Gi
    resources:
      requests:
        cpu: 200m
        memory: 512Mi
      limits:
        cpu: 2000m
        memory: 2Gi
```

**Key Features**:
- 3 replicas provide automatic failover
- Primary election happens automatically
- Survives 1 node failure without data loss

---

## Example w/ Ingress Enabled

In this example, ingress is enabled to allow external access to the MongoDB service. This is particularly useful when MongoDB needs to be accessed by clients outside the Kubernetes cluster.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: MongodbKubernetes
metadata:
  name: ingress-mongodb
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: ingress-mongodb
  container:
    replicas: 1
    persistenceEnabled: true
    diskSize: 5Gi
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 1000m
        memory: 1Gi
  ingress:
    enabled: true
    hostname: mongodb.example.com
```

**Security Note**: When enabling ingress, ensure proper network policies and authentication are configured.

---

## Example w/ Production Configuration

Optimized configuration for production workloads with sufficient resources and high availability.

```yaml
apiVersion: kubernetes.openmcf.org/v1
kind: MongodbKubernetes
metadata:
  name: production-mongodb
  org: my-org
  env: production
spec:
  targetCluster:
    clusterName: my-gke-cluster
  namespace:
    value: production-mongodb
  container:
    replicas: 3
    persistenceEnabled: true
    diskSize: 100Gi
    resources:
      requests:
        cpu: 1000m
        memory: 4Gi
      limits:
        cpu: 4000m
        memory: 16Gi
  ingress:
    isEnabled: false  # Use internal service for security
```

**Production Features**:
- 3-member replica set for HA
- Large persistent volumes (100Gi per pod)
- Substantial resources for production load
- Internal-only access for security

---

## Deployment Verification

After deploying MongoDB, verify the deployment using these commands:

### Check Percona Operator

```bash
# Verify operator is running
kubectl get pods -n mongodb-operator
kubectl logs -n mongodb-operator -l app.kubernetes.io/name=percona-server-mongodb-operator
```

### Check MongoDB Resources

```bash
# Check PerconaServerMongoDB CRD
kubectl get perconaservermongodbs -n <namespace>

# Check MongoDB pods
kubectl get pods -n <namespace>

# Check StatefulSet
kubectl get statefulset -n <namespace>

# Check services
kubectl get svc -n <namespace>
```

### Check Replica Set Status

```bash
# Port-forward to MongoDB
kubectl port-forward -n <namespace> svc/<mongodb-name> 27017:27017

# Connect with mongo shell
mongo mongodb://localhost:27017

# Check replica set status
rs.status()
```

### Expected Output

You should see:
- **PerconaServerMongoDB**: Custom resource created
- **StatefulSet**: One StatefulSet with desired replicas
- **Pods**: All pods in Running state (e.g., `mongodb-rs0-0`, `mongodb-rs0-1`, `mongodb-rs0-2`)
- **Services**: Cluster service and headless service for replica set
- **PersistentVolumeClaims**: One PVC per pod (if persistence enabled)

---

## Accessing MongoDB

### Internal Access (from within cluster)

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mongodb-client
spec:
  containers:
  - name: mongodb-client
    image: mongo:8.0
    command: ['sleep', '3600']
```

Then exec into the pod:
```bash
kubectl exec -it mongodb-client -- mongosh \
  "mongodb://<username>:<password>@<mongodb-service>:27017/?replicaSet=rs0"
```

### External Access (with ingress enabled)

```bash
# Get credentials
kubectl get secret <mongodb-name> -n <namespace> -o jsonpath='{.data.MONGODB_DATABASE_ADMIN_PASSWORD}' | base64 -d

# Connect using external hostname
mongosh "mongodb://<username>:<password>@<external-hostname>:27017/?replicaSet=rs0"
```

---

## Common Issues and Solutions

### Issue: Operator Not Found

**Error**: `CustomResourceDefinition "perconaservermongodbs.psmdb.percona.com" not found`

**Solution**: Install the Percona Server for MongoDB Operator first:
```bash
planton pulumi up --manifest percona-operator.yaml \
  --module-dir apis/project/planton/provider/kubernetes/perconaservermongodboperator/v1/iac/pulumi
```

### Issue: Pods Not Starting

**Solution**: Check operator logs and pod events:
```bash
kubectl logs -n mongodb-operator -l app.kubernetes.io/name=percona-server-mongodb-operator
kubectl describe pod <pod-name> -n <namespace>
```

### Issue: Replica Set Not Initializing

**Solution**: Check that all replica pods are running and can communicate:
```bash
kubectl get pods -n <namespace>
kubectl logs <pod-name> -n <namespace>
```

---

## Next Steps

After successfully deploying MongoDB:

1. **Configure Backups**: Use Percona's backup CRD to configure automated backups
2. **Set Up Monitoring**: Deploy monitoring agents for Percona Monitoring and Management (PMM)
3. **Configure Access Control**: Set up additional users and roles as needed
4. **Test Failover**: Verify high availability by simulating node failures
5. **Optimize Performance**: Adjust resources based on workload requirements

For more information, see the [Percona Operator Documentation](https://docs.percona.com/percona-operator-for-mongodb/).
